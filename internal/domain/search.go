package domain

// VacancySearchParams содержит параметры полнотекстового поиска вакансий
type VacancySearchParams struct {
	Query           string
	EmploymentTypes []EmploymentType
	WorkSchedules   []WorkSchedule
	CompensationMin *float64
	CompensationMax *float64
	IsVerified      *bool
	Sort            string // "date_desc", "date_asc", "salary_desc", "salary_asc"
	Limit           int
	Offset          int
}
