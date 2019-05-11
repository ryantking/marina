package blob

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/db/models/blob"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/store"
)

var (
	getBlob         = store.GetBlob
	deleteBlob      = store.DeleteBlob
	deleteBlobEntry = blob.Delete
)

// Exists returns a status 200 or 404 depending on if the blob exists or not
func Exists(c echo.Context) error {
	digest, _, _, err := parsePath(c)
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentLength, fmt.Sprint(len(digest)))
	c.Response().Header().Set(docker.HeaderContentDigest, digest)
	return c.NoContent(http.StatusOK)
}

// Get downloads a blob to the response body, responding 404 if the blob is not found
func Get(c echo.Context) error {
	digest, repoName, orgName, err := parsePath(c)
	if err != nil {
		return err
	}

	r, err := getBlob(digest, repoName, orgName)
	if err != nil {
		return err
	}
	defer r.Close()

	return c.Stream(http.StatusOK, echo.MIMEOctetStream, r)
}

// Delete deletes a blob from the repository
func Delete(c echo.Context) error {
	digest, repoName, orgName, err := parsePath(c)
	if err != nil {
		return err
	}

	err = deleteBlob(digest, repoName, orgName)
	if err != nil {
		return err
	}
	err = deleteBlobEntry(digest)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}
