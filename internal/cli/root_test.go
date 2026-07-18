package cli

import (
	"context"
	"strings"
	"testing"
)

func TestExecuteHelpAndVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     []string
		wantCode int
		wantOut  string
	}{
		{name: "root help", args: []string{"--help"}, wantCode: exitSuccess, wantOut: "Usage:"},
		{name: "analyze help", args: []string{"analyze", "--help"}, wantCode: exitSuccess, wantOut: "Analyze a GitHub repository."},
		{name: "version", args: []string{"version"}, wantCode: exitSuccess, wantOut: "Einblick dev"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var stdout strings.Builder
			var stderr strings.Builder
			code := Execute(context.Background(), tt.args, &stdout, &stderr)
			got := stdout.String() + stderr.String()
			if code != tt.wantCode {
				t.Fatalf("got code %d, want %d; output:\n%s", code, tt.wantCode, got)
			}
			if !strings.Contains(got, tt.wantOut) {
				t.Fatalf("output missing %q:\n%s", tt.wantOut, got)
			}
		})
	}
}
