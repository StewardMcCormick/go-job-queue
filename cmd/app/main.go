package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := NewApp()

	app.Run()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	app.Shutdown()
}
