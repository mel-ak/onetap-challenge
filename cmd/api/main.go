package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/mel-ak/onetap-challenge/internal/adapters/notification"
	"github.com/mel-ak/onetap-challenge/internal/adapters/repository/postgres"
	"github.com/mel-ak/onetap-challenge/internal/usecases"
)

func main() {
	// Initialize database connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewRepository(db)

	// Initialize notification service
	notifier := notification.NewEmailNotifier(
		os.Getenv("SMTP_FROM"),
		os.Getenv("SMTP_TO"),
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)

	// Initialize services
	billService := usecases.NewBillService(repo)

	// Start periodic updates
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start background jobs
	go func() {
		if err := billService.RefreshBills(ctx, "all"); err != nil {
			notifier.NotifyError(ctx, err, "Periodic bill refresh failed")
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down...")
	cancel()
}
