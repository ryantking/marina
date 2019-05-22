package upload

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/prisma"
	"github.com/ryantking/marina/pkg/store"
	"github.com/ryantking/marina/pkg/web"
)

const (
	headerRange = "Range"
)

var (
	client = prisma.New(nil)

	storeChunk = store.UploadChunk
)

// Get returns an upload
func Get(c echo.Context) error {
	uuid, _, _, err := parsePath(c)
	if err != nil {
		return err
	}

	lastRange, err := getLastRange(c.Request().Context(), uuid)
	if err != nil {
		return err
	}

	c.Response().Header().Set(headerRange, lastRange)
	c.Response().Header().Set(docker.HeaderUploadUUID, uuid)
	return c.NoContent(http.StatusNoContent)
}

// Start starts the process of uploading a new blob
func Start(c echo.Context) error {
	repo, org, err := web.ParsePath(c)
	if err != nil {
		return err
	}

	uuid := uuid.New().String()
	_, err = client.CreateUpload(prisma.UploadCreateInput{Uuid: uuid}).Exec(c.Request().Context())
	if err != nil {
		return err
	}
	lastRange, err := getLastRange(c.Request().Context(), uuid)
	if err != nil {
		return err
	}

	loc := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%s", org, repo, uuid)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(headerRange, lastRange)
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	return c.NoContent(http.StatusAccepted)
}

// Blob is used to save a chunk of data to a blob
func Blob(c echo.Context) error {
	defer c.Request().Body.Close()
	uuid, repoName, orgName, err := parsePath(c)
	if err != nil {
		return err
	}
	sz, err := parseLength(c)
	if err != nil {
		return err
	}
	start, err := parseRange(c)
	if err != nil {
		return err
	}
	sz, err = storeChunk(uuid, c.Request().Body, sz, start)
	if err != nil {
		return err
	}
	_, err = client.CreateChunk(prisma.ChunkCreateInput{
		RangeStart: start,
		RangeEnd:   start + sz - 1,
		Upload: prisma.UploadCreateOneWithoutChunksInput{
			Connect: &prisma.UploadWhereUniqueInput{Uuid: &uuid},
		},
	}).Exec(c.Request().Context())
	if err != nil {
		return err
	}

	loc := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%s", orgName, repoName, uuid)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set(docker.HeaderUploadUUID, fmt.Sprint(uuid))
	lastRange, err := getLastRange(c.Request().Context(), uuid)
	if err != nil {
		return err
	}
	c.Response().Header().Set(headerRange, lastRange)
	return c.NoContent(http.StatusNoContent)
}

// Finish is used to signal when the upload has completed
func Finish(c echo.Context) error {
	repoName, orgName, uuid, err := parsePath(c)
	if err != nil {
		return err
	}
	digest := c.QueryParam("digest")
	sz, err := parseLength(c)
	if err != nil {
		return err
	}

	if sz > 0 {
		defer c.Request().Body.Close()
		start, err := parseRange(c)
		if err != nil {
			return err
		}

		_, err = storeChunk(uuid, c.Request().Body, sz, start)
		if err != nil {
			return err
		}
	}

	err = store.FinishUpload(digest, uuid, repoName, orgName)
	if err != nil {
		return errors.Wrap(err, "error finishing upload")
	}

	_, err = client.UpdateUpload(prisma.UploadUpdateParams{
		Where: prisma.UploadWhereUniqueInput{Uuid: &uuid},
		Data:  prisma.UploadUpdateInput{Done: prisma.Bool(true)},
	}).Exec(c.Request().Context())

	_, err = client.CreateBlob(prisma.BlobCreateInput{Digest: digest}).Exec(c.Request().Context())
	if err != nil {
		return err
	}

	loc := fmt.Sprintf("/v2/%s/%s/blobs/%s", orgName, repoName, digest)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set(docker.HeaderContentDigest, digest)
	return c.NoContent(http.StatusCreated)
}

// // Cancel cancels a currently running upload
// func Cancel(c echo.Context) error {
// 	_, _, uuid, err := parsePath(c)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = store.DeleteUpload(uuid)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = upload.Delete(uuid)
// 	if err != nil {
// 		return err
// 	}
//
// 	return c.NoContent(http.StatusNoContent)
// }
//
// func storeChunk(c echo.Context, uuid string, sz, start int64) error {
// 	sz, err := store.UploadChunk(uuid, c.Request().Body, sz, start)
// 	if err != nil {
// 		return errors.Wrap(err, "error storing chunk")
// 	}
// 	end := start + sz - 1
// 	err = chunk.New(uuid, start, end)
// 	if err != nil {
// 		return errors.Wrap(err, "error creating upload chunk")
// 	}
//
// 	return nil
// }
