package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	WebexAPIKey     string
	WebexAPIBaseURL string
	Port            string
}

func Load() (*Config, error) {
	cfg := &Config{
		WebexAPIKey:     os.Getenv("WEBEX_PUBLIC_WORKSPACE_API_KEY"),
		WebexAPIBaseURL: getEnvWithDefault("WEBEX_API_BASE_URL", "https://webexapis.com/v1"),
		Port:            getEnvWithDefault("PORT", "3001"),
	}

	// Clean up API key (remove Bearer prefix if present)
	cfg.WebexAPIKey = strings.TrimPrefix(cfg.WebexAPIKey, "Bearer ")
	cfg.WebexAPIKey = strings.TrimSpace(cfg.WebexAPIKey)

	return cfg, nil
}

func Validate() error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	if cfg.WebexAPIKey == "" {
		return fmt.Errorf("WEBEX_PUBLIC_WORKSPACE_API_KEY environment variable is not set")
	}

	return nil
}

func GetWebexHeaders() map[string]string {
	cfg, _ := Load()
	return map[string]string{
		"Accept":        "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", cfg.WebexAPIKey),
	}
}

func GetWebexJSONHeaders() map[string]string {
	headers := GetWebexHeaders()
	headers["Content-Type"] = "application/json"
	return headers
}

func GetWebexURL(endpoint string) string {
	cfg, _ := Load()
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	return cfg.WebexAPIBaseURL + endpoint
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetWebexToken() string {
	cfg, _ := Load()
	return cfg.WebexAPIKey
}

func GetWebexBaseURL() string {
	cfg, _ := Load()
	return cfg.WebexAPIBaseURL
}

func GetUseFastHTTP() bool {
	return os.Getenv("USE_FASTHTTP") != "false"
}