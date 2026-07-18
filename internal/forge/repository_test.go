package forge

import "testing"

func TestParseRepositoryRef(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    RepositoryRef
		wantErr bool
	}{
		{name: "canonical", input: "owner/repo", want: RepositoryRef{Owner: "owner", Name: "repo"}},
		{name: "url", input: "https://github.com/owner/repo", want: RepositoryRef{Owner: "owner", Name: "repo"}},
		{name: "git suffix", input: "https://github.com/owner/repo.git", want: RepositoryRef{Owner: "owner", Name: "repo"}},
		{name: "empty", input: "", wantErr: true},
		{name: "too many parts", input: "a/b/c", wantErr: true},
		{name: "space", input: "owner/my repo", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseRepositoryRef(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
