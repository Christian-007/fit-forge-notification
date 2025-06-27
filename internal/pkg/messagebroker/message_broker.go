package messagebroker

import (
	"log/slog"
	"sync"
)

type Message struct {
	ID     string // uuid
	Topic  string
	Points int
	UserID string
}

type InMemoryMessageBroker struct {
	logger      *slog.Logger
	subscribers map[string][]chan Message
	lock        sync.RWMutex
}

func NewInMemoryMessageBroker(logger *slog.Logger) *InMemoryMessageBroker {
	return &InMemoryMessageBroker{
		subscribers: make(map[string][]chan Message),
		logger:      logger,
	}
}

func (i *InMemoryMessageBroker) Subscribe(key string) chan Message {
	i.lock.RLock()
	defer i.lock.RUnlock()

	ch := make(chan Message, 10)
	i.subscribers[key] = append(i.subscribers[key], ch)
	return ch
}

func (i *InMemoryMessageBroker) Publish(key string, message Message) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	subscribers := i.subscribers[key]
	for _, ch := range subscribers {
		select {
		case ch <- message:
			i.logger.Info("[InMemoryMessageBroker] published a message", slog.String("key", key), slog.Any("payload", message))
		default:
			i.logger.Info("[InMemoryMessageBroker] unsuccessful message to publish", slog.String("key", key), slog.Any("payload", message))
		}
	}
}
