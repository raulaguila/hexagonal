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
	profileCacheKeyPrefix = "profile:"
	profileCacheTTL       = 15 * time.Minute
)

// CachedProfileRepository decorates a ProfileRepository with caching logic
type CachedProfileRepository struct {
	delegate output.ProfileRepository
	redis    *redis.Service
}

// NewCachedProfileRepository creates a new cached repository
func NewCachedProfileRepository(delegate output.ProfileRepository, redis *redis.Service) output.ProfileRepository {
	return &CachedProfileRepository{
		delegate: delegate,
		redis:    redis,
	}
}

func (r *CachedProfileRepository) cacheKey(id uint) string {
	return fmt.Sprintf("%sid:%d", profileCacheKeyPrefix, id)
}

// FindByID method with caching
func (r *CachedProfileRepository) FindByID(ctx context.Context, id uint) (*entity.Profile, error) {
	key := r.cacheKey(id)
	client := r.redis.GetClient()

	// Try cache
	val, err := client.Get(ctx, key).Result()
	if err == nil {
		var profile entity.Profile
		if err := json.Unmarshal([]byte(val), &profile); err == nil {
			return &profile, nil
		}
	}

	// Cache miss, call delegate
	profile, err := r.delegate.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Set cache asynchronously to not block response
	go func() {
		if data, err := json.Marshal(profile); err == nil {
			_ = client.Set(context.Background(), key, data, profileCacheTTL).Err()
		}
	}()

	return profile, nil
}

// Pass-through methods (invalidate cache on write)

func (r *CachedProfileRepository) Create(ctx context.Context, profile *entity.Profile) error {
	return r.delegate.Create(ctx, profile)
}

func (r *CachedProfileRepository) Update(ctx context.Context, profile *entity.Profile) error {
	if err := r.delegate.Update(ctx, profile); err != nil {
		return err
	}
	// Invalidate cache
	_ = r.redis.GetClient().Del(ctx, r.cacheKey(profile.ID)).Err()
	return nil
}

func (r *CachedProfileRepository) Delete(ctx context.Context, ids []uint) error {
	if err := r.delegate.Delete(ctx, ids); err != nil {
		return err
	}
	// Invalidate keys
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = r.cacheKey(id)
	}
	if len(keys) > 0 {
		_ = r.redis.GetClient().Del(ctx, keys...).Err()
	}
	return nil
}

// Read-only methods without caching (for now) or complex query caching strategy needed

func (r *CachedProfileRepository) FindAll(ctx context.Context, filter *dto.ProfileFilter) ([]*entity.Profile, error) {
	return r.delegate.FindAll(ctx, filter)
}

func (r *CachedProfileRepository) FindByName(ctx context.Context, name string) (*entity.Profile, error) {
	return r.delegate.FindByName(ctx, name)
}

func (r *CachedProfileRepository) Count(ctx context.Context, filter *dto.ProfileFilter) (int64, error) {
	return r.delegate.Count(ctx, filter)
}
