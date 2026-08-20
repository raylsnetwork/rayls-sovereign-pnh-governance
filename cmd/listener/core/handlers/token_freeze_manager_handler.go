package handlers

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/events"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TokenFreezeManagerV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

// TokenFreezeManagerEventHandler handles events from the TokenFreezeManager smart contract
type TokenFreezeManagerEventHandler struct {
	tokenFreezeRepo core.TokenFreezeRepository
	provider        core.Provider
	log             logger.Logger
}

// NewTokenFreezeManagerEventHandler creates a new TokenFreezeManagerEventHandler instance
func NewTokenFreezeManagerEventHandler(
	tokenFreezeRepo core.TokenFreezeRepository,
	provider core.Provider,
	log logger.Logger,
) *TokenFreezeManagerEventHandler {
	return &TokenFreezeManagerEventHandler{
		tokenFreezeRepo: tokenFreezeRepo,
		provider:        provider,
		log:             log,
	}
}

// ContractName returns the contract name this handler processes.
func (h *TokenFreezeManagerEventHandler) ContractName() string {
	return events.ContractTokenFreezeManager
}

// Handle processes a TokenFreezeManager contract event by routing it to the appropriate handler method
func (h *TokenFreezeManagerEventHandler) Handle(ctx context.Context, log core.ContractLog) error {
	switch log.EventName {
	case events.TokenFreezeStatusChanged:
		if err := h.processTokenFreezeStatusChanged(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	default:
		h.log.Debug("No handler for TokenFreezeManager event", "event", log.EventName)
	}
	return nil
}

// Name returns the handler identifier for logging purposes
func (h *TokenFreezeManagerEventHandler) Name() string {
	return "TokenFreezeManagerHandler"
}

// processTokenFreezeStatusChanged processes TokenFreezeStatusChanged events
func (h *TokenFreezeManagerEventHandler) processTokenFreezeStatusChanged(
	ctx context.Context,
	log core.ContractLog,
) error {
	// Unmarshal event data to the expected event type
	event, err := core.UnmarshalEventData[*TokenFreezeManagerV1.TokenFreezeManagerV1TokenFreezeStatusChanged](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for TokenFreezeStatusChanged: %w", err)
	}

	resourceIdStr := hex.EncodeToString(event.ResourceId[:])

	h.log.Info("TokenFreezeStatusChanged event found",
		"resourceId", resourceIdStr,
		"action", event.Action,
		"chainIdsCount", len(event.ChainIds))

	// Fetch block to get timestamp
	block, err := h.provider.GetBlockByNumber(ctx, log.BlockNumber)
	if err != nil {
		return withstack.Wrap(err)
	}

	blockTime := block.Time()
	if blockTime > math.MaxInt64 {
		return withstack.Wrap(fmt.Errorf("block timestamp overflows int64"))
	}
	blockTimestamp := time.Unix(int64(blockTime), 0)

	// Process the freeze/unfreeze event
	blockNumber := new(big.Int).SetUint64(log.BlockNumber)
	if err := h.tokenFreezeRepo.UpdateTokenFreezeStatus(
		ctx,
		resourceIdStr,
		event.ChainIds,
		event.Action,
		blockNumber,
		log.TransactionHash,
		blockTimestamp,
	); err != nil {
		return withstack.Wrap(err)
	}

	h.log.Info("Token freeze status processed successfully",
		"resourceId", resourceIdStr,
		"action", types.FreezeAction(event.Action).String(),
		"chainIdsCount", len(event.ChainIds))

	return nil
}
