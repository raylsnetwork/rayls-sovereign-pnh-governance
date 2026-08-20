package core

import (
	"context"
	"fmt"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

// BalanceRepository defines the interface for balance operations
type BalanceRepository interface {
	CreateBalance(ctx context.Context, tx *domain.Balance) error
	UpdateBalance(ctx context.Context, chainId string, resourceId string, ercId domain.BigInt, amount string) error
	GetResourceBalance(ctx context.Context, chainId string, resourceId string, ercId string) (*domain.Balance, error)
	UpdateSenderReceiverBalances(
		ctx context.Context,
		senderChainId string,
		senderResourceId string,
		senderErcId domain.BigInt,
		senderNewAmount string,
		receiverChainId string,
		receiverResourceId string,
		receiverErcId domain.BigInt,
		receiverNewAmount string,
	) error
}

// TransactionRepository defines the interface for transaction operations
type TransactionRepository interface {
	GetTransactions(ctx context.Context, limit int) ([]domain.Transaction, error)
	MarkTransactionAsProcessed(ctx context.Context, transactionId string) error
	FlagTransaction(ctx context.Context, transactionId string) error
	Transact(ctx context.Context, fn func(ctx context.Context) error) error
	// FindPendingCounterpartByReferenceId looks up a PENDING, unprocessed Enygma transaction
	// that shares the same reference_id as the given transaction. When a cross-transfer fails,
	// the revert carries the same reference_id as the original — this is the correlation key.
	FindPendingCounterpartByReferenceId(ctx context.Context, transactionId string) (*domain.Transaction, error)
}

// TransactionProcessor is the core business logic component responsible for verifying if storage proofs are valid
type TransactionProcessor struct {
	balanceRepo BalanceRepository
	txRepo      TransactionRepository
	log         logger.Logger
	batchSize   int
}

// NewTransactionProcessor creates a new TransactionProcessor instance
func NewTransactionProcessor(
	balanceRepo BalanceRepository,
	txRepo TransactionRepository,
	log logger.Logger,
	batchSize int,
) *TransactionProcessor {
	return &TransactionProcessor{
		balanceRepo: balanceRepo,
		txRepo:      txRepo,
		log:         log,
		batchSize:   batchSize,
	}
}

// Run processes transactions continuously based on ticker events
func (tc *TransactionProcessor) Run(ctx context.Context, ticker <-chan time.Time) error {
	tc.log.Info("Transaction processor started")
	for {
		select {
		case <-ctx.Done():
			tc.log.Info("TransactionProcessor stopped due to context cancellation")
			return nil
		case <-ticker:
			if err := tc.Start(ctx); err != nil {
				tc.log.Error("Error in transaction verifying", "error", err)
			}
		}
	}
}

// Start retrieves and processes unprocessed transactions
func (tc *TransactionProcessor) Start(ctx context.Context) error {
	// Get EXECUTED or NULL transactions that haven't been processed yet
	transactionsToProcess, err := tc.txRepo.GetTransactions(ctx, tc.batchSize)
	if err != nil {
		return withstack.Wrap(err)
	}

	if len(transactionsToProcess) == 0 {
		tc.log.Debug("No transactions to process")
		return nil
	}

	// Check if balances are enough for the transaction and update their status accordingly
	err = tc.processTransactions(ctx, transactionsToProcess)
	if err != nil {
		return withstack.Wrap(err)
	}

	return nil
}

// processTransactions processes a batch of transactions, updating balances and flagging invalid ones
func (tc *TransactionProcessor) processTransactions(ctx context.Context, transactions []domain.Transaction) error {
	tc.log.Info("Processing transaction set", "count", fmt.Sprint(len(transactions)))
	var processedCount, errorCount int

	for _, transaction := range transactions {
		err := tc.processTransaction(ctx, transaction)
		if err != nil {
			tc.log.Error("Error processing transaction",
				"tx_id", transaction.ID.String(),
				"tx_type", transaction.MsgType,
				"error", err,
			)
			errorCount++
			continue
		}

		processedCount++
	}

	tc.log.Info("Transaction processing completed",
		"processed", processedCount,
		"errors", errorCount,
		"total", fmt.Sprint(len(transactions)),
	)

	if errorCount > 0 {
		return fmt.Errorf("failed to process %d out of %d transactions", errorCount, len(transactions))
	}

	return nil
}

// processTransaction processes a single transaction: updates balance and flags if source balance becomes negative
func (tc *TransactionProcessor) processTransaction(ctx context.Context, tx domain.Transaction) error {
	if tc.isBalanceSkipped(tx) {
		if err := tc.txRepo.MarkTransactionAsProcessed(ctx, tx.ID.String()); err != nil {
			return withstack.Wrap(err)
		}
		tc.log.Info("Transaction marked as processed (balance update skipped)",
			"tx_id", tx.ID.String(),
			"protocol", tx.Protocol,
			"msg_type", tx.MsgType,
		)
		return nil
	}

	switch tx.TxType {
	case types.CrossChain:
		return tc.processCrossChain(ctx, tx, getErcIdString(tx))
	case types.Mint:
		return tc.processMint(ctx, tx, getErcIdString(tx))
	case types.Burn:
		return tc.processBurn(ctx, tx, getErcIdString(tx))
	default:
		tc.log.Error("Unknown transaction type", "tx_type", tx.TxType)
		return nil
	}
}

// processCrossChain handles cross-chain transaction balance updates
func (tc *TransactionProcessor) processCrossChain(ctx context.Context, tx domain.Transaction, ercId string) error {
	// For Enygma cross-chain transfers, check if this is a revert of a failed transfer.
	// When an Enygma crossTransfer fails on the destination, the relayer creates a revert
	// transfer in the reverse direction (B→A). The original A→B stays PENDING and is never
	// processed. Processing the revert alone would incorrectly subtract from the destination
	// chain's balance. If we find a matching PENDING counterpart, both are no-ops.
	if tc.isZeroKnowledgeProtected(tx.MsgType) && tx.TxType == types.CrossChain {
		pending, err := tc.txRepo.FindPendingCounterpartByReferenceId(ctx, tx.ID.String())
		if err != nil {
			return withstack.Wrap(err)
		}
		if pending != nil {
			// Mark both the revert and the original pending tx as processed without balance changes
			if err := tc.txRepo.MarkTransactionAsProcessed(ctx, tx.ID.String()); err != nil {
				return withstack.Wrap(err)
			}
			if err := tc.txRepo.MarkTransactionAsProcessed(ctx, pending.ID.String()); err != nil {
				return withstack.Wrap(err)
			}
			tc.log.Info("Enygma revert pair detected, both transactions marked as processed (no balance change)",
				"revert_tx_id", tx.ID.String(),
				"pending_tx_id", pending.ID.String(),
				"resource_id", tx.ResourceId,
			)
			return nil
		}
	}

	// DB transaction starts here
	return tc.txRepo.Transact(ctx, func(txCtx context.Context) error {
		sourceBalance, err := tc.balanceRepo.GetResourceBalance(
			txCtx,
			tx.GetFromChainId().String(),
			tx.ResourceId,
			ercId,
		)
		if err != nil {
			return withstack.Wrap(err)
		}

		destBalance, err := tc.balanceRepo.GetResourceBalance(
			txCtx,
			tx.GetToChainId().String(),
			tx.ResourceId,
			ercId,
		)
		if err != nil {
			return withstack.Wrap(err)
		}

		// Calculate new balances
		newSourceAmount := sourceBalance.Amount.Sub(tx.Amount)

		// Flag and skip balance update if the transfer would cause a negative balance.
		// Zero Knowledge protected assets (Enygma, DVP) are exempt since ZK proofs prevent fraud.
		if newSourceAmount.IsNegative() && !tc.isZeroKnowledgeProtected(tx.MsgType) {
			if flagErr := tc.flagInvalidTransaction(txCtx, tx.ID.String()); flagErr != nil {
				return withstack.Wrap(flagErr)
			}
			if markErr := tc.txRepo.MarkTransactionAsProcessed(txCtx, tx.ID.String()); markErr != nil {
				return withstack.Wrap(markErr)
			}
			tc.log.Warn("Cross-chain transaction flagged, balance update skipped",
				"tx_id", tx.ID.String(),
				"chain_id", sourceBalance.ChainId,
				"resource_id", sourceBalance.ResourceId,
			)
			return nil
		}

		newDestAmount := destBalance.Amount.Add(tx.Amount)

		err = tc.balanceRepo.UpdateSenderReceiverBalances(
			txCtx,
			sourceBalance.ChainId,
			sourceBalance.ResourceId,
			sourceBalance.ErcId,
			newSourceAmount.String(),
			destBalance.ChainId,
			destBalance.ResourceId,
			destBalance.ErcId,
			newDestAmount.String(),
		)
		if err != nil {
			return withstack.Wrap(err)
		}

		// Mark transaction as processed
		if err := tc.txRepo.MarkTransactionAsProcessed(txCtx, tx.ID.String()); err != nil {
			return withstack.Wrap(err)
		}

		tc.log.Info("Cross-chain transaction processed successfully",
			"tx_id", tx.ID.String(),
			"chain_id", sourceBalance.ChainId,
			"resource_id", sourceBalance.ResourceId,
		)
		// Commit transaction
		return nil
	})
}

func (tc *TransactionProcessor) processMint(ctx context.Context, tx domain.Transaction, ercId string) error {
	return tc.txRepo.Transact(ctx, func(txCtx context.Context) error {
		sourceBalance, err := tc.balanceRepo.GetResourceBalance(
			txCtx,
			tx.GetFromChainId().String(),
			tx.ResourceId,
			ercId,
		)
		if err != nil {
			return withstack.Wrap(err)
		}

		newAmount := sourceBalance.Amount.Add(tx.Amount)

		err = tc.balanceRepo.UpdateBalance(
			txCtx,
			sourceBalance.ChainId,
			sourceBalance.ResourceId,
			sourceBalance.ErcId,
			newAmount.String(),
		)
		if err != nil {
			return withstack.Wrap(err)
		}

		// Mark transaction as processed
		if err := tc.txRepo.MarkTransactionAsProcessed(txCtx, tx.ID.String()); err != nil {
			return withstack.Wrap(err)
		}

		tc.log.Info("Transaction processed successfully",
			"tx_id", tx.ID.String(),
			"chain_id", sourceBalance.ChainId,
			"resource_id", sourceBalance.ResourceId,
			"operation", "mint",
		)

		// Commit transaction
		return nil
	})
}

func (tc *TransactionProcessor) processBurn(ctx context.Context, tx domain.Transaction, ercId string) error {
	return tc.txRepo.Transact(ctx, func(txCtx context.Context) error {
		sourceBalance, err := tc.balanceRepo.GetResourceBalance(
			txCtx,
			tx.GetFromChainId().String(),
			tx.ResourceId,
			ercId,
		)
		if err != nil {
			return withstack.Wrap(err)
		}

		newAmount := sourceBalance.Amount.Sub(tx.Amount)

		// Flag and skip balance update if the burn would cause a negative balance.
		// Zero Knowledge protected assets (Enygma, DVP) are exempt since ZK proofs prevent fraud.
		if newAmount.IsNegative() && !tc.isZeroKnowledgeProtected(tx.MsgType) {
			if flagErr := tc.flagInvalidTransaction(txCtx, tx.ID.String()); flagErr != nil {
				return withstack.Wrap(flagErr)
			}
			if markErr := tc.txRepo.MarkTransactionAsProcessed(txCtx, tx.ID.String()); markErr != nil {
				return withstack.Wrap(markErr)
			}
			tc.log.Warn("Burn transaction flagged, balance update skipped",
				"tx_id", tx.ID.String(),
				"chain_id", sourceBalance.ChainId,
				"resource_id", sourceBalance.ResourceId,
			)
			return nil
		}

		err = tc.balanceRepo.UpdateBalance(
			txCtx,
			sourceBalance.ChainId,
			sourceBalance.ResourceId,
			sourceBalance.ErcId,
			newAmount.String(),
		)
		if err != nil {
			return withstack.Wrap(err)
		}

		// Mark transaction as processed
		if err := tc.txRepo.MarkTransactionAsProcessed(txCtx, tx.ID.String()); err != nil {
			return withstack.Wrap(err)
		}

		tc.log.Info("Transaction processed successfully",
			"tx_id", tx.ID.String(),
			"chain_id", sourceBalance.ChainId,
			"resource_id", sourceBalance.ResourceId,
			"operation", "burn",
		)

		// Commit transaction
		return nil
	})
}

// flagInvalidTransaction flags a transaction due to negative balance
func (tc *TransactionProcessor) flagInvalidTransaction(ctx context.Context, txId string) error {
	if err := tc.txRepo.FlagTransaction(ctx, txId); err != nil {
		tc.log.Error("Error flagging transaction",
			"tx_id", txId,
			"error", err,
		)
		return withstack.Wrap(err)
	}

	tc.log.Warn("Transaction flagged due to negative balance",
		"tx_id", txId,
	)

	return nil
}

// isBalanceSkipped returns true for transaction types where balance tracking does not apply:
//   - DvP deposit/withdraw of DvP ERC721 or ERC1155 tokens (balance only updates during swap operations).
//   - DvP swap involving Enygma tokens (balance only updates during deposit/withdraw).
func (tc *TransactionProcessor) isBalanceSkipped(tx domain.Transaction) bool {
	isDvpNFTDepositWithdraw := (tx.Protocol == types.DvpDeposit || tx.Protocol == types.DvpWithdraw) &&
		(tx.MsgType == uint8(types.AssetTypeDvpERC721) || tx.MsgType == uint8(types.AssetTypeDvpERC1155))

	isDvpEnygmaSwap := tx.Protocol == types.DvpSwap && tx.MsgType == uint8(types.AssetTypeEnygma)

	return isDvpNFTDepositWithdraw || isDvpEnygmaSwap
}

// isZeroKnowledgeProtected returns true if the asset type is protected by Zero Knowledge proofs
// These asset types (Enygma, DVP ERC721/1155) don't need balance validation as ZK prevents issues
func (tc *TransactionProcessor) isZeroKnowledgeProtected(msgType uint8) bool {
	return msgType == uint8(types.AssetTypeEnygma) ||
		msgType == uint8(types.AssetTypeDvpERC721) ||
		msgType == uint8(types.AssetTypeDvpERC1155)
}

// getErcIdString returns the ercId as a string or empty string if nil
func getErcIdString(tx domain.Transaction) string {
	if tx.GetErcId() != nil {
		return tx.GetErcId().String()
	}
	return ""
}
