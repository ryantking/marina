package store

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/prisma"
)

var c = prisma.New(nil)

func getChunks(uuid string) ([]prisma.Chunk, error) {
	chunks, err := c.Upload(prisma.UploadWhereUniqueInput{Uuid: &uuid}).Chunks(nil).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return chunks, nil
}

func mergeChunks(client Client, uuid string, chunks []prisma.Chunk, loc string) error {
	fromPaths := make([]string, len(chunks))
	for i, chunk := range chunks {
		fromPaths[i] = fmt.Sprintf("uploads/%s/%d.tar.gz", uuid, chunk.RangeStart)
	}
	err := client.Merge(loc, fromPaths...)
	if err != nil {
		return err
	}

	err = deleteChunks(client, uuid, chunks)
	if err != nil {
		return errors.Wrap(err, "error deleting chunks")
	}

	return nil
}

func deleteChunks(client Client, uuid string, chunks []prisma.Chunk) error {
	for _, chunk := range chunks {
		loc := fmt.Sprintf("uploads/%s/%d.tar.gz", uuid, chunk.RangeStart)
		err := client.Remove(loc)
		if err != nil {
			return err
		}
	}

	return nil
}
