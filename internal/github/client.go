package github

import (
	"net/http"
	"strings"
)

const defaultBaseURL = "https://api.github.com"

// ProviderOption customizes the GitHub provider.
type ProviderOption func(*Provider)

// WithHTTPClient sets the HTTP client used for API requests.
func WithHTTPClient(client *http.Client) ProviderOption {
	return func(p *Provider) {
		if client != nil {
			p.client = client
		}
	}
}

// WithBaseURL sets the API base URL, primarily for tests.
func WithBaseURL(baseURL string) ProviderOption {
	return func(p *Provider) {
		p.baseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithToken sets the GitHub API token. The token is never logged by this package.
func WithToken(token string) ProviderOption {
	return func(p *Provider) {
		p.token = strings.TrimSpace(token)
	}
}
