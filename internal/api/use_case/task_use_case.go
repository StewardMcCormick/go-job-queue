package uc

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	"github.com/StewardMcCormick/go-job-queue/internal/api/domain/helpers"
	errs "github.com/StewardMcCormick/go-job-queue/internal/api/error"
)

var ErrInternal = errors.New("internal error")

type TaskService interface {
	PublishCreateEvent(ctx context.Context, req *pb.Task) error
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
	task := helpers.TaskCreateRequestToTask(req)
	err := uc.taskService.PublishCreateEvent(ctx, task)
	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("%w - event publishing error", errs.ErrInternal)
	}

	return helpers.TaskToCreateTaskResponse(task), nil
}
