package domain

import (
	"time"

	"github.com/google/uuid"
)

// TokenFreezeState represents the current freeze state for a token-participant pair
type TokenFreezeState struct {
	ResourceId string    `json:"resourceId" gorm:"column:resource_id;primaryKey;not null"`
	ChainId    string    `json:"chainId"    gorm:"column:chain_id;primaryKey;not null"`
	CreatedAt  time.Time `json:"createdAt"  gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `json:"updatedAt"  gorm:"column:updated_at;not null"`
	IsFrozen   bool      `json:"isFrozen"   gorm:"column:is_frozen;not null;default:false"`
}

// TokenFreezeAudit represents a historical record of freeze/unfreeze operations
type TokenFreezeAudit struct {
	ID              uuid.UUID `json:"id"              gorm:"column:id;type:uuid;primaryKey"`
	ResourceId      string    `json:"resourceId"      gorm:"column:resource_id;not null"`
	ChainId         string    `json:"chainId"         gorm:"column:chain_id;not null"`
	Action          uint8     `json:"action"          gorm:"column:action;not null"`
	BlockNumber     BigInt    `json:"blockNumber"     gorm:"column:block_number;type:numeric;not null"`
	TransactionHash string    `json:"transactionHash" gorm:"column:transaction_hash;not null"`
	CreatedAt       time.Time `json:"createdAt"       gorm:"column:created_at;not null;default:now()"`
}
