package model

import (
	"time"

	"github.com/google/uuid"
)

// AuthModel represents the database model for Auth
type AuthModel struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Status    bool      `gorm:"column:status;type:bool;not null;"`
	Token     *string   `gorm:"column:token;type:citext;unique;index"`
	Password  *string   `gorm:"column:password;type:citext;"`
}

// TableName returns the table name for Auth
func (AuthModel) TableName() string {
	return "usr_auth"
}
