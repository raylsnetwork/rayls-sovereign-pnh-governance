package core

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/testutil"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
)

// Helpers

func buildToken(resourceId, name, symbol string) domain.TokenWithBalancesAndFreezeState {
	return domain.TokenWithBalancesAndFreezeState{
		Token: domain.Token{
			ResourceId: resourceId,
			Name:       name,
			Symbol:     symbol,
		},
	}
}

func withTokenStatus(t domain.TokenWithBalancesAndFreezeState, status uint8) domain.TokenWithBalancesAndFreezeState {
	t.Status = status
	return t
}

func TestTokenService_GetTokenByResourceId_ReturnsTokenWithBalances(t *testing.T) {
	// Querying a token by resourceId returns the token with its balances
	repo := testutil.NewFakeTokenRepository()
	repo.Tokens = []domain.TokenWithBalancesAndFreezeState{
		buildToken("abcd", "Token A", "TKA"),
	}

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokenByResourceId(context.Background(), "abcd")

	require.NoError(t, err)
	assert.Equal(t, "abcd", result.ResourceId)
	assert.Equal(t, "Token A", result.Name)
	assert.Equal(t, "TKA", result.Symbol)
}

func TestTokenService_GetTokenByResourceId_NotFound(t *testing.T) {
	// Querying a non-existent token returns not found error
	repo := testutil.NewFakeTokenRepository()

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokenByResourceId(context.Background(), "abcd")

	require.Error(t, err)
	assert.Nil(t, result)
	var notFoundErr *NotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestTokenService_GetTokenRegistryStatus_ReturnsDto(t *testing.T) {
	// Querying token registry status returns a TokenRegistryStatusDto with the correct fields
	repo := testutil.NewFakeTokenRepository()
	repo.Tokens = []domain.TokenWithBalancesAndFreezeState{
		buildToken("abcd1234", "Coin", "ABC"),
	}

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokenRegistryStatus(context.Background(), "abcd1234")

	require.NoError(t, err)
	assert.Equal(t, "abcd1234", result.ResourceId)
	assert.Equal(t, "Coin", result.Name)
	assert.Equal(t, "ABC", result.Symbol)
}

func TestTokenService_GetTokenRegistryStatus_NotFound(t *testing.T) {
	// Querying a non-existent token returns not found error
	repo := testutil.NewFakeTokenRepository()

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokenRegistryStatus(context.Background(), "abcdef12")

	require.Error(t, err)
	assert.Nil(t, result)
	var notFoundErr *NotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestTokenService_GetTokenRegistryStatus_InvalidHex(t *testing.T) {
	// Invalid hex resourceId returns a validation error
	repo := testutil.NewFakeTokenRepository()

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokenRegistryStatus(context.Background(), "not-valid-hex!")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "must be a valid hex string")
}

func TestTokenService_GetTokenRegistryStatus_EmptyResourceId(t *testing.T) {
	// Empty resourceId returns a validation error
	repo := testutil.NewFakeTokenRepository()

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokenRegistryStatus(context.Background(), "")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestTokenService_GetTokensList_ReturnsAllTokens(t *testing.T) {
	// Listing tokens without filters returns all tokens with pagination metadata
	repo := testutil.NewFakeTokenRepository()
	repo.Tokens = []domain.TokenWithBalancesAndFreezeState{
		withTokenStatus(buildToken("a1", "Token A", ""), 1),
		withTokenStatus(buildToken("b2", "Token B", ""), 1),
	}

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokensList(context.Background(), dto.TokenListFilters{})

	require.NoError(t, err)
	require.Len(t, result.Data, 2)
	assert.Equal(t, "Token A", result.Data[0].Name)
	assert.Equal(t, "Token B", result.Data[1].Name)
	assert.Equal(t, int64(2), result.Total)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 1, result.Page)
}

func TestTokenService_GetTokensList_ReturnsTokensWithNameFilter(t *testing.T) {
	// Listing tokens with name filter returns matching tokens
	repo := testutil.NewFakeTokenRepository()
	repo.Tokens = []domain.TokenWithBalancesAndFreezeState{
		buildToken("a1", "Alpha", ""),
		buildToken("a2", "Alpine", ""),
		buildToken("b1", "Beta", ""),
	}

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokensList(context.Background(), dto.TokenListFilters{Name: "alp"})

	require.NoError(t, err)
	require.Len(t, result.Data, 2)
	assert.Equal(t, "Alpha", result.Data[0].Name)
	assert.Equal(t, "Alpine", result.Data[1].Name)
}

func TestTokenService_GetTokensList_ReturnsTokensWithStatusFilter(t *testing.T) {
	// Listing tokens with status filter returns matching tokens
	repo := testutil.NewFakeTokenRepository()
	repo.Tokens = []domain.TokenWithBalancesAndFreezeState{
		withTokenStatus(buildToken("x1", "Token X", ""), 1),
		withTokenStatus(buildToken("y1", "Token Y", ""), 0),
		withTokenStatus(buildToken("z1", "Token Z", ""), 2),
	}

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokensList(context.Background(), dto.TokenListFilters{Status: "active"})

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, "Token X", result.Data[0].Name)
	assert.Equal(t, "x1", result.Data[0].ResourceId)
}

func TestTokenService_GetTokensList_ReturnsEmptyWhenNoData(t *testing.T) {
	// Listing tokens when repository is empty returns empty list
	repo := testutil.NewFakeTokenRepository()

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokensList(context.Background(), dto.TokenListFilters{Name: "zzz"})

	require.NoError(t, err)
	assert.Empty(t, result.Data)
	assert.Equal(t, int64(0), result.Total)
}

func TestTokenService_GetTokensList_ReturnsTokensWithSpacesInName(t *testing.T) {
	// Tokens with spaces in name are matched correctly
	repo := testutil.NewFakeTokenRepository()
	repo.Tokens = []domain.TokenWithBalancesAndFreezeState{
		buildToken("b2", "Token One", ""),
		buildToken("c3", "Token Two", ""),
	}

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokensList(context.Background(), dto.TokenListFilters{Name: " "})

	require.NoError(t, err)
	require.Len(t, result.Data, 2)
	assert.Equal(t, "Token One", result.Data[0].Name)
	assert.Equal(t, "Token Two", result.Data[1].Name)
}

func TestTokenService_GetTokenByResourceId_EmptyResourceIdReturnsValidationError(t *testing.T) {
	// Empty resourceId is rejected before any database call
	repo := testutil.NewFakeTokenRepository()

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokenByResourceId(context.Background(), "")

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestTokenService_GetTokenByResourceId_InvalidHexReturnsValidationError(t *testing.T) {
	// A resourceId containing non-hex characters is rejected
	repo := testutil.NewFakeTokenRepository()

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokenByResourceId(context.Background(), "not-valid-hex!")

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Contains(t, err.Error(), "hex")
}

func TestTokenService_GetTokensList_ReturnsTokensWithSymbolFilter(t *testing.T) {
	// Listing tokens filtered by symbol prefix returns only matching tokens
	repo := testutil.NewFakeTokenRepository()
	repo.Tokens = []domain.TokenWithBalancesAndFreezeState{
		buildToken("a1", "Token A", "TKA"),
		buildToken("b1", "Token B", "TKB"),
		buildToken("c1", "Ether", "ETH"),
	}

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokensList(context.Background(), dto.TokenListFilters{Symbol: "TK"})

	require.NoError(t, err)
	require.Len(t, result.Data, 2)
	assert.Equal(t, "TKA", result.Data[0].Symbol)
	assert.Equal(t, "TKB", result.Data[1].Symbol)
}

func TestTokenService_GetTokensList_ReturnsTokensWithIssuerIdFilter(t *testing.T) {
	// Listing tokens filtered by issuerId returns only tokens issued by that chain
	repo := testutil.NewFakeTokenRepository()
	t1 := buildToken("a1", "Token A", "TKA")
	t1.IssuerId = "issuer-abc"
	t2 := buildToken("b1", "Token B", "TKB")
	t2.IssuerId = "issuer-xyz"
	repo.Tokens = []domain.TokenWithBalancesAndFreezeState{t1, t2}

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokensList(context.Background(), dto.TokenListFilters{IssuerId: "issuer-abc"})

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, "Token A", result.Data[0].Name)
	assert.Equal(t, "issuer-abc", result.Data[0].IssuerId)
}

func TestTokenService_GetTokensList_InvalidDateRangeReturnsValidationError(t *testing.T) {
	// createdAfter must be before createdBefore, otherwise a validation error is returned
	repo := testutil.NewFakeTokenRepository()

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokensList(context.Background(), dto.TokenListFilters{
		CreatedAfter:  "2025-12-31T00:00:00Z",
		CreatedBefore: "2025-01-01T00:00:00Z",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
}

func TestTokenService_GetTokensList_PaginationNormalizesOutOfRangeValues(t *testing.T) {
	// Limit > 100 is normalized to 10, and page < 1 is normalized to 1
	repo := testutil.NewFakeTokenRepository()
	repo.Tokens = []domain.TokenWithBalancesAndFreezeState{
		buildToken("a1", "Token A", "TKA"),
	}

	svc := NewTokenService(repo, &testutil.StubLogger{})

	result, err := svc.GetTokensList(context.Background(), dto.TokenListFilters{
		Limit: 9999,
		Page:  0,
	})

	require.NoError(t, err)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 1, result.Page)
}
