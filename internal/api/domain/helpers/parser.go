package helpers

import pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"

func TaskCreateRequestToTask(req *pb.CreateTaskRequest) *pb.Task {
	return &pb.Task{
		Priority:    req.Priority,
		Type:        req.Type,
		Payload:     req.Payload,
		RetryNumber: req.RetryNumber,
		Deadline:    req.Deadline,
		DependsOn:   req.DependsOn,
	}
}

func TaskToCreateTaskResponse(task *pb.Task) *pb.CreateTaskResponse {
	return &pb.CreateTaskResponse{
		Task: task,
	}
}
