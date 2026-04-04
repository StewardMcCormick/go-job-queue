package service

import (
	"context"
	"fmt"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	"github.com/StewardMcCormick/go-job-queue/pkg/event_bus/events"
)

type EventBus interface {
	Publish(ctx context.Context, event events.Event) error
}

type taskService struct {
	eventBus EventBus
}

func NewTaskService(eventBus EventBus) *taskService {
	return &taskService{
		eventBus: eventBus,
	}
}

func (s *taskService) PublishCreateEvent(ctx context.Context, req *pb.Task) error {
	event := events.NewCreateTaskEvent(req)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return fmt.Errorf("%w - publish event error", err)
	}

	return nil
}
