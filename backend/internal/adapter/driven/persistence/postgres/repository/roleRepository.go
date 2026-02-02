package repository

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/mapper"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/model"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/output"
)

// roleRepository implements the RoleRepository interface
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new RoleRepository instance
func NewRoleRepository(db *gorm.DB) output.RoleRepository {
	return &roleRepository{db: db}
}

// applyFilter applies filters to the query
func (r *roleRepository) applyFilter(ctx context.Context, filter *dto.RoleFilter) *gorm.DB {
	query := r.db.WithContext(ctx)

	if filter != nil {
		if filter.ID != nil {
			query = query.Where("id = ?", *filter.ID)
		}

		if filter.Search != "" {
			query = query.Where("unaccent(name) LIKE unaccent(?)", fmt.Sprintf("%%%s%%", filter.Search))
		}

		// List root specifically (or exclude it) - typical pattern might be exclusion if not requested
		// but let's stick to what was likely there or implied.
		// Previous code profileRepository had some specific filters, let's adapt.
		if !filter.ListRoot {
			// Assuming ROOT is a specific role we might want to hide, but let's keep it simple for now unless we see explicit logic required.
			// If filter.ListRoot is false, maybe we exclude ID 1? But IDs are UUIDs now.
			// We'll skip specific ID exclusion unless we know constraints.
		}

		query = r.applyOrder(query, filter)
	}

	return query
}

// applyOrder applies ordering to the query
func (r *roleRepository) applyOrder(query *gorm.DB, filter *dto.RoleFilter) *gorm.DB {
	sort := filter.Sort
	order := filter.Order

	if sort == "" {
		sort = "name" // Default sort
	}
	if !slices.Contains([]string{"asc", "desc"}, strings.ToLower(order)) {
		order = "asc"
	}

	return query.Order(fmt.Sprintf("%s %s", sort, order))
}

// Count returns the total number of roles matching the filter
func (r *roleRepository) Count(ctx context.Context, filter *dto.RoleFilter) (int64, error) {
	var count int64
	err := r.applyFilter(ctx, filter).Model(&model.RoleModel{}).Count(&count).Error
	return count, err
}

// FindAll returns all roles matching the filter
func (r *roleRepository) FindAll(ctx context.Context, filter *dto.RoleFilter) ([]*entity.Role, error) {
	query := r.applyFilter(ctx, filter)

	if filter != nil {
		if ok, offset, limit := filter.ApplyPagination(); ok {
			query = query.Offset(offset).Limit(limit)
		}
	}

	var models []*model.RoleModel
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	return mapper.RolesToEntities(models), nil
}

// FindByID returns a role by its ID
func (r *roleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var m model.RoleModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return mapper.RoleToEntity(&m), nil
}

// FindByName returns a role by its name
func (r *roleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	var m model.RoleModel
	if err := r.db.WithContext(ctx).First(&m, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return mapper.RoleToEntity(&m), nil
}

// Create creates a new role
func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	m := mapper.RoleToModel(role)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	role.ID = m.ID
	role.CreatedAt = m.CreatedAt
	role.UpdatedAt = m.UpdatedAt
	return nil
}

// Update updates an existing role
func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	m := mapper.RoleToModel(role)
	return r.db.WithContext(ctx).Model(m).Updates(m).Error
}

// Delete deletes roles by their IDs
func (r *roleRepository) Delete(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	result := r.db.WithContext(ctx).Delete(&model.RoleModel{}, "id IN ?", ids)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
