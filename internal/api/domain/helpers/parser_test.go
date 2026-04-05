package helpers

import (
	"testing"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestTaskCreateRequestToTask(t *testing.T) {
	tests := []struct {
		name string
		req  *pb.CreateTaskRequest
		want *pb.Task
	}{
		{
			name: "valid request with all fields",
			req: &pb.CreateTaskRequest{
				Priority:          pb.TaskPriority_TASK_PRIORITY_HIGH,
				Type:              "email",
				Payload:           map[string][]byte{"to": []byte("test@example.com"), "subject": []byte("Hello")},
				ShouldRetryNumber: 3,
				Deadline:          nil,
				DependsOn:         []string{"task-1", "task-2"},
			},
			want: &pb.Task{
				Priority:          pb.TaskPriority_TASK_PRIORITY_HIGH,
				Type:              "email",
				Payload:           map[string][]byte{"to": []byte("test@example.com"), "subject": []byte("Hello")},
				ShouldRetryNumber: 3,
				Deadline:          nil,
				DependsOn:         []string{"task-1", "task-2"},
			},
		},
		{
			name: "request with empty fields",
			req: &pb.CreateTaskRequest{
				Priority:          pb.TaskPriority_TASK_PRIORITY_UNSPECIFIED,
				Type:              "",
				Payload:           nil,
				ShouldRetryNumber: 0,
				Deadline:          nil,
				DependsOn:         nil,
			},
			want: &pb.Task{
				Priority:          pb.TaskPriority_TASK_PRIORITY_UNSPECIFIED,
				Type:              "",
				Payload:           nil,
				ShouldRetryNumber: 0,
				Deadline:          nil,
				DependsOn:         nil,
			},
		},
		{
			name: "request with empty payload map",
			req: &pb.CreateTaskRequest{
				Priority:          pb.TaskPriority_TASK_PRIORITY_NORMAL,
				Type:              "notification",
				Payload:           map[string][]byte{},
				ShouldRetryNumber: 1,
				Deadline:          nil,
				DependsOn:         []string{"task-3"},
			},
			want: &pb.Task{
				Priority:          pb.TaskPriority_TASK_PRIORITY_NORMAL,
				Type:              "notification",
				Payload:           map[string][]byte{},
				ShouldRetryNumber: 1,
				Deadline:          nil,
				DependsOn:         []string{"task-3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TaskCreateRequestToTask(tt.req)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTaskCreateRequestToTask_NilInput(t *testing.T) {
	assert.Panics(t, func() {
		TaskCreateRequestToTask(nil)
	})
}

func TestTaskToCreateTaskResponse(t *testing.T) {
	tests := []struct {
		name string
		task *pb.Task
		want *pb.CreateTaskResponse
	}{
		{
			name: "valid task with all fields",
			task: &pb.Task{
				Id:                "task-123",
				Status:            pb.TaskStatus_TASK_STATUS_CREATED,
				Priority:          pb.TaskPriority_TASK_PRIORITY_HIGH,
				Type:              "notification",
				Payload:           map[string][]byte{"user_id": []byte("123")},
				ShouldRetryNumber: 5,
				Deadline:          nil,
				DependsOn:         []string{"task-3"},
			},
			want: &pb.CreateTaskResponse{
				Task: &pb.Task{
					Id:                "task-123",
					Status:            pb.TaskStatus_TASK_STATUS_CREATED,
					Priority:          pb.TaskPriority_TASK_PRIORITY_HIGH,
					Type:              "notification",
					Payload:           map[string][]byte{"user_id": []byte("123")},
					ShouldRetryNumber: 5,
					Deadline:          nil,
					DependsOn:         []string{"task-3"},
				},
			},
		},
		{
			name: "nil task",
			task: nil,
			want: &pb.CreateTaskResponse{
				Task: nil,
			},
		},
		{
			name: "task with empty payload",
			task: &pb.Task{
				Id:                "task-456",
				Status:            pb.TaskStatus_TASK_STATUS_PENDING,
				Priority:          pb.TaskPriority_TASK_PRIORITY_NORMAL,
				Type:              "email",
				Payload:           map[string][]byte{},
				ShouldRetryNumber: 0,
				Deadline:          nil,
				DependsOn:         nil,
			},
			want: &pb.CreateTaskResponse{
				Task: &pb.Task{
					Id:                "task-456",
					Status:            pb.TaskStatus_TASK_STATUS_PENDING,
					Priority:          pb.TaskPriority_TASK_PRIORITY_NORMAL,
					Type:              "email",
					Payload:           map[string][]byte{},
					ShouldRetryNumber: 0,
					Deadline:          nil,
					DependsOn:         nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TaskToCreateTaskResponse(tt.task)
			assert.Equal(t, tt.want, got)
		})
	}
}
