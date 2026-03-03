.PHONY: help up down restart logs ps clean build test

# Default target
help:
	@echo "Authentication Project - Available Commands:"
	@echo ""
	@echo "  make up          - Start all services"
	@echo "  make down        - Stop all services"
	@echo "  make restart     - Restart all services"
	@echo "  make logs        - Show logs (all services)"
	@echo "  make logs-api    - Show API logs"
	@echo "  make logs-client - Show client logs"
	@echo "  make logs-db     - Show database logs"
	@echo "  make ps          - Show running containers"
	@echo "  make clean       - Stop and remove everything"
	@echo "  make build       - Rebuild all services"
	@echo "  make test        - Run all tests"
	@echo "  make test-api    - Run API tests"
	@echo "  make test-client - Run client tests"
	@echo ""

# Start all services
up:
	@echo "Starting all services..."
	cd deployment/compose && docker-compose --env-file ../../.env up -d
	@echo "Services started! Access:"
	@echo "  - Frontend: http://localhost:5173"
	@echo "  - Backend:  http://localhost:8080"

# Stop all services
down:
	@echo "Stopping all services..."
	cd deployment/compose && docker-compose --env-file ../../.env down

# Restart all services
restart:
	@echo "Restarting all services..."
	cd deployment/compose && docker-compose --env-file ../../.env restart

# Show logs
logs:
	cd deployment/compose && docker-compose --env-file ../../.env logs -f

logs-api:
	cd deployment/compose && docker-compose --env-file ../../.env logs -f api

logs-client:
	cd deployment/compose && docker-compose --env-file ../../.env logs -f client

logs-db:
	cd deployment/compose && docker-compose --env-file ../../.env logs -f database

logs-migration:
	cd deployment/compose && docker-compose --env-file ../../.env logs migration

# Show running containers
ps:
	cd deployment/compose && docker-compose --env-file ../../.env ps

# Clean everything
clean:
	@echo "Cleaning up..."
	cd deployment/compose && docker-compose --env-file ../../.env down -v --remove-orphans
	@echo "Cleaned!"

# Rebuild services
build:
	@echo "Rebuilding all services..."
	cd deployment/compose && docker-compose --env-file ../../.env up -d --build

build-api:
	@echo "Rebuilding API..."
	cd deployment/compose && docker-compose --env-file ../../.env up -d --build api

build-client:
	@echo "Rebuilding client..."
	cd deployment/compose && docker-compose --env-file ../../.env up -d --build client

# Run tests
test: test-api test-client

test-api:
	@echo "Running API tests..."
	cd apps/api && go test ./...

test-client:
	@echo "Running client tests..."
	cd apps/client && npm run test:unit:run

# Database commands
db-shell:
	cd deployment/compose && docker-compose --env-file ../../.env exec database psql -U auth_user -d auth_db

db-migrate:
	cd deployment/compose && docker-compose --env-file ../../.env up migration

# Development
dev-api:
	cd apps/api && go run cmd/server/main.go

dev-client:
	cd apps/client && npm run dev

# Install dependencies
install:
	@echo "Installing dependencies..."
	cd apps/api && go mod download
	cd apps/client && npm install
	@echo "Dependencies installed!"
