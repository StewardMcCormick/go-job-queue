package app

import (
	"context"
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/StewardMcCormick/go-job-queue/config"
	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	"github.com/StewardMcCormick/go-job-queue/internal/adapter/postgres"
	appredis "github.com/StewardMcCormick/go-job-queue/internal/adapter/redis"
	"github.com/StewardMcCormick/go-job-queue/internal/api/handlers"
	"github.com/StewardMcCormick/go-job-queue/internal/api/service"
	uc "github.com/StewardMcCormick/go-job-queue/internal/api/use_case"
	"github.com/StewardMcCormick/go-job-queue/internal/storage"
	bus "github.com/StewardMcCormick/go-job-queue/pkg/event_bus"
	"github.com/StewardMcCormick/go-job-queue/pkg/log"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type DIContainer interface {
	Db() postgres.DB
	RedisClient(db int) *redis.Client
	EventBus() bus.EventBus
	Logger() *zap.Logger
	TaskRedisStorage() storage.TaskRedisStorage
	TaskPostgresStorage() storage.TaskPostgresStorage
	TaskService() service.TaskService
	TaskUseCase() uc.TaskUseCase
	Handlers() pb.JobQueueServiceServer

	Close(ctx context.Context) error
}

type diContainer struct {
	closer *closer

	db          postgres.DB
	redisClient *redis.Client
	eventBus    bus.EventBus
	logger      *zap.Logger

	taskRedisStorage    storage.TaskRedisStorage
	taskPostgresStorage storage.TaskPostgresStorage

	taskService service.TaskService

	taskUseCase uc.TaskUseCase

	handlers pb.JobQueueServiceServer
}

func NewDIContainer() *diContainer {
	di := &diContainer{}

	di.closer = NewCloser()
	di.closer.log = di.Logger()

	return di
}

func (d *diContainer) Db() postgres.DB {
	if d.db == nil {
		start := time.Now()
		d.logger.Info("[START] DB connection initialization...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		pg, err := postgres.NewPool(ctx, config.AppConfig().Postgres)
		if err != nil {
			d.logger.Fatal(fmt.Sprintf("[START] DB connection initialization error: %v, duration: %d ms",
				err, time.Since(start).Milliseconds()),
			)
		}

		d.db = pg
		d.logger.Info(fmt.Sprintf("[START] DB connection initialization completed, duration: %d ms",
			time.Since(start).Milliseconds()),
		)

		d.closer.Add("DB", func(ctx context.Context) error {
			d.db.Close()
			return nil
		})
	}

	return d.db
}

func (d *diContainer) RedisClient(db int) *redis.Client {
	if d.redisClient == nil {
		start := time.Now()
		d.logger.Info("[START] Redis connection initialization...")
		rc, err := appredis.NewConnection(config.AppConfig().Redis, db)
		if err != nil {
			d.logger.Fatal(fmt.Sprintf("[START] Redis connection initialization error: %v, duration: %d ms",
				err, time.Since(start).Milliseconds()),
			)
		}

		d.redisClient = rc
		d.logger.Info(fmt.Sprintf("[START] Redis connection initialization completed, duration: %d ms",
			time.Since(start).Milliseconds()),
		)

		d.closer.Add("Redis", func(ctx context.Context) error {
			return d.redisClient.Close()
		})
	}

	return d.redisClient
}

func (d *diContainer) EventBus() bus.EventBus {
	if d.eventBus == nil {
		start := time.Now()
		d.logger.Info("[START] Event Bus initialization...")
		d.eventBus = bus.NewEventBus()
		d.logger.Info(fmt.Sprintf("[START] Event Bus initialization completed, duration: %d ms",
			time.Since(start).Milliseconds()),
		)
	}

	return d.eventBus
}

func (d *diContainer) Logger() *zap.Logger {
	if d.logger == nil {
		cfg := config.AppConfig()
		l, err := log.NewLogger(cfg.Log, string(cfg.App.Env), cfg.App.Name, cfg.App.Version)
		if err != nil {
			panic(err)
		}

		d.logger = l

		d.closer.Add("Logger", func(ctx context.Context) error {
			if err := l.Sync(); !errors.Is(err, syscall.ENOTTY) && !errors.Is(err, syscall.EINVAL) && !errors.Is(err, syscall.EBADF) {
				return err
			}

			return nil
		})
	}

	return d.logger
}

func (d *diContainer) TaskRedisStorage() storage.TaskRedisStorage {
	if d.taskRedisStorage == nil {
		d.taskRedisStorage = storage.NewTaskRedisStorage(d.RedisClient(0))
	}

	return d.taskRedisStorage
}

func (d *diContainer) TaskPostgresStorage() storage.TaskPostgresStorage {
	if d.taskPostgresStorage == nil {
		d.taskPostgresStorage = storage.NewTaskPostgresStorage(d.Db())
	}

	return d.taskPostgresStorage
}

func (d *diContainer) TaskService() service.TaskService {
	if d.taskService == nil {
		d.taskService = service.NewTaskService(d.EventBus(), d.TaskRedisStorage(), d.TaskPostgresStorage())
	}

	return d.taskService
}

func (d *diContainer) TaskUseCase() uc.TaskUseCase {
	if d.taskUseCase == nil {
		d.taskUseCase = uc.NewTaskUseCase(d.TaskService())
	}

	return d.taskUseCase
}

func (d *diContainer) Handlers() pb.JobQueueServiceServer {
	if d.handlers == nil {
		d.handlers = handlers.NewHandler(d.TaskUseCase())
	}

	return d.handlers
}

func (d *diContainer) Close(ctx context.Context) error {
	return d.closer.Close(ctx)
}
