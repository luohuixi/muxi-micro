package discovery

import (
	"context"
)

type DiscoverCenter interface {
	Discover(ctx context.Context, serviceName string) ([]string, error)
	Watch(ctx context.Context, serviceName string) <-chan *Event
}

type Event struct {
	Type    string
	Address string
}
