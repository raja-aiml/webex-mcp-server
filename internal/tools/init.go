package tools

import (
	"fmt"
	"sync"

	"github.com/raja-aiml/webex-mcp-server/internal/config"
	"github.com/raja-aiml/webex-mcp-server/internal/webex"
)

var (
	defaultClient webex.HTTPClient
	clientOnce    sync.Once
	clientErr     error
)

// InitializeDefaultClient initializes the default HTTP client
// This should be called at application startup
func InitializeDefaultClient() error {
	clientOnce.Do(func() {
		cfg, err := config.Load()
		if err != nil {
			clientErr = err
			return
		}
		defaultClient, clientErr = webex.NewClientWithConfig(cfg)
	})
	return clientErr
}

// MustInitializeDefaultClient initializes the default client and panics on error
func MustInitializeDefaultClient() {
	if err := InitializeDefaultClient(); err != nil {
		panic(fmt.Sprintf("failed to initialize default client: %v", err))
	}
}

// getDefaultClient returns the default client or an error if not initialized
func getDefaultClient() (webex.HTTPClient, error) {
	if defaultClient == nil {
		if err := InitializeDefaultClient(); err != nil {
			return nil, fmt.Errorf("failed to initialize client: %w", err)
		}
	}
	return defaultClient, nil
}
