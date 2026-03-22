package service

import (
	"context"
	"errors"
	"time"

	"tsuskills-dbmanager/internal/domain"
	"tsuskills-dbmanager/internal/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type VacancyService struct {
	repo   IVacancyRepository
	search IVacancySearch
	log    logger.Logger
}

func NewVacancyService(repo IVacancyRepository, search IVacancySearch, log logger.Logger) *VacancyService {
	return &VacancyService{repo: repo, search: search, log: log}
}

func (s *VacancyService) CreateVacancy(ctx context.Context, v *domain.Vacancy) (uuid.UUID, domain.ErrorCode) {
	now := time.Now()
	v.ID = uuid.New()
	v.CreatedAt = now
	v.UpdatedAt = now

	id, err := s.repo.Create(ctx, v)
	if err != nil {
		s.log.Error(ctx, "CreateVacancy: repo", zap.Error(err))
		return uuid.Nil, domain.CodeInternal
	}

	// async-safe: если OpenSearch недоступен, вакансия всё равно сохранена в PG
	if err := s.search.IndexVacancy(ctx, v); err != nil {
		s.log.Warn(ctx, "CreateVacancy: opensearch index failed (non-fatal)", zap.Error(err))
	}

	s.log.Info(ctx, "CreateVacancy: success", zap.String("id", id.String()))
	return id, domain.CodeOK
}

func (s *VacancyService) GetVacancy(ctx context.Context, id uuid.UUID) (*domain.Vacancy, domain.ErrorCode) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.CodeNotFound
		}
		s.log.Error(ctx, "GetVacancy: repo", zap.Error(err))
		return nil, domain.CodeInternal
	}
	return v, domain.CodeOK
}

func (s *VacancyService) UpdateVacancy(ctx context.Context, v *domain.Vacancy) domain.ErrorCode {
	v.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, v); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.CodeNotFound
		}
		s.log.Error(ctx, "UpdateVacancy: repo", zap.Error(err))
		return domain.CodeInternal
	}

	// переиндексируем — для этого получим полную запись из PG
	full, err := s.repo.GetByID(ctx, v.ID)
	if err == nil {
		if idxErr := s.search.IndexVacancy(ctx, full); idxErr != nil {
			s.log.Warn(ctx, "UpdateVacancy: opensearch reindex failed (non-fatal)", zap.Error(idxErr))
		}
	}

	s.log.Info(ctx, "UpdateVacancy: success", zap.String("id", v.ID.String()))
	return domain.CodeOK
}

func (s *VacancyService) DeleteVacancy(ctx context.Context, id uuid.UUID) domain.ErrorCode {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.CodeNotFound
		}
		s.log.Error(ctx, "DeleteVacancy: repo", zap.Error(err))
		return domain.CodeInternal
	}

	if err := s.search.DeleteVacancy(ctx, id); err != nil {
		s.log.Warn(ctx, "DeleteVacancy: opensearch delete failed (non-fatal)", zap.Error(err))
	}

	s.log.Info(ctx, "DeleteVacancy: success", zap.String("id", id.String()))
	return domain.CodeOK
}

func (s *VacancyService) ListByEmployer(ctx context.Context, employerID uuid.UUID, limit, offset int) ([]domain.Vacancy, domain.ErrorCode) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	list, err := s.repo.ListByEmployer(ctx, employerID, limit, offset)
	if err != nil {
		s.log.Error(ctx, "ListByEmployer: repo", zap.Error(err))
		return nil, domain.CodeInternal
	}
	return list, domain.CodeOK
}

func (s *VacancyService) ListAll(ctx context.Context, limit, offset int) ([]domain.Vacancy, domain.ErrorCode) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	list, err := s.repo.ListAll(ctx, limit, offset)
	if err != nil {
		s.log.Error(ctx, "ListAll: repo", zap.Error(err))
		return nil, domain.CodeInternal
	}
	return list, domain.CodeOK
}

func (s *VacancyService) SearchVacancies(ctx context.Context, params domain.VacancySearchParams) ([]domain.Vacancy, int, domain.ErrorCode) {
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	ids, total, err := s.search.SearchVacancies(ctx, params)
	if err != nil {
		s.log.Error(ctx, "SearchVacancies: opensearch", zap.Error(err))
		return nil, 0, domain.CodeInternal
	}

	if len(ids) == 0 {
		return []domain.Vacancy{}, total, domain.CodeOK
	}

	// гидрация из Postgres
	vacancies := make([]domain.Vacancy, 0, len(ids))
	for _, id := range ids {
		v, err := s.repo.GetByID(ctx, id)
		if err != nil {
			continue // пропускаем удалённые из PG, но оставшиеся в индексе
		}
		vacancies = append(vacancies, *v)
	}

	return vacancies, total, domain.CodeOK
}
