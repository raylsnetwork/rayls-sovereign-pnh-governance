package repositories

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
)

// Ensure HeaderProofEventRepository implements core.HeaderProofEventRepository at compile time
var _ core.HeaderProofEventRepository = (*HeaderProofEventRepository)(nil)

// HeaderProofEventRepository implements core.HeaderProofEventRepository
type HeaderProofEventRepository struct {
	db *gorm.DB
}

// NewHeaderProofEventRepository creates a new header proof event repository adapter
func NewHeaderProofEventRepository(dbClient *gorm.DB) core.HeaderProofEventRepository {
	if dbClient == nil {
		panic("dbClient is nil")
	}

	return &HeaderProofEventRepository{
		db: dbClient,
	}
}

// Create stores a new header proof event
func (r *HeaderProofEventRepository) Create(ctx context.Context, event *domain.HeaderProofEvent) error {
	// Use OnConflict to ensure idempotency, does nothing if chain_id+block_number already exists
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(event).Error; err != nil {
		return fmt.Errorf("failed to create header proof event: %w", err)
	}
	return nil
}

// CreateBatch inserts multiple header proof events in a single statement.
// Duplicate (chain_id, block_number) pairs are silently ignored via ON CONFLICT DO NOTHING.
func (r *HeaderProofEventRepository) CreateBatch(ctx context.Context, events []*domain.HeaderProofEvent) error {
	if len(events) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&events).Error; err != nil {
			return fmt.Errorf("failed to batch-create header proof events: %w", err)
		}
		return nil
	})
}

// GetByBlockNumber retrieves a header proof event by chain ID and block number
func (r *HeaderProofEventRepository) GetByBlockNumber(
	ctx context.Context,
	chainID *big.Int,
	blockNumber *big.Int,
) (*domain.HeaderProofEvent, error) {
	var event domain.HeaderProofEvent
	err := r.db.WithContext(ctx).Where("chain_id = ? AND block_number = ?",
		chainID.String(), blockNumber.String()).
		First(&event).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // nil signals record not found without wrapping a sentinel error
		}
		return nil, fmt.Errorf(
			"failed to get header proof event for chain_id=%s block_number=%s: %w",
			chainID.String(),
			blockNumber.String(),
			err,
		)
	}
	return &event, nil
}
