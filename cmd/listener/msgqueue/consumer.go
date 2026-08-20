package msgqueue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type JetStreamConsumer interface {
	Messages(opts ...jetstream.PullMessagesOpt) (jetstream.MessagesContext, error)
}

type Message[T MessageWithID] struct {
	V   T
	Ack func(context.Context) error
}

type Consumer[T MessageWithID] struct {
	iter jetstream.MessagesContext
}

func newConsumer[T MessageWithID](jsCons JetStreamConsumer) (*Consumer[T], error) {
	iter, err := jsCons.Messages()
	if err != nil {
		return nil, fmt.Errorf("failed to create messages iterator: %w", err)
	}
	return &Consumer[T]{
		iter: iter,
	}, nil
}

func (c *Consumer[T]) Next(ctx context.Context) (Message[T], error) {
	for {
		jsMsg, err := c.iter.Next(
			jetstream.NextContext(ctx),
		)
		if err != nil {
			// If it's the internal MaxWait timeout, just loop again.
			if errors.Is(err, context.DeadlineExceeded) {
				// no message during this poll; try again
				continue
			}

			// If ctx is done, surface that.
			if ctx.Err() != nil {
				return Message[T]{}, ctx.Err()
			}
			return Message[T]{}, fmt.Errorf("next message failed: %w", err)
		}

		var empty Message[T]

		obj, err := unmarshalInto[T](jsMsg.Data())
		if err != nil {
			return empty, fmt.Errorf("failed to unmarshall message: %w", err)
		}

		return Message[T]{
			V:   obj,
			Ack: jsMsg.DoubleAck,
		}, nil
	}
}

func unmarshalInto[T any](data []byte) (obj T, err error) {
	err = json.Unmarshal(data, &obj)
	return obj, err
}
