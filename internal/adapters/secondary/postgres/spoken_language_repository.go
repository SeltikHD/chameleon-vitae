package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
)

// SpokenLanguageRepository implements ports.SpokenLanguageRepository using PostgreSQL.
type SpokenLanguageRepository struct {
	pool *pgxpool.Pool
}

// Create creates a new spoken language.
func (r *SpokenLanguageRepository) Create(ctx context.Context, language *domain.SpokenLanguage) error {
	if language.ID == "" {
		language.ID = uuid.New().String()
	}

	language.CreatedAt = time.Now().UTC()

	query := `
		INSERT INTO spoken_languages (
			id, user_id, language, proficiency, display_order, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`

	_, err := r.pool.Exec(ctx, query,
		language.ID,
		language.UserID,
		language.Language,
		string(language.Proficiency),
		language.DisplayOrder,
		language.CreatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("create spoken language", err)
	}

	return nil
}

// GetByID retrieves a spoken language by ID.
func (r *SpokenLanguageRepository) GetByID(ctx context.Context, id string) (*domain.SpokenLanguage, error) {
	query := `
		SELECT id, user_id, language, proficiency, display_order, created_at
		FROM spoken_languages
		WHERE id = $1
	`

	return r.scanLanguage(r.pool.QueryRow(ctx, query, id))
}

// ListByUserID lists all spoken languages for a user.
func (r *SpokenLanguageRepository) ListByUserID(ctx context.Context, userID string) ([]domain.SpokenLanguage, error) {
	query := `
		SELECT id, user_id, language, proficiency, display_order, created_at
		FROM spoken_languages
		WHERE user_id = $1
		ORDER BY display_order ASC, language ASC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, domain.NewDatabaseError("list spoken languages", err)
	}
	defer rows.Close()

	return r.scanLanguages(rows)
}

// Update updates an existing spoken language.
func (r *SpokenLanguageRepository) Update(ctx context.Context, language *domain.SpokenLanguage) error {
	query := `
		UPDATE spoken_languages SET
			language = $2,
			proficiency = $3,
			display_order = $4
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		language.ID,
		language.Language,
		string(language.Proficiency),
		language.DisplayOrder,
	)
	if err != nil {
		return domain.NewDatabaseError("update spoken language", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSpokenLanguageNotFound
	}

	return nil
}

// Delete removes a spoken language.
func (r *SpokenLanguageRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM spoken_languages WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return domain.NewDatabaseError("delete spoken language", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSpokenLanguageNotFound
	}

	return nil
}

// scanLanguage scans a single spoken language row.
func (r *SpokenLanguageRepository) scanLanguage(row pgx.Row) (*domain.SpokenLanguage, error) {
	lang := &domain.SpokenLanguage{}
	var proficiency string

	err := row.Scan(
		&lang.ID,
		&lang.UserID,
		&lang.Language,
		&proficiency,
		&lang.DisplayOrder,
		&lang.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSpokenLanguageNotFound
		}
		return nil, domain.NewDatabaseError("scan spoken language", err)
	}

	lang.Proficiency = domain.LanguageProficiency(proficiency)

	return lang, nil
}

// scanLanguages scans multiple spoken language rows.
func (r *SpokenLanguageRepository) scanLanguages(rows pgx.Rows) ([]domain.SpokenLanguage, error) {
	languages := make([]domain.SpokenLanguage, 0)

	for rows.Next() {
		lang := domain.SpokenLanguage{}
		var proficiency string

		err := rows.Scan(
			&lang.ID,
			&lang.UserID,
			&lang.Language,
			&proficiency,
			&lang.DisplayOrder,
			&lang.CreatedAt,
		)
		if err != nil {
			return nil, domain.NewDatabaseError("scan spoken language row", err)
		}

		lang.Proficiency = domain.LanguageProficiency(proficiency)
		languages = append(languages, lang)
	}

	if err := rows.Err(); err != nil {
		return nil, domain.NewDatabaseError("iterate spoken languages", err)
	}

	return languages, nil
}
