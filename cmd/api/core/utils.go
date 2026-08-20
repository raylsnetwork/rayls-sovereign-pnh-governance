package core

import (
	"math/big"
	"strconv"
	"time"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

// bigIntToString converts *big.Int to string, returns empty string if nil
func bigIntToString(n *big.Int) string {
	if n == nil {
		return ""
	}
	return n.String()
}

// formatTimestampToUnix converts time.Time to Unix timestamp string
func formatTimestampToUnix(t time.Time) string {
	if t.IsZero() {
		return "0"
	}
	return strconv.FormatInt(t.Unix(), 10)
}

// mapTeleportStatus maps database teleport status to string representation according to used protocol
func (m *TransactionMapper) mapTeleportStatus(teleportStatus *uint8, protocol types.ProtocolType) string {
	if teleportStatus == nil {
		return ""
	}

	var statusMap map[uint8]string
	switch protocol {
	case types.DvpSwap:
		statusMap = types.DvpSwapStateToString
	case types.Enygma:
		statusMap = types.EnygmaTransferStatusToString
	default:
		statusMap = types.AtomicTeleportStatusToString
	}

	if statusStr, exists := statusMap[*teleportStatus]; exists {
		return statusStr
	}

	return "unknown"
}
