package entity

import (
	"slices"
	"time"

	"github.com/raulaguila/go-api/pkg/validator"
)

// Profile represents a user profile with permissions in the domain
type Profile struct {
	ID          uint
	Name        string
	Permissions []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewProfile creates a new Profile entity
func NewProfile(name string, permissions []string) *Profile {
	now := time.Now()
	return &Profile{
		Name:        name,
		Permissions: permissions,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Validate validates the profile entity
func (p *Profile) Validate() error {
	if len(p.Name) < validator.MinProfileNameLength {
		return ErrProfileNameTooShort()
	}
	return nil
}

// UpdateName updates the profile name
func (p *Profile) UpdateName(name string) {
	p.Name = name
	p.UpdatedAt = time.Now()
}

// UpdatePermissions updates the profile permissions
func (p *Profile) UpdatePermissions(permissions []string) {
	p.Permissions = permissions
	p.UpdatedAt = time.Now()
}

// HasPermission checks if profile has a specific permission
func (p *Profile) HasPermission(permission string) bool {
	return slices.Contains(p.Permissions, permission)
}

// IsRoot checks if this is the root profile (ID = 1)
func (p *Profile) IsRoot() bool {
	return p.ID == 1
}
