package config

// MockProvider is a test implementation of the Provider interface
type MockProvider struct {
	Token       string
	BaseURL     string
	Headers     map[string]string
	JSONHeaders map[string]string
	UseFastHTTP bool
}

func (m *MockProvider) GetWebexToken() (string, error) {
	return m.Token, nil
}

func (m *MockProvider) GetWebexBaseURL() (string, error) {
	return m.BaseURL, nil
}

func (m *MockProvider) GetWebexURL(endpoint string) (string, error) {
	baseURL, err := m.GetWebexBaseURL()
	if err != nil {
		return "", err
	}
	return baseURL + endpoint, nil
}

func (m *MockProvider) GetWebexHeaders() (map[string]string, error) {
	if m.Headers == nil {
		token, err := m.GetWebexToken()
		if err != nil {
			return nil, err
		}
		return map[string]string{
			"Authorization": "Bearer " + token,
			"Accept":        "application/json",
		}, nil
	}
	return m.Headers, nil
}

func (m *MockProvider) GetWebexJSONHeaders() (map[string]string, error) {
	if m.JSONHeaders == nil {
		headers, err := m.GetWebexHeaders()
		if err != nil {
			return nil, err
		}
		headers["Content-Type"] = "application/json"
		return headers, nil
	}
	return m.JSONHeaders, nil
}

func (m *MockProvider) GetUseFastHTTP() bool {
	return m.UseFastHTTP
}