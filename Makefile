.PHONY: build run clean install test fmt vet lint help

# Variables
BINARY_NAME=gob
INSTALL_PATH=/usr/local/bin
CMD_PATH=./cmd/gob

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BINARY_NAME)"

# Run the application
run: build
	@./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "Clean complete"

# Install the application
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo mv $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "Installation complete"

# Uninstall the application
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstall complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@golangci-lint run
	@echo "Lint complete"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	@GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME)-linux-arm64 $(CMD_PATH)
	@echo "Multi-platform build complete"

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .
	@echo "Docker build complete"

# Docker run
docker-run:
	@echo "Running Docker container..."
	@docker run --rm -it --privileged --net=host \
		-v /var/run/dbus:/var/run/dbus \
		$(BINARY_NAME):latest

# Show help
help:
	@echo "GOB - Gestor de Bluetooth para Linux"
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Build and run the application"
	@echo "  clean        - Remove build artifacts"
	@echo "  install      - Install to $(INSTALL_PATH)"
	@echo "  uninstall    - Uninstall from $(INSTALL_PATH)"
	@echo "  test         - Run tests"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  lint         - Run linter (requires golangci-lint)"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run in Docker container"
	@echo "  help         - Show this help message"
