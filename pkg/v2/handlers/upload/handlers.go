package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db/models/layer"
	"github.com/ryantking/marina/pkg/db/models/repo"
	"github.com/ryantking/marina/pkg/db/models/upload"
	"github.com/ryantking/marina/pkg/docker"
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
	c.Response().Header().Set("Content-Range", "0-0")
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	return c.NoContent(http.StatusAccepted)
}

// Chunk is used to save a chunk of data to a layer
func Chunk(c echo.Context) error {
	repoName, orgName, uuid, err := parsePath(c)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s.tar.gz", uuid))
	if err != nil {
		return err
	}
	n, err := io.Copy(f, c.Request().Body)
	if err != nil {
		return err
	}

	loc := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%s", orgName, repoName, uuid)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set("Range", fmt.Sprintf("0-%d", n))
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set(docker.HeaderUploadUUID, fmt.Sprint(uuid))
	return c.NoContent(http.StatusAccepted)
}

// Finish is used to signal when the upload has completed
func Finish(c echo.Context) error {
	repoName, orgName, uuid, err := parsePath(c)
	if err != nil {
		return err
	}
	digest := c.QueryParam("digest")

	os.Rename(fmt.Sprintf("%s.tar.gz", uuid), fmt.Sprintf("%s.tar.gz", digest))
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
