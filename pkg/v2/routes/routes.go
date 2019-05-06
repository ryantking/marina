package routes

import (
	"net/http"

	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/v2/handlers/registry"
	"github.com/ryantking/marina/pkg/v2/handlers/repository"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

func Registry() *restful.WebService {
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

	ws.Route(ws.HEAD("/{org}/{repo}/blobs/{digest}").To(repository.LayerExists))
	ws.Route(ws.HEAD("/{repo}/blobs/{digest}").To(repository.LayerExists))
	ws.Route(ws.GET("/{org}/{repo}/blobs/{digest}").To(repository.GetLayer))
	ws.Route(ws.GET("/{repo}/blobs/{digest}").To(repository.GetLayer))

	ws.Route(ws.POST("/{org}/{repo}/blobs/uploads").To(repository.StartUpload))
	ws.Route(ws.POST("/{repo}/blobs/uploads").To(repository.StartUpload))
	ws.Route(ws.PATCH("/{org}/{repo}/blobs/uploads/{uuid}").To(repository.UploadChunk))
	ws.Route(ws.PUT("/{org}/{repo}/blobs/uploads/{uuid}").To(repository.FinishUpload))

	ws.Route(ws.HEAD("/{repo}/manifests/{ref}").To(repository.ManifestExists))
	ws.Route(ws.HEAD("/{org}/{repo}/manifests/{ref}").To(repository.ManifestExists))
	ws.Route(ws.GET("/{repo}/manifests/{ref}").To(repository.GetManifest))
	ws.Route(ws.GET("/{org}/{repo}/manifests/{ref}").To(repository.GetManifest))
	ws.Route(ws.PUT("/{repo}/manifests/{ref}").To(repository.UpdateManifest))
	ws.Route(ws.PUT("/{org}/{repo}/manifests/{ref}").To(repository.UpdateManifest))

	return ws
}
