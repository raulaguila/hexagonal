package profile

import (
	"context"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/input"
	"github.com/raulaguila/go-api/internal/core/port/output"
	"github.com/raulaguila/go-api/pkg/apperror"
	"github.com/raulaguila/go-api/pkg/utils"
)

// profileUseCase implements the ProfileUseCase interface
type profileUseCase struct {
	profileRepo output.ProfileRepository
}

// NewProfileUseCase creates a new ProfileUseCase instance
func NewProfileUseCase(profileRepo output.ProfileRepository) input.ProfileUseCase {
	return &profileUseCase{
		profileRepo: profileRepo,
	}
}

// GetProfiles returns a paginated list of profiles
func (uc *profileUseCase) GetProfiles(ctx context.Context, filter *dto.ProfileFilter) (*dto.PaginatedOutput[dto.ProfileOutput], error) {
	profiles, err := uc.profileRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	count, err := uc.profileRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	includePerms := filter.WithPermissions == nil || *filter.WithPermissions
	outputs := dto.EntitiesToProfileOutputs(profiles, includePerms)
	return dto.NewPaginatedOutput(outputs, filter.Page, filter.Limit, count), nil
}

// ListProfiles returns a simple list of profiles (id + name)
func (uc *profileUseCase) ListProfiles(ctx context.Context, filter *dto.ProfileFilter) ([]dto.ItemOutput, error) {
	// Disable pagination for list
	filter.Page = 0
	filter.Limit = 0

	profiles, err := uc.profileRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	outputs := make([]dto.ItemOutput, len(profiles))
	for i, profile := range profiles {
		outputs[i] = dto.ItemOutput{
			ID:   &profile.ID,
			Name: &profile.Name,
		}
	}

	return outputs, nil
}

// GetProfileByID returns a profile by its ID
func (uc *profileUseCase) GetProfileByID(ctx context.Context, id uint) (*dto.ProfileOutput, error) {
	profile, err := uc.profileRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.ProfileNotFound()
	}

	return dto.EntityToProfileOutput(profile, true), nil
}

// CreateProfile creates a new profile
func (uc *profileUseCase) CreateProfile(ctx context.Context, input *dto.ProfileInput) (*dto.ProfileOutput, error) {
	profile := entity.NewProfile(
		utils.Deref(input.Name, ""),
		utils.Deref(input.Permissions, []string{}),
	)

	if err := profile.Validate(); err != nil {
		return nil, err
	}

	if err := uc.profileRepo.Create(ctx, profile); err != nil {
		return nil, err
	}

	return dto.EntityToProfileOutput(profile, true), nil
}

// UpdateProfile updates an existing profile
func (uc *profileUseCase) UpdateProfile(ctx context.Context, id uint, input *dto.ProfileInput) (*dto.ProfileOutput, error) {
	profile, err := uc.profileRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.ProfileNotFound()
	}

	if input.Name != nil {
		profile.UpdateName(*input.Name)
	}
	if input.Permissions != nil {
		profile.UpdatePermissions(*input.Permissions)
	}

	if err := profile.Validate(); err != nil {
		return nil, err
	}

	if err := uc.profileRepo.Update(ctx, profile); err != nil {
		return nil, err
	}

	return dto.EntityToProfileOutput(profile, true), nil
}

// DeleteProfiles deletes profiles by their IDs
func (uc *profileUseCase) DeleteProfiles(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return uc.profileRepo.Delete(ctx, ids)
}
