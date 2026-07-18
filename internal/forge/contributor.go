package forge

// Actor identifies a person or automation account on a forge.
type Actor struct {
	Login string `json:"login"`
	Bot   bool   `json:"bot"`
}

// AuthorRole describes the author's relationship to the repository.
type AuthorRole string

const (
	AuthorRoleUnknown      AuthorRole = "unknown"
	AuthorRoleOwner        AuthorRole = "owner"
	AuthorRoleMember       AuthorRole = "member"
	AuthorRoleCollaborator AuthorRole = "collaborator"
	AuthorRoleContributor  AuthorRole = "contributor"
	AuthorRoleFirstTimer   AuthorRole = "first_timer"
)
