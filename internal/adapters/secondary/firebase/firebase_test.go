package firebase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SeltikHD/chameleon-vitae/internal/adapters/secondary/firebase"
)

func TestNew(t *testing.T) {
	ctx := context.Background()

	t.Run("requires project ID", func(t *testing.T) {
		cfg := firebase.Config{}
		adapter, err := firebase.New(ctx, cfg)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.ErrorIs(t, err, firebase.ErrMissingProjectID)
	})

	t.Run("fails with invalid credentials file", func(t *testing.T) {
		cfg := firebase.Config{
			ProjectID:       "test-project",
			CredentialsFile: "/non/existent/path.json",
		}
		adapter, err := firebase.New(ctx, cfg)
		require.Error(t, err)
		assert.Nil(t, adapter)
	})

	t.Run("fails with invalid credentials JSON", func(t *testing.T) {
		cfg := firebase.Config{
			ProjectID:       "test-project",
			CredentialsJSON: []byte("invalid json"),
		}
		adapter, err := firebase.New(ctx, cfg)
		require.Error(t, err)
		assert.Nil(t, adapter)
	})
}

func TestErrors(t *testing.T) {
	t.Run("ErrInvalidToken is defined", func(t *testing.T) {
		assert.NotNil(t, firebase.ErrInvalidToken)
		assert.Error(t, firebase.ErrInvalidToken)
	})

	t.Run("ErrTokenExpired is defined", func(t *testing.T) {
		assert.NotNil(t, firebase.ErrTokenExpired)
		assert.Error(t, firebase.ErrTokenExpired)
	})

	t.Run("ErrMissingProjectID is defined", func(t *testing.T) {
		assert.NotNil(t, firebase.ErrMissingProjectID)
		assert.Error(t, firebase.ErrMissingProjectID)
	})
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	t.Skip("Integration tests require Firebase credentials")
}
