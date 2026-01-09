package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

// ExperienceRepository implements ports.ExperienceRepository using PostgreSQL.
type ExperienceRepository struct {
	pool *pgxpool.Pool
}

// Create creates a new experience.
func (r *ExperienceRepository) Create(ctx context.Context, experience *domain.Experience) error {
	if experience.ID == "" {
		experience.ID = uuid.New().String()
	}

	now := time.Now().UTC()
	experience.CreatedAt = now
	experience.UpdatedAt = now

	metadataJSON, err := json.Marshal(experience.Metadata)
	if err != nil {
		return domain.NewDatabaseError("marshal experience metadata", err)
	}

	query := `
		INSERT INTO experiences (
			id, user_id, type, title, organization, location,
			start_date, end_date, is_current, description, url,
			metadata, display_order, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`

	var endDate *time.Time
	if experience.EndDate != nil {
		t := experience.EndDate.Time
		endDate = &t
	}

	_, err = r.pool.Exec(ctx, query,
		experience.ID,
		experience.UserID,
		string(experience.Type),
		experience.Title,
		experience.Organization,
		experience.Location,
		experience.StartDate.Time,
		endDate,
		experience.IsCurrent,
		experience.Description,
		experience.URL,
		metadataJSON,
		experience.DisplayOrder,
		experience.CreatedAt,
		experience.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("create experience", err)
	}

	return nil
}

// GetByID retrieves an experience by ID.
func (r *ExperienceRepository) GetByID(ctx context.Context, id string) (*domain.Experience, error) {
	query := `
		SELECT id, user_id, type, title, organization, location,
			   start_date, end_date, is_current, description, url,
			   metadata, display_order, created_at, updated_at
		FROM experiences
		WHERE id = $1
	`

	return r.scanExperience(ctx, r.pool.QueryRow(ctx, query, id))
}

// GetByIDWithBullets retrieves an experience with all its bullets.
func (r *ExperienceRepository) GetByIDWithBullets(ctx context.Context, id string) (*domain.Experience, error) {
	exp, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	bulletQuery := `
		SELECT id, experience_id, content, impact_score, keywords,
			   metadata, display_order, created_at, updated_at
		FROM bullets
		WHERE experience_id = $1
		ORDER BY display_order ASC, created_at ASC
	`

	rows, err := r.pool.Query(ctx, bulletQuery, id)
	if err != nil {
		return nil, domain.NewDatabaseError("get bullets for experience", err)
	}
	defer rows.Close()

	bullets := make([]domain.Bullet, 0)
	for rows.Next() {
		bullet, err := scanBulletRow(rows)
		if err != nil {
			return nil, err
		}
		bullets = append(bullets, *bullet)
	}

	if err := rows.Err(); err != nil {
		return nil, domain.NewDatabaseError("iterate bullets", err)
	}

	exp.Bullets = bullets
	return exp, nil
}

// ListByUserID lists all experiences for a user.
func (r *ExperienceRepository) ListByUserID(ctx context.Context, userID string, opts ports.ListOptions) ([]domain.Experience, int, error) {
	countQuery := `SELECT COUNT(*) FROM experiences WHERE user_id = $1`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, domain.NewDatabaseError("count experiences", err)
	}

	query := `
		SELECT id, user_id, type, title, organization, location,
			   start_date, end_date, is_current, description, url,
			   metadata, display_order, created_at, updated_at
		FROM experiences
		WHERE user_id = $1
		ORDER BY display_order ASC, start_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, opts.Limit, opts.Offset)
	if err != nil {
		return nil, 0, domain.NewDatabaseError("list experiences", err)
	}
	defer rows.Close()

	experiences, err := r.scanExperiences(ctx, rows)
	if err != nil {
		return nil, 0, err
	}

	return experiences, total, nil
}

// ListByUserIDAndType lists experiences filtered by type.
func (r *ExperienceRepository) ListByUserIDAndType(ctx context.Context, userID string, expType domain.ExperienceType, opts ports.ListOptions) ([]domain.Experience, int, error) {
	countQuery := `SELECT COUNT(*) FROM experiences WHERE user_id = $1 AND type = $2`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, userID, string(expType)).Scan(&total); err != nil {
		return nil, 0, domain.NewDatabaseError("count experiences by type", err)
	}

	query := `
		SELECT id, user_id, type, title, organization, location,
			   start_date, end_date, is_current, description, url,
			   metadata, display_order, created_at, updated_at
		FROM experiences
		WHERE user_id = $1 AND type = $2
		ORDER BY display_order ASC, start_date DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.pool.Query(ctx, query, userID, string(expType), opts.Limit, opts.Offset)
	if err != nil {
		return nil, 0, domain.NewDatabaseError("list experiences by type", err)
	}
	defer rows.Close()

	experiences, err := r.scanExperiences(ctx, rows)
	if err != nil {
		return nil, 0, err
	}

	return experiences, total, nil
}

// Update updates an existing experience.
func (r *ExperienceRepository) Update(ctx context.Context, experience *domain.Experience) error {
	experience.UpdatedAt = time.Now().UTC()

	metadataJSON, err := json.Marshal(experience.Metadata)
	if err != nil {
		return domain.NewDatabaseError("marshal experience metadata", err)
	}

	query := `
		UPDATE experiences SET
			type = $2,
			title = $3,
			organization = $4,
			location = $5,
			start_date = $6,
			end_date = $7,
			is_current = $8,
			description = $9,
			url = $10,
			metadata = $11,
			display_order = $12,
			updated_at = $13
		WHERE id = $1
	`

	var endDate *time.Time
	if experience.EndDate != nil {
		t := experience.EndDate.Time
		endDate = &t
	}

	result, err := r.pool.Exec(ctx, query,
		experience.ID,
		string(experience.Type),
		experience.Title,
		experience.Organization,
		experience.Location,
		experience.StartDate.Time,
		endDate,
		experience.IsCurrent,
		experience.Description,
		experience.URL,
		metadataJSON,
		experience.DisplayOrder,
		experience.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("update experience", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrExperienceNotFound
	}

	return nil
}

// Delete removes an experience and all its bullets.
func (r *ExperienceRepository) Delete(ctx context.Context, id string) error {
	// Start transaction to delete bullets and experience atomically.
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.NewDatabaseError("begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Delete bullets first (foreign key constraint).
	_, err = tx.Exec(ctx, `DELETE FROM bullets WHERE experience_id = $1`, id)
	if err != nil {
		return domain.NewDatabaseError("delete bullets", err)
	}

	// Delete experience.
	result, err := tx.Exec(ctx, `DELETE FROM experiences WHERE id = $1`, id)
	if err != nil {
		return domain.NewDatabaseError("delete experience", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrExperienceNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.NewDatabaseError("commit transaction", err)
	}

	return nil
}

// UpdateDisplayOrder updates the display order of experiences.
func (r *ExperienceRepository) UpdateDisplayOrder(ctx context.Context, orders []ports.DisplayOrderUpdate) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.NewDatabaseError("begin transaction", err)
	}
	defer tx.Rollback(ctx)

	query := `UPDATE experiences SET display_order = $2, updated_at = $3 WHERE id = $1`
	now := time.Now().UTC()

	for _, order := range orders {
		_, err := tx.Exec(ctx, query, order.ID, order.DisplayOrder, now)
		if err != nil {
			return domain.NewDatabaseError("update display order", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.NewDatabaseError("commit transaction", err)
	}

	return nil
}

// scanExperience scans a single experience row.
func (r *ExperienceRepository) scanExperience(ctx context.Context, row pgx.Row) (*domain.Experience, error) {
	exp := &domain.Experience{}
	var expType string
	var startDate time.Time
	var endDate *time.Time
	var metadataJSON []byte

	err := row.Scan(
		&exp.ID,
		&exp.UserID,
		&expType,
		&exp.Title,
		&exp.Organization,
		&exp.Location,
		&startDate,
		&endDate,
		&exp.IsCurrent,
		&exp.Description,
		&exp.URL,
		&metadataJSON,
		&exp.DisplayOrder,
		&exp.CreatedAt,
		&exp.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrExperienceNotFound
		}
		return nil, domain.NewDatabaseError("scan experience", err)
	}

	exp.Type = domain.ExperienceType(expType)
	exp.StartDate = domain.Date{Time: startDate}
	if endDate != nil {
		exp.EndDate = &domain.Date{Time: *endDate}
	}

	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &exp.Metadata); err != nil {
			return nil, domain.NewDatabaseError("unmarshal experience metadata", err)
		}
	}
	if exp.Metadata == nil {
		exp.Metadata = make(map[string]any)
	}
	exp.Bullets = make([]domain.Bullet, 0)

	return exp, nil
}

// scanExperiences scans multiple experience rows.
func (r *ExperienceRepository) scanExperiences(ctx context.Context, rows pgx.Rows) ([]domain.Experience, error) {
	experiences := make([]domain.Experience, 0)

	for rows.Next() {
		exp := domain.Experience{}
		var expType string
		var startDate time.Time
		var endDate *time.Time
		var metadataJSON []byte

		err := rows.Scan(
			&exp.ID,
			&exp.UserID,
			&expType,
			&exp.Title,
			&exp.Organization,
			&exp.Location,
			&startDate,
			&endDate,
			&exp.IsCurrent,
			&exp.Description,
			&exp.URL,
			&metadataJSON,
			&exp.DisplayOrder,
			&exp.CreatedAt,
			&exp.UpdatedAt,
		)
		if err != nil {
			return nil, domain.NewDatabaseError("scan experience row", err)
		}

		exp.Type = domain.ExperienceType(expType)
		exp.StartDate = domain.Date{Time: startDate}
		if endDate != nil {
			exp.EndDate = &domain.Date{Time: *endDate}
		}

		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &exp.Metadata); err != nil {
				return nil, domain.NewDatabaseError("unmarshal experience metadata", err)
			}
		}
		if exp.Metadata == nil {
			exp.Metadata = make(map[string]any)
		}
		exp.Bullets = make([]domain.Bullet, 0)

		experiences = append(experiences, exp)
	}

	if err := rows.Err(); err != nil {
		return nil, domain.NewDatabaseError("iterate experiences", err)
	}

	return experiences, nil
}
