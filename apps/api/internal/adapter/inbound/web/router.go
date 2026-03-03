package web

import (
	"authentication-project-exam/internal/adapter/inbound/web/handler"
	"authentication-project-exam/internal/adapter/inbound/web/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type ServerDependencies struct {
	AuthHandler         *handler.AuthHandler
	AuthMiddleware      fiber.Handler
	HealthHandler       *handler.HealthHandler
	RequestIDMiddleware fiber.Handler
	LogMiddleware       fiber.Handler
}

func (d *ServerDependencies) Build() *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	app.Use(d.RequestIDMiddleware)
	app.Use(d.LogMiddleware)

	api := app.Group("/api/v1")

	route.RegisterAuthRoutes(api, d.AuthHandler, d.AuthMiddleware)
	route.RegisterHealthRoutes(api, d.HealthHandler)
	return app
}
