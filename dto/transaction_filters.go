package dto

type MergedTransactionsFilters struct {
	FromChainId     string `form:"sourceChainId"`
	ToChainId       string `form:"destinationChainId"`
	From            string `form:"fromAddress"`
	To              string `form:"toAddress"`
	ResourceId      string `form:"resourceId"`
	MessageId       string `form:"messageId"`
	MessageType     string `form:"messageType" enums:"custom,erc20,erc721,erc1155,enygma,dvp_erc721,dvp_erc1155"`
	InitiatedAfter  string `form:"initiatedAfter"`
	InitiatedBefore string `form:"initiatedBefore"`
	Limit           int    `form:"limit"`
	Page            int    `form:"page"`
}
