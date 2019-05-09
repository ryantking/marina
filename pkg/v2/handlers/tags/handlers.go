package tags

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/db/models/repo"
	"github.com/ryantking/marina/pkg/db/models/tag"
	"github.com/ryantking/marina/pkg/docker"
)

// List lists all tags for a given repository
func List(c echo.Context) error {
	repoName, orgName := parsePath(c)
	exists, err := repo.Exists(repoName, orgName)
	if err != nil {
		return err
	}
	if !exists {
		c.Set("docker_err_code", docker.CodeNameUnknown)
		return echo.NewHTTPError(http.StatusNotFound, "no such repository")
	}
	pageSize, last, err := pageAndLast(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not parse query parameters")
	}
	tags, nextLast, err := tag.List(repoName, orgName, pageSize, last)
	if err != nil {
		return err
	}

	if nextLast != "" {
		link := fmt.Sprintf("/%s/%s/tags/list?n=%d&last=%s", orgName, repoName, pageSize, nextLast)
		c.Response().Header().Set("Link", link)
	}
	res := map[string]interface{}{}
	res["name"] = fmt.Sprintf("%s/%s", orgName, repoName)
	res["tags"] = tags
	return c.JSON(http.StatusOK, res)
}

func parsePath(c echo.Context) (string, string) {
	return c.Param("repo"), c.Param("org")
}

func pageAndLast(c echo.Context) (uint, string, error) {
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
