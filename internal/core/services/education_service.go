// Package services contains the application services (use cases).
package services

import (
	"context"
	"fmt"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

// EducationService handles education-related use cases.
type EducationService struct {
	educationRepo ports.EducationRepository
}

// NewEducationService creates a new EducationService with required dependencies.
func NewEducationService(educationRepo ports.EducationRepository) *EducationService {
	return &EducationService{
		educationRepo: educationRepo,
	}
}

// CreateEducationRequest contains the parameters for creating an education entry.
type CreateEducationRequest struct {
	UserID       string
	Institution  string
	Degree       string
	FieldOfStudy *string
	Location     *string
	StartDate    *domain.Date
	EndDate      *domain.Date
	GPA          *string
	Honors       []string
	DisplayOrder int
}

// CreateEducation creates a new education entry for a user.
func (s *EducationService) CreateEducation(ctx context.Context, req CreateEducationRequest) (*domain.Education, error) {
	// Create education.
	education, err := domain.NewEducation(req.UserID, req.Institution, req.Degree)
	if err != nil {
		return nil, err
	}

	// Set optional fields.
	if req.FieldOfStudy != nil {
		education.SetFieldOfStudy(*req.FieldOfStudy)
	}

	if req.Location != nil {
		education.SetLocation(*req.Location)
	}

	if req.StartDate != nil || req.EndDate != nil {
		if err := education.SetDates(req.StartDate, req.EndDate); err != nil {
			return nil, err
		}
	}

	if req.GPA != nil {
		education.SetGPA(*req.GPA)
	}

	for _, honor := range req.Honors {
		education.AddHonor(honor)
	}

	education.DisplayOrder = req.DisplayOrder

	// Validate.
	if err := education.Validate(); err != nil {
		return nil, err
	}

	// Save.
	if err := s.educationRepo.Create(ctx, education); err != nil {
		return nil, fmt.Errorf("failed to create education: %w", err)
	}

	return education, nil
}

// GetEducation retrieves an education entry by ID.
func (s *EducationService) GetEducation(ctx context.Context, educationID string) (*domain.Education, error) {
	education, err := s.educationRepo.GetByID(ctx, educationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get education: %w", err)
	}
	return education, nil
}

// ListEducationRequest contains parameters for listing education entries.
type ListEducationRequest struct {
	UserID string
}

// ListEducation lists all education entries for a user.
func (s *EducationService) ListEducation(ctx context.Context, req ListEducationRequest) ([]domain.Education, error) {
	education, err := s.educationRepo.ListByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to list education: %w", err)
	}
	return education, nil
}

// UpdateEducationRequest contains parameters for updating an education entry.
type UpdateEducationRequest struct {
	EducationID  string
	Institution  *string
	Degree       *string
	FieldOfStudy *string
	Location     *string
	StartDate    *domain.Date
	EndDate      *domain.Date
	GPA          *string
	Honors       []string
	DisplayOrder *int
}

// UpdateEducation updates an existing education entry.
func (s *EducationService) UpdateEducation(ctx context.Context, req UpdateEducationRequest) (*domain.Education, error) {
	// Get existing.
	education, err := s.educationRepo.GetByID(ctx, req.EducationID)
	if err != nil {
		return nil, err
	}

	// Update fields.
	if req.Institution != nil {
		education.Institution = *req.Institution
	}

	if req.Degree != nil {
		education.Degree = *req.Degree
	}

	if req.FieldOfStudy != nil {
		education.SetFieldOfStudy(*req.FieldOfStudy)
	}

	if req.Location != nil {
		education.SetLocation(*req.Location)
	}

	if req.StartDate != nil || req.EndDate != nil {
		startDate := education.StartDate
		endDate := education.EndDate

		if req.StartDate != nil {
			startDate = req.StartDate
		}
		if req.EndDate != nil {
			endDate = req.EndDate
		}

		if err := education.SetDates(startDate, endDate); err != nil {
			return nil, err
		}
	}

	if req.GPA != nil {
		education.SetGPA(*req.GPA)
	}

	if req.Honors != nil {
		education.Honors = req.Honors
	}

	if req.DisplayOrder != nil {
		education.DisplayOrder = *req.DisplayOrder
	}

	// Validate.
	if err := education.Validate(); err != nil {
		return nil, err
	}

	// Save.
	if err := s.educationRepo.Update(ctx, education); err != nil {
		return nil, fmt.Errorf("failed to update education: %w", err)
	}

	return education, nil
}

// DeleteEducation removes an education entry.
func (s *EducationService) DeleteEducation(ctx context.Context, educationID string) error {
	if err := s.educationRepo.Delete(ctx, educationID); err != nil {
		return fmt.Errorf("failed to delete education: %w", err)
	}
	return nil
}

// UpdateEducationOrderRequest contains parameters for updating education display order.
type UpdateEducationOrderRequest struct {
	Orders []ports.DisplayOrderUpdate
}

// UpdateEducationOrder updates the display order of education entries.
func (s *EducationService) UpdateEducationOrder(ctx context.Context, req UpdateEducationOrderRequest) error {
	if err := s.educationRepo.UpdateDisplayOrder(ctx, req.Orders); err != nil {
		return fmt.Errorf("failed to update education order: %w", err)
	}
	return nil
}
