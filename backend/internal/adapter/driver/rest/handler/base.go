package handler

import (
	"net/url"

	"github.com/gofiber/fiber/v2"

	"github.com/raulaguila/go-api/internal/adapter/driver/rest/middleware"
)

// GetLocal retrieves a value from Fiber locals safely cast to T.
func GetLocal[T any](c *fiber.Ctx, key middleware.CtxKey) *T {
	val := c.Locals(key)
	if val == nil {
		return nil
	}
	return val.(*T)
}

// GetQuery retrieves and unescapes a query parameter.
func GetQuery(c *fiber.Ctx, key string) (string, error) {
	return url.QueryUnescape(c.Query(key, ""))
}
