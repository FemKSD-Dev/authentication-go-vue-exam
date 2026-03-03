package bootstrap

import (
	"authentication-project-exam/internal/adapter/inbound/web/handler"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Container struct {
	AuthHandler         *handler.AuthHandler
	AuthMiddleware      fiber.Handler
	HealthHandler       *handler.HealthHandler
	RequestIDMiddleware fiber.Handler
	LogMiddleware       fiber.Handler
	AppLogger           *zap.Logger
}
