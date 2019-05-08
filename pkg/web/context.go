package web

import "github.com/labstack/echo"

const defaultOrg = "library"

// OrgAndRepo returns the organization and repo from the gin context, defaulting org to library
func OrgAndRepo(c echo.Context) (string, string) {
	org := c.Param("org")
	if org == "" {
		org = defaultOrg
	}
	repo := c.Param("repo")
	return org, repo
}
