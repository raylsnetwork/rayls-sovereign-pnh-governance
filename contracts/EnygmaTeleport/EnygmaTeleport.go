// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package EnygmaTeleport

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

// IEnygmaV1EnygmaPointWithChainId is an auto generated low-level Go binding around an user-defined struct.
type IEnygmaV1EnygmaPointWithChainId struct {
	C1      *big.Int
	C2      *big.Int
	ChainId *big.Int
}

// IEnygmaV1SupplyUpdateTx is an auto generated low-level Go binding around an user-defined struct.
type IEnygmaV1SupplyUpdateTx struct {
	Amount *big.Int
	TxType uint8
}

// EnygmaTeleportMetaData contains all meta data concerning the EnygmaTeleport contract.
var EnygmaTeleportMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"enygmaDvpBalanceUpdated\",\"inputs\":[{\"name\":\"encryptedMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"enygmaSupplyUpdated\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"update\",\"type\":\"tuple\",\"internalType\":\"structIEnygmaV1.SupplyUpdateTx\",\"components\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"txType\",\"type\":\"uint8\",\"internalType\":\"enumIEnygmaV1.TxType\"}]},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"enygmaTransferCompleted\",\"inputs\":[{\"name\":\"encryptedMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"finalizeBalances\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"finalizedBlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pendingBlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"balances\",\"type\":\"tuple[]\",\"internalType\":\"structIEnygmaV1.EnygmaPointWithChainId[]\",\"components\":[{\"name\":\"c1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"c2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"encryptedMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"messageTag\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"anonymitySet\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"arrayHashSecrets\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"toChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BalancesFinalized\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"finalizedBlockNumber\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"pendingBlockNumber\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"balances\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structIEnygmaV1.EnygmaPointWithChainId[]\",\"components\":[{\"name\":\"c1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"c2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaDvpBalanceUpdated\",\"inputs\":[{\"name\":\"encryptedMessage\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaSupplyUpdated\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"update\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIEnygmaV1.SupplyUpdateTx\",\"components\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"txType\",\"type\":\"uint8\",\"internalType\":\"enumIEnygmaV1.TxType\"}]},{\"name\":\"chainId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaTransfer\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"encryptedMessage\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"messageTag\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"anonymitySet\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"arrayHashSecrets\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"toChainId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaTransferCompleted\",\"inputs\":[{\"name\":\"encryptedMessage\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EnygmaTeleport__InvalidEndpointAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "EnygmaTeleport",
	Bin: "0x608060405234801561001057600080fd5b50604051610aa5380380610aa583398101604081905261002f9161013d565b6100388161003e565b50610194565b6001600160a01b03811661007457604051638944034760e01b81526001600160a01b038216600482015260240160405180910390fd5b600061007e6100db565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b60008060ff1961010c60017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f3561016d565b60405160200161011e91815260200190565b60408051601f1981840301815291905280516020909101201692915050565b60006020828403121561014f57600080fd5b81516001600160a01b038116811461016657600080fd5b9392505050565b8181038181111561018e57634e487b7160e01b600052601160045260246000fd5b92915050565b610902806101a36000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80635ff64892146100675780636924a82f1461007c578063952acada1461008f578063bf7e214f146100a2578063c651be9b146100c0578063c77b6c7f146100d3575b600080fd5b61007a6100753660046104b2565b6100e6565b005b61007a61008a3660046104f3565b610139565b61007a61009d3660046105c6565b610193565b6100aa6101fd565b6040516100b79190610688565b60405180910390f35b61007a6100ce36600461069c565b610216565b61007a6100e13660046104b2565b61026d565b6100fc336000356001600160e01b0319166102b4565b7fc41e4c4f147a5116dfa9b8c8c11200b617cb26b0515661df9741bb298cc8ad00828260405161012d92919061070d565b60405180910390a15050565b61014f336000356001600160e01b0319166102b4565b83857fcab9daef74ea3374e6a930e337db120037cc08e628a458f0cd2c756dc74946b985858560405161018493929190610729565b60405180910390a35050505050565b6101a9336000356001600160e01b0319166102b4565b897f2b925a0de145bd51c571f8da5994d78f6cce29dc9f967e20a4251396993756838a8a8a8a8a8a8a8a8a6040516101e9999897969594939291906107b6565b60405180910390a250505050505050505050565b6000610207610408565b546001600160a01b0316919050565b61022c336000356001600160e01b0319166102b4565b82847f34cfbaa12047c291fffeecb25bfc1e6f14aaac101875c7abaab60c5f2bc70cb8848460405161025f929190610815565b60405180910390a350505050565b610283336000356001600160e01b0319166102b4565b7fb01fdaeaafdf95de8d02de23df1f10cabc3ce2fa7b001936dfda154de3f0dff0828260405161012d92919061070d565b60006102be610408565b80549091506001600160a01b0316806102f6576000604051638944034760e01b81526004016102ed9190610688565b60405180910390fd5b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa15801561035a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061037e9190610857565b925092509250826103ff5780156103a85760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff8216156103e45760405163a426878960e01b81526001600160a01b038816600482015263ffffffff831660248201526044016102ed565b86604051632ecd3d0360e21b81526004016102ed9190610688565b50505050505050565b60008060ff1961043960017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f356108a5565b60405160200161044b91815260200190565b60408051601f1981840301815291905280516020909101201692915050565b60008083601f84011261047c57600080fd5b5081356001600160401b0381111561049357600080fd5b6020830191508360208285010111156104ab57600080fd5b9250929050565b600080602083850312156104c557600080fd5b82356001600160401b038111156104db57600080fd5b6104e78582860161046a565b90969095509350505050565b60008060008060006080868803121561050b57600080fd5b85359450602086013593506040860135925060608601356001600160401b038082111561053757600080fd5b818801915088601f83011261054b57600080fd5b81358181111561055a57600080fd5b89602060608302850101111561056f57600080fd5b9699959850939650602001949392505050565b60008083601f84011261059457600080fd5b5081356001600160401b038111156105ab57600080fd5b6020830191508360208260051b85010111156104ab57600080fd5b60008060008060008060008060008060e08b8d0312156105e557600080fd5b8a35995060208b01356001600160401b038082111561060357600080fd5b61060f8e838f0161046a565b909b50995060408d0135985060608d0135975060808d013591508082111561063657600080fd5b6106428e838f01610582565b909750955060a08d013591508082111561065b57600080fd5b506106688d828e01610582565b9150809450508092505060c08b013590509295989b9194979a5092959850565b6001600160a01b0391909116815260200190565b60008060008084860360a08112156106b357600080fd5b85359450602086013593506040603f19820112156106d057600080fd5b509295919450506040830192608001359150565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6020815260006107216020830184866106e4565b949350505050565b838152604060208083018290528282018490526000919060609081850187855b888110156107755781358352838201358484015285820135868401529184019190840190600101610749565b50909998505050505050505050565b81835260006001600160fb1b0383111561079d57600080fd5b8260051b80836020870137939093016020019392505050565b60c0815260006107ca60c083018b8d6106e4565b89602084015288604084015282810360608401526107e981888a610784565b905082810360808401526107fe818688610784565b9150508260a08301529a9950505050505050505050565b823581526060810160208401356006811061082f57600080fd5b6020830152604090910191909152919050565b8051801515811461085257600080fd5b919050565b60008060006060848603121561086c57600080fd5b61087584610842565b9250602084015163ffffffff8116811461088e57600080fd5b915061089c60408501610842565b90509250925092565b818103818111156108c657634e487b7160e01b600052601160045260246000fd5b9291505056fea2646970667358221220116e9f1949c81ace8b9cc8eddf8c530ab4041711f747d26e0d3a3c9de9baa52e64736f6c63430008180033",
}

// EnygmaTeleport is an auto generated Go binding around an Ethereum contract.
type EnygmaTeleport struct {
	abi abi.ABI
}

// NewEnygmaTeleport creates a new instance of EnygmaTeleport.
func NewEnygmaTeleport() *EnygmaTeleport {
	parsed, err := EnygmaTeleportMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &EnygmaTeleport{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *EnygmaTeleport) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address authority_) returns()
func (enygmaTeleport *EnygmaTeleport) PackConstructor(authority_ common.Address) []byte {
	enc, err := enygmaTeleport.abi.Pack("", authority_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (enygmaTeleport *EnygmaTeleport) PackAuthority() []byte {
	enc, err := enygmaTeleport.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (enygmaTeleport *EnygmaTeleport) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := enygmaTeleport.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackEnygmaDvpBalanceUpdated is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc77b6c7f.
//
// Solidity: function enygmaDvpBalanceUpdated(bytes encryptedMessage) returns()
func (enygmaTeleport *EnygmaTeleport) PackEnygmaDvpBalanceUpdated(encryptedMessage []byte) []byte {
	enc, err := enygmaTeleport.abi.Pack("enygmaDvpBalanceUpdated", encryptedMessage)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEnygmaSupplyUpdated is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc651be9b.
//
// Solidity: function enygmaSupplyUpdated(bytes32 resourceId, uint256 blockNumber, (uint256,uint8) update, uint256 chainId) returns()
func (enygmaTeleport *EnygmaTeleport) PackEnygmaSupplyUpdated(resourceId [32]byte, blockNumber *big.Int, update IEnygmaV1SupplyUpdateTx, chainId *big.Int) []byte {
	enc, err := enygmaTeleport.abi.Pack("enygmaSupplyUpdated", resourceId, blockNumber, update, chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEnygmaTransferCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5ff64892.
//
// Solidity: function enygmaTransferCompleted(bytes encryptedMessage) returns()
func (enygmaTeleport *EnygmaTeleport) PackEnygmaTransferCompleted(encryptedMessage []byte) []byte {
	enc, err := enygmaTeleport.abi.Pack("enygmaTransferCompleted", encryptedMessage)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackFinalizeBalances is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6924a82f.
//
// Solidity: function finalizeBalances(bytes32 resourceId, uint256 finalizedBlockNumber, uint256 pendingBlockNumber, (uint256,uint256,uint256)[] balances) returns()
func (enygmaTeleport *EnygmaTeleport) PackFinalizeBalances(resourceId [32]byte, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int, balances []IEnygmaV1EnygmaPointWithChainId) []byte {
	enc, err := enygmaTeleport.abi.Pack("finalizeBalances", resourceId, finalizedBlockNumber, pendingBlockNumber, balances)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x952acada.
//
// Solidity: function transfer(bytes32 resourceId, bytes encryptedMessage, uint256 messageTag, uint256 blockNumber, uint256[] anonymitySet, uint256[] arrayHashSecrets, uint256 toChainId) returns()
func (enygmaTeleport *EnygmaTeleport) PackTransfer(resourceId [32]byte, encryptedMessage []byte, messageTag *big.Int, blockNumber *big.Int, anonymitySet []*big.Int, arrayHashSecrets []*big.Int, toChainId *big.Int) []byte {
	enc, err := enygmaTeleport.abi.Pack("transfer", resourceId, encryptedMessage, messageTag, blockNumber, anonymitySet, arrayHashSecrets, toChainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// EnygmaTeleportAuthorityUpdated represents a AuthorityUpdated event raised by the EnygmaTeleport contract.
type EnygmaTeleportAuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const EnygmaTeleportAuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (EnygmaTeleportAuthorityUpdated) ContractEventName() string {
	return EnygmaTeleportAuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (enygmaTeleport *EnygmaTeleport) UnpackAuthorityUpdatedEvent(log *types.Log) (*EnygmaTeleportAuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != enygmaTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaTeleportAuthorityUpdated)
	if len(log.Data) > 0 {
		if err := enygmaTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaTeleport.abi.Events[event].Inputs {
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

// EnygmaTeleportBalancesFinalized represents a BalancesFinalized event raised by the EnygmaTeleport contract.
type EnygmaTeleportBalancesFinalized struct {
	ResourceId           [32]byte
	FinalizedBlockNumber *big.Int
	PendingBlockNumber   *big.Int
	Balances             []IEnygmaV1EnygmaPointWithChainId
	Raw                  *types.Log // Blockchain specific contextual infos
}

const EnygmaTeleportBalancesFinalizedEventName = "BalancesFinalized"

// ContractEventName returns the user-defined event name.
func (EnygmaTeleportBalancesFinalized) ContractEventName() string {
	return EnygmaTeleportBalancesFinalizedEventName
}

// UnpackBalancesFinalizedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BalancesFinalized(bytes32 indexed resourceId, uint256 indexed finalizedBlockNumber, uint256 pendingBlockNumber, (uint256,uint256,uint256)[] balances)
func (enygmaTeleport *EnygmaTeleport) UnpackBalancesFinalizedEvent(log *types.Log) (*EnygmaTeleportBalancesFinalized, error) {
	event := "BalancesFinalized"
	if log.Topics[0] != enygmaTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaTeleportBalancesFinalized)
	if len(log.Data) > 0 {
		if err := enygmaTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaTeleport.abi.Events[event].Inputs {
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

// EnygmaTeleportEnygmaDvpBalanceUpdated represents a EnygmaDvpBalanceUpdated event raised by the EnygmaTeleport contract.
type EnygmaTeleportEnygmaDvpBalanceUpdated struct {
	EncryptedMessage []byte
	Raw              *types.Log // Blockchain specific contextual infos
}

const EnygmaTeleportEnygmaDvpBalanceUpdatedEventName = "EnygmaDvpBalanceUpdated"

// ContractEventName returns the user-defined event name.
func (EnygmaTeleportEnygmaDvpBalanceUpdated) ContractEventName() string {
	return EnygmaTeleportEnygmaDvpBalanceUpdatedEventName
}

// UnpackEnygmaDvpBalanceUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaDvpBalanceUpdated(bytes encryptedMessage)
func (enygmaTeleport *EnygmaTeleport) UnpackEnygmaDvpBalanceUpdatedEvent(log *types.Log) (*EnygmaTeleportEnygmaDvpBalanceUpdated, error) {
	event := "EnygmaDvpBalanceUpdated"
	if log.Topics[0] != enygmaTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaTeleportEnygmaDvpBalanceUpdated)
	if len(log.Data) > 0 {
		if err := enygmaTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaTeleport.abi.Events[event].Inputs {
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

// EnygmaTeleportEnygmaSupplyUpdated represents a EnygmaSupplyUpdated event raised by the EnygmaTeleport contract.
type EnygmaTeleportEnygmaSupplyUpdated struct {
	ResourceId  [32]byte
	BlockNumber *big.Int
	Update      IEnygmaV1SupplyUpdateTx
	ChainId     *big.Int
	Raw         *types.Log // Blockchain specific contextual infos
}

const EnygmaTeleportEnygmaSupplyUpdatedEventName = "EnygmaSupplyUpdated"

// ContractEventName returns the user-defined event name.
func (EnygmaTeleportEnygmaSupplyUpdated) ContractEventName() string {
	return EnygmaTeleportEnygmaSupplyUpdatedEventName
}

// UnpackEnygmaSupplyUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaSupplyUpdated(bytes32 indexed resourceId, uint256 indexed blockNumber, (uint256,uint8) update, uint256 chainId)
func (enygmaTeleport *EnygmaTeleport) UnpackEnygmaSupplyUpdatedEvent(log *types.Log) (*EnygmaTeleportEnygmaSupplyUpdated, error) {
	event := "EnygmaSupplyUpdated"
	if log.Topics[0] != enygmaTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaTeleportEnygmaSupplyUpdated)
	if len(log.Data) > 0 {
		if err := enygmaTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaTeleport.abi.Events[event].Inputs {
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

// EnygmaTeleportEnygmaTransfer represents a EnygmaTransfer event raised by the EnygmaTeleport contract.
type EnygmaTeleportEnygmaTransfer struct {
	ResourceId       [32]byte
	EncryptedMessage []byte
	MessageTag       *big.Int
	BlockNumber      *big.Int
	AnonymitySet     []*big.Int
	ArrayHashSecrets []*big.Int
	ToChainId        *big.Int
	Raw              *types.Log // Blockchain specific contextual infos
}

const EnygmaTeleportEnygmaTransferEventName = "EnygmaTransfer"

// ContractEventName returns the user-defined event name.
func (EnygmaTeleportEnygmaTransfer) ContractEventName() string {
	return EnygmaTeleportEnygmaTransferEventName
}

// UnpackEnygmaTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaTransfer(bytes32 indexed resourceId, bytes encryptedMessage, uint256 messageTag, uint256 blockNumber, uint256[] anonymitySet, uint256[] arrayHashSecrets, uint256 toChainId)
func (enygmaTeleport *EnygmaTeleport) UnpackEnygmaTransferEvent(log *types.Log) (*EnygmaTeleportEnygmaTransfer, error) {
	event := "EnygmaTransfer"
	if log.Topics[0] != enygmaTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaTeleportEnygmaTransfer)
	if len(log.Data) > 0 {
		if err := enygmaTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaTeleport.abi.Events[event].Inputs {
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

// EnygmaTeleportEnygmaTransferCompleted represents a EnygmaTransferCompleted event raised by the EnygmaTeleport contract.
type EnygmaTeleportEnygmaTransferCompleted struct {
	EncryptedMessage []byte
	Raw              *types.Log // Blockchain specific contextual infos
}

const EnygmaTeleportEnygmaTransferCompletedEventName = "EnygmaTransferCompleted"

// ContractEventName returns the user-defined event name.
func (EnygmaTeleportEnygmaTransferCompleted) ContractEventName() string {
	return EnygmaTeleportEnygmaTransferCompletedEventName
}

// UnpackEnygmaTransferCompletedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaTransferCompleted(bytes encryptedMessage)
func (enygmaTeleport *EnygmaTeleport) UnpackEnygmaTransferCompletedEvent(log *types.Log) (*EnygmaTeleportEnygmaTransferCompleted, error) {
	event := "EnygmaTransferCompleted"
	if log.Topics[0] != enygmaTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaTeleportEnygmaTransferCompleted)
	if len(log.Data) > 0 {
		if err := enygmaTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaTeleport.abi.Events[event].Inputs {
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
func (enygmaTeleport *EnygmaTeleport) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], enygmaTeleport.abi.Errors["EnygmaTeleportInvalidEndpointAddress"].ID.Bytes()[:4]) {
		return enygmaTeleport.UnpackEnygmaTeleportInvalidEndpointAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTeleport.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return enygmaTeleport.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTeleport.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return enygmaTeleport.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTeleport.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return enygmaTeleport.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaTeleport.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return enygmaTeleport.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// EnygmaTeleportEnygmaTeleportInvalidEndpointAddress represents a EnygmaTeleport__InvalidEndpointAddress error raised by the EnygmaTeleport contract.
type EnygmaTeleportEnygmaTeleportInvalidEndpointAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EnygmaTeleport__InvalidEndpointAddress()
func EnygmaTeleportEnygmaTeleportInvalidEndpointAddressErrorID() common.Hash {
	return common.HexToHash("0x9e5ccb78e4f28886cf38604a48f9e9d6a16cc0236493271139be9bbe2ea19d69")
}

// UnpackEnygmaTeleportInvalidEndpointAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EnygmaTeleport__InvalidEndpointAddress()
func (enygmaTeleport *EnygmaTeleport) UnpackEnygmaTeleportInvalidEndpointAddressError(raw []byte) (*EnygmaTeleportEnygmaTeleportInvalidEndpointAddress, error) {
	out := new(EnygmaTeleportEnygmaTeleportInvalidEndpointAddress)
	if err := enygmaTeleport.abi.UnpackIntoInterface(out, "EnygmaTeleportInvalidEndpointAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTeleportRaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the EnygmaTeleport contract.
type EnygmaTeleportRaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func EnygmaTeleportRaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (enygmaTeleport *EnygmaTeleport) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*EnygmaTeleportRaylsAccessManagedContractPaused, error) {
	out := new(EnygmaTeleportRaylsAccessManagedContractPaused)
	if err := enygmaTeleport.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTeleportRaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the EnygmaTeleport contract.
type EnygmaTeleportRaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func EnygmaTeleportRaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (enygmaTeleport *EnygmaTeleport) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*EnygmaTeleportRaylsAccessManagedInvalidAuthority, error) {
	out := new(EnygmaTeleportRaylsAccessManagedInvalidAuthority)
	if err := enygmaTeleport.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTeleportRaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the EnygmaTeleport contract.
type EnygmaTeleportRaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func EnygmaTeleportRaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (enygmaTeleport *EnygmaTeleport) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*EnygmaTeleportRaylsAccessManagedMustSchedule, error) {
	out := new(EnygmaTeleportRaylsAccessManagedMustSchedule)
	if err := enygmaTeleport.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaTeleportRaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the EnygmaTeleport contract.
type EnygmaTeleportRaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func EnygmaTeleportRaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (enygmaTeleport *EnygmaTeleport) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*EnygmaTeleportRaylsAccessManagedUnauthorized, error) {
	out := new(EnygmaTeleportRaylsAccessManagedUnauthorized)
	if err := enygmaTeleport.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}
