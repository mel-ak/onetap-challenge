.PHONY: build run test lint migrate-up migrate-down docker-build docker-run docker-stop docker-clean generate-mocks

# Build the application
build:
	go build -o bin/bill-aggregator ./cmd/api

# Run the application
run:
	go run ./cmd/api

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Run database migrations up
migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5434/bill_aggregator?sslmode=disable" up

# Run database migrations down
migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5434/bill_aggregator?sslmode=disable" down

# Build Docker image
docker-build:
	docker-compose build

# Run Docker containers
docker-run:
	docker-compose up

# Stop Docker containers
docker-stop:
	docker-compose down

# Clean up Docker resources
docker-clean:
	docker-compose down -v

# Generate mocks for testing
generate-mocks:
	mockgen -source=internal/ports/repositories.go -destination=internal/mocks/repositories_mock.go
	mockgen -source=internal/ports/services.go -destination=internal/mocks/services_mock.go