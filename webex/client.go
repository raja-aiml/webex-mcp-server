package webex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/raja-aiml/webex-mcp-server-go/config"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	headers    map[string]string
}

// handleHTTPError processes HTTP error responses in a consistent way
func handleHTTPError(resp *http.Response, body []byte) error {
	var errorData map[string]interface{}
	if err := json.Unmarshal(body, &errorData); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return fmt.Errorf("webex API error: %v", errorData)
}

// processResponse handles common response processing for all HTTP methods
func processResponse(resp *http.Response) (map[string]interface{}, error) {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to close response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, handleHTTPError(resp, body)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: config.GetWebexURL(""),
		headers: config.GetWebexHeaders(),
	}
}

func (c *Client) Get(endpoint string, params map[string]string) (map[string]interface{}, error) {
	fullURL := config.GetWebexURL(endpoint)

	if len(params) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return nil, err
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

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return processResponse(resp)
}

func (c *Client) Post(endpoint string, data interface{}) (map[string]interface{}, error) {
	fullURL := config.GetWebexURL(endpoint)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	headers := config.GetWebexJSONHeaders()
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return processResponse(resp)
}

func (c *Client) Put(endpoint string, data interface{}) (map[string]interface{}, error) {
	fullURL := config.GetWebexURL(endpoint)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	headers := config.GetWebexJSONHeaders()
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return processResponse(resp)
}

func (c *Client) Delete(endpoint string) error {
	fullURL := config.GetWebexURL(endpoint)

	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return err
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return handleHTTPError(resp, body)
	}

	return nil
}

// DefaultClient returns the optimized client that automatically selects backend
func DefaultClient() HTTPClient {
	return NewOptimizedClient()
}
