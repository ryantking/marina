package store

import (
	"fmt"
	"io"

	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/config"
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
func CreateUpload(uuid, repoName, orgName string, r io.Reader, sz int64) (int64, error) {
	client, err := getClient()
	if err != nil {
		return 0, err
	}

	path := fmt.Sprintf("uploads/%s/%s/%s.tar.gz", orgName, repoName, uuid)
	n, err := client.PutObject(config.Get().S3.Bucket, path, r, sz, minio.PutObjectOptions{})
	if err != nil {
		return 0, errors.Wrap(err, "error uploading object")
	}
	return n, err
}

// FinishUpload finalizes an upload
func FinishUpload(digest, uuid, repoName, orgName string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	bucket := config.Get().S3.Bucket
	srcPath := fmt.Sprintf("uploads/%s/%s/%s.tar.gz", orgName, repoName, uuid)
	src := minio.NewSourceInfo(bucket, srcPath, nil)
	dstPath := fmt.Sprintf("layers/%s/%s/%s.tar.gz", orgName, repoName, digest)
	dst, err := minio.NewDestinationInfo(bucket, dstPath, nil, nil)
	if err != nil {
		return err
	}
	err = client.CopyObject(dst, src)
	if err != nil {
		return err
	}

	return client.RemoveObject(bucket, srcPath)
}
