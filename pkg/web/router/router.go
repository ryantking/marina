package router

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echolog "github.com/onrik/logrus/echo"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/config"
	"github.com/ryantking/marina/pkg/docker"
	v2 "github.com/ryantking/marina/pkg/v2"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var logConfig = middleware.LoggerConfig{Format: "[ECHO] ${method} ${path} | ${status}\n"}

// New creates a router for the v2 API with all routes
func New() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Logger = echolog.NewLogger(logrus.StandardLogger(), "")
	e.Binder = new(docker.Binder)
	if config.Get().Environment == "DEVELOPMENT" {
		e.Use(middleware.LoggerWithConfig(logConfig))
	}
	e.HTTPErrorHandler = errHandler
	v2.RegisterRoutes(e)

	return e
}

func errHandler(err error, c echo.Context) {
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
