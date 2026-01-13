// Package services contains the application services (use cases).
package services

import (
	"context"
	"fmt"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

// ResumeService handles resume generation and management use cases.
type ResumeService struct {
	resumeRepo     ports.ResumeRepository
	userRepo       ports.UserRepository
	experienceRepo ports.ExperienceRepository
	bulletRepo     ports.BulletRepository
	skillRepo      ports.SkillRepository
	languageRepo   ports.SpokenLanguageRepository
	educationRepo  ports.EducationRepository
	projectRepo    ports.ProjectRepository
	aiProvider     ports.AIProvider
	pdfEngine      ports.PDFEngine
	jobParser      ports.JobParser
	fileStorage    ports.FileStorage
}

// NewResumeService creates a new ResumeService with required dependencies.
func NewResumeService(
	resumeRepo ports.ResumeRepository,
	userRepo ports.UserRepository,
	experienceRepo ports.ExperienceRepository,
	bulletRepo ports.BulletRepository,
	skillRepo ports.SkillRepository,
	languageRepo ports.SpokenLanguageRepository,
	educationRepo ports.EducationRepository,
	projectRepo ports.ProjectRepository,
	aiProvider ports.AIProvider,
	pdfEngine ports.PDFEngine,
	jobParser ports.JobParser,
	fileStorage ports.FileStorage,
) *ResumeService {
	return &ResumeService{
		resumeRepo:     resumeRepo,
		userRepo:       userRepo,
		experienceRepo: experienceRepo,
		bulletRepo:     bulletRepo,
		skillRepo:      skillRepo,
		languageRepo:   languageRepo,
		educationRepo:  educationRepo,
		projectRepo:    projectRepo,
		aiProvider:     aiProvider,
		pdfEngine:      pdfEngine,
		jobParser:      jobParser,
		fileStorage:    fileStorage,
	}
}

// ParseJobURLRequest contains parameters for parsing a job URL.
type ParseJobURLRequest struct {
	URL string
}

// ParseJobURL parses a job description from a URL using Jina Reader.
func (s *ResumeService) ParseJobURL(ctx context.Context, req ParseJobURLRequest) (*ports.ParsedJob, error) {
	parsedJob, err := s.jobParser.ParseJobURL(ctx, req.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse job URL: %w", err)
	}
	return parsedJob, nil
}

// CreateResumeRequest contains parameters for creating a resume.
type CreateResumeRequest struct {
	UserID         string
	JobDescription string
	JobTitle       *string
	CompanyName    *string
	JobURL         *string
	TargetLanguage string
}

// CreateResume creates a new resume draft.
func (s *ResumeService) CreateResume(ctx context.Context, req CreateResumeRequest) (*domain.Resume, error) {
	// Verify user exists.
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Create resume.
	resume, err := domain.NewResume(req.UserID, req.JobDescription)
	if err != nil {
		return nil, err
	}

	if req.TargetLanguage != "" {
		resume.TargetLanguage = req.TargetLanguage
	}

	if req.JobTitle != nil || req.CompanyName != nil || req.JobURL != nil {
		title := ""
		company := ""
		url := ""
		if req.JobTitle != nil {
			title = *req.JobTitle
		}
		if req.CompanyName != nil {
			company = *req.CompanyName
		}
		if req.JobURL != nil {
			url = *req.JobURL
		}
		resume.SetJobDetails(title, company, url)
	}

	if err := resume.Validate(); err != nil {
		return nil, err
	}

	if err := s.resumeRepo.Create(ctx, resume); err != nil {
		return nil, fmt.Errorf("failed to create resume: %w", err)
	}

	return resume, nil
}

// GetResume retrieves a resume by ID.
func (s *ResumeService) GetResume(ctx context.Context, resumeID string) (*domain.Resume, error) {
	resume, err := s.resumeRepo.GetByID(ctx, resumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get resume: %w", err)
	}
	return resume, nil
}

// ListResumesRequest contains parameters for listing resumes.
type ListResumesRequest struct {
	UserID string
	Status *string
	Limit  int
	Offset int
}

// ListResumesResponse contains the result of listing resumes.
type ListResumesResponse struct {
	Resumes []domain.Resume
	Total   int
}

// ListResumes lists resumes for a user with optional status filter.
func (s *ResumeService) ListResumes(ctx context.Context, req ListResumesRequest) (*ListResumesResponse, error) {
	opts := ports.ListOptions{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	if opts.Limit == 0 {
		opts = ports.DefaultListOptions()
	}

	var resumes []domain.Resume
	var total int
	var err error

	if req.Status != nil && *req.Status != "" {
		status, parseErr := domain.ParseResumeStatus(*req.Status)
		if parseErr != nil {
			return nil, parseErr
		}
		resumes, total, err = s.resumeRepo.ListByUserIDAndStatus(ctx, req.UserID, status, opts)
	} else {
		resumes, total, err = s.resumeRepo.ListByUserID(ctx, req.UserID, opts)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list resumes: %w", err)
	}

	return &ListResumesResponse{
		Resumes: resumes,
		Total:   total,
	}, nil
}

// TailorResumeRequest contains parameters for tailoring a resume.
type TailorResumeRequest struct {
	ResumeID   string
	MaxBullets int
}

// TailorResume generates AI-tailored content for a resume.
func (s *ResumeService) TailorResume(ctx context.Context, req TailorResumeRequest) (*domain.Resume, error) {
	// Get the resume.
	resume, err := s.resumeRepo.GetByID(ctx, req.ResumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get resume: %w", err)
	}

	// Get user profile.
	user, err := s.userRepo.GetByID(ctx, resume.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get all user's bullets.
	allBullets, err := s.bulletRepo.ListByUserID(ctx, resume.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bullets: %w", err)
	}

	if len(allBullets) == 0 {
		return nil, domain.ErrNoBulletsAvailable
	}

	// Get user's skills.
	skills, err := s.skillRepo.ListByUserID(ctx, resume.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}

	// Analyze job description.
	jobAnalysis, err := s.aiProvider.AnalyzeJob(ctx, ports.AnalyzeJobRequest{
		JobDescription: resume.JobDescription,
		TargetLanguage: resume.TargetLanguage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to analyze job: %w", err)
	}

	// Update job details from analysis if not already set.
	if resume.JobTitle == nil && jobAnalysis.Title != "" {
		resume.SetJobDetails(jobAnalysis.Title, jobAnalysis.Company, "")
	}

	// Select the most relevant bullets.
	maxBullets := req.MaxBullets
	if maxBullets == 0 {
		maxBullets = 15 // Default.
	}

	bulletSelection, err := s.aiProvider.SelectBullets(ctx, ports.SelectBulletsRequest{
		JobAnalysis:      jobAnalysis,
		AvailableBullets: allBullets,
		MaxBullets:       maxBullets,
		TargetLanguage:   resume.TargetLanguage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to select bullets: %w", err)
	}

	resume.SelectedBullets = bulletSelection.SelectedBulletIDs

	// Get the selected bullets.
	selectedBullets, err := s.bulletRepo.ListByIDs(ctx, bulletSelection.SelectedBulletIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get selected bullets: %w", err)
	}

	// Tailor each bullet.
	tailoredBulletResults := make([]ports.TailoredBulletResult, 0, len(selectedBullets))
	for _, bullet := range selectedBullets {
		tailored, err := s.aiProvider.TailorBullet(ctx, ports.TailorBulletRequest{
			Bullet:         bullet,
			JobAnalysis:    jobAnalysis,
			TargetLanguage: resume.TargetLanguage,
			Style:          "professional",
		})
		if err != nil {
			// Log error but continue with other bullets.
			continue
		}
		tailoredBulletResults = append(tailoredBulletResults, *tailored)
	}

	// Generate professional summary.
	summaryResult, err := s.aiProvider.GenerateSummary(ctx, ports.GenerateSummaryRequest{
		User:            user,
		JobAnalysis:     jobAnalysis,
		SelectedBullets: selectedBullets,
		TargetLanguage:  resume.TargetLanguage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}

	// Group tailored bullets by experience.
	bulletsByExp := make(map[string][]domain.TailoredBullet)
	for i, bullet := range selectedBullets {
		if i >= len(tailoredBulletResults) {
			break
		}
		tb := domain.TailoredBullet{
			BulletID:        bullet.ID,
			OriginalContent: bullet.Content,
			TailoredContent: tailoredBulletResults[i].TailoredContent,
		}
		bulletsByExp[bullet.ExperienceID] = append(bulletsByExp[bullet.ExperienceID], tb)
	}

	// Get experiences for the selected bullets.
	expIDs := make([]string, 0, len(bulletsByExp))
	for expID := range bulletsByExp {
		expIDs = append(expIDs, expID)
	}

	tailoredExperiences := make([]domain.TailoredExperience, 0, len(expIDs))
	for _, expID := range expIDs {
		exp, err := s.experienceRepo.GetByID(ctx, expID)
		if err != nil {
			continue
		}

		te := domain.TailoredExperience{
			ExperienceID: exp.ID,
			Title:        exp.Title,
			Organization: exp.Organization,
			StartDate:    exp.StartDate.String(),
			IsCurrent:    exp.IsCurrent,
			Bullets:      bulletsByExp[expID],
		}
		if exp.EndDate != nil {
			endStr := exp.EndDate.String()
			te.EndDate = &endStr
		}
		tailoredExperiences = append(tailoredExperiences, te)
	}

	// Build skill list.
	skillNames := make([]string, 0, len(skills))
	for _, skill := range skills {
		skillNames = append(skillNames, skill.Name)
	}

	// Calculate match score.
	matchScore, err := s.aiProvider.ScoreMatch(ctx, ports.ScoreMatchRequest{
		JobAnalysis: jobAnalysis,
		Resume: &domain.ResumeContent{
			Summary:     summaryResult.Summary,
			Experiences: tailoredExperiences,
			Skills:      skillNames,
		},
		UserSkills: skills,
	})
	if err != nil {
		// Use default score if scoring fails.
		defaultScore, _ := domain.NewMatchScore(0)
		matchScore = &defaultScore
	}

	// Build the generated content.
	generatedContent := &domain.ResumeContent{
		Summary:     summaryResult.Summary,
		Experiences: tailoredExperiences,
		Skills:      skillNames,
		Analysis: &domain.ResumeAnalysis{
			MatchedKeywords: jobAnalysis.RequiredSkills,
			MissingKeywords: jobAnalysis.PreferredSkills,
			StrengthAreas:   []string{},
		},
	}

	resume.SetGeneratedContent(generatedContent)
	if err := resume.SetScore(matchScore.Int()); err != nil {
		// Ignore score setting error.
	}

	// Save the updated resume.
	if err := s.resumeRepo.Update(ctx, resume); err != nil {
		return nil, fmt.Errorf("failed to update resume: %w", err)
	}

	return resume, nil
}

// GeneratePDFRequest contains parameters for generating a PDF.
type GeneratePDFRequest struct {
	ResumeID     string
	TemplateName string
}

// GeneratePDF generates a PDF for a resume.
func (s *ResumeService) GeneratePDF(ctx context.Context, req GeneratePDFRequest) (*domain.Resume, error) {
	resume, err := s.resumeRepo.GetByID(ctx, req.ResumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get resume: %w", err)
	}

	if !resume.CanGeneratePDF() {
		return nil, domain.ErrResumeNotReady
	}

	// Get user for personal info.
	user, err := s.userRepo.GetByID(ctx, resume.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get spoken languages.
	languages, err := s.languageRepo.ListByUserID(ctx, resume.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get languages: %w", err)
	}

	// Get education entries.
	education, err := s.educationRepo.ListByUserID(ctx, resume.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get education: %w", err)
	}

	// Get projects with bullets.
	projects, err := s.projectRepo.ListByUserID(ctx, resume.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	// Get user skills for categorization.
	skills, err := s.skillRepo.ListByUserID(ctx, resume.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}

	// Build HTML using Jake's Resume template.
	template := NewJakeResumeTemplate()
	html := template.Render(ResumeTemplateData{
		User:      user,
		Resume:    resume,
		Education: education,
		Projects:  projects,
		Languages: languages,
		Skills:    skills,
		FontSize:  11, // Default to 11pt
	})

	// Generate PDF.
	templateName := req.TemplateName
	if templateName == "" {
		templateName = "jake"
	}

	pdfResult, err := s.pdfEngine.GeneratePDF(ctx, ports.GeneratePDFRequest{
		HTML:         html,
		TemplateName: templateName,
		Options:      ports.DefaultPDFOptions(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	defer pdfResult.Content.Close()

	// Upload PDF to storage.
	filename := fmt.Sprintf("resumes/%s/%s.pdf", resume.UserID, resume.ID)
	uploadResult, err := s.fileStorage.Upload(ctx, ports.UploadRequest{
		Key:         filename,
		Content:     pdfResult.Content,
		ContentType: "application/pdf",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload PDF: %w", err)
	}

	// Update resume with PDF URL.
	resume.PDFURL = &uploadResult.URL
	if err := resume.TransitionStatus(domain.ResumeStatusReviewed); err != nil {
		// Ignore status transition error.
	}

	if err := s.resumeRepo.Update(ctx, resume); err != nil {
		return nil, fmt.Errorf("failed to update resume: %w", err)
	}

	return resume, nil
}

// DownloadPDFRequest contains parameters for downloading a resume PDF.
type DownloadPDFRequest struct {
	ResumeID     string
	TemplateName string
}

// DownloadPDFResult contains the result of downloading a PDF.
type DownloadPDFResult struct {
	Content     []byte
	Filename    string
	ContentType string
}

// DownloadPDF generates (if needed) and returns the PDF bytes for a resume.
func (s *ResumeService) DownloadPDF(ctx context.Context, req DownloadPDFRequest) (*DownloadPDFResult, error) {
	resume, err := s.resumeRepo.GetByID(ctx, req.ResumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get resume: %w", err)
	}

	if !resume.CanGeneratePDF() {
		return nil, domain.ErrResumeNotReady
	}

	// Get user for filename generation.
	user, err := s.userRepo.GetByID(ctx, resume.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if PDF already exists.
	filename := fmt.Sprintf("resumes/%s/%s.pdf", resume.UserID, resume.ID)

	// Try to download existing PDF.
	reader, err := s.fileStorage.Download(ctx, filename)
	if err == nil && reader != nil {
		defer reader.Close()
		content, readErr := readAll(reader)
		if readErr == nil && len(content) > 0 {
			return &DownloadPDFResult{
				Content:     content,
				Filename:    s.generatePDFFilename(user, resume),
				ContentType: "application/pdf",
			}, nil
		}
	}

	// PDF doesn't exist, generate it.
	languages, err := s.languageRepo.ListByUserID(ctx, resume.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get languages: %w", err)
	}

	// Get education entries.
	education, err := s.educationRepo.ListByUserID(ctx, resume.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get education: %w", err)
	}

	// Get projects with bullets.
	projects, err := s.projectRepo.ListByUserID(ctx, resume.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	// Get user skills for categorization.
	skills, err := s.skillRepo.ListByUserID(ctx, resume.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}

	// Build HTML using Jake's Resume template.
	template := NewJakeResumeTemplate()
	htmlContent := template.Render(ResumeTemplateData{
		User:      user,
		Resume:    resume,
		Education: education,
		Projects:  projects,
		Languages: languages,
		Skills:    skills,
		FontSize:  11,
	})

	templateName := req.TemplateName
	if templateName == "" {
		templateName = "jake"
	}

	pdfResult, err := s.pdfEngine.GeneratePDF(ctx, ports.GeneratePDFRequest{
		HTML:         htmlContent,
		TemplateName: templateName,
		Options:      ports.DefaultPDFOptions(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	defer pdfResult.Content.Close()

	// Read PDF content.
	pdfBytes, err := readAll(pdfResult.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF content: %w", err)
	}

	// Upload for caching (best effort, don't fail if upload fails).
	go func() {
		uploadCtx := context.Background()
		_, uploadErr := s.fileStorage.Upload(uploadCtx, ports.UploadRequest{
			Key:         filename,
			Content:     newBytesReader(pdfBytes),
			ContentType: "application/pdf",
		})
		if uploadErr != nil {
			// Log but don't fail.
			fmt.Printf("Warning: failed to cache PDF: %v\n", uploadErr)
		}
	}()

	return &DownloadPDFResult{
		Content:     pdfBytes,
		Filename:    s.generatePDFFilename(user, resume),
		ContentType: "application/pdf",
	}, nil
}

// generatePDFFilename generates a descriptive filename for the PDF.
func (s *ResumeService) generatePDFFilename(user *domain.User, resume *domain.Resume) string {
	name := user.GetDisplayName()
	if resume.CompanyName != nil && *resume.CompanyName != "" {
		return fmt.Sprintf("%s_Resume_%s.pdf", sanitizeFilename(name), sanitizeFilename(*resume.CompanyName))
	}
	if resume.JobTitle != nil && *resume.JobTitle != "" {
		return fmt.Sprintf("%s_Resume_%s.pdf", sanitizeFilename(name), sanitizeFilename(*resume.JobTitle))
	}
	return fmt.Sprintf("%s_Resume.pdf", sanitizeFilename(name))
}

// UpdateResumeStatusRequest contains parameters for updating resume status.
type UpdateResumeStatusRequest struct {
	ResumeID  string
	NewStatus string
	Notes     *string
}

// UpdateResumeStatus updates the status of a resume.
func (s *ResumeService) UpdateResumeStatus(ctx context.Context, req UpdateResumeStatusRequest) (*domain.Resume, error) {
	resume, err := s.resumeRepo.GetByID(ctx, req.ResumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get resume: %w", err)
	}

	newStatus, err := domain.ParseResumeStatus(req.NewStatus)
	if err != nil {
		return nil, err
	}

	if err := resume.TransitionStatus(newStatus); err != nil {
		return nil, err
	}

	if req.Notes != nil {
		resume.Notes = req.Notes
	}

	if err := s.resumeRepo.Update(ctx, resume); err != nil {
		return nil, fmt.Errorf("failed to update resume: %w", err)
	}

	return resume, nil
}

// DeleteResume removes a resume.
func (s *ResumeService) DeleteResume(ctx context.Context, resumeID string) error {
	// Get resume to check for PDF.
	resume, err := s.resumeRepo.GetByID(ctx, resumeID)
	if err != nil {
		return fmt.Errorf("failed to get resume: %w", err)
	}

	// Delete PDF from storage if exists.
	if resume.PDFURL != nil {
		filename := fmt.Sprintf("resumes/%s/%s.pdf", resume.UserID, resume.ID)
		// Ignore delete errors for storage.
		_ = s.fileStorage.Delete(ctx, filename)
	}

	if err := s.resumeRepo.Delete(ctx, resumeID); err != nil {
		return fmt.Errorf("failed to delete resume: %w", err)
	}

	return nil
}

// Helper functions for PDF handling.

// readAll reads all bytes from a reader.
func readAll(r interface{ Read([]byte) (int, error) }) ([]byte, error) {
	var buf []byte
	chunk := make([]byte, 4096)
	for {
		n, err := r.Read(chunk)
		if n > 0 {
			buf = append(buf, chunk[:n]...)
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return buf, err
		}
	}
	return buf, nil
}

// newBytesReader creates a reader from bytes.
func newBytesReader(b []byte) *bytesReader {
	return &bytesReader{data: b}
}

type bytesReader struct {
	data   []byte
	offset int
}

func (r *bytesReader) Read(p []byte) (int, error) {
	if r.offset >= len(r.data) {
		return 0, fmt.Errorf("EOF")
	}
	n := copy(p, r.data[r.offset:])
	r.offset += n
	return n, nil
}

// sanitizeFilename removes or replaces characters that are invalid in filenames.
func sanitizeFilename(name string) string {
	var result []byte
	for i := 0; i < len(name); i++ {
		c := name[i]
		switch c {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			result = append(result, '_')
		case ' ':
			result = append(result, '_')
		default:
			if c >= 32 && c < 127 {
				result = append(result, c)
			}
		}
	}
	return string(result)
}
