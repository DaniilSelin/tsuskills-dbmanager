package handler

import (
	"net/http"

	"tsuskills-dbmanager/internal/delivery/dto"
	"tsuskills-dbmanager/internal/delivery/mapper"
	"tsuskills-dbmanager/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateVacancy обрабатывает POST /api/v1/vacancies
func (h *Handler) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "CreateVacancy"

	var req dto.VacancyCreateDTO
	if !h.decodeAndValidate(ctx, w, r, &req) {
		return
	}

	vacancy, err := mapper.VacancyFromCreateDTO(req)
	if err != nil {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, err.Error())
		return
	}

	id, svcCode := h.vacancySV.CreateVacancy(ctx, vacancy)
	if svcCode != domain.CodeOK {
		h.handleServiceError(ctx, w, svcCode, op)
		return
	}

	h.log.Info(ctx, op+": success", zap.String("id", id.String()))
	h.writeJSON(ctx, w, http.StatusCreated, map[string]string{
		"id":      id.String(),
		"message": "vacancy created",
	})
}

// GetVacancy обрабатывает GET /api/v1/vacancies/{id}
func (h *Handler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "GetVacancy"

	id, ok := h.extractUUIDParam(r, "id")
	if !ok {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid vacancy ID")
		return
	}

	vacancy, svcCode := h.vacancySV.GetVacancy(ctx, id)
	if svcCode != domain.CodeOK {
		h.handleServiceError(ctx, w, svcCode, op)
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, mapper.VacancyToDTO(*vacancy))
}

// UpdateVacancy обрабатывает PUT /api/v1/vacancies/{id}
func (h *Handler) UpdateVacancy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "UpdateVacancy"

	id, ok := h.extractUUIDParam(r, "id")
	if !ok {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid vacancy ID")
		return
	}

	var req dto.VacancyUpdateDTO
	if !h.decodeAndValidate(ctx, w, r, &req) {
		return
	}

	// получаем текущую вакансию для partial update
	existing, svcCode := h.vacancySV.GetVacancy(ctx, id)
	if svcCode != domain.CodeOK {
		h.handleServiceError(ctx, w, svcCode, op)
		return
	}

	mapper.ApplyUpdateDTO(existing, req)

	svcCode = h.vacancySV.UpdateVacancy(ctx, existing)
	if svcCode != domain.CodeOK {
		h.handleServiceError(ctx, w, svcCode, op)
		return
	}

	h.log.Info(ctx, op+": success", zap.String("id", id.String()))
	h.writeJSON(ctx, w, http.StatusOK, map[string]string{"message": "vacancy updated"})
}

// DeleteVacancy обрабатывает DELETE /api/v1/vacancies/{id}
func (h *Handler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "DeleteVacancy"

	id, ok := h.extractUUIDParam(r, "id")
	if !ok {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid vacancy ID")
		return
	}

	svcCode := h.vacancySV.DeleteVacancy(ctx, id)
	if svcCode != domain.CodeOK {
		h.handleServiceError(ctx, w, svcCode, op)
		return
	}

	h.log.Info(ctx, op+": success", zap.String("id", id.String()))
	h.writeJSON(ctx, w, http.StatusOK, map[string]string{"message": "vacancy deleted"})
}

// ListVacancies обрабатывает GET /api/v1/vacancies
func (h *Handler) ListVacancies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "ListVacancies"

	limit := h.queryInt(r, "limit", 20)
	offset := h.queryInt(r, "offset", 0)

	// если передан employer_id — показываем вакансии работодателя
	if eid := r.URL.Query().Get("employer_id"); eid != "" {
		employerID, err := uuid.Parse(eid)
		if err != nil {
			h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid employer_id")
			return
		}
		list, svcCode := h.vacancySV.ListByEmployer(ctx, employerID, limit, offset)
		if svcCode != domain.CodeOK {
			h.handleServiceError(ctx, w, svcCode, op)
			return
		}
		h.writeJSON(ctx, w, http.StatusOK, dto.VacancyListResponseDTO{
			Vacancies: mapper.VacanciesToDTO(list),
			Total:     len(list),
		})
		return
	}

	list, svcCode := h.vacancySV.ListAll(ctx, limit, offset)
	if svcCode != domain.CodeOK {
		h.handleServiceError(ctx, w, svcCode, op)
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, dto.VacancyListResponseDTO{
		Vacancies: mapper.VacanciesToDTO(list),
		Total:     len(list),
	})
}

// SearchVacancies обрабатывает POST /api/v1/vacancies/search
func (h *Handler) SearchVacancies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "SearchVacancies"

	var req dto.VacancySearchRequest
	if !h.decodeAndValidate(ctx, w, r, &req) {
		return
	}

	params := mapper.SearchRequestToDomain(req)

	vacancies, total, svcCode := h.vacancySV.SearchVacancies(ctx, params)
	if svcCode != domain.CodeOK {
		h.handleServiceError(ctx, w, svcCode, op)
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, dto.VacancyListResponseDTO{
		Vacancies: mapper.VacanciesToDTO(vacancies),
		Total:     total,
	})
}
