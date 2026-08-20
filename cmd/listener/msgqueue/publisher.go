package msgqueue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type MessageWithID interface {
	GetID() string
}

type JetStreamPublisher interface {
	PublishMsg(ctx context.Context, msg *nats.Msg, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error)
}

type Publisher[T MessageWithID] struct {
	subject string
	jsPub   JetStreamPublisher
}

func newPublisher[T MessageWithID](subject string, jsPub JetStreamPublisher) *Publisher[T] {
	return &Publisher[T]{
		subject: subject,
		jsPub:   jsPub,
	}
}

func (q *Publisher[T]) Push(ctx context.Context, msg T) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	natsMsg := &nats.Msg{
		Subject: q.subject,
		Header:  nats.Header{},
		Data:    data,
	}

	// Set Msg-Id header to avoid message duplication
	withMsgIDOpt := jetstream.WithMsgID(q.getMsgIDWithSubject(msg.GetID()))
	_, err = q.jsPub.PublishMsg(ctx, natsMsg, withMsgIDOpt)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

// getMsgIDWithSubject scopes the dedup ID per subject to avoid collisions
func (q *Publisher[T]) getMsgIDWithSubject(id string) string {
	return q.subject + "." + id
}
