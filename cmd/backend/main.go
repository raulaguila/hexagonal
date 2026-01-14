package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/raulaguila/go-api/docs" // Swagger docs

	"github.com/raulaguila/go-api/config"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres"
	"github.com/raulaguila/go-api/internal/adapter/driven/storage/minio"
	"github.com/raulaguila/go-api/internal/adapter/driver/rest"
	"github.com/raulaguila/go-api/internal/di"
	"github.com/raulaguila/go-api/pkg/logger"
)

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

	// Initialize logger
	log := initLogger(cfg)

	log.Info("Configuration loaded", slog.String("port", cfg.Port), slog.String("environment", cfg.Environment))

	// Connect to PostgreSQL
	log.Info("Connecting to PostgreSQL...")
	db := postgres.MustConnect(postgres.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Database: cfg.DBName,
		Timezone: cfg.Timezone,
	})
	log.Info("Database connected", slog.String("host", cfg.DBHost), slog.String("database", cfg.DBName))

	// Connect to MinIO
	log.Info("Connecting to MinIO...")
	_ = minio.MustConnect(minio.Config{
		Host:       cfg.MinioHost,
		Port:       cfg.MinioPort,
		User:       cfg.MinioUser,
		Password:   cfg.MinioPassword,
		BucketName: cfg.MinioBucketName,
	})
	log.Info("Storage connected", slog.String("host", cfg.MinioHost), slog.String("bucket", cfg.MinioBucketName))

	// Initialize dependency container
	log.Info("Initializing dependencies...")
	container := di.NewContainer(cfg, log, db)

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

	// Start server
	log.Startup(cfg.Port, cfg.Version)
	if err := server.Start(); err != nil {
		log.Error("Server failed to start",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}

// initLogger creates and configures the application logger
func initLogger(cfg *config.Config) *logger.Logger {
	return logger.Init(logger.Config{
		Level:       parseLogLevel(cfg.LogLevel),
		Format:      cfg.LogFormat,
		ServiceName: "go-api",
		Version:     cfg.Version,
		Environment: cfg.Environment,
		AddSource:   cfg.Environment != "production",
	})
}

// handleShutdown handles graceful shutdown on SIGINT/SIGTERM
func handleShutdown(log *logger.Logger, server *rest.Server, container *di.Container) {
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

// parseLogLevel converts string to logger.Level
func parseLogLevel(level string) logger.Level {
	switch level {
	case "debug":
		return logger.LevelDebug
	case "warn":
		return logger.LevelWarn
	case "error":
		return logger.LevelError
	default:
		return logger.LevelInfo
	}
}
