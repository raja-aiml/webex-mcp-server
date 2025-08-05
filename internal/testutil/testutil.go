package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// SetEnv sets an environment variable and returns a cleanup function
func SetEnv(t *testing.T, key, value string) func() {
	t.Helper()
	old, exists := os.LookupEnv(key)
	os.Setenv(key, value)
	
	return func() {
		if exists {
			os.Setenv(key, old)
		} else {
			os.Unsetenv(key)
		}
	}
}

// MockHTTPServer creates a test HTTP server with custom handlers
func MockHTTPServer(t *testing.T, handlers map[string]func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	t.Helper()
	
	mux := http.NewServeMux()
	for path, handler := range handlers {
		mux.HandleFunc(path, handler)
	}
	
	return httptest.NewServer(mux)
}

// JSONResponse writes a JSON response
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// AssertJSON compares two JSON strings
func AssertJSON(t *testing.T, expected, actual string) {
	t.Helper()
	
	var expectedData, actualData interface{}
	if err := json.Unmarshal([]byte(expected), &expectedData); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(actual), &actualData); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}
	
	expectedJSON, _ := json.MarshalIndent(expectedData, "", "  ")
	actualJSON, _ := json.MarshalIndent(actualData, "", "  ")
	
	if string(expectedJSON) != string(actualJSON) {
		t.Errorf("JSON mismatch:\nExpected:\n%s\nActual:\n%s", expectedJSON, actualJSON)
	}
}