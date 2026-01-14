package repository

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/mapper"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/model"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/output"
)

const (
	userTable    = "usr_user"
	authTable    = "usr_auth"
	profileTable = "usr_profile"
)

// userRepository implements the UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) output.UserRepository {
	return &userRepository{db: db}
}

// applyFilter applies filters to the query
func (r *userRepository) applyFilter(ctx context.Context, filter *dto.UserFilter) *gorm.DB {
	query := r.db.WithContext(ctx)

	if filter != nil {
		if filter.ID != nil {
			query = query.Where(userTable+".id = ?", *filter.ID)
		}

		if filter.Status != nil {
			query = query.Where(authTable+".status = ?", *filter.Status)
		}

		if filter.ProfileID != 0 {
			query = query.Where(authTable+".profile_id = ?", filter.ProfileID)
		}

		query = query.Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.auth_id", authTable, authTable, userTable))
		query = query.Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.profile_id", profileTable, profileTable, authTable))

		if filter.Search != "" {
			columns := []string{
				userTable + ".name",
				userTable + ".username",
				userTable + ".mail",
				profileTable + ".name",
			}
			var conditions []string
			for _, col := range columns {
				conditions = append(conditions, fmt.Sprintf("unaccent(LOWER(%s)) LIKE unaccent(LOWER('%%%s%%'))", col, filter.Search))
			}
			query = query.Where(strings.Join(conditions, " OR "))
		}

		query = r.applyOrder(query, filter)
	}

	return query.Group(userTable + ".id")
}

// applyOrder applies ordering to the query
func (r *userRepository) applyOrder(query *gorm.DB, filter *dto.UserFilter) *gorm.DB {
	sort := filter.Sort
	order := filter.Order

	// Use default values instead of reading from environment
	// This makes the repository pure and dependency-free
	if sort == "" {
		sort = "id" // Default sort field
	}
	if !slices.Contains([]string{"asc", "desc"}, strings.ToLower(order)) {
		order = "asc" // Default order
	}

	if !strings.Contains(sort, ".") {
		sort = userTable + "." + sort
	}

	return query.Order(fmt.Sprintf("%s %s", sort, order))
}

// Count returns the total number of users matching the filter
func (r *userRepository) Count(ctx context.Context, filter *dto.UserFilter) (int64, error) {
	var count int64
	err := r.applyFilter(ctx, filter).Model(&model.UserModel{}).Count(&count).Error
	return count, err
}

// FindAll returns all users matching the filter
func (r *userRepository) FindAll(ctx context.Context, filter *dto.UserFilter) ([]*entity.User, error) {
	query := r.applyFilter(ctx, filter)

	if filter != nil {
		if ok, offset, limit := filter.ApplyPagination(); ok {
			query = query.Offset(offset).Limit(limit)
		}
	}

	var models []*model.UserModel
	if err := query.Preload("Auth.Profile").Find(&models).Error; err != nil {
		return nil, err
	}

	return mapper.UsersToEntities(models), nil
}

// FindByID returns a user by its ID
func (r *userRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).Preload("Auth.Profile").First(&m, id).Error; err != nil {
		return nil, err
	}
	return mapper.UserToEntity(&m), nil
}

// FindByUsername returns a user by its username
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).Preload("Auth.Profile").Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.UserToEntity(&m), nil
}

// FindByEmail returns a user by its email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).Preload("Auth.Profile").Where("mail = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.UserToEntity(&m), nil
}

// FindByToken returns a user by its authentication token
func (r *userRepository) FindByToken(ctx context.Context, token string) (*entity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.auth_id", authTable, authTable, userTable)).
		Preload("Auth.Profile").
		First(&m, authTable+".token = ?", token).Error; err != nil {
		return nil, err
	}
	return mapper.UserToEntity(&m), nil
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	m := mapper.UserToModel(user)
	if err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	user.ID = m.ID
	user.AuthID = m.AuthID
	if m.Auth != nil {
		user.Auth.ID = m.Auth.ID
	}
	user.CreatedAt = m.CreatedAt
	user.UpdatedAt = m.UpdatedAt
	return nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	m := mapper.UserToModel(user)

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update Auth first
		if m.Auth != nil {
			if err := tx.Model(m.Auth).Updates(map[string]any{
				"status":     m.Auth.Status,
				"profile_id": m.Auth.ProfileID,
				"token":      m.Auth.Token,
				"password":   m.Auth.Password,
			}).Error; err != nil {
				return err
			}
		}

		// Update User
		return tx.Model(m).Updates(map[string]any{
			"name":     m.Name,
			"username": m.Username,
			"mail":     m.Email,
			"auth_id":  m.AuthID,
		}).Error
	})
}

// Delete deletes users by their IDs
func (r *userRepository) Delete(ctx context.Context, ids []uint) error {
	var models []*model.UserModel
	if err := r.db.WithContext(ctx).Find(&models, ids).Error; err != nil {
		return err
	}
	if len(models) == 0 {
		return gorm.ErrRecordNotFound
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Select(clause.Associations).Where("id IN ?", ids).Delete(&model.UserModel{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
