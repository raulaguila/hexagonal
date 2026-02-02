package entity

import (
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/raulaguila/go-api/pkg/validator"
)

// Role represents a user role with permissions in the domain
type Role struct {
	ID          uuid.UUID
	Name        string
	Permissions []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewRole creates a new Role entity
func NewRole(name string, permissions []string) *Role {
	now := time.Now()
	return &Role{
		ID:          uuid.New(), // Generate new UUID
		Name:        name,
		Permissions: permissions,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Validate validates the role entity
func (r *Role) Validate() error {
	if len(r.Name) < validator.MinProfileNameLength { // Assuming keeping same constant or need to rename constant too
		return ErrRoleNameTooShort()
	}
	return nil
}

// UpdateName updates the role name
func (r *Role) UpdateName(name string, now time.Time) {
	r.Name = name
	r.UpdatedAt = now
}

// UpdatePermissions updates the role permissions
func (r *Role) UpdatePermissions(permissions []string, now time.Time) {
	r.Permissions = permissions
	r.UpdatedAt = now
}

// HasPermission checks if role has a specific permission
func (r *Role) HasPermission(permission string) bool {
	return slices.Contains(r.Permissions, permission)
}

// IsRoot checks if this is the root role
func (r *Role) IsRoot() bool {
	return r.Name == "ROOT"
}
