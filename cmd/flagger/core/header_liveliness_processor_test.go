package core

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/flagger/testutil"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

const expirationPeriod = 5 * time.Minute

// Helpers
func buildHeaderProof(chainId, blockNumber int64, createdAt time.Time) domain.HeaderProofEvent {
	return domain.HeaderProofEvent{
		ChainID:     domain.BigInt{Int: big.NewInt(chainId)},
		BlockNumber: domain.BigInt{Int: big.NewInt(blockNumber)},
		CreatedAt:   createdAt,
	}
}

func expiredHeader(chainId, blockNumber int64) domain.HeaderProofEvent {
	return buildHeaderProof(chainId, blockNumber, time.Now().Add(-6*time.Minute))
}

func validHeader(chainId, blockNumber int64) domain.HeaderProofEvent {
	return buildHeaderProof(chainId, blockNumber, time.Now().Add(-1*time.Minute))
}

func TestHeaderLivelinessProcessor_FlagsExpiredHeader(t *testing.T) {
	proofRepo := testutil.NewFakeHeaderProofRepository()
	proofRepo.HeaderProofs = []domain.HeaderProofEvent{
		expiredHeader(1, 100),
	}
	flagRepo := testutil.NewFakeHeaderFlagEventRepository()

	processor := NewHeaderLivelinessProcessor(proofRepo, flagRepo, expirationPeriod, &testutil.StubLogger{})

	err := processor.Start(context.Background())
	require.NoError(t, err)
	assert.True(t, flagRepo.IsFlagged(big.NewInt(1), big.NewInt(100)))
}

func TestHeaderLivelinessProcessor_DoesNotFlagNonExpiredHeader(t *testing.T) {
	proofRepo := testutil.NewFakeHeaderProofRepository()
	proofRepo.HeaderProofs = []domain.HeaderProofEvent{
		validHeader(1, 100),
	}
	flagRepo := testutil.NewFakeHeaderFlagEventRepository()

	processor := NewHeaderLivelinessProcessor(proofRepo, flagRepo, expirationPeriod, &testutil.StubLogger{})

	err := processor.Start(context.Background())
	require.NoError(t, err)
	assert.Empty(t, flagRepo.FlaggedParticipants)
}

func TestHeaderLivelinessProcessor_HandlesEmptyHeaders(t *testing.T) {
	proofRepo := testutil.NewFakeHeaderProofRepository()
	flagRepo := testutil.NewFakeHeaderFlagEventRepository()

	processor := NewHeaderLivelinessProcessor(proofRepo, flagRepo, expirationPeriod, &testutil.StubLogger{})

	err := processor.Start(context.Background())
	require.NoError(t, err)
	assert.Empty(t, flagRepo.FlaggedParticipants)
}

func TestHeaderLivelinessProcessor_ReturnsErrorWhenGetHeadersFails(t *testing.T) {
	proofRepo := testutil.NewFakeHeaderProofRepository()
	proofRepo.Error = errors.New("database error")
	flagRepo := testutil.NewFakeHeaderFlagEventRepository()

	processor := NewHeaderLivelinessProcessor(proofRepo, flagRepo, expirationPeriod, &testutil.StubLogger{})

	err := processor.Start(context.Background())
	require.Error(t, err)
}

func TestHeaderLivelinessProcessor_ContinuesWhenFlagFails(t *testing.T) {
	proofRepo := testutil.NewFakeHeaderProofRepository()
	proofRepo.HeaderProofs = []domain.HeaderProofEvent{
		expiredHeader(1, 100),
	}
	flagRepo := testutil.NewFakeHeaderFlagEventRepository()
	flagRepo.Error = errors.New("flag error")

	processor := NewHeaderLivelinessProcessor(proofRepo, flagRepo, expirationPeriod, &testutil.StubLogger{})

	err := processor.Start(context.Background())
	require.NoError(t, err)
}

func TestHeaderLivelinessProcessor_FlagsMultipleExpiredHeaders(t *testing.T) {
	proofRepo := testutil.NewFakeHeaderProofRepository()
	proofRepo.HeaderProofs = []domain.HeaderProofEvent{
		expiredHeader(1, 100),
		expiredHeader(2, 200),
		expiredHeader(3, 300),
	}
	flagRepo := testutil.NewFakeHeaderFlagEventRepository()

	processor := NewHeaderLivelinessProcessor(proofRepo, flagRepo, expirationPeriod, &testutil.StubLogger{})

	err := processor.Start(context.Background())

	// All participants are flagged
	require.NoError(t, err)
	assert.True(t, flagRepo.IsFlagged(big.NewInt(1), big.NewInt(100)))
	assert.True(t, flagRepo.IsFlagged(big.NewInt(2), big.NewInt(200)))
	assert.True(t, flagRepo.IsFlagged(big.NewInt(3), big.NewInt(300)))
	assert.Len(t, flagRepo.FlaggedParticipants, 3)
}

func TestHeaderLivelinessProcessor_FlagsOnlyExpiredHeaders(t *testing.T) {
	proofRepo := testutil.NewFakeHeaderProofRepository()
	proofRepo.HeaderProofs = []domain.HeaderProofEvent{
		expiredHeader(1, 100), // expired
		validHeader(2, 200),   // valid
		expiredHeader(3, 300), // expired
	}
	flagRepo := testutil.NewFakeHeaderFlagEventRepository()

	processor := NewHeaderLivelinessProcessor(proofRepo, flagRepo, expirationPeriod, &testutil.StubLogger{})

	err := processor.Start(context.Background())

	// Only expired headers are flagged
	require.NoError(t, err)
	assert.True(t, flagRepo.IsFlagged(big.NewInt(1), big.NewInt(100)))
	assert.False(t, flagRepo.IsFlagged(big.NewInt(2), big.NewInt(200)))
	assert.True(t, flagRepo.IsFlagged(big.NewInt(3), big.NewInt(300)))
	assert.Len(t, flagRepo.FlaggedParticipants, 2)
}

func TestHeaderLivelinessProcessor_RunStopsOnContextCancellation(t *testing.T) {
	proofRepo := testutil.NewFakeHeaderProofRepository()
	flagRepo := testutil.NewFakeHeaderFlagEventRepository()
	processor := NewHeaderLivelinessProcessor(proofRepo, flagRepo, expirationPeriod, &testutil.StubLogger{})

	ctx, cancel := context.WithCancel(context.Background())
	ticker := make(chan time.Time)

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

func TestHeaderLivelinessProcessor_UnflagsRecoveredParticipant(t *testing.T) {
	proofRepo := testutil.NewFakeHeaderProofRepository()
	proofRepo.HeaderProofs = []domain.HeaderProofEvent{
		validHeader(1, 200), // node live again; block advanced past the one it was flagged at
	}
	flagRepo := testutil.NewFakeHeaderFlagEventRepository()
	// Pre-existing liveliness flag from an earlier gap (e.g. a VPN outage).
	_, err := flagRepo.FlagParticipant(context.Background(), big.NewInt(1), big.NewInt(100),
		uint8(types.HeaderFlagReasonLiveliness), uint8(types.HeaderFlagInitiatorAutomaticSystem))
	require.NoError(t, err)
	require.True(t, flagRepo.IsFlagged(big.NewInt(1), big.NewInt(100)))

	processor := NewHeaderLivelinessProcessor(proofRepo, flagRepo, expirationPeriod, &testutil.StubLogger{})

	err = processor.Start(context.Background())

	require.NoError(t, err)
	assert.False(t, flagRepo.IsFlagged(big.NewInt(1), big.NewInt(100)),
		"a participant submitting fresh proofs should be auto-unflagged")
	assert.Empty(t, flagRepo.FlaggedParticipants)
}

func TestHeaderLivelinessProcessor_DoesNotUnflagNonLivelinessFlag(t *testing.T) {
	// A flag raised for a non-liveliness reason must survive a liveliness recovery: Start clears
	// only flags scoped to HeaderFlagReasonLiveliness, so other-reason flags stay for review.
	const otherReason = uint8(99) // stand-in for a future non-liveliness reason (e.g. manual/compliance)

	proofRepo := testutil.NewFakeHeaderProofRepository()
	proofRepo.HeaderProofs = []domain.HeaderProofEvent{
		validHeader(1, 200), // fresh proofs, so Start takes the recovery (unflag) path for chain 1
	}
	flagRepo := testutil.NewFakeHeaderFlagEventRepository()
	// Pre-seed a flag raised for a reason other than liveliness.
	_, err := flagRepo.FlagParticipant(context.Background(), big.NewInt(1), big.NewInt(100),
		otherReason, uint8(types.HeaderFlagInitiatorAutomaticSystem))
	require.NoError(t, err)
	require.True(t, flagRepo.IsFlagged(big.NewInt(1), big.NewInt(100)))

	processor := NewHeaderLivelinessProcessor(proofRepo, flagRepo, expirationPeriod, &testutil.StubLogger{})

	err = processor.Start(context.Background())

	require.NoError(t, err)
	assert.True(t, flagRepo.IsFlagged(big.NewInt(1), big.NewInt(100)),
		"a flag raised for a non-liveliness reason must not be cleared by liveliness self-heal")
}

func TestHeaderLivelinessProcessor_ContinuesWhenUnflagFails(t *testing.T) {
	proofRepo := testutil.NewFakeHeaderProofRepository()
	proofRepo.HeaderProofs = []domain.HeaderProofEvent{
		validHeader(1, 100),
	}
	flagRepo := testutil.NewFakeHeaderFlagEventRepository()
	flagRepo.Error = errors.New("unflag error")

	processor := NewHeaderLivelinessProcessor(proofRepo, flagRepo, expirationPeriod, &testutil.StubLogger{})

	err := processor.Start(context.Background())

	// The unflag error is logged per-chain, not propagated.
	require.NoError(t, err)
}
