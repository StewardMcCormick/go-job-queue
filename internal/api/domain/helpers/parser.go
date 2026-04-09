package helpers

import pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"

func TaskCreateRequestToTask(req *pb.CreateTaskRequest) *pb.Task {
	return &pb.Task{
		Priority:          req.Priority,
		Type:              req.Type,
		Payload:           req.Payload,
		ShouldRetryNumber: req.ShouldRetryNumber,
		Deadline:          req.Deadline,
		DependsOn:         req.DependsOn,
	}
}

func TaskToCreateTaskResponse(task *pb.Task) *pb.CreateTaskResponse {
	return &pb.CreateTaskResponse{
		Id:                task.Id,
		Status:            task.Status,
		Priority:          task.Priority,
		Type:              task.Type,
		Payload:           task.Payload,
		ShouldRetryNumber: task.ShouldRetryNumber,
		Retries:           task.Retries,
		Deadline:          task.Deadline,
		DependsOn:         task.DependsOn,
		DependencyFor:     task.DependencyFor,
		CreatedAt:         task.CreatedAt,
		UpdatedAt:         task.UpdatedAt,
		StartedAt:         task.StartedAt,
		CompletedAt:       task.CompletedAt,
	}
}

func TaskToGetTaskByIdResponse(task *pb.Task) *pb.GetTaskByIdResponse {
	return &pb.GetTaskByIdResponse{
		Id:                task.Id,
		Status:            task.Status,
		Priority:          task.Priority,
		Type:              task.Type,
		Payload:           task.Payload,
		ShouldRetryNumber: task.ShouldRetryNumber,
		Retries:           task.Retries,
		Deadline:          task.Deadline,
		DependsOn:         task.DependsOn,
		DependencyFor:     task.DependencyFor,
		CreatedAt:         task.CreatedAt,
		UpdatedAt:         task.UpdatedAt,
		StartedAt:         task.StartedAt,
		CompletedAt:       task.CompletedAt,
	}
}
