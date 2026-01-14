package input

import (
	"context"

	"github.com/raulaguila/go-api/internal/core/dto"
)

// ProfileUseCase defines the interface for profile operations
type ProfileUseCase interface {
	// GetProfiles returns a paginated list of profiles
	GetProfiles(ctx context.Context, filter *dto.ProfileFilter) (*dto.PaginatedOutput[dto.ProfileOutput], error)

	// ListProfiles returns a simple list of profiles (id + name)
	ListProfiles(ctx context.Context, filter *dto.ProfileFilter) ([]dto.ItemOutput, error)

	// GetProfileByID returns a profile by its ID
	GetProfileByID(ctx context.Context, id uint) (*dto.ProfileOutput, error)

	// CreateProfile creates a new profile
	CreateProfile(ctx context.Context, input *dto.ProfileInput) (*dto.ProfileOutput, error)

	// UpdateProfile updates an existing profile
	UpdateProfile(ctx context.Context, id uint, input *dto.ProfileInput) (*dto.ProfileOutput, error)

	// DeleteProfiles deletes profiles by their IDs
	DeleteProfiles(ctx context.Context, ids []uint) error
}
