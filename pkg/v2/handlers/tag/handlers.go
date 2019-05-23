package tag

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/prisma"
	"github.com/ryantking/marina/pkg/web"
)

const (
	headerLink = "Link"
)

var (
	client  = prisma.New(nil)
	orderBy = prisma.TagOrderByInputUpdatedAtDesc
)

// List lists all tags for a given repository
func List(c echo.Context) error {
	repo, org, err := web.ParsePath(c)
	if err != nil {
		return err
	}
	n, last, err := web.ParsePagination(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	tagsParams := &prisma.TagsParams{
		OrderBy: &orderBy,
		Where: &prisma.TagWhereInput{
			Image: &prisma.ImageWhereInput{
				Repo: &prisma.RepositoryWhereInput{
					Name: &repo,
					Org:  &prisma.OrganizationWhereInput{Name: &org},
				},
			},
		},
	}
	if n != 0 {
		tagsParams.First = &n
	}
	if last != "" {
		var tags []prisma.Tag
		tags, err = client.Tags(&prisma.TagsParams{
			Where: &prisma.TagWhereInput{
				Ref: &last,
				Image: &prisma.ImageWhereInput{
					Repo: &prisma.RepositoryWhereInput{
						Name: &repo,
						Org:  &prisma.OrganizationWhereInput{Name: &org},
					},
				},
			},
		}).Exec(c.Request().Context())
		if err != nil {
			return err
		}
		tagsParams.After = &tags[0].ID
	}
	tags, err := client.Tags(tagsParams).Exec(c.Request().Context())
	if err != nil {
		return err
	}
	if n > 0 && len(tags) == int(n) {
		link := fmt.Sprintf("/v2/%s/%s/tags/list?n=%d&last=%s", org, repo, n, tags[n-1].Ref)
		c.Response().Header().Set(headerLink, link)
	}

	tagRefs := make([]string, len(tags))
	for i, tag := range tags {
		tagRefs[i] = tag.Ref
	}
	res := map[string]interface{}{}
	res["name"] = fmt.Sprintf("%s/%s", org, repo)
	res["tags"] = tagRefs
	return c.JSON(http.StatusOK, res)
}
