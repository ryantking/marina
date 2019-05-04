package v2

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/ryantking/marina/pkg/v2/routes"
)

// NewRouter creates a router for the v2 API with all routes
func NewRouter() http.Handler {
	container := restful.NewContainer()
	container.Add(routes.Registry())

	return container
}
