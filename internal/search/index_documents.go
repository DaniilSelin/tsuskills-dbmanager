package search

import (
	"os"
)

// GetVacanciesIndexDocument - создает []byte вид индекса вакансии
func GetVacanciesIndexDocument() ([]byte, error) {
	content, err := os.ReadFile("./internal/search/documents/vacancies_index.json")
	if err != nil {
		return nil, err
	}
	return content, nil
} 