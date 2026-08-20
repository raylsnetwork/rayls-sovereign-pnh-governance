package domain

import (
	"time"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

// HeaderFlagEvent represents a header flagging event in the database
type HeaderFlagEvent struct {
	ID          int       `gorm:"primaryKey"     json:"id"`
	ChainID     BigInt    `gorm:"type:numeric"   json:"chainId"`
	BlockNumber BigInt    `gorm:"type:numeric"   json:"blockNumber"`
	Reason      uint8     `                      json:"reason"`
	Initiator   uint8     `                      json:"initiator"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

// GetReasonString returns the string representation of Reason
func (h HeaderFlagEvent) GetReasonString() string {
	return types.HeaderFlagReasonToString[h.Reason]
}

// GetInitiatorString returns the string representation of Initiator
func (h HeaderFlagEvent) GetInitiatorString() string {
	return types.HeaderFlagInitiatorToString[h.Initiator]
}
