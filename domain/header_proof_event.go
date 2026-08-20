package domain

import (
	"time"
)

// HeaderProofEvent represents a blockchain header proof submitted to the Private Network Hub
type HeaderProofEvent struct {
	ID          int       `gorm:"primaryKey"                json:"id"`
	ChainID     BigInt    `gorm:"type:numeric;not null"     json:"chainId"`
	BlockNumber BigInt    `gorm:"type:numeric;not null"     json:"blockNumber"`
	BlockHash   string    `gorm:"type:varchar(66);not null" json:"blockHash"`
	CreatedAt   time.Time `gorm:"not null;default:now()"    json:"createdAt"`
}
