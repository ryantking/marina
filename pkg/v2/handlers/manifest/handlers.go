package manifest

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db/models/repo"
	"github.com/ryantking/marina/pkg/db/models/tag"
	"github.com/ryantking/marina/pkg/docker"
)

// Exists returns whether or not a manifest is present for the given digest
func Exists(c echo.Context) error {
	ref, repoName, orgName := parsePath(c)
	manifest, _, err := tag.GetManifest(ref, repoName, orgName)
	if err == tag.ErrManifestNotFound {
		c.Set("docker_err_code", "NAME_UNKNOWN")
		return echo.NewHTTPError(http.StatusNotFound, "could not find manifest")
	}
	if err != nil {
		return err
	}
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(manifest))
	c.Response().Header().Set(echo.HeaderContentLength, fmt.Sprint(len(manifest)))
	c.Response().Header().Set(docker.HeaderContentDigest, digest)
	return c.NoContent(http.StatusOK)
}

// Get returns the manifest to the response, giving a 404 if it cannot be found
func Get(c echo.Context) error {
	ref, repoName, orgName := parsePath(c)
	manifest, manifestType, err := tag.GetManifest(ref, repoName, orgName)
	if err == tag.ErrManifestNotFound {
		c.Set("docker_err_code", "NAME_UNKNOWN")
		return echo.NewHTTPError(http.StatusNotFound, "could not find manifest")
	}
	if err != nil {
		return err
	}

	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(manifest))
	c.Response().Header().Set(docker.HeaderContentDigest, digest)
	return c.Blob(http.StatusOK, manifestType, manifest)
}

// Update updates a manifest in the database, creating if it does not currently exist
func Update(c echo.Context) error {
	ref, repoName, orgName := parsePath(c)
	exists, err := repo.Exists(repoName, orgName)
	if err != nil {
		return err
	}
	if !exists {
		c.Set("docker_err_code", "NAME_UNKNOWN")
		return echo.NewHTTPError(http.StatusNotFound, "could not find the given repository")
	}
	manifestType := c.Request().Header.Get(echo.HeaderContentType)
	manifest, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.Set("docker_err_code", "MANIFEST_INVALID")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(manifest))
	err = tag.UpdateManifest(ref, repoName, orgName, manifest, manifestType)
	if err != nil {
		return errors.Wrap(err, "error updating manifest in database")
	}
	loc := fmt.Sprintf("/v2/%s/%s/manifests/%s", orgName, repoName, ref)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set(docker.HeaderContentDigest, digest)
	return c.NoContent(http.StatusCreated)
}

func parsePath(c echo.Context) (string, string, string) {
	return c.Param("ref"), c.Param("repo"), c.Param("org")
}
