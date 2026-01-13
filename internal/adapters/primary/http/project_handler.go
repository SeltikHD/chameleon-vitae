package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

// ProjectHandler handles project-related HTTP requests.
type ProjectHandler struct {
	projectService *services.ProjectService
}

// NewProjectHandler creates a new ProjectHandler.
func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// List returns all projects for the authenticated user.
//
//	@Summary		List projects
//	@Description	Returns all projects for the authenticated user, ordered by display_order
//	@Tags			projects
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	ListProjectsResponse
//	@Failure		401	{object}	ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/projects [get]
func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	req := services.ListProjectsRequest{
		UserID: authUser.ID,
	}

	projects, err := h.projectService.ListProjects(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to list projects")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve projects")
		return
	}

	data := make([]ProjectResponse, 0, len(projects))
	for _, proj := range projects {
		data = append(data, mapProjectToResponse(&proj))
	}

	respondJSON(w, http.StatusOK, ListProjectsResponse{
		Data:  data,
		Total: len(data),
	})
}

// Create creates a new project.
//
//	@Summary		Create project
//	@Description	Creates a new project for the authenticated user
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateProjectRequest	true	"Project data"
//	@Success		201		{object}	ProjectResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid request body"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		422		{object}	ErrorResponse	"Validation failed"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/projects [post]
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var req CreateProjectRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Validate required fields.
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name is required")
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

	techStack := req.TechStack
	if techStack == nil {
		techStack = make([]string, 0)
	}

	// Convert non-pointer strings to pointers for optional fields.
	var description, url, repositoryURL *string
	if req.Description != "" {
		description = &req.Description
	}
	if req.URL != "" {
		url = &req.URL
	}
	if req.RepositoryURL != "" {
		repositoryURL = &req.RepositoryURL
	}

	svcReq := services.CreateProjectRequest{
		UserID:        authUser.ID,
		Name:          req.Name,
		Description:   description,
		TechStack:     techStack,
		URL:           url,
		RepositoryURL: repositoryURL,
		StartDate:     startDate,
		EndDate:       endDate,
		DisplayOrder:  req.DisplayOrder,
		Bullets:       req.Bullets,
	}

	project, err := h.projectService.CreateProject(r.Context(), svcReq)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		log.Error().Err(err).Str("user_id", authUser.ID).Msg("Failed to create project")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create project")
		return
	}

	respondJSON(w, http.StatusCreated, mapProjectToResponse(project))
}

// Get retrieves a single project by ID.
//
//	@Summary		Get project
//	@Description	Retrieves a specific project by ID with its bullets
//	@Tags			projects
//	@Produce		json
//	@Security		BearerAuth
//	@Param			projectID	path		string	true	"Project ID"
//	@Success		200			{object}	ProjectResponse
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Project not found"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/projects/{projectID} [get]
func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	projectID := chi.URLParam(r, "projectID")
	if projectID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Project ID is required")
		return
	}

	project, err := h.projectService.GetProject(r.Context(), projectID)
	if err != nil {
		if errors.Is(err, domain.ErrProjectNotFound) {
			respondError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found")
			return
		}
		log.Error().Err(err).Str("project_id", projectID).Msg("Failed to get project")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve project")
		return
	}

	// Verify ownership.
	if project.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found")
		return
	}

	respondJSON(w, http.StatusOK, mapProjectToResponse(project))
}

// Update updates an existing project.
//
//	@Summary		Update project
//	@Description	Updates an existing project
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			projectID	path		string				true	"Project ID"
//	@Param			request		body		UpdateProjectRequest	true	"Project data"
//	@Success		200			{object}	ProjectResponse
//	@Failure		400			{object}	ErrorResponse	"Invalid request body"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Project not found"
//	@Failure		422			{object}	ErrorResponse	"Validation failed"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/projects/{projectID} [put]
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	projectID := chi.URLParam(r, "projectID")
	if projectID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Project ID is required")
		return
	}

	// Verify ownership.
	existing, err := h.projectService.GetProject(r.Context(), projectID)
	if err != nil {
		if errors.Is(err, domain.ErrProjectNotFound) {
			respondError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve project")
		return
	}
	if existing.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found")
		return
	}

	var req UpdateProjectRequest
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

	svcReq := services.UpdateProjectRequest{
		ProjectID:     projectID,
		Name:          req.Name,
		Description:   req.Description,
		TechStack:     req.TechStack,
		URL:           req.URL,
		RepositoryURL: req.RepositoryURL,
		StartDate:     startDate,
		EndDate:       endDate,
		DisplayOrder:  req.DisplayOrder,
	}

	project, err := h.projectService.UpdateProject(r.Context(), svcReq)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		log.Error().Err(err).Str("project_id", projectID).Msg("Failed to update project")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update project")
		return
	}

	respondJSON(w, http.StatusOK, mapProjectToResponse(project))
}

// Delete removes a project.
//
//	@Summary		Delete project
//	@Description	Deletes a project and all its bullets
//	@Tags			projects
//	@Produce		json
//	@Security		BearerAuth
//	@Param			projectID	path	string	true	"Project ID"
//	@Success		204			"No Content"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Project not found"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/projects/{projectID} [delete]
func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	projectID := chi.URLParam(r, "projectID")
	if projectID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Project ID is required")
		return
	}

	// Verify ownership.
	existing, err := h.projectService.GetProject(r.Context(), projectID)
	if err != nil {
		if errors.Is(err, domain.ErrProjectNotFound) {
			respondError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve project")
		return
	}
	if existing.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found")
		return
	}

	if err := h.projectService.DeleteProject(r.Context(), projectID); err != nil {
		log.Error().Err(err).Str("project_id", projectID).Msg("Failed to delete project")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete project")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddBullet adds a bullet to a project.
//
//	@Summary		Add project bullet
//	@Description	Adds a new bullet to a project
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			projectID	path		string					true	"Project ID"
//	@Param			request		body		CreateProjectBulletRequest	true	"Bullet data"
//	@Success		201			{object}	ProjectBulletResponse
//	@Failure		400			{object}	ErrorResponse	"Invalid request body"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Project not found"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/projects/{projectID}/bullets [post]
func (h *ProjectHandler) AddBullet(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	projectID := chi.URLParam(r, "projectID")
	if projectID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Project ID is required")
		return
	}

	// Verify ownership.
	project, err := h.projectService.GetProject(r.Context(), projectID)
	if err != nil {
		if errors.Is(err, domain.ErrProjectNotFound) {
			respondError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve project")
		return
	}
	if project.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found")
		return
	}

	var req CreateProjectBulletRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Content == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Content is required")
		return
	}

	svcReq := services.AddProjectBulletRequest{
		ProjectID:    projectID,
		Content:      req.Content,
		DisplayOrder: req.DisplayOrder,
	}

	bullet, err := h.projectService.AddProjectBullet(r.Context(), svcReq)
	if err != nil {
		log.Error().Err(err).Str("project_id", projectID).Msg("Failed to add project bullet")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to add bullet")
		return
	}

	respondJSON(w, http.StatusCreated, mapProjectBulletToResponse(bullet))
}

// DeleteBullet removes a bullet from a project.
//
//	@Summary		Delete project bullet
//	@Description	Deletes a bullet from a project
//	@Tags			projects
//	@Produce		json
//	@Security		BearerAuth
//	@Param			projectID	path	string	true	"Project ID"
//	@Param			bulletID	path	string	true	"Bullet ID"
//	@Success		204			"No Content"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Project or bullet not found"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/v1/projects/{projectID}/bullets/{bulletID} [delete]
func (h *ProjectHandler) DeleteBullet(w http.ResponseWriter, r *http.Request) {
	authUser, ok := GetAuthenticatedUser(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	projectID := chi.URLParam(r, "projectID")
	bulletID := chi.URLParam(r, "bulletID")
	if projectID == "" || bulletID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Project ID and Bullet ID are required")
		return
	}

	// Verify ownership.
	project, err := h.projectService.GetProject(r.Context(), projectID)
	if err != nil {
		if errors.Is(err, domain.ErrProjectNotFound) {
			respondError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve project")
		return
	}
	if project.UserID != authUser.ID {
		respondError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found")
		return
	}

	if err := h.projectService.DeleteProjectBullet(r.Context(), bulletID); err != nil {
		if errors.Is(err, domain.ErrProjectBulletNotFound) {
			respondError(w, http.StatusNotFound, "BULLET_NOT_FOUND", "Bullet not found")
			return
		}
		log.Error().Err(err).Str("bullet_id", bulletID).Msg("Failed to delete project bullet")
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete bullet")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// mapProjectToResponse maps a domain project to a response DTO.
func mapProjectToResponse(project *domain.Project) ProjectResponse {
	resp := ProjectResponse{
		ID:           project.ID,
		Name:         project.Name,
		TechStack:    project.TechStack,
		DisplayOrder: project.DisplayOrder,
		CreatedAt:    project.CreatedAt,
		UpdatedAt:    project.UpdatedAt,
	}

	// Convert *string to string for optional fields.
	if project.Description != nil {
		resp.Description = *project.Description
	}
	if project.URL != nil {
		resp.URL = *project.URL
	}
	if project.RepositoryURL != nil {
		resp.RepositoryURL = *project.RepositoryURL
	}

	if project.StartDate != nil {
		s := project.StartDate.String()
		resp.StartDate = &s
	}
	if project.EndDate != nil {
		s := project.EndDate.String()
		resp.EndDate = &s
	}

	if resp.TechStack == nil {
		resp.TechStack = make([]string, 0)
	}

	// Map bullets.
	resp.Bullets = make([]ProjectBulletResponse, 0, len(project.Bullets))
	for _, bullet := range project.Bullets {
		resp.Bullets = append(resp.Bullets, mapProjectBulletToResponse(&bullet))
	}

	return resp
}

// mapProjectBulletToResponse maps a domain project bullet to a response DTO.
func mapProjectBulletToResponse(bullet *domain.ProjectBullet) ProjectBulletResponse {
	return ProjectBulletResponse{
		ID:           bullet.ID,
		ProjectID:    bullet.ProjectID,
		Content:      bullet.Content,
		DisplayOrder: bullet.DisplayOrder,
		CreatedAt:    bullet.CreatedAt,
	}
}
