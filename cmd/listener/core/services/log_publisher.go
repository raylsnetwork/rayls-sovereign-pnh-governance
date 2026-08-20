package services

import (
	"context"
	"fmt"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/events"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
)

// LogPublisher publishes parsed contract logs to a message queue.
type LogPublisher struct {
	mq  core.ContractMQ
	log logger.Logger
}

// NewLogPublisher creates a new LogPublisher with the given message queue and logger.
func NewLogPublisher(mq core.ContractMQ, log logger.Logger) *LogPublisher {
	return &LogPublisher{
		mq:  mq,
		log: log,
	}
}

// Publish sends each contract log to the message queue sequentially.
// It returns an error on the first failed publish (fail-fast).
func (p *LogPublisher) Publish(ctx context.Context, logs []core.ContractLog) error {
	for i, l := range logs {
		if err := p.mq.Push(ctx, l); err != nil {
			return fmt.Errorf("failed to publish log %d/%d: %w", i+1, len(logs), err)
		}

		if l.ContractName != events.ContractProofs {
			p.log.Debug("Published contract log",
				"contract", l.ContractName,
				"event", l.EventName,
				"block", fmt.Sprint(l.BlockNumber))
		}
	}

	if len(logs) > 0 {
		p.log.Info("Published contract logs", "count", len(logs))
	}

	return nil
}
