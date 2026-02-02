package dto

import (
	"github.com/raulaguila/go-api/pkg/apperror"
	"github.com/raulaguila/go-api/pkg/validator"
)

// RoleInput represents input data for creating/updating a role
type RoleInput struct {
	Name        *string   `json:"name" validate:"omitempty,min=4,max=100"`
	Permissions *[]string `json:"permissions"`
}

// Validate validates the RoleInput
func (p *RoleInput) Validate() error {
	if p.Name != nil && len(*p.Name) < validator.MinProfileNameLength { // Assuming constant name remains or update it
		return apperror.InvalidInput("name", "name must be at least 4 characters")
	}
	return nil
}

// UserInput represents input data for creating/updating a user
type UserInput struct {
	Name     *string   `json:"name" validate:"omitempty,min=5,max=100"`
	Username *string   `json:"username" validate:"omitempty,min=5,max=50"`
	Email    *string   `json:"email" validate:"omitempty,email"`
	Status   *bool     `json:"status"`
	RoleIDs  *[]string `json:"role_ids" validate:"omitempty,dive,uuid"`
}

// Validate validates the UserInput
func (u *UserInput) Validate() error {
	if u.Name != nil && len(*u.Name) < validator.MinNameLength {
		return apperror.InvalidInput("name", "name must be at least 5 characters")
	}
	if u.Username != nil && len(*u.Username) < validator.MinUsernameLength {
		return apperror.InvalidInput("username", "username must be at least 5 characters")
	}
	if u.Email != nil && !validator.IsValidEmail(*u.Email) {
		return apperror.InvalidInput("email", "invalid email format")
	}
	return nil
}

// PasswordInput represents input data for setting a password
type PasswordInput struct {
	Password        string `json:"password" validate:"required,min=6,max=128"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
}

// Validate validates password input
func (p *PasswordInput) Validate() error {
	if p.Password == "" {
		return apperror.InvalidInput("password", "password is required")
	}
	if len(p.Password) < validator.MinPasswordLength {
		return apperror.InvalidInput("password", "password must be at least 6 characters")
	}
	if p.Password != p.PasswordConfirm {
		return apperror.InvalidInput("password_confirm", "passwords do not match")
	}
	return nil
}

// LoginInput represents input data for login
type LoginInput struct {
	Login      string `json:"login" validate:"required"`
	Password   string `json:"password" validate:"required"`
	Expiration bool   `json:"expiration"`
}

// Validate validates the LoginInput
func (l *LoginInput) Validate() error {
	if l.Login == "" {
		return apperror.InvalidInput("login", "login is required")
	}
	if l.Password == "" {
		return apperror.InvalidInput("password", "password is required")
	}
	return nil
}

// IDsInput represents multiple IDs input
type IDsInput struct {
	IDs []string `json:"ids" validate:"required,min=1,dive,uuid"`
}

// Validate validates the IDsInput
func (i *IDsInput) Validate() error {
	if len(i.IDs) == 0 {
		return apperror.InvalidInput("ids", "at least one id is required")
	}
	// uuid validation is handled by 'dive,uuid' tag mainly, but could manual check empty strings
	return nil
}
