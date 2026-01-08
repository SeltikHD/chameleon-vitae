// Package domain contains the core business entities and value objects.
package domain

import "time"

// Skill represents a user's skill with category and proficiency.
type Skill struct {
	ID                string           `json:"id"`
	UserID            string           `json:"user_id"`
	Name              string           `json:"name"`
	Category          *string          `json:"category,omitempty"`
	ProficiencyLevel  ProficiencyLevel `json:"proficiency_level"`
	YearsOfExperience *float64         `json:"years_of_experience,omitempty"`
	IsHighlighted     bool             `json:"is_highlighted"`
	DisplayOrder      int              `json:"display_order"`
	CreatedAt         time.Time        `json:"created_at"`
}

// NewSkill creates a new skill with required fields.
func NewSkill(userID, name string) (*Skill, error) {
	if name == "" {
		return nil, ErrEmptySkillName
	}

	level, _ := NewProficiencyLevel(50) // Default to 50%

	return &Skill{
		UserID:           userID,
		Name:             name,
		ProficiencyLevel: level,
		IsHighlighted:    false,
		CreatedAt:        time.Now().UTC(),
	}, nil
}

// Validate validates the skill entity.
func (s *Skill) Validate() error {
	v := &ValidationErrors{}

	if s.UserID == "" {
		v.AddFieldError("user_id", "user ID is required")
	}

	if s.Name == "" {
		v.AddFieldError("name", "name is required")
	}

	if s.ProficiencyLevel.Int() < 0 || s.ProficiencyLevel.Int() > 100 {
		v.AddFieldError("proficiency_level", "must be between 0 and 100")
	}

	if s.YearsOfExperience != nil && *s.YearsOfExperience < 0 {
		v.AddFieldError("years_of_experience", "cannot be negative")
	}

	return v.ToError()
}

// SetProficiency sets the proficiency level after validation.
func (s *Skill) SetProficiency(level int) error {
	proficiency, err := NewProficiencyLevel(level)
	if err != nil {
		return err
	}
	s.ProficiencyLevel = proficiency
	return nil
}

// SetCategory sets the skill category.
func (s *Skill) SetCategory(category string) {
	if category == "" {
		s.Category = nil
	} else {
		s.Category = &category
	}
}

// SetYearsOfExperience sets the years of experience.
func (s *Skill) SetYearsOfExperience(years float64) {
	if years <= 0 {
		s.YearsOfExperience = nil
	} else {
		s.YearsOfExperience = &years
	}
}

// Highlight marks the skill as highlighted.
func (s *Skill) Highlight() {
	s.IsHighlighted = true
}

// Unhighlight removes the highlighted status.
func (s *Skill) Unhighlight() {
	s.IsHighlighted = false
}

// IsExpert returns true if proficiency is >= 80.
func (s *Skill) IsExpert() bool {
	return s.ProficiencyLevel.Int() >= 80
}

// IsBeginner returns true if proficiency is < 30.
func (s *Skill) IsBeginner() bool {
	return s.ProficiencyLevel.Int() < 30
}

// SpokenLanguage represents a natural language spoken by the user.
type SpokenLanguage struct {
	ID           string              `json:"id"`
	UserID       string              `json:"user_id"`
	Language     string              `json:"language"`
	Proficiency  LanguageProficiency `json:"proficiency"`
	DisplayOrder int                 `json:"display_order"`
	CreatedAt    time.Time           `json:"created_at"`
}

// NewSpokenLanguage creates a new spoken language with required fields.
func NewSpokenLanguage(userID, language string, proficiency LanguageProficiency) (*SpokenLanguage, error) {
	if language == "" {
		return nil, NewFieldError(ErrValidation, "language", "language is required")
	}

	if !proficiency.IsValid() {
		return nil, ErrInvalidProficiency
	}

	return &SpokenLanguage{
		UserID:      userID,
		Language:    language,
		Proficiency: proficiency,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

// Validate validates the spoken language entity.
func (l *SpokenLanguage) Validate() error {
	v := &ValidationErrors{}

	if l.UserID == "" {
		v.AddFieldError("user_id", "user ID is required")
	}

	if l.Language == "" {
		v.AddFieldError("language", "language is required")
	}

	if !l.Proficiency.IsValid() {
		v.AddFieldError("proficiency", "invalid proficiency level")
	}

	return v.ToError()
}

// IsNative returns true if the language is native.
func (l *SpokenLanguage) IsNative() bool {
	return l.Proficiency == ProficiencyNative
}

// IsFluent returns true if the language is fluent or native.
func (l *SpokenLanguage) IsFluent() bool {
	return l.Proficiency == ProficiencyNative || l.Proficiency == ProficiencyFluent
}
