package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestRunHTTPServer(t *testing.T) {
	// Create a test server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "test-server",
		Version: "1.0.0",
	}, nil)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Run server (should exit when context is cancelled)
	err := RunHTTPServer(ctx, ":0", server, "test-service", "1.0.0")
	if err != nil && err != context.DeadlineExceeded {
		t.Errorf("RunHTTPServer() error = %v", err)
	}
}

func TestRunHTTPServer_HealthCheck(t *testing.T) {
	// Create a test server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "test-server",
		Version: "1.0.0",
	}, nil)

	// Create a context with cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- RunHTTPServer(ctx, ":9999", server, "test-service", "1.0.0")
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test health endpoint
	resp, err := http.Get("http://localhost:9999/health")
	if err != nil {
		t.Fatalf("Failed to reach health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Health check status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	// Cancel context to stop server
	cancel()

	// Wait for server to shut down
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("RunHTTPServer() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("Server didn't shut down in time")
	}
}

func TestRunStdioServer(t *testing.T) {
	// Skip this test as it interferes with stdin/stdout
	t.Skip("Skipping stdio test - interferes with test runner stdin/stdout")
}
