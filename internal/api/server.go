// Package api provides functionality for initializing a server.
package api

import (
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"time"

	"webServerGoLiftover/internal/api/handlers"
	"webServerGoLiftover/internal/api/middleware"
	"webServerGoLiftover/internal/config"
	"webServerGoLiftover/internal/service/liftover/v1"
)

// InitServer returns a http.Server object ready to be listening and serving.
func InitServer(cfg *config.ServerConfig, app *config.Application, logger *zerolog.Logger) *http.Server {
	liftoverService := liftover.InitLiftover(logger)
	urlHandler := handlers.InitURLHandler(cfg, app, liftoverService, logger)
	mux := http.NewServeMux()
	logger.Info().Msg(fmt.Sprintf("Starting serving files from %s", app.Path.WD+"/web-server-go-liftover/public"))
	mux.Handle("/index/", http.StripPrefix("/index/", http.FileServer(http.Dir(app.Path.WD+"/web-server-go-liftover/public"))))
	mux.HandleFunc("/index/hg19to38/", middleware.Conveyor(urlHandler.MainHandler19to38, urlHandler.Authorizer))
	mux.HandleFunc("/index/hg38to19/", middleware.Conveyor(urlHandler.MainHandler38to19, urlHandler.Authorizer))
	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", cfg.Constants.ServerHost, cfg.Constants.ServerPort),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	logger.Info().Msg(fmt.Sprintf("Running on %s", srv.Addr))
	return srv
}
