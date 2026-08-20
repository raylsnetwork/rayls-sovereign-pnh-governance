package dto

import "time"

// TokenRegistryStatusDto represents the response from the blockchain TokenRegistry contract.
// It omits the database-specific ID field since this data comes directly from the smart contract.
type TokenRegistryStatusDto struct {
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `json:"name"`
	Symbol      string    `json:"symbol"`
	ResourceId  string    `json:"resourceId"`
	MetadataUrl string    `json:"metadataUrl"`
	ErcStandard uint8     `json:"ercStandard"`
	Decimals    uint8     `json:"decimals"`
	IssuerId    string    `json:"issuerId"`
	Status      uint8     `json:"status"`
}
