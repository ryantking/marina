package tag

const (
	headerLink = "Link"
)

// List lists all tags for a given repository
// func List(c echo.Context) error {
// 	repoName, orgName, err := parsePath(c)
// 	if err != nil {
// 		return err
// 	}
// 	pageSize, last, err := docker.ParsePagination(c)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}
// 	tags, nextLast, err := tagList(repoName, orgName, pageSize, last)
// 	if err != nil {
// 		return errors.Wrap(err, "error getting list of tags")
// 	}
//
// 	if nextLast != "" {
// 		link := fmt.Sprintf("/v2/%s/%s/tags/list?n=%d&last=%s", orgName, repoName, pageSize, nextLast)
// 		c.Response().Header().Set(headerLink, link)
// 	}
// 	res := map[string]interface{}{}
// 	res["name"] = fmt.Sprintf("%s/%s", orgName, repoName)
// 	res["tags"] = tags
// 	return c.JSON(http.StatusOK, res)
// }
