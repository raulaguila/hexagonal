// Package validator provides centralized validation functions and constants
// for the entire application. This package ensures consistent validation
// rules across all layers (entities, DTOs, handlers).
//
// Usage:
//
//	import "github.com/raulaguila/go-api/pkg/validator"
//
//	// Validate email
//	if !validator.IsValidEmail(email) {
//	    return errors.New("invalid email")
//	}
//
//	// Use constants for validation
//	if len(name) < validator.MinNameLength {
//	    return errors.New("name too short")
//	}
package validator

import "regexp"

// Validation constants define the minimum and maximum lengths for various fields.
// These values are used consistently across entities and DTOs to ensure
// validation rules are centralized and easy to maintain.
const (
	// MinNameLength is the minimum length for a user's display name.
	// This ensures names are meaningful and identifiable.
	MinNameLength = 5

	// MaxNameLength is the maximum length for a user's display name.
	// This prevents excessively long names that could affect UI/UX.
	MaxNameLength = 100

	// MinUsernameLength is the minimum length for a username (login identifier).
	// Usernames need to be long enough to be unique but memorable.
	MinUsernameLength = 5

	// MaxUsernameLength is the maximum length for a username.
	// This aligns with common platform standards and database constraints.
	MaxUsernameLength = 50

	// MinProfileNameLength is the minimum length for a profile/role name.
	// Profile names like "Admin" or "User" need at least 4 characters.
	MinProfileNameLength = 4

	// MaxProfileNameLength is the maximum length for a profile name.
	MaxProfileNameLength = 100

	// MinPasswordLength is the minimum length for a password.
	// 6 characters is the minimum for basic security; consider increasing for production.
	MinPasswordLength = 6

	// MaxPasswordLength is the maximum length for a password.
	// 128 characters allows for passphrases while preventing DoS via bcrypt.
	MaxPasswordLength = 128

	// MinEmailLength is the minimum length for a valid email address.
	// The shortest valid email is "a@b.c" (5 chars), but we use 3 as a pre-check.
	MinEmailLength = 3
)

// EmailRegex is the compiled regular expression for email validation.
// It matches standard email formats: user@domain.tld
//
// Pattern breakdown:
//   - ^[a-zA-Z0-9._%+-]+ : Local part (before @)
//   - @                  : Required @ symbol
//   - [a-zA-Z0-9.-]+     : Domain name
//   - \.[a-zA-Z]{2,}$    : TLD (at least 2 characters)
var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail checks if the provided string is a valid email address.
// It performs a length check first (for performance) then validates against
// the EmailRegex pattern.
//
// Example:
//
//	validator.IsValidEmail("user@example.com")  // returns true
//	validator.IsValidEmail("invalid")           // returns false
//	validator.IsValidEmail("")                  // returns false
func IsValidEmail(email string) bool {
	if len(email) < MinEmailLength {
		return false
	}
	return EmailRegex.MatchString(email)
}
