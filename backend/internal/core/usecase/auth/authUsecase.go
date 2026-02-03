package auth

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/input"
	"github.com/raulaguila/go-api/internal/core/port/output"
	"github.com/raulaguila/go-api/pkg/apperror"
)

// Config holds JWT configuration
type Config struct {
	AccessPrivateKey  *rsa.PrivateKey
	AccessExpiration  time.Duration
	RefreshPrivateKey *rsa.PrivateKey
	RefreshExpiration time.Duration
}

// authUseCase implements the AuthUseCase interface
type authUseCase struct {
	userRepo  output.UserRepository
	tokenRepo output.TokenRepository
	config    Config
}

// NewAuthUseCase creates a new AuthUseCase instance
func NewAuthUseCase(userRepo output.UserRepository, tokenRepo output.TokenRepository, config Config) input.AuthUseCase {
	return &authUseCase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		config:    config,
	}
}

// Login authenticates a user and returns tokens
func (uc *authUseCase) Login(ctx context.Context, input *dto.LoginInput) (*dto.AuthOutput, error) {
	user, err := uc.userRepo.FindByUsername(ctx, input.Login)
	if err != nil {
		return nil, apperror.UserNotFound()
	}

	if !user.ValidatePassword(input.Password) {
		return nil, apperror.InvalidCredentials()
	}

	if user.Auth == nil || !user.Auth.Status || user.Auth.Password == nil {
		return nil, apperror.DisabledUser()
	}

	return uc.generateAuthOutput(user, input.Expiration)
}

// Refresh refreshes the authentication tokens
func (uc *authUseCase) Refresh(ctx context.Context, userID string, expiration bool) (*dto.AuthOutput, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperror.InvalidInput("user_id", "invalid uuid format")
	}

	user, err := uc.userRepo.FindByID(ctx, uid)
	if err != nil {
		return nil, apperror.UserNotFound()
	}

	return uc.generateAuthOutput(user, expiration)
}

// Me returns the current authenticated user information
func (uc *authUseCase) Me(ctx context.Context, userID string) (*dto.UserOutput, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperror.InvalidInput("user_id", "invalid uuid format")
	}

	user, err := uc.userRepo.FindByID(ctx, uid)
	if err != nil {
		return nil, apperror.UserNotFound()
	}

	return dto.EntityToUserOutput(user, true), nil
}

// Logout invalidates the provided token
func (uc *authUseCase) Logout(ctx context.Context, token string) error {
	if uc.tokenRepo == nil {
		// Log warning or return error if redis is not configured
		return nil // Or apperror.Internal("Token repository not configured")
	}
	// We blacklist the token for the duration of the refresh token expiration
	// conservatively, or access token expiration.
	// Since we don't know if it's access or refresh without parsing, and we might blacklist both.
	// For simplicity, blacklist for RefreshExpiration (longest).
	return uc.tokenRepo.BlacklistToken(ctx, token, uc.config.RefreshExpiration)
}

// generateAuthOutput creates authentication output with tokens
func (uc *authUseCase) generateAuthOutput(user *entity.User, expiration bool) (*dto.AuthOutput, error) {
	// Generate new token if not exists
	if user.Auth.Token == nil {
		token := uuid.New().String()
		user.Auth.SetToken(token, time.Now())
		if err := uc.userRepo.Update(context.Background(), user); err != nil {
			return nil, err
		}
	}

	accessToken, err := uc.generateToken(user, uc.config.AccessPrivateKey, func() *time.Duration {
		if expiration {
			return &uc.config.AccessExpiration
		}
		return nil
	}())
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.generateToken(user, uc.config.RefreshPrivateKey, func() *time.Duration {
		if expiration {
			return &uc.config.RefreshExpiration
		}
		return nil
	}())
	if err != nil {
		return nil, err
	}

	return &dto.AuthOutput{
		User:         dto.EntityToUserOutput(user, true),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateToken generates a JWT token
func (uc *authUseCase) generateToken(user *entity.User, privateKey *rsa.PrivateKey, expire *time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"id":    user.ID.String(), // Add ID to claims for easier retrieval?
		"token": user.Auth.Token,
		"iat":   now.Unix(),
	}

	if expire != nil {
		claims["exp"] = now.Add(*expire).Unix()
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
}
