package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/config"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/handler"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/middleware"
)

type Handlers struct {
	Health    *handler.HealthHandler
	Swagger   fiber.Handler
	Auth      *handler.AuthHandler
	User       *handler.UserHandler
	Workspace  *handler.WorkspaceHandler
	File       *handler.FileHandler
	Comment    *handler.CommentHandler
	AIJob      *handler.AIJobHandler
}

func Register(app *fiber.App, cfg *config.Config, h Handlers) {
	app.Get("/healthz", h.Health.HealthCheck)

	app.Get("/swagger/*", h.Swagger)

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register/request", h.Auth.RegisterRequest)
	auth.Post("/register/verify", h.Auth.RegisterVerify)
	auth.Post("/register/resend", h.Auth.RegisterResend)
	auth.Post("/login", h.Auth.Login)
	auth.Post("/login/verify", h.Auth.LoginVerify)

	protected := api.Group("", middleware.JWTAuth(cfg))
	protected.Get("/users/me", h.User.Me)

	workspaces := protected.Group("/workspaces")
	workspaces.Post("/", h.Workspace.Create)
	workspaces.Get("/", h.Workspace.List)
	workspaces.Get("/:id", h.Workspace.Get)
	workspaces.Post("/:workspaceId/files", h.File.Upload)
	workspaces.Get("/:workspaceId/files", h.File.ListByWorkspace)
	workspaces.Get("/:workspaceId/ai-jobs", h.AIJob.ListByWorkspace)

	files := protected.Group("/files")
	files.Get("/:id", h.File.Get)
	files.Get("/:id/download", h.File.Download)

	comments := protected.Group("/files/:fileId/comments")
	comments.Post("/", h.Comment.Create)
	comments.Get("/", h.Comment.List)

	aiJobs := protected.Group("/ai-jobs")
	aiJobs.Post("/", h.AIJob.Create)
	aiJobs.Get("/:id", h.AIJob.Get)
}
