package testutil

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/config"
)

// MustMarshal JSON-marshals v, strips the "Raw" field (go-ethereum types.Log
// has strict JSON unmarshalling that breaks round-trips), and returns the result.
func MustMarshal(t *testing.T, v any) json.RawMessage {
	t.Helper()
	data, err := json.Marshal(v)
	require.NoError(t, err)

	var m map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &m))
	delete(m, "Raw")
	stripped, err := json.Marshal(m)
	require.NoError(t, err)
	return stripped
}

// ============================================================================
// Logger Stub
// ============================================================================

// StubLogger is a no-op logger implementation for testing
type StubLogger struct{}

func (l *StubLogger) Debug(msg string, args ...any) {}
func (l *StubLogger) Info(msg string, args ...any)  {}
func (l *StubLogger) Warn(msg string, args ...any)  {}
func (l *StubLogger) Error(msg string, args ...any) {}

// ============================================================================
// ConfigProvider Stub
// ============================================================================

// StubConfigProvider is a manual stub for the ConfigProvider interface.
type StubConfigProvider struct {
	// Return values
	StartingBlockNumber *big.Int
	BatchSize           int64
	Config              *config.Config
	ConfigErr           error

	// Spy fields - for verifying method calls
	GetConfigCalled        bool
	GetConfigCallCount     int
	GetStartingBlockCalled bool
	GetBatchSizeCalled     bool
	LastGetConfigContext   context.Context
}

// NewStubConfigProvider creates a StubConfigProvider
func NewStubConfigProvider() *StubConfigProvider {
	return &StubConfigProvider{
		StartingBlockNumber: big.NewInt(0),
		BatchSize:           10,
		Config:              &config.Config{},
		ConfigErr:           nil,
	}
}

func (s *StubConfigProvider) GetStartingBlockNumber() *big.Int {
	s.GetStartingBlockCalled = true
	return s.StartingBlockNumber
}

func (s *StubConfigProvider) GetBatchSize() int64 {
	s.GetBatchSizeCalled = true
	return s.BatchSize
}

func (s *StubConfigProvider) GetConfig(ctx context.Context) (*config.Config, error) {
	s.GetConfigCalled = true
	s.GetConfigCallCount++
	s.LastGetConfigContext = ctx
	return s.Config, s.ConfigErr
}

// ============================================================================
// Provider Stub
// ============================================================================

// StubProvider is a manual stub for the core.Provider interface.
type StubProvider struct {
	Timestamp uint64
}

func (s *StubProvider) GetLatestBlock(_ context.Context) (*types.Block, error) {
	return types.NewBlockWithHeader(&types.Header{Time: s.Timestamp}), nil
}

func (s *StubProvider) GetBlockByNumber(_ context.Context, _ uint64) (*types.Block, error) {
	return types.NewBlockWithHeader(&types.Header{Time: s.Timestamp}), nil
}
