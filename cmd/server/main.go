package main

import (
	"tsuskills-dbmanager/config"
	"tsuskills-dbmanager/internal/search"
	_ "github.com/davecgh/go-spew/spew"
	"log"
) 

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("FAILED: error when load config: %v", err)
	}

	clientOpenSearch, err := search.NewClient(cfg.Search)
	if err != nil {
		log.Fatalf("FAILED: error when connect to opensearch: %v", err)
	}
	err = search.Ping(clientOpenSearch)
	if err != nil {
		log.Fatalf("FAILED: error when ping opensearch: %v", err)
	}
	//spew.Dump(clientOpenSearch)
}