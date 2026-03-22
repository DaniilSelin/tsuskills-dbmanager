package mapper

import (
	"tsuskills-dbmanager/internal/delivery/dto"
	"tsuskills-dbmanager/internal/domain"

	"github.com/google/uuid"
)

func VacancyFromCreateDTO(d dto.VacancyCreateDTO) (*domain.Vacancy, error) {
	employerID, err := uuid.Parse(d.EmployerID)
	if err != nil {
		return nil, err
	}

	return &domain.Vacancy{
		EmployerID:       employerID,
		Title:            d.Title,
		ActivityType:     domain.ActivityType{ID: d.ActivityType.ID, Name: d.ActivityType.Name},
		EmploymentType:   d.EmploymentType,
		WorkSchedule:     d.WorkSchedule,
		IsVerified:       d.IsVerified,
		Skills:           skillsFromDTO(d.Skills),
		CompensationType: d.CompensationType,
		CompensationMin:  d.CompensationMin,
		CompensationMax:  d.CompensationMax,
		Description:      d.Description,
	}, nil
}

func ApplyUpdateDTO(existing *domain.Vacancy, d dto.VacancyUpdateDTO) {
	if d.Title != "" {
		existing.Title = d.Title
	}
	if d.ActivityType != nil {
		existing.ActivityType = domain.ActivityType{ID: d.ActivityType.ID, Name: d.ActivityType.Name}
	}
	if d.EmploymentType != "" {
		existing.EmploymentType = d.EmploymentType
	}
	if d.WorkSchedule != "" {
		existing.WorkSchedule = d.WorkSchedule
	}
	if d.IsVerified != nil {
		existing.IsVerified = *d.IsVerified
	}
	if d.IsArchived != nil {
		existing.IsArchived = *d.IsArchived
	}
	if d.Skills != nil {
		existing.Skills = skillsFromDTO(d.Skills)
	}
	if d.CompensationType != "" {
		existing.CompensationType = d.CompensationType
	}
	if d.CompensationMin != nil {
		existing.CompensationMin = *d.CompensationMin
	}
	if d.CompensationMax != nil {
		existing.CompensationMax = *d.CompensationMax
	}
	if d.Description != "" {
		existing.Description = d.Description
	}
}

func VacancyToDTO(v domain.Vacancy) dto.VacancyResponseDTO {
	return dto.VacancyResponseDTO{
		ID:               v.ID.String(),
		EmployerID:       v.EmployerID.String(),
		Title:            v.Title,
		ActivityType:     dto.ActivityTypeDTO{ID: v.ActivityType.ID, Name: v.ActivityType.Name},
		EmploymentType:   string(v.EmploymentType),
		WorkSchedule:     string(v.WorkSchedule),
		IsVerified:       v.IsVerified,
		IsArchived:       v.IsArchived,
		Skills:           skillsToDTO(v.Skills),
		CompensationType: string(v.CompensationType),
		CompensationMin:  v.CompensationMin,
		CompensationMax:  v.CompensationMax,
		Description:      v.Description,
		CreatedAt:        v.CreatedAt,
		UpdatedAt:        v.UpdatedAt,
	}
}

func VacanciesToDTO(list []domain.Vacancy) []dto.VacancyResponseDTO {
	result := make([]dto.VacancyResponseDTO, 0, len(list))
	for _, v := range list {
		result = append(result, VacancyToDTO(v))
	}
	return result
}

func SearchRequestToDomain(d dto.VacancySearchRequest) domain.VacancySearchParams {
	params := domain.VacancySearchParams{
		Query:           d.Query,
		CompensationMin: d.CompensationMin,
		CompensationMax: d.CompensationMax,
		IsVerified:      d.IsVerified,
		Sort:            d.Sort,
		Limit:           d.Limit,
		Offset:          d.Offset,
	}

	for _, et := range d.EmploymentTypes {
		params.EmploymentTypes = append(params.EmploymentTypes, domain.EmploymentType(et))
	}
	for _, ws := range d.WorkSchedules {
		params.WorkSchedules = append(params.WorkSchedules, domain.WorkSchedule(ws))
	}

	return params
}

func skillsFromDTO(list []dto.SkillDTO) []domain.Skill {
	result := make([]domain.Skill, 0, len(list))
	for _, s := range list {
		result = append(result, domain.Skill{ID: s.ID, Name: s.Name})
	}
	return result
}

func skillsToDTO(list []domain.Skill) []dto.SkillDTO {
	result := make([]dto.SkillDTO, 0, len(list))
	for _, s := range list {
		result = append(result, dto.SkillDTO{ID: s.ID, Name: s.Name})
	}
	return result
}
