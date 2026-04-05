package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/StewardMcCormick/go-job-queue/config"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Cannot load configuration: %v", err)
	}

	app, err := InitApp(cfg)
	if err != nil {
		log.Fatalf("Cannot init application: %v", err)
	}

	app.Run()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	app.log.Info("[SHUTDOWN] Start shutting down...")
	err = app.Shutdown()
	if err != nil {
		app.log.Error(fmt.Sprintf("[SHUTDOWN] Shoutdown error: %v", err))
	}
	app.log.Info("[SHUTDOWN] Shutdown completed")
}
