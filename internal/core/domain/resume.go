// Package domain contains the core business entities and value objects.
package domain

import "time"

// Resume represents a generated resume tailored to a specific job application.
type Resume struct {
	ID               string         `json:"id"`
	UserID           string         `json:"user_id"`
	JobDescription   string         `json:"job_description"`
	JobTitle         *string        `json:"job_title,omitempty"`
	CompanyName      *string        `json:"company_name,omitempty"`
	JobURL           *string        `json:"job_url,omitempty"`
	TargetLanguage   string         `json:"target_language"`
	SelectedBullets  []string       `json:"selected_bullets"`
	GeneratedContent *ResumeContent `json:"generated_content,omitempty"`
	PDFURL           *string        `json:"pdf_url,omitempty"`
	Score            MatchScore     `json:"score"`
	Notes            *string        `json:"notes,omitempty"`
	Status           ResumeStatus   `json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// ResumeContent represents the AI-generated content for a resume.
type ResumeContent struct {
	Summary     string               `json:"summary"`
	Experiences []TailoredExperience `json:"experiences"`
	Skills      []string             `json:"skills"`
	Analysis    *ResumeAnalysis      `json:"analysis,omitempty"`
}

// TailoredExperience represents an experience entry tailored for a specific job.
type TailoredExperience struct {
	ExperienceID string           `json:"experience_id"`
	Title        string           `json:"title"`
	Organization string           `json:"organization"`
	StartDate    string           `json:"start_date"`
	EndDate      *string          `json:"end_date,omitempty"`
	IsCurrent    bool             `json:"is_current"`
	Bullets      []TailoredBullet `json:"bullets"`
}

// TailoredBullet represents a bullet point that has been tailored for a specific job.
type TailoredBullet struct {
	BulletID        string `json:"bullet_id"`
	OriginalContent string `json:"original_content"`
	TailoredContent string `json:"tailored_content"`
}

// ResumeAnalysis contains the AI analysis of how well the resume matches the job.
type ResumeAnalysis struct {
	MatchedKeywords  []string `json:"matched_keywords"`
	MissingKeywords  []string `json:"missing_keywords"`
	Recommendations  []string `json:"recommendations"`
	StrengthAreas    []string `json:"strength_areas"`
	ImprovementAreas []string `json:"improvement_areas"`
}

// NewResume creates a new resume draft with required fields.
func NewResume(userID, jobDescription string) (*Resume, error) {
	if jobDescription == "" {
		return nil, ErrEmptyJobDescription
	}

	now := time.Now().UTC()
	score, _ := NewMatchScore(0)

	return &Resume{
		UserID:          userID,
		JobDescription:  jobDescription,
		TargetLanguage:  "en",
		SelectedBullets: make([]string, 0),
		Score:           score,
		Status:          ResumeStatusDraft,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Validate validates the resume entity.
func (r *Resume) Validate() error {
	v := &ValidationErrors{}

	if r.UserID == "" {
		v.AddFieldError("user_id", "user ID is required")
	}

	if r.JobDescription == "" {
		v.AddFieldError("job_description", "job description is required")
	}

	if !r.Status.IsValid() {
		v.AddFieldError("status", "invalid resume status")
	}

	if r.TargetLanguage != "en" && r.TargetLanguage != "pt-br" {
		v.AddFieldError("target_language", "must be 'en' or 'pt-br'")
	}

	return v.ToError()
}

// SetJobDetails sets the job-related metadata.
func (r *Resume) SetJobDetails(title, company, url string) {
	if title != "" {
		r.JobTitle = &title
	}
	if company != "" {
		r.CompanyName = &company
	}
	if url != "" {
		r.JobURL = &url
	}
	r.UpdatedAt = time.Now().UTC()
}

// SetGeneratedContent sets the AI-generated content.
func (r *Resume) SetGeneratedContent(content *ResumeContent) {
	r.GeneratedContent = content
	r.Status = ResumeStatusGenerated
	r.UpdatedAt = time.Now().UTC()
}

// SetScore sets the match score.
func (r *Resume) SetScore(score int) error {
	matchScore, err := NewMatchScore(score)
	if err != nil {
		return err
	}
	r.Score = matchScore
	r.UpdatedAt = time.Now().UTC()
	return nil
}

// SelectBullets sets the selected bullet IDs.
func (r *Resume) SelectBullets(bulletIDs []string) {
	r.SelectedBullets = bulletIDs
	r.UpdatedAt = time.Now().UTC()
}

// AddSelectedBullet adds a bullet to the selection.
func (r *Resume) AddSelectedBullet(bulletID string) {
	for _, id := range r.SelectedBullets {
		if id == bulletID {
			return // Already selected
		}
	}
	r.SelectedBullets = append(r.SelectedBullets, bulletID)
	r.UpdatedAt = time.Now().UTC()
}

// RemoveSelectedBullet removes a bullet from the selection.
func (r *Resume) RemoveSelectedBullet(bulletID string) {
	for i, id := range r.SelectedBullets {
		if id == bulletID {
			r.SelectedBullets = append(r.SelectedBullets[:i], r.SelectedBullets[i+1:]...)
			r.UpdatedAt = time.Now().UTC()
			return
		}
	}
}

// SetPDFURL sets the URL of the generated PDF.
func (r *Resume) SetPDFURL(url string) {
	r.PDFURL = &url
	r.UpdatedAt = time.Now().UTC()
}

// TransitionStatus transitions the resume to a new status.
func (r *Resume) TransitionStatus(newStatus ResumeStatus) error {
	if !newStatus.IsValid() {
		return ErrInvalidResumeStatus
	}

	// Define valid transitions.
	validTransitions := map[ResumeStatus][]ResumeStatus{
		ResumeStatusDraft:     {ResumeStatusGenerated},
		ResumeStatusGenerated: {ResumeStatusReviewed, ResumeStatusDraft},
		ResumeStatusReviewed:  {ResumeStatusSubmitted, ResumeStatusGenerated},
		ResumeStatusSubmitted: {ResumeStatusInterview, ResumeStatusRejected},
		ResumeStatusInterview: {ResumeStatusAccepted, ResumeStatusRejected},
		ResumeStatusRejected:  {}, // Terminal state
		ResumeStatusAccepted:  {}, // Terminal state
	}

	allowed, exists := validTransitions[r.Status]
	if !exists {
		return ErrInvalidStatusTransition
	}

	for _, status := range allowed {
		if status == newStatus {
			r.Status = newStatus
			r.UpdatedAt = time.Now().UTC()
			return nil
		}
	}

	return ErrInvalidStatusTransition
}

// IsDraft returns true if the resume is a draft.
func (r *Resume) IsDraft() bool {
	return r.Status == ResumeStatusDraft
}

// IsGenerated returns true if the resume has been generated.
func (r *Resume) IsGenerated() bool {
	return r.Status == ResumeStatusGenerated || r.Status == ResumeStatusReviewed
}

// IsSubmitted returns true if the resume has been submitted.
func (r *Resume) IsSubmitted() bool {
	return r.Status == ResumeStatusSubmitted ||
		r.Status == ResumeStatusInterview ||
		r.Status == ResumeStatusRejected ||
		r.Status == ResumeStatusAccepted
}

// CanGeneratePDF returns true if the resume can be exported to PDF.
func (r *Resume) CanGeneratePDF() bool {
	return r.GeneratedContent != nil &&
		(r.Status == ResumeStatusGenerated || r.Status == ResumeStatusReviewed)
}

// GetJobDisplayName returns a display name for the job.
func (r *Resume) GetJobDisplayName() string {
	if r.JobTitle != nil && r.CompanyName != nil {
		return *r.JobTitle + " at " + *r.CompanyName
	}
	if r.JobTitle != nil {
		return *r.JobTitle
	}
	if r.CompanyName != nil {
		return "Position at " + *r.CompanyName
	}
	return "Untitled Resume"
}
