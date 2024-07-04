package gql

import (
	"ozon-test/internal/models"
	"ozon-test/internal/pubsub"
)

type Resolver struct {
	Storage models.Storage
	PubSub  pubsub.PubSub
}
