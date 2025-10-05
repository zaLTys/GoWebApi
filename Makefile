.PHONY: help build run test test-unit test-integration test-coverage clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building application..."
	@go build -o books-api .

run: ## Run the application
	@echo "Starting application..."
	@go run .

test: ## Run all tests
	@echo "Running all tests..."
	@go test -v ./...

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@go test -v ./tests/controllers/... ./tests/services/... ./tests/repositories/... ./tests/models/...

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	@go test -v ./tests/integration/...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	@go test -race -v ./...

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run || echo "Install golangci-lint: https://golangci-lint.run/usage/install/"

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@~/go/bin/swag init || swag init

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f books-api
	@rm -f coverage.out coverage.html
	@rm -rf tmp/

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t books-api:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 books-api:latest

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

dev: ## Run in development mode with hot reload
	@echo "Starting development server..."
	@air || echo "Install air: go install github.com/air-verse/air@latest"
