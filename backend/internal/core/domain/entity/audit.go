package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID             uuid.UUID
	ActorID        *uuid.UUID // Nullable for system actions
	Action         string     // e.g., "USER_CREATE", "ROLE_UPDATE"
	ResourceEntity string     // e.g., "USER", "ROLE"
	ResourceID     string     // ID of the resource
	Metadata       json.RawMessage
	IPAddress      string
	UserAgent      string
	CreatedAt      time.Time
}

// NewAuditLog creates a new audit log entry
func NewAuditLog(actorID *uuid.UUID, action, entity, resourceID, ip, ua string, metadata map[string]interface{}) (*AuditLog, error) {
	metaJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	return &AuditLog{
		ID:             uuid.New(),
		ActorID:        actorID,
		Action:         action,
		ResourceEntity: entity,
		ResourceID:     resourceID,
		Metadata:       metaJSON,
		IPAddress:      ip,
		UserAgent:      ua,
		CreatedAt:      time.Now(),
	}, nil
}
