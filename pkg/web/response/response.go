package response

import (
	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

// WriteString writes a string to the response writer, handling any errro
func WriteString(resp *restful.Response, s string) {
	_, err := resp.Write([]byte(s))
	if err != nil {
		log.WithError(err).Error("error writing to response writer")
	}
}
