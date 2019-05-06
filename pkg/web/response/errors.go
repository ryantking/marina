package response

import (
	"net/http"

	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

// WriteError writes an error to the response with the given HTTP status code
func WriteError(resp *restful.Response, httpStatus int, err error) {
	err = resp.WriteError(httpStatus, err)
	if err != nil {
		log.WithError(err).Error("error writing to response")
	}
}

// SendServerError sends a server error response to the client
func SendServerError(resp *restful.Response, err error, msg string) {
	WriteError(resp, http.StatusInternalServerError, err)
	log.WithError(err).Error(msg)
}

// SendBadRequest sends a server error with a bad request to the client
func SendBadRequest(resp *restful.Response, err error) {
	WriteError(resp, http.StatusBadRequest, err)
}
