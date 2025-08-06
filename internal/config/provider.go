package config

// Provider defines the interface for configuration access
type Provider interface {
	GetWebexURL(endpoint string) (string, error)
	GetWebexHeaders() (map[string]string, error)
	GetWebexJSONHeaders() (map[string]string, error)
	GetWebexToken() (string, error)
	GetWebexBaseURL() (string, error)
	GetUseFastHTTP() bool
}

// DefaultProvider implements Provider using environment variables
type DefaultProvider struct{}

func NewDefaultProvider() Provider {
	return &DefaultProvider{}
}

func (p *DefaultProvider) GetWebexURL(endpoint string) (string, error) {
	return GetWebexURL(endpoint)
}

func (p *DefaultProvider) GetWebexHeaders() (map[string]string, error) {
	return GetWebexHeaders()
}

func (p *DefaultProvider) GetWebexJSONHeaders() (map[string]string, error) {
	return GetWebexJSONHeaders()
}

func (p *DefaultProvider) GetWebexToken() (string, error) {
	return GetWebexToken()
}

func (p *DefaultProvider) GetWebexBaseURL() (string, error) {
	return GetWebexBaseURL()
}

func (p *DefaultProvider) GetUseFastHTTP() bool {
	return GetUseFastHTTP()
}