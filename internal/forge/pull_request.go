package forge

import "time"

// PullRequest is the provider-neutral pull request shape used by metrics.
type PullRequest struct {
	Number       int           `json:"number"`
	Author       Actor         `json:"author"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	ClosedAt     *time.Time    `json:"closed_at,omitempty"`
	MergedAt     *time.Time    `json:"merged_at,omitempty"`
	Draft        bool          `json:"draft"`
	AuthorRole   AuthorRole    `json:"author_role"`
	ReviewEvents []ReviewEvent `json:"review_events,omitempty"`
}

// ReviewEvent is reserved for future review-latency metrics.
type ReviewEvent struct {
	Reviewer  Actor     `json:"reviewer"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}
