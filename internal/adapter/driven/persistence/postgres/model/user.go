package model

import (
	"time"

	"gorm.io/gorm"
)

// UserModel represents the database model for User
type UserModel struct {
	ID        uint       `gorm:"primarykey"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	Name      string     `gorm:"column:name;"`
	Username  string     `gorm:"column:username;"`
	Email     string     `gorm:"column:mail;"`
	AuthID    uint       `gorm:"column:auth_id;"`
	Auth      *AuthModel `gorm:"constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for User
func (UserModel) TableName() string {
	return "usr_user"
}

// AfterDelete hook to delete associated Auth
func (u *UserModel) AfterDelete(tx *gorm.DB) error {
	return tx.Delete(&AuthModel{}, u.AuthID).Error
}
