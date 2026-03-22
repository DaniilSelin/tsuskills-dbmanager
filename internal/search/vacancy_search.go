package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"tsuskills-dbmanager/internal/domain"
	"tsuskills-dbmanager/internal/logger"

	osclient "github.com/opensearch-project/opensearch-go"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type VacancySearch struct {
	client *osclient.Client
	log    logger.Logger
}

func NewVacancySearch(client *osclient.Client, log logger.Logger) *VacancySearch {
	return &VacancySearch{client: client, log: log}
}

// IndexVacancy индексирует вакансию в OpenSearch
func (vs *VacancySearch) IndexVacancy(ctx context.Context, vacancy *domain.Vacancy) error {
	body, err := json.Marshal(vacancy)
	if err != nil {
		return fmt.Errorf("marshal vacancy: %w", err)
	}

	res, err := vs.client.Index(
		"vacancy_v1",
		bytes.NewReader(body),
		vs.client.Index.WithDocumentID(vacancy.ID.String()),
		vs.client.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("opensearch index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("opensearch error: %s", res.Status())
	}

	vs.log.Info(ctx, "Vacancy indexed", zap.String("id", vacancy.ID.String()))
	return nil
}

// DeleteVacancy удаляет вакансию из индекса OpenSearch
func (vs *VacancySearch) DeleteVacancy(ctx context.Context, id uuid.UUID) error {
	res, err := vs.client.Delete(
		"vacancy_v1",
		id.String(),
		vs.client.Delete.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("opensearch delete: %w", err)
	}
	defer res.Body.Close()

	return nil
}

// SearchVacancies выполняет полнотекстовый поиск вакансий
func (vs *VacancySearch) SearchVacancies(ctx context.Context, params domain.VacancySearchParams) ([]uuid.UUID, int, error) {
	query := vs.buildQuery(params)

	body, err := json.Marshal(query)
	if err != nil {
		return nil, 0, fmt.Errorf("marshal query: %w", err)
	}

	limit := params.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	res, err := vs.client.Search(
		vs.client.Search.WithContext(ctx),
		vs.client.Search.WithIndex("vacancy_v1"),
		vs.client.Search.WithBody(bytes.NewReader(body)),
		vs.client.Search.WithSize(limit),
		vs.client.Search.WithFrom(offset),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("opensearch search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("opensearch search error: %s", res.Status())
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("decode response: %w", err)
	}

	ids := make([]uuid.UUID, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		id, err := uuid.Parse(hit.ID)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	return ids, result.Hits.Total.Value, nil
}

func (vs *VacancySearch) buildQuery(params domain.VacancySearchParams) map[string]interface{} {
	musts := []map[string]interface{}{}

	if params.Query != "" {
		musts = append(musts, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  params.Query,
				"fields": []string{"Title^3", "Description", "Skills.Name"},
			},
		})
	}

	if len(params.EmploymentTypes) > 0 {
		strs := make([]string, len(params.EmploymentTypes))
		for i, et := range params.EmploymentTypes {
			strs[i] = string(et)
		}
		musts = append(musts, map[string]interface{}{
			"terms": map[string]interface{}{
				"EmploymentType": strs,
			},
		})
	}

	if len(params.WorkSchedules) > 0 {
		strs := make([]string, len(params.WorkSchedules))
		for i, ws := range params.WorkSchedules {
			strs[i] = string(ws)
		}
		musts = append(musts, map[string]interface{}{
			"terms": map[string]interface{}{
				"WorkSchedule": strs,
			},
		})
	}

	filters := []map[string]interface{}{}

	if params.CompensationMin != nil {
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{
				"CompensationMax": map[string]interface{}{"gte": *params.CompensationMin},
			},
		})
	}
	if params.CompensationMax != nil {
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{
				"CompensationMin": map[string]interface{}{"lte": *params.CompensationMax},
			},
		})
	}
	if params.IsVerified != nil {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"IsVerified": *params.IsVerified,
			},
		})
	}

	// не показываем архивные
	filters = append(filters, map[string]interface{}{
		"term": map[string]interface{}{"IsArchived": false},
	})

	if len(musts) == 0 {
		musts = append(musts, map[string]interface{}{"match_all": map[string]interface{}{}})
	}

	boolQ := map[string]interface{}{
		"must": musts,
	}
	if len(filters) > 0 {
		boolQ["filter"] = filters
	}

	q := map[string]interface{}{
		"query": map[string]interface{}{"bool": boolQ},
	}

	// Сортировка
	sort := params.Sort
	if sort == "" {
		sort = "date_desc"
	}
	switch strings.ToLower(sort) {
	case "date_asc":
		q["sort"] = []map[string]interface{}{{"CreatedAt": "asc"}}
	case "salary_desc":
		q["sort"] = []map[string]interface{}{{"CompensationMax": "desc"}}
	case "salary_asc":
		q["sort"] = []map[string]interface{}{{"CompensationMin": "asc"}}
	default: // date_desc
		q["sort"] = []map[string]interface{}{{"CreatedAt": "desc"}}
	}

	return q
}
