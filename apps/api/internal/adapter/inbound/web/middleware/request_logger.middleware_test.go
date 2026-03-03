package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func TestNewRequestLoggerMiddleware_Success(t *testing.T) {
	logger := zap.NewNop()
	app := fiber.New()
	app.Use(NewRequestIDMiddleware())
	app.Use(NewRequestLoggerMiddleware(logger))
	app.Get("/", func(c *fiber.Ctx) error {
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
}

func TestNewRequestLoggerMiddleware_ClientError(t *testing.T) {
	logger := zap.NewNop()
	app := fiber.New()
	app.Use(NewRequestIDMiddleware())
	app.Use(NewRequestLoggerMiddleware(logger))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusBadRequest).SendString("bad")
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestNewRequestLoggerMiddleware_ServerError(t *testing.T) {
	logger := zap.NewNop()
	app := fiber.New()
	app.Use(NewRequestIDMiddleware())
	app.Use(NewRequestLoggerMiddleware(logger))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusInternalServerError).SendString("error")
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
	}
}

func TestNewRequestLoggerMiddleware_HandlerError(t *testing.T) {
	logger := zap.NewNop()
	app := fiber.New()
	app.Use(NewRequestIDMiddleware())
	app.Use(NewRequestLoggerMiddleware(logger))
	app.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusInternalServerError, "handler error")
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
	}
}
