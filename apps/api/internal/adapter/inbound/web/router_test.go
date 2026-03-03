package web

import (
	"context"
	"net/http/httptest"
	"testing"

	"authentication-project-exam/internal/adapter/inbound/web/handler"
	"authentication-project-exam/internal/config"
	"authentication-project-exam/internal/core/port/inbound"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestServerDependencies_Build(t *testing.T) {
	cfg := &config.Config{JWTExpireMinutes: 30, RefreshTokenExpireMinutes: 60}
	authHandler := handler.NewAuthHandler(&fakeAuthPort{}, cfg)
	healthHandler := handler.NewHealthHandler()

	deps := &ServerDependencies{
		AuthHandler:         authHandler,
		HealthHandler:       healthHandler,
		RequestIDMiddleware: func(c *fiber.Ctx) error { c.Locals("requestID", uuid.New().String()); return c.Next() },
		LogMiddleware:      func(c *fiber.Ctx) error { return c.Next() },
		AuthMiddleware:     func(c *fiber.Ctx) error { return c.Next() },
	}

	app := deps.Build()
	if app == nil {
		t.Fatal("Build returned nil")
	}

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
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
