package core

import (
	"context"
	"time"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
)

// HeaderProofPurgeProcessor periodically deletes header proof events older than the configured retention period,
// preserving the most recent proof per chain for liveness checking.
type HeaderProofPurgeProcessor struct {
	repo          HeaderProofRepository
	retentionTime time.Duration
	log           logger.Logger
}

// NewHeaderProofPurgeProcessor creates a new HeaderProofPurgeProcessor.
func NewHeaderProofPurgeProcessor(
	repo HeaderProofRepository,
	retentionTime time.Duration,
	log logger.Logger,
) *HeaderProofPurgeProcessor {
	return &HeaderProofPurgeProcessor{
		repo:          repo,
		retentionTime: retentionTime,
		log:           log,
	}
}

// Run starts the purge loop. On each tick it deletes records older than retentionTime,
// always preserving the latest proof per chain. Returns nil when ctx is cancelled.
func (p *HeaderProofPurgeProcessor) Run(ctx context.Context, ticker <-chan time.Time) error {
	p.log.Info("Header proof purge processor started", "retention", p.retentionTime)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker:
			cutoff := time.Now().Add(-p.retentionTime)
			deleted, err := p.repo.DeleteOlderThan(ctx, cutoff)
			if err != nil {
				p.log.Error("Failed to purge header proof events", "error", err)
				continue
			}
			if deleted > 0 {
				p.log.Info("Purged header proof events", "deleted", deleted, "cutoff", cutoff)
			}
		}
	}
}
