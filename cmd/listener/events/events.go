package events

// Contract names - only contracts that are actively processed
const (
	ContractTokenCore          = "TokenCore"
	ContractTokenFreezeManager = "TokenFreezeManager"
	ContractEnygmaTokenManager = "EnygmaTokenManager"
	ContractParticipantCore    = "ParticipantCore"
	ContractAuditManager       = "AuditManager"
	ContractTeleport           = "Teleport"
	ContractEnygmaTeleport     = "EnygmaTeleport"
	ContractProofs             = "Proofs"
	ContractDvpTeleport        = "DvpTeleport"
)

// Event names and signatures - only for active contracts
const (
	// TokenCore events
	Erc20TokenRegistered         = "Erc20TokenRegistered"
	Erc20TokenRegisteredSig      = "Erc20TokenRegistered(bytes32,uint256,uint256,string,uint256)"
	Erc721TokenRegistered        = "Erc721TokenRegistered"
	Erc721TokenRegisteredSig     = "Erc721TokenRegistered(bytes32,uint256,uint256,string,uint256[])"
	Erc1155TokenRegistered       = "Erc1155TokenRegistered"
	Erc1155TokenRegisteredSig    = "Erc1155TokenRegistered(bytes32,uint256,uint256,string,(uint256,uint256)[])"
	DvpErc721TokenRegistered     = "DvpErc721TokenRegistered"
	DvpErc721TokenRegisteredSig  = "DvpErc721TokenRegistered(bytes32,uint256,uint256,string,uint256[])"
	DvpErc1155TokenRegistered    = "DvpErc1155TokenRegistered"
	DvpErc1155TokenRegisteredSig = "DvpErc1155TokenRegistered(bytes32,uint256,uint256,string,(uint256,uint256)[])"
	TokenStatusUpdated           = "TokenStatusUpdated"
	TokenStatusUpdatedSig        = "TokenStatusUpdated(uint256,string,uint8)"
	TokenBalanceUpdated          = "TokenBalanceUpdated"
	TokenBalanceUpdatedSig       = "TokenBalanceUpdated(bytes32,uint256,uint8,(uint256,uint256))"

	// TokenFreezeManager events
	TokenFreezeStatusChanged    = "TokenFreezeStatusChanged"
	TokenFreezeStatusChangedSig = "TokenFreezeStatusChanged(bytes32,uint256[],uint8)"

	// EnygmaTokenManager events
	EnygmaTokenRegistered    = "EnygmaTokenRegistered"
	EnygmaTokenRegisteredSig = "EnygmaTokenRegistered(bytes32,uint256,uint256,string,uint256)"

	// Teleport events
	AtomicMessageAdditionalDataBatch    = "AtomicMessageAdditionalDataBatch"
	AtomicMessageAdditionalDataBatchSig = "AtomicMessageAdditionalDataBatch(string[],string)"
	AtomicMessageStatusChangedBatch     = "AtomicMessageStatusChangedBatch"
	AtomicMessageStatusChangedBatchSig  = "AtomicMessageStatusChangedBatch(string[],uint8)"
	EncryptedDataBatchStored            = "EncryptedDataBatchStored"
	EncryptedDataBatchStoredSig         = "EncryptedDataBatchStored(string,bytes,uint256)"

	// ParticipantCore events
	ParticipantRegistered    = "ParticipantRegistered"
	ParticipantRegisteredSig = "ParticipantRegistered((uint256,uint8,uint8,string,string,uint256,uint256,bool))"
	ParticipantUpdated       = "ParticipantUpdated"
	ParticipantUpdatedSig    = "ParticipantUpdated((uint256,uint8,uint8,string,string,uint256,uint256,bool))"

	// AuditManager events
	NewAuditOrChainInfo    = "NewAuditOrChainInfo"
	NewAuditOrChainInfoSig = "NewAuditOrChainInfo()"

	// Proofs events
	HeaderProofSubmitted    = "HeaderProofSubmitted"
	HeaderProofSubmittedSig = "HeaderProofSubmitted(uint256,uint256,bytes32)"

	// EnygmaTeleport
	EnygmaTransfer             = "EnygmaTransfer"
	EnygmaTransferSig          = "EnygmaTransfer(bytes32,bytes,uint256,uint256,uint256[],uint256[],uint256)"
	EnygmaTransferCompleted    = "EnygmaTransferCompleted"
	EnygmaTransferCompletedSig = "EnygmaTransferCompleted(bytes)"
	EnygmaSupplyUpdated        = "EnygmaSupplyUpdated"
	EnygmaSupplyUpdatedSig     = "EnygmaSupplyUpdated(bytes32,uint256,(uint256,uint8),uint256)"
	EnygmaDvpBalanceUpdated    = "EnygmaDvpBalanceUpdated"
	EnygmaDvpBalanceUpdatedSig = "EnygmaDvpBalanceUpdated(bytes)"

	// DvpTeleport events
	ERCDvpBalanceUpdated    = "ERCDvpBalanceUpdated"
	ERCDvpBalanceUpdatedSig = "ERCDvpBalanceUpdated(bytes)"
	SwapInitiated           = "SwapInitiated"
	SwapInitiatedSig        = "SwapInitiated(bytes32,bytes,bytes,uint256,uint256)"
	SwapCompleted           = "SwapCompleted"
	SwapCompletedSig        = "SwapCompleted(bytes32,bytes)"
	SwapCancelled           = "SwapCancelled"
	SwapCancelledSig        = "SwapCancelled(bytes32)"
	SwapTimedOut            = "SwapTimedOut"
	SwapTimedOutSig         = "SwapTimedOut(bytes32)"
)
