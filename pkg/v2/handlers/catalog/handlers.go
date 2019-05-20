package catalog

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/prisma"
	"github.com/ryantking/marina/pkg/web"
)

const (
	headerLink = "Link"
)

var (
	client  = prisma.New(nil)
	orderBy = prisma.RepositoryOrderByInputFullNameAsc
)

// Get returns the catalog for the registry
func Get(c echo.Context) error {
	n, last, err := web.ParsePagination(c)
	if err != nil {
		return err
	}

	if n == 0 {
		return getAll(c)
	}

	return getPaginated(c, n, last)
}

func getAll(c echo.Context) error {
	repos, err := client.Repositories(&prisma.RepositoriesParams{OrderBy: &orderBy}).Exec(c.Request().Context())
	if err != nil {
		return err
	}

	return writeNames(c, repos)
}

func getPaginated(c echo.Context, n int32, last string) error {
	lastRepo, err := client.Repositories(&prisma.RepositoriesParams{
		Where: &prisma.RepositoryWhereInput{
			FullName: prisma.Str(last),
		},
	}).Exec(c.Request().Context())
	if err != nil {
		return err
	}

	repos, err := client.Repositories(&prisma.RepositoriesParams{
		OrderBy: &orderBy,
		After:   prisma.Str(lastRepo[0].ID),
		First:   prisma.Int32(n),
	}).Exec(c.Request().Context())
	if err != nil {
		return err
	}
	link := fmt.Sprintf("/v2/_catalog?n=%d&last=%s", n, repos[len(repos)-1].FullName)
	c.Response().Header().Set(headerLink, link)
	return writeNames(c, repos)
}

func writeNames(c echo.Context, repos []prisma.Repository) error {
	names := make([]string, len(repos))
	for i, repo := range repos {
		names[i] = repo.FullName
	}
	res := make(map[string][]string)
	res["repositories"] = names
	return c.JSON(http.StatusOK, res)
}
