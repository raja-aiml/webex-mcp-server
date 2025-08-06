package config

import (
	"testing"

	"github.com/raja-aiml/webex-mcp-server-go/internal/testutil"
)

func TestDefaultProvider(t *testing.T) {
	// Set up test environment
	cleanup1 := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-key")
	defer cleanup1()
	cleanup2 := testutil.SetEnv(t, "WEBEX_API_BASE_URL", "https://test.api.com")
	defer cleanup2()
	cleanup3 := testutil.SetEnv(t, "USE_FASTHTTP", "false")
	defer cleanup3()

	provider := NewDefaultProvider()

	t.Run("GetWebexToken", func(t *testing.T) {
		got, _ := provider.GetWebexToken()
		if got != "test-key" {
			t.Errorf("GetWebexToken() = %v, want %v", got, "test-key")
		}
	})

	t.Run("GetWebexURL", func(t *testing.T) {
		got, _ := provider.GetWebexURL("/test")
		if got != "https://test.api.com/test" {
			t.Errorf("GetWebexURL() = %v, want %v", got, "https://test.api.com/test")
		}
	})

	t.Run("GetWebexHeaders", func(t *testing.T) {
		headers, _ := provider.GetWebexHeaders()
		if headers["Authorization"] != "Bearer test-key" {
			t.Errorf("Authorization header = %v, want %v", headers["Authorization"], "Bearer test-key")
		}
	})

	t.Run("GetWebexJSONHeaders", func(t *testing.T) {
		headers, _ := provider.GetWebexJSONHeaders()
		if headers["Content-Type"] != "application/json" {
			t.Errorf("Content-Type header = %v, want %v", headers["Content-Type"], "application/json")
		}
		if headers["Authorization"] != "Bearer test-key" {
			t.Errorf("Authorization header = %v, want %v", headers["Authorization"], "Bearer test-key")
		}
	})

	t.Run("GetWebexBaseURL", func(t *testing.T) {
		got, _ := provider.GetWebexBaseURL()
		if got != "https://test.api.com" {
			t.Errorf("GetWebexBaseURL() = %v, want %v", got, "https://test.api.com")
		}
	})

	t.Run("GetUseFastHTTP", func(t *testing.T) {
		if got := provider.GetUseFastHTTP(); got != false {
			t.Errorf("GetUseFastHTTP() = %v, want %v", got, false)
		}
	})
}

func TestProviderInterface(t *testing.T) {
	// Ensure MockProvider implements Provider interface
	var _ Provider = &MockProvider{}
	
	mock := &MockProvider{
		Token:   "mock-key",
		BaseURL: "https://mock.api",
		Headers: map[string]string{
			"Authorization": "Bearer mock-key",
		},
		JSONHeaders: map[string]string{
			"Authorization": "Bearer mock-key",
			"Content-Type":  "application/json",
		},
		UseFastHTTP: false,
	}

	got, _ := mock.GetWebexToken()
	if got != "mock-key" {
		t.Errorf("GetWebexToken() = %v, want %v", got, "mock-key")
	}

	got2, _ := mock.GetWebexURL("/endpoint")
	if got2 != "https://mock.api/endpoint" {
		t.Errorf("GetWebexURL() = %v, want %v", got2, "https://mock.api/endpoint")
	}

	headers, _ := mock.GetWebexHeaders()
	if got := headers["Authorization"]; got != "Bearer mock-key" {
		t.Errorf("GetWebexHeaders()[Authorization] = %v, want %v", got, "Bearer mock-key")
	}
}