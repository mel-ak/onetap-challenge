package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/mel-ak/onetap-challenge/internal/adapters/cache"
	"github.com/mel-ak/onetap-challenge/internal/adapters/provider"
	"github.com/mel-ak/onetap-challenge/internal/adapters/repository"
	"github.com/mel-ak/onetap-challenge/internal/config"
	"github.com/mel-ak/onetap-challenge/internal/usecases"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(dbConn string) error {
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	// Run migrations and handle the case where there are no new migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to run")
			return nil
		}
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}

func main() {
	// Load configuration
	cfg := config.NewDefaultConfig()

	// Run migrations
	if err := runMigrations(cfg.DBConn()); err != nil {
		log.Printf("Warning: Migration error: %v", err)
		// Continue execution even if migrations fail
		// This allows the application to start even if tables already exist
	}

	// Initialize adapters
	dbRepo, err := repository.NewPostgresRepository(cfg.DBConn())
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	redisClient := cache.NewRedisClient(cfg.Redis.Host + ":" + cfg.Redis.Port)
	providerSvc := provider.NewHTTPProvider()

	// Initialize use cases
	userUsecase := usecases.NewUserUsecase(dbRepo)
	accountUsecase := usecases.NewAccountUsecase(dbRepo, redisClient)
	billUsecase := usecases.NewBillUsecase(dbRepo, providerSvc, redisClient)
	providerUsecase := usecases.NewProviderUsecase(dbRepo)

	// Setup router
	router := mux.NewRouter()
	router.HandleFunc("/health", usecases.HealthCheck).Methods(http.MethodGet)
	router.HandleFunc("/users", userUsecase.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users", userUsecase.ListUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{user_id}", userUsecase.GetUser).Methods(http.MethodGet)

	router.HandleFunc("/users/{user_id}", userUsecase.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/users/{user_id}", userUsecase.DeleteUser).Methods(http.MethodDelete)

	// Provider routes
	router.HandleFunc("/providers", providerUsecase.CreateProvider).Methods(http.MethodPost)
	router.HandleFunc("/providers", providerUsecase.ListProviders).Methods(http.MethodGet)
	router.HandleFunc("/providers/{provider_id}", providerUsecase.GetProvider).Methods(http.MethodGet)

	router.HandleFunc("/accounts/link", accountUsecase.LinkAccount).Methods(http.MethodPost)
	router.HandleFunc("/bills", billUsecase.FetchBills).Methods(http.MethodGet)
	router.HandleFunc("/accounts/{account_id}", accountUsecase.DeleteAccount).Methods(http.MethodDelete)

	// Create and start server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
