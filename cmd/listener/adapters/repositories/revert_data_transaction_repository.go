package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
)

// Ensure RevertDataTransactionRepository implements core.RevertDataTransactionRepository at compile time
var _ core.RevertDataTransactionRepository = (*RevertDataTransactionRepository)(nil)

// RevertDataTransactionRepository implements core.RevertDataTransactionRepository
type RevertDataTransactionRepository struct {
	db *gorm.DB
}

// NewRevertDataTransactionRepository creates a new adapter
func NewRevertDataTransactionRepository(dbClient *gorm.DB) core.RevertDataTransactionRepository {
	if dbClient == nil {
		panic("dbClient is nil")
	}
	return &RevertDataTransactionRepository{db: dbClient}
}

// CreateRevertTransactions creates multiple revert transaction rows (idempotent based on transaction_id)
func (r *RevertDataTransactionRepository) CreateRevertTransactions(
	ctx context.Context,
	revert_transactions []*domain.RevertDataTransaction,
) error {
	if len(revert_transactions) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "transaction_id"}},
			DoNothing: true,
		}).Create(&revert_transactions).Error; err != nil {
			return fmt.Errorf("failed to create revert transactions: %w", err)
		}
		return nil
	})
}
