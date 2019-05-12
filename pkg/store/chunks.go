package store

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db/models/upload/chunk"
)

var (
	getChunks = chunk.GetAll
)

func chunksToBlob(client Client, chunks []*chunk.Model, loc string) error {
	r, n, err := mergeChunks(client, chunks)
	if err != nil {
		return errors.Wrap(err, "error getting upload reader")
	}

	_, err = client.Put(loc, r, n)
	if err != nil {
		return errors.Wrap(err, "error uploading blob")
	}

	err = deleteChunks(client, chunks)
	if err != nil {
		return errors.Wrap(err, "error deleting chunks")
	}

	return nil
}

func mergeChunks(client Client, chunks []*chunk.Model) (io.Reader, int64, error) {
	var sz int64
	readers := make([]io.Reader, len(chunks))
	for i, chunk := range chunks {
		loc := fmt.Sprintf("uploads/%s/%d.tar.gz", chunk.UUID, chunk.RangeStart)
		obj, err := client.Get(loc)
		if err != nil {
			return nil, 0, err
		}
		readers[i] = obj
		sz += chunk.RangeEnd - chunk.RangeStart + 1
	}

	return io.MultiReader(readers...), sz, nil
}

func deleteChunks(client Client, chunks []*chunk.Model) error {
	for _, chunk := range chunks {
		loc := fmt.Sprintf("uploads/%s/%d.tar.gz", chunk.UUID, chunk.RangeStart)
		err := client.Remove(loc)
		if err != nil {
			return err
		}
	}

	return nil
}
