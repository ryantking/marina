package web

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ErrorHandler handles errors and writes them in the docker format when needed
func ErrorHandler(err error, c echo.Context) {
	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		c.String(http.StatusInternalServerError, err.Error())
		log.WithError(errors.Cause(err)).Errorf(err.Error())
		return
	}

	dockerErrCode := c.Get("docker_err_code")
	if dockerErrCode == nil {
		c.String(httpErr.Code, httpErr.Message.(string))
		return
	}

	dockerErr := map[string]interface{}{"code": dockerErrCode, "message": httpErr.Message}
	if detail := c.Get("docker_err_detail"); detail != nil {
		dockerErr["detail"] = detail
	}
	c.JSON(httpErr.Code, map[string]interface{}{"errors": []interface{}{dockerErr}})
}
