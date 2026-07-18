package github

import (
	"strings"
	"time"

	"github.com/einblick/einblick/internal/forge"
)

type repositoryResponse struct {
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	Name          string    `json:"name"`
	DefaultBranch string    `json:"default_branch"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Archived      bool      `json:"archived"`
	Fork          bool      `json:"fork"`
}

type pullRequestResponse struct {
	Number                int        `json:"number"`
	User                  user       `json:"user"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	ClosedAt              *time.Time `json:"closed_at"`
	MergedAt              *time.Time `json:"merged_at"`
	Draft                 bool       `json:"draft"`
	AuthorAssociation     string     `json:"author_association"`
}

type user struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

func mapRepository(response repositoryResponse) forge.Repository {
	return forge.Repository{
		Ref: forge.RepositoryRef{
			Owner: response.Owner.Login,
			Name:  response.Name,
		},
		DefaultBranch: response.DefaultBranch,
		CreatedAt:     response.CreatedAt.UTC(),
		UpdatedAt:     response.UpdatedAt.UTC(),
		Archived:      response.Archived,
		Fork:          response.Fork,
	}
}

func mapPullRequests(response []pullRequestResponse) []forge.PullRequest {
	pulls := make([]forge.PullRequest, 0, len(response))
	for _, item := range response {
		pulls = append(pulls, forge.PullRequest{
			Number: item.Number,
			Author: forge.Actor{
				Login: item.User.Login,
				Bot:   strings.EqualFold(item.User.Type, "Bot") || strings.HasSuffix(strings.ToLower(item.User.Login), "[bot]"),
			},
			CreatedAt:  item.CreatedAt.UTC(),
			UpdatedAt:  item.UpdatedAt.UTC(),
			ClosedAt:   utcPointer(item.ClosedAt),
			MergedAt:   utcPointer(item.MergedAt),
			Draft:      item.Draft,
			AuthorRole: mapAuthorRole(item.AuthorAssociation),
		})
	}
	return pulls
}

func utcPointer(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	utc := value.UTC()
	return &utc
}

func mapAuthorRole(value string) forge.AuthorRole {
	switch strings.ToUpper(value) {
	case "OWNER":
		return forge.AuthorRoleOwner
	case "MEMBER":
		return forge.AuthorRoleMember
	case "COLLABORATOR":
		return forge.AuthorRoleCollaborator
	case "CONTRIBUTOR":
		return forge.AuthorRoleContributor
	case "FIRST_TIMER", "FIRST_TIME_CONTRIBUTOR":
		return forge.AuthorRoleFirstTimer
	default:
		return forge.AuthorRoleUnknown
	}
}
