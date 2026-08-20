package core

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/api/testutil"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
)

// Helpers
func buildBalance(chainId, resourceId string, amount int64) domain.Balance {
	return domain.Balance{
		ChainId:    chainId,
		ResourceId: resourceId,
		Amount:     decimal.NewFromInt(amount),
	}
}

func TestBalanceService_GetBalancesInChain_ReturnsAllBalances(t *testing.T) {
	// Querying all balances in chain 1 returns all resources with their amounts
	repo := testutil.NewFakeBalanceRepository()
	repo.Balances = []domain.Balance{
		buildBalance("1", "resource-a", 100),
		buildBalance("1", "resource-b", 200),
		buildBalance("1", "resource-c", 300),
	}

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalancesInChain(context.Background(), "1", "/")

	require.NoError(t, err)
	balances := result.([]domain.Balance)
	require.Len(t, balances, 3)
	assert.Equal(t, "resource-a", balances[0].ResourceId)
	assert.Equal(t, "100", balances[0].Amount.String())
	assert.Equal(t, "resource-b", balances[1].ResourceId)
	assert.Equal(t, "200", balances[1].Amount.String())
	assert.Equal(t, "resource-c", balances[2].ResourceId)
	assert.Equal(t, "300", balances[2].Amount.String())
}

func TestBalanceService_GetBalancesInChain_ReturnsSpecificResource(t *testing.T) {
	// Querying a specific resource returns only that balance
	repo := testutil.NewFakeBalanceRepository()
	repo.Balances = []domain.Balance{
		buildBalance("1", "resource-b", 200),
	}

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalancesInChain(context.Background(), "1", "resource-b")

	require.NoError(t, err)
	balance := result.(*domain.Balance)
	assert.Equal(t, "resource-b", balance.ResourceId)
	assert.Equal(t, "200", balance.Amount.String())
}

func TestBalanceService_GetBalancesInChain_NotFoundWhenChainEmpty(t *testing.T) {
	// Querying a chain with no balances returns not found error
	repo := testutil.NewFakeBalanceRepository()

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalancesInChain(context.Background(), "1", "/")

	require.Error(t, err)
	assert.Nil(t, result)
	var notFoundErr *NotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestBalanceService_GetBalancesInChain_NotFoundWhenResourceMissing(t *testing.T) {
	// Querying a non-existent resource returns not found error
	repo := testutil.NewFakeBalanceRepository()
	repo.NotFoundErr = ErrRecordNotFound
	repo.Balances = []domain.Balance{
		buildBalance("1", "resource-a", 100),
	}

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalancesInChain(context.Background(), "1", "resource-x")

	require.Error(t, err)
	assert.Nil(t, result)
	var notFoundErr *NotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestBalanceService_GetBalanceAcrossAllChains_ReturnsFromAllChains(t *testing.T) {
	// Querying a token across all chains returns balances from each chain
	repo := testutil.NewFakeBalanceRepository()
	repo.Balances = []domain.Balance{
		buildBalance("1", "token-a", 100),
		buildBalance("2", "token-a", 200),
		buildBalance("3", "token-a", 300),
	}

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalanceAcrossAllChains(context.Background(), "token-a")

	require.NoError(t, err)
	require.Len(t, result, 3)
	assert.Equal(t, "1", result[0].ChainId)
	assert.Equal(t, "100", result[0].Amount.String())
	assert.Equal(t, "2", result[1].ChainId)
	assert.Equal(t, "200", result[1].Amount.String())
	assert.Equal(t, "3", result[2].ChainId)
	assert.Equal(t, "300", result[2].Amount.String())
}

func TestBalanceService_GetBalanceAcrossAllChains_NotFoundWhenResourceMissing(t *testing.T) {
	// Querying a non-existent resource across chains returns not found error
	repo := testutil.NewFakeBalanceRepository()

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalanceAcrossAllChains(context.Background(), "token-x")

	require.Error(t, err)
	assert.Nil(t, result)
	var notFoundErr *NotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestBalanceService_GetBalanceAcrossSpecificChains_ReturnsFromSpecifiedChains(t *testing.T) {
	// Querying specific chains returns only balances from those chains
	repo := testutil.NewFakeBalanceRepository()
	repo.Balances = []domain.Balance{
		buildBalance("1", "token-a", 100),
		buildBalance("2", "token-a", 200),
	}

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalanceAcrossSpecificChains(context.Background(), "token-a", []string{"1", "2"})

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "1", result[0].ChainId)
	assert.Equal(t, "100", result[0].Amount.String())
	assert.Equal(t, "2", result[1].ChainId)
	assert.Equal(t, "200", result[1].Amount.String())
}

func TestBalanceService_GetBalanceAcrossSpecificChains_EmptyChainsReturnsAll(t *testing.T) {
	// Passing empty chain list returns balances from all chains
	repo := testutil.NewFakeBalanceRepository()
	repo.Balances = []domain.Balance{
		buildBalance("1", "token-a", 100),
		buildBalance("2", "token-a", 200),
		buildBalance("3", "token-a", 300),
	}

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalanceAcrossSpecificChains(context.Background(), "token-a", []string{})

	require.NoError(t, err)
	require.Len(t, result, 3)
	assert.Equal(t, "1", result[0].ChainId)
	assert.Equal(t, "100", result[0].Amount.String())
	assert.Equal(t, "2", result[1].ChainId)
	assert.Equal(t, "200", result[1].Amount.String())
	assert.Equal(t, "3", result[2].ChainId)
	assert.Equal(t, "300", result[2].Amount.String())
}

func TestBalanceService_GetBalancesInChain_StripsLeadingSlashFromResourceId(t *testing.T) {
	// Resource IDs with leading slash are normalized before lookup
	repo := testutil.NewFakeBalanceRepository()
	repo.Balances = []domain.Balance{
		buildBalance("1", "abc123", 500),
	}

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalancesInChain(context.Background(), "1", "/abc123")

	require.NoError(t, err)
	balance := result.(*domain.Balance)
	assert.Equal(t, "abc123", balance.ResourceId)
}

func TestBalanceService_GetBalanceAcrossSpecificChains_NotFoundInSpecificChains(t *testing.T) {
	// Querying chains that don't have the resource returns not found error
	repo := testutil.NewFakeBalanceRepository()

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalanceAcrossSpecificChains(context.Background(), "token-a", []string{"3", "4"})

	require.Error(t, err)
	assert.Nil(t, result)
	var notFoundErr *NotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestBalanceService_GetBalancesInChain_EmptyResourceIdReturnsAll(t *testing.T) {
	// Empty resource ID returns all balances in the chain
	repo := testutil.NewFakeBalanceRepository()
	repo.Balances = []domain.Balance{
		buildBalance("1", "resource-a", 100),
		buildBalance("1", "resource-b", 200),
	}

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalancesInChain(context.Background(), "1", "")

	require.NoError(t, err)
	balances := result.([]domain.Balance)
	require.Len(t, balances, 2)
	assert.Equal(t, "100", balances[0].Amount.String())
	assert.Equal(t, "200", balances[1].Amount.String())
}

func TestBalanceService_GetBalancesInChain_ReturnsOnlyRequestedChain(t *testing.T) {
	// Balances are scoped to the requested chain - other chains are excluded
	repo := testutil.NewFakeBalanceRepository()
	repo.Balances = []domain.Balance{
		buildBalance("1", "token-a", 100),
		buildBalance("2", "token-b", 200),
		buildBalance("3", "token-c", 300),
	}

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalancesInChain(context.Background(), "1", "/")

	require.NoError(t, err)
	balances := result.([]domain.Balance)
	require.Len(t, balances, 1)
	assert.Equal(t, "1", balances[0].ChainId)
	assert.Equal(t, "token-a", balances[0].ResourceId)
}

func TestBalanceService_GetBalanceAcrossAllChains_ReturnsIndependentBalancesPerChain(t *testing.T) {
	// Each chain has an independent balance for the same token
	repo := testutil.NewFakeBalanceRepository()
	repo.Balances = []domain.Balance{
		buildBalance("1", "token-a", 1000),
		buildBalance("2", "token-a", 500),
	}

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalanceAcrossAllChains(context.Background(), "token-a")

	require.NoError(t, err)
	assert.Len(t, result, 2)
	totalAmount := result[0].Amount.Add(result[1].Amount)
	assert.Equal(t, "1500", totalAmount.String())
}

func TestBalanceService_GetBalancesInChain_EmptyChainIdCannotQueryBalances(t *testing.T) {
	// Empty chain ID returns validation error
	repo := testutil.NewFakeBalanceRepository()

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalancesInChain(context.Background(), "", "/")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Chain ID")
}

func TestBalanceService_GetBalanceAcrossAllChains_EmptyResourceIdCannotQuery(t *testing.T) {
	// Empty resource ID returns validation error
	repo := testutil.NewFakeBalanceRepository()

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalanceAcrossAllChains(context.Background(), "")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Resource ID")
}

func TestBalanceService_GetBalanceAcrossSpecificChains_EmptyResourceIdCannotQuery(t *testing.T) {
	// Empty resource ID returns validation error
	repo := testutil.NewFakeBalanceRepository()

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalanceAcrossSpecificChains(context.Background(), "", []string{"1", "2"})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Resource ID")
}

func TestBalanceService_GetBalancesInChain_DatabaseErrorWhenFetchingSpecificResource(t *testing.T) {
	// A non-not-found database error from FindByChainAndResource is returned as an InternalError,
	// not confused with a NotFoundError
	repo := testutil.NewFakeBalanceRepository()
	repo.Error = errors.New("connection refused")

	service := NewBalanceService(repo, &testutil.StubLogger{})

	result, err := service.GetBalancesInChain(context.Background(), "1", "resource-abc")

	require.Error(t, err)
	assert.Nil(t, result)
	var internalErr *InternalError
	assert.True(t, errors.As(err, &internalErr))
}
