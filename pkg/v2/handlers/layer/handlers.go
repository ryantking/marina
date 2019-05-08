package layer

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/db/models/layer"
	"github.com/ryantking/marina/pkg/docker"
)

// Exists returns a status 200 or 404 depending on if the layer exists or not
func Exists(c echo.Context) error {
	digest, repoName, orgName := parsePath(c)
	exists, err := layer.Exists(digest, repoName, orgName)
	if err != nil {
		return err
	}
	if !exists {
		c.Set("docker_err_code", "BLOB_UNKNOWN")
		return echo.NewHTTPError(http.StatusNotFound, "could not find layer")
	}

	c.Response().Header().Set(echo.HeaderContentLength, fmt.Sprint(len(digest)))
	c.Response().Header().Set(docker.HeaderContentDigest, digest)
	return c.NoContent(http.StatusOK)
}

// Get downloads a layer to the response body, responding 404 if the layer is not found
func Get(c echo.Context) error {
	digest, repoName, orgName := parsePath(c)
	exists, err := layer.Exists(digest, repoName, orgName)
	if err != nil {
		return err
	}
	if !exists {
		c.Set("docker_err_code", "BLOB_UNKNOWN")
		return echo.NewHTTPError(http.StatusNotFound, "could not find layer")
	}

	f, err := os.Open(fmt.Sprintf("%s_%s_%s.tar.gz", orgName, repoName, digest))
	if err != nil {
		return err
	}

	return c.Stream(http.StatusOK, echo.MIMEOctetStream, f)
}

func parsePath(c echo.Context) (string, string, string) {
	return c.Param("digest"), c.Param("repo"), c.Param("org")
}
