package dto

import "github.com/shopspring/decimal"

// BatchTransactionDto is an optimized DTO for batch transaction endpoints
type BatchTransactionDto struct {
	MessageId          string          `json:"messageId"`          // Transaction message ID
	CreatedAt          string          `json:"createdAt"`          // DB record creation timestamp
	SourceChainId      string          `json:"sourceChainId"`      // Source blockchain ID
	SourceAddress      string          `json:"sourceAddress"`      // Sender address
	DestinationChainId string          `json:"destinationChainId"` // Destination blockchain ID
	DestinationAddress string          `json:"destinationAddress"` // Receiver address
	Amount             decimal.Decimal `json:"amount"`             // Transfer amount
	TokenSymbol        string          `json:"tokenSymbol"`        // Token symbol
}
