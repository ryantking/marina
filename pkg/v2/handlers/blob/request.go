package blob

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db/models/blob"
	"github.com/ryantking/marina/pkg/db/models/repo"
	"github.com/ryantking/marina/pkg/docker"
)

var (
	repoExists = repo.Exists
	blobExists = blob.Exists
)

func parsePath(c echo.Context) (string, string, string, error) {
	digest := c.Param("digest")
	repoName := c.Param("repo")
	orgName := c.Param("org")
	exists, err := repoExists(repoName, orgName)
	if err != nil {
		return "", "", "", errors.Wrap(err, "error checking if repository exists")
	}
	if !exists {
		c.Set("docker_err_code", docker.CodeNameUnknown)
		return "", "", "", echo.NewHTTPError(http.StatusNotFound, "no such repository")
	}

	exists, err = blobExists(digest)
	if err != nil {
		return "", "", "", errors.Wrap(err, "error checking if layer exists")
	}
	if !exists {
		c.Set("docker_err_code", docker.CodeBlobUnknown)
		return "", "", "", echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	return digest, repoName, orgName, nil
}
