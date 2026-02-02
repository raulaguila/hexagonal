package output

import (
	"context"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
)

// AuditRepository defines the interface for audit log storage
type AuditRepository interface {
	Create(ctx context.Context, log *entity.AuditLog) error
}
