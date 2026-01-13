package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

// EducationHandler handles education-related HTTP requests.
type EducationHandler struct {
	educationService *services.EducationService
}

// NewEducationHandler creates a new EducationHandler.
func NewEducationHandler(educationService *services.EducationService) *EducationHandler {
	return &EducationHandler{
		educationService: educationService,
	}
}

// List returns all education entries for the authenticated user.
//
//	@Summary		List education entries
//	@Description	Returns all education entries for the authenticated user, ordered by display_order
//	@Tags			education
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	ListEducationResponse
//	@Failure		401	{object}	ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/education [get]
func (h *EducationHandler) List(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	req := services.ListEducationRequest{
		UserID: authUser.ID,
	}

	education, err := h.educationService.ListEducation(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to list education")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve education")
		return
	}

	data := make([]EducationResponse, 0, len(education))
	for _, edu := range education {
		data = append(data, mapEducationToResponse(&edu))
	}

	respondJSON(w, http.StatusOK, ListEducationResponse{
		Data:  data,
		Total: len(data),
	})
}

// Create creates a new education entry.
//
//	@Summary		Create education entry
//	@Description	Creates a new education entry for the authenticated user
//	@Tags			education
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateEducationRequest	true	"Education data"
//	@Success		201		{object}	EducationResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid request body"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		422		{object}	ErrorResponse	"Validation failed"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/education [post]
func (h *EducationHandler) Create(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var req CreateEducationRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Validate required fields.
	if req.Institution == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Institution is required")
		return
	}
	if req.Degree == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Degree is required")
		return
	}

	// Parse dates.
	var startDate, endDate *domain.Date
	if req.StartDate != nil {
		d, err := domain.ParseDate(*req.StartDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid start_date format")
			return
		}
		startDate = &d
	}
	if req.EndDate != nil {
		d, err := domain.ParseDate(*req.EndDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid end_date format")
			return
		}
		endDate = &d
	}

	// Convert non-pointer strings to pointers for optional fields.
	var fieldOfStudy *string
	if req.FieldOfStudy != "" {
		fieldOfStudy = &req.FieldOfStudy
	}

	svcReq := services.CreateEducationRequest{
		UserID:       authUser.ID,
		Institution:  req.Institution,
		Degree:       req.Degree,
		FieldOfStudy: fieldOfStudy,
		Location:     req.Location,
		StartDate:    startDate,
		EndDate:      endDate,
		GPA:          req.GPA,
		Honors:       req.Honors,
		DisplayOrder: req.DisplayOrder,
	}

	education, err := h.educationService.CreateEducation(r.Context(), svcReq)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to create education")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create education")
		return
	}

	respondJSON(w, http.StatusCreated, mapEducationToResponse(education))
}

// Get retrieves a single education entry by ID.
//
//	@Summary		Get education entry
//	@Description	Retrieves a specific education entry by ID
//	@Tags			education
//	@Produce		json
//	@Security		BearerAuth
//	@Param			educationID	path		string	true	"Education ID"
//	@Success		200			{object}	EducationResponse
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Education not found"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/education/{educationID} [get]
func (h *EducationHandler) Get(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	educationID := chi.URLParam(r, "educationID")
	if educationID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Education ID is required")
		return
	}

	education, err := h.educationService.GetEducation(r.Context(), educationID)
	if err != nil {
		if errors.Is(err, domain.ErrEducationNotFound) {
			respondError(w, http.StatusNotFound, "EDUCATION_NOT_FOUND", "Education entry not found")
			return
		}
		log.Error().Err(err).Str("education_id", educationID).Msg("Failed to get education")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve education")
		return
	}

	// Verify ownership.
	if education.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "EDUCATION_NOT_FOUND", "Education entry not found")
		return
	}

	respondJSON(w, http.StatusOK, mapEducationToResponse(education))
}

// Update updates an existing education entry.
//
//	@Summary		Update education entry
//	@Description	Updates an existing education entry
//	@Tags			education
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			educationID	path		string					true	"Education ID"
//	@Param			request		body		UpdateEducationRequest	true	"Education data"
//	@Success		200			{object}	EducationResponse
//	@Failure		400			{object}	ErrorResponse	"Invalid request body"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Education not found"
//	@Failure		422			{object}	ErrorResponse	"Validation failed"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/education/{educationID} [put]
func (h *EducationHandler) Update(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	educationID := chi.URLParam(r, "educationID")
	if educationID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Education ID is required")
		return
	}

	// Verify ownership.
	existing, err := h.educationService.GetEducation(r.Context(), educationID)
	if err != nil {
		if errors.Is(err, domain.ErrEducationNotFound) {
			respondError(w, http.StatusNotFound, "EDUCATION_NOT_FOUND", "Education entry not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve education")
		return
	}
	if existing.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "EDUCATION_NOT_FOUND", "Education entry not found")
		return
	}

	var req UpdateEducationRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Parse dates.
	var startDate, endDate *domain.Date
	if req.StartDate != nil {
		d, err := domain.ParseDate(*req.StartDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid start_date format")
			return
		}
		startDate = &d
	}
	if req.EndDate != nil {
		d, err := domain.ParseDate(*req.EndDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid end_date format")
			return
		}
		endDate = &d
	}

	svcReq := services.UpdateEducationRequest{
		EducationID:  educationID,
		Institution:  req.Institution,
		Degree:       req.Degree,
		FieldOfStudy: req.FieldOfStudy,
		Location:     req.Location,
		StartDate:    startDate,
		EndDate:      endDate,
		GPA:          req.GPA,
		Honors:       req.Honors,
		DisplayOrder: req.DisplayOrder,
	}

	education, err := h.educationService.UpdateEducation(r.Context(), svcReq)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		log.Error().Err(err).Str("education_id", educationID).Msg("Failed to update education")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update education")
		return
	}

	respondJSON(w, http.StatusOK, mapEducationToResponse(education))
}

// Delete removes an education entry.
//
//	@Summary		Delete education entry
//	@Description	Deletes an education entry
//	@Tags			education
//	@Produce		json
//	@Security		BearerAuth
//	@Param			educationID	path	string	true	"Education ID"
//	@Success		204			"No Content"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Education not found"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/education/{educationID} [delete]
func (h *EducationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	educationID := chi.URLParam(r, "educationID")
	if educationID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Education ID is required")
		return
	}

	// Verify ownership.
	existing, err := h.educationService.GetEducation(r.Context(), educationID)
	if err != nil {
		if errors.Is(err, domain.ErrEducationNotFound) {
			respondError(w, http.StatusNotFound, "EDUCATION_NOT_FOUND", "Education entry not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve education")
		return
	}
	if existing.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "EDUCATION_NOT_FOUND", "Education entry not found")
		return
	}

	if err := h.educationService.DeleteEducation(r.Context(), educationID); err != nil {
		log.Error().Err(err).Str("education_id", educationID).Msg("Failed to delete education")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete education")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// mapEducationToResponse maps a domain education to a response DTO.
func mapEducationToResponse(education *domain.Education) EducationResponse {
	resp := EducationResponse{
		ID:           education.ID,
		Institution:  education.Institution,
		Degree:       education.Degree,
		Location:     education.Location,
		GPA:          education.GPA,
		Honors:       education.Honors,
		DisplayOrder: education.DisplayOrder,
		CreatedAt:    education.CreatedAt,
		UpdatedAt:    education.UpdatedAt,
	}

	// Convert *string to string.
	if education.FieldOfStudy != nil {
		resp.FieldOfStudy = *education.FieldOfStudy
	}

	if education.StartDate != nil {
		s := education.StartDate.String()
		resp.StartDate = &s
	}
	if education.EndDate != nil {
		s := education.EndDate.String()
		resp.EndDate = &s
	}

	if resp.Honors == nil {
		resp.Honors = make([]string, 0)
	}

	return resp
}
