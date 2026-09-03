.PHONY: all build run test test-coverage clean docker-up docker-down

BINARY_NAME=api

all: test build

build:
	@echo "Building binary..."
	go build -o bin/$(BINARY_NAME) ./cmd/api

run:
	@echo "Running MariaDiezmaBack API..."
	go run ./cmd/api

test:
	@echo "Running tests..."
	go test -v -race ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

clean:
	@echo "Cleaning up build artifacts..."
	rm -rf bin/ coverage.out

docker-up:
	@echo "Starting docker services..."
	docker compose up --build -d

docker-down:
	@echo "Stopping docker services..."
	docker compose down

seed:
	@echo "Injecting seed data..."
	go run ./cmd/seed

seed-docker:
	@echo "Injecting seed data into PostgreSQL container..."
	docker exec -i mariadiezma_db psql -U postgres -d mariadiezma < migrations/000002_seed_data.sql

