package main

import (
	"fmt"

	"github.com/StewardMcCormick/go-job-queue/config"
	"github.com/StewardMcCormick/go-job-queue/internal/api/handlers"
	"github.com/StewardMcCormick/go-job-queue/internal/api/server"
	"github.com/StewardMcCormick/go-job-queue/internal/api/service"
	uc "github.com/StewardMcCormick/go-job-queue/internal/api/use_case"
	bus "github.com/StewardMcCormick/go-job-queue/pkg/event_bus"
	"github.com/StewardMcCormick/go-job-queue/pkg/log"
	"go.uber.org/zap"
)

type Server interface {
	Run() error
	Stop() error
	Addr() string
}

type App struct {
	server Server
	log    *zap.Logger
}

func InitApp(cfg config.Config) (*App, error) {
	a := &App{}

	a.InitLogger(cfg.Log, cfg.App.Env, cfg.App.Name, cfg.App.Version)

	a.log.Info("[START] Server initialization...")
	err := a.InitServer(cfg.Server)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) InitLogger(cfg log.Config, env config.AppEnv, appName, appVersion string) {
	logger, err := log.NewLogger(cfg, string(env), appName, appVersion)
	if err != nil {
		panic(err)
	}

	a.log = logger
}

func (a *App) InitServer(cfg server.Config) error {
	eventBus := bus.NewEventBus()
	taskService := service.NewTaskService(eventBus)
	taskUseCase := uc.NewTaskUseCase(taskService)
	jobQueueHandler := handlers.NewHandler(taskUseCase)

	s, err := server.NewServer(cfg, a.log, jobQueueHandler)
	if err != nil {
		return err
	}

	a.server = s
	return nil
}

func (a *App) Run() {
	go func() {
		a.log.Info(fmt.Sprintf("[START] Server starts on: %s", a.server.Addr()))
		err := a.server.Run()
		if err != nil {
			a.log.Error(fmt.Sprintf("[START] Server start error: %v", err))
		}
	}()
}

func (a *App) Shutdown() error {
	err := a.server.Stop()
	if err != nil {
		return fmt.Errorf("server stop error: %w", err)
	}
	return nil
}
