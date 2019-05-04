package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/ryantking/marina/pkg/config"
	"github.com/ryantking/marina/pkg/logging"
	v2 "github.com/ryantking/marina/pkg/v2"

	"github.com/facebookgo/grace/gracehttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	logging.Init()
	cfg := config.Get()
	log.Infof("starting HTTP server at %s", cfg.RootURL)
	http.Handle("/v2", v2.NewRouter())
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
