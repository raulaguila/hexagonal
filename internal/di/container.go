// Package di provides dependency injection container for the application.
// It centralizes the creation and wiring of all dependencies.
package di

import (
	"log/slog"

	"gorm.io/gorm"

	"github.com/raulaguila/go-api/config"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/repository"
	"github.com/raulaguila/go-api/internal/adapter/driven/storage/redis"
	"github.com/raulaguila/go-api/internal/app"
	"github.com/raulaguila/go-api/internal/core/usecase/auth"
	"github.com/raulaguila/go-api/internal/core/usecase/profile"
	"github.com/raulaguila/go-api/internal/core/usecase/user"
	"github.com/raulaguila/go-api/pkg/loggerx"
)

// Container holds all application dependencies.
// It's the single source of truth for dependency injection.
type Container struct {
	// Infrastructure
	Config *config.Environment
	Log    *loggerx.Logger
	DB     *gorm.DB
	Redis  *redis.Service

	// Repositories
	repositories *app.Repositories
}

// NewContainer creates and initializes a new dependency container
func NewContainer(cfg *config.Environment, log *loggerx.Logger, db *gorm.DB, redis *redis.Service) *Container {
	c := &Container{
		Config: cfg,
		Log:    log,
		DB:     db,
		Redis:  redis,
	}

	c.initRepositories()

	log.Info("Dependency container initialized", slog.Int("repositories", 2), slog.Int("use_cases", 3))

	return c
}

// initRepositories initializes all repository implementations
func (c *Container) initRepositories() {
	profileRepo := repository.NewProfileRepository(c.DB)
	userRepo := repository.NewUserRepository(c.DB)

	// Apply caching decorator if Redis is available
	if c.Redis != nil {
		profileRepo = repository.NewCachedProfileRepository(profileRepo, c.Redis)
		userRepo = repository.NewCachedUserRepository(userRepo, c.Redis)
	}

	c.repositories = &app.Repositories{
		User:    userRepo,
		Profile: profileRepo,
	}
}

// Application returns a fully configured Application instance
func (c *Container) Application() *app.Application {
	return app.New(
		c.Config,
		c.Log,
		auth.NewAuthUseCase(c.repositories.User, auth.Config{
			AccessPrivateKey:  c.Config.AccessPrivateKey,
			AccessExpiration:  c.Config.AccessExpiration,
			RefreshPrivateKey: c.Config.RefreshPrivateKey,
			RefreshExpiration: c.Config.RefreshExpiration,
		}),
		profile.NewProfileUseCase(c.repositories.Profile),
		user.NewUserUseCase(c.repositories.User),
		c.repositories,
	)
}
