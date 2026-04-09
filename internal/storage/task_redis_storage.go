package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	"github.com/redis/go-redis/v9"
)

type taskRedisStorage struct {
	client *redis.Client
}

func NewTaskRedisStorage(client *redis.Client) *taskRedisStorage {
	return &taskRedisStorage{
		client: client,
	}
}

func (s *taskRedisStorage) Save(ctx context.Context, task *pb.Task) error {
	dependsOn, err := json.Marshal(task.DependsOn)
	if err != nil {
		return fmt.Errorf("depends on marshalling error: %w", err)
	}
	dependencyFor, err := json.Marshal(task.DependencyFor)
	if err != nil {
		return fmt.Errorf("dependency for marshalling error: %w", err)
	}
	payload, err := json.Marshal(task.Payload)
	if err != nil {
		return fmt.Errorf("payload marshalling error: %w", err)
	}

	err = s.client.HSet(ctx, fmt.Sprintf("task:%s", task.Id),
		"Id", task.Id,
		"Status", task.Status.String(),
		"Priority", task.Priority.String(),
		"Type", task.Type,
		"Payload", payload,
		"ShouldRetryNumber", task.ShouldRetryNumber,
		"Retries", task.Retries,
		"Deadline", task.Deadline.AsTime(),
		"DependsOn", dependsOn,
		"DependencyFor", dependencyFor,
		"CreatedAt", task.CompletedAt.AsTime(),
		"UpdatedAt", task.UpdatedAt.AsTime(),
		"StartedAt", task.StartedAt.AsTime(),
		"CompletedAt", task.CompletedAt.AsTime(),
	).Err()
	if err != nil {
		return fmt.Errorf("saving task in redis error: %w", err)
	}

	return nil
}

func (s *taskRedisStorage) Exists(ctx context.Context, id string) (bool, error) {
	redisId := "task:" + id

	if err := s.client.HGet(ctx, redisId, "Id").Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("redis error: %w", err)
	}

	return true, nil
}

func (s *taskRedisStorage) Remove(ctx context.Context, id string) error {
	return s.client.Del(ctx, fmt.Sprintf("task:%s", id)).Err()
}

func (s *taskRedisStorage) UpdateDependencyFor(ctx context.Context, dependencyId, taskId string) error {
	key := fmt.Sprintf("task:%s", dependencyId)
	res, err := s.client.HGet(ctx, key, "DependencyFor").Result()
	if err != nil {
		return fmt.Errorf("cannot update dependencyFor field in redis: %w", err)
	}

	var data []string
	err = json.Unmarshal([]byte(res), &data)
	if err != nil {
		return fmt.Errorf("cannot update dependencyFor field in redis: %w", err)
	}

	data = append(data, taskId)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot update dependencyFor field in redis: %w", err)
	}

	err = s.client.HSet(ctx, key, "DependencyFor", jsonData).Err()
	if err != nil {
		return fmt.Errorf("cannot update dependencyFor field in redis: %w", err)
	}

	return nil
}
