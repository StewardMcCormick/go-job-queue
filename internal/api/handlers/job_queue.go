package handlers

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	errs "github.com/StewardMcCormick/go-job-queue/internal/api/error"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrValidation = errors.New("validation error")

type TaskUseCase interface {
	Create(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error)
}

type JobHandler struct {
	pb.UnimplementedJobQueueServiceServer
	taskUseCase TaskUseCase
}

func NewHandler(taskUseCase TaskUseCase) *JobHandler {
	return &JobHandler{
		taskUseCase: taskUseCase,
	}
}

func (h *JobHandler) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{RepeatedNum: req.Num}, nil
}

func (h *JobHandler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	if err := h.validateCreateTaskRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response, err := h.taskUseCase.Create(ctx, req)
	if err != nil {
		if errors.Is(err, errs.ErrInternal) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("%s - cannot create task", err.Error()))
		}
		return nil, status.Error(codes.Unknown, fmt.Sprintf("%s - cannot create task", err.Error()))
	}

	return response, nil
}

func (h *JobHandler) validateCreateTaskRequest(req *pb.CreateTaskRequest) error {
	if req.Priority == pb.TaskPriority_TASK_PRIORITY_UNSPECIFIED {
		return fmt.Errorf(
			`%w - task priority must be one of:
					- TASK_PRIORITY_BACKGROUND;
					- TASK_PRIORITY_NORMAL;
    				- TASK_PRIORITY_HIGH;
    				- TASK_PRIORITY_IMMEDIATE;`, ErrValidation)
	}

	return nil
}
