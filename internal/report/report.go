package report

import (
	"time"

	"github.com/einblick/einblick/internal/forge"
	"github.com/einblick/einblick/internal/metrics"
)

// Report is the presentation-neutral result of an analysis.
type Report struct {
	SchemaVersion string           `json:"schema_version"`
	ToolVersion   string           `json:"tool_version"`
	Repository    forge.Repository `json:"repository"`
	ObservedAt    time.Time        `json:"observed_at"`
	Metrics       []metrics.Result `json:"metrics"`
	Warnings      []string         `json:"warnings"`
}
