package forge

import (
	"context"
	"errors"
)

var (
	ErrInvalidRepository = errors.New("invalid repository identifier")
	ErrAuthentication    = errors.New("authentication or authorization failed")
	ErrNotFound          = errors.New("repository not found or inaccessible")
	ErrRateLimited       = errors.New("provider rate limit exceeded")
	ErrTransient         = errors.New("transient provider failure")
	ErrMalformedResponse = errors.New("malformed provider response")
)

// Provider retrieves provider-neutral source forge data.
type Provider interface {
	GetRepository(ctx context.Context, ref RepositoryRef) (Repository, error)
	ListPullRequests(ctx context.Context, ref RepositoryRef, query PullRequestQuery) ([]PullRequest, error)
}

// PullRequestState selects the pull request population returned by a provider.
type PullRequestState string

const (
	PullRequestStateOpen   PullRequestState = "open"
	PullRequestStateClosed PullRequestState = "closed"
	PullRequestStateAll    PullRequestState = "all"
)

// PullRequestQuery contains the current pull request filters needed by metrics.
type PullRequestQuery struct {
	State PullRequestState
	Limit int
}
