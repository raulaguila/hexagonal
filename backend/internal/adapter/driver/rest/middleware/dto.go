package middleware

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

// Lookup determines where to look for the DTO
type Lookup uint8

const (
	Body Lookup = iota
	Query
	Params
	Cookie
)

// DTOConfig holds configuration for DTO parsing middleware
type DTOConfig struct {
	ContextKey   any
	OnLookup     Lookup
	Model        any
	ErrorHandler fiber.ErrorHandler
}

// DefaultDTOConfig provides default configuration
var DefaultDTOConfig = DTOConfig{
	ContextKey: CtxKeyDTO,
	OnLookup:   Body,
	Model:      new(map[string]any),
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": err.Error(),
		})
	},
}

// ParseDTO creates a middleware for parsing DTOs
func ParseDTO(config ...DTOConfig) fiber.Handler {
	cfg := DefaultDTOConfig
	if len(config) > 0 {
		cfg = config[0]
		if cfg.ContextKey == nil {
			cfg.ContextKey = DefaultDTOConfig.ContextKey
		}
		if cfg.Model == nil {
			cfg.Model = DefaultDTOConfig.Model
		}
		if cfg.ErrorHandler == nil {
			cfg.ErrorHandler = DefaultDTOConfig.ErrorHandler
		}
	}

	parser := func(c *fiber.Ctx, obj any) (any, error) {
		switch cfg.OnLookup {
		case Body:
			return obj, c.BodyParser(obj)
		case Query:
			return obj, c.QueryParser(obj)
		case Params:
			return obj, c.ParamsParser(obj)
		default:
			return obj, c.CookieParser(obj)
		}
	}

	return func(c *fiber.Ctx) error {
		obj := reflect.New(reflect.TypeOf(cfg.Model).Elem()).Interface()
		obj, err := parser(c, obj)
		if err != nil {
			fmt.Printf("Error parsing DTO: %v - %v\n", err, reflect.TypeOf(cfg.Model))
			return cfg.ErrorHandler(c, err)
		}

		c.Locals(cfg.ContextKey, obj)
		return c.Next()
	}
}
