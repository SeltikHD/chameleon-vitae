# Chameleon Vitae Makefile
# ==============================================================================
# Usage: make <target>
# Run `make help` to see all available targets
# ==============================================================================

.PHONY: help dev build test lint clean infra-up infra-down

# Default target
.DEFAULT_GOAL := help

# ==============================================================================
# Variables
# ==============================================================================

BINARY_NAME := chameleon-vitae
BUILD_DIR := ./bin
GO := go
GOFLAGS := -v

# ==============================================================================
# Help
# ==============================================================================

help: ## Show this help message
	@echo "Chameleon Vitae - Available targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ==============================================================================
# Development
# ==============================================================================

dev: ## Run the server in development mode
	$(GO) run ./cmd/server/main.go

build: ## Build the binary
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server/main.go

test: ## Run all tests
	$(GO) test -v -race -cover ./...

lint: ## Run linters
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

fmt: ## Format Go code
	$(GO) fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	$(GO) vet ./...

# ==============================================================================
# Infrastructure
# ==============================================================================

infra-up: ## Start infrastructure (PostgreSQL + Gotenberg)
	podman-compose up -d

infra-down: ## Stop infrastructure
	podman-compose down

infra-logs: ## Show infrastructure logs
	podman-compose logs -f

infra-ps: ## Show running containers
	podman-compose ps

# ==============================================================================
# Database
# ==============================================================================

db-connect: ## Connect to PostgreSQL database
	podman exec -it chameleon-postgres psql -U chameleon -d chameleon_vitae

db-reset: ## Reset database (WARNING: destroys all data)
	podman-compose down -v
	podman-compose up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Database reset complete"

# ==============================================================================
# Frontend
# ==============================================================================

frontend-install: ## Install frontend dependencies
	cd frontend && pnpm install

frontend-dev: ## Run frontend in development mode
	cd frontend && pnpm dev

frontend-build: ## Build frontend for production
	cd frontend && pnpm build

# ==============================================================================
# Cleanup
# ==============================================================================

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)
	$(GO) clean -cache -testcache

clean-all: clean infra-down ## Clean everything including infrastructure
	podman volume rm chameleon-vitae_postgres_data 2>/dev/null || true

# ==============================================================================
# Release
# ==============================================================================

version: ## Show version
	@echo "Chameleon Vitae v0.1.0"
