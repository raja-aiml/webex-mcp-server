package config

import (
	"os"
	"testing"

	"github.com/raja-aiml/webex-mcp-server/internal/testutil"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "valid configuration",
			setup: func() {
				os.Setenv("WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-token")
			},
			wantErr: false,
		},
		{
			name: "missing API key",
			setup: func() {
				os.Unsetenv("WEBEX_PUBLIC_WORKSPACE_API_KEY")
			},
			wantErr: true,
		},
		{
			name: "empty API key",
			setup: func() {
				os.Setenv("WEBEX_PUBLIC_WORKSPACE_API_KEY", "")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset singleton state for isolated testing
			ResetForTesting()

			// Save original env
			origKey := os.Getenv("WEBEX_PUBLIC_WORKSPACE_API_KEY")
			defer func() {
				if origKey != "" {
					os.Setenv("WEBEX_PUBLIC_WORKSPACE_API_KEY", origKey)
				} else {
					os.Unsetenv("WEBEX_PUBLIC_WORKSPACE_API_KEY")
				}
				ResetForTesting() // Reset again after test
			}()

			tt.setup()
			err := Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetWebexToken(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "returns API key",
			envValue: "test-api-key",
			want:     "test-api-key",
		},
		{
			name:     "returns empty string when not set",
			envValue: "",
			want:     "",
		},
		{
			name:     "removes Bearer prefix",
			envValue: "Bearer test-api-key",
			want:     "test-api-key",
		},
		{
			name:     "trims whitespace",
			envValue: "  test-api-key  ",
			want:     "test-api-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetForTesting()
			cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", tt.envValue)
			defer func() {
				cleanup()
				ResetForTesting()
			}()

			got, _ := GetWebexToken()
			if got != tt.want {
				t.Errorf("GetWebexToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWebexURL(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		endpoint string
		want     string
	}{
		{
			name:     "uses default URL",
			envValue: "",
			endpoint: "/messages",
			want:     "https://webexapis.com/v1/messages",
		},
		{
			name:     "uses custom URL",
			envValue: "https://custom.webex.com/api",
			endpoint: "/rooms",
			want:     "https://custom.webex.com/api/rooms",
		},
		{
			name:     "handles empty endpoint",
			envValue: "",
			endpoint: "",
			want:     "https://webexapis.com/v1/",
		},
		{
			name:     "handles trailing slash",
			envValue: "https://webexapis.com/v1",
			endpoint: "/messages",
			want:     "https://webexapis.com/v1/messages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetForTesting()
			cleanup1 := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-token") // Required for config Load()
			cleanup2 := testutil.SetEnv(t, "WEBEX_API_BASE_URL", tt.envValue)
			defer func() {
				cleanup1()
				cleanup2()
				ResetForTesting()
			}()

			got, _ := GetWebexURL(tt.endpoint)
			if got != tt.want {
				t.Errorf("GetWebexURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWebexHeaders(t *testing.T) {
	ResetForTesting()
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-token")
	defer func() {
		cleanup()
		ResetForTesting()
	}()

	headers, _ := GetWebexHeaders()

	expectedHeaders := map[string]string{
		"Authorization": "Bearer test-token",
		"Accept":        "application/json",
	}

	for key, expectedValue := range expectedHeaders {
		if value, ok := headers[key]; !ok {
			t.Errorf("Missing header %s", key)
		} else if value != expectedValue {
			t.Errorf("Header %s = %v, want %v", key, value, expectedValue)
		}
	}

	if len(headers) != len(expectedHeaders) {
		t.Errorf("Got %d headers, want %d", len(headers), len(expectedHeaders))
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		want    *Config
		wantErr bool
	}{
		{
			name: "loads configuration with defaults",
			setup: func() {
				os.Setenv("WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-key")
				os.Unsetenv("WEBEX_API_BASE_URL")
				os.Unsetenv("PORT")
			},
			want: &Config{
				WebexAPIKey:     "test-key",
				WebexAPIBaseURL: "https://webexapis.com/v1",
				Port:            "3001",
			},
			wantErr: false,
		},
		{
			name: "loads configuration with custom values",
			setup: func() {
				os.Setenv("WEBEX_PUBLIC_WORKSPACE_API_KEY", "custom-key")
				os.Setenv("WEBEX_API_BASE_URL", "https://custom.api.com")
				os.Setenv("PORT", "8080")
			},
			want: &Config{
				WebexAPIKey:     "custom-key",
				WebexAPIBaseURL: "https://custom.api.com",
				Port:            "8080",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetForTesting()
			// Save original env
			origKey := os.Getenv("WEBEX_PUBLIC_WORKSPACE_API_KEY")
			origURL := os.Getenv("WEBEX_API_BASE_URL")
			origPort := os.Getenv("PORT")
			defer func() {
				os.Setenv("WEBEX_PUBLIC_WORKSPACE_API_KEY", origKey)
				os.Setenv("WEBEX_API_BASE_URL", origURL)
				os.Setenv("PORT", origPort)
				ResetForTesting()
			}()

			tt.setup()
			got, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if got.WebexAPIKey != tt.want.WebexAPIKey {
					t.Errorf("Load() APIKey = %v, want %v", got.WebexAPIKey, tt.want.WebexAPIKey)
				}
				if got.WebexAPIBaseURL != tt.want.WebexAPIBaseURL {
					t.Errorf("Load() BaseURL = %v, want %v", got.WebexAPIBaseURL, tt.want.WebexAPIBaseURL)
				}
				if got.Port != tt.want.Port {
					t.Errorf("Load() Port = %v, want %v", got.Port, tt.want.Port)
				}
			}
		})
	}
}

func TestGetWebexJSONHeaders(t *testing.T) {
	ResetForTesting()
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-token")
	defer func() {
		cleanup()
		ResetForTesting()
	}()

	headers, _ := GetWebexJSONHeaders()

	expectedHeaders := map[string]string{
		"Authorization": "Bearer test-token",
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	for key, expectedValue := range expectedHeaders {
		if value, ok := headers[key]; !ok {
			t.Errorf("Missing header %s", key)
		} else if value != expectedValue {
			t.Errorf("Header %s = %v, want %v", key, value, expectedValue)
		}
	}
}

func TestGetWebexBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "returns default URL",
			envValue: "",
			want:     "https://webexapis.com/v1",
		},
		{
			name:     "returns custom URL",
			envValue: "https://custom.api.com",
			want:     "https://custom.api.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetForTesting()
			cleanup1 := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-token") // Required for config Load()
			cleanup2 := testutil.SetEnv(t, "WEBEX_API_BASE_URL", tt.envValue)
			defer func() {
				cleanup1()
				cleanup2()
				ResetForTesting()
			}()

			got, _ := GetWebexBaseURL()
			if got != tt.want {
				t.Errorf("GetWebexBaseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustLoad(t *testing.T) {
	t.Run("panics when config load fails", func(t *testing.T) {
		ResetForTesting()
		cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "")
		defer func() {
			cleanup()
			ResetForTesting()
		}()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustLoad() should have panicked")
			}
		}()

		MustLoad()
	})

	t.Run("returns config when successful", func(t *testing.T) {
		ResetForTesting()
		cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-token")
		defer func() {
			cleanup()
			ResetForTesting()
		}()

		cfg := MustLoad()
		if cfg == nil {
			t.Errorf("MustLoad() returned nil config")
		}
		if cfg.WebexAPIKey != "test-token" {
			t.Errorf("MustLoad() WebexAPIKey = %v, want test-token", cfg.WebexAPIKey)
		}
	})
}

func TestGetWebexHeaders_ErrorCase(t *testing.T) {
	ResetForTesting()
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "")
	defer func() {
		cleanup()
		ResetForTesting()
	}()

	_, err := GetWebexHeaders()
	if err == nil {
		t.Errorf("GetWebexHeaders() should return error when config fails")
	}
}

func TestGetWebexURL_ErrorCase(t *testing.T) {
	ResetForTesting()
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "")
	defer func() {
		cleanup()
		ResetForTesting()
	}()

	_, err := GetWebexURL("/test")
	if err == nil {
		t.Errorf("GetWebexURL() should return error when config fails")
	}
}

func TestGetWebexBaseURL_ErrorCase(t *testing.T) {
	ResetForTesting()
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "")
	defer func() {
		cleanup()
		ResetForTesting()
	}()

	_, err := GetWebexBaseURL()
	if err == nil {
		t.Errorf("GetWebexBaseURL() should return error when config fails")
	}
}
