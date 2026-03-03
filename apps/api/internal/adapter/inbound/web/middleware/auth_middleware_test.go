package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"authentication-project-exam/internal/core/port/outbound"

	"github.com/gofiber/fiber/v2"
)

type fakeTokenManager struct {
	verifyFn func(token string) (*outbound.TokenClaims, error)
}

func (f *fakeTokenManager) Issue(userID, username string) (string, string, error) {
	return "access", "refresh", nil
}

func (f *fakeTokenManager) Verify(token string) (*outbound.TokenClaims, error) {
	if f.verifyFn != nil {
		return f.verifyFn(token)
	}
	return &outbound.TokenClaims{UserID: "user-1", Username: "john"}, nil
}

func TestNewAuthMiddleware(t *testing.T) {
	m := NewAuthMiddleware(&fakeTokenManager{})
	if m == nil {
		t.Fatal("NewAuthMiddleware returned nil")
	}
}

func TestAuthMiddleware_VerifyToken_NoToken(t *testing.T) {
	app := fiber.New()
	app.Get("/protected", NewAuthMiddleware(&fakeTokenManager{}).VerifyToken, func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAuthMiddleware_VerifyToken_FromCookie(t *testing.T) {
	fake := &fakeTokenManager{
		verifyFn: func(token string) (*outbound.TokenClaims, error) {
			if token != "cookie-token" {
				t.Errorf("expected cookie-token, got %s", token)
			}
			return &outbound.TokenClaims{UserID: "user-1", Username: "john"}, nil
		},
	}
	app := fiber.New()
	app.Get("/protected", NewAuthMiddleware(fake).VerifyToken, func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		if userID != "user-1" {
			t.Errorf("expected user_id user-1, got %v", userID)
		}
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "cookie-token"})
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestAuthMiddleware_VerifyToken_FromBearerHeader(t *testing.T) {
	fake := &fakeTokenManager{
		verifyFn: func(token string) (*outbound.TokenClaims, error) {
			if token != "bearer-token" {
				t.Errorf("expected bearer-token, got %s", token)
			}
			return &outbound.TokenClaims{UserID: "user-2", Username: "jane"}, nil
		},
	}
	app := fiber.New()
	app.Get("/protected", NewAuthMiddleware(fake).VerifyToken, func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		if userID != "user-2" {
			t.Errorf("expected user_id user-2, got %v", userID)
		}
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer bearer-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestAuthMiddleware_VerifyToken_InvalidToken(t *testing.T) {
	fake := &fakeTokenManager{
		verifyFn: func(token string) (*outbound.TokenClaims, error) {
			return nil, fiber.ErrUnauthorized
		},
	}
	app := fiber.New()
	app.Get("/protected", NewAuthMiddleware(fake).VerifyToken, func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}
