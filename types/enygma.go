package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// EnygmaTxType represents Enygma transaction types (Mint or Burn)
type EnygmaTxType uint8

const (
	EnygmaTxTypeMint EnygmaTxType = 1
	EnygmaTxTypeBurn EnygmaTxType = 2
)

// EnygmaTransferStatus represents the lifecycle state of an Enygma cross-chain transfer
type EnygmaTransferStatus uint8

const (
	EnygmaTransferPending  EnygmaTransferStatus = iota // 0 - transfer initiated, awaiting destination confirmation
	EnygmaTransferExecuted                             // 1 - transfer confirmed at destination
)

// EnygmaTransferStatusToString maps EnygmaTransferStatus values to their string representations
var EnygmaTransferStatusToString = map[uint8]string{
	uint8(EnygmaTransferPending):  "pending",
	uint8(EnygmaTransferExecuted): "executed",
}

// EnygmaTransferCompleted represents a completed Enygma transfer
type EnygmaTransferCompleted struct {
	MessageId       string
	TransactionHash string
	ChainId         *big.Int
}

// EnygmaTransferBatch represents a batch of Enygma transfers
type EnygmaTransferBatch struct {
	ResourceId     string
	HubBlockNumber *big.Int
	FromChainID    *big.Int
	ToChainID      *big.Int
	ToRValueToAdd  *big.Int
	Transactions   []*EnygmaTransferBatchTx
	BatchId        string
}

// EnygmaTransferBatchTx represents a single transaction in an Enygma transfer batch
type EnygmaTransferBatchTx struct {
	MessageId   string
	ReferenceId [32]byte
	FromAddress common.Address
	ToAmount    *big.Int
	ToAddress   common.Address
	// ToCallables []EnygmaPNEvents.SharedObjectsEnygmaCrossTransferCallable
}
