// Package dto contains Data Transfer Objects for the application.
// DTOs are used to transfer data between layers and provide a clean
// interface for input/output operations.
//
// This file contains mapper functions that convert between domain entities
// and DTOs. Mappers are centralized here to:
//   - Avoid duplication across use cases
//   - Ensure consistent transformation logic
//   - Make it easy to update mappings in one place
//
// Usage:
//
//	import "github.com/raulaguila/go-api/internal/core/dto"
//
//	// Convert a single entity
//	output := dto.EntityToUserOutput(user)
//
//	// Convert a slice of entities
//	outputs := dto.EntitiesToUserOutputs(users)
package dto

import "github.com/raulaguila/go-api/internal/core/domain/entity"

// EntityToUserOutput converts a User entity to UserOutput DTO.
// Returns nil if the input user is nil.
//
// The mapping includes:
//   - Basic user info (ID, Name, Username, Email)
//   - Status from Auth
//   - Profile info if available
//   - IsNew flag (true if user has no password set)
//
// Example:
//
//	user := &entity.User{ID: 1, Name: "John"}
//	output := dto.EntityToUserOutput(user)
//	// output.ID = &1, output.Name = &"John"
func EntityToUserOutput(user *entity.User) *UserOutput {
	if user == nil {
		return nil
	}

	isNew := user.IsNew()
	output := &UserOutput{
		ID:       &user.ID,
		Name:     &user.Name,
		Username: &user.Username,
		Email:    &user.Email,
		New:      &isNew,
	}

	if user.Auth != nil {
		output.Status = &user.Auth.Status
		if user.Auth.Profile != nil {
			output.Profile = EntityToProfileOutput(user.Auth.Profile, true)
		}
	}

	return output
}

// EntityToProfileOutput converts a Profile entity to ProfileOutput DTO.
// Returns nil if the input profile is nil.
//
// Parameters:
//   - profile: The entity to convert
//   - includePermissions: If true, includes the permissions array in output
//
// Example:
//
//	profile := &entity.Profile{ID: 1, Name: "Admin"}
//	output := dto.EntityToProfileOutput(profile, true)
func EntityToProfileOutput(profile *entity.Profile, includePermissions bool) *ProfileOutput {
	if profile == nil {
		return nil
	}

	output := &ProfileOutput{
		ID:   &profile.ID,
		Name: &profile.Name,
	}

	if includePermissions {
		perms := profile.Permissions
		output.Permissions = &perms
	}

	return output
}

// EntitiesToUserOutputs converts a slice of User entities to UserOutput DTOs.
// This function is optimized for use with PaginatedOutput which requires []UserOutput.
//
// Note: Returns []UserOutput (values) instead of []*UserOutput (pointers) because
// PaginatedOutput[T paginableOutput] uses []T and paginableOutput constraint
// specifies value types (ProfileOutput | UserOutput | ItemOutput).
//
// Example:
//
//	users := []*entity.User{{ID: 1}, {ID: 2}}
//	outputs := dto.EntitiesToUserOutputs(users)
//	// len(outputs) == 2
func EntitiesToUserOutputs(users []*entity.User) []UserOutput {
	outputs := make([]UserOutput, len(users))
	for i, user := range users {
		if out := EntityToUserOutput(user); out != nil {
			outputs[i] = *out
		}
	}
	return outputs
}

// EntitiesToProfileOutputs converts a slice of Profile entities to ProfileOutput DTOs.
// See EntitiesToUserOutputs for notes on why this returns []ProfileOutput.
//
// Parameters:
//   - profiles: The entities to convert
//   - includePermissions: If true, includes permissions in each output
//
// Example:
//
//	profiles := []*entity.Profile{{ID: 1}, {ID: 2}}
//	outputs := dto.EntitiesToProfileOutputs(profiles, true)
func EntitiesToProfileOutputs(profiles []*entity.Profile, includePermissions bool) []ProfileOutput {
	outputs := make([]ProfileOutput, len(profiles))
	for i, profile := range profiles {
		if out := EntityToProfileOutput(profile, includePermissions); out != nil {
			outputs[i] = *out
		}
	}
	return outputs
}
