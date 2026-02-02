package output

import (
	"context"
	"time"
)

// TokenRepository defines the interface for token storage operations (blacklist)
type TokenRepository interface {
	BlacklistToken(ctx context.Context, token string, duration time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}
