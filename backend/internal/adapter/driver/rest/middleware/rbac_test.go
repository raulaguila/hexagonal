package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/raulaguila/go-api/config"
	"github.com/raulaguila/go-api/internal/adapter/driver/rest/middleware"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestRequirePermission(t *testing.T) {
	app := fiber.New()

	// Setup i18n for testing
	app.Use(fiberi18n.New(&fiberi18n.Config{
		Loader: &fiberi18n.EmbedLoader{
			FS: config.Locales,
		},
		RootPath:        "locales",
		Next:            func(c *fiber.Ctx) bool { return false },
		AcceptLanguages: []language.Tag{language.English},
		DefaultLanguage: language.English,
	}))

	// Mock Auth Middleware to set user in context
	mockAuth := func(user *entity.User) fiber.Handler {
		return func(c *fiber.Ctx) error {
			if user != nil {
				c.Locals(middleware.LocalUser, user)
			}
			return c.Next()
		}
	}

	app.Get("/protected/view", mockAuth(nil), middleware.RequirePermission("users:view"), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	t.Run("Unauthorized when no user", func(t *testing.T) {
		app.Get("/test/no-user", mockAuth(nil), middleware.RequirePermission("users:view"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test/no-user", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Forbidden when user missing permission", func(t *testing.T) {
		user := &entity.User{
			ID: uuid.New(),
			Roles: []*entity.Role{
				{Name: "GUEST", Permissions: []string{"other:view"}},
			},
		}

		app.Get("/test/forbidden", mockAuth(user), middleware.RequirePermission("users:view"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test/forbidden", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})

	t.Run("Allowed when user has permission", func(t *testing.T) {
		user := &entity.User{
			ID: uuid.New(),
			Roles: []*entity.Role{
				{Name: "EDITOR", Permissions: []string{"users:view"}},
			},
		}

		app.Get("/test/allowed", mockAuth(user), middleware.RequirePermission("users:view"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test/allowed", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Allowed when user has * permission", func(t *testing.T) {
		user := &entity.User{
			ID: uuid.New(),
			Roles: []*entity.Role{
				{Name: "ADMIN", Permissions: []string{"*"}},
			},
		}

		app.Get("/test/wildcard", mockAuth(user), middleware.RequirePermission("users:delete"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test/wildcard", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Allowed when user is ROOT", func(t *testing.T) {
		user := &entity.User{
			ID: uuid.New(),
			Roles: []*entity.Role{
				{Name: "ROOT", Permissions: []string{}},
			},
		}

		app.Get("/test/root", mockAuth(user), middleware.RequirePermission("users:delete"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test/root", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}
