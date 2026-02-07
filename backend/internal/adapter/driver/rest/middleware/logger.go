package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/godeh/sloggergo"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

// RequestLogger returns a middleware that logs HTTP requests using sloggergo
func RequestLogger(log *sloggergo.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Call next handler
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Generate RequestID if not present
		reqID := c.Get(fiber.HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
			c.Set(fiber.HeaderXRequestID, reqID)
		}

		// Get TraceID from OpenTelemetry context if available
		traceID := ""
		spanCtx := trace.SpanContextFromContext(c.UserContext())
		if spanCtx.HasTraceID() {
			traceID = spanCtx.TraceID().String()
		}

		// Determine log level
		status := c.Response().StatusCode()

		// Build log fields
		fields := map[string]any{
			"status":     status,
			"method":     c.Method(),
			"path":       c.OriginalURL(),
			"ip":         c.IP(),
			"latency":    duration.String(),
			"request_id": reqID,
			"pid":        os.Getpid(),
		}

		// Get Authorization if available
		authorization := c.Get(fiber.HeaderAuthorization)
		if authorization != "" {
			fields["authorization"] = strings.ReplaceAll(authorization, "Bearer ", "")[:min(len(authorization), 10)]
		}

		if traceID != "" {
			fields["trace_id"] = traceID
		}

		if err != nil {
			fields["error"] = err.Error()
		}

		// Log with context
		logWithFields := log.With(flattenFields(fields)...)

		msg := "HTTP Request"

		switch {
		case status >= 500:
			logWithFields.Error(msg)
		case status >= 400:
			logWithFields.Warn(msg)
		default:
			logWithFields.Info(msg)
		}

		return err
	}
}

// flattenFields converts map to alternating key/value slice
func flattenFields(m map[string]any) []any {
	out := make([]any, 0, len(m)*2)
	for k, v := range m {
		out = append(out, k, v)
	}
	return out
}
