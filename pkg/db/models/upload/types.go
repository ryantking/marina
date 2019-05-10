package upload

// Model represents an entry in the database
type Model struct {
	UUID string `db:"uuid,omitempty"`
	Done bool   `db:"done"`
}
