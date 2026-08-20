package services

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/events"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/testutil"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/contracts/Proofs"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
)

// fakeHeaderProofBatchRepo records CreateBatch calls and can inject errors.
type fakeHeaderProofBatchRepo struct {
	batches [][]*domain.HeaderProofEvent
	err     error
}

var _ core.HeaderProofEventRepository = (*fakeHeaderProofBatchRepo)(nil)

func (f *fakeHeaderProofBatchRepo) Create(_ context.Context, _ *domain.HeaderProofEvent) error {
	return nil
}

func (f *fakeHeaderProofBatchRepo) CreateBatch(_ context.Context, evts []*domain.HeaderProofEvent) error {
	if f.err != nil {
		return f.err
	}
	cp := make([]*domain.HeaderProofEvent, len(evts))
	copy(cp, evts)
	f.batches = append(f.batches, cp)
	return nil
}

func (f *fakeHeaderProofBatchRepo) GetByBlockNumber(
	_ context.Context,
	_ *big.Int,
	_ *big.Int,
) (*domain.HeaderProofEvent, error) {
	return nil, nil //nolint:nilnil // test stub; nil means "not found" with no error
}

func makeHeaderProofMessage(t *testing.T, chainId, blockNumber int64, tracker *ackTracker) core.Message {
	t.Helper()
	event := &Proofs.ProofsHeaderProofSubmitted{
		ChainId:     big.NewInt(chainId),
		BlockNumber: big.NewInt(blockNumber),
		HeaderHash:  [32]byte{0x01, 0x02, 0x03},
	}
	return core.Message{
		Log: core.ContractLog{
			ContractName: events.ContractProofs,
			EventName:    events.HeaderProofSubmitted,
			RawEventData: testutil.MustMarshal(t, event),
		},
		Ack: tracker.ack,
	}
}

func newBatchProcessor(repo core.HeaderProofEventRepository) *ProofsBatchProcessor {
	return NewProofsBatchProcessor(nil, repo, &testutil.StubLogger{}, 10, 0)
}

func TestProofsBatchProcessor_Flush_SuccessfulBatch(t *testing.T) {
	// All parseable messages are persisted and acked after CreateBatch succeeds
	repo := &fakeHeaderProofBatchRepo{}
	p := newBatchProcessor(repo)

	ack1, ack2 := &ackTracker{}, &ackTracker{}
	messages := []core.Message{
		makeHeaderProofMessage(t, 1, 100, ack1),
		makeHeaderProofMessage(t, 2, 200, ack2),
	}

	p.flush(context.Background(), messages)

	require.Len(t, repo.batches, 1)
	assert.Len(t, repo.batches[0], 2)
	assert.True(t, ack1.called)
	assert.True(t, ack2.called)
}

func TestProofsBatchProcessor_Flush_ParseFailure(t *testing.T) {
	// Unparseable messages are acked immediately to avoid redelivery loops; parseable ones are batched and acked
	repo := &fakeHeaderProofBatchRepo{}
	p := newBatchProcessor(repo)

	ackBad, ackGood := &ackTracker{}, &ackTracker{}
	messages := []core.Message{
		{
			Log: core.ContractLog{
				EventName:    events.HeaderProofSubmitted,
				RawEventData: []byte(`not-valid-json`),
			},
			Ack: ackBad.ack,
		},
		makeHeaderProofMessage(t, 1, 100, ackGood),
	}

	p.flush(context.Background(), messages)

	require.Len(t, repo.batches, 1)
	assert.Len(t, repo.batches[0], 1)
	assert.True(t, ackBad.called)
	assert.True(t, ackGood.called)
}

func TestProofsBatchProcessor_Flush_DBFailure(t *testing.T) {
	// When CreateBatch fails, parseable messages are NOT acked so NATS redelivers them
	repo := &fakeHeaderProofBatchRepo{err: errors.New("db unavailable")}
	p := newBatchProcessor(repo)

	ack1, ack2 := &ackTracker{}, &ackTracker{}
	messages := []core.Message{
		makeHeaderProofMessage(t, 1, 100, ack1),
		makeHeaderProofMessage(t, 2, 200, ack2),
	}

	p.flush(context.Background(), messages)

	assert.Empty(t, repo.batches)
	assert.False(t, ack1.called)
	assert.False(t, ack2.called)
}

func TestProofsBatchProcessor_Flush_MixedBatch(t *testing.T) {
	// Unknown-event messages are acked and skipped immediately; parseable ones proceed to batch insert
	repo := &fakeHeaderProofBatchRepo{}
	p := newBatchProcessor(repo)

	ackUnknown, ackValid := &ackTracker{}, &ackTracker{}
	messages := []core.Message{
		{
			Log: core.ContractLog{EventName: "UnknownEvent"},
			Ack: ackUnknown.ack,
		},
		makeHeaderProofMessage(t, 1, 100, ackValid),
	}

	p.flush(context.Background(), messages)

	require.Len(t, repo.batches, 1)
	assert.Len(t, repo.batches[0], 1)
	assert.True(t, ackUnknown.called)
	assert.True(t, ackValid.called)
}
