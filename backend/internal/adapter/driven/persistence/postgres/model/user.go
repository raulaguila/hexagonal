package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserModel represents the database model for User
type UserModel struct {
	ID        uuid.UUID    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time    `gorm:"autoCreateTime"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime"`
	Name      string       `gorm:"column:name;type:citext;not null;"`
	Username  string       `gorm:"column:username;type:citext;not null;"`
	Email     string       `gorm:"column:mail;type:citext;not null;"`
	AuthID    uuid.UUID    `gorm:"column:auth_id;type:uuid;not null;"`
	Auth      *AuthModel   `gorm:"constraint:OnDelete:CASCADE"`
	Roles     []*RoleModel `gorm:"many2many:usr_user_role;joinForeignKey:user_id;joinReferences:role_id"`
}

// TableName returns the table name for User
func (UserModel) TableName() string {
	return "usr_user"
}

// AfterDelete hook to delete associated Auth
func (u *UserModel) AfterDelete(tx *gorm.DB) error {
	return tx.Delete(&AuthModel{}, u.AuthID).Error
}
