package request

import (
	"strconv"

	"github.com/emicklei/go-restful"
)

const (
	defaultOrg = "library"

	// Range is the header used to specify range
	Range = "Range"
)

// GetRepoAndOrg returns the repository and organization from the request
func GetRepoAndOrg(req *restful.Request) (string, string) {
	repoName := req.PathParameter("repo")
	orgName := req.PathParameter("org")
	if orgName == "" {
		orgName = defaultOrg
	}
	return repoName, orgName
}

// GetUUID returns UUID from a request, parsed as an integer
func GetUUID(req *restful.Request) (uint64, error) {
	s := req.PathParameter("uuid")
	uuid, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return uuid, nil
}
