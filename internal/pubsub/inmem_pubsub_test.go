// internal/pubsub/inmemory_test.go
package pubsub

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestInMemoryPubSub_SubscribeAndPublish(t *testing.T) {
	ps := NewInMemoryPubSub()
	postID := uuid.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := ps.Subscribe(ctx, postID)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	message := "Hello, World!"
	go func() {
		time.Sleep(100 * time.Millisecond)
		if err := ps.Publish(context.Background(), postID, message); err != nil {
			t.Errorf("Publish() error = %v", err)
		}
	}()

	// Check if the subscriber receives the message
	select {
	case msg := <-ch:
		if msg != message {
			t.Errorf("Expected message %q, got %q", message, msg)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Did not receive message in time")
	}
}

func TestInMemoryPubSub_Unsubscribe(t *testing.T) {
	ps := NewInMemoryPubSub()
	postID := uuid.New()

	ctx, cancel := context.WithCancel(context.Background())

	ch, err := ps.Subscribe(ctx, postID)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	cancel()

	time.Sleep(100 * time.Millisecond)

	select {
	case _, ok := <-ch:
		if ok {
			t.Errorf("Expected channel to be closed, but it was open")
		}
	default:
		t.Errorf("Expected channel to be closed, but it was open")
	}

	message := "Hello, World!"
	if err := ps.Publish(context.Background(), postID, message); err != nil {
		t.Errorf("Publish() error = %v", err)
	}
}
