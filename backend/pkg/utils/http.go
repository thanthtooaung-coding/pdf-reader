package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/response"
)

func OK(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	return c.Status(fiber.StatusOK).JSON(response.APIResponse{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, status int, err error) error {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	return c.Status(status).JSON(response.APIResponse{
		Success: false,
		Error:   msg,
	})
}

func BadRequest(c *fiber.Ctx, err error) error  { return Error(c, fiber.StatusBadRequest, err) }
func NotFound(c *fiber.Ctx, err error) error    { return Error(c, fiber.StatusNotFound, err) }
func Internal(c *fiber.Ctx, err error) error    { return Error(c, fiber.StatusInternalServerError, err) }
func Unauthorized(c *fiber.Ctx, err error) error { return Error(c, fiber.StatusUnauthorized, err) }
func Forbidden(c *fiber.Ctx, err error) error   { return Error(c, fiber.StatusForbidden, err) }
func TooManyRequests(c *fiber.Ctx, err error) error {
	return Error(c, fiber.StatusTooManyRequests, err)
}

func ActorID(c *fiber.Ctx) *uuid.UUID {
	if v := c.Locals("userID"); v != nil {
		if id, ok := v.(uuid.UUID); ok {
			return &id
		}
	}
	return nil
}

func MustActorID(c *fiber.Ctx) (uuid.UUID, error) {
	id := ActorID(c)
	if id == nil {
		return uuid.Nil, fiber.ErrUnauthorized
	}
	return *id, nil
}
