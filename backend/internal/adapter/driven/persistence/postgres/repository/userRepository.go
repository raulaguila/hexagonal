package repository

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/mapper"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/model"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/raulaguila/go-api/internal/core/port/output"
)

const (
	userTable     = "usr_user"
	authTable     = "usr_auth"
	roleTable     = "usr_role"
	userRoleTable = "usr_user_role"
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
		if filter.ID != nil && *filter.ID != "" {
			query = query.Where(userTable+".id = ?", *filter.ID)
		}

		if filter.Status != nil {
			query = query.Where(authTable+".status = ?", *filter.Status)
		}

		// Join with Auth table
		query = query.Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.auth_id", authTable, authTable, userTable))

		if filter.RoleID != nil && *filter.RoleID != "" {
			// Join with UserRole table to filter by RoleID
			query = query.Joins(fmt.Sprintf("JOIN %s ON %s.user_id = %s.id", userRoleTable, userRoleTable, userTable))
			query = query.Where(userRoleTable+".role_id = ?", *filter.RoleID)
		}

		if filter.Search != "" {
			searchPattern := "%" + filter.Search + "%"
			query = query.Where(
				"unaccent("+userTable+".name) ILIKE unaccent(?) OR "+
					"unaccent("+userTable+".username) ILIKE unaccent(?) OR "+
					"unaccent("+userTable+".email) ILIKE unaccent(?)",
				searchPattern, searchPattern, searchPattern,
			)
		}

		query = r.applyOrder(query, filter)
	}

	return query.Group(userTable + ".id")
}

// applyOrder applies ordering to the query
func (r *userRepository) applyOrder(query *gorm.DB, filter *dto.UserFilter) *gorm.DB {
	sort := filter.Sort
	order := filter.Order

	if sort == "" {
		sort = "name" // Default sort
	}
	if !slices.Contains([]string{"asc", "desc"}, strings.ToLower(order)) {
		order = "asc"
	}

	if !strings.Contains(sort, ".") {
		sort = userTable + "." + sort
	}

	return query.Order(fmt.Sprintf("%s %s", sort, order))
}

// Count returns the total number of users matching the filter
func (r *userRepository) Count(ctx context.Context, filter *dto.UserFilter) (int64, error) {
	var count int64
	// Need to be careful with Count distinct if joining many-to-many
	err := r.applyFilter(ctx, filter).Model(&model.UserModel{}).Distinct(userTable + ".id").Count(&count).Error
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
	// Preload Auth and Roles
	if err := query.Preload("Auth").Preload("Roles").Find(&models).Error; err != nil {
		return nil, err
	}

	return mapper.UsersToEntities(models), nil
}

// FindByID returns a user by its ID
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).Preload("Auth").Preload("Roles").First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return mapper.UserToEntity(&m), nil
}

// FindByUsername returns a user by its username
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).Preload("Auth").Preload("Roles").Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.UserToEntity(&m), nil
}

// FindByEmail returns a user by its email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).Preload("Auth").Preload("Roles").Where("email = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.UserToEntity(&m), nil
}

// FindByToken returns a user by its authentication token
func (r *userRepository) FindByToken(ctx context.Context, token string) (*entity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.auth_id", authTable, authTable, userTable)).
		Preload("Auth").Preload("Roles").
		First(&m, authTable+".token = ?", token).Error; err != nil {
		return nil, err
	}
	return mapper.UserToEntity(&m), nil
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	m := mapper.UserToModel(user)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// First create the user without roles
		if err := tx.WithContext(ctx).Omit("Roles").Create(m).Error; err != nil {
			return err
		}

		// Then add role associations if any exist
		if len(m.Roles) > 0 {
			if err := tx.Model(m).Association("Roles").Replace(m.Roles); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
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
				"status":   m.Auth.Status,
				"token":    m.Auth.Token,
				"password": m.Auth.Password,
			}).Error; err != nil {
				return err
			}
		}

		// Update User fields
		if err := tx.Model(m).Updates(map[string]any{
			"name":     m.Name,
			"username": m.Username,
			"email":    m.Email,
			"auth_id":  m.AuthID,
		}).Error; err != nil {
			return err
		}

		// Update Roles Association
		// Replace current roles with the new set
		return tx.Model(m).Association("Roles").Replace(m.Roles)
	})
}

// Delete deletes users by their IDs
func (r *userRepository) Delete(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
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
