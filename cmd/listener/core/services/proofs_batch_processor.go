package services

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/events"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/contracts/Proofs"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
)

const (
	DefaultProofsBatchMaxSize  = 650
	DefaultProofsFlushInterval = 2 * time.Second
	proofsInternalChannelCap   = 1300
	shutdownFlushTimeout       = 5 * time.Second
)

// ProofsBatchProcessor consumes HeaderProofSubmitted events from a message queue,
// accumulates them in a buffer, and persists them in batch to reduce DB round trips.
// A burst of 606 events (6 chains × 101 PN blocks) is flushed in a single INSERT.
type ProofsBatchProcessor struct {
	consumer      core.EventConsumer
	repo          core.HeaderProofEventRepository
	log           logger.Logger
	maxBatch      int
	flushInterval time.Duration
}

// NewProofsBatchProcessor creates a new ProofsBatchProcessor.
func NewProofsBatchProcessor(
	consumer core.EventConsumer,
	repo core.HeaderProofEventRepository,
	log logger.Logger,
	maxBatch int,
	flushInterval time.Duration,
) *ProofsBatchProcessor {
	return &ProofsBatchProcessor{
		consumer:      consumer,
		repo:          repo,
		log:           log,
		maxBatch:      maxBatch,
		flushInterval: flushInterval,
	}
}

// Run starts the batch processor. It spawns a producer goroutine that feeds
// incoming messages into an internal channel, then drains that channel in a
// select loop, flushing on maxBatch size or flushInterval — whichever comes first.
// Returns nil when ctx is cancelled.
func (p *ProofsBatchProcessor) Run(ctx context.Context) error {
	ch := make(chan core.Message, proofsInternalChannelCap)

	// Producer: feeds consumer.Next() into ch without blocking the flush loop.
	go func() {
		for {
			msg, err := p.consumer.Next(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				p.log.Error("ProofsBatchProcessor: failed to get next message", "error", err)
				continue
			}
			select {
			case ch <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()

	ticker := time.NewTicker(p.flushInterval)
	defer ticker.Stop()

	buf := make([]core.Message, 0, p.maxBatch)

	for {
		select {
		case msg := <-ch:
			buf = append(buf, msg)
			if len(buf) >= p.maxBatch {
				p.flush(ctx, buf)
				buf = buf[:0]
			}

		case <-ticker.C:
			if len(buf) > 0 {
				p.flush(ctx, buf)
				buf = buf[:0]
			}

		case <-ctx.Done():
			if len(buf) > 0 {
				// Best-effort flush on shutdown; parent ctx is already cancelled so a
				// fresh context is required here.
				flushCtx, cancel := context.WithTimeout(
					context.Background(),
					shutdownFlushTimeout,
				) //nolint:contextcheck
				p.flush(flushCtx, buf)
				cancel()
			}
			return nil
		}
	}
}

// flush parses a batch of messages and persists them in a single DB call.
// On success all messages are acked; on failure none are acked (NATS redelivers).
// Messages that cannot be parsed are acked immediately to avoid infinite redelivery.
func (p *ProofsBatchProcessor) flush(ctx context.Context, messages []core.Message) {
	domainEvents := make([]*domain.HeaderProofEvent, 0, len(messages))
	toAck := make([]core.Message, 0, len(messages))
	toSkip := make([]core.Message, 0)

	for _, msg := range messages {
		if msg.Log.EventName != events.HeaderProofSubmitted {
			p.log.Debug("ProofsBatchProcessor: unexpected event, acking and skipping",
				"event", msg.Log.EventName)
			toSkip = append(toSkip, msg)
			continue
		}

		event, err := core.UnmarshalEventData[*Proofs.ProofsHeaderProofSubmitted](msg.Log)
		if err != nil {
			p.log.Error("ProofsBatchProcessor: failed to unmarshal event, acking to prevent redelivery loop",
				"error", err)
			toSkip = append(toSkip, msg)
			continue
		}

		domainEvents = append(domainEvents, &domain.HeaderProofEvent{
			ChainID:     domain.BigInt{Int: event.ChainId},
			BlockNumber: domain.BigInt{Int: event.BlockNumber},
			BlockHash:   common.BytesToHash(event.HeaderHash[:]).Hex(),
			CreatedAt:   msg.Log.BlockTimestamp,
		})
		toAck = append(toAck, msg)
	}

	// Ack unparseable messages immediately.
	for _, msg := range toSkip {
		if err := msg.Ack(ctx); err != nil {
			p.log.Error("ProofsBatchProcessor: failed to ack skipped message", "error", err)
		}
	}

	if len(domainEvents) == 0 {
		return
	}

	if err := p.repo.CreateBatch(ctx, domainEvents); err != nil {
		p.log.Error("ProofsBatchProcessor: failed to batch-insert proofs, messages will be redelivered",
			"count", len(domainEvents), "error", err)
		return
	}

	p.log.Info("ProofsBatchProcessor: flushed batch", "count", len(domainEvents))

	for _, msg := range toAck {
		if err := msg.Ack(ctx); err != nil {
			p.log.Error("ProofsBatchProcessor: failed to ack message after successful insert", "error", err)
		}
	}
}
