package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/raulaguila/go-api/internal/adapter/driven/storage/redis"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/output"
)

// userRepositoryCached implements the UserRepository interface with caching
type userRepositoryCached struct {
	repo  output.UserRepository
	redis *redis.Service
}

// NewUserRepositoryCached creates a new cached UserRepository instance
func NewUserRepositoryCached(repo output.UserRepository, redis *redis.Service) output.UserRepository {
	return &userRepositoryCached{
		repo:  repo,
		redis: redis,
	}
}

func (r *userRepositoryCached) keyByID(id uuid.UUID) string {
	return fmt.Sprintf("user:id:%s", id.String())
}

func (r *userRepositoryCached) keyByEmail(email string) string {
	return fmt.Sprintf("user:email:%s", email)
}

func (r *userRepositoryCached) keyByToken(token string) string {
	return fmt.Sprintf("user:token:%s", token)
}

func (r *userRepositoryCached) keyByUsername(username string) string {
	return fmt.Sprintf("user:username:%s", username)
}

// Count returns the total number of users matching the filter
func (r *userRepositoryCached) Count(ctx context.Context, filter *dto.UserFilter) (int64, error) {
	return r.repo.Count(ctx, filter)
}

// FindAll returns all users matching the filter
func (r *userRepositoryCached) FindAll(ctx context.Context, filter *dto.UserFilter) ([]*entity.User, error) {
	return r.repo.FindAll(ctx, filter)
}

// FindByID returns a user by its ID
func (r *userRepositoryCached) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	key := r.keyByID(id)
	user := &entity.User{}

	if err := r.redis.Get(ctx, key, user); err == nil {
		return user, nil
	}

	user, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := r.redis.Set(ctx, key, user, time.Hour); err != nil {
		return nil, err
	}

	return user, nil
}

// FindByUsername returns a user by its username
func (r *userRepositoryCached) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	key := r.keyByUsername(username)
	user := &entity.User{}

	if err := r.redis.Get(ctx, key, user); err == nil {
		return user, nil
	}

	user, err := r.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if err := r.redis.Set(ctx, key, user, time.Hour); err != nil {
		return nil, err
	}

	return user, nil
}

// FindByEmail returns a user by its email
func (r *userRepositoryCached) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	key := r.keyByEmail(email)
	user := &entity.User{}

	if err := r.redis.Get(ctx, key, user); err == nil {
		return user, nil
	}

	user, err := r.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := r.redis.Set(ctx, key, user, time.Hour); err != nil {
		return nil, err
	}

	return user, nil
}

// FindByToken returns a user by its authentication token
func (r *userRepositoryCached) FindByToken(ctx context.Context, token string) (*entity.User, error) {
	// Tokens might change often or be short lived. Caching depends on strategy.
	// Assuming we cache for checking active sessions quickly.
	key := r.keyByToken(token)
	user := &entity.User{}

	if err := r.redis.Get(ctx, key, user); err == nil {
		return user, nil
	}

	user, err := r.repo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if err := r.redis.Set(ctx, key, user, time.Hour); err != nil {
		return nil, err
	}

	return user, nil
}

// Create creates a new user
func (r *userRepositoryCached) Create(ctx context.Context, user *entity.User) error {
	if err := r.repo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

// Update updates an existing user
func (r *userRepositoryCached) Update(ctx context.Context, user *entity.User) error {
	if err := r.repo.Update(ctx, user); err != nil {
		return err
	}

	// Invalidate all keys
	keys := []string{
		r.keyByID(user.ID),
		r.keyByEmail(user.Email),
		r.keyByUsername(user.Username),
	}
	if user.Auth != nil && user.Auth.Token != nil {
		keys = append(keys, r.keyByToken(*user.Auth.Token))
	}

	return r.redis.Del(ctx, keys...)
}

// Delete deletes users by their IDs
func (r *userRepositoryCached) Delete(ctx context.Context, ids []uuid.UUID) error {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = r.keyByID(id)
	}

	if err := r.repo.Delete(ctx, ids); err != nil {
		return err
	}

	return r.redis.Del(ctx, keys...)
}
