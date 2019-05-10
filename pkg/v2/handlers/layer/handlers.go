package layer

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db/models/layer"
	"github.com/ryantking/marina/pkg/db/models/repo"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/store"
)

// Exists returns a status 200 or 404 depending on if the layer exists or not
func Exists(c echo.Context) error {
	digest, _, _, err := parsePath(c)
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentLength, fmt.Sprint(len(digest)))
	c.Response().Header().Set(docker.HeaderContentDigest, digest)
	return c.NoContent(http.StatusOK)
}

// Get downloads a layer to the response body, responding 404 if the layer is not found
func Get(c echo.Context) error {
	digest, repoName, orgName, err := parsePath(c)
	if err != nil {
		return err
	}

	r, err := store.GetLayer(digest, repoName, orgName)
	if err != nil {
		return err
	}
	defer r.Close()

	return c.Stream(http.StatusOK, echo.MIMEOctetStream, r)
}

// Delete deletes a layer from the repository
func Delete(c echo.Context) error {
	digest, repoName, orgName, err := parsePath(c)
	if err != nil {
		return err
	}

	err = store.DeleteLayer(digest, repoName, orgName)
	if err != nil {
		return err
	}
	err = layer.Delete(digest)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}

func parsePath(c echo.Context) (string, string, string, error) {
	orgName := c.Param("org")
	repoName := c.Param("repo")
	exists, err := repo.Exists(repoName, orgName)
	if err != nil {
		return "", "", "", errors.Wrap(err, "error checking if repository exists")
	}
	if !exists {
		c.Set("docker_err_code", "NAME_UNKNOWN")
		return "", "", "", echo.NewHTTPError(http.StatusNotFound, "no such repository")
	}

	digest := c.Param("digest")
	exists, err = layer.Exists(digest)
	if err != nil {
		return "", "", "", errors.Wrap(err, "error checking if layer exists")
	}
	if !exists {
		c.Set("docker_err_code", "BLOB_UNKNOWN")
		return "", "", "", echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	return digest, repoName, orgName, nil
}
