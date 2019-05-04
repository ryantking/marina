package main

import (
	"net/http"

	"github.com/ryantking/marina/pkg/config"
	"github.com/ryantking/marina/pkg/logging"

	"github.com/facebookgo/grace/gracehttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	logging.Init()
	cfg := config.Get()
	log.Infof("starting HTTP server at %s", cfg.RootURL)
	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      http.DefaultServeMux,
		WriteTimeout: cfg.Server.WriteTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
	}
	if err := gracehttp.Serve(srv); err != nil {
		log.WithError(err).Error("fatal server error")
	}
	log.Info("HTTP server shut down")
}
