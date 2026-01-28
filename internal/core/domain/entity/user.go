package entity

import (
	"time"

	"github.com/raulaguila/go-api/pkg/validator"
)

// User represents a user in the domain
type User struct {
	ID        uint
	Name      string
	Username  string
	Email     string
	AuthID    uint
	Auth      *Auth
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new User entity
func NewUser(name, username, email string, auth *Auth) (*User, error) {
	now := time.Now()
	u := &User{
		Name:      name,
		Username:  username,
		Email:     email,
		Auth:      auth,
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
	if u.Auth == nil || u.Auth.ProfileID == 0 {
		return ErrProfileRequired()
	}
	return nil
}

// UpdateName updates the user's name
func (u *User) UpdateName(name string) {
	u.Name = name
	u.UpdatedAt = time.Now()
}

// UpdateUsername updates the user's username
func (u *User) UpdateUsername(username string) {
	u.Username = username
	u.UpdatedAt = time.Now()
}

// UpdateEmail updates the user's email
func (u *User) UpdateEmail(email string) {
	u.Email = email
	u.UpdatedAt = time.Now()
}

// SetPassword sets the user's password through Auth
func (u *User) SetPassword(password string) error {
	if u.Auth == nil {
		var err error
		if u.Auth, err = NewAuth(0, true); err != nil {
			return err
		}
	}
	return u.Auth.SetPassword(password)
}

// ValidatePassword validates the password through Auth
func (u *User) ValidatePassword(password string) bool {
	if u.Auth == nil {
		return false
	}
	return u.Auth.ValidatePassword(password)
}

// ResetPassword resets the user's password through Auth
func (u *User) ResetPassword() {
	if u.Auth != nil {
		u.Auth.ResetPassword()
	}
	u.UpdatedAt = time.Now()
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Auth != nil && u.Auth.IsActive()
}

// IsNew checks if the user is new (has no password set)
func (u *User) IsNew() bool {
	return u.Auth == nil || !u.Auth.HasPassword()
}

// GetProfileID returns the profile ID
func (u *User) GetProfileID() uint {
	if u.Auth == nil {
		return 0
	}
	return u.Auth.ProfileID
}

// GetProfile returns the user's profile
func (u *User) GetProfile() *Profile {
	if u.Auth == nil {
		return nil
	}
	return u.Auth.Profile
}
