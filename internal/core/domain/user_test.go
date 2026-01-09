package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
)

func TestNewUser(t *testing.T) {
	t.Run("creates user with valid firebase UID", func(t *testing.T) {
		user, err := domain.NewUser("firebase-uid-123")
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, "firebase-uid-123", user.FirebaseUID)
		assert.Equal(t, "en", user.PreferredLanguage)
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)
	})

	t.Run("fails with empty firebase UID", func(t *testing.T) {
		user, err := domain.NewUser("")
		require.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrInvalidFirebaseUID)
	})
}

func TestUserValidate(t *testing.T) {
	t.Run("valid user", func(t *testing.T) {
		user, _ := domain.NewUser("test-uid")
		err := user.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid with empty firebase UID", func(t *testing.T) {
		user := &domain.User{}
		err := user.Validate()
		require.Error(t, err)
	})

	t.Run("invalid preferred language", func(t *testing.T) {
		user, _ := domain.NewUser("test-uid")
		user.PreferredLanguage = "invalid"
		err := user.Validate()
		require.Error(t, err)
	})

	t.Run("valid pt-br language", func(t *testing.T) {
		user, _ := domain.NewUser("test-uid")
		user.PreferredLanguage = "pt-br"
		err := user.Validate()
		assert.NoError(t, err)
	})
}

func TestUserSetEmail(t *testing.T) {
	user, _ := domain.NewUser("test-uid")

	t.Run("sets email", func(t *testing.T) {
		user.SetEmail("test@example.com")
		require.NotNil(t, user.Email)
		assert.Equal(t, "test@example.com", *user.Email)
	})

	t.Run("clears email with empty string", func(t *testing.T) {
		user.SetEmail("")
		assert.Nil(t, user.Email)
	})
}

func TestUserSetName(t *testing.T) {
	user, _ := domain.NewUser("test-uid")

	t.Run("sets name", func(t *testing.T) {
		user.SetName("John Doe")
		require.NotNil(t, user.Name)
		assert.Equal(t, "John Doe", *user.Name)
	})

	t.Run("clears name with empty string", func(t *testing.T) {
		user.SetName("")
		assert.Nil(t, user.Name)
	})
}

func TestUserGetDisplayName(t *testing.T) {
	t.Run("returns name when set", func(t *testing.T) {
		user, _ := domain.NewUser("test-uid")
		user.SetName("John Doe")
		assert.Equal(t, "John Doe", user.GetDisplayName())
	})

	t.Run("returns email when name not set", func(t *testing.T) {
		user, _ := domain.NewUser("test-uid")
		user.SetEmail("john@example.com")
		assert.Equal(t, "john@example.com", user.GetDisplayName())
	})

	t.Run("returns Anonymous when nothing set", func(t *testing.T) {
		user, _ := domain.NewUser("test-uid")
		assert.Equal(t, "Anonymous User", user.GetDisplayName())
	})
}
