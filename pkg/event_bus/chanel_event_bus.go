package bus

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/StewardMcCormick/go-job-queue/pkg/event_bus/events"
)

var (
	ErrNoSubscribers = errors.New("no one subscribers")
	ErrExitByTimeout = errors.New("exit by timeout")
)

type eventBus struct {
	channels map[events.EventType]chan events.Event
	timeout  time.Duration
	mu       *sync.Mutex
}

func NewEventBus() *eventBus {
	return &eventBus{
		channels: make(map[events.EventType]chan events.Event),
		timeout:  5 * time.Second,
		mu:       &sync.Mutex{},
	}
}

func (b *eventBus) Publish(ctx context.Context, event events.Event) error {
	// b.mu.Lock()
	// ch, exist := b.channels[event.Type]
	// b.mu.Unlock()
	//
	// if !exist {
	//	return fmt.Errorf("%w - no subscibers for event type %v", ErrNoSubscribers, event.Type)
	// }
	//
	// if err := b.publishBlocked(ctx, event, ch); err != nil {
	//	return fmt.Errorf("%w - event publishing error", err)
	// }
	//
	// log.Printf("new event: %v", event)

	return nil
}

// func (b *eventBus) publishBlocked(ctx context.Context, event events.Event, ch chan<- events.Event) error {
//	select {
//	case ch <- event:
//		return nil
//	case <-time.After(b.timeout):
//		return ErrExitByTimeout
//	case <-ctx.Done():
//		return ctx.Err()
//	}
// }
