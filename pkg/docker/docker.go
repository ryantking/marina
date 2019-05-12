package docker

import (
	"strconv"

	"github.com/labstack/echo"
)

const (
	// HeaderAPIVersion is the header used to announce to clients which API version the server supports
	HeaderAPIVersion = "Docker-Distribution-API-Version"

	// APIVersion is the version  of the API that the server supports
	APIVersion = "registry/2.0"

	// HeaderContentDigest is the header used to pass the digest
	HeaderContentDigest = "Docker-Content-Digest"

	// HeaderUploadUUID is used to pass the upload UUID back to the client
	HeaderUploadUUID = "Docker-Upload-UUID"

	// Docker MIME types for responses
	MIMEManifestV1      = "application/vnd.docker.distribution.manifest.v1+prettyjws"
	MIMEManifestV2      = "application/vnd.docker.distribution.manifest.v2+json"
	MIMEManifestListV2  = "application/vnd.docker.distribution.manifest.list.v2+json"
	MIMEImageManifestV1 = "application/json application/vnd.oci.image.manifest.v1+json"
	MIMEImageIndexV1    = "application/vnd.oci.image.index.v1+json"
)

// ParsePagination parses the docker page number and last from the request
func ParsePagination(c echo.Context) (uint, string, error) {
	s := c.QueryParam("n")
	if s == "" {
		return 0, "", nil
	}
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, "", err
	}

	return uint(n), c.QueryParam("last"), nil
}
