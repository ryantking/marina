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
	digest, repo, org, err := parsePath(c)
	if err != nil {
		return err
	}

	r, err := getBlob(digest, repo, org)
	if err != nil {
		return err
	}
	defer r.Close()

	c.Response().Header().Set(docker.HeaderContentDigest, digest)
	return c.Stream(http.StatusOK, echo.MIMEOctetStream, r)
}

// Delete deletes a blob from the repository
func Delete(c echo.Context) error {
	digest, repo, org, err := parsePath(c)
	if err != nil {
		return err
	}

	err = deleteBlob(digest, repo, org)
	if err != nil {
		return err
	}
	_, err = client.DeleteBlob(prisma.BlobWhereUniqueInput{Digest: &digest}).Exec(c.Request().Context())
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}

func parsePath(c echo.Context) (string, string, string, error) {
	repo, org, err := web.ParsePath(c)
	if err != nil {
		return "", "", "", err
	}

	digest := c.Param("digest")
	blobs, err := client.Blobs(&prisma.BlobsParams{
		Where: &prisma.BlobWhereInput{
			Digest: prisma.Str(digest),
			Repo: &prisma.RepositoryWhereInput{
				Name: prisma.Str(repo),
				Org: &prisma.OrganizationWhereInput{
					Name: prisma.Str(org),
				},
			},
		},
	}).Exec(c.Request().Context())
	if err != nil {
		return "", "", "", err
	}
	if len(blobs) == 0 {
		c.Set("docker_err_code", docker.CodeBlobUnknown)
		return "", "", "", echo.NewHTTPError(http.StatusNotFound, "blob unknown to registry")
	}

	return digest, repo, org, nil
}
