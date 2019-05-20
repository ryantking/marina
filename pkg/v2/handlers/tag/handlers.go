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
	if n == 0 {
		return listAll(c, repo, org)
	}

	return listPaginated(c, repo, org, n, last)
}

func listAll(c echo.Context, repo, org string) error {
	tags, err := client.Tags(&prisma.TagsParams{
		OrderBy: &orderBy,
		Where: &prisma.TagWhereInput{
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

	return writeTags(c, repo, org, tags)
}

func listPaginated(c echo.Context, repo, org string, n int32, last string) error {
	tags, err := client.Tags(&prisma.TagsParams{
		OrderBy: &orderBy,
		After:   &last,
		First:   &n,
		Where: &prisma.TagWhereInput{
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

	link := fmt.Sprintf("/v2/%s/%s/tags/list?n=%d&last=%s", org, repo, n, tags[len(tags)-1].Ref)
	c.Response().Header().Set(headerLink, link)
	return writeTags(c, repo, org, tags)
}

func writeTags(c echo.Context, repo, org string, tags []prisma.Tag) error {
	tagRefs := make([]string, len(tags))
	for i, tag := range tags {
		tagRefs[i] = tag.Ref
	}
	res := map[string]interface{}{}
	res["name"] = fmt.Sprintf("%s/%s", org, repo)
	res["tags"] = tagRefs
	return c.JSON(http.StatusOK, res)
}
