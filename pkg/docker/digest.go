package docker

import (
	"crypto/sha256"
	"fmt"
)

// MakeDigest takes in a series of bytes and returns it as a SHA256 docker digest
func MakeDigest(data []byte) string {
	return fmt.Sprintf("sha256:%x", sha256.Sum256(data))
}
