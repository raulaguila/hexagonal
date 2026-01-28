package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/model"
	"github.com/raulaguila/go-api/internal/adapter/driven/persistence/postgres/repository"
	"github.com/raulaguila/go-api/internal/core/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupPostgresContainer starts a Postgres container for testing
func setupPostgresContainer(ctx context.Context) (*tcpostgres.PostgresContainer, string, error) {
	dbName := "testdb"
	dbUser := "user"
	dbPassword := "password"

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase(dbName),
		tcpostgres.WithUsername(dbUser),
		tcpostgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		return nil, "", err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, "", err
	}

	return pgContainer, connStr, nil
}

func TestUserRepository_Integration(t *testing.T) {
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
	db := postgres.MustConnect(&postgres.Config{Dsn: connStr})

	// Migrate schema using Models
	err = db.AutoMigrate(&model.UserModel{}, &model.AuthModel{}, &model.ProfileModel{})
	require.NoError(t, err)

	// Seed required data (Profiles) using Model
	profile := &model.ProfileModel{Name: "Admin", Permissions: []string{"all"}}
	err = db.Create(profile).Error
	require.NoError(t, err)

	repo := repository.NewUserRepository(db)

	t.Run("Create and Find User", func(t *testing.T) {
		auth, _ := entity.NewAuth(profile.ID, true)
		users, _ := entity.NewUser("John Doe", "johndoe", "john@test.com", auth)

		// Test Create
		err := repo.Create(ctx, users)
		assert.NoError(t, err)
		assert.NotZero(t, users.ID)

		// Test FindByID
		found, err := repo.FindByID(ctx, users.ID)
		assert.NoError(t, err)
		assert.Equal(t, users.Email, found.Email)
		assert.Equal(t, users.Name, found.Name)
	})
}
