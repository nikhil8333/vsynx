# Testing Guide

This document explains how to run and write tests for Vsynx Manager.

## Running Tests

### Test Internal Packages Only

The internal packages (core business logic) can be tested without building the frontend:

```bash
# Run all internal package tests
go test ./internal/...

# Run with verbose output
go test -v ./internal/...

# Run with coverage
go test -cover ./internal/...

# Run specific package
go test ./internal/models
go test ./internal/validation
go test ./internal/marketplace
go test ./internal/openvsx
```

### Test All Packages (Requires Frontend Build)

To test the entire application including GUI components, you must first build the frontend:

```bash
# Install frontend dependencies
cd frontend
npm install

# Build frontend
npm run build
cd ..

# Now run all tests
go test ./...
```

Alternatively, use the provided test script:

**Linux/macOS:**
```bash
chmod +x scripts/test.sh
./scripts/test.sh
```

**Windows:**
```powershell
.\scripts\test.ps1
```

## Test Coverage

### Generate Coverage Report

```bash
# Generate coverage for internal packages
go test -coverprofile=coverage.out ./internal/...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Current Test Coverage

- `internal/models` - Data structures and JSON serialization
- `internal/validation` - Core validation logic and trust classification
- `internal/marketplace` - Microsoft Marketplace API client
- `internal/openvsx` - OpenVSX Registry API client

## Writing Tests

### Test File Naming

Test files should be named `*_test.go` and placed in the same directory as the code they test.

### Test Function Naming

```go
func TestFunctionName(t *testing.T) { ... }
func TestStructName_MethodName(t *testing.T) { ... }
func BenchmarkFunctionName(b *testing.B) { ... }
```

### Table-Driven Tests

Use table-driven tests for testing multiple scenarios:

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected TrustLevel
    }{
        {"Valid extension", "ms-python.python", TrustLevelLegitimate},
        {"Invalid format", "bad-format", TrustLevelUnknown},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Validate(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Mocking External Dependencies

For tests that would make real HTTP requests, use test doubles or skip:

```go
func TestFetchMetadata(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    // Test implementation
}
```

Run without integration tests:
```bash
go test -short ./...
```

### Testing With Temporary Files

Use `t.TempDir()` for temporary directories:

```go
func TestScanner(t *testing.T) {
    tempDir := t.TempDir() // Automatically cleaned up
    // Create test files in tempDir
    // Run tests
}
```

## Continuous Integration

Tests are run automatically on:
- Pull requests
- Commits to main branch
- Release tags

See `.github/workflows/test.yml` for CI configuration.

## Troubleshooting

### "pattern all:frontend/dist: no matching files found"

This error occurs when testing the main package before building the frontend. Solutions:

1. **Option 1**: Test internal packages only
   ```bash
   go test ./internal/...
   ```

2. **Option 2**: Build the frontend first
   ```bash
   cd frontend && npm install && npm run build && cd ..
   go test ./...
   ```

3. **Option 3**: Skip main package tests
   ```bash
   go test $(go list ./... | grep -v "^github.com/yourusername/secureopenvsx$")
   ```

### "missing go.sum entry"

Run:
```bash
go mod tidy
go mod download
```

### Frontend Tests

Frontend tests use Vitest (if configured):

```bash
cd frontend
npm test
npm run test:coverage
```

## Test Best Practices

1. **Test behavior, not implementation**
2. **Use descriptive test names**
3. **Keep tests independent** - each test should work in isolation
4. **Use table-driven tests** for multiple scenarios
5. **Mock external dependencies** - don't make real API calls
6. **Test error cases** - not just happy paths
7. **Use t.Helper()** for test helper functions
8. **Clean up resources** - use defer or t.Cleanup()

## Example Test Structure

```go
package mypackage

import (
    "testing"
)

// TestNewClient tests client initialization
func TestNewClient(t *testing.T) {
    client := NewClient()
    
    if client == nil {
        t.Fatal("NewClient returned nil")
    }
    
    if client.httpClient == nil {
        t.Error("httpClient is nil")
    }
}

// TestValidation demonstrates table-driven testing
func TestValidation(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        expectError bool
    }{
        {"valid input", "test.extension", false},
        {"invalid input", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Validate(tt.input)
            if (err != nil) != tt.expectError {
                t.Errorf("Validate() error = %v, expectError %v", err, tt.expectError)
            }
        })
    }
}
```

## Running Specific Tests

```bash
# Run tests matching pattern
go test -run TestValidation ./...

# Run single test
go test -run TestNewClient ./internal/marketplace

# Run with race detector
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Verbose output with test names
go test -v ./...
```

## Code Coverage Goals

- **Critical packages** (validation, security): 80%+
- **Business logic**: 70%+
- **Utilities**: 60%+

---

For build instructions, see [BUILD.md](BUILD.md)
