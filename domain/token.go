package domain

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

type Token struct {
	Model
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	ResourceId  string `json:"resourceId"`
	MetadataUrl string `json:"metadataUrl"`
	ErcStandard uint8  `json:"ercStandard"`
	Decimals    uint8  `json:"decimals"`
	IssuerId    string `json:"issuerId"`
	Status      uint8  `json:"status"`
}

// CirculatingSupplyEntry represents a single balance entry in circulatingSupply
type CirculatingSupplyEntry struct {
	ParticipantId string          `json:"participantId"`
	TokenId       string          `json:"tokenId,omitempty"` // Only present for non-fungible tokens (ERC721, ERC1155)
	Balance       decimal.Decimal `json:"balance"`
}

type TokenWithBalancesAndFreezeState struct {
	Token
	TotalSupply       decimal.Decimal          `json:"totalSupply"`
	CirculatingSupply []CirculatingSupplyEntry `json:"circulatingSupply"`
	FrozenChainIds    []string                 `json:"frozenChainIds"`
}

// circulatingSupplyJSON is a custom scanner for PostgreSQL json_agg results
type CirculatingSupplyJSON []CirculatingSupplyEntry

func (c *CirculatingSupplyJSON) Scan(value any) error {
	if value == nil {
		*c = []CirculatingSupplyEntry{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("circulatingSupplyJSON: expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, c)
}

// FrozenChainIdsJSON is a custom scanner for PostgreSQL json_agg results containing chain IDs
type FrozenChainIdsJSON []string

func (f *FrozenChainIdsJSON) Scan(value any) error {
	if value == nil {
		*f = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("frozenChainIdsJSON: expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, f)
}
