package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
)

func TestNewExperience(t *testing.T) {
	startDate := domain.NewDate(2020, 1, 15)

	t.Run("creates experience with valid fields", func(t *testing.T) {
		exp, err := domain.NewExperience(
			"user-123",
			domain.ExperienceTypeWork,
			"Software Engineer",
			"Tech Corp",
			startDate,
		)
		require.NoError(t, err)
		require.NotNil(t, exp)
		assert.Equal(t, "user-123", exp.UserID)
		assert.Equal(t, domain.ExperienceTypeWork, exp.Type)
		assert.Equal(t, "Software Engineer", exp.Title)
		assert.Equal(t, "Tech Corp", exp.Organization)
		assert.Equal(t, startDate, exp.StartDate)
		assert.False(t, exp.IsCurrent)
		assert.NotZero(t, exp.CreatedAt)
	})

	t.Run("fails with invalid experience type", func(t *testing.T) {
		exp, err := domain.NewExperience(
			"user-123",
			domain.ExperienceType("invalid"),
			"Title",
			"Org",
			startDate,
		)
		require.Error(t, err)
		assert.Nil(t, exp)
		assert.ErrorIs(t, err, domain.ErrInvalidExperienceType)
	})
}

func TestExperienceValidate(t *testing.T) {
	startDate := domain.NewDate(2020, 1, 15)

	t.Run("valid experience", func(t *testing.T) {
		exp, _ := domain.NewExperience("user-123", domain.ExperienceTypeWork, "Title", "Org", startDate)
		err := exp.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid without user ID", func(t *testing.T) {
		exp := &domain.Experience{
			Type:         domain.ExperienceTypeWork,
			Title:        "Title",
			Organization: "Org",
			StartDate:    startDate,
		}
		err := exp.Validate()
		require.Error(t, err)
	})

	t.Run("invalid without title", func(t *testing.T) {
		exp := &domain.Experience{
			UserID:       "user-123",
			Type:         domain.ExperienceTypeWork,
			Organization: "Org",
			StartDate:    startDate,
		}
		err := exp.Validate()
		require.Error(t, err)
	})
}

func TestExperienceSetEndDate(t *testing.T) {
	startDate := domain.NewDate(2020, 1, 15)
	exp, _ := domain.NewExperience("user-123", domain.ExperienceTypeWork, "Title", "Org", startDate)

	t.Run("sets valid end date", func(t *testing.T) {
		endDate := domain.NewDate(2022, 12, 31)
		err := exp.SetEndDate(&endDate)
		require.NoError(t, err)
		assert.False(t, exp.IsCurrent)
	})

	t.Run("rejects end date before start date", func(t *testing.T) {
		endDate := domain.NewDate(2019, 1, 1)
		err := exp.SetEndDate(&endDate)
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidDateRange)
	})
}

func TestExperienceTypeIsValid(t *testing.T) {
	t.Run("valid types", func(t *testing.T) {
		validTypes := domain.ValidExperienceTypes()
		for _, expType := range validTypes {
			assert.True(t, expType.IsValid(), "expected %s to be valid", expType)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		invalid := domain.ExperienceType("not_a_type")
		assert.False(t, invalid.IsValid())
	})
}

func TestParseExperienceType(t *testing.T) {
	t.Run("parses valid type", func(t *testing.T) {
		expType, err := domain.ParseExperienceType("work")
		require.NoError(t, err)
		assert.Equal(t, domain.ExperienceTypeWork, expType)
	})

	t.Run("fails for invalid type", func(t *testing.T) {
		_, err := domain.ParseExperienceType("invalid")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidExperienceType)
	})
}
