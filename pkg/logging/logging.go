package logging

import (
	"os"

	"github.com/ryantking/marina/pkg/config"
	log "github.com/sirupsen/logrus"
)

// Init initializes logrus
func Init() {
	log.SetOutput(os.Stdout)
	if config.Get().Environment == "DEVELOPMENT" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
