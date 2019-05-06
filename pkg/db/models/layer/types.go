package layer

// Model represents an entry in the database
type Model struct {
	Digest   string `db:"digest"`
	RepoName string `db:"repo_name"`
	OrgName  string `db:"org_name"`
}
