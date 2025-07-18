package migrations

import (
    "bytes"
    "context"
    "embed"
    "fmt"

    "github.com/pressly/goose/v3"
    "github.com/opensearch-project/opensearch-go"
)

func init() {
    goose.AddMigration(UpCreateVacanciesIndex, DownCreateVacanciesIndex)
}

//go:embed 01_create_vacancies_index.json
var vacancyMapping []byte

// UpCreateVacanciesIndex создаёт индекс "vacancies" в OpenSearch.
func UpCreateVacanciesIndex(tx goose.Dialect) error {
    // Создаём OpenSearch клиент с дефолтными настройками (можно использовать env-конфиг)
    client, err := opensearch.NewClient(opensearch.Config{})
    if err != nil {
        return fmt.Errorf("failed to create OpenSearch client: %w", err)
    }

    ctx := context.Background()
    // Проверяем существование индекса
    existsRes, err := client.Indices.Exists([]string{"vacancies"}, client.Indices.Exists.WithContext(ctx))
    if err != nil {
        return fmt.Errorf("error checking vacancies index existence: %w", err)
    }
    if existsRes.StatusCode == 200 {
        // Индекс уже существует — пропускаем создание
        return nil
    }

    // Создаём индекс, используя внешне хранимый JSON
    res, err := client.Indices.Create(
        "vacancies",
        client.Indices.Create.WithBody(bytes.NewReader(vacancyMapping)),
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

// DownCreateVacanciesIndex удаляет индекс "vacancies" в OpenSearch.
func DownCreateVacanciesIndex(tx goose.Dialect) error {
    client, err := opensearch.NewClient(opensearch.Config{})
    if err != nil {
        return fmt.Errorf("failed to create OpenSearch client: %w", err)
    }
    ctx := context.Background()
    res, err := client.Indices.Delete([]string{"vacancies"}, client.Indices.Delete.WithContext(ctx))
    if err != nil {
        return fmt.Errorf("error deleting vacancies index: %w", err)
    }
    defer res.Body.Close()
    if res.IsError() {
        return fmt.Errorf("OpenSearch error deleting index: %s", res.String())
    }
    return nil
}
