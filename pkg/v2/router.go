package v2

import (
	"net/http"

	"github.com/emicklei/go-restful"
)

// NewRouter creates a router for the v2 API with all routes
func NewRouter() http.Handler {
	container := restful.NewContainer()
	return container
}
