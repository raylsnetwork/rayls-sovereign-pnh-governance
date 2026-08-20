package types

import (
	"math/big"
	"time"
)

// AssetType represents different types of blockchain assets (ERC20, ERC721, etc.)
type AssetType uint8

const (
	AssetTypeCustom AssetType = iota
	AssetTypeERC20
	AssetTypeERC721
	AssetTypeERC1155
	AssetTypeEnygma
	AssetTypeDvpERC721
	AssetTypeDvpERC1155
)

// The on-chain SharedObjects.ErcStandard enum appends test-only variants immediately after
// DvpERC1155, byte-identical on the wire to their base standard. FACTORY-mode registers example
// bytecode under these ordinals, so the gov-api receives e.g. 7 (ERC20TEST) for what is
// functionally an ERC20. These ordinals are not in AssetTypeToString, so without normalization a
// FACTORY-mode example token yields an empty ercStandard string and breaks ercId normalization /
// balance indexing.
const (
	AssetTypeERC20Test      AssetType = iota + AssetTypeDvpERC1155 + 1 // 7
	AssetTypeERC721Test                                                // 8
	AssetTypeERC1155Test                                               // 9
	AssetTypeEnygmaTest                                                // 10
	AssetTypeDvpERC721Test                                             // 11
	AssetTypeDvpERC1155Test                                            // 12
)

// testToProductionAssetType maps each test-only ordinal to the production token it folds into.
// Anything not in this map is left unchanged.
var testToProductionAssetType = map[uint8]uint8{
	uint8(AssetTypeERC20Test):      uint8(AssetTypeERC20),
	uint8(AssetTypeERC721Test):     uint8(AssetTypeERC721),
	uint8(AssetTypeERC1155Test):    uint8(AssetTypeERC1155),
	uint8(AssetTypeEnygmaTest):     uint8(AssetTypeEnygma),
	uint8(AssetTypeDvpERC721Test):  uint8(AssetTypeDvpERC721),
	uint8(AssetTypeDvpERC1155Test): uint8(AssetTypeDvpERC1155),
}

// NormalizeAssetType folds a test-only ErcStandard ordinal (ERC20TEST..DvpERC1155TEST) back to its
// production base. Production, custom, and unknown ordinals pass through unchanged.
func NormalizeAssetType(ercStandard uint8) uint8 {
	if production, isTest := testToProductionAssetType[ercStandard]; isTest {
		return production
	}
	return ercStandard
}

var AssetTypeToString = map[uint8]string{
	uint8(AssetTypeCustom):     "custom",
	uint8(AssetTypeERC20):      "erc20",
	uint8(AssetTypeERC721):     "erc721",
	uint8(AssetTypeERC1155):    "erc1155",
	uint8(AssetTypeEnygma):     "enygma",
	uint8(AssetTypeDvpERC721):  "dvp_erc721",
	uint8(AssetTypeDvpERC1155): "dvp_erc1155",
}

var StringToAssetType = map[string]uint8{
	"custom":      uint8(AssetTypeCustom),
	"erc20":       uint8(AssetTypeERC20),
	"erc721":      uint8(AssetTypeERC721),
	"erc1155":     uint8(AssetTypeERC1155),
	"enygma":      uint8(AssetTypeEnygma),
	"dvp_erc721":  uint8(AssetTypeDvpERC721),
	"dvp_erc1155": uint8(AssetTypeDvpERC1155),
}

// BridgeTransactionType represents the type of bridge transaction (Transfer or Proof)
type BridgeTransactionType int

// UpdateTransaction represents a transaction that updates token balances
type UpdateTransaction struct {
	ResourceId   [32]byte  `json:"resource_id"`
	Amount       *big.Int  `json:"amount"`
	UpdateType   TxType    `json:"update_type"`
	TxHash       string    `json:"tx_hash"`
	BlockNumber  string    `json:"block_number"`
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	ErcId        *big.Int  `json:"erc_id"`
	MsgType      uint8     `json:"msg_type"`
	BatchId      string    `json:"batch_id"`
	LogIndex     uint64    `json:"log_index"`
	HubTimestamp time.Time `json:"hub_timestamp"`
}

// TxType represents the type of transaction operation (Burn, Mint, CrossChain)
type TxType uint8

const (
	Burn TxType = iota
	Mint
	CrossChain
)

// ProtocolType defines the type of protocol being perfomed in an operation.
type ProtocolType uint8

const (
	Custom ProtocolType = iota
	Vanilla
	Atomic
	Enygma
	DvpDeposit
	DvpWithdraw
	DvpSwap
)

var ProtocolTypeToString = map[ProtocolType]string{
	Custom:      "CUSTOM",
	Vanilla:     "VANILLA",
	Atomic:      "ATOMIC",
	Enygma:      "ENYGMA",
	DvpDeposit:  "DVP_DEPOSIT",
	DvpWithdraw: "DVP_WITHDRAW",
	DvpSwap:     "DVP_SWAP",
}
