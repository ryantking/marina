package upload

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db/models/blob"
	"github.com/ryantking/marina/pkg/db/models/upload"
	"github.com/ryantking/marina/pkg/db/models/upload/chunk"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/store"
)

const (
	headerRange = "Range"
)

// Get returns an upload
func Get(c echo.Context) error {
	_, _, uuid, err := parsePath(c)
	if err != nil {
		return err
	}

	lastRange, err := chunk.GetLastRange(uuid)
	if err != nil {
		return err
	}

	c.Response().Header().Set(headerRange, lastRange)
	c.Response().Header().Set(docker.HeaderUploadUUID, uuid)
	return c.NoContent(http.StatusNoContent)
}

// Start starts the process of uploading a new blob
func Start(c echo.Context) error {
	repoName, orgName, _, err := parsePath(c)
	if err != nil {
		return err
	}

	uuid, err := upload.New()
	if err != nil {
		return errors.Wrap(err, "error creating new upload")
	}

	loc := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%s", orgName, repoName, uuid)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(headerRange, "0-0")
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	return c.NoContent(http.StatusAccepted)
}

// Blob is used to save a chunk of data to a blob
func Blob(c echo.Context) error {
	defer c.Request().Body.Close()
	repoName, orgName, uuid, err := parsePath(c)
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
	err = storeChunk(c, uuid, sz, start)
	if err != nil {
		return err
	}

	loc := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%s", orgName, repoName, uuid)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set(docker.HeaderUploadUUID, fmt.Sprint(uuid))
	lastRange, err := chunk.GetLastRange(uuid)
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
		err := storeChunk(c, uuid, sz, 0)
		if err != nil {
			return err
		}
	}

	err = store.FinishUpload(digest, uuid, repoName, orgName)
	if err != nil {
		return errors.Wrap(err, "error finishing upload")
	}

	err = upload.Finish(uuid)
	if err != nil {
		return errors.Wrap(err, "error updating upload status")
	}

	_, err = blob.New(digest, repoName, orgName)
	if err != nil {
		return errors.Wrap(err, "error creating new blob in database")
	}

	loc := fmt.Sprintf("/v2/%s/%s/blobs/%s", orgName, repoName, digest)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set(docker.HeaderContentDigest, digest)
	return c.NoContent(http.StatusCreated)
}

// Cancel cancels a currently running upload
func Cancel(c echo.Context) error {
	_, _, uuid, err := parsePath(c)
	if err != nil {
		return err
	}

	err = store.DeleteUpload(uuid)
	if err != nil {
		return err
	}

	err = upload.Delete(uuid)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func storeChunk(c echo.Context, uuid string, sz, start int64) error {
	sz, err := store.UploadChunk(uuid, c.Request().Body, sz, start)
	if err != nil {
		return errors.Wrap(err, "error storing chunk")
	}
	end := start + sz - 1
	err = chunk.New(uuid, start, end)
	if err != nil {
		return errors.Wrap(err, "error creating upload chunk")
	}

	return nil
}
