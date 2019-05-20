package manifest

import (
	"encoding/json"
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
	image, err := parsePath(c)
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
	image, err := parsePath(c)
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
	repo, org, err := web.ParsePath(c)
	if err != nil {
		return err
	}

	manifestType := c.Request().Header.Get(echo.HeaderContentType)
	var manifest docker.Manifest
	err = c.Bind(&manifest)
	if err != nil {
		return err
	}
	b, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	repos, err := client.Repositories(&prisma.RepositoriesParams{
		Where: &prisma.RepositoryWhereInput{
			Name: &repo,
			Org: &prisma.OrganizationWhereInput{
				Name: &org,
			},
		},
	}).Exec(c.Request().Context())
	image, err := client.UpsertImage(prisma.ImageUpsertParams{
		Where: prisma.ImageWhereUniqueInput{Digest: prisma.Str(manifest.Digest())},
		Update: prisma.ImageUpdateInput{
			Manifest:     prisma.Str(string(b)),
			ManifestType: &manifestType,
		},
		Create: prisma.ImageCreateInput{
			Digest:       manifest.Digest(),
			Manifest:     string(b),
			ManifestType: manifestType,
			Repo: prisma.RepositoryCreateOneWithoutImagesInput{
				Connect: &prisma.RepositoryWhereUniqueInput{ID: &repos[0].ID},
			},
		},
	}).Exec(c.Request().Context())
	if err != nil {
		return err
	}
	if len(ref) != 71 || !strings.HasPrefix(ref, "sha256:") {
		_, err := client.DeleteManyTags(&prisma.TagWhereInput{
			Ref: &ref,
			Image: &prisma.ImageWhereInput{
				Repo: &prisma.RepositoryWhereInput{
					Name: &repo,
					Org:  &prisma.OrganizationWhereInput{Name: &org},
				},
			},
		}).Exec(c.Request().Context())
		if err != nil {
			return err
		}

		_, err = client.CreateTag(prisma.TagCreateInput{
			Ref: ref,
			Image: prisma.ImageCreateOneWithoutTagsInput{
				Connect: &prisma.ImageWhereUniqueInput{ID: &image.ID},
			},
		}).Exec(c.Request().Context())
	}

	loc := fmt.Sprintf("/v2/%s/%s/manifests/%s", org, repo, ref)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set(docker.HeaderContentDigest, manifest.Digest())
	return c.NoContent(http.StatusCreated)
}

// Delete deletes a manifest for a given digest
func Delete(c echo.Context) error {
	repo, org, err := web.ParsePath(c)
	if err != nil {
		return err
	}
	image, err := parsePath(c)
	if err != nil {
		return err
	}

	ref := c.Param("ref")
	if len(ref) != 71 || !strings.HasPrefix(ref, "sha256:") {
		_, err := client.DeleteManyTags(&prisma.TagWhereInput{
			Ref: &ref,
			Image: &prisma.ImageWhereInput{
				Repo: &prisma.RepositoryWhereInput{
					Name: &repo,
					Org:  &prisma.OrganizationWhereInput{Name: &org},
				},
			},
		}).Exec(c.Request().Context())
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
