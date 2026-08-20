package handlers

import (
	"context"
	"fmt"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/events"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/AuditManagerV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

// AuditManagerEventHandler handles events from the AuditManager smart contract.
type AuditManagerEventHandler struct {
	configProvider core.ConfigProvider
	decryptor      core.Decryptor
	log            logger.Logger
	pnData         *core.PNodeDataAndSecrets
}

// NewAuditManagerEventHandler creates a new AuditManagerEventHandler instance.
// The pnData parameter should be a pointer to the shared PNodeDataAndSecrets that will be updated
// when there are the keys for a participant are updated in the Private Network Hub.
func NewAuditManagerEventHandler(
	configProvider core.ConfigProvider,
	decryptor core.Decryptor,
	log logger.Logger,
	pnData *core.PNodeDataAndSecrets,
) *AuditManagerEventHandler {
	return &AuditManagerEventHandler{
		configProvider: configProvider,
		decryptor:      decryptor,
		log:            log,
		pnData:         pnData,
	}
}

// ContractName returns the contract name this handler processes.
func (h *AuditManagerEventHandler) ContractName() string {
	return events.ContractAuditManager
}

// Handle processes an AuditManager contract event by routing it to the appropriate handler method.
func (h *AuditManagerEventHandler) Handle(ctx context.Context, log core.ContractLog) error {
	switch log.EventName {
	case events.NewAuditOrChainInfo:
		if err := h.processNewAuditOrChainInfo(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	default:
		h.log.Debug("No handler for AuditManager event", "event", log.EventName)
	}
	return nil
}

// Name returns the handler identifier for logging purposes.
func (h *AuditManagerEventHandler) Name() string {
	return "AuditManagerHandler"
}

// processNewAuditOrChainInfo processes NewAuditOrChainInfo events
func (h *AuditManagerEventHandler) processNewAuditOrChainInfo(ctx context.Context, log core.ContractLog) error {
	_, err := core.UnmarshalEventData[*AuditManagerV1.AuditManagerV1NewAuditOrChainInfo](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for NewAuditOrChainInfo: %w", err)
	}

	h.log.Info("NewAuditOrChainInfo event found, refreshing participant data...", "txHash", log.TransactionHash)

	config, err := h.configProvider.GetConfig(ctx)
	if err != nil {
		return withstack.Wrap(err)
	}

	pnData, err := h.decryptor.GatherParticipantsData(ctx, config)
	if err != nil {
		return withstack.Wrap(err)
	}

	// Store the participant secrets in memory for subsequent decryption tasks
	*h.pnData = pnData

	h.log.Info("Participant data gathered and stored successfully", "chainCount", fmt.Sprint(len(pnData)))

	return nil
}
