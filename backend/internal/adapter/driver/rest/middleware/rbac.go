package middleware

import (
	"slices"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/raulaguila/go-api/internal/adapter/driver/rest/presenter"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
)

// RequirePermission creates a middleware to check if user has a specific permission
// The User must have been set in context by the Auth middleware previously
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals(LocalUser).(*entity.User)
		if !ok || user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(presenter.Response{
				Message: fiberi18n.MustLocalize(c, "unauthorized"),
			})
		}

		// Root role usually has access to everything
		// Skip disabled roles - only enabled roles grant permissions
		for _, role := range user.Roles {
			if !role.IsEnabled() {
				continue // Skip disabled roles
			}
			if role.IsRoot() || slices.Contains(role.Permissions, "*") || role.HasPermission(permission) {
				goto next
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(presenter.Response{
			Message: fiberi18n.MustLocalize(c, "forbidden"),
		})

	next:
		return c.Next()
	}
}
