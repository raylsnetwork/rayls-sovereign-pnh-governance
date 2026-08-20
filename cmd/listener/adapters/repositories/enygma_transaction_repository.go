package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
)

// Ensure EnygmaTransactionRepository implements core.EnygmaTransactionRepository at compile time
var _ core.EnygmaTransactionRepository = (*EnygmaTransactionRepository)(nil)

// EnygmaTransactionRepository implements core.EnygmaTransactionRepository
type EnygmaTransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new adapter
func NewEnygmaTransactionRepository(dbClient *gorm.DB) core.EnygmaTransactionRepository {
	if dbClient == nil {
		panic("dbClient is nil")
	}
	return &EnygmaTransactionRepository{db: dbClient}
}

// UpsertEnygmaTransactionsBulk performs atomic bulk updates on transactions
func (r *EnygmaTransactionRepository) UpsertEnygmaTransactionsBulk(
	ctx context.Context,
	enygma_transactions []*domain.EnygmaTransaction,
) error {
	if len(enygma_transactions) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "transaction_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
			}).
			Create(&enygma_transactions).Error; err != nil {
			return fmt.Errorf("failed to upsert enygma transactions: %w", err)
		}
		return nil
	})
}
