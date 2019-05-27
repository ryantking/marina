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
func (c *Client) Get(path string, start, end int64) (io.ReadCloser, error) {
	var opts minio.GetObjectOptions
	if start >= 0 && end > 0 {
		err := opts.SetRange(start, end)
		if err != nil {
			return nil, err
		}
	}
	obj, err := c.GetObject(c.bucket, path, opts)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving object from S3")
	}

	return obj, nil
}

// Put saves an object to the bucket
func (c *Client) Put(path string, r io.Reader, sz int32) (int32, error) {
	n, err := c.PutObject(c.bucket, path, r, int64(sz), minio.PutObjectOptions{})
	if err != nil {
		return 0, errors.Wrap(err, "error uploading object to S3")
	}

	return int32(n), nil
}

// Remove deletes an object from the bucket
func (c *Client) Remove(path string) error {
	return c.RemoveObject(c.bucket, path)
}

// Merge combines multiple objects into a single one
func (c *Client) Merge(toPath string, fromPaths ...string) error {
	srcs := make([]minio.SourceInfo, len(fromPaths))
	for i, fromPath := range fromPaths {
		srcs[i] = minio.NewSourceInfo(c.bucket, fromPath, nil)
	}
	dst, err := minio.NewDestinationInfo(c.bucket, toPath, nil, nil)
	if err != nil {
		return err
	}

	return c.ComposeObject(dst, srcs)
}
