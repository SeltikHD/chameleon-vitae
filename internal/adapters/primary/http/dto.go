package http

import (
	"encoding/json"
	"net/http"
	"time"
)

// ===============================
// Common Response Types
// ===============================

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status  string `json:"status" example:"healthy"`
	Service string `json:"service" example:"chameleon-vitae"`
}

// ErrorDetail represents a single field error.
type ErrorDetail struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Invalid email format"`
}

// ErrorResponse represents the standard error response format.
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

// ErrorBody contains error details.
type ErrorBody struct {
	Code    string        `json:"code" example:"VALIDATION_ERROR"`
	Message string        `json:"message" example:"Validation failed"`
	Details []ErrorDetail `json:"details,omitempty"`
}

// SuccessResponse represents a generic success response with data.
type SuccessResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty" example:"Operation successful"`
}

// PaginatedResponse represents a paginated list response.
type PaginatedResponse struct {
	Data   any `json:"data"`
	Total  int `json:"total" example:"100"`
	Limit  int `json:"limit" example:"50"`
	Offset int `json:"offset" example:"0"`
}

// ===============================
// Auth DTOs
// ===============================

// SyncUserRequest represents the request body for user synchronization.
type SyncUserRequest struct {
	FirebaseUID string  `json:"firebase_uid" example:"abc123xyz"`
	Email       *string `json:"email,omitempty" example:"user@example.com"`
	Name        *string `json:"name,omitempty" example:"John Doe"`
	Picture     *string `json:"picture,omitempty" example:"https://example.com/photo.jpg"`
}

// SyncUserResponse represents the response after user synchronization.
type SyncUserResponse struct {
	ID          string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FirebaseUID string    `json:"firebase_uid" example:"abc123xyz"`
	Email       *string   `json:"email,omitempty" example:"user@example.com"`
	Name        *string   `json:"name,omitempty" example:"John Doe"`
	CreatedAt   time.Time `json:"created_at" example:"2026-01-09T10:00:00Z"`
}

// ===============================
// User DTOs
// ===============================

// UserResponse represents the user profile response.
type UserResponse struct {
	ID                string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FirebaseUID       string    `json:"firebase_uid" example:"abc123xyz"`
	PictureURL        *string   `json:"picture_url,omitempty" example:"https://example.com/photo.jpg"`
	Email             *string   `json:"email,omitempty" example:"user@example.com"`
	Name              *string   `json:"name,omitempty" example:"John Doe"`
	Headline          *string   `json:"headline,omitempty" example:"Senior Software Engineer"`
	Summary           *string   `json:"summary,omitempty" example:"Experienced developer..."`
	Location          *string   `json:"location,omitempty" example:"San Francisco, CA"`
	Phone             *string   `json:"phone,omitempty" example:"+1-555-123-4567"`
	Website           *string   `json:"website,omitempty" example:"https://johndoe.dev"`
	LinkedInURL       *string   `json:"linkedin_url,omitempty" example:"https://linkedin.com/in/johndoe"`
	GitHubURL         *string   `json:"github_url,omitempty" example:"https://github.com/johndoe"`
	PortfolioURL      *string   `json:"portfolio_url,omitempty" example:"https://portfolio.johndoe.dev"`
	PreferredLanguage string    `json:"preferred_language" example:"en"`
	CreatedAt         time.Time `json:"created_at" example:"2026-01-09T10:00:00Z"`
	UpdatedAt         time.Time `json:"updated_at" example:"2026-01-09T10:00:00Z"`
}

// UpdateUserRequest represents the request body for updating user profile.
type UpdateUserRequest struct {
	Name              *string `json:"name,omitempty" example:"John Doe"`
	Headline          *string `json:"headline,omitempty" example:"Senior Software Engineer"`
	Summary           *string `json:"summary,omitempty" example:"Experienced developer..."`
	Location          *string `json:"location,omitempty" example:"San Francisco, CA"`
	Phone             *string `json:"phone,omitempty" example:"+1-555-123-4567"`
	Website           *string `json:"website,omitempty" example:"https://johndoe.dev"`
	LinkedInURL       *string `json:"linkedin_url,omitempty" example:"https://linkedin.com/in/johndoe"`
	GitHubURL         *string `json:"github_url,omitempty" example:"https://github.com/johndoe"`
	PortfolioURL      *string `json:"portfolio_url,omitempty" example:"https://portfolio.johndoe.dev"`
	PreferredLanguage *string `json:"preferred_language,omitempty" example:"en"`
}

// ===============================
// Experience DTOs
// ===============================

// ExperienceResponse represents an experience in API responses.
type ExperienceResponse struct {
	ID           string           `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Type         string           `json:"type" example:"work"`
	Title        string           `json:"title" example:"Senior Software Engineer"`
	Organization string           `json:"organization" example:"Tech Company Inc."`
	Location     *string          `json:"location,omitempty" example:"Remote"`
	StartDate    string           `json:"start_date" example:"2022-01-15"`
	EndDate      *string          `json:"end_date,omitempty" example:"2024-06-30"`
	IsCurrent    bool             `json:"is_current" example:"false"`
	Description  *string          `json:"description,omitempty" example:"Led backend development..."`
	URL          *string          `json:"url,omitempty" example:"https://techcompany.com"`
	Metadata     map[string]any   `json:"metadata,omitempty"`
	DisplayOrder int              `json:"display_order" example:"0"`
	Bullets      []BulletResponse `json:"bullets,omitempty"`
	CreatedAt    time.Time        `json:"created_at" example:"2026-01-09T10:00:00Z"`
	UpdatedAt    time.Time        `json:"updated_at" example:"2026-01-09T10:00:00Z"`
}

// CreateExperienceRequest represents the request body for creating an experience.
type CreateExperienceRequest struct {
	Type         string         `json:"type" example:"work"`
	Title        string         `json:"title" example:"Senior Software Engineer"`
	Organization string         `json:"organization" example:"Tech Company Inc."`
	Location     *string        `json:"location,omitempty" example:"Remote"`
	StartDate    string         `json:"start_date" example:"2022-01-15"`
	EndDate      *string        `json:"end_date,omitempty" example:"2024-06-30"`
	IsCurrent    bool           `json:"is_current" example:"false"`
	Description  *string        `json:"description,omitempty" example:"Led backend development..."`
	URL          *string        `json:"url,omitempty" example:"https://techcompany.com"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

// UpdateExperienceRequest represents the request body for updating an experience.
type UpdateExperienceRequest struct {
	Type         *string        `json:"type,omitempty" example:"work"`
	Title        *string        `json:"title,omitempty" example:"Senior Software Engineer"`
	Organization *string        `json:"organization,omitempty" example:"Tech Company Inc."`
	Location     *string        `json:"location,omitempty" example:"Remote"`
	StartDate    *string        `json:"start_date,omitempty" example:"2022-01-15"`
	EndDate      *string        `json:"end_date,omitempty" example:"2024-06-30"`
	IsCurrent    *bool          `json:"is_current,omitempty" example:"false"`
	Description  *string        `json:"description,omitempty" example:"Led backend development..."`
	URL          *string        `json:"url,omitempty" example:"https://techcompany.com"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	DisplayOrder *int           `json:"display_order,omitempty" example:"1"`
}

// ListExperiencesResponse represents the paginated list of experiences.
type ListExperiencesResponse struct {
	Data   []ExperienceResponse `json:"data"`
	Total  int                  `json:"total" example:"15"`
	Limit  int                  `json:"limit" example:"50"`
	Offset int                  `json:"offset" example:"0"`
}

// ===============================
// Bullet DTOs
// ===============================

// BulletResponse represents a bullet in API responses.
type BulletResponse struct {
	ID           string         `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ExperienceID string         `json:"experience_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Content      string         `json:"content" example:"Reduced API latency by 40%"`
	ImpactScore  int            `json:"impact_score" example:"75"`
	Keywords     []string       `json:"keywords" example:"performance,optimization"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	DisplayOrder int            `json:"display_order" example:"0"`
	CreatedAt    time.Time      `json:"created_at" example:"2026-01-09T10:00:00Z"`
	UpdatedAt    time.Time      `json:"updated_at" example:"2026-01-09T10:00:00Z"`
}

// CreateBulletRequest represents the request body for creating a bullet.
type CreateBulletRequest struct {
	Content      string   `json:"content" example:"Reduced API latency by 40%"`
	Keywords     []string `json:"keywords,omitempty" example:"performance,optimization"`
	DisplayOrder *int     `json:"display_order,omitempty" example:"0"`
}

// UpdateBulletRequest represents the request body for updating a bullet.
type UpdateBulletRequest struct {
	Content      *string  `json:"content,omitempty" example:"Updated bullet content"`
	Keywords     []string `json:"keywords,omitempty" example:"new,keywords"`
	DisplayOrder *int     `json:"display_order,omitempty" example:"1"`
}

// ScoreBulletResponse represents the response after recalculating bullet score.
type ScoreBulletResponse struct {
	ID             string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Content        string `json:"content" example:"Reduced API latency by 40%"`
	ImpactScore    int    `json:"impact_score" example:"85"`
	ScoreReasoning string `json:"score_reasoning,omitempty" example:"Contains quantifiable metrics..."`
}

// AnalyzeBulletRequest represents the request for analyzing a bullet's impact.
type AnalyzeBulletRequest struct {
	JobDescription string `json:"job_description" example:"We are looking for a Senior Backend Engineer..."`
}

// ===============================
// Education DTOs
// ===============================

// EducationResponse represents an education entry in API responses.
type EducationResponse struct {
	ID           string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Institution  string    `json:"institution" example:"Massachusetts Institute of Technology"`
	Degree       string    `json:"degree" example:"Bachelor of Science"`
	FieldOfStudy string    `json:"field_of_study" example:"Computer Science"`
	Location     *string   `json:"location,omitempty" example:"Cambridge, MA"`
	StartDate    *string   `json:"start_date,omitempty" example:"2018-09-01"`
	EndDate      *string   `json:"end_date,omitempty" example:"2022-05-15"`
	GPA          *string   `json:"gpa,omitempty" example:"3.85/4.0"`
	Honors       []string  `json:"honors,omitempty" example:"Magna Cum Laude,Dean's List"`
	DisplayOrder int       `json:"display_order" example:"0"`
	CreatedAt    time.Time `json:"created_at" example:"2026-01-09T10:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2026-01-09T10:00:00Z"`
}

// CreateEducationRequest represents the request body for creating an education entry.
type CreateEducationRequest struct {
	Institution  string   `json:"institution" example:"Massachusetts Institute of Technology"`
	Degree       string   `json:"degree" example:"Bachelor of Science"`
	FieldOfStudy string   `json:"field_of_study" example:"Computer Science"`
	Location     *string  `json:"location,omitempty" example:"Cambridge, MA"`
	StartDate    *string  `json:"start_date,omitempty" example:"2018-09-01"`
	EndDate      *string  `json:"end_date,omitempty" example:"2022-05-15"`
	GPA          *string  `json:"gpa,omitempty" example:"3.85/4.0"`
	Honors       []string `json:"honors,omitempty" example:"Magna Cum Laude,Dean's List"`
	DisplayOrder int      `json:"display_order,omitempty" example:"0"`
}

// UpdateEducationRequest represents the request body for updating an education entry.
type UpdateEducationRequest struct {
	Institution  *string  `json:"institution,omitempty" example:"MIT"`
	Degree       *string  `json:"degree,omitempty" example:"Master of Science"`
	FieldOfStudy *string  `json:"field_of_study,omitempty" example:"Artificial Intelligence"`
	Location     *string  `json:"location,omitempty" example:"Cambridge, MA"`
	StartDate    *string  `json:"start_date,omitempty" example:"2022-09-01"`
	EndDate      *string  `json:"end_date,omitempty" example:"2024-05-15"`
	GPA          *string  `json:"gpa,omitempty" example:"3.95/4.0"`
	Honors       []string `json:"honors,omitempty" example:"Summa Cum Laude"`
	DisplayOrder *int     `json:"display_order,omitempty" example:"1"`
}

// ListEducationResponse represents the list of education entries.
type ListEducationResponse struct {
	Data  []EducationResponse `json:"data"`
	Total int                 `json:"total" example:"2"`
}

// ===============================
// Project DTOs
// ===============================

// ProjectResponse represents a project in API responses.
type ProjectResponse struct {
	ID            string                  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name          string                  `json:"name" example:"Chameleon Vitae"`
	Description   string                  `json:"description,omitempty" example:"AI-powered resume engineering tool"`
	TechStack     []string                `json:"tech_stack" example:"Go,Vue.js,PostgreSQL"`
	URL           string                  `json:"url,omitempty" example:"https://chameleon-vitae.dev"`
	RepositoryURL string                  `json:"repository_url,omitempty" example:"https://github.com/SeltikHD/chameleon-vitae"`
	StartDate     *string                 `json:"start_date,omitempty" example:"2024-01-01"`
	EndDate       *string                 `json:"end_date,omitempty" example:"2024-06-30"`
	DisplayOrder  int                     `json:"display_order" example:"0"`
	Bullets       []ProjectBulletResponse `json:"bullets,omitempty"`
	CreatedAt     time.Time               `json:"created_at" example:"2026-01-09T10:00:00Z"`
	UpdatedAt     time.Time               `json:"updated_at" example:"2026-01-09T10:00:00Z"`
}

// ProjectBulletResponse represents a project bullet in API responses.
type ProjectBulletResponse struct {
	ID           string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProjectID    string    `json:"project_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Content      string    `json:"content" example:"Implemented AI-driven bullet selection algorithm"`
	DisplayOrder int       `json:"display_order" example:"0"`
	CreatedAt    time.Time `json:"created_at" example:"2026-01-09T10:00:00Z"`
}

// CreateProjectRequest represents the request body for creating a project.
type CreateProjectRequest struct {
	Name          string   `json:"name" example:"Chameleon Vitae"`
	Description   string   `json:"description,omitempty" example:"AI-powered resume engineering tool"`
	TechStack     []string `json:"tech_stack,omitempty" example:"Go,Vue.js,PostgreSQL"`
	URL           string   `json:"url,omitempty" example:"https://chameleon-vitae.dev"`
	RepositoryURL string   `json:"repository_url,omitempty" example:"https://github.com/SeltikHD/chameleon-vitae"`
	StartDate     *string  `json:"start_date,omitempty" example:"2024-01-01"`
	EndDate       *string  `json:"end_date,omitempty" example:"2024-06-30"`
	DisplayOrder  int      `json:"display_order,omitempty" example:"0"`
	Bullets       []string `json:"bullets,omitempty" example:"Implemented feature X,Optimized performance by 50%"`
}

// UpdateProjectRequest represents the request body for updating a project.
type UpdateProjectRequest struct {
	Name          *string  `json:"name,omitempty" example:"Chameleon Vitae v2"`
	Description   *string  `json:"description,omitempty" example:"Updated description"`
	TechStack     []string `json:"tech_stack,omitempty" example:"Go,Vue.js,PostgreSQL,Redis"`
	URL           *string  `json:"url,omitempty" example:"https://new-url.dev"`
	RepositoryURL *string  `json:"repository_url,omitempty" example:"https://github.com/new/repo"`
	StartDate     *string  `json:"start_date,omitempty" example:"2024-01-01"`
	EndDate       *string  `json:"end_date,omitempty" example:"2025-01-01"`
	DisplayOrder  *int     `json:"display_order,omitempty" example:"1"`
}

// CreateProjectBulletRequest represents the request body for creating a project bullet.
type CreateProjectBulletRequest struct {
	Content      string `json:"content" example:"Implemented new feature"`
	DisplayOrder int    `json:"display_order,omitempty" example:"0"`
}

// ListProjectsResponse represents the list of projects.
type ListProjectsResponse struct {
	Data  []ProjectResponse `json:"data"`
	Total int               `json:"total" example:"5"`
}

// ===============================
// Skill DTOs
// ===============================

// SkillResponse represents a skill in API responses.
type SkillResponse struct {
	ID                string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name              string    `json:"name" example:"Go"`
	Category          *string   `json:"category,omitempty" example:"Programming Languages"`
	ProficiencyLevel  int       `json:"proficiency_level" example:"85"`
	YearsOfExperience *float64  `json:"years_of_experience,omitempty" example:"3.5"`
	IsHighlighted     bool      `json:"is_highlighted" example:"true"`
	DisplayOrder      int       `json:"display_order" example:"0"`
	CreatedAt         time.Time `json:"created_at" example:"2026-01-09T10:00:00Z"`
}

// SkillInput represents a single skill in batch operations.
type SkillInput struct {
	Name              string   `json:"name" example:"Go"`
	Category          *string  `json:"category,omitempty" example:"Programming Languages"`
	ProficiencyLevel  *int     `json:"proficiency_level,omitempty" example:"85"`
	YearsOfExperience *float64 `json:"years_of_experience,omitempty" example:"3.5"`
	IsHighlighted     *bool    `json:"is_highlighted,omitempty" example:"true"`
}

// BatchUpsertSkillsRequest represents the request for batch skill upsert.
type BatchUpsertSkillsRequest struct {
	Skills []SkillInput `json:"skills"`
}

// BatchUpsertSkillsResponse represents the response for batch skill upsert.
type BatchUpsertSkillsResponse struct {
	Created int             `json:"created" example:"1"`
	Updated int             `json:"updated" example:"2"`
	Data    []SkillResponse `json:"data"`
}

// ListSkillsResponse represents the list of skills.
type ListSkillsResponse struct {
	Data  []SkillResponse `json:"data"`
	Total int             `json:"total" example:"25"`
}

// ===============================
// Spoken Language DTOs
// ===============================

// SpokenLanguageResponse represents a spoken language in API responses.
type SpokenLanguageResponse struct {
	ID           string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Language     string    `json:"language" example:"English"`
	Proficiency  string    `json:"proficiency" example:"native"`
	DisplayOrder int       `json:"display_order" example:"0"`
	CreatedAt    time.Time `json:"created_at" example:"2026-01-09T10:00:00Z"`
}

// CreateSpokenLanguageRequest represents the request for creating a spoken language.
type CreateSpokenLanguageRequest struct {
	Language     string `json:"language" example:"Portuguese"`
	Proficiency  string `json:"proficiency" example:"native"`
	DisplayOrder *int   `json:"display_order,omitempty" example:"0"`
}

// ListSpokenLanguagesResponse represents the list of spoken languages.
type ListSpokenLanguagesResponse struct {
	Data []SpokenLanguageResponse `json:"data"`
}

// ===============================
// Resume DTOs
// ===============================

// ResumeResponse represents a resume in API responses.
type ResumeResponse struct {
	ID               string            `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	JobTitle         string            `json:"job_title,omitempty" example:"Senior Backend Engineer"`
	CompanyName      string            `json:"company_name,omitempty" example:"Awesome Corp"`
	JobURL           string            `json:"job_url,omitempty" example:"https://linkedin.com/jobs/..."`
	JobDescription   string            `json:"job_description,omitempty"`
	TargetLanguage   string            `json:"target_language" example:"en"`
	SelectedBullets  []string          `json:"selected_bullets,omitempty"`
	GeneratedContent *ResumeContentDTO `json:"generated_content,omitempty"`
	PDFURL           string            `json:"pdf_url,omitempty" example:"https://storage.../resume.pdf"`
	Score            int               `json:"score" example:"85"`
	Notes            string            `json:"notes,omitempty"`
	Status           string            `json:"status" example:"draft"`
	CreatedAt        time.Time         `json:"created_at" example:"2026-01-09T10:00:00Z"`
	UpdatedAt        time.Time         `json:"updated_at" example:"2026-01-09T10:00:00Z"`
}

// ResumeContentDTO represents the AI-generated resume content.
type ResumeContentDTO struct {
	Summary     string                  `json:"summary"`
	Experiences []TailoredExperienceDTO `json:"experiences"`
	Skills      []string                `json:"skills"`
	Analysis    *ResumeAnalysisDTO      `json:"analysis,omitempty"`
}

// TailoredExperienceDTO represents a tailored experience entry.
type TailoredExperienceDTO struct {
	ExperienceID string              `json:"experience_id"`
	Title        string              `json:"title"`
	Organization string              `json:"organization"`
	StartDate    string              `json:"start_date"`
	EndDate      *string             `json:"end_date,omitempty"`
	IsCurrent    bool                `json:"is_current"`
	Bullets      []TailoredBulletDTO `json:"bullets"`
}

// TailoredBulletDTO represents a tailored bullet point.
type TailoredBulletDTO struct {
	BulletID        string `json:"bullet_id"`
	OriginalContent string `json:"original_content"`
	TailoredContent string `json:"tailored_content"`
}

// ResumeAnalysisDTO contains the AI analysis of how well the resume matches.
type ResumeAnalysisDTO struct {
	MatchedKeywords []string `json:"matched_keywords"`
	MissingKeywords []string `json:"missing_keywords"`
	Recommendations []string `json:"recommendations"`
}

// ResumeListItem represents a resume in list responses (without full content).
type ResumeListItem struct {
	ID             string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	JobTitle       string    `json:"job_title" example:"Senior Backend Engineer"`
	CompanyName    string    `json:"company_name" example:"Awesome Corp"`
	JobURL         *string   `json:"job_url,omitempty" example:"https://linkedin.com/jobs/..."`
	TargetLanguage string    `json:"target_language" example:"en"`
	Score          *int      `json:"score,omitempty" example:"85"`
	Status         string    `json:"status" example:"draft"`
	CreatedAt      time.Time `json:"created_at" example:"2026-01-09T10:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2026-01-09T10:00:00Z"`
}

// CreateResumeRequest represents the request for creating a resume.
type CreateResumeRequest struct {
	JobDescription string `json:"job_description" example:"## Senior Backend Engineer..."`
	JobTitle       string `json:"job_title,omitempty" example:"Senior Backend Engineer"`
	CompanyName    string `json:"company_name,omitempty" example:"Awesome Corp"`
	JobURL         string `json:"job_url,omitempty" example:"https://linkedin.com/jobs/12345"`
	TargetLanguage string `json:"target_language,omitempty" example:"en"`
}

// TailorResumeRequest represents the request for tailoring a resume.
type TailorResumeRequest struct {
	MaxBulletsPerJob int `json:"max_bullets_per_job,omitempty" example:"15"`
}

// TailorResumeResponse represents the response after tailoring a resume.
type TailorResumeResponse struct {
	ID               string          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status           string          `json:"status" example:"generated"`
	Score            int             `json:"score" example:"85"`
	SelectedBullets  []string        `json:"selected_bullets"`
	GeneratedContent map[string]any  `json:"generated_content"`
	Analysis         *TailorAnalysis `json:"analysis,omitempty"`
}

// TailorAnalysis contains the analysis result from tailoring.
type TailorAnalysis struct {
	MatchedKeywords []string `json:"matched_keywords" example:"golang,microservices"`
	MissingKeywords []string `json:"missing_keywords" example:"terraform"`
	Recommendations []string `json:"recommendations" example:"Consider adding cloud experience"`
}

// UpdateResumeContentRequest represents the request for updating resume content.
type UpdateResumeContentRequest struct {
	Status string  `json:"status" example:"reviewed"`
	Notes  *string `json:"notes,omitempty" example:"Made adjustments to summary"`
}

// ListResumesResponse represents the paginated list of resumes.
type ListResumesResponse struct {
	Data   []ResumeResponse `json:"data"`
	Total  int              `json:"total" example:"10"`
	Limit  int              `json:"limit" example:"20"`
	Offset int              `json:"offset" example:"0"`
}

// ===============================
// Tools DTOs
// ===============================

// ParseJobURLRequest represents the request for parsing a job URL.
type ParseJobURLRequest struct {
	URL string `json:"url" example:"https://linkedin.com/jobs/view/12345"`
}

// ParseJobURLResponse represents the parsed job posting.
type ParseJobURLResponse struct {
	URL      string            `json:"url" example:"https://linkedin.com/jobs/view/12345"`
	Title    string            `json:"title" example:"Senior Backend Engineer at Awesome Corp"`
	Markdown string            `json:"markdown" example:"## Senior Backend Engineer..."`
	Metadata *ParseJobMetadata `json:"metadata,omitempty"`
}

// ParseJobMetadata contains metadata about the parsed job.
type ParseJobMetadata struct {
	Source    string    `json:"source" example:"linkedin.com"`
	FetchedAt time.Time `json:"fetched_at" example:"2026-01-09T10:00:00Z"`
}

// ===============================
// Helper Functions
// ===============================

// respondJSON writes a JSON response with the given status code.
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// respondError writes an error response.
func respondError(w http.ResponseWriter, status int, code, message string) {
	respondJSON(w, status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

// respondErrorWithDetails writes an error response with field details.
func respondErrorWithDetails(w http.ResponseWriter, status int, code, message string, details []ErrorDetail) {
	respondJSON(w, status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// decodeJSON decodes a JSON request body into the target struct.
func decodeJSON(r *http.Request, target any) error {
	if r.Body == nil {
		return ErrEmptyRequestBody
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}
