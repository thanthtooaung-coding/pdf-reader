package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/config"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/mapper"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/otp"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/repository"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/request"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/response"
	"github.com/thanthtooaung-coding/pdf-reader/backend/pkg/email"
	"github.com/thanthtooaung-coding/pdf-reader/backend/pkg/utils"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid email, username or password")
	ErrOTPNotFound        = errors.New("otp expired or not found")
	ErrInvalidOTP         = errors.New("invalid otp")
	ErrTooManyAttempts    = errors.New("too many attempts")
	ErrTooManyOTPRequests = errors.New("too many otp requests")
	ErrPendingNotFound    = errors.New("no pending registration found")
	ErrRoleNotFound       = errors.New("default role not configured")
	ErrEmailUnavailable   = errors.New("email service unavailable")
)

type AuthService interface {
	RegisterRequest(ctx context.Context, req request.RegisterRequest) (*response.OTPSentResponse, error)
	RegisterVerify(ctx context.Context, req request.VerifyOTPRequest) (*response.AuthResponse, error)
	RegisterResend(ctx context.Context, req request.ResendOTPRequest) (*response.OTPSentResponse, error)
	LoginRequest(ctx context.Context, req request.LoginRequest) (*response.OTPSentResponse, error)
	LoginVerify(ctx context.Context, req request.VerifyOTPRequest) (*response.AuthResponse, error)
}

type authServiceImpl struct {
	logger   *logrus.Logger
	cfg      *config.Config
	db       *gorm.DB
	users    repository.UserRepository
	roles    repository.RoleRepository
	otpStore *otp.Store
	mailer   email.Sender
}

func NewAuthService(
	logger *logrus.Logger,
	cfg *config.Config,
	db *gorm.DB,
	users repository.UserRepository,
	roles repository.RoleRepository,
	otpStore *otp.Store,
	mailer email.Sender,
) AuthService {
	return &authServiceImpl{
		logger: logger, cfg: cfg, db: db,
		users: users, roles: roles, otpStore: otpStore, mailer: mailer,
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (s *authServiceImpl) RegisterRequest(ctx context.Context, req request.RegisterRequest) (*response.OTPSentResponse, error) {
	emailAddr := normalizeEmail(req.Email)
	username := strings.TrimSpace(req.Username)

	if existing, _ := s.users.GetByEmail(emailAddr); existing != nil {
		return nil, ErrEmailExists
	}
	if existing, _ := s.users.GetByUsername(username); existing != nil {
		return nil, ErrUsernameExists
	}

	rateWindow := time.Duration(s.cfg.OTPRateLimitWindowMin) * time.Minute
	if err := s.otpStore.CheckRegRateLimit(ctx, emailAddr, s.cfg.OTPRateLimitMax, rateWindow); err != nil {
		return nil, ErrTooManyOTPRequests
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	code, err := utils.GenerateOTP()
	if err != nil {
		return nil, errors.New("failed to generate otp")
	}

	pending := otp.PendingRegistration{
		OTPHash:      utils.HashOTP(code, s.otpStore.Pepper()),
		Fullname:     req.Fullname,
		Username:     username,
		Email:        emailAddr,
		PasswordHash: hashedPassword,
	}

	ttl := time.Duration(s.cfg.OTPTTLSeconds) * time.Second
	if err := s.otpStore.SavePending(ctx, emailAddr, pending, ttl); err != nil {
		return nil, errors.New("failed to start registration")
	}

	if err := s.sendOTP(emailAddr, code); err != nil {
		_ = s.otpStore.DeletePending(ctx, emailAddr)
		return nil, err
	}

	return s.otpResponse(code), nil
}

func (s *authServiceImpl) RegisterResend(ctx context.Context, req request.ResendOTPRequest) (*response.OTPSentResponse, error) {
	emailAddr := normalizeEmail(req.Email)
	pending, err := s.otpStore.GetPending(ctx, emailAddr)
	if err != nil || pending == nil {
		return nil, ErrPendingNotFound
	}

	rateWindow := time.Duration(s.cfg.OTPRateLimitWindowMin) * time.Minute
	if err := s.otpStore.CheckRegRateLimit(ctx, emailAddr, s.cfg.OTPRateLimitMax, rateWindow); err != nil {
		return nil, ErrTooManyOTPRequests
	}

	code, err := utils.GenerateOTP()
	if err != nil {
		return nil, errors.New("failed to generate otp")
	}

	pending.OTPHash = utils.HashOTP(code, s.otpStore.Pepper())
	pending.Attempts = 0
	ttl := time.Duration(s.cfg.OTPTTLSeconds) * time.Second
	if err := s.otpStore.SavePending(ctx, emailAddr, *pending, ttl); err != nil {
		return nil, errors.New("failed to update pending registration")
	}

	if err := s.sendOTP(emailAddr, code); err != nil {
		return nil, err
	}
	return s.otpResponse(code), nil
}

func (s *authServiceImpl) RegisterVerify(ctx context.Context, req request.VerifyOTPRequest) (*response.AuthResponse, error) {
	emailAddr := normalizeEmail(req.Email)
	pending, err := s.otpStore.GetPending(ctx, emailAddr)
	if err != nil || pending == nil {
		return nil, ErrOTPNotFound
	}

	pending.Attempts++
	if pending.Attempts > s.cfg.OTPMaxVerifyAttempts {
		_ = s.otpStore.DeletePending(ctx, emailAddr)
		return nil, ErrTooManyAttempts
	}

	if !utils.VerifyOTP(req.OTP, s.otpStore.Pepper(), pending.OTPHash) {
		ttl := time.Duration(s.cfg.OTPTTLSeconds) * time.Second
		_ = s.otpStore.SavePending(ctx, emailAddr, *pending, ttl)
		return nil, ErrInvalidOTP
	}

	role, err := s.roles.GetByCode("USER")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	var created models.User
	err = s.db.Transaction(func(tx *gorm.DB) error {
		u := &models.User{
			RoleID:   role.ID,
			Fullname: pending.Fullname,
			Username: pending.Username,
			Email:    emailAddr,
			Password: pending.PasswordHash,
			Base:     models.Base{IsEnable: true},
		}
		if err := tx.Create(u).Error; err != nil {
			return err
		}
		created = *u
		return nil
	})
	if err != nil {
		return nil, errors.New("failed to create account")
	}

	_ = s.otpStore.DeletePending(ctx, emailAddr)

	full, err := s.users.GetByID(created.ID)
	if err != nil {
		return nil, err
	}

	return s.authResponse(*full)
}

func (s *authServiceImpl) LoginRequest(ctx context.Context, req request.LoginRequest) (*response.OTPSentResponse, error) {
	emailAddr := normalizeEmail(req.Email)
	username := strings.TrimSpace(req.Username)

	u, err := s.users.GetByEmail(emailAddr)
	if err != nil || u == nil || !u.IsEnable {
		return nil, ErrInvalidCredentials
	}
	if username != "" && !strings.EqualFold(u.Username, username) {
		return nil, ErrInvalidCredentials
	}
	if err := utils.CheckPassword(req.Password, u.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	rateWindow := time.Duration(s.cfg.OTPRateLimitWindowMin) * time.Minute
	if err := s.otpStore.CheckLoginRateLimit(ctx, emailAddr, s.cfg.OTPRateLimitMax, rateWindow); err != nil {
		return nil, ErrTooManyOTPRequests
	}

	code, err := utils.GenerateOTP()
	if err != nil {
		return nil, errors.New("failed to generate otp")
	}

	pending := otp.PendingLogin{
		OTPHash:  utils.HashOTP(code, s.otpStore.Pepper()),
		Email:    emailAddr,
		UserID:   u.ID.String(),
		Attempts: 0,
	}
	ttl := time.Duration(s.cfg.OTPTTLSeconds) * time.Second
	if err := s.otpStore.SaveLogin(ctx, emailAddr, pending, ttl); err != nil {
		return nil, errors.New("failed to start login")
	}

	if err := s.sendOTP(emailAddr, code); err != nil {
		_ = s.otpStore.DeleteLogin(ctx, emailAddr)
		return nil, err
	}
	return s.otpResponse(code), nil
}

func (s *authServiceImpl) LoginVerify(ctx context.Context, req request.VerifyOTPRequest) (*response.AuthResponse, error) {
	emailAddr := normalizeEmail(req.Email)
	pending, err := s.otpStore.GetLogin(ctx, emailAddr)
	if err != nil || pending == nil {
		return nil, ErrOTPNotFound
	}

	pending.Attempts++
	if pending.Attempts > s.cfg.OTPMaxVerifyAttempts {
		_ = s.otpStore.DeleteLogin(ctx, emailAddr)
		return nil, ErrTooManyAttempts
	}

	if !utils.VerifyOTP(req.OTP, s.otpStore.Pepper(), pending.OTPHash) {
		ttl := time.Duration(s.cfg.OTPTTLSeconds) * time.Second
		_ = s.otpStore.SaveLogin(ctx, emailAddr, *pending, ttl)
		return nil, ErrInvalidOTP
	}

	u, err := s.users.GetByEmail(emailAddr)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	_ = s.otpStore.DeleteLogin(ctx, emailAddr)
	return s.authResponse(*u)
}

func (s *authServiceImpl) authResponse(u models.User) (*response.AuthResponse, error) {
	roleCode := ""
	if u.Role != nil {
		roleCode = u.Role.Code
	}
	token, exp, err := utils.IssueAccessToken(s.cfg.JWTSecret, s.cfg.JWTExpiresHours, u.ID, u.Email, roleCode)
	if err != nil {
		return nil, err
	}
	return &response.AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(time.Until(exp).Seconds()),
		User:        mapper.ToUserResponse(u),
	}, nil
}

func (s *authServiceImpl) otpResponse(code string) *response.OTPSentResponse {
	resp := &response.OTPSentResponse{
		Message:          "otp_sent",
		ExpiresInSeconds: s.cfg.OTPTTLSeconds,
	}
	if s.cfg.OTPDebugReturn {
		resp.OTP = code
	}
	return resp
}

func (s *authServiceImpl) sendOTP(to, code string) error {
	if !s.mailer.Enabled() {
		if s.cfg.OTPDebugReturn {
			s.logger.WithField("email", to).WithField("otp", code).Warn("email not configured; otp logged")
			return nil
		}
		return ErrEmailUnavailable
	}
	if err := s.mailer.SendOTP(to, code); err != nil {
		s.logger.WithError(err).Error("send otp email")
		return ErrEmailUnavailable
	}
	return nil
}
