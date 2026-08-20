package handlers

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/events"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TeleportV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

// TeleportEventHandler handles events from the Teleport smart contract.
type TeleportEventHandler struct {
	txRepo     core.TransactionRepository
	revertRepo core.RevertDataTransactionRepository
	decryptor  core.Decryptor
	log        logger.Logger
	pnData     *core.PNodeDataAndSecrets
}

// NewTeleportEventHandler creates a new TeleportEventHandler instance.
func NewTeleportEventHandler(
	txRepo core.TransactionRepository,
	revertRepo core.RevertDataTransactionRepository,
	decryptor core.Decryptor,
	log logger.Logger,
	pnData *core.PNodeDataAndSecrets,
) *TeleportEventHandler {
	return &TeleportEventHandler{
		txRepo:     txRepo,
		revertRepo: revertRepo,
		decryptor:  decryptor,
		log:        log,
		pnData:     pnData,
	}
}

// ContractName returns the contract name this handler processes.
func (h *TeleportEventHandler) ContractName() string {
	return events.ContractTeleport
}

// Handle processes a Teleport contract event by routing it to the appropriate handler method.
func (h *TeleportEventHandler) Handle(ctx context.Context, log core.ContractLog) error {
	switch log.EventName {
	case events.AtomicMessageAdditionalDataBatch:
		if err := h.processAtomicMessageAdditionalDataBatch(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.AtomicMessageStatusChangedBatch:
		if err := h.processAtomicMessageStatusChangedBatch(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	case events.EncryptedDataBatchStored:
		if err := h.processEncryptedDataBatchStored(ctx, log); err != nil {
			return withstack.Wrap(err)
		}
	default:
		h.log.Debug("No handler for Teleport event", "event", log.EventName)
	}
	return nil
}

// Name returns the handler identifier for logging purposes.
func (h *TeleportEventHandler) Name() string {
	return "TeleportHandler"
}

// processAtomicMessageAdditionalDataBatch processes AtomicMessageAdditionalDataBatch events
func (h *TeleportEventHandler) processAtomicMessageAdditionalDataBatch(
	ctx context.Context,
	log core.ContractLog,
) error {
	// Unmarshal to the expected event type
	event, err := core.UnmarshalEventData[*TeleportV1.TeleportV1AtomicMessageAdditionalDataBatch](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for AtomicMessageAdditionalDataBatch: %w", err)
	}

	if len(event.EncryptedData) == 0 {
		h.log.Debug("No encrypted data in AtomicMessageAdditionalDataBatch event")
		return nil
	}

	decodedPayload, err := hex.DecodeString(event.EncryptedData)
	if err != nil {
		return fmt.Errorf("failed to decode hex for EncryptedData: %w", err)
	}

	decryptedData, err := core.DecryptPayload[[]types.AtomicTeleportAdditionalData](
		h.decryptor,
		decodedPayload,
		log.BlockNumber,
		*h.pnData,
		types.AtomicSecret,
	)
	if err != nil {
		return withstack.Wrap(err)
	}

	// Skip if no decrypted data
	if len(decryptedData) == 0 {
		h.log.Warn("Decrypted atomic data is empty")
		return nil
	}

	h.log.Debug("Successfully decrypted AtomicMessageAdditionalDataBatch",
		"decrypted_records", len(decryptedData),
		"message_ids", len(event.MsgIds))

	// Use the helper to extract shared IDs and create transaction map
	transactionMap, err := core.BuildTransactionMap(
		ctx,
		h.txRepo,
		decryptedData,
		func(data types.AtomicTeleportAdditionalData) string { return data.SharedId },
		core.SharedID,
	)
	if err != nil {
		return withstack.Wrap(err)
	}

	transactionsToUpdate := make([]*domain.Transaction, 0, len(decryptedData))
	var revertTransactionsToCreate []*domain.RevertDataTransaction
	for _, data := range decryptedData {
		transaction, exists := transactionMap[data.SharedId]
		if !exists {
			h.log.Warn("Transaction not found", "shared_id", data.SharedId)
			continue
		}

		applyTransactionUpdates(transaction, data)

		// Build and persist revert data if applicable
		if revert := buildRevertDataIfPresent(data); revert != nil {
			revert.TransactionId = transaction.ID
			revertTransactionsToCreate = append(revertTransactionsToCreate, revert)
		}

		transactionsToUpdate = append(transactionsToUpdate, transaction)
	}

	if err = h.txRepo.UpdateTransactionsBulk(
		ctx,
		transactionsToUpdate,
		"tx_hash_destination",
		"destination_timestamp",
		"hub_tx_hash",
		"updated_at",
	); err != nil {
		return withstack.Wrap(err)
	}

	if len(revertTransactionsToCreate) > 0 {
		if err = h.revertRepo.CreateRevertTransactions(ctx, revertTransactionsToCreate); err != nil {
			return withstack.Wrap(err)
		}
	}

	h.log.Info("AtomicMessageAdditionalDataBatch event processed successfully",
		"total_updated_transactions", len(transactionsToUpdate))

	return nil
}

// processAtomicMessageStatusChangedBatch processes AtomicMessageStatusChangedBatch events
func (h *TeleportEventHandler) processAtomicMessageStatusChangedBatch(ctx context.Context, log core.ContractLog) error {
	// Unmarshal to the expected event type
	event, err := core.UnmarshalEventData[*TeleportV1.TeleportV1AtomicMessageStatusChangedBatch](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for AtomicMessageStatusChangedBatch: %w", err)
	}

	h.log.Info("AtomicMessageStatusChangedBatch event found")

	if len(event.MsgIds) == 0 {
		h.log.Debug("No message IDs in AtomicMessageStatusChangedBatch event")
		return nil
	}

	status := types.AtomicTeleportStatus(event.Status)

	transactionMap, err := core.BuildTransactionMap(
		ctx,
		h.txRepo,
		event.MsgIds,
		func(id string) string { return id },
		core.SharedID,
	)
	if err != nil {
		return fmt.Errorf("failed to retrieve transactions from database: %w", err)
	}

	transactionsToUpdate := make([]*domain.Transaction, 0, len(transactionMap))
	for _, msgId := range event.MsgIds {
		transaction, exists := transactionMap[msgId]
		if !exists || transaction.TeleportStatus == nil ||
			*transaction.TeleportStatus != uint8(types.AtomicTeleportPending) {
			h.log.Warn("Transaction not found or invalid for AtomicMessageStatusChangedBatch", "messageId", msgId)
			continue
		}
		statusValue := uint8(status)
		transaction.TeleportStatus = &statusValue
		transactionsToUpdate = append(transactionsToUpdate, transaction)
		if status == types.AtomicTeleportRejected {
			if createErr := h.txRepo.CreateCreditMemo(ctx, transaction); createErr != nil {
				return fmt.Errorf("failed to create credit memo: %w", createErr)
			}
		}
	}

	err = h.txRepo.UpdateTransactionsBulk(ctx, transactionsToUpdate, "teleport_status", "updated_at")
	if err != nil {
		return withstack.Wrap(err)
	}

	h.log.Info("AtomicMessageStatusChangedBatch event processed successfully",
		"total_updated_transactions", len(transactionsToUpdate))

	return nil
}

// processEncryptedDataBatchStored processes EncryptedDataBatchStored events
func (h *TeleportEventHandler) processEncryptedDataBatchStored(ctx context.Context, log core.ContractLog) error {
	// Unmarshal to the expected event type
	event, err := core.UnmarshalEventData[*TeleportV1.TeleportV1EncryptedDataBatchStored](log)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event data for EncryptedDataBatchStored: %w", err)
	}

	h.log.Info("EncryptedDataBatchStored event found")

	decryptedBatch, err := core.DecryptPayload[[]types.DispatchedMessageToPrivateHub](
		h.decryptor,
		event.Data,
		log.BlockNumber,
		*h.pnData,
		types.ParticipantSecret,
	)
	if err != nil {
		return withstack.Wrap(err)
	}

	if len(decryptedBatch) == 0 {
		h.log.Warn("EncryptedDataBatch: Decrypted batch is empty")
		return nil
	}

	h.log.Debug("Processed EncryptedDataBatchStored",
		"total_decrypted", fmt.Sprint(len(decryptedBatch)),
		"block_number", strconv.FormatUint(log.BlockNumber, 10))

	// Persist transfer messages to database (filtering happens inside)
	if err := h.persistTransferMessages(
		ctx,
		decryptedBatch,
		log.BlockNumber,
		log.TransactionHash,
		log.BlockTimestamp,
	); err != nil {
		return withstack.Wrap(err)
	}

	return nil
}

// persistTransferMessages filters and converts transfer messages to transactions, then persists them
func (h *TeleportEventHandler) persistTransferMessages(
	ctx context.Context,
	batch []types.DispatchedMessageToPrivateHub,
	blockNumber uint64,
	hash string,
	blockTimestamp time.Time,
) error {
	transactions := make([]*domain.Transaction, 0, len(batch))
	batchIds := make([]string, 0)
	seenBatchIds := make(map[string]bool)

	for _, msg := range batch {
		// Skip non-transfer messages or DVP-related arbitrary messages
		if msg.TransactionType != types.Transfer ||
			msg.Data.MessageMetadata.ResourceId == types.ResourceIDPNCommunicator {
			continue
		}

		h.log.Debug("Processing transfer message", "shared-id", msg.SharedId)

		tx, err := h.convertToTransaction(&msg, strconv.FormatUint(blockNumber, 10), hash, blockTimestamp)
		if err != nil {
			h.log.Error("Failed to convert message to transaction",
				"shared_id", msg.SharedId,
				"error", err)
			continue
		}
		transactions = append(transactions, tx)

		// Collect unique batch IDs
		if msg.BatchId != "" && !seenBatchIds[msg.BatchId] {
			batchIds = append(batchIds, msg.BatchId)
			seenBatchIds[msg.BatchId] = true
		}
	}

	if len(transactions) == 0 {
		h.log.Debug("No transfer messages to persist")
		return nil
	}

	// Create transactions and promote batches atomically
	if err := h.txRepo.CreateTransactionsWithPromotion(ctx, transactions, batchIds); err != nil {
		return err
	}

	h.log.Info("Successfully persisted transfer messages", "count", len(transactions))
	return nil
}

// convertToTransaction converts a DispatchedMessageToPrivateHub to a domain.Transaction
func (h *TeleportEventHandler) convertToTransaction(
	msg *types.DispatchedMessageToPrivateHub,
	blockNumber,
	hash string,
	hubTimestamp time.Time,
) (*domain.Transaction, error) {
	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}

	msgAmount := msg.Data.MessageMetadata.TransferMetadata.Amount
	var amount decimal.Decimal
	if msgAmount == nil {
		amount = decimal.Zero
	} else {
		amount = decimal.NewFromBigInt(msgAmount, 0)
	}

	blockNumberDecimal, err := decimal.NewFromString(blockNumber)
	if err != nil {
		return nil, err
	}

	// Set erc_id only for ERC721/ERC1155 assets. For ERC20 and others, keep it nil to avoid misleading 0 values
	var ercId *big.Int
	if msg.Data.MessageMetadata.TransferMetadata.AssetType == types.AssetTypeERC1155 ||
		msg.Data.MessageMetadata.TransferMetadata.AssetType == types.AssetTypeERC721 {
		ercId = msg.Data.MessageMetadata.TransferMetadata.Id
	}

	resourceId := hex.EncodeToString(msg.Data.MessageMetadata.ResourceId[:])
	messageId := hex.EncodeToString(msg.MessageId[:])
	protocolType, teleportStatus := h.classifyProtocol(msg)

	// Determine aggregation type and key
	// For standalone transactions, always use message_id as aggregation_key
	// The PromoteBatchToRegular call will update aggregation_key to batch_id when the batch reaches 2+ transactions
	aggregationType := domain.AggregationTypeTransaction
	aggregationKey := messageId

	transaction := &domain.Transaction{
		ToChainId:         domain.NewBigInt(msg.ToChainId),
		FromChainId:       domain.NewBigInt(msg.FromChainId),
		ResourceId:        resourceId,
		BlockNumber:       blockNumberDecimal,
		HubTxHash:         hash,
		HubTimestamp:      hubTimestamp,
		Amount:            amount,
		TxType:            types.CrossChain,
		MsgType:           uint8(msg.Data.MessageMetadata.TransferMetadata.AssetType),
		MessageId:         messageId,
		From:              msg.Data.MessageMetadata.TransferMetadata.From,
		To:                msg.Data.MessageMetadata.TransferMetadata.To,
		Protocol:          protocolType,
		SharedId:          msg.SharedId,
		BatchId:           msg.BatchId,
		LogIndex:          uint64(msg.LogIdx),
		TeleportStatus:    teleportStatus,
		Payload:           hex.EncodeToString(msg.Data.Payload),
		TxHashDestination: msg.TxHashDestination.String(),
		ErcId:             domain.NewBigInt(ercId),
		SourceTimestamp:   core.SafeUnixTime(msg.TxHashSourceTimestamp),
		AggregationType:   aggregationType,
		AggregationKey:    aggregationKey,
	}

	return transaction, nil
}

// classifyProtocol classifies the protocol and returns appropriate TeleportStatus
func (h *TeleportEventHandler) classifyProtocol(
	msg *types.DispatchedMessageToPrivateHub,
) (types.ProtocolType, *uint8) {
	if msg.IsAtomic {
		pendingStatus := uint8(types.AtomicTeleportPending)
		return types.Atomic, &pendingStatus
	}
	// Vanilla protocol doesn't have TeleportStatus
	return types.Vanilla, nil
}

// applyTransactionUpdates updates transaction fields
func applyTransactionUpdates(tx *domain.Transaction, data types.AtomicTeleportAdditionalData) {
	emptyHash := common.Hash{}

	// Destination unlock tx
	if data.TxHashDestination != emptyHash {
		tx.TxHashDestination = data.TxHashDestination.String()
	}
	if data.TxHashDestinationTimestamp != 0 {
		tx.DestinationTimestamp = core.SafeUnixTime(data.TxHashDestinationTimestamp)
	}

	// Hub hash
	if data.BatchHubTxHash != emptyHash {
		tx.HubTxHash = data.BatchHubTxHash.String()
	}
}

// buildRevertDataIfPresent constructs a RevertDataTransaction if any revert fields are present.
// Returns nil if no revert data exists.
func buildRevertDataIfPresent(data types.AtomicTeleportAdditionalData) *domain.RevertDataTransaction {
	emptyHash := common.Hash{}
	revert := &domain.RevertDataTransaction{}

	// Destination revert
	if data.TxHashDestinationRevert != emptyHash {
		revert.TxHashDestinationRevert = data.TxHashDestinationRevert.String()
	}
	if data.TxHashDestinationRevertStatus > 0 {
		revert.TxHashDestinationRevertStatus = uint8(data.TxHashDestinationRevertStatus)
	}

	// Source revert
	if data.TxHashSourceRevert != emptyHash {
		revert.TxHashSourceRevert = data.TxHashSourceRevert.String()
	}
	if data.TxHashSourceRevertStatus > 0 {
		revert.TxHashSourceRevertStatus = uint8(data.TxHashSourceRevertStatus)
	}

	// Only return revert record if at least one revert field was set
	hasRevert := data.TxHashDestinationRevert != emptyHash ||
		data.TxHashDestinationRevertStatus != 0 ||
		data.TxHashSourceRevert != emptyHash ||
		data.TxHashSourceRevertStatus != 0

	if hasRevert {
		return revert
	}

	return nil
}
