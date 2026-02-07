package auditor_test

import (
	"context"
	"testing"
	"time"

	"github.com/godeh/sloggergo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/usecase/auditor"
)

// MockAuditRepo
type MockAuditRepo struct {
	mock.Mock
}

func (m *MockAuditRepo) Create(ctx context.Context, log *entity.AuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func TestAuditorUseCase_Log(t *testing.T) {
	mockRepo := new(MockAuditRepo)
	logger := sloggergo.New() // default logger

	service := auditor.NewAuditorUseCase(mockRepo, logger)

	actorID := uuid.New()
	metadata := map[string]interface{}{
		"ip":         "127.0.0.1",
		"user_agent": "Go-Test",
		"detail":     "test",
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(l *entity.AuditLog) bool {
		return l.Action == "TEST_ACTION" && l.IPAddress == "127.0.0.1" && l.ResourceEntity == "USER"
	})).Return(nil)

	service.Log(&actorID, "TEST_ACTION", "USER", "123", metadata)

	time.Sleep(100 * time.Millisecond) // Wait for goroutine

	mockRepo.AssertExpectations(t)
}
