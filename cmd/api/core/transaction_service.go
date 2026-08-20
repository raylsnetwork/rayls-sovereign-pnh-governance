package core

import (
	"context"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/utils"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

var _ TransactionService = (*transactionService)(nil)

// transactionService implements the TransactionService interface
type transactionService struct {
	txRepo          TransactionRepository
	metadataService TokenMetadataService
	mapper          *TransactionMapper
	log             logger.Logger
}

// NewTransactionService creates a new transaction service
func NewTransactionService(
	txRepo TransactionRepository,
	metadataService TokenMetadataService,
	log logger.Logger,
) *transactionService {
	return &transactionService{
		txRepo:          txRepo,
		metadataService: metadataService,
		mapper:          NewTransactionMapper(metadataService),
		log:             log,
	}
}

// GetTransactionByMessageId retrieves a single transaction by messageId
func (s *transactionService) GetTransactionByMessageId(
	ctx context.Context,
	messageId string,
) (*dto.TransactionDetailDto, error) {
	if messageId == "" {
		return nil, NewValidationError("messageId", "messageId cannot be empty")
	}

	transaction, err := s.txRepo.FindByMessageId(ctx, messageId)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, NewNotFoundError("transaction", messageId)
		}
		return nil, NewInternalError("FindByMessageId", err)
	}

	return s.mapper.ToTransactionDetailDto(ctx, *transaction, "transaction")
}

// GetTransactionByTransactionId retrieves a single transaction by transactionId
func (s *transactionService) GetTransactionByTransactionId(
	ctx context.Context,
	transactionId string,
) (*dto.TransactionDetailDto, error) {
	if transactionId == "" {
		return nil, NewValidationError("transactionId", "transactionId cannot be empty")
	}

	transaction, err := s.txRepo.FindByTransactionId(ctx, transactionId)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, NewNotFoundError("transaction", transactionId)
		}
		return nil, NewInternalError("FindByTransactionId", err)
	}

	return s.mapper.ToTransactionDetailDto(ctx, *transaction, "transaction")
}

// GetTransactionsList retrieves a paginated list of transactions
func (s *transactionService) GetTransactionsList(
	ctx context.Context,
	filters dto.MergedTransactionsFilters,
) (*types.Paginated[dto.TransactionListDto], error) {
	// Validate and sanitize filters
	if err := s.validateFilters(&filters); err != nil {
		return nil, err
	}

	transactions, total, err := s.txRepo.FindWithFilters(ctx, filters)
	if err != nil {
		return nil, NewInternalError("FindWithFilters", err)
	}

	s.log.Debug("Fetched aggregated transactions from database",
		"count", len(transactions),
		"total", total,
		"page", filters.Page,
		"limit", filters.Limit)

	dtos := s.mapTransactionsToListDtos(transactions)

	result := &types.Paginated[dto.TransactionListDto]{
		Data:  dtos,
		Total: total,
		Limit: filters.Limit,
		Page:  filters.Page,
	}

	return result, nil
}

// GetFlaggedTransactions retrieves all flagged transactions
func (s *transactionService) GetFlaggedTransactions(ctx context.Context) ([]domain.FlaggedTransaction, error) {
	s.log.Info("Fetching all flagged transactions")

	// Call repository to get flagged transactions
	flagged, err := s.txRepo.FindFlagged(ctx)
	if err != nil {
		s.log.Error("Failed to fetch flagged transactions", "error", err)
		return nil, err
	}

	s.log.Info("Successfully fetched flagged transactions", "count", len(flagged))
	return flagged, nil
}

// mapTransactionsToListDtos converts domain transactions to list DTOs
func (s *transactionService) mapTransactionsToListDtos(transactions []domain.Transaction) []dto.TransactionListDto {
	dtos := make([]dto.TransactionListDto, 0, len(transactions))

	for _, tx := range transactions {
		idType, id := determineIdAndType(tx)

		dto := dto.TransactionListDto{
			Type:          tx.GetProtocolString(),
			Id:            id,
			IdType:        idType,
			CreatedAt:     tx.CreatedAt.UTC().String(),
			SourceChainId: tx.GetFromChainId().String(),
			SourceAddress: tx.From,
		}
		dtos = append(dtos, dto)
	}

	return dtos
}

// GetEnygmaBatchTransactions retrieves paginated transactions in an enygma batch
func (s *transactionService) GetEnygmaBatchTransactions(
	ctx context.Context,
	batchId string,
	page, limit int,
) (*types.Paginated[dto.BatchTransactionDto], error) {
	if batchId == "" {
		return nil, NewValidationError("batchId", "batchId cannot be empty")
	}

	transactions, total, err := s.txRepo.FindByEnygmaBatchIdPaginated(ctx, batchId, page, limit)
	if err != nil {
		return nil, NewInternalError("FindByEnygmaBatchIdPaginated", err)
	}

	if total == 0 {
		return nil, NewNotFoundError("enygma batch", batchId)
	}

	return &types.Paginated[dto.BatchTransactionDto]{
		Data:  s.mapper.ToBatchTransactionDtos(transactions),
		Total: total,
		Limit: limit,
		Page:  page,
	}, nil
}

// GetRegularBatchTransactions retrieves paginated transactions in a regular batch
func (s *transactionService) GetRegularBatchTransactions(
	ctx context.Context,
	batchId string,
	page, limit int,
) (*types.Paginated[dto.BatchTransactionDto], error) {
	if batchId == "" {
		return nil, NewValidationError("batchId", "batchId cannot be empty")
	}

	transactions, total, err := s.txRepo.FindByBatchIdPaginated(ctx, batchId, page, limit)
	if err != nil {
		return nil, NewInternalError("FindByBatchIdPaginated", err)
	}

	if total == 0 {
		return nil, NewNotFoundError("batch", batchId)
	}

	return &types.Paginated[dto.BatchTransactionDto]{
		Data:  s.mapper.ToBatchTransactionDtos(transactions),
		Total: total,
		Limit: limit,
		Page:  page,
	}, nil
}

// GetTransactionsBySharedId retrieves transactions by shared_id (zkdvp swaps)
func (s *transactionService) GetTransactionsBySharedId(
	ctx context.Context,
	sharedId string,
) ([]dto.DvpSwapTransactionDto, error) {
	if sharedId == "" {
		return nil, NewValidationError("sharedId", "sharedId cannot be empty")
	}

	transactions, err := s.txRepo.FindBySharedId(ctx, sharedId)
	if err != nil {
		return nil, NewInternalError("FindBySharedId", err)
	}

	if len(transactions) == 0 {
		return nil, NewNotFoundError("transactions with shared_id", sharedId)
	}

	return s.mapper.ToDvpSwapTransactionDtos(transactions), nil
}

// validateFilters validates and sanitizes query parameters
func (s *transactionService) validateFilters(filters *dto.MergedTransactionsFilters) error {
	// Normalize pagination
	if filters.Limit < 1 || filters.Limit > 1000 {
		filters.Limit = 10
	}
	if filters.Page < 1 {
		filters.Page = 1
	}

	// Validate chain IDs and addresses
	if err := s.validateChainIdsAndAddresses(filters); err != nil {
		return err
	}

	// Validate filters
	if err := s.validateFilterParameters(filters); err != nil {
		return err
	}

	// Validate timestamps
	if _, _, err := utils.ValidateTimestampRange(filters.InitiatedAfter, filters.InitiatedBefore); err != nil {
		return NewValidationError("initiatedAfter/initiatedBefore", err.Error())
	}

	return nil
}

// validateChainIdsAndAddresses validates chain IDs and addresses together
func (s *transactionService) validateChainIdsAndAddresses(filters *dto.MergedTransactionsFilters) error {
	// Validate chain IDs
	if filters.FromChainId != "" {
		if _, err := strconv.ParseInt(filters.FromChainId, 10, 64); err != nil {
			return NewValidationError("sourceChainId", "must be a numeric chain ID")
		}
	}
	if filters.ToChainId != "" {
		if _, err := strconv.ParseInt(filters.ToChainId, 10, 64); err != nil {
			return NewValidationError("destinationChainId", "must be a numeric chain ID")
		}
	}

	// Validate and normalize addresses
	if filters.From != "" {
		if !isValidHexAddress(filters.From) {
			return NewValidationError("fromAddress", "must be a valid hex address (0x + 40 hex chars)")
		}
	}
	if filters.To != "" {
		if !isValidHexAddress(filters.To) {
			return NewValidationError("toAddress", "must be a valid hex address (0x + 40 hex chars)")
		}
	}

	return nil
}

// validateFilterParameters validates resourceId and messageType
func (s *transactionService) validateFilterParameters(filters *dto.MergedTransactionsFilters) error {
	// Validate resourceId
	if filters.ResourceId != "" {
		buff := make([]byte, hex.DecodedLen(len(filters.ResourceId)))
		if _, err := hex.Decode(buff, []byte(filters.ResourceId)); err != nil {
			return NewValidationError("resourceId", "must be a valid hex string")
		}
	}

	// Validate messageType
	if filters.MessageType != "" {
		if _, ok := types.StringToAssetType[filters.MessageType]; !ok {
			return NewValidationError(
				"messageType",
				"must be one of: erc20, erc721, erc1155, enygma, zkdvp_erc721, zkdvp_erc1155, custom",
			)
		}
	}

	return nil
}

// ethAddressHexLength is the expected length of an Ethereum address in hex (20 bytes = 40 hex chars)
const ethAddressHexLength = 40

// isValidHexAddress checks if a string is a valid Ethereum hex address
func isValidHexAddress(addr string) bool {
	if !strings.HasPrefix(addr, "0x") && !strings.HasPrefix(addr, "0X") {
		return false
	}

	// Remove 0x prefix
	hexPart := addr[2:]

	// Must be exactly 40 characters (20 bytes)
	if len(hexPart) != ethAddressHexLength {
		return false
	}

	// Must be valid hex
	_, err := hex.DecodeString(hexPart)
	return err == nil
}

// determineIdAndType returns both the identifier type and the actual ID value
// that should be used by the frontend to query transaction details based on the aggregation type.
func determineIdAndType(tx domain.Transaction) (idType string, id string) {
	switch tx.AggregationType {
	case domain.AggregationTypeTransaction:
		if tx.MessageId != "" {
			return "message_id", tx.MessageId
		}
		return "transaction_id", tx.ID.String()
	case domain.AggregationTypeDvpSwap:
		return "shared_id", tx.SharedId
	case domain.AggregationTypeEnygmaBatch:
		if tx.MessageId != "" {
			return "batch_id", tx.AggregationKey
		}
		return "transaction_id", tx.ID.String()
	case domain.AggregationTypeRegularBatch:
		return "batch_id", tx.AggregationKey
	default:
		return "transaction_id", tx.ID.String()
	}
}
