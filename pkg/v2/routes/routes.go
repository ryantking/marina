package routes

import (
	"net/http"

	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/v2/handlers/layer"
	"github.com/ryantking/marina/pkg/v2/handlers/manifest"
	"github.com/ryantking/marina/pkg/v2/handlers/registry"
	"github.com/ryantking/marina/pkg/v2/handlers/upload"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

// WebService returns a web service with all the routes for the docker V2 API
func WebService() *restful.WebService {
	tags := []string{"registry"}
	ws := new(restful.WebService)
	ws.Path("/v2").
		Consumes(restful.MIME_OCTET, docker.MIMEManifestV2).
		Produces(
			restful.MIME_OCTET, restful.MIME_JSON, docker.MIMEManifestV1, docker.MIMEManifestV2,
			docker.MIMEManifestListV2, docker.MIMEImageManifestV1, docker.MIMEImageIndexV1,
		)

	ws.Route(ws.GET("").To(registry.APIVersion).
		Doc("Find API version").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "Version supported", "true").
		Returns(http.StatusUnauthorized, "Unauthorized", nil))

	ws.Route(ws.HEAD("/{org}/{repo}/blobs/{digest}").To(layer.Exists))
	ws.Route(ws.HEAD("/{repo}/blobs/{digest}").To(layer.Exists))
	ws.Route(ws.GET("/{org}/{repo}/blobs/{digest}").To(layer.Get))
	ws.Route(ws.GET("/{repo}/blobs/{digest}").To(layer.Get))

	ws.Route(ws.POST("/{org}/{repo}/blobs/uploads").To(upload.Start))
	ws.Route(ws.POST("/{repo}/blobs/uploads").To(upload.Start))
	ws.Route(ws.PATCH("/{org}/{repo}/blobs/uploads/{uuid}").To(upload.Chunk))
	ws.Route(ws.PATCH("/{repo}/blobs/uploads/{uuid}").To(upload.Chunk))
	ws.Route(ws.PUT("/{org}/{repo}/blobs/uploads/{uuid}").To(upload.Finish))
	ws.Route(ws.PUT("/{repo}/blobs/uploads/{uuid}").To(upload.Finish))

	ws.Route(ws.HEAD("/{repo}/manifests/{ref}").To(manifest.Exists))
	ws.Route(ws.HEAD("/{org}/{repo}/manifests/{ref}").To(manifest.Exists))
	ws.Route(ws.GET("/{repo}/manifests/{ref}").To(manifest.Get))
	ws.Route(ws.GET("/{org}/{repo}/manifests/{ref}").To(manifest.Get))
	ws.Route(ws.PUT("/{repo}/manifests/{ref}").To(manifest.Update))
	ws.Route(ws.PUT("/{org}/{repo}/manifests/{ref}").To(manifest.Update))

	return ws
}
