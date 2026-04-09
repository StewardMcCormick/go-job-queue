package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	errs "github.com/StewardMcCormick/go-job-queue/internal/api/error"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *taskRedisStorage) GetById(ctx context.Context, id string) (*pb.Task, error) {
	id = fmt.Sprintf("task:%s", id)

	res, err := s.client.HGetAll(ctx, id).Result()
	if len(res) == 0 {
		return nil, fmt.Errorf("%w - task with id '%s' was not found", errs.ErrNotFound, id)
	}
	if err != nil {
		return nil, fmt.Errorf("get task by id error: %w", err)
	}

	task, err := s.parseTaskFromMap(res)
	if err != nil {
		return nil, fmt.Errorf("cannot get task with id '%s': %w", id, err)
	}

	return task, nil
}

func (s *taskRedisStorage) parseTaskFromMap(m map[string]string) (*pb.Task, error) {
	createdAt, err := time.Parse(time.RFC3339, m["CreatedAt"])
	if err != nil {
		return nil, fmt.Errorf("time parsing error: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, m["UpdatedAt"])
	if err != nil {
		return nil, fmt.Errorf("time parsing error: %w", err)
	}

	startedAt, err := time.Parse(time.RFC3339, m["StartedAt"])
	if err != nil {
		return nil, fmt.Errorf("time parsing error: %w", err)
	}

	completedAt, err := time.Parse(time.RFC3339, m["CompletedAt"])
	if err != nil {
		return nil, fmt.Errorf("time parsing error: %w", err)
	}

	deadline, err := time.Parse(time.RFC3339, m["Deadline"])
	if err != nil {
		return nil, fmt.Errorf("time parsing error: %w", err)
	}

	shouldRetryNumber, err := strconv.ParseUint(m["ShouldRetryNumber"], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("uint parse error: %w", err)
	}

	retries, err := strconv.ParseUint(m["Retries"], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("uint parse error: %w", err)
	}

	var dependsOn []string
	err = json.Unmarshal([]byte(m["DependsOn"]), &dependsOn)
	if err != nil {
		return nil, fmt.Errorf("dependsOn parse error: %w", err)
	}

	var dependencyFor []string
	err = json.Unmarshal([]byte(m["DependencyFor"]), &dependencyFor)
	if err != nil {
		return nil, fmt.Errorf("dependencyFor parse error: %w", err)
	}

	return &pb.Task{
		Id:                m["Id"],
		Status:            pb.TaskStatus(pb.TaskStatus_value[m["Status"]]),
		Priority:          pb.TaskPriority(pb.TaskPriority_value[m["Priority"]]),
		Type:              m["type"],
		Payload:           []byte(m["payload"]),
		ShouldRetryNumber: uint32(shouldRetryNumber),
		Retries:           uint32(retries),
		Deadline:          timestamppb.New(deadline),
		DependsOn:         dependsOn,
		DependencyFor:     dependencyFor,
		CreatedAt:         timestamppb.New(createdAt),
		UpdatedAt:         timestamppb.New(updatedAt),
		StartedAt:         timestamppb.New(startedAt),
		CompletedAt:       timestamppb.New(completedAt),
	}, nil
}
