package logger

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// FiberMiddleware creates a Fiber middleware for HTTP request logging
func FiberMiddleware(log *Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		status := c.Response().StatusCode()

		// Log the request
		log.HTTPRequest(
			c.Method(),
			c.OriginalURL(),
			status,
			latency,
			c.IP(),
		)

		return err
	}
}

// RequestLogger returns a logger with request context
func RequestLogger(log *Logger, c *fiber.Ctx) *Logger {
	requestID := c.Get("X-Request-ID", "")
	if requestID == "" {
		requestID = c.GetRespHeader("X-Request-ID", "")
	}

	return log.With(
		slog.String("request_id", requestID),
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.String("ip", c.IP()),
	)
}

// ContextWithLogger adds logger to context
func ContextWithLogger(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, log)
}

// FromContext retrieves logger from context
func FromContext(ctx context.Context) *Logger {
	if log, ok := ctx.Value(loggerKey{}).(*Logger); ok {
		return log
	}
	return Default()
}

type loggerKey struct{}
