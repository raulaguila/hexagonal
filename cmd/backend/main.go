package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	_ "github.com/raulaguila/go-api/docs" // Swagger docs

	"github.com/raulaguila/go-api/config"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres"
	"github.com/raulaguila/go-api/internal/adapter/driven/storage/minio"
	"github.com/raulaguila/go-api/internal/adapter/driven/storage/redis"
	"github.com/raulaguila/go-api/internal/adapter/driver/rest"
	"github.com/raulaguila/go-api/internal/di"
	"github.com/raulaguila/go-api/pkg/loggerx"
	"github.com/raulaguila/go-api/pkg/loggerx/formatter"
	"github.com/raulaguila/go-api/pkg/loggerx/sink"
	"github.com/raulaguila/go-api/pkg/telemetry"
)

// initLogger creates and configures the application logger
func initLogger(cfg *config.Environment) *loggerx.Logger {
	return loggerx.New(
		// loggerx.WithTimeFormat(time.DateTime),
		loggerx.WithLevel(loggerx.ParseLogLevel(cfg.LogLevel)),
		loggerx.WithSink(sink.NewStdout(
			sink.WithFormatter(formatter.NewJSON()),
		)),
		loggerx.WithFields(map[string]any{
			"app":     cfg.ServiceName,
			"version": cfg.Version,
		}),
		loggerx.WithSink(sink.NewOTLP()),
	)
}

// @title 						Go API
// @version						1.0.0
// @description 				This API is a user-friendly solution designed to serve as the foundation for more complex APIs.

// @contact.name				Raul del Aguila
// @contact.url 				https://github.com/raulaguila
// @contact.email				email@email.com

// @BasePath					/

// @securityDefinitions.apiKey	Bearer
// @in							header
// @name						Authorization
// @description 				Type "Bearer" followed by a space and the JWT token.
func main() {
	// Load configuration
	cfg := config.MustLoad()

	// Initialize OpenTelemetry
	ctx := context.Background()
	shutdown, err := telemetry.SetupOTelSDK(ctx, cfg.ServiceName, cfg.Version, cfg.OtelExporterOtlpEndpoint)
	if err != nil {
		fmt.Printf("Error setting up OpenTelemetry: %v\n", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down OpenTelemetry: %v\n", err)
		}
	}()

	// Initialize logger
	log := initLogger(cfg)

	log.Info("Configuration loaded", slog.Int("port", cfg.Port), slog.String("environment", cfg.Environment))

	// Connect to PostgreSQL
	log.Info("Connecting to PostgreSQL...")
	db := postgres.MustConnect(&postgres.Config{Dsn: cfg.PGDSN})
	log.Info("Database connected", slog.String("host", cfg.PGHost), slog.String("database", cfg.PGBase))

	// Connect to MinIO
	log.Info("Connecting to MinIO...")
	_ = minio.MustConnect(&minio.Config{Url: cfg.MinioUrl, User: cfg.MinioUser, Password: cfg.MinioPassword, BucketName: cfg.MinioBucketName})
	log.Info("Storage connected", slog.String("host", cfg.MinioHost), slog.String("bucket", cfg.MinioBucketName))

	// Connect to Redis
	log.Info("Connecting to Redis...")
	redisSvc := redis.MustNew(redis.Config{Host: cfg.RedisHost, Port: cfg.RedisPort, Password: cfg.RedisPass, DB: cfg.RedisDB})
	log.Info("Redis connected", slog.String("host", cfg.RedisHost))

	// Initialize dependency container
	log.Info("Initializing dependencies...")
	container := di.NewContainer(cfg, log, db, redisSvc)

	// Get application instance
	application := container.Application()

	// Create REST server using the Application
	server := rest.NewServer(
		rest.Config{
			Port:              cfg.Port,
			EnablePrefork:     cfg.EnablePrefork,
			EnableLogger:      cfg.EnableLogger,
			EnableSwagger:     cfg.EnableSwagger,
			Version:           cfg.Version,
			AccessPrivateKey:  cfg.AccessPrivateKey,
			RefreshPrivateKey: cfg.RefreshPrivateKey,
		},
		application,
		log,
	)

	// Handle graceful shutdown
	go handleShutdown(log, server, container)

	log.Info("Application starting",
		slog.Int("port", cfg.Port),
		slog.String("go_version", runtime.Version()),
		slog.Int("num_cpu", runtime.NumCPU()),
		slog.Int("gomaxprocs", runtime.GOMAXPROCS(0)),
	)

	// Start server
	if err := server.Start(); err != nil {
		log.Error("Server failed to start",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}

// handleShutdown handles graceful shutdown on SIGINT/SIGTERM
func handleShutdown(log *loggerx.Logger, server *rest.Server, container *di.Container) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Shutdown(sig.String())

	// Shutdown server
	_ = server.Shutdown()

	// Close database connection
	if container.DB != nil {
		if sqlDB, err := container.DB.DB(); err == nil {
			_ = sqlDB.Close()
			log.Info("Database connection closed")
		}
	}

	os.Exit(0)
}
