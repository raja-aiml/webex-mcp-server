package webex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/raja-aiml/webex-mcp-server/internal/config"
)

// handleHTTPError processes HTTP error responses in a consistent way
func handleHTTPError(resp *http.Response, body []byte) error {
	var errorData map[string]interface{}
	if err := json.Unmarshal(body, &errorData); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return fmt.Errorf("webex API error: %v", errorData)
}

// Client provides a simple HTTP client for Webex API calls
type Client struct {
	httpClient *http.Client
	baseURL    string
	headers    map[string]string
}

// NewClient creates a client with configuration from environment
func NewClient() (HTTPClient, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	headers, err := config.GetWebexHeaders()
	if err != nil {
		return nil, fmt.Errorf("failed to get headers: %w", err)
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: cfg.WebexAPIBaseURL,
		headers: headers,
	}, nil
}

// NewClientWithConfig creates a client with provided configuration
func NewClientWithConfig(cfg *config.Config) (HTTPClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	headers := map[string]string{
		"Accept":        "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", cfg.WebexAPIKey),
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: cfg.WebexAPIBaseURL,
		headers: headers,
	}, nil
}

// Get performs a GET request
func (c *Client) Get(endpoint string, params map[string]string) (map[string]interface{}, error) {
	fullURL, err := c.buildURL(endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}
	return c.doRequest("GET", fullURL, nil)
}

// Post performs a POST request
func (c *Client) Post(endpoint string, data interface{}) (map[string]interface{}, error) {
	fullURL := c.buildSimpleURL(endpoint)
	return c.doRequest("POST", fullURL, data)
}

// Put performs a PUT request
func (c *Client) Put(endpoint string, data interface{}) (map[string]interface{}, error) {
	fullURL := c.buildSimpleURL(endpoint)
	return c.doRequest("PUT", fullURL, data)
}

// Delete performs a DELETE request
func (c *Client) Delete(endpoint string) error {
	fullURL := c.buildSimpleURL(endpoint)
	_, err := c.doRequest("DELETE", fullURL, nil)
	return err
}

// buildSimpleURL constructs URL without query parameters
func (c *Client) buildSimpleURL(endpoint string) string {
	if endpoint != "" && endpoint[0] != '/' {
		endpoint = "/" + endpoint
	}
	return c.baseURL + endpoint
}

// buildURL constructs URL with query parameters
func (c *Client) buildURL(endpoint string, params map[string]string) (string, error) {
	fullURL := c.buildSimpleURL(endpoint)

	if len(params) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return "", err
		}
		q := u.Query()
		for key, value := range params {
			if value != "" {
				q.Set(key, value)
			}
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}
	return fullURL, nil
}

// doRequest executes the HTTP request
func (c *Client) doRequest(method, url string, data interface{}) (map[string]interface{}, error) {
	var reqBody io.Reader

	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to close response body: %v\n", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return c.handleResponse(resp, respBody)
}

// handleResponse processes the HTTP response
func (c *Client) handleResponse(resp *http.Response, body []byte) (map[string]interface{}, error) {
	if resp.StatusCode >= 400 {
		return nil, handleHTTPError(resp, body)
	}

	if len(body) == 0 {
		return nil, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
