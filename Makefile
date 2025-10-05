.PHONY: help test test-race test-coverage fmt vet lint tidy build clean

# Default PostgreSQL connection string for tests
export POSTGRES_CONNECTION_STRING ?= postgres://postgres:password@localhost:5432/postgres?sslmode=disable

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

test: ## Run tests
	go test -v ./...

test-race: ## Run tests with race detector
	go test -race -v ./...

test-coverage: ## Run tests with coverage report
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run staticcheck
	@command -v staticcheck >/dev/null 2>&1 || { echo "Installing staticcheck..."; go install honnef.co/go/tools/cmd/staticcheck@latest; }
	staticcheck ./...

tidy: ## Tidy dependencies
	go mod tidy

build: ## Build the package
	go build ./...

clean: ## Clean generated files
	rm -f coverage.out coverage.html
	go clean

security: ## Run security checks
	@command -v gosec >/dev/null 2>&1 || { echo "Installing gosec..."; go install github.com/securego/gosec/v2/cmd/gosec@latest; }
	gosec ./...
	@command -v govulncheck >/dev/null 2>&1 || { echo "Installing govulncheck..."; go install golang.org/x/vuln/cmd/govulncheck@latest; }
	govulncheck ./...

ci: fmt vet lint test-race test-coverage ## Run all CI checks locally

all: tidy fmt vet lint test-coverage ## Run all checks and tests
