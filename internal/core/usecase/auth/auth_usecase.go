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
	userRepo output.UserRepository
	config   Config
}

// NewAuthUseCase creates a new AuthUseCase instance
func NewAuthUseCase(userRepo output.UserRepository, config Config) input.AuthUseCase {
	return &authUseCase{
		userRepo: userRepo,
		config:   config,
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
func (uc *authUseCase) Refresh(ctx context.Context, userID uint, expiration bool) (*dto.AuthOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.UserNotFound()
	}

	return uc.generateAuthOutput(user, expiration)
}

// Me returns the current authenticated user information
func (uc *authUseCase) Me(ctx context.Context, userID uint) (*dto.UserOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.UserNotFound()
	}

	return dto.EntityToUserOutput(user), nil
}

// generateAuthOutput creates authentication output with tokens
func (uc *authUseCase) generateAuthOutput(user *entity.User, expiration bool) (*dto.AuthOutput, error) {
	// Generate new token if not exists
	if user.Auth.Token == nil {
		token := uuid.New().String()
		user.Auth.SetToken(token)
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
		User:         dto.EntityToUserOutput(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateToken generates a JWT token
func (uc *authUseCase) generateToken(user *entity.User, privateKey *rsa.PrivateKey, expire *time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"token": user.Auth.Token,
		"iat":   now.Unix(),
	}

	if expire != nil {
		claims["exp"] = now.Add(*expire).Unix()
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
}
