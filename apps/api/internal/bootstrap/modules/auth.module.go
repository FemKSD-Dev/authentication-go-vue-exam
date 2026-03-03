package modules

import (
	"authentication-project-exam/internal/adapter/inbound/web/handler"
	"authentication-project-exam/internal/adapter/inbound/web/middleware"
	"authentication-project-exam/internal/adapter/outbound/persistence/postgres"
	"authentication-project-exam/internal/adapter/outbound/security"
	"authentication-project-exam/internal/config"
	"authentication-project-exam/internal/core/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthModule struct {
	Handler        *handler.AuthHandler
	AuthMiddleware fiber.Handler
}

func NewAuthModule(db *gorm.DB, cfg *config.Config) (*AuthModule, error) {
	usersRepository := postgres.NewUsersRepository(db)
	passwordManager := security.NewArgon2IDPasswordEncoder()
	tokenManager := security.NewJWTTokenIssuer(cfg.JWTSecret, cfg.JWTExpireMinutes, cfg.RefreshTokenExpireMinutes)

	authService := service.NewAuthService(usersRepository, passwordManager, tokenManager)
	authHandler := handler.NewAuthHandler(authService, cfg)
	authMiddleware := middleware.NewAuthMiddleware(tokenManager)

	return &AuthModule{
		Handler:        authHandler,
		AuthMiddleware: authMiddleware.VerifyToken,
	}, nil
}
