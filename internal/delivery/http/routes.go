package http

import (
	httpBase "net/http"
	"tsuskills-dbmanager/internal/logger"

	"github.com/gorilla/mux"
)

type IHandler interface {
	CreateVacancy(w httpBase.ResponseWriter, r *httpBase.Request)
	GetVacancy(w httpBase.ResponseWriter, r *httpBase.Request)
	UpdateVacancy(w httpBase.ResponseWriter, r *httpBase.Request)
	DeleteVacancy(w httpBase.ResponseWriter, r *httpBase.Request)
	ListVacancies(w httpBase.ResponseWriter, r *httpBase.Request)
	SearchVacancies(w httpBase.ResponseWriter, r *httpBase.Request)
}

func NewRouter(h IHandler, log logger.Logger) *mux.Router {
	r := mux.NewRouter()

	r.Use(RequestIDMiddleware)
	r.Use(CORSMiddleware)
	r.Use(LoggingMiddleware(log))
	r.Use(RecoveryMiddleware(log))

	api := r.PathPrefix("/api/v1/vacancies").Subrouter()

	api.HandleFunc("", h.CreateVacancy).Methods(httpBase.MethodPost, httpBase.MethodOptions)
	api.HandleFunc("", h.ListVacancies).Methods(httpBase.MethodGet, httpBase.MethodOptions)
	api.HandleFunc("/search", h.SearchVacancies).Methods(httpBase.MethodPost, httpBase.MethodOptions)
	api.HandleFunc("/{id}", h.GetVacancy).Methods(httpBase.MethodGet, httpBase.MethodOptions)
	api.HandleFunc("/{id}", h.UpdateVacancy).Methods(httpBase.MethodPut, httpBase.MethodOptions)
	api.HandleFunc("/{id}", h.DeleteVacancy).Methods(httpBase.MethodDelete, httpBase.MethodOptions)

	// Health check
	r.HandleFunc("/health", func(w httpBase.ResponseWriter, r *httpBase.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpBase.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(httpBase.MethodGet)

	return r
}
