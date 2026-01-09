package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

// ExperienceHandler handles experience-related HTTP requests.
type ExperienceHandler struct {
	experienceService *services.ExperienceService
}

// NewExperienceHandler creates a new ExperienceHandler.
func NewExperienceHandler(experienceService *services.ExperienceService) *ExperienceHandler {
	return &ExperienceHandler{
		experienceService: experienceService,
	}
}

// List returns all experiences for the authenticated user.
//
//	@Summary		List experiences
//	@Description	Returns a paginated list of experiences for the authenticated user
//	@Tags			experiences
//	@Produce		json
//	@Security		BearerAuth
//	@Param			type	query		string	false	"Filter by experience type"
//	@Param			limit	query		int		false	"Pagination limit"	default(50)
//	@Param			offset	query		int		false	"Pagination offset"	default(0)
//	@Success		200		{object}	ListExperiencesResponse
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/experiences [get]
func (h *ExperienceHandler) List(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Parse query parameters
	expType := r.URL.Query().Get("type")
	limit := parseIntParam(r, "limit", 50)
	offset := parseIntParam(r, "offset", 0)

	req := services.ListExperiencesRequest{
		UserID: authUser.ID,
		Limit:  limit,
		Offset: offset,
	}
	if expType != "" {
		req.Type = &expType
	}

	result, err := h.experienceService.ListExperiences(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to list experiences")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve experiences")
		return
	}

	data := make([]ExperienceResponse, 0, len(result.Experiences))
	for _, exp := range result.Experiences {
		data = append(data, mapExperienceToResponse(&exp))
	}

	respondJSON(w, http.StatusOK, ListExperiencesResponse{
		Data:   data,
		Total:  result.Total,
		Limit:  limit,
		Offset: offset,
	})
}

// Get returns a specific experience by ID.
//
//	@Summary		Get experience
//	@Description	Returns a specific experience with all its bullets
//	@Tags			experiences
//	@Produce		json
//	@Security		BearerAuth
//	@Param			experienceID	path		string	true	"Experience ID"
//	@Success		200				{object}	ExperienceResponse
//	@Failure		401				{object}	ErrorResponse	"Unauthorized"
//	@Failure		404				{object}	ErrorResponse	"Experience not found"
//	@Failure		500				{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/experiences/{experienceID} [get]
func (h *ExperienceHandler) Get(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	experienceID := chi.URLParam(r, "experienceID")
	if experienceID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Experience ID is required")
		return
	}

	experience, err := h.experienceService.GetExperienceWithBullets(r.Context(), experienceID)
	if err != nil {
		if errors.Is(err, domain.ErrExperienceNotFound) {
			respondError(w, http.StatusNotFound, "EXPERIENCE_NOT_FOUND", "Experience not found")
			return
		}
		log.Error().Err(err).Str("experience_id", experienceID).Msg("Failed to get experience")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve experience")
		return
	}

	// Verify ownership
	if experience.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "EXPERIENCE_NOT_FOUND", "Experience not found")
		return
	}

	response := mapExperienceToResponse(experience)
	respondJSON(w, http.StatusOK, response)
}

// Create creates a new experience.
//
//	@Summary		Create experience
//	@Description	Creates a new experience for the authenticated user
//	@Tags			experiences
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateExperienceRequest	true	"Experience data"
//	@Success		201		{object}	ExperienceResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid request body"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		422		{object}	ErrorResponse	"Validation failed"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/experiences [post]
func (h *ExperienceHandler) Create(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var req CreateExperienceRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Build service request
	createReq := services.CreateExperienceRequest{
		UserID:       authUser.ID,
		Type:         req.Type,
		Title:        req.Title,
		Organization: req.Organization,
		Location:     req.Location,
		StartDate:    req.StartDate,
		IsCurrent:    req.IsCurrent,
		Description:  req.Description,
		URL:          req.URL,
	}
	if req.EndDate != nil {
		createReq.EndDate = req.EndDate
	}

	experience, err := h.experienceService.CreateExperience(r.Context(), createReq)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to create experience")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create experience")
		return
	}

	response := mapExperienceToResponse(experience)
	respondJSON(w, http.StatusCreated, response)
}

// Update updates an existing experience.
//
//	@Summary		Update experience
//	@Description	Updates an existing experience
//	@Tags			experiences
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			experienceID	path		string					true	"Experience ID"
//	@Param			request			body		UpdateExperienceRequest	true	"Experience data"
//	@Success		200				{object}	ExperienceResponse
//	@Failure		400				{object}	ErrorResponse	"Invalid request body"
//	@Failure		401				{object}	ErrorResponse	"Unauthorized"
//	@Failure		404				{object}	ErrorResponse	"Experience not found"
//	@Failure		422				{object}	ErrorResponse	"Validation failed"
//	@Failure		500				{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/experiences/{experienceID} [put]
func (h *ExperienceHandler) Update(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	experienceID := chi.URLParam(r, "experienceID")
	if experienceID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Experience ID is required")
		return
	}

	var req UpdateExperienceRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Verify ownership first
	existing, err := h.experienceService.GetExperience(r.Context(), experienceID)
	if err != nil {
		if errors.Is(err, domain.ErrExperienceNotFound) {
			respondError(w, http.StatusNotFound, "EXPERIENCE_NOT_FOUND", "Experience not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to verify experience")
		return
	}
	if existing.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "EXPERIENCE_NOT_FOUND", "Experience not found")
		return
	}

	// Build update request
	updateReq := services.UpdateExperienceRequest{
		ExperienceID: experienceID,
		Type:         req.Type,
		Title:        req.Title,
		Organization: req.Organization,
		Location:     req.Location,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		IsCurrent:    req.IsCurrent,
		Description:  req.Description,
		URL:          req.URL,
		DisplayOrder: req.DisplayOrder,
	}

	experience, err := h.experienceService.UpdateExperience(r.Context(), updateReq)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		log.Error().Err(err).Str("experience_id", experienceID).Msg("Failed to update experience")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update experience")
		return
	}

	response := mapExperienceToResponse(experience)
	respondJSON(w, http.StatusOK, response)
}

// Delete removes an experience and all its bullets.
//
//	@Summary		Delete experience
//	@Description	Deletes an experience and all its bullets
//	@Tags			experiences
//	@Security		BearerAuth
//	@Param			experienceID	path	string	true	"Experience ID"
//	@Success		204				"No content"
//	@Failure		401				{object}	ErrorResponse	"Unauthorized"
//	@Failure		404				{object}	ErrorResponse	"Experience not found"
//	@Failure		500				{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/experiences/{experienceID} [delete]
func (h *ExperienceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	experienceID := chi.URLParam(r, "experienceID")
	if experienceID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Experience ID is required")
		return
	}

	// Verify ownership first
	existing, err := h.experienceService.GetExperience(r.Context(), experienceID)
	if err != nil {
		if errors.Is(err, domain.ErrExperienceNotFound) {
			respondError(w, http.StatusNotFound, "EXPERIENCE_NOT_FOUND", "Experience not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to verify experience")
		return
	}
	if existing.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "EXPERIENCE_NOT_FOUND", "Experience not found")
		return
	}

	if err := h.experienceService.DeleteExperience(r.Context(), experienceID); err != nil {
		log.Error().Err(err).Str("experience_id", experienceID).Msg("Failed to delete experience")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete experience")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// mapExperienceToResponse maps a domain Experience to an ExperienceResponse.
func mapExperienceToResponse(exp *domain.Experience) ExperienceResponse {
	response := ExperienceResponse{
		ID:           exp.ID,
		Type:         string(exp.Type),
		Title:        exp.Title,
		Organization: exp.Organization,
		Location:     exp.Location,
		StartDate:    exp.StartDate.String(),
		IsCurrent:    exp.IsCurrent,
		Description:  exp.Description,
		URL:          exp.URL,
		Metadata:     exp.Metadata,
		DisplayOrder: exp.DisplayOrder,
		CreatedAt:    exp.CreatedAt,
		UpdatedAt:    exp.UpdatedAt,
	}

	if exp.EndDate != nil && !exp.EndDate.IsZero() {
		endDate := exp.EndDate.String()
		response.EndDate = &endDate
	}

	// Map bullets
	response.Bullets = make([]BulletResponse, 0, len(exp.Bullets))
	for _, b := range exp.Bullets {
		response.Bullets = append(response.Bullets, mapBulletToResponse(&b))
	}

	return response
}

// parseIntParam parses an integer query parameter with a default value.
func parseIntParam(r *http.Request, name string, defaultVal int) int {
	val := r.URL.Query().Get(name)
	if val == "" {
		return defaultVal
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return intVal
}

// handleValidationError checks if the error is a validation error and responds accordingly.
func handleValidationError(w http.ResponseWriter, err error) bool {
	var validationErr *domain.ValidationErrors
	if errors.As(err, &validationErr) {
		details := make([]ErrorDetail, 0, len(validationErr.Errors))
		for _, fieldErr := range validationErr.Errors {
			details = append(details, ErrorDetail{
				Field:   fieldErr.Field,
				Message: fieldErr.Message,
			})
		}
		respondErrorWithDetails(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Validation failed", details)
		return true
	}
	return false
}
