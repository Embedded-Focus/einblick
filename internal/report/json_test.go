package report

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/einblick/einblick/internal/forge"
	"github.com/einblick/einblick/internal/metrics"
)

func TestRenderJSON(t *testing.T) {
	t.Parallel()

	var out strings.Builder
	err := RenderJSON(&out, Report{
		SchemaVersion: "1",
		ToolVersion:   "dev",
		Repository: forge.Repository{
			Ref: forge.RepositoryRef{Owner: "owner", Name: "repo"},
		},
		ObservedAt: time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC),
		Metrics: []metrics.Result{{
			Definition: metrics.Definition{
				ID:    "pull_requests.open.count",
				Title: "Open pull requests",
				Unit:  "pull_requests",
			},
			Value:  4,
			Status: metrics.StatusOK,
		}},
		Warnings: []string{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got struct {
		Repository string `json:"repository"`
		Metrics    []struct {
			ID    string `json:"id"`
			Value int    `json:"value"`
			Unit  string `json:"unit"`
		} `json:"metrics"`
	}
	if err := json.Unmarshal([]byte(out.String()), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Repository != "owner/repo" {
		t.Fatalf("got repository %q, want owner/repo", got.Repository)
	}
	if len(got.Metrics) != 1 || got.Metrics[0].ID != "pull_requests.open.count" || got.Metrics[0].Value != 4 || got.Metrics[0].Unit != "pull_requests" {
		t.Fatalf("unexpected metrics: %#v", got.Metrics)
	}
}
