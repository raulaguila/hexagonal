package di

import (
	"log/slog"

	"gorm.io/gorm"

	"github.com/raulaguila/go-api/config"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/repository"
	"github.com/raulaguila/go-api/internal/adapter/driven/storage/redis"
	"github.com/raulaguila/go-api/internal/app"
	"github.com/raulaguila/go-api/internal/core/port/output"
	"github.com/raulaguila/go-api/internal/core/usecase/auth"
	"github.com/raulaguila/go-api/internal/core/usecase/role"
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
	roleRepo := repository.NewRoleRepository(c.DB)
	userRepo := repository.NewUserRepository(c.DB)

	// Apply caching decorator and initialize token repo if Redis is available
	var tokenRepo output.TokenRepository
	if c.Redis != nil {
		roleRepo = repository.NewRoleRepositoryCached(roleRepo, c.Redis)
		userRepo = repository.NewUserRepositoryCached(userRepo, c.Redis)
		tokenRepo = redis.NewTokenRepository(c.Redis)
	}

	c.repositories = &app.Repositories{
		User:  userRepo,
		Role:  roleRepo,
		Token: tokenRepo,
	}
}

// Application returns a fully configured Application instance
func (c *Container) Application() *app.Application {
	return app.New(
		c.Config,
		c.Log,
		auth.NewAuthUseCase(c.repositories.User, c.repositories.Token, auth.Config{
			AccessPrivateKey:  c.Config.AccessPrivateKey,
			AccessExpiration:  c.Config.AccessExpiration,
			RefreshPrivateKey: c.Config.RefreshPrivateKey,
			RefreshExpiration: c.Config.RefreshExpiration,
		}),
		role.NewRoleUseCase(c.repositories.Role),
		user.NewUserUseCase(c.repositories.User, c.repositories.Role), // Updated: UserUseCase needs RoleRepo
		c.repositories,
	)
}
