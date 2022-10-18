package main

import (
	"context"
	"github.com/prybintsev/validation_cloud/internal/db/sqlite"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/prybintsev/validation_cloud/internal/router"
)

func listenToSignals(cancel context.CancelFunc) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	log.Info("Gracefully shutting down the http server")
	cancel()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go listenToSignals(cancel)

	dbCon, err := sqlite.ConnectAndMigrate(ctx)
	if err != nil {
		log.WithError(err).Error("Authentication server has stopped unexpectedly")
		return
	}
	err = router.StartHttpServer(ctx, dbCon)
	if err != nil {
		log.WithError(err).Error("Authentication server has stopped unexpectedly")
		return
	}
	log.Info("Exiting")
}
