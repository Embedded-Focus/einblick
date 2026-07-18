package metrics

// Status reports whether a metric is complete enough to interpret.
type Status string

const (
	StatusOK         Status = "ok"
	StatusIncomplete Status = "incomplete"
)

// Result is a computed metric value plus caveats.
type Result struct {
	Definition Definition `json:"definition"`
	Value      any        `json:"value"`
	Status     Status     `json:"status"`
	Notes      []string   `json:"notes,omitempty"`
}
