package org

// Model represents a single entry in the database
type Model struct {
	Name string `db:"name"`
}
