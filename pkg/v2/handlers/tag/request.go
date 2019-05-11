package tag

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db/models/repo"
	"github.com/ryantking/marina/pkg/docker"
)

var (
	repoExists = repo.Exists
)

func parsePath(c echo.Context) (string, string, error) {
	repoName := c.Param("repo")
	orgName := c.Param("org")
	exists, err := repoExists(repoName, orgName)
	if err != nil {
		return "", "", errors.Wrap(err, "error checking if repository exists")
	}
	if !exists {
		c.Set("docker_err_code", docker.CodeNameUnknown)
		return "", "", echo.NewHTTPError(http.StatusNotFound, "no such repository")
	}

	return repoName, orgName, nil
}

func parsePagination(c echo.Context) (uint, string, error) {
	s := c.QueryParam("n")
	if s == "" {
		return 0, "", nil
	}
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, "", err
	}
	return uint(n), c.QueryParam("last"), nil
}
