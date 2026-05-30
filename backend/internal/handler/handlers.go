package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/request"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/service"
	"github.com/thanthtooaung-coding/pdf-reader/backend/pkg/utils"
)

type AuthHandler struct{ svc service.AuthService }

func NewAuthHandler(svc service.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

func (h *AuthHandler) RegisterRequest(c *fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, err)
	}
	if err := utils.ValidateStruct(req); err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.RegisterRequest(c.Context(), req)
	if err != nil {
		return mapAuthError(c, err)
	}
	return utils.OK(c, res, "otp_sent")
}

func (h *AuthHandler) RegisterVerify(c *fiber.Ctx) error {
	var req request.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, err)
	}
	if err := utils.ValidateStruct(req); err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.RegisterVerify(c.Context(), req)
	if err != nil {
		return mapAuthError(c, err)
	}
	return utils.Created(c, res)
}

func (h *AuthHandler) RegisterResend(c *fiber.Ctx) error {
	var req request.ResendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, err)
	}
	if err := utils.ValidateStruct(req); err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.RegisterResend(c.Context(), req)
	if err != nil {
		return mapAuthError(c, err)
	}
	return utils.OK(c, res, "otp_sent")
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, err)
	}
	if err := utils.ValidateStruct(req); err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.LoginRequest(c.Context(), req)
	if err != nil {
		return mapAuthError(c, err)
	}
	return utils.OK(c, res, "otp_sent")
}

func (h *AuthHandler) LoginVerify(c *fiber.Ctx) error {
	var req request.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, err)
	}
	if err := utils.ValidateStruct(req); err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.LoginVerify(c.Context(), req)
	if err != nil {
		return mapAuthError(c, err)
	}
	return utils.OK(c, res)
}

func mapAuthError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		return utils.Unauthorized(c, err)
	case errors.Is(err, service.ErrTooManyAttempts), errors.Is(err, service.ErrTooManyOTPRequests):
		return utils.TooManyRequests(c, err)
	case errors.Is(err, service.ErrEmailUnavailable):
		return utils.Internal(c, err)
	default:
		return utils.BadRequest(c, err)
	}
}

type UserHandler struct{ svc service.UserService }

func NewUserHandler(svc service.UserService) *UserHandler { return &UserHandler{svc: svc} }

func (h *UserHandler) Me(c *fiber.Ctx) error {
	userID, err := utils.MustActorID(c)
	if err != nil {
		return utils.Unauthorized(c, err)
	}
	res, err := h.svc.GetMe(userID)
	if err != nil {
		return utils.NotFound(c, err)
	}
	return utils.OK(c, res)
}

type WorkspaceHandler struct{ svc service.WorkspaceService }

func NewWorkspaceHandler(svc service.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{svc: svc}
}

func (h *WorkspaceHandler) Create(c *fiber.Ctx) error {
	userID, err := utils.MustActorID(c)
	if err != nil {
		return utils.Unauthorized(c, err)
	}
	var req request.CreateWorkspaceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, err)
	}
	if err := utils.ValidateStruct(req); err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.Create(userID, req)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return utils.Created(c, res)
}

func (h *WorkspaceHandler) Get(c *fiber.Ctx) error {
	userID, err := utils.MustActorID(c)
	if err != nil {
		return utils.Unauthorized(c, err)
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.Get(userID, id)
	if err != nil {
		return mapDomainError(c, err)
	}
	return utils.OK(c, res)
}

func (h *WorkspaceHandler) List(c *fiber.Ctx) error {
	userID, err := utils.MustActorID(c)
	if err != nil {
		return utils.Unauthorized(c, err)
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	res, err := h.svc.List(userID, page, limit)
	if err != nil {
		return utils.Internal(c, err)
	}
	return utils.OK(c, res)
}

type FileHandler struct {
	svc     service.FileService
	maxSize int
}

func NewFileHandler(svc service.FileService, maxUploadMB int) *FileHandler {
	return &FileHandler{svc: svc, maxSize: maxUploadMB << 20}
}

func (h *FileHandler) Upload(c *fiber.Ctx) error {
	userID, err := utils.MustActorID(c)
	if err != nil {
		return utils.Unauthorized(c, err)
	}
	workspaceID, err := uuid.Parse(c.Params("workspaceId"))
	if err != nil {
		return utils.BadRequest(c, err)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return utils.BadRequest(c, errors.New("file is required"))
	}
	if fileHeader.Size > int64(h.maxSize) {
		return utils.BadRequest(c, errors.New("file too large"))
	}

	f, err := fileHeader.Open()
	if err != nil {
		return utils.Internal(c, err)
	}
	defer f.Close()

	res, err := h.svc.Upload(userID, workspaceID, fileHeader.Filename, f)
	if err != nil {
		return mapDomainError(c, err)
	}
	return utils.Created(c, res)
}

func (h *FileHandler) Get(c *fiber.Ctx) error {
	userID, err := utils.MustActorID(c)
	if err != nil {
		return utils.Unauthorized(c, err)
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.Get(userID, id)
	if err != nil {
		return mapDomainError(c, err)
	}
	return utils.OK(c, res)
}

func (h *FileHandler) ListByWorkspace(c *fiber.Ctx) error {
	userID, err := utils.MustActorID(c)
	if err != nil {
		return utils.Unauthorized(c, err)
	}
	workspaceID, err := uuid.Parse(c.Params("workspaceId"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.ListByWorkspace(userID, workspaceID)
	if err != nil {
		return mapDomainError(c, err)
	}
	return utils.OK(c, res)
}

func (h *FileHandler) Download(c *fiber.Ctx) error {
	userID, err := utils.MustActorID(c)
	if err != nil {
		return utils.Unauthorized(c, err)
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	path, name, err := h.svc.DownloadPath(userID, id)
	if err != nil {
		return mapDomainError(c, err)
	}
	return c.Download(path, name)
}

type CommentHandler struct{ svc service.CommentService }

func NewCommentHandler(svc service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

func (h *CommentHandler) Create(c *fiber.Ctx) error {
	userID, err := utils.MustActorID(c)
	if err != nil {
		return utils.Unauthorized(c, err)
	}
	fileID, err := uuid.Parse(c.Params("fileId"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	var req request.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, err)
	}
	if err := utils.ValidateStruct(req); err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.Create(userID, fileID, req)
	if err != nil {
		return mapDomainError(c, err)
	}
	return utils.Created(c, res)
}

func (h *CommentHandler) List(c *fiber.Ctx) error {
	userID, err := utils.MustActorID(c)
	if err != nil {
		return utils.Unauthorized(c, err)
	}
	fileID, err := uuid.Parse(c.Params("fileId"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.List(userID, fileID)
	if err != nil {
		return mapDomainError(c, err)
	}
	return utils.OK(c, res)
}

type AIJobHandler struct{ svc service.AIJobService }

func NewAIJobHandler(svc service.AIJobService) *AIJobHandler { return &AIJobHandler{svc: svc} }

func (h *AIJobHandler) Create(c *fiber.Ctx) error {
	userID, err := utils.MustActorID(c)
	if err != nil {
		return utils.Unauthorized(c, err)
	}
	var req request.CreateAIJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, err)
	}
	if err := utils.ValidateStruct(req); err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.Create(userID, req)
	if err != nil {
		return mapDomainError(c, err)
	}
	return utils.Created(c, res)
}

func (h *AIJobHandler) Get(c *fiber.Ctx) error {
	userID, err := utils.MustActorID(c)
	if err != nil {
		return utils.Unauthorized(c, err)
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.Get(userID, id)
	if err != nil {
		return mapDomainError(c, err)
	}
	return utils.OK(c, res)
}

func (h *AIJobHandler) ListByWorkspace(c *fiber.Ctx) error {
	userID, err := utils.MustActorID(c)
	if err != nil {
		return utils.Unauthorized(c, err)
	}
	workspaceID, err := uuid.Parse(c.Params("workspaceId"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	res, err := h.svc.ListByWorkspace(userID, workspaceID)
	if err != nil {
		return mapDomainError(c, err)
	}
	return utils.OK(c, res)
}

func mapDomainError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, service.ErrForbidden):
		return utils.Forbidden(c, err)
	case errors.Is(err, service.ErrWorkspaceNotFound),
		errors.Is(err, service.ErrFileNotFound),
		errors.Is(err, service.ErrAIJobNotFound),
		errors.Is(err, service.ErrUserNotFound):
		return utils.NotFound(c, err)
	case errors.Is(err, service.ErrInvalidFileType):
		return utils.BadRequest(c, err)
	default:
		return utils.Internal(c, err)
	}
}
