//go:build integration

package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SeltikHD/chameleon-vitae/internal/adapters/secondary/postgres"
	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
)

var testDB *postgres.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	cfg := postgres.Config{
		Host:     getEnv("TEST_DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("TEST_DB_USER", "chameleon"),
		Password: getEnv("TEST_DB_PASSWORD", "chameleon_secret"), // pragma: allowlist secret
		Database: getEnv("TEST_DB_NAME", "chameleon_test"),
		SSLMode:  "disable",
		MaxConns: 5,
		MinConns: 1,
	}

	db, err := postgres.New(ctx, cfg)
	if err != nil {
		os.Exit(0)
	}
	testDB = db

	code := m.Run()
	testDB.Close()
	os.Exit(code)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func TestUserRepositoryCRUD(t *testing.T) {
	ctx := context.Background()
	repo := testDB.UserRepository()
	firebaseUID := "test-firebase-" + time.Now().Format("20060102150405")

	t.Run("Create", func(t *testing.T) {
		user, err := domain.NewUser(firebaseUID)
		require.NoError(t, err)
		user.SetName("Test User")
		user.SetEmail("test@example.com")
		err = repo.Create(ctx, user)
		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)
	})

	t.Run("GetByFirebaseUID", func(t *testing.T) {
		user, err := repo.GetByFirebaseUID(ctx, firebaseUID)
		require.NoError(t, err)
		assert.Equal(t, firebaseUID, user.FirebaseUID)
		assert.Equal(t, "Test User", *user.Name)
	})

	t.Run("GetByID", func(t *testing.T) {
		user, err := repo.GetByFirebaseUID(ctx, firebaseUID)
		require.NoError(t, err)
		fetchedUser, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, fetchedUser.ID)
	})

	t.Run("Update", func(t *testing.T) {
		user, err := repo.GetByFirebaseUID(ctx, firebaseUID)
		require.NoError(t, err)
		user.SetName("Updated Name")
		err = repo.Update(ctx, user)
		require.NoError(t, err)
		updated, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", *updated.Name)
	})

	t.Run("Delete", func(t *testing.T) {
		user, err := repo.GetByFirebaseUID(ctx, firebaseUID)
		require.NoError(t, err)
		err = repo.Delete(ctx, user.ID)
		require.NoError(t, err)
		_, err = repo.GetByID(ctx, user.ID)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})
}

func TestUserRepositoryUpsert(t *testing.T) {
	ctx := context.Background()
	repo := testDB.UserRepository()
	firebaseUID := "test-upsert-" + time.Now().Format("20060102150405")

	t.Run("Upsert_Create", func(t *testing.T) {
		user, err := domain.NewUser(firebaseUID)
		require.NoError(t, err)
		user.SetName("Upsert User")
		err = repo.Upsert(ctx, user)
		require.NoError(t, err)
		created, err := repo.GetByFirebaseUID(ctx, firebaseUID)
		require.NoError(t, err)
		assert.Equal(t, "Upsert User", *created.Name)
	})

	t.Run("Upsert_Update", func(t *testing.T) {
		user, err := repo.GetByFirebaseUID(ctx, firebaseUID)
		require.NoError(t, err)
		user.SetName("Updated Upsert User")
		err = repo.Upsert(ctx, user)
		require.NoError(t, err)
		updated, err := repo.GetByFirebaseUID(ctx, firebaseUID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Upsert User", *updated.Name)
	})

	user, _ := repo.GetByFirebaseUID(ctx, firebaseUID)
	if user != nil {
		_ = repo.Delete(ctx, user.ID)
	}
}

func TestExperienceRepositoryCRUD(t *testing.T) {
	ctx := context.Background()
	userRepo := testDB.UserRepository()
	expRepo := testDB.ExperienceRepository()

	firebaseUID := "test-exp-user-" + time.Now().Format("20060102150405")
	user, err := domain.NewUser(firebaseUID)
	require.NoError(t, err)
	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	defer func() {
		_ = userRepo.Delete(ctx, user.ID)
	}()

	var experienceID string

	t.Run("Create", func(t *testing.T) {
		startDate := domain.NewDate(2020, 1, 1)
		exp, err := domain.NewExperience(
			user.ID,
			domain.ExperienceTypeWork,
			"Software Engineer",
			"Test Company",
			startDate,
		)
		require.NoError(t, err)
		err = expRepo.Create(ctx, exp)
		require.NoError(t, err)
		assert.NotEmpty(t, exp.ID)
		experienceID = exp.ID
	})

	t.Run("GetByID", func(t *testing.T) {
		exp, err := expRepo.GetByID(ctx, experienceID)
		require.NoError(t, err)
		assert.Equal(t, "Software Engineer", exp.Title)
		assert.Equal(t, "Test Company", exp.Organization)
	})

	t.Run("ListByUserID", func(t *testing.T) {
		experiences, total, err := expRepo.ListByUserID(ctx, user.ID, nil)
		require.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Len(t, experiences, 1)
	})

	t.Run("Update", func(t *testing.T) {
		exp, err := expRepo.GetByID(ctx, experienceID)
		require.NoError(t, err)
		exp.Title = "Senior Software Engineer"
		err = expRepo.Update(ctx, exp)
		require.NoError(t, err)
		updated, err := expRepo.GetByID(ctx, experienceID)
		require.NoError(t, err)
		assert.Equal(t, "Senior Software Engineer", updated.Title)
	})

	t.Run("Delete", func(t *testing.T) {
		err := expRepo.Delete(ctx, experienceID)
		require.NoError(t, err)
		_, err = expRepo.GetByID(ctx, experienceID)
		assert.ErrorIs(t, err, domain.ErrExperienceNotFound)
	})
}

func TestDBHealthCheck(t *testing.T) {
	ctx := context.Background()
	err := testDB.HealthCheck(ctx)
	require.NoError(t, err)
}
