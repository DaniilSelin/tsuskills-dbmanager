package main

import (
	"context"
	"log"

	"tsuskills-dbmanager/config"
	"tsuskills-dbmanager/migrations"
	"tsuskills-dbmanager/internal/search"
	_ "github.com/davecgh/go-spew/spew"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("FAILED: error when load config: %v", err)
	}

	// Создаем клиент для OpenSearch 
	clientOpenSearch, err := search.NewClient(cfg.Search)
	if err != nil {
		log.Fatalf("FAILED: error when connect to opensearch: %v", err)
	}

	err = search.Ping(clientOpenSearch)
	if err != nil {
		log.Fatalf("FAILED: error when ping opensearch: %v", err)
	}

	// Запускаем миграции
	migrations.SetOpenSearchClient(clientOpenSearch)
	index, err := search.GetVacanciesIndexDocument()
	if err != nil {
		log.Fatalf("FAILED: error when get []byte vacancy index: %v", err)
	}
	migrations.UpCreateVacanciesIndex(ctx, index)

	log.Println("OpenSearch client connected and migrations applied successfully.")
}