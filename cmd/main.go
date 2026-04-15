package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/StewardMcCormick/go-job-queue/cmd/app"
)

func main() {
	a := app.NewApp()

	a.Run()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	a.Shutdown()
}
