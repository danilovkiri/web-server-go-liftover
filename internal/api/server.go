package api

import (
	"fmt"
	"net/http"
	"time"
	"webServerGoLiftover/internal/api/handlers"
	"webServerGoLiftover/internal/api/middleware"
	"webServerGoLiftover/internal/config"
)

func InitServer(cfg *config.ServerConfig, app *config.Application) *http.Server {
	urlHandler := handlers.InitURLHandler(cfg, app)
	mux := http.NewServeMux()
	mux.Handle("/index/", http.StripPrefix("/index/", http.FileServer(http.Dir("../../public"))))
	mux.HandleFunc("/index/hg19to38/", middleware.Conveyor(urlHandler.MainHandler19to38, urlHandler.Authorizer))
	mux.HandleFunc("/index/hg38to19/", middleware.Conveyor(urlHandler.MainHandler38to19, urlHandler.Authorizer))
	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", cfg.Constants.ServerIP, cfg.Constants.ServerPort),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	return srv
}
