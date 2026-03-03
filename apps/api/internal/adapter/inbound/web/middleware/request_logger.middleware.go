package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	iconSuccess = "✓"
	iconClient  = "⚠"
	iconServer  = "✗"
)

func NewRequestLoggerMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()

		requestID, _ := c.Locals(RequestIDKey).(string)

		// Compact, readable format: METHOD path → STATUS latency
		msg := fmt.Sprintf("%s %s → %d %s",
			c.Method(),
			c.OriginalURL(),
			status,
			latency.Round(time.Microsecond),
		)

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.OriginalURL()),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("ip", c.IP()),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Error(iconServer+" "+msg, fields...)
			return err
		}

		switch {
		case status >= 500:
			logger.Error(iconServer+" "+msg, fields...)
		case status >= 400:
			logger.Warn(iconClient+" "+msg, fields...)
		default:
			logger.Info(iconSuccess+" "+msg, fields...)
		}

		return nil
	}
}
