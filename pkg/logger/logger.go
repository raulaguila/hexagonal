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

	// ExternalHandler optional handler for external systems (e.g., Elasticsearch)
	ExternalHandler slog.Handler
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
	config   Config
	handlers []slog.Handler
	mu       sync.RWMutex
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
			// Customize time format for Elasticsearch compatibility
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

	// Create multi-handler if external handler is provided
	handlers := []slog.Handler{handler}
	if cfg.ExternalHandler != nil {
		handlers = append(handlers, cfg.ExternalHandler)
	}

	var finalHandler slog.Handler
	if len(handlers) > 1 {
		finalHandler = &multiHandler{handlers: handlers}
	} else {
		finalHandler = handler
	}

	// Add default attributes
	baseLogger := slog.New(finalHandler).With(
		slog.String("service", cfg.ServiceName),
		slog.String("version", cfg.Version),
		slog.String("environment", cfg.Environment),
	)

	return &Logger{
		Logger:   baseLogger,
		config:   cfg,
		handlers: handlers,
	}
}

// multiHandler allows logging to multiple handlers
type multiHandler struct {
	handlers []slog.Handler
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &multiHandler{handlers: newHandlers}
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &multiHandler{handlers: newHandlers}
}

// AddHandler adds an external handler dynamically
func (l *Logger) AddHandler(handler slog.Handler) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.handlers = append(l.handlers, handler)
}

// With returns a new Logger with additional attributes
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger:   l.Logger.With(args...),
		config:   l.config,
		handlers: l.handlers,
	}
}

// WithContext returns a new Logger with context-derived attributes
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Extract common context values (request ID, user ID, etc.)
	newLogger := l
	if requestID := ctx.Value("request_id"); requestID != nil {
		newLogger = newLogger.With(slog.String("request_id", requestID.(string)))
	}
	if userID := ctx.Value("user_id"); userID != nil {
		newLogger = newLogger.With(slog.Any("user_id", userID))
	}
	return newLogger
}

// WithError returns a logger with error information
func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}
	return l.With(
		slog.String("error", err.Error()),
		slog.String("error_type", getErrorType(err)),
	)
}

// WithRequest returns a logger with HTTP request information
func (l *Logger) WithRequest(method, path, ip, userAgent string) *Logger {
	return l.With(
		slog.Group("request",
			slog.String("method", method),
			slog.String("path", path),
			slog.String("ip", ip),
			slog.String("user_agent", userAgent),
		),
	)
}

// WithDuration returns a logger with duration information
func (l *Logger) WithDuration(d time.Duration) *Logger {
	return l.With(
		slog.Duration("duration", d),
		slog.Float64("duration_ms", float64(d.Nanoseconds())/1e6),
	)
}

// WithComponent returns a logger for a specific component
func (l *Logger) WithComponent(name string) *Logger {
	return l.With(slog.String("component", name))
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

// DatabaseConnected logs successful database connection
func (l *Logger) DatabaseConnected(host, port, database string) {
	l.Info("Database connected",
		slog.Group("database",
			slog.String("host", host),
			slog.String("port", port),
			slog.String("name", database),
		),
	)
}

// DatabaseError logs database errors
func (l *Logger) DatabaseError(operation string, err error) {
	l.WithError(err).Error("Database operation failed",
		slog.String("operation", operation),
	)
}

// StorageConnected logs successful storage connection
func (l *Logger) StorageConnected(host, bucket string) {
	l.Info("Object storage connected",
		slog.Group("storage",
			slog.String("host", host),
			slog.String("bucket", bucket),
		),
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

// AuthSuccess logs successful authentication
func (l *Logger) AuthSuccess(userID uint, username string) {
	l.Info("Authentication successful",
		slog.Uint64("user_id", uint64(userID)),
		slog.String("username", username),
	)
}

// AuthFailure logs failed authentication
func (l *Logger) AuthFailure(username, reason string) {
	l.Warn("Authentication failed",
		slog.String("username", username),
		slog.String("reason", reason),
	)
}

// UserCreated logs user creation
func (l *Logger) UserCreated(userID uint, email string) {
	l.Info("User created",
		slog.Uint64("user_id", uint64(userID)),
		slog.String("email", maskEmail(email)),
	)
}

// UserDeleted logs user deletion
func (l *Logger) UserDeleted(userIDs []uint) {
	ids := make([]uint64, len(userIDs))
	for i, id := range userIDs {
		ids[i] = uint64(id)
	}
	l.Info("Users deleted",
		slog.Any("user_ids", ids),
	)
}

// Helper functions

func getErrorType(err error) string {
	if err == nil {
		return ""
	}
	return runtime.FuncForPC(0).Name()
}

func maskEmail(email string) string {
	if len(email) < 5 {
		return "***"
	}
	atIndex := -1
	for i, c := range email {
		if c == '@' {
			atIndex = i
			break
		}
	}
	if atIndex <= 2 {
		return email[:1] + "***" + email[atIndex:]
	}
	return email[:2] + "***" + email[atIndex:]
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
