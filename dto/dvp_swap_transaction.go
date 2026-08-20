package dto

import (
	"github.com/shopspring/decimal"
)

// DvpSwapTransactionDto is an optimized DTO for dvp swap transaction endpoints
type DvpSwapTransactionDto struct {
	TransactionId              string          `json:"transactionId"`              // Transaction ID
	CreatedAt                  string          `json:"createdAt"`                  // DB record creation timestamp
	SourceChainId              string          `json:"sourceChainId"`              // Source blockchain ID
	SourceAddress              string          `json:"sourceAddress"`              // Sender address
	DestinationChainId         string          `json:"destinationChainId"`         // Destination blockchain ID
	DestinationAddress         string          `json:"destinationAddress"`         // Receiver address
	Amount                     decimal.Decimal `json:"amount"`                     // Transfer amount
	TokenName                  string          `json:"tokenName"`                  // Token name
	TokenSymbol                string          `json:"tokenSymbol"`                // Token symbol
	SourceTransactionHash      string          `json:"sourceTransactionHash"`      // Tx hash on source chain
	DestinationTransactionHash string          `json:"destinationTransactionHash"` // Tx hash on destination chain
	HubTimestamp               string          `json:"hubTimestamp,omitempty"`     // Hub timestamp (Unix seconds)
	ResourceId                 string          `json:"resourceId"`                 // Token resource ID
	MessageType                string          `json:"messageType"`                // Asset type string
	Protocol                   string          `json:"protocol"`                   // Protocol string
	Status                     string          `json:"status"`                     // DvP swap status
	ErcId                      string          `json:"ercId,omitempty"`            // Token ID for ERC721/ERC1155
	Payload                    string          `json:"payload,omitempty"`          // Encoded calldata
	HubHash                    string          `json:"hubHash,omitempty"`          // Private Network Hub transaction hash
}
