package core

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/flagger/testutil"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

// Helpers
func buildTransaction(txType types.TxType, amount int64, fromChain, toChain int64) domain.Transaction {
	return domain.Transaction{
		Model:       domain.Model{ID: uuid.New()},
		ResourceId:  "resource-123",
		FromChainId: domain.BigInt{Int: big.NewInt(fromChain)},
		ToChainId:   domain.BigInt{Int: big.NewInt(toChain)},
		Amount:      decimal.NewFromInt(amount),
		TxType:      txType,
		MsgType:     uint8(types.AssetTypeERC20),
	}
}

func addBalance(repo *testutil.FakeBalanceRepository, chainId, resourceId string, amount int64) {
	key := chainId + ":" + resourceId + ":"
	repo.Balances[key] = &domain.Balance{
		ChainId:    chainId,
		ResourceId: resourceId,
		Amount:     decimal.NewFromInt(amount),
		ErcId:      domain.BigInt{},
	}
}

// getBalance retrieves a balance from the fake repository
func getBalance(repo *testutil.FakeBalanceRepository, chainId, resourceId string) decimal.Decimal {
	key := chainId + ":" + resourceId + ":"
	if b, ok := repo.Balances[key]; ok {
		return b.Amount
	}
	return decimal.Zero
}

// addTransactions adds transactions to the fake repository
func addTransactions(repo *testutil.FakeTransactionRepository, txs ...domain.Transaction) {
	repo.Transactions = append(repo.Transactions, txs...)
}

func TestTransactionProcessor_ProcessesCrossChainTransfer(t *testing.T) {
	// A cross-chain transaction transferring 1 from chain 1 to chain 2
	// with initial balances of 1 (source) and 1 (destination)
	tx := buildTransaction(types.CrossChain, 1, 1, 2)

	balanceRepo := testutil.NewFakeBalanceRepository()
	addBalance(balanceRepo, "1", "resource-123", 1)
	addBalance(balanceRepo, "2", "resource-123", 1)

	txRepo := testutil.NewFakeTransactionRepository()
	addTransactions(txRepo, tx)

	processor := NewTransactionProcessor(balanceRepo, txRepo, &testutil.StubLogger{}, 100)

	err := processor.Start(context.Background())
	require.NoError(t, err)

	// The Source balance decreased by 1 (1 -> 0)
	assert.Equal(t, "0", getBalance(balanceRepo, "1", "resource-123").String())

	// The Destination balance increased by 1 (1 -> 2)
	assert.Equal(t, "2", getBalance(balanceRepo, "2", "resource-123").String())

	// The Transaction is marked as processed and in not flagged
	assert.True(t, txRepo.ProcessedIDs[tx.ID.String()])
	assert.False(t, txRepo.FlaggedIDs[tx.ID.String()])
}

func TestTransactionProcessor_FlagsInsufficientBalance(t *testing.T) {
	// Participant A transfers 2 but only has 1 according the auditor;
	// balances must remain unchanged when the transaction is flagged
	tx := buildTransaction(types.CrossChain, 2, 1, 2)

	balanceRepo := testutil.NewFakeBalanceRepository()
	addBalance(balanceRepo, "1", "resource-123", 1)
	addBalance(balanceRepo, "2", "resource-123", 0)

	txRepo := testutil.NewFakeTransactionRepository()
	addTransactions(txRepo, tx)

	processor := NewTransactionProcessor(balanceRepo, txRepo, &testutil.StubLogger{}, 100)

	err := processor.Start(context.Background())
	require.NoError(t, err)

	// Balances are not modified — the fraudulent transfer is rejected
	assert.Equal(t, "1", getBalance(balanceRepo, "1", "resource-123").String())
	assert.Equal(t, "0", getBalance(balanceRepo, "2", "resource-123").String())

	// The transaction is flagged due to negative balance
	assert.True(t, txRepo.FlaggedIDs[tx.ID.String()])

	// The transaction is still marked as processed
	assert.True(t, txRepo.ProcessedIDs[tx.ID.String()])
}

func TestTransactionProcessor_ProcessesMint(t *testing.T) {
	// A mint transaction adding 10 to chain 1
	tx := buildTransaction(types.Mint, 10, 1, 1)

	balanceRepo := testutil.NewFakeBalanceRepository()
	addBalance(balanceRepo, "1", "resource-123", 10)

	txRepo := testutil.NewFakeTransactionRepository()
	addTransactions(txRepo, tx)

	processor := NewTransactionProcessor(balanceRepo, txRepo, &testutil.StubLogger{}, 100)

	err := processor.Start(context.Background())
	require.NoError(t, err)

	// The balance increased by 10 (10 -> 20)
	assert.Equal(t, "20", getBalance(balanceRepo, "1", "resource-123").String())

	// The Transaction is marked as processed
	assert.True(t, txRepo.ProcessedIDs[tx.ID.String()])
	assert.False(t, txRepo.FlaggedIDs[tx.ID.String()])
}

func TestTransactionProcessor_ProcessesBurn(t *testing.T) {
	// A burn transaction removing 10 from chain 1
	tx := buildTransaction(types.Burn, 10, 1, 1)

	balanceRepo := testutil.NewFakeBalanceRepository()
	addBalance(balanceRepo, "1", "resource-123", 20)

	txRepo := testutil.NewFakeTransactionRepository()
	addTransactions(txRepo, tx)

	processor := NewTransactionProcessor(balanceRepo, txRepo, &testutil.StubLogger{}, 100)

	err := processor.Start(context.Background())
	require.NoError(t, err)

	// The balance decreased by 10 (20 -> 10)
	assert.Equal(t, "10", getBalance(balanceRepo, "1", "resource-123").String())

	// The transaction is marked as processed
	assert.True(t, txRepo.ProcessedIDs[tx.ID.String()])
	assert.False(t, txRepo.FlaggedIDs[tx.ID.String()])
}

func TestTransactionProcessor_FlagsBurnWithInsufficientBalance(t *testing.T) {
	// A burn transaction removing 20 but the balance is only 10 (according to the auditor knowledge);
	// balance must remain unchanged when the transaction is flagged
	tx := buildTransaction(types.Burn, 20, 1, 1)

	balanceRepo := testutil.NewFakeBalanceRepository()
	addBalance(balanceRepo, "1", "resource-123", 10)

	txRepo := testutil.NewFakeTransactionRepository()
	addTransactions(txRepo, tx)

	processor := NewTransactionProcessor(balanceRepo, txRepo, &testutil.StubLogger{}, 100)

	err := processor.Start(context.Background())
	require.NoError(t, err)

	// Balance is not modified — the fraudulent burn is rejected
	assert.Equal(t, "10", getBalance(balanceRepo, "1", "resource-123").String())

	// The transaction is flagged due to negative balance
	assert.True(t, txRepo.FlaggedIDs[tx.ID.String()])

	// The transaction is still marked as processed
	assert.True(t, txRepo.ProcessedIDs[tx.ID.String()])
}

func TestTransactionProcessor_HandlesEmptyQueue(t *testing.T) {
	// No transactions to process
	balanceRepo := testutil.NewFakeBalanceRepository()
	txRepo := testutil.NewFakeTransactionRepository()

	processor := NewTransactionProcessor(balanceRepo, txRepo, &testutil.StubLogger{}, 100)

	err := processor.Start(context.Background())
	require.NoError(t, err)
	assert.Empty(t, txRepo.ProcessedIDs)
	assert.Empty(t, txRepo.FlaggedIDs)
}

func TestTransactionProcessor_ContinuesOnUnknownType(t *testing.T) {
	// A transaction with unknown type followed by a valid mint
	unknownTxType := types.TxType(99)
	unknownTx := buildTransaction(unknownTxType, 10, 1, 1)
	mintTx := buildTransaction(types.Mint, 10, 1, 1)

	balanceRepo := testutil.NewFakeBalanceRepository()
	addBalance(balanceRepo, "1", "resource-123", 10)

	txRepo := testutil.NewFakeTransactionRepository()
	addTransactions(txRepo, unknownTx, mintTx)

	processor := NewTransactionProcessor(balanceRepo, txRepo, &testutil.StubLogger{}, 100)

	err := processor.Start(context.Background())
	require.NoError(t, err)

	// The valid mint transaction was still processed
	assert.Equal(t, "20", getBalance(balanceRepo, "1", "resource-123").String())
	assert.True(t, txRepo.ProcessedIDs[mintTx.ID.String()])
}

func TestTransactionProcessor_RunStopsOnContextCancellation(t *testing.T) {
	balanceRepo := testutil.NewFakeBalanceRepository()
	txRepo := testutil.NewFakeTransactionRepository()
	processor := NewTransactionProcessor(balanceRepo, txRepo, &testutil.StubLogger{}, 100)

	ctx, cancel := context.WithCancel(context.Background())
	ticker := make(chan time.Time)

	// Running the processor and immediately cancelling
	done := make(chan error, 1)
	go func() {
		done <- processor.Run(ctx, ticker)
	}()

	cancel()

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(1 * time.Second):
		t.Fatal("Run did not return within timeout after context cancellation")
	}
}

func TestTransactionProcessor_EnygmaRevertSkipsBalanceUpdate(t *testing.T) {
	// Simulates the Enygma failed cross-transfer scenario:
	// 1. A→B(10) executed (test 1)
	// 2. A→B(1) pending — failed on destination (test 2 original)
	// 3. B→A(1) executed — revert of the failed transfer (test 2 revert)
	//
	// The revert B→A(1) should be detected as a counterpart to the pending A→B(1)
	// and both should be marked as processed without modifying balances.

	pendingStatus := uint8(0)  // EnygmaTransferPending
	executedStatus := uint8(1) // EnygmaTransferExecuted
	sharedRefId := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"

	// Test 1: normal executed transfer A→B(10)
	normalTx := domain.Transaction{
		Model:          domain.Model{ID: uuid.New()},
		ResourceId:     "resource-enygma",
		FromChainId:    domain.BigInt{Int: big.NewInt(1)},
		ToChainId:      domain.BigInt{Int: big.NewInt(2)},
		Amount:         decimal.NewFromInt(10),
		TxType:         types.CrossChain,
		MsgType:        uint8(types.AssetTypeEnygma),
		Protocol:       types.Enygma,
		TeleportStatus: &executedStatus,
		EnygmaTransaction: &domain.EnygmaTransaction{
			ReferenceId: "different-ref-id",
		},
	}

	// Test 2 original: A→B(1) PENDING (failed on destination)
	pendingOriginal := domain.Transaction{
		Model:          domain.Model{ID: uuid.New()},
		ResourceId:     "resource-enygma",
		FromChainId:    domain.BigInt{Int: big.NewInt(1)},
		ToChainId:      domain.BigInt{Int: big.NewInt(2)},
		Amount:         decimal.NewFromInt(1),
		TxType:         types.CrossChain,
		MsgType:        uint8(types.AssetTypeEnygma),
		Protocol:       types.Enygma,
		TeleportStatus: &pendingStatus,
		EnygmaTransaction: &domain.EnygmaTransaction{
			ReferenceId: sharedRefId, // same reference_id as the revert
		},
	}

	// Test 2 revert: B→A(1) EXECUTED (revert of the failed transfer)
	revertTx := domain.Transaction{
		Model:          domain.Model{ID: uuid.New()},
		ResourceId:     "resource-enygma",
		FromChainId:    domain.BigInt{Int: big.NewInt(2)},
		ToChainId:      domain.BigInt{Int: big.NewInt(1)},
		Amount:         decimal.NewFromInt(1),
		TxType:         types.CrossChain,
		MsgType:        uint8(types.AssetTypeEnygma),
		Protocol:       types.Enygma,
		TeleportStatus: &executedStatus,
		EnygmaTransaction: &domain.EnygmaTransaction{
			ReferenceId: sharedRefId, // same reference_id as the original
		},
	}

	balanceRepo := testutil.NewFakeBalanceRepository()
	addBalance(balanceRepo, "1", "resource-enygma", 1000)
	addBalance(balanceRepo, "2", "resource-enygma", 0)

	txRepo := testutil.NewFakeTransactionRepository()
	// Only EXECUTED transactions are returned by GetTransactions
	// but the PENDING one exists in the repository for counterpart lookup
	txRepo.Transactions = []domain.Transaction{normalTx, revertTx}
	// Add the pending tx separately so FindPendingCounterpart can find it
	// but GetTransactions won't return it (it filters by EXECUTED status)
	txRepo.PendingTransactions = []domain.Transaction{pendingOriginal}

	processor := NewTransactionProcessor(balanceRepo, txRepo, &testutil.StubLogger{}, 100)

	err := processor.Start(context.Background())
	require.NoError(t, err)

	// Normal A→B(10) was processed: A=990, B=10
	assert.Equal(t, "990", getBalance(balanceRepo, "1", "resource-enygma").String())
	assert.Equal(t, "10", getBalance(balanceRepo, "2", "resource-enygma").String())

	// The normal transfer is processed
	assert.True(t, txRepo.ProcessedIDs[normalTx.ID.String()])
	// The revert is marked as processed (no-op)
	assert.True(t, txRepo.ProcessedIDs[revertTx.ID.String()])
	// The pending original is also marked as processed (no-op)
	assert.True(t, txRepo.ProcessedIDs[pendingOriginal.ID.String()])
}

func buildTxWithProtocol(protocol types.ProtocolType, msgType types.AssetType) domain.Transaction {
	return domain.Transaction{
		Model:    domain.Model{ID: uuid.New()},
		Protocol: protocol,
		MsgType:  uint8(msgType),
	}
}

func TestIsBalanceSkipped(t *testing.T) {
	processor := &TransactionProcessor{}

	tests := []struct {
		name     string
		protocol types.ProtocolType
		msgType  types.AssetType
		want     bool
	}{
		// Should be skipped
		{"DvpDeposit ERC721", types.DvpDeposit, types.AssetTypeDvpERC721, true},
		{"DvpDeposit ERC1155", types.DvpDeposit, types.AssetTypeDvpERC1155, true},
		{"DvpWithdraw ERC721", types.DvpWithdraw, types.AssetTypeDvpERC721, true},
		{"DvpWithdraw ERC1155", types.DvpWithdraw, types.AssetTypeDvpERC1155, true},
		{"DvpSwap Enygma", types.DvpSwap, types.AssetTypeEnygma, true},
		// Should NOT be skipped
		{"DvpSwap ERC721", types.DvpSwap, types.AssetTypeDvpERC721, false},
		{"DvpSwap ERC1155", types.DvpSwap, types.AssetTypeDvpERC1155, false},
		{"DvpDeposit Enygma", types.DvpDeposit, types.AssetTypeEnygma, false},
		{"Vanilla ERC20", types.Vanilla, types.AssetTypeERC20, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tx := buildTxWithProtocol(tc.protocol, tc.msgType)
			assert.Equal(t, tc.want, processor.isBalanceSkipped(tx))
		})
	}
}
