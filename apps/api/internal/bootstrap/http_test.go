package bootstrap

import (
	"context"
	"testing"

	"authentication-project-exam/internal/adapter/inbound/web/handler"
	"authentication-project-exam/internal/config"
	"authentication-project-exam/internal/core/port/inbound"

	"github.com/gofiber/fiber/v2"
)

func TestNewHTTPServer_NilPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when container is nil")
		}
	}()
	NewHTTPServer(nil)
}

func TestNewHTTPServer_Success(t *testing.T) {
	cfg := &config.Config{JWTExpireMinutes: 30, RefreshTokenExpireMinutes: 60}
	authHandler := handler.NewAuthHandler(&fakeAuthPort{}, cfg)
	healthHandler := handler.NewHealthHandler()

	container := &Container{
		AuthHandler:         authHandler,
		HealthHandler:       healthHandler,
		RequestIDMiddleware: func(c *fiber.Ctx) error { return c.Next() },
		LogMiddleware:      func(c *fiber.Ctx) error { return c.Next() },
		AuthMiddleware:      func(c *fiber.Ctx) error { return c.Next() },
	}

	deps := NewHTTPServer(container)
	if deps == nil {
		t.Fatal("expected non-nil ServerDependencies")
	}
	if deps.AuthHandler != authHandler {
		t.Error("AuthHandler not set correctly")
	}
	if deps.HealthHandler != healthHandler {
		t.Error("HealthHandler not set correctly")
	}
}

type fakeAuthPort struct{}

func (fakeAuthPort) Register(ctx context.Context, p inbound.RegisterPayload) (*inbound.RegisterResult, error) {
	return nil, nil
}
func (fakeAuthPort) Login(ctx context.Context, p inbound.LoginPayload) (*inbound.LoginResult, error) {
	return nil, nil
}
func (fakeAuthPort) Me(ctx context.Context, q inbound.MeQuery) (*inbound.MeResult, error) {
	return nil, nil
}
func (fakeAuthPort) RefreshToken(ctx context.Context, p inbound.RefreshTokenPayload) (*inbound.RefreshTokenResult, error) {
	return nil, nil
}
