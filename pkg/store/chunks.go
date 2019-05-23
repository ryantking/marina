package store

import (
	"context"
	"fmt"
	"io"
	"strings"

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
	var sz int32
	readers := make([]io.Reader, len(chunks))
	for i, chunk := range chunks {
		chunkLoc := fmt.Sprintf("uploads/%s/%d.tar.gz", uuid, chunk.RangeStart)
		obj, err := client.Get(chunkLoc)
		if err != nil {
			return err
		}
		readers[i] = obj
		sz += chunk.RangeEnd - chunk.RangeStart + 1
	}

	_, err := client.Put(loc, strings.NewReader("fucker"), sz)
	if err != nil {
		return errors.Wrap(err, "error uploading blob")
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
