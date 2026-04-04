package events

import pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"

type EventType string

const (
	EventTypeCreateTask EventType = "task.create"
)

type Event struct {
	Type    EventType
	Payload *pb.Task
}

func NewCreateTaskEvent(payload *pb.Task) Event {
	return Event{
		Type:    EventTypeCreateTask,
		Payload: payload,
	}
}
