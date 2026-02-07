package dto

// RoleInput represents input data for creating/updating a role
type RoleInput struct {
	Name        *string   `json:"name" validate:"omitempty,min=4,max=100"`
	Permissions *[]string `json:"permissions"`
	Enabled     *bool     `json:"enabled"`
}

// UserInput represents input data for creating/updating a user
type UserInput struct {
	Name     *string   `json:"name" validate:"omitempty,min=5,max=100"`
	Username *string   `json:"username" validate:"omitempty,min=5,max=50"`
	Email    *string   `json:"email" validate:"omitempty,email"`
	Status   *bool     `json:"status"`
	RoleIDs  *[]string `json:"role_ids" validate:"omitempty,dive,uuid"`
}

// PasswordInput represents input data for setting a password
type PasswordInput struct {
	Password        string `json:"password" validate:"required,min=6,max=128"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
}

// LoginInput represents input data for login
type LoginInput struct {
	Login      string `json:"login" validate:"required"`
	Password   string `json:"password" validate:"required"`
	Expiration bool   `json:"expiration"`
}

// IDsInput represents multiple IDs input
type IDsInput struct {
	IDs []string `json:"ids" validate:"required,min=1,dive,uuid"`
}
