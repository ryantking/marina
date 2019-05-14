package catalog

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/db/models/repo"
	"github.com/ryantking/marina/pkg/docker"
)

const (
	headerLink = "Link"
)

var (
	getNames          = repo.GetNames
	getNamesPaginated = repo.GetNamesPaginated
)

// Get returns the catalog for the registry
func Get(c echo.Context) error {
	n, last, err := docker.ParsePagination(c)
	if err != nil {
		return err
	}
	if n > 0 {
		return getPaginated(c, n, last)
	}

	return getAll(c)
}

func getAll(c echo.Context) error {
	names, err := getNames()
	if err != nil {
		return err
	}

	return writeNames(c, names)
}

func getPaginated(c echo.Context, n uint, last string) error {
	names, nextLast, err := getNamesPaginated(n, last)
	if err != nil {
		return err
	}
	if nextLast != "" {
		link := fmt.Sprintf("/v2/_catalog?n=%d&last=%s", n, nextLast)
		c.Response().Header().Set(headerLink, link)
	}

	return writeNames(c, names)
}

func writeNames(c echo.Context, names []string) error {
	res := make(map[string][]string)
	res["repositories"] = names
	return c.JSON(http.StatusOK, res)
}
