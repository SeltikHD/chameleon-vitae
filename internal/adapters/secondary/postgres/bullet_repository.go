package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
)

// BulletRepository implements ports.BulletRepository using PostgreSQL.
type BulletRepository struct {
	pool *pgxpool.Pool
}

// Create creates a new bullet.
func (r *BulletRepository) Create(ctx context.Context, bullet *domain.Bullet) error {
	if bullet.ID == "" {
		bullet.ID = uuid.New().String()
	}

	now := time.Now().UTC()
	bullet.CreatedAt = now
	bullet.UpdatedAt = now

	metadataJSON, err := json.Marshal(bullet.Metadata)
	if err != nil {
		return domain.NewDatabaseError("marshal bullet metadata", err)
	}

	query := `
		INSERT INTO bullets (
			id, experience_id, content, impact_score, keywords,
			metadata, display_order, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	_, err = r.pool.Exec(ctx, query,
		bullet.ID,
		bullet.ExperienceID,
		bullet.Content,
		bullet.ImpactScore.Int(),
		bullet.Keywords,
		metadataJSON,
		bullet.DisplayOrder,
		bullet.CreatedAt,
		bullet.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("create bullet", err)
	}

	return nil
}

// GetByID retrieves a bullet by ID.
func (r *BulletRepository) GetByID(ctx context.Context, id string) (*domain.Bullet, error) {
	query := `
		SELECT id, experience_id, content, impact_score, keywords,
			   metadata, display_order, created_at, updated_at
		FROM bullets
		WHERE id = $1
	`

	return r.scanBullet(r.pool.QueryRow(ctx, query, id))
}

// ListByExperienceID lists all bullets for an experience.
func (r *BulletRepository) ListByExperienceID(ctx context.Context, experienceID string) ([]domain.Bullet, error) {
	query := `
		SELECT id, experience_id, content, impact_score, keywords,
			   metadata, display_order, created_at, updated_at
		FROM bullets
		WHERE experience_id = $1
		ORDER BY display_order ASC, created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, experienceID)
	if err != nil {
		return nil, domain.NewDatabaseError("list bullets by experience", err)
	}
	defer rows.Close()

	return r.scanBullets(rows)
}

// ListByIDs retrieves multiple bullets by their IDs.
func (r *BulletRepository) ListByIDs(ctx context.Context, ids []string) ([]domain.Bullet, error) {
	if len(ids) == 0 {
		return []domain.Bullet{}, nil
	}

	// Build parameterized query for multiple IDs.
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, experience_id, content, impact_score, keywords,
			   metadata, display_order, created_at, updated_at
		FROM bullets
		WHERE id IN (%s)
		ORDER BY display_order ASC
	`, strings.Join(placeholders, ", "))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, domain.NewDatabaseError("list bullets by ids", err)
	}
	defer rows.Close()

	return r.scanBullets(rows)
}

// ListByUserID lists all bullets for a user (across all experiences).
func (r *BulletRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Bullet, error) {
	query := `
		SELECT b.id, b.experience_id, b.content, b.impact_score, b.keywords,
			   b.metadata, b.display_order, b.created_at, b.updated_at
		FROM bullets b
		INNER JOIN experiences e ON b.experience_id = e.id
		WHERE e.user_id = $1
		ORDER BY e.display_order ASC, b.display_order ASC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, domain.NewDatabaseError("list bullets by user", err)
	}
	defer rows.Close()

	return r.scanBullets(rows)
}

// Update updates an existing bullet.
func (r *BulletRepository) Update(ctx context.Context, bullet *domain.Bullet) error {
	bullet.UpdatedAt = time.Now().UTC()

	metadataJSON, err := json.Marshal(bullet.Metadata)
	if err != nil {
		return domain.NewDatabaseError("marshal bullet metadata", err)
	}

	query := `
		UPDATE bullets SET
			content = $2,
			impact_score = $3,
			keywords = $4,
			metadata = $5,
			display_order = $6,
			updated_at = $7
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		bullet.ID,
		bullet.Content,
		bullet.ImpactScore.Int(),
		bullet.Keywords,
		metadataJSON,
		bullet.DisplayOrder,
		bullet.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("update bullet", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrBulletNotFound
	}

	return nil
}

// Delete removes a bullet.
func (r *BulletRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM bullets WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return domain.NewDatabaseError("delete bullet", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrBulletNotFound
	}

	return nil
}

// SearchByKeywords searches bullets by keywords.
func (r *BulletRepository) SearchByKeywords(ctx context.Context, userID string, keywords []string) ([]domain.Bullet, error) {
	if len(keywords) == 0 {
		return []domain.Bullet{}, nil
	}

	// Use PostgreSQL array overlap operator for keyword matching.
	query := `
		SELECT b.id, b.experience_id, b.content, b.impact_score, b.keywords,
			   b.metadata, b.display_order, b.created_at, b.updated_at
		FROM bullets b
		INNER JOIN experiences e ON b.experience_id = e.id
		WHERE e.user_id = $1
		  AND (b.keywords && $2 OR b.content ILIKE ANY($3))
		ORDER BY b.impact_score DESC
	`

	// Create patterns for ILIKE.
	patterns := make([]string, len(keywords))
	for i, kw := range keywords {
		patterns[i] = "%" + kw + "%"
	}

	rows, err := r.pool.Query(ctx, query, userID, keywords, patterns)
	if err != nil {
		return nil, domain.NewDatabaseError("search bullets by keywords", err)
	}
	defer rows.Close()

	return r.scanBullets(rows)
}

// GetHighImpactBullets retrieves bullets with impact score >= threshold.
func (r *BulletRepository) GetHighImpactBullets(ctx context.Context, userID string, minScore int, limit int) ([]domain.Bullet, error) {
	query := `
		SELECT b.id, b.experience_id, b.content, b.impact_score, b.keywords,
			   b.metadata, b.display_order, b.created_at, b.updated_at
		FROM bullets b
		INNER JOIN experiences e ON b.experience_id = e.id
		WHERE e.user_id = $1 AND b.impact_score >= $2
		ORDER BY b.impact_score DESC
		LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, userID, minScore, limit)
	if err != nil {
		return nil, domain.NewDatabaseError("get high impact bullets", err)
	}
	defer rows.Close()

	return r.scanBullets(rows)
}

// scanBullet scans a single bullet row.
func (r *BulletRepository) scanBullet(row pgx.Row) (*domain.Bullet, error) {
	bullet := &domain.Bullet{}
	var impactScore int
	var metadataJSON []byte

	err := row.Scan(
		&bullet.ID,
		&bullet.ExperienceID,
		&bullet.Content,
		&impactScore,
		&bullet.Keywords,
		&metadataJSON,
		&bullet.DisplayOrder,
		&bullet.CreatedAt,
		&bullet.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBulletNotFound
		}
		return nil, domain.NewDatabaseError("scan bullet", err)
	}

	bullet.ImpactScore = domain.ImpactScore(impactScore)

	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &bullet.Metadata); err != nil {
			return nil, domain.NewDatabaseError("unmarshal bullet metadata", err)
		}
	}
	if bullet.Metadata == nil {
		bullet.Metadata = make(map[string]any)
	}
	if bullet.Keywords == nil {
		bullet.Keywords = make([]string, 0)
	}

	return bullet, nil
}

// scanBullets scans multiple bullet rows.
func (r *BulletRepository) scanBullets(rows pgx.Rows) ([]domain.Bullet, error) {
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

	return bullets, nil
}

// scanBulletRow is a helper for scanning a single row from a Rows result.
func scanBulletRow(rows pgx.Rows) (*domain.Bullet, error) {
	bullet := &domain.Bullet{}
	var impactScore int
	var metadataJSON []byte

	err := rows.Scan(
		&bullet.ID,
		&bullet.ExperienceID,
		&bullet.Content,
		&impactScore,
		&bullet.Keywords,
		&metadataJSON,
		&bullet.DisplayOrder,
		&bullet.CreatedAt,
		&bullet.UpdatedAt,
	)
	if err != nil {
		return nil, domain.NewDatabaseError("scan bullet row", err)
	}

	bullet.ImpactScore = domain.ImpactScore(impactScore)

	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &bullet.Metadata); err != nil {
			return nil, domain.NewDatabaseError("unmarshal bullet metadata", err)
		}
	}
	if bullet.Metadata == nil {
		bullet.Metadata = make(map[string]any)
	}
	if bullet.Keywords == nil {
		bullet.Keywords = make([]string, 0)
	}

	return bullet, nil
}
