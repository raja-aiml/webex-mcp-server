package tools

import (
	"os"
	"sync"
	"testing"

	"github.com/raja-aiml/webex-mcp-server/internal/config"
	"github.com/raja-aiml/webex-mcp-server/internal/testutil"
)

func TestMustInitializeDefaultClient(t *testing.T) {
	// Reset config for testing
	config.ResetForTesting()

	// Test panic case - no API key
	cleanup1 := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "")
	defer func() {
		cleanup1()
		config.ResetForTesting()
	}()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustInitializeDefaultClient() should have panicked")
		}
	}()

	MustInitializeDefaultClient()
}

func TestMustInitializeDefaultClient_Success(t *testing.T) {
	// Save original env var
	originalKey := os.Getenv("WEBEX_PUBLIC_WORKSPACE_API_KEY")

	// Reset config and client
	config.ResetForTesting()
	defaultClient = nil
	clientOnce = sync.Once{}
	clientErr = nil
	clientOnce = sync.Once{}
	clientErr = nil

	// Set up valid environment
	os.Setenv("WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-key")
	defer func() {
		// Restore original env
		if originalKey != "" {
			os.Setenv("WEBEX_PUBLIC_WORKSPACE_API_KEY", originalKey)
		} else {
			os.Unsetenv("WEBEX_PUBLIC_WORKSPACE_API_KEY")
		}
		config.ResetForTesting()
		defaultClient = nil
		clientOnce = sync.Once{}
		clientErr = nil
		clientOnce = sync.Once{}
		clientErr = nil
		clientOnce = sync.Once{}
		clientErr = nil
	}()

	// Should not panic with valid config
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustInitializeDefaultClient() panicked unexpectedly: %v", r)
		}
	}()

	MustInitializeDefaultClient()

	if defaultClient == nil {
		t.Error("MustInitializeDefaultClient() did not initialize default client")
	}
}

func TestInitializeDefaultClient_Error(t *testing.T) {
	// Reset config and client
	config.ResetForTesting()
	defaultClient = nil
	clientOnce = sync.Once{}
	clientErr = nil

	// Test without API key
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "")
	defer func() {
		cleanup()
		config.ResetForTesting()
		defaultClient = nil
		clientOnce = sync.Once{}
		clientErr = nil
		clientOnce = sync.Once{}
		clientErr = nil
	}()

	err := InitializeDefaultClient()
	if err == nil {
		t.Error("InitializeDefaultClient() should return error without API key")
	}
}

func TestGetDefaultClient_Error(t *testing.T) {
	// Reset client
	defaultClient = nil
	clientOnce = sync.Once{}
	clientErr = nil
	config.ResetForTesting()

	// Test without API key
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "")
	defer func() {
		cleanup()
		config.ResetForTesting()
		defaultClient = nil
		clientOnce = sync.Once{}
		clientErr = nil
		clientOnce = sync.Once{}
		clientErr = nil
	}()

	_, err := getDefaultClient()
	if err == nil {
		t.Error("getDefaultClient() should return error when client not initialized")
	}
}
