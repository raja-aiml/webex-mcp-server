# Testing Guide

## Running Tests

### Run all tests
```bash
make test
# or
make dev test
```

### Run tests with coverage
```bash
make test-coverage
```

This will:
1. Run all tests with race detection
2. Generate `coverage.out` file
3. Create `coverage.html` for viewing in browser
4. Display total coverage percentage

### Run specific package tests
```bash
go test ./internal/config/...
go test ./internal/webex/...
go test ./internal/tools/...
```

## Test Coverage Summary

The project includes comprehensive unit tests with the following coverage:

| Package | Coverage | Description |
|---------|----------|-------------|
| `internal/config` | ~80% | Configuration management and providers |
| `internal/handlers` | ~92% | HTTP handlers and health checks |
| `internal/webex` | ~74% | Webex API client |
| `internal/server` | ~38% | Server creation and setup |
| `internal/app` | ~70% | Application lifecycle |
| `internal/tools` | Partial | Tool implementations |

## Test Structure

```
internal/
├── testutil/           # Shared test utilities
│   └── testutil.go    # Helper functions for tests
├── config/
│   ├── config_test.go      # Config function tests
│   ├── provider_test.go    # Provider interface tests
│   └── mock_provider.go    # Mock implementation for testing
├── webex/
│   └── client_test.go      # HTTP client tests
├── handlers/
│   └── handlers_test.go    # HTTP handler tests
├── server/
│   ├── server_test.go      # Server creation tests
│   └── config_test.go      # Config initialization tests
├── app/
│   └── app_test.go         # Application lifecycle tests
└── tools/
    └── base_test.go        # Base tool functionality tests
```

## Writing Tests

### 1. Use Table-Driven Tests
```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "test",
            want:    "TEST",
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### 2. Use Test Utilities
```go
// Set environment variable with cleanup
cleanup := testutil.SetEnv(t, "KEY", "value")
defer cleanup()

// Create mock HTTP server
server := testutil.MockHTTPServer(t, handlers)
defer server.Close()
```

### 3. Mock External Dependencies
```go
// Use MockProvider for testing
mockProvider := &config.MockProvider{
    Token:   "test-token",
    BaseURL: "https://test.api.com",
}
```

## Integration Tests

For integration testing with actual Webex APIs:

1. Set real API credentials in environment
2. Run with integration tag:
```bash
go test -tags=integration ./...
```

## Continuous Integration

The tests are designed to run in CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run tests
  run: make test-coverage
  
- name: Upload coverage
  uses: codecov/codecov-action@v3
  with:
    file: ./coverage.out
```

## Best Practices

1. **Isolation**: Each test should be independent
2. **Cleanup**: Always clean up resources (use defer)
3. **Mocking**: Mock external dependencies
4. **Coverage**: Aim for >80% coverage
5. **Speed**: Keep unit tests fast (<1s per test)
6. **Naming**: Use descriptive test names
7. **Assertions**: One logical assertion per test

## Debugging Tests

### Run single test
```bash
go test -run TestName ./internal/package/...
```

### Verbose output
```bash
go test -v ./...
```

### Debug with delve
```bash
dlv test ./internal/package/...
```