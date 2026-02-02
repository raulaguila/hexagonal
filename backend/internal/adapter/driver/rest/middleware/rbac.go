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
		// We can check if implicit '*' permission or check role name
		// Implementation plan said: Check if User.Role.IsRoot()
		// But User entity has Roles slice. So we check if ANY role is root.

		isRoot := false
		hasPerm := false

		for _, role := range user.Roles {
			if role.IsRoot() || slices.Contains(role.Permissions, "*") {
				isRoot = true
				break
			}
			if role.HasPermission(permission) {
				hasPerm = true
			}
		}

		if isRoot || hasPerm {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(presenter.Response{
			Message: fiberi18n.MustLocalize(c, "forbidden"),
		})
	}
}
