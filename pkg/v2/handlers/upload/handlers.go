package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ryantking/marina/pkg/db/models/layer"
	"github.com/ryantking/marina/pkg/db/models/repo"
	"github.com/ryantking/marina/pkg/db/models/upload"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/web/request"
	"github.com/ryantking/marina/pkg/web/response"

	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

// Start starts the process of uploading a new layer
func Start(req *restful.Request, resp *restful.Response) {
	log.Debug("start upload")
	repoName, orgName := request.GetRepoAndOrg(req)
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

	resp.AddHeader(response.Location, fmt.Sprintf("/v2/%s/%s/blobs/uploads/%d", orgName, repoName, upl.UUID))
	resp.AddHeader(response.ContentRange, "0-0")
	resp.AddHeader(response.ContentLength, "0")
	resp.WriteHeader(http.StatusAccepted)
}

// Chunk is used to save a chunk of data to a layer
func Chunk(req *restful.Request, resp *restful.Response) {
	repoName, orgName := request.GetRepoAndOrg(req)
	uuid, err := request.GetUUID(req)
	if err != nil {
		response.SendBadRequest(resp, err)
		fmt.Println(err.Error())
		return
	}
	f, err := os.Create(fmt.Sprintf("upload_%d.tar.gz", uuid))
	if err != nil {
		response.SendServerError(resp, err, "error creating file")
		return
	}
	n, err := io.Copy(f, req.Request.Body)
	if err != nil {
		response.SendServerError(resp, err, "error writing to file")
		return
	}

	resp.AddHeader(response.Location, fmt.Sprintf("/v2/%s/%s/blobs/uploads/%d", orgName, repoName, uuid))
	resp.AddHeader(request.Range, fmt.Sprintf("0-%d", n))
	resp.AddHeader(response.ContentLength, "0")
	resp.AddHeader(docker.HeaderUploadUUID, fmt.Sprint(uuid))
	resp.WriteHeader(http.StatusAccepted)
}

// Finish is used to signal when the upload has completed
func Finish(req *restful.Request, resp *restful.Response) {
	digest := req.QueryParameter("digest")
	repoName, orgName := request.GetRepoAndOrg(req)
	uuid, err := request.GetUUID(req)
	if err != nil {
		response.SendBadRequest(resp, err)
		return
	}

	os.Rename(fmt.Sprintf("upload_%d.tar.gz", uuid), fmt.Sprintf("%s_%s_%s.tar.gz", orgName, repoName, digest))
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

	resp.AddHeader(response.Location, fmt.Sprintf("/v2/%s/%s/blobs/%s", orgName, repoName, digest))
	resp.AddHeader(response.ContentLength, "0")
	resp.AddHeader(docker.HeaderContentDigest, digest)
	resp.WriteHeader(http.StatusCreated)
}
