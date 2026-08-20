package indexer

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/events"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/config"
)

// MapKeys returns a slice with all keys from a generic map
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Keccak256Hash computes the Keccak 256 hash of a string
func Keccak256Hash(data string) common.Hash {
	return crypto.Keccak256Hash([]byte(data))
}

// NewParserRegistry constructs the address -> ContractParsers map using config and instantiated contracts.
func NewParserRegistry(cfg *config.Config, c *Contracts) map[common.Address]ContractParsers {
	return map[common.Address]ContractParsers{
		common.HexToAddress(cfg.PrivateHub.TokenCore): {
			Name:    events.ContractTokenCore,
			Parsers: tokenCoreParsers(c),
		},
		common.HexToAddress(cfg.PrivateHub.TokenFreezeManager): {
			Name:    events.ContractTokenFreezeManager,
			Parsers: tokenFreezeManagerParsers(c),
		},
		common.HexToAddress(cfg.PrivateHub.EnygmaTokenManager): {
			Name:    events.ContractEnygmaTokenManager,
			Parsers: enygmaTokenManagerParsers(c),
		},
		common.HexToAddress(cfg.PrivateHub.Teleport): {
			Name:    events.ContractTeleport,
			Parsers: teleportParsers(c),
		},
		common.HexToAddress(cfg.PrivateHub.ParticipantCore): {
			Name:    events.ContractParticipantCore,
			Parsers: participantCoreParsers(c),
		},
		common.HexToAddress(cfg.PrivateHub.AuditManager): {
			Name:    events.ContractAuditManager,
			Parsers: auditManagerParsers(c),
		},
		common.HexToAddress(cfg.PrivateHub.ProofsAddress): {
			Name:    events.ContractProofs,
			Parsers: proofsParsers(c),
		},
		common.HexToAddress(cfg.PrivateHub.EnygmaTeleport): {
			Name:    events.ContractEnygmaTeleport,
			Parsers: enygmaTeleportParsers(c),
		},
		common.HexToAddress(cfg.PrivateHub.DvpTeleport): {
			Name:    events.ContractDvpTeleport,
			Parsers: dvpTeleportParsers(c),
		},
	}
}

// tokenCoreParsers returns the event parsers for the TokenCore contract
func tokenCoreParsers(c *Contracts) []EventParser {
	return []EventParser{
		{events.Erc20TokenRegistered, Keccak256Hash(events.Erc20TokenRegisteredSig), func(log types.Log) (any, error) {
			event, err := c.TokenCore.UnpackErc20TokenRegisteredEvent(&log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse Erc20TokenRegistered event: %w", err)
			}
			return event, nil
		}},
		{
			events.Erc721TokenRegistered,
			Keccak256Hash(events.Erc721TokenRegisteredSig),
			func(log types.Log) (any, error) {
				event, err := c.TokenCore.UnpackErc721TokenRegisteredEvent(&log)
				if err != nil {
					return nil, fmt.Errorf("failed to parse Erc721TokenRegistered event: %w", err)
				}
				return event, nil
			},
		},
		{
			events.Erc1155TokenRegistered,
			Keccak256Hash(events.Erc1155TokenRegisteredSig),
			func(log types.Log) (any, error) {
				event, err := c.TokenCore.UnpackErc1155TokenRegisteredEvent(&log)
				if err != nil {
					return nil, fmt.Errorf("failed to parse Erc1155TokenRegistered event: %w", err)
				}
				return event, nil
			},
		},
		{events.TokenStatusUpdated, Keccak256Hash(events.TokenStatusUpdatedSig), func(log types.Log) (any, error) {
			event, err := c.TokenCore.UnpackTokenStatusUpdatedEvent(&log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse TokenStatusUpdated event: %w", err)
			}
			return event, nil
		}},
		{events.TokenBalanceUpdated, Keccak256Hash(events.TokenBalanceUpdatedSig), func(log types.Log) (any, error) {
			event, err := c.TokenCore.UnpackTokenBalanceUpdatedEvent(&log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse TokenBalanceUpdated event: %w", err)
			}
			return event, nil
		}},
		{
			events.DvpErc721TokenRegistered,
			Keccak256Hash(events.DvpErc721TokenRegisteredSig),
			func(log types.Log) (any, error) {
				event, err := c.TokenCore.UnpackDvpErc721TokenRegisteredEvent(&log)
				if err != nil {
					return nil, fmt.Errorf("failed to parse DvpErc721TokenRegistered event: %w", err)
				}
				return event, nil
			},
		},
		{
			events.DvpErc1155TokenRegistered,
			Keccak256Hash(events.DvpErc1155TokenRegisteredSig),
			func(log types.Log) (any, error) {
				event, err := c.TokenCore.UnpackDvpErc1155TokenRegisteredEvent(&log)
				if err != nil {
					return nil, fmt.Errorf("failed to parse DvpErc1155TokenRegistered event: %w", err)
				}
				return event, nil
			},
		},
	}
}

// tokenFreezeManagerParsers returns the event parsers for the TokenFreezeManager contract
func tokenFreezeManagerParsers(c *Contracts) []EventParser {
	return []EventParser{
		{
			events.TokenFreezeStatusChanged,
			Keccak256Hash(events.TokenFreezeStatusChangedSig),
			func(log types.Log) (any, error) {
				event, err := c.TokenFreezeManager.UnpackTokenFreezeStatusChangedEvent(&log)
				if err != nil {
					return nil, fmt.Errorf("failed to parse TokenFreezeStatusChanged event: %w", err)
				}
				return event, nil
			},
		},
	}
}

// enygmaTokenManagerParsers returns the event parsers for the EnygmaTokenManager contract
func enygmaTokenManagerParsers(c *Contracts) []EventParser {
	return []EventParser{
		{
			events.EnygmaTokenRegistered,
			Keccak256Hash(events.EnygmaTokenRegisteredSig),
			func(log types.Log) (any, error) {
				event, err := c.EnygmaTokenManager.UnpackEnygmaTokenRegisteredEvent(&log)
				if err != nil {
					return nil, fmt.Errorf("failed to parse EnygmaTokenRegistered event: %w", err)
				}
				return event, nil
			},
		},
	}
}

// teleportParsers returns the event parsers for the Teleport contract
func teleportParsers(c *Contracts) []EventParser {
	return []EventParser{
		{
			events.AtomicMessageAdditionalDataBatch,
			Keccak256Hash(events.AtomicMessageAdditionalDataBatchSig),
			func(log types.Log) (any, error) {
				event, err := c.Teleport.UnpackAtomicMessageAdditionalDataBatchEvent(&log)
				if err != nil {
					return nil, fmt.Errorf("failed to parse AtomicMessageAdditionalDataBatch event: %w", err)
				}
				return event, nil
			},
		},
		{
			events.AtomicMessageStatusChangedBatch,
			Keccak256Hash(events.AtomicMessageStatusChangedBatchSig),
			func(log types.Log) (any, error) {
				event, err := c.Teleport.UnpackAtomicMessageStatusChangedBatchEvent(&log)
				if err != nil {
					return nil, fmt.Errorf("failed to parse AtomicMessageStatusChangedBatch event: %w", err)
				}
				return event, nil
			},
		},
		{
			events.EncryptedDataBatchStored,
			Keccak256Hash(events.EncryptedDataBatchStoredSig),
			func(log types.Log) (any, error) {
				event, err := c.Teleport.UnpackEncryptedDataBatchStoredEvent(&log)
				if err != nil {
					return nil, fmt.Errorf("failed to parse EncryptedDataBatchStored event: %w", err)
				}
				return event, nil
			},
		},
	}
}

// participantCoreParsers returns the event parsers for the ParticipantCore contract
func participantCoreParsers(c *Contracts) []EventParser {
	return []EventParser{
		{
			events.ParticipantRegistered,
			Keccak256Hash(events.ParticipantRegisteredSig),
			func(log types.Log) (any, error) {
				event, err := c.ParticipantCore.UnpackParticipantRegisteredEvent(&log)
				if err != nil {
					return nil, fmt.Errorf("failed to parse ParticipantRegistered event: %w", err)
				}
				return event, nil
			},
		},
		{events.ParticipantUpdated, Keccak256Hash(events.ParticipantUpdatedSig), func(log types.Log) (any, error) {
			event, err := c.ParticipantCore.UnpackParticipantUpdatedEvent(&log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse ParticipantUpdated event: %w", err)
			}
			return event, nil
		}},
	}
}

// auditManagerParsers returns the event parsers for the AuditManager contract
func auditManagerParsers(c *Contracts) []EventParser {
	return []EventParser{
		{events.NewAuditOrChainInfo, Keccak256Hash(events.NewAuditOrChainInfoSig), func(log types.Log) (any, error) {
			event, err := c.AuditManager.UnpackNewAuditOrChainInfoEvent(&log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse NewAuditOrChainInfo event: %w", err)
			}
			return event, nil
		}},
	}
}

// proofsParsers returns the event parsers for the Proofs contract
func proofsParsers(c *Contracts) []EventParser {
	return []EventParser{
		{events.HeaderProofSubmitted, Keccak256Hash(events.HeaderProofSubmittedSig), func(log types.Log) (any, error) {
			event, err := c.Proofs.UnpackHeaderProofSubmittedEvent(&log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse HeaderProofSubmitted event: %w", err)
			}
			return event, nil
		}},
	}
}

// enygmaTeleportParsers returns the event parsers for the EnygmaTeleport contract
func enygmaTeleportParsers(c *Contracts) []EventParser {
	return []EventParser{
		{events.EnygmaTransfer, Keccak256Hash(events.EnygmaTransferSig), func(log types.Log) (any, error) {
			event, err := c.EnygmaTeleport.UnpackEnygmaTransferEvent(&log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse EnygmaTransfer event: %w", err)
			}
			return event, nil
		}},
		{
			events.EnygmaTransferCompleted,
			Keccak256Hash(events.EnygmaTransferCompletedSig),
			func(log types.Log) (any, error) {
				event, err := c.EnygmaTeleport.UnpackEnygmaTransferCompletedEvent(&log)
				if err != nil {
					return nil, fmt.Errorf("failed to parse EnygmaTransferCompleted event: %w", err)
				}
				return event, nil
			},
		},
		{events.EnygmaSupplyUpdated, Keccak256Hash(events.EnygmaSupplyUpdatedSig), func(log types.Log) (any, error) {
			event, err := c.EnygmaTeleport.UnpackEnygmaSupplyUpdatedEvent(&log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse EnygmaSupplyUpdated event: %w", err)
			}
			return event, nil
		}},
		{
			events.EnygmaDvpBalanceUpdated,
			Keccak256Hash(events.EnygmaDvpBalanceUpdatedSig),
			func(log types.Log) (any, error) {
				event, err := c.EnygmaTeleport.UnpackEnygmaDvpBalanceUpdatedEvent(&log)
				if err != nil {
					return nil, fmt.Errorf("failed to parse DvpBalanceUpdated event: %w", err)
				}
				return event, nil
			},
		},
	}
}

// dvpTeleportParsers returns the event parsers for the DvpTeleport contract
func dvpTeleportParsers(c *Contracts) []EventParser {
	return []EventParser{
		{events.ERCDvpBalanceUpdated, Keccak256Hash(events.ERCDvpBalanceUpdatedSig), func(log types.Log) (any, error) {
			event, err := c.DvpTeleport.UnpackERCDvpBalanceUpdatedEvent(&log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse ERCDvpBalanceUpdated event: %w", err)
			}
			return event, nil
		}},
		{events.SwapInitiated, Keccak256Hash(events.SwapInitiatedSig), func(log types.Log) (any, error) {
			event, err := c.DvpTeleport.UnpackSwapInitiatedEvent(&log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse SwapInitiated event: %w", err)
			}
			return event, nil
		}},
		{events.SwapCompleted, Keccak256Hash(events.SwapCompletedSig), func(log types.Log) (any, error) {
			event, err := c.DvpTeleport.UnpackSwapCompletedEvent(&log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse SwapCompleted event: %w", err)
			}
			return event, nil
		}},
		{events.SwapCancelled, Keccak256Hash(events.SwapCancelledSig), func(log types.Log) (any, error) {
			event, err := c.DvpTeleport.UnpackSwapCancelledEvent(&log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse SwapCancelled event: %w", err)
			}
			return event, nil
		}},
		{events.SwapTimedOut, Keccak256Hash(events.SwapTimedOutSig), func(log types.Log) (any, error) {
			event, err := c.DvpTeleport.UnpackSwapTimedOutEvent(&log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse SwapTimedOut event: %w", err)
			}
			return event, nil
		}},
	}
}
