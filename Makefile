# =========================
# Configs
# =========================

# Docker
DOCKER_BASE     ?= $(or ${base}, built)
DOCKER_SERVICE  ?= $(or ${service}, backend)
COMPOSE_COMMAND := BASE=${DOCKER_BASE} docker compose --env-file config/.env -f build/docker/compose.yml

# Go build
GO          := go
GOOS        := linux
GOARCH      := amd64
BUILD_FLAGS := -ldflags "-w -s"
GOBUILD     := CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build $(BUILD_FLAGS)

# Tool versions
SWAG_VERSION         := v1.16.3
GOCOV_VERSION        := v1.2.1
GOCOV_HTML_VERSION   := v1.4.0

# Colors
GREEN  := \033[1;32m
YELLOW := \033[1;33m
BLUE   := \033[1;34m
CYAN   := \033[1;36m
MAGENTA:= \033[1;35m
RED    := \033[1;31m
RESET  := \033[0m

# =========================
# Helpers
# =========================

define clean_dangling_images
	@echo "$(YELLOW)🧹 Cleaning up dangling Docker images...$(RESET)"
	@if test -n "$$(docker images -f "dangling=true" -q)"; then \
		docker rmi $$(docker images -f "dangling=true" -q); \
	fi > /dev/null
endef

# =========================
# Dependências
# =========================

.PHONY: check-docker check-swag doctor

check-docker:
	@command -v docker >/dev/null || (echo "❌ Docker não encontrado!" && exit 1)
	@docker compose version >/dev/null || (echo "❌ Docker Compose não encontrado!" && exit 1)

check-swag:
	@command -v swag >/dev/null 2>&1 || { \
		echo "❌ swag não encontrado. Instale com:"; \
		echo "   go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)"; \
		exit 1; \
	}

doctor: ## Check required tools
	@for cmd in go docker swag; do \
	  if ! command -v $$cmd >/dev/null 2>&1; then \
	    echo "❌ Missing: $$cmd"; \
	  else \
	    echo "✅ Found: $$cmd"; \
	  fi; \
	done

# =========================
# Ajuda
# =========================

.PHONY: all help
all: help
help: ## Display available commands and their descriptions
	@echo "$(CYAN)Usage:$(RESET)"
	@echo "  make [COMMAND]\n"
	@echo "$(CYAN)Example:$(RESET)"
	@echo "  make build\n"
	@echo "$(CYAN)Commands:$(RESET)\n"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(CYAN)%-30s$(RESET) %s\n", $$1, $$2}'

# =========================
# Go commands
# =========================

.PHONY: init test run build swag format tidy lint audit benchmark ci

init: ## Create environment file
	@echo "$(BLUE)⚙️  Initializing environment setup...$(RESET)"
	@chmod +x config/env.sh && config/env.sh && mv .env config/
	@echo "$(GREEN)✅ Environment file successfully created!$(RESET)\n"

test: ## Run tests and generate coverage report
	@echo "$(BLUE)🔍 Running tests...$(RESET)"
	@-$(GO) install github.com/axw/gocov/gocov@$(GOCOV_VERSION)
	@-$(GO) install github.com/matm/gocov-html/cmd/gocov-html@$(GOCOV_HTML_VERSION)
	@-$(GO) clean -testcache
	@-$(GO) test ./... -coverprofile cover.out
	@-$(GO) tool cover -html=cover.out
	@gocov convert cover.out | gocov-html -t kit > report.html
	@echo "$(GREEN)✅ Tests completed!$(RESET)\n"

run: ## Run application from source code
	@echo "$(CYAN)▶️  Running the application...$(RESET)"
	@$(GO) run cmd/backend/main.go
	@echo "$(GREEN)✅ Application stopped.$(RESET)"

build: ## Build all applications from source code
	@echo "$(BLUE)🚀 Building application...$(RESET)"
	@${GOBUILD} -o bin/backend cmd/backend/main.go
	@echo "$(GREEN)✅ Build completed successfully!$(RESET)\n"

swag: check-swag ## Update swagger files
	@echo "$(BLUE)📄 Updating Swagger API documentation...$(RESET)"
	@swag init -g cmd/backend/main.go --parseDependency -o ./docs
	@echo "$(GREEN)✅ Swagger files updated successfully.$(RESET)\n"

format: ## Fix code format issues
	@echo "$(YELLOW)📝 Formatting code...$(RESET)"
	@$(GO) run mvdan.cc/gofumpt@latest -w -l . 2>&1 > /dev/null
	@echo "$(GREEN)✅ Code formatting complete!$(RESET)\n"

tidy: ## Clean and tidy dependencies
	@echo "$(YELLOW)🔧 Cleaning and tidying Go dependencies...$(RESET)"
	@$(GO) mod tidy -v 2>&1 > /dev/null
	@echo "$(GREEN)✅ Dependencies tidied successfully.$(RESET)\n"

audit: ## Conduct quality checks
	@echo "$(YELLOW)🔎 Running code audit...$(RESET)"
	@$(GO) mod verify 2>&1 > /dev/null
	@$(GO) vet ./... 2>&1 > /dev/null
	@$(GO) run golang.org/x/vuln/cmd/govulncheck@latest -show verbose ./... 2>&1 > /dev/null
	@echo "$(GREEN)✅ Code audit finished!$(RESET)\n"

benchmark: ## Benchmark code performance
	@echo "$(MAGENTA)⚡ Running benchmarks...$(RESET)"
	@$(GO) test ./... -benchmem -bench=. -run=^Benchmark_$ 2>&1 > /dev/null
	@echo "$(GREEN)✅ Benchmark completed!$(RESET)\n"

ci: tidy lint test audit ## Run all quality checks

# =========================
# Docker compose commands
# =========================

.PHONY: compose-up compose-build compose-down compose-clean compose-remove compose-exec compose-log compose-top compose-stats

compose-up: check-docker ## Create and start containers
	@echo "$(BLUE)🚀 Starting Docker containers...$(RESET)"
	@${COMPOSE_COMMAND} up -d 2>&1 > /dev/null
	@echo "$(GREEN)✅ Containers are up and running!$(RESET)\n"

compose-build: check-docker ## Build, create and start containers
	@echo "$(BLUE)🚢 Building and starting Docker containers...$(RESET)"
	@${COMPOSE_COMMAND} up -d --build 2>&1 > /dev/null
	@echo "$(GREEN)✅ Containers are up and running!$(RESET)\n"
	@$(clean_dangling_images)

compose-down: check-docker ## Stop and remove containers and networks
	@echo "$(YELLOW)🛑 Stopping and removing containers...$(RESET)"
	@${COMPOSE_COMMAND} down 2>&1 > /dev/null
	@echo "$(GREEN)✅ Containers stopped.$(RESET)\n"

compose-clean: check-docker ## Clear dangling Docker images
	$(clean_dangling_images)

compose-remove: check-docker ## Stop and remove containers, networks and volumes
	@echo "$(RED)⚠️  WARNING: This will permanently delete all containers, networks, and VOLUMES!$(RESET)"
	@echo -n "$(RED)❌ All data will be lost. Are you sure? [y/N] $(RESET)" && read ans && [ $${ans:-N} = y ]
	@echo "$(YELLOW)\n🛑 Stopping and removing all Docker resources...$(RESET)"
	@${COMPOSE_COMMAND} down -v --remove-orphans 2>&1 > /dev/null
	@echo "$(GREEN)✅ Containers, networks, and volumes removed successfully.$(RESET)\n"

compose-exec: check-docker ## Access container bash
	@echo "$(BLUE)🔑 Accessing the container shell...$(RESET)"
	@${COMPOSE_COMMAND} exec -it ${DOCKER_SERVICE} bash

compose-log: check-docker ## Show container logger
	@echo "$(BLUE)📜 Fetching container logs...$(RESET)"
	@${COMPOSE_COMMAND} logs -f ${DOCKER_SERVICE}

compose-top: check-docker ## Display containers processes
	@echo "$(BLUE)🔍 Displaying container processes...$(RESET)"
	@${COMPOSE_COMMAND} top

compose-stats: check-docker ## Display containers stats
	@echo "$(CYAN)📊 Showing container statistics...$(RESET)"
	@${COMPOSE_COMMAND} stats
