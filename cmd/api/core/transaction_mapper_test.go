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

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/testutil"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

// Helpers

const testMetadataUrl = "https://example.com/metadata/"

func buildTransaction(msgType types.AssetType) domain.Transaction {
	return domain.Transaction{
		MessageId:   "msg-" + types.AssetTypeToString[uint8(msgType)],
		MsgType:     uint8(msgType),
		FromChainId: domain.BigInt{Int: big.NewInt(1)},
		ToChainId:   domain.BigInt{Int: big.NewInt(2)},
		Amount:      decimal.NewFromInt(100),
	}
}

func buildEnygmaTransaction(ccTimestamp time.Time) domain.Transaction {
	tx := buildTransaction(types.AssetTypeEnygma)
	tx.HubTxHash = "0xcc123abc"
	tx.HubTimestamp = ccTimestamp
	tx.BlockNumber = decimal.NewFromInt(12345)
	return tx
}

func buildNFTTransaction(msgType types.AssetType, ercId int64, metadataUrl string) domain.Transaction {
	tx := buildTransaction(msgType)
	tx.ErcId = domain.BigInt{Int: big.NewInt(ercId)}
	tx.Token = domain.Token{
		MetadataUrl: metadataUrl,
	}
	return tx
}

func buildSimpleToken(name, symbol string) domain.Token {
	return domain.Token{
		Name:   name,
		Symbol: symbol,
	}
}

func withTimestamps(tx domain.Transaction, source, created time.Time) domain.Transaction {
	tx.SourceTimestamp = source
	tx.CreatedAt = created
	return tx
}

func withStatus(tx domain.Transaction, status types.AtomicTeleportStatus) domain.Transaction {
	s := uint8(status)
	tx.TeleportStatus = &s
	return tx
}

func withAddresses(tx domain.Transaction, from, to string) domain.Transaction {
	tx.From = from
	tx.To = to
	return tx
}

func TestTransactionMapper_ToTransactionDetailDto_EnygmaTransactionPopulatesCcFields(t *testing.T) {
	// An Enygma transaction with EnygmaTransaction data should populate Hub-specific fields
	mapper := NewTransactionMapper(&testutil.StubTokenMetadataService{})

	ccTimestamp := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	tx := buildEnygmaTransaction(ccTimestamp)

	result, err := mapper.ToTransactionDetailDto(context.Background(), tx, "enygma")

	require.NoError(t, err)
	assert.Equal(t, "0xcc123abc", result.HubHash)
	assert.NotEmpty(t, result.HubTimestamp)
	assert.Equal(t, "12345", result.HubBlockNumber)
	assert.Equal(t, "enygma", result.MessageType)
}

func TestTransactionMapper_ToTransactionDetailDto_EnygmaTransactionOmitsSourceAndDestinationTimestamps(t *testing.T) {
	// Enygma transactions do not carry source/destination timestamps — they must be omitted from the response
	mapper := NewTransactionMapper(&testutil.StubTokenMetadataService{})

	ccTimestamp := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	tx := buildEnygmaTransaction(ccTimestamp)
	tx.SourceTimestamp = time.Date(2024, 6, 15, 10, 25, 0, 0, time.UTC)
	tx.DestinationTimestamp = time.Date(2024, 6, 15, 10, 28, 0, 0, time.UTC)

	result, err := mapper.ToTransactionDetailDto(context.Background(), tx, "enygma")

	require.NoError(t, err)
	assert.Empty(t, result.SourceTimestamp)
	assert.Empty(t, result.DestinationTimestamp)
}

func TestTransactionMapper_ToTransactionDetailDto_NonEnygmaTransactionPopulatesStatusFields(t *testing.T) {
	// A regular ERC20 transaction should populate status and timestamp fields instead of Hub fields
	mapper := NewTransactionMapper(&testutil.StubTokenMetadataService{})

	sourceTimestamp := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	createdAt := time.Date(2024, 6, 15, 10, 35, 0, 0, time.UTC)
	tx := buildTransaction(types.AssetTypeERC20)
	tx = withStatus(tx, types.AtomicTeleportExecuted)
	tx = withTimestamps(tx, sourceTimestamp, createdAt)

	result, err := mapper.ToTransactionDetailDto(context.Background(), tx, "transaction")

	require.NoError(t, err)
	assert.Equal(t, "executed", result.Status)
	assert.NotEmpty(t, result.SourceTimestamp) // Unix timestamp string (e.g. "1718447400")
	assert.Contains(t, result.CreatedAt, "2024-06-15")
	// Hub fields should be empty for non-Enygma
	assert.Empty(t, result.HubHash)
	assert.Empty(t, result.HubTimestamp)
}

func TestTransactionMapper_ToTransactionDetailDto_ERC721TransactionFetchesNFTMetadata(t *testing.T) {
	// An ERC721 NFT transaction with token metadata URL should fetch and populate metadata
	metadataService := &testutil.StubTokenMetadataService{
		Metadata: &dto.TokenMetadataInfoDto{
			Name:        "CryptoPunk #1234",
			Description: "A rare punk",
			ImageUrl:    "https://example.com/punk1234.png",
		},
	}
	mapper := NewTransactionMapper(metadataService)

	tx := buildNFTTransaction(types.AssetTypeERC721, 1234, testMetadataUrl)

	result, err := mapper.ToTransactionDetailDto(context.Background(), tx, "transaction")

	require.NoError(t, err)
	assert.Equal(t, "1234", result.ErcId)
	assert.Equal(t, "CryptoPunk #1234", result.TokenMetadataName)
	assert.Equal(t, "A rare punk", result.TokenMetadataDescription)
	assert.Equal(t, "https://example.com/punk1234.png", result.TokenMetadataImageUrl)
}

func TestTransactionMapper_ToTransactionDetailDto_ERC1155TransactionFetchesNFTMetadata(t *testing.T) {
	// An ERC1155 multi-token transaction should also fetch NFT metadata
	metadataService := &testutil.StubTokenMetadataService{
		Metadata: &dto.TokenMetadataInfoDto{
			Name:     "Game Item #42",
			ImageUrl: "https://example.com/item42.png",
		},
	}
	mapper := NewTransactionMapper(metadataService)

	tx := buildNFTTransaction(types.AssetTypeERC1155, 42, testMetadataUrl)

	result, err := mapper.ToTransactionDetailDto(context.Background(), tx, "transaction")

	require.NoError(t, err)
	assert.Equal(t, "42", result.ErcId)
	assert.Equal(t, "Game Item #42", result.TokenMetadataName)
}

func TestTransactionMapper_ToTransactionDetailDto_NFTWithoutMetadataUrlSkipsMetadataFetch(t *testing.T) {
	// An NFT transaction without metadata URL should not attempt to fetch metadata
	metadataService := &testutil.StubTokenMetadataService{
		Metadata: &dto.TokenMetadataInfoDto{Name: "Should not be fetched"},
	}
	mapper := NewTransactionMapper(metadataService)

	tx := buildNFTTransaction(types.AssetTypeERC721, 999, "") // No metadata URL

	result, err := mapper.ToTransactionDetailDto(context.Background(), tx, "transaction")

	require.NoError(t, err)
	assert.Equal(t, "999", result.ErcId)
	// Metadata fields should be empty since no URL was provided
	assert.Empty(t, result.TokenMetadataName)
	assert.Empty(t, result.TokenMetadataImageUrl)
}

func TestTransactionMapper_ToTransactionDetailDto_MetadataFetchErrorLeavesFieldsEmpty(t *testing.T) {
	// When metadata fetch fails, the transaction should still succeed with empty metadata fields
	metadataService := &testutil.StubTokenMetadataService{
		ShouldError: true,
	}
	mapper := NewTransactionMapper(metadataService)

	tx := buildNFTTransaction(types.AssetTypeERC721, 555, testMetadataUrl)

	result, err := mapper.ToTransactionDetailDto(context.Background(), tx, "transaction")

	// Should succeed even when metadata fetch fails
	require.NoError(t, err)
	assert.Equal(t, "555", result.ErcId)
	assert.Empty(t, result.TokenMetadataName)
}

func TestTransactionMapper_ToTransactionDetailDto_ERC20DoesNotFetchMetadata(t *testing.T) {
	// An ERC20 fungible token should not attempt to fetch NFT metadata even with ErcId
	metadataService := &testutil.StubTokenMetadataService{
		Metadata: &dto.TokenMetadataInfoDto{Name: "Should not be fetched"},
	}
	mapper := NewTransactionMapper(metadataService)

	tx := buildTransaction(types.AssetTypeERC20)
	tx.Token.MetadataUrl = testMetadataUrl

	result, err := mapper.ToTransactionDetailDto(context.Background(), tx, "transaction")

	require.NoError(t, err)
	// Metadata fields should be empty for ERC20
	assert.Empty(t, result.TokenMetadataName)
	assert.Empty(t, result.ErcId)
}

func TestTransactionMapper_ToBatchTransactionDtos_ConvertsMutipleTransactions(t *testing.T) {
	// A batch of 3 transactions should all be converted to DTOs with correct field mapping
	mapper := NewTransactionMapper(&testutil.StubTokenMetadataService{})

	tx1 := buildTransaction(types.AssetTypeERC20)
	tx1.MessageId = "msg-1"
	tx1 = withAddresses(tx1, "0x1111", "0x2222")
	tx1.Token = buildSimpleToken("", "TKN1")

	tx2 := buildTransaction(types.AssetTypeERC20)
	tx2.MessageId = "msg-2"
	tx2.FromChainId = domain.BigInt{Int: big.NewInt(3)}
	tx2.ToChainId = domain.BigInt{Int: big.NewInt(4)}
	tx2 = withAddresses(tx2, "0x3333", "0x4444")
	tx2.Amount = decimal.NewFromInt(200)
	tx2.Token = buildSimpleToken("", "TKN2")

	tx3 := buildTransaction(types.AssetTypeERC20)
	tx3.MessageId = "msg-3"
	tx3.FromChainId = domain.BigInt{Int: big.NewInt(5)}
	tx3.ToChainId = domain.BigInt{Int: big.NewInt(6)}
	tx3 = withAddresses(tx3, "0x5555", "0x6666")
	tx3.Amount = decimal.NewFromInt(300)
	tx3.Token = buildSimpleToken("", "TKN3")

	transactions := []domain.Transaction{tx1, tx2, tx3}

	result := mapper.ToBatchTransactionDtos(transactions)

	require.Len(t, result, 3)

	assert.Equal(t, "msg-1", result[0].MessageId)
	assert.Equal(t, "1", result[0].SourceChainId)
	assert.Equal(t, "2", result[0].DestinationChainId)
	assert.Equal(t, "TKN1", result[0].TokenSymbol)

	assert.Equal(t, "msg-2", result[1].MessageId)
	assert.Equal(t, "3", result[1].SourceChainId)

	assert.Equal(t, "msg-3", result[2].MessageId)
	assert.Equal(t, "5", result[2].SourceChainId)
}

func TestTransactionMapper_ToTransactionDetailDto_ZeroDestinationTimestampIsOmitted(t *testing.T) {
	// A transaction with zero destination timestamp should omit the field entirely (empty string, omitempty)
	mapper := NewTransactionMapper(&testutil.StubTokenMetadataService{})

	tx := buildTransaction(types.AssetTypeERC20)
	tx.DestinationTimestamp = time.Time{} // Zero time

	result, err := mapper.ToTransactionDetailDto(context.Background(), tx, "transaction")

	require.NoError(t, err)
	assert.Empty(t, result.DestinationTimestamp)
}

func TestTransactionMapper_ToTransactionDetailDto_PreservesAllCoreFields(t *testing.T) {
	// All core transaction fields should be mapped correctly to the DTO
	mapper := NewTransactionMapper(&testutil.StubTokenMetadataService{})

	destTimestamp := time.Date(2024, 6, 15, 14, 0, 0, 0, time.UTC)
	tx := buildTransaction(types.AssetTypeERC20)
	tx.MessageId = "msg-full"
	tx.TxHashSource = "0xsourcehash"
	tx.TxHashDestination = "0xdesthash"
	tx.FromChainId = domain.BigInt{Int: big.NewInt(100)}
	tx.ToChainId = domain.BigInt{Int: big.NewInt(200)}
	tx.DestinationTimestamp = destTimestamp
	tx = withAddresses(tx, "0xfromaddr", "0xtoaddr")
	tx.Amount = decimal.NewFromFloat(123.456)
	tx.ResourceId = "resource-abc"
	tx.Payload = "0xpayloaddata"
	tx.Token = buildSimpleToken("Test Token", "TST")

	result, err := mapper.ToTransactionDetailDto(context.Background(), tx, "regular")

	require.NoError(t, err)
	assert.Equal(t, "msg-full", result.MessageId)
	assert.Equal(t, "0xsourcehash", result.SourceTransactionHash)
	assert.Equal(t, "0xdesthash", result.DestinationTransactionHash)
	assert.Equal(t, "100", result.SourceChainId)
	assert.Equal(t, "200", result.DestinationChainId)
	assert.Equal(t, "0xfromaddr", result.SourceAddress)
	assert.Equal(t, "0xtoaddr", result.DestinationAddress)
	assert.Equal(t, "resource-abc", result.ResourceId)
	assert.Equal(t, "erc20", result.MessageType)
	assert.Equal(t, "0xpayloaddata", result.Payload)
	assert.Equal(t, "regular", result.Type)
	assert.Equal(t, "Test Token", result.TokenName)
	assert.Equal(t, "TST", result.TokenSymbol)
}

func TestTransactionMapper_ToDvpSwapTransactionDtos_MapsFieldsCorrectly(t *testing.T) {
	// DvP swap DTO maps transactionId, addresses, amounts, and token symbols from domain transactions
	mapper := NewTransactionMapper(&testutil.StubTokenMetadataService{})

	id1 := uuid.New()
	id2 := uuid.New()

	tx1 := buildTransaction(types.AssetTypeERC20)
	tx1.ID = id1
	tx1 = withAddresses(tx1, "0xaaaa", "0xbbbb")
	tx1.Amount = decimal.NewFromInt(500)
	tx1.Token = buildSimpleToken("", "SYM1")

	tx2 := buildTransaction(types.AssetTypeERC20)
	tx2.ID = id2
	tx2.FromChainId = domain.BigInt{Int: big.NewInt(3)}
	tx2.ToChainId = domain.BigInt{Int: big.NewInt(4)}
	tx2 = withAddresses(tx2, "0xcccc", "0xdddd")
	tx2.Amount = decimal.NewFromInt(1000)
	tx2.Token = buildSimpleToken("", "SYM2")

	result := mapper.ToDvpSwapTransactionDtos([]domain.Transaction{tx1, tx2})

	require.Len(t, result, 2)

	assert.Equal(t, id1.String(), result[0].TransactionId)
	assert.Equal(t, "0xaaaa", result[0].SourceAddress)
	assert.Equal(t, "0xbbbb", result[0].DestinationAddress)
	assert.Equal(t, "1", result[0].SourceChainId)
	assert.Equal(t, "2", result[0].DestinationChainId)
	assert.Equal(t, "SYM1", result[0].TokenSymbol)

	assert.Equal(t, id2.String(), result[1].TransactionId)
	assert.Equal(t, "3", result[1].SourceChainId)
	assert.Equal(t, "4", result[1].DestinationChainId)
	assert.Equal(t, "SYM2", result[1].TokenSymbol)
}
