package v2

import (
	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/v2/handlers/base"
	"github.com/ryantking/marina/pkg/v2/handlers/blob"
	"github.com/ryantking/marina/pkg/v2/handlers/catalog"
	"github.com/ryantking/marina/pkg/v2/handlers/manifest"
	"github.com/ryantking/marina/pkg/v2/handlers/tag"
	"github.com/ryantking/marina/pkg/v2/handlers/upload"
)

func RegisterRoutes(e *echo.Echo) {
	v2 := e.Group("/v2")
	v2.GET("/", base.Get)
	v2.GET("/:org/:repo/tags/list", tag.List)
	v2.HEAD("/:org/:repo/blobs/:digest", blob.Exists)
	v2.GET("/:org/:repo/blobs/:digest", blob.Get)
	v2.DELETE("/:org/:repo/blobs/:digest", blob.Delete)
	v2.POST("/:org/:repo/blobs/uploads/", upload.Start)
	v2.PUT("/:org/:repo/blobs/uploads/:uuid", upload.Finish)
	v2.PATCH("/:org/:repo/blobs/uploads/:uuid", upload.Blob)
	v2.HEAD("/:org/:repo/manifests/:ref", manifest.Exists)
	v2.GET("/:org/:repo/manifests/:ref", manifest.Get)
	v2.PUT("/:org/:repo/manifests/:ref", manifest.Update)
	v2.DELETE("/:org/:repo/manifests/:ref", manifest.Delete)
	v2.GET("/_catalog", catalog.Get)
}
