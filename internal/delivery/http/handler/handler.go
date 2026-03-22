package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"tsuskills-dbmanager/internal/delivery/validator"
	"tsuskills-dbmanager/internal/domain"
	"tsuskills-dbmanager/internal/logger"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type IVacancyService interface {
	CreateVacancy(ctx context.Context, v *domain.Vacancy) (uuid.UUID, domain.ErrorCode)
	GetVacancy(ctx context.Context, id uuid.UUID) (*domain.Vacancy, domain.ErrorCode)
	UpdateVacancy(ctx context.Context, v *domain.Vacancy) domain.ErrorCode
	DeleteVacancy(ctx context.Context, id uuid.UUID) domain.ErrorCode
	ListByEmployer(ctx context.Context, employerID uuid.UUID, limit, offset int) ([]domain.Vacancy, domain.ErrorCode)
	ListAll(ctx context.Context, limit, offset int) ([]domain.Vacancy, domain.ErrorCode)
	SearchVacancies(ctx context.Context, params domain.VacancySearchParams) ([]domain.Vacancy, int, domain.ErrorCode)
}

type Handler struct {
	vacancySV IVacancyService
	log       logger.Logger
}

func NewHandler(svc IVacancyService, l logger.Logger) *Handler {
	return &Handler{vacancySV: svc, log: l}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

func (h *Handler) writeJSON(ctx context.Context, w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error(ctx, "Failed to encode JSON", zap.Error(err))
	}
}

func (h *Handler) writeError(ctx context.Context, w http.ResponseWriter, statusCode int, code domain.ErrorCode, message string) {
	h.writeJSON(ctx, w, statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Code:    string(code),
		Message: message,
	})
}

func (h *Handler) handleServiceError(ctx context.Context, w http.ResponseWriter, code domain.ErrorCode, operation string) {
	switch code {
	case domain.CodeOK:
		return
	case domain.CodeNotFound:
		h.writeError(ctx, w, http.StatusNotFound, code, "Resource not found")
	case domain.CodeConflict:
		h.writeError(ctx, w, http.StatusConflict, code, "Resource already exists")
	case domain.CodeInvalidRequestBody:
		h.writeError(ctx, w, http.StatusBadRequest, code, "Invalid request")
	default:
		h.log.Error(ctx, operation+": internal error", zap.String("code", string(code)))
		h.writeError(ctx, w, http.StatusInternalServerError, domain.CodeInternal, "Internal server error")
	}
}

func (h *Handler) decodeAndValidate(ctx context.Context, w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid JSON")
		return false
	}
	if err := validator.ValidateStruct(dst); err != nil {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, err.Error())
		return false
	}
	return true
}

func (h *Handler) extractUUIDParam(r *http.Request, name string) (uuid.UUID, bool) {
	vars := mux.Vars(r)
	raw, ok := vars[name]
	if !ok || raw == "" {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func (h *Handler) queryInt(r *http.Request, key string, defaultVal int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}
