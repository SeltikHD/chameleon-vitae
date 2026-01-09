package http

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

// ToolsHandler handles utility tool HTTP requests.
type ToolsHandler struct {
	resumeService *services.ResumeService
}

// NewToolsHandler creates a new ToolsHandler.
func NewToolsHandler(resumeService *services.ResumeService) *ToolsHandler {
	return &ToolsHandler{
		resumeService: resumeService,
	}
}

// ParseJobURL parses a job posting URL and extracts structured data.
//
//	@Summary		Parse job URL
//	@Description	Fetches a job posting URL and extracts structured data using Jina Reader
//	@Tags			tools
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		ParseJobURLRequest	true	"Job URL to parse"
//	@Success		200		{object}	ParseJobURLResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid request body"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		422		{object}	ErrorResponse	"Failed to parse job posting"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/tools/parse-job [post]
func (h *ToolsHandler) ParseJobURL(w http.ResponseWriter, r *http.Request) {
	_, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var req ParseJobURLRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.URL == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "URL is required")
		return
	}

	parseReq := services.ParseJobURLRequest{
		URL: req.URL,
	}

	result, err := h.resumeService.ParseJobURL(r.Context(), parseReq)
	if err != nil {
		log.Error().Err(err).Str("url", req.URL).Msg("Failed to parse job URL")
		respondError(w, http.StatusUnprocessableEntity, "PARSE_FAILED", "Failed to parse job posting")
		return
	}

	response := ParseJobURLResponse{
		URL:      result.URL,
		Title:    result.Title,
		Markdown: result.Content,
		Metadata: &ParseJobMetadata{
			Source:    extractDomain(req.URL),
			FetchedAt: time.Now(),
		},
	}

	respondJSON(w, http.StatusOK, response)
}

// extractDomain extracts the domain from a URL.
func extractDomain(rawURL string) string {
	// Simple extraction - get the host from the URL.
	// A more robust implementation could use url.Parse.
	start := 0
	if len(rawURL) > 8 && rawURL[:8] == "https://" {
		start = 8
	} else if len(rawURL) > 7 && rawURL[:7] == "http://" {
		start = 7
	}
	end := start
	for end < len(rawURL) && rawURL[end] != '/' {
		end++
	}
	return rawURL[start:end]
}
