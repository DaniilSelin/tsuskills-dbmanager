package service

import (
	"tsuskills-dbmanager/internal/interfaces"
)

type VacancyService struct {
	cfg config.Config
	search interfaces.Search
	repo interfaces.Repository
}

func NewVacancyService(
		cfg confgi.Config, 
		search interfaces.Search,
		repo interfaces.Repository
	) {
	return VacancyService{
		cfg: cfg,
		search: search,
		repo: repo
	}
}

// реализуй CRUD

func Create