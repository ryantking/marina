package tag

// Model represents an entry in the database
type Model struct {
	Name        string `db:"name"`
	RepoName    string `db:"repo_name"`
	OrgName     string `db:"org_name"`
	ImageDigest string `db:"image_digest"`
}
