// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TokenFreezeManagerV1

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

// TokenStructsFrozenToken is an auto generated low-level Go binding around an user-defined struct.
type TokenStructsFrozenToken struct {
	ResourceId         [32]byte
	FrozenParticipants []*big.Int
}

// TokenFreezeManagerV1MetaData contains all meta data concerning the TokenFreezeManagerV1 contract.
var TokenFreezeManagerV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"broadcastCurrentFrozenResourcesForNewParticipant\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"broadcastFrozenToken\",\"inputs\":[{\"name\":\"frozenToken\",\"type\":\"tuple\",\"internalType\":\"structTokenStructs.FrozenToken\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"frozenParticipants\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"broadcastUnfrozenToken\",\"inputs\":[{\"name\":\"unfrozenToken\",\"type\":\"tuple\",\"internalType\":\"structTokenStructs.FrozenToken\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"frozenParticipants\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"endpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRaylsEndpoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"freezeToken\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"chainIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllFrozenTokens\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structTokenStructs.FrozenToken[]\",\"components\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"frozenParticipants\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_tokenRegistryAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isTokenFrozenForParticipant\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setEndpoint\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTokenRegistryAddress\",\"inputs\":[{\"name\":\"_tokenRegistryAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tokenRegistryAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreezeToken\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"chainIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenFreezeStatusChanged\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"chainIds\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"action\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumTokenStructs.FreezeAction\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenFreezeManagerV1__UnauthorizedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "TokenFreezeManagerV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b50608051611f0361003e60003960008181610f9c01528181610fc501526110f90152611f036000f3fe6080604052600436106100d35760003560e01c8063a615beb21161007a578063a615beb2146101f2578063ad3cb1cc14610214578063bf7e214f14610252578063c0c53b8b14610267578063dbbb415514610287578063e08ac15e146102a7578063e66cdd2b146102c7578063f7bb5fad146102e757600080fd5b8063054e3537146100d85780630566585a146100fa578063338042e41461011a5780634f1ef2861461014f57806352d1902d1461016257806358555ab8146101855780635be2aca0146101a55780635e280f11146101d2575b600080fd5b3480156100e457600080fd5b506100f86100f33660046116f1565b610307565b005b34801561010657600080fd5b506100f86101153660046117cf565b610431565b34801561012657600080fd5b5061013a61013536600461185b565b61049b565b60405190151581526020015b60405180910390f35b6100f861015d366004611899565b610575565b34801561016e57600080fd5b50610177610594565b604051908152602001610146565b34801561019157600080fd5b506100f86101a036600461193e565b6105b1565b3480156101b157600080fd5b506000546101c5906001600160a01b031681565b6040516101469190611959565b3480156101de57600080fd5b506001546101c5906001600160a01b031681565b3480156101fe57600080fd5b50610207610660565b60405161014691906119c6565b34801561022057600080fd5b50610245604051806040016040528060058152602001640352e302e360dc1b81525081565b6040516101469190611a7a565b34801561025e57600080fd5b506101c5610721565b34801561027357600080fd5b506100f8610282366004611a8d565b61073a565b34801561029357600080fd5b506100f86102a236600461193e565b61086e565b3480156102b357600080fd5b506100f86102c2366004611ad0565b610918565b3480156102d357600080fd5b506100f86102e23660046117cf565b610c54565b3480156102f357600080fd5b506100f8610302366004611b4e565b610cbe565b6000546001600160a01b0316331461033d573360405163300b738360e01b81526004016103349190611959565b60405180910390fd5b610345611600565b6001546040516001600160a01b03909116906341d717449084906003906311170d8b60e11b9061037a90600290602401611b94565b60408051601f19818403018152918152602080830180516001600160e01b03166001600160e01b03199586161790528151808201835260008082528351808401855281815284519384018552908352925160e089901b90951685526103e9969594909291908a90600401611c54565b6020604051808303816000875af1158015610408573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061042c9190611d1a565b505050565b6000546001600160a01b0316331461045e573360405163300b738360e01b81526004016103349190611959565b610466611600565b6001546040516001600160a01b03909116906341d7174490600090600390639eb7b4eb60e01b9061037a908890602401611d33565b6000805b6002548110156105695783600282815481106104bd576104bd611d46565b906000526020600020906002020160000154036105615760005b600282815481106104ea576104ea611d46565b90600052602060002090600202016001018054905081101561055f57836002838154811061051a5761051a611d46565b9060005260206000209060020201600101828154811061053c5761053c611d46565b9060005260206000200154036105575760019250505061056f565b6001016104d7565b505b60010161049f565b50600090505b92915050565b61057d610f91565b61058682611021565b610590828261103a565b5050565b600061059e6110ee565b50600080516020611eae83398151915290565b6105c7336000356001600160e01b031916611137565b6001600160a01b03811661063e5760405162461bcd60e51b815260206004820152603860248201527f546f6b656e467265657a654d616e616765723a20546f6b656e526567697374726044820152777920616464726573732063616e6e6f74206265207a65726f60401b6064820152608401610334565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b60606002805480602002602001604051908101604052809291908181526020016000905b828210156107185783829060005260206000209060020201604051806040016040529081600082015481526020016001820180548060200260200160405190810160405280929190818152602001828054801561070057602002820191906000526020600020905b8154815260200190600101908083116106ec575b50505050508152505081526020019060010190610684565b50505050905090565b600061072b611282565b546001600160a01b0316919050565b60006107446112e4565b805490915060ff600160401b82041615906001600160401b031660008115801561076b5750825b90506000826001600160401b031660011480156107875750303b155b905081158015610795575080155b156107b35760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156107dd57845460ff60401b1916600160401b1785555b6107e561130d565b600080546001600160a01b03808a166001600160a01b03199283161790925560018054928b169290911691909117905561081e86611315565b831561086457845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050505050565b610884336000356001600160e01b031916611137565b6001600160a01b0381166108f65760405162461bcd60e51b815260206004820152603360248201527f546f6b656e467265657a654d616e616765723a20456e64706f696e7420616464604482015272726573732063616e6e6f74206265207a65726f60681b6064820152608401610334565b600180546001600160a01b0319166001600160a01b0392909216919091179055565b6000546001600160a01b03163314610945573360405163300b738360e01b81526004016103349190611959565b600080805b60025481101561099357856002828154811061096857610968611d46565b9060005260206000209060020201600001540361098b5780925060019150610993565b60010161094a565b508015610ae5576000600283815481106109af576109af611d46565b9060005260206000209060020201905060005b84811015610a71576000805b6001840154811015610a29578787848181106109ec576109ec611d46565b90506020020135846001018281548110610a0857610a08611d46565b906000526020600020015403610a215760019150610a29565b6001016109ce565b5080610a685782600101878784818110610a4557610a45611d46565b835460018101855560009485526020948590209190940292909201359190920155505b506001016109c2565b50604080518082018252825481526001830180548351602082810282018101909552818152610adf948693818601939091830182828015610ad157602002820191906000526020600020905b815481526020019060010190808311610abd575b505050505081525050610431565b50610c10565b6040805180820190915260606020820152858152836001600160401b03811115610b1157610b1161170a565b604051908082528060200260200182016040528015610b3a578160200160208202803683370190505b50602082015260005b84811015610b8d57858582818110610b5d57610b5d611d46565b9050602002013582602001518281518110610b7a57610b7a611d46565b6020908102919091010152600101610b43565b50600280546001810182556000829052825191027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace81019182556020808401518051859493610c02937f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5acf909101920190611637565b505050610c0e81610431565b505b847fd20267c111a6f5fe591300c4d39881c83d67fdeefca9c994b9e1939064ff758085856001604051610c4593929190611d70565b60405180910390a25050505050565b6000546001600160a01b03163314610c81573360405163300b738360e01b81526004016103349190611959565b610c89611600565b6001546040516001600160a01b03909116906341d7174490600090600390633d6bad0d60e21b9061037a908890602401611d33565b6000546001600160a01b03163314610ceb573360405163300b738360e01b81526004016103349190611959565b600080805b600254811015610d39578460028281548110610d0e57610d0e611d46565b90600052602060002090600202016000015403610d315780925060019150610d39565b600101610cf0565b5080610d795760405162461bcd60e51b815260206004820152600f60248201526e151bdad95b881b9bdd08199bdd5b99608a1b6044820152606401610334565b600060028381548110610d8e57610d8e611d46565b600091825260208220600290910201915060018201905b8551811015610e7d5760005b8254811015610e7457868281518110610dcc57610dcc611d46565b6020026020010151838281548110610de657610de6611d46565b906000526020600020015403610e6c5782548390610e0690600190611db7565b81548110610e1657610e16611d46565b9060005260206000200154838281548110610e3357610e33611d46565b906000526020600020018190555082805480610e5157610e51611dd8565b60019003818190600052602060002001600090559055610e74565b600101610db1565b50600101610da5565b508054600003610f315760028054610e9790600190611db7565b81548110610ea757610ea7611d46565b906000526020600020906002020160028581548110610ec857610ec8611d46565b60009182526020909120825460029092020190815560018083018054610ef19284019190611682565b509050506002805480610f0657610f06611dd8565b60008281526020812060026000199093019283020181815590610f2c60018301826116c2565b505090555b610f4e604051806040016040528088815260200187815250610c54565b857fd20267c111a6f5fe591300c4d39881c83d67fdeefca9c994b9e1939064ff7580866000604051610f81929190611dee565b60405180910390a2505050505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016148061100157507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610ff5611356565b6001600160a01b031614155b1561101f5760405163703e46dd60e11b815260040160405180910390fd5b565b611037336000356001600160e01b031916611137565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015611094575060408051601f3d908101601f1916820190925261109191810190611d1a565b60015b6110b35781604051634c9c8ce360e01b81526004016103349190611959565b600080516020611eae83398151915281146110e457604051632a87526960e21b815260048101829052602401610334565b61042c838361136c565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461101f5760405163703e46dd60e11b815260040160405180910390fd5b6000611141611282565b80549091506001600160a01b031680611170576000604051638944034760e01b81526004016103349190611959565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa1580156111d4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111f89190611e4c565b925092509250826112795780156112225760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff82161561125e5760405163a426878960e01b81526001600160a01b038816600482015263ffffffff83166024820152604401610334565b86604051632ecd3d0360e21b81526004016103349190611959565b50505050505050565b60008060ff196112b360017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35611db7565b6040516020016112c591815260200190565b60408051601f1981840301815291905280516020909101201692915050565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061056f565b61101f6113c2565b600061131f611282565b80549091506001600160a01b03161561134d5781604051638944034760e01b81526004016103349190611959565b610590826113e7565b6000600080516020611eae83398151915261072b565b61137582611477565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a28051156113ba5761042c82826114d3565b610590611549565b6113ca611568565b61101f57604051631afcd79f60e31b815260040160405180910390fd5b6001600160a01b0381166114105780604051638944034760e01b81526004016103349190611959565b600061141a611282565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b806001600160a01b03163b6000036114a45780604051634c9c8ce360e01b81526004016103349190611959565b600080516020611eae83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b0316846040516114f09190611e91565b600060405180830381855af49150503d806000811461152b576040519150601f19603f3d011682016040523d82523d6000602084013e611530565b606091505b5091509150611540858383611582565b95945050505050565b341561101f5760405163b398979f60e01b815260040160405180910390fd5b60006115726112e4565b54600160401b900460ff16919050565b60608261159757611592826115d8565b6115d1565b81511580156115ae57506001600160a01b0384163b155b156115ce5783604051639996b31560e01b81526004016103349190611959565b50805b9392505050565b8051156115e757805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b6040805160c08101909152806000815260006020820181905260408201819052606082018190526080820181905260a09091015290565b828054828255906000526020600020908101928215611672579160200282015b82811115611672578251825591602001919060010190611657565b5061167e9291506116dc565b5090565b8280548282559060005260206000209081019282156116725760005260206000209182015b828111156116725782548255916001019190600101906116a7565b508054600082559060005260206000209081019061103791905b5b8082111561167e57600081556001016116dd565b60006020828403121561170357600080fd5b5035919050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b03811182821017156117485761174861170a565b604052919050565b600082601f83011261176157600080fd5b813560206001600160401b0382111561177c5761177c61170a565b8160051b61178b828201611720565b92835284810182019282810190878511156117a557600080fd5b83870192505b848310156117c4578235825291830191908301906117ab565b979650505050505050565b6000602082840312156117e157600080fd5b81356001600160401b03808211156117f857600080fd5b908301906040828603121561180c57600080fd5b6040516040810181811083821117156118275761182761170a565b6040528235815260208301358281111561184057600080fd5b61184c87828601611750565b60208301525095945050505050565b6000806040838503121561186e57600080fd5b50508035926020909101359150565b80356001600160a01b038116811461189457600080fd5b919050565b600080604083850312156118ac57600080fd5b6118b58361187d565b91506020808401356001600160401b03808211156118d257600080fd5b818601915086601f8301126118e657600080fd5b8135818111156118f8576118f861170a565b61190a601f8201601f19168501611720565b9150808252878482850101111561192057600080fd5b80848401858401376000848284010152508093505050509250929050565b60006020828403121561195057600080fd5b6115d18261187d565b6001600160a01b0391909116815260200190565b6000604083018251845260208084015160406020870152828151808552606088019150602083019450600092505b808310156119bb578451825293830193600192909201919083019061199b565b509695505050505050565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015611a1d57603f19888603018452611a0b85835161196d565b945092850192908501906001016119ef565b5092979650505050505050565b60005b83811015611a45578181015183820152602001611a2d565b50506000910152565b60008151808452611a66816020860160208601611a2a565b601f01601f19169290920160200192915050565b6020815260006115d16020830184611a4e565b600080600060608486031215611aa257600080fd5b611aab8461187d565b9250611ab96020850161187d565b9150611ac76040850161187d565b90509250925092565b600080600060408486031215611ae557600080fd5b8335925060208401356001600160401b0380821115611b0357600080fd5b818601915086601f830112611b1757600080fd5b813581811115611b2657600080fd5b8760208260051b8501011115611b3b57600080fd5b6020830194508093505050509250925092565b60008060408385031215611b6157600080fd5b8235915060208301356001600160401b03811115611b7e57600080fd5b611b8a85828601611750565b9150509250929050565b600060208083018184528085548083526040925060408601915060408160051b8701016000888152858120815b84811015611c2f57898403603f1901865281548452878401879052600180830180548987018190529085528985209190859060608801905b80831015611c1657845482529383019391830191908c0190611bf9565b50988b0198965050506002929092019150600101611bc1565b50919998505050505050505050565b634e487b7160e01b600052602160045260246000fd5b6000610180898352886020840152806040840152611c7481840189611a4e565b90508281036060840152611c888188611a4e565b90508281036080840152611c9c8187611a4e565b905082810360a0840152611cb08186611a4e565b9150508251600d8110611cc557611cc5611c3e565b60c0830152602083015160e083015260408301516001600160a01b039081166101008401526060840151811661012084015260808401511661014083015260a090920151610160909101529695505050505050565b600060208284031215611d2c57600080fd5b5051919050565b6020815260006115d1602083018461196d565b634e487b7160e01b600052603260045260246000fd5b60028110611d6c57611d6c611c3e565b9052565b6040808252810183905260006001600160fb1b03841115611d9057600080fd5b8360051b8086606085013782016060019050611daf6020830184611d5c565b949350505050565b8181038181111561056f57634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b604080825283519082018190526000906020906060840190828701845b82811015611e2757815184529284019290840190600101611e0b565b50505080925050506115d16020830184611d5c565b8051801515811461189457600080fd5b600080600060608486031215611e6157600080fd5b611e6a84611e3c565b9250602084015163ffffffff81168114611e8357600080fd5b9150611ac760408501611e3c565b60008251611ea3818460208701611a2a565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca26469706673582212204f0f8f8f7b44394bb63fbfd438bbf234aa6ced1b8d51aa22ce20a29e35d6857d64736f6c63430008180033",
}

// TokenFreezeManagerV1 is an auto generated Go binding around an Ethereum contract.
type TokenFreezeManagerV1 struct {
	abi abi.ABI
}

// NewTokenFreezeManagerV1 creates a new instance of TokenFreezeManagerV1.
func NewTokenFreezeManagerV1() *TokenFreezeManagerV1 {
	parsed, err := TokenFreezeManagerV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &TokenFreezeManagerV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *TokenFreezeManagerV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := tokenFreezeManagerV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
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
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackAuthority() []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := tokenFreezeManagerV1.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackBroadcastCurrentFrozenResourcesForNewParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x054e3537.
//
// Solidity: function broadcastCurrentFrozenResourcesForNewParticipant(uint256 chainId) returns()
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackBroadcastCurrentFrozenResourcesForNewParticipant(chainId *big.Int) []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("broadcastCurrentFrozenResourcesForNewParticipant", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackBroadcastFrozenToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0566585a.
//
// Solidity: function broadcastFrozenToken((bytes32,uint256[]) frozenToken) returns()
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackBroadcastFrozenToken(frozenToken TokenStructsFrozenToken) []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("broadcastFrozenToken", frozenToken)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackBroadcastUnfrozenToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe66cdd2b.
//
// Solidity: function broadcastUnfrozenToken((bytes32,uint256[]) unfrozenToken) returns()
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackBroadcastUnfrozenToken(unfrozenToken TokenStructsFrozenToken) []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("broadcastUnfrozenToken", unfrozenToken)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackEndpoint() []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("endpoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEndpoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackEndpoint(data []byte) (common.Address, error) {
	out, err := tokenFreezeManagerV1.abi.Unpack("endpoint", data)
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
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackFreezeToken(resourceId [32]byte, chainIds []*big.Int) []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("freezeToken", resourceId, chainIds)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetAllFrozenTokens is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa615beb2.
//
// Solidity: function getAllFrozenTokens() view returns((bytes32,uint256[])[])
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackGetAllFrozenTokens() []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("getAllFrozenTokens")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAllFrozenTokens is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa615beb2.
//
// Solidity: function getAllFrozenTokens() view returns((bytes32,uint256[])[])
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackGetAllFrozenTokens(data []byte) ([]TokenStructsFrozenToken, error) {
	out, err := tokenFreezeManagerV1.abi.Unpack("getAllFrozenTokens", data)
	if err != nil {
		return *new([]TokenStructsFrozenToken), err
	}
	out0 := *abi.ConvertType(out[0], new([]TokenStructsFrozenToken)).(*[]TokenStructsFrozenToken)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc0c53b8b.
//
// Solidity: function initialize(address _endpoint, address _tokenRegistryAddress, address authority_) returns()
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackInitialize(endpoint common.Address, tokenRegistryAddress common.Address, authority common.Address) []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("initialize", endpoint, tokenRegistryAddress, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackIsTokenFrozenForParticipant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x338042e4.
//
// Solidity: function isTokenFrozenForParticipant(bytes32 resourceId, uint256 chainId) view returns(bool)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackIsTokenFrozenForParticipant(resourceId [32]byte, chainId *big.Int) []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("isTokenFrozenForParticipant", resourceId, chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsTokenFrozenForParticipant is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x338042e4.
//
// Solidity: function isTokenFrozenForParticipant(bytes32 resourceId, uint256 chainId) view returns(bool)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackIsTokenFrozenForParticipant(data []byte) (bool, error) {
	out, err := tokenFreezeManagerV1.abi.Unpack("isTokenFrozenForParticipant", data)
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
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackProxiableUUID() []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := tokenFreezeManagerV1.abi.Unpack("proxiableUUID", data)
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
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackSetEndpoint(endpoint common.Address) []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("setEndpoint", endpoint)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetTokenRegistryAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x58555ab8.
//
// Solidity: function setTokenRegistryAddress(address _tokenRegistryAddress) returns()
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackSetTokenRegistryAddress(tokenRegistryAddress common.Address) []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("setTokenRegistryAddress", tokenRegistryAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackTokenRegistryAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5be2aca0.
//
// Solidity: function tokenRegistryAddress() view returns(address)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackTokenRegistryAddress() []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("tokenRegistryAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenRegistryAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5be2aca0.
//
// Solidity: function tokenRegistryAddress() view returns(address)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackTokenRegistryAddress(data []byte) (common.Address, error) {
	out, err := tokenFreezeManagerV1.abi.Unpack("tokenRegistryAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackUnfreezeToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf7bb5fad.
//
// Solidity: function unfreezeToken(bytes32 resourceId, uint256[] chainIds) returns()
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackUnfreezeToken(resourceId [32]byte, chainIds []*big.Int) []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("unfreezeToken", resourceId, chainIds)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := tokenFreezeManagerV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TokenFreezeManagerV1AuthorityUpdated represents a AuthorityUpdated event raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const TokenFreezeManagerV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (TokenFreezeManagerV1AuthorityUpdated) ContractEventName() string {
	return TokenFreezeManagerV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*TokenFreezeManagerV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != tokenFreezeManagerV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TokenFreezeManagerV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tokenFreezeManagerV1.abi.Events[event].Inputs {
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

// TokenFreezeManagerV1Initialized represents a Initialized event raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const TokenFreezeManagerV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (TokenFreezeManagerV1Initialized) ContractEventName() string {
	return TokenFreezeManagerV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackInitializedEvent(log *types.Log) (*TokenFreezeManagerV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != tokenFreezeManagerV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TokenFreezeManagerV1Initialized)
	if len(log.Data) > 0 {
		if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tokenFreezeManagerV1.abi.Events[event].Inputs {
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

// TokenFreezeManagerV1TokenFreezeStatusChanged represents a TokenFreezeStatusChanged event raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1TokenFreezeStatusChanged struct {
	ResourceId [32]byte
	ChainIds   []*big.Int
	Action     uint8
	Raw        *types.Log // Blockchain specific contextual infos
}

const TokenFreezeManagerV1TokenFreezeStatusChangedEventName = "TokenFreezeStatusChanged"

// ContractEventName returns the user-defined event name.
func (TokenFreezeManagerV1TokenFreezeStatusChanged) ContractEventName() string {
	return TokenFreezeManagerV1TokenFreezeStatusChangedEventName
}

// UnpackTokenFreezeStatusChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenFreezeStatusChanged(bytes32 indexed resourceId, uint256[] chainIds, uint8 action)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackTokenFreezeStatusChangedEvent(log *types.Log) (*TokenFreezeManagerV1TokenFreezeStatusChanged, error) {
	event := "TokenFreezeStatusChanged"
	if log.Topics[0] != tokenFreezeManagerV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TokenFreezeManagerV1TokenFreezeStatusChanged)
	if len(log.Data) > 0 {
		if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tokenFreezeManagerV1.abi.Events[event].Inputs {
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

// TokenFreezeManagerV1Upgraded represents a Upgraded event raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const TokenFreezeManagerV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (TokenFreezeManagerV1Upgraded) ContractEventName() string {
	return TokenFreezeManagerV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackUpgradedEvent(log *types.Log) (*TokenFreezeManagerV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != tokenFreezeManagerV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TokenFreezeManagerV1Upgraded)
	if len(log.Data) > 0 {
		if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tokenFreezeManagerV1.abi.Events[event].Inputs {
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
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], tokenFreezeManagerV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return tokenFreezeManagerV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenFreezeManagerV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return tokenFreezeManagerV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenFreezeManagerV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return tokenFreezeManagerV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenFreezeManagerV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return tokenFreezeManagerV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenFreezeManagerV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return tokenFreezeManagerV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenFreezeManagerV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return tokenFreezeManagerV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenFreezeManagerV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return tokenFreezeManagerV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenFreezeManagerV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return tokenFreezeManagerV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenFreezeManagerV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return tokenFreezeManagerV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenFreezeManagerV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return tokenFreezeManagerV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenFreezeManagerV1.abi.Errors["TokenFreezeManagerV1UnauthorizedCaller"].ID.Bytes()[:4]) {
		return tokenFreezeManagerV1.UnpackTokenFreezeManagerV1UnauthorizedCallerError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenFreezeManagerV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return tokenFreezeManagerV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], tokenFreezeManagerV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return tokenFreezeManagerV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// TokenFreezeManagerV1AddressEmptyCode represents a AddressEmptyCode error raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func TokenFreezeManagerV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackAddressEmptyCodeError(raw []byte) (*TokenFreezeManagerV1AddressEmptyCode, error) {
	out := new(TokenFreezeManagerV1AddressEmptyCode)
	if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenFreezeManagerV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func TokenFreezeManagerV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackERC1967InvalidImplementationError(raw []byte) (*TokenFreezeManagerV1ERC1967InvalidImplementation, error) {
	out := new(TokenFreezeManagerV1ERC1967InvalidImplementation)
	if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenFreezeManagerV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func TokenFreezeManagerV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackERC1967NonPayableError(raw []byte) (*TokenFreezeManagerV1ERC1967NonPayable, error) {
	out := new(TokenFreezeManagerV1ERC1967NonPayable)
	if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenFreezeManagerV1FailedCall represents a FailedCall error raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func TokenFreezeManagerV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackFailedCallError(raw []byte) (*TokenFreezeManagerV1FailedCall, error) {
	out := new(TokenFreezeManagerV1FailedCall)
	if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenFreezeManagerV1InvalidInitialization represents a InvalidInitialization error raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func TokenFreezeManagerV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackInvalidInitializationError(raw []byte) (*TokenFreezeManagerV1InvalidInitialization, error) {
	out := new(TokenFreezeManagerV1InvalidInitialization)
	if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenFreezeManagerV1NotInitializing represents a NotInitializing error raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func TokenFreezeManagerV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackNotInitializingError(raw []byte) (*TokenFreezeManagerV1NotInitializing, error) {
	out := new(TokenFreezeManagerV1NotInitializing)
	if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenFreezeManagerV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func TokenFreezeManagerV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*TokenFreezeManagerV1RaylsAccessManagedContractPaused, error) {
	out := new(TokenFreezeManagerV1RaylsAccessManagedContractPaused)
	if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenFreezeManagerV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func TokenFreezeManagerV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*TokenFreezeManagerV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(TokenFreezeManagerV1RaylsAccessManagedInvalidAuthority)
	if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenFreezeManagerV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func TokenFreezeManagerV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*TokenFreezeManagerV1RaylsAccessManagedMustSchedule, error) {
	out := new(TokenFreezeManagerV1RaylsAccessManagedMustSchedule)
	if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenFreezeManagerV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func TokenFreezeManagerV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*TokenFreezeManagerV1RaylsAccessManagedUnauthorized, error) {
	out := new(TokenFreezeManagerV1RaylsAccessManagedUnauthorized)
	if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenFreezeManagerV1TokenFreezeManagerV1UnauthorizedCaller represents a TokenFreezeManagerV1__UnauthorizedCaller error raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1TokenFreezeManagerV1UnauthorizedCaller struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TokenFreezeManagerV1__UnauthorizedCaller(address caller)
func TokenFreezeManagerV1TokenFreezeManagerV1UnauthorizedCallerErrorID() common.Hash {
	return common.HexToHash("0x300b7383dbbb11ac05102a85adec0be8ce7f46654deafffe1c39abe93ff25dc5")
}

// UnpackTokenFreezeManagerV1UnauthorizedCallerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TokenFreezeManagerV1__UnauthorizedCaller(address caller)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackTokenFreezeManagerV1UnauthorizedCallerError(raw []byte) (*TokenFreezeManagerV1TokenFreezeManagerV1UnauthorizedCaller, error) {
	out := new(TokenFreezeManagerV1TokenFreezeManagerV1UnauthorizedCaller)
	if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, "TokenFreezeManagerV1UnauthorizedCaller", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenFreezeManagerV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func TokenFreezeManagerV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*TokenFreezeManagerV1UUPSUnauthorizedCallContext, error) {
	out := new(TokenFreezeManagerV1UUPSUnauthorizedCallContext)
	if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// TokenFreezeManagerV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the TokenFreezeManagerV1 contract.
type TokenFreezeManagerV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func TokenFreezeManagerV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (tokenFreezeManagerV1 *TokenFreezeManagerV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*TokenFreezeManagerV1UUPSUnsupportedProxiableUUID, error) {
	out := new(TokenFreezeManagerV1UUPSUnsupportedProxiableUUID)
	if err := tokenFreezeManagerV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
