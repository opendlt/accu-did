.PHONY: build run test lint fmt clean help

# Build settings
BINARY_NAME=registrar
BUILD_DIR=bin
MAIN_PATH=./cmd/registrar

# Go settings
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BUILD_DIR)

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the registrar binary
	@echo "Building registrar..."
	@mkdir -p $(GOBIN)
	@go build -o $(GOBIN)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary built at $(GOBIN)/$(BINARY_NAME)"

run: ## Run the registrar (port 8081)
	@echo "Starting registrar on :8081..."
	@go run $(MAIN_PATH)

test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@gofumpt -l -w .
	@go mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

install-tools: ## Install development tools
	@echo "Installing tools..."
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

dev: ## Development setup
	@make install-tools
	@make deps

all: clean fmt lint test build ## Run all checks and build