package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"tsuskills-dbmanager/config"
	router "tsuskills-dbmanager/internal/delivery/http"
	"tsuskills-dbmanager/internal/delivery/http/handler"
	"tsuskills-dbmanager/internal/domain"
	"tsuskills-dbmanager/internal/infra/opensearch"
	"tsuskills-dbmanager/internal/infra/postgres"
	"tsuskills-dbmanager/internal/logger"
	"tsuskills-dbmanager/internal/repository"
	"tsuskills-dbmanager/internal/search"
	"tsuskills-dbmanager/internal/service"

	"github.com/google/uuid"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Logger
	appLogger, err := logger.New(&cfg.Logger.Logger)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	appLogger.Info(ctx, "Starting dbmanager service...")

	// Postgres
	pool, err := postgres.Connect(ctx, &cfg.Postgres)
	if err != nil {
		appLogger.Fatal(ctx, fmt.Sprintf("Failed to connect to Postgres: %v", err))
	}
	defer pool.Close()
	appLogger.Info(ctx, "Connected to PostgreSQL")

	// Postgres migrations
	connString := cfg.Postgres.Pool.ConnConfig.ConnString()
	if err := postgres.RunMigrations(connString, cfg.Postgres.MigrationsPath); err != nil {
		appLogger.Fatal(ctx, fmt.Sprintf("Failed to run PG migrations: %v", err))
	}
	appLogger.Info(ctx, "PostgreSQL migrations applied")

	// OpenSearch
	osClient, err := opensearch.NewClient(cfg.Search.Client)
	if err != nil {
		appLogger.Fatal(ctx, fmt.Sprintf("Failed to create OpenSearch client: %v", err))
	}

	if err := opensearch.Ping(osClient, cfg.Search.Connect); err != nil {
		// OpenSearch недоступен — не фатально, поиск просто не будет работать
		appLogger.Warn(ctx, fmt.Sprintf("OpenSearch unavailable (search disabled): %v", err))
		osClient = nil
	} else {
		appLogger.Info(ctx, "Connected to OpenSearch")
	}

	// Dependencies
	vacancyRepo := repository.NewVacancyRepository(pool)

	var vacancySearch *search.VacancySearch
	if osClient != nil {
		vacancySearch = search.NewVacancySearch(osClient, appLogger)
	}

	// Используем noopSearch если OpenSearch недоступен
	var searchImpl service.IVacancySearch
	if vacancySearch != nil {
		searchImpl = vacancySearch
	} else {
		searchImpl = &noopSearch{}
	}

	vacancyService := service.NewVacancyService(vacancyRepo, searchImpl, appLogger)
	h := handler.NewHandler(vacancyService, appLogger)
	r := router.NewRouter(h, appLogger)

	// HTTP Server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		appLogger.Info(ctx, fmt.Sprintf("Received signal: %v, shutting down...", sig))

		shutCtx, shutCancel := context.WithTimeout(context.Background(), cfg.Server.ShutDownTimeOut)
		defer shutCancel()

		if err := httpServer.Shutdown(shutCtx); err != nil {
			appLogger.Error(ctx, fmt.Sprintf("Server shutdown error: %v", err))
		}
		cancel()
	}()

	appLogger.Info(ctx, fmt.Sprintf("Server starting on %s:%d", cfg.Server.Host, cfg.Server.Port))
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.Fatal(ctx, fmt.Sprintf("Server failed: %v", err))
	}

	appLogger.Info(ctx, "Server stopped")
}

// noopSearch — заглушка для работы без OpenSearch
type noopSearch struct{}

func (n *noopSearch) IndexVacancy(_ context.Context, _ *domain.Vacancy) error { return nil }
func (n *noopSearch) DeleteVacancy(_ context.Context, _ uuid.UUID) error      { return nil }
func (n *noopSearch) SearchVacancies(_ context.Context, _ domain.VacancySearchParams) ([]uuid.UUID, int, error) {
	return nil, 0, fmt.Errorf("search unavailable: OpenSearch is not connected")
}
