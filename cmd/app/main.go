package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/StewardMcCormick/go-job-queue/cmd/app"
)

func main() {
	a, err := app.Init()
	if err != nil {
		panic(err)
	}

	err = a.Run()
	if err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig
	err = a.Shutdown()
	if err != nil {
		panic(err)
	}
}
