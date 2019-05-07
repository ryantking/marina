package manifest

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"

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
func Update(req *restful.Request, resp *restful.Response) {
	ref := req.PathParameter("ref")
	repoName, orgName := request.GetRepoAndOrg(req)
	manifest, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		response.SendBadRequest(resp, err)
		return
	}
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(manifest))
	manifestType := req.HeaderParameter(response.ContentType)
	if err != nil {
		response.SendBadRequest(resp, err)
		return
	}

	err = tag.UpdateManifest(ref, repoName, orgName, manifest, manifestType)
	if err != nil {
		response.SendServerError(resp, err, "error updating manifest")
		return
	}

	resp.AddHeader(response.Location, fmt.Sprintf("/v2/%s/%s/manifests/%s", orgName, repoName, ref))
	resp.AddHeader(response.ContentLength, "0")
	resp.AddHeader(docker.HeaderContentDigest, digest)
	resp.WriteHeader(http.StatusCreated)
}
