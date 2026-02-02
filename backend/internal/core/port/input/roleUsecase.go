package input

import (
	"context"

	"github.com/raulaguila/go-api/internal/core/dto"
)

// RoleUseCase defines the interface for role operations
type RoleUseCase interface {
	// GetRoles returns a paginated list of roles
	GetRoles(ctx context.Context, filter *dto.RoleFilter) (*dto.PaginatedOutput[dto.RoleOutput], error)

	// ListRoles returns a simple list of roles (id + name)
	ListRoles(ctx context.Context, filter *dto.RoleFilter) ([]dto.ItemOutput, error)

	// GetRoleByID returns a role by its ID
	GetRoleByID(ctx context.Context, id string) (*dto.RoleOutput, error)

	// CreateRole creates a new role
	CreateRole(ctx context.Context, input *dto.RoleInput) (*dto.RoleOutput, error)

	// UpdateRole updates an existing role
	UpdateRole(ctx context.Context, id string, input *dto.RoleInput) (*dto.RoleOutput, error)

	// DeleteRoles deletes roles by their IDs
	DeleteRoles(ctx context.Context, ids []string) error
}
