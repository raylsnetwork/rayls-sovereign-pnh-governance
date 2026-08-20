package repositories

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

// Ensure TransactionRepository implements core.TransactionRepository at compile time
var _ core.TransactionRepository = (*TransactionRepository)(nil)

// TransactionRepository implements core.TransactionRepository
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new adapter
func NewTransactionRepository(dbClient *gorm.DB) core.TransactionRepository {
	if dbClient == nil {
		panic("dbClient is nil")
	}
	return &TransactionRepository{db: dbClient}
}

// helper to optionally preload common associations
func (r *TransactionRepository) applyPreloads(q *gorm.DB, preload bool) *gorm.DB {
	if !preload {
		return q
	}

	return q.Preload("Token")
}

// CreateTransaction creates a single transaction with idempotency based on composite unique key
func (r *TransactionRepository) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	if tx == nil {
		return fmt.Errorf("failed to create transaction: tx is nil")
	}

	// This prevents duplicate transactions if the same transaction is reprocessed
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "message_id"},
			{Name: "block_number"},
			{Name: "log_index"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
	}).Create(tx).Error
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

// CreateTransactions creates multiple transactions in a single batch operation
func (r *TransactionRepository) CreateTransactions(ctx context.Context, transactions []*domain.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return r.createTransactions(tx, transactions)
	})
}

// CreateTransactionsWithEnygmaData creates both transactions and enygma transactions in a single database transaction
func (r *TransactionRepository) CreateTransactionsWithEnygmaData(
	ctx context.Context,
	transactions []*domain.Transaction,
	enygmaTransactions []*domain.EnygmaTransaction,
) error {
	if len(transactions) == 0 {
		return nil
	}
	if len(enygmaTransactions) != len(transactions) {
		return fmt.Errorf(
			"transaction count (%d) must equal enygma count (%d)",
			len(transactions),
			len(enygmaTransactions),
		)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.createTransactions(tx, transactions); err != nil {
			return err
		}

		if len(enygmaTransactions) > 0 {
			for i := range enygmaTransactions {
				enygmaTransactions[i].TransactionId = transactions[i].ID
			}

			if err := tx.Clauses(
				clause.OnConflict{
					Columns:   []clause.Column{{Name: "transaction_id"}},
					DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
				},
			).Create(&enygmaTransactions).Error; err != nil {
				return fmt.Errorf("failed to create enygma transactions: %w", err)
			}
		}

		return nil
	})
}

func (r *TransactionRepository) GetTransactionByMessageID(
	ctx context.Context,
	messageID string,
	preload bool,
) (*domain.Transaction, error) {
	var t domain.Transaction
	q := r.applyPreloads(r.db.WithContext(ctx), preload)
	if err := q.First(&t, "message_id = ?", messageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction not found (message_id=%s): %w", messageID, err)
		}
		return nil, fmt.Errorf("failed to get transaction (message_id=%s): %w", messageID, err)
	}
	return &t, nil
}

func (r *TransactionRepository) GetTransactionsBySharedIDs(
	ctx context.Context,
	sharedIDs []string,
	preload bool,
) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	q := r.applyPreloads(r.db.WithContext(ctx), preload)
	if err := q.Where("shared_id IN (?)", sharedIDs).Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction not found (shared_id=%s): %w", sharedIDs, err)
		}
		return nil, fmt.Errorf("failed to get transaction (message_id=%s): %w", sharedIDs, err)
	}
	return transactions, nil
}

func (r *TransactionRepository) GetTransactionsByMessageIDs(
	ctx context.Context,
	messageIDs []string,
	preload bool,
) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	q := r.applyPreloads(r.db.WithContext(ctx), preload)
	if err := q.Where("message_id IN (?)", messageIDs).Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction not found (message_ids=%v): %w", messageIDs, err)
		}
		return nil, fmt.Errorf("failed to get transactions (message_ids=%v): %w", messageIDs, err)
	}
	return transactions, nil
}

func (r *TransactionRepository) UpdateTransaction(ctx context.Context, tx *domain.Transaction) error {
	if tx == nil {
		return errors.New("UpdateTransaction: tx is nil")
	}
	if err := r.db.WithContext(ctx).Save(tx).Error; err != nil {
		return fmt.Errorf("failed to update transaction (id=%v): %w", tx.ID, err)
	}
	return nil
}

// CreateCreditMemo creates a credit memo transaction (idempotent based on message_id)
func (r *TransactionRepository) CreateCreditMemo(ctx context.Context, tx *domain.Transaction) error {
	creditMemoStatus := uint8(types.AtomicTeleportCreditMemo)

	transactionCredit := &domain.Transaction{
		ResourceId:     tx.ResourceId,
		FromChainId:    domain.NewBigInt(tx.GetToChainId()),   // Reverse the chain ids
		ToChainId:      domain.NewBigInt(tx.GetFromChainId()), // Reverse the chain ids
		Amount:         tx.Amount,
		TxType:         tx.TxType,
		SharedId:       tx.SharedId,
		TeleportStatus: &creditMemoStatus,
		HubTxHash:      tx.HubTxHash,
		BlockNumber:    tx.BlockNumber,
		MsgType:        tx.MsgType,
		MessageId:      tx.MessageId,
		From:           tx.To,   // Reverse the addresses
		To:             tx.From, // Reverse the addresses
		Protocol:       tx.Protocol,
		ErcId:          domain.NewBigInt(tx.GetErcId()),
	}

	if err := r.CreateTransaction(ctx, transactionCredit); err != nil {
		return fmt.Errorf("failed to create credit memo transaction: %w", err)
	}

	return nil
}

func (r *TransactionRepository) PersistUpdateTransaction(
	ctx context.Context,
	ut *types.UpdateTransaction,
	fromChain, toChain string,
) (uuid.UUID, error) {
	if ut == nil {
		return uuid.Nil, errors.New("PersistUpdateTransaction: update transaction is nil")
	}

	fromChainBigInt, _ := new(big.Int).SetString(fromChain, 10)
	toChainBigInt, _ := new(big.Int).SetString(toChain, 10)
	blockNumberStr, err := decimal.NewFromString(ut.BlockNumber)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse block number: %w", err)
	}

	msg := &domain.Transaction{
		FromChainId:  domain.NewBigInt(fromChainBigInt),
		ToChainId:    domain.NewBigInt(toChainBigInt),
		ResourceId:   hex.EncodeToString(ut.ResourceId[:]),
		BlockNumber:  blockNumberStr,
		Amount:       decimal.NewFromBigInt(ut.Amount, 0),
		HubTxHash:    ut.TxHash,
		HubTimestamp: ut.HubTimestamp,
		TxType:       ut.UpdateType,
		ErcId:        domain.NewBigInt(ut.ErcId),
		MsgType:      ut.MsgType,
		LogIndex:     ut.LogIndex,
	}

	if err := r.CreateTransaction(ctx, msg); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create update transaction: %w", err)
	}

	return msg.ID, nil
}

func (r *TransactionRepository) PersistUpdateTransactions(
	ctx context.Context,
	uts *[]types.UpdateTransaction,
	fromChain, toChain string,
) ([]uuid.UUID, error) {
	if uts == nil || len(*uts) == 0 {
		return []uuid.UUID{}, nil
	}

	msgs := make([]*domain.Transaction, len(*uts))
	fromChainBigInt, _ := new(big.Int).SetString(fromChain, 10)
	toChainBigInt, _ := new(big.Int).SetString(toChain, 10)

	for i, ut := range *uts {
		blockNumberStr, err := decimal.NewFromString(ut.BlockNumber)
		if err != nil {
			return []uuid.UUID{}, fmt.Errorf("failed to parse block number for transaction %d: %w", i, err)
		}

		msgs[i] = &domain.Transaction{
			FromChainId:  domain.NewBigInt(fromChainBigInt),
			ToChainId:    domain.NewBigInt(toChainBigInt),
			ResourceId:   hex.EncodeToString(ut.ResourceId[:]),
			BlockNumber:  blockNumberStr,
			Amount:       decimal.NewFromBigInt(ut.Amount, 0),
			HubTxHash:    ut.TxHash,
			HubTimestamp: ut.HubTimestamp,
			TxType:       ut.UpdateType,
			ErcId:        domain.NewBigInt(ut.ErcId),
			MsgType:      ut.MsgType,
			LogIndex:     ut.LogIndex,
		}
	}

	if err := r.CreateTransactions(ctx, msgs); err != nil {
		return []uuid.UUID{}, fmt.Errorf("failed to create update transactions: %w", err)
	}

	msgIds := make([]uuid.UUID, len(msgs))
	for i, m := range msgs {
		msgIds[i] = m.ID
	}

	return msgIds, nil
}

// UpdateTeleportStatusBySharedID updates teleport_status for all transactions with the given shared_id.
//
// Unlike UpdateDvpTeleportConfirmation, RowsAffected is not checked: terminal
// status transitions (SwapCancelled, SwapTimedOut) may legitimately target a
// sharedID this listener never observed the matching SwapInitiated for —
// foreign swaps, or swaps whose init preceded the configured starting block.
// Silently no-op'ing in that case mirrors the ack-and-skip behavior the
// DvpTeleport handler applies when swap salts are missing, and avoids a
// redelivery loop for cancellations that can never match a row here.
func (r *TransactionRepository) UpdateTeleportStatusBySharedID(
	ctx context.Context,
	sharedID string,
	status uint8,
) error {
	result := r.db.WithContext(ctx).
		Model(&domain.Transaction{}).
		Where("shared_id = ?", sharedID).
		Update("teleport_status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update teleport_status for shared_id %s: %w", sharedID, result.Error)
	}

	return nil
}

// UpdateDvpTeleportConfirmation fills in the missing chain-side data on both DVP transactions
// created during initiation.
func (r *TransactionRepository) UpdateDvpTeleportConfirmation(
	ctx context.Context,
	sharedID string,
	chainID *big.Int,
	to string,
	pnTxHash string,
	pnTxTimestamp time.Time,
	calldata string,
) error {
	chainIDStr := chainID.String()

	sql := `
		UPDATE transactions
		SET
			to_chain_id           = CASE WHEN from_chain_id IS NOT NULL THEN ? ELSE to_chain_id           END,
			"to"                  = CASE WHEN from_chain_id IS NOT NULL THEN ? ELSE "to"                  END,
			from_chain_id         = CASE WHEN to_chain_id   IS NOT NULL THEN ? ELSE from_chain_id         END,
			"from"                = CASE WHEN to_chain_id   IS NOT NULL THEN ? ELSE "from"                END,
			payload               = CASE WHEN to_chain_id   IS NOT NULL THEN ? ELSE payload               END,
			tx_hash_destination   = ?,
			destination_timestamp = ?,
			updated_at            = NOW()
		WHERE shared_id = ?
	`

	args := []any{
		// TokenIn row (from_chain_id IS NOT NULL) → fill destination side
		chainIDStr, to,
		// TokenOut row (to_chain_id IS NOT NULL) → fill source side
		chainIDStr, to, calldata,
		// Both rows
		pnTxHash, pnTxTimestamp,
		sharedID,
	}

	result := r.db.WithContext(ctx).Exec(sql, args...)
	if result.Error != nil {
		return fmt.Errorf("failed to update DVP teleport confirmation for shared_id %s: %w", sharedID, result.Error)
	}
	// 0 rows: event was already processed (idempotent replay).
	// 1 row: partial update indicates a data integrity issue and must fail.
	// 2 rows: expected outcome.
	if result.RowsAffected == 1 {
		return fmt.Errorf(
			"partial update for shared_id %s: expected 2 rows affected but got 1; possible data integrity issue",
			sharedID,
		)
	}

	return nil
}

// UpdateTransactionsBulk performs atomic bulk updates on transactions.
// columns restricts which columns are written; if empty, defaults to "updated_at" only.
func (r *TransactionRepository) UpdateTransactionsBulk(
	ctx context.Context,
	transactions []*domain.Transaction,
	columns ...string,
) error {
	if len(transactions) == 0 {
		return nil
	}
	if len(columns) == 0 {
		columns = []string{"updated_at"}
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, transaction := range transactions {
			if err := tx.Model(transaction).Select(columns).Updates(transaction).Error; err != nil {
				return fmt.Errorf("failed to update transaction %s: %w", transaction.ID, err)
			}
		}
		return nil
	})
}

// CreateTransactionsWithPromotion creates transactions and promotes batches atomically.
// All operations happen in a single database transaction.
func (r *TransactionRepository) CreateTransactionsWithPromotion(
	ctx context.Context,
	transactions []*domain.Transaction,
	batchIds []string,
) error {
	if len(transactions) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create transactions
		if err := r.createTransactions(tx, transactions); err != nil {
			return err
		}

		// 2. Promote batches that have >= 2 transactions
		for _, batchId := range batchIds {
			var count int64
			if err := tx.Model(&domain.Transaction{}).
				Where("batch_id = ?", batchId).
				Count(&count).Error; err != nil {
				return fmt.Errorf("failed to count batch %s: %w", batchId, err)
			}

			if count >= 2 {
				if err := tx.Model(&domain.Transaction{}).
					Where("batch_id = ?", batchId).
					Updates(map[string]any{
						"aggregation_type": domain.AggregationTypeRegularBatch,
						"aggregation_key":  batchId,
					}).Error; err != nil {
					return fmt.Errorf("failed to promote batch %s: %w", batchId, err)
				}
			}
		}

		return nil
	})
}

// createTransactions creates transactions within an existing database transaction
func (r *TransactionRepository) createTransactions(tx *gorm.DB, transactions []*domain.Transaction) error {
	return tx.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "message_id"},
				{Name: "block_number"},
				{Name: "log_index"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
		},
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
	).Create(&transactions).Error
}
