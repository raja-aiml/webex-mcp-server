package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthHandler(t *testing.T) {
	handler := HealthHandler("test-service", "1.0.0")

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check content type
	expectedContentType := "application/json"
	if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			ct, expectedContentType)
	}

	// Check response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Verify response fields
	if response["status"] != "healthy" {
		t.Errorf("Expected status to be 'healthy', got %v", response["status"])
	}

	if response["service"] != "test-service" {
		t.Errorf("Expected service to be 'test-service', got %v", response["service"])
	}

	if response["version"] != "1.0.0" {
		t.Errorf("Expected version to be '1.0.0', got %v", response["version"])
	}

	// Check time format
	if timeStr, ok := response["time"].(string); ok {
		if _, err := time.Parse(time.RFC3339, timeStr); err != nil {
			t.Errorf("Time is not in RFC3339 format: %v", err)
		}
	} else {
		t.Error("Time field is missing or not a string")
	}
}

func TestSetupHTTPHandlers(t *testing.T) {
	// Create a mock MCP server (we'll use nil for this test)
	mux := SetupHTTPHandlers(nil, "test-service", "1.0.0")

	if mux == nil {
		t.Fatal("SetupHTTPHandlers returned nil")
	}

	// Test that health endpoint is registered
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("health endpoint returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHealthHandler_ErrorHandling(t *testing.T) {
	// This test verifies that the error logging in HealthHandler doesn't panic
	// We can't easily test the actual logging, but we can ensure it handles errors gracefully

	handler := HealthHandler("test-service", "1.0.0")

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a custom ResponseWriter that fails on Write
	failingWriter := &failingResponseWriter{
		ResponseWriter: httptest.NewRecorder(),
		failOnWrite:    true,
	}

	// This should not panic
	handler.ServeHTTP(failingWriter, req)
}

// failingResponseWriter is a ResponseWriter that can be configured to fail
type failingResponseWriter struct {
	http.ResponseWriter
	failOnWrite bool
}

func (f *failingResponseWriter) Write([]byte) (int, error) {
	if f.failOnWrite {
		return 0, http.ErrAbortHandler
	}
	return f.ResponseWriter.Write(nil)
}
