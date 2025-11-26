.DEFAULT_GOAL: help

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## Run HTTP server locally on port 8082
	@go mod tidy
	@DATABASE_URL=$${DATABASE_URL:-postgres://postgres:postgres@localhost:5435/packs?sslmode=disable} go run ./cmd/packs-service

.PHONY: infra-up
infra-up: ## Start Postgres via docker-compose
	@docker compose -f docker-compose.yml up -d db

.PHONY: infra-down
infra-down: ## Stop Postgres
	@docker compose -f docker-compose.yml down

.PHONY: db-init
db-init: ## Apply schema and seed default pack sizes
	@docker compose -f docker-compose.yml exec -T db psql -U postgres -d packs -v ON_ERROR_STOP=1 -f /schema/schema.sql
	@docker compose -f docker-compose.yml exec -T db psql -U postgres -d packs -v ON_ERROR_STOP=1 -f /schema/seeds.sql

.PHONY: db-reset
db-reset: ## Recreate database volume, start and init schema + seeds
	@docker compose -f docker-compose.yml down -v
	@docker compose -f docker-compose.yml up -d db
	@sleep 2
	@$(MAKE) db-init

.PHONY: build
build: ## Build the Go binary for production
	@mkdir -p build
	@go build -ldflags="-s -w" -o build/packs-service ./cmd/packs-service

.PHONY: test
test: ## Execute the tests in the development environment
	@go test ./... -count=1 -timeout 2m

.PHONY: docker-build
docker-build: ## Build Docker image for packs-service
	@docker build -t packs-service -f Dockerfile .

.PHONY: docker-run
docker-run: ## Run Docker container locally on port 8082
	@docker run -p 8082:8082 packs-service


