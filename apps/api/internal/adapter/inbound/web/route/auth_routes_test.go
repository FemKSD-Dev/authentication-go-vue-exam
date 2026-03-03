package route

import (
	"context"
	"net/http/httptest"
	"testing"

	"authentication-project-exam/internal/adapter/inbound/web/handler"
	"authentication-project-exam/internal/config"
	"authentication-project-exam/internal/core/port/inbound"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterAuthRoutes(t *testing.T) {
	cfg := &config.Config{JWTExpireMinutes: 30, RefreshTokenExpireMinutes: 60}
	authHandler := handler.NewAuthHandler(&fakeAuthPort{}, cfg)
	authMiddleware := func(c *fiber.Ctx) error { return c.Next() }

	app := fiber.New()
	api := app.Group("/api")
	RegisterAuthRoutes(api, authHandler, authMiddleware)

	req := httptest.NewRequest("POST", "/api/register", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected 400 for invalid body, got %d", resp.StatusCode)
	}

	req2 := httptest.NewRequest("GET", "/api/me", nil)
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("expected 401 for me without auth, got %d", resp2.StatusCode)
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
