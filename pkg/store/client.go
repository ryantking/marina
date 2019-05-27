package store

import (
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/store/s3"
)

// Client is used for the storing and retrieval of files
type Client interface {
	// Get retrieves an object from the bucket
	Get(path string, start, end int64) (io.ReadCloser, error)

	// Put saves a file to the store
	Put(path string, r io.Reader, sz int32) (int32, error)

	// Remove deletes an object from the store
	Remove(path string) error

	// Merge combines multiple objects into one
	Merge(toPath string, fromPaths ...string) error
}

var (
	client  Client
	clientL sync.Mutex
)

func getClient() (Client, error) {
	clientL.Lock()
	defer clientL.Unlock()

	if client != nil {
		return client, nil
	}

	c, err := s3.New()
	if err != nil {
		return nil, errors.Wrap(err, "error creating S3 client")
	}
	client = c
	return client, nil
}
