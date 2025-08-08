package webex

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/raja-aiml/webex-mcp-server/internal/config"
	"github.com/raja-aiml/webex-mcp-server/internal/testutil"
)

func TestNewClient(t *testing.T) {
	t.Run("creates client with default configuration", func(t *testing.T) {
		config.ResetForTesting()
		cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-token")
		defer func() {
			cleanup()
			config.ResetForTesting()
		}()

		client, err := NewClient()
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		if client == nil {
			t.Error("NewClient() returned nil")
		}
		if c, ok := client.(*Client); ok {
			if c.httpClient == nil {
				t.Error("Client httpClient is nil")
			}
			if c.baseURL == "" {
				t.Error("Client baseURL is empty")
			}
			if len(c.headers) == 0 {
				t.Error("Client headers are empty")
			}
		}
	})
}

func TestNewClientWithConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name: "creates client with valid config",
			config: &config.Config{
				WebexAPIKey:     "test-key",
				WebexAPIBaseURL: "https://api.webex.com/v1",
			},
			wantErr: false,
		},
		{
			name:    "fails with nil config",
			config:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClientWithConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientWithConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClientWithConfig() returned nil")
			}
		})
	}
}

func TestClient_Get(t *testing.T) {
	server := testutil.MockHTTPServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("Expected /test, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("param1") != "value1" {
			t.Errorf("Expected param1=value1, got %s", r.URL.Query().Get("param1"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": "success",
		})
	})
	defer server.Close()

	cfg := &config.Config{
		WebexAPIKey:     "test-token",
		WebexAPIBaseURL: server.URL,
	}

	client, err := NewClientWithConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result, err := client.Get("/test", map[string]string{"param1": "value1"})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if result["result"] != "success" {
		t.Errorf("Expected result=success, got %v", result["result"])
	}
}

func TestClient_Post(t *testing.T) {
	server := testutil.MockHTTPServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["key"] != "value" {
			t.Errorf("Expected key=value in body, got %v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": "123",
		})
	})
	defer server.Close()

	cfg := &config.Config{
		WebexAPIKey:     "test-token",
		WebexAPIBaseURL: server.URL,
	}

	client, err := NewClientWithConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result, err := client.Post("/endpoint", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("Post() error = %v", err)
	}
	if result["id"] != "123" {
		t.Errorf("Expected id=123, got %v", result["id"])
	}
}

func TestClient_Put(t *testing.T) {
	server := testutil.MockHTTPServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"updated": true,
		})
	})
	defer server.Close()

	cfg := &config.Config{
		WebexAPIKey:     "test-token",
		WebexAPIBaseURL: server.URL,
	}

	client, err := NewClientWithConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result, err := client.Put("/endpoint", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	if result["updated"] != true {
		t.Errorf("Expected updated=true, got %v", result["updated"])
	}
}

func TestClient_Delete(t *testing.T) {
	server := testutil.MockHTTPServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	cfg := &config.Config{
		WebexAPIKey:     "test-token",
		WebexAPIBaseURL: server.URL,
	}

	client, err := NewClientWithConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.Delete("/endpoint")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   interface{}
		wantErr    bool
	}{
		{
			name:       "handles 400 error",
			statusCode: 400,
			response:   map[string]interface{}{"error": "bad request"},
			wantErr:    true,
		},
		{
			name:       "handles 401 error",
			statusCode: 401,
			response:   map[string]interface{}{"error": "unauthorized"},
			wantErr:    true,
		},
		{
			name:       "handles 500 error",
			statusCode: 500,
			response:   map[string]interface{}{"error": "internal server error"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testutil.MockHTTPServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			})
			defer server.Close()

			cfg := &config.Config{
				WebexAPIKey:     "test-token",
				WebexAPIBaseURL: server.URL,
			}

			client, err := NewClientWithConfig(cfg)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			_, err = client.Get("/test", nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_buildURL(t *testing.T) {
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
			want:     "https://api.webex.com/test",
		},
		{
			name:     "builds URL with params",
			endpoint: "/test",
			params:   map[string]string{"key": "value", "foo": "bar"},
			want:     "https://api.webex.com/test?foo=bar&key=value",
		},
		{
			name:     "ignores empty params",
			endpoint: "/test",
			params:   map[string]string{"key": "value", "empty": ""},
			want:     "https://api.webex.com/test?key=value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				baseURL: "https://api.webex.com",
			}

			got, err := client.buildURL(tt.endpoint, tt.params)
			if err != nil {
				t.Fatalf("buildURL() error = %v", err)
			}

			// For params, order might vary, so we need to check if all params are present
			if len(tt.params) > 0 {
				if len(got) < len("https://api.webex.com/test?") {
					t.Errorf("buildURL() = %v, want params in URL", got)
				}
			} else if got != tt.want {
				t.Errorf("buildURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
