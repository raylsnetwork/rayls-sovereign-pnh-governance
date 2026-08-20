package domain

import "github.com/google/uuid"

type FlaggedTransaction struct {
	Model
	TransactionId uuid.UUID `json:"transactionId"`
}
