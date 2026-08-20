package domain

import (
	"github.com/google/uuid"
)

type RevertDataTransaction struct {
	Transaction                   Transaction `gorm:"foreignKey:TransactionId;references:ID" json:"transaction"`
	TransactionId                 uuid.UUID   `                                              json:"transactionId"`
	TxHashDestinationRevert       string      `                                              json:"txHashDestinationRevert"`
	TxHashDestinationRevertStatus uint8       `                                              json:"txHashDestinationRevertStatus"`
	TxHashSourceRevert            string      `                                              json:"txHashSourceRevert"`
	TxHashSourceRevertStatus      uint8       `                                              json:"txHashSourceRevertStatus"`
}
