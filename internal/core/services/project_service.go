// Package services contains the application services (use cases).
package services

import (
	"context"
	"fmt"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

// ProjectService handles project-related use cases.
type ProjectService struct {
	projectRepo       ports.ProjectRepository
	projectBulletRepo ports.ProjectBulletRepository
}

// NewProjectService creates a new ProjectService with required dependencies.
func NewProjectService(
	projectRepo ports.ProjectRepository,
	projectBulletRepo ports.ProjectBulletRepository,
) *ProjectService {
	return &ProjectService{
		projectRepo:       projectRepo,
		projectBulletRepo: projectBulletRepo,
	}
}

// CreateProjectRequest contains the parameters for creating a project.
type CreateProjectRequest struct {
	UserID        string
	Name          string
	Description   *string
	TechStack     []string
	URL           *string
	RepositoryURL *string
	StartDate     *domain.Date
	EndDate       *domain.Date
	DisplayOrder  int
	Bullets       []string // Initial bullet contents
}

// CreateProject creates a new project for a user.
func (s *ProjectService) CreateProject(ctx context.Context, req CreateProjectRequest) (*domain.Project, error) {
	// Create project.
	project, err := domain.NewProject(req.UserID, req.Name, req.TechStack)
	if err != nil {
		return nil, err
	}

	// Set optional fields.
	if req.Description != nil {
		project.SetDescription(*req.Description)
	}

	if req.URL != nil {
		project.SetURL(*req.URL)
	}

	if req.RepositoryURL != nil {
		project.SetRepositoryURL(*req.RepositoryURL)
	}

	if req.StartDate != nil || req.EndDate != nil {
		if err := project.SetDates(req.StartDate, req.EndDate); err != nil {
			return nil, err
		}
	}

	project.DisplayOrder = req.DisplayOrder

	// Validate.
	if err := project.Validate(); err != nil {
		return nil, err
	}

	// Save project.
	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Create initial bullets if provided.
	for i, content := range req.Bullets {
		bullet, err := domain.NewProjectBullet(project.ID, content)
		if err != nil {
			continue
		}
		bullet.DisplayOrder = i

		if err := s.projectBulletRepo.Create(ctx, bullet); err != nil {
			// Log but don't fail - project is already created
			continue
		}

		project.Bullets = append(project.Bullets, *bullet)
	}

	return project, nil
}

// GetProject retrieves a project by ID.
func (s *ProjectService) GetProject(ctx context.Context, projectID string) (*domain.Project, error) {
	project, err := s.projectRepo.GetByIDWithBullets(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	return project, nil
}

// ListProjectsRequest contains parameters for listing projects.
type ListProjectsRequest struct {
	UserID string
}

// ListProjects lists all projects for a user.
func (s *ProjectService) ListProjects(ctx context.Context, req ListProjectsRequest) ([]domain.Project, error) {
	projects, err := s.projectRepo.ListByUserIDWithBullets(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	return projects, nil
}

// UpdateProjectRequest contains parameters for updating a project.
type UpdateProjectRequest struct {
	ProjectID     string
	Name          *string
	Description   *string
	TechStack     []string
	URL           *string
	RepositoryURL *string
	StartDate     *domain.Date
	EndDate       *domain.Date
	DisplayOrder  *int
}

// UpdateProject updates an existing project.
func (s *ProjectService) UpdateProject(ctx context.Context, req UpdateProjectRequest) (*domain.Project, error) {
	// Get existing.
	project, err := s.projectRepo.GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}

	// Update fields.
	if req.Name != nil {
		project.Name = *req.Name
	}

	if req.Description != nil {
		project.SetDescription(*req.Description)
	}

	if req.TechStack != nil {
		project.SetTechStack(req.TechStack)
	}

	if req.URL != nil {
		project.SetURL(*req.URL)
	}

	if req.RepositoryURL != nil {
		project.SetRepositoryURL(*req.RepositoryURL)
	}

	if req.StartDate != nil || req.EndDate != nil {
		startDate := project.StartDate
		endDate := project.EndDate

		if req.StartDate != nil {
			startDate = req.StartDate
		}
		if req.EndDate != nil {
			endDate = req.EndDate
		}

		if err := project.SetDates(startDate, endDate); err != nil {
			return nil, err
		}
	}

	if req.DisplayOrder != nil {
		project.DisplayOrder = *req.DisplayOrder
	}

	// Validate.
	if err := project.Validate(); err != nil {
		return nil, err
	}

	// Save.
	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	// Reload with bullets.
	return s.projectRepo.GetByIDWithBullets(ctx, project.ID)
}

// DeleteProject removes a project and all its bullets.
func (s *ProjectService) DeleteProject(ctx context.Context, projectID string) error {
	if err := s.projectRepo.Delete(ctx, projectID); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	return nil
}

// UpdateProjectOrderRequest contains parameters for updating project display order.
type UpdateProjectOrderRequest struct {
	Orders []ports.DisplayOrderUpdate
}

// UpdateProjectOrder updates the display order of projects.
func (s *ProjectService) UpdateProjectOrder(ctx context.Context, req UpdateProjectOrderRequest) error {
	if err := s.projectRepo.UpdateDisplayOrder(ctx, req.Orders); err != nil {
		return fmt.Errorf("failed to update project order: %w", err)
	}
	return nil
}

// AddProjectBulletRequest contains parameters for adding a bullet to a project.
type AddProjectBulletRequest struct {
	ProjectID    string
	Content      string
	DisplayOrder int
}

// AddProjectBullet adds a bullet to a project.
func (s *ProjectService) AddProjectBullet(ctx context.Context, req AddProjectBulletRequest) (*domain.ProjectBullet, error) {
	// Verify project exists.
	_, err := s.projectRepo.GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}

	// Create bullet.
	bullet, err := domain.NewProjectBullet(req.ProjectID, req.Content)
	if err != nil {
		return nil, err
	}
	bullet.DisplayOrder = req.DisplayOrder

	// Save.
	if err := s.projectBulletRepo.Create(ctx, bullet); err != nil {
		return nil, fmt.Errorf("failed to create project bullet: %w", err)
	}

	return bullet, nil
}

// UpdateProjectBulletRequest contains parameters for updating a project bullet.
type UpdateProjectBulletRequest struct {
	BulletID     string
	Content      *string
	DisplayOrder *int
}

// UpdateProjectBullet updates an existing project bullet.
func (s *ProjectService) UpdateProjectBullet(ctx context.Context, req UpdateProjectBulletRequest) (*domain.ProjectBullet, error) {
	// Get existing.
	bullet, err := s.projectBulletRepo.GetByID(ctx, req.BulletID)
	if err != nil {
		return nil, err
	}

	// Update fields.
	if req.Content != nil {
		bullet.Content = *req.Content
	}

	if req.DisplayOrder != nil {
		bullet.DisplayOrder = *req.DisplayOrder
	}

	// Save.
	if err := s.projectBulletRepo.Update(ctx, bullet); err != nil {
		return nil, fmt.Errorf("failed to update project bullet: %w", err)
	}

	return bullet, nil
}

// DeleteProjectBullet removes a project bullet.
func (s *ProjectService) DeleteProjectBullet(ctx context.Context, bulletID string) error {
	if err := s.projectBulletRepo.Delete(ctx, bulletID); err != nil {
		return fmt.Errorf("failed to delete project bullet: %w", err)
	}
	return nil
}

// SearchProjectsByTechRequest contains parameters for searching projects by tech stack.
type SearchProjectsByTechRequest struct {
	UserID       string
	Technologies []string
}

// SearchProjectsByTech searches for projects containing specific technologies.
func (s *ProjectService) SearchProjectsByTech(ctx context.Context, req SearchProjectsByTechRequest) ([]domain.Project, error) {
	projects, err := s.projectRepo.SearchByTechStack(ctx, req.UserID, req.Technologies)
	if err != nil {
		return nil, fmt.Errorf("failed to search projects: %w", err)
	}
	return projects, nil
}
