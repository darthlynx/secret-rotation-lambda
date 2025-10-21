.PHONY: build clean test coverage deploy local-test

# Variables
BINARY_NAME=bootstrap # required name for AWS Lambda custom runtime
LAMBDA_ZIP=lambda.zip

BUILD_DIR ?= bin
CMD_DIR ?= ./cmd/lambda

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo 'v0.0.0')
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

PLATFORMS ?= linux/amd64 darwin/amd64 linux/arm64

.DEFAULT_GOAL := help

# Build for Lambda (local platform by default)
build: deps
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Build for Linux AMD64
build-linux:
	@$(MAKE) build GOOS=linux GOARCH=amd64

# Build for multiple platforms
build-all:
	@mkdir -p $(BUILD_DIR)
	@set -e; \
	for p in $(PLATFORMS); do \
		os=$$(echo $$p | cut -d/ -f1); arch=$$(echo $$p | cut -d/ -f2); \
		out=$(BUILD_DIR)/$(BINARY_NAME)-$$os-$$arch; \
		echo "building $$out ..."; \
		GOOS=$$os GOARCH=$$arch go build -ldflags "$(LDFLAGS)" -o $$out $(CMD_DIR); \
	done

# Create deployment package
package: build
	@echo "Creating deployment package..."
	cd $(BUILD_DIR) && \
	zip $(LAMBDA_ZIP) $(BINARY_NAME) && \
	cd - && \
	@echo "Package created: $(LAMBDA_ZIP)"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f $(BUILD_DIR)/$(LAMBDA_ZIP)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

# Lint code
lint:
	@echo "Running linter..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint not found; install from https://golangci-lint.run/"; exit 1; \
	fi
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@if ! command -v goimports >/dev/null 2>&1; then \
		echo "goimports not found; installing..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	goimports -w .

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Local testing with sample event
local-test: build
	@echo "Testing locally..."
	@echo '{"secret_arn":"arn:aws:secretsmanager:us-east-1:123456789012:secret:test","secret_type":"plaintext","generator_options":{"length":16,"include_digits":true,"include_uppercase":true,"include_lowercase":true,"include_special_chars":false}}' | \
	go run cmd/lambda/main.go

# Deploy to AWS (requires AWS CLI and appropriate permissions)
# deploy: package
# 	@echo "Deploying to AWS Lambda..."
# 	@read -p "Enter Lambda function name: " func_name; \
# 	aws lambda update-function-code \
# 		--function-name $$func_name \
# 		--zip-file fileb://$(LAMBDA_ZIP)
# 	@echo "Deployment complete"

# Create new Lambda function (first time setup)
# create-function: package
# 	@echo "Creating new Lambda function..."
# 	@read -p "Enter Lambda function name: " func_name; \
# 	read -p "Enter IAM role ARN: " role_arn; \
# 	aws lambda create-function \
# 		--function-name $$func_name \
# 		--runtime provided.al2 \
# 		--handler bootstrap \
# 		--role $$role_arn \
# 		--zip-file fileb://$(LAMBDA_ZIP) \
# 		--timeout 30 \
# 		--memory-size 256
# 	@echo "Function created"

# Show Help
help:
	@echo "Available targets:"
	@echo "  build           - Build the Lambda function binary"
	@echo "  build-linux     - Build the Lambda function binary for Linux AMD64"
	@echo "  build-all       - Build the Lambda function binary for all platforms"
	@echo "  package         - Create deployment ZIP package"
	@echo "  test            - Run all tests"
	@echo "  coverage        - Run tests with coverage report"
	@echo "  test-race       - Run tests with race detector"
	@echo "  clean           - Remove build artifacts"
	@echo "  lint            - Run linter"
	@echo "  fmt             - Format code"
	@echo "  deps            - Download and tidy dependencies"
	@echo "  local-test      - Test locally with sample event"
# 	@echo "  deploy          - Deploy to existing Lambda function"
# 	@echo "  create-function - Create new Lambda function"
