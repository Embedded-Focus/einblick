package metrics

import (
	"context"

	"github.com/einblick/einblick/internal/forge"
)

// Definition documents a stable measurement.
type Definition struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
}

// Calculator computes one metric from provider-neutral data.
type Calculator interface {
	Definition() Definition
	Calculate(ctx context.Context, provider forge.Provider, repo forge.Repository) (Result, error)
}
