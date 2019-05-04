package routes

import (
	"net/http"

	"github.com/ryantking/marina/pkg/v2/handlers/registry"
	"github.com/ryantking/marina/pkg/v2/handlers/repository"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

func Registry() *restful.WebService {
	tags := []string{"registry"}
	ws := new(restful.WebService)
	ws.Path("/v2").
		Consumes(restful.MIME_OCTET, restful.MIME_JSON, "application/vnd.docker.distribution.manifest.v1+prettyjws").
		Produces(restful.MIME_OCTET, restful.MIME_JSON)

	ws.Route(ws.GET("").To(registry.APIVersion).
		Doc("Find API version").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "Version supported", "true").
		Returns(http.StatusUnauthorized, "Unauthorized", nil))

	ws.Route(ws.POST("/{name}/blobs/uploads").To(repository.StartUpload))
	ws.Route(ws.PUT("/{name}/blobs/uploads/{uuid}").To(repository.FinishUpload))
	ws.Route(ws.GET("/{name}/manifests/{reference}").To(repository.GetManifest))
	ws.Route(ws.PUT("/{name}/manifests/{reference}").To(repository.UpdateManifest))
	ws.Route(ws.PATCH("/{name}/blobs/uploads/{uuid}").To(repository.DoUpload))
	ws.Route(ws.HEAD("/{name}/blobs/{digest}").To(repository.CheckDigest))

	return ws
}
