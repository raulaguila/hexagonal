package repository

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/mapper"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/model"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/output"
)

// profileRepository implements the ProfileRepository interface
type profileRepository struct {
	db *gorm.DB
}

// NewProfileRepository creates a new ProfileRepository instance
func NewProfileRepository(db *gorm.DB) output.ProfileRepository {
	return &profileRepository{db: db}
}

// applyFilter applies filters to the query
func (r *profileRepository) applyFilter(ctx context.Context, filter *dto.ProfileFilter) *gorm.DB {
	query := r.db.WithContext(ctx)

	if filter != nil {
		if filter.ID != nil {
			query = query.Where("id = ?", *filter.ID)
		}

		if filter.Search != "" {
			searchLike := fmt.Sprintf("unaccent(LOWER(name)) LIKE unaccent(LOWER('%%%s%%'))", filter.Search)
			query = query.Where(searchLike)
		}

		if filter.WithPermissions != nil && !*filter.WithPermissions {
			query = query.Omit("permissions")
		}

		if !filter.ListRoot {
			query = query.Where("name != ?", "ROOT")
		}

		query = r.applyOrder(query, filter)
	}

	return query.Group("id")
}

// applyOrder applies ordering to the query
func (r *profileRepository) applyOrder(query *gorm.DB, filter *dto.ProfileFilter) *gorm.DB {
	sort := filter.Sort
	order := filter.Order

	if sort == "" {
		sort = os.Getenv("API_DEFAULT_SORT")
	}
	if !slices.Contains([]string{"asc", "desc"}, strings.ToLower(order)) {
		order = os.Getenv("API_DEFAULT_ORDER")
	}

	return query.Order(fmt.Sprintf("%s %s", sort, order))
}

// Count returns the total number of profiles matching the filter
func (r *profileRepository) Count(ctx context.Context, filter *dto.ProfileFilter) (int64, error) {
	var count int64
	err := r.applyFilter(ctx, filter).Model(&model.ProfileModel{}).Count(&count).Error
	return count, err
}

// FindAll returns all profiles matching the filter
func (r *profileRepository) FindAll(ctx context.Context, filter *dto.ProfileFilter) ([]*entity.Profile, error) {
	query := r.applyFilter(ctx, filter)

	if filter != nil {
		if ok, offset, limit := filter.ApplyPagination(); ok {
			query = query.Offset(offset).Limit(limit)
		}
	}

	var models []*model.ProfileModel
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	return mapper.ProfilesToEntities(models), nil
}

// FindByID returns a profile by its ID
func (r *profileRepository) FindByID(ctx context.Context, id uint) (*entity.Profile, error) {
	var m model.ProfileModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return mapper.ProfileToEntity(&m), nil
}

// FindByName returns a profile by its name
func (r *profileRepository) FindByName(ctx context.Context, name string) (*entity.Profile, error) {
	var m model.ProfileModel
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.ProfileToEntity(&m), nil
}

// Create creates a new profile
func (r *profileRepository) Create(ctx context.Context, profile *entity.Profile) error {
	m := mapper.ProfileToModel(profile)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	profile.ID = m.ID
	profile.CreatedAt = m.CreatedAt
	profile.UpdatedAt = m.UpdatedAt
	return nil
}

// Update updates an existing profile
func (r *profileRepository) Update(ctx context.Context, profile *entity.Profile) error {
	m := mapper.ProfileToModel(profile)
	return r.db.WithContext(ctx).Model(m).Updates(map[string]any{
		"name":        m.Name,
		"permissions": m.Permissions,
	}).Error
}

// Delete deletes profiles by their IDs
func (r *profileRepository) Delete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&model.ProfileModel{}, ids)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
