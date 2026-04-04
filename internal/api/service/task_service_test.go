package service

import (
	"context"
	"errors"
	"testing"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	"github.com/StewardMcCormick/go-job-queue/internal/api/service/mocks"
	"github.com/StewardMcCormick/go-job-queue/pkg/event_bus/events"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TaskServiceTestSuite struct {
	suite.Suite
	mockEventBus *mocks.MockEventBus
	service      *taskService
}

func (s *TaskServiceTestSuite) TestNewTaskService() {
	s.Run("creates new service instance", func() {
		s.NotNil(s.service)
	})
}

func (s *TaskServiceTestSuite) SetupTest() {
	s.mockEventBus = mocks.NewMockEventBus(s.T())
	s.service = NewTaskService(s.mockEventBus)
}

func TestTaskServiceSuite(t *testing.T) {
	suite.Run(t, new(TaskServiceTestSuite))
}

func (s *TaskServiceTestSuite) TestPublishCreateEvent_Success() {
	testCases := []struct {
		name          string
		task          *pb.Task
		expectedEvent events.Event
	}{
		{
			name: "successful publish with all fields",
			task: &pb.Task{
				Id:          "task-123",
				Status:      pb.TaskStatus_TASK_STATUS_CREATED,
				Priority:    pb.TaskPriority_TASK_PRIORITY_HIGH,
				Type:        "email",
				Payload:     map[string][]byte{"to": []byte("test@example.com")},
				RetryNumber: 3,
				DependsOn:   []string{"task-1"},
			},
			expectedEvent: events.Event{
				Type: events.EventTypeCreateTask,
				Payload: &pb.Task{
					Id:          "task-123",
					Status:      pb.TaskStatus_TASK_STATUS_CREATED,
					Priority:    pb.TaskPriority_TASK_PRIORITY_HIGH,
					Type:        "email",
					Payload:     map[string][]byte{"to": []byte("test@example.com")},
					RetryNumber: 3,
					DependsOn:   []string{"task-1"},
				},
			},
		},
		{
			name: "successful publish with empty fields",
			task: &pb.Task{
				Priority:    pb.TaskPriority_TASK_PRIORITY_UNSPECIFIED,
				Type:        "",
				Payload:     nil,
				RetryNumber: 0,
				DependsOn:   nil,
			},
			expectedEvent: events.Event{
				Type: events.EventTypeCreateTask,
				Payload: &pb.Task{
					Priority:    pb.TaskPriority_TASK_PRIORITY_UNSPECIFIED,
					Type:        "",
					Payload:     nil,
					RetryNumber: 0,
					DependsOn:   nil,
				},
			},
		},
		{
			name: "successful publish with empty payload map",
			task: &pb.Task{
				Id:          "task-456",
				Status:      pb.TaskStatus_TASK_STATUS_PENDING,
				Priority:    pb.TaskPriority_TASK_PRIORITY_NORMAL,
				Type:        "notification",
				Payload:     map[string][]byte{},
				RetryNumber: 1,
				DependsOn:   []string{"task-2", "task-3"},
			},
			expectedEvent: events.Event{
				Type: events.EventTypeCreateTask,
				Payload: &pb.Task{
					Id:          "task-456",
					Status:      pb.TaskStatus_TASK_STATUS_PENDING,
					Priority:    pb.TaskPriority_TASK_PRIORITY_NORMAL,
					Type:        "notification",
					Payload:     map[string][]byte{},
					RetryNumber: 1,
					DependsOn:   []string{"task-2", "task-3"},
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.mockEventBus.EXPECT().Publish(mock.Anything, tc.expectedEvent).Return(nil)

			err := s.service.PublishCreateEvent(context.Background(), tc.task)

			s.NoError(err)
		})
	}
}

func (s *TaskServiceTestSuite) TestPublishCreateEvent_Error() {
	testCases := []struct {
		name          string
		task          *pb.Task
		expectedEvent events.Event
		mockError     error
	}{
		{
			name: "event bus publish error",
			task: &pb.Task{
				Id:       "task-789",
				Priority: pb.TaskPriority_TASK_PRIORITY_HIGH,
				Type:     "email",
				Payload:  map[string][]byte{"to": []byte("test@example.com")},
			},
			expectedEvent: events.Event{
				Type: events.EventTypeCreateTask,
				Payload: &pb.Task{
					Id:       "task-789",
					Priority: pb.TaskPriority_TASK_PRIORITY_HIGH,
					Type:     "email",
					Payload:  map[string][]byte{"to": []byte("test@example.com")},
				},
			},
			mockError: errors.New("event bus connection failed"),
		},
		{
			name: "context cancellation error",
			task: &pb.Task{
				Priority: pb.TaskPriority_TASK_PRIORITY_NORMAL,
				Type:     "task",
			},
			expectedEvent: events.Event{
				Type: events.EventTypeCreateTask,
				Payload: &pb.Task{
					Priority: pb.TaskPriority_TASK_PRIORITY_NORMAL,
					Type:     "task",
				},
			},
			mockError: context.Canceled,
		},
		{
			name: "nil task",
			task: nil,
			expectedEvent: events.Event{
				Type:    events.EventTypeCreateTask,
				Payload: nil,
			},
			mockError: errors.New("nil task error"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.mockEventBus.EXPECT().Publish(mock.Anything, tc.expectedEvent).Return(tc.mockError)

			err := s.service.PublishCreateEvent(context.Background(), tc.task)

			s.Error(err)
			s.Contains(err.Error(), "publish event error")
		})
	}
}
