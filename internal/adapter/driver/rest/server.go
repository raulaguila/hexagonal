package rest

import (
	"crypto/rsa"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"golang.org/x/text/language"

	"github.com/raulaguila/go-api/config"
	"github.com/raulaguila/go-api/internal/adapter/driver/rest/handler"
	"github.com/raulaguila/go-api/internal/adapter/driver/rest/middleware"
	"github.com/raulaguila/go-api/internal/adapter/driver/rest/presenter"
	"github.com/raulaguila/go-api/internal/app"
	"github.com/raulaguila/go-api/pkg/loggerx"
)

// Config holds server configuration
type Config struct {
	Port              int
	EnablePrefork     bool
	EnableLogger      bool
	EnableSwagger     bool
	Version           string
	AccessPrivateKey  *rsa.PrivateKey
	RefreshPrivateKey *rsa.PrivateKey
	LocalesFS         interface {
		Open(name string) (interface{ Close() error }, error)
	}
}

// Server represents the REST API server
type Server struct {
	app    *fiber.App
	config Config
	log    *Logger
	appCtx *app.Application
}

// Logger is an alias for the logger package
type Logger = loggerx.Logger

// NewServer creates a new REST API server
func NewServer(
	config Config,
	application *app.Application,
	log *Logger,
) *Server {
	return &Server{
		config: config,
		appCtx: application,
		log:    log,
	}
}

// Start starts the REST API server
func (s *Server) Start() error {
	s.app = fiber.New(fiber.Config{
		EnablePrintRoutes:     false,
		Prefork:               s.config.EnablePrefork,
		CaseSensitive:         true,
		StrictRouting:         true,
		DisableStartupMessage: false,
		AppName:               "API Backend",
		ReduceMemoryUsage:     false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return presenter.InternalServerError(c, err.Error())
		},
		BodyLimit:      4 * 1024 * 1024,
		ReadBufferSize: 1024 * 1024,
	})

	s.setupMiddlewares()
	s.setupRoutes()

	return s.app.Listen(fmt.Sprintf(":%d", s.config.Port))
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

// setupMiddlewares configures global middlewares
func (s *Server) setupMiddlewares() {
	s.app.Use(recover.New())
	s.app.Use(otelfiber.Middleware())

	// Use structured logger with TraceID support
	if s.config.EnableLogger {
		s.app.Use(middleware.RequestLogger(s.log))
	}

	s.app.Use(
		cors.New(cors.Config{
			AllowOrigins:  "*",
			AllowMethods:  strings.Join([]string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete, fiber.MethodOptions}, ","),
			AllowHeaders:  "*",
			ExposeHeaders: "*",
			MaxAge:        -1,
		}),
		fiberi18n.New(&fiberi18n.Config{
			Next: func(c *fiber.Ctx) bool {
				return false
			},
			RootPath:        "./locales",
			AcceptLanguages: []language.Tag{language.AmericanEnglish, language.BrazilianPortuguese},
			DefaultLanguage: language.AmericanEnglish,
			Loader:          &fiberi18n.EmbedLoader{FS: config.Locales},
		}),
		limiter.New(limiter.Config{
			Next:       nil,
			Max:        300,
			Expiration: 30 * time.Second,
			LimitReached: func(c *fiber.Ctx) error {
				return presenter.New(c, fiber.StatusTooManyRequests, fiberi18n.MustLocalize(c, "manyRequests"), nil)
			},
		}),
	)
}

// setupRoutes configures API routes
func (s *Server) setupRoutes() {
	// Swagger
	if s.config.EnableSwagger {
		s.app.Get("/swagger/*", swagger.New(swagger.Config{
			DisplayRequestDuration: true,
			DocExpansion:           "none",
			ValidatorUrl:           "none",
			SyntaxHighlight: &swagger.SyntaxHighlightConfig{
				Activate: true,
				Theme:    "arta",
			},
			CustomStyle: template.CSS(""),
		}))
	}

	// Auth middlewares
	accessAuth := middleware.Auth(middleware.AuthConfig{
		PrivateKey:    s.config.AccessPrivateKey,
		UserRepo:      s.appCtx.Repositories.User,
		AllowSkipAuth: s.appCtx.Config.Environment == "development",
		Log:           s.appCtx.Log,
	})

	refreshAuth := middleware.Auth(middleware.AuthConfig{
		PrivateKey:    s.config.RefreshPrivateKey,
		UserRepo:      s.appCtx.Repositories.User,
		AllowSkipAuth: s.appCtx.Config.Environment == "development",
		Log:           s.appCtx.Log,
	})

	// Register handlers
	handler.NewHealthHandler(s.app.Group(""), s.appCtx)
	handler.NewAuthHandler(s.app.Group("/auth"), s.appCtx.Auth, accessAuth, refreshAuth)
	handler.NewProfileHandler(s.app.Group("/profile"), s.appCtx.Profile, accessAuth)
	handler.NewUserHandler(s.app.Group("/user"), s.appCtx.User, accessAuth)

	// 404 handler
	s.app.All("*", func(c *fiber.Ctx) error {
		return presenter.NotFound(c, fiberi18n.MustLocalize(c, "nonExistentRoute"))
	})
}

// Port returns the server port
func (s *Server) Port() int {
	return s.config.Port
}

// GetEnvBool returns a boolean from environment variable
func GetEnvBool(key string) bool {
	return os.Getenv(key) == "1"
}

// GetEnvString returns a string from environment variable with default
func GetEnvString(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
