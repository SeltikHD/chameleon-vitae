// Package services contains the application services (use cases).
package services

import (
	"context"
	"fmt"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

// SkillService handles skill and spoken language-related use cases.
type SkillService struct {
	skillRepo    ports.SkillRepository
	languageRepo ports.SpokenLanguageRepository
}

// NewSkillService creates a new SkillService with required dependencies.
func NewSkillService(
	skillRepo ports.SkillRepository,
	languageRepo ports.SpokenLanguageRepository,
) *SkillService {
	return &SkillService{
		skillRepo:    skillRepo,
		languageRepo: languageRepo,
	}
}

// CreateSkillRequest contains the parameters for creating a skill.
type CreateSkillRequest struct {
	UserID            string
	Name              string
	Category          *string
	ProficiencyLevel  int
	YearsOfExperience *float64
	IsHighlighted     bool
	DisplayOrder      int
}

// CreateSkill creates a new skill for a user.
func (s *SkillService) CreateSkill(ctx context.Context, req CreateSkillRequest) (*domain.Skill, error) {
	// Check if skill already exists.
	existing, _ := s.skillRepo.GetByUserIDAndName(ctx, req.UserID, req.Name)
	if existing != nil {
		return nil, domain.ErrSkillAlreadyExists
	}

	// Create skill.
	skill, err := domain.NewSkill(req.UserID, req.Name)
	if err != nil {
		return nil, err
	}

	if req.Category != nil {
		skill.SetCategory(*req.Category)
	}

	if req.ProficiencyLevel > 0 {
		if err := skill.SetProficiency(req.ProficiencyLevel); err != nil {
			return nil, err
		}
	}

	if req.YearsOfExperience != nil {
		skill.SetYearsOfExperience(*req.YearsOfExperience)
	}

	if req.IsHighlighted {
		skill.Highlight()
	}

	skill.DisplayOrder = req.DisplayOrder

	// Validate.
	if err := skill.Validate(); err != nil {
		return nil, err
	}

	// Save.
	if err := s.skillRepo.Create(ctx, skill); err != nil {
		return nil, fmt.Errorf("failed to create skill: %w", err)
	}

	return skill, nil
}

// GetSkill retrieves a skill by ID.
func (s *SkillService) GetSkill(ctx context.Context, skillID string) (*domain.Skill, error) {
	skill, err := s.skillRepo.GetByID(ctx, skillID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}
	return skill, nil
}

// ListSkillsRequest contains parameters for listing skills.
type ListSkillsRequest struct {
	UserID   string
	Category *string
}

// ListSkills lists all skills for a user with optional category filter.
func (s *SkillService) ListSkills(ctx context.Context, req ListSkillsRequest) ([]domain.Skill, error) {
	var skills []domain.Skill
	var err error

	if req.Category != nil && *req.Category != "" {
		skills, err = s.skillRepo.ListByUserIDAndCategory(ctx, req.UserID, *req.Category)
	} else {
		skills, err = s.skillRepo.ListByUserID(ctx, req.UserID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list skills: %w", err)
	}
	return skills, nil
}

// ListHighlightedSkills lists highlighted skills for a user.
func (s *SkillService) ListHighlightedSkills(ctx context.Context, userID string) ([]domain.Skill, error) {
	skills, err := s.skillRepo.ListHighlighted(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list highlighted skills: %w", err)
	}
	return skills, nil
}

// UpdateSkillRequest contains the parameters for updating a skill.
type UpdateSkillRequest struct {
	SkillID           string
	Name              *string
	Category          *string
	ProficiencyLevel  *int
	YearsOfExperience *float64
	IsHighlighted     *bool
	DisplayOrder      *int
}

// UpdateSkill updates an existing skill.
func (s *SkillService) UpdateSkill(ctx context.Context, req UpdateSkillRequest) (*domain.Skill, error) {
	skill, err := s.skillRepo.GetByID(ctx, req.SkillID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}

	// Apply updates.
	if req.Name != nil {
		skill.Name = *req.Name
	}

	if req.Category != nil {
		skill.SetCategory(*req.Category)
	}

	if req.ProficiencyLevel != nil {
		if err := skill.SetProficiency(*req.ProficiencyLevel); err != nil {
			return nil, err
		}
	}

	if req.YearsOfExperience != nil {
		skill.SetYearsOfExperience(*req.YearsOfExperience)
	}

	if req.IsHighlighted != nil {
		if *req.IsHighlighted {
			skill.Highlight()
		} else {
			skill.Unhighlight()
		}
	}

	if req.DisplayOrder != nil {
		skill.DisplayOrder = *req.DisplayOrder
	}

	// Validate.
	if err := skill.Validate(); err != nil {
		return nil, err
	}

	// Save.
	if err := s.skillRepo.Update(ctx, skill); err != nil {
		return nil, fmt.Errorf("failed to update skill: %w", err)
	}

	return skill, nil
}

// BatchUpsertSkillsRequest contains parameters for batch upserting skills.
type BatchUpsertSkillsRequest struct {
	UserID string
	Skills []CreateSkillRequest
}

// BatchUpsertSkillsResponse contains the result of batch upsert.
type BatchUpsertSkillsResponse struct {
	Created int
	Updated int
}

// BatchUpsertSkills creates or updates multiple skills at once.
func (s *SkillService) BatchUpsertSkills(ctx context.Context, req BatchUpsertSkillsRequest) (*BatchUpsertSkillsResponse, error) {
	skills := make([]domain.Skill, 0, len(req.Skills))

	for _, skillReq := range req.Skills {
		skill, err := domain.NewSkill(req.UserID, skillReq.Name)
		if err != nil {
			return nil, err
		}

		if skillReq.Category != nil {
			skill.SetCategory(*skillReq.Category)
		}

		if skillReq.ProficiencyLevel > 0 {
			if err := skill.SetProficiency(skillReq.ProficiencyLevel); err != nil {
				return nil, err
			}
		}

		if skillReq.YearsOfExperience != nil {
			skill.SetYearsOfExperience(*skillReq.YearsOfExperience)
		}

		if skillReq.IsHighlighted {
			skill.Highlight()
		}

		skill.DisplayOrder = skillReq.DisplayOrder

		if err := skill.Validate(); err != nil {
			return nil, err
		}

		skills = append(skills, *skill)
	}

	created, updated, err := s.skillRepo.BatchUpsert(ctx, skills)
	if err != nil {
		return nil, fmt.Errorf("failed to batch upsert skills: %w", err)
	}

	return &BatchUpsertSkillsResponse{
		Created: created,
		Updated: updated,
	}, nil
}

// DeleteSkill removes a skill.
func (s *SkillService) DeleteSkill(ctx context.Context, skillID string) error {
	if err := s.skillRepo.Delete(ctx, skillID); err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}
	return nil
}

// SearchSkills searches skills by name.
func (s *SkillService) SearchSkills(ctx context.Context, userID, query string) ([]domain.Skill, error) {
	skills, err := s.skillRepo.SearchByName(ctx, userID, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search skills: %w", err)
	}
	return skills, nil
}

// ===== Spoken Languages =====

// CreateSpokenLanguageRequest contains parameters for creating a spoken language.
type CreateSpokenLanguageRequest struct {
	UserID       string
	Language     string
	Proficiency  string
	DisplayOrder int
}

// CreateSpokenLanguage creates a new spoken language for a user.
func (s *SkillService) CreateSpokenLanguage(ctx context.Context, req CreateSpokenLanguageRequest) (*domain.SpokenLanguage, error) {
	proficiency, err := domain.ParseLanguageProficiency(req.Proficiency)
	if err != nil {
		return nil, err
	}

	language, err := domain.NewSpokenLanguage(req.UserID, req.Language, proficiency)
	if err != nil {
		return nil, err
	}

	language.DisplayOrder = req.DisplayOrder

	if err := language.Validate(); err != nil {
		return nil, err
	}

	if err := s.languageRepo.Create(ctx, language); err != nil {
		return nil, fmt.Errorf("failed to create spoken language: %w", err)
	}

	return language, nil
}

// ListSpokenLanguages lists all spoken languages for a user.
func (s *SkillService) ListSpokenLanguages(ctx context.Context, userID string) ([]domain.SpokenLanguage, error) {
	languages, err := s.languageRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list spoken languages: %w", err)
	}
	return languages, nil
}

// UpdateSpokenLanguageRequest contains parameters for updating a spoken language.
type UpdateSpokenLanguageRequest struct {
	LanguageID   string
	Language     *string
	Proficiency  *string
	DisplayOrder *int
}

// UpdateSpokenLanguage updates an existing spoken language.
func (s *SkillService) UpdateSpokenLanguage(ctx context.Context, req UpdateSpokenLanguageRequest) (*domain.SpokenLanguage, error) {
	language, err := s.languageRepo.GetByID(ctx, req.LanguageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spoken language: %w", err)
	}

	if req.Language != nil {
		language.Language = *req.Language
	}

	if req.Proficiency != nil {
		proficiency, err := domain.ParseLanguageProficiency(*req.Proficiency)
		if err != nil {
			return nil, err
		}
		language.Proficiency = proficiency
	}

	if req.DisplayOrder != nil {
		language.DisplayOrder = *req.DisplayOrder
	}

	if err := language.Validate(); err != nil {
		return nil, err
	}

	if err := s.languageRepo.Update(ctx, language); err != nil {
		return nil, fmt.Errorf("failed to update spoken language: %w", err)
	}

	return language, nil
}

// DeleteSpokenLanguage removes a spoken language.
func (s *SkillService) DeleteSpokenLanguage(ctx context.Context, languageID string) error {
	if err := s.languageRepo.Delete(ctx, languageID); err != nil {
		return fmt.Errorf("failed to delete spoken language: %w", err)
	}
	return nil
}
