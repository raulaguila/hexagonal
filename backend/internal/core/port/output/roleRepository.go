package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
)

// RoleRepository defines the interface for role persistence operations
type RoleRepository interface {
	// Count returns the total number of roles matching the filter
	Count(ctx context.Context, filter *dto.RoleFilter) (int64, error)

	// FindAll returns all roles matching the filter
	FindAll(ctx context.Context, filter *dto.RoleFilter) ([]*entity.Role, error)

	// FindByID returns a role by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)

	// FindByName returns a role by its name
	FindByName(ctx context.Context, name string) (*entity.Role, error)

	// Create creates a new role
	Create(ctx context.Context, role *entity.Role) error

	// Update updates an existing role
	Update(ctx context.Context, role *entity.Role) error

	// Delete deletes roles by their IDs
	Delete(ctx context.Context, ids []uuid.UUID) error
}
