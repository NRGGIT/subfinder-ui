.PHONY: build run test clean help

# Default target
.DEFAULT_GOAL := help

# Variables
APP_NAME := subfinder-service
DOCKER_IMAGE := $(APP_NAME)
DOCKER_TAG := latest

# Help target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  docker-build   - Build the Docker image"
	@echo "  docker-run     - Run the Docker container"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  help           - Show this help message"

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) ./cmd/server

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	./$(APP_NAME)

# Build the Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run the Docker container
docker-run: docker-build
	@echo "Running Docker container $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(APP_NAME)
	go clean
