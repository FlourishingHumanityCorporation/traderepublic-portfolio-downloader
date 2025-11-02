//go:generate go tool mockgen -source=eventbus.go -destination eventbus_mock_gen.go -package=bus

package bus

import (
	"log/slog"
	"sync"
)

type Event struct {
	Topic string
	ID    string
	Data  any
}

func NewEvent(topic, id string, data any) Event {
	return Event{
		Topic: topic,
		ID:    id,
		Data:  data,
	}
}

type EventHandler func(Event)

type EventBusInterface interface {
	Subscribe(string, EventHandler)
	Publish(Event)
}

type EventBus struct {
	subscribers map[string][]EventHandler
	mu          sync.RWMutex
}

func New() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

// Subscribe registers an event handler for a specific topic.
// The method is thread-safe and uses a write lock to prevent concurrent modifications.
// Multiple handlers can be registered for the same topic, and they will be called in the order they were added.
func (b *EventBus) Subscribe(topic string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[topic] = append(b.subscribers[topic], handler)
}

// Publish sends an event to all registered subscribers for the event's topic.
// The method is thread-safe and uses a read lock to allow concurrent publishing.
// Each subscriber handler is executed in a separate goroutine to prevent blocking.
// If no subscribers are found for the event topic, the event is silently ignored.
// The event publication is logged at debug level with the topic and event ID.
func (b *EventBus) Publish(event Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if handlers, found := b.subscribers[event.Topic]; found {
		for _, handler := range handlers {
			go handler(event)
		}
	}

	slog.Debug("event published", "topic", event.Topic, "id", event.ID)
}

// Retry republishes the given event to the event bus.
func (b *EventBus) Retry(event Event) {
	b.Publish(event)
}
