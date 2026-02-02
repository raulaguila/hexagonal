package middleware

import (
	"errors"
	"log"
	"reflect"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/raulaguila/go-api/internal/adapter/driver/rest/presenter"
	"github.com/raulaguila/go-api/pkg/apperror"
	"github.com/raulaguila/go-api/pkg/pgerror"
)

// ErrorMapping maps errors to HTTP responses
type ErrorMapping map[string]map[error][]any

// NewErrorHandler creates an error handler with error mapping
func NewErrorHandler(possibleErrors ErrorMapping) func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		// First, check if it's an AppError and handle it directly
		var appErr *apperror.Error
		if errors.As(err, &appErr) {
			return handleAppError(c, appErr)
		}

		// Legacy error handling with mapping
		for method, mapper := range possibleErrors {
			if method == c.Method() || method == "*" {
				for key, value := range mapper {
					pgErr := pgerror.HandlerError(err)
					if errors.Is(pgErr, key) {
						return presenter.New(c, value[0].(int), fiberi18n.MustLocalize(c, value[1].(string)), nil)
					}
				}
			}
		}

		log.Printf("Undetected error '%v': %s\n", reflect.TypeOf(err), err.Error())
		return presenter.InternalServerError(c, fiberi18n.MustLocalize(c, "errGeneric"))
	}
}

// DefaultErrorHandler returns a simple error handler that handles apperror.Error
func DefaultErrorHandler() func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		var appErr *apperror.Error
		if errors.As(err, &appErr) {
			return handleAppError(c, appErr)
		}

		log.Printf("Undetected error '%v': %s\n", reflect.TypeOf(err), err.Error())
		return presenter.InternalServerError(c, fiberi18n.MustLocalize(c, "errGeneric"))
	}
}

// handleAppError converts AppError to HTTP response
func handleAppError(c *fiber.Ctx, err *apperror.Error) error {
	status := mapAppErrorToStatus(err.Code)
	message := err.Message

	// Try to localize if message is a key
	if localized := fiberi18n.MustLocalize(c, string(err.Code)); localized != string(err.Code) {
		message = localized
	}

	return presenter.New(c, status, message, nil)
}

// mapAppErrorToStatus maps apperror codes to HTTP status codes
func mapAppErrorToStatus(code apperror.Code) int {
	switch code {
	// Auth errors
	case apperror.CodeUnauthorized, apperror.CodeInvalidCredentials, apperror.CodeDisabledUser, apperror.CodeTokenExpired:
		return fiber.StatusUnauthorized
	case apperror.CodeForbidden:
		return fiber.StatusForbidden

	// Resource errors
	case apperror.CodeNotFound, apperror.CodeUserNotFound, apperror.CodeRoleNotFound:
		return fiber.StatusNotFound
	case apperror.CodeAlreadyExists, apperror.CodeConflict:
		return fiber.StatusConflict
	case apperror.CodeResourceInUse:
		return fiber.StatusBadRequest

	// Validation errors
	case apperror.CodeInvalidInput, apperror.CodeValidationFailed, apperror.CodeUserHasPassword, apperror.CodePasswordMismatch:
		return fiber.StatusBadRequest

	// System errors
	case apperror.CodeInternal, apperror.CodeDatabaseError, apperror.CodeExternalService:
		return fiber.StatusInternalServerError

	default:
		return fiber.StatusInternalServerError
	}
}
