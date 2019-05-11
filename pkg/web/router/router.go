package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echolog "github.com/onrik/logrus/echo"
	"github.com/ryantking/marina/pkg/config"
	"github.com/ryantking/marina/pkg/docker"
	v2 "github.com/ryantking/marina/pkg/v2"
	"github.com/ryantking/marina/pkg/web"
	"github.com/sirupsen/logrus"
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
	e.HTTPErrorHandler = web.ErrorHandler
	v2.RegisterRoutes(e)

	return e
}
