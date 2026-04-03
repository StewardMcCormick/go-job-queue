package main

import (
	"log"

	"github.com/StewardMcCormick/go-job-queue/config"
	"github.com/StewardMcCormick/go-job-queue/internal/api/server"
)

type Server interface {
	Run() error
	Stop() error
	Addr() string
}

type App struct {
	server Server
}

func InitApp(cfg config.Config) (*App, error) {
	a := &App{}

	err := a.InitServer(cfg.Server)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) InitServer(cfg server.Config) error {
	s, err := server.NewServer(cfg)
	if err != nil {
		return err
	}

	a.server = s
	return nil
}

func (a *App) Run() {
	go func() {
		log.Printf("[START] Server starts on: %s", a.server.Addr())
		err := a.server.Run()
		if err != nil {
			panic(err)
		}
	}()
}

func (a *App) Shutdown() error {
	return nil
}
