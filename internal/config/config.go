package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type Config struct {
	WebexAPIKey     string
	WebexAPIBaseURL string
	Port            string
}

var (
	once     sync.Once
	instance *Config
	loadErr  error
)

// Load loads the configuration from environment variables
// It uses a singleton pattern to ensure config is loaded only once
func Load() (*Config, error) {
	once.Do(func() {
		instance = &Config{
			WebexAPIKey:     os.Getenv("WEBEX_PUBLIC_WORKSPACE_API_KEY"),
			WebexAPIBaseURL: getEnvWithDefault("WEBEX_API_BASE_URL", "https://webexapis.com/v1"),
			Port:            getEnvWithDefault("PORT", "3001"),
		}

		// Clean up API key (remove Bearer prefix if present)
		instance.WebexAPIKey = strings.TrimPrefix(instance.WebexAPIKey, "Bearer ")
		instance.WebexAPIKey = strings.TrimSpace(instance.WebexAPIKey)

		// Validate during load
		if instance.WebexAPIKey == "" {
			loadErr = fmt.Errorf("WEBEX_PUBLIC_WORKSPACE_API_KEY environment variable is not set")
		}
	})

	return instance, loadErr
}

// MustLoad loads the configuration and panics if there's an error
// Use this in initialization code where errors are fatal
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

func GetWebexHeaders() (map[string]string, error) {
	cfg, err := Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	
	return map[string]string{
		"Accept":        "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", cfg.WebexAPIKey),
	}, nil
}

func GetWebexJSONHeaders() (map[string]string, error) {
	headers, err := GetWebexHeaders()
	if err != nil {
		return nil, err
	}
	headers["Content-Type"] = "application/json"
	return headers, nil
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

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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

func GetUseFastHTTP() bool {
	return os.Getenv("USE_FASTHTTP") != "false"
}