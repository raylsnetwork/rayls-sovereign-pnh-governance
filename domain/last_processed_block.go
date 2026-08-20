package domain

import "time"

type LastProcessedBlock struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Number    BigInt `json:"block_number"`
}

func (LastProcessedBlock) TableName() string {
	return "last_processed_block"
}
