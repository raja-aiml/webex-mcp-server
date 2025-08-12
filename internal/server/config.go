package server

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/raja-aiml/webex-mcp-server/internal/config"
)

// InitializeConfig loads environment variables and validates configuration
func InitializeConfig(envPath string) error {
	// Load environment variables from specified path or default
	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			return fmt.Errorf("error loading .env file from %s: %w", envPath, err)
		}
	} else {
		// Try loading from default location
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file not found in default location: %v", err)
		}
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	return nil
}
