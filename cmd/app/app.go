package app

import (
	"context"
	"time"

	"github.com/StewardMcCormick/go-job-queue/config"
	"github.com/StewardMcCormick/go-job-queue/internal/api/server"
	"go.uber.org/zap"
)

type App struct {
	deps   DIContainer
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
	go func() {
		if err := a.server.Run(); err != nil {
			panic(err)
		}
	}()
}

func (a *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.deps.Close(ctx); err != nil {
		panic(err)
	}
}
