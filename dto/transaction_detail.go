package dto

import (
	"github.com/shopspring/decimal"
)

// TransactionDetailDto represents a single transaction
type TransactionDetailDto struct {
	Id                         string          `json:"id"`
	MessageId                  string          `json:"messageId"`
	SourceTransactionHash      string          `json:"sourceTransactionHash,omitempty"`
	DestinationTransactionHash string          `json:"destinationTransactionHash,omitempty"`
	Status                     string          `json:"status,omitempty"`
	SourceChainId              string          `json:"sourceChainId,omitempty"`
	DestinationChainId         string          `json:"destinationChainId,omitempty"`
	SourceTimestamp            string          `json:"sourceTimestamp,omitempty"`
	DestinationTimestamp       string          `json:"destinationTimestamp,omitempty"`
	SourceAddress              string          `json:"sourceAddress"`
	DestinationAddress         string          `json:"destinationAddress"`
	Amount                     decimal.Decimal `json:"amount"`
	ResourceId                 string          `json:"resourceId"`
	MessageType                string          `json:"messageType"`
	TokenName                  string          `json:"tokenName,omitempty"`
	TokenSymbol                string          `json:"tokenSymbol,omitempty"`
	TokenDecimals              uint8           `json:"tokenDecimals"`
	Protocol                   string          `json:"protocol"`
	ErcId                      string          `json:"ercId,omitempty"`
	TokenMetadataImageUrl      string          `json:"tokenMetadataImageUrl,omitempty"`
	TokenMetadataName          string          `json:"tokenMetadataName,omitempty"`
	TokenMetadataDescription   string          `json:"tokenMetadataDescription,omitempty"`
	Payload                    string          `json:"payload,omitempty"`
	// Enygma-specific fields (omitted for non-Enygma)
	HubHash        string `json:"hubHash,omitempty"`
	HubTimestamp   string `json:"hubTimestamp,omitempty"`
	HubBlockNumber string `json:"hubBlockNumber,omitempty"`
	// Revert data (omitted for non-reverted transactions)
	RevertDataTransaction *RevertDataTransactionDto `json:"revertDataTransaction,omitempty"`
	// Metadata
	Type      string `json:"type,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// RevertDataTransactionDto represents revert-specific data for a reverted transaction
type RevertDataTransactionDto struct {
	TxHashDestinationRevert       string `json:"txHashDestinationRevert,omitempty"`
	TxHashDestinationRevertStatus string `json:"txHashDestinationRevertStatus,omitempty"`
	TxHashSourceRevert            string `json:"txHashSourceRevert,omitempty"`
	TxHashSourceRevertStatus      string `json:"txHashSourceRevertStatus,omitempty"`
}
