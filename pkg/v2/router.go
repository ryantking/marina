package v2

import (
	"fmt"
	"net/http"

	"github.com/ryantking/marina/pkg/v2/routes"

	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

// NewRouter creates a router for the v2 API with all routes
func NewRouter() http.Handler {
	container := restful.NewContainer()
	container.Add(routes.Registry())
	container.Filter(func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		log.Debugf("%s %s", req.Request.Method, req.Request.URL)
		fmt.Println(req.Request.Header)
		chain.ProcessFilter(req, resp)
	})

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders: []string{"Location", "Docker-Upload-UUID", "Content-Length", "Range"},
		// ExposeHeaders: []string{},
		AllowedHeaders: []string{"Content-Type", "Content-Length", "Content-Range", "Accept"},
		// AllowedHeaders: []string{},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		CookiesAllowed: false,
		Container:      container,
	}
	container.Filter(cors.Filter)

	return container
}
