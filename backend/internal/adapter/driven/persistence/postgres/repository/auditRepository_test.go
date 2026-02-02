package repository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/model"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/repository"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	container, connStr, err := setupPostgresContainer(ctx)
	require.NoError(t, err)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Connect to DB
	db := postgres.MustConnect(&postgres.Config{Url: connStr})

	// Enable extensions
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"citext\";")
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";")

	// Migrate schema
	err = db.AutoMigrate(&model.AuditModel{})
	require.NoError(t, err)

	repo := repository.NewAuditRepository(db)

	t.Run("Create Audit Log", func(t *testing.T) {
		actorID := uuid.New()
		resourceID := uuid.New()
		meta := map[string]string{"key": "value"}
		metaJSON, _ := json.Marshal(meta)

		log := &entity.AuditLog{
			ID:             uuid.New(),
			ActorID:        &actorID,
			Action:         "CREATE",
			ResourceEntity: "USER",
			ResourceID:     resourceID.String(),
			Metadata:       metaJSON,
			IPAddress:      "127.0.0.1",
			UserAgent:      "Go-Test-Client",
			CreatedAt:      time.Now(),
		}

		err := repo.Create(ctx, log)
		assert.NoError(t, err)

		// Verify direct DB persistence since repository might not have FindByID (Audits are usually append-only or searched via OLAP/logs)
		// But let's check if we can query it using GORM direct
		var saved model.AuditModel
		err = db.First(&saved, "id = ?", log.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, log.Action, saved.Action)
		assert.Equal(t, log.IPAddress, saved.IPAddress)
	})
}
