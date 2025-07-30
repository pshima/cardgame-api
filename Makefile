# Makefile for Card Game API

# Variables
APP_NAME := cardgame-api
DOCKER_IMAGE := $(APP_NAME):latest
DOCKER_IMAGE_DEV := $(APP_NAME):dev
PORT := 8080

# Cross-compilation variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitHash=$(COMMIT_HASH) -w -s"
BUILD_FLAGS = -trimpath $(LDFLAGS)
DIST_DIR = dist

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: help build run test clean docker-build docker-run docker-stop docker-clean compose-up compose-down compose-logs
.PHONY: build-all build-linux build-windows build-darwin cross-compile
.PHONY: build-linux-amd64 build-linux-arm64 build-windows-amd64 build-windows-arm64 build-darwin-amd64 build-darwin-arm64

help: ## Show this help message
	@echo "$(GREEN)Card Game API Makefile$(NC)"
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

build: ## Build the Go application for current platform
	@echo "$(GREEN)Building application for current platform...$(NC)"
	go build $(BUILD_FLAGS) -o $(APP_NAME) .

run: ## Run the application locally
	@echo "$(GREEN)Running application...$(NC)"
	go run .

test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v ./...

clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	rm -f $(APP_NAME)
	rm -rf $(DIST_DIR)
	go clean

# Cross-compilation targets
$(DIST_DIR):
	mkdir -p $(DIST_DIR)

cross-compile: build-all ## Alias for build-all

build-all: build-linux build-windows build-darwin ## Build for all supported platforms
	@echo "$(GREEN)✅ All cross-platform builds completed successfully!$(NC)"
	@echo "$(YELLOW)Built binaries:$(NC)"
	@ls -la $(DIST_DIR)/ 2>/dev/null || echo "No binaries found"

# Linux builds
build-linux: build-linux-amd64 build-linux-arm64 ## Build for Linux (amd64 and arm64)

build-linux-amd64: $(DIST_DIR) ## Build for Linux amd64
	@echo "$(GREEN)Building for Linux amd64...$(NC)"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 .

build-linux-arm64: $(DIST_DIR) ## Build for Linux arm64
	@echo "$(GREEN)Building for Linux arm64...$(NC)"
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME)-linux-arm64 .

# Windows builds
build-windows: build-windows-amd64 build-windows-arm64 ## Build for Windows (amd64 and arm64)

build-windows-amd64: $(DIST_DIR) ## Build for Windows amd64
	@echo "$(GREEN)Building for Windows amd64...$(NC)"
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe .

build-windows-arm64: $(DIST_DIR) ## Build for Windows arm64
	@echo "$(GREEN)Building for Windows arm64...$(NC)"
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME)-windows-arm64.exe .

# macOS builds
build-darwin: build-darwin-amd64 build-darwin-arm64 ## Build for macOS (Intel and Apple Silicon)

build-darwin-amd64: $(DIST_DIR) ## Build for macOS amd64 (Intel)
	@echo "$(GREEN)Building for macOS amd64 (Intel)...$(NC)"
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 .

build-darwin-arm64: $(DIST_DIR) ## Build for macOS arm64 (Apple Silicon)
	@echo "$(GREEN)Building for macOS arm64 (Apple Silicon)...$(NC)"
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME)-darwin-arm64 .

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