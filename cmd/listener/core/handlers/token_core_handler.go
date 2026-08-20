package handlers

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/events"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TokenCoreV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

// TokenCoreEventHandler handles events from the TokenCore smart contract.
type TokenCoreEventHandler struct {
	txRepo       core.TransactionRepository
	tokenRepo    core.TokenRepository
	tokenService core.TokenService
	log          logger.Logger
}

// MintSupply represents a mint item (ID and amount) for mint transaction persistence.
type MintSupply struct {
	ercID  *big.Int
	amount *big.Int
}

// NewTokenCoreEventHandler creates a new TokenCoreEventHandler instance.
func NewTokenCoreEventHandler(
	txRepo core.TransactionRepository,
	tokenRepo core.TokenRepository,
	tokenService core.TokenService,
	log logger.Logger,
) *TokenCoreEventHandler {
	return &TokenCoreEventHandler{
		txRepo:       txRepo,
		tokenRepo:    tokenRepo,
		tokenService: tokenService,
		log:          log,
	}
}

// ContractName returns the contract name this handler processes.
func (h *TokenCoreEventHandler) ContractName() string {
	return events.ContractTokenCore
}

// Handle processes a TokenCore contract event by routing it to the appropriate handler method.
func (h *TokenCoreEventHandler) Handle(ctx context.Context, log core.ContractLog) error {
	switch log.EventName {
	case events.Erc20TokenRegistered:
		if err := h.processErc20TokenRegistered(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.Erc721TokenRegistered:
		if err := h.processErc721TokenRegistered(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.Erc1155TokenRegistered:
		if err := h.processErc1155TokenRegistered(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.DvpErc721TokenRegistered:
		if err := h.processDvpErc721TokenRegistered(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.DvpErc1155TokenRegistered:
		if err := h.processDvpErc1155TokenRegistered(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.TokenStatusUpdated:
		if err := h.processTokenStatusUpdated(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.TokenBalanceUpdated:
		if err := h.processTokenBalanceUpdated(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	default:
		h.log.Debug("No handler for TokenCore event", "event", log.EventName)
	}
	return nil
}

// Name returns the handler identifier for logging purposes.
func (h *TokenCoreEventHandler) Name() string {
	return "TokenCoreHandler"
}

// processErc20TokenRegistered processes ERC20 token registration events
func (h *TokenCoreEventHandler) processErc20TokenRegistered(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*TokenCoreV1.TokenCoreV1Erc20TokenRegistered](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for Erc20TokenRegistered: %w", err)
	}

	h.log.Info("Erc20TokenRegistered event found")

	if _, err := h.persistTokenRegistry(ctx, event.ResourceId, event.IssuerChainId); err != nil {
		if errors.Is(err, core.ErrIssuerMismatch) {
			h.log.Info("Skipping Erc20TokenRegistered - token belongs to another network",
				"resourceID", hex.EncodeToString(event.ResourceId[:]), "eventIssuer", event.IssuerChainId.String())
			return nil
		}
		return withstack.Wrap(err)
	}

	supplies := []MintSupply{{
		amount: event.InitialSupply,
	}}
	if err := h.persistInitialMintTransactions(
		ctx,
		log,
		event.ResourceId,
		uint8(types.AssetTypeERC20),
		event.IssuerChainId,
		supplies,
	); err != nil {
		return withstack.Wrap(err)
	}

	return nil
}

// processErc721TokenRegistered processes ERC721 token registration events
func (h *TokenCoreEventHandler) processErc721TokenRegistered(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*TokenCoreV1.TokenCoreV1Erc721TokenRegistered](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for Erc721TokenRegistered: %w", err)
	}

	h.log.Info("Erc721TokenRegistered event found")

	if _, err := h.persistTokenRegistry(ctx, event.ResourceId, event.IssuerChainId); err != nil {
		if errors.Is(err, core.ErrIssuerMismatch) {
			h.log.Info("Skipping Erc721TokenRegistered - token belongs to another network",
				"resourceID", hex.EncodeToString(event.ResourceId[:]), "eventIssuer", event.IssuerChainId.String())
			return nil
		}
		return withstack.Wrap(err)
	}

	supplies := make([]MintSupply, len(event.InitialSupply))
	for i, id := range event.InitialSupply {
		supplies[i] = MintSupply{
			ercID:  id,
			amount: big.NewInt(1),
		}
	}
	if err := h.persistInitialMintTransactions(
		ctx,
		log,
		event.ResourceId,
		uint8(types.AssetTypeERC721),
		event.IssuerChainId,
		supplies,
	); err != nil {
		return withstack.Wrap(err)
	}

	return nil
}

// processErc1155TokenRegistered processes ERC1155 token registration events
func (h *TokenCoreEventHandler) processErc1155TokenRegistered(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*TokenCoreV1.TokenCoreV1Erc1155TokenRegistered](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for Erc1155TokenRegistered: %w", err)
	}

	h.log.Info("Erc1155TokenRegistered event found")

	if _, err := h.persistTokenRegistry(ctx, event.ResourceId, event.IssuerChainId); err != nil {
		if errors.Is(err, core.ErrIssuerMismatch) {
			h.log.Info("Skipping Erc1155TokenRegistered - token belongs to another network",
				"resourceID", hex.EncodeToString(event.ResourceId[:]), "eventIssuer", event.IssuerChainId.String())
			return nil
		}
		return withstack.Wrap(err)
	}

	supplies := make([]MintSupply, len(event.InitialSupply))
	for i, supply := range event.InitialSupply {
		supplies[i] = MintSupply{
			ercID:  supply.Id,
			amount: supply.Amount,
		}
	}

	if err := h.persistInitialMintTransactions(
		ctx,
		log,
		event.ResourceId,
		uint8(types.AssetTypeERC1155),
		event.IssuerChainId,
		supplies,
	); err != nil {
		return withstack.Wrap(err)
	}

	return nil
}

// processDvpErc721TokenRegistered processes Dvp ERC721 token registration events
func (h *TokenCoreEventHandler) processDvpErc721TokenRegistered(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*TokenCoreV1.TokenCoreV1DvpErc721TokenRegistered](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for DvpErc721TokenRegistered: %w", err)
	}

	h.log.Info("DvpErc721TokenRegistered event found")

	if _, err := h.persistTokenRegistry(ctx, event.ResourceId, event.IssuerChainId); err != nil {
		if errors.Is(err, core.ErrIssuerMismatch) {
			h.log.Info("Skipping DvpErc721TokenRegistered - token belongs to another network",
				"resourceID", hex.EncodeToString(event.ResourceId[:]), "eventIssuer", event.IssuerChainId.String())
			return nil
		}
		return withstack.Wrap(err)
	}

	supplies := make([]MintSupply, len(event.InitialSupply))
	for i, id := range event.InitialSupply {
		supplies[i] = MintSupply{
			ercID:  id,
			amount: big.NewInt(1),
		}
	}
	if err := h.persistInitialMintTransactions(
		ctx,
		log,
		event.ResourceId,
		uint8(types.AssetTypeDvpERC721),
		event.IssuerChainId,
		supplies,
	); err != nil {
		return withstack.Wrap(err)
	}

	return nil
}

// processDvpErc1155TokenRegistered processes Dvp ERC1155 token registration events
func (h *TokenCoreEventHandler) processDvpErc1155TokenRegistered(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*TokenCoreV1.TokenCoreV1DvpErc1155TokenRegistered](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for DvpErc1155TokenRegistered: %w", err)
	}

	h.log.Info("DvpErc1155TokenRegistered event found")

	if _, err := h.persistTokenRegistry(ctx, event.ResourceId, event.IssuerChainId); err != nil {
		if errors.Is(err, core.ErrIssuerMismatch) {
			h.log.Info("Skipping DvpErc1155TokenRegistered - token belongs to another network",
				"resourceID", hex.EncodeToString(event.ResourceId[:]), "eventIssuer", event.IssuerChainId.String())
			return nil
		}
		return withstack.Wrap(err)
	}

	supplies := make([]MintSupply, len(event.InitialSupply))
	for i, supply := range event.InitialSupply {
		supplies[i] = MintSupply{
			ercID:  supply.Id,
			amount: supply.Amount,
		}
	}
	if err := h.persistInitialMintTransactions(
		ctx,
		log,
		event.ResourceId,
		uint8(types.AssetTypeDvpERC1155),
		event.IssuerChainId,
		supplies,
	); err != nil {
		return withstack.Wrap(err)
	}

	return nil
}

// processTokenStatusUpdated updates the status of a token in the database.
func (h *TokenCoreEventHandler) processTokenStatusUpdated(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*TokenCoreV1.TokenCoreV1TokenStatusUpdated](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for TokenStatusUpdated: %w", err)
	}

	h.log.Info("TokenStatusUpdated event found")

	tokenIssuer := event.IssuerChainId.String()
	tokenName := event.Name

	token, err := h.tokenRepo.GetByIssuerAndName(ctx, tokenIssuer, tokenName)
	if err != nil {
		return withstack.Wrap(err)
	}

	token.Status = uint8(event.Status)

	err = h.tokenRepo.Upsert(ctx, token)
	if err != nil {
		return withstack.Wrap(err)
	}

	h.log.Info("Token status updated successfully",
		"ID", token.ID.String(),
		"name", token.Name,
		"symbol", token.Symbol)

	return nil
}

// processTokenBalanceUpdated handles token balance update events and creates corresponding transactions.
func (h *TokenCoreEventHandler) processTokenBalanceUpdated(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*TokenCoreV1.TokenCoreV1TokenBalanceUpdated](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for TokenBalanceUpdated: %w", err)
	}

	h.log.Info("TokenBalanceUpdated event found")

	token, err := h.persistTokenRegistry(ctx, event.ResourceId, event.IssuerChainId)
	if err != nil {
		if errors.Is(err, core.ErrIssuerMismatch) {
			h.log.Info("Skipping TokenBalanceUpdated - token belongs to another network",
				"resourceID", hex.EncodeToString(event.ResourceId[:]), "eventIssuer", event.IssuerChainId.String())
			return nil
		}
		return withstack.Wrap(err)
	}

	// Normalize ercId, for ERC20 tokens use NULL, not 0
	var ercId *big.Int
	if token.ErcStandard != uint8(types.AssetTypeERC20) ||
		(event.Payload.ErcId != nil && event.Payload.ErcId.Sign() != 0) {
		ercId = event.Payload.ErcId
	}

	updateTransaction := types.UpdateTransaction{
		ResourceId:   event.ResourceId,
		ErcId:        ercId,
		Amount:       event.Payload.Amount,
		UpdateType:   types.TxType(event.UpdateType),
		TxHash:       log.TransactionHash,
		BlockNumber:  strconv.FormatUint(log.BlockNumber, 10),
		HubTimestamp: log.BlockTimestamp,
		MsgType:      token.ErcStandard,
		LogIndex:     log.LogIndex,
	}

	_, err = h.txRepo.PersistUpdateTransaction(ctx, &updateTransaction, token.IssuerId, token.IssuerId)
	if err != nil {
		return withstack.Wrap(err)
	}

	h.log.Info("Token balance updated successfully",
		"ID", token.ID.String(),
		"name", token.Name,
		"symbol", token.Symbol)

	return nil
}

// persistTokenRegistry fetches token metadata from the registry and persists it to the database.
func (h *TokenCoreEventHandler) persistTokenRegistry(
	ctx context.Context,
	eventResourceId [32]byte,
	eventIssuerId *big.Int,
) (*domain.Token, error) {
	resourceIdStr := hex.EncodeToString(eventResourceId[:])
	eventIssuerIdStr := eventIssuerId.String()

	// Get specific token by resource ID from contract registry
	matchingToken, err := h.tokenService.GetTokenByResourceId(ctx, resourceIdStr)
	if err != nil {
		return nil, withstack.Wrap(err)
	}

	// Determine whether the blockchain returned valid token data
	registryHasValidData := matchingToken != nil && matchingToken.IssuerId != "0"

	if registryHasValidData {
		// Verify the issuer matches the event to avoid processing tokens that belong to other networks stored in the shared TokenRegistry
		if matchingToken.IssuerId != eventIssuerIdStr {
			h.log.Info("Token issuer mismatch - skipping (token belongs to another network)",
				"resourceID", resourceIdStr,
				"registryIssuer", matchingToken.IssuerId,
				"eventIssuer", eventIssuerIdStr)
			return nil, core.ErrIssuerMismatch
		}
		// Use the event's resourceId as the canonical DB key in case the registry returns a different representation
		matchingToken.ResourceId = resourceIdStr
	} else {
		h.log.Warn("Token not found in blockchain registry; persisting stub token from event data",
			"resourceID", resourceIdStr,
			"eventIssuer", eventIssuerIdStr)
		matchingToken = &domain.Token{
			ResourceId: resourceIdStr,
			IssuerId:   eventIssuerIdStr,
			Name:       "Unknown",
			Symbol:     "Unknown",
		}
	}

	// Persist the token (upserts on resource_id conflict)
	if err := h.tokenRepo.Upsert(ctx, matchingToken); err != nil {
		return nil, withstack.Wrap(err)
	}

	h.log.Info("Token registered successfully",
		"resourceID", resourceIdStr,
		"name", matchingToken.Name,
		"symbol", matchingToken.Symbol)

	return matchingToken, nil
}

// persistInitialMintTransactions creates mint transaction records for initial token supplies from registration events.
func (h *TokenCoreEventHandler) persistInitialMintTransactions(
	ctx context.Context,
	log core.ContractLog,
	resourceID [32]byte,
	msgType uint8,
	issuerChainID *big.Int,
	supplies []MintSupply,
) error {
	if len(supplies) == 0 {
		return nil
	}

	updateTransactions := make([]types.UpdateTransaction, len(supplies))
	baseLogIndex := log.LogIndex
	issuerChainIdStr := issuerChainID.String()

	for i, supply := range supplies {
		updateTransactions[i] = types.UpdateTransaction{
			ResourceId:   resourceID,
			ErcId:        supply.ercID,
			Amount:       supply.amount,
			UpdateType:   types.Mint,
			TxHash:       log.TransactionHash,
			BlockNumber:  strconv.FormatUint(log.BlockNumber, 10),
			HubTimestamp: log.BlockTimestamp,
			MsgType:      msgType,
			// Increment LogIndex to create unique composite keys for each mint transaction.
			// Initial mints from a single TokenRegistered event share the same base LogIndex,
			// so we add an offset to ensure uniqueness in the (message_id, block_number, log_index) constraint.
			LogIndex: baseLogIndex + uint64(i),
		}
	}

	_, err := h.txRepo.PersistUpdateTransactions(ctx, &updateTransactions, issuerChainIdStr, issuerChainIdStr)
	if err != nil {
		return withstack.Wrap(err)
	}

	h.log.Info("Update transaction(s) persisted successfully",
		"ResourceID", hex.EncodeToString(resourceID[:]),
		"count", len(updateTransactions))

	return nil
}
