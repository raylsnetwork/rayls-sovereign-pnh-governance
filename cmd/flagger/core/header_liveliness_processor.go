package core

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/withstack"
)

// HeaderProofRepository defines the interface for header proof operations
type HeaderProofRepository interface {
	GetLatestHeaderProofs(ctx context.Context) ([]domain.HeaderProofEvent, error)
	DeleteOlderThan(ctx context.Context, cutoff time.Time) (int64, error)
}

// HeaderFlagEventRepository defines the interface for header flag event operations
type HeaderFlagEventRepository interface {
	FlagParticipant(
		ctx context.Context,
		chainID *big.Int,
		blockNumber *big.Int,
		reason uint8,
		initiator uint8,
	) (flagged bool, err error)

	// UnflagParticipant clears a participant's flag if and only if it is currently flagged
	// with the given reason (e.g. liveliness). Returns true if a participant was unflagged.
	UnflagParticipant(
		ctx context.Context,
		chainID *big.Int,
		reason uint8,
	) (unflagged bool, err error)
}

// HeaderLivelinessProcessor is the core business logic for monitoring header proof submissions and flagging participants that fail to submit on time
type HeaderLivelinessProcessor struct {
	headerProofRepo     HeaderProofRepository
	headerFlagEventRepo HeaderFlagEventRepository
	expirationPeriod    time.Duration
	log                 logger.Logger
}

// NewHeaderLivelinessProcessor creates a new instance of the liveliness processor
func NewHeaderLivelinessProcessor(
	headerProofRepo HeaderProofRepository,
	headerFlagEventRepo HeaderFlagEventRepository,
	expirationPeriod time.Duration,
	log logger.Logger,
) *HeaderLivelinessProcessor {
	return &HeaderLivelinessProcessor{
		headerProofRepo:     headerProofRepo,
		headerFlagEventRepo: headerFlagEventRepo,
		expirationPeriod:    expirationPeriod,
		log:                 log,
	}
}

// Run continuously processes header liveliness checks on ticker events
func (h *HeaderLivelinessProcessor) Run(ctx context.Context, ticker <-chan time.Time) error {
	h.log.Info("Header liveliness handler started")

	for {
		select {
		case <-ctx.Done():
			h.log.Info("Header liveliness handler stopped")
			return nil
		case <-ticker:
			if err := h.Start(ctx); err != nil {
				h.log.Error("Error in header liveliness check",
					"error", err)
			}
		}
	}
}

// Start fetches the most recent headers and flags those that are expired
func (h *HeaderLivelinessProcessor) Start(ctx context.Context) error {
	// Get the most recent header from each chain in the DB
	headers, err := h.headerProofRepo.GetLatestHeaderProofs(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to fetch latest header proofs: %w", err))
	}

	if len(headers) == 0 {
		h.log.Debug("No header proofs found in database")
		return nil
	}

	now := time.Now()

	// For each chain, check if enough time has passed since last submission
	for _, header := range headers {
		// Calculate time since last header submission
		timeSinceLastSubmission := now.Sub(header.CreatedAt)

		// If more time has passed than allowed, flag the participant
		if timeSinceLastSubmission > h.expirationPeriod {
			flagged, err := h.headerFlagEventRepo.FlagParticipant(
				ctx,
				header.ChainID.Int,
				header.BlockNumber.Int,
				uint8(types.HeaderFlagReasonLiveliness),
				uint8(types.HeaderFlagInitiatorAutomaticSystem),
			)
			if err != nil {
				h.log.Error("Failed to flag participant for expired header",
					"chain_id", header.ChainID.String(),
					"block_number", header.BlockNumber.String(),
					"error", err)
				continue
			}

			if flagged {
				h.log.Warn("Participant flagged for missing header submission",
					"chain_id", header.ChainID.String(),
					"last_block_number", header.BlockNumber.String(),
					"reason", types.HeaderFlagReasonLiveliness)
			}
		} else {
			// Liveliness restored: the latest proof is within the expiration window again, so
			// clear any prior automatic liveliness flag (self-heal after a transient gap such as
			// a node restart or VPN outage). is_flagged must reflect CURRENT liveliness — matching
			// the expectation that an actively-proving participant is never shown as flagged.
			unflagged, err := h.headerFlagEventRepo.UnflagParticipant(
				ctx,
				header.ChainID.Int,
				uint8(types.HeaderFlagReasonLiveliness),
			)
			// A failed unflag is non-fatal: log it and keep going so a transient error on one
			// chain doesn't block liveliness processing for the others. The next tick retries.
			if err != nil {
				h.log.Error("Failed to clear liveliness flag for recovered participant",
					"chain_id", header.ChainID.String(),
					"error", err)
			} else if unflagged {
				h.log.Info("Participant unflagged: header proofs resumed within expiration window",
					"chain_id", header.ChainID.String(),
					"last_block_number", header.BlockNumber.String())
			}
		}
	}

	// Log the check result
	h.log.Info("Liveliness check completed",
		"chains_checked", len(headers))

	return nil
}
