package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/raulaguila/go-api/internal/adapter/driven/storage/redis"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/output"
)

const (
	userCacheKeyPrefix = "user:"
	userCacheTTL       = 10 * time.Minute
)

// CachedUserRepository decorates a UserRepository with caching logic
type CachedUserRepository struct {
	delegate output.UserRepository
	redis    *redis.Service
}

// NewCachedUserRepository creates a new cached repository
func NewCachedUserRepository(delegate output.UserRepository, redis *redis.Service) output.UserRepository {
	return &CachedUserRepository{
		delegate: delegate,
		redis:    redis,
	}
}

// Helper to generate cache keys
func (r *CachedUserRepository) keyByID(id uint) string {
	return fmt.Sprintf("%sid:%d", userCacheKeyPrefix, id)
}

func (r *CachedUserRepository) keyByEmail(email string) string {
	return fmt.Sprintf("%semail:%s", userCacheKeyPrefix, email)
}

func (r *CachedUserRepository) keyByUsername(username string) string {
	return fmt.Sprintf("%susername:%s", userCacheKeyPrefix, username)
}

func (r *CachedUserRepository) keyByToken(token string) string {
	return fmt.Sprintf("%stoken:%s", userCacheKeyPrefix, token)
}

// Generic get method to handle cache logic
func (r *CachedUserRepository) getCached(ctx context.Context, key string, fetcher func() (*entity.User, error)) (*entity.User, error) {
	client := r.redis.GetClient()

	// Try cache
	val, err := client.Get(ctx, key).Result()
	if err == nil {
		var user entity.User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			return &user, nil
		}
	}

	// Cache miss
	user, err := fetcher()
	if err != nil {
		return nil, err
	}

	// Set cache async
	go func() {
		if data, err := json.Marshal(user); err == nil {
			_ = client.Set(context.Background(), key, data, userCacheTTL).Err()
		}
	}()

	return user, nil
}

// FindByID returns a user by its ID with caching
func (r *CachedUserRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	return r.getCached(ctx, r.keyByID(id), func() (*entity.User, error) {
		return r.delegate.FindByID(ctx, id)
	})
}

// FindByEmail with caching
func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.getCached(ctx, r.keyByEmail(email), func() (*entity.User, error) {
		return r.delegate.FindByEmail(ctx, email)
	})
}

// FindByUsername with caching
func (r *CachedUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	return r.getCached(ctx, r.keyByUsername(username), func() (*entity.User, error) {
		return r.delegate.FindByUsername(ctx, username)
	})
}

// FindByToken with caching
func (r *CachedUserRepository) FindByToken(ctx context.Context, token string) (*entity.User, error) {
	return r.getCached(ctx, r.keyByToken(token), func() (*entity.User, error) {
		return r.delegate.FindByToken(ctx, token)
	})
}

// Pass-through methods that invalidate cache

func (r *CachedUserRepository) Create(ctx context.Context, user *entity.User) error {
	return r.delegate.Create(ctx, user)
}

func (r *CachedUserRepository) Update(ctx context.Context, user *entity.User) error {
	if err := r.delegate.Update(ctx, user); err != nil {
		return err
	}

	// Invalidate all potential keys for this user
	client := r.redis.GetClient()
	pipe := client.Pipeline()
	pipe.Del(ctx, r.keyByID(user.ID))
	pipe.Del(ctx, r.keyByEmail(user.Email))
	pipe.Del(ctx, r.keyByUsername(user.Username))
	if user.Auth != nil && user.Auth.Token != nil {
		pipe.Del(ctx, r.keyByToken(*user.Auth.Token))
	}
	_, _ = pipe.Exec(ctx)

	return nil
}

func (r *CachedUserRepository) Delete(ctx context.Context, ids []uint) error {
	if err := r.delegate.Delete(ctx, ids); err != nil {
		return err
	}

	// We only have IDs here, so we invalidate ID keys.
	// Note: Ideally we should invalidate Email/Username keys too, but we don't have them here without fetching.
	// For now, TTL handles eventual consistency, or we could fetch before delete if strict consistency is required.
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = r.keyByID(id)
	}
	if len(keys) > 0 {
		_ = r.redis.GetClient().Del(ctx, keys...).Err()
	}
	return nil
}

// Read-only pass-through

func (r *CachedUserRepository) Count(ctx context.Context, filter *dto.UserFilter) (int64, error) {
	return r.delegate.Count(ctx, filter)
}

func (r *CachedUserRepository) FindAll(ctx context.Context, filter *dto.UserFilter) ([]*entity.User, error) {
	return r.delegate.FindAll(ctx, filter)
}
