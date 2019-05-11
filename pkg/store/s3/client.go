package s3

import (
	"io"

	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/config"
)

type Client struct {
	minio.Client
	bucket string
}

// New returns a new S3 client
func New() (*Client, error) {
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

	return &Client{*c, cfg.Bucket}, nil
}

// Get returns an object from a given path
func (c *Client) Get(path string) (io.ReadCloser, error) {
	obj, err := c.GetObject(c.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving object from S3")
	}

	return obj, nil
}

// Put saves an object to the bucket
func (c *Client) Put(path string, r io.Reader, sz int64) (int64, error) {
	n, err := c.PutObject(c.bucket, path, r, sz, minio.PutObjectOptions{})
	if err != nil {
		return 0, errors.Wrap(err, "error uploading object to S3")
	}

	return n, nil
}

// Remove deletes an object from the bucket
func (c *Client) Remove(path string) error {
	return c.RemoveObject(c.bucket, path)
}
