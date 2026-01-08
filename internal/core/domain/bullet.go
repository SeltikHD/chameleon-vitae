// Package domain contains the core business entities and value objects.
package domain

import "time"

// Bullet represents an atomic unit of experience that can be selected for resume tailoring.
// Each bullet is a single achievement or responsibility that can be independently
// selected and potentially rewritten by AI for a specific job application.
type Bullet struct {
	ID           string         `json:"id"`
	ExperienceID string         `json:"experience_id"`
	Content      string         `json:"content"`
	ImpactScore  ImpactScore    `json:"impact_score"`
	Keywords     []string       `json:"keywords"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	DisplayOrder int            `json:"display_order"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// NewBullet creates a new bullet with required fields.
func NewBullet(experienceID, content string) (*Bullet, error) {
	if content == "" {
		return nil, ErrEmptyBulletContent
	}

	now := time.Now().UTC()
	return &Bullet{
		ExperienceID: experienceID,
		Content:      content,
		ImpactScore:  DefaultImpactScore(),
		Keywords:     make([]string, 0),
		Metadata:     make(map[string]any),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Validate validates the bullet entity.
func (b *Bullet) Validate() error {
	v := &ValidationErrors{}

	if b.ExperienceID == "" {
		v.AddFieldError("experience_id", "experience ID is required")
	}

	if b.Content == "" {
		v.AddFieldError("content", "content is required")
	}

	if b.ImpactScore.Int() < 0 || b.ImpactScore.Int() > 100 {
		v.AddFieldError("impact_score", "must be between 0 and 100")
	}

	return v.ToError()
}

// SetImpactScore sets the impact score after validation.
func (b *Bullet) SetImpactScore(score int) error {
	impactScore, err := NewImpactScore(score)
	if err != nil {
		return err
	}
	b.ImpactScore = impactScore
	b.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateContent updates the bullet content.
func (b *Bullet) UpdateContent(content string) error {
	if content == "" {
		return ErrEmptyBulletContent
	}
	b.Content = content
	b.UpdatedAt = time.Now().UTC()
	return nil
}

// SetKeywords sets the keywords for the bullet.
func (b *Bullet) SetKeywords(keywords []string) {
	b.Keywords = keywords
	b.UpdatedAt = time.Now().UTC()
}

// AddKeyword adds a keyword to the bullet.
func (b *Bullet) AddKeyword(keyword string) {
	for _, k := range b.Keywords {
		if k == keyword {
			return // Already exists
		}
	}
	b.Keywords = append(b.Keywords, keyword)
	b.UpdatedAt = time.Now().UTC()
}

// HasKeyword checks if the bullet has a specific keyword.
func (b *Bullet) HasKeyword(keyword string) bool {
	for _, k := range b.Keywords {
		if k == keyword {
			return true
		}
	}
	return false
}

// IsHighImpact returns true if the bullet has a high impact score (>= 70).
func (b *Bullet) IsHighImpact() bool {
	return b.ImpactScore.Int() >= 70
}

// IsLowImpact returns true if the bullet has a low impact score (< 40).
func (b *Bullet) IsLowImpact() bool {
	return b.ImpactScore.Int() < 40
}
