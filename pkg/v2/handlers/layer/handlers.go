package layer

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ryantking/marina/pkg/db/models/layer"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/web/request"
	"github.com/ryantking/marina/pkg/web/response"

	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

// Exists returns a status 200 or 404 depending on if the layer exists or not
func Exists(req *restful.Request, resp *restful.Response) {
	digest := req.PathParameter("digest")
	repoName, orgName := request.GetRepoAndOrg(req)
	exists, err := layer.Exists(digest, repoName, orgName)
	if err != nil {
		response.SendServerError(resp, err, "error checking if layer exists")
		return
	}
	if !exists {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	resp.AddHeader(response.ContentLength, fmt.Sprint(len(digest)))
	resp.AddHeader(docker.HeaderContentDigest, digest)
	resp.WriteHeader(http.StatusOK)
}

// Get downloads a layer to the response body, responding 404 if the layer is not found
func Get(req *restful.Request, resp *restful.Response) {
	repoName, orgName := request.GetRepoAndOrg(req)
	digest := req.PathParameter("digest")
	exists, err := layer.Exists(digest, repoName, orgName)
	if err != nil {
		response.SendServerError(resp, err, "error checking if layer exists")
		return
	}
	if !exists {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	f, err := os.Open(fmt.Sprintf("%s_%s_%s.tar.gz", orgName, repoName, digest))
	if err != nil {
		response.SendServerError(resp, err, "error opening layer file")
		return
	}

	resp.AddHeader(response.ContentType, restful.MIME_OCTET)
	_, err = io.Copy(resp, f)
	if err != nil {
		log.WithError(err).Error("error copying layer to response body")
	}
}
