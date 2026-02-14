.PHONY: build run clean install test fmt lint help

# Binary name
BINARY_NAME=mystart
BUILD_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application for current platform
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/mystart

build-linux: ## Build for Linux (amd64)
	@echo "Building $(BINARY_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux ./cmd/mystart

build-darwin: ## Build for macOS (amd64 and arm64)
	@echo "Building $(BINARY_NAME) for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/mystart
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/mystart

build-all: build-linux build-darwin ## Build for all platforms

run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

clean: ## Remove build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

install: ## Install the application
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install ./cmd/mystart

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

deps-upgrade: ## Upgrade dependencies
	@echo "Upgrading dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

vendor: ## Create vendor directory
	@echo "Vendoring dependencies..."
	$(GOMOD) vendor
