package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// Config represents application configuration - TypeScript interface pattern
type Config struct {
	WebexAPIKey     string `json:"webexApiKey"`     // TypeScript camelCase naming
	WebexToken      string `json:"webexToken"`      // Alias for compatibility
	WebexAPIBaseURL string `json:"webexApiBaseUrl"` // TypeScript camelCase naming
	Port            string `json:"port"`
	NodeEnv         string `json:"nodeEnv,omitempty"` // TypeScript NODE_ENV equivalent
}

// ConfigError represents configuration errors - TypeScript error pattern
type ConfigError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// Error implements the error interface - TypeScript Error pattern
func (e *ConfigError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s error in field '%s': %s", e.Type, e.Field, e.Message)
	}
	return fmt.Sprintf("%s error: %s", e.Type, e.Message)
}

// ConfigResult represents the result of configuration loading - TypeScript result pattern
type ConfigResult struct {
	Config  *Config `json:"config,omitempty"`
	Success bool    `json:"success"`
	Error   string  `json:"error,omitempty"`
}

var (
	once     sync.Once
	instance *Config
	loadErr  error
)

// ResetForTesting resets the singleton state for testing purposes
// This function should only be used in tests - TypeScript pattern
func ResetForTesting() {
	once = sync.Once{}
	instance = nil
	loadErr = nil
}

// Load loads the configuration from environment variables - TypeScript async pattern
// It uses a singleton pattern to ensure config is loaded only once
func Load() (*Config, error) {
	once.Do(func() {
		instance = &Config{
			WebexAPIKey:     os.Getenv("WEBEX_PUBLIC_WORKSPACE_API_KEY"),
			WebexAPIBaseURL: getEnvWithDefault("WEBEX_API_BASE_URL", "https://webexapis.com/v1"),
			Port:            getEnvWithDefault("PORT", "3001"),
			NodeEnv:         getEnvWithDefault("NODE_ENV", "development"), // TypeScript NODE_ENV
		}

		// Clean up API key (remove Bearer prefix if present) - TypeScript string processing
		instance.WebexAPIKey = strings.TrimPrefix(instance.WebexAPIKey, "Bearer ")
		instance.WebexAPIKey = strings.TrimSpace(instance.WebexAPIKey)

		// Set token alias for compatibility
		instance.WebexToken = instance.WebexAPIKey

		// Validate during load - TypeScript validation pattern
		if instance.WebexAPIKey == "" {
			loadErr = &ConfigError{
				Type:    "validation",
				Message: "WEBEX_PUBLIC_WORKSPACE_API_KEY environment variable is not set",
				Field:   "WebexAPIKey",
			}
		}
	})

	return instance, loadErr
}

// LoadConfig loads configuration with TypeScript-style result pattern
func LoadConfig() ConfigResult {
	cfg, err := Load()
	if err != nil {
		return ConfigResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ConfigResult{
		Config:  cfg,
		Success: true,
	}
}

// MustLoad loads the configuration and panics if there's an error
// Use this in initialization code where errors are fatal - TypeScript pattern
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}
	return cfg
}

// Validate validates the current configuration - TypeScript validation pattern
func Validate() error {
	_, err := Load()
	return err
}

// ValidateWebexConfig validates Webex-specific configuration - TypeScript pattern
func ValidateWebexConfig() error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	// TypeScript-style validation checks
	requiredFields := map[string]string{
		"WebexAPIKey": cfg.WebexAPIKey,
	}

	for field, value := range requiredFields {
		if value == "" {
			return &ConfigError{
				Type:    "validation",
				Message: fmt.Sprintf("required field %s is empty", field),
				Field:   field,
			}
		}
	}

	return nil
}

// GetWebexHeaders returns headers for Webex API requests - TypeScript object pattern
// contentType is optional - if provided, it sets the Content-Type header
func GetWebexHeaders(contentType ...string) (map[string]string, error) {
	cfg, err := Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// TypeScript-style object creation
	headers := map[string]string{
		"Accept":        "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", cfg.WebexAPIKey),
	}

	// TypeScript-style optional parameter handling
	if len(contentType) > 0 && contentType[0] != "" {
		headers["Content-Type"] = contentType[0]
	}

	return headers, nil
}

// GetWebexJSONHeaders is a convenience function for JSON requests - TypeScript pattern
func GetWebexJSONHeaders() (map[string]string, error) {
	return GetWebexHeaders("application/json")
}

// GetWebexURL constructs Webex API URLs - TypeScript URL construction pattern
func GetWebexURL(endpoint string) (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	// TypeScript-style string manipulation
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	return cfg.WebexAPIBaseURL + endpoint, nil
}

// GetWebexToken returns the Webex API token - TypeScript getter pattern
func GetWebexToken() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}
	return cfg.WebexAPIKey, nil
}

// GetWebexBaseURL returns the Webex API base URL - TypeScript getter pattern
func GetWebexBaseURL() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}
	return cfg.WebexAPIBaseURL, nil
}

// IsProduction checks if running in production - TypeScript environment check pattern
func IsProduction() bool {
	cfg, _ := Load()
	if cfg == nil {
		return false
	}
	return cfg.NodeEnv == "production"
}

// IsDevelopment checks if running in development - TypeScript environment check pattern
func IsDevelopment() bool {
	cfg, _ := Load()
	if cfg == nil {
		return true // Default to development
	}
	return cfg.NodeEnv == "development" || cfg.NodeEnv == ""
}

// getEnvWithDefault gets environment variable with default value - TypeScript utility pattern
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
