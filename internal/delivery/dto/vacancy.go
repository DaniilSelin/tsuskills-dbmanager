package dto

import (
	"time"
	"tsuskills-dbmanager/internal/domain"
)

type VacancyCreateDTO struct {
	EmployerID       string                  `json:"employer_id" validate:"required,uuid4"`
	Title            string                  `json:"title" validate:"required,min=1,max=255"`
	ActivityType     ActivityTypeDTO         `json:"activity_type" validate:"required"`
	EmploymentType   domain.EmploymentType   `json:"employment_type" validate:"required,employment_type"`
	WorkSchedule     domain.WorkSchedule     `json:"work_schedule" validate:"required,work_schedule"`
	IsVerified       bool                    `json:"is_verified"`
	Skills           []SkillDTO              `json:"skills" validate:"required"`
	CompensationType domain.CompensationType `json:"compensation_type" validate:"omitempty,compensation_type"`
	CompensationMin  float64                 `json:"compensation_min" validate:"gte=0"`
	CompensationMax  float64                 `json:"compensation_max" validate:"gtefield=CompensationMin"`
	Description      string                  `json:"description" validate:"omitempty,min=10"`
}

type VacancyUpdateDTO struct {
	Title            string                  `json:"title" validate:"omitempty,min=1,max=255"`
	ActivityType     *ActivityTypeDTO        `json:"activity_type" validate:"omitempty"`
	EmploymentType   domain.EmploymentType   `json:"employment_type" validate:"omitempty,employment_type"`
	WorkSchedule     domain.WorkSchedule     `json:"work_schedule" validate:"omitempty,work_schedule"`
	IsVerified       *bool                   `json:"is_verified"`
	IsArchived       *bool                   `json:"is_archived"`
	Skills           []SkillDTO              `json:"skills"`
	CompensationType domain.CompensationType `json:"compensation_type" validate:"omitempty,compensation_type"`
	CompensationMin  *float64                `json:"compensation_min" validate:"omitempty,gte=0"`
	CompensationMax  *float64                `json:"compensation_max"`
	Description      string                  `json:"description" validate:"omitempty,min=10"`
}

type VacancyResponseDTO struct {
	ID               string          `json:"id"`
	EmployerID       string          `json:"employer_id"`
	Title            string          `json:"title"`
	ActivityType     ActivityTypeDTO `json:"activity_type"`
	EmploymentType   string          `json:"employment_type"`
	WorkSchedule     string          `json:"work_schedule"`
	IsVerified       bool            `json:"is_verified"`
	IsArchived       bool            `json:"is_archived"`
	Skills           []SkillDTO      `json:"skills"`
	CompensationType string          `json:"compensation_type"`
	CompensationMin  float64         `json:"compensation_min"`
	CompensationMax  float64         `json:"compensation_max"`
	Description      string          `json:"description"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type VacancyListResponseDTO struct {
	Vacancies []VacancyResponseDTO `json:"vacancies"`
	Total     int                  `json:"total"`
}

type ActivityTypeDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type SkillDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required,min=1,max=100"`
}
