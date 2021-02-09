package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func waitShutdown() {
	/* wait for shutdown signal */
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Warn().Msg("Shutting down...")
}
