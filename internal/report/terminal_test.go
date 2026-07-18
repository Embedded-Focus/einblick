package report

import (
	"strings"
	"testing"
	"time"

	"github.com/einblick/einblick/internal/forge"
	"github.com/einblick/einblick/internal/metrics"
)

func TestRenderTerminal(t *testing.T) {
	t.Parallel()

	var out strings.Builder
	err := RenderTerminal(&out, Report{
		SchemaVersion: "1",
		ToolVersion:   "dev",
		Repository: forge.Repository{
			Ref:           forge.RepositoryRef{Owner: "owner", Name: "repo"},
			DefaultBranch: "main",
		},
		ObservedAt: time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC),
		Metrics: []metrics.Result{{
			Definition: metrics.Definition{Title: "Open pull requests"},
			Value:      3,
			Status:     metrics.StatusOK,
		}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	for _, want := range []string{"Einblick: owner/repo", "Default branch:    main", "Open pull requests: 3"} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q:\n%s", want, got)
		}
	}
}
