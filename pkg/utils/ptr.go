// Package utils provides common utility functions used throughout the application.
// These are generic helpers that don't belong to any specific domain.
//
// Usage:
//
//	import "github.com/raulaguila/go-api/pkg/utils"
//
//	// Safely dereference a pointer with a default value
//	name := utils.Deref(input.Name, "default")
//
//	// Create a pointer to a value
//	statusPtr := utils.Ptr(true)
package utils

// Deref returns the value of a pointer or a default value if the pointer is nil.
// This is useful for safely accessing optional pointer values in DTOs.
//
// Example:
//
//	var name *string = nil
//	result := utils.Deref(name, "default")  // returns "default"
//
//	value := "hello"
//	name = &value
//	result = utils.Deref(name, "default")   // returns "hello"
func Deref[T any](ptr *T, defaultVal T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultVal
}

// Ptr returns a pointer to the given value.
// This is useful for creating pointers to literal values,
// especially when building DTOs or test data.
//
// Example:
//
//	// Instead of:
//	name := "John"
//	input := &UserInput{Name: &name}
//
//	// You can write:
//	input := &UserInput{Name: utils.Ptr("John")}
func Ptr[T any](v T) *T {
	return &v
}
