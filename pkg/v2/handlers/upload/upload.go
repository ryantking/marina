package upload

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/prisma"
	"github.com/ryantking/marina/pkg/store"
)

var (
	orderBy = prisma.ChunkOrderByInputRangeEndDesc

	storeChunk = store.UploadChunk
)

func getLastRange(ctx context.Context, uuid string) (string, error) {
	chunks, err := client.Chunks(&prisma.ChunksParams{
		First:   prisma.Int32(1),
		OrderBy: &orderBy,
		Where: &prisma.ChunkWhereInput{
			Upload: &prisma.UploadWhereInput{
				Uuid: &uuid,
			},
		},
	}).Exec(ctx)
	if len(chunks) == 0 {
		return "0-0", nil
	}
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d-%d", chunks[0].RangeStart, chunks[0].RangeEnd), nil
}

func createChunk(ctx context.Context, uuid string, r io.Reader, sz, start int32) (int32, error) {
	sz, err := storeChunk(uuid, r, sz, start)
	if err != nil {
		return 0, err
	}
	_, err = client.CreateChunk(prisma.ChunkCreateInput{
		RangeStart: start,
		RangeEnd:   start + sz - 1,
		Upload: prisma.UploadCreateOneWithoutChunksInput{
			Connect: &prisma.UploadWhereUniqueInput{Uuid: &uuid},
		},
	}).Exec(ctx)
	if err != nil {
		return 0, err
	}

	return sz, nil
}

func createBlob(ctx context.Context, repo *prisma.Repository, uuid, digest, orgName string) error {
	err := finishUpload(digest, uuid, repo.Name, orgName)
	if err != nil {
		return errors.Wrap(err, "error finishing upload")
	}

	_, err = client.UpdateUpload(prisma.UploadUpdateParams{
		Where: prisma.UploadWhereUniqueInput{Uuid: &uuid},
		Data:  prisma.UploadUpdateInput{Done: prisma.Bool(true)},
	}).Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.CreateBlob(prisma.BlobCreateInput{
		Digest: digest,
		Repo: prisma.RepositoryCreateOneWithoutBlobsInput{
			Connect: &prisma.RepositoryWhereUniqueInput{ID: &repo.ID},
		},
	}).Exec(ctx)
	return err
}
