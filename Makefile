.PHONY: test test-internal test-all coverage build-frontend clean install-deps

# Test only internal packages (no frontend required)
test-internal:
	go test -v ./internal/...

# Test internal packages with coverage
test-coverage:
	go test -cover -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out

# Build frontend and run all tests
test-all: build-frontend
	go test -v ./...

# Build the frontend
build-frontend:
	cd frontend && npm install && npm run build

# Install Go dependencies
install-deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -rf build/
	rm -rf frontend/dist/
	rm -rf frontend/node_modules/
	rm -f coverage.out

# Run the application in dev mode (requires Wails CLI)
dev:
	wails dev

# Build the application (requires Wails CLI)
build: build-frontend
	wails build

# Build CLI only
build-cli:
	go build -o vsynx.exe .

# Run linter
lint:
	go vet ./...
	cd frontend && npm run lint

# Format code
fmt:
	go fmt ./...
	cd frontend && npm run format

# Quick test (internal packages only)
test: test-internal

# Help target
help:
	@echo "Available targets:"
	@echo "  test-internal   - Test internal packages only (fast, no frontend needed)"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-all        - Build frontend and test everything"
	@echo "  build-frontend  - Build React frontend"
	@echo "  install-deps    - Install Go dependencies"
	@echo "  clean           - Clean build artifacts"
	@echo "  dev             - Run in development mode (requires Wails)"
	@echo "  build           - Build production application"
	@echo "  build-cli       - Build CLI only"
	@echo "  lint            - Run linters"
	@echo "  fmt             - Format code"
