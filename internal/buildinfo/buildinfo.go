package buildinfo

import "fmt"

var (
	// Version is injected by release builds with -ldflags.
	Version = "dev"
	// Commit is injected by release builds with -ldflags.
	Commit = "unknown"
	// Built is injected by release builds with -ldflags.
	Built = "unknown"
)

// String returns the human-readable build identity.
func String() string {
	return fmt.Sprintf("Einblick %s (commit %s, built %s)", Version, Commit, Built)
}
