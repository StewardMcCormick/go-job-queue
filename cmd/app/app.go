package main

import (
	"github.com/StewardMcCormick/go-job-queue/config"
	"github.com/StewardMcCormick/go-job-queue/internal/api/server"
	"go.uber.org/zap"
)

type App struct {
	deps   *diContainer
	server server.Server
	logger *zap.Logger
}

func NewApp() *App {
	a := &App{
		deps: NewDIContainer(),
	}

	a.initDeps()
	a.logger = a.deps.Logger()

	return a
}

func (a *App) initDeps() {
	inits := []func(){
		a.initServer,
	}

	for _, fn := range inits {
		fn()
	}
}

func (a *App) initServer() {
	s, err := server.NewServer(
		config.AppConfig().Server,
		a.deps.Logger(),
		a.deps.Handlers(),
	)
	if err != nil {
		panic(err)
	}

	a.server = s
}

func (a *App) Run() {
	if err := a.server.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Shutdown() {}
