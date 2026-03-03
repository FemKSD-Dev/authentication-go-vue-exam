package bootstrap

import (
	"authentication-project-exam/internal/adapter/inbound/web/middleware"
	"authentication-project-exam/internal/bootstrap/modules"
	"authentication-project-exam/internal/config"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	Config *config.Config
	HTTP   *fiber.App
}

func (a *App) Run() error {
	return a.HTTP.Listen(":" + a.Config.HTTPPort)
}

func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	appLogger, err := NewLogger()
	if err != nil {
		return nil, err
	}

	db, err := NewDatabase(cfg, appLogger)
	if err != nil {
		return nil, err
	}

	authModule, err := modules.NewAuthModule(db, cfg)
	if err != nil {
		return nil, err
	}

	healthModule := modules.NewHealthModule()

	container := &Container{
		AuthHandler:         authModule.Handler,
		AuthMiddleware:      authModule.AuthMiddleware,
		HealthHandler:       healthModule.Handler,
		AppLogger:           appLogger,
		RequestIDMiddleware: middleware.NewRequestIDMiddleware(),
		LogMiddleware:       middleware.NewRequestLoggerMiddleware(appLogger),
	}

	serverDependencies := NewHTTPServer(container)

	httpApp := serverDependencies.Build()

	return &App{
		Config: cfg,
		HTTP:   httpApp,
	}, nil
}
