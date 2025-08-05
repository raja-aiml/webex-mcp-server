package webex

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/raja-aiml/webex-mcp-server-go/internal/config"
	"github.com/raja-aiml/webex-mcp-server-go/internal/testutil"
)

func TestNewClient(t *testing.T) {
	t.Run("creates client with default provider", func(t *testing.T) {
		client := NewClient()
		if client == nil {
			t.Error("NewClient() returned nil")
		}
		if c, ok := client.(*Client); ok {
			if c.configProvider == nil {
				t.Error("Client configProvider is nil")
			}
		}
	})
}

func TestNewClientWithConfig(t *testing.T) {
	tests := []struct {
		name        string
		useFastHTTP bool
		wantFast    bool
	}{
		{
			name:        "creates net/http client",
			useFastHTTP: false,
			wantFast:    false,
		},
		{
			name:        "creates fasthttp client",
			useFastHTTP: true,
			wantFast:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := &config.MockProvider{
				Token:       "test-token",
				BaseURL:     "https://test.api.com",
				UseFastHTTP: tt.useFastHTTP,
				Headers: map[string]string{
					"Authorization": "Bearer test-token",
				},
			}

			client := NewClientWithConfig(mockProvider)
			c, ok := client.(*Client)
			if !ok {
				t.Fatal("NewClientWithConfig() did not return *Client")
			}

			if c.useFastHTTP != tt.wantFast {
				t.Errorf("Client.useFastHTTP = %v, want %v", c.useFastHTTP, tt.wantFast)
			}

			if tt.wantFast && c.fastClient == nil {
				t.Error("FastHTTP client is nil")
			}
			if !tt.wantFast && c.httpClient == nil {
				t.Error("HTTP client is nil")
			}
		})
	}
}

func TestClient_Get(t *testing.T) {
	// Create test server
	server := testutil.MockHTTPServer(t, map[string]func(w http.ResponseWriter, r *http.Request){
		"/test": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("Expected GET, got %s", r.Method)
			}
			if r.Header.Get("Authorization") != "Bearer test-token" {
				t.Errorf("Missing or incorrect Authorization header")
			}
			
			// Check query params
			if r.URL.Query().Get("param1") != "value1" {
				t.Errorf("Missing query param: param1")
			}
			
			testutil.JSONResponse(w, http.StatusOK, map[string]interface{}{
				"result": "success",
			})
		},
	})
	defer server.Close()

	mockProvider := &config.MockProvider{
		Token:   "test-token",
		BaseURL: server.URL,
		Headers: map[string]string{
			"Authorization": "Bearer test-token",
			"Accept":        "application/json",
		},
		UseFastHTTP: false,
	}

	client := NewClientWithConfig(mockProvider)
	
	result, err := client.Get("/test", map[string]string{
		"param1": "value1",
	})
	
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	
	if result["result"] != "success" {
		t.Errorf("Get() result = %v, want success", result["result"])
	}
}

func TestClient_Post(t *testing.T) {
	// Create test server
	server := testutil.MockHTTPServer(t, map[string]func(w http.ResponseWriter, r *http.Request){
		"/test": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Expected POST, got %s", r.Method)
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Missing Content-Type header")
			}
			
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("Failed to decode body: %v", err)
			}
			
			if body["key"] != "value" {
				t.Errorf("Incorrect body data")
			}
			
			testutil.JSONResponse(w, http.StatusOK, map[string]interface{}{
				"created": true,
			})
		},
	})
	defer server.Close()

	mockProvider := &config.MockProvider{
		Token:   "test-token",
		BaseURL: server.URL,
		Headers: map[string]string{
			"Authorization": "Bearer test-token",
		},
		UseFastHTTP: false,
	}

	client := NewClientWithConfig(mockProvider)
	
	result, err := client.Post("/test", map[string]interface{}{
		"key": "value",
	})
	
	if err != nil {
		t.Fatalf("Post() error = %v", err)
	}
	
	if result["created"] != true {
		t.Errorf("Post() result = %v, want true", result["created"])
	}
}

func TestClient_Put(t *testing.T) {
	// Create test server
	server := testutil.MockHTTPServer(t, map[string]func(w http.ResponseWriter, r *http.Request){
		"/test": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Errorf("Expected PUT, got %s", r.Method)
			}
			
			testutil.JSONResponse(w, http.StatusOK, map[string]interface{}{
				"updated": true,
			})
		},
	})
	defer server.Close()

	mockProvider := &config.MockProvider{
		Token:   "test-token",
		BaseURL: server.URL,
		Headers: map[string]string{
			"Authorization": "Bearer test-token",
		},
		UseFastHTTP: false,
	}

	client := NewClientWithConfig(mockProvider)
	
	result, err := client.Put("/test", map[string]interface{}{
		"update": "data",
	})
	
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	
	if result["updated"] != true {
		t.Errorf("Put() result = %v, want true", result["updated"])
	}
}

func TestClient_Delete(t *testing.T) {
	// Create test server
	server := testutil.MockHTTPServer(t, map[string]func(w http.ResponseWriter, r *http.Request){
		"/test": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "DELETE" {
				t.Errorf("Expected DELETE, got %s", r.Method)
			}
			
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer server.Close()

	mockProvider := &config.MockProvider{
		Token:   "test-token",
		BaseURL: server.URL,
		Headers: map[string]string{
			"Authorization": "Bearer test-token",
		},
		UseFastHTTP: false,
	}

	client := NewClientWithConfig(mockProvider)
	
	err := client.Delete("/test")
	
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		wantErrContains string
	}{
		{
			name:           "handles 400 error",
			statusCode:     http.StatusBadRequest,
			responseBody:   `{"error": "bad request"}`,
			wantErrContains: "webex API error",
		},
		{
			name:           "handles 401 error",
			statusCode:     http.StatusUnauthorized,
			responseBody:   `{"message": "unauthorized"}`,
			wantErrContains: "webex API error",
		},
		{
			name:           "handles 500 error",
			statusCode:     http.StatusInternalServerError,
			responseBody:   `Internal Server Error`,
			wantErrContains: "HTTP 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testutil.MockHTTPServer(t, map[string]func(w http.ResponseWriter, r *http.Request){
				"/test": func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.statusCode)
					w.Write([]byte(tt.responseBody))
				},
			})
			defer server.Close()

			mockProvider := &config.MockProvider{
				Token:   "test-token",
				BaseURL: server.URL,
				Headers: map[string]string{
					"Authorization": "Bearer test-token",
				},
				UseFastHTTP: false,
			}

			client := NewClientWithConfig(mockProvider)
			
			_, err := client.Get("/test", nil)
			
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			
			if !contains(err.Error(), tt.wantErrContains) {
				t.Errorf("Error = %v, want to contain %s", err, tt.wantErrContains)
			}
		})
	}
}

func TestClient_buildURL(t *testing.T) {
	mockProvider := &config.MockProvider{
		BaseURL: "https://api.test.com",
	}
	
	client := &Client{
		configProvider: mockProvider,
	}
	
	tests := []struct {
		name     string
		endpoint string
		params   map[string]string
		want     string
	}{
		{
			name:     "builds URL without params",
			endpoint: "/test",
			params:   nil,
			want:     "https://api.test.com/test",
		},
		{
			name:     "builds URL with params",
			endpoint: "/test",
			params: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			want: "https://api.test.com/test?key1=value1&key2=value2",
		},
		{
			name:     "ignores empty params",
			endpoint: "/test",
			params: map[string]string{
				"key1": "value1",
				"key2": "",
			},
			want: "https://api.test.com/test?key1=value1",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := client.buildURL(tt.endpoint, tt.params)
			
			// For params test, just check that both params are present
			// since map iteration order is not guaranteed
			if tt.name == "builds URL with params" {
				if !contains(got, "key1=value1") || !contains(got, "key2=value2") {
					t.Errorf("buildURL() = %v, missing expected params", got)
				}
			} else if got != tt.want {
				t.Errorf("buildURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}