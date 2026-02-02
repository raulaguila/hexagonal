package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/input"
	"github.com/raulaguila/go-api/internal/core/port/output"
	"github.com/raulaguila/go-api/pkg/apperror"
	"github.com/raulaguila/go-api/pkg/utils"
)

type userUseCase struct {
	repo     output.UserRepository
	roleRepo output.RoleRepository // Needed to validate roles
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(repo output.UserRepository, roleRepo output.RoleRepository) input.UserUseCase {
	return &userUseCase{
		repo:     repo,
		roleRepo: roleRepo,
	}
}

// GetUsers returns a paginated list of users
func (u *userUseCase) GetUsers(ctx context.Context, filter *dto.UserFilter) (*dto.PaginatedOutput[dto.UserOutput], error) {
	users, err := u.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, apperror.Internal(err.Error(), err)
	}

	total, err := u.repo.Count(ctx, filter)
	if err != nil {
		return nil, apperror.Internal(err.Error(), err)
	}

	return dto.NewPaginatedOutput(
		dto.EntitiesToUserOutputs(users),
		filter.Page,
		filter.Limit,
		total,
	), nil
}

// GetUserByID returns a user by its ID
func (u *userUseCase) GetUserByID(ctx context.Context, id string) (*dto.UserOutput, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, apperror.InvalidInput("id", "invalid uuid format")
	}

	user, err := u.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, apperror.UserNotFound()
	}

	return dto.EntityToUserOutput(user), nil
}

// CreateUser creates a new user
func (u *userUseCase) CreateUser(ctx context.Context, input *dto.UserInput) (*dto.UserOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	existing, _ := u.repo.FindByEmail(ctx, *input.Email)
	if existing != nil {
		return nil, apperror.Conflict("email", "email already exists")
	}

	username := *input.Username
	existing, _ = u.repo.FindByUsername(ctx, username)
	if existing != nil {
		return nil, apperror.Conflict("username", "username already exists")
	}

	// Validate Roles
	var roles []*entity.Role
	if input.RoleIDs != nil && len(*input.RoleIDs) > 0 {
		for _, ridStr := range *input.RoleIDs {
			rid, err := uuid.Parse(ridStr)
			if err != nil {
				return nil, apperror.InvalidInput("role_ids", "invalid role id format")
			}
			role, err := u.roleRepo.FindByID(ctx, rid)
			if err != nil {
				return nil, apperror.InvalidInput("role_ids", "role not found")
			}
			roles = append(roles, role)
		}
	}

	auth, err := entity.NewAuth(utils.Deref(input.Status, true))
	if err != nil {
		return nil, apperror.New("ENTITY_ERROR", err.Error())
	}

	user, err := entity.NewUser(*input.Name, username, *input.Email, auth)
	if err != nil {
		// Entity validation error
		return nil, apperror.New("ENTITY_ERROR", err.Error())
	}

	// Add Roles
	for _, role := range roles {
		user.AddRole(role)
	}

	if input.Status != nil && *input.Status {
		user.Auth.Enable(time.Now())
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return nil, apperror.Internal(err.Error(), err)
	}

	output := dto.EntityToUserOutput(user)
	output.New = nil
	return output, nil
}

// UpdateUser updates an existing user
func (u *userUseCase) UpdateUser(ctx context.Context, id string, input *dto.UserInput) (*dto.UserOutput, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, apperror.InvalidInput("id", "invalid uuid format")
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, err := u.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, apperror.UserNotFound()
	}

	if input.Email != nil && *input.Email != user.Email {
		existing, _ := u.repo.FindByEmail(ctx, *input.Email)
		if existing != nil {
			return nil, apperror.Conflict("email", "email already exists")
		}
		user.UpdateEmail(*input.Email, time.Now())
	}

	if input.Username != nil && *input.Username != user.Username {
		existing, _ := u.repo.FindByUsername(ctx, *input.Username)
		if existing != nil {
			return nil, apperror.Conflict("username", "username already exists")
		}
		user.UpdateUsername(*input.Username, time.Now())
	}

	if input.Name != nil {
		user.UpdateName(*input.Name, time.Now())
	}

	if input.Status != nil {
		if *input.Status {
			user.Auth.Enable(time.Now())
		} else {
			user.Auth.Disable(time.Now())
		}
	}

	if input.RoleIDs != nil {
		// Replace roles logic (same as Create)
		var newRoles []*entity.Role
		for _, ridStr := range *input.RoleIDs {
			rid, err := uuid.Parse(ridStr)
			if err != nil {
				return nil, apperror.InvalidInput("role_ids", "invalid role id format")
			}
			role, err := u.roleRepo.FindByID(ctx, rid)
			if err != nil {
				return nil, apperror.InvalidInput("role_ids", "role not found")
			}
			newRoles = append(newRoles, role)
		}
		user.Roles = newRoles
		user.UpdatedAt = time.Now()
	}

	if err := u.repo.Update(ctx, user); err != nil {
		return nil, apperror.Internal(err.Error(), err)
	}

	return dto.EntityToUserOutput(user), nil
}

// DeleteUsers deletes users by their IDs
func (u *userUseCase) DeleteUsers(ctx context.Context, ids []string) error {
	var uids []uuid.UUID
	for _, id := range ids {
		uid, err := uuid.Parse(id)
		if err != nil {
			return apperror.InvalidInput("ids", "invalid uuid format")
		}
		uids = append(uids, uid)
	}

	if err := u.repo.Delete(ctx, uids); err != nil {
		return apperror.Internal(err.Error(), err)
	}

	return nil
}

// ResetPassword resets a user's password
func (u *userUseCase) ResetPassword(ctx context.Context, email string) error {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil
	}

	u.ResetPasswordHelper(user) // Helper defined locally or direct call
	// Logic from View:
	user.ResetPassword(time.Now())
	return u.repo.Update(ctx, user)
}

func (u *userUseCase) ResetPasswordHelper(user *entity.User) {
	user.ResetPassword(time.Now())
}

// SetPassword sets a user's password
func (u *userUseCase) SetPassword(ctx context.Context, email string, input *dto.PasswordInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return apperror.UserNotFound()
	}

	// Logic: check user has password
	// entity.Auth HasPassword check?
	// Let's assume user.Auth.HasPassword() exists if entity definition had it or I should check user.Auth directly.
	// Entity view showed ValidatePassword/etc.
	// I will write:
	if user.Auth != nil && user.Auth.Password != nil {
		return apperror.UserHasPassword()
	}

	if err := user.SetPassword(input.Password, time.Now()); err != nil {
		return apperror.New("ENTITY_ERROR", err.Error())
	}

	// Set Token (random UUID)
	token := uuid.New().String()
	user.Auth.SetToken(token, time.Now())

	return u.repo.Update(ctx, user)
}
