package migrations

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	_ "embed"
	"github.com/opensearch-project/opensearch-go"
)

var client *opensearch.Client

// SetOpenSearchClient устанавливает клиент OpenSearch для выполнения миграций.
func SetOpenSearchClient(c *opensearch.Client) {
	client = c
}

// UpCreateVacanciesIndex проверяет соответствие индекса "vacancies" переданной схеме.
// Если индекс не существует, создаёт его. Если существует, но схема отличается,
// логирует предупреждение, удаляет и создаёт заново.
func UpCreateVacanciesIndex(ctx context.Context, mapping []byte) error {
	if client == nil {
		return fmt.Errorf("OpenSearch client is not set. Call SetOpenSearchClient first.")
	}

	// Проверяем наличие индекса
	existsRes, err := client.Indices.Exists([]string{"vacancies"}, client.Indices.Exists.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error checking vacancies index existence: %w", err)
	}

	if existsRes.StatusCode == 404 {
		// Индекс отсутствует: создаём
		log.Println("WARNING: vacancies index not found, creating new index.")
		return createVacanciesIndex(ctx, mapping)
	}

	// Индекс есть: получаем текущую схему
	getRes, err := client.Indices.GetMapping(
	    client.Indices.GetMapping.WithIndex("vacancies"),
	    client.Indices.GetMapping.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("error fetching current mapping: %w", err)
	}
	defer getRes.Body.Close()

	current, err := io.ReadAll(getRes.Body)
	if err != nil {
		return fmt.Errorf("error reading current mapping: %w", err)
	}

	// Сравниваем схемы
	if !bytes.Equal(current, mapping) {
		// Схемы не совпадают: удаляем и создаём заново
		log.Println("WARNING: vacancies index mapping differs from expected, recreating index.")
		if err := DownCreateVacanciesIndex(ctx); err != nil {
			return fmt.Errorf("error deleting old vacancies index: %w", err)
		}
		return createVacanciesIndex(ctx, mapping)
	}

	// Схемы совпадают: ничего не делать
	return nil
}

// DownCreateVacanciesIndex удаляет индекс "vacancies".
func DownCreateVacanciesIndex(ctx context.Context) error {
	if client == nil {
		return fmt.Errorf("OpenSearch client is not set. Call SetOpenSearchClient first.")
	}

	res, err := client.Indices.Delete(
		[]string{"vacancies"},
		client.Indices.Delete.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("error deleting vacancies index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("OpenSearch error deleting index: %s", res.String())
	}
	return nil
}

// createVacanciesIndex создаёт индекс "vacancies" по переданной схеме.
func createVacanciesIndex(ctx context.Context, mapping []byte) error {
	res, err := client.Indices.Create(
		"vacancies",
		client.Indices.Create.WithBody(bytes.NewReader(mapping)),
		client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to create vacancies index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("OpenSearch error creating index: %s", res.String())
	}
	return nil
}
