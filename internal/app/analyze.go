package app

import (
	"context"
	"fmt"
	"time"

	"github.com/einblick/einblick/internal/forge"
	"github.com/einblick/einblick/internal/metrics"
	"github.com/einblick/einblick/internal/report"
)

// Analyzer orchestrates repository retrieval, metric calculation and report assembly.
type Analyzer struct {
	provider forge.Provider
	metrics  []metrics.Calculator
	now      func() time.Time
	version  string
}

// NewAnalyzer constructs an Analyzer with explicit dependencies.
func NewAnalyzer(provider forge.Provider, calculators []metrics.Calculator, now func() time.Time, version string) *Analyzer {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Analyzer{
		provider: provider,
		metrics:  calculators,
		now:      now,
		version:  version,
	}
}

// Analyze retrieves repository evidence and computes configured metrics.
func (a *Analyzer) Analyze(ctx context.Context, ref forge.RepositoryRef) (report.Report, error) {
	if err := ref.Validate(); err != nil {
		return report.Report{}, err
	}

	repo, err := a.provider.GetRepository(ctx, ref)
	if err != nil {
		return report.Report{}, fmt.Errorf("get repository %s: %w", ref.String(), err)
	}

	results := make([]metrics.Result, 0, len(a.metrics))
	warnings := make([]string, 0)
	for _, calculator := range a.metrics {
		result, err := calculator.Calculate(ctx, a.provider, repo)
		if err != nil {
			return report.Report{}, fmt.Errorf("calculate metric %s: %w", calculator.Definition().ID, err)
		}
		if result.Status == metrics.StatusIncomplete {
			warnings = append(warnings, result.Notes...)
		}
		results = append(results, result)
	}

	return report.Report{
		SchemaVersion: "1",
		ToolVersion:   a.version,
		Repository:    repo,
		ObservedAt:    a.now().UTC(),
		Metrics:       results,
		Warnings:      warnings,
	}, nil
}
