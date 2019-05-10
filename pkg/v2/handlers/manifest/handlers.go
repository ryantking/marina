package manifest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db/models/image"
	"github.com/ryantking/marina/pkg/db/models/repo"
	"github.com/ryantking/marina/pkg/docker"
)

// Exists returns whether or not a manifest is present for the given digest
func Exists(c echo.Context) error {
	ref, repoName, orgName, err := parsePath(c)
	if err != nil {
		return err
	}
	manifest, _, err := image.GetManifest(ref, repoName, orgName)
	if err == image.ErrManifestNotFound {
		c.Set("docker_err_code", "NAME_UNKNOWN")
		return echo.NewHTTPError(http.StatusNotFound, "could not find manifest")
	}
	if err != nil {
		return errors.Wrap(err, "error retrieving manifest from database")
	}
	b, err := json.Marshal(manifest)
	if err != nil {
		return errors.Wrap(err, "error marshalling manifest")
	}
	c.Response().Header().Set(echo.HeaderContentLength, fmt.Sprint(len(b)))
	c.Response().Header().Set(docker.HeaderContentDigest, manifest.Digest())
	return c.NoContent(http.StatusOK)
}

// Get returns the manifest to the response, giving a 404 if it cannot be found
func Get(c echo.Context) error {
	ref, repoName, orgName, err := parsePath(c)
	if err != nil {
		return err
	}
	manifest, manifestType, err := image.GetManifest(ref, repoName, orgName)
	if err == image.ErrManifestNotFound {
		c.Set("docker_err_code", "NAME_UNKNOWN")
		return echo.NewHTTPError(http.StatusNotFound, "could not find manifest")
	}
	if err != nil {
		return errors.Wrap(err, "error retrieving manifest from database")
	}

	b, err := json.Marshal(manifest)
	if err != nil {
		return errors.Wrap(err, "error marshalling manifest")
	}
	c.Response().Header().Set(docker.HeaderContentDigest, manifest.Digest())
	return c.Blob(http.StatusOK, manifestType, b)
}

// Update updates a manifest in the database, creating if it does not currently exist
func Update(c echo.Context) error {
	ref, repoName, orgName, err := parsePath(c)
	if err != nil {
		return err
	}

	var manifest docker.Manifest
	manifestType := c.Request().Header.Get(echo.HeaderContentType)
	err = c.Bind(&manifest)
	if err != nil {
		return err
	}
	err = image.UpdateManifest(ref, repoName, orgName, manifest, manifestType)
	if err != nil {
		return errors.Wrap(err, "error updating manifest in database")
	}
	loc := fmt.Sprintf("/v2/%s/%s/manifests/%s", orgName, repoName, ref)
	c.Response().Header().Set(echo.HeaderLocation, loc)
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set(docker.HeaderContentDigest, manifest.Digest())
	return c.NoContent(http.StatusCreated)
}

// Delete deletes a manifest for a given digest
func Delete(c echo.Context) error {
	ref, _, _, err := parsePath(c)
	if err != nil {
		return err
	}

	err = image.Delete(ref)
	if err == image.ErrDeleteOnTag {
		c.Set("docker_err_code", "TAG_INVALID")
		return echo.NewHTTPError(http.StatusBadRequest, "cannot delete on tag")
	}
	if err != nil {
		return errors.Wrap(err, "error deleting manifest")
	}

	return c.NoContent(http.StatusAccepted)
}

func parsePath(c echo.Context) (string, string, string, error) {
	ref := c.Param("ref")
	repoName := c.Param("repo")
	orgName := c.Param("org")
	exists, err := repo.Exists(repoName, orgName)
	if err != nil {
		return "", "", "", errors.Wrap(err, "error checking if repository exists")
	}
	if !exists {
		c.Set("docker_err_code", "NAME_UNKNOWN")
		return "", "", "", echo.NewHTTPError(http.StatusNotFound, "no such repository")
	}

	return ref, repoName, orgName, nil
}
