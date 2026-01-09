package http

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

// UserHandler handles user profile HTTP requests.
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetMe returns the authenticated user's profile.
//
//	@Summary		Get current user profile
//	@Description	Returns the authenticated user's complete profile
//	@Tags			user
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	UserResponse
//	@Failure		401	{object}	ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	user, err := h.userService.GetUser(r.Context(), authUser.ID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			respondError(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
			return
		}
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to get user")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve user profile")
		return
	}

	response := mapUserToResponse(user)
	respondJSON(w, http.StatusOK, response)
}

// UpdateMe updates the authenticated user's profile.
//
//	@Summary		Update current user profile
//	@Description	Updates the authenticated user's profile fields
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		UpdateUserRequest	true	"Update data"
//	@Success		200		{object}	UserResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid request body"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		422		{object}	ErrorResponse	"Validation failed"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/me [patch]
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var req UpdateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	updateReq := services.UpdateProfileRequest{
		UserID:            authUser.ID,
		Name:              req.Name,
		Headline:          req.Headline,
		Summary:           req.Summary,
		Location:          req.Location,
		Phone:             req.Phone,
		Website:           req.Website,
		LinkedInURL:       req.LinkedInURL,
		GitHubURL:         req.GitHubURL,
		PortfolioURL:      req.PortfolioURL,
		PreferredLanguage: req.PreferredLanguage,
	}

	user, err := h.userService.UpdateProfile(r.Context(), updateReq)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		if errors.Is(err, domain.ErrUserNotFound) {
			respondError(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
			return
		}
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to update user")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update user profile")
		return
	}

	response := mapUserToResponse(user)
	respondJSON(w, http.StatusOK, response)
}

// mapUserToResponse maps a domain User to a UserResponse.
func mapUserToResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:                user.ID,
		FirebaseUID:       user.FirebaseUID,
		PictureURL:        user.PictureURL,
		Email:             user.Email,
		Name:              user.Name,
		Headline:          user.Headline,
		Summary:           user.Summary,
		Location:          user.Location,
		Phone:             user.Phone,
		Website:           user.Website,
		LinkedInURL:       user.LinkedInURL,
		GitHubURL:         user.GitHubURL,
		PortfolioURL:      user.PortfolioURL,
		PreferredLanguage: user.PreferredLanguage,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}
}
