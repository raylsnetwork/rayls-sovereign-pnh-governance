package domain

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

type EnygmaTransaction struct {
	TransactionId uuid.UUID   `json:"transactionId"`
	Transaction   Transaction `json:"transaction"       gorm:"foreignKey:TransactionId;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	ToRValueToAdd BigInt      `json:"toRValueToAdd"`
	ReferenceId   string      `json:"referenceId"`
	UpdatedAt     time.Time   `json:"updatedAt"`
}

// GetToRValueToAdd returns the ToRValueToAdd as *big.Int (unwraps from BigInt)
func (t EnygmaTransaction) GetToRValueToAdd() *big.Int {
	return t.ToRValueToAdd.Unwrap()
}
