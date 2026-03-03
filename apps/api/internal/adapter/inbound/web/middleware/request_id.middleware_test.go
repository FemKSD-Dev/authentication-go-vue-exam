package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestNewRequestIDMiddleware_WithHeader(t *testing.T) {
	app := fiber.New()
	app.Use(NewRequestIDMiddleware())
	app.Get("/", func(c *fiber.Ctx) error {
		requestID := c.Locals(RequestIDKey)
		if requestID != "custom-request-id" {
			t.Errorf("expected custom-request-id, got %v", requestID)
		}
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", "custom-request-id")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
	if resp.Header.Get("X-Request-ID") != "custom-request-id" {
		t.Errorf("expected X-Request-ID header custom-request-id, got %s", resp.Header.Get("X-Request-ID"))
	}
}

func TestNewRequestIDMiddleware_WithoutHeader(t *testing.T) {
	app := fiber.New()
	app.Use(NewRequestIDMiddleware())
	app.Get("/", func(c *fiber.Ctx) error {
		requestID := c.Locals(RequestIDKey).(string)
		if requestID == "" {
			t.Error("expected non-empty request ID")
		}
		if _, err := uuid.Parse(requestID); err != nil {
			t.Errorf("expected valid UUID, got %s: %v", requestID, err)
		}
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
	if resp.Header.Get("X-Request-ID") == "" {
		t.Error("expected X-Request-ID header to be set")
	}
}
