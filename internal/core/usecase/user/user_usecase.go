package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/input"
	"github.com/raulaguila/go-api/internal/core/port/output"
	"github.com/raulaguila/go-api/pkg/apperror"
	"github.com/raulaguila/go-api/pkg/utils"
)

// userUseCase implements the UserUseCase interface
type userUseCase struct {
	userRepo output.UserRepository
}

// NewUserUseCase creates a new UserUseCase instance
func NewUserUseCase(userRepo output.UserRepository) input.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

// GetUsers returns a paginated list of users
func (uc *userUseCase) GetUsers(ctx context.Context, filter *dto.UserFilter) (*dto.PaginatedOutput[dto.UserOutput], error) {
	users, err := uc.userRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	count, err := uc.userRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	outputs := dto.EntitiesToUserOutputs(users)
	return dto.NewPaginatedOutput(outputs, filter.Page, filter.Limit, count), nil
}

// GetUserByID returns a user by its ID
func (uc *userUseCase) GetUserByID(ctx context.Context, id uint) (*dto.UserOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.UserNotFound()
	}

	return dto.EntityToUserOutput(user), nil
}

// CreateUser creates a new user
func (uc *userUseCase) CreateUser(ctx context.Context, input *dto.UserInput) (*dto.UserOutput, error) {
	auth, err := entity.NewAuth(
		utils.Deref(input.ProfileID, uint(0)),
		utils.Deref(input.Status, true),
	)
	if err != nil {
		return nil, err
	}

	user, err := entity.NewUser(
		utils.Deref(input.Name, ""),
		utils.Deref(input.Username, ""),
		utils.Deref(input.Email, ""),
		auth,
	)
	if err != nil {
		return nil, err
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Reload user with relations
	user, err = uc.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return dto.EntityToUserOutput(user), nil
}

// UpdateUser updates an existing user
func (uc *userUseCase) UpdateUser(ctx context.Context, id uint, input *dto.UserInput) (*dto.UserOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.UserNotFound()
	}

	if input.Name != nil {
		user.UpdateName(*input.Name)
	}
	if input.Username != nil {
		user.UpdateUsername(*input.Username)
	}
	if input.Email != nil {
		user.UpdateEmail(*input.Email)
	}
	if input.Status != nil && user.Auth != nil {
		if *input.Status {
			user.Auth.Enable()
		} else {
			user.Auth.Disable()
		}
	}
	if input.ProfileID != nil && user.Auth != nil {
		user.Auth.ProfileID = *input.ProfileID
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Reload user with relations
	user, err = uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.EntityToUserOutput(user), nil
}

// DeleteUsers deletes users by their IDs
func (uc *userUseCase) DeleteUsers(ctx context.Context, ids []uint) error {
	return uc.userRepo.Delete(ctx, ids)
}

// ResetPassword resets a user's password
func (uc *userUseCase) ResetPassword(ctx context.Context, email string) error {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return apperror.UserNotFound()
	}

	if user.Auth == nil || (user.Auth.Password == nil && user.Auth.Token == nil) {
		return nil // Already reset
	}

	user.ResetPassword()
	return uc.userRepo.Update(ctx, user)
}

// SetPassword sets a user's password
func (uc *userUseCase) SetPassword(ctx context.Context, email string, input *dto.PasswordInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return apperror.UserNotFound()
	}

	if user.Auth != nil && user.Auth.Password != nil {
		return apperror.UserHasPassword()
	}

	if err := user.SetPassword(input.Password); err != nil {
		return err
	}

	// Generate token for the user
	token := uuid.New().String()
	user.Auth.SetToken(token)

	return uc.userRepo.Update(ctx, user)
}
