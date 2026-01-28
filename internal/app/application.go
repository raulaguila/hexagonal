// Package app provides the Application layer that unifies all use cases
// and can be used by any interface (REST, gRPC, CLI, etc.)
package app

import (
	"github.com/raulaguila/go-api/config"
	"github.com/raulaguila/go-api/internal/core/port/input"
	"github.com/raulaguila/go-api/internal/core/port/output"
	"github.com/raulaguila/go-api/pkg/loggerx"
)

// Application is the main entry point for all business operations.
// It holds all use cases and can be injected into any interface adapter.
type Application struct {
	// Configuration
	Config *config.Environment

	// Logging
	Log *loggerx.Logger

	// Use Cases (Input Ports)
	Auth    input.AuthUseCase
	Profile input.ProfileUseCase
	User    input.UserUseCase

	// Repositories (Output Ports) - exposed for adapters that need direct access
	Repositories *Repositories
}

// Repositories holds all repository implementations
type Repositories struct {
	User    output.UserRepository
	Profile output.ProfileRepository
}

// Options holds optional dependencies for the application
type Options struct {
	// ExternalLogHandler for sending logs to external systems
	ExternalLogHandler any
}

// Option is a function that configures the Application
type Option func(*Application)

// WithLogger sets a custom logger
func WithLogger(log *loggerx.Logger) Option {
	return func(a *Application) {
		a.Log = log
	}
}

// New creates a new Application instance with all dependencies wired up
func New(
	cfg *config.Environment,
	log *loggerx.Logger,
	authUC input.AuthUseCase,
	profileUC input.ProfileUseCase,
	userUC input.UserUseCase,
	repos *Repositories,
	opts ...Option,
) *Application {
	app := &Application{
		Config:       cfg,
		Log:          log,
		Auth:         authUC,
		Profile:      profileUC,
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
