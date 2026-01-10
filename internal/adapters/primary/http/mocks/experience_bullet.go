// Package mocks provides mock implementations of port interfaces for testing.
package mocks

import (
	"context"
	"sync"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

// InMemoryExperienceRepository is an in-memory mock implementation of ExperienceRepository.
type InMemoryExperienceRepository struct {
	mu          sync.RWMutex
	experiences map[string]*domain.Experience
}

// NewInMemoryExperienceRepository creates a new in-memory experience repository.
func NewInMemoryExperienceRepository() *InMemoryExperienceRepository {
	return &InMemoryExperienceRepository{
		experiences: make(map[string]*domain.Experience),
	}
}

// Create creates a new experience.
func (r *InMemoryExperienceRepository) Create(ctx context.Context, experience *domain.Experience) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	clone := *experience
	r.experiences[experience.ID] = &clone
	return nil
}

// GetByID retrieves an experience by ID.
func (r *InMemoryExperienceRepository) GetByID(ctx context.Context, id string) (*domain.Experience, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	exp, exists := r.experiences[id]
	if !exists {
		return nil, domain.ErrExperienceNotFound
	}

	clone := *exp
	return &clone, nil
}

// GetByIDWithBullets retrieves an experience with all its bullets.
func (r *InMemoryExperienceRepository) GetByIDWithBullets(ctx context.Context, id string) (*domain.Experience, error) {
	return r.GetByID(ctx, id)
}

// ListByUserIDWithBullets lists all experiences for a user.
func (r *InMemoryExperienceRepository) ListByUserIDWithBullets(ctx context.Context, userID string, opts ports.ListOptions) ([]domain.Experience, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.Experience
	for _, exp := range r.experiences {
		if exp.UserID == userID {
			clone := *exp
			result = append(result, clone)
		}
	}

	return result, len(result), nil
}

// ListByUserIDAndTypeWithBullets lists experiences filtered by type.
func (r *InMemoryExperienceRepository) ListByUserIDAndTypeWithBullets(ctx context.Context, userID string, expType domain.ExperienceType, opts ports.ListOptions) ([]domain.Experience, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.Experience
	for _, exp := range r.experiences {
		if exp.UserID == userID && exp.Type == expType {
			clone := *exp
			result = append(result, clone)
		}
	}

	return result, len(result), nil
}

// Update updates an existing experience.
func (r *InMemoryExperienceRepository) Update(ctx context.Context, experience *domain.Experience) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.experiences[experience.ID]; !exists {
		return domain.ErrExperienceNotFound
	}

	clone := *experience
	r.experiences[experience.ID] = &clone
	return nil
}

// Delete removes an experience and all its bullets.
func (r *InMemoryExperienceRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.experiences[id]; !exists {
		return domain.ErrExperienceNotFound
	}

	delete(r.experiences, id)
	return nil
}

// UpdateDisplayOrder updates the display order of experiences.
func (r *InMemoryExperienceRepository) UpdateDisplayOrder(ctx context.Context, orders []ports.DisplayOrderUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, order := range orders {
		if exp, exists := r.experiences[order.ID]; exists {
			exp.DisplayOrder = order.DisplayOrder
		}
	}

	return nil
}

// Seed adds experiences for testing.
func (r *InMemoryExperienceRepository) Seed(experiences ...*domain.Experience) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, exp := range experiences {
		clone := *exp
		r.experiences[exp.ID] = &clone
	}
}

// Reset clears all data.
func (r *InMemoryExperienceRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.experiences = make(map[string]*domain.Experience)
}

// Verify interface compliance.
var _ ports.ExperienceRepository = (*InMemoryExperienceRepository)(nil)

// InMemoryBulletRepository is an in-memory mock implementation of BulletRepository.
type InMemoryBulletRepository struct {
	mu      sync.RWMutex
	bullets map[string]*domain.Bullet
}

// NewInMemoryBulletRepository creates a new in-memory bullet repository.
func NewInMemoryBulletRepository() *InMemoryBulletRepository {
	return &InMemoryBulletRepository{
		bullets: make(map[string]*domain.Bullet),
	}
}

// Create creates a new bullet.
func (r *InMemoryBulletRepository) Create(ctx context.Context, bullet *domain.Bullet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	clone := *bullet
	r.bullets[bullet.ID] = &clone
	return nil
}

// GetByID retrieves a bullet by ID.
func (r *InMemoryBulletRepository) GetByID(ctx context.Context, id string) (*domain.Bullet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bullet, exists := r.bullets[id]
	if !exists {
		return nil, domain.ErrBulletNotFound
	}

	clone := *bullet
	return &clone, nil
}

// ListByExperienceID lists all bullets for an experience.
func (r *InMemoryBulletRepository) ListByExperienceID(ctx context.Context, experienceID string) ([]domain.Bullet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.Bullet
	for _, bullet := range r.bullets {
		if bullet.ExperienceID == experienceID {
			clone := *bullet
			result = append(result, clone)
		}
	}

	return result, nil
}

// ListByIDs retrieves multiple bullets by their IDs.
func (r *InMemoryBulletRepository) ListByIDs(ctx context.Context, ids []string) ([]domain.Bullet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.Bullet
	for _, id := range ids {
		if bullet, exists := r.bullets[id]; exists {
			clone := *bullet
			result = append(result, clone)
		}
	}

	return result, nil
}

// ListByUserID lists all bullets for a user (across all experiences).
func (r *InMemoryBulletRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Bullet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// This mock returns all bullets - in reality it would filter by user
	var result []domain.Bullet
	for _, bullet := range r.bullets {
		clone := *bullet
		result = append(result, clone)
	}

	return result, nil
}

// Update updates an existing bullet.
func (r *InMemoryBulletRepository) Update(ctx context.Context, bullet *domain.Bullet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.bullets[bullet.ID]; !exists {
		return domain.ErrBulletNotFound
	}

	clone := *bullet
	r.bullets[bullet.ID] = &clone
	return nil
}

// Delete removes a bullet.
func (r *InMemoryBulletRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.bullets[id]; !exists {
		return domain.ErrBulletNotFound
	}

	delete(r.bullets, id)
	return nil
}

// SearchByKeywords searches bullets by keywords.
func (r *InMemoryBulletRepository) SearchByKeywords(ctx context.Context, userID string, keywords []string) ([]domain.Bullet, error) {
	return r.ListByUserID(ctx, userID)
}

// GetHighImpactBullets retrieves bullets with impact score >= threshold.
func (r *InMemoryBulletRepository) GetHighImpactBullets(ctx context.Context, userID string, minScore int, limit int) ([]domain.Bullet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.Bullet
	for _, bullet := range r.bullets {
		if bullet.ImpactScore.Int() >= minScore {
			clone := *bullet
			result = append(result, clone)
			if len(result) >= limit {
				break
			}
		}
	}

	return result, nil
}

// Seed adds bullets for testing.
func (r *InMemoryBulletRepository) Seed(bullets ...*domain.Bullet) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, bullet := range bullets {
		clone := *bullet
		r.bullets[bullet.ID] = &clone
	}
}

// Reset clears all data.
func (r *InMemoryBulletRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bullets = make(map[string]*domain.Bullet)
}

// Verify interface compliance.
var _ ports.BulletRepository = (*InMemoryBulletRepository)(nil)
