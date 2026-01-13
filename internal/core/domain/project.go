// Package domain contains the core business entities and value objects.
package domain

import (
	"errors"
	"strings"
	"time"
)

// Project represents a side project or personal work.
type Project struct {
	ID            string          `json:"id"`
	UserID        string          `json:"user_id"`
	Name          string          `json:"name"`
	Description   *string         `json:"description,omitempty"`
	TechStack     []string        `json:"tech_stack"`
	URL           *string         `json:"url,omitempty"`
	RepositoryURL *string         `json:"repository_url,omitempty"`
	StartDate     *Date           `json:"start_date,omitempty"`
	EndDate       *Date           `json:"end_date,omitempty"`
	DisplayOrder  int             `json:"display_order"`
	Bullets       []ProjectBullet `json:"bullets,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// ProjectBullet represents an achievement/feature bullet for a project.
type ProjectBullet struct {
	ID           string    `json:"id"`
	ProjectID    string    `json:"project_id"`
	Content      string    `json:"content"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewProject creates a new project with required fields.
func NewProject(userID, name string, techStack []string) (*Project, error) {
	if userID == "" {
		return nil, ErrValidation
	}
	if name == "" {
		return nil, ErrValidation
	}

	if techStack == nil {
		techStack = make([]string, 0)
	}

	now := time.Now().UTC()
	return &Project{
		UserID:       userID,
		Name:         name,
		TechStack:    techStack,
		Bullets:      make([]ProjectBullet, 0),
		DisplayOrder: 0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Validate validates the project entity.
func (p *Project) Validate() error {
	v := &ValidationErrors{}

	if p.UserID == "" {
		v.AddFieldError("user_id", "user ID is required")
	}

	if p.Name == "" {
		v.AddFieldError("name", "project name is required")
	}

	// Validate date range if both dates are provided.
	if p.StartDate != nil && p.EndDate != nil {
		if !p.StartDate.IsZero() && !p.EndDate.IsZero() {
			if p.EndDate.Before(*p.StartDate) {
				v.AddFieldError("end_date", "end date must be after start date")
			}
		}
	}

	return v.ToError()
}

// SetDates sets the start and end dates for the project.
func (p *Project) SetDates(startDate, endDate *Date) error {
	if startDate != nil && endDate != nil {
		if !startDate.IsZero() && !endDate.IsZero() && endDate.Before(*startDate) {
			return ErrInvalidDateRange
		}
	}

	p.StartDate = startDate
	p.EndDate = endDate
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// SetDescription sets the project description.
func (p *Project) SetDescription(description string) {
	if description == "" {
		p.Description = nil
	} else {
		p.Description = &description
	}
	p.UpdatedAt = time.Now().UTC()
}

// SetURL sets the project URL.
func (p *Project) SetURL(url string) {
	if url == "" {
		p.URL = nil
	} else {
		p.URL = &url
	}
	p.UpdatedAt = time.Now().UTC()
}

// SetRepositoryURL sets the repository URL.
func (p *Project) SetRepositoryURL(url string) {
	if url == "" {
		p.RepositoryURL = nil
	} else {
		p.RepositoryURL = &url
	}
	p.UpdatedAt = time.Now().UTC()
}

// SetTechStack sets the tech stack.
func (p *Project) SetTechStack(techStack []string) {
	if techStack == nil {
		p.TechStack = make([]string, 0)
	} else {
		p.TechStack = techStack
	}
	p.UpdatedAt = time.Now().UTC()
}

// AddBullet adds a bullet to the project.
func (p *Project) AddBullet(bullet ProjectBullet) {
	bullet.ProjectID = p.ID
	p.Bullets = append(p.Bullets, bullet)
	p.UpdatedAt = time.Now().UTC()
}

// TechStackString returns the tech stack as a comma-separated string.
func (p *Project) TechStackString() string {
	if len(p.TechStack) == 0 {
		return ""
	}

	var result strings.Builder
	for i, tech := range p.TechStack {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(tech)
	}
	return result.String()
}

// DateRange returns a formatted date range string.
func (p *Project) DateRange() string {
	if p.StartDate == nil || p.StartDate.IsZero() {
		return ""
	}

	start := p.StartDate.Format("Jan. 2006")

	if p.EndDate == nil || p.EndDate.IsZero() {
		return start + " -- Present"
	}

	return start + " -- " + p.EndDate.Format("Jan. 2006")
}

// NewProjectBullet creates a new project bullet.
func NewProjectBullet(projectID, content string) (*ProjectBullet, error) {
	if projectID == "" {
		return nil, ErrValidation
	}
	if content == "" {
		return nil, ErrValidation
	}

	now := time.Now().UTC()
	return &ProjectBullet{
		ProjectID:    projectID,
		Content:      content,
		DisplayOrder: 0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Error definitions for Project.
var (
	ErrProjectNotFound       = errors.New("project not found")
	ErrProjectBulletNotFound = errors.New("project bullet not found")
)
