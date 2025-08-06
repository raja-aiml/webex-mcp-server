package webex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/raja-aiml/webex-mcp-server-go/internal/config"
	"github.com/valyala/fasthttp"
)

// handleHTTPError processes HTTP error responses in a consistent way
func handleHTTPError(resp *http.Response, body []byte) error {
	var errorData map[string]interface{}
	if err := json.Unmarshal(body, &errorData); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return fmt.Errorf("webex API error: %v", errorData)
}

// Client provides a single, configurable HTTP client
// that can use either net/http or fasthttp based on configuration
type Client struct {
	useFastHTTP    bool
	httpClient     *http.Client
	fastClient     *fasthttp.Client
	baseURL        string
	headers        map[string]string
	configProvider config.Provider
}

// NewClient creates a client with automatic backend selection
func NewClient() (HTTPClient, error) {
	return NewClientWithConfig(config.NewDefaultProvider())
}

// NewClientWithConfig creates a client with dependency injection
func NewClientWithConfig(configProvider config.Provider) (HTTPClient, error) {
	useFastHTTP := configProvider.GetUseFastHTTP()

	baseURL, err := configProvider.GetWebexURL("")
	if err != nil {
		return nil, fmt.Errorf("failed to get base URL: %w", err)
	}

	headers, err := configProvider.GetWebexHeaders()
	if err != nil {
		return nil, fmt.Errorf("failed to get headers: %w", err)
	}

	client := &Client{
		useFastHTTP:    useFastHTTP,
		baseURL:        baseURL,
		headers:        headers,
		configProvider: configProvider,
	}

	if useFastHTTP {
		client.fastClient = &fasthttp.Client{
			MaxConnsPerHost:     100,
			MaxIdleConnDuration: 10 * time.Second,
			ReadTimeout:         30 * time.Second,
			WriteTimeout:        30 * time.Second,
		}
	} else {
		client.httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return client, nil
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
	fullURL, err := c.configProvider.GetWebexURL(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}
	return c.doRequest("POST", fullURL, data)
}

// Put performs a PUT request
func (c *Client) Put(endpoint string, data interface{}) (map[string]interface{}, error) {
	fullURL, err := c.configProvider.GetWebexURL(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}
	return c.doRequest("PUT", fullURL, data)
}

// Delete performs a DELETE request
func (c *Client) Delete(endpoint string) error {
	fullURL, err := c.configProvider.GetWebexURL(endpoint)
	if err != nil {
		return fmt.Errorf("failed to get URL: %w", err)
	}
	_, err = c.doRequest("DELETE", fullURL, nil)
	return err
}

// buildURL constructs URL with query parameters
func (c *Client) buildURL(endpoint string, params map[string]string) (string, error) {
	fullURL, err := c.configProvider.GetWebexURL(endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to get URL: %w", err)
	}
	if len(params) > 0 {
		u, _ := url.Parse(fullURL)
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

// doRequest executes the HTTP request using the configured backend
func (c *Client) doRequest(method, url string, data interface{}) (map[string]interface{}, error) {
	var body []byte
	if data != nil {
		var err error
		body, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
	}

	if c.useFastHTTP {
		return c.doFastHTTPRequest(method, url, body)
	}
	return c.doNetHTTPRequest(method, url, body)
}

// doNetHTTPRequest executes request using net/http
func (c *Client) doNetHTTPRequest(method, url string, body []byte) (map[string]interface{}, error) {
	var reqBody io.Reader
	if body != nil {
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
	if body != nil {
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

	return c.handleResponse(resp.StatusCode, respBody)
}

// doFastHTTPRequest executes request using fasthttp
func (c *Client) doFastHTTPRequest(method, url string, body []byte) (map[string]interface{}, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(method)

	// Set headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
		req.SetBody(body)
	}

	if err := c.fastClient.Do(req, resp); err != nil {
		return nil, err
	}

	return c.handleResponse(resp.StatusCode(), resp.Body())
}

// handleResponse processes the HTTP response
func (c *Client) handleResponse(statusCode int, body []byte) (map[string]interface{}, error) {
	if statusCode >= 400 {
		// Create a mock response for handleHTTPError
		resp := &http.Response{StatusCode: statusCode}
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
