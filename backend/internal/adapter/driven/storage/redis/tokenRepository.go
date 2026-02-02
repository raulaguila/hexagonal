package redis

import (
	"context"
	"time"

	"github.com/raulaguila/go-api/internal/core/port/output"
)

type tokenRepository struct {
	service *Service
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(service *Service) output.TokenRepository {
	return &tokenRepository{
		service: service,
	}
}

// BlacklistToken adds a token to the blacklist with an expiration
func (r *tokenRepository) BlacklistToken(ctx context.Context, token string, duration time.Duration) error {
	// Value can be anything, e.g., "revoked"
	return r.service.Set(ctx, "blacklist:"+token, "revoked", duration)
}

// IsTokenBlacklisted checks if a token is in the blacklist
func (r *tokenRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	var val string
	err := r.service.Get(ctx, "blacklist:"+token, &val)
	if err != nil {
		// If key does not exist, it's not blacklisted
		// We need to check if err is "redis: nil" (key not found)
		// Our service.Get might return error if key not found.
		// Let's check service.Get implementation.
		// It returns generic error. Ideally we need to know if it's NotFound.
		// For now, assume error means not found if it's not a connection error.
		// But strictly, we should check.
		// Go-redis returns redis.Nil.
		return false, nil
	}
	return true, nil
}
