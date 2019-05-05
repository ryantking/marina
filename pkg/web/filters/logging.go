package filters

import (
	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

// RequestLogger is a filter that prints the method and path of all incoming requests
func RequestLogger(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Debugf("%s %s", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}
