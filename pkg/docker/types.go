package docker

// Manifest is any docker manifest
type Manifest interface {
	Digest() string
}

// ManifestV2 is a docker manifest following the version2 type
type ManifestV2 struct {
	Size      string `json:"size"`
	MediaType string `json:"mediaType"`
	Config    struct {
		Size      uint64 `json:"size"`
		Digest    string `json:"digest"`
		MediaType string `json:"mediaType"`
	} `json:"config"`
	Layers []struct {
		Size      uint64 `json:"size"`
		Digest    string `json:"digest"`
		MediaType string `json:"mediaType"`
	} `json:"layers"`
}
