package input

import (
	"context"

	"github.com/raulaguila/go-api/internal/core/dto"
)

// UserReader defines read operations for users
type UserReader interface {
	GetUsers(ctx context.Context, filter *dto.UserFilter) (*dto.PaginatedOutput[dto.UserOutput], error)
	GetUserByID(ctx context.Context, id string) (*dto.UserOutput, error)
}

// UserWriter defines write operations for users
type UserWriter interface {
	CreateUser(ctx context.Context, input *dto.UserInput) (*dto.UserOutput, error)
	UpdateUser(ctx context.Context, id string, input *dto.UserInput) (*dto.UserOutput, error)
	DeleteUsers(ctx context.Context, ids []string) error
}

// PasswordManager defines password operations
type PasswordManager interface {
	ResetPassword(ctx context.Context, email string) error
	SetPassword(ctx context.Context, email string, input *dto.PasswordInput) error
}

// UserUseCase defines the interface for user operations (Composition)
type UserUseCase interface {
	UserReader
	UserWriter
	PasswordManager
}
