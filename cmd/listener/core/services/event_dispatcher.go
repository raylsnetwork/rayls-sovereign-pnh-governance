package services

import (
	"context"
	"fmt"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
)

// EventDispatcher consumes contract logs from a message queue and dispatches
// them to the appropriate EventHandler based on contract name.
type EventDispatcher struct {
	consumer core.EventConsumer
	handlers map[string]core.EventHandler
	log      logger.Logger
}

// NewEventDispatcher creates a new EventDispatcher. It builds a handler lookup map
// from the provided handlers and panics if two handlers register the same contract name.
func NewEventDispatcher(consumer core.EventConsumer, handlers []core.EventHandler, log logger.Logger) *EventDispatcher {
	handlerMap := make(map[string]core.EventHandler, len(handlers))
	for _, h := range handlers {
		name := h.ContractName()
		if _, exists := handlerMap[name]; exists {
			panic(fmt.Sprintf("duplicate handler for contract %q", name))
		}
		handlerMap[name] = h
	}

	return &EventDispatcher{
		consumer: consumer,
		handlers: handlerMap,
		log:      log,
	}
}

// Run is a blocking loop that consumes messages and dispatches them to handlers.
// It returns nil when the context is cancelled.
func (d *EventDispatcher) Run(ctx context.Context) error {
	for {
		msg, err := d.consumer.Next(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil //nolint:nilerr // context cancellation is expected shutdown; err is an artifact of the cancelled ctx
			}
			d.log.Error("Failed to get next message", "error", err)
			continue
		}

		handler, ok := d.handlers[msg.Log.ContractName]
		if !ok {
			d.log.Debug("No handler registered for contract, skipping",
				"contract", msg.Log.ContractName,
				"event", msg.Log.EventName,
				"block", fmt.Sprint(msg.Log.BlockNumber))
			if err := msg.Ack(ctx); err != nil {
				d.log.Error("Failed to ack skipped message", "error", err)
			}
			continue
		}

		if err := handler.Handle(ctx, msg.Log); err != nil {
			d.log.Error("Handler failed, message will be redelivered",
				"handler", handler.Name(),
				"contract", msg.Log.ContractName,
				"event", msg.Log.EventName,
				"error", err)
			continue
		}

		if err := msg.Ack(ctx); err != nil {
			d.log.Error("Failed to ack message", "error", err)
		}
	}
}
