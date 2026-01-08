// Package ports defines the interfaces (ports) that adapters must implement.
package ports

import (
	"context"
	"io"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
)

// Note: AuthProvider and AuthClaims are defined in auth.go

// AIProvider defines the interface for AI-powered resume operations.
// Implementations should handle communication with LLM APIs (e.g., Groq).
type AIProvider interface {
	// AnalyzeJob analyzes a job description and extracts key requirements.
	AnalyzeJob(ctx context.Context, req AnalyzeJobRequest) (*JobAnalysis, error)

	// SelectBullets selects the most relevant bullets for a job description.
	SelectBullets(ctx context.Context, req SelectBulletsRequest) (*BulletSelection, error)

	// TailorBullet rewrites a bullet to better match job requirements.
	TailorBullet(ctx context.Context, req TailorBulletRequest) (*TailoredBulletResult, error)

	// GenerateSummary generates a professional summary tailored to the job.
	GenerateSummary(ctx context.Context, req GenerateSummaryRequest) (*SummaryResult, error)

	// ScoreMatch calculates a match score between resume and job.
	ScoreMatch(ctx context.Context, req ScoreMatchRequest) (*domain.MatchScore, error)

	// Close releases any resources held by the AI provider.
	Close() error
}

// AnalyzeJobRequest contains parameters for job analysis.
type AnalyzeJobRequest struct {
	// JobDescription is the parsed job description text.
	JobDescription string

	// TargetLanguage is the language for the analysis output.
	TargetLanguage string
}

// JobAnalysis contains the result of job description analysis.
type JobAnalysis struct {
	// Title is the extracted job title.
	Title string

	// Company is the extracted company name.
	Company string

	// RequiredSkills are skills explicitly required.
	RequiredSkills []string

	// PreferredSkills are nice-to-have skills.
	PreferredSkills []string

	// Keywords are important keywords from the description.
	Keywords []string

	// SeniorityLevel is the detected seniority level.
	SeniorityLevel string

	// YearsExperience is the required years of experience.
	YearsExperience *int

	// Summary is a brief summary of the job requirements.
	Summary string
}

// SelectBulletsRequest contains parameters for bullet selection.
type SelectBulletsRequest struct {
	// JobAnalysis is the analyzed job description.
	JobAnalysis *JobAnalysis

	// AvailableBullets are all bullets to select from.
	AvailableBullets []domain.Bullet

	// MaxBullets is the maximum number of bullets to select.
	MaxBullets int

	// TargetLanguage is the language for selection.
	TargetLanguage string
}

// BulletSelection contains the result of bullet selection.
type BulletSelection struct {
	// SelectedBulletIDs are the IDs of selected bullets in order of relevance.
	SelectedBulletIDs []string

	// Reasoning explains why these bullets were selected.
	Reasoning string
}

// TailorBulletRequest contains parameters for bullet tailoring.
type TailorBulletRequest struct {
	// Bullet is the original bullet to tailor.
	Bullet domain.Bullet

	// JobAnalysis is the analyzed job description.
	JobAnalysis *JobAnalysis

	// TargetLanguage is the output language.
	TargetLanguage string

	// Style is the writing style (e.g., "professional", "technical").
	Style string
}

// TailoredBulletResult contains the result of bullet tailoring.
type TailoredBulletResult struct {
	// OriginalID is the ID of the original bullet.
	OriginalID string

	// TailoredContent is the rewritten bullet content.
	TailoredContent string

	// Keywords are the matched keywords from the job.
	Keywords []string
}

// GenerateSummaryRequest contains parameters for summary generation.
type GenerateSummaryRequest struct {
	// User is the user's profile information.
	User *domain.User

	// JobAnalysis is the analyzed job description.
	JobAnalysis *JobAnalysis

	// SelectedBullets are the bullets selected for this resume.
	SelectedBullets []domain.Bullet

	// TargetLanguage is the output language.
	TargetLanguage string
}

// SummaryResult contains the generated professional summary.
type SummaryResult struct {
	// Summary is the generated professional summary.
	Summary string
}

// ScoreMatchRequest contains parameters for match scoring.
type ScoreMatchRequest struct {
	// JobAnalysis is the analyzed job description.
	JobAnalysis *JobAnalysis

	// Resume is the generated resume content.
	Resume *domain.ResumeContent

	// UserSkills are the user's skills.
	UserSkills []domain.Skill
}

// PDFEngine defines the interface for PDF generation.
// Implementations should handle communication with Gotenberg.
type PDFEngine interface {
	// GeneratePDF generates a PDF from HTML content.
	GeneratePDF(ctx context.Context, req GeneratePDFRequest) (*PDFResult, error)

	// GetTemplates returns available resume templates.
	GetTemplates(ctx context.Context) ([]PDFTemplate, error)

	// HealthCheck checks if the PDF engine is available.
	HealthCheck(ctx context.Context) error

	// Close releases any resources held by the PDF engine.
	Close() error
}

// GeneratePDFRequest contains parameters for PDF generation.
type GeneratePDFRequest struct {
	// HTML is the HTML content to convert.
	HTML string

	// TemplateName is the name of the template to use.
	TemplateName string

	// CSS is optional custom CSS to apply.
	CSS string

	// Options are PDF generation options.
	Options PDFOptions
}

// PDFOptions contains options for PDF generation.
type PDFOptions struct {
	// PaperWidth is the paper width in inches.
	PaperWidth float64

	// PaperHeight is the paper height in inches.
	PaperHeight float64

	// MarginTop is the top margin in inches.
	MarginTop float64

	// MarginBottom is the bottom margin in inches.
	MarginBottom float64

	// MarginLeft is the left margin in inches.
	MarginLeft float64

	// MarginRight is the right margin in inches.
	MarginRight float64

	// Scale is the scale factor (0.1 to 2.0).
	Scale float64
}

// DefaultPDFOptions returns sensible defaults for PDF generation.
func DefaultPDFOptions() PDFOptions {
	return PDFOptions{
		PaperWidth:   8.5,
		PaperHeight:  11,
		MarginTop:    0.4,
		MarginBottom: 0.4,
		MarginLeft:   0.4,
		MarginRight:  0.4,
		Scale:        1.0,
	}
}

// PDFResult contains the generated PDF.
type PDFResult struct {
	// Content is the PDF file content.
	Content io.ReadCloser

	// Size is the size of the PDF in bytes.
	Size int64

	// Filename is the suggested filename.
	Filename string
}

// PDFTemplate represents an available PDF template.
type PDFTemplate struct {
	// Name is the unique template identifier.
	Name string

	// DisplayName is the human-readable name.
	DisplayName string

	// Description describes the template style.
	Description string

	// PreviewURL is a URL to a preview image.
	PreviewURL string
}

// JobParser defines the interface for parsing job descriptions from URLs.
// Implementations should handle communication with Jina Reader API.
type JobParser interface {
	// ParseJobURL fetches and parses a job description from a URL.
	ParseJobURL(ctx context.Context, url string) (*ParsedJob, error)

	// HealthCheck checks if the parser is available.
	HealthCheck(ctx context.Context) error

	// Close releases any resources held by the job parser.
	Close() error
}

// ParsedJob contains the parsed job description.
type ParsedJob struct {
	// URL is the original job URL.
	URL string

	// Title is the page title (often includes job title).
	Title string

	// Content is the main content in markdown format.
	Content string

	// Description is the extracted job description text.
	Description string

	// PublishedDate is when the job was published (if available).
	PublishedDate string

	// Metadata contains additional extracted metadata.
	Metadata map[string]string
}

// FileStorage defines the interface for file storage operations.
// Implementations could use local storage, S3, Azure Blob, etc.
type FileStorage interface {
	// Upload uploads a file and returns its URL.
	Upload(ctx context.Context, req UploadRequest) (*UploadResult, error)

	// Download downloads a file by its key.
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// Delete removes a file.
	Delete(ctx context.Context, key string) error

	// GetURL returns a URL for accessing a file.
	GetURL(ctx context.Context, key string) (string, error)

	// Close releases any resources held by the storage.
	Close() error
}

// UploadRequest contains parameters for file upload.
type UploadRequest struct {
	// Key is the storage key/path for the file.
	Key string

	// Content is the file content.
	Content io.Reader

	// ContentType is the MIME type.
	ContentType string

	// Metadata is optional metadata to store with the file.
	Metadata map[string]string
}

// UploadResult contains the result of a file upload.
type UploadResult struct {
	// Key is the storage key.
	Key string

	// URL is the public URL to access the file.
	URL string

	// Size is the size in bytes.
	Size int64
}
