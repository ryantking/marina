package routes

import (
	"net/http"

	"github.com/ryantking/marina/pkg/v2/handlers/registry"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

func Registry() *restful.WebService {
	tags := []string{"registry"}
	ws := new(restful.WebService)
	ws.Path("/v2")

	ws.Route(ws.GET("").To(registry.APIVersion).
		Doc("Find API version").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "Version supported", "true").
		Returns(http.StatusUnauthorized, "Unauthorized", nil))

	return ws
}
