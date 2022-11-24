package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"webServerGoLiftover/internal/api"
	"webServerGoLiftover/internal/config"
	"webServerGoLiftover/internal/logger"
	"webServerGoLiftover/internal/utils"
)

func main() {
	// initialize logger
	flog, err := os.OpenFile(`server.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer flog.Close()
	loggerInstance := logger.InitLog(flog)

	// set context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// define a waitGroup
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// parse configuration
	cfg, app, err := config.NewConfiguration()
	if err != nil {
		loggerInstance.Fatal().Err(err).Msg("Configuration parsing failed")
	}

	// make necessary directories
	dirs := []string{cfg.Constants.FileStorage, app.Path.UploadedDir, app.Path.ProcessedDir}
	for _, dir := range dirs {
		loggerInstance.Info().Msg(fmt.Sprintf("Attempting directory %s creation", dir))
		err = utils.MakeDir(dir)
		if err != nil {
			loggerInstance.Fatal().Err(err).Msg("Directory creation failed")
		}
	}

	// initialize server
	server := api.InitServer(cfg, app, loggerInstance)

	// set a graceful shutdown listener
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer wg.Done()
		<-done
		loggerInstance.Warn().Msg("Server shutdown attempted")
		ctxTO, cancelTO := context.WithTimeout(ctx, 30*time.Second)
		defer cancelTO()
		if err := server.Shutdown(ctxTO); err != nil {
			loggerInstance.Fatal().Err(err).Msg("Server shutdown failed")
		}
		cancel()
	}()

	// start up the server
	loggerInstance.Warn().Msg("Server start attempted")
	if err := server.ListenAndServeTLS(cfg.Constants.CertFile, cfg.Constants.KeyFile); err != nil && err != http.ErrServerClosed {
		loggerInstance.Fatal().Err(err).Msg("Server start failed")
	}
	wg.Wait()
	loggerInstance.Info().Msg("Server shutdown succeeded")
}
