package repo

// Model represents an entry in the database
type Model struct {
	Name    string `db:"name"`
	OrgName string `db:"org_name"`
}
