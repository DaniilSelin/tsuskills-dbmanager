package domain

import (
	"time"

	"github.com/google/uuid"
)

type Vacancy struct {
	ID               uuid.UUID
	EmployerID       uuid.UUID
	IsArchived       bool
	Title            string
	ActivityType     ActivityType
	EmploymentType   EmploymentType
	WorkSchedule     WorkSchedule
	IsVerified       bool
	Skills           []Skill
	CompensationType CompensationType
	CompensationMin  float64
	CompensationMax  float64
	Description      string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type ActivityType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Skill struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
