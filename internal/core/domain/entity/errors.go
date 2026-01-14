package entity

import "github.com/raulaguila/go-api/pkg/apperror"

// Entity validation errors - using apperror for consistent error handling

// ErrProfileNameTooShort returns error for short profile name
func ErrProfileNameTooShort() *apperror.Error {
	return apperror.InvalidInput("name", "profile name must be at least 4 characters")
}

// ErrUserNameTooShort returns error for short user name
func ErrUserNameTooShort() *apperror.Error {
	return apperror.InvalidInput("name", "name must be at least 5 characters")
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

// ErrProfileRequired returns error when profile is missing
func ErrProfileRequired() *apperror.Error {
	return apperror.InvalidInput("profile_id", "profile is required")
}

// ErrPasswordTooShort returns error for short password
func ErrPasswordTooShort() *apperror.Error {
	return apperror.InvalidInput("password", "password must be at least 6 characters")
}
