package docker

const (
	// HeaderAPIVersion is the header used to announce to clients which API version the server supports
	HeaderAPIVersion = "Docker-Distribution-API-Version"

	// APIVersion is the version  of the API that the server supports
	APIVersion = "registry/2.0"

	// Docker MIME types for responses
	MIMEManifestV1      = "application/vnd.docker.distribution.manifest.v1+prettyjws"
	MIMEManifestV2      = "application/vnd.docker.distribution.manifest.v2+json"
	MIMEManifestListV2  = "application/vnd.docker.distribution.manifest.list.v2+json"
	MIMEImageManifestV1 = "application/json application/vnd.oci.image.manifest.v1+json"
	MIMEImageIndexV1    = "application/vnd.oci.image.index.v1+json"
)
