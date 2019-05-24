package manifest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/prisma"
	"github.com/ryantking/marina/pkg/web"
)

var client = prisma.New(nil)

// Exists returns whether or not a manifest is present for the given digest
func Exists(c echo.Context) error {
	repo, org, err := web.GetRepo(c)
	if err != nil {
		return err
	}
	image, err := web.GetImage(c, repo.ID, org.ID)
	if err != nil {
		return err
	}
	manifest, err := docker.ParseManifest(image.ManifestType, []byte(image.Manifest))
	if err != nil {
		return err
	}
	c.Response().Header().Set(echo.HeaderContentLength, fmt.Sprint(len(image.Manifest)))
	c.Response().Header().Set(docker.HeaderContentDigest, manifest.Digest())
	return c.NoContent(http.StatusOK)
}

// Get returns the manifest to the response, giving a 404 if it cannot be found
func Get(c echo.Context) error {
	repo, org, err := web.GetRepo(c)
	if err != nil {
		return err
	}
	image, err := web.GetImage(c, repo.ID, org.ID)
	if err != nil {
		return err
	}
	manifest, err := docker.ParseManifest(image.ManifestType, []byte(image.Manifest))
	if err != nil {
		return err
	}

	c.Response().Header().Set(docker.HeaderContentDigest, manifest.Digest())
	return c.Blob(http.StatusOK, image.ManifestType, []byte(image.Manifest))
}

// Update updates a manifest in the database, creating if it does not currently exist
func Update(c echo.Context) error {
	ref := c.Param("ref")
	repo, org, err := web.GetRepo(c)
	if err != nil {
		return err
	}

	var manifest docker.Manifest
	err = c.Bind(&manifest)
	if err != nil {
		return err
	}
	image, err := updateImage(c, manifest, repo.ID)
	if err != nil {
		return err
	}
	if len(ref) != 71 || !strings.HasPrefix(ref, "sha256:") {
		err = replaceTag(c.Request().Context(), ref, image.ID)
		if err != nil {
			return err
		}
	}

	loc := fmt.Sprintf("/v2/%s/%s/manifests/%s", org.Name, repo.Name, ref)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set(docker.HeaderContentDigest, manifest.Digest())
	return c.NoContent(http.StatusCreated)
}

// Delete deletes a manifest for a given digest
func Delete(c echo.Context) error {
	repo, org, err := web.GetRepo(c)
	if err != nil {
		return err
	}
	image, err := web.GetImage(c, repo.ID, org.ID)
	if err != nil {
		return err
	}

	ref := c.Param("ref")
	if len(ref) != 71 || !strings.HasPrefix(ref, "sha256:") {
		err = deleteTag(c.Request().Context(), ref, image.ID)
		if err != nil {
			return err
		}
	}

	_, err = client.DeleteImage(prisma.ImageWhereUniqueInput{ID: &image.ID}).Exec(c.Request().Context())
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}
