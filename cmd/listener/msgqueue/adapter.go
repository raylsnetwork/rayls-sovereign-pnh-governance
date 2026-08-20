package msgqueue

import (
	"context"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
)

// Compile-time interface checks
var (
	_ core.ContractMQ    = (*PublisherAdapter)(nil)
	_ core.ContractMQ    = (*RoutingPublisherAdapter)(nil)
	_ core.EventConsumer = (*ConsumerAdapter)(nil)
)

// contractLogPublisher is the subset of Publisher[ContractLog] the adapter needs.
type contractLogPublisher interface {
	Push(ctx context.Context, msg core.ContractLog) error
}

// contractLogConsumer is the subset of Consumer[ContractLog] the adapter needs.
type contractLogConsumer interface {
	Next(ctx context.Context) (Message[core.ContractLog], error)
}

// PublisherAdapter wraps a Publisher[ContractLog] to satisfy core.ContractMQ.
type PublisherAdapter struct {
	pub contractLogPublisher
}

// NewPublisherAdapter creates a new PublisherAdapter.
func NewPublisherAdapter(pub *Publisher[core.ContractLog]) *PublisherAdapter {
	return &PublisherAdapter{pub: pub}
}

// Push delegates to the underlying publisher.
func (a *PublisherAdapter) Push(ctx context.Context, log core.ContractLog) error {
	return a.pub.Push(ctx, log)
}

// RoutingPublisherAdapter routes logs to different publishers based on a predicate.
// Logs for which toSecondary returns true go to the secondary publisher; all others to primary.
type RoutingPublisherAdapter struct {
	primary     contractLogPublisher
	secondary   contractLogPublisher
	toSecondary func(core.ContractLog) bool
}

// NewRoutingPublisherAdapter creates a RoutingPublisherAdapter.
func NewRoutingPublisherAdapter(
	primary *Publisher[core.ContractLog],
	secondary *Publisher[core.ContractLog],
	toSecondary func(core.ContractLog) bool,
) *RoutingPublisherAdapter {
	return &RoutingPublisherAdapter{
		primary:     primary,
		secondary:   secondary,
		toSecondary: toSecondary,
	}
}

// Push routes the log to the appropriate publisher.
func (a *RoutingPublisherAdapter) Push(ctx context.Context, log core.ContractLog) error {
	if a.toSecondary(log) {
		return a.secondary.Push(ctx, log)
	}
	return a.primary.Push(ctx, log)
}

// ConsumerAdapter wraps a Consumer[ContractLog] to satisfy core.EventConsumer.
type ConsumerAdapter struct {
	cons contractLogConsumer
}

// NewConsumerAdapter creates a new ConsumerAdapter.
func NewConsumerAdapter(cons *Consumer[core.ContractLog]) *ConsumerAdapter {
	return &ConsumerAdapter{cons: cons}
}

// Next consumes the next message and converts it to a core.Message.
func (a *ConsumerAdapter) Next(ctx context.Context) (core.Message, error) {
	msg, err := a.cons.Next(ctx)
	if err != nil {
		return core.Message{}, err
	}
	return core.Message{
		Log: msg.V,
		Ack: msg.Ack,
	}, nil
}
