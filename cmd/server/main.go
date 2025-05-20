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

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {
	cfg := config.Config{
		Server: config.ServerConfig{
			Port: "8081",
		},
		Database: config.DatabaseConfig{
			Host:     "postgres",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			DBName:   "bill_aggregator",
			SSLMode:  "disable",
		},
		Redis: config.RedisConfig{
			Host: "localhost",
			Port: "6379",
		},
	}

	// Run migrations
	if err := runMigrations(cfg.DBConn()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
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

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(cfg.Server.Port, r))
}
