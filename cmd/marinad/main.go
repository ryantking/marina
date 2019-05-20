package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/ryantking/marina/pkg/config"
	"github.com/ryantking/marina/pkg/logging"
	"github.com/ryantking/marina/pkg/web/router"
	log "github.com/sirupsen/logrus"
)

func main() {
	logging.Init()
	cfg := config.Get()
	log.Infof("starting HTTP server at %s", cfg.RootURL)
	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router.New(),
		WriteTimeout: cfg.Server.WriteTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
	}
	if err := gracehttp.Serve(srv); err != nil {
		log.WithError(err).Error("fatal server error")
	}
	log.Info("HTTP server shut down")
}
