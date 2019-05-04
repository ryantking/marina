package registry

import (
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/web/response"

	"github.com/emicklei/go-restful"
)

// APIVersion returns whether a 200 status to indicate that this API supports v2
func APIVersion(req *restful.Request, resp *restful.Response) {
	resp.AddHeader(docker.HeaderAPIVersion, docker.APIVersion)
	response.WriteString(resp, "true")
}
