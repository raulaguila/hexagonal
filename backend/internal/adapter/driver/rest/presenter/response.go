package presenter

import "github.com/gofiber/fiber/v2"

// Response represents a standard HTTP response
type Response struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
	Object  any    `json:"object,omitempty"`
}

// New creates and sends a new HTTP response
func New(c *fiber.Ctx, status int, message string, object any) error {
	return c.Status(status).JSON(&Response{
		Code:    status,
		Message: message,
		Object:  object,
	})
}

// Success sends a success response with data
func Success(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(data)
}

// Created sends a created response
func Created(c *fiber.Ctx, message string, data any) error {
	return New(c, fiber.StatusCreated, message, data)
}

// BadRequest sends a bad request response
func BadRequest(c *fiber.Ctx, message string) error {
	return New(c, fiber.StatusBadRequest, message, nil)
}

// Unauthorized sends an unauthorized response
func Unauthorized(c *fiber.Ctx, message string) error {
	return New(c, fiber.StatusUnauthorized, message, nil)
}

// NotFound sends a not found response
func NotFound(c *fiber.Ctx, message string) error {
	return New(c, fiber.StatusNotFound, message, nil)
}

// Conflict sends a conflict response
func Conflict(c *fiber.Ctx, message string) error {
	return New(c, fiber.StatusConflict, message, nil)
}

// InternalServerError sends an internal server error response
func InternalServerError(c *fiber.Ctx, message string) error {
	return New(c, fiber.StatusInternalServerError, message, nil)
}
