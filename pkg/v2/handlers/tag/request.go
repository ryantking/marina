package tag

// func parsePath(c echo.Context) (string, string, error) {
// 	repoName := c.Param("repo")
// 	orgName := c.Param("org")
// 	exists, err := repoExists(repoName, orgName)
// 	if err != nil {
// 		return "", "", errors.Wrap(err, "error checking if repository exists")
// 	}
// 	if !exists {
// 		c.Set("docker_err_code", docker.CodeNameUnknown)
// 		return "", "", echo.NewHTTPError(http.StatusNotFound, "no such repository")
// 	}
//
// 	return repoName, orgName, nil
// }
