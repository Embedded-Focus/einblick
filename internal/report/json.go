package report

import (
	"encoding/json"
	"io"
	"time"
)

// RenderJSON writes a stable JSON representation of a report.
func RenderJSON(w io.Writer, r Report) error {
	output := jsonReport{
		SchemaVersion: r.SchemaVersion,
		ToolVersion:   r.ToolVersion,
		Repository:    r.Repository.Ref.String(),
		ObservedAt:    r.ObservedAt.UTC(),
		Metrics:       make([]jsonMetric, 0, len(r.Metrics)),
		Warnings:      r.Warnings,
	}
	for _, metric := range r.Metrics {
		output.Metrics = append(output.Metrics, jsonMetric{
			ID:          metric.Definition.ID,
			Title:       metric.Definition.Title,
			Description: metric.Definition.Description,
			Value:       metric.Value,
			Status:      string(metric.Status),
			Unit:        metric.Definition.Unit,
			Notes:       metric.Notes,
		})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

type jsonReport struct {
	SchemaVersion string       `json:"schema_version"`
	ToolVersion   string       `json:"tool_version"`
	Repository    string       `json:"repository"`
	ObservedAt    time.Time    `json:"observed_at"`
	Metrics       []jsonMetric `json:"metrics"`
	Warnings      []string     `json:"warnings"`
}

type jsonMetric struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Value       any      `json:"value"`
	Status      string   `json:"status"`
	Unit        string   `json:"unit"`
	Notes       []string `json:"notes,omitempty"`
}
