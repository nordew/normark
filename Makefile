.PHONY: help docker-up docker-down docker-logs docker-clean docker-ps docker-build docker-run docker-stop docker-app-logs run build test migrate-up migrate-down migrate-create dev

help:
	@echo "Available commands:"
	@echo ""
	@echo "Docker - Infrastructure Only (PostgreSQL, Redis):"
	@echo "  make docker-up        - Start database services"
	@echo "  make docker-down      - Stop all services"
	@echo "  make docker-logs      - View database logs"
	@echo "  make docker-clean     - Stop services and remove volumes"
	@echo "  make docker-ps        - Show running containers"
	@echo ""
	@echo "Docker - Full Stack (App + Databases):"
	@echo "  make docker-build     - Build application Docker image"
	@echo "  make docker-run       - Run full stack (app + databases)"
	@echo "  make docker-stop      - Stop full stack"
	@echo "  make docker-app-logs  - View application logs"
	@echo ""
	@echo "Local Development:"
	@echo "  make dev              - Start databases + run app locally"
	@echo "  make run              - Run app locally"
	@echo "  make build            - Build application binary"
	@echo "  make test             - Run tests"
	@echo ""
	@echo "Database Migrations:"
	@echo "  make migrate-up       - Run database migrations up"
	@echo "  make migrate-down     - Run database migrations down"
	@echo "  make migrate-create   - Create new migration (usage: make migrate-create name=migration_name)"

docker-up:
	docker compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 3
	@docker compose ps

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-clean:
	docker compose down -v
	@echo "All containers stopped and volumes removed"

docker-ps:
	docker compose ps

docker-build:
	@echo "Building application Docker image..."
	docker compose build app

docker-run:
	@echo "Starting full stack (app + databases)..."
	docker compose --profile full up -d
	@echo "Waiting for services to be healthy..."
	@sleep 3
	@docker compose ps
	@echo ""
	@echo "Application is running at http://localhost:${SERVER_PORT:-8080}"
	@echo "PostgreSQL is running at localhost:${POSTGRES_PORT:-5432}"
	@echo "Redis is running at localhost:${REDIS_PORT:-6379}"

docker-stop:
	docker compose --profile full down

docker-app-logs:
	docker compose logs -f app

run:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Copy .env.example to .env first."; \
		exit 1; \
	fi
	go run cmd/app/main.go

build:
	@mkdir -p bin
	go build -ldflags="-w -s" -o bin/normark cmd/app/main.go
	@echo "Binary built at bin/normark"

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

migrate-up:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Copy .env.example to .env first."; \
		exit 1; \
	fi
	go run cmd/migrate/main.go up

migrate-down:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Copy .env.example to .env first."; \
		exit 1; \
	fi
	go run cmd/migrate/main.go down

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@mkdir -p migrations
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	touch migrations/$${timestamp}_$(name).up.sql migrations/$${timestamp}_$(name).down.sql; \
	echo "Created migrations/$${timestamp}_$(name).up.sql"; \
	echo "Created migrations/$${timestamp}_$(name).down.sql"

dev: docker-up
	@echo ""
	@echo "=== Development Environment Ready ==="
	@echo "PostgreSQL: localhost:${POSTGRES_PORT:-5432}"
	@echo "Redis: localhost:${REDIS_PORT:-6379}"
	@echo ""
	@echo "Starting application..."
	@echo ""
	@make run
