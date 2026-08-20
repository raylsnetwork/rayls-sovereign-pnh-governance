package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ContractLog represents a parsed log from a smart contract
type ContractLog struct {
	ContractName    string
	ContractAddress string
	EventName       string
	RawEventData    json.RawMessage
	BlockNumber     uint64
	TransactionHash string
	LogIndex        uint64
	TxIndex         uint
	BlockTimestamp  time.Time
}

// GetID returns a unique identifier for this log entry, suitable for deduplication.
func (c ContractLog) GetID() string {
	return fmt.Sprintf("%d-%d-%d", c.BlockNumber, c.TxIndex, c.LogIndex)
}

// UnmarshalEventData unmarshals the RawEventData into the specified type.
func UnmarshalEventData[T any](log ContractLog) (T, error) {
	var result T
	if err := json.Unmarshal(log.RawEventData, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal event data: %w", err)
	}
	return result, nil
}

// Error types for ContractLog
var (
	ErrContractNotFound   = errors.New("contract not found")
	ErrUnwatchedEvent     = errors.New("unwatched event")
	ErrFailedToParseEvent = errors.New("failed to parse event")
	ErrLogNoTopics        = errors.New("log has no topics")
	ErrIssuerMismatch     = errors.New("token issuer mismatch")
)
