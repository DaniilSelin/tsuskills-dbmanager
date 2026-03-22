package validator

import (
	"fmt"
	"strings"

	"tsuskills-dbmanager/internal/domain"

	validator "github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("employment_type", employmentTypeValidator)
	validate.RegisterValidation("compensation_type", compensationTypeValidator)
	validate.RegisterValidation("work_schedule", workScheduleValidator)
}

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var messages []string
	for _, ve := range err.(validator.ValidationErrors) {
		messages = append(messages, formatMessage(ve.Field(), ve.Tag(), ve.Param()))
	}

	return fmt.Errorf("validation failed: %s", strings.Join(messages, "; "))
}

func employmentTypeValidator(fl validator.FieldLevel) bool {
	return domain.EmploymentType(fl.Field().String()).IsValid()
}

func compensationTypeValidator(fl validator.FieldLevel) bool {
	return domain.CompensationType(fl.Field().String()).IsValid()
}

func workScheduleValidator(fl validator.FieldLevel) bool {
	return domain.WorkSchedule(fl.Field().String()).IsValid()
}

func formatMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be >= %s", field, param)
	case "gtefield":
		return fmt.Sprintf("%s must be >= %s", field, param)
	case "employment_type":
		return fmt.Sprintf("%s is not a valid employment type", field)
	case "compensation_type":
		return fmt.Sprintf("%s is not a valid compensation type", field)
	case "work_schedule":
		return fmt.Sprintf("%s is not a valid work schedule", field)
	default:
		return fmt.Sprintf("%s failed validation: %s", field, tag)
	}
}
