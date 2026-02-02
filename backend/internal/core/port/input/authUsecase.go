package input

import (
	"context"

	"github.com/raulaguila/go-api/internal/core/dto"
)

// AuthUseCase defines the interface for authentication operations
type AuthUseCase interface {
	// Login authenticates a user and returns tokens
	Login(ctx context.Context, input *dto.LoginInput) (*dto.AuthOutput, error)

	// Refresh refreshes the authentication tokens
	Refresh(ctx context.Context, userID string, expiration bool) (*dto.AuthOutput, error)

	// Me returns the current authenticated user information
	Me(ctx context.Context, userID string) (*dto.UserOutput, error)
	// Logout invalidates the provided token
	Logout(ctx context.Context, token string) error
}
