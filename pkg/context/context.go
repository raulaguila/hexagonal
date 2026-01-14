package context

import (
	"context"

	"github.com/raulaguila/go-api/pkg/logger"
)

// Key types for context values
type contextKey string

const (
	loggerKey    contextKey = "logger"
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
	userNameKey  contextKey = "user_name"
)

// WithLogger adds a logger to the context
func WithLogger(ctx context.Context, log *logger.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

// Logger retrieves the logger from context, returns default if not found
func Logger(ctx context.Context) *logger.Logger {
	if log, ok := ctx.Value(loggerKey).(*logger.Logger); ok {
		return log
	}
	return logger.Default()
}

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestID retrieves the request ID from context
func RequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// WithUserID adds user ID to the context
func WithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserID retrieves the user ID from context
func UserID(ctx context.Context) uint {
	if id, ok := ctx.Value(userIDKey).(uint); ok {
		return id
	}
	return 0
}

// WithUserName adds username to the context
func WithUserName(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, userNameKey, username)
}

// UserName retrieves the username from context
func UserName(ctx context.Context) string {
	if name, ok := ctx.Value(userNameKey).(string); ok {
		return name
	}
	return ""
}

// Enrich adds common context values (request ID, user info) to logger
func Enrich(ctx context.Context) *logger.Logger {
	log := Logger(ctx)

	if requestID := RequestID(ctx); requestID != "" {
		log = log.With("request_id", requestID)
	}

	if userID := UserID(ctx); userID > 0 {
		log = log.With("user_id", userID)
	}

	if userName := UserName(ctx); userName != "" {
		log = log.With("username", userName)
	}

	return log
}
