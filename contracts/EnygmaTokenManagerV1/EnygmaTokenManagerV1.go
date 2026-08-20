// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package EnygmaTokenManagerV1

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

// EnygmaTokenManagerV1MetaData contains all meta data concerning the EnygmaTokenManagerV1 contract.
var EnygmaTokenManagerV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"endpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRaylsEndpoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"enygmaFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEnygmaFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_tokenRegistryAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_enygmaFactory\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerEnygmaToken\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"tokenData\",\"type\":\"tuple\",\"internalType\":\"structSharedObjects.TokenRegistrationData\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"totalSupply\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"issuerChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pnRegistryAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bytecode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"initializerParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"isFungible\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"ercStandard\",\"type\":\"uint8\",\"internalType\":\"enumSharedObjects.ErcStandard\"},{\"name\":\"isCustom\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"participantStorage\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setEndpoint\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEnygmaFactory\",\"inputs\":[{\"name\":\"_enygmaFactory\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTokenCoreAddress\",\"inputs\":[{\"name\":\"_tokenCoreAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTokenRegistryAddress\",\"inputs\":[{\"name\":\"_tokenRegistryAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tokenCoreAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenRegistryAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaFactoryUpdated\",\"inputs\":[{\"name\":\"oldFactory\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newFactory\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaTokenRegistered\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"issuerChainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"initialSupply\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnygmaTokenManagerV1__UnauthorizedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EnygmaTokenManagerV1__ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "EnygmaTokenManagerV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b5060805161156061003e6000396000818161093b015281816109640152610a9d01526115606000f3fe6080604052600436106100d35760003560e01c80636de9a6ee1161007a5780636de9a6ee146101d8578063821db16e146101f85780638f388fc614610216578063a9ac0efd14610236578063ad3cb1cc14610256578063bf7e214f14610294578063dbbb4155146102a9578063f8c8765e146102c957600080fd5b8063153f716d146100d85780634f1ef286146100fa57806352d1902d1461010d57806358555ab8146101355780635be2aca0146101555780635e280f11146101825780635f997c5b146101a257806362abdc39146101b8575b600080fd5b3480156100e457600080fd5b506100f86100f3366004610fbb565b6102e9565b005b6100f8610108366004611063565b610398565b34801561011957600080fd5b506101226103b7565b6040519081526020015b60405180910390f35b34801561014157600080fd5b506100f8610150366004610fbb565b6103d4565b34801561016157600080fd5b50600054610175906001600160a01b031681565b60405161012c91906110d3565b34801561018e57600080fd5b50600254610175906001600160a01b031681565b3480156101ae57600080fd5b5061012260045481565b3480156101c457600080fd5b506100f86101d33660046110e7565b610433565b3480156101e457600080fd5b50600154610175906001600160a01b031681565b34801561020457600080fd5b506003546001600160a01b0316610175565b34801561022257600080fd5b506100f8610231366004610fbb565b610717565b34801561024257600080fd5b50600354610175906001600160a01b031681565b34801561026257600080fd5b50610287604051806040016040528060058152602001640352e302e360dc1b81525081565b60405161012c91906111b4565b3480156102a057600080fd5b50610175610776565b3480156102b557600080fd5b506100f86102c4366004610fbb565b61078f565b3480156102d557600080fd5b506100f86102e43660046111c7565b6107ee565b6000546001600160a01b0316331461031f57336040516309dd326b60e41b815260040161031691906110d3565b60405180910390fd5b6001600160a01b03811661034657604051630bc624a360e41b815260040160405180910390fd5b600380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907ffda14faeaf47950b60da38870ec78c6a2a9c6e5c7fc977ba9bcb9f1d2303b20590600090a35050565b6103a0610930565b6103a9826109c0565b6103b382826109d9565b5050565b60006103c1610a92565b5060008051602061150b83398151915290565b6103ea336000356001600160e01b031916610adb565b6001600160a01b03811661041157604051630bc624a360e41b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b0392909216919091179055565b6001546001600160a01b0316331461046057336040516309dd326b60e41b815260040161031691906110d3565b6003546001600160a01b0316600061047b60e0870187611223565b8101906104889190611290565b9250505060006040518061014001604052808880600001906104aa9190611223565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505050908252506020908101906104f3908a018a611223565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525093855250505060ff8516602083015260408083018c90526001600160a01b03808a16606085015260808c8101359085015288811660a0850152600254811660c08501529154821660e084015260086101009093019290925290516372dce7df60e11b815291925084169063e5b9cfbe906105a190849060040161130d565b600060405180830381600087803b1580156105bb57600080fd5b505af11580156105cf573d6000803e3d6000fd5b50506040516375ede10d60e01b8152600481018b9052600092506001600160a01b03861691506375ede10d90602401602060405180830381865afa15801561061b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061063f91906113d4565b60025460405163d4f951c760e01b8152600481018c90526001600160a01b03808416602483015292935091169063d4f951c790604401600060405180830381600087803b15801561068f57600080fd5b505af11580156106a3573d6000803e3d6000fd5b5050505060808801357f74d31871242a42afa34ddbdcd2089ee9700aa8499e2662f830955d4b457aaf308a436106d98c80611223565b6106e660608f018f611223565b8101906106f391906113f1565b60405161070495949392919061140a565b60405180910390a2505050505050505050565b61072d336000356001600160e01b031916610adb565b6001600160a01b03811661075457604051630bc624a360e41b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b0392909216919091179055565b6000610780610c26565b546001600160a01b0316919050565b6107a5336000356001600160e01b031916610adb565b6001600160a01b0381166107cc57604051630bc624a360e41b815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0392909216919091179055565b60006107f8610c88565b805490915060ff600160401b82041615906001600160401b031660008115801561081f5750825b90506000826001600160401b0316600114801561083b5750303b155b905081158015610849575080155b156108675760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561089157845460ff60401b1916600160401b1785555b610899610cb3565b600280546001600160a01b03808c166001600160a01b031992831617909255600080548b841690831617905560038054928a16929091169190911790556108df86610cbb565b831561092557845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614806109a057507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610994610cfc565b6001600160a01b031614155b156109be5760405163703e46dd60e11b815260040160405180910390fd5b565b6109d6336000356001600160e01b031916610adb565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610a33575060408051601f3d908101601f19168201909252610a3091810190611451565b60015b610a525781604051634c9c8ce360e01b815260040161031691906110d3565b60008051602061150b8339815191528114610a8357604051632a87526960e21b815260048101829052602401610316565b610a8d8383610d12565b505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146109be5760405163703e46dd60e11b815260040160405180910390fd5b6000610ae5610c26565b80549091506001600160a01b031680610b14576000604051638944034760e01b815260040161031691906110d3565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015610b78573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b9c919061147f565b92509250925082610c1d578015610bc65760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff821615610c025760405163a426878960e01b81526001600160a01b038816600482015263ffffffff83166024820152604401610316565b86604051632ecd3d0360e21b815260040161031691906110d3565b50505050505050565b60008060ff19610c5760017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f356114cd565b604051602001610c6991815260200190565b60408051601f1981840301815291905280516020909101201692915050565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b6109be610d68565b6000610cc5610c26565b80549091506001600160a01b031615610cf35781604051638944034760e01b815260040161031691906110d3565b6103b382610d8d565b600060008051602061150b833981519152610780565b610d1b82610e1d565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a2805115610d6057610a8d8282610e79565b6103b3610eef565b610d70610f0e565b6109be57604051631afcd79f60e31b815260040160405180910390fd5b6001600160a01b038116610db65780604051638944034760e01b815260040161031691906110d3565b6000610dc0610c26565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b806001600160a01b03163b600003610e4a5780604051634c9c8ce360e01b815260040161031691906110d3565b60008051602061150b83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b031684604051610e9691906114ee565b600060405180830381855af49150503d8060008114610ed1576040519150601f19603f3d011682016040523d82523d6000602084013e610ed6565b606091505b5091509150610ee6858383610f28565b95945050505050565b34156109be5760405163b398979f60e01b815260040160405180910390fd5b6000610f18610c88565b54600160401b900460ff16919050565b606082610f3d57610f3882610f7e565b610f77565b8151158015610f5457506001600160a01b0384163b155b15610f745783604051639996b31560e01b815260040161031691906110d3565b50805b9392505050565b805115610f8d57805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b6001600160a01b03811681146109d657600080fd5b600060208284031215610fcd57600080fd5b8135610f7781610fa6565b634e487b7160e01b600052604160045260246000fd5b60006001600160401b038084111561100857611008610fd8565b604051601f8501601f19908116603f0116810190828211818310171561103057611030610fd8565b8160405280935085815286868601111561104957600080fd5b858560208301376000602087830101525050509392505050565b6000806040838503121561107657600080fd5b823561108181610fa6565b915060208301356001600160401b0381111561109c57600080fd5b8301601f810185136110ad57600080fd5b6110bc85823560208401610fee565b9150509250929050565b6001600160a01b03169052565b6001600160a01b0391909116815260200190565b600080600080600060a086880312156110ff57600080fd5b8535945060208601356001600160401b0381111561111c57600080fd5b8601610180818903121561112f57600080fd5b935060408601359250606086013561114681610fa6565b9150608086013561115681610fa6565b809150509295509295909350565b60005b8381101561117f578181015183820152602001611167565b50506000910152565b600081518084526111a0816020860160208601611164565b601f01601f19169290920160200192915050565b602081526000610f776020830184611188565b600080600080608085870312156111dd57600080fd5b84356111e881610fa6565b935060208501356111f881610fa6565b9250604085013561120881610fa6565b9150606085013561121881610fa6565b939692955090935050565b6000808335601e1984360301811261123a57600080fd5b8301803591506001600160401b0382111561125457600080fd5b60200191503681900382131561126957600080fd5b9250929050565b600082601f83011261128157600080fd5b610f7783833560208501610fee565b6000806000606084860312156112a557600080fd5b83356001600160401b03808211156112bc57600080fd5b6112c887838801611270565b945060208601359150808211156112de57600080fd5b506112eb86828701611270565b925050604084013560ff8116811461130257600080fd5b809150509250925092565b602081526000825161014080602085015261132c610160850183611188565b91506020850151601f198584030160408601526113498382611188565b9250506040850151611360606086018260ff169052565b5060608501516080850152608085015161137d60a08601826110c6565b5060a085015160c085015260c085015161139a60e08601826110c6565b5060e08501516101006113af818701836110c6565b86015190506101206113c3868201836110c6565b959095015193019290925250919050565b6000602082840312156113e657600080fd5b8151610f7781610fa6565b60006020828403121561140357600080fd5b5035919050565b85815284602082015260806040820152826080820152828460a0830137600060a08483010152600060a0601f19601f86011683010190508260608301529695505050505050565b60006020828403121561146357600080fd5b5051919050565b8051801515811461147a57600080fd5b919050565b60008060006060848603121561149457600080fd5b61149d8461146a565b9250602084015163ffffffff811681146114b657600080fd5b91506114c46040850161146a565b90509250925092565b81810381811115610cad57634e487b7160e01b600052601160045260246000fd5b60008251611500818460208701611164565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca2646970667358221220b65815f937930a7fd18e38bda209f7fb70522e708dbdb6dca7e565a0584ee80164736f6c63430008180033",
}

// EnygmaTokenManagerV1 is an auto generated Go binding around an Ethereum contract.
type EnygmaTokenManagerV1 struct {
	abi abi.ABI
}

// NewEnygmaTokenManagerV1 creates a new instance of EnygmaTokenManagerV1.
func NewEnygmaTokenManagerV1() *EnygmaTokenManagerV1 {
	parsed, err := EnygmaTokenManagerV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &EnygmaTokenManagerV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *EnygmaTokenManagerV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := enygmaTokenManagerV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackAuthority() []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := enygmaTokenManagerV1.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackEndpoint() []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("endpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackEndpoint(data []byte) (common.Address, error) {
	out, err := enygmaTokenManagerV1.abi.Unpack("endpoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackEnygmaFactory is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9ac0efd.
//
// Solidity: function enygmaFactory() view returns(address)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackEnygmaFactory() []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("enygmaFactory")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEnygmaFactory is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9ac0efd.
//
// Solidity: function enygmaFactory() view returns(address)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackEnygmaFactory(data []byte) (common.Address, error) {
	out, err := enygmaTokenManagerV1.abi.Unpack("enygmaFactory", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetEnygmaFactory is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x821db16e.
//
// Solidity: function getEnygmaFactory() view returns(address)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackGetEnygmaFactory() []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("getEnygmaFactory")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEnygmaFactory is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x821db16e.
//
// Solidity: function getEnygmaFactory() view returns(address)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackGetEnygmaFactory(data []byte) (common.Address, error) {
	out, err := enygmaTokenManagerV1.abi.Unpack("getEnygmaFactory", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf8c8765e.
//
// Solidity: function initialize(address _endpoint, address _tokenRegistryAddress, address _enygmaFactory, address authority_) returns()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackInitialize(endpoint common.Address, tokenRegistryAddress common.Address, enygmaFactory common.Address, authority common.Address) []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("initialize", endpoint, tokenRegistryAddress, enygmaFactory, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackProxiableUUID() []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := enygmaTokenManagerV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackRegisterEnygmaToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x62abdc39.
//
// Solidity: function registerEnygmaToken(bytes32 _resourceId, (string,string,string,bytes,uint256,address,bytes,bytes,bool,uint8,bool,address) tokenData, uint256 , address owner, address participantStorage) returns()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackRegisterEnygmaToken(resourceId [32]byte, tokenData SharedObjectsTokenRegistrationData, arg2 *big.Int, owner common.Address, participantStorage common.Address) []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("registerEnygmaToken", resourceId, tokenData, arg2, owner, participantStorage)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackResourceId() []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := enygmaTokenManagerV1.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSetEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdbbb4155.
//
// Solidity: function setEndpoint(address _endpoint) returns()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackSetEndpoint(endpoint common.Address) []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("setEndpoint", endpoint)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetEnygmaFactory is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x153f716d.
//
// Solidity: function setEnygmaFactory(address _enygmaFactory) returns()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackSetEnygmaFactory(enygmaFactory common.Address) []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("setEnygmaFactory", enygmaFactory)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetTokenCoreAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8f388fc6.
//
// Solidity: function setTokenCoreAddress(address _tokenCoreAddress) returns()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackSetTokenCoreAddress(tokenCoreAddress common.Address) []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("setTokenCoreAddress", tokenCoreAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetTokenRegistryAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x58555ab8.
//
// Solidity: function setTokenRegistryAddress(address _tokenRegistryAddress) returns()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackSetTokenRegistryAddress(tokenRegistryAddress common.Address) []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("setTokenRegistryAddress", tokenRegistryAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackTokenCoreAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6de9a6ee.
//
// Solidity: function tokenCoreAddress() view returns(address)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackTokenCoreAddress() []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("tokenCoreAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenCoreAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6de9a6ee.
//
// Solidity: function tokenCoreAddress() view returns(address)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackTokenCoreAddress(data []byte) (common.Address, error) {
	out, err := enygmaTokenManagerV1.abi.Unpack("tokenCoreAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackTokenRegistryAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5be2aca0.
//
// Solidity: function tokenRegistryAddress() view returns(address)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackTokenRegistryAddress() []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("tokenRegistryAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenRegistryAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5be2aca0.
//
// Solidity: function tokenRegistryAddress() view returns(address)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackTokenRegistryAddress(data []byte) (common.Address, error) {
	out, err := enygmaTokenManagerV1.abi.Unpack("tokenRegistryAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := enygmaTokenManagerV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// EnygmaTokenManagerV1AuthorityUpdated represents a AuthorityUpdated event raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const EnygmaTokenManagerV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (EnygmaTokenManagerV1AuthorityUpdated) ContractEventName() string {
	return EnygmaTokenManagerV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*EnygmaTokenManagerV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != enygmaTokenManagerV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaTokenManagerV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaTokenManagerV1.abi.Events[event].Inputs {
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

// EnygmaTokenManagerV1EnygmaFactoryUpdated represents a EnygmaFactoryUpdated event raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1EnygmaFactoryUpdated struct {
	OldFactory common.Address
	NewFactory common.Address
	Raw        *types.Log // Blockchain specific contextual infos
}

const EnygmaTokenManagerV1EnygmaFactoryUpdatedEventName = "EnygmaFactoryUpdated"

// ContractEventName returns the user-defined event name.
func (EnygmaTokenManagerV1EnygmaFactoryUpdated) ContractEventName() string {
	return EnygmaTokenManagerV1EnygmaFactoryUpdatedEventName
}

// UnpackEnygmaFactoryUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaFactoryUpdated(address indexed oldFactory, address indexed newFactory)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackEnygmaFactoryUpdatedEvent(log *types.Log) (*EnygmaTokenManagerV1EnygmaFactoryUpdated, error) {
	event := "EnygmaFactoryUpdated"
	if log.Topics[0] != enygmaTokenManagerV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaTokenManagerV1EnygmaFactoryUpdated)
	if len(log.Data) > 0 {
		if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaTokenManagerV1.abi.Events[event].Inputs {
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

// EnygmaTokenManagerV1EnygmaTokenRegistered represents a EnygmaTokenRegistered event raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1EnygmaTokenRegistered struct {
	ResourceId    [32]byte
	IssuerChainId *big.Int
	BlockNumber   *big.Int
	Name          string
	InitialSupply *big.Int
	Raw           *types.Log // Blockchain specific contextual infos
}

const EnygmaTokenManagerV1EnygmaTokenRegisteredEventName = "EnygmaTokenRegistered"

// ContractEventName returns the user-defined event name.
func (EnygmaTokenManagerV1EnygmaTokenRegistered) ContractEventName() string {
	return EnygmaTokenManagerV1EnygmaTokenRegisteredEventName
}

// UnpackEnygmaTokenRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaTokenRegistered(bytes32 resourceId, uint256 indexed issuerChainId, uint256 blockNumber, string name, uint256 initialSupply)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackEnygmaTokenRegisteredEvent(log *types.Log) (*EnygmaTokenManagerV1EnygmaTokenRegistered, error) {
	event := "EnygmaTokenRegistered"
	if log.Topics[0] != enygmaTokenManagerV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaTokenManagerV1EnygmaTokenRegistered)
	if len(log.Data) > 0 {
		if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaTokenManagerV1.abi.Events[event].Inputs {
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

// EnygmaTokenManagerV1Initialized represents a Initialized event raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const EnygmaTokenManagerV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (EnygmaTokenManagerV1Initialized) ContractEventName() string {
	return EnygmaTokenManagerV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackInitializedEvent(log *types.Log) (*EnygmaTokenManagerV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != enygmaTokenManagerV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaTokenManagerV1Initialized)
	if len(log.Data) > 0 {
		if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaTokenManagerV1.abi.Events[event].Inputs {
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

// EnygmaTokenManagerV1Upgraded represents a Upgraded event raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const EnygmaTokenManagerV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (EnygmaTokenManagerV1Upgraded) ContractEventName() string {
	return EnygmaTokenManagerV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackUpgradedEvent(log *types.Log) (*EnygmaTokenManagerV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != enygmaTokenManagerV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaTokenManagerV1Upgraded)
	if len(log.Data) > 0 {
		if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaTokenManagerV1.abi.Events[event].Inputs {
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
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["EnygmaTokenManagerV1UnauthorizedCaller"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackEnygmaTokenManagerV1UnauthorizedCallerError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["EnygmaTokenManagerV1ZeroAddress"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackEnygmaTokenManagerV1ZeroAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTokenManagerV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return enygmaTokenManagerV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// EnygmaTokenManagerV1AddressEmptyCode represents a AddressEmptyCode error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func EnygmaTokenManagerV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackAddressEmptyCodeError(raw []byte) (*EnygmaTokenManagerV1AddressEmptyCode, error) {
	out := new(EnygmaTokenManagerV1AddressEmptyCode)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTokenManagerV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func EnygmaTokenManagerV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackERC1967InvalidImplementationError(raw []byte) (*EnygmaTokenManagerV1ERC1967InvalidImplementation, error) {
	out := new(EnygmaTokenManagerV1ERC1967InvalidImplementation)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTokenManagerV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func EnygmaTokenManagerV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackERC1967NonPayableError(raw []byte) (*EnygmaTokenManagerV1ERC1967NonPayable, error) {
	out := new(EnygmaTokenManagerV1ERC1967NonPayable)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTokenManagerV1EnygmaTokenManagerV1UnauthorizedCaller represents a EnygmaTokenManagerV1__UnauthorizedCaller error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1EnygmaTokenManagerV1UnauthorizedCaller struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EnygmaTokenManagerV1__UnauthorizedCaller(address caller)
func EnygmaTokenManagerV1EnygmaTokenManagerV1UnauthorizedCallerErrorID() common.Hash {
	return common.HexToHash("0x9dd326b0722594764594c09d724fc45e8766b44881d5bd3546de6cd55b3e7b0e")
}

// UnpackEnygmaTokenManagerV1UnauthorizedCallerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EnygmaTokenManagerV1__UnauthorizedCaller(address caller)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackEnygmaTokenManagerV1UnauthorizedCallerError(raw []byte) (*EnygmaTokenManagerV1EnygmaTokenManagerV1UnauthorizedCaller, error) {
	out := new(EnygmaTokenManagerV1EnygmaTokenManagerV1UnauthorizedCaller)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "EnygmaTokenManagerV1UnauthorizedCaller", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTokenManagerV1EnygmaTokenManagerV1ZeroAddress represents a EnygmaTokenManagerV1__ZeroAddress error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1EnygmaTokenManagerV1ZeroAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EnygmaTokenManagerV1__ZeroAddress()
func EnygmaTokenManagerV1EnygmaTokenManagerV1ZeroAddressErrorID() common.Hash {
	return common.HexToHash("0xbc624a30f8ded2af488b019c17bb868494b0513c84e7ade5285a6a3c302b6289")
}

// UnpackEnygmaTokenManagerV1ZeroAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EnygmaTokenManagerV1__ZeroAddress()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackEnygmaTokenManagerV1ZeroAddressError(raw []byte) (*EnygmaTokenManagerV1EnygmaTokenManagerV1ZeroAddress, error) {
	out := new(EnygmaTokenManagerV1EnygmaTokenManagerV1ZeroAddress)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "EnygmaTokenManagerV1ZeroAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTokenManagerV1FailedCall represents a FailedCall error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func EnygmaTokenManagerV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackFailedCallError(raw []byte) (*EnygmaTokenManagerV1FailedCall, error) {
	out := new(EnygmaTokenManagerV1FailedCall)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTokenManagerV1InvalidInitialization represents a InvalidInitialization error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func EnygmaTokenManagerV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackInvalidInitializationError(raw []byte) (*EnygmaTokenManagerV1InvalidInitialization, error) {
	out := new(EnygmaTokenManagerV1InvalidInitialization)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTokenManagerV1NotInitializing represents a NotInitializing error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func EnygmaTokenManagerV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackNotInitializingError(raw []byte) (*EnygmaTokenManagerV1NotInitializing, error) {
	out := new(EnygmaTokenManagerV1NotInitializing)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTokenManagerV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func EnygmaTokenManagerV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*EnygmaTokenManagerV1RaylsAccessManagedContractPaused, error) {
	out := new(EnygmaTokenManagerV1RaylsAccessManagedContractPaused)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTokenManagerV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func EnygmaTokenManagerV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*EnygmaTokenManagerV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(EnygmaTokenManagerV1RaylsAccessManagedInvalidAuthority)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTokenManagerV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func EnygmaTokenManagerV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*EnygmaTokenManagerV1RaylsAccessManagedMustSchedule, error) {
	out := new(EnygmaTokenManagerV1RaylsAccessManagedMustSchedule)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTokenManagerV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func EnygmaTokenManagerV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*EnygmaTokenManagerV1RaylsAccessManagedUnauthorized, error) {
	out := new(EnygmaTokenManagerV1RaylsAccessManagedUnauthorized)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTokenManagerV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func EnygmaTokenManagerV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*EnygmaTokenManagerV1UUPSUnauthorizedCallContext, error) {
	out := new(EnygmaTokenManagerV1UUPSUnauthorizedCallContext)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTokenManagerV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the EnygmaTokenManagerV1 contract.
type EnygmaTokenManagerV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func EnygmaTokenManagerV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (enygmaTokenManagerV1 *EnygmaTokenManagerV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*EnygmaTokenManagerV1UUPSUnsupportedProxiableUUID, error) {
	out := new(EnygmaTokenManagerV1UUPSUnsupportedProxiableUUID)
	if err := enygmaTokenManagerV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
