# =========================
# Configs
# =========================

# Project
PROJECT_NAME     := go-api
PROJECT_BIN_PATH := bin
PROJECT_VERSION  := $(shell cat config/version.txt | tr -d "[:space:]")
PROJECT_ZIP_NAME := $(PROJECT_NAME)_backend_$(PROJECT_VERSION).zip

# Docker
DOCKER			:= docker
DOCKER_COMPOSE	:= docker compose
DOCKER_BASE     ?= $(or ${base}, built)
DOCKER_SERVICE  ?= $(or ${service}, backend)
COMPOSE_SERVICE := $(DOCKER_COMPOSE) --env-file config/.env -f build/docker/services.compose.yml
COMPOSE_COMMAND := BASE=${DOCKER_BASE} $(DOCKER_COMPOSE) --env-file config/.env -f build/docker/compose.yml

# Go build
GO          := go
GOOS        := linux
GOARCH      := amd64
BUILD_FLAGS := -ldflags "-w -s"
GOBUILD     := GOEXPERIMENT=jsonv2 CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build $(BUILD_FLAGS)

# Colors
GREEN  := \033[1;32m
YELLOW := \033[1;33m
BLUE   := \033[1;34m
CYAN   := \033[1;36m
MAGENTA:= \033[1;35m
RED    := \033[1;31m
RED_L  := \033[91m
RESET  := \033[0m

GOPATH          ?= $(HOME)/go
SWAG_BIN        := $(GOPATH)/bin/swag
GOFUMPT_BIN     := $(GOPATH)/bin/gofumpt
GOIMPORTS_BIN   := $(GOPATH)/bin/goimports
GOVULNCHECK_BIN := $(GOPATH)/bin/govulncheck
LINTER_BIN      := $(GOPATH)/bin/golangci-lint
MIGRATE_BIN     := $(GOPATH)/bin/migrate

# =========================
# Helpers
# =========================
define clean_dangling_images
	@echo "$(YELLOW)🧹 Cleaning up dangling Docker images...$(RESET)\n"
	@if test -n "$$($(DOCKER) images -f "dangling=true" -q)"; then \
		docker rmi $$($(DOCKER) images -f "dangling=true" -q); \
	fi > /dev/null
endef

define check_tool
	NAME="$(1)"; \
	VERSION_CMD="$(2)"; \
	VERSION=$$($$VERSION_CMD 2>/dev/null | head -n 2 | tail -n -1); \
	if [ $$? -ne 0 ] || [ -z "$$VERSION" ]; then \
		printf "$(RED)❌ %-20s %-10s %-12s %-30s$(RESET)\n" "$$NAME" "$(3)" "$(4)" "Install: $(5)"; \
		MISSING=1; \
	else \
		printf "✅ %-20s %-10s %-12s %-30s\n" "$$NAME" "$(3)" "$(4)" "$$VERSION"; \
	fi
endef

.PHONY: all
all: help
help: ## Display available commands and their descriptions
	@echo "$(CYAN)Usage:$(RESET)"
	@echo "  make [COMMAND]\n"
	@echo "$(CYAN)Example:$(RESET)"
	@echo "  make build\n"
	@echo "$(CYAN)Commands:$(RESET)\n"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(CYAN)%-30s$(RESET) %s\n", $$1, $$2}'

.PHONY: check-system
check-system: ## Check required tools
	@MISSING=0; \
    printf "  ========================================\n"; \
    printf "   Environment Tooling Check\n"; \
    printf "  ========================================\n\n"; \
    printf "$(CYAN)%-22s %-8s %-10s %-40s$(RESET)\n" "  Tool" "Dev" "Deploy" "Version"; \
    printf "$(BLUE)%-22s %-8s %-10s %-40s$(RESET)\n" "  -------" "-------" "-------" "-------"; \
	\
	$(call check_tool,Go,go version,✅,✅,Visit: https://go.dev/doc/install); \
	$(call check_tool,Docker,docker --version,❌,✅,Visit: https://docs.docker.com/engine/install/ubuntu/); \
	$(call check_tool,Docker Compose,docker compose version,❌,✅,Visit: https://docs.docker.com/engine/install/ubuntu/); \
	$(call check_tool,Swag,$(SWAG_BIN) -v,✅,❌,go install github.com/swaggo/swag/cmd/swag@v1.16.3); \
	$(call check_tool,Go Format,$(GOFUMPT_BIN) --version,✅,❌,go install mvdan.cc/gofumpt@v0.9.2); \
	$(call check_tool,Go Imports,$(GOIMPORTS_BIN) --help,✅,❌,go install golang.org/x/tools/cmd/goimports@latest); \
	$(call check_tool,Go Vulnerability,$(GOVULNCHECK_BIN) --version,✅,❌,go install golang.org/x/vuln/cmd/govulncheck@v1.1.4); \
	$(call check_tool,GolangCI Lint,$(LINTER_BIN) version,✅,❌,go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5); \
	$(call check_tool,Migrate,$(MIGRATE_BIN) -version,✅,❌,go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest); \
	$(call check_tool,Zip,zip -v,✅,❌,sudo apt install zip); \
    \
    if [ $$MISSING -eq 0 ]; then \
    	printf "\n✅ All required tools found!"; \
    	printf "\n✅ Environment is ready 🚀\n\n"; \
    else \
    	printf "\n⚠️ Some tools are missing"; \
    	printf "\n⚠️ Check above\n\n"; \
    fi

# =========================
# 1. Environment & Setup
# =========================
.PHONY: init
init: ## Create environment file
	@echo "$(BLUE)⚙️  Initializing environment setup...$(RESET)"
	@chmod +x config/env.sh && config/env.sh && mv .env config/
	@echo "$(GREEN)✅ Environment file successfully created and configured!$(RESET)\n"

.PHONY: tidy
tidy: ## Clean and tidy dependencies
	@echo "$(YELLOW)🔧 Cleaning and tidying Go dependencies...$(RESET)"
	@$(GO) mod tidy -v 2>&1 > /dev/null
	@echo "$(GREEN)✅Dependencies tidied successfully.$(RESET)\n"

.PHONY: list-updates
list-updates: ## List available direct dependency updates
	@echo "$(BLUE)📦 Checking for available direct dependency updates...$(RESET)"
	@$(GO) list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all
	@echo "$(GREEN)✅ Dependency check complete!$(RESET)\n"

.PHONY: scaffold
scaffold: ## Generate boilerplate code for a new module (Usage: make scaffold module=MyModule)
	@echo "$(BLUE)🚀 Generating code for module $(module)...$(RESET)"
	@$(GO) run cmd/scaffold/main.go -name $(module)
	@echo "$(GREEN)✅ Module $(module) structure created!$(RESET)\n"
	@echo "$(YELLOW)Don't forget to run 'make tidy' and update DI container!$(RESET)\n"

# =========================
# 2. Quality & Testing
# =========================
.PHONY: format
format: ## Fix code format issues
	@echo "$(BLUE)📝 Formatting code to fix style issues...$(RESET)"
	@$(GOIMPORTS_BIN) -w .
	@$(GOFUMPT_BIN) -w -l .
	@echo "$(GREEN)✅  Code formatting complete! All issues fixed.$(RESET)\n"

.PHONY: lint
lint: ## Run linter
	@echo "$(BLUE)🔍 Running linter...$(RESET)"
	@$(LINTER_BIN) run ./...
	@echo "$(GREEN)✅  Lint check passed!$(RESET)\n"

.PHONY: audit
audit: ## Conduct quality checks (Vet, Vuln, Lint)
	@echo "$(BLUE)🔎 Running code audit...$(RESET)"
	@$(GO) mod verify 2>&1 > /dev/null
	@$(GO) vet ./... 2>&1 > /dev/null
	@$(GOVULNCHECK_BIN) -show verbose ./...
	@$(MAKE) lint
	@echo "$(GREEN)✅  Code audit finished!$(RESET)\n"

.PHONY: test
test: ## Run tests and generate coverage report
	@echo "$(BLUE)🔍 Running tests...$(RESET)"
	@-$(GO) clean -testcache
	@-$(GO) test $$(go list ./... | grep -v /cmd) -coverprofile cover.out
	@-$(GO) tool cover -html=cover.out -o coverage.html
	@echo "$(GREEN)✅  Tests completed! Coverage: $$(go tool cover -func=cover.out | tail -1)$(RESET)\n"
	@-xdg-open ./coverage.html 2>/dev/null || true

.PHONY: benchmark
benchmark: ## Benchmark code performance
	@echo "$(MAGENTA)⚡ Running benchmarks...$(RESET)"
	@$(GO) test ./... -benchmem -bench=. -run=^Benchmark_$ 2>&1 > /dev/null
	@echo "$(GREEN)✅ Benchmark completed!$(RESET)\n"

# =========================
# 3. Build & Artifacts
# =========================
.PHONY: swag
swag: ## Update swagger files
	@echo "$(BLUE)📄 Updating Swagger API documentation...$(RESET)"
	@$(SWAG_BIN) init -g cmd/backend/main.go --parseDependency 2>&1 > /dev/null
	@echo "$(GREEN)✅ Swagger files updated successfully.$(RESET)\n"

.PHONY: build
build: ## Build all applications from source code
	@echo "$(BLUE)🏗️ Building application...$(RESET)"
	@-rm -r $(PROJECT_BIN_PATH)/ 2> /dev/null
	@${GOBUILD} -o $(PROJECT_BIN_PATH)/backend cmd/backend/main.go
	@echo "$(GREEN)✅  Build completed successfully!$(RESET)\n"

.PHONY: zip
zip: build ## Build the application and zip it
	@echo "$(BLUE)📦 Zipping application...$(RESET)"
	@-rm $(PROJECT_ZIP_NAME) 2> /dev/null
	@cp .github/README.md README.md
	@zip -qr $(PROJECT_ZIP_NAME) config/.env build/* bin/ README.md Makefile
	@rm README.md
	@echo "$(GREEN)✅  Zipped successfully!$(RESET)\n"

# =========================
# 4. Database Migrations
# =========================
MIGRATIONS_DIR := build/db/migrations

# Extract DB vars from .env (removing single quotes, comments and whitespace)
DB_USER := $(shell grep "^POSTGRES_USER" config/.env | cut -d '=' -f2 | cut -d '\#' -f1 | tr -d "'[:space:]")
DB_PASS := $(shell grep "^POSTGRES_PASS" config/.env | cut -d '=' -f2 | cut -d '\#' -f1 | tr -d "'[:space:]")
DB_PORT := $(shell grep "^POSTGRES_PORT" config/.env | cut -d '=' -f2 | cut -d '\#' -f1 | tr -d "'[:space:]")
DB_NAME := $(shell grep "^POSTGRES_BASE" config/.env | cut -d '=' -f2 | cut -d '\#' -f1 | tr -d "'[:space:]")
# Migration URL (always localhost for running from host)
MIGRATE_URL := "postgres://$(DB_USER):$(DB_PASS)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable"

.PHONY: migrate-create
migrate-create: ## Create a new migration file (Usage: make migrate-create name=create_users_table)
	@echo "$(BLUE)🔨 Creating new migration file...$(RESET)"
	@$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
	@echo "$(GREEN)✅ Migration file created!$(RESET)\n"

.PHONY: migrate-up
migrate-up: ## Run database migrations
	@echo "$(BLUE)🔄 Running migrations...$(RESET)"
	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(MIGRATE_URL) up
	@echo "$(GREEN)✅ Migrations completed!$(RESET)\n"

.PHONY: migrate-down
migrate-down: ## Rollback the last migration
	@echo "$(BLUE)🔄 Rolling back last migration...$(RESET)"
	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(MIGRATE_URL) down 1
	@echo "$(GREEN)✅ Rollback completed!$(RESET)\n"

.PHONY: migrate-reset
migrate-reset: ## Rollback all migrations
	@echo "$(BLUE)🔄 Resetting database...$(RESET)"
	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(MIGRATE_URL) down
	@echo "$(GREEN)✅ Database reset!$(RESET)\n"

.PHONY: migrate-force
migrate-force: ## Force set migration version (Usage: make migrate-force version=1)
	@echo "$(BLUE)⚠️ Forcing migration version...$(RESET)"
	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(MIGRATE_URL) force $(version)
	@echo "$(GREEN)✅ Version forced!$(RESET)\n"

# =========================
# 5. Docker Compose
# =========================
.PHONY: compose-up
compose-up: ## Create and start containers
	@echo "$(BLUE)🚀 Starting Docker containers...$(RESET)"
	@${COMPOSE_COMMAND} up -d 2>&1 > /dev/null
	@echo "$(GREEN)✅ Containers are up and running!$(RESET)\n"

.PHONY: compose-build
compose-build: ## Build, create and start containers
	@echo "$(BLUE)🚢 Building and starting Docker containers...$(RESET)"
	@${COMPOSE_COMMAND} up -d --build 2>&1 > /dev/null
	@echo "$(GREEN)✅ Containers are up and running!$(RESET)\n"
	@$(clean_dangling_images)

.PHONY: compose-services
compose-services: ## Build, create and start services containers
	@echo "$(BLUE)🚢 Building and starting Docker services containers...$(RESET)"
	@${COMPOSE_SERVICE} up -d --build 2>&1 > /dev/null
	@echo "$(GREEN)✅ Containers are up and running!$(RESET)\n"
	@$(clean_dangling_images)

.PHONY: compose-down
compose-down: ## Stop and remove containers and networks
	@echo "$(YELLOW)🛑 Stopping and removing containers...$(RESET)"
	@${COMPOSE_COMMAND} down 2>&1 > /dev/null
	@echo "$(GREEN)✅  Containers stopped.$(RESET)\n"

.PHONY: compose-clean
compose-clean: ## Clear dangling Docker images
	$(clean_dangling_images)

.PHONY: compose-remove
compose-remove: ## Stop and remove containers, networks and volumes
	@echo "$(RED)⚠️  WARNING: This will permanently delete all containers, networks, and VOLUMES!$(RESET)"
	@echo -n "$(RED)❌ All data will be lost. Are you sure? [y/N] $(RESET)" && read ans && [ $${ans:-N} = y ]
	@echo "$(YELLOW)\n🛑 Stopping and removing all Docker resources...$(RESET)"
	@${COMPOSE_COMMAND} down -v --remove-orphans 2>&1 > /dev/null
	@echo "$(GREEN)✅ Containers, networks, and volumes removed successfully.$(RESET)\n"

.PHONY: compose-exec
compose-exec: ## Access container bash
	@echo "$(BLUE)🔑 Accessing the container shell...$(RESET)"
	@${COMPOSE_COMMAND} exec -it ${DOCKER_SERVICE} bash

.PHONY: compose-log
compose-log: ## Show container logger
	@echo "$(BLUE)📜 Fetching container logs...$(RESET)"
	@${COMPOSE_COMMAND} logs -f ${DOCKER_SERVICE}

.PHONY: compose-top
compose-top: ## Display containers processes
	@echo "$(BLUE)🔍 Displaying container processes...$(RESET)"
	@${COMPOSE_COMMAND} top

.PHONY: compose-stats
compose-stats: ## Display containers stats
	@echo "$(CYAN)📊 Showing container statistics...$(RESET)"
	@${COMPOSE_COMMAND} stats
