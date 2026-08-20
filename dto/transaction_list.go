package dto

type TransactionListDto struct {
	Type          string `json:"type"`
	Id            string `json:"id"`
	IdType        string `json:"idType"`
	CreatedAt     string `json:"createdAt"`
	SourceChainId string `json:"sourceChainId"`
	SourceAddress string `json:"sourceAddress"`
}
