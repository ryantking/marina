package base

import (
	"net/http"

	"github.com/labstack/echo"
)

// Get verifies that the server implements the Docker Registry API V2
func Get(c echo.Context) error {
	return c.String(http.StatusOK, "true")
}
