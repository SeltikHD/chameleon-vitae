// Package ports defines the interfaces (ports) that adapters must implement.
package ports

import (
	"context"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
)

// UserRepository defines the interface for user persistence operations.
type UserRepository interface {
	// Create creates a new user in the database.
	Create(ctx context.Context, user *domain.User) error

	// GetByID retrieves a user by their internal ID.
	GetByID(ctx context.Context, id string) (*domain.User, error)

	// GetByFirebaseUID retrieves a user by their Firebase UID.
	GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error)

	// Update updates an existing user.
	Update(ctx context.Context, user *domain.User) error

	// Delete removes a user from the database.
	Delete(ctx context.Context, id string) error

	// Upsert creates or updates a user based on Firebase UID.
	Upsert(ctx context.Context, user *domain.User) error
}

// ExperienceRepository defines the interface for experience persistence operations.
type ExperienceRepository interface {
	// Create creates a new experience.
	Create(ctx context.Context, experience *domain.Experience) error

	// GetByID retrieves an experience by ID.
	GetByID(ctx context.Context, id string) (*domain.Experience, error)

	// GetByIDWithBullets retrieves an experience with all its bullets.
	GetByIDWithBullets(ctx context.Context, id string) (*domain.Experience, error)

	// ListByUserID lists all experiences for a user.
	ListByUserID(ctx context.Context, userID string, opts ListOptions) ([]domain.Experience, int, error)

	// ListByUserIDAndType lists experiences filtered by type.
	ListByUserIDAndType(ctx context.Context, userID string, expType domain.ExperienceType, opts ListOptions) ([]domain.Experience, int, error)

	// Update updates an existing experience.
	Update(ctx context.Context, experience *domain.Experience) error

	// Delete removes an experience and all its bullets.
	Delete(ctx context.Context, id string) error

	// UpdateDisplayOrder updates the display order of experiences.
	UpdateDisplayOrder(ctx context.Context, orders []DisplayOrderUpdate) error
}

// BulletRepository defines the interface for bullet persistence operations.
type BulletRepository interface {
	// Create creates a new bullet.
	Create(ctx context.Context, bullet *domain.Bullet) error

	// GetByID retrieves a bullet by ID.
	GetByID(ctx context.Context, id string) (*domain.Bullet, error)

	// ListByExperienceID lists all bullets for an experience.
	ListByExperienceID(ctx context.Context, experienceID string) ([]domain.Bullet, error)

	// ListByIDs retrieves multiple bullets by their IDs.
	ListByIDs(ctx context.Context, ids []string) ([]domain.Bullet, error)

	// ListByUserID lists all bullets for a user (across all experiences).
	ListByUserID(ctx context.Context, userID string) ([]domain.Bullet, error)

	// Update updates an existing bullet.
	Update(ctx context.Context, bullet *domain.Bullet) error

	// Delete removes a bullet.
	Delete(ctx context.Context, id string) error

	// SearchByKeywords searches bullets by keywords.
	SearchByKeywords(ctx context.Context, userID string, keywords []string) ([]domain.Bullet, error)

	// GetHighImpactBullets retrieves bullets with impact score >= threshold.
	GetHighImpactBullets(ctx context.Context, userID string, minScore int, limit int) ([]domain.Bullet, error)
}

// SkillRepository defines the interface for skill persistence operations.
type SkillRepository interface {
	// Create creates a new skill.
	Create(ctx context.Context, skill *domain.Skill) error

	// GetByID retrieves a skill by ID.
	GetByID(ctx context.Context, id string) (*domain.Skill, error)

	// GetByUserIDAndName retrieves a skill by user ID and name.
	GetByUserIDAndName(ctx context.Context, userID, name string) (*domain.Skill, error)

	// ListByUserID lists all skills for a user.
	ListByUserID(ctx context.Context, userID string) ([]domain.Skill, error)

	// ListByUserIDAndCategory lists skills filtered by category.
	ListByUserIDAndCategory(ctx context.Context, userID, category string) ([]domain.Skill, error)

	// ListHighlighted lists highlighted skills for a user.
	ListHighlighted(ctx context.Context, userID string) ([]domain.Skill, error)

	// Update updates an existing skill.
	Update(ctx context.Context, skill *domain.Skill) error

	// Upsert creates or updates a skill based on user ID and name.
	Upsert(ctx context.Context, skill *domain.Skill) error

	// BatchUpsert creates or updates multiple skills.
	BatchUpsert(ctx context.Context, skills []domain.Skill) (created int, updated int, err error)

	// Delete removes a skill.
	Delete(ctx context.Context, id string) error

	// SearchByName searches skills by name (fuzzy match).
	SearchByName(ctx context.Context, userID, query string) ([]domain.Skill, error)
}

// SpokenLanguageRepository defines the interface for spoken language persistence.
type SpokenLanguageRepository interface {
	// Create creates a new spoken language.
	Create(ctx context.Context, language *domain.SpokenLanguage) error

	// GetByID retrieves a spoken language by ID.
	GetByID(ctx context.Context, id string) (*domain.SpokenLanguage, error)

	// ListByUserID lists all spoken languages for a user.
	ListByUserID(ctx context.Context, userID string) ([]domain.SpokenLanguage, error)

	// Update updates an existing spoken language.
	Update(ctx context.Context, language *domain.SpokenLanguage) error

	// Delete removes a spoken language.
	Delete(ctx context.Context, id string) error
}

// ResumeRepository defines the interface for resume persistence operations.
type ResumeRepository interface {
	// Create creates a new resume.
	Create(ctx context.Context, resume *domain.Resume) error

	// GetByID retrieves a resume by ID.
	GetByID(ctx context.Context, id string) (*domain.Resume, error)

	// ListByUserID lists all resumes for a user.
	ListByUserID(ctx context.Context, userID string, opts ListOptions) ([]domain.Resume, int, error)

	// ListByUserIDAndStatus lists resumes filtered by status.
	ListByUserIDAndStatus(ctx context.Context, userID string, status domain.ResumeStatus, opts ListOptions) ([]domain.Resume, int, error)

	// Update updates an existing resume.
	Update(ctx context.Context, resume *domain.Resume) error

	// Delete removes a resume.
	Delete(ctx context.Context, id string) error
}

// ListOptions contains pagination and filtering options.
type ListOptions struct {
	Limit  int
	Offset int
}

// DefaultListOptions returns default list options.
func DefaultListOptions() ListOptions {
	return ListOptions{
		Limit:  50,
		Offset: 0,
	}
}

// DisplayOrderUpdate represents a display order update for an entity.
type DisplayOrderUpdate struct {
	ID           string
	DisplayOrder int
}
