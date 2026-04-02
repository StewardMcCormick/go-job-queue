package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, err := Init()
	if err != nil {
		panic(err)
	}

	err = app.Run()
	if err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig
	err = app.Shutdown()
	if err != nil {
		panic(err)
	}
}
