package filters

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/ryantking/marina/pkg/web/response"
	log "github.com/sirupsen/logrus"
)

// PanicRecovery handles a panic on any incoming request gracefully
func PanicRecovery(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	defer func() {
		r := recover()
		switch rval := r.(type) {
		case nil:
			return
		case error:
			response.WriteError(resp, http.StatusInternalServerError, rval)
			log.WithError(rval).WithField("path", req.Request.URL.Path).Errorf("server panic")
		default:
			err := errors.New(fmt.Sprint(rval))
			response.WriteError(resp, http.StatusInternalServerError, err)
			log.WithError(err).WithField("path", req.Request.URL.Path).Errorf("server panic")
		}
	}()

	chain.ProcessFilter(req, resp)
}
