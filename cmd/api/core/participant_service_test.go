package core

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/api/testutil"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/dto"
)

// Helpers

func buildParticipant(chainId uint, name string) domain.Participant {
	id := chainId
	return domain.Participant{
		ChainId: &id,
		Name:    name,
	}
}

func withParticipantStatus(p domain.Participant, status string) domain.Participant {
	p.StatusStr = status
	return p
}

func withParticipantRole(p domain.Participant, role string) domain.Participant {
	p.RoleStr = role
	return p
}

func TestParticipantService_GetParticipant_ReturnsParticipantData(t *testing.T) {
	// Querying a participant by chain ID returns the participant data
	repo := testutil.NewFakeParticipantRepository()
	p := buildParticipant(123, "Participant A")
	p = withParticipantStatus(p, "active")
	p = withParticipantRole(p, "issuer")
	repo.Participants = []domain.Participant{p}

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantByChainId(context.Background(), "123")

	require.NoError(t, err)
	assert.Equal(t, "Participant A", result.Name)
	assert.Equal(t, "active", result.StatusStr)
	assert.Equal(t, "issuer", result.RoleStr)
}

func TestParticipantService_GetParticipant_NotFound(t *testing.T) {
	// Querying a non-existent participant returns not found error
	repo := testutil.NewFakeParticipantRepository()
	repo.NotFoundErr = ErrRecordNotFound

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantByChainId(context.Background(), "999")

	require.Error(t, err)
	assert.Nil(t, result)
	var notFoundErr *NotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestParticipantService_GetParticipantsList_ReturnsParticipants(t *testing.T) {
	// Listing participants with status filter returns matching participants
	repo := testutil.NewFakeParticipantRepository()
	repo.Participants = []domain.Participant{
		withParticipantStatus(buildParticipant(1, "Participant A"), "active"),
		withParticipantStatus(buildParticipant(2, "Participant B"), "active"),
	}

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantsList(context.Background(), dto.ParticipantListFilters{Status: "active"})

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "Participant A", result[0].Name)
	assert.Equal(t, "Participant B", result[1].Name)
}

func TestParticipantService_GetParticipantsList_ReturnsParticipantsWithRoleFilter(t *testing.T) {
	// Listing participants with role filter returns matching participants
	repo := testutil.NewFakeParticipantRepository()
	repo.Participants = []domain.Participant{
		withParticipantRole(buildParticipant(1, "Participant A"), "issuer"),
		withParticipantRole(buildParticipant(2, "Participant B"), "issuer"),
	}

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantsList(context.Background(), dto.ParticipantListFilters{Role: "issuer"})

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "Participant A", result[0].Name)
	assert.Equal(t, "Participant B", result[1].Name)
}

func TestParticipantService_GetParticipantsList_ReturnsParticipantsWithCombinedFilters(t *testing.T) {
	// Listing participants with both status and role filters returns matching participants
	repo := testutil.NewFakeParticipantRepository()
	p1 := withParticipantRole(withParticipantStatus(buildParticipant(1, "Participant A"), "active"), "issuer")
	p2 := withParticipantRole(withParticipantStatus(buildParticipant(4, "Participant D"), "active"), "issuer")
	repo.Participants = []domain.Participant{p1, p2}

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantsList(context.Background(), dto.ParticipantListFilters{
		Status: "active",
		Role:   "issuer",
	})

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "Participant A", result[0].Name)
	assert.Equal(t, "Participant D", result[1].Name)
}

func TestParticipantService_GetParticipantsList_ReturnsAllWhenNoFilters(t *testing.T) {
	// Listing participants without filters returns all participants
	repo := testutil.NewFakeParticipantRepository()
	repo.Participants = []domain.Participant{
		buildParticipant(1, "Participant A"),
		buildParticipant(2, "Participant B"),
		buildParticipant(3, "Participant C"),
	}

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantsList(context.Background(), dto.ParticipantListFilters{})

	require.NoError(t, err)
	require.Len(t, result, 3)
	assert.Equal(t, "Participant A", result[0].Name)
	assert.Equal(t, "Participant B", result[1].Name)
	assert.Equal(t, "Participant C", result[2].Name)
}

func TestParticipantService_GetParticipantsList_ReturnsEmptyWhenNoData(t *testing.T) {
	// Listing participants when repository is empty returns empty list
	repo := testutil.NewFakeParticipantRepository()

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantsList(context.Background(), dto.ParticipantListFilters{Status: "active"})

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestParticipantService_GetParticipantsList_StatusFilterIsCaseInsensitive(t *testing.T) {
	// Uppercase status filter matches lowercase status values
	repo := testutil.NewFakeParticipantRepository()
	repo.Participants = []domain.Participant{
		withParticipantStatus(buildParticipant(1, "Participant A"), "active"),
		withParticipantStatus(buildParticipant(2, "Participant B"), "active"),
		withParticipantStatus(buildParticipant(3, "Participant C"), "inactive"),
	}

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantsList(context.Background(), dto.ParticipantListFilters{Status: "ACTIVE"})

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "Participant A", result[0].Name)
	assert.Equal(t, "Participant B", result[1].Name)
}

func TestParticipantService_GetParticipantsList_InvalidDateRangeValidation(t *testing.T) {
	// createdAfter must be before createdBefore
	repo := testutil.NewFakeParticipantRepository()

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantsList(context.Background(), dto.ParticipantListFilters{
		CreatedAfter:  "2024-12-31T00:00:00Z",
		CreatedBefore: "2024-01-01T00:00:00Z",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Contains(t, err.Error(), "createdAfter")
}

func TestParticipantService_GetParticipantByChainId_EmptyChainIdReturnsValidationError(t *testing.T) {
	// An empty chainId is rejected before querying the database
	repo := testutil.NewFakeParticipantRepository()

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantByChainId(context.Background(), "")

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
}

func TestParticipantService_GetParticipantByChainId_NonIntegerChainIdReturnsValidationError(t *testing.T) {
	// A non-numeric chainId is rejected before querying the database
	repo := testutil.NewFakeParticipantRepository()

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantByChainId(context.Background(), "not-a-number")

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Contains(t, err.Error(), "integer")
}

func TestParticipantService_GetParticipantByChainId_DatabaseErrorIsWrapped(t *testing.T) {
	// A non-not-found database error is wrapped and returned as an InternalError
	repo := testutil.NewFakeParticipantRepository()
	repo.Error = errors.New("database query failed")

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantByChainId(context.Background(), "123")

	require.Error(t, err)
	assert.Nil(t, result)
	var internalErr *InternalError
	assert.True(t, errors.As(err, &internalErr))
}

func TestParticipantService_GetParticipantsList_InvalidStatusReturnsValidationError(t *testing.T) {
	// An unrecognized status value is rejected with a validation error listing valid options
	repo := testutil.NewFakeParticipantRepository()

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantsList(context.Background(), dto.ParticipantListFilters{
		Status: "unknown-status",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Contains(t, err.Error(), "status")
}

func TestParticipantService_GetParticipantsList_InvalidRoleReturnsValidationError(t *testing.T) {
	// An unrecognized role value is rejected with a validation error listing valid options
	repo := testutil.NewFakeParticipantRepository()

	service := NewParticipantService(repo, &testutil.StubLogger{})

	result, err := service.GetParticipantsList(context.Background(), dto.ParticipantListFilters{
		Role: "superadmin",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Contains(t, err.Error(), "role")
}
