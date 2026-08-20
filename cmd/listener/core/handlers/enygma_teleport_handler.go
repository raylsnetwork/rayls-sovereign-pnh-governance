package handlers

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/events"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/EnygmaTeleport"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

// EnygmaTeleportEventHandler handles events from the EnygmaTeleport smart contract.
type EnygmaTeleportEventHandler struct {
	txRepo    core.TransactionRepository
	decryptor core.Decryptor
	provider  core.Provider
	log       logger.Logger
	pnData    *core.PNodeDataAndSecrets
}

// NewEnygmaTeleportEventHandler creates a new EnygmaTeleportEventHandler instance.
func NewEnygmaTeleportEventHandler(
	txRepo core.TransactionRepository,
	decryptor core.Decryptor,
	provider core.Provider,
	log logger.Logger,
	pnData *core.PNodeDataAndSecrets,
) *EnygmaTeleportEventHandler {
	return &EnygmaTeleportEventHandler{
		txRepo:    txRepo,
		decryptor: decryptor,
		provider:  provider,
		log:       log,
		pnData:    pnData,
	}
}

// ContractName returns the contract name this handler processes.
func (h *EnygmaTeleportEventHandler) ContractName() string {
	return events.ContractEnygmaTeleport
}

// Handle processes an EnygmaTeleport contract event by routing it to the appropriate handler method.
func (h *EnygmaTeleportEventHandler) Handle(ctx context.Context, log core.ContractLog) error {
	switch log.EventName {
	case events.EnygmaTransfer:
		if err := h.processEnygmaTransferEvent(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.EnygmaTransferCompleted:
		if err := h.processEnygmaTransferCompletedEvent(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.EnygmaSupplyUpdated:
		if err := h.processEnygmaSupplyUpdatedEvent(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.EnygmaDvpBalanceUpdated:
		if err := h.processEnygmaDvpBalanceUpdated(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	default:
		h.log.Debug("No handler for EnygmaTeleport event", "event", log.EventName)
	}
	return nil
}

// Name returns the handler identifier for logging purposes.
func (h *EnygmaTeleportEventHandler) Name() string {
	return "EnygmaTeleportHandler"
}

// processEnygmaTransferEvent processes EnygmaTransfer events
func (h *EnygmaTeleportEventHandler) processEnygmaTransferEvent(ctx context.Context, log core.ContractLog) error {
	// Unmarshal event data to the expected event type
	event, err := core.UnmarshalEventData[*EnygmaTeleport.EnygmaTeleportEnygmaTransfer](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for EnygmaTransfer: %w", err)
	}

	if len(event.EncryptedMessage) == 0 {
		h.log.Debug("No encrypted data in TeleportEnygmaTransfer event")
		return nil
	}

	h.log.Info("EnygmaTransfer event found")

	batch, err := core.DecryptPayload[types.EnygmaTransferBatch](
		h.decryptor,
		event.EncryptedMessage,
		log.BlockNumber,
		*h.pnData,
		types.ParticipantSecret,
	)
	if err != nil {
		return withstack.Wrap(err)
	}

	if len(batch.Transactions) == 0 {
		h.log.Debug("No valid Enygma transactions to persist")
		return nil
	}

	// Convert resource ID from batch
	var resourceId [32]byte
	copy(resourceId[:], common.Hex2Bytes(batch.ResourceId))

	transactions := make([]*domain.Transaction, 0, len(batch.Transactions))
	enygmaTransactions := make([]*domain.EnygmaTransaction, 0, len(batch.Transactions))

	// Process each transaction in the batch
	// Enygma transactions are always batches, using batch_id as aggregation key
	for _, tx := range batch.Transactions {
		h.log.Debug("EnygmaTransfer event will be persisted", "referenceId", hex.EncodeToString(tx.ReferenceId[:]))

		pendingStatus := uint8(types.EnygmaTransferPending)
		transaction := &domain.Transaction{
			MessageId:       tx.MessageId,
			BatchId:         batch.BatchId,
			ResourceId:      hex.EncodeToString(resourceId[:]),
			Protocol:        types.Enygma,
			FromChainId:     domain.NewBigInt(batch.FromChainID),
			ToChainId:       domain.NewBigInt(batch.ToChainID),
			From:            tx.FromAddress.String(),
			To:              tx.ToAddress.String(),
			Amount:          decimal.NewFromBigInt(tx.ToAmount, 0),
			TxType:          types.CrossChain,
			MsgType:         uint8(types.AssetTypeEnygma),
			HubTxHash:       log.TransactionHash,
			HubTimestamp:    log.BlockTimestamp,
			BlockNumber:     decimal.RequireFromString(strconv.FormatUint(log.BlockNumber, 10)),
			LogIndex:        log.LogIndex,
			TeleportStatus:  &pendingStatus,
			AggregationType: domain.AggregationTypeEnygmaBatch,
			AggregationKey:  batch.BatchId,
		}

		enygmaTransaction := &domain.EnygmaTransaction{
			ReferenceId:   hex.EncodeToString(tx.ReferenceId[:]),
			ToRValueToAdd: domain.NewBigInt(batch.ToRValueToAdd),
		}

		transactions = append(transactions, transaction)
		enygmaTransactions = append(enygmaTransactions, enygmaTransaction)

		h.log.Debug("EnygmaTransfer entry prepared for persistence", "messageId", transaction.MessageId)
	}

	if err = h.txRepo.CreateTransactionsWithEnygmaData(ctx, transactions, enygmaTransactions); err != nil {
		return withstack.Wrap(err)
	}

	h.log.Info("EnygmaTransfer event processed", "count", len(transactions))
	return nil
}

// processEnygmaTransferCompletedEvent processes EnygmaTransferCompleted events
func (h *EnygmaTeleportEventHandler) processEnygmaTransferCompletedEvent(
	ctx context.Context,
	log core.ContractLog,
) error {
	// Unmarshal event data to the expected event type
	event, err := core.UnmarshalEventData[*EnygmaTeleport.EnygmaTeleportEnygmaTransferCompleted](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for EnygmaTransferCompleted: %w", err)
	}

	h.log.Info("EnygmaTransferCompleted event found")

	if len(event.EncryptedMessage) == 0 {
		h.log.Debug("No encrypted data in TeleportEnygmaTransferCompleted event")
		return nil
	}

	decryptedData, err := core.DecryptPayload[[]types.EnygmaTransferCompleted](
		h.decryptor,
		event.EncryptedMessage,
		log.BlockNumber,
		*h.pnData,
		types.ParticipantSecret,
	)
	if err != nil {
		return withstack.Wrap(err)
	}

	// Skip if no decrypted data
	if len(decryptedData) == 0 {
		h.log.Warn("Decrypted EnygmaTransferCompleted data is empty")
		return nil
	}

	h.log.Debug("Successfully decrypted AtomicMessageAdditionalDataBatch",
		"decrypted_records", fmt.Sprint(len(decryptedData)))

	transactionMap, err := core.BuildTransactionMap(
		ctx,
		h.txRepo,
		decryptedData,
		func(data types.EnygmaTransferCompleted) string { return data.MessageId },
		core.MessageID,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch transactions by message IDs: %w", err)
	}

	transactionsToUpdate := make([]*domain.Transaction, 0, len(transactionMap))
	for _, data := range decryptedData {
		transaction, exists := transactionMap[data.MessageId]
		if !exists {
			h.log.Warn("Transaction not found for EnygmaTransferCompleted", "messageId", data.MessageId)
			continue
		}

		executedStatus := uint8(types.EnygmaTransferExecuted)
		transaction.TxHashDestination = data.TransactionHash
		transaction.TeleportStatus = &executedStatus
		transactionsToUpdate = append(transactionsToUpdate, transaction)
	}

	if err = h.txRepo.UpdateTransactionsBulk(
		ctx,
		transactionsToUpdate,
		"tx_hash_destination",
		"teleport_status",
		"updated_at",
	); err != nil {
		return withstack.Wrap(err)
	}

	h.log.Info("EnygmaTransferCompleted processed successfully",
		"count", fmt.Sprint(len(transactionsToUpdate)))

	return nil
}

// processEnygmaSupplyUpdatedEvent processes EnygmaSupplyUpdated events
func (h *EnygmaTeleportEventHandler) processEnygmaSupplyUpdatedEvent(ctx context.Context, log core.ContractLog) error {
	// Unmarshal event data to the expected event type
	event, err := core.UnmarshalEventData[*EnygmaTeleport.EnygmaTeleportEnygmaSupplyUpdated](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for EnygmaSupplyUpdated: %w", err)
	}

	h.log.Info("EnygmaSupplyUpdated event found")

	// Extract data from the event
	resourceId := hex.EncodeToString(event.ResourceId[:])
	blockNumber := decimal.NewFromBigInt(event.BlockNumber, 0)
	amount := decimal.NewFromBigInt(event.Update.Amount, 0)
	chainId := event.ChainId
	txType := types.Mint

	// Adjust amount based on transaction type:
	// Type 1 (mint) = positive value
	// Type 2 (burn) = negative value
	if event.Update.TxType == uint8(types.EnygmaTxTypeBurn) {
		amount = amount.Neg()
		txType = types.Burn
	}

	txIdentifier := uuid.New()

	enygmaTx := &domain.Transaction{
		Model:           domain.Model{ID: txIdentifier},
		ResourceId:      resourceId,
		FromChainId:     domain.NewBigInt(chainId),
		ToChainId:       domain.NewBigInt(chainId),
		BlockNumber:     blockNumber,
		Amount:          amount,
		TxType:          txType,
		MsgType:         uint8(types.AssetTypeEnygma),
		HubTxHash:       log.TransactionHash,
		HubTimestamp:    log.BlockTimestamp,
		LogIndex:        log.LogIndex,
		AggregationType: domain.AggregationTypeTransaction,
		// Use UUID as key since no message_id
		AggregationKey: txIdentifier.String(),
	}

	if err := h.txRepo.CreateTransaction(ctx, enygmaTx); err != nil {
		return withstack.Wrap(err)
	}

	h.log.Info("EnygmaSupplyUpdated event processed successfully", "resourceID", resourceId)

	return nil
}

// processEnygmaDvpBalanceUpdated processes EnygmaDvpBalanceUpdated events
func (h *EnygmaTeleportEventHandler) processEnygmaDvpBalanceUpdated(ctx context.Context, log core.ContractLog) error {
	event, err := core.UnmarshalEventData[*EnygmaTeleport.EnygmaTeleportEnygmaDvpBalanceUpdated](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for EnygmaDvpBalanceUpdated: %w", err)
	}

	h.log.Info("EnygmaDvpBalanceUpdated event found")

	hubBlock := log.BlockNumber
	hubTxHash := log.TransactionHash
	hubTimestamp := log.BlockTimestamp

	decryptedData, err := core.DecryptPayload[types.DvpBalanceUpdated](
		h.decryptor,
		event.EncryptedMessage,
		hubBlock,
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

	h.log.Info("EnygmaDvpBalanceUpdated event processed")

	return nil
}
