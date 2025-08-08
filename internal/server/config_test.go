package server

import (
	"os"
	"testing"

	"github.com/raja-aiml/webex-mcp-server/internal/config"
	"github.com/raja-aiml/webex-mcp-server/internal/testutil"
)

func TestInitializeConfig(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "successful initialization with API key",
			setup: func() {
				os.Setenv("WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-key")
			},
			wantErr: false,
		},
		{
			name: "fails without API key",
			setup: func() {
				os.Unsetenv("WEBEX_PUBLIC_WORKSPACE_API_KEY")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.ResetForTesting()

			// Save original env
			origKey := os.Getenv("WEBEX_PUBLIC_WORKSPACE_API_KEY")
			defer func() {
				if origKey != "" {
					os.Setenv("WEBEX_PUBLIC_WORKSPACE_API_KEY", origKey)
				} else {
					os.Unsetenv("WEBEX_PUBLIC_WORKSPACE_API_KEY")
				}
				config.ResetForTesting()
			}()

			tt.setup()
			err := InitializeConfig("")

			if (err != nil) != tt.wantErr {
				t.Errorf("InitializeConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitializeConfig_LoadsEnvFile(t *testing.T) {
	config.ResetForTesting()
	defer config.ResetForTesting()

	// Create a temporary .env file
	envContent := []byte("TEST_ENV_VAR=test_value\n")
	if err := os.WriteFile(".env", envContent, 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(".env")

	// Set required API key
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-key")
	defer cleanup()

	err := InitializeConfig("")
	if err != nil {
		t.Errorf("InitializeConfig() error = %v", err)
	}

	// Check if env var from .env file was loaded
	if os.Getenv("TEST_ENV_VAR") != "test_value" {
		t.Error("Failed to load variables from .env file")
	}

	// Cleanup
	os.Unsetenv("TEST_ENV_VAR")
}
