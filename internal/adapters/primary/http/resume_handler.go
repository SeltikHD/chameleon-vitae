package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

// ResumeHandler handles resume-related HTTP requests.
type ResumeHandler struct {
	resumeService *services.ResumeService
}

// NewResumeHandler creates a new ResumeHandler.
func NewResumeHandler(resumeService *services.ResumeService) *ResumeHandler {
	return &ResumeHandler{
		resumeService: resumeService,
	}
}

// List returns all resumes for the authenticated user.
//
//	@Summary		List resumes
//	@Description	Returns a paginated list of resumes for the authenticated user
//	@Tags			resumes
//	@Produce		json
//	@Security		BearerAuth
//	@Param			status	query		string	false	"Filter by status"
//	@Param			limit	query		int		false	"Pagination limit"	default(20)
//	@Param			offset	query		int		false	"Pagination offset"	default(0)
//	@Success		200		{object}	ListResumesResponse
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/resumes [get]
func (h *ResumeHandler) List(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	status := r.URL.Query().Get("status")
	limit := parseIntParam(r, "limit", 20)
	offset := parseIntParam(r, "offset", 0)

	listReq := services.ListResumesRequest{
		UserID: authUser.ID,
		Limit:  limit,
		Offset: offset,
	}
	if status != "" {
		listReq.Status = &status
	}

	result, err := h.resumeService.ListResumes(r.Context(), listReq)
	if err != nil {
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to list resumes")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve resumes")
		return
	}

	data := make([]ResumeResponse, 0, len(result.Resumes))
	for _, resume := range result.Resumes {
		data = append(data, mapResumeToResponse(&resume))
	}

	respondJSON(w, http.StatusOK, ListResumesResponse{
		Data:   data,
		Total:  result.Total,
		Limit:  limit,
		Offset: offset,
	})
}

// Get returns a specific resume by ID.
//
//	@Summary		Get resume
//	@Description	Returns a specific resume with full content
//	@Tags			resumes
//	@Produce		json
//	@Security		BearerAuth
//	@Param			resumeID	path		string	true	"Resume ID"
//	@Success		200			{object}	ResumeResponse
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Resume not found"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/resumes/{resumeID} [get]
func (h *ResumeHandler) Get(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	resumeID := chi.URLParam(r, "resumeID")
	if resumeID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Resume ID is required")
		return
	}

	resume, err := h.resumeService.GetResume(r.Context(), resumeID)
	if err != nil {
		if errors.Is(err, domain.ErrResumeNotFound) {
			respondError(w, http.StatusNotFound, "RESUME_NOT_FOUND", "Resume not found")
			return
		}
		log.Error().Err(err).Str("resume_id", resumeID).Msg("Failed to get resume")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve resume")
		return
	}

	// Verify ownership.
	if resume.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "RESUME_NOT_FOUND", "Resume not found")
		return
	}

	response := mapResumeToResponse(resume)
	respondJSON(w, http.StatusOK, response)
}

// Create creates a new resume draft from a job description.
//
//	@Summary		Create resume
//	@Description	Creates a new resume draft from a job description
//	@Tags			resumes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateResumeRequest	true	"Resume data"
//	@Success		201		{object}	ResumeResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid request body"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		422		{object}	ErrorResponse	"Validation failed"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/resumes [post]
func (h *ResumeHandler) Create(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var req CreateResumeRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	createReq := services.CreateResumeRequest{
		UserID:         authUser.ID,
		JobDescription: req.JobDescription,
		TargetLanguage: req.TargetLanguage,
	}

	if req.JobTitle != "" {
		createReq.JobTitle = &req.JobTitle
	}
	if req.CompanyName != "" {
		createReq.CompanyName = &req.CompanyName
	}
	if req.JobURL != "" {
		createReq.JobURL = &req.JobURL
	}

	resume, err := h.resumeService.CreateResume(r.Context(), createReq)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to create resume")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create resume")
		return
	}

	response := mapResumeToResponse(resume)
	respondJSON(w, http.StatusCreated, response)
}

// Tailor triggers AI to analyze the job and generate tailored content.
//
//	@Summary		Tailor resume
//	@Description	Uses AI to select and rewrite bullets for a specific job description
//	@Tags			resumes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			resumeID	path		string				true	"Resume ID"
//	@Param			request		body		TailorResumeRequest	false	"Tailoring parameters"
//	@Success		200			{object}	ResumeResponse
//	@Failure		400			{object}	ErrorResponse	"Invalid request body"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Resume not found"
//	@Failure		422			{object}	ErrorResponse	"Validation failed"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/resumes/{resumeID}/tailor [post]
func (h *ResumeHandler) Tailor(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	resumeID := chi.URLParam(r, "resumeID")
	if resumeID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Resume ID is required")
		return
	}

	// Verify ownership first.
	existing, err := h.resumeService.GetResume(r.Context(), resumeID)
	if err != nil {
		if errors.Is(err, domain.ErrResumeNotFound) {
			respondError(w, http.StatusNotFound, "RESUME_NOT_FOUND", "Resume not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to verify resume")
		return
	}
	if existing.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "RESUME_NOT_FOUND", "Resume not found")
		return
	}

	var req TailorResumeRequest
	if r.Body != nil && r.ContentLength > 0 {
		if err := decodeJSON(r, &req); err != nil {
			respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
			return
		}
	}

	tailorReq := services.TailorResumeRequest{
		ResumeID:   resumeID,
		MaxBullets: req.MaxBulletsPerJob,
	}

	resume, err := h.resumeService.TailorResume(r.Context(), tailorReq)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		if errors.Is(err, domain.ErrNoBulletsAvailable) {
			respondError(w, http.StatusUnprocessableEntity, "NO_BULLETS", "No bullets available for tailoring")
			return
		}
		log.Error().Err(err).Str("resume_id", resumeID).Msg("Failed to tailor resume")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to tailor resume")
		return
	}

	response := mapResumeToResponse(resume)
	respondJSON(w, http.StatusOK, response)
}

// UpdateStatus updates the status of a resume.
//
//	@Summary		Update resume status/content
//	@Description	Updates the status or content of a resume for manual adjustments
//	@Tags			resumes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			resumeID	path		string						true	"Resume ID"
//	@Param			request		body		UpdateResumeContentRequest	true	"Resume update"
//	@Success		200			{object}	ResumeResponse
//	@Failure		400			{object}	ErrorResponse	"Invalid request body"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Resume not found"
//	@Failure		422			{object}	ErrorResponse	"Validation failed"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/resumes/{resumeID}/content [patch]
func (h *ResumeHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	resumeID := chi.URLParam(r, "resumeID")
	if resumeID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Resume ID is required")
		return
	}

	// Verify ownership first.
	existing, err := h.resumeService.GetResume(r.Context(), resumeID)
	if err != nil {
		if errors.Is(err, domain.ErrResumeNotFound) {
			respondError(w, http.StatusNotFound, "RESUME_NOT_FOUND", "Resume not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to verify resume")
		return
	}
	if existing.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "RESUME_NOT_FOUND", "Resume not found")
		return
	}

	var req UpdateResumeContentRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Status == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Status is required")
		return
	}

	updateReq := services.UpdateResumeStatusRequest{
		ResumeID:  resumeID,
		NewStatus: req.Status,
		Notes:     req.Notes,
	}

	resume, err := h.resumeService.UpdateResumeStatus(r.Context(), updateReq)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		log.Error().Err(err).Str("resume_id", resumeID).Msg("Failed to update resume status")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update resume")
		return
	}

	response := mapResumeToResponse(resume)
	respondJSON(w, http.StatusOK, response)
}

// GeneratePDF generates a PDF of the resume.
//
//	@Summary		Generate PDF
//
// GeneratePDF generates and returns a PDF file of the resume.
//
//	@Summary		Generate PDF
//	@Description	Generates and downloads a PDF file of the resume
//	@Tags			resumes
//	@Produce		application/pdf
//	@Security		BearerAuth
//	@Param			resumeID			path		string	true	"Resume ID"
//	@Param			template			query		string	false	"Template name"	default(modern)
//	@Param			force_regenerate	query		bool	false	"Force regeneration ignoring cache"	default(false)
//	@Success		200					{file}		binary	"PDF file"
//	@Failure		401					{object}	ErrorResponse	"Unauthorized"
//	@Failure		404					{object}	ErrorResponse	"Resume not found"
//	@Failure		422					{object}	ErrorResponse	"Resume not ready for PDF"
//	@Failure		500					{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/resumes/{resumeID}/pdf [get]
func (h *ResumeHandler) GeneratePDF(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	resumeID := chi.URLParam(r, "resumeID")
	if resumeID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Resume ID is required")
		return
	}

	// Verify ownership first.
	existing, err := h.resumeService.GetResume(r.Context(), resumeID)
	if err != nil {
		if errors.Is(err, domain.ErrResumeNotFound) {
			respondError(w, http.StatusNotFound, "RESUME_NOT_FOUND", "Resume not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to verify resume")
		return
	}
	if existing.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "RESUME_NOT_FOUND", "Resume not found")
		return
	}

	template := r.URL.Query().Get("template")
	if template == "" {
		template = "modern"
	}

	// Check for force_regenerate query parameter.
	forceRegenerate := r.URL.Query().Get("force_regenerate") == "true"

	pdfReq := services.DownloadPDFRequest{
		ResumeID:        resumeID,
		TemplateName:    template,
		ForceRegenerate: forceRegenerate,
	}

	result, err := h.resumeService.DownloadPDF(r.Context(), pdfReq)
	if err != nil {
		if errors.Is(err, domain.ErrResumeNotReady) {
			respondError(w, http.StatusUnprocessableEntity, "RESUME_NOT_READY", "Resume content must be generated before PDF")
			return
		}
		log.Error().Err(err).Str("resume_id", resumeID).Msg("Failed to generate PDF")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate PDF")
		return
	}

	// Set headers for binary PDF download.
	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+result.Filename+"\"")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(result.Content)))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusOK)

	// Write raw PDF bytes directly to the response.
	_, writeErr := w.Write(result.Content)
	if writeErr != nil {
		log.Error().Err(writeErr).Str("resume_id", resumeID).Msg("Failed to write PDF response")
	}
}

// Delete removes a resume.
//
//	@Summary		Delete resume
//	@Description	Deletes a resume
//	@Tags			resumes
//	@Security		BearerAuth
//	@Param			resumeID	path	string	true	"Resume ID"
//	@Success		204			"No content"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Resume not found"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/resumes/{resumeID} [delete]
func (h *ResumeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	resumeID := chi.URLParam(r, "resumeID")
	if resumeID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Resume ID is required")
		return
	}

	// Verify ownership first.
	existing, err := h.resumeService.GetResume(r.Context(), resumeID)
	if err != nil {
		if errors.Is(err, domain.ErrResumeNotFound) {
			respondError(w, http.StatusNotFound, "RESUME_NOT_FOUND", "Resume not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to verify resume")
		return
	}
	if existing.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "RESUME_NOT_FOUND", "Resume not found")
		return
	}

	if err := h.resumeService.DeleteResume(r.Context(), resumeID); err != nil {
		log.Error().Err(err).Str("resume_id", resumeID).Msg("Failed to delete resume")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete resume")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// mapResumeToResponse maps a domain Resume to a ResumeResponse.
func mapResumeToResponse(resume *domain.Resume) ResumeResponse {
	resp := ResumeResponse{
		ID:              resume.ID,
		JobDescription:  resume.JobDescription,
		TargetLanguage:  resume.TargetLanguage,
		SelectedBullets: resume.SelectedBullets,
		Score:           resume.Score.Int(),
		Status:          string(resume.Status),
		CreatedAt:       resume.CreatedAt,
		UpdatedAt:       resume.UpdatedAt,
	}

	if resume.JobTitle != nil {
		resp.JobTitle = *resume.JobTitle
	}
	if resume.CompanyName != nil {
		resp.CompanyName = *resume.CompanyName
	}
	if resume.JobURL != nil {
		resp.JobURL = *resume.JobURL
	}
	if resume.PDFURL != nil {
		resp.PDFURL = *resume.PDFURL
	}
	if resume.Notes != nil {
		resp.Notes = *resume.Notes
	}
	if resume.GeneratedContent != nil {
		resp.GeneratedContent = mapResumeContentToDTO(resume.GeneratedContent)
	}

	return resp
}

// mapResumeContentToDTO maps domain ResumeContent to ResumeContentDTO.
func mapResumeContentToDTO(content *domain.ResumeContent) *ResumeContentDTO {
	if content == nil {
		return nil
	}

	experiences := make([]TailoredExperienceDTO, 0, len(content.Experiences))
	for _, exp := range content.Experiences {
		bullets := make([]TailoredBulletDTO, 0, len(exp.Bullets))
		for _, b := range exp.Bullets {
			bullets = append(bullets, TailoredBulletDTO{
				BulletID:        b.BulletID,
				OriginalContent: b.OriginalContent,
				TailoredContent: b.TailoredContent,
			})
		}
		experiences = append(experiences, TailoredExperienceDTO{
			ExperienceID: exp.ExperienceID,
			Title:        exp.Title,
			Organization: exp.Organization,
			StartDate:    exp.StartDate,
			EndDate:      exp.EndDate,
			IsCurrent:    exp.IsCurrent,
			Bullets:      bullets,
		})
	}

	dto := &ResumeContentDTO{
		Summary:     content.Summary,
		Experiences: experiences,
		Skills:      content.Skills,
	}

	if content.Analysis != nil {
		dto.Analysis = &ResumeAnalysisDTO{
			MatchedKeywords: content.Analysis.MatchedKeywords,
			MissingKeywords: content.Analysis.MissingKeywords,
			Recommendations: content.Analysis.Recommendations,
		}
	}

	return dto
}
