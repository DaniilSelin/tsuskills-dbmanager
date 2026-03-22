package service

import (
	"context"

	"tsuskills-dbmanager/internal/domain"

	"github.com/google/uuid"
)

type IVacancyRepository interface {
	Create(ctx context.Context, v *domain.Vacancy) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Vacancy, error)
	Update(ctx context.Context, v *domain.Vacancy) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByEmployer(ctx context.Context, employerID uuid.UUID, limit, offset int) ([]domain.Vacancy, error)
	ListAll(ctx context.Context, limit, offset int) ([]domain.Vacancy, error)
}

type IVacancySearch interface {
	IndexVacancy(ctx context.Context, vacancy *domain.Vacancy) error
	DeleteVacancy(ctx context.Context, id uuid.UUID) error
	SearchVacancies(ctx context.Context, params domain.VacancySearchParams) ([]uuid.UUID, int, error)
}
