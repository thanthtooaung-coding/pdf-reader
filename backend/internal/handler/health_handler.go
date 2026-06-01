package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

// HealthCheck returns service liveness status.
//
// @Summary      Health check
// @Description  Returns ok when the API process is running.
// @Tags         health
// @Produce      json
// @Success      200  {object}  response.HealthResponse
// @Router       /healthz [get]
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(response.HealthResponse{
		Status:  "ok",
		Service: "pdf-reader-backend",
	})
}
