package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/events"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/DvpTeleport"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

// DvpTeleportEventHandler handles events from the DvpTeleport smart contract.
type DvpTeleportEventHandler struct {
	txRepo    core.TransactionRepository
	decryptor core.Decryptor
	provider  core.Provider
	log       logger.Logger
	pnData    *core.PNodeDataAndSecrets
	swapSalts core.SwapSaltsStore
}

// NewDvpTeleportEventHandler creates a new DvpTeleportEventHandler instance.
func NewDvpTeleportEventHandler(
	txRepo core.TransactionRepository,
	decryptor core.Decryptor,
	provider core.Provider,
	log logger.Logger,
	pnData *core.PNodeDataAndSecrets,
	swapSalts core.SwapSaltsStore,
) *DvpTeleportEventHandler {
	return &DvpTeleportEventHandler{
		txRepo:    txRepo,
		decryptor: decryptor,
		provider:  provider,
		log:       log,
		pnData:    pnData,
		swapSalts: swapSalts,
	}
}

// ContractName returns the contract name this handler processes.
func (h *DvpTeleportEventHandler) ContractName() string {
	return events.ContractDvpTeleport
}

// Handle dispatches DvpTeleport contract logs to the appropriate processor.
func (h *DvpTeleportEventHandler) Handle(ctx context.Context, log core.ContractLog) error {
	switch log.EventName {
	case events.ERCDvpBalanceUpdated:
		if err := h.processERCDvpBalanceUpdatedEvent(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.SwapInitiated:
		if err := h.processSwapInitiatedEvent(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.SwapCompleted:
		if err := h.processSwapCompletedEvent(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.SwapCancelled:
		if err := h.processSwapCancelledEvent(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.SwapTimedOut:
		if err := h.processSwapTimedOutEvent(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	default:
		h.log.Debug("No handler for DvpTeleport event", "event", log.EventName)
	}
	return nil
}

// Name returns the handler identifier for logging purposes.
func (h *DvpTeleportEventHandler) Name() string {
	return "DvpTeleportHandler"
}

// processERCDvpBalanceUpdatedEvent processes ERCDvpBalanceUpdated events.
func (h *DvpTeleportEventHandler) processERCDvpBalanceUpdatedEvent(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*DvpTeleport.DvpTeleportERCDvpBalanceUpdated](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for ERCDvpBalanceUpdated: %w", err)
	}

	h.log.Info("ERCDvpBalanceUpdated event found")

	hubBlock := log.BlockNumber
	hubTxHash := log.TransactionHash
	hubTimestamp := log.BlockTimestamp

	decryptedData, err := core.DecryptPayload[types.DvpBalanceUpdated](
		h.decryptor,
		event.EncryptedMessage,
		event.Raw.BlockNumber,
		*h.pnData,
		types.AtomicSecret,
	)
	if err != nil {
		return withstack.Wrap(err)
	}

	updateTransaction := &domain.Transaction{
		ErcId:                domain.NewBigInt(decryptedData.ErcId),
		MsgType:              decryptedData.TokenType,
		TxType:               decryptedData.UpdateType,
		ResourceId:           decryptedData.ResourceId,
		Protocol:             core.DvpProtocolFromUpdateType(decryptedData.UpdateType),
		SharedId:             decryptedData.SharedId,
		From:                 decryptedData.From,
		To:                   decryptedData.To,
		FromChainId:          domain.NewBigInt(decryptedData.SourceChainId),
		ToChainId:            domain.NewBigInt(decryptedData.DestinationChainId),
		Amount:               decimal.NewFromBigInt(decryptedData.Amount, 0),
		HubTxHash:            hubTxHash,
		HubTimestamp:         hubTimestamp,
		TxHashSource:         decryptedData.SourceTxHash,
		SourceTimestamp:      decryptedData.SourceTxTimestamp,
		TxHashDestination:    hubTxHash,
		DestinationTimestamp: hubTimestamp,
		BlockNumber:          decimal.RequireFromString(strconv.FormatUint(hubBlock, 10)),
		LogIndex:             log.LogIndex,
		AggregationType:      domain.AggregationTypeEnygmaBatch,
		AggregationKey:       uuid.New().String(),
	}

	err = h.txRepo.CreateTransaction(ctx, updateTransaction)
	if err != nil {
		return withstack.Wrap(err)
	}

	h.log.Info("ERCDvpBalanceUpdated event processed")

	return nil
}

// processSwapInitiatedEvent processes SwapInitiated events, emitted when the
// initiator locks their side of the swap. The encrypted payload carries the
// swap details used to materialise the two crosschain transaction rows
// (tokenIn on the initiator's chain, tokenOut on the responder's chain).
func (h *DvpTeleportEventHandler) processSwapInitiatedEvent(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*DvpTeleport.DvpTeleportSwapInitiated](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for SwapInitiated: %w", err)
	}

	// SwapInitiated's Ctxt is encapsulated against the destination-chain
	// participant's view public key, not the operator's. The event doesn't
	// identify which participant that is, so the decryptor tries each
	// candidate DK (same pattern as DecryptPayloadBytes).
	plaintext, initiatorCtxtSalt, err := h.decryptor.DecryptSwapPayload(
		event.Ctxt,
		event.EncryptedData,
		log.BlockNumber,
		*h.pnData,
	)
	if err != nil {
		return withstack.Wrap(err)
	}

	var decryptedData types.DvpSwapMessage
	if err := json.Unmarshal(plaintext, &decryptedData); err != nil {
		return fmt.Errorf("failed to unmarshal SwapInitiated payload: %w", err)
	}

	// Guard all payload fields we dereference below. decimal.NewFromBigInt
	// and domain.NewBigInt both panic on nil, so a malformed payload must be
	// rejected before we build the transaction rows. InitiatorSelfSalt is
	// separately required for decrypting the matching SwapCompleted.
	if err := validateSwapInitiatedPayload(decryptedData); err != nil {
		return withstack.Wrap(err)
	}
	sharedIdKey := hex.EncodeToString(event.SharedId[:])
	if err := h.swapSalts.Put(ctx, sharedIdKey, types.DvpSwapSalts{
		InitiatorSelfSalt: decryptedData.InitiatorSelfSalt.Bytes(),
		InitiatorCtxtSalt: initiatorCtxtSalt,
	}); err != nil {
		return withstack.Wrap(err)
	}

	h.log.Info("SwapInitiated event found", "sharedId", decryptedData.SharedId)

	pendingStatus := uint8(types.DvpSwapStatePending)
	chainId := domain.NewBigInt(decryptedData.ChainId)
	blockNumber := decimal.RequireFromString(strconv.FormatUint(log.BlockNumber, 10))

	baseTx := domain.Transaction{
		Protocol:        types.DvpSwap,
		TeleportStatus:  &pendingStatus,
		SharedId:        decryptedData.SharedId,
		HubTxHash:       log.TransactionHash,
		HubTimestamp:    log.BlockTimestamp,
		SourceTimestamp: decryptedData.PNTxTimestamp,
		TxHashSource:    decryptedData.PNTxHash,
		BlockNumber:     blockNumber,
		LogIndex:        log.LogIndex,
		AggregationType: domain.AggregationTypeDvpSwap,
		AggregationKey:  decryptedData.SharedId,
	}

	// decryptedData.To is the initiator's own address in the SwapInitiated
	// payload (the responder's address arrives later via SwapCompleted).
	// Hence tokenInTx.From and tokenOutTx.To both point at it.

	// TokenIn: being sent from the initiator's chain
	tokenInTx := baseTx
	tokenInTx.ErcId = core.StringToDomainBigInt(decryptedData.TokenInID)
	tokenInTx.MsgType = types.DvpTokenTypeToAssetType(decryptedData.TokenInType)
	tokenInTx.ResourceId = decryptedData.TokenInResourceID
	tokenInTx.Amount = decimal.NewFromBigInt(decryptedData.TokenInAmount, 0)
	tokenInTx.TxType = types.CrossChain
	tokenInTx.FromChainId = chainId
	tokenInTx.From = decryptedData.To

	// TokenOut: being received at the responder's chain
	tokenOutTx := baseTx
	tokenOutTx.ErcId = core.StringToDomainBigInt(decryptedData.TokenOutID)
	tokenOutTx.MsgType = types.DvpTokenTypeToAssetType(decryptedData.TokenOutType)
	tokenOutTx.ResourceId = decryptedData.TokenOutResourceID
	tokenOutTx.Amount = decimal.NewFromBigInt(decryptedData.TokenOutAmount, 0)
	tokenOutTx.TxType = types.CrossChain
	tokenOutTx.ToChainId = chainId
	tokenOutTx.To = decryptedData.To
	tokenOutTx.LogIndex++ // Avoid UNIQUE constraint on (cc_tx_hash, block_number, log_index)

	if err := h.txRepo.CreateTransactions(ctx, []*domain.Transaction{&tokenInTx, &tokenOutTx}); err != nil {
		return withstack.Wrap(err)
	}

	h.log.Info("SwapInitiated event processed", "sharedId", decryptedData.SharedId)
	return nil
}

// processSwapCompletedEvent processes SwapCompleted events, emitted when the
// responder settles the swap. The encrypted payload carries the responder
// side's chainId, address and PN tx reference needed to fill in the
// previously-created transaction rows and marks the swap as completed.
func (h *DvpTeleportEventHandler) processSwapCompletedEvent(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*DvpTeleport.DvpTeleportSwapCompleted](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for SwapCompleted: %w", err)
	}

	sharedIdKey := hex.EncodeToString(event.SharedId[:])

	salts, err := h.swapSalts.Get(ctx, sharedIdKey)
	if err != nil {
		// No stored salts means this listener never processed the matching
		// SwapInitiated (foreign swap, or init occurred before the configured
		// starting block). Ack and skip to mirror the relayer's behavior
		// (`if swap == nil { return nil }`) and avoid a redelivery loop.
		if errors.Is(err, core.ErrSwapSaltsNotFound) {
			h.log.Info("SwapCompleted skipped: no stored salts", "sharedId", sharedIdKey)
			return nil
		}
		return withstack.Wrap(err)
	}

	plaintext, err := h.decryptor.DecryptWithSalt(event.EncryptedData, salts.InitiatorSelfSalt)
	if err != nil {
		return withstack.Wrap(err)
	}

	var decryptedData types.DvpSwapMessage
	if err := json.Unmarshal(plaintext, &decryptedData); err != nil {
		return fmt.Errorf("failed to unmarshal SwapCompleted payload: %w", err)
	}

	h.log.Info("SwapCompleted event found", "sharedId", sharedIdKey)

	if err := h.txRepo.UpdateDvpTeleportConfirmation(
		ctx,
		decryptedData.SharedId,
		decryptedData.ChainId,
		decryptedData.To,
		decryptedData.PNTxHash,
		decryptedData.PNTxTimestamp,
		"",
	); err != nil {
		return fmt.Errorf("failed to update DVP transactions with confirmation data: %w", err)
	}

	if err := h.txRepo.UpdateTeleportStatusBySharedID(
		ctx,
		decryptedData.SharedId,
		uint8(types.DvpSwapStateCompleted),
	); err != nil {
		return fmt.Errorf("failed to update teleport_status: %w", err)
	}

	h.log.Info("SwapCompleted event processed", "sharedId", decryptedData.SharedId)
	return nil
}

// processSwapCancelledEvent processes SwapCancelled events by marking the
// swap as cancelled.
func (h *DvpTeleportEventHandler) processSwapCancelledEvent(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*DvpTeleport.DvpTeleportSwapCancelled](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for SwapCancelled: %w", err)
	}

	sharedId := hex.EncodeToString(event.SharedId[:])
	h.log.Info("SwapCancelled event found", "sharedId", sharedId)

	if err := h.txRepo.UpdateTeleportStatusBySharedID(ctx, sharedId, uint8(types.DvpSwapStateCancelled)); err != nil {
		return fmt.Errorf("failed to update teleport_status: %w", err)
	}

	h.log.Info("SwapCancelled event processed", "sharedId", sharedId)
	return nil
}

// validateSwapInitiatedPayload checks that the decrypted SwapInitiated
// message has every field we later dereference. domain.NewBigInt and
// decimal.NewFromBigInt both panic on a nil *big.Int, so a malformed or
// partial payload would otherwise crash the handler goroutine.
func validateSwapInitiatedPayload(m types.DvpSwapMessage) error {
	switch {
	case m.InitiatorSelfSalt == nil:
		return fmt.Errorf("SwapInitiated payload missing InitiatorSelfSalt")
	case m.ChainId == nil:
		return fmt.Errorf("SwapInitiated payload missing ChainId")
	case m.TokenInAmount == nil:
		return fmt.Errorf("SwapInitiated payload missing TokenInAmount")
	case m.TokenOutAmount == nil:
		return fmt.Errorf("SwapInitiated payload missing TokenOutAmount")
	case m.SharedId == "":
		return fmt.Errorf("SwapInitiated payload missing SharedId")
	case m.TokenInResourceID == "":
		return fmt.Errorf("SwapInitiated payload missing TokenInResourceID")
	case m.TokenOutResourceID == "":
		return fmt.Errorf("SwapInitiated payload missing TokenOutResourceID")
	}
	return nil
}

// processSwapTimedOutEvent processes SwapTimedOut events. Timed-out swaps
// share the cancelled status since the outcome is the same: initiator funds
// are reverted and the swap is terminal.
func (h *DvpTeleportEventHandler) processSwapTimedOutEvent(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*DvpTeleport.DvpTeleportSwapTimedOut](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for SwapTimedOut: %w", err)
	}

	sharedId := hex.EncodeToString(event.SharedId[:])
	h.log.Info("SwapTimedOut event found", "sharedId", sharedId)

	if err := h.txRepo.UpdateTeleportStatusBySharedID(ctx, sharedId, uint8(types.DvpSwapStateCancelled)); err != nil {
		return fmt.Errorf("failed to update teleport_status: %w", err)
	}

	h.log.Info("SwapTimedOut event processed", "sharedId", sharedId)
	return nil
}
