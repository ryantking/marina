package store

import (
	"fmt"
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db/models/upload/chunk"
	"github.com/ryantking/marina/pkg/store/s3"
)

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

// CreateUpload creates a new upload from the given reader
func UploadChunk(uuid string, r io.Reader, sz, start, end int64) (int64, error) {
	client, err := getClient()
	if err != nil {
		return 0, err
	}

	path := fmt.Sprintf("uploads/%s/%d_%d.tar.gz", uuid, start, end)
	if end == 0 {
		path = fmt.Sprintf("uploads/%s.tar.gz", uuid)
	}
	n, err := client.Put(path, r, sz)
	if err != nil {
		return 0, errors.Wrap(err, "error uploading object")
	}
	return n, nil
}

// FinishUpload finalizes an upload
func FinishUpload(digest, uuid, repoName, orgName string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	chunks, err := chunk.GetAll(uuid)
	if err != nil {
		return errors.Wrap(err, "error getting upload chunks from database")
	}
	r, n, err := mergeChunks(chunks)
	if err != nil {
		return errors.Wrap(err, "error getting upload reader")
	}

	loc := fmt.Sprintf("blobs/%s/%s/%s.tar.gz", orgName, repoName, digest)
	_, err = client.Put(loc, r, n)
	if err != nil {
		return err
	}

	err = deleteChunks(chunks)
	if err != nil {
		return errors.Wrap(err, "error deleting chunks")
	}

	return nil
}

func mergeChunks(chunks []*chunk.Model) (io.Reader, int64, error) {
	client, err := getClient()
	if err != nil {
		return nil, 0, err
	}

	if len(chunks) == 1 {
		loc := fmt.Sprintf("uploads/%s.tar.gz", chunks[0].UUID)
		obj, err := client.Get(loc)
		if err != nil {
			return nil, 0, err
		}
		return obj, chunks[0].RangeEnd - chunks[0].RangeStart + 1, nil
	}

	var sz int64
	readers := make([]io.Reader, len(chunks))
	for i, chunk := range chunks {
		loc := fmt.Sprintf("uploads/%s/%d_%d.tar.gz", chunk.UUID, chunk.RangeStart, chunk.RangeEnd)
		obj, err := client.Get(loc)
		if err != nil {
			return nil, 0, err
		}
		readers[i] = obj
		sz += chunk.RangeEnd - chunk.RangeStart + 1
	}

	return io.MultiReader(readers...), sz, nil
}

func deleteChunks(chunks []*chunk.Model) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	if len(chunks) == 1 {
		loc := fmt.Sprintf("uploads/%s.tar.gz", chunks[0].UUID)
		return client.Remove(loc)
	}

	for _, chunk := range chunks {
		loc := fmt.Sprintf("uploads/%s/%d_%d.tar.gz", chunk.UUID, chunk.RangeStart, chunk.RangeEnd)
		err := client.Remove(loc)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetBlob returns the blob object as a reader interface
func GetBlob(digest, repoName, orgName string) (io.ReadCloser, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("blobs/%s/%s/%s.tar.gz", orgName, repoName, digest)
	obj, err := client.Get(path)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving object")
	}

	return obj, nil
}

// DeleteLayer deletes a blob
func DeleteBlob(digest, repoName, orgName string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("blobs/%s/%s/%s.tar.gz", orgName, repoName, digest)
	return client.Remove(path)
}

// DeleteUpload deletes an upload
func DeleteUpload(uuid string) error {
	chunks, err := chunk.GetAll(uuid)
	if err != nil {
		return errors.Wrap(err, "error getting upload chunks from database")
	}

	err = deleteChunks(chunks)
	if err != nil {
		return errors.Wrap(err, "error deleting chunks")
	}

	return nil
}
