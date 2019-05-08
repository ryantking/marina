package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db/models/layer"
	"github.com/ryantking/marina/pkg/db/models/repo"
	"github.com/ryantking/marina/pkg/db/models/upload"
	"github.com/ryantking/marina/pkg/docker"
)

// Start starts the process of uploading a new layer
func Start(c echo.Context) error {
	repoName, orgName := parsePath(c)
	exists, err := repo.Exists(repoName, orgName)
	if err != nil {
		return errors.Wrap(err, "error checking if repo exists")
	}
	if !exists {
		c.Set("docker_err_code", "NAME_UNKNOWN")
		return echo.NewHTTPError(http.StatusNotFound, "repository not found")
	}

	upl, err := upload.New()
	if err != nil {
		return errors.Wrap(err, "error creating new upload")
	}

	loc := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%d", orgName, repoName, upl.UUID)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set("Content-Range", "0-0")
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	return c.NoContent(http.StatusAccepted)
}

// Chunk is used to save a chunk of data to a layer
func Chunk(c echo.Context) error {
	repoName, orgName := parsePath(c)
	uuid, err := parseUUID(c)
	if err != nil {
		c.Set("docker_err_code", "BLOB_UPLOAD_INVALID")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	f, err := os.Create(fmt.Sprintf("upload_%d.tar.gz", uuid))
	if err != nil {
		return err
	}
	n, err := io.Copy(f, c.Request().Body)
	if err != nil {
		return err
	}

	loc := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%d", orgName, repoName, uuid)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set("Range", fmt.Sprintf("0-%d", n))
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set(docker.HeaderUploadUUID, fmt.Sprint(uuid))
	return c.NoContent(http.StatusAccepted)
}

// Finish is used to signal when the upload has completed
func Finish(c echo.Context) error {
	repoName, orgName := parsePath(c)
	uuid, err := parseUUID(c)
	if err != nil {
		c.Set("docker_err_code", "BLOB_UPLOAD_INVALID")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	digest := c.QueryParam("digest")

	os.Rename(fmt.Sprintf("upload_%d.tar.gz", uuid), fmt.Sprintf("%s_%s_%s.tar.gz", orgName, repoName, digest))
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

func parsePath(c echo.Context) (string, string) {
	return c.Param("repo"), c.Param("org")
}

func parseUUID(c echo.Context) (uint64, error) {
	s := c.Param("uuid")
	return strconv.ParseUint(s, 10, 64)
}
