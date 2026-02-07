package app

import (
	"github.com/godeh/sloggergo"

	"github.com/raulaguila/go-api/config"
	"github.com/raulaguila/go-api/internal/adapter/driven/storage/redis"
	"github.com/raulaguila/go-api/internal/core/port/input"
	"github.com/raulaguila/go-api/internal/core/port/output"
)

// Application is the main entry point for all business operations.
// It holds all use cases and can be injected into any interface adapter.
type Application struct {
	// Configuration
	Config *config.Environment

	// Logging
	Log *sloggergo.Logger

	// Use Cases (Input Ports)
	Auth input.AuthUseCase
	Role input.RoleUseCase
	User input.UserUseCase

	// Repositories (Output Ports) - exposed for adapters that need direct access
	Repositories *Repositories

	// Redis service for cache operations and health checks
	Redis *redis.Service

	// Auditor for audit logging
	Auditor input.AuditorUseCase
}

// Repositories holds all repository implementations
type Repositories struct {
	User  output.UserRepository
	Role  output.RoleRepository
	Token output.TokenRepository
	Audit output.AuditRepository
}

// Options holds optional dependencies for the application
type Options struct {
	// ExternalLogHandler for sending logs to external systems
	ExternalLogHandler any
}

// Option is a function that configures the Application
type Option func(*Application)

// WithLogger sets a custom logger
func WithLogger(log *sloggergo.Logger) Option {
	return func(a *Application) {
		a.Log = log
	}
}

// WithRedis sets the Redis service for health checks
func WithRedis(redis *redis.Service) Option {
	return func(a *Application) {
		a.Redis = redis
	}
}

// WithAuditor sets the auditor service
func WithAuditor(aud input.AuditorUseCase) Option {
	return func(a *Application) {
		a.Auditor = aud
	}
}

// New creates a new Application instance with all dependencies wired up
func New(
	cfg *config.Environment,
	log *sloggergo.Logger,
	authUC input.AuthUseCase,
	roleUC input.RoleUseCase,
	userUC input.UserUseCase,
	repos *Repositories,
	opts ...Option,
) *Application {
	app := &Application{
		Config:       cfg,
		Log:          log,
		Auth:         authUC,
		Role:         roleUC,
		User:         userUC,
		Repositories: repos,
	}

	// Apply options
	for _, opt := range opts {
		opt(app)
	}

	return app
}

// Version returns the application version
func (a *Application) Version() string {
	return a.Config.Version
}

// IsProduction checks if running in production
func (a *Application) IsProduction() bool {
	return a.Config.Environment == "production"
}
