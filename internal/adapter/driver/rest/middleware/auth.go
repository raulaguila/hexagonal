package middleware

import (
	"context"
	"crypto/rsa"
	"errors"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/golang-jwt/jwt/v5"

	"github.com/raulaguila/go-api/internal/adapter/driver/rest/presenter"
	"github.com/raulaguila/go-api/internal/core/port/output"
	"github.com/raulaguila/go-api/pkg/logger"
)

const (
	// LocalUser is the context key for the authenticated user
	LocalUser = "localUser"
	// LocalUserID is the context key for the authenticated user ID
	LocalUserID = "localUserID"
)

// AuthConfig holds authentication middleware configuration
type AuthConfig struct {
	PrivateKey    *rsa.PrivateKey
	UserRepo      output.UserRepository
	AllowSkipAuth bool           // Injected config instead of os.Getenv
	Log           *logger.Logger // Injected logger instead of log.Println
}

// Auth creates an authentication middleware
func Auth(cfg AuthConfig) fiber.Handler {
	return keyauth.New(keyauth.Config{
		KeyLookup:  "header:" + fiber.HeaderAuthorization,
		AuthScheme: "Bearer",
		ContextKey: "token",
		Next: func(c *fiber.Ctx) bool {
			// Skip auth if allowed via config and header is set
			if cfg.AllowSkipAuth && c.Get("X-Skip-Auth", "false") == "true" {
				c.Locals(LocalUserID, uint(0))
				return true
			}
			return false
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return presenter.Unauthorized(c, err.Error())
		},
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			parsedToken, err := jwt.Parse(key, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, jwt.ErrTokenSignatureInvalid
				}
				return cfg.PrivateKey.Public(), nil
			})
			if err != nil {
				if cfg.Log != nil {
					cfg.Log.Debug("JWT parse error", "error", err.Error())
				}
				return false, errors.New(fiberi18n.MustLocalize(c, "errGeneric"))
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok || !parsedToken.Valid {
				return false, errors.New("invalid jwt token")
			}

			token, ok := claims["token"].(string)
			if !ok {
				return false, errors.New("invalid token claim")
			}

			// Use c.Context() instead of context.Background()
			user, err := cfg.UserRepo.FindByToken(c.Context(), token)
			if err != nil {
				if cfg.Log != nil {
					cfg.Log.Debug("User lookup error", "error", err.Error())
				}
				return false, errors.New(fiberi18n.MustLocalize(c, "errGeneric"))
			}

			if user.Auth == nil || !user.Auth.Status {
				return false, errors.New(fiberi18n.MustLocalize(c, "disabledUser"))
			}

			c.Locals(LocalUserID, user.ID)
			c.Locals(LocalUser, user)
			return true, nil
		},
	})
}

// GetUserID retrieves the user ID from context
func GetUserID(c *fiber.Ctx) uint {
	if id, ok := c.Locals(LocalUserID).(uint); ok {
		return id
	}
	return 0
}

// contextKey is a type for context keys to avoid collisions
type contextKey string

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, contextKey("user_id"), userID)
}

// UserIDFromContext extracts user ID from context
func UserIDFromContext(ctx context.Context) uint {
	if id, ok := ctx.Value(contextKey("user_id")).(uint); ok {
		return id
	}
	return 0
}
