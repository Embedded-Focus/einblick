package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/einblick/einblick/internal/buildinfo"
	"github.com/einblick/einblick/internal/forge"
)

const (
	exitSuccess   = 0
	exitFailure   = 1
	exitUsage     = 2
	exitAuth      = 3
	exitNotFound  = 4
	exitRateLimit = 5
)

// Execute runs the CLI and returns a process exit code.
func Execute(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		writeRootHelp(stdout)
		return exitSuccess
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeRootHelp(stdout)
		return exitSuccess
	case "analyze":
		return runAnalyze(ctx, args[1:], stdout, stderr)
	case "compare":
		return runCompare(args[1:], stdout, stderr)
	case "version":
		_, err := fmt.Fprintln(stdout, buildinfo.String())
		if err != nil {
			return exitFailure
		}
		return exitSuccess
	default:
		_, _ = fmt.Fprintf(stderr, "unknown command %q\n\n", args[0])
		writeRootHelp(stderr)
		return exitUsage
	}
}

func newFlagSet(name string, stderr io.Writer) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(stderr)
	return flags
}

func mapErrorToExitCode(err error) int {
	switch {
	case err == nil:
		return exitSuccess
	case errors.Is(err, forge.ErrInvalidRepository):
		return exitUsage
	case errors.Is(err, forge.ErrAuthentication):
		return exitAuth
	case errors.Is(err, forge.ErrNotFound):
		return exitNotFound
	case errors.Is(err, forge.ErrRateLimited):
		return exitRateLimit
	default:
		return exitFailure
	}
}

func tokenFromEnvironment() string {
	if token := os.Getenv("EINBLICK_GITHUB_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITHUB_TOKEN")
}

func writeRootHelp(w io.Writer) {
	_, _ = fmt.Fprint(w, `Einblick evaluates contributor-facing project signals.

Usage:
  einblick analyze owner/repository [--format terminal|json] [--token token] [--timeout 30s]
  einblick compare owner/repository owner/repository
  einblick version

Exit codes:
  0 success
  1 runtime failure
  2 invalid command-line input
  3 authentication or authorization failure
  4 repository not found or inaccessible
  5 provider rate limit exceeded
`)
}

func reportRuntimeError(stderr io.Writer, err error) int {
	code := mapErrorToExitCode(err)
	if code != exitSuccess {
		_, _ = fmt.Fprintf(stderr, "error: %v\n", err)
	}
	return code
}
