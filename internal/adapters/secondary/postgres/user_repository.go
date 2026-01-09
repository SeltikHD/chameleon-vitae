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

// UserRepository implements ports.UserRepository using PostgreSQL.
type UserRepository struct {
	pool *pgxpool.Pool
}

// Create creates a new user in the database.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
		INSERT INTO users (
			id, firebase_uid, picture_url, email, name, headline, summary,
			location, phone, website, linkedin_url, github_url, portfolio_url,
			preferred_language, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`

	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.FirebaseUID,
		user.PictureURL,
		user.Email,
		user.Name,
		user.Headline,
		user.Summary,
		user.Location,
		user.Phone,
		user.Website,
		user.LinkedInURL,
		user.GitHubURL,
		user.PortfolioURL,
		user.PreferredLanguage,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("create user", err)
	}

	return nil
}

// GetByID retrieves a user by their internal ID.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, firebase_uid, picture_url, email, name, headline, summary,
			   location, phone, website, linkedin_url, github_url, portfolio_url,
			   preferred_language, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FirebaseUID,
		&user.PictureURL,
		&user.Email,
		&user.Name,
		&user.Headline,
		&user.Summary,
		&user.Location,
		&user.Phone,
		&user.Website,
		&user.LinkedInURL,
		&user.GitHubURL,
		&user.PortfolioURL,
		&user.PreferredLanguage,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.NewDatabaseError("get user by id", err)
	}

	return user, nil
}

// GetByFirebaseUID retrieves a user by their Firebase UID.
func (r *UserRepository) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	query := `
		SELECT id, firebase_uid, picture_url, email, name, headline, summary,
			   location, phone, website, linkedin_url, github_url, portfolio_url,
			   preferred_language, created_at, updated_at
		FROM users
		WHERE firebase_uid = $1
	`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, firebaseUID).Scan(
		&user.ID,
		&user.FirebaseUID,
		&user.PictureURL,
		&user.Email,
		&user.Name,
		&user.Headline,
		&user.Summary,
		&user.Location,
		&user.Phone,
		&user.Website,
		&user.LinkedInURL,
		&user.GitHubURL,
		&user.PortfolioURL,
		&user.PreferredLanguage,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.NewDatabaseError("get user by firebase uid", err)
	}

	return user, nil
}

// Update updates an existing user.
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE users SET
			picture_url = $2,
			email = $3,
			name = $4,
			headline = $5,
			summary = $6,
			location = $7,
			phone = $8,
			website = $9,
			linkedin_url = $10,
			github_url = $11,
			portfolio_url = $12,
			preferred_language = $13,
			updated_at = $14
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		user.ID,
		user.PictureURL,
		user.Email,
		user.Name,
		user.Headline,
		user.Summary,
		user.Location,
		user.Phone,
		user.Website,
		user.LinkedInURL,
		user.GitHubURL,
		user.PortfolioURL,
		user.PreferredLanguage,
		user.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("update user", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete removes a user from the database.
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return domain.NewDatabaseError("delete user", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Upsert creates or updates a user based on Firebase UID.
func (r *UserRepository) Upsert(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now().UTC()
	user.UpdatedAt = now

	query := `
		INSERT INTO users (
			id, firebase_uid, picture_url, email, name, headline, summary,
			location, phone, website, linkedin_url, github_url, portfolio_url,
			preferred_language, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
		ON CONFLICT (firebase_uid) DO UPDATE SET
			picture_url = EXCLUDED.picture_url,
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			updated_at = EXCLUDED.updated_at
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(ctx, query,
		user.ID,
		user.FirebaseUID,
		user.PictureURL,
		user.Email,
		user.Name,
		user.Headline,
		user.Summary,
		user.Location,
		user.Phone,
		user.Website,
		user.LinkedInURL,
		user.GitHubURL,
		user.PortfolioURL,
		user.PreferredLanguage,
		now, // created_at for new records
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return domain.NewDatabaseError("upsert user", err)
	}

	return nil
}
