package catalog

import (
	"strconv"

	"github.com/labstack/echo"
)

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
