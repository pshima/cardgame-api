# Makefile for Card Game API

# Variables
APP_NAME := cardgame-api
DOCKER_IMAGE := $(APP_NAME):latest
DOCKER_IMAGE_DEV := $(APP_NAME):dev
PORT := 8080

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: help build run test clean docker-build docker-run docker-stop docker-clean compose-up compose-down compose-logs

help: ## Show this help message
	@echo "$(GREEN)Card Game API Makefile$(NC)"
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

build: ## Build the Go application
	@echo "$(GREEN)Building application...$(NC)"
	go build -o $(APP_NAME) .

run: ## Run the application locally
	@echo "$(GREEN)Running application...$(NC)"
	go run .

test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v ./...

clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	rm -f $(APP_NAME)
	go clean

docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run Docker container
	@echo "$(GREEN)Running Docker container...$(NC)"
	docker run -d \
		--name $(APP_NAME) \
		-p $(PORT):$(PORT) \
		-e LOG_LEVEL=INFO \
		-e GIN_MODE=release \
		--restart unless-stopped \
		$(DOCKER_IMAGE)
	@echo "$(GREEN)Container started on http://localhost:$(PORT)$(NC)"

docker-stop: ## Stop Docker container
	@echo "$(YELLOW)Stopping Docker container...$(NC)"
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

docker-clean: docker-stop ## Clean Docker images and containers
	@echo "$(RED)Cleaning Docker artifacts...$(NC)"
	docker rmi $(DOCKER_IMAGE) || true

docker-logs: ## View Docker container logs
	@echo "$(GREEN)Showing container logs...$(NC)"
	docker logs -f $(APP_NAME)

compose-up: ## Start services with Docker Compose
	@echo "$(GREEN)Starting services with Docker Compose...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Services started:$(NC)"
	@echo "  - API: http://localhost:$(PORT)"
	@echo "  - API Docs: http://localhost:$(PORT)/api-docs"
	@echo "  - Metrics: http://localhost:$(PORT)/metrics"

compose-down: ## Stop services with Docker Compose
	@echo "$(YELLOW)Stopping services...$(NC)"
	docker-compose down

compose-logs: ## View Docker Compose logs
	@echo "$(GREEN)Showing service logs...$(NC)"
	docker-compose logs -f

compose-monitoring: ## Start with monitoring stack (Prometheus + Grafana)
	@echo "$(GREEN)Starting services with monitoring...$(NC)"
	docker-compose --profile monitoring up -d
	@echo "$(GREEN)Services started:$(NC)"
	@echo "  - API: http://localhost:$(PORT)"
	@echo "  - Prometheus: http://localhost:9090"
	@echo "  - Grafana: http://localhost:3000 (admin/admin)"

# Development targets
dev-build: ## Build development Docker image with live reload
	@echo "$(GREEN)Building development image...$(NC)"
	docker build -f Dockerfile.dev -t $(DOCKER_IMAGE_DEV) .

dev-run: ## Run development container with volume mount
	@echo "$(GREEN)Running development container...$(NC)"
	docker run -it --rm \
		--name $(APP_NAME)-dev \
		-p $(PORT):$(PORT) \
		-v $(PWD):/app \
		-e LOG_LEVEL=DEBUG \
		-e GIN_MODE=debug \
		$(DOCKER_IMAGE_DEV)

# Utility targets
generate-cards: ## Generate card images
	@echo "$(GREEN)Generating card images...$(NC)"
	go run generate_cards.go

lint: ## Run linter
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run ./...

fmt: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	go fmt ./...

# Security targets
security-scan: ## Run security scans on the container
	@echo "$(GREEN)Running security scans...$(NC)"
	./scripts/security-scan.sh $(DOCKER_IMAGE)

docker-build-secure: ## Build Docker image with security scanning
	@echo "$(GREEN)Building Docker image with security checks...$(NC)"
	docker build -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)Running security scan...$(NC)"
	./scripts/security-scan.sh $(DOCKER_IMAGE)

compose-security: ## Run with security-enhanced configuration
	@echo "$(GREEN)Starting with security-enhanced configuration...$(NC)"
	docker-compose -f docker-compose.security.yml up -d
	@echo "$(GREEN)Security-enhanced services started:$(NC)"
	@echo "  - API: http://localhost:8080 (localhost only)"
	@echo "  - Note: Running with read-only filesystem and security constraints"

.DEFAULT_GOAL := help