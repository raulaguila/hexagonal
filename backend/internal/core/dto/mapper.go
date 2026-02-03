package dto

import "github.com/raulaguila/go-api/internal/core/domain/entity"

// EntityToUserOutput converts a User entity to UserOutput DTO.
// Returns nil if the input user is nil.
func EntityToUserOutput(user *entity.User, includePermissions bool) *UserOutput {
	if user == nil {
		return nil
	}

	idStr := user.ID.String()
	isNew := user.IsNew()
	output := &UserOutput{
		ID:       &idStr,
		Name:     &user.Name,
		Username: &user.Username,
		Email:    &user.Email,
		New:      &isNew,
	}

	if user.Auth != nil {
		output.Status = &user.Auth.Status
	}

	// Map Roles
	if len(user.Roles) > 0 {
		roleOutputs := make([]*RoleOutput, len(user.Roles))
		for i, r := range user.Roles {
			roleOutputs[i] = EntityToRoleOutput(r, includePermissions) // Always include permissions in user details? Or maybe not?
			// Usually valid to include basic info.
		}
		output.Roles = roleOutputs
	}

	return output
}

// EntityToRoleOutput converts a Role entity to RoleOutput DTO.
// Returns nil if the input role is nil.
func EntityToRoleOutput(role *entity.Role, includePermissions bool) *RoleOutput {
	if role == nil {
		return nil
	}

	idStr := role.ID.String()
	output := &RoleOutput{
		ID:      &idStr,
		Name:    &role.Name,
		Enabled: &role.Enabled,
	}

	if includePermissions {
		perms := role.Permissions
		output.Permissions = &perms
	}

	return output
}

// EntitiesToUserOutputs converts a slice of User entities to UserOutput DTOs.
func EntitiesToUserOutputs(users []*entity.User) []UserOutput {
	outputs := make([]UserOutput, len(users))
	for i, user := range users {
		if out := EntityToUserOutput(user, false); out != nil {
			outputs[i] = *out
		}
	}
	return outputs
}

// EntitiesToRoleOutputs converts a slice of Role entities to RoleOutput DTOs.
func EntitiesToRoleOutputs(roles []*entity.Role, includePermissions bool) []RoleOutput {
	outputs := make([]RoleOutput, len(roles))
	for i, role := range roles {
		if out := EntityToRoleOutput(role, includePermissions); out != nil {
			outputs[i] = *out
		}
	}
	return outputs
}
