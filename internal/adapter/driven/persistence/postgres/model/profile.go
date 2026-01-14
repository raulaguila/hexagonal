package model

import (
	"time"

	"github.com/lib/pq"
)

// ProfileModel represents the database model for Profile
type ProfileModel struct {
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	Name        string         `gorm:"column:name;type:varchar(100);unique;not null;"`
	Permissions pq.StringArray `gorm:"column:permissions;type:text[];not null;"`
}

// TableName returns the table name for Profile
func (ProfileModel) TableName() string {
	return "usr_profile"
}
