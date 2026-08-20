// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ParticipantStorageV1

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// ParticipantStructsAuditInfoData is an auto generated low-level Go binding around an user-defined struct.
type ParticipantStructsAuditInfoData struct {
	ChainId                      *big.Int
	RaylsViewPublicKey           string
	EncryptedRaylsViewPrivateKey []byte
	Mac                          []byte
	BlockNumber                  *big.Int
}

// ParticipantStructsKeyAgreementData is an auto generated low-level Go binding around an user-defined struct.
type ParticipantStructsKeyAgreementData struct {
	ChainId     *big.Int
	Ciphertext  []byte
	Digest      []byte
	BlockNumber *big.Int
}

// ParticipantStructsParticipant is an auto generated low-level Go binding around an user-defined struct.
type ParticipantStructsParticipant struct {
	ChainId            *big.Int
	Role               uint8
	Status             uint8
	OwnerId            string
	Name               string
	CreatedAt          *big.Int
	UpdatedAt          *big.Int
	AllowedToBroadcast bool
}

// ParticipantStructsParticipantData is an auto generated low-level Go binding around an user-defined struct.
type ParticipantStructsParticipantData struct {
	ChainId            *big.Int
	Role               uint8
	OwnerId            string
	Name               string
	AllowedToBroadcast bool
}

// ParticipantStructsPrivacyNodeSpendDataSafeReturn is an auto generated low-level Go binding around an user-defined struct.
type ParticipantStructsPrivacyNodeSpendDataSafeReturn struct {
	PaymentSpendPublicKey *big.Int
	PnAddresses           []common.Address
	ChainId               *big.Int
}

// ParticipantStructsPrivacyNodeViewData is an auto generated low-level Go binding around an user-defined struct.
type ParticipantStructsPrivacyNodeViewData struct {
	ChainId            *big.Int
	RaylsViewPublicKey string
	BlockNumber        *big.Int
}

// ParticipantStorageV1MetaData contains all meta data concerning the ParticipantStorageV1 contract.
var ParticipantStorageV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addParticipant\",\"inputs\":[{\"name\":\"_participant\",\"type\":\"tuple\",\"internalType\":\"structParticipantStructs.ParticipantData\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Role\"},{\"name\":\"ownerId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"allowedToBroadcast\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addParticipants\",\"inputs\":[{\"name\":\"_participants\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.ParticipantData[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Role\"},{\"name\":\"ownerId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"allowedToBroadcast\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"broadcastCurrentParticipants\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkEnygmaAccountAllowed\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkEnygmaIssuerAccountAllowed\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configureModules\",\"inputs\":[{\"name\":\"_participantCore\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_auditManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_enygmaManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"contractVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"endpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRaylsEndpoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllParticipants\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.Participant[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Role\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Status\"},{\"name\":\"ownerId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allowedToBroadcast\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllParticipantsChainIds\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllPaymentSpendPublicKeys\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.PrivacyNodeSpendDataSafeReturn[]\",\"components\":[{\"name\":\"paymentSpendPublicKey\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pnAddresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllPrivacyNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAuditInfo\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"data\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.AuditInfoData[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raylsViewPublicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encryptedRaylsViewPrivateKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"mac\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAuditManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIAuditManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getChainViewData\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"data\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.PrivacyNodeViewData[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raylsViewPublicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEnygmaAllParticipantsChainIds\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEnygmaManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEnygmaManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getKeyAgreements\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.KeyAgreementData[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"digest\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParticipant\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structParticipantStructs.Participant\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Role\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Status\"},{\"name\":\"ownerId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allowedToBroadcast\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParticipantCore\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIParticipantCore\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParticipantDataBatch\",\"inputs\":[],\"outputs\":[{\"name\":\"pnViewData\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.PrivacyNodeViewData[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raylsViewPublicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"auditInfo\",\"type\":\"tuple[]\",\"internalType\":\"structParticipantStructs.AuditInfoData[]\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raylsViewPublicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encryptedRaylsViewPrivateKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"mac\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pnChainIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"auditChainIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPaymentSpendPublicKeyByChainId\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initiateKeyAgreement\",\"inputs\":[{\"name\":\"initiatorChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"responderChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"digest\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeParticipant\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setAuditInfo\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raylsViewPublicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encryptedRaylsViewPrivateKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"mac\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAuditManager\",\"inputs\":[{\"name\":\"_auditManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainViewData\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raylsViewPublicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEnygmaManager\",\"inputs\":[{\"name\":\"_enygmaManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEnygmaPnEventsAddress\",\"inputs\":[{\"name\":\"_pnEnygmaEvents\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setParticipantCore\",\"inputs\":[{\"name\":\"_participantCore\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPaymentSpendPublicKey\",\"inputs\":[{\"name\":\"_chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_paymentSpendPublicKey\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_pnAddresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateBroadcastMessagesPermission\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allowed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateRole\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Role\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateStatus\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumParticipantStructs.Status\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"validateMessageParticipants\",\"inputs\":[{\"name\":\"originChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validateParticipantStatus\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyParticipant\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ModulesConfigured\",\"inputs\":[{\"name\":\"participantCore\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"auditManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"enygmaManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParticipantStorageV1__UnauthorizedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAppV1__HubNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"hubStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsAppV1__PrivacyNodeFrozen\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAppV1__PrivacyNodeNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsAppV1__PublicChainNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"publicChainStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsAppV1__ResourceNotApproved\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAppV1__TokenRegistryNotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAppV1__UnauthorizedTokenRegistry\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "ParticipantStorageV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b50608051613d4d61003e60003960008181611e6e01528181611e970152611fce0152613d4d6000f3fe60806040526004361061023d5760003560e01c80635f997c5b1161012f578063ad3cb1cc116100b1578063ad3cb1cc14610695578063b3816fdb146106d3578063bf7e214f146106f3578063c4d66de814610708578063c6885dc714610728578063c68bab0f1461074a578063c9557f561461076a578063d5c3614f1461078a578063d6b81d5e146107aa578063df5cc335146107c8578063e2446536146107f5578063ede591fd1461081557600080fd5b80635f997c5b146105485780636279aa9f1461055e578063628de2651461057e57806365ffa17e1461059e578063683f7f27146105be57806381a40de6146105de5780638c8e0c2a146105fe578063987345cb1461061e578063a01afbfb1461063c578063a0842e331461065c578063a0a8e4601461068157600080fd5b806337e0d45d116101c357806337e0d45d146103da5780633a1b3d31146103ef5780634017734d1461040f578063485cc9551461042f5780634b9bd78c1461044f5780634cc168c91461047f5780634ec53b2e146104945780634f1ef286146104b457806352bf1c8d146104c757806352d1902d146104e7578063594c5c761461050a5780635e280f111461052857600080fd5b8063035c87131461024257806306ec249d14610264578063079c2a49146102795780630d30d092146102995780630fd44407146102cf57806311f50c85146102fc578063195ec9ee146103295780631b9db2ef1461034b5780632d3358751461037857806331b1125f1461039a57806333ded4a8146103ba575b600080fd5b34801561024e57600080fd5b5061026261025d3660046123c1565b610835565b005b34801561027057600080fd5b5061026261089c565b34801561028557600080fd5b506102626102943660046123de565b610947565b3480156102a557600080fd5b506102b96102b4366004612460565b6109f2565b6040516102c691906124c9565b60405180910390f35b3480156102db57600080fd5b506102ef6102ea366004612460565b610a94565b6040516102c691906125e0565b34801561030857600080fd5b5061031c610317366004612460565b610b30565b6040516102c691906125f3565b34801561033557600080fd5b5061033e610b9e565b6040516102c691906126ce565b34801561035757600080fd5b5061036b610366366004612460565b610c49565b6040516102c69190612732565b34801561038457600080fd5b5061038d610d2d565b6040516102c69190612781565b3480156103a657600080fd5b506102626103b5366004612794565b610dc7565b3480156103c657600080fd5b506102626103d53660046127ed565b610ebd565b3480156103e657600080fd5b5061038d610f66565b3480156103fb57600080fd5b5061026261040a36600461282a565b610fe4565b34801561041b57600080fd5b5061026261042a366004612897565b611054565b34801561043b57600080fd5b5061026261044a3660046128e9565b6110c8565b34801561045b57600080fd5b5061046f61046a3660046123c1565b6111d8565b60405190151581526020016102c6565b34801561048b57600080fd5b5061038d611273565b3480156104a057600080fd5b506102626104af366004612917565b6112f1565b6102626104c2366004612af4565b6113a6565b3480156104d357600080fd5b506102626104e23660046123c1565b6113c5565b3480156104f357600080fd5b506104fc611432565b6040519081526020016102c6565b34801561051657600080fd5b506002546001600160a01b031661031c565b34801561053457600080fd5b5060005461031c906001600160a01b031681565b34801561055457600080fd5b506104fc60015481565b34801561056a57600080fd5b50610262610579366004612b43565b611450565b34801561058a57600080fd5b50610262610599366004612cb0565b611502565b3480156105aa57600080fd5b506102626105b9366004612d60565b611570565b3480156105ca57600080fd5b506102626105d9366004612460565b6115e0565b3480156105ea57600080fd5b506104fc6105f9366004612460565b61164f565b34801561060a57600080fd5b5061046f610619366004612d85565b6116e5565b34801561062a57600080fd5b506003546001600160a01b031661031c565b34801561064857600080fd5b50610262610657366004612460565b61178e565b34801561066857600080fd5b506106716117d9565b6040516102c69493929190612e53565b34801561068d57600080fd5b5060016104fc565b3480156106a157600080fd5b506106c6604051806040016040528060058152602001640352e302e360dc1b81525081565b6040516102c69190612eab565b3480156106df57600080fd5b5061046f6106ee366004612460565b611893565b3480156106ff57600080fd5b5061031c6118ef565b34801561071457600080fd5b506102626107233660046123c1565b611908565b34801561073457600080fd5b5061073d611932565b6040516102c69190612ebe565b34801561075657600080fd5b506102626107653660046123c1565b6119cc565b34801561077657600080fd5b50610262610785366004612460565b611a2a565b34801561079657600080fd5b506102626107a5366004612f74565b611aaa565b3480156107b657600080fd5b506004546001600160a01b031661031c565b3480156107d457600080fd5b506107e86107e3366004612460565b611b31565b6040516102c69190612f96565b34801561080157600080fd5b50610262610810366004612fa9565b611bcd565b34801561082157600080fd5b506102626108303660046123c1565b611c3b565b61084b336000356001600160e01b031916611c99565b6001600160a01b03811661087a5760405162461bcd60e51b815260040161087190612fe5565b60405180910390fd5b600280546001600160a01b0319166001600160a01b0392909216919091179055565b6108b2336000356001600160e01b031916611c99565b6002546001600160a01b03166108da5760405162461bcd60e51b815260040161087190613042565b60006108e4611ddb565b60025460405163e91638eb60e01b8152600481018390529192506001600160a01b03169063e91638eb906024015b600060405180830381600087803b15801561092c57600080fd5b505af1158015610940573d6000803e3d6000fd5b5050505050565b61095d336000356001600160e01b031916611c99565b6004546001600160a01b03166109855760405162461bcd60e51b815260040161087190613096565b6004805460405163079c2a4960e01b81526001600160a01b039091169163079c2a49916109ba918891889188918891016130e8565b600060405180830381600087803b1580156109d457600080fd5b505af11580156109e8573d6000803e3d6000fd5b5050505050505050565b6003546060906001600160a01b0316610a1d5760405162461bcd60e51b81526004016108719061313a565b600354604051630698684960e11b8152600481018490526001600160a01b0390911690630d30d09290602401600060405180830381865afa158015610a66573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610a8e91908101906131d0565b92915050565b6003546060906001600160a01b0316610abf5760405162461bcd60e51b81526004016108719061313a565b600354604051630fd4440760e01b8152600481018490526001600160a01b0390911690630fd4440790602401600060405180830381865afa158015610b08573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610a8e91908101906133bc565b600080546040516311f50c8560e01b8152600481018490526001600160a01b03909116906311f50c8590602401602060405180830381865afa158015610b7a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a8e91906133f0565b6002546060906001600160a01b0316610bc95760405162461bcd60e51b815260040161087190613042565b600260009054906101000a90046001600160a01b03166001600160a01b031663195ec9ee6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610c1c573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610c4491908101906134f3565b905090565b610c9460408051610100810190915260008082526020820190815260200160008152602001606081526020016060815260200160008152602001600081526020016000151581525090565b6002546001600160a01b0316610cbc5760405162461bcd60e51b815260040161087190613042565b600254604051631b9db2ef60e01b8152600481018490526001600160a01b0390911690631b9db2ef90602401600060405180830381865afa158015610d05573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610a8e9190810190613596565b6004546060906001600160a01b0316610d585760405162461bcd60e51b815260040161087190613096565b6004805460408051632d33587560e01b815290516001600160a01b0390921692632d3358759282820192600092908290030181865afa158015610d9f573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610c449190810190613629565b610ddd336000356001600160e01b031916611c99565b6001600160a01b038316610e035760405162461bcd60e51b815260040161087190612fe5565b6001600160a01b038216610e295760405162461bcd60e51b81526004016108719061365d565b6001600160a01b038116610e4f5760405162461bcd60e51b8152600401610871906136b6565b600280546001600160a01b03199081166001600160a01b03868116918217909355600380548316868516908117909155600480549093169385169384179092556040517ff6c27e15b7e995e3edd056f3e9b7b01098dfe3f91cccf2af78ff33215fc1829d90600090a4505050565b610ed3336000356001600160e01b031916611c99565b6002546001600160a01b0316610efb5760405162461bcd60e51b815260040161087190613042565b60025460405163067bda9560e31b81526004810184905282151560248201526001600160a01b03909116906333ded4a8906044015b600060405180830381600087803b158015610f4a57600080fd5b505af1158015610f5e573d6000803e3d6000fd5b505050505050565b6002546060906001600160a01b0316610f915760405162461bcd60e51b815260040161087190613042565b600260009054906101000a90046001600160a01b03166001600160a01b03166337e0d45d6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610d9f573d6000803e3d6000fd5b610ffa336000356001600160e01b031916611c99565b6002546001600160a01b03166110225760405162461bcd60e51b815260040161087190613042565b600254604051633a1b3d3160e01b81526001600160a01b0390911690633a1b3d3190610f309085908590600401613710565b61106a336000356001600160e01b031916611c99565b6003546001600160a01b03166110925760405162461bcd60e51b81526004016108719061313a565b600354604051634017734d60e01b81526001600160a01b0390911690634017734d906109ba90879087908790879060040161374d565b60006110d2611def565b805490915060ff600160401b82041615906001600160401b03166000811580156110f95750825b90506000826001600160401b031660011480156111155750303b155b905081158015611123575080155b156111415760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561116b57845460ff60401b1916600160401b1785555b611173611e18565b61117c87611908565b6001805561118986611e22565b83156111cf57845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b6004546000906001600160a01b03166112035760405162461bcd60e51b815260040161087190613096565b600480546040516312e6f5e360e21b81526001600160a01b0390911691634b9bd78c91611232918691016125f3565b602060405180830381865afa15801561124f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a8e9190613778565b6003546060906001600160a01b031661129e5760405162461bcd60e51b81526004016108719061313a565b600360009054906101000a90046001600160a01b03166001600160a01b0316634cc168c96040518163ffffffff1660e01b8152600401600060405180830381865afa158015610d9f573d6000803e3d6000fd5b611307336000356001600160e01b031916611c99565b6003546001600160a01b031661132f5760405162461bcd60e51b81526004016108719061313a565b6003546040516327629d9760e11b81526001600160a01b0390911690634ec53b2e9061136b908a908a908a908a908a908a908a90600401613795565b600060405180830381600087803b15801561138557600080fd5b505af1158015611399573d6000803e3d6000fd5b5050505050505050505050565b6113ae611e63565b6113b782611ef1565b6113c18282611f0a565b5050565b6113db336000356001600160e01b031916611c99565b6004546001600160a01b03166114035760405162461bcd60e51b815260040161087190613096565b600480546040516352bf1c8d60e01b81526001600160a01b03909116916352bf1c8d91610912918591016125f3565b600061143c611fc3565b50600080516020613cf88339815191525b90565b611466336000356001600160e01b031916611c99565b6003546001600160a01b031661148e5760405162461bcd60e51b81526004016108719061313a565b600354604051636279aa9f60e01b81526001600160a01b0390911690636279aa9f906114c8908990899089908990899089906004016137dd565b600060405180830381600087803b1580156114e257600080fd5b505af11580156114f6573d6000803e3d6000fd5b50505050505050505050565b611518336000356001600160e01b031916611c99565b6002546001600160a01b03166115405760405162461bcd60e51b815260040161087190613042565b60025460405163628de26560e01b81526001600160a01b039091169063628de26590610912908490600401613895565b611586336000356001600160e01b031916611c99565b6002546001600160a01b03166115ae5760405162461bcd60e51b815260040161087190613042565b6002546040516332ffd0bf60e11b81526001600160a01b03909116906365ffa17e90610f3090859085906004016138ec565b6115f6336000356001600160e01b031916611c99565b6002546001600160a01b031661161e5760405162461bcd60e51b815260040161087190613042565b60025460405163683f7f2760e01b8152600481018390526001600160a01b039091169063683f7f2790602401610912565b6004546000906001600160a01b031661167a5760405162461bcd60e51b815260040161087190613096565b600480546040516340d206f360e11b81529182018490526001600160a01b0316906381a40de690602401602060405180830381865afa1580156116c1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a8e9190613900565b6004546000906001600160a01b03166117105760405162461bcd60e51b815260040161087190613096565b60048054604051634647061560e11b81526001600160a01b038681169382019390935260248101859052911690638c8e0c2a90604401602060405180830381865afa158015611763573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117879190613778565b9392505050565b600061179861200c565b9050336001600160a01b038216146117d357604051620d23e560e01b81523360048201526001600160a01b0382166024820152604401610871565b50600155565b6003546060908190819081906001600160a01b031661180a5760405162461bcd60e51b81526004016108719061313a565b600360009054906101000a90046001600160a01b03166001600160a01b031663a0842e336040518163ffffffff1660e01b8152600401600060405180830381865afa15801561185d573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526118859190810190613a36565b935093509350935090919293565b6002546000906001600160a01b03166118be5760405162461bcd60e51b815260040161087190613042565b60025460405163b3816fdb60e01b8152600481018490526001600160a01b039091169063b3816fdb90602401611232565b60006118f96120a3565b546001600160a01b0316919050565b611910612105565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b6004546060906001600160a01b031661195d5760405162461bcd60e51b815260040161087190613096565b600480546040805163c6885dc760e01b815290516001600160a01b039092169263c6885dc79282820192600092908290030181865afa1580156119a4573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610c449190810190613ae2565b6119e2336000356001600160e01b031916611c99565b6001600160a01b038116611a085760405162461bcd60e51b81526004016108719061365d565b600380546001600160a01b0319166001600160a01b0392909216919091179055565b6002546001600160a01b0316611a525760405162461bcd60e51b815260040161087190613042565b6002546040516364aabfab60e11b8152600481018390526001600160a01b039091169063c9557f569060240160006040518083038186803b158015611a9657600080fd5b505afa158015610940573d6000803e3d6000fd5b6002546001600160a01b0316611ad25760405162461bcd60e51b815260040161087190613042565b60025460405163d5c3614f60e01b815260048101849052602481018390526001600160a01b039091169063d5c3614f9060440160006040518083038186803b158015611b1d57600080fd5b505afa158015610f5e573d6000803e3d6000fd5b6003546060906001600160a01b0316611b5c5760405162461bcd60e51b81526004016108719061313a565b60035460405163df5cc33560e01b8152600481018490526001600160a01b039091169063df5cc33590602401600060405180830381865afa158015611ba5573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610a8e9190810190613c28565b611be3336000356001600160e01b031916611c99565b6002546001600160a01b0316611c0b5760405162461bcd60e51b815260040161087190613042565b600254604051637122329b60e11b81526001600160a01b039091169063e244653690610912908490600401613c5c565b611c51336000356001600160e01b031916611c99565b6001600160a01b038116611c775760405162461bcd60e51b8152600401610871906136b6565b600480546001600160a01b0319166001600160a01b0392909216919091179055565b6000611ca36120a3565b80549091506001600160a01b031680611cd2576000604051638944034760e01b815260040161087191906125f3565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015611d36573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d5a9190613c6f565b925092509250826111cf578015611d845760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff821615611dc05760405163a426878960e01b81526001600160a01b038816600482015263ffffffff83166024820152604401610871565b86604051632ecd3d0360e21b815260040161087191906125f3565b60006034361061144d575060331936013590565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610a8e565b611e20612105565b565b6000611e2c6120a3565b80549091506001600160a01b031615611e5a5781604051638944034760e01b815260040161087191906125f3565b6113c18261212a565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480611ed357507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316611ec76121ba565b6001600160a01b031614155b15611e205760405163703e46dd60e11b815260040160405180910390fd5b611f07336000356001600160e01b031916611c99565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015611f64575060408051601f3d908101601f19168201909252611f6191810190613900565b60015b611f835781604051634c9c8ce360e01b815260040161087191906125f3565b600080516020613cf88339815191528114611fb457604051632a87526960e21b815260048101829052602401610871565b611fbe83836121d0565b505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614611e205760405163703e46dd60e11b815260040160405180910390fd5b600080546040516311f50c8560e01b8152600360048201526001600160a01b03909116906311f50c8590602401602060405180830381865afa158015612056573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061207a91906133f0565b90506001600160a01b03811661144d57604051633eba255b60e01b815260040160405180910390fd5b60008060ff196120d460017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35613cba565b6040516020016120e691815260200190565b60408051601f1981840301815291905280516020909101201692915050565b61210d612226565b611e2057604051631afcd79f60e31b815260040160405180910390fd5b6001600160a01b0381166121535780604051638944034760e01b815260040161087191906125f3565b600061215d6120a3565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b6000600080516020613cf88339815191526118f9565b6121d982612240565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a280511561221e57611fbe828261229c565b6113c1612312565b6000612230611def565b54600160401b900460ff16919050565b806001600160a01b03163b60000361226d5780604051634c9c8ce360e01b815260040161087191906125f3565b600080516020613cf883398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b0316846040516122b99190613cdb565b600060405180830381855af49150503d80600081146122f4576040519150601f19603f3d011682016040523d82523d6000602084013e6122f9565b606091505b5091509150612309858383612331565b95945050505050565b3415611e205760405163b398979f60e01b815260040160405180910390fd5b6060826123465761234182612384565b611787565b815115801561235d57506001600160a01b0384163b155b1561237d5783604051639996b31560e01b815260040161087191906125f3565b5092915050565b80511561239357805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b6001600160a01b0381168114611f0757600080fd5b6000602082840312156123d357600080fd5b8135611787816123ac565b600080600080606085870312156123f457600080fd5b843593506020850135925060408501356001600160401b038082111561241957600080fd5b818701915087601f83011261242d57600080fd5b81358181111561243c57600080fd5b8860208260051b850101111561245157600080fd5b95989497505060200194505050565b60006020828403121561247257600080fd5b5035919050565b60005b8381101561249457818101518382015260200161247c565b50506000910152565b600081518084526124b5816020860160208601612479565b601f01601f19169290920160200192915050565b600060208083018184528085518083526040925060408601915060408160051b87010184880160005b8381101561255957603f19898403018552815160808151855288820151818a8701526125208287018261249d565b9150508782015185820389870152612538828261249d565b606093840151969093019590955250948701949250908601906001016124f2565b509098975050505050505050565b600082825180855260208086019550808260051b84010181860160005b848110156125d357601f1986840301895281516060815185528582015181878701526125b28287018261249d565b60409384015196909301959095525098840198925090830190600101612584565b5090979650505050505050565b6020815260006117876020830184612567565b6001600160a01b0391909116815260200190565b634e487b7160e01b600052602160045260246000fd5b6003811061262d5761262d612607565b9052565b6004811061262d5761262d612607565b600061010082518452602083015161265c602086018261261d565b50604083015161266f6040860182612631565b5060608301518160608601526126878286018261249d565b915050608083015184820360808601526126a1828261249d565b91505060a083015160a085015260c083015160c085015260e0830151151560e08501528091505092915050565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b8281101561272557603f19888603018452612713858351612641565b945092850192908501906001016126f7565b5092979650505050505050565b6020815260006117876020830184612641565b60008151808452602080850194506020840160005b838110156127765781518752958201959082019060010161275a565b509495945050505050565b6020815260006117876020830184612745565b6000806000606084860312156127a957600080fd5b83356127b4816123ac565b925060208401356127c4816123ac565b915060408401356127d4816123ac565b809150509250925092565b8015158114611f0757600080fd5b6000806040838503121561280057600080fd5b823591506020830135612812816127df565b809150509250929050565b60048110611f0757600080fd5b6000806040838503121561283d57600080fd5b8235915060208301356128128161281d565b60008083601f84011261286157600080fd5b5081356001600160401b0381111561287857600080fd5b60208301915083602082850101111561289057600080fd5b9250929050565b600080600080606085870312156128ad57600080fd5b8435935060208501356001600160401b038111156128ca57600080fd5b6128d68782880161284f565b9598909750949560400135949350505050565b600080604083850312156128fc57600080fd5b8235612907816123ac565b91506020830135612812816123ac565b600080600080600080600060a0888a03121561293257600080fd5b873596506020880135955060408801356001600160401b038082111561295757600080fd5b6129638b838c0161284f565b909750955060608a013591508082111561297c57600080fd5b506129898a828b0161284f565b989b979a50959894979596608090950135949350505050565b634e487b7160e01b600052604160045260246000fd5b60405160a081016001600160401b03811182821017156129da576129da6129a2565b60405290565b604051608081016001600160401b03811182821017156129da576129da6129a2565b604051606081016001600160401b03811182821017156129da576129da6129a2565b60405161010081016001600160401b03811182821017156129da576129da6129a2565b604051601f8201601f191681016001600160401b0381118282101715612a6f57612a6f6129a2565b604052919050565b60006001600160401b03821115612a9057612a906129a2565b50601f01601f191660200190565b600082601f830112612aaf57600080fd5b8135612ac2612abd82612a77565b612a47565b818152846020838601011115612ad757600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215612b0757600080fd5b8235612b12816123ac565b915060208301356001600160401b03811115612b2d57600080fd5b612b3985828601612a9e565b9150509250929050565b60008060008060008060a08789031215612b5c57600080fd5b8635955060208701356001600160401b0380821115612b7a57600080fd5b612b868a838b0161284f565b90975095506040890135915080821115612b9f57600080fd5b612bab8a838b01612a9e565b94506060890135915080821115612bc157600080fd5b50612bce89828a01612a9e565b925050608087013590509295509295509295565b60006001600160401b03821115612bfb57612bfb6129a2565b5060051b60200190565b60038110611f0757600080fd5b600060a08284031215612c2457600080fd5b612c2c6129b8565b9050813581526020820135612c4081612c05565b602082015260408201356001600160401b0380821115612c5f57600080fd5b612c6b85838601612a9e565b60408401526060840135915080821115612c8457600080fd5b50612c9184828501612a9e565b6060830152506080820135612ca5816127df565b608082015292915050565b60006020808385031215612cc357600080fd5b82356001600160401b0380821115612cda57600080fd5b818501915085601f830112612cee57600080fd5b8135612cfc612abd82612be2565b81815260059190911b83018401908481019088831115612d1b57600080fd5b8585015b83811015612d5357803585811115612d375760008081fd5b612d458b89838a0101612c12565b845250918601918601612d1f565b5098975050505050505050565b60008060408385031215612d7357600080fd5b82359150602083013561281281612c05565b60008060408385031215612d9857600080fd5b8235612da3816123ac565b946020939093013593505050565b600082825180855260208086019550808260051b84010181860160005b848110156125d357601f19868403018952815160a081518552858201518187870152612dfc8287018261249d565b91505060408083015186830382880152612e16838261249d565b9250505060608083015186830382880152612e31838261249d565b6080948501519790940196909652505098840198925090830190600101612dce565b608081526000612e666080830187612567565b8281036020840152612e788187612db1565b90508281036040840152612e8c8186612745565b90508281036060840152612ea08185612745565b979650505050505050565b602081526000611787602083018461249d565b600060208083018184528085518083526040925060408601915060408160051b8701018488016000805b84811015612f6557898403603f190186528251805185528881015160608a8701819052815190870181905260808701918b019085905b80821015612f475782516001600160a01b03168452928c0192918c019160019190910190612f1e565b50505090880151948801949094529487019491870191600101612ee8565b50919998505050505050505050565b60008060408385031215612f8757600080fd5b50508035926020909101359150565b6020815260006117876020830184612db1565b600060208284031215612fbb57600080fd5b81356001600160401b03811115612fd157600080fd5b612fdd84828501612c12565b949350505050565b6020808252603c908201527f5061727469636970616e7453746f7261676556313a205061727469636970616e60408201527f74436f726520616464726573732063616e6e6f74206265207a65726f00000000606082015260800190565b60208082526034908201527f5061727469636970616e7453746f7261676556313a205061727469636970616e6040820152731d10dbdc99481b5bd91d5b19481b9bdd081cd95d60621b606082015260800190565b60208082526032908201527f5061727469636970616e7453746f7261676556313a20456e79676d614d616e6160408201527119d95c881b5bd91d5b19481b9bdd081cd95d60721b606082015260800190565b84815260208082018590526060604083018190528201839052600090849060808401835b86811015612d5357833561311f816123ac565b6001600160a01b03168252928201929082019060010161310c565b60208082526031908201527f5061727469636970616e7453746f7261676556313a2041756469744d616e6167604082015270195c881b5bd91d5b19481b9bdd081cd95d607a1b606082015260800190565b600082601f83011261319c57600080fd5b81516131aa612abd82612a77565b8181528460208386010111156131bf57600080fd5b612fdd826020830160208701612479565b600060208083850312156131e357600080fd5b82516001600160401b03808211156131fa57600080fd5b818501915085601f83011261320e57600080fd5b815161321c612abd82612be2565b81815260059190911b8301840190848101908883111561323b57600080fd5b8585015b83811015612d535780518581111561325657600080fd5b86016080818c03601f1901121561326d5760008081fd5b6132756129e0565b8882015181526040808301518881111561328f5760008081fd5b61329d8e8c8387010161318b565b8b84015250606080840151898111156132b65760008081fd5b6132c48f8d8388010161318b565b92840192909252608093909301519282019290925284525091860191860161323f565b600082601f8301126132f857600080fd5b81516020613308612abd83612be2565b82815260059290921b8401810191818101908684111561332757600080fd5b8286015b848110156133b15780516001600160401b038082111561334b5760008081fd5b908801906060828b03601f19018113156133655760008081fd5b61336d612a02565b878401518152604080850151848111156133875760008081fd5b6133958e8b8389010161318b565b8a8401525091909301519083015250835291830191830161332b565b509695505050505050565b6000602082840312156133ce57600080fd5b81516001600160401b038111156133e457600080fd5b612fdd848285016132e7565b60006020828403121561340257600080fd5b8151611787816123ac565b805161341881612c05565b919050565b80516134188161281d565b8051613418816127df565b6000610100828403121561344657600080fd5b61344e612a24565b9050815181526134606020830161340d565b60208201526134716040830161341d565b604082015260608201516001600160401b038082111561349057600080fd5b61349c8583860161318b565b606084015260808401519150808211156134b557600080fd5b506134c28482850161318b565b60808301525060a082015160a082015260c082015160c08201526134e860e08301613428565b60e082015292915050565b6000602080838503121561350657600080fd5b82516001600160401b038082111561351d57600080fd5b818501915085601f83011261353157600080fd5b815161353f612abd82612be2565b81815260059190911b8301840190848101908883111561355e57600080fd5b8585015b83811015612d535780518581111561357a5760008081fd5b6135888b89838a0101613433565b845250918601918601613562565b6000602082840312156135a857600080fd5b81516001600160401b038111156135be57600080fd5b612fdd84828501613433565b600082601f8301126135db57600080fd5b815160206135eb612abd83612be2565b8083825260208201915060208460051b87010193508684111561360d57600080fd5b602086015b848110156133b15780518352918301918301613612565b60006020828403121561363b57600080fd5b81516001600160401b0381111561365157600080fd5b612fdd848285016135ca565b60208082526039908201527f5061727469636970616e7453746f7261676556313a2041756469744d616e6167604082015278657220616464726573732063616e6e6f74206265207a65726f60381b606082015260800190565b6020808252603a908201527f5061727469636970616e7453746f7261676556313a20456e79676d614d616e6160408201527967657220616464726573732063616e6e6f74206265207a65726f60301b606082015260800190565b828152604081016117876020830184612631565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b848152606060208201526000613767606083018587613724565b905082604083015295945050505050565b60006020828403121561378a57600080fd5b8151611787816127df565b87815286602082015260a0604082015260006137b560a083018789613724565b82810360608401526137c8818688613724565b91505082608083015298975050505050505050565b86815260a0602082015260006137f760a083018789613724565b8281036040840152613809818761249d565b9050828103606084015261381d818661249d565b915050826080830152979650505050505050565b8051825260006020820151613849602085018261261d565b50604082015160a0604085015261386360a085018261249d565b90506060830151848203606086015261387c828261249d565b9150506080830151151560808501528091505092915050565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b8281101561272557603f198886030184526138da858351613831565b945092850192908501906001016138be565b82815260408101611787602083018461261d565b60006020828403121561391257600080fd5b5051919050565b600082601f83011261392a57600080fd5b8151602061393a612abd83612be2565b82815260059290921b8401810191818101908684111561395957600080fd5b8286015b848110156133b15780516001600160401b038082111561397d5760008081fd5b9088019060a0828b03601f19018113156139975760008081fd5b61399f6129b8565b878401518152604080850151848111156139b95760008081fd5b6139c78e8b8389010161318b565b8a84015250606080860151858111156139e05760008081fd5b6139ee8f8c838a010161318b565b8385015250608091508186015185811115613a095760008081fd5b613a178f8c838a010161318b565b918401919091525091909301519083015250835291830191830161395d565b60008060008060808587031215613a4c57600080fd5b84516001600160401b0380821115613a6357600080fd5b613a6f888389016132e7565b95506020870151915080821115613a8557600080fd5b613a9188838901613919565b94506040870151915080821115613aa757600080fd5b613ab3888389016135ca565b93506060870151915080821115613ac957600080fd5b50613ad6878288016135ca565b91505092959194509250565b60006020808385031215613af557600080fd5b82516001600160401b0380821115613b0c57600080fd5b818501915085601f830112613b2057600080fd5b8151613b2e612abd82612be2565b81815260059190911b83018401908481019088831115613b4d57600080fd5b8585015b83811015612d5357805185811115613b6857600080fd5b86016060818c03601f19011215613b7e57600080fd5b613b86612a02565b888201518152604082015187811115613b9e57600080fd5b8201603f81018d13613baf57600080fd5b89810151613bbf612abd82612be2565b81815260059190911b8201604001908b8101908f831115613bdf57600080fd5b6040840193505b82841015613c08578351613bf9816123ac565b8252928c0192908c0190613be6565b848d01525050506060919091015160408201528352918601918601613b51565b600060208284031215613c3a57600080fd5b81516001600160401b03811115613c5057600080fd5b612fdd84828501613919565b6020815260006117876020830184613831565b600080600060608486031215613c8457600080fd5b8351613c8f816127df565b602085015190935063ffffffff81168114613ca957600080fd5b60408501519092506127d4816127df565b81810381811115610a8e57634e487b7160e01b600052601160045260246000fd5b60008251613ced818460208701612479565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca26469706673582212206f855e729c005eac18ab47c0c0f8962d8706f1f628cc48d67455f197a28daca164736f6c63430008180033",
}

// ParticipantStorageV1 is an auto generated Go binding around an Ethereum contract.
type ParticipantStorageV1 struct {
	abi abi.ABI
}

// NewParticipantStorageV1 creates a new instance of ParticipantStorageV1.
func NewParticipantStorageV1() *ParticipantStorageV1 {
	parsed, err := ParticipantStorageV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ParticipantStorageV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ParticipantStorageV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (participantStorageV1 *ParticipantStorageV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := participantStorageV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (participantStorageV1 *ParticipantStorageV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := participantStorageV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackAddParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe2446536.
//
// Solidity: function addParticipant((uint256,uint8,string,string,bool) _participant) returns()
func (participantStorageV1 *ParticipantStorageV1) PackAddParticipant(participant ParticipantStructsParticipantData) []byte {
	enc, err := participantStorageV1.abi.Pack("addParticipant", participant)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAddParticipants is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x628de265.
//
// Solidity: function addParticipants((uint256,uint8,string,string,bool)[] _participants) returns()
func (participantStorageV1 *ParticipantStorageV1) PackAddParticipants(participants []ParticipantStructsParticipantData) []byte {
	enc, err := participantStorageV1.abi.Pack("addParticipants", participants)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) PackAuthority() []byte {
	enc, err := participantStorageV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := participantStorageV1.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackBroadcastCurrentParticipants is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06ec249d.
//
// Solidity: function broadcastCurrentParticipants() returns()
func (participantStorageV1 *ParticipantStorageV1) PackBroadcastCurrentParticipants() []byte {
	enc, err := participantStorageV1.abi.Pack("broadcastCurrentParticipants")
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCheckEnygmaAccountAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4b9bd78c.
//
// Solidity: function checkEnygmaAccountAllowed(address _address) view returns(bool)
func (participantStorageV1 *ParticipantStorageV1) PackCheckEnygmaAccountAllowed(address common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("checkEnygmaAccountAllowed", address)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCheckEnygmaAccountAllowed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4b9bd78c.
//
// Solidity: function checkEnygmaAccountAllowed(address _address) view returns(bool)
func (participantStorageV1 *ParticipantStorageV1) UnpackCheckEnygmaAccountAllowed(data []byte) (bool, error) {
	out, err := participantStorageV1.abi.Unpack("checkEnygmaAccountAllowed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackCheckEnygmaIssuerAccountAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8c8e0c2a.
//
// Solidity: function checkEnygmaIssuerAccountAllowed(address _address, uint256 _chainId) view returns(bool)
func (participantStorageV1 *ParticipantStorageV1) PackCheckEnygmaIssuerAccountAllowed(address common.Address, chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("checkEnygmaIssuerAccountAllowed", address, chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCheckEnygmaIssuerAccountAllowed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8c8e0c2a.
//
// Solidity: function checkEnygmaIssuerAccountAllowed(address _address, uint256 _chainId) view returns(bool)
func (participantStorageV1 *ParticipantStorageV1) UnpackCheckEnygmaIssuerAccountAllowed(data []byte) (bool, error) {
	out, err := participantStorageV1.abi.Unpack("checkEnygmaIssuerAccountAllowed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackConfigureModules is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x31b1125f.
//
// Solidity: function configureModules(address _participantCore, address _auditManager, address _enygmaManager) returns()
func (participantStorageV1 *ParticipantStorageV1) PackConfigureModules(participantCore common.Address, auditManager common.Address, enygmaManager common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("configureModules", participantCore, auditManager, enygmaManager)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackContractVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (participantStorageV1 *ParticipantStorageV1) PackContractVersion() []byte {
	enc, err := participantStorageV1.abi.Pack("contractVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackContractVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (participantStorageV1 *ParticipantStorageV1) UnpackContractVersion(data []byte) (*big.Int, error) {
	out, err := participantStorageV1.abi.Unpack("contractVersion", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) PackEndpoint() []byte {
	enc, err := participantStorageV1.abi.Pack("endpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) UnpackEndpoint(data []byte) (common.Address, error) {
	out, err := participantStorageV1.abi.Unpack("endpoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetAddressByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (participantStorageV1 *ParticipantStorageV1) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := participantStorageV1.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := participantStorageV1.abi.Unpack("getAddressByResourceId", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetAllParticipants is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x195ec9ee.
//
// Solidity: function getAllParticipants() view returns((uint256,uint8,uint8,string,string,uint256,uint256,bool)[])
func (participantStorageV1 *ParticipantStorageV1) PackGetAllParticipants() []byte {
	enc, err := participantStorageV1.abi.Pack("getAllParticipants")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllParticipants is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x195ec9ee.
//
// Solidity: function getAllParticipants() view returns((uint256,uint8,uint8,string,string,uint256,uint256,bool)[])
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAllParticipants(data []byte) ([]ParticipantStructsParticipant, error) {
	out, err := participantStorageV1.abi.Unpack("getAllParticipants", data)
	if err != nil {
		return *new([]ParticipantStructsParticipant), err
	}
	out0 := *abi.ConvertType(out[0], new([]ParticipantStructsParticipant)).(*[]ParticipantStructsParticipant)
	return out0, err
}

// PackGetAllParticipantsChainIds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x37e0d45d.
//
// Solidity: function getAllParticipantsChainIds() view returns(uint256[])
func (participantStorageV1 *ParticipantStorageV1) PackGetAllParticipantsChainIds() []byte {
	enc, err := participantStorageV1.abi.Pack("getAllParticipantsChainIds")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllParticipantsChainIds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x37e0d45d.
//
// Solidity: function getAllParticipantsChainIds() view returns(uint256[])
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAllParticipantsChainIds(data []byte) ([]*big.Int, error) {
	out, err := participantStorageV1.abi.Unpack("getAllParticipantsChainIds", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, err
}

// PackGetAllPaymentSpendPublicKeys is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc6885dc7.
//
// Solidity: function getAllPaymentSpendPublicKeys() view returns((uint256,address[],uint256)[])
func (participantStorageV1 *ParticipantStorageV1) PackGetAllPaymentSpendPublicKeys() []byte {
	enc, err := participantStorageV1.abi.Pack("getAllPaymentSpendPublicKeys")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllPaymentSpendPublicKeys is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc6885dc7.
//
// Solidity: function getAllPaymentSpendPublicKeys() view returns((uint256,address[],uint256)[])
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAllPaymentSpendPublicKeys(data []byte) ([]ParticipantStructsPrivacyNodeSpendDataSafeReturn, error) {
	out, err := participantStorageV1.abi.Unpack("getAllPaymentSpendPublicKeys", data)
	if err != nil {
		return *new([]ParticipantStructsPrivacyNodeSpendDataSafeReturn), err
	}
	out0 := *abi.ConvertType(out[0], new([]ParticipantStructsPrivacyNodeSpendDataSafeReturn)).(*[]ParticipantStructsPrivacyNodeSpendDataSafeReturn)
	return out0, err
}

// PackGetAllPrivacyNodes is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4cc168c9.
//
// Solidity: function getAllPrivacyNodes() view returns(uint256[])
func (participantStorageV1 *ParticipantStorageV1) PackGetAllPrivacyNodes() []byte {
	enc, err := participantStorageV1.abi.Pack("getAllPrivacyNodes")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllPrivacyNodes is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4cc168c9.
//
// Solidity: function getAllPrivacyNodes() view returns(uint256[])
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAllPrivacyNodes(data []byte) ([]*big.Int, error) {
	out, err := participantStorageV1.abi.Unpack("getAllPrivacyNodes", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, err
}

// PackGetAuditInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdf5cc335.
//
// Solidity: function getAuditInfo(uint256 chainId) view returns((uint256,string,bytes,bytes,uint256)[] data)
func (participantStorageV1 *ParticipantStorageV1) PackGetAuditInfo(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("getAuditInfo", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAuditInfo is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdf5cc335.
//
// Solidity: function getAuditInfo(uint256 chainId) view returns((uint256,string,bytes,bytes,uint256)[] data)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAuditInfo(data []byte) ([]ParticipantStructsAuditInfoData, error) {
	out, err := participantStorageV1.abi.Unpack("getAuditInfo", data)
	if err != nil {
		return *new([]ParticipantStructsAuditInfoData), err
	}
	out0 := *abi.ConvertType(out[0], new([]ParticipantStructsAuditInfoData)).(*[]ParticipantStructsAuditInfoData)
	return out0, err
}

// PackGetAuditManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x987345cb.
//
// Solidity: function getAuditManager() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) PackGetAuditManager() []byte {
	enc, err := participantStorageV1.abi.Pack("getAuditManager")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAuditManager is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x987345cb.
//
// Solidity: function getAuditManager() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetAuditManager(data []byte) (common.Address, error) {
	out, err := participantStorageV1.abi.Unpack("getAuditManager", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetChainViewData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0fd44407.
//
// Solidity: function getChainViewData(uint256 chainId) view returns((uint256,string,uint256)[] data)
func (participantStorageV1 *ParticipantStorageV1) PackGetChainViewData(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("getChainViewData", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetChainViewData is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0fd44407.
//
// Solidity: function getChainViewData(uint256 chainId) view returns((uint256,string,uint256)[] data)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetChainViewData(data []byte) ([]ParticipantStructsPrivacyNodeViewData, error) {
	out, err := participantStorageV1.abi.Unpack("getChainViewData", data)
	if err != nil {
		return *new([]ParticipantStructsPrivacyNodeViewData), err
	}
	out0 := *abi.ConvertType(out[0], new([]ParticipantStructsPrivacyNodeViewData)).(*[]ParticipantStructsPrivacyNodeViewData)
	return out0, err
}

// PackGetEnygmaAllParticipantsChainIds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2d335875.
//
// Solidity: function getEnygmaAllParticipantsChainIds() view returns(uint256[])
func (participantStorageV1 *ParticipantStorageV1) PackGetEnygmaAllParticipantsChainIds() []byte {
	enc, err := participantStorageV1.abi.Pack("getEnygmaAllParticipantsChainIds")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEnygmaAllParticipantsChainIds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2d335875.
//
// Solidity: function getEnygmaAllParticipantsChainIds() view returns(uint256[])
func (participantStorageV1 *ParticipantStorageV1) UnpackGetEnygmaAllParticipantsChainIds(data []byte) ([]*big.Int, error) {
	out, err := participantStorageV1.abi.Unpack("getEnygmaAllParticipantsChainIds", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, err
}

// PackGetEnygmaManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd6b81d5e.
//
// Solidity: function getEnygmaManager() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) PackGetEnygmaManager() []byte {
	enc, err := participantStorageV1.abi.Pack("getEnygmaManager")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEnygmaManager is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd6b81d5e.
//
// Solidity: function getEnygmaManager() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetEnygmaManager(data []byte) (common.Address, error) {
	out, err := participantStorageV1.abi.Unpack("getEnygmaManager", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetKeyAgreements is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0d30d092.
//
// Solidity: function getKeyAgreements(uint256 chainId) view returns((uint256,bytes,bytes,uint256)[])
func (participantStorageV1 *ParticipantStorageV1) PackGetKeyAgreements(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("getKeyAgreements", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetKeyAgreements is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0d30d092.
//
// Solidity: function getKeyAgreements(uint256 chainId) view returns((uint256,bytes,bytes,uint256)[])
func (participantStorageV1 *ParticipantStorageV1) UnpackGetKeyAgreements(data []byte) ([]ParticipantStructsKeyAgreementData, error) {
	out, err := participantStorageV1.abi.Unpack("getKeyAgreements", data)
	if err != nil {
		return *new([]ParticipantStructsKeyAgreementData), err
	}
	out0 := *abi.ConvertType(out[0], new([]ParticipantStructsKeyAgreementData)).(*[]ParticipantStructsKeyAgreementData)
	return out0, err
}

// PackGetParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1b9db2ef.
//
// Solidity: function getParticipant(uint256 chainId) view returns((uint256,uint8,uint8,string,string,uint256,uint256,bool))
func (participantStorageV1 *ParticipantStorageV1) PackGetParticipant(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("getParticipant", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetParticipant is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1b9db2ef.
//
// Solidity: function getParticipant(uint256 chainId) view returns((uint256,uint8,uint8,string,string,uint256,uint256,bool))
func (participantStorageV1 *ParticipantStorageV1) UnpackGetParticipant(data []byte) (ParticipantStructsParticipant, error) {
	out, err := participantStorageV1.abi.Unpack("getParticipant", data)
	if err != nil {
		return *new(ParticipantStructsParticipant), err
	}
	out0 := *abi.ConvertType(out[0], new(ParticipantStructsParticipant)).(*ParticipantStructsParticipant)
	return out0, err
}

// PackGetParticipantCore is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x594c5c76.
//
// Solidity: function getParticipantCore() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) PackGetParticipantCore() []byte {
	enc, err := participantStorageV1.abi.Pack("getParticipantCore")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetParticipantCore is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x594c5c76.
//
// Solidity: function getParticipantCore() view returns(address)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetParticipantCore(data []byte) (common.Address, error) {
	out, err := participantStorageV1.abi.Unpack("getParticipantCore", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetParticipantDataBatch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa0842e33.
//
// Solidity: function getParticipantDataBatch() view returns((uint256,string,uint256)[] pnViewData, (uint256,string,bytes,bytes,uint256)[] auditInfo, uint256[] pnChainIds, uint256[] auditChainIds)
func (participantStorageV1 *ParticipantStorageV1) PackGetParticipantDataBatch() []byte {
	enc, err := participantStorageV1.abi.Pack("getParticipantDataBatch")
	if err != nil {
		panic(err)
	}
	return enc
}

// GetParticipantDataBatchOutput serves as a container for the return parameters of contract
// method GetParticipantDataBatch.
type GetParticipantDataBatchOutput struct {
	PnViewData    []ParticipantStructsPrivacyNodeViewData
	AuditInfo     []ParticipantStructsAuditInfoData
	PnChainIds    []*big.Int
	AuditChainIds []*big.Int
}

// UnpackGetParticipantDataBatch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0842e33.
//
// Solidity: function getParticipantDataBatch() view returns((uint256,string,uint256)[] pnViewData, (uint256,string,bytes,bytes,uint256)[] auditInfo, uint256[] pnChainIds, uint256[] auditChainIds)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetParticipantDataBatch(data []byte) (GetParticipantDataBatchOutput, error) {
	out, err := participantStorageV1.abi.Unpack("getParticipantDataBatch", data)
	outstruct := new(GetParticipantDataBatchOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.PnViewData = *abi.ConvertType(out[0], new([]ParticipantStructsPrivacyNodeViewData)).(*[]ParticipantStructsPrivacyNodeViewData)
	outstruct.AuditInfo = *abi.ConvertType(out[1], new([]ParticipantStructsAuditInfoData)).(*[]ParticipantStructsAuditInfoData)
	outstruct.PnChainIds = *abi.ConvertType(out[2], new([]*big.Int)).(*[]*big.Int)
	outstruct.AuditChainIds = *abi.ConvertType(out[3], new([]*big.Int)).(*[]*big.Int)
	return *outstruct, err

}

// PackGetPaymentSpendPublicKeyByChainId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x81a40de6.
//
// Solidity: function getPaymentSpendPublicKeyByChainId(uint256 chainId) view returns(uint256)
func (participantStorageV1 *ParticipantStorageV1) PackGetPaymentSpendPublicKeyByChainId(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("getPaymentSpendPublicKeyByChainId", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPaymentSpendPublicKeyByChainId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x81a40de6.
//
// Solidity: function getPaymentSpendPublicKeyByChainId(uint256 chainId) view returns(uint256)
func (participantStorageV1 *ParticipantStorageV1) UnpackGetPaymentSpendPublicKeyByChainId(data []byte) (*big.Int, error) {
	out, err := participantStorageV1.abi.Unpack("getPaymentSpendPublicKeyByChainId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x485cc955.
//
// Solidity: function initialize(address _endpoint, address authority_) returns()
func (participantStorageV1 *ParticipantStorageV1) PackInitialize(endpoint common.Address, authority common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("initialize", endpoint, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackInitialize0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4d66de8.
//
// Solidity: function initialize(address _endpoint) returns()
func (participantStorageV1 *ParticipantStorageV1) PackInitialize0(endpoint common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("initialize0", endpoint)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackInitiateKeyAgreement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4ec53b2e.
//
// Solidity: function initiateKeyAgreement(uint256 initiatorChainId, uint256 responderChainId, bytes ciphertext, bytes digest, uint256 blockNumber) returns()
func (participantStorageV1 *ParticipantStorageV1) PackInitiateKeyAgreement(initiatorChainId *big.Int, responderChainId *big.Int, ciphertext []byte, digest []byte, blockNumber *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("initiateKeyAgreement", initiatorChainId, responderChainId, ciphertext, digest, blockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (participantStorageV1 *ParticipantStorageV1) PackProxiableUUID() []byte {
	enc, err := participantStorageV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (participantStorageV1 *ParticipantStorageV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := participantStorageV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackRemoveParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x683f7f27.
//
// Solidity: function removeParticipant(uint256 chainId) returns()
func (participantStorageV1 *ParticipantStorageV1) PackRemoveParticipant(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("removeParticipant", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (participantStorageV1 *ParticipantStorageV1) PackResourceId() []byte {
	enc, err := participantStorageV1.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (participantStorageV1 *ParticipantStorageV1) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := participantStorageV1.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSetAuditInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6279aa9f.
//
// Solidity: function setAuditInfo(uint256 chainId, string raylsViewPublicKey, bytes encryptedRaylsViewPrivateKey, bytes mac, uint256 blockNumber) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetAuditInfo(chainId *big.Int, raylsViewPublicKey string, encryptedRaylsViewPrivateKey []byte, mac []byte, blockNumber *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("setAuditInfo", chainId, raylsViewPublicKey, encryptedRaylsViewPrivateKey, mac, blockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetAuditManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc68bab0f.
//
// Solidity: function setAuditManager(address _auditManager) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetAuditManager(auditManager common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("setAuditManager", auditManager)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetChainViewData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4017734d.
//
// Solidity: function setChainViewData(uint256 chainId, string raylsViewPublicKey, uint256 blockNumber) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetChainViewData(chainId *big.Int, raylsViewPublicKey string, blockNumber *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("setChainViewData", chainId, raylsViewPublicKey, blockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetEnygmaManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xede591fd.
//
// Solidity: function setEnygmaManager(address _enygmaManager) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetEnygmaManager(enygmaManager common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("setEnygmaManager", enygmaManager)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetEnygmaPnEventsAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52bf1c8d.
//
// Solidity: function setEnygmaPnEventsAddress(address _pnEnygmaEvents) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetEnygmaPnEventsAddress(pnEnygmaEvents common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("setEnygmaPnEventsAddress", pnEnygmaEvents)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetParticipantCore is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x035c8713.
//
// Solidity: function setParticipantCore(address _participantCore) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetParticipantCore(participantCore common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("setParticipantCore", participantCore)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetPaymentSpendPublicKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x079c2a49.
//
// Solidity: function setPaymentSpendPublicKey(uint256 _chainId, uint256 _paymentSpendPublicKey, address[] _pnAddresses) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetPaymentSpendPublicKey(chainId *big.Int, paymentSpendPublicKey *big.Int, pnAddresses []common.Address) []byte {
	enc, err := participantStorageV1.abi.Pack("setPaymentSpendPublicKey", chainId, paymentSpendPublicKey, pnAddresses)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa01afbfb.
//
// Solidity: function setResourceId(bytes32 _resourceId) returns()
func (participantStorageV1 *ParticipantStorageV1) PackSetResourceId(resourceId [32]byte) []byte {
	enc, err := participantStorageV1.abi.Pack("setResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdateBroadcastMessagesPermission is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x33ded4a8.
//
// Solidity: function updateBroadcastMessagesPermission(uint256 chainId, bool allowed) returns()
func (participantStorageV1 *ParticipantStorageV1) PackUpdateBroadcastMessagesPermission(chainId *big.Int, allowed bool) []byte {
	enc, err := participantStorageV1.abi.Pack("updateBroadcastMessagesPermission", chainId, allowed)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdateRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x65ffa17e.
//
// Solidity: function updateRole(uint256 chainId, uint8 role) returns()
func (participantStorageV1 *ParticipantStorageV1) PackUpdateRole(chainId *big.Int, role uint8) []byte {
	enc, err := participantStorageV1.abi.Pack("updateRole", chainId, role)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdateStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3a1b3d31.
//
// Solidity: function updateStatus(uint256 chainId, uint8 status) returns()
func (participantStorageV1 *ParticipantStorageV1) PackUpdateStatus(chainId *big.Int, status uint8) []byte {
	enc, err := participantStorageV1.abi.Pack("updateStatus", chainId, status)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (participantStorageV1 *ParticipantStorageV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := participantStorageV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackValidateMessageParticipants is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd5c3614f.
//
// Solidity: function validateMessageParticipants(uint256 originChainId, uint256 destinationChainId) view returns()
func (participantStorageV1 *ParticipantStorageV1) PackValidateMessageParticipants(originChainId *big.Int, destinationChainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("validateMessageParticipants", originChainId, destinationChainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackValidateParticipantStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9557f56.
//
// Solidity: function validateParticipantStatus(uint256 chainId) view returns()
func (participantStorageV1 *ParticipantStorageV1) PackValidateParticipantStatus(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("validateParticipantStatus", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackVerifyParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb3816fdb.
//
// Solidity: function verifyParticipant(uint256 chainId) view returns(bool)
func (participantStorageV1 *ParticipantStorageV1) PackVerifyParticipant(chainId *big.Int) []byte {
	enc, err := participantStorageV1.abi.Pack("verifyParticipant", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVerifyParticipant is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb3816fdb.
//
// Solidity: function verifyParticipant(uint256 chainId) view returns(bool)
func (participantStorageV1 *ParticipantStorageV1) UnpackVerifyParticipant(data []byte) (bool, error) {
	out, err := participantStorageV1.abi.Unpack("verifyParticipant", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// ParticipantStorageV1AuthorityUpdated represents a AuthorityUpdated event raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const ParticipantStorageV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (ParticipantStorageV1AuthorityUpdated) ContractEventName() string {
	return ParticipantStorageV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (participantStorageV1 *ParticipantStorageV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*ParticipantStorageV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != participantStorageV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ParticipantStorageV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := participantStorageV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range participantStorageV1.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ParticipantStorageV1Initialized represents a Initialized event raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const ParticipantStorageV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (ParticipantStorageV1Initialized) ContractEventName() string {
	return ParticipantStorageV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (participantStorageV1 *ParticipantStorageV1) UnpackInitializedEvent(log *types.Log) (*ParticipantStorageV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != participantStorageV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ParticipantStorageV1Initialized)
	if len(log.Data) > 0 {
		if err := participantStorageV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range participantStorageV1.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ParticipantStorageV1ModulesConfigured represents a ModulesConfigured event raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1ModulesConfigured struct {
	ParticipantCore common.Address
	AuditManager    common.Address
	EnygmaManager   common.Address
	Raw             *types.Log // Blockchain specific contextual infos
}

const ParticipantStorageV1ModulesConfiguredEventName = "ModulesConfigured"

// ContractEventName returns the user-defined event name.
func (ParticipantStorageV1ModulesConfigured) ContractEventName() string {
	return ParticipantStorageV1ModulesConfiguredEventName
}

// UnpackModulesConfiguredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ModulesConfigured(address indexed participantCore, address indexed auditManager, address indexed enygmaManager)
func (participantStorageV1 *ParticipantStorageV1) UnpackModulesConfiguredEvent(log *types.Log) (*ParticipantStorageV1ModulesConfigured, error) {
	event := "ModulesConfigured"
	if log.Topics[0] != participantStorageV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ParticipantStorageV1ModulesConfigured)
	if len(log.Data) > 0 {
		if err := participantStorageV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range participantStorageV1.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ParticipantStorageV1Upgraded represents a Upgraded event raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const ParticipantStorageV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (ParticipantStorageV1Upgraded) ContractEventName() string {
	return ParticipantStorageV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (participantStorageV1 *ParticipantStorageV1) UnpackUpgradedEvent(log *types.Log) (*ParticipantStorageV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != participantStorageV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ParticipantStorageV1Upgraded)
	if len(log.Data) > 0 {
		if err := participantStorageV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range participantStorageV1.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (participantStorageV1 *ParticipantStorageV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["ParticipantStorageV1UnauthorizedCaller"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackParticipantStorageV1UnauthorizedCallerError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAppV1HubNotActive"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAppV1HubNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAppV1PrivacyNodeFrozen"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAppV1PrivacyNodeFrozenError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAppV1PrivacyNodeNotActive"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAppV1PrivacyNodeNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAppV1PublicChainNotActive"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAppV1PublicChainNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAppV1ResourceNotApproved"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAppV1ResourceNotApprovedError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAppV1TokenRegistryNotConfigured"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAppV1TokenRegistryNotConfiguredError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["RaylsAppV1UnauthorizedTokenRegistry"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackRaylsAppV1UnauthorizedTokenRegistryError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], participantStorageV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return participantStorageV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ParticipantStorageV1AddressEmptyCode represents a AddressEmptyCode error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func ParticipantStorageV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (participantStorageV1 *ParticipantStorageV1) UnpackAddressEmptyCodeError(raw []byte) (*ParticipantStorageV1AddressEmptyCode, error) {
	out := new(ParticipantStorageV1AddressEmptyCode)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func ParticipantStorageV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (participantStorageV1 *ParticipantStorageV1) UnpackERC1967InvalidImplementationError(raw []byte) (*ParticipantStorageV1ERC1967InvalidImplementation, error) {
	out := new(ParticipantStorageV1ERC1967InvalidImplementation)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func ParticipantStorageV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (participantStorageV1 *ParticipantStorageV1) UnpackERC1967NonPayableError(raw []byte) (*ParticipantStorageV1ERC1967NonPayable, error) {
	out := new(ParticipantStorageV1ERC1967NonPayable)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1FailedCall represents a FailedCall error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func ParticipantStorageV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (participantStorageV1 *ParticipantStorageV1) UnpackFailedCallError(raw []byte) (*ParticipantStorageV1FailedCall, error) {
	out := new(ParticipantStorageV1FailedCall)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1InvalidInitialization represents a InvalidInitialization error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func ParticipantStorageV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (participantStorageV1 *ParticipantStorageV1) UnpackInvalidInitializationError(raw []byte) (*ParticipantStorageV1InvalidInitialization, error) {
	out := new(ParticipantStorageV1InvalidInitialization)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1NotInitializing represents a NotInitializing error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func ParticipantStorageV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (participantStorageV1 *ParticipantStorageV1) UnpackNotInitializingError(raw []byte) (*ParticipantStorageV1NotInitializing, error) {
	out := new(ParticipantStorageV1NotInitializing)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1ParticipantStorageV1UnauthorizedCaller represents a ParticipantStorageV1__UnauthorizedCaller error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1ParticipantStorageV1UnauthorizedCaller struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ParticipantStorageV1__UnauthorizedCaller(address caller)
func ParticipantStorageV1ParticipantStorageV1UnauthorizedCallerErrorID() common.Hash {
	return common.HexToHash("0x462e5c5b72c0fc8d71c0efe37e91aba3df63f77db44f129a892a044de14a6df8")
}

// UnpackParticipantStorageV1UnauthorizedCallerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ParticipantStorageV1__UnauthorizedCaller(address caller)
func (participantStorageV1 *ParticipantStorageV1) UnpackParticipantStorageV1UnauthorizedCallerError(raw []byte) (*ParticipantStorageV1ParticipantStorageV1UnauthorizedCaller, error) {
	out := new(ParticipantStorageV1ParticipantStorageV1UnauthorizedCaller)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "ParticipantStorageV1UnauthorizedCaller", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func ParticipantStorageV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*ParticipantStorageV1RaylsAccessManagedContractPaused, error) {
	out := new(ParticipantStorageV1RaylsAccessManagedContractPaused)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func ParticipantStorageV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*ParticipantStorageV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(ParticipantStorageV1RaylsAccessManagedInvalidAuthority)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func ParticipantStorageV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*ParticipantStorageV1RaylsAccessManagedMustSchedule, error) {
	out := new(ParticipantStorageV1RaylsAccessManagedMustSchedule)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func ParticipantStorageV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*ParticipantStorageV1RaylsAccessManagedUnauthorized, error) {
	out := new(ParticipantStorageV1RaylsAccessManagedUnauthorized)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAppV1HubNotActive represents a RaylsAppV1__HubNotActive error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAppV1HubNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	HubStatus         uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func ParticipantStorageV1RaylsAppV1HubNotActiveErrorID() common.Hash {
	return common.HexToHash("0x3fae5bbd70277aa1cd008dceb93b19a7055c2a6d29b84733371e1c41b2048b15")
}

// UnpackRaylsAppV1HubNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAppV1HubNotActiveError(raw []byte) (*ParticipantStorageV1RaylsAppV1HubNotActive, error) {
	out := new(ParticipantStorageV1RaylsAppV1HubNotActive)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAppV1HubNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAppV1PrivacyNodeFrozen represents a RaylsAppV1__PrivacyNodeFrozen error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAppV1PrivacyNodeFrozen struct {
	TokenAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__PrivacyNodeFrozen(address tokenAddress)
func ParticipantStorageV1RaylsAppV1PrivacyNodeFrozenErrorID() common.Hash {
	return common.HexToHash("0xc80bd255e67000277f5aed4960b64f92e2d5a652f07a22fba7d044de6add8f0e")
}

// UnpackRaylsAppV1PrivacyNodeFrozenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__PrivacyNodeFrozen(address tokenAddress)
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAppV1PrivacyNodeFrozenError(raw []byte) (*ParticipantStorageV1RaylsAppV1PrivacyNodeFrozen, error) {
	out := new(ParticipantStorageV1RaylsAppV1PrivacyNodeFrozen)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAppV1PrivacyNodeFrozen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAppV1PrivacyNodeNotActive represents a RaylsAppV1__PrivacyNodeNotActive error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAppV1PrivacyNodeNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func ParticipantStorageV1RaylsAppV1PrivacyNodeNotActiveErrorID() common.Hash {
	return common.HexToHash("0xfdcdd2a6e576bf1f342ce493560565ef686a59cd3e0486f6869151efb2c7853f")
}

// UnpackRaylsAppV1PrivacyNodeNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAppV1PrivacyNodeNotActiveError(raw []byte) (*ParticipantStorageV1RaylsAppV1PrivacyNodeNotActive, error) {
	out := new(ParticipantStorageV1RaylsAppV1PrivacyNodeNotActive)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAppV1PrivacyNodeNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAppV1PublicChainNotActive represents a RaylsAppV1__PublicChainNotActive error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAppV1PublicChainNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	PublicChainStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func ParticipantStorageV1RaylsAppV1PublicChainNotActiveErrorID() common.Hash {
	return common.HexToHash("0xb607961611e6e4126e09c80bcd1e35e7a1e886888daa292eecc27cd9d4e37f3f")
}

// UnpackRaylsAppV1PublicChainNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAppV1PublicChainNotActiveError(raw []byte) (*ParticipantStorageV1RaylsAppV1PublicChainNotActive, error) {
	out := new(ParticipantStorageV1RaylsAppV1PublicChainNotActive)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAppV1PublicChainNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAppV1ResourceNotApproved represents a RaylsAppV1__ResourceNotApproved error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAppV1ResourceNotApproved struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__ResourceNotApproved()
func ParticipantStorageV1RaylsAppV1ResourceNotApprovedErrorID() common.Hash {
	return common.HexToHash("0x8f144935367c131b72d26b0320b764f69ba3639e65abb1c811084bbd46e5c731")
}

// UnpackRaylsAppV1ResourceNotApprovedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__ResourceNotApproved()
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAppV1ResourceNotApprovedError(raw []byte) (*ParticipantStorageV1RaylsAppV1ResourceNotApproved, error) {
	out := new(ParticipantStorageV1RaylsAppV1ResourceNotApproved)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAppV1ResourceNotApproved", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAppV1TokenRegistryNotConfigured represents a RaylsAppV1__TokenRegistryNotConfigured error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAppV1TokenRegistryNotConfigured struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__TokenRegistryNotConfigured()
func ParticipantStorageV1RaylsAppV1TokenRegistryNotConfiguredErrorID() common.Hash {
	return common.HexToHash("0x3eba255b70fc7afd9cc5be90de2023dae8350ac3c29cbd5eaf139cadd9c4292e")
}

// UnpackRaylsAppV1TokenRegistryNotConfiguredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__TokenRegistryNotConfigured()
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAppV1TokenRegistryNotConfiguredError(raw []byte) (*ParticipantStorageV1RaylsAppV1TokenRegistryNotConfigured, error) {
	out := new(ParticipantStorageV1RaylsAppV1TokenRegistryNotConfigured)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAppV1TokenRegistryNotConfigured", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1RaylsAppV1UnauthorizedTokenRegistry represents a RaylsAppV1__UnauthorizedTokenRegistry error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1RaylsAppV1UnauthorizedTokenRegistry struct {
	Caller        common.Address
	TokenRegistry common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func ParticipantStorageV1RaylsAppV1UnauthorizedTokenRegistryErrorID() common.Hash {
	return common.HexToHash("0x000d23e5a298a9951b289bd8f5eece62aa717c000d6b0509a9f77d16f67a5b7d")
}

// UnpackRaylsAppV1UnauthorizedTokenRegistryError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func (participantStorageV1 *ParticipantStorageV1) UnpackRaylsAppV1UnauthorizedTokenRegistryError(raw []byte) (*ParticipantStorageV1RaylsAppV1UnauthorizedTokenRegistry, error) {
	out := new(ParticipantStorageV1RaylsAppV1UnauthorizedTokenRegistry)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "RaylsAppV1UnauthorizedTokenRegistry", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func ParticipantStorageV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (participantStorageV1 *ParticipantStorageV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*ParticipantStorageV1UUPSUnauthorizedCallContext, error) {
	out := new(ParticipantStorageV1UUPSUnauthorizedCallContext)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ParticipantStorageV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the ParticipantStorageV1 contract.
type ParticipantStorageV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func ParticipantStorageV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (participantStorageV1 *ParticipantStorageV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*ParticipantStorageV1UUPSUnsupportedProxiableUUID, error) {
	out := new(ParticipantStorageV1UUPSUnsupportedProxiableUUID)
	if err := participantStorageV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
