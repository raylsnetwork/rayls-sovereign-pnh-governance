package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/testutil"
)

// stubContractMQ records Push calls and can return a configured error at a specific index.
type stubContractMQ struct {
	pushed  []core.ContractLog
	failAt  int
	failErr error
}

func (s *stubContractMQ) Push(_ context.Context, log core.ContractLog) error {
	if s.failErr != nil && len(s.pushed) == s.failAt {
		return s.failErr
	}
	s.pushed = append(s.pushed, log)
	return nil
}

func TestLogPublisher_Publish_EmptyLogs(t *testing.T) {
	// Publishing an empty slice returns nil and makes zero Push calls
	mq := &stubContractMQ{}
	pub := NewLogPublisher(mq, &testutil.StubLogger{})

	err := pub.Publish(context.Background(), nil)

	require.NoError(t, err)
	assert.Empty(t, mq.pushed)
}

func TestLogPublisher_Publish_SingleLog(t *testing.T) {
	// A single log is published via Push
	mq := &stubContractMQ{}
	pub := NewLogPublisher(mq, &testutil.StubLogger{})
	logs := []core.ContractLog{
		{ContractName: "TokenCore", EventName: "TokenRegistered", BlockNumber: 100},
	}

	err := pub.Publish(context.Background(), logs)

	require.NoError(t, err)
	require.Len(t, mq.pushed, 1)
	assert.Equal(t, "TokenCore", mq.pushed[0].ContractName)
}

func TestLogPublisher_Publish_MultipleLogs(t *testing.T) {
	// All logs are published in order
	mq := &stubContractMQ{}
	pub := NewLogPublisher(mq, &testutil.StubLogger{})
	logs := []core.ContractLog{
		{ContractName: "TokenCore", EventName: "TokenRegistered", BlockNumber: 100},
		{ContractName: "Teleport", EventName: "MessageSent", BlockNumber: 101},
		{ContractName: "Proofs", EventName: "HeaderProofSubmitted", BlockNumber: 102},
	}

	err := pub.Publish(context.Background(), logs)

	require.NoError(t, err)
	require.Len(t, mq.pushed, 3)
	assert.Equal(t, "TokenCore", mq.pushed[0].ContractName)
	assert.Equal(t, "Teleport", mq.pushed[1].ContractName)
	assert.Equal(t, "Proofs", mq.pushed[2].ContractName)
}

func TestLogPublisher_Publish_FirstPushFails(t *testing.T) {
	// Error on first push is returned immediately, no logs published
	mq := &stubContractMQ{failAt: 0, failErr: errors.New("connection refused")}
	pub := NewLogPublisher(mq, &testutil.StubLogger{})
	logs := []core.ContractLog{
		{ContractName: "TokenCore", EventName: "TokenRegistered", BlockNumber: 100},
		{ContractName: "Teleport", EventName: "MessageSent", BlockNumber: 101},
	}

	err := pub.Publish(context.Background(), logs)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
	assert.Contains(t, err.Error(), "1/2")
	assert.Empty(t, mq.pushed)
}

func TestLogPublisher_Publish_MiddlePushFails(t *testing.T) {
	// Error on second of three logs stops publishing; only first log was pushed
	mq := &stubContractMQ{failAt: 1, failErr: errors.New("timeout")}
	pub := NewLogPublisher(mq, &testutil.StubLogger{})
	logs := []core.ContractLog{
		{ContractName: "TokenCore", EventName: "TokenRegistered", BlockNumber: 100},
		{ContractName: "Teleport", EventName: "MessageSent", BlockNumber: 101},
		{ContractName: "Proofs", EventName: "HeaderProofSubmitted", BlockNumber: 102},
	}

	err := pub.Publish(context.Background(), logs)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
	assert.Contains(t, err.Error(), "2/3")
	require.Len(t, mq.pushed, 1)
	assert.Equal(t, "TokenCore", mq.pushed[0].ContractName)
}
