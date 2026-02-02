package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/raulaguila/go-api/pkg/validator"
)

// User represents a user in the domain
type User struct {
	ID        uuid.UUID
	Name      string
	Username  string
	Email     string
	AuthID    uuid.UUID
	Auth      *Auth
	Roles     []*Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new User entity
func NewUser(name, username, email string, auth *Auth) (*User, error) {
	now := time.Now()
	u := &User{
		ID:        uuid.New(),
		Name:      name,
		Username:  username,
		Email:     email,
		Auth:      auth,
		Roles:     []*Role{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}

// Validate validates the user entity
func (u *User) Validate() error {
	if len(u.Name) < validator.MinNameLength {
		return ErrUserNameTooShort()
	}
	if len(u.Username) < validator.MinUsernameLength {
		return ErrUsernameTooShort()
	}
	if u.Email == "" {
		return ErrEmailRequired()
	}
	if !validator.IsValidEmail(u.Email) {
		return ErrInvalidEmailFormat()
	}
	if u.Auth == nil {
		return ErrAuthRequired()
	}
	// Removed Profile check as it's now Roles, and roles might be empty on creation or assigned later
	return nil
}

// UpdateName updates the user's name
func (u *User) UpdateName(name string, now time.Time) {
	u.Name = name
	u.UpdatedAt = now
}

// UpdateUsername updates the user's username
func (u *User) UpdateUsername(username string, now time.Time) {
	u.Username = username
	u.UpdatedAt = now
}

// UpdateEmail updates the user's email
func (u *User) UpdateEmail(email string, now time.Time) {
	u.Email = email
	u.UpdatedAt = now
}

// SetPassword sets the user's password through Auth
func (u *User) SetPassword(password string, now time.Time) error {
	if u.Auth == nil {
		var err error
		if u.Auth, err = NewAuth(true); err != nil {
			return err
		}
	}
	return u.Auth.SetPassword(password, now)
}

// ValidatePassword validates the password through Auth
func (u *User) ValidatePassword(password string) bool {
	if u.Auth == nil {
		return false
	}
	return u.Auth.ValidatePassword(password)
}

// ResetPassword resets the user's password through Auth
func (u *User) ResetPassword(now time.Time) {
	if u.Auth != nil {
		u.Auth.ResetPassword(now)
	}
	u.UpdatedAt = now
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Auth != nil && u.Auth.IsActive()
}

// IsNew checks if the user is new (has no password set)
func (u *User) IsNew() bool {
	return u.Auth == nil || !u.Auth.HasPassword()
}

// AddRole adds a role to the user
func (u *User) AddRole(role *Role) {
	if role == nil {
		return
	}
	// Verify if already exists to avoid duplicates
	for _, r := range u.Roles {
		if r.ID == role.ID {
			return
		}
	}
	u.Roles = append(u.Roles, role)
}

// HasRole checks if user has a specific role
func (u *User) HasRole(roleName string) bool {
	for _, r := range u.Roles {
		if r.Name == roleName {
			return true
		}
	}
	return false
}

// HasPermission checks if user has a specific permission via any role
func (u *User) HasPermission(permission string) bool {
	for _, r := range u.Roles {
		if r.HasPermission(permission) {
			return true
		}
	}
	return false
}
