// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TokenRegistryV1

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

// SharedObjectsTokenRegistrationData is an auto generated low-level Go binding around an user-defined struct.
type SharedObjectsTokenRegistrationData struct {
	Name              string
	Symbol            string
	Uri               string
	TotalSupply       []byte
	IssuerChainId     *big.Int
	PnRegistryAddress common.Address
	Bytecode          []byte
	InitializerParams []byte
	IsFungible        bool
	ErcStandard       uint8
	IsCustom          bool
	TokenAddress      common.Address
}

// TokenStructsFrozenToken is an auto generated low-level Go binding around an user-defined struct.
type TokenStructsFrozenToken struct {
	ResourceId         [32]byte
	FrozenParticipants []*big.Int
}

// TokenStructsToken is an auto generated low-level Go binding around an user-defined struct.
type TokenStructsToken struct {
	ResourceId                  [32]byte
	Name                        string
	Symbol                      string
	IssuerChainId               *big.Int
	IssuerImplementationAddress common.Address
	PnRegistryAddress           common.Address
	TokenAddress                common.Address
	IsFungible                  bool
	Status                      uint8
	CreatedAt                   *big.Int
	UpdatedAt                   *big.Int
	Metadata                    TokenStructsTokenMetadata
	ErcStandard                 uint8
}

// TokenStructsTokenMetadata is an auto generated low-level Go binding around an user-defined struct.
type TokenStructsTokenMetadata struct {
	Url      string
	Decimals uint8
}

// TokenRegistryV1MetaData contains all meta data concerning the TokenRegistryV1 contract.
var TokenRegistryV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addToken\",\"inputs\":[{\"name\":\"tokenData\",\"type\":\"tuple\",\"internalType\":\"structSharedObjects.TokenRegistrationData\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"totalSupply\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"issuerChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pnRegistryAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bytecode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"initializerParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"isFungible\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"ercStandard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"},{\"name\":\"isCustom\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"broadcastCurrentFrozenResourcesForNewParticipant\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"broadcastFrozenToken\",\"inputs\":[{\"name\":\"frozenToken\",\"type\":\"tuple\",\"internalType\":\"structTokenStructs.FrozenToken\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"frozenParticipants\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"broadcastUnfrozenToken\",\"inputs\":[{\"name\":\"unfrozenToken\",\"type\":\"tuple\",\"internalType\":\"structTokenStructs.FrozenToken\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"frozenParticipants\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"configureModules\",\"inputs\":[{\"name\":\"_tokenCore\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_tokenFreezeManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_enygmaTokenManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"contractVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"endpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRaylsEndpoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"freezeToken\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"chainIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllTokens\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structTokenStructs.Token[]\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"issuerChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"issuerImplementationAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pnRegistryAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isFungible\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.TokenStatus\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"metadata\",\"type\":\"tuple\",\"internalType\":\"structTokenStructs.TokenMetadata\",\"components\":[{\"name\":\"url\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"ercStandard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEnygmaTokenManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEnygmaTokenManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenByResourceId\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structTokenStructs.Token\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"issuerChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"issuerImplementationAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pnRegistryAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isFungible\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.TokenStatus\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"metadata\",\"type\":\"tuple\",\"internalType\":\"structTokenStructs.TokenMetadata\",\"components\":[{\"name\":\"url\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"ercStandard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenCore\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractITokenCore\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenFreezeManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractITokenFreezeManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isTokenFrozenForParticipant\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setEnygmaFactory\",\"inputs\":[{\"name\":\"_enygmaFactory\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEnygmaTokenManager\",\"inputs\":[{\"name\":\"_enygmaTokenManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTokenCore\",\"inputs\":[{\"name\":\"_tokenCore\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTokenFreezeManager\",\"inputs\":[{\"name\":\"_tokenFreezeManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unfreezeToken\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"chainIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateStatus\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumTokenStructs.TokenStatus\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateTokenBalance\",\"inputs\":[{\"name\":\"issuerChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"updateType\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.BalanceUpdateType\"},{\"name\":\"metadata\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ModulesConfigured\",\"inputs\":[{\"name\":\"tokenCore\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenFreezeManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"enygmaTokenManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAppV1__HubNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"hubStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsAppV1__PrivacyNodeFrozen\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAppV1__PrivacyNodeNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsAppV1__PublicChainNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"publicChainStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsAppV1__ResourceNotApproved\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAppV1__TokenRegistryNotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAppV1__UnauthorizedTokenRegistry\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "TokenRegistryV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b5060805161266d61003e600039600081816111660152818161118f01526112c6015261266d6000f3fe60806040526004361061016d5760003560e01c806372d14e30116100c757806372d14e301461035d57806380387be21461037d5780638a58b79f1461039d578063a01afbfb146103ca578063a0a8e460146103ea578063a0b3f1ee146103fe578063ad3cb1cc1461041c578063bf7e214f1461045a578063c4d66de81461046f578063e08ac15e1461048f578063e66cdd2b146104af578063f0a599bb146104cf578063f356cc03146104ed578063f7bb5fad1461050b578063f867cd521461052b57600080fd5b8063054372ed146101725780630566585a1461019457806311f50c85146101b457806313aa1a1d146101ea57806314ae646b1461020a578063153f716d1461022a5780632a5c792a1461024a57806331b1125f1461026c578063338042e41461028c578063485cc955146102bc57806349f6ee3c146102dc5780634f1ef286146102f157806352d1902d146103045780635e280f11146103275780635f997c5b14610347575b600080fd5b34801561017e57600080fd5b5061019261018d3660046116c7565b61054b565b005b3480156101a057600080fd5b506101926101af36600461181a565b6105c9565b3480156101c057600080fd5b506101d46101cf366004611890565b610644565b6040516101e191906118b6565b60405180910390f35b3480156101f657600080fd5b506101926102053660046118ea565b6106b8565b34801561021657600080fd5b506101926102253660046118ea565b61071f565b34801561023657600080fd5b506101926102453660046118ea565b61077d565b34801561025657600080fd5b5061025f6107c2565b6040516101e19190611ab5565b34801561027857600080fd5b50610192610287366004611b19565b610839565b34801561029857600080fd5b506102ac6102a7366004611b64565b61092f565b60405190151581526020016101e1565b3480156102c857600080fd5b506101926102d7366004611b86565b6109ab565b3480156102e857600080fd5b50610192610ad7565b6101926102ff366004611c2c565b610b58565b34801561031057600080fd5b50610319610b77565b6040519081526020016101e1565b34801561033357600080fd5b506000546101d4906001600160a01b031681565b34801561035357600080fd5b5061031960015481565b34801561036957600080fd5b50610319610378366004611c7b565b610b95565b34801561038957600080fd5b50610192610398366004611cb6565b610c2c565b3480156103a957600080fd5b506103bd6103b8366004611890565b610cb0565b6040516101e19190611d1c565b3480156103d657600080fd5b506101926103e5366004611890565b610d9f565b3480156103f657600080fd5b506001610319565b34801561040a57600080fd5b506004546001600160a01b03166101d4565b34801561042857600080fd5b5061044d604051806040016040528060058152602001640352e302e360dc1b81525081565b6040516101e19190611d2f565b34801561046657600080fd5b506101d4610dea565b34801561047b57600080fd5b5061019261048a3660046118ea565b610e03565b34801561049b57600080fd5b506101926104aa366004611d42565b610e2d565b3480156104bb57600080fd5b506101926104ca36600461181a565b610ea5565b3480156104db57600080fd5b506002546001600160a01b03166101d4565b3480156104f957600080fd5b506003546001600160a01b03166101d4565b34801561051757600080fd5b50610192610526366004611dc0565b610eeb565b34801561053757600080fd5b506101926105463660046118ea565b610f33565b610561336000356001600160e01b031916610f91565b60025460405163054372ed60e01b81526001600160a01b039091169063054372ed906105939085908590600401611dfc565b600060405180830381600087803b1580156105ad57600080fd5b505af11580156105c1573d6000803e3d6000fd5b505050505050565b6105df336000356001600160e01b031916610f91565b6003546040516302b32c2d60e11b81526001600160a01b0390911690630566585a9061060f908490600401611e4c565b600060405180830381600087803b15801561062957600080fd5b505af115801561063d573d6000803e3d6000fd5b5050505050565b600080546040516311f50c8560e01b8152600481018490526001600160a01b03909116906311f50c8590602401602060405180830381865afa15801561068e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106b29190611e84565b92915050565b6106ce336000356001600160e01b031916610f91565b6001600160a01b0381166106fd5760405162461bcd60e51b81526004016106f490611ea1565b60405180910390fd5b600380546001600160a01b0319166001600160a01b0392909216919091179055565b610735336000356001600160e01b031916610f91565b6001600160a01b03811661075b5760405162461bcd60e51b81526004016106f490611efb565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b610793336000356001600160e01b031916610f91565b6004805460405163153f716d60e01b81526001600160a01b039091169163153f716d9161060f918591016118b6565b6002546040805163152e3c9560e11b815290516060926001600160a01b031691632a5c792a9160048083019260009291908290030181865afa15801561080c573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610834919081019061215b565b905090565b61084f336000356001600160e01b031916610f91565b6001600160a01b0383166108755760405162461bcd60e51b81526004016106f490611efb565b6001600160a01b03821661089b5760405162461bcd60e51b81526004016106f490611ea1565b6001600160a01b0381166108c15760405162461bcd60e51b81526004016106f49061220b565b600280546001600160a01b03199081166001600160a01b03868116918217909355600380548316868516908117909155600480549093169385169384179092556040517ff6c27e15b7e995e3edd056f3e9b7b01098dfe3f91cccf2af78ff33215fc1829d90600090a4505050565b600354604051630ce010b960e21b815260048101849052602481018390526000916001600160a01b03169063338042e490604401602060405180830381865afa158015610980573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109a49190612265565b9392505050565b60006109b56110d3565b805490915060ff600160401b82041615906001600160401b03166000811580156109dc5750825b90506000826001600160401b031660011480156109f85750303b155b905081158015610a06575080155b15610a245760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff191660011785558315610a4e57845460ff60401b1916600160401b1785555b610a566110fc565b610a5f87610e03565b600080546001600160a01b0319166001600160a01b0389161790556003600155610a8886611106565b8315610ace57845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b610aed336000356001600160e01b031916610f91565b6003546001600160a01b031663054e3537610b06611147565b6040518263ffffffff1660e01b8152600401610b2491815260200190565b600060405180830381600087803b158015610b3e57600080fd5b505af1158015610b52573d6000803e3d6000fd5b50505050565b610b6061115b565b610b69826111e9565b610b738282611202565b5050565b6000610b816112bb565b506000805160206126188339815191525b90565b6000610bad336000356001600160e01b031916610f91565b6002546001600160a01b031663c711ce0883610bc7611304565b6040518363ffffffff1660e01b8152600401610be492919061230d565b6020604051808303816000875af1158015610c03573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106b291906124a9565b919050565b610c42336000356001600160e01b031916610f91565b60025460405163401c3df160e11b81526001600160a01b03909116906380387be290610c789087908790879087906004016124c2565b600060405180830381600087803b158015610c9257600080fd5b505af1158015610ca6573d6000803e3d6000fd5b5050505050505050565b610d2e604080516101a0810182526000808252606060208084018290528385018290528184018390526080840183905260a0840183905260c0840183905260e0840183905261010084018390526101208401839052610140840183905284518086019095529084528301529061016082019081526020016000905290565b600254604051638a58b79f60e01b8152600481018490526001600160a01b0390911690638a58b79f90602401600060405180830381865afa158015610d77573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526106b29190810190612501565b6000610da961131a565b9050336001600160a01b03821614610de457604051620d23e560e01b81523360048201526001600160a01b03821660248201526044016106f4565b50600155565b6000610df46113b1565b546001600160a01b0316919050565b610e0b611413565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b610e43336000356001600160e01b031916610f91565b60035460405163704560af60e11b81526001600160a01b039091169063e08ac15e90610e7790869086908690600401612535565b600060405180830381600087803b158015610e9157600080fd5b505af1158015610ace573d6000803e3d6000fd5b610ebb336000356001600160e01b031916610f91565b60035460405163e66cdd2b60e01b81526001600160a01b039091169063e66cdd2b9061060f908490600401611e4c565b610f01336000356001600160e01b031916610f91565b60035460405163f7bb5fad60e01b81526001600160a01b039091169063f7bb5fad906105939085908590600401612576565b610f49336000356001600160e01b031916610f91565b6001600160a01b038116610f6f5760405162461bcd60e51b81526004016106f49061220b565b600480546001600160a01b0319166001600160a01b0392909216919091179055565b6000610f9b6113b1565b80549091506001600160a01b031680610fca576000604051638944034760e01b81526004016106f491906118b6565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa15801561102e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611052919061258f565b92509250925082610ace57801561107c5760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff8216156110b85760405163a426878960e01b81526001600160a01b038816600482015263ffffffff831660248201526044016106f4565b86604051632ecd3d0360e21b81526004016106f491906118b6565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a006106b2565b611104611413565b565b60006111106113b1565b80549091506001600160a01b03161561113e5781604051638944034760e01b81526004016106f491906118b6565b610b7382611438565b600060343610610b92575060331936013590565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614806111cb57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166111bf6114c8565b6001600160a01b031614155b156111045760405163703e46dd60e11b815260040160405180910390fd5b6111ff336000356001600160e01b031916610f91565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa92505050801561125c575060408051601f3d908101601f19168201909252611259918101906124a9565b60015b61127b5781604051634c9c8ce360e01b81526004016106f491906118b6565b60008051602061261883398151915281146112ac57604051632a87526960e21b8152600481018290526024016106f4565b6112b683836114de565b505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146111045760405163703e46dd60e11b815260040160405180910390fd5b3360143610610b92575060131936013560601c90565b600080546040516311f50c8560e01b8152600360048201526001600160a01b03909116906311f50c8590602401602060405180830381865afa158015611364573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113889190611e84565b90506001600160a01b038116610b9257604051633eba255b60e01b815260040160405180910390fd5b60008060ff196113e260017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f356125da565b6040516020016113f491815260200190565b60408051601f1981840301815291905280516020909101201692915050565b61141b611534565b61110457604051631afcd79f60e31b815260040160405180910390fd5b6001600160a01b0381166114615780604051638944034760e01b81526004016106f491906118b6565b600061146b6113b1565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b6000600080516020612618833981519152610df4565b6114e78261154e565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a280511561152c576112b682826115aa565b610b73611620565b600061153e6110d3565b54600160401b900460ff16919050565b806001600160a01b03163b60000361157b5780604051634c9c8ce360e01b81526004016106f491906118b6565b60008051602061261883398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b0316846040516115c791906125fb565b600060405180830381855af49150503d8060008114611602576040519150601f19603f3d011682016040523d82523d6000602084013e611607565b606091505b509150915061161785838361163f565b95945050505050565b34156111045760405163b398979f60e01b815260040160405180910390fd5b6060826116545761164f82611692565b6109a4565b815115801561166b57506001600160a01b0384163b155b1561168b5783604051639996b31560e01b81526004016106f491906118b6565b5092915050565b8051156116a157805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b600381106111ff57600080fd5b600080604083850312156116da57600080fd5b8235915060208301356116ec816116ba565b809150509250929050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b038111828210171561172f5761172f6116f7565b60405290565b6040516101a081016001600160401b038111828210171561172f5761172f6116f7565b604051601f8201601f191681016001600160401b0381118282101715611780576117806116f7565b604052919050565b60006001600160401b038211156117a1576117a16116f7565b5060051b60200190565b600082601f8301126117bc57600080fd5b813560206117d16117cc83611788565b611758565b8083825260208201915060208460051b8701019350868411156117f357600080fd5b602086015b8481101561180f57803583529183019183016117f8565b509695505050505050565b60006020828403121561182c57600080fd5b81356001600160401b038082111561184357600080fd5b908301906040828603121561185757600080fd5b61185f61170d565b8235815260208301358281111561187557600080fd5b611881878286016117ab565b60208301525095945050505050565b6000602082840312156118a257600080fd5b5035919050565b6001600160a01b03169052565b6001600160a01b0391909116815260200190565b6001600160a01b03811681146111ff57600080fd5b8035610c27816118ca565b6000602082840312156118fc57600080fd5b81356109a4816118ca565b60005b8381101561192257818101518382015260200161190a565b50506000910152565b60008151808452611943816020860160208601611907565b601f01601f19169290920160200192915050565b634e487b7160e01b600052602160045260246000fd5b6003811061197d5761197d611957565b9052565b6000815160408452611996604085018261192b565b60209384015160ff16949093019390935250919050565b600d811061197d5761197d611957565b60006101a08251845260208301518160208601526119dd8286018261192b565b915050604083015184820360408601526119f7828261192b565b915050606083015160608501526080830151611a1660808601826118a9565b5060a0830151611a2960a08601826118a9565b5060c0830151611a3c60c08601826118a9565b5060e0830151611a5060e086018215159052565b5061010080840151611a648287018261196d565b5050610120838101519085015261014080840151908501526101608084015185830382870152611a948382611981565b9250505061018080840151611aab828701826119ad565b5090949350505050565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015611b0c57603f19888603018452611afa8583516119bd565b94509285019290850190600101611ade565b5092979650505050505050565b600080600060608486031215611b2e57600080fd5b8335611b39816118ca565b92506020840135611b49816118ca565b91506040840135611b59816118ca565b809150509250925092565b60008060408385031215611b7757600080fd5b50508035926020909101359150565b60008060408385031215611b9957600080fd5b8235611ba4816118ca565b915060208301356116ec816118ca565b60006001600160401b03821115611bcd57611bcd6116f7565b50601f01601f191660200190565b600082601f830112611bec57600080fd5b8135611bfa6117cc82611bb4565b818152846020838601011115611c0f57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215611c3f57600080fd5b8235611c4a816118ca565b915060208301356001600160401b03811115611c6557600080fd5b611c7185828601611bdb565b9150509250929050565b600060208284031215611c8d57600080fd5b81356001600160401b03811115611ca357600080fd5b820161018081850312156109a457600080fd5b60008060008060808587031215611ccc57600080fd5b8435935060208501359250604085013560028110611ce957600080fd5b915060608501356001600160401b03811115611d0457600080fd5b611d1087828801611bdb565b91505092959194509250565b6020815260006109a460208301846119bd565b6020815260006109a4602083018461192b565b600080600060408486031215611d5757600080fd5b8335925060208401356001600160401b0380821115611d7557600080fd5b818601915086601f830112611d8957600080fd5b813581811115611d9857600080fd5b8760208260051b8501011115611dad57600080fd5b6020830194508093505050509250925092565b60008060408385031215611dd357600080fd5b8235915060208301356001600160401b03811115611df057600080fd5b611c71858286016117ab565b828152604081016109a4602083018461196d565b60008151808452602080850194506020840160005b83811015611e4157815187529582019590820190600101611e25565b509495945050505050565b602081528151602082015260006020830151604080840152611e716060840182611e10565b949350505050565b8051610c27816118ca565b600060208284031215611e9657600080fd5b81516109a4816118ca565b6020808252603a908201527f546f6b656e526567697374727956313a20546f6b656e467265657a654d616e6160408201527967657220616464726573732063616e6e6f74206265207a65726f60301b606082015260800190565b60208082526031908201527f546f6b656e526567697374727956313a20546f6b656e436f726520616464726560408201527073732063616e6e6f74206265207a65726f60781b606082015260800190565b600082601f830112611f5d57600080fd5b8151611f6b6117cc82611bb4565b818152846020838601011115611f8057600080fd5b611e71826020830160208701611907565b80151581146111ff57600080fd5b8051610c2781611f91565b8051610c27816116ba565b600060408284031215611fc757600080fd5b611fcf61170d565b905081516001600160401b03811115611fe757600080fd5b611ff384828501611f4c565b825250602082015160ff8116811461200a57600080fd5b602082015292915050565b600d81106111ff57600080fd5b8051610c2781612015565b60006101a0828403121561204057600080fd5b612048611735565b90508151815260208201516001600160401b038082111561206857600080fd5b61207485838601611f4c565b6020840152604084015191508082111561208d57600080fd5b61209985838601611f4c565b6040840152606084015160608401526120b460808501611e79565b60808401526120c560a08501611e79565b60a08401526120d660c08501611e79565b60c08401526120e760e08501611f9f565b60e084015261010091506120fc828501611faa565b8284015261012091508184015182840152610140915081840151828401526101609150818401518181111561213057600080fd5b61213c86828701611fb5565b83850152505050610180612151818401612022565b9082015292915050565b6000602080838503121561216e57600080fd5b82516001600160401b038082111561218557600080fd5b818501915085601f83011261219957600080fd5b81516121a76117cc82611788565b81815260059190911b830184019084810190888311156121c657600080fd5b8585015b838110156121fe578051858111156121e25760008081fd5b6121f08b89838a010161202d565b8452509186019186016121ca565b5098975050505050505050565b6020808252603a908201527f546f6b656e526567697374727956313a20456e79676d61546f6b656e4d616e6160408201527967657220616464726573732063616e6e6f74206265207a65726f60301b606082015260800190565b60006020828403121561227757600080fd5b81516109a481611f91565b6000808335601e1984360301811261229957600080fd5b83016020810192503590506001600160401b038111156122b857600080fd5b8036038213156122c757600080fd5b9250929050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b8035610c2781611f91565b8035610c2781612015565b60408152600061231d8485612282565b6101808060408601526123356101c0860183856122ce565b92506123446020880188612282565b9250603f198087860301606088015261235e8585846122ce565b945061236d60408a018a612282565b94509150808786030160808801526123868585846122ce565b945061239560608a018a612282565b94509150808786030160a08801526123ae8585846122ce565b9450608089013560c08801526123c660a08a016118df565b93506123d560e08801856118a9565b6123e260c08a018a612282565b945091506101008188870301818901526123fd8686856122ce565b955061240c60e08b018b612282565b955092506101208289880301818a01526124278787866122ce565b9650612434828c016122f7565b95506101409350612448848a018715159052565b612453818c01612302565b9550505050610160612467818801856119ad565b612472828a016122f7565b801515888501529350612486818a016118df565b93505050506124996101a08501826118a9565b5090506109a460208301846118a9565b6000602082840312156124bb57600080fd5b5051919050565b8481528360208201526000600284106124dd576124dd611957565b836040830152608060608301526124f7608083018461192b565b9695505050505050565b60006020828403121561251357600080fd5b81516001600160401b0381111561252957600080fd5b611e718482850161202d565b838152604060208201819052810182905260006001600160fb1b0383111561255c57600080fd5b8260051b8085606085013791909101606001949350505050565b828152604060208201526000611e716040830184611e10565b6000806000606084860312156125a457600080fd5b83516125af81611f91565b602085015190935063ffffffff811681146125c957600080fd5b6040850151909250611b5981611f91565b818103818111156106b257634e487b7160e01b600052601160045260246000fd5b6000825161260d818460208701611907565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca26469706673582212206ac938642fe1984f48e43f444c9b1aa3d14dcdf11fe53cec0ab347066074f16d64736f6c63430008180033",
}

// TokenRegistryV1 is an auto generated Go binding around an Ethereum contract.
type TokenRegistryV1 struct {
	abi abi.ABI
}

// NewTokenRegistryV1 creates a new instance of TokenRegistryV1.
func NewTokenRegistryV1() *TokenRegistryV1 {
	parsed, err := TokenRegistryV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &TokenRegistryV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *TokenRegistryV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (tokenRegistryV1 *TokenRegistryV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := tokenRegistryV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (tokenRegistryV1 *TokenRegistryV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := tokenRegistryV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackAddToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x72d14e30.
//
// Solidity: function addToken((string,string,string,bytes,uint256,address,bytes,bytes,bool,uint8,bool,address) tokenData) returns(bytes32)
func (tokenRegistryV1 *TokenRegistryV1) PackAddToken(tokenData SharedObjectsTokenRegistrationData) []byte {
	enc, err := tokenRegistryV1.abi.Pack("addToken", tokenData)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAddToken is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x72d14e30.
//
// Solidity: function addToken((string,string,string,bytes,uint256,address,bytes,bytes,bool,uint8,bool,address) tokenData) returns(bytes32)
func (tokenRegistryV1 *TokenRegistryV1) UnpackAddToken(data []byte) ([32]byte, error) {
	out, err := tokenRegistryV1.abi.Unpack("addToken", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (tokenRegistryV1 *TokenRegistryV1) PackAuthority() []byte {
	enc, err := tokenRegistryV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (tokenRegistryV1 *TokenRegistryV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := tokenRegistryV1.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackBroadcastCurrentFrozenResourcesForNewParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x49f6ee3c.
//
// Solidity: function broadcastCurrentFrozenResourcesForNewParticipant() returns()
func (tokenRegistryV1 *TokenRegistryV1) PackBroadcastCurrentFrozenResourcesForNewParticipant() []byte {
	enc, err := tokenRegistryV1.abi.Pack("broadcastCurrentFrozenResourcesForNewParticipant")
	if err != nil {
		panic(err)
	}
	return enc
}

// PackBroadcastFrozenToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0566585a.
//
// Solidity: function broadcastFrozenToken((bytes32,uint256[]) frozenToken) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackBroadcastFrozenToken(frozenToken TokenStructsFrozenToken) []byte {
	enc, err := tokenRegistryV1.abi.Pack("broadcastFrozenToken", frozenToken)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackBroadcastUnfrozenToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe66cdd2b.
//
// Solidity: function broadcastUnfrozenToken((bytes32,uint256[]) unfrozenToken) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackBroadcastUnfrozenToken(unfrozenToken TokenStructsFrozenToken) []byte {
	enc, err := tokenRegistryV1.abi.Pack("broadcastUnfrozenToken", unfrozenToken)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackConfigureModules is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x31b1125f.
//
// Solidity: function configureModules(address _tokenCore, address _tokenFreezeManager, address _enygmaTokenManager) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackConfigureModules(tokenCore common.Address, tokenFreezeManager common.Address, enygmaTokenManager common.Address) []byte {
	enc, err := tokenRegistryV1.abi.Pack("configureModules", tokenCore, tokenFreezeManager, enygmaTokenManager)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackContractVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (tokenRegistryV1 *TokenRegistryV1) PackContractVersion() []byte {
	enc, err := tokenRegistryV1.abi.Pack("contractVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackContractVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (tokenRegistryV1 *TokenRegistryV1) UnpackContractVersion(data []byte) (*big.Int, error) {
	out, err := tokenRegistryV1.abi.Unpack("contractVersion", data)
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
func (tokenRegistryV1 *TokenRegistryV1) PackEndpoint() []byte {
	enc, err := tokenRegistryV1.abi.Pack("endpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (tokenRegistryV1 *TokenRegistryV1) UnpackEndpoint(data []byte) (common.Address, error) {
	out, err := tokenRegistryV1.abi.Unpack("endpoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackFreezeToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe08ac15e.
//
// Solidity: function freezeToken(bytes32 resourceId, uint256[] chainIds) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackFreezeToken(resourceId [32]byte, chainIds []*big.Int) []byte {
	enc, err := tokenRegistryV1.abi.Pack("freezeToken", resourceId, chainIds)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetAddressByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (tokenRegistryV1 *TokenRegistryV1) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := tokenRegistryV1.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (tokenRegistryV1 *TokenRegistryV1) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := tokenRegistryV1.abi.Unpack("getAddressByResourceId", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetAllTokens is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2a5c792a.
//
// Solidity: function getAllTokens() view returns((bytes32,string,string,uint256,address,address,address,bool,uint8,uint256,uint256,(string,uint8),uint8)[])
func (tokenRegistryV1 *TokenRegistryV1) PackGetAllTokens() []byte {
	enc, err := tokenRegistryV1.abi.Pack("getAllTokens")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllTokens is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2a5c792a.
//
// Solidity: function getAllTokens() view returns((bytes32,string,string,uint256,address,address,address,bool,uint8,uint256,uint256,(string,uint8),uint8)[])
func (tokenRegistryV1 *TokenRegistryV1) UnpackGetAllTokens(data []byte) ([]TokenStructsToken, error) {
	out, err := tokenRegistryV1.abi.Unpack("getAllTokens", data)
	if err != nil {
		return *new([]TokenStructsToken), err
	}
	out0 := *abi.ConvertType(out[0], new([]TokenStructsToken)).(*[]TokenStructsToken)
	return out0, err
}

// PackGetEnygmaTokenManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa0b3f1ee.
//
// Solidity: function getEnygmaTokenManager() view returns(address)
func (tokenRegistryV1 *TokenRegistryV1) PackGetEnygmaTokenManager() []byte {
	enc, err := tokenRegistryV1.abi.Pack("getEnygmaTokenManager")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEnygmaTokenManager is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0b3f1ee.
//
// Solidity: function getEnygmaTokenManager() view returns(address)
func (tokenRegistryV1 *TokenRegistryV1) UnpackGetEnygmaTokenManager(data []byte) (common.Address, error) {
	out, err := tokenRegistryV1.abi.Unpack("getEnygmaTokenManager", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetTokenByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a58b79f.
//
// Solidity: function getTokenByResourceId(bytes32 resourceId) view returns((bytes32,string,string,uint256,address,address,address,bool,uint8,uint256,uint256,(string,uint8),uint8))
func (tokenRegistryV1 *TokenRegistryV1) PackGetTokenByResourceId(resourceId [32]byte) []byte {
	enc, err := tokenRegistryV1.abi.Pack("getTokenByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8a58b79f.
//
// Solidity: function getTokenByResourceId(bytes32 resourceId) view returns((bytes32,string,string,uint256,address,address,address,bool,uint8,uint256,uint256,(string,uint8),uint8))
func (tokenRegistryV1 *TokenRegistryV1) UnpackGetTokenByResourceId(data []byte) (TokenStructsToken, error) {
	out, err := tokenRegistryV1.abi.Unpack("getTokenByResourceId", data)
	if err != nil {
		return *new(TokenStructsToken), err
	}
	out0 := *abi.ConvertType(out[0], new(TokenStructsToken)).(*TokenStructsToken)
	return out0, err
}

// PackGetTokenCore is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf0a599bb.
//
// Solidity: function getTokenCore() view returns(address)
func (tokenRegistryV1 *TokenRegistryV1) PackGetTokenCore() []byte {
	enc, err := tokenRegistryV1.abi.Pack("getTokenCore")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenCore is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf0a599bb.
//
// Solidity: function getTokenCore() view returns(address)
func (tokenRegistryV1 *TokenRegistryV1) UnpackGetTokenCore(data []byte) (common.Address, error) {
	out, err := tokenRegistryV1.abi.Unpack("getTokenCore", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetTokenFreezeManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf356cc03.
//
// Solidity: function getTokenFreezeManager() view returns(address)
func (tokenRegistryV1 *TokenRegistryV1) PackGetTokenFreezeManager() []byte {
	enc, err := tokenRegistryV1.abi.Pack("getTokenFreezeManager")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenFreezeManager is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf356cc03.
//
// Solidity: function getTokenFreezeManager() view returns(address)
func (tokenRegistryV1 *TokenRegistryV1) UnpackGetTokenFreezeManager(data []byte) (common.Address, error) {
	out, err := tokenRegistryV1.abi.Unpack("getTokenFreezeManager", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x485cc955.
//
// Solidity: function initialize(address _endpoint, address authority_) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackInitialize(endpoint common.Address, authority common.Address) []byte {
	enc, err := tokenRegistryV1.abi.Pack("initialize", endpoint, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackInitialize0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4d66de8.
//
// Solidity: function initialize(address _endpoint) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackInitialize0(endpoint common.Address) []byte {
	enc, err := tokenRegistryV1.abi.Pack("initialize0", endpoint)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackIsTokenFrozenForParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x338042e4.
//
// Solidity: function isTokenFrozenForParticipant(bytes32 resourceId, uint256 chainId) view returns(bool)
func (tokenRegistryV1 *TokenRegistryV1) PackIsTokenFrozenForParticipant(resourceId [32]byte, chainId *big.Int) []byte {
	enc, err := tokenRegistryV1.abi.Pack("isTokenFrozenForParticipant", resourceId, chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsTokenFrozenForParticipant is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x338042e4.
//
// Solidity: function isTokenFrozenForParticipant(bytes32 resourceId, uint256 chainId) view returns(bool)
func (tokenRegistryV1 *TokenRegistryV1) UnpackIsTokenFrozenForParticipant(data []byte) (bool, error) {
	out, err := tokenRegistryV1.abi.Unpack("isTokenFrozenForParticipant", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (tokenRegistryV1 *TokenRegistryV1) PackProxiableUUID() []byte {
	enc, err := tokenRegistryV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (tokenRegistryV1 *TokenRegistryV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := tokenRegistryV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (tokenRegistryV1 *TokenRegistryV1) PackResourceId() []byte {
	enc, err := tokenRegistryV1.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (tokenRegistryV1 *TokenRegistryV1) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := tokenRegistryV1.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSetEnygmaFactory is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x153f716d.
//
// Solidity: function setEnygmaFactory(address _enygmaFactory) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackSetEnygmaFactory(enygmaFactory common.Address) []byte {
	enc, err := tokenRegistryV1.abi.Pack("setEnygmaFactory", enygmaFactory)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetEnygmaTokenManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf867cd52.
//
// Solidity: function setEnygmaTokenManager(address _enygmaTokenManager) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackSetEnygmaTokenManager(enygmaTokenManager common.Address) []byte {
	enc, err := tokenRegistryV1.abi.Pack("setEnygmaTokenManager", enygmaTokenManager)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa01afbfb.
//
// Solidity: function setResourceId(bytes32 _resourceId) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackSetResourceId(resourceId [32]byte) []byte {
	enc, err := tokenRegistryV1.abi.Pack("setResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetTokenCore is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14ae646b.
//
// Solidity: function setTokenCore(address _tokenCore) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackSetTokenCore(tokenCore common.Address) []byte {
	enc, err := tokenRegistryV1.abi.Pack("setTokenCore", tokenCore)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetTokenFreezeManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x13aa1a1d.
//
// Solidity: function setTokenFreezeManager(address _tokenFreezeManager) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackSetTokenFreezeManager(tokenFreezeManager common.Address) []byte {
	enc, err := tokenRegistryV1.abi.Pack("setTokenFreezeManager", tokenFreezeManager)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUnfreezeToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf7bb5fad.
//
// Solidity: function unfreezeToken(bytes32 resourceId, uint256[] chainIds) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackUnfreezeToken(resourceId [32]byte, chainIds []*big.Int) []byte {
	enc, err := tokenRegistryV1.abi.Pack("unfreezeToken", resourceId, chainIds)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdateStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x054372ed.
//
// Solidity: function updateStatus(bytes32 resourceId, uint8 status) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackUpdateStatus(resourceId [32]byte, status uint8) []byte {
	enc, err := tokenRegistryV1.abi.Pack("updateStatus", resourceId, status)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdateTokenBalance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x80387be2.
//
// Solidity: function updateTokenBalance(uint256 issuerChainId, bytes32 resourceId, uint8 updateType, bytes metadata) returns()
func (tokenRegistryV1 *TokenRegistryV1) PackUpdateTokenBalance(issuerChainId *big.Int, resourceId [32]byte, updateType uint8, metadata []byte) []byte {
	enc, err := tokenRegistryV1.abi.Pack("updateTokenBalance", issuerChainId, resourceId, updateType, metadata)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (tokenRegistryV1 *TokenRegistryV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := tokenRegistryV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TokenRegistryV1AuthorityUpdated represents a AuthorityUpdated event raised by the TokenRegistryV1 contract.
type TokenRegistryV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const TokenRegistryV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (TokenRegistryV1AuthorityUpdated) ContractEventName() string {
	return TokenRegistryV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (tokenRegistryV1 *TokenRegistryV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*TokenRegistryV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != tokenRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TokenRegistryV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := tokenRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tokenRegistryV1.abi.Events[event].Inputs {
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

// TokenRegistryV1Initialized represents a Initialized event raised by the TokenRegistryV1 contract.
type TokenRegistryV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const TokenRegistryV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (TokenRegistryV1Initialized) ContractEventName() string {
	return TokenRegistryV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (tokenRegistryV1 *TokenRegistryV1) UnpackInitializedEvent(log *types.Log) (*TokenRegistryV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != tokenRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TokenRegistryV1Initialized)
	if len(log.Data) > 0 {
		if err := tokenRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tokenRegistryV1.abi.Events[event].Inputs {
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

// TokenRegistryV1ModulesConfigured represents a ModulesConfigured event raised by the TokenRegistryV1 contract.
type TokenRegistryV1ModulesConfigured struct {
	TokenCore          common.Address
	TokenFreezeManager common.Address
	EnygmaTokenManager common.Address
	Raw                *types.Log // Blockchain specific contextual infos
}

const TokenRegistryV1ModulesConfiguredEventName = "ModulesConfigured"

// ContractEventName returns the user-defined event name.
func (TokenRegistryV1ModulesConfigured) ContractEventName() string {
	return TokenRegistryV1ModulesConfiguredEventName
}

// UnpackModulesConfiguredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ModulesConfigured(address indexed tokenCore, address indexed tokenFreezeManager, address indexed enygmaTokenManager)
func (tokenRegistryV1 *TokenRegistryV1) UnpackModulesConfiguredEvent(log *types.Log) (*TokenRegistryV1ModulesConfigured, error) {
	event := "ModulesConfigured"
	if log.Topics[0] != tokenRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TokenRegistryV1ModulesConfigured)
	if len(log.Data) > 0 {
		if err := tokenRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tokenRegistryV1.abi.Events[event].Inputs {
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

// TokenRegistryV1Upgraded represents a Upgraded event raised by the TokenRegistryV1 contract.
type TokenRegistryV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const TokenRegistryV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (TokenRegistryV1Upgraded) ContractEventName() string {
	return TokenRegistryV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (tokenRegistryV1 *TokenRegistryV1) UnpackUpgradedEvent(log *types.Log) (*TokenRegistryV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != tokenRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TokenRegistryV1Upgraded)
	if len(log.Data) > 0 {
		if err := tokenRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tokenRegistryV1.abi.Events[event].Inputs {
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
func (tokenRegistryV1 *TokenRegistryV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["RaylsAppV1HubNotActive"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackRaylsAppV1HubNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["RaylsAppV1PrivacyNodeFrozen"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackRaylsAppV1PrivacyNodeFrozenError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["RaylsAppV1PrivacyNodeNotActive"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackRaylsAppV1PrivacyNodeNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["RaylsAppV1PublicChainNotActive"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackRaylsAppV1PublicChainNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["RaylsAppV1ResourceNotApproved"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackRaylsAppV1ResourceNotApprovedError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["RaylsAppV1TokenRegistryNotConfigured"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackRaylsAppV1TokenRegistryNotConfiguredError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["RaylsAppV1UnauthorizedTokenRegistry"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackRaylsAppV1UnauthorizedTokenRegistryError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenRegistryV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return tokenRegistryV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// TokenRegistryV1AddressEmptyCode represents a AddressEmptyCode error raised by the TokenRegistryV1 contract.
type TokenRegistryV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func TokenRegistryV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (tokenRegistryV1 *TokenRegistryV1) UnpackAddressEmptyCodeError(raw []byte) (*TokenRegistryV1AddressEmptyCode, error) {
	out := new(TokenRegistryV1AddressEmptyCode)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the TokenRegistryV1 contract.
type TokenRegistryV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func TokenRegistryV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (tokenRegistryV1 *TokenRegistryV1) UnpackERC1967InvalidImplementationError(raw []byte) (*TokenRegistryV1ERC1967InvalidImplementation, error) {
	out := new(TokenRegistryV1ERC1967InvalidImplementation)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the TokenRegistryV1 contract.
type TokenRegistryV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func TokenRegistryV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (tokenRegistryV1 *TokenRegistryV1) UnpackERC1967NonPayableError(raw []byte) (*TokenRegistryV1ERC1967NonPayable, error) {
	out := new(TokenRegistryV1ERC1967NonPayable)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1FailedCall represents a FailedCall error raised by the TokenRegistryV1 contract.
type TokenRegistryV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func TokenRegistryV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (tokenRegistryV1 *TokenRegistryV1) UnpackFailedCallError(raw []byte) (*TokenRegistryV1FailedCall, error) {
	out := new(TokenRegistryV1FailedCall)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1InvalidInitialization represents a InvalidInitialization error raised by the TokenRegistryV1 contract.
type TokenRegistryV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func TokenRegistryV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (tokenRegistryV1 *TokenRegistryV1) UnpackInvalidInitializationError(raw []byte) (*TokenRegistryV1InvalidInitialization, error) {
	out := new(TokenRegistryV1InvalidInitialization)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1NotInitializing represents a NotInitializing error raised by the TokenRegistryV1 contract.
type TokenRegistryV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func TokenRegistryV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (tokenRegistryV1 *TokenRegistryV1) UnpackNotInitializingError(raw []byte) (*TokenRegistryV1NotInitializing, error) {
	out := new(TokenRegistryV1NotInitializing)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the TokenRegistryV1 contract.
type TokenRegistryV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func TokenRegistryV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (tokenRegistryV1 *TokenRegistryV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*TokenRegistryV1RaylsAccessManagedContractPaused, error) {
	out := new(TokenRegistryV1RaylsAccessManagedContractPaused)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the TokenRegistryV1 contract.
type TokenRegistryV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func TokenRegistryV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (tokenRegistryV1 *TokenRegistryV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*TokenRegistryV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(TokenRegistryV1RaylsAccessManagedInvalidAuthority)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the TokenRegistryV1 contract.
type TokenRegistryV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func TokenRegistryV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (tokenRegistryV1 *TokenRegistryV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*TokenRegistryV1RaylsAccessManagedMustSchedule, error) {
	out := new(TokenRegistryV1RaylsAccessManagedMustSchedule)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the TokenRegistryV1 contract.
type TokenRegistryV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func TokenRegistryV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (tokenRegistryV1 *TokenRegistryV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*TokenRegistryV1RaylsAccessManagedUnauthorized, error) {
	out := new(TokenRegistryV1RaylsAccessManagedUnauthorized)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1RaylsAppV1HubNotActive represents a RaylsAppV1__HubNotActive error raised by the TokenRegistryV1 contract.
type TokenRegistryV1RaylsAppV1HubNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	HubStatus         uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func TokenRegistryV1RaylsAppV1HubNotActiveErrorID() common.Hash {
	return common.HexToHash("0x3fae5bbd70277aa1cd008dceb93b19a7055c2a6d29b84733371e1c41b2048b15")
}

// UnpackRaylsAppV1HubNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func (tokenRegistryV1 *TokenRegistryV1) UnpackRaylsAppV1HubNotActiveError(raw []byte) (*TokenRegistryV1RaylsAppV1HubNotActive, error) {
	out := new(TokenRegistryV1RaylsAppV1HubNotActive)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAppV1HubNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1RaylsAppV1PrivacyNodeFrozen represents a RaylsAppV1__PrivacyNodeFrozen error raised by the TokenRegistryV1 contract.
type TokenRegistryV1RaylsAppV1PrivacyNodeFrozen struct {
	TokenAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__PrivacyNodeFrozen(address tokenAddress)
func TokenRegistryV1RaylsAppV1PrivacyNodeFrozenErrorID() common.Hash {
	return common.HexToHash("0xc80bd255e67000277f5aed4960b64f92e2d5a652f07a22fba7d044de6add8f0e")
}

// UnpackRaylsAppV1PrivacyNodeFrozenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__PrivacyNodeFrozen(address tokenAddress)
func (tokenRegistryV1 *TokenRegistryV1) UnpackRaylsAppV1PrivacyNodeFrozenError(raw []byte) (*TokenRegistryV1RaylsAppV1PrivacyNodeFrozen, error) {
	out := new(TokenRegistryV1RaylsAppV1PrivacyNodeFrozen)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAppV1PrivacyNodeFrozen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1RaylsAppV1PrivacyNodeNotActive represents a RaylsAppV1__PrivacyNodeNotActive error raised by the TokenRegistryV1 contract.
type TokenRegistryV1RaylsAppV1PrivacyNodeNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func TokenRegistryV1RaylsAppV1PrivacyNodeNotActiveErrorID() common.Hash {
	return common.HexToHash("0xfdcdd2a6e576bf1f342ce493560565ef686a59cd3e0486f6869151efb2c7853f")
}

// UnpackRaylsAppV1PrivacyNodeNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func (tokenRegistryV1 *TokenRegistryV1) UnpackRaylsAppV1PrivacyNodeNotActiveError(raw []byte) (*TokenRegistryV1RaylsAppV1PrivacyNodeNotActive, error) {
	out := new(TokenRegistryV1RaylsAppV1PrivacyNodeNotActive)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAppV1PrivacyNodeNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1RaylsAppV1PublicChainNotActive represents a RaylsAppV1__PublicChainNotActive error raised by the TokenRegistryV1 contract.
type TokenRegistryV1RaylsAppV1PublicChainNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	PublicChainStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func TokenRegistryV1RaylsAppV1PublicChainNotActiveErrorID() common.Hash {
	return common.HexToHash("0xb607961611e6e4126e09c80bcd1e35e7a1e886888daa292eecc27cd9d4e37f3f")
}

// UnpackRaylsAppV1PublicChainNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func (tokenRegistryV1 *TokenRegistryV1) UnpackRaylsAppV1PublicChainNotActiveError(raw []byte) (*TokenRegistryV1RaylsAppV1PublicChainNotActive, error) {
	out := new(TokenRegistryV1RaylsAppV1PublicChainNotActive)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAppV1PublicChainNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1RaylsAppV1ResourceNotApproved represents a RaylsAppV1__ResourceNotApproved error raised by the TokenRegistryV1 contract.
type TokenRegistryV1RaylsAppV1ResourceNotApproved struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__ResourceNotApproved()
func TokenRegistryV1RaylsAppV1ResourceNotApprovedErrorID() common.Hash {
	return common.HexToHash("0x8f144935367c131b72d26b0320b764f69ba3639e65abb1c811084bbd46e5c731")
}

// UnpackRaylsAppV1ResourceNotApprovedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__ResourceNotApproved()
func (tokenRegistryV1 *TokenRegistryV1) UnpackRaylsAppV1ResourceNotApprovedError(raw []byte) (*TokenRegistryV1RaylsAppV1ResourceNotApproved, error) {
	out := new(TokenRegistryV1RaylsAppV1ResourceNotApproved)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAppV1ResourceNotApproved", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1RaylsAppV1TokenRegistryNotConfigured represents a RaylsAppV1__TokenRegistryNotConfigured error raised by the TokenRegistryV1 contract.
type TokenRegistryV1RaylsAppV1TokenRegistryNotConfigured struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__TokenRegistryNotConfigured()
func TokenRegistryV1RaylsAppV1TokenRegistryNotConfiguredErrorID() common.Hash {
	return common.HexToHash("0x3eba255b70fc7afd9cc5be90de2023dae8350ac3c29cbd5eaf139cadd9c4292e")
}

// UnpackRaylsAppV1TokenRegistryNotConfiguredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__TokenRegistryNotConfigured()
func (tokenRegistryV1 *TokenRegistryV1) UnpackRaylsAppV1TokenRegistryNotConfiguredError(raw []byte) (*TokenRegistryV1RaylsAppV1TokenRegistryNotConfigured, error) {
	out := new(TokenRegistryV1RaylsAppV1TokenRegistryNotConfigured)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAppV1TokenRegistryNotConfigured", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1RaylsAppV1UnauthorizedTokenRegistry represents a RaylsAppV1__UnauthorizedTokenRegistry error raised by the TokenRegistryV1 contract.
type TokenRegistryV1RaylsAppV1UnauthorizedTokenRegistry struct {
	Caller        common.Address
	TokenRegistry common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAppV1__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func TokenRegistryV1RaylsAppV1UnauthorizedTokenRegistryErrorID() common.Hash {
	return common.HexToHash("0x000d23e5a298a9951b289bd8f5eece62aa717c000d6b0509a9f77d16f67a5b7d")
}

// UnpackRaylsAppV1UnauthorizedTokenRegistryError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAppV1__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func (tokenRegistryV1 *TokenRegistryV1) UnpackRaylsAppV1UnauthorizedTokenRegistryError(raw []byte) (*TokenRegistryV1RaylsAppV1UnauthorizedTokenRegistry, error) {
	out := new(TokenRegistryV1RaylsAppV1UnauthorizedTokenRegistry)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "RaylsAppV1UnauthorizedTokenRegistry", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the TokenRegistryV1 contract.
type TokenRegistryV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func TokenRegistryV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (tokenRegistryV1 *TokenRegistryV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*TokenRegistryV1UUPSUnauthorizedCallContext, error) {
	out := new(TokenRegistryV1UUPSUnauthorizedCallContext)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRegistryV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the TokenRegistryV1 contract.
type TokenRegistryV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func TokenRegistryV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (tokenRegistryV1 *TokenRegistryV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*TokenRegistryV1UUPSUnsupportedProxiableUUID, error) {
	out := new(TokenRegistryV1UUPSUnsupportedProxiableUUID)
	if err := tokenRegistryV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
