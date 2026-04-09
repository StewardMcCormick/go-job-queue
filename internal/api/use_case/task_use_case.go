package uc

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	"github.com/StewardMcCormick/go-job-queue/internal/api/domain/helpers"
	errs "github.com/StewardMcCormick/go-job-queue/internal/api/error"
	"github.com/StewardMcCormick/go-job-queue/pkg/app_context"
	bus "github.com/StewardMcCormick/go-job-queue/pkg/event_bus"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskService interface {
	PublishCreateEvent(ctx context.Context, req *pb.Task) error
	SaveInRedis(ctx context.Context, req *pb.Task) error
	DeleteFromRedis(ctx context.Context, id string) error

	ValidateDependencies(ctx context.Context, req *pb.Task) error
}

type taskUseCase struct {
	taskService TaskService
}

func NewTaskUseCase(taskService TaskService) *taskUseCase {
	return &taskUseCase{
		taskService: taskService,
	}
}

func (uc *taskUseCase) Create(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	log := appctx.GetLogger(ctx)

	task := helpers.TaskCreateRequestToTask(req)
	id, err := uc.generateId()
	if err != nil {
		return nil, fmt.Errorf("task id generation error: %w", err)
	}
	task.Id = id
	uc.setTime(task, time.Now())

	err = uc.taskService.ValidateDependencies(ctx, task)
	if err != nil {
		log.Error(fmt.Sprintf("dependencies validation error: %v", err))
		return nil, fmt.Errorf("dependencies validation error: %w", err)
	}

	if err := uc.taskService.SaveInRedis(ctx, task); err != nil {
		log.Error(fmt.Sprintf("task saving in redis error: %v", err))
		return nil, fmt.Errorf("%w - task saving in redis error", errs.ErrInternal)
	}

	err = uc.taskService.PublishCreateEvent(ctx, task)
	if err != nil {
		log.Error(err.Error())
		if err := uc.taskService.DeleteFromRedis(ctx, task.Id); err != nil {
			return nil, fmt.Errorf("cannot delete task after event publishing error: %w", err)
		}
		if errors.Is(err, bus.ErrNoSubscribers) {
			return nil, fmt.Errorf("%w - no available subscribers", errs.ErrBadRequest)
		}
		return nil, fmt.Errorf("%w - event publishing error", errs.ErrInternal)
	}

	return helpers.TaskToCreateTaskResponse(task), nil
}

func (uc *taskUseCase) generateId() (string, error) {
	s, err := gonanoid.New(32)
	if err != nil {
		return "", err
	}

	return s, nil
}

func (uc *taskUseCase) setTime(task *pb.Task, t time.Time) {
	now := timestamppb.New(t)
	task.CreatedAt, task.UpdatedAt = now, now
}
