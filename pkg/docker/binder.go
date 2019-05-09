package docker

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
)

// Binder is the echo binder that allows docker manifests to be read
type Binder struct{}

func (b *Binder) Bind(i interface{}, c echo.Context) (err error) {
	// You may use default binder
	db := new(echo.DefaultBinder)
	if err = db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
		return
	}

	manifestType := c.Request().Header.Get(echo.HeaderContentType)
	switch manifestType {
	case MIMEManifestV2:
		man := i.(*Manifest)
		*man = new(ManifestV2)
		err = json.NewDecoder(c.Request().Body).Decode(i)
	default:
		c.Set("docker_err_code", CodeManifestInvalid)
		err = echo.NewHTTPError(http.StatusUnsupportedMediaType, "manifest type not supported")
	}
	return
}
