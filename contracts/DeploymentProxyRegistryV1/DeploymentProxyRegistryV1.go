// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package DeploymentProxyRegistryV1

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

// DeploymentProxyRegistryV1MetaData contains all meta data concerning the DeploymentProxyRegistryV1 contract.
var DeploymentProxyRegistryV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllContractNames\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllContracts\",\"inputs\":[],\"outputs\":[{\"name\":\"names\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"addresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getContract\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerContract\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerContracts\",\"inputs\":[{\"name\":\"names\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"contractAddresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeContract\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateContract\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"newAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ContractRegistered\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ContractRemoved\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ContractUpdated\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"oldAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "DeploymentProxyRegistryV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b5061001d610022565b6100d4565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156100725760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100d15780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b608051611f4c6100fd600039600081816110050152818161102e01526111620152611f4c6000f3fe6080604052600436106100975760003560e01c806318d3ce961461009c57806335817773146100c85780633e8552af146100f55780634f1ef2861461011757806352d1902d1461012a578063697a60b31461014d5780637f3c2c281461016d57806397623b581461018d578063ad3cb1cc146101ad578063bf7e214f146101eb578063c4d66de814610200578063d7a9b81614610220575b600080fd5b3480156100a857600080fd5b506100b1610242565b6040516100bf929190611618565b60405180910390f35b3480156100d457600080fd5b506100e86100e3366004611728565b6103f5565b6040516100bf9190611764565b34801561010157600080fd5b5061011561011036600461182d565b610425565b005b610115610125366004611900565b61079d565b34801561013657600080fd5b5061013f6107bc565b6040519081526020016100bf565b34801561015957600080fd5b50610115610168366004611957565b6107d9565b34801561017957600080fd5b50610115610188366004611957565b610935565b34801561019957600080fd5b506101156101a8366004611728565b610ad5565b3480156101b957600080fd5b506101de604051806040016040528060058152602001640352e302e360dc1b81525081565b6040516100bf91906119a4565b3480156101f757600080fd5b506100e8610cbb565b34801561020c57600080fd5b5061011561021b3660046119b7565b610cd4565b34801561022c57600080fd5b50610235610dd6565b6040516100bf91906119d2565b6001546060908190806001600160401b038111156102625761026261166b565b60405190808252806020026020018201604052801561028b578160200160208202803683370190505b50915060005b81811015610315576000600182815481106102ae576102ae6119e5565b906000526020600020016040516102c59190611aa8565b9081526040519081900360200190205483516001600160a01b03909116908490839081106102f5576102f56119e5565b6001600160a01b0390921660209283029190910190910152600101610291565b5060018281805480602002602001604051908101604052809291908181526020016000905b828210156103e6578382906000526020600020018054610359906119fb565b80601f0160208091040260200160405190810160405280929190818152602001828054610385906119fb565b80156103d25780601f106103a7576101008083540402835291602001916103d2565b820191906000526020600020905b8154815290600101906020018083116103b557829003601f168201915b50505050508152602001906001019061033a565b50505050915092509250509091565b600080826040516104069190611ab4565b908152604051908190036020019020546001600160a01b031692915050565b61043b336000356001600160e01b031916610eaf565b60008251116104b05760405162461bcd60e51b815260206004820152603660248201527f4465706c6f796d656e7450726f7879526567697374727956313a204e616d65736044820152752061727261792063616e6e6f7420626520656d70747960501b60648201526084015b60405180910390fd5b80518251146105255760405162461bcd60e51b815260206004820152603b60248201527f4465706c6f796d656e7450726f7879526567697374727956313a20417272617960448201527a0e640daeae6e840d0c2ecca40e8d0ca40e6c2daca40d8cadccee8d602b1b60648201526084016104a7565b60005b8251811015610798576000838281518110610545576105456119e5565b6020026020010151511161056b5760405162461bcd60e51b81526004016104a790611ad0565b60006001600160a01b0316828281518110610588576105886119e5565b60200260200101516001600160a01b0316036105b65760405162461bcd60e51b81526004016104a790611b1f565b60006001600160a01b031660008483815181106105d5576105d56119e5565b60200260200101516040516105ea9190611ab4565b908152604051908190036020019020546001600160a01b0316146106205760405162461bcd60e51b81526004016104a790611b72565b6000828281518110610634576106346119e5565b60200260200101516001600160a01b03163b116106635760405162461bcd60e51b81526004016104a790611bbb565b818181518110610675576106756119e5565b60200260200101516000848381518110610691576106916119e5565b60200260200101516040516106a69190611ab4565b908152602001604051809103902060006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060018382815181106106ec576106ec6119e5565b602090810291909101810151825460018101845560009384529190922001906107159082611c5b565b50818181518110610728576107286119e5565b60200260200101516001600160a01b031683828151811061074b5761074b6119e5565b60200260200101516040516107609190611ab4565b604051908190038120907f2a88e68a891ddb61f7aebfdeefd9fb74964fcd5371b692ea59cca73fc58f480590600090a3600101610528565b505050565b6107a5610ffa565b6107ae8261108a565b6107b882826110a3565b5050565b60006107c6611157565b50600080516020611ed783398151915290565b6107ef336000356001600160e01b031916610eaf565b60008251116108105760405162461bcd60e51b81526004016104a790611ad0565b6001600160a01b0381166108365760405162461bcd60e51b81526004016104a790611b1f565b60006001600160a01b03166000836040516108519190611ab4565b908152604051908190036020019020546001600160a01b0316036108875760405162461bcd60e51b81526004016104a790611d14565b600080836040516108989190611ab4565b908152604051908190036020018120546001600160a01b0316915082906000906108c3908690611ab4565b90815260405190819003602001812080546001600160a01b039384166001600160a01b031990911617905583821691831690610900908690611ab4565b604051908190038120907f508230d9e989b49c651e1c7be29de727f16d9c1027f7660466911e442c3121fb90600090a4505050565b61094b336000356001600160e01b031916610eaf565b600082511161096c5760405162461bcd60e51b81526004016104a790611ad0565b6001600160a01b0381166109925760405162461bcd60e51b81526004016104a790611b1f565b60006001600160a01b03166000836040516109ad9190611ab4565b908152604051908190036020019020546001600160a01b0316146109e35760405162461bcd60e51b81526004016104a790611b72565b6000816001600160a01b03163b11610a0d5760405162461bcd60e51b81526004016104a790611bbb565b80600083604051610a1e9190611ab4565b90815260405190819003602001902080546001600160a01b03929092166001600160a01b03199092169190911790556001805480820182556000919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf601610a888382611c5b565b50806001600160a01b031682604051610aa19190611ab4565b604051908190038120907f2a88e68a891ddb61f7aebfdeefd9fb74964fcd5371b692ea59cca73fc58f480590600090a35050565b610aeb336000356001600160e01b031916610eaf565b60006001600160a01b0316600082604051610b069190611ab4565b908152604051908190036020019020546001600160a01b031603610b3c5760405162461bcd60e51b81526004016104a790611d14565b60008082604051610b4d9190611ab4565b908152604051908190036020018120546001600160a01b03169150600090610b76908490611ab4565b90815260405190819003602001902080546001600160a01b031916905560005b600154811015610c6e57828051906020012060018281548110610bbb57610bbb6119e5565b90600052602060002001604051610bd29190611aa8565b604051809103902003610c665760018054610bee908290611d54565b81548110610bfe57610bfe6119e5565b9060005260206000200160018281548110610c1b57610c1b6119e5565b906000526020600020019081610c319190611d75565b506001805480610c4357610c43611e49565b600190038181906000526020600020016000610c5f9190611520565b9055610c6e565b600101610b96565b50806001600160a01b031682604051610c879190611ab4565b604051908190038120907f7ad6d7c73213a93d033801d72c729d90fb72104796b75ad799858b26f2b9872f90600090a35050565b6000610cc56111a0565b546001600160a01b0316919050565b6000610cde611202565b805490915060ff600160401b82041615906001600160401b0316600081158015610d055750825b90506000826001600160401b03166001148015610d215750303b155b905081158015610d2f575080155b15610d4d5760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff191660011785558315610d7757845460ff60401b1916600160401b1785555b610d7f61122d565b610d8886611235565b8315610dce57845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050565b60606001805480602002602001604051908101604052809291908181526020016000905b82821015610ea6578382906000526020600020018054610e19906119fb565b80601f0160208091040260200160405190810160405280929190818152602001828054610e45906119fb565b8015610e925780601f10610e6757610100808354040283529160200191610e92565b820191906000526020600020905b815481529060010190602001808311610e7557829003601f168201915b505050505081526020019060010190610dfa565b50505050905090565b6000610eb96111a0565b80549091506001600160a01b031680610ee8576000604051638944034760e01b81526004016104a79190611764565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015610f4c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f709190611e6f565b92509250925082610ff1578015610f9a5760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff821615610fd65760405163a426878960e01b81526001600160a01b038816600482015263ffffffff831660248201526044016104a7565b86604051632ecd3d0360e21b81526004016104a79190611764565b50505050505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016148061106a57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661105e611276565b6001600160a01b031614155b156110885760405163703e46dd60e11b815260040160405180910390fd5b565b6110a0336000356001600160e01b031916610eaf565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa9250505080156110fd575060408051601f3d908101601f191682019092526110fa91810190611ebd565b60015b61111c5781604051634c9c8ce360e01b81526004016104a79190611764565b600080516020611ed7833981519152811461114d57604051632a87526960e21b8152600481018290526024016104a7565b610798838361128c565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146110885760405163703e46dd60e11b815260040160405180910390fd5b60008060ff196111d160017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35611d54565b6040516020016111e391815260200190565b60408051601f1981840301815291905280516020909101201692915050565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b6110886112e2565b600061123f6111a0565b80549091506001600160a01b03161561126d5781604051638944034760e01b81526004016104a79190611764565b6107b882611307565b6000600080516020611ed7833981519152610cc5565b61129582611397565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a28051156112da5761079882826113f3565b6107b8611469565b6112ea611488565b61108857604051631afcd79f60e31b815260040160405180910390fd5b6001600160a01b0381166113305780604051638944034760e01b81526004016104a79190611764565b600061133a6111a0565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b806001600160a01b03163b6000036113c45780604051634c9c8ce360e01b81526004016104a79190611764565b600080516020611ed783398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b0316846040516114109190611ab4565b600060405180830381855af49150503d806000811461144b576040519150601f19603f3d011682016040523d82523d6000602084013e611450565b606091505b50915091506114608583836114a2565b95945050505050565b34156110885760405163b398979f60e01b815260040160405180910390fd5b6000611492611202565b54600160401b900460ff16919050565b6060826114b7576114b2826114f8565b6114f1565b81511580156114ce57506001600160a01b0384163b155b156114ee5783604051639996b31560e01b81526004016104a79190611764565b50805b9392505050565b80511561150757805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b50805461152c906119fb565b6000825580601f1061153c575050565b601f0160209004906000526020600020908101906110a091905b8082111561156a5760008155600101611556565b5090565b60005b83811015611589578181015183820152602001611571565b50506000910152565b600081518084526115aa81602086016020860161156e565b601f01601f19169290920160200192915050565b60008282518085526020808601955060208260051b8401016020860160005b8481101561160b57601f198684030189526115f9838351611592565b988401989250908301906001016115dd565b5090979650505050505050565b60408152600061162b60408301856115be565b82810360208481019190915284518083528582019282019060005b8181101561160b5784516001600160a01b031683529383019391830191600101611646565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b03811182821017156116a9576116a961166b565b604052919050565b60006001600160401b038311156116ca576116ca61166b565b6116dd601f8401601f1916602001611681565b90508281528383830111156116f157600080fd5b828260208301376000602084830101529392505050565b600082601f83011261171957600080fd5b6114f1838335602085016116b1565b60006020828403121561173a57600080fd5b81356001600160401b0381111561175057600080fd5b61175c84828501611708565b949350505050565b6001600160a01b0391909116815260200190565b60006001600160401b038211156117915761179161166b565b5060051b60200190565b80356001600160a01b03811681146117b257600080fd5b919050565b600082601f8301126117c857600080fd5b813560206117dd6117d883611778565b611681565b8083825260208201915060208460051b8701019350868411156117ff57600080fd5b602086015b84811015611822576118158161179b565b8352918301918301611804565b509695505050505050565b6000806040838503121561184057600080fd5b82356001600160401b038082111561185757600080fd5b818501915085601f83011261186b57600080fd5b8135602061187b6117d883611778565b82815260059290921b8401810191818101908984111561189a57600080fd5b8286015b848110156118d2578035868111156118b65760008081fd5b6118c48c86838b0101611708565b84525091830191830161189e565b50965050860135925050808211156118e957600080fd5b506118f6858286016117b7565b9150509250929050565b6000806040838503121561191357600080fd5b61191c8361179b565b915060208301356001600160401b0381111561193757600080fd5b8301601f8101851361194857600080fd5b6118f6858235602084016116b1565b6000806040838503121561196a57600080fd5b82356001600160401b0381111561198057600080fd5b61198c85828601611708565b92505061199b6020840161179b565b90509250929050565b6020815260006114f16020830184611592565b6000602082840312156119c957600080fd5b6114f18261179b565b6020815260006114f160208301846115be565b634e487b7160e01b600052603260045260246000fd5b600181811c90821680611a0f57607f821691505b602082108103611a2f57634e487b7160e01b600052602260045260246000fd5b50919050565b60008154611a42816119fb565b60018281168015611a5a5760018114611a6f57611a9e565b60ff1984168752821515830287019450611a9e565b8560005260208060002060005b85811015611a955781548a820152908401908201611a7c565b50505082870194505b5050505092915050565b60006114f18284611a35565b60008251611ac681846020870161156e565b9190910192915050565b6020808252602f908201527f4465706c6f796d656e7450726f7879526567697374727956313a204e616d652060408201526e63616e6e6f7420626520656d70747960881b606082015260800190565b60208082526033908201527f4465706c6f796d656e7450726f7879526567697374727956313a20496e76616c604082015272696420636f6e7472616374206164647265737360681b606082015260800190565b6020808252603b90820152600080516020611ef783398151915260408201527a1858dd081b985b5948185b1c9958591e481c9959da5cdd195c9959602a1b606082015260800190565b6020808252603590820152600080516020611ef78339815191526040820152746163742062797465636f646520697320656d70747960581b606082015260800190565b601f821115610798576000816000526020600020601f850160051c81016020861015611c275750805b601f850160051c820191505b81811015610dce57828155600101611c33565b600019600383901b1c191660019190911b1790565b81516001600160401b03811115611c7457611c7461166b565b611c8881611c8284546119fb565b84611bfe565b602080601f831160018114611cb75760008415611ca55750858301515b611caf8582611c46565b865550610dce565b600085815260208120601f198616915b82811015611ce657888601518255948401946001909101908401611cc7565b5085821015611d045787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6020808252603290820152600080516020611ef78339815191526040820152711858dd081b9bdd081c9959da5cdd195c995960721b606082015260800190565b8181038181111561122757634e487b7160e01b600052601160045260246000fd5b818103611d80575050565b611d8a82546119fb565b6001600160401b03811115611da157611da161166b565b611daf81611c8284546119fb565b6000601f821160018114611ddd5760008315611dcb5750848201545b611dd58482611c46565b855550611e42565b600085815260209020601f19841690600086815260209020845b83811015611e175782860154825560019586019590910190602001611df7565b5085831015611e355781850154600019600388901b60f8161c191681555b50505060018360011b0184555b5050505050565b634e487b7160e01b600052603160045260246000fd5b805180151581146117b257600080fd5b600080600060608486031215611e8457600080fd5b611e8d84611e5f565b9250602084015163ffffffff81168114611ea657600080fd5b9150611eb460408501611e5f565b90509250925092565b600060208284031215611ecf57600080fd5b505191905056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc4465706c6f796d656e7450726f7879526567697374727956313a20436f6e7472a2646970667358221220b18a2768de47005a7ed45da07d5002cd8d92a9230f6937b076561d7997dd8af464736f6c63430008180033",
}

// DeploymentProxyRegistryV1 is an auto generated Go binding around an Ethereum contract.
type DeploymentProxyRegistryV1 struct {
	abi abi.ABI
}

// NewDeploymentProxyRegistryV1 creates a new instance of DeploymentProxyRegistryV1.
func NewDeploymentProxyRegistryV1() *DeploymentProxyRegistryV1 {
	parsed, err := DeploymentProxyRegistryV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &DeploymentProxyRegistryV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *DeploymentProxyRegistryV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := deploymentProxyRegistryV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := deploymentProxyRegistryV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
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
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) PackAuthority() []byte {
	enc, err := deploymentProxyRegistryV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := deploymentProxyRegistryV1.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetAllContractNames is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd7a9b816.
//
// Solidity: function getAllContractNames() view returns(string[])
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) PackGetAllContractNames() []byte {
	enc, err := deploymentProxyRegistryV1.abi.Pack("getAllContractNames")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllContractNames is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd7a9b816.
//
// Solidity: function getAllContractNames() view returns(string[])
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackGetAllContractNames(data []byte) ([]string, error) {
	out, err := deploymentProxyRegistryV1.abi.Unpack("getAllContractNames", data)
	if err != nil {
		return *new([]string), err
	}
	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)
	return out0, err
}

// PackGetAllContracts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18d3ce96.
//
// Solidity: function getAllContracts() view returns(string[] names, address[] addresses)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) PackGetAllContracts() []byte {
	enc, err := deploymentProxyRegistryV1.abi.Pack("getAllContracts")
	if err != nil {
		panic(err)
	}
	return enc
}

// GetAllContractsOutput serves as a container for the return parameters of contract
// method GetAllContracts.
type GetAllContractsOutput struct {
	Names     []string
	Addresses []common.Address
}

// UnpackGetAllContracts is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18d3ce96.
//
// Solidity: function getAllContracts() view returns(string[] names, address[] addresses)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackGetAllContracts(data []byte) (GetAllContractsOutput, error) {
	out, err := deploymentProxyRegistryV1.abi.Unpack("getAllContracts", data)
	outstruct := new(GetAllContractsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Names = *abi.ConvertType(out[0], new([]string)).(*[]string)
	outstruct.Addresses = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	return *outstruct, err

}

// PackGetContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x35817773.
//
// Solidity: function getContract(string name) view returns(address)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) PackGetContract(name string) []byte {
	enc, err := deploymentProxyRegistryV1.abi.Pack("getContract", name)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetContract is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x35817773.
//
// Solidity: function getContract(string name) view returns(address)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackGetContract(data []byte) (common.Address, error) {
	out, err := deploymentProxyRegistryV1.abi.Unpack("getContract", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4d66de8.
//
// Solidity: function initialize(address authority_) returns()
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) PackInitialize(authority common.Address) []byte {
	enc, err := deploymentProxyRegistryV1.abi.Pack("initialize", authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) PackProxiableUUID() []byte {
	enc, err := deploymentProxyRegistryV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := deploymentProxyRegistryV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackRegisterContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7f3c2c28.
//
// Solidity: function registerContract(string name, address contractAddress) returns()
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) PackRegisterContract(name string, contractAddress common.Address) []byte {
	enc, err := deploymentProxyRegistryV1.abi.Pack("registerContract", name, contractAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRegisterContracts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3e8552af.
//
// Solidity: function registerContracts(string[] names, address[] contractAddresses) returns()
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) PackRegisterContracts(names []string, contractAddresses []common.Address) []byte {
	enc, err := deploymentProxyRegistryV1.abi.Pack("registerContracts", names, contractAddresses)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRemoveContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x97623b58.
//
// Solidity: function removeContract(string name) returns()
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) PackRemoveContract(name string) []byte {
	enc, err := deploymentProxyRegistryV1.abi.Pack("removeContract", name)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdateContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x697a60b3.
//
// Solidity: function updateContract(string name, address newAddress) returns()
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) PackUpdateContract(name string, newAddress common.Address) []byte {
	enc, err := deploymentProxyRegistryV1.abi.Pack("updateContract", name, newAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := deploymentProxyRegistryV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// DeploymentProxyRegistryV1AuthorityUpdated represents a AuthorityUpdated event raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const DeploymentProxyRegistryV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (DeploymentProxyRegistryV1AuthorityUpdated) ContractEventName() string {
	return DeploymentProxyRegistryV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*DeploymentProxyRegistryV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != deploymentProxyRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DeploymentProxyRegistryV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range deploymentProxyRegistryV1.abi.Events[event].Inputs {
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

// DeploymentProxyRegistryV1ContractRegistered represents a ContractRegistered event raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1ContractRegistered struct {
	Name            common.Hash
	ContractAddress common.Address
	Raw             *types.Log // Blockchain specific contextual infos
}

const DeploymentProxyRegistryV1ContractRegisteredEventName = "ContractRegistered"

// ContractEventName returns the user-defined event name.
func (DeploymentProxyRegistryV1ContractRegistered) ContractEventName() string {
	return DeploymentProxyRegistryV1ContractRegisteredEventName
}

// UnpackContractRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ContractRegistered(string indexed name, address indexed contractAddress)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackContractRegisteredEvent(log *types.Log) (*DeploymentProxyRegistryV1ContractRegistered, error) {
	event := "ContractRegistered"
	if log.Topics[0] != deploymentProxyRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DeploymentProxyRegistryV1ContractRegistered)
	if len(log.Data) > 0 {
		if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range deploymentProxyRegistryV1.abi.Events[event].Inputs {
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

// DeploymentProxyRegistryV1ContractRemoved represents a ContractRemoved event raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1ContractRemoved struct {
	Name            common.Hash
	ContractAddress common.Address
	Raw             *types.Log // Blockchain specific contextual infos
}

const DeploymentProxyRegistryV1ContractRemovedEventName = "ContractRemoved"

// ContractEventName returns the user-defined event name.
func (DeploymentProxyRegistryV1ContractRemoved) ContractEventName() string {
	return DeploymentProxyRegistryV1ContractRemovedEventName
}

// UnpackContractRemovedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ContractRemoved(string indexed name, address indexed contractAddress)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackContractRemovedEvent(log *types.Log) (*DeploymentProxyRegistryV1ContractRemoved, error) {
	event := "ContractRemoved"
	if log.Topics[0] != deploymentProxyRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DeploymentProxyRegistryV1ContractRemoved)
	if len(log.Data) > 0 {
		if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range deploymentProxyRegistryV1.abi.Events[event].Inputs {
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

// DeploymentProxyRegistryV1ContractUpdated represents a ContractUpdated event raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1ContractUpdated struct {
	Name       common.Hash
	OldAddress common.Address
	NewAddress common.Address
	Raw        *types.Log // Blockchain specific contextual infos
}

const DeploymentProxyRegistryV1ContractUpdatedEventName = "ContractUpdated"

// ContractEventName returns the user-defined event name.
func (DeploymentProxyRegistryV1ContractUpdated) ContractEventName() string {
	return DeploymentProxyRegistryV1ContractUpdatedEventName
}

// UnpackContractUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ContractUpdated(string indexed name, address indexed oldAddress, address indexed newAddress)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackContractUpdatedEvent(log *types.Log) (*DeploymentProxyRegistryV1ContractUpdated, error) {
	event := "ContractUpdated"
	if log.Topics[0] != deploymentProxyRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DeploymentProxyRegistryV1ContractUpdated)
	if len(log.Data) > 0 {
		if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range deploymentProxyRegistryV1.abi.Events[event].Inputs {
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

// DeploymentProxyRegistryV1Initialized represents a Initialized event raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const DeploymentProxyRegistryV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (DeploymentProxyRegistryV1Initialized) ContractEventName() string {
	return DeploymentProxyRegistryV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackInitializedEvent(log *types.Log) (*DeploymentProxyRegistryV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != deploymentProxyRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DeploymentProxyRegistryV1Initialized)
	if len(log.Data) > 0 {
		if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range deploymentProxyRegistryV1.abi.Events[event].Inputs {
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

// DeploymentProxyRegistryV1Upgraded represents a Upgraded event raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const DeploymentProxyRegistryV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (DeploymentProxyRegistryV1Upgraded) ContractEventName() string {
	return DeploymentProxyRegistryV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackUpgradedEvent(log *types.Log) (*DeploymentProxyRegistryV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != deploymentProxyRegistryV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DeploymentProxyRegistryV1Upgraded)
	if len(log.Data) > 0 {
		if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range deploymentProxyRegistryV1.abi.Events[event].Inputs {
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
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], deploymentProxyRegistryV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return deploymentProxyRegistryV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], deploymentProxyRegistryV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return deploymentProxyRegistryV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], deploymentProxyRegistryV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return deploymentProxyRegistryV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], deploymentProxyRegistryV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return deploymentProxyRegistryV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], deploymentProxyRegistryV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return deploymentProxyRegistryV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], deploymentProxyRegistryV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return deploymentProxyRegistryV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], deploymentProxyRegistryV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return deploymentProxyRegistryV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], deploymentProxyRegistryV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return deploymentProxyRegistryV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], deploymentProxyRegistryV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return deploymentProxyRegistryV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], deploymentProxyRegistryV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return deploymentProxyRegistryV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], deploymentProxyRegistryV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return deploymentProxyRegistryV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], deploymentProxyRegistryV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return deploymentProxyRegistryV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// DeploymentProxyRegistryV1AddressEmptyCode represents a AddressEmptyCode error raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func DeploymentProxyRegistryV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackAddressEmptyCodeError(raw []byte) (*DeploymentProxyRegistryV1AddressEmptyCode, error) {
	out := new(DeploymentProxyRegistryV1AddressEmptyCode)
	if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DeploymentProxyRegistryV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func DeploymentProxyRegistryV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackERC1967InvalidImplementationError(raw []byte) (*DeploymentProxyRegistryV1ERC1967InvalidImplementation, error) {
	out := new(DeploymentProxyRegistryV1ERC1967InvalidImplementation)
	if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DeploymentProxyRegistryV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func DeploymentProxyRegistryV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackERC1967NonPayableError(raw []byte) (*DeploymentProxyRegistryV1ERC1967NonPayable, error) {
	out := new(DeploymentProxyRegistryV1ERC1967NonPayable)
	if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DeploymentProxyRegistryV1FailedCall represents a FailedCall error raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func DeploymentProxyRegistryV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackFailedCallError(raw []byte) (*DeploymentProxyRegistryV1FailedCall, error) {
	out := new(DeploymentProxyRegistryV1FailedCall)
	if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DeploymentProxyRegistryV1InvalidInitialization represents a InvalidInitialization error raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func DeploymentProxyRegistryV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackInvalidInitializationError(raw []byte) (*DeploymentProxyRegistryV1InvalidInitialization, error) {
	out := new(DeploymentProxyRegistryV1InvalidInitialization)
	if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DeploymentProxyRegistryV1NotInitializing represents a NotInitializing error raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func DeploymentProxyRegistryV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackNotInitializingError(raw []byte) (*DeploymentProxyRegistryV1NotInitializing, error) {
	out := new(DeploymentProxyRegistryV1NotInitializing)
	if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DeploymentProxyRegistryV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func DeploymentProxyRegistryV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*DeploymentProxyRegistryV1RaylsAccessManagedContractPaused, error) {
	out := new(DeploymentProxyRegistryV1RaylsAccessManagedContractPaused)
	if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DeploymentProxyRegistryV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func DeploymentProxyRegistryV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*DeploymentProxyRegistryV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(DeploymentProxyRegistryV1RaylsAccessManagedInvalidAuthority)
	if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DeploymentProxyRegistryV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func DeploymentProxyRegistryV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*DeploymentProxyRegistryV1RaylsAccessManagedMustSchedule, error) {
	out := new(DeploymentProxyRegistryV1RaylsAccessManagedMustSchedule)
	if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DeploymentProxyRegistryV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func DeploymentProxyRegistryV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*DeploymentProxyRegistryV1RaylsAccessManagedUnauthorized, error) {
	out := new(DeploymentProxyRegistryV1RaylsAccessManagedUnauthorized)
	if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DeploymentProxyRegistryV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func DeploymentProxyRegistryV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*DeploymentProxyRegistryV1UUPSUnauthorizedCallContext, error) {
	out := new(DeploymentProxyRegistryV1UUPSUnauthorizedCallContext)
	if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DeploymentProxyRegistryV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the DeploymentProxyRegistryV1 contract.
type DeploymentProxyRegistryV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func DeploymentProxyRegistryV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (deploymentProxyRegistryV1 *DeploymentProxyRegistryV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*DeploymentProxyRegistryV1UUPSUnsupportedProxiableUUID, error) {
	out := new(DeploymentProxyRegistryV1UUPSUnsupportedProxiableUUID)
	if err := deploymentProxyRegistryV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
