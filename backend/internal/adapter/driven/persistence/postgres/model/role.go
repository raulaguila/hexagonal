package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// RoleModel represents the database model for Role
type RoleModel struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	Name        string         `gorm:"column:name;type:citext;unique;not null;"`
	Permissions pq.StringArray `gorm:"column:permissions;type:citext[];not null;"`
	Enabled     bool           `gorm:"column:enabled;type:boolean;default:true;not null;"`
}

// TableName returns the table name for Role
func (RoleModel) TableName() string {
	return "usr_role"
}
