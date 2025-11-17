# Makefile for poc-gracefully-goroutine

# Variables
BINARY_NAME=poc-gracefully-goroutine
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*')
COVERAGE_DIR=coverage

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	go build -o $(BINARY_NAME) .

# Run tests
.PHONY: test
t test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

# Clean build artifacts
.PHONY: clean
c clean:
	rm -f $(BINARY_NAME)
	rm -rf $(COVERAGE_DIR)/*.out $(COVERAGE_DIR)/*.html

# Run the application
.PHONY: run
r run:
	go run .

# Format code
.PHONY: fmt
f fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
.PHONY: lint
l lint:
	golangci-lint run

# Tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy
