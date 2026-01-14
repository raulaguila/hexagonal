package mapper

import (
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/model"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
)

// MapSlice is a generic helper to map slices
func MapSlice[T any, U any](items []T, fn func(T) U) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}

// ProfileToModel converts a Profile entity to a ProfileModel
func ProfileToModel(e *entity.Profile) *model.ProfileModel {
	if e == nil {
		return nil
	}
	return &model.ProfileModel{
		ID:          e.ID,
		Name:        e.Name,
		Permissions: e.Permissions,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// ProfileToEntity converts a ProfileModel to a Profile entity
func ProfileToEntity(m *model.ProfileModel) *entity.Profile {
	if m == nil {
		return nil
	}
	return &entity.Profile{
		ID:          m.ID,
		Name:        m.Name,
		Permissions: m.Permissions,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// AuthToModel converts an Auth entity to an AuthModel
func AuthToModel(e *entity.Auth) *model.AuthModel {
	if e == nil {
		return nil
	}
	return &model.AuthModel{
		ID:        e.ID,
		Status:    e.Status,
		ProfileID: e.ProfileID,
		Profile:   ProfileToModel(e.Profile),
		Token:     e.Token,
		Password:  e.Password,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// AuthToEntity converts an AuthModel to an Auth entity
func AuthToEntity(m *model.AuthModel) *entity.Auth {
	if m == nil {
		return nil
	}
	return &entity.Auth{
		ID:        m.ID,
		Status:    m.Status,
		ProfileID: m.ProfileID,
		Profile:   ProfileToEntity(m.Profile),
		Token:     m.Token,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// UserToModel converts a User entity to a UserModel
func UserToModel(e *entity.User) *model.UserModel {
	if e == nil {
		return nil
	}
	return &model.UserModel{
		ID:        e.ID,
		Name:      e.Name,
		Username:  e.Username,
		Email:     e.Email,
		AuthID:    e.AuthID,
		Auth:      AuthToModel(e.Auth),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// UserToEntity converts a UserModel to a User entity
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
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// UsersToEntities converts a slice of UserModels to User entities
func UsersToEntities(models []*model.UserModel) []*entity.User {
	return MapSlice(models, UserToEntity)
}

// ProfilesToEntities converts a slice of ProfileModels to Profile entities
func ProfilesToEntities(models []*model.ProfileModel) []*entity.Profile {
	return MapSlice(models, ProfileToEntity)
}

// UsersToModels converts a slice of User entities to UserModels
func UsersToModels(entities []*entity.User) []*model.UserModel {
	return MapSlice(entities, UserToModel)
}

// ProfilesToModels converts a slice of Profile entities to ProfileModels
func ProfilesToModels(entities []*entity.Profile) []*model.ProfileModel {
	return MapSlice(entities, ProfileToModel)
}
