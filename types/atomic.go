package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// RaylzMessage represents a message in the Rayls bridge protocol
type RaylzMessage struct {
	MessageMetadata RaylzMessageMetadata
	Payload         []byte
}

// RaylzMessageMetadata contains metadata for a Rayls bridge message
type RaylzMessageMetadata struct {
	Valid            bool
	Nonce            *big.Int
	ResourceId       [32]byte
	LockData         []byte
	TransferMetadata BridgedTransferMetadata
	IgnoresNonce     bool
}

// BridgedTransferMetadata contains transfer details for bridged assets
type BridgedTransferMetadata struct {
	Id        *big.Int
	Amount    *big.Int
	From      string
	To        string
	AssetType AssetType
}

// DispatchedMessageToPrivateHub represents a message dispatched from a Privacy Node to the Private Network Hub
type DispatchedMessageToPrivateHub struct {
	MessageId             [32]byte       `json:"message_id"`
	From                  common.Address `json:"from"`
	ToChainId             *big.Int       `json:"to_chain_id"`
	To                    common.Address `json:"to"`
	Data                  RaylzMessage   `json:"data"`
	FromChainId           *big.Int       `json:"from_chain_id"`
	SharedId              string         `json:"shared_id"`
	TxHashSource          common.Hash    `json:"tx_hash_source"`
	TxHashSourceTimestamp uint64         `json:"tx_hash_source_timestamp"`
	TxHashSourceStatus    int8           `json:"tx_hash_source_status"`
	// only used when vanilla; in case of Atomic - the ones in additionalData are used
	TxHashDestination          common.Hash `json:"tx_hash_destination"`
	TxHashDestinationTimestamp uint64      `json:"tx_hash_destination_timestamp"`
	TxHashDestinationStatus    int8        `json:"tx_hash_destination_status"`
	// proofs
	Proofs                     []byte                `json:"proofs"`
	TxTrieProof                common.Hash           `json:"tx_trie_proof"`
	BlockHash                  common.Hash           `json:"block_hash"`
	TxLocation                 int                   `json:"tx_location"`
	TransactionType            BridgeTransactionType `json:"transaction_type"`
	LogIdx                     uint                  `json:"log_idx"`
	IsAtomic                   bool                  `json:"is_atomic"`
	TxSentToDestinationSuccess bool                  `json:"tx_sent_to_destination"`
	BatchId                    string                `json:"batch_id"`
	BatchHubTxHash             common.Hash           `json:"batch_hub_tx_hash"`
}

// AtomicTeleportAdditionalData contains additional data for atomic teleport transactions
type AtomicTeleportAdditionalData struct {
	TxHashDestinationRevert       common.Hash `json:"tx_hash_destination_revert,omitempty"`
	TxHashDestinationRevertStatus int8        `json:"tx_hash_destination_revert_status,omitempty"`
	TxHashSourceRevert            common.Hash `json:"tx_hash_source_revert,omitempty"`
	TxHashSourceRevertStatus      int8        `json:"tx_hash_source_revert_status,omitempty"`
	TxHashDestination             common.Hash `json:"tx_hash_destination,omitempty"`
	TxHashDestinationTimestamp    uint64      `json:"tx_hash_destination_timestamp,omitempty"`
	TxHashDestinationStatus       int8        `json:"tx_hash_destination_status,omitempty"`
	RevertReason                  string      `json:"revert_reason,omitempty"`
	SharedId                      string      `json:"shared_id,omitempty"`
	BatchHubTxHash                common.Hash `json:"batch_hub_tx_hash,omitempty"`
}

// AtomicTeleportStatus represents the state of an Atomic transaction
type AtomicTeleportStatus uint8

const (
	AtomicTeleportPending AtomicTeleportStatus = iota
	AtomicTeleportExecuted
	AtomicTeleportRejected
	AtomicTeleportReverted
	AtomicTeleportCreditMemo
)

var AtomicTeleportStatusToString = map[uint8]string{
	uint8(AtomicTeleportPending):    "pending",
	uint8(AtomicTeleportExecuted):   "executed",
	uint8(AtomicTeleportRejected):   "rejected",
	uint8(AtomicTeleportReverted):   "reverted",
	uint8(AtomicTeleportCreditMemo): "credit_memo", // Status used only inside Listener, to credit back amounts on reverted transactions.
}

const (
	Transfer BridgeTransactionType = 1
	Proof    BridgeTransactionType = 2
)
