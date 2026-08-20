package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/events"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/contracts/ParticipantCoreV1"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/withstack"
)

// ParticipantCoreEventHandler handles events from the ParticipantCore smart contract.
type ParticipantCoreEventHandler struct {
	participantRepo core.ParticipantRepository
	log             logger.Logger
}

// NewParticipantCoreEventHandler creates a new ParticipantCoreEventHandler instance.
func NewParticipantCoreEventHandler(
	participantRepo core.ParticipantRepository,
	log logger.Logger,
) *ParticipantCoreEventHandler {
	return &ParticipantCoreEventHandler{
		participantRepo: participantRepo,
		log:             log,
	}
}

// ContractName returns the contract name this handler processes.
func (h *ParticipantCoreEventHandler) ContractName() string {
	return events.ContractParticipantCore
}

// Handle processes a ParticipantCore contract event by routing it to the appropriate handler method.
func (h *ParticipantCoreEventHandler) Handle(ctx context.Context, log core.ContractLog) error {
	switch log.EventName {
	case events.ParticipantRegistered:
		if err := h.processParticipantRegistered(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.ParticipantUpdated:
		if err := h.processParticipantUpdated(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	default:
		h.log.Debug("No handler for ParticipantCore event", "event", log.EventName)
	}
	return nil
}

// Name returns the handler identifier for logging purposes.
func (h *ParticipantCoreEventHandler) Name() string {
	return "ParticipantCoreHandler"
}

// processParticipantRegistered processes ParticipantRegistered events
func (h *ParticipantCoreEventHandler) processParticipantRegistered(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*ParticipantCoreV1.ParticipantCoreV1ParticipantRegistered](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for ParticipantRegistered: %w", err)
	}

	h.log.Info("ParticipantRegistered event found", "txHash", log.TransactionHash)

	p := event.Participant
	if !p.ChainId.IsUint64() {
		return fmt.Errorf("chainId exceeds uint64 max: %s", p.ChainId.String())
	}

	participantOnChainCreationTime := time.Unix(p.CreatedAt.Int64(), 0).UTC()
	chainId := uint(p.ChainId.Uint64())

	participant := domain.Participant{
		Model: domain.Model{
			CreatedAt: participantOnChainCreationTime,
			UpdatedAt: participantOnChainCreationTime,
		},
		OwnerId:            p.OwnerId,
		ChainId:            &chainId,
		Name:               p.Name,
		Status:             p.Status,
		Role:               p.Role,
		AllowedToBroadcast: p.AllowedToBroadcast,
	}

	participant.StatusStr = domain.MemberStatusToString[int(participant.Status)]
	participant.RoleStr = domain.MemberRoleToString[int(participant.Role)]

	if err := h.participantRepo.Upsert(ctx, participant); err != nil {
		return fmt.Errorf("failed to upsert new participant registration: %w", err)
	}

	h.log.Info("Participant registered successfully",
		"ChainID", fmt.Sprintf("%d", *participant.ChainId),
		"Name", participant.Name,
		"OwnerId", participant.OwnerId)

	return nil
}

// processParticipantUpdated processes ParticipantUpdated events
func (h *ParticipantCoreEventHandler) processParticipantUpdated(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*ParticipantCoreV1.ParticipantCoreV1ParticipantUpdated](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for ParticipantUpdated: %w", err)
	}

	h.log.Info("ParticipantUpdated event found", "txHash", log.TransactionHash)

	p := event.Participant
	if !p.ChainId.IsUint64() {
		return fmt.Errorf("chainId exceeds uint64 max: %s", p.ChainId.String())
	}
	chainId := uint(p.ChainId.Uint64())

	participant := domain.Participant{
		Model:              domain.Model{UpdatedAt: time.Unix(p.UpdatedAt.Int64(), 0).UTC()},
		OwnerId:            p.OwnerId,
		ChainId:            &chainId,
		Name:               p.Name,
		Status:             p.Status,
		Role:               p.Role,
		AllowedToBroadcast: p.AllowedToBroadcast,
	}

	participant.StatusStr = domain.MemberStatusToString[int(participant.Status)]
	participant.RoleStr = domain.MemberRoleToString[int(participant.Role)]

	if err := h.participantRepo.Upsert(ctx, participant); err != nil {
		return fmt.Errorf("failed to upsert participant update: %w", err)
	}

	h.log.Info("Participant updated successfully",
		"ChainID", fmt.Sprintf("%d", *participant.ChainId),
		"Name", participant.Name,
		"OwnerId", participant.OwnerId)

	return nil
}
