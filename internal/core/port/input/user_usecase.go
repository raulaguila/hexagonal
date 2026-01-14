package input

import (
	"context"

	"github.com/raulaguila/go-api/internal/core/dto"
)

// UserUseCase defines the interface for user operations
type UserUseCase interface {
	// GetUsers returns a paginated list of users
	GetUsers(ctx context.Context, filter *dto.UserFilter) (*dto.PaginatedOutput[dto.UserOutput], error)

	// GetUserByID returns a user by its ID
	GetUserByID(ctx context.Context, id uint) (*dto.UserOutput, error)

	// CreateUser creates a new user
	CreateUser(ctx context.Context, input *dto.UserInput) (*dto.UserOutput, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, id uint, input *dto.UserInput) (*dto.UserOutput, error)

	// DeleteUsers deletes users by their IDs
	DeleteUsers(ctx context.Context, ids []uint) error

	// ResetPassword resets a user's password
	ResetPassword(ctx context.Context, email string) error

	// SetPassword sets a user's password
	SetPassword(ctx context.Context, email string, input *dto.PasswordInput) error
}
