package auditor

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/port/output"
	"github.com/raulaguila/go-api/pkg/loggerx"
)

// Auditor defines the interface for the auditor service
type Auditor interface {
	Log(ctx context.Context, actorID *uuid.UUID, action, resType, resID string, metadata map[string]interface{})
}

type auditor struct {
	repo   output.AuditRepository
	logger *loggerx.Logger
}

// NewAuditor creates a new auditor service
func NewAuditor(repo output.AuditRepository, logger *loggerx.Logger) Auditor {
	return &auditor{
		repo:   repo,
		logger: logger,
	}
}

// Log records an action asynchronously
func (s *auditor) Log(ctx context.Context, actorID *uuid.UUID, action, resType, resID string, metadata map[string]any) {
	// Copy context values if needed, or use Background for async
	// In strict Hexagonal, we might want to use a specific context
	go func() {
		// Use a detached context with timeout to ensure we don't hang indefinitely but also don't get cancelled immediately if parent request finishes
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get IP and UserAgent from context if available?
		// For now, let's assume they are empty or passed in metadata?
		// Better design: Accept a Struct "AuditInput" instead of many args.
		// But sticking to the plan's simplicity. We will extract IP/UA from metadata if present.

		ip, _ := metadata["ip"].(string)
		ua, _ := metadata["user_agent"].(string)

		// Remove technical metadata fields from the JSON payload we store
		cleanMeta := make(map[string]any)
		for k, v := range metadata {
			if k != "ip" && k != "user_agent" {
				cleanMeta[k] = v
			}
		}

		log, err := entity.NewAuditLog(actorID, action, resType, resID, ip, ua, cleanMeta)
		if err != nil {
			s.logger.Error("failed to create audit log entity", slog.String("error", err.Error()))
			return
		}

		if err := s.repo.Create(bgCtx, log); err != nil {
			s.logger.Error("failed to persist audit log", slog.String("error", err.Error()))
		}
	}()
}
