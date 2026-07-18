package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/einblick/einblick/internal/forge"
)

// Provider retrieves repository data from GitHub's REST API.
type Provider struct {
	client    *http.Client
	baseURL   string
	token     string
	userAgent string
}

// NewProvider constructs a GitHub REST provider.
func NewProvider(version string, options ...ProviderOption) *Provider {
	p := &Provider{
		client:    http.DefaultClient,
		baseURL:   defaultBaseURL,
		userAgent: "einblick/" + version,
	}
	for _, option := range options {
		option(p)
	}
	return p
}

// GetRepository retrieves provider-neutral repository metadata.
func (p *Provider) GetRepository(ctx context.Context, ref forge.RepositoryRef) (forge.Repository, error) {
	if err := ref.Validate(); err != nil {
		return forge.Repository{}, err
	}

	req, err := p.newRequest(ctx, http.MethodGet, "/repos/"+url.PathEscape(ref.Owner)+"/"+url.PathEscape(ref.Name), nil)
	if err != nil {
		return forge.Repository{}, err
	}

	var response repositoryResponse
	if err := p.doJSON(req, &response); err != nil {
		return forge.Repository{}, err
	}
	return mapRepository(response), nil
}

// ListPullRequests retrieves provider-neutral pull requests.
func (p *Provider) ListPullRequests(ctx context.Context, ref forge.RepositoryRef, query forge.PullRequestQuery) ([]forge.PullRequest, error) {
	if err := ref.Validate(); err != nil {
		return nil, err
	}

	limit := query.Limit
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	state := string(query.State)
	if state == "" {
		state = string(forge.PullRequestStateOpen)
	}

	path := "/repos/" + url.PathEscape(ref.Owner) + "/" + url.PathEscape(ref.Name) + "/pulls"
	req, err := p.newRequest(ctx, http.MethodGet, path+"?state="+url.QueryEscape(state)+"&per_page="+strconv.Itoa(limit), nil)
	if err != nil {
		return nil, err
	}

	var response []pullRequestResponse
	if err := p.doJSON(req, &response); err != nil {
		return nil, err
	}
	return mapPullRequests(response), nil
}

func (p *Provider) newRequest(ctx context.Context, method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, p.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("%w: build GitHub request", forge.ErrInvalidRepository)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", p.userAgent)
	if p.token != "" {
		req.Header.Set("Authorization", "Bearer "+p.token)
	}
	return req, nil
}

func (p *Provider) doJSON(req *http.Request, target any) error {
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", forge.ErrTransient, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusUnauthorized, http.StatusForbidden:
		if resp.Header.Get("X-RateLimit-Remaining") == "0" {
			return forge.ErrRateLimited
		}
		return forge.ErrAuthentication
	case http.StatusNotFound:
		return forge.ErrNotFound
	case http.StatusTooManyRequests:
		return forge.ErrRateLimited
	default:
		if resp.StatusCode >= 500 {
			return fmt.Errorf("%w: GitHub status %d", forge.ErrTransient, resp.StatusCode)
		}
		return fmt.Errorf("%w: GitHub status %d", forge.ErrMalformedResponse, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("%w: decode GitHub response: %v", forge.ErrMalformedResponse, err)
	}
	return nil
}
