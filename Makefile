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

# Show help
help:
	@echo "Available targets:"
	@echo "  build       - Build the binary"
	@echo "  run         - Run the application"
	@echo "  run-http    - Run in HTTP/SSE mode"
	@echo "  test        - Run tests"
	@echo "  clean       - Remove build artifacts"
	@echo "  deps        - Download dependencies"
	@echo "  install     - Install the binary"
	@echo "  update-deps - Update dependencies"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  build-all   - Build for multiple platforms"
	@echo "  help        - Show this help message"