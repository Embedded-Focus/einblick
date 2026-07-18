package forge

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// RepositoryRef identifies a repository within an owner namespace.
type RepositoryRef struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

// String returns the canonical owner/repository identifier.
func (r RepositoryRef) String() string {
	return r.Owner + "/" + r.Name
}

// Validate checks that the repository reference is non-empty and unambiguous.
func (r RepositoryRef) Validate() error {
	if !validPathPart(r.Owner) || !validPathPart(r.Name) {
		return fmt.Errorf("%w: expected owner/repository", ErrInvalidRepository)
	}
	return nil
}

// ParseRepositoryRef accepts owner/repository and common GitHub repository URLs.
func ParseRepositoryRef(input string) (RepositoryRef, error) {
	value := strings.TrimSpace(input)
	if value == "" {
		return RepositoryRef{}, fmt.Errorf("%w: repository is required", ErrInvalidRepository)
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		parsed, err := url.Parse(value)
		if err != nil {
			return RepositoryRef{}, fmt.Errorf("%w: %q", ErrInvalidRepository, input)
		}
		parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		if len(parts) < 2 {
			return RepositoryRef{}, fmt.Errorf("%w: expected GitHub URL with owner and repository", ErrInvalidRepository)
		}
		value = parts[0] + "/" + strings.TrimSuffix(parts[1], ".git")
	}

	parts := strings.Split(value, "/")
	if len(parts) != 2 {
		return RepositoryRef{}, fmt.Errorf("%w: expected owner/repository", ErrInvalidRepository)
	}
	ref := RepositoryRef{Owner: parts[0], Name: strings.TrimSuffix(parts[1], ".git")}
	if err := ref.Validate(); err != nil {
		return RepositoryRef{}, err
	}
	return ref, nil
}

func validPathPart(value string) bool {
	if value == "" || value == "." || value == ".." {
		return false
	}
	return !strings.ContainsAny(value, "/ \t\n\r")
}

// Repository is the provider-neutral repository metadata used by reports.
type Repository struct {
	Ref           RepositoryRef `json:"ref"`
	DefaultBranch string        `json:"default_branch"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	Archived      bool          `json:"archived"`
	Fork          bool          `json:"fork"`
}
