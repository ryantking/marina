package upload

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
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

	finishUpload = store.FinishUpload
	deleteUpload = store.DeleteUpload
)

// Get returns an upload
func Get(c echo.Context) error {
	_, _, err := web.GetRepo(c)
	if err != nil {
		return err
	}

	uuid := c.Param("uuid")
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
	repo, org, err := web.GetRepo(c)
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

	loc := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%s", org.Name, repo.Name, uuid)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(headerRange, lastRange)
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	return c.NoContent(http.StatusAccepted)
}

// Chunk is used to save a chunk of data to a blob
func Chunk(c echo.Context) error {
	defer c.Request().Body.Close()

	repo, org, err := web.GetRepo(c)
	if err != nil {
		return err
	}

	uuid := c.Param("uuid")
	sz, start, err := web.ParseLengthAndStart(c)
	if err != nil {
		return err
	}
	sz, err = createChunk(c.Request().Context(), uuid, c.Request().Body, sz, start)
	if err != nil {
		return err
	}

	loc := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%s", org.Name, repo.Name, uuid)
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
	repo, org, err := web.GetRepo(c)
	if err != nil {
		return err
	}
	uuid := c.Param("uuid")
	digest := c.QueryParam("digest")
	sz, start, err := web.ParseLengthAndStart(c)
	if err != nil {
		return err
	}

	if sz > 0 {
		defer c.Request().Body.Close()
		sz, err = createChunk(c.Request().Context(), uuid, c.Request().Body, sz, start)
		if err != nil {
			return err
		}
	}

	err = createBlob(c.Request().Context(), repo, uuid, digest, org.Name)
	if err != nil {
		return err
	}

	loc := fmt.Sprintf("/v2/%s/%s/blobs/%s", org.Name, repo.Name, digest)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set(docker.HeaderContentDigest, digest)
	return c.NoContent(http.StatusCreated)
}

// Cancel cancels a currently running upload
func Cancel(c echo.Context) error {
	_, _, err := web.GetRepo(c)
	if err != nil {
		return err
	}

	uuid := c.Param("uuid")
	err = deleteUpload(uuid)
	if err != nil {
		return err
	}

	_, err = client.DeleteUpload(prisma.UploadWhereUniqueInput{Uuid: &uuid}).Exec(c.Request().Context())
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
