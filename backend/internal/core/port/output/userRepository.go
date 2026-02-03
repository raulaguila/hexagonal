package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	// Count returns the total number of users matching the filter
	Count(ctx context.Context, filter *dto.UserFilter) (int64, error)

	// FindAll returns all users matching the filter
	FindAll(ctx context.Context, filter *dto.UserFilter) ([]*entity.User, error)

	// FindByID returns a user by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// FindByUsername returns a user by its username
	FindByUsername(ctx context.Context, username string) (*entity.User, error)

	// FindByEmail returns a user by its email
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// FindByToken returns a user by its authentication token
	FindByToken(ctx context.Context, token string) (*entity.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete deletes users by their IDs
	Delete(ctx context.Context, ids []uuid.UUID) error
}
