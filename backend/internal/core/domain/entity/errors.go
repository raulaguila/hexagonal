package entity

import "github.com/raulaguila/go-api/pkg/apperror"

// Entity validation errors - using apperror for consistent error handling

// ErrRoleNameTooShort returns error for short role name
func ErrRoleNameTooShort() *apperror.Error {
	return apperror.InvalidInput("name", "role name must be at least 4 characters")
}

// ErrRoleRequired returns error when role is missing
func ErrRoleRequired() *apperror.Error {
	return apperror.InvalidInput("role", "at least one role is required")
}

// ErrUserNameTooShort returns error for short user name
func ErrUserNameTooShort() *apperror.Error {
	return apperror.InvalidInput("name", "user name must be at least 5 characters")
}

// ErrUsernameTooShort returns error for short username
func ErrUsernameTooShort() *apperror.Error {
	return apperror.InvalidInput("username", "username must be at least 5 characters")
}

// ErrEmailRequired returns error when email is missing
func ErrEmailRequired() *apperror.Error {
	return apperror.InvalidInput("email", "email is required")
}

// ErrInvalidEmailFormat returns error for invalid email format
func ErrInvalidEmailFormat() *apperror.Error {
	return apperror.InvalidInput("email", "invalid email format")
}

// ErrPasswordTooShort returns error for short password
func ErrPasswordTooShort() *apperror.Error {
	return apperror.InvalidInput("password", "password must be at least 6 characters")
}

// ErrAuthRequired returns error when auth is missing
func ErrAuthRequired() *apperror.Error {
	return apperror.InvalidInput("auth", "auth is required")
}
