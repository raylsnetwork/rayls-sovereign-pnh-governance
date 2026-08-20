package indexer

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/events"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/AuditManagerV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/DvpTeleport"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/EnygmaTeleport"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/EnygmaTokenManagerV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/ParticipantCoreV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/Proofs"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TeleportV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TokenCoreV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TokenFreezeManagerV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
)

// eventSigEntry maps an event name and its manually-written signature to the contract ABI
// that contains the canonical (source-of-truth) definition.
type eventSigEntry struct {
	eventName string
	sigConst  string // the signature string from events.go
	abiJSON   string // the contract ABI JSON containing the event
}

// allEventSignatures returns every event signature registered in events.go paired with
// the contract ABI that defines the canonical signature.
func allEventSignatures() []eventSigEntry {
	return []eventSigEntry{
		// TokenCore
		{events.Erc20TokenRegistered, events.Erc20TokenRegisteredSig, TokenCoreV1.TokenCoreV1MetaData.ABI},
		{events.Erc721TokenRegistered, events.Erc721TokenRegisteredSig, TokenCoreV1.TokenCoreV1MetaData.ABI},
		{events.Erc1155TokenRegistered, events.Erc1155TokenRegisteredSig, TokenCoreV1.TokenCoreV1MetaData.ABI},
		{events.DvpErc721TokenRegistered, events.DvpErc721TokenRegisteredSig, TokenCoreV1.TokenCoreV1MetaData.ABI},
		{events.DvpErc1155TokenRegistered, events.DvpErc1155TokenRegisteredSig, TokenCoreV1.TokenCoreV1MetaData.ABI},
		{events.TokenStatusUpdated, events.TokenStatusUpdatedSig, TokenCoreV1.TokenCoreV1MetaData.ABI},
		{events.TokenBalanceUpdated, events.TokenBalanceUpdatedSig, TokenCoreV1.TokenCoreV1MetaData.ABI},

		// EnygmaTokenManager
		{
			events.EnygmaTokenRegistered,
			events.EnygmaTokenRegisteredSig,
			EnygmaTokenManagerV1.EnygmaTokenManagerV1MetaData.ABI,
		},

		// Teleport
		{
			events.AtomicMessageAdditionalDataBatch,
			events.AtomicMessageAdditionalDataBatchSig,
			TeleportV1.TeleportV1MetaData.ABI,
		},
		{
			events.AtomicMessageStatusChangedBatch,
			events.AtomicMessageStatusChangedBatchSig,
			TeleportV1.TeleportV1MetaData.ABI,
		},
		{events.EncryptedDataBatchStored, events.EncryptedDataBatchStoredSig, TeleportV1.TeleportV1MetaData.ABI},

		// ParticipantCore
		{
			events.ParticipantRegistered,
			events.ParticipantRegisteredSig,
			ParticipantCoreV1.ParticipantCoreV1MetaData.ABI,
		},
		{events.ParticipantUpdated, events.ParticipantUpdatedSig, ParticipantCoreV1.ParticipantCoreV1MetaData.ABI},

		// AuditManager
		{events.NewAuditOrChainInfo, events.NewAuditOrChainInfoSig, AuditManagerV1.AuditManagerV1MetaData.ABI},

		// Proofs
		{events.HeaderProofSubmitted, events.HeaderProofSubmittedSig, Proofs.ProofsMetaData.ABI},

		// EnygmaTeleport
		{events.EnygmaTransfer, events.EnygmaTransferSig, EnygmaTeleport.EnygmaTeleportMetaData.ABI},
		{events.EnygmaTransferCompleted, events.EnygmaTransferCompletedSig, EnygmaTeleport.EnygmaTeleportMetaData.ABI},
		{events.EnygmaSupplyUpdated, events.EnygmaSupplyUpdatedSig, EnygmaTeleport.EnygmaTeleportMetaData.ABI},
		{events.EnygmaDvpBalanceUpdated, events.EnygmaDvpBalanceUpdatedSig, EnygmaTeleport.EnygmaTeleportMetaData.ABI},

		// DvpTeleport
		{events.ERCDvpBalanceUpdated, events.ERCDvpBalanceUpdatedSig, DvpTeleport.DvpTeleportMetaData.ABI},
		{events.SwapInitiated, events.SwapInitiatedSig, DvpTeleport.DvpTeleportMetaData.ABI},
		{events.SwapCompleted, events.SwapCompletedSig, DvpTeleport.DvpTeleportMetaData.ABI},
		{events.SwapCancelled, events.SwapCancelledSig, DvpTeleport.DvpTeleportMetaData.ABI},
		{events.SwapTimedOut, events.SwapTimedOutSig, DvpTeleport.DvpTeleportMetaData.ABI},

		// TokenFreezeManager
		{
			events.TokenFreezeStatusChanged,
			events.TokenFreezeStatusChangedSig,
			TokenFreezeManagerV1.TokenFreezeManagerV1MetaData.ABI,
		},
	}
}

// ValidateEventSignatures checks that every event signature in events.go produces
// the same keccak256 hash as the canonical signature from the contract ABI.
// Returns an error describing all mismatches, or nil if everything is correct.
// This prevents silent event-dropping bugs caused by stale or incorrect signature strings.
func ValidateEventSignatures(log logger.Logger) error {
	entries := allEventSignatures()
	var mismatches []string

	for _, e := range entries {
		parsedABI, err := abi.JSON(strings.NewReader(e.abiJSON))
		if err != nil {
			mismatches = append(mismatches, fmt.Sprintf(
				"%s: failed to parse contract ABI: %v", e.eventName, err,
			))
			continue
		}

		ev, ok := parsedABI.Events[e.eventName]
		if !ok {
			mismatches = append(mismatches, fmt.Sprintf(
				"%s: event not found in contract ABI", e.eventName,
			))
			continue
		}

		abiHash := crypto.Keccak256Hash([]byte(ev.Sig))
		eventsHash := crypto.Keccak256Hash([]byte(e.sigConst))

		if abiHash != eventsHash {
			log.Error("Event signature mismatch - events will be silently dropped",
				"event", e.eventName,
				"events.go", e.sigConst,
				"contract", ev.Sig,
				"expected_hash", abiHash.Hex(),
				"actual_hash", eventsHash.Hex(),
			)
			mismatches = append(mismatches, fmt.Sprintf(
				"%s: events.go has %q but contract ABI has %q",
				e.eventName, e.sigConst, ev.Sig,
			))
		}
	}

	if len(mismatches) > 0 {
		return fmt.Errorf("event signature validation failed (%d mismatch(es)):\n  %s",
			len(mismatches), strings.Join(mismatches, "\n  "))
	}

	log.Info("All event signatures validated successfully against contract ABIs")
	return nil
}
