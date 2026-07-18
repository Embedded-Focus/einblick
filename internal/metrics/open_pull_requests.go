package metrics

import (
	"context"
	"fmt"

	"github.com/einblick/einblick/internal/forge"
)

// OpenPullRequests counts currently open pull requests.
type OpenPullRequests struct{}

// Definition returns the metric identity and methodology summary.
func (OpenPullRequests) Definition() Definition {
	return Definition{
		ID:          "pull_requests.open.count",
		Title:       "Open pull requests",
		Description: "Count of currently open pull requests returned by the provider. This is not a health judgment by itself.",
		Unit:        "pull_requests",
	}
}

// Calculate counts open pull requests for the repository.
func (m OpenPullRequests) Calculate(ctx context.Context, provider forge.Provider, repo forge.Repository) (Result, error) {
	pulls, err := provider.ListPullRequests(ctx, repo.Ref, forge.PullRequestQuery{
		State: forge.PullRequestStateOpen,
		Limit: 100,
	})
	if err != nil {
		return Result{}, fmt.Errorf("list open pull requests: %w", err)
	}

	return Result{
		Definition: m.Definition(),
		Value:      len(pulls),
		Status:     StatusOK,
		Notes:      []string{"The count is capped at the first 100 open pull requests in this bootstrap implementation."},
	}, nil
}
