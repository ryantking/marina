package repository

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/ryantking/marina/pkg/db/models/layer"
	"github.com/ryantking/marina/pkg/db/models/repo"
	"github.com/ryantking/marina/pkg/db/models/tag"
	"github.com/ryantking/marina/pkg/db/models/upload"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/web/response"

	"github.com/emicklei/go-restful"
)

const (
	defaultOrg = "library"
)

func StartUpload(req *restful.Request, resp *restful.Response) {
	repoName, orgName := getOrgAndRepo(req)
	err := repo.CreateIfNotExists(repoName, orgName)
	if err != nil {
		response.SendServerError(resp, err, "error checking if organization and repo exist")
		return
	}

	upl, err := upload.New()
	if err != nil {
		response.SendServerError(resp, err, "error creating new upload")
		return
	}

	resp.AddHeader("Location", fmt.Sprintf("/v2/%s/%s/blobs/uploads/%d", orgName, repoName, upl.UUID))
	resp.AddHeader("Range", "0-0")
	resp.AddHeader("Content-Length", "0")
	resp.WriteHeader(http.StatusAccepted)
}

func UploadChunk(req *restful.Request, resp *restful.Response) {
	repoName, orgName := getOrgAndRepo(req)
	uuid := req.PathParameter("uuid")

	f, err := os.Create(fmt.Sprintf("%s_%s_%s.tar.gz", orgName, repoName, uuid))
	if err != nil {
		response.SendServerError(resp, err, "error creating file")
		return
	}
	n, err := io.Copy(f, req.Request.Body)
	if err != nil {
		response.SendServerError(resp, err, "error writing to file")
		return
	}

	resp.AddHeader("Location", fmt.Sprintf("/v2/%s/%s/blobs/uploads/%s", orgName, repoName, uuid))
	resp.AddHeader("Range", fmt.Sprintf("0-%d", n))
	resp.AddHeader("Content-Length", "0")
	resp.AddHeader("Docker-Upload-UUID", uuid)
	resp.WriteHeader(http.StatusAccepted)
}

func FinishUpload(req *restful.Request, resp *restful.Response) {
	digest := req.QueryParameter("digest")
	repoName, orgName := getOrgAndRepo(req)
	uuid, err := getUUID(req)
	if err != nil {
		response.SendBadRequest(resp, err)
		return
	}

	upl := &upload.Model{UUID: uuid, Done: true}
	err = upl.Save()
	if err != nil {
		response.SendServerError(resp, err, "error updating upload status")
		return
	}

	_, err = layer.New(digest, repoName, orgName)
	if err != nil {
		response.SendServerError(resp, err, "error creating new layer")
		return
	}

	resp.AddHeader("Location", fmt.Sprintf("/v2/%s/%s/blobs/%s", orgName, repoName, digest))
	resp.AddHeader("Content-Length", "0")
	resp.AddHeader("Docker-Content-Digest", digest)
	resp.WriteHeader(http.StatusCreated)
}

func LayerExists(req *restful.Request, resp *restful.Response) {
	digest := req.PathParameter("digest")
	repoName, orgName := getOrgAndRepo(req)
	exists, err := layer.Exists(digest, repoName, orgName)
	if err != nil {
		response.SendServerError(resp, err, "error checking if layer exists")
		return
	}
	if !exists {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	resp.AddHeader("Content-Length", fmt.Sprint(len(digest)))
	resp.AddHeader("Docker-Content-Digest", digest)
	resp.WriteHeader(http.StatusOK)
}

func GetManifest(req *restful.Request, resp *restful.Response) {
	repoName, orgName := getOrgAndRepo(req)
	ref := req.PathParameter("ref")
	manifest, err := tag.GetManifest(ref, repoName, orgName)
	if err == tag.ErrManifestNotFound {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		response.SendServerError(resp, err, "error finding tag manifest")
		return
	}

	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(manifest))
	resp.AddHeader("Docker-Content-Digest", digest)
	resp.AddHeader("Content-Type", docker.MIMEManifestV2)
	resp.WriteHeader(http.StatusOK)
	_, err = resp.Write(manifest)
	if err != nil {
		response.SendServerError(resp, err, "error writing manifest to response")
		return
	}
}

func UpdateManifest(req *restful.Request, resp *restful.Response) {
	ref := req.PathParameter("ref")
	repoName, orgName := getOrgAndRepo(req)
	manifest, err := ioutil.ReadAll(req.Request.Body)
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(manifest))
	if err != nil {
		response.SendBadRequest(resp, err)
		return
	}

	err = tag.UpdateManifest(ref, repoName, orgName, manifest)
	if err != nil {
		response.SendServerError(resp, err, "error updating manifest")
		return
	}

	resp.AddHeader("Location", fmt.Sprintf("/v2/%s/%s/manifests/%s", orgName, repoName, ref))
	resp.AddHeader("Content-Length", "0")
	resp.AddHeader("Docker-Content-Digest", digest)
	resp.WriteHeader(http.StatusCreated)
}

func getOrgAndRepo(req *restful.Request) (string, string) {
	repoName := req.PathParameter("repo")
	orgName := req.PathParameter("org")
	if orgName == "" {
		orgName = defaultOrg
	}
	return repoName, orgName
}

func getUUID(req *restful.Request) (uint64, error) {
	s := req.PathParameter("uuid")
	uuid, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return uuid, nil
}
