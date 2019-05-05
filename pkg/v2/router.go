package v2

import (
	"net/http"

	"github.com/ryantking/marina/pkg/v2/routes"
	"github.com/ryantking/marina/pkg/web/filters"

	"github.com/emicklei/go-restful"
)

// NewRouter creates a router for the v2 API with all routes
func NewRouter() http.Handler {
	container := restful.NewContainer()
	container.Add(routes.Registry())
	container.Filter(filters.RequestLogger)
	container.Filter(filters.PanicRecovery)
	container.Filter(container.OPTIONSFilter)

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"Location", "Docker-Upload-UUID", "Content-Length", "Range"},
		AllowedHeaders: []string{"Content-Type", "Content-Length", "Content-Range", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		CookiesAllowed: false,
		Container:      container,
	}
	container.Filter(cors.Filter)

	return container
}
