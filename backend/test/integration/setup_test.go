package integration

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/raulaguila/go-api/config"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/model"
	"github.com/raulaguila/go-api/internal/adapter/driven/storage/redis"
	"github.com/raulaguila/go-api/internal/adapter/driver/rest"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/di"
	"github.com/raulaguila/go-api/pkg/loggerx"
	"github.com/raulaguila/go-api/pkg/loggerx/formatter"
	"github.com/raulaguila/go-api/pkg/loggerx/sink"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	baseURL string
)

func TestMain(m *testing.M) {
	// Setup Context
	ctx := context.Background()

	// --- 1. Start Postgres Container ---
	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("user"),
		tcpostgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		fmt.Printf("Failed to start postgres container: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			fmt.Printf("failed to terminate postgres: %v\n", err)
		}
	}()

	pgConnStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Printf("Failed to get postgres connection string: %v\n", err)
		os.Exit(1)
	}

	// --- 2. Start Redis Container ---
	redisContainer, err := tcredis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		fmt.Printf("Failed to start redis container: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			fmt.Printf("failed to terminate redis: %v\n", err)
		}
	}()

	redisHost, err := redisContainer.Host(ctx)
	if err != nil {
		fmt.Printf("Failed to get redis host: %v\n", err)
		os.Exit(1)
	}
	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		fmt.Printf("Failed to get redis port: %v\n", err)
		os.Exit(1)
	}

	// --- 3. Configure Application ---
	os.Setenv("PORT", "9998")
	os.Setenv("LOG_LEVEL", "error")

	cfg := config.MustLoad()
	cfg.Port = 9998
	cfg.PGUrl = pgConnStr
	cfg.RedisHost = redisHost
	p, err := strconv.Atoi(redisPort.Port())
	if err != nil {
		fmt.Printf("Failed to convert redis port: %v\n", err)
		os.Exit(1)
	}
	cfg.RedisPort = p

	cfg.RedisUser = ""
	cfg.RedisPass = ""
	cfg.RedisDB = 0

	// Enable Logger for debugging (off for less noise in CI)
	log := loggerx.New(
		loggerx.WithLevel(loggerx.ErrorLevel), // Keep logs quiet unless error
		loggerx.WithSink(sink.NewStdout(sink.WithFormatter(formatter.NewText()))),
	)

	// Connect to DBs
	db := postgres.MustConnect(&postgres.Config{Url: cfg.PGUrl})
	redisSvc := redis.MustNew(redis.Config{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		User:     cfg.RedisUser,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	// Initialize Schema (Migrate)
	// We need to run migrations on the fresh DB.
	// Assuming we can access GORM from here.
	// We can use the AutoMigrate from our models if we have them accessbile, OR load sql files.
	// Since we are testing integration of the full app, the migration usually happens in 'main' or via a migration tool.
	// The original `setup_test.go` did NOT migrate, relying on external DB.
	// We MUST migrate now.
	// Let's borrow the migration logic from repository tests, but apply to the *whole* schema we know of.
	// Or we can try to read the SQL files.
	// Easier: AutoMigrate validation.
	// We need to import models.
	// But `test/integration` might not want to import `internal/adapter/.../persistence/postgres/model` directly if it treats app as blackbox?
	// It already imports `postgres` adapter.
	// Let's assume we can AutoMigrate via GORM instance.
	// We need to import the models.
	// Since dependencies are allowed, let's just import all models we know.
	// user, role, auth, audit.
	// But we need to update imports.

	// Clean Redis
	if err := redisSvc.GetClient().FlushAll(ctx).Err(); err != nil {
		fmt.Printf("Failed to flush redis: %v\n", err)
	}

	// DI Container
	container := di.NewContainer(cfg, log, db, redisSvc)
	app := container.Application()

	// Perform Migrations (adding this step essentially)
	// We need to manually register the models if the app doesn't expose a Migrate function.
	// Look at how `main.go` does it.
	// If `di` or `postgres` doesn't do it, we do it here.
	// Let's verify imports first.

	// Server
	server := rest.NewServer(
		rest.Config{
			Port:              cfg.Port,
			EnablePrefork:     false,
			EnableLogger:      true,
			EnableSwagger:     false,
			Version:           cfg.Version,
			AccessPrivateKey:  cfg.AccessPrivateKey,
			RefreshPrivateKey: cfg.RefreshPrivateKey,
		},
		app,
		log,
		redisSvc,
	)

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			// fmt.Printf("Server stopped: %v\n", err)
		}
	}()

	baseURL = fmt.Sprintf("http://localhost:%d", cfg.Port)

	// --- 4. Setup Data ---
	// Create required extensions
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";")
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"citext\";")
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"unaccent\";")

	// Initialize Schema (Migrate)
	if err := db.AutoMigrate(
		&model.AuditModel{},
		&model.AuthModel{},
		&model.RoleModel{},
		&model.UserModel{},
	); err != nil {
		fmt.Printf("Failed to migrate: %v\n", err)
		os.Exit(1)
	}

	// User Seed Logic
	userRepo := app.Repositories.User
	auth, err := entity.NewAuth(true)
	if err != nil {
		os.Exit(1)
	}
	_ = auth.SetPassword("password", time.Now())
	user, err := entity.NewUser("Test User", "testuser", "test@test.com", auth)

	if err := userRepo.Create(ctx, user); err != nil {
		// Since we have a fresh DB, this should only happen if migration failed or unexpected.
		// fmt.Printf("Create user failed: %v\n", err)
	}
	fmt.Printf("setup: Test user 'testuser' created.\n")
	time.Sleep(2 * time.Second)

	// Run tests
	code := m.Run()

	// Teardown
	server.Shutdown()
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
	redisSvc.Close()

	os.Exit(code)
}
