package handlers

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/events"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/contracts/EnygmaTokenManagerV1"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/withstack"
)

// ErrEnygmaTokenNotMatched indicates that a token was found but its issuer doesn't match the event
var ErrEnygmaTokenNotMatched = errors.New("no matching token found for registration")

// EnygmaTokenManagerEventHandler handles events from the EnygmaTokenManager smart contract.
type EnygmaTokenManagerEventHandler struct {
	tokenRepo    core.TokenRepository
	tokenService core.TokenService
	log          logger.Logger
}

// NewEnygmaTokenManagerEventHandler creates a new EnygmaTokenManagerEventHandler instance.
func NewEnygmaTokenManagerEventHandler(
	tokenRepo core.TokenRepository,
	tokenService core.TokenService,
	log logger.Logger,
) *EnygmaTokenManagerEventHandler {
	return &EnygmaTokenManagerEventHandler{
		tokenRepo:    tokenRepo,
		tokenService: tokenService,
		log:          log,
	}
}

// ContractName returns the contract name this handler processes.
func (h *EnygmaTokenManagerEventHandler) ContractName() string {
	return events.ContractEnygmaTokenManager
}

// Handle processes an EnygmaTokenManager contract event by routing it to the appropriate handler method.
func (h *EnygmaTokenManagerEventHandler) Handle(ctx context.Context, log core.ContractLog) error {
	switch log.EventName {
	case events.EnygmaTokenRegistered:
		if err := h.processEnygmaTokenRegistered(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	default:
		h.log.Debug("No handler for EnygmaTokenManager event", "event", log.EventName)
	}
	return nil
}

// Name returns the handler identifier for logging purposes.
func (h *EnygmaTokenManagerEventHandler) Name() string {
	return "EnygmaTokenManagerHandler"
}

// processEnygmaTokenRegistered processes Enygma token registration events
func (h *EnygmaTokenManagerEventHandler) processEnygmaTokenRegistered(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*EnygmaTokenManagerV1.EnygmaTokenManagerV1EnygmaTokenRegistered](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for EnygmaTokenRegistered: %w", err)
	}

	h.log.Info("EnygmaTokenRegistered event found")

	eventResourceId := hex.EncodeToString(event.ResourceId[:])
	eventIssuerId := event.IssuerChainId.String()

	_, err = h.processTokenRegistered(ctx, eventResourceId, eventIssuerId)
	if errors.Is(err, core.ErrIssuerMismatch) {
		h.log.Info("Skipping EnygmaTokenRegistered - token belongs to another network",
			"resourceID", eventResourceId, "eventIssuer", eventIssuerId)
		return nil
	}
	return err
}

// processTokenRegistered handles the core token registration logic.
func (h *EnygmaTokenManagerEventHandler) processTokenRegistered(
	ctx context.Context,
	eventResourceId string,
	eventIssuerId string,
) (*domain.Token, error) {
	// Get specific token by resource ID from contract registry
	matchingToken, err := h.tokenService.GetTokenByResourceId(ctx, eventResourceId)
	if err != nil {
		return nil, withstack.Wrap(err)
	}

	h.log.Info("TokenRegistered event found")

	// Determine whether the blockchain returned valid token data
	registryHasValidData := matchingToken != nil && matchingToken.IssuerId != "0"

	if registryHasValidData {
		// Token found in registry: verify the issuer matches the event to avoid processing tokens that belong to other networks stored in the shared TokenRegistry.
		if matchingToken.IssuerId != eventIssuerId {
			h.log.Info("Token issuer mismatch - skipping (token belongs to another network)",
				"resourceID", eventResourceId,
				"registryIssuer", matchingToken.IssuerId,
				"eventIssuer", eventIssuerId)
			return nil, core.ErrIssuerMismatch
		}
		// Use the event's resourceId as the canonical DB key.
		matchingToken.ResourceId = eventResourceId
	} else {
		// Persist a stub token from event data to ensure consistency with transaction records
		h.log.Warn("Token not found in blockchain registry; persisting stub token from event data",
			"resourceID", eventResourceId,
			"eventIssuer", eventIssuerId)
		matchingToken = &domain.Token{
			ResourceId: eventResourceId,
			IssuerId:   eventIssuerId,
			Name:       "Unknown",
			Symbol:     "Unknown",
		}
	}

	// Persist the token (upserts on resource_id conflict)
	if err := h.tokenRepo.Upsert(ctx, matchingToken); err != nil {
		return nil, withstack.Wrap(err)
	}

	h.log.Info("Token registered successfully",
		"resourceID", eventResourceId,
		"name", matchingToken.Name,
		"symbol", matchingToken.Symbol)

	return matchingToken, nil
}
