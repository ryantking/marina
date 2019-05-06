package upload

// Model represents an entry in the database
type Model struct {
	UUID uint64 `db:"uuid,omitempty"`
	Done bool   `db:"done"`
}
