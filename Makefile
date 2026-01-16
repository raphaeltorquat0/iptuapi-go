.PHONY: all test lint build clean examples help

# Default target
all: lint test build

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Build
build:
	go build -v ./...

# Run security scan
security:
	@which govulncheck > /dev/null || (echo "Installing govulncheck..." && go install golang.org/x/vuln/cmd/govulncheck@latest)
	govulncheck ./...

# Tidy dependencies
tidy:
	go mod tidy
	go mod verify

# Clean build artifacts
clean:
	rm -f coverage.out coverage.html
	go clean -cache -testcache

# Run examples (requires IPTU_API_KEY)
examples:
	@if [ -z "$(IPTU_API_KEY)" ]; then echo "IPTU_API_KEY is required"; exit 1; fi
	go run ./examples/basic/

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Help
help:
	@echo "Available targets:"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make lint          - Run linter"
	@echo "  make build         - Build package"
	@echo "  make security      - Run security scan"
	@echo "  make tidy          - Tidy dependencies"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make fmt           - Format code"
	@echo "  make examples      - Run examples (requires IPTU_API_KEY)"
	@echo "  make help          - Show this help"
