package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

// ProjectRepository implements ports.ProjectRepository using PostgreSQL.
type ProjectRepository struct {
	pool *pgxpool.Pool
}

// NewProjectRepository creates a new ProjectRepository.
func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{pool: pool}
}

// Create creates a new project.
func (r *ProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	if project.ID == "" {
		project.ID = uuid.New().String()
	}

	project.CreatedAt = time.Now().UTC()
	project.UpdatedAt = project.CreatedAt

	query := `
		INSERT INTO projects (
			id, user_id, name, description, tech_stack,
			url, repository_url, start_date, end_date,
			display_order, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`

	var startDate, endDate interface{}
	if project.StartDate != nil {
		startDate = project.StartDate.Time
	}
	if project.EndDate != nil {
		endDate = project.EndDate.Time
	}

	_, err := r.pool.Exec(ctx, query,
		project.ID,
		project.UserID,
		project.Name,
		project.Description,
		project.TechStack,
		project.URL,
		project.RepositoryURL,
		startDate,
		endDate,
		project.DisplayOrder,
		project.CreatedAt,
		project.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("create project", err)
	}

	return nil
}

// GetByID retrieves a project by ID (without bullets).
func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	query := `
		SELECT id, user_id, name, description, tech_stack,
			   url, repository_url, start_date, end_date,
			   display_order, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	return r.scanProject(r.pool.QueryRow(ctx, query, id))
}

// GetByIDWithBullets retrieves a project with all its bullets.
func (r *ProjectRepository) GetByIDWithBullets(ctx context.Context, id string) (*domain.Project, error) {
	project, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	bullets, err := r.getBulletsByProjectID(ctx, id)
	if err != nil {
		return nil, err
	}

	project.Bullets = bullets
	return project, nil
}

// ListByUserID lists all projects for a user, ordered by display_order.
func (r *ProjectRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Project, error) {
	query := `
		SELECT id, user_id, name, description, tech_stack,
			   url, repository_url, start_date, end_date,
			   display_order, created_at, updated_at
		FROM projects
		WHERE user_id = $1
		ORDER BY display_order ASC, end_date DESC NULLS FIRST, start_date DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, domain.NewDatabaseError("list projects", err)
	}
	defer rows.Close()

	return r.scanProjectList(rows)
}

// ListByUserIDWithBullets lists all projects with bullets for a user.
func (r *ProjectRepository) ListByUserIDWithBullets(ctx context.Context, userID string) ([]domain.Project, error) {
	projects, err := r.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Load bullets for each project.
	for i := range projects {
		bullets, err := r.getBulletsByProjectID(ctx, projects[i].ID)
		if err != nil {
			return nil, err
		}
		projects[i].Bullets = bullets
	}

	return projects, nil
}

// Update updates an existing project.
func (r *ProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	project.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE projects SET
			name = $2,
			description = $3,
			tech_stack = $4,
			url = $5,
			repository_url = $6,
			start_date = $7,
			end_date = $8,
			display_order = $9,
			updated_at = $10
		WHERE id = $1
	`

	var startDate, endDate interface{}
	if project.StartDate != nil {
		startDate = project.StartDate.Time
	}
	if project.EndDate != nil {
		endDate = project.EndDate.Time
	}

	result, err := r.pool.Exec(ctx, query,
		project.ID,
		project.Name,
		project.Description,
		project.TechStack,
		project.URL,
		project.RepositoryURL,
		startDate,
		endDate,
		project.DisplayOrder,
		project.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("update project", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProjectNotFound
	}

	return nil
}

// Delete removes a project and all its bullets (CASCADE).
func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM projects WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return domain.NewDatabaseError("delete project", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProjectNotFound
	}

	return nil
}

// UpdateDisplayOrder updates the display order of projects.
func (r *ProjectRepository) UpdateDisplayOrder(ctx context.Context, orders []ports.DisplayOrderUpdate) error {
	if len(orders) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.NewDatabaseError("begin transaction", err)
	}
	defer tx.Rollback(ctx)

	query := `UPDATE projects SET display_order = $2, updated_at = $3 WHERE id = $1`

	for _, order := range orders {
		_, err := tx.Exec(ctx, query, order.ID, order.DisplayOrder, time.Now().UTC())
		if err != nil {
			return domain.NewDatabaseError("update project order", err)
		}
	}

	return tx.Commit(ctx)
}

// SearchByTechStack searches projects containing any of the given technologies.
func (r *ProjectRepository) SearchByTechStack(ctx context.Context, userID string, technologies []string) ([]domain.Project, error) {
	if len(technologies) == 0 {
		return r.ListByUserID(ctx, userID)
	}

	query := `
		SELECT id, user_id, name, description, tech_stack,
			   url, repository_url, start_date, end_date,
			   display_order, created_at, updated_at
		FROM projects
		WHERE user_id = $1 AND tech_stack && $2
		ORDER BY display_order ASC, end_date DESC NULLS FIRST
	`

	rows, err := r.pool.Query(ctx, query, userID, technologies)
	if err != nil {
		return nil, domain.NewDatabaseError("search projects by tech", err)
	}
	defer rows.Close()

	return r.scanProjectList(rows)
}

// getBulletsByProjectID retrieves bullets for a project.
func (r *ProjectRepository) getBulletsByProjectID(ctx context.Context, projectID string) ([]domain.ProjectBullet, error) {
	query := `
		SELECT id, project_id, content, display_order, created_at, updated_at
		FROM project_bullets
		WHERE project_id = $1
		ORDER BY display_order ASC
	`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, domain.NewDatabaseError("get project bullets", err)
	}
	defer rows.Close()

	var bullets []domain.ProjectBullet
	for rows.Next() {
		var bullet domain.ProjectBullet
		err := rows.Scan(
			&bullet.ID,
			&bullet.ProjectID,
			&bullet.Content,
			&bullet.DisplayOrder,
			&bullet.CreatedAt,
			&bullet.UpdatedAt,
		)
		if err != nil {
			return nil, domain.NewDatabaseError("scan project bullet", err)
		}
		bullets = append(bullets, bullet)
	}

	if err := rows.Err(); err != nil {
		return nil, domain.NewDatabaseError("iterate project bullets", err)
	}

	return bullets, nil
}

// scanProject scans a single project row.
func (r *ProjectRepository) scanProject(row pgx.Row) (*domain.Project, error) {
	var project domain.Project
	var startDate, endDate *time.Time

	err := row.Scan(
		&project.ID,
		&project.UserID,
		&project.Name,
		&project.Description,
		&project.TechStack,
		&project.URL,
		&project.RepositoryURL,
		&startDate,
		&endDate,
		&project.DisplayOrder,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProjectNotFound
		}
		return nil, domain.NewDatabaseError("scan project", err)
	}

	if startDate != nil {
		d := domain.Date{Time: *startDate}
		project.StartDate = &d
	}
	if endDate != nil {
		d := domain.Date{Time: *endDate}
		project.EndDate = &d
	}

	if project.TechStack == nil {
		project.TechStack = make([]string, 0)
	}
	if project.Bullets == nil {
		project.Bullets = make([]domain.ProjectBullet, 0)
	}

	return &project, nil
}

// scanProjectList scans multiple project rows.
func (r *ProjectRepository) scanProjectList(rows pgx.Rows) ([]domain.Project, error) {
	var projects []domain.Project

	for rows.Next() {
		var project domain.Project
		var startDate, endDate *time.Time

		err := rows.Scan(
			&project.ID,
			&project.UserID,
			&project.Name,
			&project.Description,
			&project.TechStack,
			&project.URL,
			&project.RepositoryURL,
			&startDate,
			&endDate,
			&project.DisplayOrder,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, domain.NewDatabaseError("scan project list", err)
		}

		if startDate != nil {
			d := domain.Date{Time: *startDate}
			project.StartDate = &d
		}
		if endDate != nil {
			d := domain.Date{Time: *endDate}
			project.EndDate = &d
		}

		if project.TechStack == nil {
			project.TechStack = make([]string, 0)
		}
		if project.Bullets == nil {
			project.Bullets = make([]domain.ProjectBullet, 0)
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, domain.NewDatabaseError("iterate project rows", err)
	}

	return projects, nil
}
