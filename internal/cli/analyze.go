package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/einblick/einblick/internal/app"
	"github.com/einblick/einblick/internal/buildinfo"
	"github.com/einblick/einblick/internal/forge"
	"github.com/einblick/einblick/internal/github"
	"github.com/einblick/einblick/internal/metrics"
	"github.com/einblick/einblick/internal/report"
)

func runAnalyze(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) int {
	flags := newFlagSet("analyze", stderr)
	format := flags.String("format", "terminal", "output format: terminal or json")
	token := flags.String("token", "", "GitHub token; defaults to EINBLICK_GITHUB_TOKEN, then GITHUB_TOKEN")
	timeout := flags.Duration("timeout", 30*time.Second, "overall analysis timeout")
	flags.Usage = func() { writeAnalyzeHelp(stderr) }

	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return exitSuccess
		}
		return mapErrorToExitCode(err)
	}
	if flags.NArg() != 1 {
		_, _ = fmt.Fprintln(stderr, "error: analyze expects exactly one repository")
		writeAnalyzeHelp(stderr)
		return exitUsage
	}

	ref, err := forge.ParseRepositoryRef(flags.Arg(0))
	if err != nil {
		return reportRuntimeError(stderr, err)
	}
	if *format != "terminal" && *format != "json" {
		return reportRuntimeError(stderr, fmt.Errorf("%w: unsupported format %q", forge.ErrInvalidRepository, *format))
	}

	resolvedToken := *token
	if resolvedToken == "" {
		resolvedToken = tokenFromEnvironment()
	}
	if resolvedToken == "" && *format == "terminal" {
		_, _ = fmt.Fprintln(stderr, "warning: no GitHub token configured; unauthenticated rate limits apply")
	}

	ctx, cancel := context.WithTimeout(ctx, *timeout)
	defer cancel()

	provider := github.NewProvider(buildinfo.Version, github.WithToken(resolvedToken))
	analyzer := app.NewAnalyzer(provider, []metrics.Calculator{metrics.OpenPullRequests{}}, time.Now, buildinfo.Version)

	analysis, err := analyzer.Analyze(ctx, ref)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return reportRuntimeError(stderr, fmt.Errorf("analysis canceled or timed out: %w", err))
		}
		return reportRuntimeError(stderr, err)
	}

	if *format == "json" {
		err = report.RenderJSON(stdout, analysis)
	} else {
		err = report.RenderTerminal(stdout, analysis)
	}
	if err != nil {
		return reportRuntimeError(stderr, err)
	}
	return exitSuccess
}

func writeAnalyzeHelp(w io.Writer) {
	_, _ = fmt.Fprint(w, `Analyze a GitHub repository.

Usage:
  einblick analyze owner/repository [--format terminal|json] [--token token] [--timeout 30s]
  einblick analyze https://github.com/owner/repository --format json

Authentication:
  Credentials are discovered from --token, EINBLICK_GITHUB_TOKEN, then GITHUB_TOKEN.
  If none are present, Einblick uses unauthenticated GitHub API access.
`)
}
