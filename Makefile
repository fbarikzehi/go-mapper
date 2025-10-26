.PHONY: help test test-coverage test-race lint fmt vet bench clean install-tools

# Variables
GOBIN ?= $(shell go env GOPATH)/bin
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOFMT = gofmt
GOLINT = golangci-lint

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $1, $2}' $(MAKEFILE_LIST)

test: ## Run tests
	$(GOTEST) -v -race ./...

test-coverage: ## Run tests with coverage
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detector
	$(GOTEST) -race -short ./...

bench: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem -run=^$ ./...

bench-compare: ## Run benchmarks and save to file
	$(GOTEST) -bench=. -benchmem -run=^$ ./... | tee bench.txt

lint: ## Run linters
	$(GOLINT) run --timeout=5m

fmt: ## Format code
	$(GOFMT) -s -w .
	$(GOCMD) fmt ./...

vet: ## Run go vet
	$(GOCMD) vet ./...

tidy: ## Tidy go modules
	$(GOMOD) tidy
	$(GOMOD) verify

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f coverage.out coverage.html bench.txt

install-tools: ## Install development tools
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

pre-commit: fmt vet lint test ## Run pre-commit checks

ci: tidy fmt vet lint test-coverage ## Run CI pipeline locally

.DEFAULT_GOAL := help
