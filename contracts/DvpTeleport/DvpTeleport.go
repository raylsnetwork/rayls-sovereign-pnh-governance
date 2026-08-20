// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package DvpTeleport

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

// DvpTeleportMetaData contains all meta data concerning the DvpTeleport contract.
var DvpTeleportMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"emitCommitments\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenType\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"treeNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitNullifiers\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenType\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitSwapCancelled\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitSwapCompleted\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitSwapInitiated\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ctxt\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"responderCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expiresAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitSwapTimedOut\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ercDvpBalanceUpdated\",\"inputs\":[{\"name\":\"encryptedMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Commitments\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenType\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"treeNumber\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ERCDvpBalanceUpdated\",\"inputs\":[{\"name\":\"encryptedMessage\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Nullifiers\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenType\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SwapCancelled\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SwapCompleted\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SwapInitiated\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"ctxt\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"responderCommitment\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"expiresAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SwapTimedOut\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "DvpTeleport",
	Bin: "0x608060405234801561001057600080fd5b50604051610b0f380380610b0f83398101604081905261002f9161013d565b6100388161003e565b50610194565b6001600160a01b03811661007457604051638944034760e01b81526001600160a01b038216600482015260240160405180910390fd5b600061007e6100db565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b60008060ff1961010c60017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f3561016d565b60405160200161011e91815260200190565b60408051601f1981840301815291905280516020909101201692915050565b60006020828403121561014f57600080fd5b81516001600160a01b038116811461016657600080fd5b9392505050565b8181038181111561018e57634e487b7160e01b600052601160045260246000fd5b92915050565b61096c806101a36000396000f3fe608060405234801561001057600080fd5b50600436106100775760003560e01c8062e599681461007c5780633cad97741461009157806345c73e6b146100a45780638d7a0c1a146100b7578063bf7e214f146100ca578063d43ff248146100e8578063d6863696146100fb578063f00536921461010e575b600080fd5b61008f61008a36600461058d565b610121565b005b61008f61009f366004610615565b610182565b61008f6100b2366004610660565b6101d7565b61008f6100c5366004610660565b61021b565b6100d261025f565b6040516100df9190610679565b60405180910390f35b61008f6100f63660046106ed565b610278565b61008f610109366004610754565b6102db565b61008f61011c3660046107ad565b61033c565b610137336000356001600160e01b03191661038f565b867f6b0b79924df5c2d0907b21c558dae85b33ec7ffa79d50065d83b309b9744386a87878787878760405161017196959493929190610817565b60405180910390a250505050505050565b610198336000356001600160e01b03191661038f565b827fc8c1414976453b7f8359322b6c644cbb667a30bab1ae7bc95e7b856d17d5ca8b83836040516101ca929190610854565b60405180910390a2505050565b6101ed336000356001600160e01b03191661038f565b60405181907f5a399591cfd74c375a1ffd61c20221db0db82381f65516889ed13c8cd1f99d5b90600090a250565b610231336000356001600160e01b03191661038f565b60405181907f5d468243e9af041dda86441da7983b63fe758c86c95c487b7df71364b6c5591090600090a250565b60006102696104e3565b546001600160a01b0316919050565b61028e336000356001600160e01b03191661038f565b82856001600160a01b03167f586ab18efd7bdd12e2b2064a7a74c9fabbb8c5a98d483ec5ec25899388231cc48685856040516102cc93929190610870565b60405180910390a35050505050565b6102f1336000356001600160e01b03191661038f565b836001600160a01b03167f259f270c5a1d564ded25ccac004efc36cde097fc3840157651761576f40e983584848460405161032e93929190610870565b60405180910390a250505050565b610352336000356001600160e01b03191661038f565b7ffb0e9ee6fd84f8120c19ce15ce8d96dee834b491f3449e96e6f31b06a7faab218282604051610383929190610854565b60405180910390a15050565b60006103996104e3565b80549091506001600160a01b0316806103d1576000604051638944034760e01b81526004016103c89190610679565b60405180910390fd5b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015610435573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061045991906108c1565b925092509250826104da5780156104835760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff8216156104bf5760405163a426878960e01b81526001600160a01b038816600482015263ffffffff831660248201526044016103c8565b86604051632ecd3d0360e21b81526004016103c89190610679565b50505050505050565b60008060ff1961051460017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f3561090f565b60405160200161052691815260200190565b60408051601f1981840301815291905280516020909101201692915050565b60008083601f84011261055757600080fd5b5081356001600160401b0381111561056e57600080fd5b60208301915083602082850101111561058657600080fd5b9250929050565b600080600080600080600060a0888a0312156105a857600080fd5b8735965060208801356001600160401b03808211156105c657600080fd5b6105d28b838c01610545565b909850965060408a01359150808211156105eb57600080fd5b506105f88a828b01610545565b989b979a5095989597966060870135966080013595509350505050565b60008060006040848603121561062a57600080fd5b8335925060208401356001600160401b0381111561064757600080fd5b61065386828701610545565b9497909650939450505050565b60006020828403121561067257600080fd5b5035919050565b6001600160a01b0391909116815260200190565b80356001600160a01b03811681146106a457600080fd5b919050565b60008083601f8401126106bb57600080fd5b5081356001600160401b038111156106d257600080fd5b6020830191508360208260051b850101111561058657600080fd5b60008060008060006080868803121561070557600080fd5b61070e8661068d565b9450602086013593506040860135925060608601356001600160401b0381111561073757600080fd5b610743888289016106a9565b969995985093965092949392505050565b6000806000806060858703121561076a57600080fd5b6107738561068d565b93506020850135925060408501356001600160401b0381111561079557600080fd5b6107a1878288016106a9565b95989497509550505050565b600080602083850312156107c057600080fd5b82356001600160401b038111156107d657600080fd5b6107e285828601610545565b90969095509350505050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b60808152600061082b60808301888a6107ee565b828103602084015261083e8187896107ee565b6040840195909552505060600152949350505050565b6020815260006108686020830184866107ee565b949350505050565b838152604060208201819052810182905260006001600160fb1b0383111561089757600080fd5b8260051b8085606085013791909101606001949350505050565b805180151581146106a457600080fd5b6000806000606084860312156108d657600080fd5b6108df846108b1565b9250602084015163ffffffff811681146108f857600080fd5b9150610906604085016108b1565b90509250925092565b8181038181111561093057634e487b7160e01b600052601160045260246000fd5b9291505056fea264697066735822122037d08a8e7366df26d6ae552169bfce67c9c25d655fe1fb6fb2fbf0e41865d87f64736f6c63430008180033",
}

// DvpTeleport is an auto generated Go binding around an Ethereum contract.
type DvpTeleport struct {
	abi abi.ABI
}

// NewDvpTeleport creates a new instance of DvpTeleport.
func NewDvpTeleport() *DvpTeleport {
	parsed, err := DvpTeleportMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &DvpTeleport{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *DvpTeleport) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address authority_) returns()
func (dvpTeleport *DvpTeleport) PackConstructor(authority_ common.Address) []byte {
	enc, err := dvpTeleport.abi.Pack("", authority_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (dvpTeleport *DvpTeleport) PackAuthority() []byte {
	enc, err := dvpTeleport.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (dvpTeleport *DvpTeleport) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := dvpTeleport.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackEmitCommitments is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd43ff248.
//
// Solidity: function emitCommitments(address tokenAddress, uint256 tokenType, uint256 treeNumber, uint256[] commitments) returns()
func (dvpTeleport *DvpTeleport) PackEmitCommitments(tokenAddress common.Address, tokenType *big.Int, treeNumber *big.Int, commitments []*big.Int) []byte {
	enc, err := dvpTeleport.abi.Pack("emitCommitments", tokenAddress, tokenType, treeNumber, commitments)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEmitNullifiers is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd6863696.
//
// Solidity: function emitNullifiers(address tokenAddress, uint256 tokenType, uint256[] nullifiers) returns()
func (dvpTeleport *DvpTeleport) PackEmitNullifiers(tokenAddress common.Address, tokenType *big.Int, nullifiers []*big.Int) []byte {
	enc, err := dvpTeleport.abi.Pack("emitNullifiers", tokenAddress, tokenType, nullifiers)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEmitSwapCancelled is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x45c73e6b.
//
// Solidity: function emitSwapCancelled(bytes32 sharedId) returns()
func (dvpTeleport *DvpTeleport) PackEmitSwapCancelled(sharedId [32]byte) []byte {
	enc, err := dvpTeleport.abi.Pack("emitSwapCancelled", sharedId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEmitSwapCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3cad9774.
//
// Solidity: function emitSwapCompleted(bytes32 sharedId, bytes encryptedData) returns()
func (dvpTeleport *DvpTeleport) PackEmitSwapCompleted(sharedId [32]byte, encryptedData []byte) []byte {
	enc, err := dvpTeleport.abi.Pack("emitSwapCompleted", sharedId, encryptedData)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEmitSwapInitiated is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x00e59968.
//
// Solidity: function emitSwapInitiated(bytes32 sharedId, bytes encryptedData, bytes ctxt, uint256 responderCommitment, uint256 expiresAt) returns()
func (dvpTeleport *DvpTeleport) PackEmitSwapInitiated(sharedId [32]byte, encryptedData []byte, ctxt []byte, responderCommitment *big.Int, expiresAt *big.Int) []byte {
	enc, err := dvpTeleport.abi.Pack("emitSwapInitiated", sharedId, encryptedData, ctxt, responderCommitment, expiresAt)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEmitSwapTimedOut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8d7a0c1a.
//
// Solidity: function emitSwapTimedOut(bytes32 sharedId) returns()
func (dvpTeleport *DvpTeleport) PackEmitSwapTimedOut(sharedId [32]byte) []byte {
	enc, err := dvpTeleport.abi.Pack("emitSwapTimedOut", sharedId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackErcDvpBalanceUpdated is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf0053692.
//
// Solidity: function ercDvpBalanceUpdated(bytes encryptedMessage) returns()
func (dvpTeleport *DvpTeleport) PackErcDvpBalanceUpdated(encryptedMessage []byte) []byte {
	enc, err := dvpTeleport.abi.Pack("ercDvpBalanceUpdated", encryptedMessage)
	if err != nil {
		panic(err)
	}
	return enc
}

// DvpTeleportAuthorityUpdated represents a AuthorityUpdated event raised by the DvpTeleport contract.
type DvpTeleportAuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const DvpTeleportAuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (DvpTeleportAuthorityUpdated) ContractEventName() string {
	return DvpTeleportAuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (dvpTeleport *DvpTeleport) UnpackAuthorityUpdatedEvent(log *types.Log) (*DvpTeleportAuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != dvpTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpTeleportAuthorityUpdated)
	if len(log.Data) > 0 {
		if err := dvpTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvpTeleport.abi.Events[event].Inputs {
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

// DvpTeleportCommitments represents a Commitments event raised by the DvpTeleport contract.
type DvpTeleportCommitments struct {
	TokenAddress common.Address
	TokenType    *big.Int
	TreeNumber   *big.Int
	Commitments  []*big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const DvpTeleportCommitmentsEventName = "Commitments"

// ContractEventName returns the user-defined event name.
func (DvpTeleportCommitments) ContractEventName() string {
	return DvpTeleportCommitmentsEventName
}

// UnpackCommitmentsEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Commitments(address indexed tokenAddress, uint256 tokenType, uint256 indexed treeNumber, uint256[] commitments)
func (dvpTeleport *DvpTeleport) UnpackCommitmentsEvent(log *types.Log) (*DvpTeleportCommitments, error) {
	event := "Commitments"
	if log.Topics[0] != dvpTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpTeleportCommitments)
	if len(log.Data) > 0 {
		if err := dvpTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvpTeleport.abi.Events[event].Inputs {
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

// DvpTeleportERCDvpBalanceUpdated represents a ERCDvpBalanceUpdated event raised by the DvpTeleport contract.
type DvpTeleportERCDvpBalanceUpdated struct {
	EncryptedMessage []byte
	Raw              *types.Log // Blockchain specific contextual infos
}

const DvpTeleportERCDvpBalanceUpdatedEventName = "ERCDvpBalanceUpdated"

// ContractEventName returns the user-defined event name.
func (DvpTeleportERCDvpBalanceUpdated) ContractEventName() string {
	return DvpTeleportERCDvpBalanceUpdatedEventName
}

// UnpackERCDvpBalanceUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ERCDvpBalanceUpdated(bytes encryptedMessage)
func (dvpTeleport *DvpTeleport) UnpackERCDvpBalanceUpdatedEvent(log *types.Log) (*DvpTeleportERCDvpBalanceUpdated, error) {
	event := "ERCDvpBalanceUpdated"
	if log.Topics[0] != dvpTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpTeleportERCDvpBalanceUpdated)
	if len(log.Data) > 0 {
		if err := dvpTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvpTeleport.abi.Events[event].Inputs {
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

// DvpTeleportNullifiers represents a Nullifiers event raised by the DvpTeleport contract.
type DvpTeleportNullifiers struct {
	TokenAddress common.Address
	TokenType    *big.Int
	Nullifiers   []*big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const DvpTeleportNullifiersEventName = "Nullifiers"

// ContractEventName returns the user-defined event name.
func (DvpTeleportNullifiers) ContractEventName() string {
	return DvpTeleportNullifiersEventName
}

// UnpackNullifiersEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Nullifiers(address indexed tokenAddress, uint256 tokenType, uint256[] nullifiers)
func (dvpTeleport *DvpTeleport) UnpackNullifiersEvent(log *types.Log) (*DvpTeleportNullifiers, error) {
	event := "Nullifiers"
	if log.Topics[0] != dvpTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpTeleportNullifiers)
	if len(log.Data) > 0 {
		if err := dvpTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvpTeleport.abi.Events[event].Inputs {
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

// DvpTeleportSwapCancelled represents a SwapCancelled event raised by the DvpTeleport contract.
type DvpTeleportSwapCancelled struct {
	SharedId [32]byte
	Raw      *types.Log // Blockchain specific contextual infos
}

const DvpTeleportSwapCancelledEventName = "SwapCancelled"

// ContractEventName returns the user-defined event name.
func (DvpTeleportSwapCancelled) ContractEventName() string {
	return DvpTeleportSwapCancelledEventName
}

// UnpackSwapCancelledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SwapCancelled(bytes32 indexed sharedId)
func (dvpTeleport *DvpTeleport) UnpackSwapCancelledEvent(log *types.Log) (*DvpTeleportSwapCancelled, error) {
	event := "SwapCancelled"
	if log.Topics[0] != dvpTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpTeleportSwapCancelled)
	if len(log.Data) > 0 {
		if err := dvpTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvpTeleport.abi.Events[event].Inputs {
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

// DvpTeleportSwapCompleted represents a SwapCompleted event raised by the DvpTeleport contract.
type DvpTeleportSwapCompleted struct {
	SharedId      [32]byte
	EncryptedData []byte
	Raw           *types.Log // Blockchain specific contextual infos
}

const DvpTeleportSwapCompletedEventName = "SwapCompleted"

// ContractEventName returns the user-defined event name.
func (DvpTeleportSwapCompleted) ContractEventName() string {
	return DvpTeleportSwapCompletedEventName
}

// UnpackSwapCompletedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SwapCompleted(bytes32 indexed sharedId, bytes encryptedData)
func (dvpTeleport *DvpTeleport) UnpackSwapCompletedEvent(log *types.Log) (*DvpTeleportSwapCompleted, error) {
	event := "SwapCompleted"
	if log.Topics[0] != dvpTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpTeleportSwapCompleted)
	if len(log.Data) > 0 {
		if err := dvpTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvpTeleport.abi.Events[event].Inputs {
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

// DvpTeleportSwapInitiated represents a SwapInitiated event raised by the DvpTeleport contract.
type DvpTeleportSwapInitiated struct {
	SharedId            [32]byte
	EncryptedData       []byte
	Ctxt                []byte
	ResponderCommitment *big.Int
	ExpiresAt           *big.Int
	Raw                 *types.Log // Blockchain specific contextual infos
}

const DvpTeleportSwapInitiatedEventName = "SwapInitiated"

// ContractEventName returns the user-defined event name.
func (DvpTeleportSwapInitiated) ContractEventName() string {
	return DvpTeleportSwapInitiatedEventName
}

// UnpackSwapInitiatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SwapInitiated(bytes32 indexed sharedId, bytes encryptedData, bytes ctxt, uint256 responderCommitment, uint256 expiresAt)
func (dvpTeleport *DvpTeleport) UnpackSwapInitiatedEvent(log *types.Log) (*DvpTeleportSwapInitiated, error) {
	event := "SwapInitiated"
	if log.Topics[0] != dvpTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpTeleportSwapInitiated)
	if len(log.Data) > 0 {
		if err := dvpTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvpTeleport.abi.Events[event].Inputs {
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

// DvpTeleportSwapTimedOut represents a SwapTimedOut event raised by the DvpTeleport contract.
type DvpTeleportSwapTimedOut struct {
	SharedId [32]byte
	Raw      *types.Log // Blockchain specific contextual infos
}

const DvpTeleportSwapTimedOutEventName = "SwapTimedOut"

// ContractEventName returns the user-defined event name.
func (DvpTeleportSwapTimedOut) ContractEventName() string {
	return DvpTeleportSwapTimedOutEventName
}

// UnpackSwapTimedOutEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SwapTimedOut(bytes32 indexed sharedId)
func (dvpTeleport *DvpTeleport) UnpackSwapTimedOutEvent(log *types.Log) (*DvpTeleportSwapTimedOut, error) {
	event := "SwapTimedOut"
	if log.Topics[0] != dvpTeleport.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpTeleportSwapTimedOut)
	if len(log.Data) > 0 {
		if err := dvpTeleport.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvpTeleport.abi.Events[event].Inputs {
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
func (dvpTeleport *DvpTeleport) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], dvpTeleport.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return dvpTeleport.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvpTeleport.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return dvpTeleport.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvpTeleport.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return dvpTeleport.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvpTeleport.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return dvpTeleport.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// DvpTeleportRaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the DvpTeleport contract.
type DvpTeleportRaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func DvpTeleportRaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (dvpTeleport *DvpTeleport) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*DvpTeleportRaylsAccessManagedContractPaused, error) {
	out := new(DvpTeleportRaylsAccessManagedContractPaused)
	if err := dvpTeleport.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpTeleportRaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the DvpTeleport contract.
type DvpTeleportRaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func DvpTeleportRaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (dvpTeleport *DvpTeleport) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*DvpTeleportRaylsAccessManagedInvalidAuthority, error) {
	out := new(DvpTeleportRaylsAccessManagedInvalidAuthority)
	if err := dvpTeleport.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpTeleportRaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the DvpTeleport contract.
type DvpTeleportRaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func DvpTeleportRaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (dvpTeleport *DvpTeleport) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*DvpTeleportRaylsAccessManagedMustSchedule, error) {
	out := new(DvpTeleportRaylsAccessManagedMustSchedule)
	if err := dvpTeleport.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpTeleportRaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the DvpTeleport contract.
type DvpTeleportRaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func DvpTeleportRaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (dvpTeleport *DvpTeleport) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*DvpTeleportRaylsAccessManagedUnauthorized, error) {
	out := new(DvpTeleportRaylsAccessManagedUnauthorized)
	if err := dvpTeleport.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}
