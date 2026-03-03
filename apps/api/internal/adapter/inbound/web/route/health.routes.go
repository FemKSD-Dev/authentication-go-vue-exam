package route

import (
	"authentication-project-exam/internal/adapter/inbound/web/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterHealthRoutes(r fiber.Router, healthHandler *handler.HealthHandler) {
	r.Get("/health", healthHandler.HealthCheck)
}
