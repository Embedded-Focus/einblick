package github

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/einblick/einblick/internal/forge"
)

func TestProviderMapsRepositoryAndPullRequests(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "einblick/test" {
			http.Error(w, "missing user agent", http.StatusBadRequest)
			return
		}
		switch r.URL.Path {
		case "/repos/owner/repo":
			_, _ = w.Write([]byte(`{
				"owner": {"login": "owner"},
				"name": "repo",
				"default_branch": "main",
				"created_at": "2025-01-01T00:00:00Z",
				"updated_at": "2026-01-01T00:00:00Z",
				"archived": false,
				"fork": false
			}`))
		case "/repos/owner/repo/pulls":
			if r.URL.Query().Get("state") != "open" {
				http.Error(w, "unexpected state query", http.StatusBadRequest)
				return
			}
			_, _ = w.Write([]byte(`[{
				"number": 7,
				"user": {"login": "contributor", "type": "User"},
				"created_at": "2026-01-02T00:00:00Z",
				"updated_at": "2026-01-03T00:00:00Z",
				"draft": true,
				"author_association": "CONTRIBUTOR"
			}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := NewProvider("test", WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	ref := forge.RepositoryRef{Owner: "owner", Name: "repo"}

	repo, err := provider.GetRepository(context.Background(), ref)
	if err != nil {
		t.Fatalf("GetRepository: %v", err)
	}
	if repo.Ref != ref || repo.DefaultBranch != "main" {
		t.Fatalf("unexpected repository: %#v", repo)
	}

	pulls, err := provider.ListPullRequests(context.Background(), ref, forge.PullRequestQuery{State: forge.PullRequestStateOpen})
	if err != nil {
		t.Fatalf("ListPullRequests: %v", err)
	}
	if len(pulls) != 1 || pulls[0].Number != 7 || pulls[0].AuthorRole != forge.AuthorRoleContributor || !pulls[0].Draft {
		t.Fatalf("unexpected pulls: %#v", pulls)
	}
}

func TestProviderClassifiesErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code int
		want error
	}{
		{name: "auth", code: http.StatusUnauthorized, want: forge.ErrAuthentication},
		{name: "not found", code: http.StatusNotFound, want: forge.ErrNotFound},
		{name: "rate limited", code: http.StatusTooManyRequests, want: forge.ErrRateLimited},
		{name: "server", code: http.StatusBadGateway, want: forge.ErrTransient},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.code)
			}))
			defer server.Close()

			provider := NewProvider("test", WithBaseURL(server.URL), WithHTTPClient(server.Client()))
			_, err := provider.GetRepository(context.Background(), forge.RepositoryRef{Owner: "owner", Name: "repo"})
			if !errors.Is(err, tt.want) {
				t.Fatalf("got error %v, want %v", err, tt.want)
			}
		})
	}
}
