package pubsub

import (
	"context"

	"github.com/google/uuid"
)

type PubSub interface {
	Subscribe(ctx context.Context, postID uuid.UUID) (<-chan string, error)
	Publish(ctx context.Context, postID uuid.UUID, message string) error
}
