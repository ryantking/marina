package chunk

// Model represents an entry in the database
type Model struct {
	UUID       string `db:"uuid,omitempty"`
	RangeStart int64  `db:"range_start"`
	RangeEnd   int64  `db:"range_end"`
}
