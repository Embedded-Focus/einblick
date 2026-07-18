package metrics

import (
	"context"
	"testing"

	"github.com/einblick/einblick/internal/forge"
)

func TestOpenPullRequests(t *testing.T) {
	t.Parallel()

	provider := fakeProvider{
		pulls: []forge.PullRequest{{Number: 1}, {Number: 2}},
	}
	metric := OpenPullRequests{}
	result, err := metric.Calculate(context.Background(), provider, forge.Repository{Ref: forge.RepositoryRef{Owner: "o", Name: "r"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Definition.ID != "pull_requests.open.count" {
		t.Fatalf("unexpected metric ID: %s", result.Definition.ID)
	}
	if result.Value != 2 {
		t.Fatalf("got value %v, want 2", result.Value)
	}
}

type fakeProvider struct {
	forge.Provider
	pulls []forge.PullRequest
}

func (f fakeProvider) ListPullRequests(ctx context.Context, ref forge.RepositoryRef, query forge.PullRequestQuery) ([]forge.PullRequest, error) {
	return f.pulls, nil
}
