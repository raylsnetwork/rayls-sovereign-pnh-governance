package types

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// DvpSwapState represents the state of a DvP swap
type DvpSwapState uint8

const (
	DvpSwapStateCompleted DvpSwapState = iota
	DvpSwapStateFailed
	DvpSwapStateCancelled
	DvpSwapStateExpired
	// Governance-specific state for pending/initiated swaps
	DvpSwapStatePending
)

var DvpSwapStateToString = map[uint8]string{
	uint8(DvpSwapStateCompleted): "completed",
	uint8(DvpSwapStateFailed):    "failed",
	uint8(DvpSwapStateCancelled): "cancelled",
	uint8(DvpSwapStateExpired):   "expired",
	uint8(DvpSwapStatePending):   "pending",
}

// DvpBalanceUpdated represents a decrypted DVP balance update event
type DvpBalanceUpdated struct {
	ErcId      *big.Int
	TokenType  uint8
	ResourceId string
	SharedId   string

	From string
	To   string

	SourceChainId     *big.Int
	SourceTxHash      string
	SourceTxTimestamp time.Time

	DestinationTxHash      string
	DestinationTxTimestamp time.Time
	DestinationChainId     *big.Int

	Amount     *big.Int
	UpdateType TxType
}

// DvpTokenType represents token types used in DVP protocol
type DvpTokenType uint8

const (
	DvpCustom DvpTokenType = iota
	DvpERC20
	DvpERC721
	DvpERC1155
	DvpEnygma
)

// DvpTokenTypeToAssetType converts a DvpTokenType to AssetType
func DvpTokenTypeToAssetType(dvpType DvpTokenType) uint8 {
	switch dvpType {
	case DvpCustom:
		return uint8(AssetTypeCustom)
	case DvpERC20:
		return uint8(AssetTypeERC20)
	case DvpERC721:
		return uint8(AssetTypeDvpERC721)
	case DvpERC1155:
		return uint8(AssetTypeDvpERC1155)
	case DvpEnygma:
		return uint8(AssetTypeEnygma)
	default:
		return uint8(AssetTypeCustom)
	}
}

// DvpSwapMessage is the plaintext payload carried (encrypted) by the
// SwapInitiated and SwapCompleted events emitted by the DvpTeleport contract.
// It mirrors the struct produced by the relayer that encrypts the payload.
type DvpSwapMessage struct {
	SharedId      string
	To            string
	ChainId       *big.Int
	PNTxHash      string
	PNTxTimestamp time.Time

	TokenInAmount      *big.Int
	TokenInAddress     string
	TokenInResourceID  string
	TokenInType        DvpTokenType
	TokenInID          string
	TokenOutAmount     *big.Int
	TokenOutAddress    string
	TokenOutResourceID string
	TokenOutType       DvpTokenType
	TokenOutID         string

	InitiatorSelfSalt *big.Int
}

// DvpSwapSalts bundles the two per-swap salts recovered on SwapInitiated.
//
//   - InitiatorSelfSalt: from the DvpSwapMessage payload; decrypts SwapCompleted.
//   - InitiatorCtxtSalt: reduced from the ML-KEM secret in the SwapInitiated
//     Ctxt; decrypts the SwapInitiated payload and any follow-up keyed off it.
type DvpSwapSalts struct {
	InitiatorSelfSalt []byte
	InitiatorCtxtSalt []byte
}

// ResourceIDPNCommunicator represents the Resource ID used by the Privacy Node Communicator.
// This constant is defined in src/rayls-protocol-sdk/Constants.sol and is used by
// DVP middleware arbitrary messages.
var ResourceIDPNCommunicator = common.HexToHash(
	"0x0000000000000000000000000000000000000000000000000000000000000005",
)
