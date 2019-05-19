package web

import (
	"context"
	"net/http"

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
