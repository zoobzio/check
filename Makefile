.PHONY: test test-unit test-integration test-bench lint lint-fix coverage clean check ci install-tools install-hooks help

.DEFAULT_GOAL := help

# Run all tests with race detector
test: ## Run all tests with race detector
	go test -race -v ./...

# Run unit tests only (short mode)
test-unit: ## Run unit tests only (short mode)
	go test -short -v ./...

# Run integration tests
test-integration: ## Run integration tests
	go test -v -run Integration ./...

# Run benchmarks
test-bench: ## Run benchmarks
	go test -bench=. -benchmem ./...

# Run linters
lint: ## Run linters
	golangci-lint run

# Run linters with auto-fix
lint-fix: ## Run linters with auto-fix
	golangci-lint run --fix

# Generate coverage report
coverage: ## Generate coverage report
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Remove generated files
clean: ## Remove generated files
	rm -f coverage.out coverage.html
	rm -rf dist/

# Quick validation (test + lint)
check: test lint ## Quick validation (test + lint)

# Full CI simulation
ci: clean check coverage ## Full CI simulation

# Install dev tools
install-tools: ## Install dev tools
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.7.2

# Install git hooks
install-hooks: ## Install git hooks
	@echo "Installing pre-commit hook..."
	@echo '#!/bin/sh\nmake check' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Pre-commit hook installed"

# Display available commands
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
