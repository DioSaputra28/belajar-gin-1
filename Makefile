.PHONY: help migrate-up migrate-down migrate-create migrate-force migrate-version migrate-drop run build test clean

# Load environment variables from .env file
include .env
export

# Database connection string
DB_URL := mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

# Colors for output (using printf for better compatibility)
GREEN  := \\033[0;32m
YELLOW := \\033[0;33m
RED    := \\033[0;31m
NC     := \\033[0m # No Color

help: ## Show this help message
	@printf '$(GREEN)Available commands:$(NC)\n'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

# ==================== Migration Commands ====================

migrate-up: ## Run all pending migrations
	@printf "$(GREEN)Running migrations...$(NC)\n"
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" -verbose up
	@printf "$(GREEN)✓ Migrations completed successfully!$(NC)\n"

migrate-down: ## Rollback the last migration
	@printf "$(YELLOW)Rolling back last migration...$(NC)\n"
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" -verbose down 1
	@printf "$(GREEN)✓ Rollback completed!$(NC)\n"

migrate-down-all: ## Rollback all migrations
	@printf "$(RED)WARNING: This will rollback ALL migrations!$(NC)\n"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" -verbose down -all; \
		printf "$(GREEN)✓ All migrations rolled back!$(NC)\n"; \
	else \
		printf "$(YELLOW)Cancelled.$(NC)\n"; \
	fi

migrate-create: ## Create a new migration file (usage: make migrate-create name=create_users_table)
	@if [ -z "$(name)" ]; then \
		printf "$(RED)Error: Please provide a migration name$(NC)\n"; \
		printf "Usage: make migrate-create name=create_users_table\n"; \
		exit 1; \
	fi
	@printf "$(GREEN)Creating migration: $(name)$(NC)\n"
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)
	@printf "$(GREEN)✓ Migration files created!$(NC)\n"

migrate-force: ## Force set migration version (usage: make migrate-force version=1)
	@if [ -z "$(version)" ]; then \
		printf "$(RED)Error: Please provide a version number$(NC)\n"; \
		printf "Usage: make migrate-force version=1\n"; \
		exit 1; \
	fi
	@printf "$(YELLOW)Forcing migration version to $(version)...$(NC)\n"
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" force $(version)
	@printf "$(GREEN)✓ Version forced to $(version)!$(NC)\n"

migrate-version: ## Show current migration version
	@printf "$(GREEN)Current migration version:$(NC)\n"
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" version

migrate-drop: ## Drop everything in database (DANGEROUS!)
	@printf "$(RED)WARNING: This will DROP ALL TABLES in the database!$(NC)\n"
	@read -p "Are you sure? Type 'yes' to confirm: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" drop -f; \
		printf "$(GREEN)✓ Database dropped!$(NC)\n"; \
	else \
		printf "$(YELLOW)Cancelled.$(NC)\n"; \
	fi

# ==================== Application Commands ====================

run: ## Run the application
	@printf "$(GREEN)Starting application...$(NC)\n"
	go run cmd/api/main.go

build: ## Build the application
	@printf "$(GREEN)Building application...$(NC)\n"
	go build -o bin/app cmd/api/main.go
	@printf "$(GREEN)✓ Build completed! Binary: bin/app$(NC)\n"

test: ## Run tests
	@printf "$(GREEN)Running tests...$(NC)\n"
	go test -v ./...

clean: ## Clean build artifacts
	@printf "$(YELLOW)Cleaning build artifacts...$(NC)\n"
	rm -rf bin/
	@printf "$(GREEN)✓ Cleaned!$(NC)\n"

# ==================== Development Commands ====================

install-tools: ## Install required development tools
	@printf "$(GREEN)Installing golang-migrate...$(NC)\n"
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@printf "$(GREEN)✓ Tools installed!$(NC)\n"

setup: install-tools ## Setup development environment
	@printf "$(GREEN)Setting up development environment...$(NC)\n"
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		printf "$(GREEN)✓ Created .env file from .env.example$(NC)\n"; \
		printf "$(YELLOW)⚠ Please update .env with your configuration$(NC)\n"; \
	else \
		printf "$(YELLOW).env file already exists$(NC)\n"; \
	fi
	@printf "$(GREEN)Installing Go dependencies...$(NC)\n"
	go mod download
	@printf "$(GREEN)✓ Setup completed!$(NC)\n"

db-create: ## Create database
	@printf "$(GREEN)Creating database $(DB_NAME)...$(NC)\n"
	mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASSWORD) -e "CREATE DATABASE IF NOT EXISTS $(DB_NAME);"
	@printf "$(GREEN)✓ Database created!$(NC)\n"

db-drop: ## Drop database
	@printf "$(RED)WARNING: This will DROP the database $(DB_NAME)!$(NC)\n"
	@read -p "Are you sure? Type 'yes' to confirm: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASSWORD) -e "DROP DATABASE IF EXISTS $(DB_NAME);"; \
		printf "$(GREEN)✓ Database dropped!$(NC)\n"; \
	else \
		printf "$(YELLOW)Cancelled.$(NC)\n"; \
	fi

db-reset: db-drop db-create migrate-up ## Reset database (drop, create, migrate)
	@printf "$(GREEN)✓ Database reset completed!$(NC)\n"

# ==================== Utility Commands ====================

mod-tidy: ## Tidy go.mod and go.sum
	@printf "$(GREEN)Tidying Go modules...$(NC)\n"
	go mod tidy
	@printf "$(GREEN)✓ Modules tidied!$(NC)\n"

fmt: ## Format Go code
	@printf "$(GREEN)Formatting code...$(NC)\n"
	go fmt ./...
	@printf "$(GREEN)✓ Code formatted!$(NC)\n"

vet: ## Run go vet
	@printf "$(GREEN)Running go vet...$(NC)\n"
	go vet ./...
	@printf "$(GREEN)✓ Vet completed!$(NC)\n"

lint: fmt vet ## Run all linters
	@printf "$(GREEN)✓ All linting completed!$(NC)\n"
