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

# Frontend unit tests (Vitest)
test-frontend:
	cd frontend && npm run test -- --run

# Frontend unit tests with coverage
test-frontend-coverage:
	cd frontend && npm run test:coverage

# E2E tests (Playwright)
test-e2e:
	cd frontend && npm run test:e2e

# E2E tests with UI
test-e2e-ui:
	cd frontend && npm run test:e2e:ui

# Install Playwright browsers
install-playwright:
	cd frontend && npx playwright install

# Run all tests (Go + Frontend)
test-complete: test-internal test-frontend

# Help target
help:
	@echo "Available targets:"
	@echo "  test-internal       - Test internal packages only (fast, no frontend needed)"
	@echo "  test-coverage       - Run tests with coverage report"
	@echo "  test-all            - Build frontend and test everything"
	@echo "  test-frontend       - Run frontend unit tests (Vitest, single run)"
	@echo "  test-frontend-coverage - Run frontend tests with coverage"
	@echo "  test-e2e            - Run E2E tests (Playwright)"
	@echo "  test-e2e-ui         - Run E2E tests with UI"
	@echo "  test-complete       - Run all tests (Go + Frontend)"
	@echo "  install-playwright  - Install Playwright browsers"
	@echo "  build-frontend      - Build React frontend"
	@echo "  install-deps        - Install Go dependencies"
	@echo "  clean               - Clean build artifacts"
	@echo "  dev                 - Run in development mode (requires Wails)"
	@echo "  build               - Build production application"
	@echo "  build-cli           - Build CLI only"
	@echo "  lint                - Run linters"
	@echo "  fmt                 - Format code"
