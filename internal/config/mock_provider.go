package config

// MockProvider is a test implementation of the Provider interface
type MockProvider struct {
	Token       string
	BaseURL     string
	Headers     map[string]string
	JSONHeaders map[string]string
	UseFastHTTP bool
}

func (m *MockProvider) GetWebexToken() string {
	return m.Token
}

func (m *MockProvider) GetWebexBaseURL() string {
	return m.BaseURL
}

func (m *MockProvider) GetWebexURL(endpoint string) string {
	return m.BaseURL + endpoint
}

func (m *MockProvider) GetWebexHeaders() map[string]string {
	if m.Headers == nil {
		return map[string]string{
			"Authorization": "Bearer " + m.Token,
			"Accept":        "application/json",
		}
	}
	return m.Headers
}

func (m *MockProvider) GetWebexJSONHeaders() map[string]string {
	if m.JSONHeaders == nil {
		headers := m.GetWebexHeaders()
		headers["Content-Type"] = "application/json"
		return headers
	}
	return m.JSONHeaders
}

func (m *MockProvider) GetUseFastHTTP() bool {
	return m.UseFastHTTP
}