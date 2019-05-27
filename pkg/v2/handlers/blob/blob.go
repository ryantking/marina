package blob

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

func parseRange(c echo.Context) (int64, int64, error) {
	s := c.Request().Header.Get("Range")
	if s == "" {
		return -1, -1, nil
	}
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return 0, 0, echo.NewHTTPError(http.StatusBadRequest, "invalid Range header")
	}
	start, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return 0, 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	end, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return 0, 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return start, end, nil
}
