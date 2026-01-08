// Package domain contains the core business entities and value objects.
package domain

import "time"

// Experience represents a professional experience entry (work, education, project, etc.).
type Experience struct {
	ID           string         `json:"id"`
	UserID       string         `json:"user_id"`
	Type         ExperienceType `json:"type"`
	Title        string         `json:"title"`
	Organization string         `json:"organization"`
	Location     *string        `json:"location,omitempty"`
	StartDate    Date           `json:"start_date"`
	EndDate      *Date          `json:"end_date,omitempty"`
	IsCurrent    bool           `json:"is_current"`
	Description  *string        `json:"description,omitempty"`
	URL          *string        `json:"url,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	DisplayOrder int            `json:"display_order"`
	Bullets      []Bullet       `json:"bullets,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// NewExperience creates a new experience with required fields.
func NewExperience(userID string, expType ExperienceType, title, organization string, startDate Date) (*Experience, error) {
	if !expType.IsValid() {
		return nil, ErrInvalidExperienceType
	}

	now := time.Now().UTC()
	return &Experience{
		UserID:       userID,
		Type:         expType,
		Title:        title,
		Organization: organization,
		StartDate:    startDate,
		IsCurrent:    false,
		Metadata:     make(map[string]any),
		Bullets:      make([]Bullet, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Validate validates the experience entity.
func (e *Experience) Validate() error {
	v := &ValidationErrors{}

	if e.UserID == "" {
		v.AddFieldError("user_id", "user ID is required")
	}

	if !e.Type.IsValid() {
		v.AddFieldError("type", "invalid experience type")
	}

	if e.Title == "" {
		v.AddFieldError("title", "title is required")
	}

	if e.Organization == "" {
		v.AddFieldError("organization", "organization is required")
	}

	if e.StartDate.IsZero() {
		v.AddFieldError("start_date", "start date is required")
	}

	// Validate date range.
	if e.EndDate != nil && !e.EndDate.IsZero() {
		if e.EndDate.Before(e.StartDate) {
			v.AddFieldError("end_date", "end date must be after start date")
		}
	}

	// Current experience should not have an end date.
	if e.IsCurrent && e.EndDate != nil && !e.EndDate.IsZero() {
		v.AddFieldError("is_current", "current experience cannot have an end date")
	}

	return v.ToError()
}

// SetEndDate sets the end date and updates is_current accordingly.
func (e *Experience) SetEndDate(endDate *Date) error {
	if endDate != nil && !endDate.IsZero() && endDate.Before(e.StartDate) {
		return ErrInvalidDateRange
	}

	e.EndDate = endDate
	if endDate != nil && !endDate.IsZero() {
		e.IsCurrent = false
	}
	e.UpdatedAt = time.Now().UTC()
	return nil
}

// MarkAsCurrent marks the experience as current (ongoing).
func (e *Experience) MarkAsCurrent() {
	e.IsCurrent = true
	e.EndDate = nil
	e.UpdatedAt = time.Now().UTC()
}

// AddBullet adds a new bullet to the experience.
func (e *Experience) AddBullet(bullet Bullet) {
	bullet.ExperienceID = e.ID
	e.Bullets = append(e.Bullets, bullet)
	e.UpdatedAt = time.Now().UTC()
}

// Duration returns the duration of the experience in months.
// Returns -1 for current/ongoing experiences.
func (e *Experience) Duration() int {
	if e.IsCurrent || e.EndDate == nil || e.EndDate.IsZero() {
		return -1
	}

	months := 0
	start := e.StartDate.Time
	end := e.EndDate.Time

	years := end.Year() - start.Year()
	months = int(end.Month() - start.Month())
	months += years * 12

	if months < 0 {
		return 0
	}
	return months
}
