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

// ProjectBulletRepository implements ports.ProjectBulletRepository using PostgreSQL.
type ProjectBulletRepository struct {
	pool *pgxpool.Pool
}

// NewProjectBulletRepository creates a new ProjectBulletRepository.
func NewProjectBulletRepository(pool *pgxpool.Pool) *ProjectBulletRepository {
	return &ProjectBulletRepository{pool: pool}
}

// Create creates a new project bullet.
func (r *ProjectBulletRepository) Create(ctx context.Context, bullet *domain.ProjectBullet) error {
	if bullet.ID == "" {
		bullet.ID = uuid.New().String()
	}

	bullet.CreatedAt = time.Now().UTC()
	bullet.UpdatedAt = bullet.CreatedAt

	query := `
		INSERT INTO project_bullets (
			id, project_id, content, display_order, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`

	_, err := r.pool.Exec(ctx, query,
		bullet.ID,
		bullet.ProjectID,
		bullet.Content,
		bullet.DisplayOrder,
		bullet.CreatedAt,
		bullet.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("create project bullet", err)
	}

	return nil
}

// GetByID retrieves a project bullet by ID.
func (r *ProjectBulletRepository) GetByID(ctx context.Context, id string) (*domain.ProjectBullet, error) {
	query := `
		SELECT id, project_id, content, display_order, created_at, updated_at
		FROM project_bullets
		WHERE id = $1
	`

	var bullet domain.ProjectBullet
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&bullet.ID,
		&bullet.ProjectID,
		&bullet.Content,
		&bullet.DisplayOrder,
		&bullet.CreatedAt,
		&bullet.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProjectBulletNotFound
		}
		return nil, domain.NewDatabaseError("get project bullet", err)
	}

	return &bullet, nil
}

// ListByProjectID lists all bullets for a project, ordered by display_order.
func (r *ProjectBulletRepository) ListByProjectID(ctx context.Context, projectID string) ([]domain.ProjectBullet, error) {
	query := `
		SELECT id, project_id, content, display_order, created_at, updated_at
		FROM project_bullets
		WHERE project_id = $1
		ORDER BY display_order ASC
	`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, domain.NewDatabaseError("list project bullets", err)
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

// Update updates an existing project bullet.
func (r *ProjectBulletRepository) Update(ctx context.Context, bullet *domain.ProjectBullet) error {
	bullet.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE project_bullets SET
			content = $2,
			display_order = $3,
			updated_at = $4
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		bullet.ID,
		bullet.Content,
		bullet.DisplayOrder,
		bullet.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("update project bullet", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProjectBulletNotFound
	}

	return nil
}

// Delete removes a project bullet.
func (r *ProjectBulletRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM project_bullets WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return domain.NewDatabaseError("delete project bullet", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProjectBulletNotFound
	}

	return nil
}

// UpdateDisplayOrder updates the display order of project bullets.
func (r *ProjectBulletRepository) UpdateDisplayOrder(ctx context.Context, orders []ports.DisplayOrderUpdate) error {
	if len(orders) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.NewDatabaseError("begin transaction", err)
	}
	defer tx.Rollback(ctx)

	query := `UPDATE project_bullets SET display_order = $2, updated_at = $3 WHERE id = $1`

	for _, order := range orders {
		_, err := tx.Exec(ctx, query, order.ID, order.DisplayOrder, time.Now().UTC())
		if err != nil {
			return domain.NewDatabaseError("update project bullet order", err)
		}
	}

	return tx.Commit(ctx)
}
