// Package apperror provides centralized error handling for the application.
// These errors are domain-agnostic and can be mapped to different protocols
// (HTTP, gRPC, CLI) by their respective adapters.
package apperror

import (
	"errors"
	"fmt"
)

// Code represents an application error code
type Code string

const (
	// Authentication & Authorization
	CodeUnauthorized       Code = "unauthorized"
	CodeForbidden          Code = "forbidden"
	CodeInvalidCredentials Code = "incorrectCredentials"
	CodeDisabledUser       Code = "disabledUser"
	CodeTokenExpired       Code = "TOKEN_EXPIRED"

	// Resource errors
	CodeNotFound      Code = "NOT_FOUND"
	CodeAlreadyExists Code = "ALREADY_EXISTS"
	CodeConflict      Code = "CONFLICT"
	CodeResourceInUse Code = "RESOURCE_IN_USE"

	// Validation errors
	CodeInvalidInput     Code = "INVALID_INPUT"
	CodeValidationFailed Code = "VALIDATION_FAILED"

	// System errors
	CodeInternal        Code = "INTERNAL_ERROR"
	CodeDatabaseError   Code = "DATABASE_ERROR"
	CodeExternalService Code = "EXTERNAL_SERVICE_ERROR"
)

// Error represents an application error with rich context
type Error struct {
	// Code is the machine-readable error code
	Code Code

	// Message is the human-readable error message
	Message string

	// Field is the field name that caused the error (for validation)
	Field string

	// Details contains additional error context
	Details map[string]any

	// Cause is the underlying error
	Cause error
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Field, e.Message)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches a target
func (e *Error) Is(target error) bool {
	var appErr *Error
	if errors.As(target, &appErr) {
		return e.Code == appErr.Code
	}
	return false
}

// WithField adds a field name to the error
func (e *Error) WithField(field string) *Error {
	e.Field = field
	return e
}

// WithDetails adds details to the error
func (e *Error) WithDetails(key string, value any) *Error {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[key] = value
	return e
}

// WithCause wraps an underlying error
func (e *Error) WithCause(err error) *Error {
	e.Cause = err
	return e
}

// Constructor functions for common errors

// New creates a new application error
func New(code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an error with application context
func Wrap(code Code, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// NotFound creates a not found error
func NotFound(resource string) *Error {
	return &Error{
		Code:    CodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

// AlreadyExists creates an already exists error
func AlreadyExists(resource string) *Error {
	return &Error{
		Code:    CodeAlreadyExists,
		Message: fmt.Sprintf("%s already exists", resource),
	}
}

// InvalidInput creates an invalid input error
func InvalidInput(field, message string) *Error {
	return &Error{
		Code:    CodeInvalidInput,
		Message: message,
		Field:   field,
	}
}

// ValidationFailed creates a validation failed error
func ValidationFailed(message string) *Error {
	return &Error{
		Code:    CodeValidationFailed,
		Message: message,
	}
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *Error {
	return &Error{
		Code:    CodeUnauthorized,
		Message: message,
	}
}

// Forbidden creates a forbidden error
func Forbidden(message string) *Error {
	return &Error{
		Code:    CodeForbidden,
		Message: message,
	}
}

// Internal creates an internal error
func Internal(message string, cause error) *Error {
	return &Error{
		Code:    CodeInternal,
		Message: message,
		Cause:   cause,
	}
}

// Conflict creates a conflict error
func Conflict(resource, reason string) *Error {
	return &Error{
		Code:    CodeConflict,
		Message: fmt.Sprintf("%s conflict: %s", resource, reason),
	}
}

// ResourceInUse creates a resource in use error
func ResourceInUse(resource string) *Error {
	return &Error{
		Code:    CodeResourceInUse,
		Message: fmt.Sprintf("%s is in use and cannot be deleted", resource),
	}
}

// Helper functions

// IsCode checks if an error has a specific code
func IsCode(err error, code Code) bool {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	return IsCode(err, CodeNotFound)
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	return IsCode(err, CodeInvalidInput) || IsCode(err, CodeValidationFailed)
}

// GetCode extracts the error code from an error
func GetCode(err error) Code {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return CodeInternal
}

// Domain-specific error codes
const (
	// User errors
	CodeUserNotFound    Code = "userNotFound"
	CodeUserHasPassword Code = "userHasPassword"

	// Role errors
	CodeRoleNotFound Code = "ROLE_NOT_FOUND"

	// Password errors
	CodePasswordMismatch Code = "PASSWORD_MISMATCH"
)

// Domain-specific error constructors

// UserNotFound creates a user not found error
func UserNotFound() *Error {
	return &Error{
		Code:    CodeUserNotFound,
		Message: "user not found",
	}
}

// RoleNotFound creates a role not found error
func RoleNotFound() *Error {
	return &Error{
		Code:    CodeRoleNotFound,
		Message: "role not found",
	}
}

// UserHasPassword creates an error when user already has a password
func UserHasPassword() *Error {
	return &Error{
		Code:    CodeUserHasPassword,
		Message: "user already has a password",
	}
}

// PasswordMismatch creates a password mismatch error
func PasswordMismatch() *Error {
	return &Error{
		Code:    CodePasswordMismatch,
		Message: "passwords do not match",
	}
}

// DisabledUser creates a disabled user error
func DisabledUser() *Error {
	return &Error{
		Code:    CodeDisabledUser,
		Message: "user is disabled",
	}
}

// InvalidCredentials creates an invalid credentials error
func InvalidCredentials() *Error {
	return &Error{
		Code:    CodeInvalidCredentials,
		Message: "invalid credentials",
	}
}
