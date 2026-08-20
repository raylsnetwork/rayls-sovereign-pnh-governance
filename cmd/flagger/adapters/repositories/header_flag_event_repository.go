package repositories

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/flagger/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
)

// HeaderFlagEventRepository is the GORM adapter for header flag events and participant flagging
type HeaderFlagEventRepository struct {
	db *gorm.DB
}

// Compile-time check to ensure HeaderFlagEventRepository implements the interface
var _ core.HeaderFlagEventRepository = (*HeaderFlagEventRepository)(nil)

// NewHeaderFlagEventRepository creates a new GORM-based header flag event repository
func NewHeaderFlagEventRepository(db *gorm.DB) core.HeaderFlagEventRepository {
	return &HeaderFlagEventRepository{db: db}
}

// FlagParticipant creates an audit trail event and updates participant status if not already flagged
// Returns true if the participant was newly flagged, false if already flagged
func (r *HeaderFlagEventRepository) FlagParticipant(
	ctx context.Context,
	chainID *big.Int,
	blockNumber *big.Int,
	reason uint8,
	initiator uint8,
) (bool, error) {
	if chainID == nil || blockNumber == nil {
		return false, fmt.Errorf("chainID or blockNumber cannot be nil")
	}

	var flagged bool

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Insert audit trail record with ON CONFLICT DO NOTHING for idempotency
		event := domain.HeaderFlagEvent{
			ChainID:     domain.BigInt{Int: chainID},
			BlockNumber: domain.BigInt{Int: blockNumber},
			Reason:      reason,
			Initiator:   initiator,
			CreatedAt:   now,
		}

		// Use raw SQL with ON CONFLICT to handle duplicates
		result := tx.Exec(`
			INSERT INTO header_flag_events (chain_id, block_number, reason, initiator, created_at)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT (chain_id, block_number) DO NOTHING
		`, event.ChainID.String(), event.BlockNumber.String(), event.Reason, event.Initiator, event.CreatedAt)

		if result.Error != nil {
			return fmt.Errorf("failed to create header flag event: %w", result.Error)
		}

		// Update participant to set flagged status (only if not already flagged)
		flagReasonStr, ok := types.HeaderFlagReasonToString[reason]
		if !ok {
			return fmt.Errorf("unknown flag reason: %d", reason)
		}
		result = tx.Model(&domain.Participant{}).
			Where("chain_id = ? AND is_flagged = false", chainID.String()).
			Updates(map[string]interface{}{
				"is_flagged":  true,
				"flag_reason": flagReasonStr,
				"flagged_at":  now,
			})

		if result.Error != nil {
			return fmt.Errorf("failed to update participant flag status: %w", result.Error)
		}

		flagged = result.RowsAffected > 0

		return nil
	})

	return flagged, err
}

// UnflagParticipant clears a participant's flag if and only if it is currently flagged with
// the given reason. Scoping to the reason leaves flags raised for other reasons in place for
// manual review. The header_flag_events audit trail is intentionally preserved. Returns true
// if a participant row was unflagged.
func (r *HeaderFlagEventRepository) UnflagParticipant(
	ctx context.Context,
	chainID *big.Int,
	reason uint8,
) (bool, error) {
	if chainID == nil {
		return false, fmt.Errorf("chainID cannot be nil")
	}

	flagReasonStr, ok := types.HeaderFlagReasonToString[reason]
	if !ok {
		return false, fmt.Errorf("unknown flag reason: %d", reason)
	}

	result := r.db.WithContext(ctx).
		Model(&domain.Participant{}).
		Where("chain_id = ? AND is_flagged = true AND flag_reason = ?", chainID.String(), flagReasonStr).
		Updates(map[string]interface{}{
			"is_flagged":  false,
			"flag_reason": nil,
			"flagged_at":  nil,
		})

	if result.Error != nil {
		return false, fmt.Errorf("failed to clear participant flag status: %w", result.Error)
	}

	return result.RowsAffected > 0, nil
}
