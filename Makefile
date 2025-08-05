.PHONY: build run test clean install deps

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

# Build the project
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v .

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run main.go

# Run in HTTP/SSE mode
run-http:
	@echo "Running $(BINARY_NAME) in HTTP/SSE mode..."
	$(GOCMD) run main.go -http :3001

# Test the application
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 -v
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 -v
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 -v
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe -v

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -f deployment/docker/Dockerfile -t $(BINARY_NAME):latest .

docker-run:
	@echo "Running in Docker..."
	docker-compose -f deployment/local/docker-compose.yaml up

docker-run-dev:
	@echo "Running in Docker development mode..."
	docker-compose -f deployment/local/docker-compose.yaml --profile dev up

docker-stop:
	@echo "Stopping Docker containers..."
	docker-compose -f deployment/local/docker-compose.yaml down

docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose -f deployment/local/docker-compose.yaml down -v
	docker rmi $(BINARY_NAME):latest || true
	docker rmi $(BINARY_NAME):dev || true

# Release commands
release:
	@echo "Creating release..."
	@if [ -z "$(VERSION)" ]; then \
		echo "VERSION not specified. Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	./scripts/release.sh $(VERSION)

# Health check
health:
	@echo "Checking health..."
	curl -f http://localhost:$(PORT:-3001)/health || exit 1

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  run           - Run the application"
	@echo "  run-http      - Run in HTTP/SSE mode"
	@echo "  test          - Run tests"
	@echo "  clean         - Remove build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  install       - Install the binary"
	@echo "  update-deps   - Update dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  build-all     - Build for multiple platforms"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run in Docker"
	@echo "  docker-run-dev - Run in Docker development mode"
	@echo "  docker-stop   - Stop Docker containers"
	@echo "  docker-clean  - Clean Docker resources"
	@echo ""
	@echo "Release targets:"
	@echo "  release       - Create a release (requires VERSION=v1.0.0)"
	@echo "  health        - Check service health"
	@echo ""
	@echo "  help          - Show this help message"