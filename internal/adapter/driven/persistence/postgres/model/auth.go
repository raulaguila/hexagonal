package model

import (
	"time"
)

// AuthModel represents the database model for Auth
type AuthModel struct {
	ID        uint          `gorm:"primarykey"`
	CreatedAt time.Time     `gorm:"autoCreateTime"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime"`
	Status    bool          `gorm:"column:status;type:bool;not null;"`
	ProfileID uint          `gorm:"column:profile_id;type:bigint;not null;index;"`
	Profile   *ProfileModel `gorm:"foreignKey:ProfileID"`
	Token     *string       `gorm:"column:token;type:varchar(255);unique;index"`
	Password  *string       `gorm:"column:password;type:varchar(255);"`
}

// TableName returns the table name for Auth
func (AuthModel) TableName() string {
	return "usr_auth"
}
