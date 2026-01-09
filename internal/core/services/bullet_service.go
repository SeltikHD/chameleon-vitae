// Package services contains the application services (use cases).
package services

import (
	"context"
	"fmt"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

// BulletService handles bullet-related use cases.
type BulletService struct {
	bulletRepo     ports.BulletRepository
	experienceRepo ports.ExperienceRepository
	aiProvider     ports.AIProvider
}

// NewBulletService creates a new BulletService with required dependencies.
func NewBulletService(
	bulletRepo ports.BulletRepository,
	experienceRepo ports.ExperienceRepository,
	aiProvider ports.AIProvider,
) *BulletService {
	return &BulletService{
		bulletRepo:     bulletRepo,
		experienceRepo: experienceRepo,
		aiProvider:     aiProvider,
	}
}

// CreateBulletRequest contains the parameters for creating a bullet.
type CreateBulletRequest struct {
	ExperienceID string
	Content      string
	ImpactScore  *int
	Keywords     []string
	DisplayOrder int
}

// CreateBullet creates a new bullet for an experience.
func (s *BulletService) CreateBullet(ctx context.Context, req CreateBulletRequest) (*domain.Bullet, error) {
	// Verify the experience exists.
	_, err := s.experienceRepo.GetByID(ctx, req.ExperienceID)
	if err != nil {
		return nil, fmt.Errorf("experience not found: %w", err)
	}

	// Create bullet.
	bullet, err := domain.NewBullet(req.ExperienceID, req.Content)
	if err != nil {
		return nil, err
	}

	bullet.DisplayOrder = req.DisplayOrder

	if req.ImpactScore != nil {
		if err := bullet.SetImpactScore(*req.ImpactScore); err != nil {
			return nil, err
		}
	}

	if len(req.Keywords) > 0 {
		bullet.SetKeywords(req.Keywords)
	}

	// Validate.
	if err := bullet.Validate(); err != nil {
		return nil, err
	}

	// Save.
	if err := s.bulletRepo.Create(ctx, bullet); err != nil {
		return nil, fmt.Errorf("failed to create bullet: %w", err)
	}

	return bullet, nil
}

// GetBullet retrieves a bullet by ID.
func (s *BulletService) GetBullet(ctx context.Context, bulletID string) (*domain.Bullet, error) {
	bullet, err := s.bulletRepo.GetByID(ctx, bulletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bullet: %w", err)
	}
	return bullet, nil
}

// ListBulletsByExperience lists all bullets for an experience.
func (s *BulletService) ListBulletsByExperience(ctx context.Context, experienceID string) ([]domain.Bullet, error) {
	bullets, err := s.bulletRepo.ListByExperienceID(ctx, experienceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list bullets: %w", err)
	}
	return bullets, nil
}

// ListBulletsByUser lists all bullets for a user across all experiences.
func (s *BulletService) ListBulletsByUser(ctx context.Context, userID string) ([]domain.Bullet, error) {
	bullets, err := s.bulletRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list bullets: %w", err)
	}
	return bullets, nil
}

// UpdateBulletRequest contains the parameters for updating a bullet.
type UpdateBulletRequest struct {
	BulletID     string
	Content      *string
	ImpactScore  *int
	Keywords     []string
	DisplayOrder *int
}

// UpdateBullet updates an existing bullet.
func (s *BulletService) UpdateBullet(ctx context.Context, req UpdateBulletRequest) (*domain.Bullet, error) {
	bullet, err := s.bulletRepo.GetByID(ctx, req.BulletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bullet: %w", err)
	}

	// Apply updates.
	if req.Content != nil {
		if err := bullet.UpdateContent(*req.Content); err != nil {
			return nil, err
		}
	}

	if req.ImpactScore != nil {
		if err := bullet.SetImpactScore(*req.ImpactScore); err != nil {
			return nil, err
		}
	}

	if req.Keywords != nil {
		bullet.SetKeywords(req.Keywords)
	}

	if req.DisplayOrder != nil {
		bullet.DisplayOrder = *req.DisplayOrder
	}

	// Validate.
	if err := bullet.Validate(); err != nil {
		return nil, err
	}

	// Save.
	if err := s.bulletRepo.Update(ctx, bullet); err != nil {
		return nil, fmt.Errorf("failed to update bullet: %w", err)
	}

	return bullet, nil
}

// DeleteBullet removes a bullet.
func (s *BulletService) DeleteBullet(ctx context.Context, bulletID string) error {
	if err := s.bulletRepo.Delete(ctx, bulletID); err != nil {
		return fmt.Errorf("failed to delete bullet: %w", err)
	}
	return nil
}

// SearchBulletsRequest contains the parameters for searching bullets.
type SearchBulletsRequest struct {
	UserID   string
	Keywords []string
}

// SearchBullets searches bullets by keywords.
func (s *BulletService) SearchBullets(ctx context.Context, req SearchBulletsRequest) ([]domain.Bullet, error) {
	bullets, err := s.bulletRepo.SearchByKeywords(ctx, req.UserID, req.Keywords)
	if err != nil {
		return nil, fmt.Errorf("failed to search bullets: %w", err)
	}
	return bullets, nil
}

// GetHighImpactBulletsRequest contains parameters for getting high-impact bullets.
type GetHighImpactBulletsRequest struct {
	UserID   string
	MinScore int
	Limit    int
}

// GetHighImpactBullets retrieves bullets with high impact scores.
func (s *BulletService) GetHighImpactBullets(ctx context.Context, req GetHighImpactBulletsRequest) ([]domain.Bullet, error) {
	if req.MinScore == 0 {
		req.MinScore = 70 // Default threshold.
	}
	if req.Limit == 0 {
		req.Limit = 20 // Default limit.
	}

	bullets, err := s.bulletRepo.GetHighImpactBullets(ctx, req.UserID, req.MinScore, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get high impact bullets: %w", err)
	}
	return bullets, nil
}

// AnalyzeBulletImpactRequest contains parameters for analyzing bullet impact.
type AnalyzeBulletImpactRequest struct {
	BulletID       string
	JobDescription string
}

// AnalyzeBulletImpact uses AI to analyze and score a bullet's impact.
func (s *BulletService) AnalyzeBulletImpact(ctx context.Context, req AnalyzeBulletImpactRequest) (*domain.Bullet, error) {
	bullet, err := s.bulletRepo.GetByID(ctx, req.BulletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bullet: %w", err)
	}

	// Analyze the job description to get context.
	jobAnalysis, err := s.aiProvider.AnalyzeJob(ctx, ports.AnalyzeJobRequest{
		JobDescription: req.JobDescription,
		TargetLanguage: "en",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to analyze job: %w", err)
	}

	// Score the bullet against the job.
	matchScore, err := s.aiProvider.ScoreMatch(ctx, ports.ScoreMatchRequest{
		JobAnalysis: jobAnalysis,
		Resume: &domain.ResumeContent{
			Experiences: []domain.TailoredExperience{
				{
					Bullets: []domain.TailoredBullet{
						{
							BulletID:        bullet.ID,
							OriginalContent: bullet.Content,
							TailoredContent: bullet.Content,
						},
					},
				},
			},
		},
		UserSkills: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to score bullet: %w", err)
	}

	// Update bullet with new score and keywords.
	if err := bullet.SetImpactScore(matchScore.Int()); err != nil {
		return nil, err
	}
	bullet.SetKeywords(jobAnalysis.Keywords)

	if err := s.bulletRepo.Update(ctx, bullet); err != nil {
		return nil, fmt.Errorf("failed to update bullet: %w", err)
	}

	return bullet, nil
}
