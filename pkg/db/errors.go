package db

import "fmt"

// ErrCollectionNotFound is thrown when a collection does not exist in the database
type ErrCollectionNotFound struct {
	colName string
}

// Error returns the string representation of the error
func (err *ErrCollectionNotFound) Error() string {
	return fmt.Sprintf("collection '%s' does not exist", err.colName)
}
