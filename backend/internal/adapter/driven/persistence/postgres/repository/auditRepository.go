package repository

import (
	"context"

	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/model"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/port/output"
	"github.com/raulaguila/go-api/pkg/apperror"
	"gorm.io/gorm"
)

type auditRepository struct {
	db *gorm.DB
}

// NewAuditRepository creates a new audit repository
func NewAuditRepository(db *gorm.DB) output.AuditRepository {
	return &auditRepository{
		db: db,
	}
}

// Create persists an audit log
func (r *auditRepository) Create(ctx context.Context, log *entity.AuditLog) error {
	m := &model.AuditModel{
		ID:             log.ID,
		ActorID:        log.ActorID,
		Action:         log.Action,
		ResourceEntity: log.ResourceEntity,
		ResourceID:     log.ResourceID,
		Metadata:       log.Metadata, // Assuming log.Metadata is json.RawMessage (which is []byte)
		IPAddress:      log.IPAddress,
		UserAgent:      log.UserAgent,
		CreatedAt:      log.CreatedAt,
	}

	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return apperror.Internal("failed to create audit log", err)
	}

	return nil
}
