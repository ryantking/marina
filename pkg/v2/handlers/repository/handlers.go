package repository

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

var (
	globalUUID = 1
	uuidL      sync.Mutex
	uploaded   = map[string]int{}
	uuids      = map[string]int{}
	digests    = map[string]string{}
	manifests  = map[string][]byte{}
)

func StartUpload(req *restful.Request, resp *restful.Response) {
	name := req.PathParameter("name")
	uuidL.Lock()
	uuid := globalUUID
	globalUUID++
	uuidL.Unlock()
	resp.AddHeader("Location", fmt.Sprintf("/v2/%s/blobs/uploads/%d", name, uuid))
	resp.AddHeader("Range", "bytes=0-0")
	resp.AddHeader("Content-Length", "0")
	resp.AddHeader("Docker-Upload-UUID", fmt.Sprint(uuid))
	resp.WriteHeader(http.StatusAccepted)
}

func DoUpload(req *restful.Request, resp *restful.Response) {
	name := req.PathParameter("name")
	uuid := req.PathParameter("uuid")
	f, err := os.Create(fmt.Sprintf("%s.%s.tar.gz", name, uuid))
	if err != nil {
		log.WithError(err).Errorf("error opening file for writing")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	n, err := io.Copy(f, req.Request.Body)
	if err != nil {
		log.WithError(err).Errorf("error copying body")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.AddHeader("Location", fmt.Sprintf("/v2/%s/blobs/uploads/%s", name, uuid))
	resp.AddHeader("Range", fmt.Sprintf("0-%d", n))
	resp.AddHeader("Content-Length", "0")
	resp.AddHeader("Docker-Upload-UUID", uuid)
	resp.WriteHeader(http.StatusAccepted)
	uuids[uuid] = int(n)
}

func FinishUpload(req *restful.Request, resp *restful.Response) {
	name := req.PathParameter("name")
	digest := req.QueryParameter("digest")
	uuid := req.PathParameter("uuid")
	resp.AddHeader("Location", fmt.Sprintf("/v2/%s/blobs/%s", name, digest))
	resp.AddHeader("Content-Length", "0")
	resp.AddHeader("Docker-Content-Digest", digest)
	resp.WriteHeader(http.StatusCreated)
	uploaded[digest] = uuids[uuid]
	digests[name] = digest
}

func CheckDigest(req *restful.Request, resp *restful.Response) {
	n, ok := uploaded[req.PathParameter("digest")]
	if ok {
		resp.AddHeader("Content-Length", fmt.Sprint(n))
		resp.WriteHeader(http.StatusOK)
		return
	}
	resp.WriteHeader(http.StatusNotFound)
}

func GetManifest(req *restful.Request, resp *restful.Response) {
	name := req.PathParameter("name")
	manifest, ok := manifests[name]
	if !ok {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	resp.AddHeader("Content-Type", restful.MIME_OCTET)
	resp.Write(manifest)
}

func UpdateManifest(req *restful.Request, resp *restful.Response) {
	name := req.PathParameter("name")
	reference := req.PathParameter("reference")
	b, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		log.WithError(err).Error("error reading body")
		resp.WriteHeader(http.StatusInternalServerError)
	}
	manifests[name] = b
	resp.AddHeader("Location", fmt.Sprintf("/v2/%s/manifests/%s", name, reference))
	resp.AddHeader("Content-Length", "0")
	resp.AddHeader("Docker-Content-Digest", digests[name])
	resp.WriteHeader(http.StatusCreated)
}
