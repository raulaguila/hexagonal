package auditor

import (
	"context"
	"log/slog"
	"time"

	"github.com/godeh/sloggergo"
	"github.com/google/uuid"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/port/input"
	"github.com/raulaguila/go-api/internal/core/port/output"
)

type auditorUseCase struct {
	repo   output.AuditRepository
	logger *sloggergo.Logger
}

// NewAuditorUseCase creates a new auditor use case
func NewAuditorUseCase(repo output.AuditRepository, logger *sloggergo.Logger) input.AuditorUseCase {
	return &auditorUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Log records an action asynchronously
func (s *auditorUseCase) Log(actorID *uuid.UUID, action, resType, resID string, metadata map[string]any) {
	// Use a detached context with timeout to ensure we don't hang indefinitely
	bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
}
