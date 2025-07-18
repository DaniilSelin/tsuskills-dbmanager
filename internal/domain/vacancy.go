package domain

import (
	"fmt"
	"time"

	"example.com/project/internal/enums"
)

// Vacancy описывает бизнес-сущность Вакансии
type Vacancy struct {
	ID               string                 `json:"id"`
	EmployerURL      uuid.UUID              `json:"employer_url"`
	IsArchived       bool                   `json:"is_archived"`
	Title            string                 `json:"title"`
	ActivityType     ActivityType           `json:"activity_type"`
	EmploymentType   enums.EmploymentType   `json:"employment_type"`
	WorkSchedule     enums.WorkSchedule     `json:"work_schedule"`
	IsVerified       bool                   `json:"is_verified"`
	Skills           []Skill                `json:"skills"`
	CompensationType enums.CompensationType `json:"compensation_type"`
	CompensationMin  float64                `json:"compensation_min"`
	CompensationMax  float64                `json:"compensation_max"`
	Description      string                 `json:"description"`
	CreatedAt        time.Time              `json:"created_at"`
}
