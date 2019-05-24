package tag

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/web"
)

const (
	headerLink = "Link"
)

// List lists all tags for a given repository
func List(c echo.Context) error {
	repo, org, err := web.GetRepo(c)
	if err != nil {
		return err
	}
	n, last, err := web.ParsePagination(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	tags, err := getTags(c.Request().Context(), repo.ID, n, last)
	if err != nil {
		return err
	}
	if n > 0 && len(tags) == int(n) {
		link := fmt.Sprintf("/v2/%s/%s/tags/list?n=%d&last=%s", org.Name, repo.Name, n, tags[n-1].Ref)
		c.Response().Header().Set(headerLink, link)
	}

	tagRefs := make([]string, len(tags))
	for i, tag := range tags {
		tagRefs[i] = tag.Ref
	}
	res := map[string]interface{}{}
	res["name"] = fmt.Sprintf("%s/%s", org.Name, repo.Name)
	res["tags"] = tagRefs
	return c.JSON(http.StatusOK, res)
}
