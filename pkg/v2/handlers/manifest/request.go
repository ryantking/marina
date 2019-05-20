package manifest

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/prisma"
	"github.com/ryantking/marina/pkg/web"
)

func parsePath(c echo.Context) (*prisma.Image, error) {
	repo, org, err := web.ParsePath(c)
	if err != nil {
		return nil, err
	}

	var image *prisma.Image
	ref := c.Param("ref")
	if len(ref) == 71 && strings.HasPrefix(ref, "sha256:") {
		image, err = client.Image(prisma.ImageWhereUniqueInput{
			Digest: &ref,
		}).Exec(c.Request().Context())
	} else {
		image, err = digestFromTag(c, ref, repo, org)
	}

	if err == prisma.ErrNoResult {
		c.Set("docker_err_code", docker.CodeManifestUnknown)
		return nil, echo.NewHTTPError(http.StatusNotFound, "manifest unknown")
	}
	if err != nil {
		return nil, err
	}

	return image, nil
}

func digestFromTag(c echo.Context, ref, repo, org string) (*prisma.Image, error) {
	tags, err := client.Tags(&prisma.TagsParams{
		Where: &prisma.TagWhereInput{
			Ref: prisma.Str(ref),
			Image: &prisma.ImageWhereInput{
				Repo: &prisma.RepositoryWhereInput{
					Name: prisma.Str(repo),
					Org:  &prisma.OrganizationWhereInput{Name: prisma.Str(org)},
				},
			},
		},
	}).Exec(c.Request().Context())
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		c.Set("docker_err_code", docker.CodeNameUnknown)
		return nil, echo.NewHTTPError(http.StatusNotFound, "repository name not known to registry")
	}

	return client.Tag(prisma.TagWhereUniqueInput{ID: &tags[0].ID}).Image().Exec(c.Request().Context())
}
