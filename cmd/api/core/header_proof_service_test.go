package core

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/testutil"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
)

// Helpers

func buildHeaderProof(chainId, blockNumber int64, blockHash string) domain.HeaderProofEvent {
	return domain.HeaderProofEvent{
		ChainID:     domain.BigInt{Int: big.NewInt(chainId)},
		BlockNumber: domain.BigInt{Int: big.NewInt(blockNumber)},
		BlockHash:   blockHash,
		CreatedAt:   time.Now(),
	}
}

func TestHeaderProofService_GetHeaderProofsList_ReturnsProofsInBlockRange(t *testing.T) {
	// Querying header proofs within a block range returns matching proofs
	repo := testutil.NewFakeHeaderProofRepository()
	repo.HeaderProofs = []domain.HeaderProofEvent{
		buildHeaderProof(1, 20, "0xabc123"),
		buildHeaderProof(1, 30, "0xdef456"),
	}

	svc := NewHeaderProofService(repo, &testutil.StubLogger{})

	result, err := svc.GetHeaderProofsList(context.Background(), dto.HeaderProofFilters{
		ChainID:    "1",
		StartBlock: "15",
		EndBlock:   "35",
		Page:       1,
		PageSize:   10,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 2)
	assert.Equal(t, int64(2), result.Pagination.TotalRecords)
	assert.Equal(t, "20", result.Data[0].BlockNumber)
	assert.Equal(t, "30", result.Data[1].BlockNumber)
}

func TestHeaderProofService_GetHeaderProofsList_ReturnsEmptyWhenNoProofsInRange(t *testing.T) {
	// Querying a block range with no proofs returns empty list
	repo := testutil.NewFakeHeaderProofRepository()

	svc := NewHeaderProofService(repo, &testutil.StubLogger{})

	result, err := svc.GetHeaderProofsList(context.Background(), dto.HeaderProofFilters{
		ChainID:    "1",
		StartBlock: "1",
		EndBlock:   "50",
		Page:       1,
		PageSize:   10,
	})

	require.NoError(t, err)
	assert.Empty(t, result.Data)
	assert.Equal(t, int64(0), result.Pagination.TotalRecords)
}

func TestHeaderProofService_GetHeaderProofsList_ReturnsChainSpecificProofs(t *testing.T) {
	// Proofs are filtered by chain ID
	repo := testutil.NewFakeHeaderProofRepository()
	repo.HeaderProofs = []domain.HeaderProofEvent{
		buildHeaderProof(2, 15, "0x3"),
		buildHeaderProof(2, 25, "0x4"),
	}

	svc := NewHeaderProofService(repo, &testutil.StubLogger{})

	result, err := svc.GetHeaderProofsList(context.Background(), dto.HeaderProofFilters{
		ChainID:    "2",
		StartBlock: "1",
		EndBlock:   "100",
		Page:       1,
		PageSize:   10,
	})

	require.NoError(t, err)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, "2", result.Data[0].ChainID)
}

func TestHeaderProofService_GetHeaderProofsList_EndBlockLessThanStartBlockValidation(t *testing.T) {
	// endBlock must be greater than or equal to startBlock
	repo := testutil.NewFakeHeaderProofRepository()

	svc := NewHeaderProofService(repo, &testutil.StubLogger{})

	result, err := svc.GetHeaderProofsList(context.Background(), dto.HeaderProofFilters{
		ChainID:    "1",
		StartBlock: "100",
		EndBlock:   "50",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "endBlock")
}

func TestHeaderProofService_GetHeaderProofsList_Pagination(t *testing.T) {
	// Pagination metadata is calculated correctly
	repo := testutil.NewFakeHeaderProofRepository()
	repo.HeaderProofs = []domain.HeaderProofEvent{
		buildHeaderProof(1, 10, "0xabc"),
		buildHeaderProof(1, 20, "0xabc"),
		buildHeaderProof(1, 30, "0xabc"),
		buildHeaderProof(1, 40, "0xabc"),
		buildHeaderProof(1, 50, "0xabc"),
	}

	svc := NewHeaderProofService(repo, &testutil.StubLogger{})

	result, err := svc.GetHeaderProofsList(context.Background(), dto.HeaderProofFilters{
		ChainID:    "1",
		StartBlock: "1",
		EndBlock:   "100",
		Page:       2,
		PageSize:   2,
	})

	require.NoError(t, err)
	assert.Len(t, result.Data, 5)
	assert.Equal(t, int64(5), result.Pagination.TotalRecords)
	assert.Equal(t, 2, result.Pagination.CurrentPage)
	assert.Equal(t, 3, result.Pagination.TotalPages)
}

func TestHeaderProofService_GetHeaderProofsList_PageSizeExceedsMaxValidation(t *testing.T) {
	// pageSize cannot exceed maximum allowed value
	repo := testutil.NewFakeHeaderProofRepository()

	svc := NewHeaderProofService(repo, &testutil.StubLogger{})

	result, err := svc.GetHeaderProofsList(context.Background(), dto.HeaderProofFilters{
		ChainID:    "1",
		StartBlock: "1",
		EndBlock:   "100",
		PageSize:   1001,
	})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "pageSize")
}

func TestHeaderProofService_GetHeaderProofsList_NegativeStartBlockValidation(t *testing.T) {
	// startBlock cannot be negative
	repo := testutil.NewFakeHeaderProofRepository()

	svc := NewHeaderProofService(repo, &testutil.StubLogger{})

	result, err := svc.GetHeaderProofsList(context.Background(), dto.HeaderProofFilters{
		ChainID:    "1",
		StartBlock: "-1",
		EndBlock:   "100",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "startBlock")
}

func TestHeaderProofService_GetHeaderProofsList_TotalPagesCalculatesRemainderCorrectly(t *testing.T) {
	// Total pages calculation handles remainder correctly (10 records / 3 per page = 4 pages)
	repo := testutil.NewFakeHeaderProofRepository()
	repo.HeaderProofs = []domain.HeaderProofEvent{
		buildHeaderProof(1, 10, "0xabc"),
		buildHeaderProof(1, 20, "0xabc"),
		buildHeaderProof(1, 30, "0xabc"),
		buildHeaderProof(1, 40, "0xabc"),
		buildHeaderProof(1, 50, "0xabc"),
		buildHeaderProof(1, 60, "0xabc"),
		buildHeaderProof(1, 70, "0xabc"),
		buildHeaderProof(1, 80, "0xabc"),
		buildHeaderProof(1, 90, "0xabc"),
		buildHeaderProof(1, 100, "0xabc"),
	}

	svc := NewHeaderProofService(repo, &testutil.StubLogger{})

	result, err := svc.GetHeaderProofsList(context.Background(), dto.HeaderProofFilters{
		ChainID:    "1",
		StartBlock: "1",
		EndBlock:   "1000",
		Page:       1,
		PageSize:   3,
	})

	require.NoError(t, err)
	assert.Equal(t, int64(10), result.Pagination.TotalRecords)
	assert.Equal(t, 4, result.Pagination.TotalPages)
}

func TestHeaderProofService_GetHeaderProofsList_DefaultPaginationValuesAreApplied(t *testing.T) {
	// Default page=1 and pageSize=50 are applied when not specified
	repo := testutil.NewFakeHeaderProofRepository()
	repo.HeaderProofs = []domain.HeaderProofEvent{
		buildHeaderProof(1, 10, "0xabc"),
		buildHeaderProof(1, 20, "0xabc"),
		buildHeaderProof(1, 30, "0xabc"),
	}

	svc := NewHeaderProofService(repo, &testutil.StubLogger{})

	result, err := svc.GetHeaderProofsList(context.Background(), dto.HeaderProofFilters{
		ChainID:    "1",
		StartBlock: "1",
		EndBlock:   "100",
	})

	require.NoError(t, err)
	assert.Len(t, result.Data, 3)
	assert.Equal(t, 1, result.Pagination.CurrentPage)
	assert.Equal(t, 50, result.Pagination.PageSize)
}

func TestHeaderProofService_GetHeaderProofsList_NonNumericStartBlockReturnsValidationError(t *testing.T) {
	// A non-parseable startBlock returns a validation error
	repo := testutil.NewFakeHeaderProofRepository()

	svc := NewHeaderProofService(repo, &testutil.StubLogger{})

	result, err := svc.GetHeaderProofsList(context.Background(), dto.HeaderProofFilters{
		ChainID:    "1",
		StartBlock: "not-a-number",
		EndBlock:   "100",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "startBlock")
}

func TestHeaderProofService_GetHeaderProofsList_NonNumericEndBlockReturnsValidationError(t *testing.T) {
	// A non-parseable endBlock returns a validation error
	repo := testutil.NewFakeHeaderProofRepository()

	svc := NewHeaderProofService(repo, &testutil.StubLogger{})

	result, err := svc.GetHeaderProofsList(context.Background(), dto.HeaderProofFilters{
		ChainID:    "1",
		StartBlock: "1",
		EndBlock:   "not-a-number",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "endBlock")
}

func TestHeaderProofService_GetHeaderProofsList_EqualStartAndEndBlockIsAllowed(t *testing.T) {
	// startBlock == endBlock is a valid single-block range
	repo := testutil.NewFakeHeaderProofRepository()
	repo.HeaderProofs = []domain.HeaderProofEvent{
		buildHeaderProof(1, 50, "0xabc"),
	}

	svc := NewHeaderProofService(repo, &testutil.StubLogger{})

	result, err := svc.GetHeaderProofsList(context.Background(), dto.HeaderProofFilters{
		ChainID:    "1",
		StartBlock: "50",
		EndBlock:   "50",
		Page:       1,
		PageSize:   10,
	})

	require.NoError(t, err)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "50", result.Data[0].BlockNumber)
}
