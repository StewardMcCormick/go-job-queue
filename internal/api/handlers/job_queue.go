package handlers

import (
	"context"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
)

type JobHandler struct {
	pb.UnimplementedJobQueueServiceServer
}

func NewHandler() *JobHandler {
	return &JobHandler{}
}

func (h *JobHandler) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{RepeatedNum: req.Num}, nil
}
