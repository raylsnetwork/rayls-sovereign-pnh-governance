package domain

import (
	"github.com/shopspring/decimal"
)

type Balance struct {
	Model
	ResourceId string          `json:"resourceId"`
	ChainId    string          `json:"chainId"`
	Amount     decimal.Decimal `json:"amount"`
	ErcId      BigInt          `json:"ercId"`
}
