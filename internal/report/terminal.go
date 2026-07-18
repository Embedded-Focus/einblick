package report

import (
	"fmt"
	"io"
	"strings"
)

// RenderTerminal writes compact human-readable analysis output.
func RenderTerminal(w io.Writer, r Report) error {
	if _, err := fmt.Fprintf(w, "Einblick: %s\n\n", r.Repository.Ref.String()); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "Repository"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  Archived:          %s\n", yesNo(r.Repository.Archived)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  Fork:              %s\n", yesNo(r.Repository.Fork)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  Default branch:    %s\n\n", emptyValue(r.Repository.DefaultBranch)); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "Contributor signals"); err != nil {
		return err
	}
	for _, metric := range r.Metrics {
		if _, err := fmt.Fprintf(w, "  %s: %v\n", metric.Definition.Title, metric.Value); err != nil {
			return err
		}
		for _, note := range metric.Notes {
			if _, err := fmt.Fprintf(w, "    Note: %s\n", note); err != nil {
				return err
			}
		}
	}
	for _, warning := range r.Warnings {
		if _, err := fmt.Fprintf(w, "Warning: %s\n", warning); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(w, "\nNote: Individual measurements are not an overall health judgment.")
	return err
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func emptyValue(value string) string {
	if strings.TrimSpace(value) == "" {
		return "unknown"
	}
	return value
}
