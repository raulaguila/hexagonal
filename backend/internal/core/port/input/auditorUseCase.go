package input

import "github.com/google/uuid"

// AuditorUseCase defines the input port for the auditor use case
type AuditorUseCase interface {
	Log(actorID *uuid.UUID, action, resType, resID string, metadata map[string]any)
}
