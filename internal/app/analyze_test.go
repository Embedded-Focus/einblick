package app

import (
	"context"
	"testing"
	"time"

	"github.com/einblick/einblick/internal/forge"
	"github.com/einblick/einblick/internal/metrics"
)

func TestAnalyzerAnalyze(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	provider := analyzerFakeProvider{
		repo: forge.Repository{
			Ref:           forge.RepositoryRef{Owner: "owner", Name: "repo"},
			DefaultBranch: "main",
		},
		pulls: []forge.PullRequest{{Number: 1}},
	}
	analyzer := NewAnalyzer(provider, []metrics.Calculator{metrics.OpenPullRequests{}}, func() time.Time { return now }, "dev")

	got, err := analyzer.Analyze(context.Background(), forge.RepositoryRef{Owner: "owner", Name: "repo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Repository.Ref.String() != "owner/repo" {
		t.Fatalf("unexpected repository: %#v", got.Repository.Ref)
	}
	if !got.ObservedAt.Equal(now) {
		t.Fatalf("got observed_at %s, want %s", got.ObservedAt, now)
	}
	if len(got.Metrics) != 1 || got.Metrics[0].Value != 1 {
		t.Fatalf("unexpected metrics: %#v", got.Metrics)
	}
}

type analyzerFakeProvider struct {
	repo  forge.Repository
	pulls []forge.PullRequest
}

func (f analyzerFakeProvider) GetRepository(ctx context.Context, ref forge.RepositoryRef) (forge.Repository, error) {
	return f.repo, nil
}

func (f analyzerFakeProvider) ListPullRequests(ctx context.Context, ref forge.RepositoryRef, query forge.PullRequestQuery) ([]forge.PullRequest, error) {
	return f.pulls, nil
}
