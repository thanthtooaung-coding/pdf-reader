package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/config"
	"github.com/thanthtooaung-coding/pdf-reader/backend/pkg/utils"
)

func JWTAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return utils.Unauthorized(c, fiber.ErrUnauthorized)
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		userID, _, _, err := utils.ParseAccessToken(cfg.JWTSecret, token)
		if err != nil {
			return utils.Unauthorized(c, err)
		}
		c.Locals("userID", userID)
		return c.Next()
	}
}

func OptionalUserID(c *fiber.Ctx) uuid.UUID {
	if v := c.Locals("userID"); v != nil {
		if id, ok := v.(uuid.UUID); ok {
			return id
		}
	}
	return uuid.Nil
}
