package core

import (
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/types"
)

// ContractLog re-export types from the shared types package
type ContractLog = types.ContractLog

// UnmarshalEventData unmarshals the RawEventData from a ContractLog into the specified type.
func UnmarshalEventData[T any](log ContractLog) (T, error) {
	return types.UnmarshalEventData[T](log)
}

var (
	ErrContractNotFound   = types.ErrContractNotFound
	ErrUnwatchedEvent     = types.ErrUnwatchedEvent
	ErrFailedToParseEvent = types.ErrFailedToParseEvent
	ErrLogNoTopics        = types.ErrLogNoTopics
	ErrIssuerMismatch     = types.ErrIssuerMismatch
)
