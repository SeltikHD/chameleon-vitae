package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

// SpokenLanguageHandler handles spoken language-related HTTP requests.
type SpokenLanguageHandler struct {
	skillService *services.SkillService
}

// NewSpokenLanguageHandler creates a new SpokenLanguageHandler.
func NewSpokenLanguageHandler(skillService *services.SkillService) *SpokenLanguageHandler {
	return &SpokenLanguageHandler{
		skillService: skillService,
	}
}

// List returns all spoken languages for the authenticated user.
//
//	@Summary		List spoken languages
//	@Description	Returns all spoken languages for the authenticated user
//	@Tags			languages
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	ListSpokenLanguagesResponse
//	@Failure		401	{object}	ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/languages [get]
func (h *SpokenLanguageHandler) List(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	languages, err := h.skillService.ListSpokenLanguages(r.Context(), authUser.ID)
	if err != nil {
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to list spoken languages")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve spoken languages")
		return
	}

	data := make([]SpokenLanguageResponse, 0, len(languages))
	for _, lang := range languages {
		data = append(data, mapSpokenLanguageToResponse(&lang))
	}

	respondJSON(w, http.StatusOK, ListSpokenLanguagesResponse{
		Data: data,
	})
}

// Create creates a new spoken language.
//
//	@Summary		Create spoken language
//	@Description	Creates a new spoken language for the authenticated user
//	@Tags			languages
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateSpokenLanguageRequest	true	"Spoken language data"
//	@Success		201		{object}	SpokenLanguageResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid request body"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		409		{object}	ErrorResponse	"Language already exists"
//	@Failure		422		{object}	ErrorResponse	"Validation failed"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/languages [post]
func (h *SpokenLanguageHandler) Create(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var req CreateSpokenLanguageRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	displayOrder := 0
	if req.DisplayOrder != nil {
		displayOrder = *req.DisplayOrder
	}

	createReq := services.CreateSpokenLanguageRequest{
		UserID:       authUser.ID,
		Language:     req.Language,
		Proficiency:  req.Proficiency,
		DisplayOrder: displayOrder,
	}

	language, err := h.skillService.CreateSpokenLanguage(r.Context(), createReq)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to create spoken language")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create spoken language")
		return
	}

	response := mapSpokenLanguageToResponse(language)
	respondJSON(w, http.StatusCreated, response)
}

// Delete removes a spoken language.
//
//	@Summary		Delete spoken language
//	@Description	Deletes a spoken language
//	@Tags			languages
//	@Security		BearerAuth
//	@Param			languageID	path	string	true	"Spoken language ID"
//	@Success		204			"No content"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Language not found"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/languages/{languageID} [delete]
func (h *SpokenLanguageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	languageID := chi.URLParam(r, "languageID")
	if languageID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Language ID is required")
		return
	}

	if err := h.skillService.DeleteSpokenLanguage(r.Context(), languageID); err != nil {
		if errors.Is(err, domain.ErrSpokenLanguageNotFound) {
			respondError(w, http.StatusNotFound, "LANGUAGE_NOT_FOUND", "Language not found")
			return
		}
		log.Error().Err(err).Str("language_id", languageID).Msg("Failed to delete spoken language")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete spoken language")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// mapSpokenLanguageToResponse maps a domain SpokenLanguage to a SpokenLanguageResponse.
func mapSpokenLanguageToResponse(lang *domain.SpokenLanguage) SpokenLanguageResponse {
	return SpokenLanguageResponse{
		ID:           lang.ID,
		Language:     lang.Language,
		Proficiency:  string(lang.Proficiency),
		DisplayOrder: lang.DisplayOrder,
		CreatedAt:    lang.CreatedAt,
	}
}
