// Package utils provides common utility functions used throughout the application.
package utils

// Deref returns the value of a pointer or a default value if the pointer is nil.
// This is useful for safely accessing optional pointer values in DTOs.
//
// Example:
//
//	var name *string = nil
//	result := utils.Deref(name, "default")  // returns "default"
func Deref[T any](ptr *T, defaultVal T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultVal
}
