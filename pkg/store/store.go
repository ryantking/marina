package store

import (
	"fmt"
	"io"

	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/config"
	"github.com/ryantking/marina/pkg/db/models/upload/chunk"
)

var client *minio.Client

func getClient() (*minio.Client, error) {
	if client != nil {
		return client, nil
	}

	cfg := config.Get().S3
	c, err := minio.NewWithRegion(cfg.Endpoint, cfg.AccessKeyID, cfg.SecretAccessKey, cfg.UseSSL, cfg.Region)
	if err != nil {
		return nil, errors.Wrap(err, "error creating minio client")
	}
	exists, err := c.BucketExists(cfg.Bucket)
	if err != nil {
		return nil, errors.Wrap(err, "error checking if bucket exists")
	}
	if !exists {
		panic("bucket does not exist")
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
	n, err := client.PutObject(config.Get().S3.Bucket, path, r, sz, minio.PutObjectOptions{})
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
	bucket := config.Get().S3.Bucket
	chunks, err := chunk.GetAll(uuid)
	if err != nil {
		return errors.Wrap(err, "error getting upload chunks from database")
	}
	r, n, err := mergeChunks(chunks)
	if err != nil {
		return errors.Wrap(err, "error getting upload reader")
	}

	loc := fmt.Sprintf("blobs/%s/%s/%s.tar.gz", orgName, repoName, digest)
	_, err = client.PutObject(bucket, loc, r, n, minio.PutObjectOptions{})
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

	bucket := config.Get().S3.Bucket
	if len(chunks) == 1 {
		loc := fmt.Sprintf("uploads/%s.tar.gz", chunks[0].UUID)
		obj, err := client.GetObject(bucket, loc, minio.GetObjectOptions{})
		if err != nil {
			return nil, 0, err
		}
		return obj, chunks[0].RangeEnd - chunks[0].RangeStart + 1, nil
	}

	var sz int64
	readers := make([]io.Reader, len(chunks))
	for i, chunk := range chunks {
		loc := fmt.Sprintf("uploads/%s/%d_%d.tar.gz", chunk.UUID, chunk.RangeStart, chunk.RangeEnd)
		obj, err := client.GetObject(bucket, loc, minio.GetObjectOptions{})
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

	bucket := config.Get().S3.Bucket
	if len(chunks) == 1 {
		loc := fmt.Sprintf("uploads/%s.tar.gz", chunks[0].UUID)
		return client.RemoveObject(bucket, loc)
	}

	for _, chunk := range chunks {
		loc := fmt.Sprintf("uploads/%s/%d_%d.tar.gz", chunk.UUID, chunk.RangeStart, chunk.RangeEnd)
		err := client.RemoveObject(bucket, loc)
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
	obj, err := client.GetObject(config.Get().S3.Bucket, path, minio.GetObjectOptions{})
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
	return client.RemoveObject(config.Get().S3.Bucket, path)
}
