package web

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/prisma"
)

var client = prisma.New(nil)

// GetRepo parses the repository and organization out of the path and returns the object
func GetRepo(c echo.Context) (*prisma.Repository, *prisma.Organization, error) {
	orgName := c.Param("org")
	org, err := client.Organization(prisma.OrganizationWhereUniqueInput{Name: &orgName}).Exec(c.Request().Context())
	if err == prisma.ErrNoResult {
		c.Set("docker_err_code", docker.CodeNameUnknown)
		return nil, nil, echo.NewHTTPError(http.StatusNotFound, "repository name not known to registry")
	}
	if err != nil {
		return nil, nil, err
	}
	repoName := c.Param("repo")
	repos, err := client.Repositories(&prisma.RepositoriesParams{
		Where: &prisma.RepositoryWhereInput{
			Name: &repoName,
			Org: &prisma.OrganizationWhereInput{
				Name: &orgName,
			},
		},
	}).Exec(context.Background())
	if err != nil {
		return nil, nil, err
	}
	if len(repos) == 0 {
		c.Set("docker_err_code", docker.CodeNameUnknown)
		return nil, nil, echo.NewHTTPError(http.StatusNotFound, "repository name not known to registry")
	}

	return &repos[0], org, nil
}

// GetImage gets the image from the request tag/digest
func GetImage(c echo.Context, repoID, orgID string) (*prisma.Image, error) {
	ref := c.Param("ref")
	var image *prisma.Image
	var err error
	if len(ref) == 71 && strings.HasPrefix(ref, "sha256:") {
		image, err = client.Image(prisma.ImageWhereUniqueInput{Digest: &ref}).Exec(c.Request().Context())
	} else {
		image, err = digestFromTag(c, ref, repoID, orgID)
	}
	if err != nil {
		return nil, err
	}

	return image, nil
}

func digestFromTag(c echo.Context, ref, repoID, orgID string) (*prisma.Image, error) {
	tags, err := client.Tags(&prisma.TagsParams{
		Where: &prisma.TagWhereInput{
			Ref: &ref,
			Image: &prisma.ImageWhereInput{
				Repo: &prisma.RepositoryWhereInput{
					ID:  &repoID,
					Org: &prisma.OrganizationWhereInput{ID: &orgID},
				},
			},
		},
	}).Exec(c.Request().Context())
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		c.Set("docker_err_code", docker.CodeNameUnknown)
		return nil, echo.NewHTTPError(http.StatusNotFound, "tag is unknown to registry")
	}

	return client.Tag(prisma.TagWhereUniqueInput{ID: &tags[0].ID}).Image().Exec(c.Request().Context())
}

// GetBlob parses the digest and returns the associated blob
func GetBlob(c echo.Context, repoID, orgID string) (*prisma.Blob, error) {
	digest := c.Param("digest")
	blobs, err := client.Blobs(&prisma.BlobsParams{
		Where: &prisma.BlobWhereInput{
			Digest: &digest,
			Repo: &prisma.RepositoryWhereInput{
				ID:  &repoID,
				Org: &prisma.OrganizationWhereInput{ID: &orgID},
			},
		},
	}).Exec(c.Request().Context())
	if err != nil {
		return nil, err
	}
	if len(blobs) == 0 {
		c.Set("docker_err_code", docker.CodeBlobUnknown)
		return nil, echo.NewHTTPError(http.StatusNotFound, "blob unknown to registry")
	}

	return &blobs[0], nil
}

// ParsePagination parses the docker page number and last from the request
func ParsePagination(c echo.Context) (int32, string, error) {
	s := c.QueryParam("n")
	if s == "" {
		return 0, "", nil
	}
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, "", echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return int32(n), c.QueryParam("last"), nil
}
