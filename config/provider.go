package config

// Provider defines the interface for configuration access
type Provider interface {
	GetWebexURL(endpoint string) string
	GetWebexHeaders() map[string]string
	GetWebexJSONHeaders() map[string]string
	GetWebexToken() string
	GetWebexBaseURL() string
	GetUseFastHTTP() bool
}

// DefaultProvider implements Provider using environment variables
type DefaultProvider struct{}

func NewDefaultProvider() Provider {
	return &DefaultProvider{}
}

func (p *DefaultProvider) GetWebexURL(endpoint string) string {
	return GetWebexURL(endpoint)
}

func (p *DefaultProvider) GetWebexHeaders() map[string]string {
	return GetWebexHeaders()
}

func (p *DefaultProvider) GetWebexJSONHeaders() map[string]string {
	return GetWebexJSONHeaders()
}

func (p *DefaultProvider) GetWebexToken() string {
	return GetWebexToken()
}

func (p *DefaultProvider) GetWebexBaseURL() string {
	return GetWebexBaseURL()
}

func (p *DefaultProvider) GetUseFastHTTP() bool {
	return GetUseFastHTTP()
}
