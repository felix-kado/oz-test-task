// internal/pubsub/inmemory.go
package pubsub

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type InMemoryPubSub struct {
	subscribers map[uuid.UUID]map[chan string]struct{}
	mu          sync.RWMutex
}

// NewInMemoryPubSub creates a new instance of InMemoryPubSub.
func NewInMemoryPubSub() *InMemoryPubSub {
	return &InMemoryPubSub{
		subscribers: make(map[uuid.UUID]map[chan string]struct{}),
	}
}

// Subscribe allows a client to subscribe to a specific postID.
// Returns a channel to receive messages for the given postID.
func (ps *InMemoryPubSub) Subscribe(ctx context.Context, postID uuid.UUID) (<-chan string, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	slog.Info("Subscribing to post", "postID", postID)

	ch := make(chan string, 1)
	if _, ok := ps.subscribers[postID]; !ok {
		ps.subscribers[postID] = make(map[chan string]struct{})
	}
	ps.subscribers[postID][ch] = struct{}{}

	// Goroutine to handle cleanup when context is done.
	go func() {
		<-ctx.Done()
		ps.mu.Lock()
		delete(ps.subscribers[postID], ch)
		ps.mu.Unlock()
		close(ch)
		slog.Info("Unsubscribed from post", "postID", postID)
	}()

	return ch, nil
}

// Publish sends a message to all subscribers of the given postID.
func (ps *InMemoryPubSub) Publish(ctx context.Context, postID uuid.UUID, message string) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	slog.Info("Publishing message to post", "postID", postID, "message", message)

	if subscribers, ok := ps.subscribers[postID]; ok {
		for ch := range subscribers {
			select {
			case ch <- message:
			case <-ctx.Done():
				slog.Warn("Publishing interrupted by context done", "postID", postID)
				return ctx.Err()
			}
		}
		slog.Info("Message published to all subscribers", "postID", postID)
	} else {
		slog.Warn("No subscribers for post", "postID", postID)
	}
	return nil
}
