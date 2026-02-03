package mapper

import (
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/model"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
)

// UserToModel converts a User entity to a User model
func UserToModel(u *entity.User) *model.UserModel {
	if u == nil {
		return nil
	}

	return &model.UserModel{
		ID:        u.ID,
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		AuthID:    u.AuthID,
		Auth:      AuthToModel(u.Auth),
		Roles:     RolesToModels(u.Roles),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UserToEntity converts a User model to a User entity
func UserToEntity(m *model.UserModel) *entity.User {
	if m == nil {
		return nil
	}

	return &entity.User{
		ID:        m.ID,
		Name:      m.Name,
		Username:  m.Username,
		Email:     m.Email,
		AuthID:    m.AuthID,
		Auth:      AuthToEntity(m.Auth),
		Roles:     RolesToEntities(m.Roles),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// AuthToModel converts an Auth entity to an Auth model
func AuthToModel(a *entity.Auth) *model.AuthModel {
	if a == nil {
		return nil
	}

	return &model.AuthModel{
		ID:        a.ID,
		Status:    a.Status,
		Token:     a.Token,
		Password:  a.Password,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

// AuthToEntity converts an Auth model to an Auth entity
func AuthToEntity(m *model.AuthModel) *entity.Auth {
	if m == nil {
		return nil
	}

	return &entity.Auth{
		ID:        m.ID,
		Status:    m.Status,
		Token:     m.Token,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// RoleToModel converts a Role entity to a Role model
func RoleToModel(r *entity.Role) *model.RoleModel {
	if r == nil {
		return nil
	}
	return &model.RoleModel{
		ID:          r.ID,
		Name:        r.Name,
		Permissions: r.Permissions,
		Enabled:     r.Enabled,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// RoleToEntity converts a Role model to a Role entity
func RoleToEntity(m *model.RoleModel) *entity.Role {
	if m == nil {
		return nil
	}
	return &entity.Role{
		ID:          m.ID,
		Name:        m.Name,
		Permissions: m.Permissions,
		Enabled:     m.Enabled,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// UsersToEntities converts a slice of User models to User entities
func UsersToEntities(models []*model.UserModel) []*entity.User {
	users := make([]*entity.User, len(models))
	for i, m := range models {
		users[i] = UserToEntity(m)
	}
	return users
}

// RolesToEntities converts a slice of Role models to Role entities
func RolesToEntities(models []*model.RoleModel) []*entity.Role {
	roles := make([]*entity.Role, len(models))
	for i, m := range models {
		roles[i] = RoleToEntity(m)
	}
	return roles
}

// RolesToModels converts a slice of Role entities to Role models
func RolesToModels(entities []*entity.Role) []*model.RoleModel {
	models := make([]*model.RoleModel, len(entities))
	for i, e := range entities {
		models[i] = RoleToModel(e)
	}
	return models
}
