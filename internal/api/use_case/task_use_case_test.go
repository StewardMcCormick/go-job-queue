package uc

import (
	"context"
	"errors"
	"testing"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	errs "github.com/StewardMcCormick/go-job-queue/internal/api/error"
	"github.com/StewardMcCormick/go-job-queue/internal/api/use_case/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TaskUseCaseTestSuite struct {
	suite.Suite
	mockService *mocks.MockTaskService
	useCase     *taskUseCase
}

func (s *TaskUseCaseTestSuite) TestNewTaskUseCase() {
	s.Run("creates new use case instance", func() {
		s.NotNil(s.useCase)
	})
}

func (s *TaskUseCaseTestSuite) SetupTest() {
	s.mockService = mocks.NewMockTaskService(s.T())
	s.useCase = NewTaskUseCase(s.mockService)
}

func TestTaskUseCaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUseCaseTestSuite))
}

func (s *TaskUseCaseTestSuite) TestCreate_Success() {
	testCases := []struct {
		name         string
		req          *pb.CreateTaskRequest
		expectedTask *pb.Task
	}{
		{
			name: "successful task creation with all fields",
			req: &pb.CreateTaskRequest{
				Priority:    pb.TaskPriority_TASK_PRIORITY_HIGH,
				Type:        "email",
				Payload:     map[string][]byte{"to": []byte("test@example.com")},
				RetryNumber: 3,
				DependsOn:   []string{"task-1"},
			},
			expectedTask: &pb.Task{
				Priority:    pb.TaskPriority_TASK_PRIORITY_HIGH,
				Type:        "email",
				Payload:     map[string][]byte{"to": []byte("test@example.com")},
				RetryNumber: 3,
				DependsOn:   []string{"task-1"},
			},
		},
		{
			name: "task creation with empty fields",
			req: &pb.CreateTaskRequest{
				Priority:    pb.TaskPriority_TASK_PRIORITY_UNSPECIFIED,
				Type:        "",
				Payload:     nil,
				RetryNumber: 0,
				DependsOn:   nil,
			},
			expectedTask: &pb.Task{
				Priority:    pb.TaskPriority_TASK_PRIORITY_UNSPECIFIED,
				Type:        "",
				Payload:     nil,
				RetryNumber: 0,
				DependsOn:   nil,
			},
		},
		{
			name: "task creation with empty payload map",
			req: &pb.CreateTaskRequest{
				Priority:    pb.TaskPriority_TASK_PRIORITY_NORMAL,
				Type:        "notification",
				Payload:     map[string][]byte{},
				RetryNumber: 1,
				DependsOn:   []string{"task-2", "task-3"},
			},
			expectedTask: &pb.Task{
				Priority:    pb.TaskPriority_TASK_PRIORITY_NORMAL,
				Type:        "notification",
				Payload:     map[string][]byte{},
				RetryNumber: 1,
				DependsOn:   []string{"task-2", "task-3"},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.mockService.EXPECT().PublishCreateEvent(mock.Anything, tc.expectedTask).Return(nil)

			resp, err := s.useCase.Create(context.Background(), tc.req)

			s.NoError(err)
			s.NotNil(resp)
			s.Equal(tc.expectedTask, resp.Task)
		})
	}
}

func (s *TaskUseCaseTestSuite) TestCreate_Error() {
	testCases := []struct {
		name         string
		req          *pb.CreateTaskRequest
		expectedTask *pb.Task
		mockError    error
	}{
		{
			name: "event publishing error",
			req: &pb.CreateTaskRequest{
				Priority:    pb.TaskPriority_TASK_PRIORITY_HIGH,
				Type:        "email",
				Payload:     map[string][]byte{"to": []byte("test@example.com")},
				RetryNumber: 3,
				DependsOn:   nil,
			},
			expectedTask: &pb.Task{
				Priority:    pb.TaskPriority_TASK_PRIORITY_HIGH,
				Type:        "email",
				Payload:     map[string][]byte{"to": []byte("test@example.com")},
				RetryNumber: 3,
				DependsOn:   nil,
			},
			mockError: errors.New("event bus error"),
		},
		{
			name: "context cancellation",
			req: &pb.CreateTaskRequest{
				Priority: pb.TaskPriority_TASK_PRIORITY_NORMAL,
				Type:     "task",
			},
			expectedTask: &pb.Task{
				Priority: pb.TaskPriority_TASK_PRIORITY_NORMAL,
				Type:     "task",
			},
			mockError: context.Canceled,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.mockService.EXPECT().PublishCreateEvent(mock.Anything, tc.expectedTask).Return(tc.mockError)

			resp, err := s.useCase.Create(context.Background(), tc.req)

			s.Error(err)
			s.Nil(resp)
			s.ErrorIs(err, errs.ErrInternal)
		})
	}
}
