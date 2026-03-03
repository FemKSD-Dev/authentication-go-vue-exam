package route

import (
	"net/http/httptest"
	"testing"

	"authentication-project-exam/internal/adapter/inbound/web/handler"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterHealthRoutes(t *testing.T) {
	healthHandler := handler.NewHealthHandler()
	app := fiber.New()
	api := app.Group("/api")
	RegisterHealthRoutes(api, healthHandler)

	req := httptest.NewRequest("GET", "/api/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
