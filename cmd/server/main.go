package main

import (
	"database/sql"
	"log"
	"net/http"

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

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func main() {
	cfg := config.Config{
		DBConn:    "postgres://onetapuser:onetappassword@localhost:5433/onetapdb?sslmode=disable",
		RedisAddr: "localhost:6379",
		Port:      ":8080",
	}

	// Run migrations
	if err := runMigrations(cfg.DBConn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize adapters
	dbRepo, err := repository.NewPostgresRepository(cfg.DBConn)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	redisClient := cache.NewRedisClient(cfg.RedisAddr)
	providerSvc := provider.NewHTTPProvider()

	// Initialize use cases
	userUsecase := usecases.NewUserUsecase(dbRepo)
	accountUsecase := usecases.NewAccountUsecase(dbRepo, redisClient)
	billUsecase := usecases.NewBillUsecase(dbRepo, providerSvc, redisClient)

	// Set up HTTP router
	r := mux.NewRouter()
	r.HandleFunc("/health", usecases.HealthCheck).Methods("GET")
	r.HandleFunc("/users", userUsecase.CreateUser).Methods("POST")
	r.HandleFunc("/users/{user_id}", userUsecase.GetUser).Methods("GET")
	r.HandleFunc("/users/{user_id}", userUsecase.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{user_id}", userUsecase.DeleteUser).Methods("DELETE")
	r.HandleFunc("/accounts/link", accountUsecase.LinkAccount).Methods("POST")
	r.HandleFunc("/bills", billUsecase.FetchBills).Methods("GET")
	r.HandleFunc("/accounts/{account_id}", accountUsecase.DeleteAccount).Methods("DELETE")

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, r))
}
