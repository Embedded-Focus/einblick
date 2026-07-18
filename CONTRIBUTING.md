# Contributing

Einblick is built as a single Go command with provider-neutral core packages and small adapters around external services.

## Prerequisites

- Go 1.24 or newer
- `golangci-lint`

## Local Checks

Run these before opening a change:

```console
gofmt -w .
go test ./...
go test -race ./...
go vet ./...
golangci-lint run
CGO_ENABLED=0 go build -trimpath ./cmd/einblick
```

## Workflow

Keep changes small and evidence-oriented. For new metrics, include:

- stable metric ID
- human-readable title
- purpose
- population
- observation window
- exclusions
- calculation
- interpretation caveats
- tests for the calculation and any provider mapping it relies on

## Architecture Boundaries

- CLI code parses commands and maps errors to exit codes; it does not calculate metrics.
- Application code orchestrates use cases and depends on interfaces.
- Metric code consumes provider-neutral data and does not import CLI or GitHub packages.
- Provider adapters map external API responses into `internal/forge` types before metrics run.
- Tests for provider behavior should use `httptest.Server`, not live GitHub calls.
