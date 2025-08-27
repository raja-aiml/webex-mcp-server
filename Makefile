.PHONY: help build run docker deps dev clean install release health test fmt lint

# Binary name
BINARY_NAME=webex-mcp-server

# Build directory
BUILD_DIR=build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Default target
.DEFAULT_GOAL := help

# Check if we have a subcommand
ifdef MAKECMDGOALS
CMD = $(word 1, $(MAKECMDGOALS))
SUBCMD = $(word 2, $(MAKECMDGOALS))
endif

#=============================================================================
# Main Commands
#=============================================================================

## help: Show this help message
help:
	@echo "Webex MCP Server - Makefile Commands"
	@echo ""
	@echo "Usage: make <command> [<subcommand>]"
	@echo ""
	@echo "BUILD COMMANDS:"
	@echo "  build          Build the binary for current platform"
	@echo "  build all      Build binaries for all platforms"
	@echo ""
	@echo "RUN COMMANDS:"
	@echo "  run            Run locally (stdio mode, matching mcp.json local config)"
	@echo "  run http       Run the application in HTTP/SSE mode"
	@echo ""
	@echo "MCP COMMANDS:"
	@echo "  mcp-local      Run matching mcp.json webex-mcp-local-stdio config"
	@echo "  mcp-docker     Run matching mcp.json webex-mcp-docker config"
	@echo ""
	@echo "DOCKER COMMANDS:"
	@echo "  docker build     Build Docker image"
	@echo "  docker run       Run matching mcp.json docker config (stdio mode)"
	@echo "  docker run-http  Run in Docker HTTP mode (port 8084)"
	@echo "  docker run-dev   Run in Docker development mode"
	@echo "  docker stop      Stop Docker containers"
	@echo "  docker clean     Clean Docker resources"
	@echo ""
	@echo "DEPENDENCY COMMANDS:"
	@echo "  deps           Install dependencies"
	@echo "  deps update    Update dependencies"
	@echo ""
	@echo "DEVELOPMENT COMMANDS:"
	@echo "  dev fmt        Format Go code"
	@echo "  dev lint       Lint code (requires golangci-lint)"
	@echo "  dev test       Run tests"
	@echo "  test-coverage  Run tests with coverage report"
	@echo ""
	@echo "OTHER COMMANDS:"
	@echo "  clean          Remove build artifacts"
	@echo "  install        Install binary to GOPATH/bin"
	@echo "  release        Create a release (VERSION=v1.0.0)"
	@echo "  health         Check service health"
	@echo ""
	@echo "QUICK COMMANDS (shortcuts):"
	@echo "  make           Show this help"
	@echo "  make build     Build for current platform"
	@echo "  make run       Run locally (matching mcp.json local config)"
	@echo "  make test      Run tests (alias for dev test)"

#=============================================================================
# Build Commands
#=============================================================================

## build: Build command handler
build:
ifeq ($(SUBCMD),all)
	@$(MAKE) -s build-all
else ifeq ($(SUBCMD),)
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v .
else
	@echo "Unknown build subcommand: $(SUBCMD)"
	@echo "Available: build, build all"
endif

## build-all: Build binaries for all platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 -v
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 -v
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 -v
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe -v

# Support for "make build all"
all:
	@# This target exists only to support "make build all" syntax

#=============================================================================
# Run Commands
#=============================================================================

## run: Run command handler
run:
ifeq ($(SUBCMD),http)
	@$(MAKE) -s run-http-impl
else ifeq ($(SUBCMD),)
	@echo "Running $(BINARY_NAME) locally (matching mcp.json local config)..."
	$(GOCMD) run . -env .env -all-tools
else
	@echo "Unknown run subcommand: $(SUBCMD)"
	@echo "Available: run, run http"
endif

## mcp-local: Run matching mcp.json webex-mcp-local-stdio config
mcp-local:
	@echo "Running MCP server locally (stdio mode, -all-tools)..."
	$(GOCMD) run . -env .env -all-tools

## mcp-docker: Run matching mcp.json webex-mcp-docker config
mcp-docker: docker-build
	@echo "Running MCP server in Docker (stdio mode, -all-tools)..."
	docker run --rm -i --env-file .env webex-mcp-server:latest

# Support for "make run http"
http:
	@# This target exists only to support "make run http" syntax

#=============================================================================
# Docker Commands
#=============================================================================

## docker: Docker command handler
docker:
ifeq ($(SUBCMD),build)
	@$(MAKE) -s docker-build MAKECMDGOALS=docker-build
else ifeq ($(SUBCMD),run)
	@$(MAKE) -s docker-run MAKECMDGOALS=docker-run
else ifeq ($(SUBCMD),run-http)
	@$(MAKE) -s docker-run-http MAKECMDGOALS=docker-run-http
else ifeq ($(SUBCMD),run-dev)
	@$(MAKE) -s docker-run-dev MAKECMDGOALS=docker-run-dev
else ifeq ($(SUBCMD),stop)
	@$(MAKE) -s docker-stop MAKECMDGOALS=docker-stop
else ifeq ($(SUBCMD),clean)
	@$(MAKE) -s docker-clean MAKECMDGOALS=docker-clean
else ifeq ($(SUBCMD),)
	@echo "Docker commands:"
	@echo "  make docker build     - Build Docker image"
	@echo "  make docker run       - Run matching mcp.json docker config (stdio)"
	@echo "  make docker run-http  - Run in Docker HTTP mode"
	@echo "  make docker run-dev   - Run in Docker development mode"
	@echo "  make docker stop      - Stop Docker containers"
	@echo "  make docker clean     - Clean Docker resources"
else
	@echo "Unknown docker subcommand: $(SUBCMD)"
	@echo "Available: build, run, run-http, run-dev, stop, clean"
endif

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -f Dockerfile -t $(BINARY_NAME):latest .

## docker-run: Run matching mcp.json webex-mcp-docker config
docker-run: docker-build
	@echo "Running in Docker (stdio mode, matching mcp.json config)..."
	docker run --rm -i --env-file .env $(BINARY_NAME):latest

## docker-run-http: Run application in Docker HTTP mode
docker-run-http:
	@echo "Running in Docker HTTP mode..."
	docker-compose --profile http up

## docker-run-dev: Run in Docker development mode
docker-run-dev:
	@echo "Running in Docker development mode..."
	docker-compose --profile dev up

## docker-stop: Stop Docker containers
docker-stop:
	@echo "Stopping Docker containers..."
	docker-compose down

## docker-clean: Clean Docker resources
docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker rmi $(BINARY_NAME):latest || true
	docker rmi $(BINARY_NAME):dev || true

# Support for docker subcommands
run-http:
ifeq ($(CMD),docker)
	@# This target exists only to support "make docker run-http" syntax
else
	@$(MAKE) -s run-http-impl
endif

run-http-impl:
	@echo "Running $(BINARY_NAME) in HTTP/SSE mode..."
	$(GOCMD) run . -http :3001 -all-tools

run-dev:
	@# This target exists only to support "make docker run-dev" syntax

#=============================================================================
# Dependency Commands
#=============================================================================

## deps: Dependency command handler
deps:
ifeq ($(SUBCMD),update)
	@$(MAKE) -s deps-update
else ifeq ($(SUBCMD),)
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
else
	@echo "Unknown deps subcommand: $(SUBCMD)"
	@echo "Available: deps, deps update"
endif

## deps-update: Update dependencies
deps-update:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Support for "make deps update"
update:
	@# This target exists only to support "make deps update" syntax

#=============================================================================
# Development Commands
#=============================================================================

## dev: Development command handler
dev:
ifeq ($(SUBCMD),fmt)
	@$(MAKE) -s dev-fmt
else ifeq ($(SUBCMD),lint)
	@$(MAKE) -s dev-lint
else ifeq ($(SUBCMD),test)
	@$(MAKE) -s dev-test
else ifeq ($(SUBCMD),all)
	@$(MAKE) -s dev-all
else ifeq ($(SUBCMD),)
	@echo "Development commands:"
	@echo "  make dev fmt    - Format Go code"
	@echo "  make dev lint   - Lint code"
	@echo "  make dev test   - Run tests"
	@echo "  make dev all    - Run fmt, lint, and test"
else
	@echo "Unknown dev subcommand: $(SUBCMD)"
	@echo "Available: fmt, lint, test, all"
endif

## dev-fmt: Format Go code
dev-fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

## dev-lint: Lint code (requires golangci-lint)
dev-lint:
	@echo "Linting code..."
	golangci-lint run

## dev-test: Run tests
dev-test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'

## dev-all: Run format, lint, and test
dev-all: dev-fmt dev-lint dev-test
	@echo "Development checks complete!"

## check-token: Check Webex token and get user details
check-token:
	@echo "Checking Webex token and getting user details..."
	@if [ -n "$(ENV_PATH)" ]; then \
		echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"name":"test-client","version":"1.0.0","capabilities":{}}}{"jsonrpc":"2.0","id":2,"method":"mcp.tools","params":{"name":"get_me"}}' | ./build/webex-mcp-server -env $(ENV_PATH); \
	else \
		echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"name":"test-client","version":"1.0.0","capabilities":{}}}{"jsonrpc":"2.0","id":2,"method":"mcp.tools","params":{"name":"get_me"}}' | ./build/webex-mcp-server; \
	fi

## security-scan: Run security scans
security-scan:
	@echo "Running security scans..."
	@if command -v gitleaks >/dev/null 2>&1; then \
		echo "🔍 Running gitleaks scan..."; \
		gitleaks detect --config .gitleaks.toml --verbose --no-git; \
	else \
		echo "⚠️  gitleaks not installed. Install with: brew install gitleaks"; \
	fi
	@echo "🔍 Running manual secret checks..."
	@echo "Checking for .env files that would be committed to git..."
	@if git ls-files | grep -E "\.env$$|\.env\.local$$|\.env\.production$$" >/dev/null 2>&1; then \
		echo "❌ .env files found in git index (these would be committed)!"; \
		git ls-files | grep -E "\.env$$|\.env\.local$$|\.env\.production$$"; \
		exit 1; \
	else \
		echo "✅ No .env files in git index"; \
	fi
	@echo "Checking for hardcoded secrets in tracked files..."
	@if git ls-files | xargs grep -l "WEBEX_PUBLIC_WORKSPACE_API_KEY.*=.*[a-zA-Z0-9]{30,}" 2>/dev/null; then \
		echo "❌ Potential hardcoded Webex API key found in tracked files!"; \
		git ls-files | xargs grep -l "WEBEX_PUBLIC_WORKSPACE_API_KEY.*=.*[a-zA-Z0-9]{30,}" 2>/dev/null; \
		exit 1; \
	fi
	@if git ls-files | xargs grep -l "sk-[a-zA-Z0-9-_]{40,}" 2>/dev/null; then \
		echo "❌ Potential OpenAI API key found in tracked files!"; \
		git ls-files | xargs grep -l "sk-[a-zA-Z0-9-_]{40,}" 2>/dev/null; \
		exit 1; \
	fi
	@echo "✅ No hardcoded secrets found in tracked files"
	@echo "✅ Security scan completed - repository is secure!"

## setup-hooks: Set up git hooks for security
setup-hooks:
	@echo "Setting up git hooks..."
	@cp scripts/pre-commit.sh .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "✅ Pre-commit hook installed"

# Support for dev subcommands
fmt:
ifeq ($(CMD),dev)
	@:
else
	@$(MAKE) -s dev-fmt
endif

lint:
ifeq ($(CMD),dev)
	@:
else
	@$(MAKE) -s dev-lint
endif

test:
ifeq ($(CMD),dev)
	@:
else
	@$(MAKE) -s dev-test
endif

#=============================================================================
# Other Commands
#=============================================================================

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf releases/

## install: Install binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@if [ -z "$(GOPATH)" ]; then \
		echo "GOPATH is not set"; \
		exit 1; \
	fi
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

## release: Create a release (VERSION=v1.0.0)
release:
	@echo "Creating release..."
	@if [ -z "$(VERSION)" ]; then \
		echo "VERSION not specified. Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Release process would need to be implemented directly in Makefile or GitHub Actions"

## health: Check service health
health:
	@echo "Checking health..."
	@curl -f http://localhost:$${PORT:-3001}/health || exit 1

#=============================================================================
# Grouped Command Shortcuts
#=============================================================================

## docker-rebuild: Clean and rebuild Docker image
docker-rebuild: docker-clean docker-build
	@echo "Docker image rebuilt!"

## full-build: Clean, install deps, and build all platforms
full-build: clean deps build-all
	@echo "Full build complete!"