package server

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/raja-aiml/webex-mcp-server-go/config"
)

// InitializeConfig loads environment variables and validates configuration
func InitializeConfig() error {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	return nil
}
