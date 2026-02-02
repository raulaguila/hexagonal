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

// roleRepositoryCached implements the RoleRepository interface with caching
type roleRepositoryCached struct {
	repo  output.RoleRepository
	redis *redis.Service
}

// NewRoleRepositoryCached creates a new cached RoleRepository instance
func NewRoleRepositoryCached(repo output.RoleRepository, redis *redis.Service) output.RoleRepository {
	return &roleRepositoryCached{
		repo:  repo,
		redis: redis,
	}
}

func (r *roleRepositoryCached) keyByID(id uuid.UUID) string {
	return fmt.Sprintf("role:id:%s", id.String())
}

func (r *roleRepositoryCached) keyByName(name string) string {
	return fmt.Sprintf("role:name:%s", name)
}

// Count returns the total number of roles matching the filter
func (r *roleRepositoryCached) Count(ctx context.Context, filter *dto.RoleFilter) (int64, error) {
	return r.repo.Count(ctx, filter)
}

// FindAll returns all roles matching the filter
func (r *roleRepositoryCached) FindAll(ctx context.Context, filter *dto.RoleFilter) ([]*entity.Role, error) {
	return r.repo.FindAll(ctx, filter)
}

// FindByID returns a role by its ID
func (r *roleRepositoryCached) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	key := r.keyByID(id)
	role := &entity.Role{}

	if err := r.redis.Get(ctx, key, role); err == nil {
		return role, nil
	}

	role, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := r.redis.Set(ctx, key, role, time.Hour); err != nil {
		return nil, err
	}

	return role, nil
}

// FindByName returns a role by its name
func (r *roleRepositoryCached) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	key := r.keyByName(name)
	role := &entity.Role{}

	if err := r.redis.Get(ctx, key, role); err == nil {
		return role, nil
	}

	role, err := r.repo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if err := r.redis.Set(ctx, key, role, time.Hour); err != nil {
		return nil, err
	}

	return role, nil
}

// Create creates a new role
func (r *roleRepositoryCached) Create(ctx context.Context, role *entity.Role) error {
	if err := r.repo.Create(ctx, role); err != nil {
		return err
	}
	return nil
}

// Update updates an existing role
func (r *roleRepositoryCached) Update(ctx context.Context, role *entity.Role) error {
	if err := r.repo.Update(ctx, role); err != nil {
		return err
	}
	// Invalidate caches
	return r.redis.Del(ctx, r.keyByID(role.ID), r.keyByName(role.Name))
}

// Delete deletes roles by their IDs
func (r *roleRepositoryCached) Delete(ctx context.Context, ids []uuid.UUID) error {
	// Invalidate caches for deleted IDs
	// Since we don't have names here, we can only invalidate ID keys.
	// Names keys will expire naturally or inconsistent until then.
	// Better practice: Find roles first to get names, but that's expensive for delete.
	// For now clear IDs.
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = r.keyByID(id)
	}

	if err := r.repo.Delete(ctx, ids); err != nil {
		return err
	}

	return r.redis.Del(ctx, keys...)
}
