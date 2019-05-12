package store

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

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

// DeleteBlob deletes a blob from the storage backend
func DeleteBlob(digest, repoName, orgName string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("blobs/%s/%s/%s.tar.gz", orgName, repoName, digest)
	return client.Remove(path)
}

// CreateUpload creates a new upload from the given reader
func UploadChunk(uuid string, r io.Reader, sz, start int64) (int64, error) {
	client, err := getClient()
	if err != nil {
		return 0, err
	}

	path := fmt.Sprintf("uploads/%s/%d.tar.gz", uuid, start)
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
	chunks, err := getChunks(uuid)
	if err != nil {
		return errors.Wrap(err, "error getting upload chunks from database")
	}

	loc := fmt.Sprintf("blobs/%s/%s/%s.tar.gz", orgName, repoName, digest)
	err = chunksToBlob(client, chunks, loc)
	if err != nil {
		return errors.Wrap(err, "error merging and uploading chunks")
	}

	return nil
}

// DeleteUpload deletes an upload
func DeleteUpload(uuid string) error {
	chunks, err := getChunks(uuid)
	if err != nil {
		return errors.Wrap(err, "error getting upload chunks from database")
	}

	err = deleteChunks(client, chunks)
	if err != nil {
		return errors.Wrap(err, "error deleting chunks")
	}

	return nil
}
