package middleware

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/raulaguila/go-api/pkg/loggerx"
	"go.opentelemetry.io/otel/trace"
)

// RequestLogger returns a middleware that logs HTTP requests using loggerx
func RequestLogger(log *loggerx.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

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

		// Call next handler
		err := c.Next()

		// Determine log level
		duration := time.Since(start)
		status := c.Response().StatusCode()

		// Build log fields
		fields := map[string]any{
			"status":     status,
			"method":     c.Method(),
			"path":       c.Path(),
			"ip":         c.IP(),
			"latency":    duration.String(),
			"request_id": reqID,
			"pid":        os.Getpid(),
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
