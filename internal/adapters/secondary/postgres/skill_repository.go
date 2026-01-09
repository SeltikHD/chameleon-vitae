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

// SkillRepository implements ports.SkillRepository using PostgreSQL.
type SkillRepository struct {
	pool *pgxpool.Pool
}

// Create creates a new skill.
func (r *SkillRepository) Create(ctx context.Context, skill *domain.Skill) error {
	if skill.ID == "" {
		skill.ID = uuid.New().String()
	}

	skill.CreatedAt = time.Now().UTC()

	query := `
		INSERT INTO skills (
			id, user_id, name, category, proficiency_level,
			years_of_experience, is_highlighted, display_order, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	_, err := r.pool.Exec(ctx, query,
		skill.ID,
		skill.UserID,
		skill.Name,
		skill.Category,
		skill.ProficiencyLevel.Int(),
		skill.YearsOfExperience,
		skill.IsHighlighted,
		skill.DisplayOrder,
		skill.CreatedAt,
	)
	if err != nil {
		return domain.NewDatabaseError("create skill", err)
	}

	return nil
}

// GetByID retrieves a skill by ID.
func (r *SkillRepository) GetByID(ctx context.Context, id string) (*domain.Skill, error) {
	query := `
		SELECT id, user_id, name, category, proficiency_level,
			   years_of_experience, is_highlighted, display_order, created_at
		FROM skills
		WHERE id = $1
	`

	return r.scanSkill(r.pool.QueryRow(ctx, query, id))
}

// GetByUserIDAndName retrieves a skill by user ID and name.
func (r *SkillRepository) GetByUserIDAndName(ctx context.Context, userID, name string) (*domain.Skill, error) {
	query := `
		SELECT id, user_id, name, category, proficiency_level,
			   years_of_experience, is_highlighted, display_order, created_at
		FROM skills
		WHERE user_id = $1 AND LOWER(name) = LOWER($2)
	`

	return r.scanSkill(r.pool.QueryRow(ctx, query, userID, name))
}

// ListByUserID lists all skills for a user.
func (r *SkillRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Skill, error) {
	query := `
		SELECT id, user_id, name, category, proficiency_level,
			   years_of_experience, is_highlighted, display_order, created_at
		FROM skills
		WHERE user_id = $1
		ORDER BY is_highlighted DESC, display_order ASC, name ASC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, domain.NewDatabaseError("list skills by user", err)
	}
	defer rows.Close()

	return r.scanSkills(rows)
}

// ListByUserIDAndCategory lists skills filtered by category.
func (r *SkillRepository) ListByUserIDAndCategory(ctx context.Context, userID, category string) ([]domain.Skill, error) {
	query := `
		SELECT id, user_id, name, category, proficiency_level,
			   years_of_experience, is_highlighted, display_order, created_at
		FROM skills
		WHERE user_id = $1 AND LOWER(category) = LOWER($2)
		ORDER BY is_highlighted DESC, display_order ASC, name ASC
	`

	rows, err := r.pool.Query(ctx, query, userID, category)
	if err != nil {
		return nil, domain.NewDatabaseError("list skills by category", err)
	}
	defer rows.Close()

	return r.scanSkills(rows)
}

// ListHighlighted lists highlighted skills for a user.
func (r *SkillRepository) ListHighlighted(ctx context.Context, userID string) ([]domain.Skill, error) {
	query := `
		SELECT id, user_id, name, category, proficiency_level,
			   years_of_experience, is_highlighted, display_order, created_at
		FROM skills
		WHERE user_id = $1 AND is_highlighted = true
		ORDER BY display_order ASC, name ASC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, domain.NewDatabaseError("list highlighted skills", err)
	}
	defer rows.Close()

	return r.scanSkills(rows)
}

// Update updates an existing skill.
func (r *SkillRepository) Update(ctx context.Context, skill *domain.Skill) error {
	query := `
		UPDATE skills SET
			name = $2,
			category = $3,
			proficiency_level = $4,
			years_of_experience = $5,
			is_highlighted = $6,
			display_order = $7
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		skill.ID,
		skill.Name,
		skill.Category,
		skill.ProficiencyLevel.Int(),
		skill.YearsOfExperience,
		skill.IsHighlighted,
		skill.DisplayOrder,
	)
	if err != nil {
		return domain.NewDatabaseError("update skill", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSkillNotFound
	}

	return nil
}

// Upsert creates or updates a skill based on user ID and name.
func (r *SkillRepository) Upsert(ctx context.Context, skill *domain.Skill) error {
	if skill.ID == "" {
		skill.ID = uuid.New().String()
	}

	skill.CreatedAt = time.Now().UTC()

	query := `
		INSERT INTO skills (
			id, user_id, name, category, proficiency_level,
			years_of_experience, is_highlighted, display_order, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		ON CONFLICT (user_id, LOWER(name)) DO UPDATE SET
			category = EXCLUDED.category,
			proficiency_level = EXCLUDED.proficiency_level,
			years_of_experience = EXCLUDED.years_of_experience,
			is_highlighted = EXCLUDED.is_highlighted,
			display_order = EXCLUDED.display_order
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(ctx, query,
		skill.ID,
		skill.UserID,
		skill.Name,
		skill.Category,
		skill.ProficiencyLevel.Int(),
		skill.YearsOfExperience,
		skill.IsHighlighted,
		skill.DisplayOrder,
		skill.CreatedAt,
	).Scan(&skill.ID, &skill.CreatedAt)
	if err != nil {
		return domain.NewDatabaseError("upsert skill", err)
	}

	return nil
}

// BatchUpsert creates or updates multiple skills.
func (r *SkillRepository) BatchUpsert(ctx context.Context, skills []domain.Skill) (created int, updated int, err error) {
	if len(skills) == 0 {
		return 0, 0, nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, 0, domain.NewDatabaseError("begin transaction", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO skills (
			id, user_id, name, category, proficiency_level,
			years_of_experience, is_highlighted, display_order, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		ON CONFLICT (user_id, LOWER(name)) DO UPDATE SET
			category = EXCLUDED.category,
			proficiency_level = EXCLUDED.proficiency_level,
			years_of_experience = EXCLUDED.years_of_experience,
			is_highlighted = EXCLUDED.is_highlighted,
			display_order = EXCLUDED.display_order
		RETURNING (xmax = 0) as is_insert
	`

	now := time.Now().UTC()

	for i := range skills {
		skill := &skills[i]
		if skill.ID == "" {
			skill.ID = uuid.New().String()
		}
		skill.CreatedAt = now

		var isInsert bool
		err := tx.QueryRow(ctx, query,
			skill.ID,
			skill.UserID,
			skill.Name,
			skill.Category,
			skill.ProficiencyLevel.Int(),
			skill.YearsOfExperience,
			skill.IsHighlighted,
			skill.DisplayOrder,
			skill.CreatedAt,
		).Scan(&isInsert)
		if err != nil {
			return 0, 0, domain.NewDatabaseError("batch upsert skill", err)
		}

		if isInsert {
			created++
		} else {
			updated++
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, 0, domain.NewDatabaseError("commit transaction", err)
	}

	return created, updated, nil
}

// Delete removes a skill.
func (r *SkillRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM skills WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return domain.NewDatabaseError("delete skill", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSkillNotFound
	}

	return nil
}

// SearchByName searches skills by name (fuzzy match).
func (r *SkillRepository) SearchByName(ctx context.Context, userID, query string) ([]domain.Skill, error) {
	sqlQuery := `
		SELECT id, user_id, name, category, proficiency_level,
			   years_of_experience, is_highlighted, display_order, created_at
		FROM skills
		WHERE user_id = $1 AND name ILIKE $2
		ORDER BY is_highlighted DESC, display_order ASC, name ASC
	`

	rows, err := r.pool.Query(ctx, sqlQuery, userID, "%"+query+"%")
	if err != nil {
		return nil, domain.NewDatabaseError("search skills by name", err)
	}
	defer rows.Close()

	return r.scanSkills(rows)
}

// scanSkill scans a single skill row.
func (r *SkillRepository) scanSkill(row pgx.Row) (*domain.Skill, error) {
	skill := &domain.Skill{}
	var proficiencyLevel int

	err := row.Scan(
		&skill.ID,
		&skill.UserID,
		&skill.Name,
		&skill.Category,
		&proficiencyLevel,
		&skill.YearsOfExperience,
		&skill.IsHighlighted,
		&skill.DisplayOrder,
		&skill.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSkillNotFound
		}
		return nil, domain.NewDatabaseError("scan skill", err)
	}

	skill.ProficiencyLevel = domain.ProficiencyLevel(proficiencyLevel)

	return skill, nil
}

// scanSkills scans multiple skill rows.
func (r *SkillRepository) scanSkills(rows pgx.Rows) ([]domain.Skill, error) {
	skills := make([]domain.Skill, 0)

	for rows.Next() {
		skill := domain.Skill{}
		var proficiencyLevel int

		err := rows.Scan(
			&skill.ID,
			&skill.UserID,
			&skill.Name,
			&skill.Category,
			&proficiencyLevel,
			&skill.YearsOfExperience,
			&skill.IsHighlighted,
			&skill.DisplayOrder,
			&skill.CreatedAt,
		)
		if err != nil {
			return nil, domain.NewDatabaseError("scan skill row", err)
		}

		skill.ProficiencyLevel = domain.ProficiencyLevel(proficiencyLevel)
		skills = append(skills, skill)
	}

	if err := rows.Err(); err != nil {
		return nil, domain.NewDatabaseError("iterate skills", err)
	}

	return skills, nil
}
