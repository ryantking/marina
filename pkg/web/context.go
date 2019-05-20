package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/prisma"
)

var client = prisma.New(nil)

// ParsePath parses the repository and organization out of the path
func ParsePath(c echo.Context) (string, string, error) {
	repoName := c.Param("repo")
	orgName := c.Param("org")
	repos, err := client.Repositories(&prisma.RepositoriesParams{
		Where: &prisma.RepositoryWhereInput{
			Name: &repoName,
			Org: &prisma.OrganizationWhereInput{
				Name: &orgName,
			},
		},
	}).Exec(context.Background())
	if err != nil {
		return "", "", err
	}
	if len(repos) == 0 {
		c.Set("docker_err_code", docker.CodeNameUnknown)
		return "", "", echo.NewHTTPError(http.StatusNotFound, "repository name not known to registry")
	}

	return repoName, orgName, nil
}

// ParsePagination parses the docker page number and last from the request
func ParsePagination(c echo.Context) (int32, string, error) {
	s := c.QueryParam("n")
	if s == "" {
		return 0, "", nil
	}
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, "", echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return int32(n), c.QueryParam("last"), nil
}
