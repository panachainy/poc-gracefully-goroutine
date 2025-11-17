# Makefile for poc-gracefully-goroutine

# Variables
BINARY_NAME=poc-gracefully-goroutine
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*')
COVERAGE_DIR=coverage

# Run the application
.PHONY: run
r run:
	go run .
	# go run main.go & sleep 2; kill -INT $!; wait $!

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
