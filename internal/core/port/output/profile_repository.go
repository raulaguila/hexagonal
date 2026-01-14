package output

import (
	"context"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
)

// ProfileRepository defines the interface for profile persistence operations
type ProfileRepository interface {
	// Count returns the total number of profiles matching the filter
	Count(ctx context.Context, filter *dto.ProfileFilter) (int64, error)

	// FindAll returns all profiles matching the filter
	FindAll(ctx context.Context, filter *dto.ProfileFilter) ([]*entity.Profile, error)

	// FindByID returns a profile by its ID
	FindByID(ctx context.Context, id uint) (*entity.Profile, error)

	// FindByName returns a profile by its name
	FindByName(ctx context.Context, name string) (*entity.Profile, error)

	// Create creates a new profile
	Create(ctx context.Context, profile *entity.Profile) error

	// Update updates an existing profile
	Update(ctx context.Context, profile *entity.Profile) error

	// Delete deletes profiles by their IDs
	Delete(ctx context.Context, ids []uint) error
}
