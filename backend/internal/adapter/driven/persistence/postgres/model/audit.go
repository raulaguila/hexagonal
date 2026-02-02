package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditModel represents the database model for AuditLog
type AuditModel struct {
	ID             uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ActorID        *uuid.UUID      `gorm:"type:uuid;index"` // Nullable
	Action         string          `gorm:"type:text;not null;index"`
	ResourceEntity string          `gorm:"type:text;not null;index"`
	ResourceID     string          `gorm:"type:text;not null;index"`
	Metadata       json.RawMessage `gorm:"type:jsonb"`
	IPAddress      string          `gorm:"type:text"`
	UserAgent      string          `gorm:"type:text"`
	CreatedAt      time.Time       `gorm:"autoCreateTime;index"`
}

// TableName returns the table name for AuditModel
func (AuditModel) TableName() string {
	return "sys_audit_logs"
}
