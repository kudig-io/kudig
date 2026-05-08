# Makefile for kudig

# Build variables
BINARY_NAME=kudig
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build -buildvcs=false
GOTEST=$(GOCMD) test -buildvcs=false
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Directories
CMD_DIR=./cmd/kudig
PKG_DIR=./pkg/...
BUILD_DIR=./build

.PHONY: all build clean test fmt vet lint deps tidy help

# Default target
all: deps fmt vet lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)

build-darwin:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -cover $(PKG_DIR)

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out $(PKG_DIR)
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Upload coverage to Codecov (requires CODECOV_TOKEN env var)
upload-coverage:
	@echo "Uploading coverage to Codecov..."
	@if [ -z "$(CODECOV_TOKEN)" ]; then \
		echo "Error: CODECOV_TOKEN is not set"; \
		exit 1; \
	fi
	@which codecov > /dev/null || go install github.com/codecov/codecov-cli@latest
	codecov upload-process -t $(CODECOV_TOKEN) -f coverage.out

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

# Run linter with fix
lint-fix:
	@echo "Running golangci-lint with auto-fix..."
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run --fix ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Build kubectl plugin
build-kubectl-plugin:
	@echo "Building kubectl-kudig plugin..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/kubectl-kudig ./cmd/kubectl-kudig
	@echo "Built $(BUILD_DIR)/kubectl-kudig"

# Install kubectl plugin
install-kubectl-plugin: build-kubectl-plugin
	@echo "Installing kubectl-kudig plugin..."
	@cp $(BUILD_DIR)/kubectl-kudig $(GOPATH)/bin/kubectl-kudig 2>/dev/null || \
		cp $(BUILD_DIR)/kubectl-kudig /usr/local/bin/kubectl-kudig
	@echo "Installed kubectl-kudig. Run 'kubectl plugin list' to verify."

# Install binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME) 2>/dev/null || \
		cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# Run the application (for development)
run: build
	$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Run offline analysis on test data
run-test: build
	$(BUILD_DIR)/$(BINARY_NAME) offline reference/diagnose_k8s/diagnose_1765626516

# Run legacy mode on test data
run-legacy: build
	$(BUILD_DIR)/$(BINARY_NAME) legacy reference/diagnose_k8s/diagnose_1765626516

# Generate SBOM (requires syft)
sbom:
	@echo "Generating SBOM..."
	@which syft > /dev/null || (echo "Error: syft not found. Install from https://github.com/anchore/syft" && exit 1)
	syft scan dir:. -o spdx-json=sbom.spdx.json --source-name "kudig-go" --source-version "$(VERSION)"
	syft scan dir:. -o cyclonedx-json=sbom.cyclonedx.json --source-name "kudig-go" --source-version "$(VERSION)"
	@echo "SBOM generated: sbom.spdx.json, sbom.cyclonedx.json"

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build --build-arg VERSION=$(VERSION) --build-arg BUILD_TIME=$(BUILD_TIME) -t kudig/kudig:$(VERSION) -t kudig/kudig:latest .

# Docker push
docker-push: docker-build
	@echo "Pushing Docker image..."
	docker push kudig/kudig:$(VERSION)
	docker push kudig/kudig:latest

# Show help
help:
	@echo "Kudig Makefile Commands:"
	@echo ""
	@echo "  make                    - Build with deps, fmt, vet"
	@echo "  make build              - Build the binary"
	@echo "  make build-kubectl-plugin - Build kubectl plugin"
	@echo "  make build-all          - Build for all platforms"
	@echo "  make test               - Run tests"
	@echo "  make test-coverage      - Run tests with coverage"
	@echo "  make fmt                - Format code"
	@echo "  make vet                - Run go vet"
	@echo "  make lint               - Run linter"
	@echo "  make deps               - Download dependencies"
	@echo "  make tidy               - Tidy dependencies"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make install            - Install binary"
	@echo "  make install-kubectl-plugin - Install kubectl plugin"
	@echo "  make run-test           - Run on test data"
	@echo "  make run-legacy         - Run legacy mode on test data"
	@echo ""
