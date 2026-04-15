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

type TaskUseCase interface {
	Create(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error)
	GetById(ctx context.Context, id string) (*pb.GetTaskByIdResponse, error)
}

type JobHandler interface {
	Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error)
	CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error)
	GetTaskById(ctx context.Context, req *pb.GetTaskByIdRequest) (*pb.GetTaskByIdResponse, error)
}

type jobHandler struct {
	pb.UnimplementedJobQueueServiceServer
	taskUseCase TaskUseCase
}

func NewHandler(taskUseCase TaskUseCase) *jobHandler {
	return &jobHandler{
		taskUseCase: taskUseCase,
	}
}

func (h *jobHandler) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{RepeatedNum: req.Num}, nil
}

func (h *jobHandler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	if err := h.validateCreateTaskRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response, err := h.taskUseCase.Create(ctx, req)
	if err != nil {
		if errors.Is(err, errs.ErrBadRequest) {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("cannot create task - %s", err.Error()))
		}
		return nil, status.Errorf(codes.Internal, "cannot create task - %s", err.Error())
	}

	return response, nil
}

func (h *jobHandler) validateCreateTaskRequest(req *pb.CreateTaskRequest) error {
	if req.Priority == pb.TaskPriority_TASK_PRIORITY_UNSPECIFIED {
		return fmt.Errorf(
			`%w - task priority must be one of:
					- TASK_PRIORITY_BACKGROUND;
					- TASK_PRIORITY_NORMAL;
    				- TASK_PRIORITY_HIGH;
    				- TASK_PRIORITY_IMMEDIATE;`, errs.ErrValidation)
	}
	if req.Type == "" {
		return fmt.Errorf("%w - task type cannot be empty", errs.ErrValidation)
	}

	return nil
}

func (h *jobHandler) GetTaskById(ctx context.Context, req *pb.GetTaskByIdRequest) (*pb.GetTaskByIdResponse, error) {
	resp, err := h.taskUseCase.GetById(ctx, req.Id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "cannot get task with id %s: %s", req.Id, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "cannot get task with id %s: %s", req.Id, err.Error())
	}

	return resp, nil
}
