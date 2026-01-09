package http

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	userService *services.UserService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// SyncUser synchronizes a Firebase user with the local database.
//
//	@Summary		Sync user from Firebase
//	@Description	Synchronizes a Firebase authenticated user with the local PostgreSQL database. Creates a new user if not exists, updates if exists (upsert behavior).
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Bearer token from Firebase"
//	@Param			request			body		SyncUserRequest	true	"User sync request"
//	@Success		201				{object}	SyncUserResponse
//	@Failure		400				{object}	ErrorResponse	"Invalid request body"
//	@Failure		401				{object}	ErrorResponse	"Invalid or expired token"
//	@Failure		500				{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/auth/sync [post]
func (h *AuthHandler) SyncUser(w http.ResponseWriter, r *http.Request) {
	// Extract the authorization header for token verification
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header")
		return
	}

	// Parse "Bearer <token>"
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization header format")
		return
	}

	// Call the service to sync the user
	result, err := h.userService.SyncUser(r.Context(), services.SyncUserRequest{
		IDToken: token,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to sync user")
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Failed to verify token or sync user")
		return
	}

	// Build response
	response := SyncUserResponse{
		ID:          result.User.ID,
		FirebaseUID: result.User.FirebaseUID,
		Email:       result.User.Email,
		Name:        result.User.Name,
		CreatedAt:   result.User.CreatedAt,
	}

	status := http.StatusOK
	if result.IsNewUser {
		status = http.StatusCreated
	}

	respondJSON(w, status, response)
}
