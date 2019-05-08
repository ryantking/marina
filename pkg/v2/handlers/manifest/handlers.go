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
	"github.com/ryantking/marina/pkg/web/request"
	"github.com/ryantking/marina/pkg/web/response"

	"github.com/emicklei/go-restful"
)

// Exists returns whether or not a manifest is present for the given digest
func Exists(req *restful.Request, resp *restful.Response) {
	ref := req.PathParameter("ref")
	repoName, orgName := request.GetRepoAndOrg(req)
	manifest, _, err := tag.GetManifest(ref, repoName, orgName)
	if err == tag.ErrManifestNotFound {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		response.SendServerError(resp, err, "error getting manifest from database")
		return
	}
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(manifest))
	resp.AddHeader(response.ContentLength, fmt.Sprint(len(manifest)))
	resp.AddHeader(docker.HeaderContentDigest, digest)
	resp.WriteHeader(http.StatusOK)
}

// Get returns the manifest to the response, giving a 404 if it cannot be found
func Get(req *restful.Request, resp *restful.Response) {
	ref := req.PathParameter("ref")
	repoName, orgName := request.GetRepoAndOrg(req)
	manifest, manifestType, err := tag.GetManifest(ref, repoName, orgName)
	if err == tag.ErrManifestNotFound {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		response.SendServerError(resp, err, "error finding tag manifest")
		return
	}

	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(manifest))
	resp.AddHeader(docker.HeaderContentDigest, digest)
	resp.AddHeader(response.ContentType, manifestType)
	_, err = resp.Write(manifest)
	if err != nil {
		response.SendServerError(resp, err, "error writing manifest to response")
		return
	}
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
