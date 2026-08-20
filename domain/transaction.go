package domain

import (
	"math/big"
	"time"

	"github.com/shopspring/decimal"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
)

// Aggregation type constants for transaction grouping
const (
	AggregationTypeTransaction  = "transaction"
	AggregationTypeRegularBatch = "regular_batch"
	AggregationTypeEnygmaBatch  = "enygma_batch"
	AggregationTypeDvpSwap      = "dvp_swap"
)

type Transaction struct {
	Model
	ResourceId           string             `json:"resourceId"                gorm:"index;not null"`
	Token                Token              `json:"token"                      gorm:"foreignKey:ResourceId;references:ResourceId"`
	MessageId            string             `json:"messageId"                 gorm:"index"`
	FromChainId          BigInt             `json:"fromChainId"`
	ToChainId            BigInt             `json:"toChainId"`
	Amount               decimal.Decimal    `json:"amount"`
	From                 string             `json:"from"`
	To                   string             `json:"to"`
	TxType               types.TxType       `json:"txType"`
	MsgType              uint8              `json:"msgType"`
	TeleportStatus       *uint8             `json:"teleportStatus"            gorm:"default:null"`
	Protocol             types.ProtocolType `json:"protocol"                   gorm:"default:null"`
	IsFlagged            bool               `json:"isFlagged"`
	HubTxHash            string             `json:"hubTxHash"`
	HubTimestamp         time.Time          `json:"hubTimestamp"`
	BlockNumber          decimal.Decimal    `json:"blockNumber"`
	LogIndex             uint64             `json:"logIndex"`
	TxHashSource         string             `json:"txHashSource"`
	TxHashDestination    string             `json:"txHashDestination"`
	SourceTimestamp      time.Time          `json:"sourceTimestamp"`
	DestinationTimestamp time.Time          `json:"destinationTimestamp"`
	SharedId             string             `json:"sharedId"`
	BatchId              string             `json:"batchId"`
	ErcId                BigInt             `json:"ercId"`
	Payload              string             `json:"payload"`

	// Aggregation fields for scalable pagination
	AggregationType string `json:"aggregationType" gorm:"default:transaction"` // transaction, regular_batch, enygma_batch
	AggregationKey  string `json:"aggregationKey"`                             // message_id, batch_id, or enygma batch_id

	// Reverse relationships
	EnygmaTransaction     *EnygmaTransaction     `gorm:"foreignKey:TransactionId;references:ID" json:"enygmaTransaction,omitempty"`
	RevertDataTransaction *RevertDataTransaction `gorm:"foreignKey:TransactionId;references:ID" json:"revertDataTransaction,omitempty"`
}

// GetFromChainId returns the source chain ID as *big.Int (unwraps from BigInt)
func (t Transaction) GetFromChainId() *big.Int {
	return t.FromChainId.Unwrap()
}

// GetToChainId returns the destination chain ID as *big.Int (unwraps from BigInt)
func (t Transaction) GetToChainId() *big.Int {
	return t.ToChainId.Unwrap()
}

// GetErcId returns the ERC token ID as *big.Int (unwraps from BigInt)
func (t Transaction) GetErcId() *big.Int {
	return t.ErcId.Unwrap()
}

// GetMsgTypeString returns the string representation of MsgType
func (t Transaction) GetMsgTypeString() string {
	return types.AssetTypeToString[uint8(t.MsgType)]
}

// GetProtocolString returns the string representation of Protocol
func (t Transaction) GetProtocolString() string {
	return types.ProtocolTypeToString[t.Protocol]
}
