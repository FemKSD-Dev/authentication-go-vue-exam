package handler

import (
	webResponse "authentication-project-exam/internal/adapter/inbound/web/response"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(webResponse.Success[any](nil, "OK"))
}
