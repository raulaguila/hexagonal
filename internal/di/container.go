// Package di provides dependency injection container for the application.
// It centralizes the creation and wiring of all dependencies.
package di

import (
	"log/slog"

	"gorm.io/gorm"

	"github.com/raulaguila/go-api/config"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/repository"
	"github.com/raulaguila/go-api/internal/app"
	"github.com/raulaguila/go-api/internal/core/port/input"
	"github.com/raulaguila/go-api/internal/core/port/output"
	"github.com/raulaguila/go-api/internal/core/usecase/auth"
	"github.com/raulaguila/go-api/internal/core/usecase/profile"
	"github.com/raulaguila/go-api/internal/core/usecase/user"
	"github.com/raulaguila/go-api/pkg/logger"
)

// Container holds all application dependencies.
// It's the single source of truth for dependency injection.
type Container struct {
	// Infrastructure
	Config *config.Config
	Log    *logger.Logger
	DB     *gorm.DB

	// Repositories (Output Ports)
	repositories *app.Repositories

	// Use Cases (Input Ports)
	useCases *useCases
}

// useCases holds all use case implementations
type useCases struct {
	Auth    input.AuthUseCase
	Profile input.ProfileUseCase
	User    input.UserUseCase
}

// NewContainer creates and initializes a new dependency container
func NewContainer(cfg *config.Config, log *logger.Logger, db *gorm.DB) *Container {
	c := &Container{
		Config: cfg,
		Log:    log,
		DB:     db,
	}

	c.initRepositories()
	c.initUseCases()

	log.Info("Dependency container initialized", slog.Int("repositories", 2), slog.Int("use_cases", 3))

	return c
}

// initRepositories initializes all repository implementations
func (c *Container) initRepositories() {
	c.repositories = &app.Repositories{
		User:    repository.NewUserRepository(c.DB),
		Profile: repository.NewProfileRepository(c.DB),
	}
}

// initUseCases initializes all use case implementations
func (c *Container) initUseCases() {
	c.useCases = &useCases{
		Auth: auth.NewAuthUseCase(c.repositories.User, auth.Config{
			AccessPrivateKey:  c.Config.AccessPrivateKey,
			AccessExpiration:  c.Config.AccessExpiration,
			RefreshPrivateKey: c.Config.RefreshPrivateKey,
			RefreshExpiration: c.Config.RefreshExpiration,
		}),
		Profile: profile.NewProfileUseCase(c.repositories.Profile),
		User:    user.NewUserUseCase(c.repositories.User),
	}
}

// Application returns a fully configured Application instance
func (c *Container) Application() *app.Application {
	return app.New(
		c.Config,
		c.Log,
		c.useCases.Auth,
		c.useCases.Profile,
		c.useCases.User,
		c.repositories,
	)
}

// Repositories returns the repositories
func (c *Container) Repositories() *app.Repositories {
	return c.repositories
}

// UserRepository returns the user repository
func (c *Container) UserRepository() output.UserRepository {
	return c.repositories.User
}

// ProfileRepository returns the profile repository
func (c *Container) ProfileRepository() output.ProfileRepository {
	return c.repositories.Profile
}

// AuthUseCase returns the auth use case
func (c *Container) AuthUseCase() input.AuthUseCase {
	return c.useCases.Auth
}

// ProfileUseCase returns the profile use case
func (c *Container) ProfileUseCase() input.ProfileUseCase {
	return c.useCases.Profile
}

// UserUseCase returns the user use case
func (c *Container) UserUseCase() input.UserUseCase {
	return c.useCases.User
}
