package upload

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db/models/repo"
	"github.com/ryantking/marina/pkg/docker"
)

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

func parseRange(c echo.Context) (int64, error) {
	s := c.Request().Header.Get("Content-Range")
	if s == "" {
		return 0, nil
	}
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		c.Set("docker_err_code", "BLOB_UPLOAD_INVALID")
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid Content-Range header")
	}
	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		c.Set("docker_err_code", "BLOB_UPLOAD_INVALID")
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return start, nil
}
