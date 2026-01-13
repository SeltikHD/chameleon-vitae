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

// EducationRepository implements ports.EducationRepository using PostgreSQL.
type EducationRepository struct {
	pool *pgxpool.Pool
}

// NewEducationRepository creates a new EducationRepository.
func NewEducationRepository(pool *pgxpool.Pool) *EducationRepository {
	return &EducationRepository{pool: pool}
}

// Create creates a new education entry.
func (r *EducationRepository) Create(ctx context.Context, education *domain.Education) error {
	if education.ID == "" {
		education.ID = uuid.New().String()
	}

	education.CreatedAt = time.Now().UTC()
	education.UpdatedAt = education.CreatedAt

	query := `
		INSERT INTO education (
			id, user_id, institution, degree, field_of_study,
			location, start_date, end_date, gpa, honors,
			display_order, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	var startDate, endDate interface{}
	if education.StartDate != nil {
		startDate = education.StartDate.Time
	}
	if education.EndDate != nil {
		endDate = education.EndDate.Time
	}

	_, err := r.pool.Exec(ctx, query,
		education.ID,
		education.UserID,
		education.Institution,
		education.Degree,
		education.FieldOfStudy,
		education.Location,
		startDate,
		endDate,
		education.GPA,
		education.Honors,
		education.DisplayOrder,
		education.CreatedAt,
		education.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("create education", err)
	}

	return nil
}

// GetByID retrieves an education entry by ID.
func (r *EducationRepository) GetByID(ctx context.Context, id string) (*domain.Education, error) {
	query := `
		SELECT id, user_id, institution, degree, field_of_study,
			   location, start_date, end_date, gpa, honors,
			   display_order, created_at, updated_at
		FROM education
		WHERE id = $1
	`

	return r.scanEducation(r.pool.QueryRow(ctx, query, id))
}

// ListByUserID lists all education entries for a user, ordered by display_order.
func (r *EducationRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Education, error) {
	query := `
		SELECT id, user_id, institution, degree, field_of_study,
			   location, start_date, end_date, gpa, honors,
			   display_order, created_at, updated_at
		FROM education
		WHERE user_id = $1
		ORDER BY display_order ASC, end_date DESC NULLS FIRST, start_date DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, domain.NewDatabaseError("list education", err)
	}
	defer rows.Close()

	return r.scanEducationList(rows)
}

// Update updates an existing education entry.
func (r *EducationRepository) Update(ctx context.Context, education *domain.Education) error {
	education.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE education SET
			institution = $2,
			degree = $3,
			field_of_study = $4,
			location = $5,
			start_date = $6,
			end_date = $7,
			gpa = $8,
			honors = $9,
			display_order = $10,
			updated_at = $11
		WHERE id = $1
	`

	var startDate, endDate interface{}
	if education.StartDate != nil {
		startDate = education.StartDate.Time
	}
	if education.EndDate != nil {
		endDate = education.EndDate.Time
	}

	result, err := r.pool.Exec(ctx, query,
		education.ID,
		education.Institution,
		education.Degree,
		education.FieldOfStudy,
		education.Location,
		startDate,
		endDate,
		education.GPA,
		education.Honors,
		education.DisplayOrder,
		education.UpdatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("update education", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrEducationNotFound
	}

	return nil
}

// Delete removes an education entry.
func (r *EducationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM education WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return domain.NewDatabaseError("delete education", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrEducationNotFound
	}

	return nil
}

// UpdateDisplayOrder updates the display order of education entries.
func (r *EducationRepository) UpdateDisplayOrder(ctx context.Context, orders []ports.DisplayOrderUpdate) error {
	if len(orders) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.NewDatabaseError("begin transaction", err)
	}
	defer tx.Rollback(ctx)

	query := `UPDATE education SET display_order = $2, updated_at = $3 WHERE id = $1`

	for _, order := range orders {
		_, err := tx.Exec(ctx, query, order.ID, order.DisplayOrder, time.Now().UTC())
		if err != nil {
			return domain.NewDatabaseError("update education order", err)
		}
	}

	return tx.Commit(ctx)
}

// scanEducation scans a single education row.
func (r *EducationRepository) scanEducation(row pgx.Row) (*domain.Education, error) {
	var education domain.Education
	var startDate, endDate *time.Time

	err := row.Scan(
		&education.ID,
		&education.UserID,
		&education.Institution,
		&education.Degree,
		&education.FieldOfStudy,
		&education.Location,
		&startDate,
		&endDate,
		&education.GPA,
		&education.Honors,
		&education.DisplayOrder,
		&education.CreatedAt,
		&education.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEducationNotFound
		}
		return nil, domain.NewDatabaseError("scan education", err)
	}

	if startDate != nil {
		d := domain.Date{Time: *startDate}
		education.StartDate = &d
	}
	if endDate != nil {
		d := domain.Date{Time: *endDate}
		education.EndDate = &d
	}

	if education.Honors == nil {
		education.Honors = make([]string, 0)
	}

	return &education, nil
}

// scanEducationList scans multiple education rows.
func (r *EducationRepository) scanEducationList(rows pgx.Rows) ([]domain.Education, error) {
	var educationList []domain.Education

	for rows.Next() {
		var education domain.Education
		var startDate, endDate *time.Time

		err := rows.Scan(
			&education.ID,
			&education.UserID,
			&education.Institution,
			&education.Degree,
			&education.FieldOfStudy,
			&education.Location,
			&startDate,
			&endDate,
			&education.GPA,
			&education.Honors,
			&education.DisplayOrder,
			&education.CreatedAt,
			&education.UpdatedAt,
		)
		if err != nil {
			return nil, domain.NewDatabaseError("scan education list", err)
		}

		if startDate != nil {
			d := domain.Date{Time: *startDate}
			education.StartDate = &d
		}
		if endDate != nil {
			d := domain.Date{Time: *endDate}
			education.EndDate = &d
		}

		if education.Honors == nil {
			education.Honors = make([]string, 0)
		}

		educationList = append(educationList, education)
	}

	if err := rows.Err(); err != nil {
		return nil, domain.NewDatabaseError("iterate education rows", err)
	}

	return educationList, nil
}
