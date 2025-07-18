package search

import (
	"tsuskills-dbmanager/config"
	_ "tsuskills-dbmanager/internal/interfaces"
)

type VacancySearch struct {
	cfg config.Config
}

func NewVacancySearch(
		cfg config.Config, 
	) VacancySearch {
	return VacancySearch{
		cfg: cfg,
	}
}

