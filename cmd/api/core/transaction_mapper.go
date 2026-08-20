package core

import (
	"context"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

// TransactionMapper handles the domain model to DTO transformation
// Includes external metadata fetching for NFTs (will be removed in a later version)
type TransactionMapper struct {
	metadataService TokenMetadataService
}

// NewTransactionMapper creates a new mapper
func NewTransactionMapper(metadataService TokenMetadataService) *TransactionMapper {
	return &TransactionMapper{
		metadataService: metadataService,
	}
}

// ToBatchTransactionDtos converts domain transactions to BatchTransactionDto
func (m *TransactionMapper) ToBatchTransactionDtos(transactions []domain.Transaction) []dto.BatchTransactionDto {
	dtos := make([]dto.BatchTransactionDto, len(transactions))

	for i, tx := range transactions {
		dtos[i] = dto.BatchTransactionDto{
			MessageId:          tx.MessageId,
			CreatedAt:          tx.CreatedAt.UTC().String(),
			SourceChainId:      bigIntToString(tx.GetFromChainId()),
			SourceAddress:      tx.From,
			DestinationChainId: bigIntToString(tx.GetToChainId()),
			DestinationAddress: tx.To,
			Amount:             tx.Amount,
			TokenSymbol:        tx.Token.Symbol,
		}
	}

	return dtos
}

// ToDvpSwapTransactionDtos converts domain transactions to DvpSwapTransactionDto for DvP Swap transactions
func (m *TransactionMapper) ToDvpSwapTransactionDtos(transactions []domain.Transaction) []dto.DvpSwapTransactionDto {
	dtos := make([]dto.DvpSwapTransactionDto, len(transactions))

	for i, tx := range transactions {
		d := dto.DvpSwapTransactionDto{
			TransactionId:              tx.ID.String(),
			CreatedAt:                  tx.CreatedAt.UTC().String(),
			SourceChainId:              bigIntToString(tx.GetFromChainId()),
			SourceAddress:              tx.From,
			DestinationChainId:         bigIntToString(tx.GetToChainId()),
			DestinationAddress:         tx.To,
			Amount:                     tx.Amount,
			TokenName:                  tx.Token.Name,
			TokenSymbol:                tx.Token.Symbol,
			SourceTransactionHash:      tx.TxHashSource,
			DestinationTransactionHash: tx.TxHashDestination,
			ResourceId:                 tx.ResourceId,
			MessageType:                tx.GetMsgTypeString(),
			Protocol:                   tx.GetProtocolString(),
			Status:                     m.mapTeleportStatus(tx.TeleportStatus, tx.Protocol),
			Payload:                    tx.Payload,
			HubHash:                    tx.HubTxHash,
		}

		if !tx.HubTimestamp.IsZero() {
			d.HubTimestamp = formatTimestampToUnix(tx.HubTimestamp)
		}

		ercId := tx.GetErcId()
		if ercId != nil && ercId.String() != "" {
			d.ErcId = ercId.String()
		}

		dtos[i] = d
	}

	return dtos
}

// ToTransactionDetailDto maps a single domain.Transaction to a flat TransactionDetailDto
// Handles Enygma vs non-Enygma logic and fetches NFT metadata for ERC721/ERC1155
func (m *TransactionMapper) ToTransactionDetailDto(
	ctx context.Context,
	tx domain.Transaction,
	txType string,
) (*dto.TransactionDetailDto, error) {
	isEnygma := tx.MsgType == uint8(types.AssetTypeEnygma)

	result := &dto.TransactionDetailDto{
		Id:                         tx.ID.String(),
		MessageId:                  tx.MessageId,
		SourceTransactionHash:      tx.TxHashSource,
		DestinationTransactionHash: tx.TxHashDestination,
		SourceChainId:              bigIntToString(tx.GetFromChainId()),
		DestinationChainId:         bigIntToString(tx.GetToChainId()),
		SourceAddress:              tx.From,
		DestinationAddress:         tx.To,
		Amount:                     tx.Amount,
		ResourceId:                 tx.ResourceId,
		MessageType:                tx.GetMsgTypeString(),
		TokenName:                  tx.Token.Name,
		TokenSymbol:                tx.Token.Symbol,
		TokenDecimals:              tx.Token.Decimals,
		Protocol:                   tx.GetProtocolString(),
		Payload:                    tx.Payload,
		Type:                       txType,
	}

	// Handle Enygma-specific fields
	if isEnygma {
		result.HubHash = tx.HubTxHash
		result.HubBlockNumber = tx.BlockNumber.String()
		if !tx.HubTimestamp.IsZero() {
			result.HubTimestamp = formatTimestampToUnix(tx.HubTimestamp)
		}
	}

	// Omit source/destination timestamps from the response of Enygma transactions
	if !isEnygma {
		if !tx.SourceTimestamp.IsZero() {
			result.SourceTimestamp = formatTimestampToUnix(tx.SourceTimestamp)
		}
		if !tx.DestinationTimestamp.IsZero() {
			result.DestinationTimestamp = formatTimestampToUnix(tx.DestinationTimestamp)
		}
	}
	result.CreatedAt = tx.CreatedAt.UTC().String()

	// Normalize teleport status for all protocols
	result.Status = m.mapTeleportStatus(tx.TeleportStatus, tx.Protocol)

	// Handle ErcId and fetch ERC721/ERC1155 metadata if applicable
	ercId := tx.GetErcId()
	isNft := tx.MsgType == uint8(types.AssetTypeERC721) || tx.MsgType == uint8(types.AssetTypeERC1155) ||
		tx.MsgType == uint8(types.AssetTypeDvpERC721) || tx.MsgType == uint8(types.AssetTypeDvpERC1155)

	if isNft && ercId != nil && ercId.String() != "" {
		result.ErcId = ercId.String()
		tokenBaseUrl := tx.Token.MetadataUrl
		if tokenBaseUrl != "" {
			tokenInfo, err := m.metadataService.GetMetadata(ctx, tokenBaseUrl, ercId.String())
			if err == nil {
				result.TokenMetadataImageUrl = tokenInfo.ImageUrl
				result.TokenMetadataName = tokenInfo.Name
				result.TokenMetadataDescription = tokenInfo.Description
			}
			// If metadata fetch fails, fields remain empty (omitempty will exclude them)
		}
	}

	// Map revert data if present
	if tx.RevertDataTransaction != nil {
		result.RevertDataTransaction = &dto.RevertDataTransactionDto{
			TxHashDestinationRevert:       tx.RevertDataTransaction.TxHashDestinationRevert,
			TxHashDestinationRevertStatus: types.AtomicTeleportStatusToString[tx.RevertDataTransaction.TxHashDestinationRevertStatus],
			TxHashSourceRevert:            tx.RevertDataTransaction.TxHashSourceRevert,
			TxHashSourceRevertStatus:      types.AtomicTeleportStatusToString[tx.RevertDataTransaction.TxHashSourceRevertStatus],
		}
	}

	return result, nil
}
