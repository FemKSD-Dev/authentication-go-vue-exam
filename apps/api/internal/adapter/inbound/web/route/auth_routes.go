package route

import (
	"authentication-project-exam/internal/adapter/inbound/web/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(
	r fiber.Router,
	authHandler *handler.AuthHandler,
	authMiddleware fiber.Handler,
) {
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Get("/me", authMiddleware, authHandler.Me)
}
