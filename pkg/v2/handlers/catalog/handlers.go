package catalog

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/db/models/repo"
)

const (
	headerLink = "Link"
)

var (
	getNames = repo.GetNames
)

// Get returns the catalog for the registry
func Get(c echo.Context) error {
	n, last, err := parsePagination(c)
	if err != nil {
		return err
	}
	names, nextLast, err := getNames(n, last)
	if err != nil {
		return err
	}
	if nextLast != "" {
		link := fmt.Sprintf("/v2/_catalog?n=%d&last=%s", n, nextLast)
		c.Response().Header().Set(headerLink, link)
	}

	res := make(map[string][]string)
	res["repositories"] = names
	return c.JSON(http.StatusOK, res)
}
