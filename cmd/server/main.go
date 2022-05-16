package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"webServerGoLiftover/internal/api"
	"webServerGoLiftover/internal/config"
	"webServerGoLiftover/internal/utils"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	cfg, app, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	err = utils.MakeDir(cfg.Constants.FileStorage)
	if err != nil {
		log.Fatal(err)
	}
	err = utils.MakeDir(app.Path.UploadedDir)
	if err != nil {
		log.Fatal(err)
	}
	err = utils.MakeDir(app.Path.ProcessedDir)
	if err != nil {
		log.Fatal(err)
	}
	server := api.InitServer(cfg, app)

	// set a listener for os.Signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer wg.Done()
		<-done
		log.Print("Server shutdown attempted")
		ctxTO, cancelTO := context.WithTimeout(ctx, 30*time.Second)
		defer cancelTO()
		if err := server.Shutdown(ctxTO); err != nil {
			log.Fatal("Server shutdown failed:", err)
		}
		log.Print("Server shutdown succeeded")
		cancel()
	}()

	// start up the server
	log.Println("Server start attempted")
	if err := server.ListenAndServeTLS(cfg.Constants.CertFile, cfg.Constants.KeyFile); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	wg.Wait()
}
