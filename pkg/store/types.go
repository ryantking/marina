package store

import "io"

// Client is used for the storing and retrieval of files
type Client interface {
	// Get retrieves an object from the bucket
	Get(path string) (io.ReadCloser, error)

	// Put saves a file to the store
	Put(path string, r io.Reader, sz int64) (int64, error)

	// Remove deletes an object from the store
	Remove(path string) error
}
