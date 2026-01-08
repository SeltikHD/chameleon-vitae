// Package domain contains the core business entities and value objects.
// This package must have ZERO external dependencies - only standard library.
package domain

import "time"

// ExperienceType represents the type of professional experience.
type ExperienceType string

// Experience type constants.
const (
	ExperienceTypeWork              ExperienceType = "work"
	ExperienceTypeEducation         ExperienceType = "education"
	ExperienceTypeCertification     ExperienceType = "certification"
	ExperienceTypeProject           ExperienceType = "project"
	ExperienceTypeFreelance         ExperienceType = "freelance"
	ExperienceTypeVolunteer         ExperienceType = "volunteer"
	ExperienceTypeOpenSource        ExperienceType = "open_source"
	ExperienceTypeHackathon         ExperienceType = "hackathon"
	ExperienceTypeSideProject       ExperienceType = "side_project"
	ExperienceTypeEventOrganization ExperienceType = "event_organization"
	ExperienceTypePublication       ExperienceType = "publication"
	ExperienceTypeAward             ExperienceType = "award"
)

// ValidExperienceTypes returns all valid experience types.
func ValidExperienceTypes() []ExperienceType {
	return []ExperienceType{
		ExperienceTypeWork,
		ExperienceTypeEducation,
		ExperienceTypeCertification,
		ExperienceTypeProject,
		ExperienceTypeFreelance,
		ExperienceTypeVolunteer,
		ExperienceTypeOpenSource,
		ExperienceTypeHackathon,
		ExperienceTypeSideProject,
		ExperienceTypeEventOrganization,
		ExperienceTypePublication,
		ExperienceTypeAward,
	}
}

// IsValid checks if the experience type is valid.
func (t ExperienceType) IsValid() bool {
	for _, valid := range ValidExperienceTypes() {
		if t == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of the experience type.
func (t ExperienceType) String() string {
	return string(t)
}

// LanguageProficiency represents proficiency level in a spoken language.
type LanguageProficiency string

// Language proficiency constants.
const (
	ProficiencyNative       LanguageProficiency = "native"
	ProficiencyFluent       LanguageProficiency = "fluent"
	ProficiencyAdvanced     LanguageProficiency = "advanced"
	ProficiencyIntermediate LanguageProficiency = "intermediate"
	ProficiencyBasic        LanguageProficiency = "basic"
)

// ValidProficiencies returns all valid language proficiencies.
func ValidProficiencies() []LanguageProficiency {
	return []LanguageProficiency{
		ProficiencyNative,
		ProficiencyFluent,
		ProficiencyAdvanced,
		ProficiencyIntermediate,
		ProficiencyBasic,
	}
}

// IsValid checks if the proficiency is valid.
func (p LanguageProficiency) IsValid() bool {
	for _, valid := range ValidProficiencies() {
		if p == valid {
			return true
		}
	}
	return false
}

// ResumeStatus represents the lifecycle status of a generated resume.
type ResumeStatus string

// Resume status constants.
const (
	ResumeStatusDraft     ResumeStatus = "draft"
	ResumeStatusGenerated ResumeStatus = "generated"
	ResumeStatusReviewed  ResumeStatus = "reviewed"
	ResumeStatusSubmitted ResumeStatus = "submitted"
	ResumeStatusInterview ResumeStatus = "interview"
	ResumeStatusRejected  ResumeStatus = "rejected"
	ResumeStatusAccepted  ResumeStatus = "accepted"
)

// ValidResumeStatuses returns all valid resume statuses.
func ValidResumeStatuses() []ResumeStatus {
	return []ResumeStatus{
		ResumeStatusDraft,
		ResumeStatusGenerated,
		ResumeStatusReviewed,
		ResumeStatusSubmitted,
		ResumeStatusInterview,
		ResumeStatusRejected,
		ResumeStatusAccepted,
	}
}

// IsValid checks if the resume status is valid.
func (s ResumeStatus) IsValid() bool {
	for _, valid := range ValidResumeStatuses() {
		if s == valid {
			return true
		}
	}
	return false
}

// ImpactScore represents a bullet's impact score (0-100).
type ImpactScore int

// NewImpactScore creates a new impact score with validation.
func NewImpactScore(score int) (ImpactScore, error) {
	if score < 0 || score > 100 {
		return 0, ErrInvalidImpactScore
	}
	return ImpactScore(score), nil
}

// DefaultImpactScore returns the default impact score for new bullets.
func DefaultImpactScore() ImpactScore {
	return ImpactScore(50)
}

// Int returns the integer value of the impact score.
func (s ImpactScore) Int() int {
	return int(s)
}

// ProficiencyLevel represents a skill proficiency level (0-100).
type ProficiencyLevel int

// NewProficiencyLevel creates a new proficiency level with validation.
func NewProficiencyLevel(level int) (ProficiencyLevel, error) {
	if level < 0 || level > 100 {
		return 0, ErrInvalidProficiencyLevel
	}
	return ProficiencyLevel(level), nil
}

// Int returns the integer value of the proficiency level.
func (l ProficiencyLevel) Int() int {
	return int(l)
}

// MatchScore represents a resume's match score with a job description (0-100).
type MatchScore int

// NewMatchScore creates a new match score with validation.
func NewMatchScore(score int) (MatchScore, error) {
	if score < 0 || score > 100 {
		return 0, ErrInvalidMatchScore
	}
	return MatchScore(score), nil
}

// Int returns the integer value of the match score.
func (s MatchScore) Int() int {
	return int(s)
}

// Date represents a date without time information.
type Date struct {
	time.Time
}

// NewDate creates a new Date from year, month, and day.
func NewDate(year int, month time.Month, day int) Date {
	return Date{Time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// ParseDate parses a date string in YYYY-MM-DD format.
func ParseDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, ErrInvalidDateFormat
	}
	return Date{Time: t}, nil
}

// String returns the date in YYYY-MM-DD format.
func (d Date) String() string {
	return d.Time.Format("2006-01-02")
}

// IsZero reports whether the date is the zero value.
func (d Date) IsZero() bool {
	return d.Time.IsZero()
}

// Before reports whether d is before other.
func (d Date) Before(other Date) bool {
	return d.Time.Before(other.Time)
}

// After reports whether d is after other.
func (d Date) After(other Date) bool {
	return d.Time.After(other.Time)
}
