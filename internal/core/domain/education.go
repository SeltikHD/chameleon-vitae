// Package domain contains the core business entities and value objects.
package domain

import (
	"errors"
	"time"
)

// Education represents a formal education entry (university, college, certification).
type Education struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Institution  string    `json:"institution"`
	Degree       string    `json:"degree"`
	FieldOfStudy *string   `json:"field_of_study,omitempty"`
	Location     *string   `json:"location,omitempty"`
	StartDate    *Date     `json:"start_date,omitempty"`
	EndDate      *Date     `json:"end_date,omitempty"`
	GPA          *string   `json:"gpa,omitempty"`
	Honors       []string  `json:"honors,omitempty"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewEducation creates a new education entry with required fields.
func NewEducation(userID, institution, degree string) (*Education, error) {
	if userID == "" {
		return nil, ErrValidation
	}
	if institution == "" {
		return nil, ErrValidation
	}
	if degree == "" {
		return nil, ErrValidation
	}

	now := time.Now().UTC()
	return &Education{
		UserID:       userID,
		Institution:  institution,
		Degree:       degree,
		Honors:       make([]string, 0),
		DisplayOrder: 0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Validate validates the education entity.
func (e *Education) Validate() error {
	v := &ValidationErrors{}

	if e.UserID == "" {
		v.AddFieldError("user_id", "user ID is required")
	}

	if e.Institution == "" {
		v.AddFieldError("institution", "institution is required")
	}

	if e.Degree == "" {
		v.AddFieldError("degree", "degree is required")
	}

	// Validate date range if both dates are provided.
	if e.StartDate != nil && e.EndDate != nil {
		if !e.StartDate.IsZero() && !e.EndDate.IsZero() {
			if e.EndDate.Before(*e.StartDate) {
				v.AddFieldError("end_date", "end date must be after start date")
			}
		}
	}

	return v.ToError()
}

// SetDates sets the start and end dates for the education entry.
func (e *Education) SetDates(startDate, endDate *Date) error {
	if startDate != nil && endDate != nil {
		if !startDate.IsZero() && !endDate.IsZero() && endDate.Before(*startDate) {
			return ErrInvalidDateRange
		}
	}

	e.StartDate = startDate
	e.EndDate = endDate
	e.UpdatedAt = time.Now().UTC()
	return nil
}

// SetFieldOfStudy sets the field of study.
func (e *Education) SetFieldOfStudy(field string) {
	if field == "" {
		e.FieldOfStudy = nil
	} else {
		e.FieldOfStudy = &field
	}
	e.UpdatedAt = time.Now().UTC()
}

// SetLocation sets the location.
func (e *Education) SetLocation(location string) {
	if location == "" {
		e.Location = nil
	} else {
		e.Location = &location
	}
	e.UpdatedAt = time.Now().UTC()
}

// SetGPA sets the GPA.
func (e *Education) SetGPA(gpa string) {
	if gpa == "" {
		e.GPA = nil
	} else {
		e.GPA = &gpa
	}
	e.UpdatedAt = time.Now().UTC()
}

// AddHonor adds an honor to the education entry.
func (e *Education) AddHonor(honor string) {
	if honor == "" {
		return
	}
	e.Honors = append(e.Honors, honor)
	e.UpdatedAt = time.Now().UTC()
}

// DateRange returns a formatted date range string.
func (e *Education) DateRange() string {
	if e.StartDate == nil || e.StartDate.IsZero() {
		return ""
	}

	start := e.StartDate.Format("Jan. 2006")

	if e.EndDate == nil || e.EndDate.IsZero() {
		return start + " -- Present"
	}

	return start + " -- " + e.EndDate.Format("Jan. 2006")
}

// Error definitions for Education.
var (
	ErrEducationNotFound = errors.New("education entry not found")
)
