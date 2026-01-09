package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

// SkillHandler handles skill-related HTTP requests.
type SkillHandler struct {
	skillService *services.SkillService
}

// NewSkillHandler creates a new SkillHandler.
func NewSkillHandler(skillService *services.SkillService) *SkillHandler {
	return &SkillHandler{
		skillService: skillService,
	}
}

// List returns all skills for the authenticated user.
//
//	@Summary		List skills
//	@Description	Returns all skills for the authenticated user, optionally filtered by category
//	@Tags			skills
//	@Produce		json
//	@Security		BearerAuth
//	@Param			category	query		string	false	"Filter by category"
//	@Success		200			{object}	ListSkillsResponse
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/skills [get]
func (h *SkillHandler) List(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	category := r.URL.Query().Get("category")

	req := services.ListSkillsRequest{
		UserID: authUser.ID,
	}
	if category != "" {
		req.Category = &category
	}

	skills, err := h.skillService.ListSkills(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to list skills")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve skills")
		return
	}

	data := make([]SkillResponse, 0, len(skills))
	for _, skill := range skills {
		data = append(data, mapSkillToResponse(&skill))
	}

	respondJSON(w, http.StatusOK, ListSkillsResponse{
		Data:  data,
		Total: len(data),
	})
}

// BatchUpsert creates or updates multiple skills at once.
//
//	@Summary		Batch upsert skills
//	@Description	Creates or updates multiple skills. Skills are matched by name - existing skills are updated, new ones are created.
//	@Tags			skills
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		BatchUpsertSkillsRequest	true	"Skills data"
//	@Success		200		{object}	BatchUpsertSkillsResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid request body"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		422		{object}	ErrorResponse	"Validation failed"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/skills/batch [post]
func (h *SkillHandler) BatchUpsert(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var req BatchUpsertSkillsRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if len(req.Skills) == 0 {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "At least one skill is required")
		return
	}

	// Convert to service request
	skillInputs := make([]services.CreateSkillRequest, 0, len(req.Skills))
	for i, s := range req.Skills {
		skillReq := services.CreateSkillRequest{
			UserID:       authUser.ID,
			Name:         s.Name,
			Category:     s.Category,
			DisplayOrder: i,
		}
		if s.ProficiencyLevel != nil {
			skillReq.ProficiencyLevel = *s.ProficiencyLevel
		}
		if s.YearsOfExperience != nil {
			skillReq.YearsOfExperience = s.YearsOfExperience
		}
		if s.IsHighlighted != nil && *s.IsHighlighted {
			skillReq.IsHighlighted = true
		}
		skillInputs = append(skillInputs, skillReq)
	}

	batchReq := services.BatchUpsertSkillsRequest{
		UserID: authUser.ID,
		Skills: skillInputs,
	}

	result, err := h.skillService.BatchUpsertSkills(r.Context(), batchReq)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to batch upsert skills")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to upsert skills")
		return
	}

	// Re-fetch skills to return updated list
	listReq := services.ListSkillsRequest{UserID: authUser.ID}
	skills, err := h.skillService.ListSkills(r.Context(), listReq)
	if err != nil {
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to fetch skills after upsert")
		// Still return success with just counts
		respondJSON(w, http.StatusOK, BatchUpsertSkillsResponse{
			Created: result.Created,
			Updated: result.Updated,
			Data:    []SkillResponse{},
		})
		return
	}

	data := make([]SkillResponse, 0, len(skills))
	for _, skill := range skills {
		data = append(data, mapSkillToResponse(&skill))
	}

	respondJSON(w, http.StatusOK, BatchUpsertSkillsResponse{
		Created: result.Created,
		Updated: result.Updated,
		Data:    data,
	})
}

// Delete removes a skill.
//
//	@Summary		Delete skill
//	@Description	Deletes a skill
//	@Tags			skills
//	@Security		BearerAuth
//	@Param			skillID	path	string	true	"Skill ID"
//	@Success		204		"No content"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		404		{object}	ErrorResponse	"Skill not found"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/skills/{skillID} [delete]
func (h *SkillHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	skillID := chi.URLParam(r, "skillID")
	if skillID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Skill ID is required")
		return
	}

	if err := h.skillService.DeleteSkill(r.Context(), skillID); err != nil {
		if errors.Is(err, domain.ErrSkillNotFound) {
			respondError(w, http.StatusNotFound, "SKILL_NOT_FOUND", "Skill not found")
			return
		}
		log.Error().Err(err).Str("skill_id", skillID).Msg("Failed to delete skill")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete skill")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// mapSkillToResponse maps a domain Skill to a SkillResponse.
func mapSkillToResponse(s *domain.Skill) SkillResponse {
	return SkillResponse{
		ID:                s.ID,
		Name:              s.Name,
		Category:          s.Category,
		ProficiencyLevel:  s.ProficiencyLevel.Int(),
		YearsOfExperience: s.YearsOfExperience,
		IsHighlighted:     s.IsHighlighted,
		DisplayOrder:      s.DisplayOrder,
		CreatedAt:         s.CreatedAt,
	}
}
