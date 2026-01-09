package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

// BulletHandler handles bullet-related HTTP requests.
type BulletHandler struct {
	bulletService *services.BulletService
}

// NewBulletHandler creates a new BulletHandler.
func NewBulletHandler(bulletService *services.BulletService) *BulletHandler {
	return &BulletHandler{
		bulletService: bulletService,
	}
}

// Create creates a new bullet under an experience.
//
//	@Summary		Create bullet
//	@Description	Creates a new bullet point under an experience
//	@Tags			bullets
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			experienceID	path		string				true	"Experience ID"
//	@Param			request			body		CreateBulletRequest	true	"Bullet data"
//	@Success		201				{object}	BulletResponse
//	@Failure		400				{object}	ErrorResponse	"Invalid request body"
//	@Failure		401				{object}	ErrorResponse	"Unauthorized"
//	@Failure		404				{object}	ErrorResponse	"Experience not found"
//	@Failure		422				{object}	ErrorResponse	"Validation failed"
//	@Failure		500				{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/experiences/{experienceID}/bullets [post]
func (h *BulletHandler) Create(w http.ResponseWriter, r *http.Request) {
	_, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	experienceID := chi.URLParam(r, "experienceID")
	if experienceID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Experience ID is required")
		return
	}

	var req CreateBulletRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Build service request
	displayOrder := 0
	if req.DisplayOrder != nil {
		displayOrder = *req.DisplayOrder
	}
	createReq := services.CreateBulletRequest{
		ExperienceID: experienceID,
		Content:      req.Content,
		Keywords:     req.Keywords,
		DisplayOrder: displayOrder,
	}

	bullet, err := h.bulletService.CreateBullet(r.Context(), createReq)
	if err != nil {
		if errors.Is(err, domain.ErrExperienceNotFound) {
			respondError(w, http.StatusNotFound, "EXPERIENCE_NOT_FOUND", "Experience not found")
			return
		}
		if handleValidationError(w, err) {
			return
		}
		log.Error().Err(err).Str("experience_id", experienceID).Msg("Failed to create bullet")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create bullet")
		return
	}

	response := mapBulletToResponse(bullet)
	respondJSON(w, http.StatusCreated, response)
}

// Update updates an existing bullet.
//
//	@Summary		Update bullet
//	@Description	Updates an existing bullet
//	@Tags			bullets
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			bulletID	path		string				true	"Bullet ID"
//	@Param			request		body		UpdateBulletRequest	true	"Bullet data"
//	@Success		200			{object}	BulletResponse
//	@Failure		400			{object}	ErrorResponse	"Invalid request body"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Bullet not found"
//	@Failure		422			{object}	ErrorResponse	"Validation failed"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/bullets/{bulletID} [put]
func (h *BulletHandler) Update(w http.ResponseWriter, r *http.Request) {
	_, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	bulletID := chi.URLParam(r, "bulletID")
	if bulletID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Bullet ID is required")
		return
	}

	var req UpdateBulletRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Build update request
	updateReq := services.UpdateBulletRequest{
		BulletID:     bulletID,
		Content:      req.Content,
		Keywords:     req.Keywords,
		DisplayOrder: req.DisplayOrder,
	}

	bullet, err := h.bulletService.UpdateBullet(r.Context(), updateReq)
	if err != nil {
		if errors.Is(err, domain.ErrBulletNotFound) {
			respondError(w, http.StatusNotFound, "BULLET_NOT_FOUND", "Bullet not found")
			return
		}
		if handleValidationError(w, err) {
			return
		}
		log.Error().Err(err).Str("bullet_id", bulletID).Msg("Failed to update bullet")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update bullet")
		return
	}

	response := mapBulletToResponse(bullet)
	respondJSON(w, http.StatusOK, response)
}

// Delete removes a bullet.
//
//	@Summary		Delete bullet
//	@Description	Deletes a bullet
//	@Tags			bullets
//	@Security		BearerAuth
//	@Param			bulletID	path	string	true	"Bullet ID"
//	@Success		204			"No content"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Bullet not found"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/bullets/{bulletID} [delete]
func (h *BulletHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	bulletID := chi.URLParam(r, "bulletID")
	if bulletID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Bullet ID is required")
		return
	}

	if err := h.bulletService.DeleteBullet(r.Context(), bulletID); err != nil {
		if errors.Is(err, domain.ErrBulletNotFound) {
			respondError(w, http.StatusNotFound, "BULLET_NOT_FOUND", "Bullet not found")
			return
		}
		log.Error().Err(err).Str("bullet_id", bulletID).Msg("Failed to delete bullet")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete bullet")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RecalculateScore triggers AI recalculation of the bullet's impact score.
//
//	@Summary		Recalculate bullet score
//	@Description	Triggers AI recalculation of the bullet's impact score against a job description
//	@Tags			bullets
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			bulletID	path		string					true	"Bullet ID"
//	@Param			request		body		AnalyzeBulletRequest	true	"Job description for analysis"
//	@Success		200			{object}	ScoreBulletResponse
//	@Failure		400			{object}	ErrorResponse	"Invalid request body"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Bullet not found"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/bullets/{bulletID}/score [post]
func (h *BulletHandler) RecalculateScore(w http.ResponseWriter, r *http.Request) {
	_, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	bulletID := chi.URLParam(r, "bulletID")
	if bulletID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Bullet ID is required")
		return
	}

	var req AnalyzeBulletRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	analyzeReq := services.AnalyzeBulletImpactRequest{
		BulletID:       bulletID,
		JobDescription: req.JobDescription,
	}

	bullet, err := h.bulletService.AnalyzeBulletImpact(r.Context(), analyzeReq)
	if err != nil {
		if errors.Is(err, domain.ErrBulletNotFound) {
			respondError(w, http.StatusNotFound, "BULLET_NOT_FOUND", "Bullet not found")
			return
		}
		log.Error().Err(err).Str("bullet_id", bulletID).Msg("Failed to recalculate bullet score")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to recalculate score")
		return
	}

	response := ScoreBulletResponse{
		ID:          bullet.ID,
		Content:     bullet.Content,
		ImpactScore: bullet.ImpactScore.Int(),
	}

	respondJSON(w, http.StatusOK, response)
}

// mapBulletToResponse maps a domain Bullet to a BulletResponse.
func mapBulletToResponse(b *domain.Bullet) BulletResponse {
	return BulletResponse{
		ID:           b.ID,
		ExperienceID: b.ExperienceID,
		Content:      b.Content,
		ImpactScore:  b.ImpactScore.Int(),
		Keywords:     b.Keywords,
		Metadata:     b.Metadata,
		DisplayOrder: b.DisplayOrder,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}
}
