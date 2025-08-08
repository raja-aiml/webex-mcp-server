package tools

import (
	"testing"

	"github.com/raja-aiml/webex-mcp-server/internal/config"
	"github.com/raja-aiml/webex-mcp-server/internal/testutil"
)

func TestCorePluginImplementations(t *testing.T) {
	// Set up environment for tools that need default client
	config.ResetForTesting()
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-key")
	defer func() {
		cleanup()
		config.ResetForTesting()
		defaultClient = nil
	}()

	// Initialize default client
	_ = InitializeDefaultClient()

	tests := []struct {
		name   string
		plugin interface {
			Name() string
			Version() string
			Register(*Registry) error
		}
		wantName    string
		wantVersion string
	}{
		{
			name:        "coreMessagingPlugin",
			plugin:      &coreMessagingPlugin{},
			wantName:    "core-messaging",
			wantVersion: "1.0.0",
		},
		{
			name:        "coreWebhooksPlugin",
			plugin:      &coreWebhooksPlugin{},
			wantName:    "core-webhooks",
			wantVersion: "1.0.0",
		},
		{
			name:        "coreInfoPlugin",
			plugin:      &coreInfoPlugin{},
			wantName:    "core-info",
			wantVersion: "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.plugin.Name(); got != tt.wantName {
				t.Errorf("%s.Name() = %v, want %v", tt.name, got, tt.wantName)
			}
			if got := tt.plugin.Version(); got != tt.wantVersion {
				t.Errorf("%s.Version() = %v, want %v", tt.name, got, tt.wantVersion)
			}
		})
	}
}

func TestCorePluginRegister(t *testing.T) {
	// Set up environment for tools that need default client
	config.ResetForTesting()
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-key")
	defer func() {
		cleanup()
		config.ResetForTesting()
		defaultClient = nil
	}()

	// Initialize default client
	_ = InitializeDefaultClient()

	// Test that each plugin can register its tools without error
	plugins := []interface {
		Name() string
		Register(*Registry) error
	}{
		&coreMessagingPlugin{},
		&coreWebhooksPlugin{},
		&coreInfoPlugin{},
	}

	for _, plugin := range plugins {
		t.Run(plugin.Name(), func(t *testing.T) {
			registry := NewRegistry()
			err := plugin.Register(registry)
			if err != nil {
				t.Errorf("%s.Register() error = %v", plugin.Name(), err)
			}

			// Verify tools were registered
			tools := registry.GetTools()
			if len(tools) == 0 {
				t.Errorf("%s.Register() did not register any tools", plugin.Name())
			}
		})
	}
}
