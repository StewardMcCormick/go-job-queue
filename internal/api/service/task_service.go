package service

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	"github.com/StewardMcCormick/go-job-queue/pkg/event_bus/events"
	"github.com/jackc/pgx/v5"
)

type EventBus interface {
	Publish(ctx context.Context, event events.Event) error
}

type RedisStorage interface {
	Save(ctx context.Context, task *pb.Task) error
	Exists(ctx context.Context, id string) (bool, error)
	Remove(ctx context.Context, id string) error
	UpdateDependencyFor(ctx context.Context, id, dependencyId string) error
}

type PostgresStorage interface {
	GetById(ctx context.Context, id string) ([]*pb.Task, error)
}

type taskService struct {
	eventBus EventBus
	redis    RedisStorage
	postgres PostgresStorage
}

func NewTaskService(eventBus EventBus, redis RedisStorage, postgres PostgresStorage) *taskService {
	return &taskService{
		eventBus: eventBus,
		redis:    redis,
		postgres: postgres,
	}
}

func (s *taskService) PublishCreateEvent(ctx context.Context, req *pb.Task) error {
	event := events.NewCreateTaskEvent(req)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return fmt.Errorf("publish event error: %w", err)
	}

	return nil
}

func (s *taskService) SaveInRedis(ctx context.Context, req *pb.Task) error {
	if err := s.redis.Save(ctx, req); err != nil {
		return fmt.Errorf("cannot save task: %w", err)
	}

	return nil
}

func (s *taskService) DeleteFromRedis(ctx context.Context, id string) error {
	err := s.redis.Remove(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot delete task with id %s: %w", id, err)
	}

	return nil
}

func (s *taskService) ValidateDependencies(ctx context.Context, req *pb.Task) error {
	for i, id := range req.DependsOn {
		exist, err := s.redis.Exists(ctx, id)
		if err != nil {
			return fmt.Errorf("cannot get task: %w", err)
		}

		if !exist {
			_, err = s.postgres.GetById(ctx, id)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return fmt.Errorf("%w - dependency with id %s does not exist", err, id)
				}
			}
			req.DependsOn[i] = ""
		}

		err = s.redis.UpdateDependencyFor(ctx, id, req.Id)
		if err != nil {
			return fmt.Errorf("cannot update dependencyFor for task with id %s: %w", id, err)
		}
	}

	return nil
}
