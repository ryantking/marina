package docker

import (
	"encoding/json"
	"errors"
)

var (
	// ErrUnsupportedManifestType
	ErrUnsupportedManifestType = errors.New("unsupported manifest type")
)

func ParseManifest(manifestType string, data []byte) (Manifest, error) {
	switch manifestType {
	case MIMEManifestV2:
		man := new(ManifestV2)
		err := json.Unmarshal(data, man)
		if err != nil {
			return nil, err
		}
		return man, nil
	}

	return nil, ErrUnsupportedManifestType
}

// Digest returns the digest of a docker V2 manifest
func (man *ManifestV2) Digest() string {
	return man.Config.Digest
}
