package middleware

import (
	"strings"

	webResponse "authentication-project-exam/internal/adapter/inbound/web/response"
	"authentication-project-exam/internal/core/port/outbound"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	tokenManager outbound.TokenManager
}

func NewAuthMiddleware(tokenManager outbound.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenManager: tokenManager}
}

func (m *AuthMiddleware) VerifyToken(c *fiber.Ctx) error {
	accessToken := c.Cookies("access_token")
	if accessToken == "" {
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(webResponse.Fail[any](webResponse.CodeUnauthorized, "access token is required"))
	}

	claims, err := m.tokenManager.Verify(accessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(webResponse.Fail[any](webResponse.CodeUnauthorized, "invalid access token"))
	}

	c.Locals("user_id", claims.UserID)
	return c.Next()
}
