package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/flagger/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
)

// Ensure TransactionRepository implements core.TransactionRepository at compile time
var _ core.TransactionRepository = (*TransactionRepository)(nil)

// TransactionRepository is the adapter responsible for retrieving transactions from the database for flagging operations
type TransactionRepository struct {
	db *gorm.DB
}

// zeroAddress represents the Ethereum zero address (0x000…000)
var zeroAddress = common.Address{}.Hex()

// NewTransactionRepository creates a new TransactionRepository instance for managing transaction retrieval
func NewTransactionRepository(dbClient *gorm.DB) core.TransactionRepository {
	if dbClient == nil {
		panic("dbClient is nil")
	}
	return &TransactionRepository{db: dbClient}
}

// GetTransactions retrieves up to limit transactions that are ready to be processed
// For Atomic: status 1 (executed) can be processed
// For DVP Swap: status 0 (completed) can be processed
// For Enygma: status 1 (executed) can be processed
// For NULL status: can be processed (covers Vanilla and other protocols without explicit status)
// Excludes arbitrary messages (sent by zero address)
func (tr *TransactionRepository) GetTransactions(ctx context.Context, limit int) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	err := tr.db.WithContext(ctx).
		Preload("Token").
		Order("created_at ASC").
		Limit(limit).
		Where(`
			(protocol = ? AND teleport_status = ?) OR
			(protocol = ? AND teleport_status = ?) OR
			(protocol = ? AND teleport_status = ?) OR
			teleport_status IS NULL
		`,
			types.Atomic, uint8(types.AtomicTeleportExecuted),
			types.DvpSwap, uint8(types.DvpSwapStateCompleted),
			types.Enygma, uint8(types.EnygmaTransferExecuted),
		).
		Where("is_processed = ?", false).
		Where(`"from" != ?`, zeroAddress).
		Find(&transactions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get unprocessed transactions: %w", err)
	}

	return transactions, nil
}

// MarkTransactionAsProcessed marks a transaction as processed to prevent reprocessing
func (tr *TransactionRepository) MarkTransactionAsProcessed(ctx context.Context, transactionId string) error {
	if result := getTransaction(ctx, tr.db).WithContext(ctx).
		Model(&domain.Transaction{}).
		Where("id = ?", transactionId).
		Update("is_processed", true); result.Error != nil {
		return fmt.Errorf("failed to mark transaction as processed (id=%s): %w", transactionId, result.Error)
	}
	return nil
}

// FlagTransaction flags a transaction by creating a flagged_transaction record and updating is_flagged status
// IMPORTANT: This method must be called within Transact to ensure atomicity of both operations
func (tr *TransactionRepository) FlagTransaction(ctx context.Context, transactionId string) error {
	// Create flagged transaction record
	flaggedTransaction := domain.FlaggedTransaction{
		TransactionId: uuid.MustParse(transactionId),
	}

	if result := getTransaction(ctx, tr.db).WithContext(ctx).Create(&flaggedTransaction); result.Error != nil {
		return fmt.Errorf("failed to create flagged transaction (transaction_id=%s): %w", transactionId, result.Error)
	}

	// Update the transaction record
	if result := getTransaction(ctx, tr.db).WithContext(ctx).Model(&domain.Transaction{}).
		Where("id = ?", transactionId).
		Update("is_flagged", true); result.Error != nil {
		return fmt.Errorf("failed to update transaction flag status (id=%s): %w", transactionId, result.Error)
	}

	return nil
}

// FindPendingCounterpartByReferenceId finds a PENDING, unprocessed Enygma transaction that
// shares the same reference_id (via the enygma_transactions join table) as the given
// transaction. The reference_id is the on-chain correlation key: when a cross-transfer
// fails and is reverted, both the original and the revert carry the same reference_id.
func (tr *TransactionRepository) FindPendingCounterpartByReferenceId(
	ctx context.Context,
	transactionId string,
) (*domain.Transaction, error) {
	// 1) Get the reference_id of the given transaction
	var currentEnygma domain.EnygmaTransaction
	if err := getTransaction(ctx, tr.db).WithContext(ctx).
		Where("transaction_id = ?", transactionId).
		First(&currentEnygma).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // no enygma data means not an enygma tx
		}
		return nil, fmt.Errorf("failed to get enygma transaction data: %w", err)
	}

	if currentEnygma.ReferenceId == "" {
		return nil, nil //nolint:nilnil // no reference_id to match against
	}

	// 2) Get the current transaction to know its direction
	var currentTx domain.Transaction
	if err := getTransaction(ctx, tr.db).WithContext(ctx).
		Where("id = ?", transactionId).
		First(&currentTx).Error; err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// 3) Find a PENDING, unprocessed transaction with the same reference_id
	//    in the reverse direction (from→to swapped)
	var counterpart domain.Transaction
	err := getTransaction(ctx, tr.db).WithContext(ctx).
		Joins("JOIN enygma_transactions et ON et.transaction_id = transactions.id").
		Where("et.reference_id = ?", currentEnygma.ReferenceId).
		Where("transactions.id != ?", transactionId).
		Where("transactions.from_chain_id = ?", currentTx.ToChainId).
		Where("transactions.to_chain_id = ?", currentTx.FromChainId).
		Where("transactions.protocol = ?", types.Enygma).
		Where("transactions.teleport_status = ?", uint8(types.EnygmaTransferPending)).
		Where("transactions.is_processed = ?", false).
		First(&counterpart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // no matching counterpart
		}
		return nil, fmt.Errorf("failed to find pending counterpart by reference_id: %w", err)
	}
	return &counterpart, nil
}

// Transact wraps operations in a DB transaction
func (tr *TransactionRepository) Transact(ctx context.Context, fn func(ctx context.Context) error) error {
	return tr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := withTransaction(ctx, tx)
		return fn(txCtx)
	})
}
