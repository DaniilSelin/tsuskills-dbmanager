package dto

type VacancySearchRequest struct {
	Query           string   `json:"q" form:"q"`
	EmploymentTypes []string `json:"employment_types" form:"employment_types"`
	WorkSchedules   []string `json:"work_schedules" form:"work_schedules"`
	CompensationMin *float64 `json:"compensation_min" form:"compensation_min"`
	CompensationMax *float64 `json:"compensation_max" form:"compensation_max"`
	IsVerified      *bool    `json:"is_verified" form:"is_verified"`
	Sort            string   `json:"sort" form:"sort"`
	Limit           int      `json:"limit" form:"limit"`
	Offset          int      `json:"offset" form:"offset"`
}
