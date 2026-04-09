package main

import (
	"context"
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/StewardMcCormick/go-job-queue/config"
	"github.com/StewardMcCormick/go-job-queue/internal/adapter/postgres"
	r "github.com/StewardMcCormick/go-job-queue/internal/adapter/redis"
	"github.com/StewardMcCormick/go-job-queue/internal/api/handlers"
	"github.com/StewardMcCormick/go-job-queue/internal/api/server"
	"github.com/StewardMcCormick/go-job-queue/internal/api/service"
	uc "github.com/StewardMcCormick/go-job-queue/internal/api/use_case"
	"github.com/StewardMcCormick/go-job-queue/internal/storage"
	bus "github.com/StewardMcCormick/go-job-queue/pkg/event_bus"
	"github.com/StewardMcCormick/go-job-queue/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Server interface {
	Run() error
	Stop() error
	Addr() string
}

type App struct {
	server          Server
	log             *zap.Logger
	pgxPool         *pgxpool.Pool
	taskRedisClient *redis.Client
}

func InitApp(cfg config.Config) (*App, error) {
	a := &App{}

	a.InitLogger(cfg.Log, cfg.App.Env, cfg.App.Name, cfg.App.Version)
	err := a.InitRedis(cfg.Redis)
	if err != nil {
		return nil, err
	}

	err = a.InitPgxPool(cfg.Postgres)
	if err != nil {
		return nil, err
	}

	a.log.Info("[START] Server initialization...")
	err = a.InitServer(cfg.Server)
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

func (a *App) InitPgxPool(cfg postgres.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.log.Info("[START] PostgreSQL connection initialization...")
	pool, err := postgres.NewPool(ctx, cfg)
	if err != nil {
		return err
	}

	a.pgxPool = pool
	a.log.Info("[START] PostgreSQL connection initialization completed")
	return nil
}

func (a *App) InitRedis(cfg r.Config) error {
	a.log.Info("[START] Redis initialization...")
	taskClient, err := r.NewConnection(cfg, 0)
	if err != nil {
		return err
	}

	a.taskRedisClient = taskClient
	a.log.Info(fmt.Sprintf("[START] Redis starts on: %s:%s", cfg.Host, cfg.Port))
	return nil
}

func (a *App) InitServer(cfg server.Config) error {
	eventBus := bus.NewEventBus()
	taskRedisStorage := storage.NewTaskRedisStorage(a.taskRedisClient)
	taskPostgresStorage := storage.NewTaskPostgresStorage(a.pgxPool)
	taskService := service.NewTaskService(eventBus, taskRedisStorage, taskPostgresStorage)
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
		return fmt.Errorf("[SHUTDOWN] Server stop error: %w", err)
	}
	a.log.Info("[SHUTDOWN] Server closed")

	a.pgxPool.Close()
	a.log.Info("[SHUTDOWN] PostgreSQL connection closed")

	err = a.taskRedisClient.Close()
	if err != nil {
		return fmt.Errorf("[SHUTDOWN] Redis closing error: %w", err)
	}
	a.log.Info("[SHUTDOWN] Redis connection closed")

	err = a.log.Sync()
	if err != nil && !errors.Is(err, syscall.ENOTTY) && !errors.Is(err, syscall.EINVAL) && !errors.Is(err, syscall.EBADF) {
		return fmt.Errorf("[SHUTDOWN] Logger sync error: %w", err)
	}
	a.log.Info("[SHUTDOWN] Logger synced")

	return nil
}
