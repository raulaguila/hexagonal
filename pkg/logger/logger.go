// Package logger provides structured logging for the application.
// It wraps slog with additional functionality for production use.
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"
)

// Level represents logging level
type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// Config holds logger configuration
type Config struct {
	// Level minimum log level
	Level Level

	// Format: "json" or "text"
	Format string

	// Output destination (default: os.Stdout)
	Output io.Writer

	// AddSource adds source file information to logs
	AddSource bool

	// ServiceName for identifying the service in logs
	ServiceName string

	// Version of the application
	Version string

	// Environment (development, staging, production)
	Environment string
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		Level:       LevelInfo,
		Format:      "json",
		Output:      os.Stdout,
		AddSource:   false,
		ServiceName: "api",
		Environment: "development",
	}
}

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
	config Config
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// Init initializes the global logger with given config
func Init(cfg Config) *Logger {
	once.Do(func() {
		defaultLogger = New(cfg)
	})
	return defaultLogger
}

// Default returns the default logger (initializes with defaults if not set)
func Default() *Logger {
	if defaultLogger == nil {
		Init(DefaultConfig())
	}
	return defaultLogger
}

// New creates a new logger instance
func New(cfg Config) *Logger {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.RFC3339Nano))
				}
			}
			return a
		},
	}

	var handler slog.Handler
	if cfg.Format == "text" {
		handler = slog.NewTextHandler(cfg.Output, opts)
	} else {
		handler = slog.NewJSONHandler(cfg.Output, opts)
	}

	baseLogger := slog.New(handler).With(
		slog.String("service", cfg.ServiceName),
		slog.String("version", cfg.Version),
		slog.String("environment", cfg.Environment),
	)

	return &Logger{
		Logger: baseLogger,
		config: cfg,
	}
}

// With returns a new Logger with additional attributes
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
		config: l.config,
	}
}

// Startup logs application startup information
func (l *Logger) Startup(port, version string) {
	l.Info("Application starting",
		slog.String("port", port),
		slog.String("go_version", runtime.Version()),
		slog.Int("num_cpu", runtime.NumCPU()),
		slog.Int("gomaxprocs", runtime.GOMAXPROCS(0)),
	)
}

// Shutdown logs application shutdown
func (l *Logger) Shutdown(reason string) {
	l.Info("Application shutting down",
		slog.String("reason", reason),
	)
}

// HTTPRequest logs an HTTP request
func (l *Logger) HTTPRequest(method, path string, status int, latency time.Duration, ip string) {
	level := LevelInfo
	if status >= 500 {
		level = LevelError
	} else if status >= 400 {
		level = LevelWarn
	}

	l.Log(context.Background(), level, "HTTP request",
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status", status),
		slog.Duration("latency", latency),
		slog.Float64("latency_ms", float64(latency.Nanoseconds())/1e6),
		slog.String("ip", ip),
	)
}

// Global convenience functions

func Debug(msg string, args ...any) {
	Default().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Default().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Default().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Default().Error(msg, args...)
}

func With(args ...any) *Logger {
	return Default().With(args...)
}
