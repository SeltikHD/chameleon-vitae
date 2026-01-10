// Package services contains the application services (use cases).
package services

import (
	"context"
	"fmt"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

// ExperienceService handles experience-related use cases.
type ExperienceService struct {
	experienceRepo ports.ExperienceRepository
	bulletRepo     ports.BulletRepository
}

// NewExperienceService creates a new ExperienceService with required dependencies.
func NewExperienceService(
	experienceRepo ports.ExperienceRepository,
	bulletRepo ports.BulletRepository,
) *ExperienceService {
	return &ExperienceService{
		experienceRepo: experienceRepo,
		bulletRepo:     bulletRepo,
	}
}

// CreateExperienceRequest contains the parameters for creating an experience.
type CreateExperienceRequest struct {
	UserID       string
	Type         string
	Title        string
	Organization string
	Location     *string
	StartDate    string
	EndDate      *string
	IsCurrent    bool
	Description  *string
	URL          *string
	DisplayOrder int
}

// CreateExperience creates a new experience entry.
func (s *ExperienceService) CreateExperience(ctx context.Context, req CreateExperienceRequest) (*domain.Experience, error) {
	// Parse experience type.
	expType, err := domain.ParseExperienceType(req.Type)
	if err != nil {
		return nil, err
	}

	// Parse start date.
	startDate, err := domain.ParseDate(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	// Create experience.
	experience, err := domain.NewExperience(req.UserID, expType, req.Title, req.Organization, startDate)
	if err != nil {
		return nil, err
	}

	// Set optional fields.
	experience.Location = req.Location
	experience.Description = req.Description
	experience.URL = req.URL
	experience.DisplayOrder = req.DisplayOrder

	if req.IsCurrent {
		experience.MarkAsCurrent()
	} else if req.EndDate != nil && *req.EndDate != "" {
		endDate, err := domain.ParseDate(*req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %w", err)
		}
		if err := experience.SetEndDate(&endDate); err != nil {
			return nil, err
		}
	}

	// Validate.
	if err := experience.Validate(); err != nil {
		return nil, err
	}

	// Save.
	if err := s.experienceRepo.Create(ctx, experience); err != nil {
		return nil, fmt.Errorf("failed to create experience: %w", err)
	}

	return experience, nil
}

// GetExperience retrieves an experience by ID.
func (s *ExperienceService) GetExperience(ctx context.Context, experienceID string) (*domain.Experience, error) {
	experience, err := s.experienceRepo.GetByID(ctx, experienceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get experience: %w", err)
	}
	return experience, nil
}

// GetExperienceWithBullets retrieves an experience with all its bullets.
func (s *ExperienceService) GetExperienceWithBullets(ctx context.Context, experienceID string) (*domain.Experience, error) {
	experience, err := s.experienceRepo.GetByIDWithBullets(ctx, experienceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get experience with bullets: %w", err)
	}
	return experience, nil
}

// ListExperiencesRequest contains the parameters for listing experiences.
type ListExperiencesRequest struct {
	UserID string
	Type   *string
	Limit  int
	Offset int
}

// ListExperiencesResponse contains the result of listing experiences.
type ListExperiencesResponse struct {
	Experiences []domain.Experience
	Total       int
}

// ListExperiences lists experiences for a user with optional type filter.
func (s *ExperienceService) ListExperiences(ctx context.Context, req ListExperiencesRequest) (*ListExperiencesResponse, error) {
	opts := ports.ListOptions{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	if opts.Limit == 0 {
		opts = ports.DefaultListOptions()
	}

	var experiences []domain.Experience
	var total int
	var err error

	if req.Type != nil && *req.Type != "" {
		expType, parseErr := domain.ParseExperienceType(*req.Type)
		if parseErr != nil {
			return nil, parseErr
		}
		experiences, total, err = s.experienceRepo.ListByUserIDAndTypeWithBullets(ctx, req.UserID, expType, opts)
	} else {
		experiences, total, err = s.experienceRepo.ListByUserIDWithBullets(ctx, req.UserID, opts)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list experiences: %w", err)
	}

	return &ListExperiencesResponse{
		Experiences: experiences,
		Total:       total,
	}, nil
}

// UpdateExperienceRequest contains the parameters for updating an experience.
type UpdateExperienceRequest struct {
	ExperienceID string
	Type         *string
	Title        *string
	Organization *string
	Location     *string
	StartDate    *string
	EndDate      *string
	IsCurrent    *bool
	Description  *string
	URL          *string
	DisplayOrder *int
}

// UpdateExperience updates an existing experience.
func (s *ExperienceService) UpdateExperience(ctx context.Context, req UpdateExperienceRequest) (*domain.Experience, error) {
	experience, err := s.experienceRepo.GetByID(ctx, req.ExperienceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get experience: %w", err)
	}

	// Apply updates.
	if req.Type != nil {
		expType, err := domain.ParseExperienceType(*req.Type)
		if err != nil {
			return nil, err
		}
		experience.Type = expType
	}

	if req.Title != nil {
		experience.Title = *req.Title
	}

	if req.Organization != nil {
		experience.Organization = *req.Organization
	}

	if req.Location != nil {
		if *req.Location == "" {
			experience.Location = nil
		} else {
			experience.Location = req.Location
		}
	}

	if req.StartDate != nil {
		startDate, err := domain.ParseDate(*req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date: %w", err)
		}
		experience.StartDate = startDate
	}

	if req.IsCurrent != nil && *req.IsCurrent {
		experience.MarkAsCurrent()
	} else if req.EndDate != nil {
		if *req.EndDate == "" {
			experience.EndDate = nil
		} else {
			endDate, err := domain.ParseDate(*req.EndDate)
			if err != nil {
				return nil, fmt.Errorf("invalid end date: %w", err)
			}
			if err := experience.SetEndDate(&endDate); err != nil {
				return nil, err
			}
		}
	}

	if req.Description != nil {
		if *req.Description == "" {
			experience.Description = nil
		} else {
			experience.Description = req.Description
		}
	}

	if req.URL != nil {
		if *req.URL == "" {
			experience.URL = nil
		} else {
			experience.URL = req.URL
		}
	}

	if req.DisplayOrder != nil {
		experience.DisplayOrder = *req.DisplayOrder
	}

	// Validate.
	if err := experience.Validate(); err != nil {
		return nil, err
	}

	// Save.
	if err := s.experienceRepo.Update(ctx, experience); err != nil {
		return nil, fmt.Errorf("failed to update experience: %w", err)
	}

	return experience, nil
}

// DeleteExperience removes an experience and all its bullets.
func (s *ExperienceService) DeleteExperience(ctx context.Context, experienceID string) error {
	if err := s.experienceRepo.Delete(ctx, experienceID); err != nil {
		return fmt.Errorf("failed to delete experience: %w", err)
	}
	return nil
}

// ReorderExperiencesRequest contains the new order for experiences.
type ReorderExperiencesRequest struct {
	Orders []ports.DisplayOrderUpdate
}

// ReorderExperiences updates the display order of multiple experiences.
func (s *ExperienceService) ReorderExperiences(ctx context.Context, req ReorderExperiencesRequest) error {
	if err := s.experienceRepo.UpdateDisplayOrder(ctx, req.Orders); err != nil {
		return fmt.Errorf("failed to reorder experiences: %w", err)
	}
	return nil
}
