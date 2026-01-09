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

// ResumeRepository implements ports.ResumeRepository using PostgreSQL.
type ResumeRepository struct {
	pool *pgxpool.Pool
}

// Create creates a new resume.
func (r *ResumeRepository) Create(ctx context.Context, resume *domain.Resume) error {
	if resume.ID == "" {
		resume.ID = uuid.New().String()
	}

	now := time.Now().UTC()
	resume.CreatedAt = now
	resume.UpdatedAt = now

	var contentJSON []byte
	var err error
	if resume.GeneratedContent != nil {
		contentJSON, err = json.Marshal(resume.GeneratedContent)
		if err != nil {
			return domain.NewDatabaseError("marshal resume content", err)
		}
	}

	query := `
		INSERT INTO resumes (
			id, user_id, job_description, job_title, company_name, job_url,
			target_language, selected_bullets, generated_content, pdf_url,
			score, notes, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`

	_, err = r.pool.Exec(ctx, query,
		resume.ID,
		resume.UserID,
		resume.JobDescription,
		resume.JobTitle,
		resume.CompanyName,
		resume.JobURL,
		resume.TargetLanguage,
		resume.SelectedBullets,
		contentJSON,
		resume.PDFURL,
		resume.Score.Int(),
		resume.Notes,
		string(resume.Status),
		resume.CreatedAt,
		resume.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("create resume", err)
	}

	return nil
}

// GetByID retrieves a resume by ID.
func (r *ResumeRepository) GetByID(ctx context.Context, id string) (*domain.Resume, error) {
	query := `
		SELECT id, user_id, job_description, job_title, company_name, job_url,
			   target_language, selected_bullets, generated_content, pdf_url,
			   score, notes, status, created_at, updated_at
		FROM resumes
		WHERE id = $1
	`

	return r.scanResume(r.pool.QueryRow(ctx, query, id))
}

// ListByUserID lists all resumes for a user.
func (r *ResumeRepository) ListByUserID(ctx context.Context, userID string, opts ports.ListOptions) ([]domain.Resume, int, error) {
	countQuery := `SELECT COUNT(*) FROM resumes WHERE user_id = $1`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, domain.NewDatabaseError("count resumes", err)
	}

	query := `
		SELECT id, user_id, job_description, job_title, company_name, job_url,
			   target_language, selected_bullets, generated_content, pdf_url,
			   score, notes, status, created_at, updated_at
		FROM resumes
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, opts.Limit, opts.Offset)
	if err != nil {
		return nil, 0, domain.NewDatabaseError("list resumes", err)
	}
	defer rows.Close()

	resumes, err := r.scanResumes(rows)
	if err != nil {
		return nil, 0, err
	}

	return resumes, total, nil
}

// ListByUserIDAndStatus lists resumes filtered by status.
func (r *ResumeRepository) ListByUserIDAndStatus(ctx context.Context, userID string, status domain.ResumeStatus, opts ports.ListOptions) ([]domain.Resume, int, error) {
	countQuery := `SELECT COUNT(*) FROM resumes WHERE user_id = $1 AND status = $2`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, userID, string(status)).Scan(&total); err != nil {
		return nil, 0, domain.NewDatabaseError("count resumes by status", err)
	}

	query := `
		SELECT id, user_id, job_description, job_title, company_name, job_url,
			   target_language, selected_bullets, generated_content, pdf_url,
			   score, notes, status, created_at, updated_at
		FROM resumes
		WHERE user_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.pool.Query(ctx, query, userID, string(status), opts.Limit, opts.Offset)
	if err != nil {
		return nil, 0, domain.NewDatabaseError("list resumes by status", err)
	}
	defer rows.Close()

	resumes, err := r.scanResumes(rows)
	if err != nil {
		return nil, 0, err
	}

	return resumes, total, nil
}

// Update updates an existing resume.
func (r *ResumeRepository) Update(ctx context.Context, resume *domain.Resume) error {
	resume.UpdatedAt = time.Now().UTC()

	var contentJSON []byte
	var err error
	if resume.GeneratedContent != nil {
		contentJSON, err = json.Marshal(resume.GeneratedContent)
		if err != nil {
			return domain.NewDatabaseError("marshal resume content", err)
		}
	}

	query := `
		UPDATE resumes SET
			job_description = $2,
			job_title = $3,
			company_name = $4,
			job_url = $5,
			target_language = $6,
			selected_bullets = $7,
			generated_content = $8,
			pdf_url = $9,
			score = $10,
			notes = $11,
			status = $12,
			updated_at = $13
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		resume.ID,
		resume.JobDescription,
		resume.JobTitle,
		resume.CompanyName,
		resume.JobURL,
		resume.TargetLanguage,
		resume.SelectedBullets,
		contentJSON,
		resume.PDFURL,
		resume.Score.Int(),
		resume.Notes,
		string(resume.Status),
		resume.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("update resume", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrResumeNotFound
	}

	return nil
}

// Delete removes a resume.
func (r *ResumeRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM resumes WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return domain.NewDatabaseError("delete resume", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrResumeNotFound
	}

	return nil
}

// scanResume scans a single resume row.
func (r *ResumeRepository) scanResume(row pgx.Row) (*domain.Resume, error) {
	resume := &domain.Resume{}
	var score int
	var status string
	var contentJSON []byte

	err := row.Scan(
		&resume.ID,
		&resume.UserID,
		&resume.JobDescription,
		&resume.JobTitle,
		&resume.CompanyName,
		&resume.JobURL,
		&resume.TargetLanguage,
		&resume.SelectedBullets,
		&contentJSON,
		&resume.PDFURL,
		&score,
		&resume.Notes,
		&status,
		&resume.CreatedAt,
		&resume.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrResumeNotFound
		}
		return nil, domain.NewDatabaseError("scan resume", err)
	}

	resume.Score = domain.MatchScore(score)
	resume.Status = domain.ResumeStatus(status)

	if len(contentJSON) > 0 {
		resume.GeneratedContent = &domain.ResumeContent{}
		if err := json.Unmarshal(contentJSON, resume.GeneratedContent); err != nil {
			return nil, domain.NewDatabaseError("unmarshal resume content", err)
		}
	}

	if resume.SelectedBullets == nil {
		resume.SelectedBullets = make([]string, 0)
	}

	return resume, nil
}

// scanResumes scans multiple resume rows.
func (r *ResumeRepository) scanResumes(rows pgx.Rows) ([]domain.Resume, error) {
	resumes := make([]domain.Resume, 0)

	for rows.Next() {
		resume := domain.Resume{}
		var score int
		var status string
		var contentJSON []byte

		err := rows.Scan(
			&resume.ID,
			&resume.UserID,
			&resume.JobDescription,
			&resume.JobTitle,
			&resume.CompanyName,
			&resume.JobURL,
			&resume.TargetLanguage,
			&resume.SelectedBullets,
			&contentJSON,
			&resume.PDFURL,
			&score,
			&resume.Notes,
			&status,
			&resume.CreatedAt,
			&resume.UpdatedAt,
		)
		if err != nil {
			return nil, domain.NewDatabaseError("scan resume row", err)
		}

		resume.Score = domain.MatchScore(score)
		resume.Status = domain.ResumeStatus(status)

		if len(contentJSON) > 0 {
			resume.GeneratedContent = &domain.ResumeContent{}
			if err := json.Unmarshal(contentJSON, resume.GeneratedContent); err != nil {
				return nil, domain.NewDatabaseError("unmarshal resume content", err)
			}
		}

		if resume.SelectedBullets == nil {
			resume.SelectedBullets = make([]string, 0)
		}

		resumes = append(resumes, resume)
	}

	if err := rows.Err(); err != nil {
		return nil, domain.NewDatabaseError("iterate resumes", err)
	}

	return resumes, nil
}
