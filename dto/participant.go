package dto

// ParticipantListFilters represents the query parameters for filtering participants
type ParticipantListFilters struct {
	Name          string `form:"name,omitempty"`
	ChainId       *uint  `form:"chainId,omitempty"`
	Status        string `form:"status,omitempty"        enums:"new,active,inactive,frozen"`
	Role          string `form:"role,omitempty"          enums:"participant,issuer,auditor"`
	CreatedAfter  string `form:"createdAfter,omitempty"`
	CreatedBefore string `form:"createdBefore,omitempty"`
}
