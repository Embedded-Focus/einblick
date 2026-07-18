package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/einblick/einblick/internal/app"
)

func runCompare(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := newFlagSet("compare", stderr)
	flags.Usage = func() { writeCompareHelp(stderr) }
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return exitSuccess
		}
		return mapErrorToExitCode(err)
	}
	if flags.NArg() != 2 {
		_, _ = fmt.Fprintln(stderr, "error: compare expects exactly two repositories")
		writeCompareHelp(stderr)
		return exitUsage
	}
	_, _ = fmt.Fprintf(stderr, "error: %v\n", app.ErrCompareNotImplemented)
	return exitFailure
}

func writeCompareHelp(w io.Writer) {
	_, _ = fmt.Fprint(w, `Compare two repositories.

Usage:
  einblick compare owner/repository owner/repository

Status:
  This workflow is intentionally deferred after the initial executable skeleton.
`)
}
