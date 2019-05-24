package blob

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/prisma"
	"github.com/ryantking/marina/pkg/store"
	"github.com/ryantking/marina/pkg/web"
)

var (
	client = prisma.New(nil)

	getBlob    = store.GetBlob
	deleteBlob = store.DeleteBlob
)

// Exists returns a status 200 or 404 depending on if the blob exists or not
func Exists(c echo.Context) error {
	repo, org, err := web.GetRepo(c)
	if err != nil {
		return err
	}
	blob, err := web.GetBlob(c, repo.ID, org.ID)
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentLength, fmt.Sprint(len(blob.Digest)))
	c.Response().Header().Set(docker.HeaderContentDigest, blob.Digest)
	return c.NoContent(http.StatusOK)
}

// Get downloads a blob to the response body, responding 404 if the blob is not found
func Get(c echo.Context) error {
	repo, org, err := web.GetRepo(c)
	if err != nil {
		return err
	}
	blob, err := web.GetBlob(c, repo.ID, org.ID)
	if err != nil {
		return err
	}

	r, err := getBlob(blob.Digest, repo.Name, org.Name)
	if err != nil {
		return err
	}
	defer r.Close()

	c.Response().Header().Set(docker.HeaderContentDigest, blob.Digest)
	return c.Stream(http.StatusOK, echo.MIMEOctetStream, r)
}

// Delete deletes a blob from the repository
func Delete(c echo.Context) error {
	repo, org, err := web.GetRepo(c)
	if err != nil {
		return err
	}
	blob, err := web.GetBlob(c, repo.ID, org.ID)
	if err != nil {
		return err
	}

	err = deleteBlob(blob.Digest, repo.Name, org.Name)
	if err != nil {
		return err
	}
	_, err = client.DeleteBlob(prisma.BlobWhereUniqueInput{Digest: &blob.Digest}).Exec(c.Request().Context())
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}
