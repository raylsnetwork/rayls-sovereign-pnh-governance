package core

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/testutil"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
)

// Helpers

func buildServiceTransaction(messageId string, fromChain, toChain int64) domain.Transaction {
	return domain.Transaction{
		MessageId:   messageId,
		FromChainId: domain.BigInt{Int: big.NewInt(fromChain)},
		ToChainId:   domain.BigInt{Int: big.NewInt(toChain)},
	}
}

func withResourceId(tx domain.Transaction, resourceId string) domain.Transaction {
	tx.ResourceId = resourceId
	return tx
}

func withServiceAddresses(tx domain.Transaction, from, to string) domain.Transaction {
	tx.From = from
	tx.To = to
	return tx
}

func withAggregationKey(tx domain.Transaction, key string) domain.Transaction {
	tx.AggregationKey = key
	return tx
}

func withAggregationType(tx domain.Transaction, aggType string) domain.Transaction {
	tx.AggregationType = aggType
	return tx
}

func TestTransactionService_GetTransactionByMessageId_ReturnsTransaction(t *testing.T) {
	// Querying a transaction by messageId returns the transaction details
	repo := testutil.NewFakeTransactionRepository()
	tx := buildServiceTransaction("msg-123", 1, 2)
	tx = withResourceId(tx, "resource-abc")
	tx = withServiceAddresses(
		tx,
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222",
	)
	repo.Transactions = []domain.Transaction{tx}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionByMessageId(context.Background(), "msg-123")

	require.NoError(t, err)
	assert.Equal(t, "msg-123", result.MessageId)
	assert.Equal(t, "resource-abc", result.ResourceId)
}

func TestTransactionService_GetTransactionByMessageId_NotFound(t *testing.T) {
	// Querying a non-existent transaction returns a NotFoundError
	repo := testutil.NewFakeTransactionRepository()
	repo.NotFoundErr = ErrRecordNotFound

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionByMessageId(context.Background(), "nonexistent")

	require.Error(t, err)
	assert.Nil(t, result)
	var notFoundErr *NotFoundError
	assert.True(t, errors.As(err, &notFoundErr), "expected NotFoundError, got %T", err)
}

func TestTransactionService_GetTransactionsList_ReturnsTransactions(t *testing.T) {
	// Listing transactions with source chain filter returns matching transactions
	repo := testutil.NewFakeTransactionRepository()
	repo.Transactions = []domain.Transaction{
		buildServiceTransaction("msg-1", 1, 2),
		buildServiceTransaction("msg-2", 1, 3),
	}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		FromChainId: "1",
		Page:        1,
		Limit:       10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 2)
	assert.Equal(t, int64(2), result.Total)
	assert.Equal(t, "1", result.Data[0].SourceChainId)
	assert.Equal(t, "1", result.Data[1].SourceChainId)
}

func TestTransactionService_GetTransactionsList_Pagination(t *testing.T) {
	// Pagination metadata is calculated correctly
	repo := testutil.NewFakeTransactionRepository()
	repo.Transactions = []domain.Transaction{
		withAggregationKey(buildServiceTransaction("msg-1", 1, 2), "tx-1"),
		withAggregationKey(buildServiceTransaction("msg-2", 1, 2), "tx-2"),
		withAggregationKey(buildServiceTransaction("msg-3", 1, 2), "tx-3"),
		withAggregationKey(buildServiceTransaction("msg-4", 1, 2), "tx-4"),
		withAggregationKey(buildServiceTransaction("msg-5", 1, 2), "tx-5"),
	}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		Page:  2,
		Limit: 2,
	})

	require.NoError(t, err)
	assert.Len(t, result.Data, 5)
	assert.Equal(t, int64(5), result.Total)
	assert.Equal(t, 2, result.Page)
	assert.NotEmpty(t, result.Data[0].Id)
	assert.Equal(t, "1", result.Data[0].SourceChainId)
}

func TestTransactionService_GetFlaggedTransactions_ReturnsFlaggedList(t *testing.T) {
	// Fetching flagged transactions returns all flagged transaction records
	flagId1 := uuid.New()
	flagId2 := uuid.New()
	txId1 := uuid.New()
	txId2 := uuid.New()

	repo := testutil.NewFakeTransactionRepository()
	repo.FlaggedTxs = []domain.FlaggedTransaction{
		{Model: domain.Model{ID: flagId1}, TransactionId: txId1},
		{Model: domain.Model{ID: flagId2}, TransactionId: txId2},
	}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetFlaggedTransactions(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, flagId1, result[0].ID)
	assert.Equal(t, txId1, result[0].TransactionId)
	assert.Equal(t, flagId2, result[1].ID)
	assert.Equal(t, txId2, result[1].TransactionId)
}

func TestTransactionService_GetFlaggedTransactions_EmptyWhenNoFlagged(t *testing.T) {
	// Empty repository returns empty flagged transactions list
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetFlaggedTransactions(context.Background())

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestTransactionService_GetRegularBatchTransactions_ReturnsBatchTransactions(t *testing.T) {
	// Fetching a regular batch returns all transactions in the batch
	repo := testutil.NewFakeTransactionRepository()
	tx1 := buildServiceTransaction("msg-1", 1, 2)
	tx1.From = "0x1111111111111111111111111111111111111111"
	tx2 := buildServiceTransaction("msg-2", 1, 2)
	tx2.From = "0x2222222222222222222222222222222222222222"
	tx3 := buildServiceTransaction("msg-3", 1, 2)
	tx3.From = "0x3333333333333333333333333333333333333333"
	repo.BatchTransactions = []domain.Transaction{tx1, tx2, tx3}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetRegularBatchTransactions(context.Background(), "batch-123", 1, 10)

	require.NoError(t, err)
	require.Len(t, result.Data, 3)
	assert.Equal(t, int64(3), result.Total)
	assert.Equal(t, "msg-1", result.Data[0].MessageId)
	assert.Equal(t, "0x1111111111111111111111111111111111111111", result.Data[0].SourceAddress)
	assert.Equal(t, "1", result.Data[0].SourceChainId)
	assert.Equal(t, "msg-2", result.Data[1].MessageId)
	assert.Equal(t, "0x2222222222222222222222222222222222222222", result.Data[1].SourceAddress)
	assert.Equal(t, "msg-3", result.Data[2].MessageId)
	assert.Equal(t, "0x3333333333333333333333333333333333333333", result.Data[2].SourceAddress)
}

func TestTransactionService_GetRegularBatchTransactions_NotFoundWhenEmpty(t *testing.T) {
	// Querying a non-existent batch returns not found error
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetRegularBatchTransactions(context.Background(), "nonexistent", 1, 10)

	require.Error(t, err)
	assert.Nil(t, result)
	var notFoundErr *NotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestTransactionService_GetEnygmaBatchTransactions_ReturnsBatchTransactions(t *testing.T) {
	// Fetching an Enygma batch returns all transactions in the batch
	repo := testutil.NewFakeTransactionRepository()
	tx1 := buildServiceTransaction("msg-enygma-1", 100, 200)
	tx1.From = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	tx2 := buildServiceTransaction("msg-enygma-2", 100, 200)
	tx2.From = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	repo.EnygmaBatchTxs = []domain.Transaction{tx1, tx2}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetEnygmaBatchTransactions(context.Background(), "enygma-batch-456", 1, 10)

	require.NoError(t, err)
	require.Len(t, result.Data, 2)
	assert.Equal(t, int64(2), result.Total)
	assert.Equal(t, "msg-enygma-1", result.Data[0].MessageId)
	assert.Equal(t, "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", result.Data[0].SourceAddress)
	assert.Equal(t, "100", result.Data[0].SourceChainId)
	assert.Equal(t, "200", result.Data[0].DestinationChainId)
	assert.Equal(t, "msg-enygma-2", result.Data[1].MessageId)
	assert.Equal(t, "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", result.Data[1].SourceAddress)
}

func TestTransactionService_GetTransactionsList_ReturnsCombinedFilterResults(t *testing.T) {
	// Combined source and destination chain filters return matching transactions
	repo := testutil.NewFakeTransactionRepository()
	repo.Transactions = []domain.Transaction{
		buildServiceTransaction("msg-1", 1, 2),
		buildServiceTransaction("msg-4", 1, 2),
	}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		FromChainId: "1",
		ToChainId:   "2",
		Page:        1,
		Limit:       10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 2)
	assert.Equal(t, "1", result.Data[0].SourceChainId)
	assert.Equal(t, "1", result.Data[1].SourceChainId)
}

func TestTransactionService_GetTransactionsList_ReturnsAddressFilterResults(t *testing.T) {
	// Address filter is case-insensitive and returns matching transactions
	alice := "0xabcdef1234567890abcdef1234567890abcdef12"
	repo := testutil.NewFakeTransactionRepository()
	tx := withAggregationType(buildServiceTransaction("msg-1", 1, 2), domain.AggregationTypeTransaction)
	tx = withServiceAddresses(tx, alice, "0x1234567890abcdef1234567890abcdef12345678")
	repo.Transactions = []domain.Transaction{tx}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		From:  "0xABCDEF1234567890ABCDEF1234567890ABCDEF12",
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, "msg-1", result.Data[0].Id)
	assert.Equal(t, alice, result.Data[0].SourceAddress)
}

func TestTransactionService_GetTransactionsList_ReturnsResourceIdFilterResults(t *testing.T) {
	// ResourceId filter returns matching transactions
	repo := testutil.NewFakeTransactionRepository()
	tx1 := withAggregationKey(buildServiceTransaction("msg-1", 1, 2), "tx-1")
	tx1 = withResourceId(tx1, "abc123def456")
	tx2 := withAggregationKey(buildServiceTransaction("msg-2", 1, 2), "tx-2")
	tx2 = withResourceId(tx2, "abc123def456")
	repo.Transactions = []domain.Transaction{tx1, tx2}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		ResourceId: "abc123def456",
		Page:       1,
		Limit:      10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 2)
	assert.NotEmpty(t, result.Data[0].Id)
	assert.Equal(t, "1", result.Data[0].SourceChainId)
	assert.NotEmpty(t, result.Data[1].Id)
	assert.Equal(t, "1", result.Data[1].SourceChainId)
}

func TestTransactionService_GetTransactionsList_NoFiltersReturnsAll(t *testing.T) {
	// No filters returns all transactions
	repo := testutil.NewFakeTransactionRepository()
	repo.Transactions = []domain.Transaction{
		withAggregationKey(buildServiceTransaction("msg-1", 1, 2), "tx-1"),
		withAggregationKey(buildServiceTransaction("msg-2", 2, 3), "tx-2"),
		withAggregationKey(buildServiceTransaction("msg-3", 3, 1), "tx-3"),
	}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 3)
	assert.NotEmpty(t, result.Data[0].Id)
	assert.NotEmpty(t, result.Data[0].SourceChainId)
	assert.NotEmpty(t, result.Data[1].Id)
	assert.NotEmpty(t, result.Data[1].SourceChainId)
	assert.NotEmpty(t, result.Data[2].Id)
	assert.NotEmpty(t, result.Data[2].SourceChainId)
}

func TestTransactionService_GetRegularBatchTransactions_ReturnsPaginatedResults(t *testing.T) {
	// Batch transactions support pagination
	repo := testutil.NewFakeTransactionRepository()
	repo.BatchTransactions = []domain.Transaction{
		buildServiceTransaction("msg-1", 1, 2),
		buildServiceTransaction("msg-2", 1, 2),
		buildServiceTransaction("msg-3", 1, 2),
		buildServiceTransaction("msg-4", 1, 2),
		buildServiceTransaction("msg-5", 1, 2),
	}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetRegularBatchTransactions(context.Background(), "batch-abc", 2, 2)

	require.NoError(t, err)
	require.Len(t, result.Data, 5)
	assert.Equal(t, int64(5), result.Total)
}

func TestTransactionService_GetEnygmaBatchTransactions_NotFoundWhenEmpty(t *testing.T) {
	// Querying a non-existent Enygma batch returns not found error
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetEnygmaBatchTransactions(context.Background(), "nonexistent-enygma", 1, 10)

	require.Error(t, err)
	assert.Nil(t, result)
	var notFoundErr *NotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestTransactionService_GetTransactionsList_ReturnsMessageIdFilterResults(t *testing.T) {
	// MessageId filter returns matching transaction
	repo := testutil.NewFakeTransactionRepository()
	repo.Transactions = []domain.Transaction{
		withAggregationType(buildServiceTransaction("msg-specific", 1, 2), domain.AggregationTypeTransaction),
	}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		MessageId: "msg-specific",
		Page:      1,
		Limit:     10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, "msg-specific", result.Data[0].Id)
	assert.Equal(t, "1", result.Data[0].SourceChainId)
}

func TestTransactionService_GetTransactionsList_EmptyResultWhenNoMatch(t *testing.T) {
	// No matching transactions returns empty list
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		FromChainId: "999",
		Page:        1,
		Limit:       10,
	})

	require.NoError(t, err)
	assert.Empty(t, result.Data)
	assert.Equal(t, int64(0), result.Total)
}

func TestTransactionService_GetTransactionByTransactionId_ReturnsTransaction(t *testing.T) {
	// Querying a transaction by transactionId returns the transaction details
	repo := testutil.NewFakeTransactionRepository()
	repo.NotFoundErr = ErrRecordNotFound
	txId := uuid.New()
	tx := buildServiceTransaction("msg-456", 10, 20)
	tx.ID = txId
	tx = withResourceId(tx, "resource-xyz")
	tx = withServiceAddresses(
		tx,
		"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
	)
	repo.Transactions = []domain.Transaction{tx}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionByTransactionId(context.Background(), txId.String())

	require.NoError(t, err)
	assert.Equal(t, "msg-456", result.MessageId)
	assert.Equal(t, "resource-xyz", result.ResourceId)
	assert.Equal(t, "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", result.SourceAddress)
	assert.Equal(t, "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", result.DestinationAddress)
}

func TestTransactionService_GetTransactionByTransactionId_NotFound(t *testing.T) {
	// Querying a non-existent transactionId returns a NotFoundError
	repo := testutil.NewFakeTransactionRepository()
	repo.NotFoundErr = ErrRecordNotFound

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionByTransactionId(context.Background(), uuid.New().String())

	require.Error(t, err)
	assert.Nil(t, result)
	var notFoundErr *NotFoundError
	assert.True(t, errors.As(err, &notFoundErr), "expected NotFoundError, got %T", err)
}

func TestTransactionService_GetTransactionsBySharedId_ReturnsDvpSwapDtos(t *testing.T) {
	// Querying by sharedId returns DvP swap transaction DTOs
	repo := testutil.NewFakeTransactionRepository()
	tx1 := buildServiceTransaction("msg-dvp-1", 1, 2)
	tx1.SharedId = "shared-abc"
	tx1.ID = uuid.New()
	tx1.Token = domain.Token{Symbol: "SYM1"}
	tx2 := buildServiceTransaction("msg-dvp-2", 2, 1)
	tx2.SharedId = "shared-abc"
	tx2.ID = uuid.New()
	tx2.Token = domain.Token{Symbol: "SYM2"}
	repo.Transactions = []domain.Transaction{tx1, tx2}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsBySharedId(context.Background(), "shared-abc")

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, tx1.ID.String(), result[0].TransactionId)
	assert.Equal(t, "SYM1", result[0].TokenSymbol)
	assert.Equal(t, tx2.ID.String(), result[1].TransactionId)
	assert.Equal(t, "SYM2", result[1].TokenSymbol)
}

func TestTransactionService_GetTransactionsBySharedId_NotFoundWhenEmpty(t *testing.T) {
	// Querying a sharedId with no matching transactions returns a NotFoundError
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsBySharedId(context.Background(), "nonexistent-shared")

	require.Error(t, err)
	assert.Nil(t, result)
	var notFoundErr *NotFoundError
	assert.True(t, errors.As(err, &notFoundErr), "expected NotFoundError, got %T", err)
}

func TestTransactionService_GetTransactionsList_ReturnsDirectionalFilterResults(t *testing.T) {
	// Combined from/to address filters return matching transactions
	alice := "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	bob := "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

	repo := testutil.NewFakeTransactionRepository()
	tx := withServiceAddresses(buildServiceTransaction("msg-1", 1, 2), alice, bob)
	repo.Transactions = []domain.Transaction{tx}

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		From:  alice,
		To:    bob,
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, alice, result.Data[0].SourceAddress)
}

func TestTransactionService_GetTransactionsList_TransactionAggregationTypeReturnsMessageId(t *testing.T) {
	// Transaction aggregation type with message ID returns id_type as message_id
	repo := testutil.NewFakeTransactionRepository()
	tx := withAggregationType(buildServiceTransaction("msg-abc", 1, 2), domain.AggregationTypeTransaction)
	repo.Transactions = []domain.Transaction{tx}
	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, "msg-abc", result.Data[0].Id)
	assert.Equal(t, "message_id", result.Data[0].IdType)
}

func TestTransactionService_GetTransactionsList_DvpSwapAggregationTypeReturnsSharedId(t *testing.T) {
	// DvP Swap aggregation type returns id_type as shared_id
	repo := testutil.NewFakeTransactionRepository()
	tx := buildServiceTransaction("msg-dvp", 1, 2)
	tx.SharedId = "shared-123"
	tx.AggregationType = domain.AggregationTypeDvpSwap
	repo.Transactions = []domain.Transaction{tx}
	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, "shared-123", result.Data[0].Id)
	assert.Equal(t, "shared_id", result.Data[0].IdType)
}

func TestTransactionService_GetTransactionsList_EnygmaBatchAggregationTypeReturnsBatchId(t *testing.T) {
	// Enygma batch aggregation type with message ID returns batch_id
	repo := testutil.NewFakeTransactionRepository()
	tx := buildServiceTransaction("msg-enygma", 1, 2)
	tx = withAggregationKey(tx, "batch-key-enygma-123")
	tx.AggregationType = domain.AggregationTypeEnygmaBatch
	repo.Transactions = []domain.Transaction{tx}
	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, "batch-key-enygma-123", result.Data[0].Id)
	assert.Equal(t, "batch_id", result.Data[0].IdType)
}

func TestTransactionService_GetTransactionsList_RegularBatchAggregationTypeReturnsBatchId(t *testing.T) {
	// Regular batch aggregation type returns batch_id
	repo := testutil.NewFakeTransactionRepository()
	tx := buildServiceTransaction("msg-batch", 1, 2)
	tx = withAggregationKey(tx, "batch-key-regular-456")
	tx.AggregationType = domain.AggregationTypeRegularBatch
	repo.Transactions = []domain.Transaction{tx}
	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, "batch-key-regular-456", result.Data[0].Id)
	assert.Equal(t, "batch_id", result.Data[0].IdType)
}

func TestTransactionService_GetTransactionsList_TransactionWithoutMessageIdReturnsTransactionId(t *testing.T) {
	// Transaction aggregation type without message ID returns transaction_id
	repo := testutil.NewFakeTransactionRepository()
	tx := buildServiceTransaction("", 1, 2)
	tx.AggregationType = domain.AggregationTypeTransaction
	repo.Transactions = []domain.Transaction{tx}
	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, tx.ID.String(), result.Data[0].Id)
	assert.Equal(t, "transaction_id", result.Data[0].IdType)
}

func TestTransactionService_GetTransactionsList_EnygmaBatchWithoutMessageIdReturnsTransactionId(t *testing.T) {
	// Enygma batch without message ID returns transaction_id as fallback
	repo := testutil.NewFakeTransactionRepository()
	tx := buildServiceTransaction("", 1, 2)
	tx.AggregationType = domain.AggregationTypeEnygmaBatch
	tx = withAggregationKey(tx, "batch-key-abc")
	tx.MessageId = "" // Explicitly no message ID
	repo.Transactions = []domain.Transaction{tx}
	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, tx.ID.String(), result.Data[0].Id)
	assert.Equal(t, "transaction_id", result.Data[0].IdType)
}

func TestTransactionService_GetTransactionsList_UnknownAggregationTypeFallsBackToTransactionId(t *testing.T) {
	// Unknown aggregation type defaults to transaction_id
	repo := testutil.NewFakeTransactionRepository()
	tx := buildServiceTransaction("msg-unknown", 1, 2)
	tx.AggregationType = "unknown_type"
	repo.Transactions = []domain.Transaction{tx}
	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, tx.ID.String(), result.Data[0].Id)
	assert.Equal(t, "transaction_id", result.Data[0].IdType)
}

func TestTransactionService_GetTransactionByMessageId_EmptyMessageIdReturnsValidationError(t *testing.T) {
	// An empty messageId is rejected before any repository call
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionByMessageId(context.Background(), "")

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
}

func TestTransactionService_GetTransactionByMessageId_DatabaseErrorIsWrapped(t *testing.T) {
	// A database error from FindByMessageId is wrapped and returned as an InternalError
	repo := testutil.NewFakeTransactionRepository()
	repo.Error = errors.New("connection failed")

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionByMessageId(context.Background(), "msg-123")

	require.Error(t, err)
	assert.Nil(t, result)
	var internalErr *InternalError
	assert.True(t, errors.As(err, &internalErr))
}

func TestTransactionService_GetTransactionByTransactionId_EmptyTransactionIdReturnsValidationError(t *testing.T) {
	// An empty transactionId is rejected before any repository call
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionByTransactionId(context.Background(), "")

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
}

func TestTransactionService_GetTransactionByTransactionId_DatabaseErrorIsWrapped(t *testing.T) {
	// A database error from FindByTransactionId is wrapped and returned as an InternalError
	repo := testutil.NewFakeTransactionRepository()
	repo.Error = errors.New("connection failed")

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionByTransactionId(context.Background(), uuid.New().String())

	require.Error(t, err)
	assert.Nil(t, result)
	var internalErr *InternalError
	assert.True(t, errors.As(err, &internalErr))
}

func TestTransactionService_GetEnygmaBatchTransactions_EmptyBatchIdReturnsValidationError(t *testing.T) {
	// An empty batchId is rejected before any repository call
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetEnygmaBatchTransactions(context.Background(), "", 1, 10)

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
}

func TestTransactionService_GetRegularBatchTransactions_EmptyBatchIdReturnsValidationError(t *testing.T) {
	// An empty batchId is rejected before any repository call
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetRegularBatchTransactions(context.Background(), "", 1, 10)

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
}

func TestTransactionService_GetTransactionsBySharedId_EmptySharedIdReturnsValidationError(t *testing.T) {
	// An empty sharedId is rejected before any repository call
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsBySharedId(context.Background(), "")

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
}

func TestTransactionService_GetTransactionsList_InvalidFromChainIdReturnsValidationError(t *testing.T) {
	// A non-numeric sourceChainId is rejected with a validation error
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		FromChainId: "not-a-number",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Contains(t, err.Error(), "sourceChainId")
}

func TestTransactionService_GetTransactionsList_InvalidToChainIdReturnsValidationError(t *testing.T) {
	// A non-numeric destinationChainId is rejected with a validation error
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		ToChainId: "not-a-number",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Contains(t, err.Error(), "destinationChainId")
}

func TestTransactionService_GetTransactionsList_InvalidFromAddressReturnsValidationError(t *testing.T) {
	// A malformed from address (missing 0x prefix and wrong length) is rejected
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		From: "not-an-address",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Contains(t, err.Error(), "fromAddress")
}

func TestTransactionService_GetTransactionsList_InvalidToAddressReturnsValidationError(t *testing.T) {
	// A hex address shorter than 40 chars is rejected with a validation error
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		To: "0xbadaddr",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Contains(t, err.Error(), "toAddress")
}

func TestTransactionService_GetTransactionsList_InvalidResourceIdReturnsValidationError(t *testing.T) {
	// A non-hex resourceId is rejected with a validation error
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		ResourceId: "not-valid-hex!",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Contains(t, err.Error(), "resourceId")
}

func TestTransactionService_GetTransactionsList_InvalidMessageTypeReturnsValidationError(t *testing.T) {
	// An unrecognized messageType is rejected with a validation error
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		MessageType: "erc9999",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Contains(t, err.Error(), "messageType")
}

func TestTransactionService_GetTransactionsList_InvalidTimestampRangeReturnsValidationError(t *testing.T) {
	// initiatedAfter must be before initiatedBefore, otherwise a validation error is returned
	repo := testutil.NewFakeTransactionRepository()

	svc := NewTransactionService(repo, &testutil.StubTokenMetadataService{}, &testutil.StubLogger{})

	result, err := svc.GetTransactionsList(context.Background(), dto.MergedTransactionsFilters{
		InitiatedAfter:  "2026-12-31T00:00:00Z",
		InitiatedBefore: "2026-01-01T00:00:00Z",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
}
