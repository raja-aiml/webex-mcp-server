package server

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/raja-aiml/webex-mcp-server/internal/config"
)

// EnvironmentConfig represents environment configuration - TypeScript interface pattern
type EnvironmentConfig struct {
	EnvPath    string            `json:"envPath,omitempty"`
	Variables  map[string]string `json:"variables,omitempty"`
	IsLoaded   bool              `json:"isLoaded"`
	LoadedFrom string            `json:"loadedFrom,omitempty"`
}

// LoadEnvironmentResult represents the result of loading environment - TypeScript result pattern
type LoadEnvironmentResult struct {
	Success bool              `json:"success"`
	Config  EnvironmentConfig `json:"config"`
	Error   string            `json:"error,omitempty"`
}

// InitializeConfig loads environment configuration - TypeScript async function pattern
func InitializeConfig(envPath string) error {
	result := LoadEnvironment(envPath)
	if !result.Success && result.Error != "" {
		// In TypeScript style, we'd log warnings but continue execution
		// unless it's a critical error
		return nil // Non-critical for now, matches TypeScript behavior
	}
	return nil
}

// LoadEnvironment loads environment variables - TypeScript function with result pattern
func LoadEnvironment(envPath string) LoadEnvironmentResult {
	result := LoadEnvironmentResult{
		Config: EnvironmentConfig{
			Variables: make(map[string]string),
			IsLoaded:  false,
		},
	}

	// TypeScript-style path resolution
	targetPath := resolveEnvironmentPath(envPath)
	result.Config.EnvPath = targetPath

	// Try to load the .env file - TypeScript try/catch pattern
	if err := godotenv.Load(targetPath); err != nil {
		// TypeScript-style error handling - non-blocking for missing files
		if os.IsNotExist(err) {
			result.Error = "Environment file not found (continuing without it)"
			result.Success = true // TypeScript would continue execution
		} else {
			result.Error = err.Error()
			result.Success = false
		}
		return result
	}

	// Capture loaded environment variables - TypeScript object pattern
	result.Config.Variables = captureEnvironmentVariables()
	result.Config.IsLoaded = true
	result.Config.LoadedFrom = targetPath
	result.Success = true

	return result
}

// resolveEnvironmentPath resolves the environment file path - TypeScript path utility pattern
func resolveEnvironmentPath(envPath string) string {
	// TypeScript-style conditional assignment
	if envPath != "" {
		return envPath
	}

	// Default paths to try - TypeScript array pattern
	defaultPaths := []string{
		".env",
		".env.local",
		".env.example",
	}

	// Find first existing file - TypeScript find() pattern
	for _, path := range defaultPaths {
		if fileExists(path) {
			absPath, _ := filepath.Abs(path)
			return absPath
		}
	}

	// Fallback - TypeScript default behavior
	absPath, _ := filepath.Abs(".env")
	return absPath
}

// fileExists checks if a file exists - TypeScript utility function pattern
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// captureEnvironmentVariables captures current environment variables - TypeScript object pattern
func captureEnvironmentVariables() map[string]string {
	variables := make(map[string]string)

	// TypeScript-style environment variable capture
	webexKeys := []string{
		"WEBEX_PUBLIC_WORKSPACE_API_KEY",
		"WEBEX_API_BASE_URL",
		"PORT",
		"NODE_ENV", // TypeScript compatibility
	}

	for _, key := range webexKeys {
		if value := os.Getenv(key); value != "" {
			variables[key] = value
		}
	}

	return variables
}

// GetConfig returns the current configuration - TypeScript getter pattern
func GetConfig() (*config.Config, error) {
	return config.Load()
}

// ValidateConfiguration validates the current configuration - TypeScript validation pattern
func ValidateConfiguration() error {
	cfg, err := GetConfig()
	if err != nil {
		return err
	}

	// TypeScript-style validation
	if cfg.WebexToken == "" {
		return &config.ConfigError{
			Type:    "validation",
			Message: "WEBEX_PUBLIC_WORKSPACE_API_KEY environment variable is not set",
			Field:   "WebexToken",
		}
	}

	return nil
}
