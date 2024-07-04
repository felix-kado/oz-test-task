// internal/pubsub/inmemory.go
package pubsub

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type InMemoryPubSub struct {
	subscribers map[uuid.UUID]map[chan string]struct{}
	mu          sync.RWMutex
}

func NewInMemoryPubSub() *InMemoryPubSub {
	return &InMemoryPubSub{
		subscribers: make(map[uuid.UUID]map[chan string]struct{}),
	}
}

func (ps *InMemoryPubSub) Subscribe(ctx context.Context, postID uuid.UUID) (<-chan string, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string, 1)
	if _, ok := ps.subscribers[postID]; !ok {
		ps.subscribers[postID] = make(map[chan string]struct{})
	}
	ps.subscribers[postID][ch] = struct{}{}

	go func() {
		<-ctx.Done()
		ps.mu.Lock()
		delete(ps.subscribers[postID], ch)
		ps.mu.Unlock()
		close(ch)
	}()

	return ch, nil
}

func (ps *InMemoryPubSub) Publish(ctx context.Context, postID uuid.UUID, message string) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if subscribers, ok := ps.subscribers[postID]; ok {
		for ch := range subscribers {
			ch <- message
		}
	}
	return nil
}
