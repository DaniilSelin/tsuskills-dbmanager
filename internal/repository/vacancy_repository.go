package repository

import (
	"context"
	"errors"
	"fmt"

	"tsuskills-dbmanager/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VacancyRepository struct {
	pool *pgxpool.Pool
}

func NewVacancyRepository(pool *pgxpool.Pool) *VacancyRepository {
	return &VacancyRepository{pool: pool}
}

func (r *VacancyRepository) Create(ctx context.Context, v *domain.Vacancy) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Upsert activity type
	var activityTypeID int
	err = tx.QueryRow(ctx,
		`INSERT INTO activity_types (name) VALUES ($1)
		 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		 RETURNING id`, v.ActivityType.Name,
	).Scan(&activityTypeID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("upsert activity type: %w", err)
	}

	// Insert vacancy
	_, err = tx.Exec(ctx,
		`INSERT INTO vacancies
		 (id, employer_id, is_archived, title, activity_type_id, employment_type,
		  work_schedule, is_verified, compensation_type, compensation_min,
		  compensation_max, description, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		v.ID, v.EmployerID, v.IsArchived, v.Title, activityTypeID,
		string(v.EmploymentType), string(v.WorkSchedule), v.IsVerified,
		string(v.CompensationType), v.CompensationMin, v.CompensationMax,
		v.Description, v.CreatedAt, v.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert vacancy: %w", err)
	}

	// Upsert skills and link
	for _, s := range v.Skills {
		var skillID int
		err = tx.QueryRow(ctx,
			`INSERT INTO skills (name) VALUES ($1)
			 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			 RETURNING id`, s.Name,
		).Scan(&skillID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("upsert skill %q: %w", s.Name, err)
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO vacancy_skills (vacancy_id, skill_id) VALUES ($1, $2)
			 ON CONFLICT DO NOTHING`, v.ID, skillID,
		)
		if err != nil {
			return uuid.Nil, fmt.Errorf("link skill: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit: %w", err)
	}

	return v.ID, nil
}

func (r *VacancyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Vacancy, error) {
	var v domain.Vacancy
	var atID int
	var atName string

	err := r.pool.QueryRow(ctx,
		`SELECT v.id, v.employer_id, v.is_archived, v.title,
		        at.id, at.name, v.employment_type, v.work_schedule,
		        v.is_verified, v.compensation_type, v.compensation_min,
		        v.compensation_max, v.description, v.created_at, v.updated_at
		 FROM vacancies v
		 JOIN activity_types at ON at.id = v.activity_type_id
		 WHERE v.id = $1`, id,
	).Scan(
		&v.ID, &v.EmployerID, &v.IsArchived, &v.Title,
		&atID, &atName, &v.EmploymentType, &v.WorkSchedule,
		&v.IsVerified, &v.CompensationType, &v.CompensationMin,
		&v.CompensationMax, &v.Description, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get vacancy: %w", err)
	}
	v.ActivityType = domain.ActivityType{ID: atID, Name: atName}

	skills, err := r.getVacancySkills(ctx, id)
	if err != nil {
		return nil, err
	}
	v.Skills = skills

	return &v, nil
}

func (r *VacancyRepository) Update(ctx context.Context, v *domain.Vacancy) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var activityTypeID int
	err = tx.QueryRow(ctx,
		`INSERT INTO activity_types (name) VALUES ($1)
		 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		 RETURNING id`, v.ActivityType.Name,
	).Scan(&activityTypeID)
	if err != nil {
		return fmt.Errorf("upsert activity type: %w", err)
	}

	tag, err := tx.Exec(ctx,
		`UPDATE vacancies SET
		    employer_id = $1, title = $2, activity_type_id = $3,
		    employment_type = $4, work_schedule = $5, is_verified = $6,
		    compensation_type = $7, compensation_min = $8, compensation_max = $9,
		    description = $10, is_archived = $11, updated_at = NOW()
		 WHERE id = $12`,
		v.EmployerID, v.Title, activityTypeID,
		string(v.EmploymentType), string(v.WorkSchedule), v.IsVerified,
		string(v.CompensationType), v.CompensationMin, v.CompensationMax,
		v.Description, v.IsArchived, v.ID,
	)
	if err != nil {
		return fmt.Errorf("update vacancy: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	// Пересоздаём связи со скиллами
	_, _ = tx.Exec(ctx, `DELETE FROM vacancy_skills WHERE vacancy_id = $1`, v.ID)
	for _, s := range v.Skills {
		var skillID int
		err = tx.QueryRow(ctx,
			`INSERT INTO skills (name) VALUES ($1)
			 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			 RETURNING id`, s.Name,
		).Scan(&skillID)
		if err != nil {
			return fmt.Errorf("upsert skill: %w", err)
		}
		_, err = tx.Exec(ctx,
			`INSERT INTO vacancy_skills (vacancy_id, skill_id) VALUES ($1, $2)
			 ON CONFLICT DO NOTHING`, v.ID, skillID)
		if err != nil {
			return fmt.Errorf("link skill: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *VacancyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM vacancies WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete vacancy: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *VacancyRepository) ListByEmployer(ctx context.Context, employerID uuid.UUID, limit, offset int) ([]domain.Vacancy, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT v.id, v.employer_id, v.is_archived, v.title,
		        at.id, at.name, v.employment_type, v.work_schedule,
		        v.is_verified, v.compensation_type, v.compensation_min,
		        v.compensation_max, v.description, v.created_at, v.updated_at
		 FROM vacancies v
		 JOIN activity_types at ON at.id = v.activity_type_id
		 WHERE v.employer_id = $1
		 ORDER BY v.created_at DESC
		 LIMIT $2 OFFSET $3`, employerID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list by employer: %w", err)
	}
	defer rows.Close()

	return r.scanVacancies(ctx, rows)
}

func (r *VacancyRepository) ListAll(ctx context.Context, limit, offset int) ([]domain.Vacancy, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT v.id, v.employer_id, v.is_archived, v.title,
		        at.id, at.name, v.employment_type, v.work_schedule,
		        v.is_verified, v.compensation_type, v.compensation_min,
		        v.compensation_max, v.description, v.created_at, v.updated_at
		 FROM vacancies v
		 JOIN activity_types at ON at.id = v.activity_type_id
		 WHERE v.is_archived = FALSE
		 ORDER BY v.created_at DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list all: %w", err)
	}
	defer rows.Close()

	return r.scanVacancies(ctx, rows)
}

func (r *VacancyRepository) scanVacancies(ctx context.Context, rows pgx.Rows) ([]domain.Vacancy, error) {
	var result []domain.Vacancy
	for rows.Next() {
		var v domain.Vacancy
		var atID int
		var atName string

		if err := rows.Scan(
			&v.ID, &v.EmployerID, &v.IsArchived, &v.Title,
			&atID, &atName, &v.EmploymentType, &v.WorkSchedule,
			&v.IsVerified, &v.CompensationType, &v.CompensationMin,
			&v.CompensationMax, &v.Description, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan vacancy: %w", err)
		}
		v.ActivityType = domain.ActivityType{ID: atID, Name: atName}

		skills, err := r.getVacancySkills(ctx, v.ID)
		if err != nil {
			return nil, err
		}
		v.Skills = skills

		result = append(result, v)
	}
	return result, rows.Err()
}

func (r *VacancyRepository) getVacancySkills(ctx context.Context, vacancyID uuid.UUID) ([]domain.Skill, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT s.id, s.name FROM skills s
		 JOIN vacancy_skills vs ON vs.skill_id = s.id
		 WHERE vs.vacancy_id = $1
		 ORDER BY s.name`, vacancyID,
	)
	if err != nil {
		return nil, fmt.Errorf("get skills: %w", err)
	}
	defer rows.Close()

	var skills []domain.Skill
	for rows.Next() {
		var s domain.Skill
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, fmt.Errorf("scan skill: %w", err)
		}
		skills = append(skills, s)
	}
	return skills, rows.Err()
}
