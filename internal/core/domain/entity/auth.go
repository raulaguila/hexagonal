package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/raulaguila/go-api/pkg/validator"
)

// Auth represents the authentication information for a user
type Auth struct {
	ID        uint
	Status    bool
	ProfileID uint
	Profile   *Profile
	Token     *string
	Password  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewAuth creates a new Auth entity
func NewAuth(profileID uint, status bool) (*Auth, error) {
	if profileID == 0 {
		return nil, ErrProfileRequired()
	}
	now := time.Now()
	return &Auth{
		Status:    status,
		ProfileID: profileID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// SetPassword hashes and sets the password
func (a *Auth) SetPassword(password string) error {
	if len(password) < validator.MinPasswordLength {
		return ErrPasswordTooShort()
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hashed := string(hash)
	a.Password = &hashed
	a.UpdatedAt = time.Now()
	return nil
}

// ValidatePassword checks if the provided password matches the stored hash
func (a *Auth) ValidatePassword(password string) bool {
	if a.Password == nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(*a.Password), []byte(password)) == nil
}

// ResetPassword clears the password and token
func (a *Auth) ResetPassword() {
	a.Token = nil
	a.Password = nil
	a.UpdatedAt = time.Now()
}

// SetToken sets the authentication token
func (a *Auth) SetToken(token string) {
	a.Token = &token
	a.UpdatedAt = time.Now()
}

// HasPassword checks if password is set
func (a *Auth) HasPassword() bool {
	return a.Password != nil
}

// IsActive checks if the auth is active
func (a *Auth) IsActive() bool {
	return a.Status && a.Password != nil
}

// Enable enables the auth
func (a *Auth) Enable() {
	a.Status = true
	a.UpdatedAt = time.Now()
}

// Disable disables the auth
func (a *Auth) Disable() {
	a.Status = false
	a.UpdatedAt = time.Now()
}
