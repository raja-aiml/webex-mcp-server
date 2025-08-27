package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type Config struct {
	WebexAPIKey     string
	WebexToken      string // Alias for compatibility
	WebexAPIBaseURL string
	Port            string
	NodeEnv         string // For environment detection
}

var (
	once     sync.Once
	instance *Config
	loadErr  error
)

// ResetForTesting resets the singleton state for testing purposes
func ResetForTesting() {
	once = sync.Once{}
	instance = nil
	loadErr = nil
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	once.Do(func() {
		instance = &Config{
			WebexAPIKey:     os.Getenv("WEBEX_PUBLIC_WORKSPACE_API_KEY"),
			WebexAPIBaseURL: getEnvWithDefault("WEBEX_API_BASE_URL", "https://webexapis.com/v1"),
			Port:            getEnvWithDefault("PORT", "3001"),
			NodeEnv:         getEnvWithDefault("NODE_ENV", "development"),
		}

		// Clean up API key
		instance.WebexAPIKey = strings.TrimPrefix(instance.WebexAPIKey, "Bearer ")
		instance.WebexAPIKey = strings.TrimSpace(instance.WebexAPIKey)
		instance.WebexToken = instance.WebexAPIKey // Set alias

		// Validate
		if instance.WebexAPIKey == "" {
			loadErr = fmt.Errorf("WEBEX_PUBLIC_WORKSPACE_API_KEY environment variable is not set")
		}
	})

	return instance, loadErr
}

// MustLoad loads the configuration and panics if there's an error
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}
	return cfg
}

func Validate() error {
	_, err := Load()
	return err
}

// GetWebexHeaders returns headers for Webex API requests
func GetWebexHeaders(contentType ...string) (map[string]string, error) {
	cfg, err := Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	headers := map[string]string{
		"Accept":        "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", cfg.WebexAPIKey),
	}

	if len(contentType) > 0 && contentType[0] != "" {
		headers["Content-Type"] = contentType[0]
	}

	return headers, nil
}

// GetWebexJSONHeaders is a convenience function for JSON requests
func GetWebexJSONHeaders() (map[string]string, error) {
	return GetWebexHeaders("application/json")
}

func GetWebexURL(endpoint string) (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	return cfg.WebexAPIBaseURL + endpoint, nil
}

func GetWebexToken() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}
	return cfg.WebexAPIKey, nil
}

func GetWebexBaseURL() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}
	return cfg.WebexAPIBaseURL, nil
}

// Environment detection functions
func IsProduction() bool {
	cfg, _ := Load()
	return cfg != nil && cfg.NodeEnv == "production"
}

func IsDevelopment() bool {
	cfg, _ := Load()
	return cfg == nil || cfg.NodeEnv == "development" || cfg.NodeEnv == ""
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
