package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/model"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/repository"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/raulaguila/go-api/internal/core/dto"
)

func TestRoleRepository_Integration(t *testing.T) {
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
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"unaccent\";")

	// Migrate schema
	err = db.AutoMigrate(&model.RoleModel{})
	require.NoError(t, err)

	repo := repository.NewRoleRepository(db)

	t.Run("Create and Find Role", func(t *testing.T) {
		role := entity.NewRole("Developer", []string{"code:read", "code:write"})

		// Create
		err := repo.Create(ctx, role)
		assert.NoError(t, err)
		assert.NotZero(t, role.ID)

		// FindByID
		found, err := repo.FindByID(ctx, role.ID)
		assert.NoError(t, err)
		assert.Equal(t, role.Name, found.Name)
		assert.ElementsMatch(t, role.Permissions, found.Permissions)

		// FindByName
		foundByName, err := repo.FindByName(ctx, role.Name)
		assert.NoError(t, err)
		assert.Equal(t, role.ID, foundByName.ID)
	})

	t.Run("Update Role", func(t *testing.T) {
		role := entity.NewRole("Manager", []string{"team:manage"})
		err := repo.Create(ctx, role)
		require.NoError(t, err)

		// Update
		role.UpdateName("Senior Manager", time.Now())
		role.UpdatePermissions([]string{"team:manage", "budget:approve"}, time.Now())

		err = repo.Update(ctx, role)
		assert.NoError(t, err)

		found, err := repo.FindByID(ctx, role.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Senior Manager", found.Name)
		assert.ElementsMatch(t, []string{"team:manage", "budget:approve"}, found.Permissions)
	})

	t.Run("Delete Role", func(t *testing.T) {
		role := entity.NewRole("Temp", []string{})
		err := repo.Create(ctx, role)
		require.NoError(t, err)

		err = repo.Delete(ctx, []uuid.UUID{role.ID})
		assert.NoError(t, err)

		found, err := repo.FindByID(ctx, role.ID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})

	t.Run("FindAll with Filters", func(t *testing.T) {
		// Clean table first for predictable results?
		// Or just search specific names.
		r1 := entity.NewRole("FilterA", []string{})
		r2 := entity.NewRole("FilterB", []string{})
		err = repo.Create(ctx, r1)
		assert.NoError(t, err)
		err = repo.Create(ctx, r2)
		assert.NoError(t, err)

		filter := &dto.RoleFilter{
			Filter: dto.Filter{
				Search: "Filter",
				Sort:   "name",
				Order:  "asc",
				Page:   1,
				Limit:  10,
			},
		}

		roles, err := repo.FindAll(ctx, filter)
		assert.NoError(t, err)
		// Might contain previous tests roles too, so just check existence
		assert.True(t, len(roles) >= 2)

		count, err := repo.Count(ctx, filter)
		assert.NoError(t, err)
		assert.True(t, count >= 2)
	})
}
