package handler

import (
	"authentication-project-exam/internal/config"
	"authentication-project-exam/internal/core/port/inbound"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type fakeAuthPort struct {
	registerFn     func(ctx context.Context, payload inbound.RegisterPayload) (*inbound.RegisterResult, error)
	loginFn        func(ctx context.Context, payload inbound.LoginPayload) (*inbound.LoginResult, error)
	meFn           func(ctx context.Context, query inbound.MeQuery) (*inbound.MeResult, error)
	refreshTokenFn func(ctx context.Context, payload inbound.RefreshTokenPayload) (*inbound.RefreshTokenResult, error)
}

func (f *fakeAuthPort) Register(ctx context.Context, payload inbound.RegisterPayload) (*inbound.RegisterResult, error) {
	if f.registerFn != nil {
		return f.registerFn(ctx, payload)
	}
	return nil, nil
}

func (f *fakeAuthPort) Login(ctx context.Context, payload inbound.LoginPayload) (*inbound.LoginResult, error) {
	if f.loginFn != nil {
		return f.loginFn(ctx, payload)
	}
	return nil, nil
}

func (f *fakeAuthPort) Me(ctx context.Context, query inbound.MeQuery) (*inbound.MeResult, error) {
	if f.meFn != nil {
		return f.meFn(ctx, query)
	}
	return nil, nil
}

func (f *fakeAuthPort) RefreshToken(ctx context.Context, payload inbound.RefreshTokenPayload) (*inbound.RefreshTokenResult, error) {
	if f.refreshTokenFn != nil {
		return f.refreshTokenFn(ctx, payload)
	}
	return nil, nil
}

func defaultConfig() *config.Config {
	return &config.Config{
		JWTExpireMinutes:          30,
		RefreshTokenExpireMinutes: 60,
	}
}

func TestNewAuthHandler(t *testing.T) {
	fake := &fakeAuthPort{}
	cfg := defaultConfig()
	h := NewAuthHandler(fake, cfg)
	if h == nil {
		t.Fatal("NewAuthHandler returned nil")
	}
	if h.auth != fake {
		t.Error("auth port not set correctly")
	}
	if h.config != cfg {
		t.Error("config not set correctly")
	}
}

func TestAuthHandler_Register(t *testing.T) {
	cfg := defaultConfig()

	t.Run("invalid request body", func(t *testing.T) {
		app := fiber.New()
		app.Post("/register", NewAuthHandler(&fakeAuthPort{}, cfg).Register)

		req := httptest.NewRequest("POST", "/register", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("empty username", func(t *testing.T) {
		app := fiber.New()
		app.Post("/register", NewAuthHandler(&fakeAuthPort{}, cfg).Register)

		body, _ := json.Marshal(map[string]string{
			"username":         "",
			"password":         "password123",
			"confirm_password": "password123",
		})
		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("whitespace only username", func(t *testing.T) {
		app := fiber.New()
		app.Post("/register", NewAuthHandler(&fakeAuthPort{}, cfg).Register)

		body, _ := json.Marshal(map[string]string{
			"username":         "   ",
			"password":         "password123",
			"confirm_password": "password123",
		})
		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("empty password", func(t *testing.T) {
		app := fiber.New()
		app.Post("/register", NewAuthHandler(&fakeAuthPort{}, cfg).Register)

		body, _ := json.Marshal(map[string]string{
			"username":         "user1",
			"password":         "",
			"confirm_password": "",
		})
		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("register returns error", func(t *testing.T) {
		fake := &fakeAuthPort{
			registerFn: func(ctx context.Context, payload inbound.RegisterPayload) (*inbound.RegisterResult, error) {
				return nil, errors.New("username exists")
			},
		}
		app := fiber.New()
		app.Post("/register", NewAuthHandler(fake, cfg).Register)

		body, _ := json.Marshal(map[string]string{
			"username":         "user1",
			"password":         "password123",
			"confirm_password": "password123",
		})
		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
		}
	})

	t.Run("success", func(t *testing.T) {
		fake := &fakeAuthPort{
			registerFn: func(ctx context.Context, payload inbound.RegisterPayload) (*inbound.RegisterResult, error) {
				return nil, nil
			},
		}
		app := fiber.New()
		app.Post("/register", NewAuthHandler(fake, cfg).Register)

		body, _ := json.Marshal(map[string]string{
			"username":         "user1",
			"password":         "password123",
			"confirm_password": "password123",
		})
		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusCreated {
			t.Errorf("expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
		}
	})
}

func TestAuthHandler_Login(t *testing.T) {
	cfg := defaultConfig()

	t.Run("invalid request body", func(t *testing.T) {
		app := fiber.New()
		app.Post("/login", NewAuthHandler(&fakeAuthPort{}, cfg).Login)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("empty username", func(t *testing.T) {
		app := fiber.New()
		app.Post("/login", NewAuthHandler(&fakeAuthPort{}, cfg).Login)

		body, _ := json.Marshal(map[string]string{
			"username": "",
			"password": "password123",
		})
		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("whitespace only username", func(t *testing.T) {
		app := fiber.New()
		app.Post("/login", NewAuthHandler(&fakeAuthPort{}, cfg).Login)

		body, _ := json.Marshal(map[string]string{
			"username": "   ",
			"password": "password123",
		})
		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("empty password", func(t *testing.T) {
		app := fiber.New()
		app.Post("/login", NewAuthHandler(&fakeAuthPort{}, cfg).Login)

		body, _ := json.Marshal(map[string]string{
			"username": "user1",
			"password": "",
		})
		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("login returns error", func(t *testing.T) {
		fake := &fakeAuthPort{
			loginFn: func(ctx context.Context, payload inbound.LoginPayload) (*inbound.LoginResult, error) {
				return nil, errors.New("invalid credentials")
			},
		}
		app := fiber.New()
		app.Post("/login", NewAuthHandler(fake, cfg).Login)

		body, _ := json.Marshal(map[string]string{
			"username": "user1",
			"password": "wrongpassword",
		})
		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
		}
	})

	t.Run("success", func(t *testing.T) {
		fake := &fakeAuthPort{
			loginFn: func(ctx context.Context, payload inbound.LoginPayload) (*inbound.LoginResult, error) {
				return &inbound.LoginResult{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
				}, nil
			},
		}
		app := fiber.New()
		app.Post("/login", NewAuthHandler(fake, cfg).Login)

		body, _ := json.Marshal(map[string]string{
			"username": "user1",
			"password": "password123",
		})
		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
		}
		// Verify cookies are set
		cookies := resp.Cookies()
		if len(cookies) < 2 {
			t.Errorf("expected at least 2 cookies, got %d", len(cookies))
		}
	})
}

func TestAuthHandler_Me(t *testing.T) {
	cfg := defaultConfig()

	t.Run("user_id not in locals", func(t *testing.T) {
		app := fiber.New()
		app.Get("/me", NewAuthHandler(&fakeAuthPort{}, cfg).Me)

		req := httptest.NewRequest("GET", "/me", nil)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
		}
	})

	t.Run("user_id wrong type", func(t *testing.T) {
		app := fiber.New()
		app.Get("/me", func(c *fiber.Ctx) error {
			c.Locals("user_id", 123) // int, not string
			return NewAuthHandler(&fakeAuthPort{}, cfg).Me(c)
		})

		req := httptest.NewRequest("GET", "/me", nil)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
		}
	})

	t.Run("me returns error", func(t *testing.T) {
		fake := &fakeAuthPort{
			meFn: func(ctx context.Context, query inbound.MeQuery) (*inbound.MeResult, error) {
				return nil, errors.New("user not found")
			},
		}
		app := fiber.New()
		app.Get("/me", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user-123")
			return NewAuthHandler(fake, cfg).Me(c)
		})

		req := httptest.NewRequest("GET", "/me", nil)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
		}
	})

	t.Run("success", func(t *testing.T) {
		fake := &fakeAuthPort{
			meFn: func(ctx context.Context, query inbound.MeQuery) (*inbound.MeResult, error) {
				return &inbound.MeResult{
					UserID:   "user-123",
					Username: "user1",
				}, nil
			},
		}
		app := fiber.New()
		app.Get("/me", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user-123")
			return NewAuthHandler(fake, cfg).Me(c)
		})

		req := httptest.NewRequest("GET", "/me", nil)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
		}
	})
}
