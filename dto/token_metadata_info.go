package dto

import (
	"encoding/json"
)

type TokenMetadataInfoDto struct {
	ImageUrl    string
	Description string
	Name        string
}

type TokenMetadataStandartDto struct {
	Title      string `json:"title"`
	Type       string `json:"type"`
	Properties struct {
		Name struct {
			Type        string `json:"type"`
			Description string `json:"description"`
		} `json:"name"`
		Description struct {
			Type        string `json:"type"`
			Description string `json:"description"`
		} `json:"description"`
		Image struct {
			Type        string `json:"type"`
			Description string `json:"description"`
		} `json:"image"`
	} `json:"properties"`
}

type TokenMetadataStandart2Dto struct {
	Image      string `json:"image"`
	Name       string `json:"name"`
	Attributes []struct {
		TraitType string          `json:"trait_type"`
		Value     json.RawMessage `json:"value"`
	} `json:"attributes"`
}

type TokenListFilters struct {
	Name          string `form:"name,omitempty"`
	Symbol        string `form:"symbol,omitempty"`
	IssuerId      string `form:"issuerId,omitempty"`
	Status        string `form:"status,omitempty"        enums:"new,active,inactive"`
	ErcStandard   string `form:"ercStandard,omitempty"   enums:"custom,erc20,erc721,erc1155,enygma,dvp_erc721,dvp_erc1155"`
	Decimals      *uint8 `form:"decimals,omitempty"`
	CreatedAfter  string `form:"createdAfter,omitempty"`
	CreatedBefore string `form:"createdBefore,omitempty"`
	Limit         int    `form:"limit"`
	Page          int    `form:"page"`
}
