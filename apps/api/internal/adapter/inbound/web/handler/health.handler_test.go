package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestNewHealthHandler(t *testing.T) {
	h := NewHealthHandler()
	if h == nil {
		t.Fatal("NewHealthHandler returned nil")
	}
}

func TestHealthHandler_HealthCheck(t *testing.T) {
	app := fiber.New()
	app.Get("/health", NewHealthHandler().HealthCheck)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}
