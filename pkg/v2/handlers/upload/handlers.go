package upload

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db/models/layer"
	"github.com/ryantking/marina/pkg/db/models/repo"
	"github.com/ryantking/marina/pkg/db/models/upload"
	"github.com/ryantking/marina/pkg/db/models/upload/chunk"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/store"
)

const (
	headerRange = "Range"
)

// Start starts the process of uploading a new layer
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

// Blob is used to save a chunk of data to a layer
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
	start, end, err := parseRange(c)
	if err != nil {
		return err
	}
	err = storeChunk(c, uuid, sz, start, end)
	if err != nil {
		return err
	}

	loc := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%s", orgName, repoName, uuid)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set(docker.HeaderUploadUUID, fmt.Sprint(uuid))
	err = setRangeHeader(c, uuid)
	if err != nil {
		return errors.Wrap(err, "error setting range header")
	}
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
		err := storeChunk(c, uuid, sz, 0, 0)
		if err != nil {
			return err
		}
	}

	err = store.FinishUpload(digest, uuid, repoName, orgName)
	if err != nil {
		return errors.Wrap(err, "error finishing upload")
	}

	upl := &upload.Model{UUID: uuid, Done: true}
	err = upl.Save()
	if err != nil {
		return errors.Wrap(err, "error updating upload status")
	}

	_, err = layer.New(digest, repoName, orgName)
	if err != nil {
		return errors.Wrap(err, "error creating new layer in database")
	}

	loc := fmt.Sprintf("/v2/%s/%s/blobs/%s", orgName, repoName, digest)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set(docker.HeaderContentDigest, digest)
	return c.NoContent(http.StatusCreated)
}

func parsePath(c echo.Context) (string, string, string, error) {
	repoName := c.Param("repo")
	orgName := c.Param("org")
	uuid := c.Param("uuid")
	exists, err := repo.Exists(repoName, orgName)
	if err != nil {
		return "", "", "", errors.Wrap(err, "error checking if repository exists")
	}
	if !exists {
		c.Set("docker_err_code", docker.CodeNameUnknown)
		return "", "", "", echo.NewHTTPError(http.StatusNotFound, "no such repository")
	}

	return repoName, orgName, uuid, nil
}

func parseLength(c echo.Context) (int64, error) {
	s := c.Request().Header.Get(echo.HeaderContentLength)
	if s == "" {
		return -1, nil
	}

	sz, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		c.Set("docker_err_code", "BLOB_UPLOAD_INVALID")
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return sz, nil
}

func parseRange(c echo.Context) (int64, int64, error) {
	s := c.Request().Header.Get("Content-Range")
	if s == "" {
		return 0, 0, nil
	}
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		c.Set("docker_err_code", "BLOB_UPLOAD_INVALID")
		return 0, 0, echo.NewHTTPError(http.StatusBadRequest, "invalid Content-Range header")
	}
	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		c.Set("docker_err_code", "BLOB_UPLOAD_INVALID")
		return 0, 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	end, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		c.Set("docker_err_code", "BLOB_UPLOAD_INVALID")
		return 0, 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return start, end, nil
}

func storeChunk(c echo.Context, uuid string, sz, start, end int64) error {
	sz, err := store.UploadChunk(uuid, c.Request().Body, sz, start, end)
	if err != nil {
		return errors.Wrap(err, "error storing chunk")
	}
	end = start + sz - 1
	err = chunk.New(uuid, start, end)
	if err != nil {
		return errors.Wrap(err, "error creating upload chunk")
	}

	return nil
}

func setRangeHeader(c echo.Context, uuid string) error {
	chunks, err := chunk.GetAll(uuid)
	if err != nil {
		return errors.Wrap(err, "error retrieving upload chunks")
	}
	ranges := make([]string, len(chunks), len(chunks))
	for i, chunk := range chunks {
		ranges[i] = fmt.Sprintf("%d-%d", chunk.RangeStart, chunk.RangeEnd)
	}
	c.Response().Header().Set(headerRange, strings.Join(ranges, ","))
	return nil
}
