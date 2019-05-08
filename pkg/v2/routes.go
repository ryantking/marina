package v2

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/v2/handlers/layer"
	"github.com/ryantking/marina/pkg/v2/handlers/manifest"
	"github.com/ryantking/marina/pkg/v2/handlers/upload"
)

func RegisterRoutes(e *echo.Echo) {
	v2 := e.Group("/v2")
	v2.GET("", func(c echo.Context) error {
		c.String(http.StatusOK, "true")
		return nil
	})

	v2.HEAD("/:org/:repo/blobs/:digest", layer.Exists)

	v2.POST("/:org/:repo/blobs/uploads/", upload.Start)
	v2.PUT("/:org/:repo/blobs/uploads/:uuid", upload.Finish)
	v2.PATCH("/:org/:repo/blobs/uploads/:uuid", upload.Chunk)

	v2.PUT("/:org/:repo/manifests/:ref", manifest.Update)

	// 	ws.Route(ws.GET("/{org}/{repo}/blobs/{digest}").To(layer.Get))

	// 	ws.Route(ws.HEAD("/{org}/{repo}/manifests/{ref}").To(manifest.Exists))
	// 	ws.Route(ws.GET("/{org}/{repo}/manifests/{ref}").To(manifest.Get))
	// 	ws.Route(ws.PUT("/{org}/{repo}/manifests/{ref}").To(manifest.Update))
	//
}
