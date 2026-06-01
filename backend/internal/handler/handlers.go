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

// RegisterRequest starts user registration and sends an OTP email.
//
// @Summary      Request registration OTP
// @Description  Validates signup fields and sends a one-time password to the email address.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      request.RegisterRequest  true  "Registration payload"
// @Success      200   {object}  response.OTPEnvelope
// @Failure      400   {object}  response.ErrorEnvelope
// @Failure      429   {object}  response.ErrorEnvelope
// @Failure      500   {object}  response.ErrorEnvelope
// @Router       /api/v1/auth/register/request [post]
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

// RegisterVerify completes registration using the OTP sent to email.
//
// @Summary      Verify registration OTP
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      request.VerifyOTPRequest  true  "OTP verification payload"
// @Success      201   {object}  response.AuthEnvelope
// @Failure      400   {object}  response.ErrorEnvelope
// @Failure      429   {object}  response.ErrorEnvelope
// @Router       /api/v1/auth/register/verify [post]
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

// RegisterResend resends the registration OTP.
//
// @Summary      Resend registration OTP
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      request.ResendOTPRequest  true  "Email address"
// @Success      200   {object}  response.OTPEnvelope
// @Failure      400   {object}  response.ErrorEnvelope
// @Failure      429   {object}  response.ErrorEnvelope
// @Router       /api/v1/auth/register/resend [post]
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

// Login validates credentials and sends a login OTP.
//
// @Summary      Request login OTP
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      request.LoginRequest  true  "Login credentials"
// @Success      200   {object}  response.OTPEnvelope
// @Failure      400   {object}  response.ErrorEnvelope
// @Failure      401   {object}  response.ErrorEnvelope
// @Failure      429   {object}  response.ErrorEnvelope
// @Router       /api/v1/auth/login [post]
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

// LoginVerify completes login using the OTP sent to email.
//
// @Summary      Verify login OTP
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      request.VerifyOTPRequest  true  "OTP verification payload"
// @Success      200   {object}  response.AuthEnvelope
// @Failure      400   {object}  response.ErrorEnvelope
// @Failure      401   {object}  response.ErrorEnvelope
// @Failure      429   {object}  response.ErrorEnvelope
// @Router       /api/v1/auth/login/verify [post]
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

// Me returns the authenticated user profile.
//
// @Summary      Get current user
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.UserEnvelope
// @Failure      401  {object}  response.ErrorEnvelope
// @Failure      404  {object}  response.ErrorEnvelope
// @Router       /api/v1/users/me [get]
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

// Create creates a new workspace for the authenticated user.
//
// @Summary      Create workspace
// @Tags         workspaces
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      request.CreateWorkspaceRequest  true  "Workspace name"
// @Success      201   {object}  response.WorkspaceEnvelope
// @Failure      400   {object}  response.ErrorEnvelope
// @Failure      401   {object}  response.ErrorEnvelope
// @Router       /api/v1/workspaces [post]
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

// Get returns a workspace owned by the authenticated user.
//
// @Summary      Get workspace
// @Tags         workspaces
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Workspace ID (UUID)"
// @Success      200  {object}  response.WorkspaceEnvelope
// @Failure      401  {object}  response.ErrorEnvelope
// @Failure      403  {object}  response.ErrorEnvelope
// @Failure      404  {object}  response.ErrorEnvelope
// @Router       /api/v1/workspaces/{id} [get]
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

// List returns paginated workspace history for the authenticated user.
//
// @Summary      List workspaces
// @Tags         workspaces
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int  false  "Page number"   default(1)
// @Param        limit  query     int  false  "Items per page"  default(20)
// @Success      200    {object}  response.WorkspaceListEnvelope
// @Failure      401    {object}  response.ErrorEnvelope
// @Router       /api/v1/workspaces [get]
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

// Upload stores a PDF file in the given workspace.
//
// @Summary      Upload PDF
// @Tags         files
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        workspaceId  path      string  true  "Workspace ID (UUID)"
// @Param        file         formData  file    true  "PDF file"
// @Success      201          {object}  response.FileEnvelope
// @Failure      400          {object}  response.ErrorEnvelope
// @Failure      401          {object}  response.ErrorEnvelope
// @Failure      403          {object}  response.ErrorEnvelope
// @Router       /api/v1/workspaces/{workspaceId}/files [post]
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

// Get returns file metadata.
//
// @Summary      Get file metadata
// @Tags         files
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "File ID (UUID)"
// @Success      200  {object}  response.FileEnvelope
// @Failure      401  {object}  response.ErrorEnvelope
// @Failure      403  {object}  response.ErrorEnvelope
// @Failure      404  {object}  response.ErrorEnvelope
// @Router       /api/v1/files/{id} [get]
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

// ListByWorkspace lists files in a workspace.
//
// @Summary      List workspace files
// @Tags         files
// @Produce      json
// @Security     BearerAuth
// @Param        workspaceId  path      string  true  "Workspace ID (UUID)"
// @Success      200          {object}  response.FileListEnvelope
// @Failure      401          {object}  response.ErrorEnvelope
// @Failure      403          {object}  response.ErrorEnvelope
// @Failure      404          {object}  response.ErrorEnvelope
// @Router       /api/v1/workspaces/{workspaceId}/files [get]
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

// Download streams the stored PDF file.
//
// @Summary      Download PDF
// @Tags         files
// @Produce      application/pdf
// @Security     BearerAuth
// @Param        id   path      string  true  "File ID (UUID)"
// @Success      200  {file}    file
// @Failure      401  {object}  response.ErrorEnvelope
// @Failure      403  {object}  response.ErrorEnvelope
// @Failure      404  {object}  response.ErrorEnvelope
// @Router       /api/v1/files/{id}/download [get]
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

// Create adds a comment to a file.
//
// @Summary      Create comment
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        fileId  path      string                     true  "File ID (UUID)"
// @Param        body    body      request.CreateCommentRequest  true  "Comment message"
// @Success      201     {object}  response.CommentEnvelope
// @Failure      400     {object}  response.ErrorEnvelope
// @Failure      401     {object}  response.ErrorEnvelope
// @Failure      403     {object}  response.ErrorEnvelope
// @Failure      404     {object}  response.ErrorEnvelope
// @Router       /api/v1/files/{fileId}/comments [post]
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

// List returns comments for a file.
//
// @Summary      List comments
// @Tags         comments
// @Produce      json
// @Security     BearerAuth
// @Param        fileId  path      string  true  "File ID (UUID)"
// @Success      200     {object}  response.CommentListEnvelope
// @Failure      401     {object}  response.ErrorEnvelope
// @Failure      403     {object}  response.ErrorEnvelope
// @Failure      404     {object}  response.ErrorEnvelope
// @Router       /api/v1/files/{fileId}/comments [get]
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

// Create enqueues an AI job (TRANSLATE, SUMMARIZE, or COMMENT).
//
// @Summary      Create AI job
// @Description  TRANSLATE and SUMMARIZE use OpenAI when OPENAI_API_KEY is configured.
// @Tags         ai-jobs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      request.CreateAIJobRequest  true  "AI job payload"
// @Success      201   {object}  response.AIJobEnvelope
// @Failure      400   {object}  response.ErrorEnvelope
// @Failure      401   {object}  response.ErrorEnvelope
// @Failure      403   {object}  response.ErrorEnvelope
// @Router       /api/v1/ai-jobs [post]
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

// Get returns AI job status and output.
//
// @Summary      Get AI job
// @Tags         ai-jobs
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "AI job ID (UUID)"
// @Success      200  {object}  response.AIJobEnvelope
// @Failure      401  {object}  response.ErrorEnvelope
// @Failure      403  {object}  response.ErrorEnvelope
// @Failure      404  {object}  response.ErrorEnvelope
// @Router       /api/v1/ai-jobs/{id} [get]
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

// ListByWorkspace lists AI jobs for a workspace.
//
// @Summary      List workspace AI jobs
// @Tags         ai-jobs
// @Produce      json
// @Security     BearerAuth
// @Param        workspaceId  path      string  true  "Workspace ID (UUID)"
// @Success      200          {object}  response.AIJobListEnvelope
// @Failure      401          {object}  response.ErrorEnvelope
// @Failure      403          {object}  response.ErrorEnvelope
// @Failure      404          {object}  response.ErrorEnvelope
// @Router       /api/v1/workspaces/{workspaceId}/ai-jobs [get]
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
