// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keystone_capability_registry

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

type CapabilityRegistryCapability struct {
	CapabilityType        [32]byte
	Version               [32]byte
	ResponseType          uint8
	ConfigurationContract common.Address
}

type CapabilityRegistryNodeOperator struct {
	Admin common.Address
	Name  string
}

var CapabilityRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityAlreadyDeprecated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CapabilityAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedConfigurationContract\",\"type\":\"address\"}],\"name\":\"InvalidCapabilityConfigurationContractInterface\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeOperatorAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lengthOne\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lengthTwo\",\"type\":\"uint256\"}],\"name\":\"LengthMismatch\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDeprecated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"}],\"name\":\"NodeOperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityType\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability\",\"name\":\"capability\",\"type\":\"tuple\"}],\"name\":\"addCapability\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"addNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"deprecateCapability\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCapabilities\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityType\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityID\",\"type\":\"bytes32\"}],\"name\":\"getCapability\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityType\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityType\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"}],\"name\":\"getCapabilityID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"}],\"name\":\"getNodeOperator\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"isCapabilityDeprecated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint256[]\"}],\"name\":\"removeNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"updateNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b611a7d80620001586000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c806365c14dc71161008c5780639cb7c5f4116100665780639cb7c5f41461023c578063ae3c241c1461025c578063ddbe4f821461026f578063f2fde38b1461028457600080fd5b806365c14dc7146101ec57806379ba50971461020c5780638da5cb5b1461021457600080fd5b8063181f5a77116100c8578063181f5a771461013f5780631cdf63431461017e578063229111f514610191578063398f3773146101d957600080fd5b80630c5801e3146100ef578063117392ce146101045780631257001114610117575b600080fd5b6101026100fd36600461120e565b610297565b005b61010261011236600461127a565b6105a8565b61012a610125366004611292565b6107f3565b60405190151581526020015b60405180910390f35b604080518082018252601881527f4361706162696c697479526567697374727920312e302e30000000000000000060208201529051610136919061130f565b61010261018c366004611322565b610806565b6101cb61019f366004611364565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b604051908152602001610136565b6101026101e7366004611322565b6108c9565b6101ff6101fa366004611292565b610a62565b6040516101369190611386565b610102610b48565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610136565b61024f61024a366004611292565b610c45565b6040516101369190611468565b61010261026a366004611292565b610cef565b610277610dba565b6040516101369190611476565b6101026102923660046114e6565b610eff565b8281146102df576040517fab8b67c600000000000000000000000000000000000000000000000000000000815260048101849052602481018290526044015b60405180910390fd5b6000805473ffffffffffffffffffffffffffffffffffffffff16905b848110156105a057600086868381811061031757610317611503565b905060200201359050600085858481811061033457610334611503565b90506020028101906103469190611532565b61034f90611617565b805190915073ffffffffffffffffffffffffffffffffffffffff166103a0576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805173ffffffffffffffffffffffffffffffffffffffff1633148015906103dd57503373ffffffffffffffffffffffffffffffffffffffff851614155b15610414576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160008381526007602052604090205473ffffffffffffffffffffffffffffffffffffffff90811691161415806104c6575060208082015160405161045a920161130f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201206000868152600783529290922091926104ad926001019101611730565b6040516020818303038152906040528051906020012014155b1561058d578051600083815260076020908152604090912080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909316929092178255820151600190910190610533908261181f565b50806000015173ffffffffffffffffffffffffffffffffffffffff167f14c8f513e8a6d86d2d16b0cb64976de4e72386c4f8068eca3b7354373f8fe97a838360200151604051610584929190611939565b60405180910390a25b50508061059990611981565b90506102fb565b505050505050565b6105b0610f13565b6040805182356020828101919091528084013582840152825180830384018152606090920190925280519101206105e8600382610f96565b1561061f576040517fe288638f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061063160808401606085016114e6565b73ffffffffffffffffffffffffffffffffffffffff161461079c5761065c60808301606084016114e6565b73ffffffffffffffffffffffffffffffffffffffff163b158061073c575061068a60808301606084016114e6565b6040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527f884efe6100000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff91909116906301ffc9a790602401602060405180830381865afa158015610716573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061073a91906119b9565b155b1561079c5761075160808301606084016114e6565b6040517fabb5e3fd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911660048201526024016102d6565b6107a7600382610fb1565b50600081815260026020526040902082906107c282826119db565b505060405181907f65610e5677eedff94555572640e442f89848a109ef8593fa927ac30b2565ff0690600090a25050565b6000610800600583610f96565b92915050565b61080e610f13565b60005b818110156108c457600083838381811061082d5761082d611503565b60209081029290920135600081815260079093526040832080547fffffffffffffffffffffffff000000000000000000000000000000000000000016815590935091905061087e6001830182611174565b50506040518181527f1e5877d7b3001d1569bf733b76c7eceda58bd6c031e5b8d0b7042308ba2e9d4f9060200160405180910390a1506108bd81611981565b9050610811565b505050565b6108d1610f13565b60005b818110156108c45760008383838181106108f0576108f0611503565b90506020028101906109029190611532565b61090b90611617565b805190915073ffffffffffffffffffffffffffffffffffffffff1661095c576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600854604080518082018252835173ffffffffffffffffffffffffffffffffffffffff908116825260208086015181840190815260008681526007909252939020825181547fffffffffffffffffffffffff0000000000000000000000000000000000000000169216919091178155915190919060018201906109df908261181f565b509050506008600081546109f290611981565b909155508151602083015160405173ffffffffffffffffffffffffffffffffffffffff909216917fda6697b182650034bd205cdc2dbfabb06bdb3a0a83a2b45bfefa3c4881284e0b91610a4791859190611939565b60405180910390a2505080610a5b90611981565b90506108d4565b6040805180820190915260008152606060208201526000828152600760209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff1683526001810180549192840191610abf906116e3565b80601f0160208091040260200160405190810160405280929190818152602001828054610aeb906116e3565b8015610b385780601f10610b0d57610100808354040283529160200191610b38565b820191906000526020600020905b815481529060010190602001808311610b1b57829003601f168201915b5050505050815250509050919050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610bc9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016102d6565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b604080516080808201835260008083526020808401829052838501829052606084018290528582526002808252918590208551938401865280548452600180820154928501929092529182015493949293919284019160ff1690811115610cae57610cae6113c9565b6001811115610cbf57610cbf6113c9565b815260029190910154610100900473ffffffffffffffffffffffffffffffffffffffff1660209091015292915050565b610cf7610f13565b610d02600382610f96565b610d3b576040517fe181733f000000000000000000000000000000000000000000000000000000008152600481018290526024016102d6565b610d46600582610f96565b15610d80576040517f16950d1d000000000000000000000000000000000000000000000000000000008152600481018290526024016102d6565b610d8b600582610fb1565b5060405181907fdcea1b78b6ddc31592a94607d537543fcaafda6cc52d6d5cc7bbfca1422baf2190600090a250565b60606000610dc86003610fbd565b90506000610dd66005610fca565b8251610de29190611a5d565b67ffffffffffffffff811115610dfa57610dfa611570565b604051908082528060200260200182016040528015610e6a57816020015b6040805160808101825260008082526020808301829052928201819052606082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181610e185790505b5090506000805b8351811015610ef6576000848281518110610e8e57610e8e611503565b60200260200101519050610eac816005610f9690919063ffffffff16565b610ee557610eb981610c45565b848481518110610ecb57610ecb611503565b60200260200101819052508280610ee190611981565b9350505b50610eef81611981565b9050610e71565b50909392505050565b610f07610f13565b610f1081610fd4565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610f94576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102d6565b565b600081815260018301602052604081205415155b9392505050565b6000610faa83836110c9565b60606000610faa83611118565b6000610800825490565b3373ffffffffffffffffffffffffffffffffffffffff821603611053576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102d6565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600081815260018301602052604081205461111057508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610800565b506000610800565b60608160000180548060200260200160405190810160405280929190818152602001828054801561116857602002820191906000526020600020905b815481526020019060010190808311611154575b50505050509050919050565b508054611180906116e3565b6000825580601f10611190575050565b601f016020900490600052602060002090810190610f1091905b808211156111be57600081556001016111aa565b5090565b60008083601f8401126111d457600080fd5b50813567ffffffffffffffff8111156111ec57600080fd5b6020830191508360208260051b850101111561120757600080fd5b9250929050565b6000806000806040858703121561122457600080fd5b843567ffffffffffffffff8082111561123c57600080fd5b611248888389016111c2565b9096509450602087013591508082111561126157600080fd5b5061126e878288016111c2565b95989497509550505050565b60006080828403121561128c57600080fd5b50919050565b6000602082840312156112a457600080fd5b5035919050565b6000815180845260005b818110156112d1576020818501810151868301820152016112b5565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000610faa60208301846112ab565b6000806020838503121561133557600080fd5b823567ffffffffffffffff81111561134c57600080fd5b611358858286016111c2565b90969095509350505050565b6000806040838503121561137757600080fd5b50508035926020909101359150565b6020815273ffffffffffffffffffffffffffffffffffffffff8251166020820152600060208301516040808401526113c160608401826112ab565b949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8051825260208101516020830152604081015160028110611442577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b604083015260609081015173ffffffffffffffffffffffffffffffffffffffff16910152565b6080810161080082846113f8565b6020808252825182820181905260009190848201906040850190845b818110156114b8576114a58385516113f8565b9284019260809290920191600101611492565b50909695505050505050565b73ffffffffffffffffffffffffffffffffffffffff81168114610f1057600080fd5b6000602082840312156114f857600080fd5b8135610faa816114c4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc183360301811261156657600080fd5b9190910192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff811182821017156115c2576115c2611570565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561160f5761160f611570565b604052919050565b60006040823603121561162957600080fd5b61163161159f565b823561163c816114c4565b815260208381013567ffffffffffffffff8082111561165a57600080fd5b9085019036601f83011261166d57600080fd5b81358181111561167f5761167f611570565b6116af847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016115c8565b915080825236848285010111156116c557600080fd5b80848401858401376000908201840152918301919091525092915050565b600181811c908216806116f757607f821691505b60208210810361128c577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000602080835260008454611744816116e3565b80848701526040600180841660008114611765576001811461179d576117cb565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838a01528284151560051b8a010195506117cb565b896000528660002060005b858110156117c35781548b82018601529083019088016117a8565b8a0184019650505b509398975050505050505050565b601f8211156108c457600081815260208120601f850160051c810160208610156118005750805b601f850160051c820191505b818110156105a05782815560010161180c565b815167ffffffffffffffff81111561183957611839611570565b61184d8161184784546116e3565b846117d9565b602080601f8311600181146118a0576000841561186a5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556105a0565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156118ed578886015182559484019460019091019084016118ce565b508582101561192957878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b8281526040602082015260006113c160408301846112ab565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036119b2576119b2611952565b5060010190565b6000602082840312156119cb57600080fd5b81518015158114610faa57600080fd5b8135815560208201356001820155600281016040830135600281106119ff57600080fd5b81546060850135611a0f816114c4565b74ffffffffffffffffffffffffffffffffffffffff008160081b1660ff84167fffffffffffffffffffffff000000000000000000000000000000000000000000841617178455505050505050565b818103818111156108005761080061195256fea164736f6c6343000813000a",
}

var CapabilityRegistryABI = CapabilityRegistryMetaData.ABI

var CapabilityRegistryBin = CapabilityRegistryMetaData.Bin

func DeployCapabilityRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CapabilityRegistry, error) {
	parsed, err := CapabilityRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CapabilityRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CapabilityRegistry{address: address, abi: *parsed, CapabilityRegistryCaller: CapabilityRegistryCaller{contract: contract}, CapabilityRegistryTransactor: CapabilityRegistryTransactor{contract: contract}, CapabilityRegistryFilterer: CapabilityRegistryFilterer{contract: contract}}, nil
}

type CapabilityRegistry struct {
	address common.Address
	abi     abi.ABI
	CapabilityRegistryCaller
	CapabilityRegistryTransactor
	CapabilityRegistryFilterer
}

type CapabilityRegistryCaller struct {
	contract *bind.BoundContract
}

type CapabilityRegistryTransactor struct {
	contract *bind.BoundContract
}

type CapabilityRegistryFilterer struct {
	contract *bind.BoundContract
}

type CapabilityRegistrySession struct {
	Contract     *CapabilityRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CapabilityRegistryCallerSession struct {
	Contract *CapabilityRegistryCaller
	CallOpts bind.CallOpts
}

type CapabilityRegistryTransactorSession struct {
	Contract     *CapabilityRegistryTransactor
	TransactOpts bind.TransactOpts
}

type CapabilityRegistryRaw struct {
	Contract *CapabilityRegistry
}

type CapabilityRegistryCallerRaw struct {
	Contract *CapabilityRegistryCaller
}

type CapabilityRegistryTransactorRaw struct {
	Contract *CapabilityRegistryTransactor
}

func NewCapabilityRegistry(address common.Address, backend bind.ContractBackend) (*CapabilityRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(CapabilityRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCapabilityRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistry{address: address, abi: abi, CapabilityRegistryCaller: CapabilityRegistryCaller{contract: contract}, CapabilityRegistryTransactor: CapabilityRegistryTransactor{contract: contract}, CapabilityRegistryFilterer: CapabilityRegistryFilterer{contract: contract}}, nil
}

func NewCapabilityRegistryCaller(address common.Address, caller bind.ContractCaller) (*CapabilityRegistryCaller, error) {
	contract, err := bindCapabilityRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryCaller{contract: contract}, nil
}

func NewCapabilityRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*CapabilityRegistryTransactor, error) {
	contract, err := bindCapabilityRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryTransactor{contract: contract}, nil
}

func NewCapabilityRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*CapabilityRegistryFilterer, error) {
	contract, err := bindCapabilityRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryFilterer{contract: contract}, nil
}

func bindCapabilityRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CapabilityRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_CapabilityRegistry *CapabilityRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CapabilityRegistry.Contract.CapabilityRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_CapabilityRegistry *CapabilityRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.CapabilityRegistryTransactor.contract.Transfer(opts)
}

func (_CapabilityRegistry *CapabilityRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.CapabilityRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_CapabilityRegistry *CapabilityRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CapabilityRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.contract.Transfer(opts)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) GetCapabilities(opts *bind.CallOpts) ([]CapabilityRegistryCapability, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getCapabilities")

	if err != nil {
		return *new([]CapabilityRegistryCapability), err
	}

	out0 := *abi.ConvertType(out[0], new([]CapabilityRegistryCapability)).(*[]CapabilityRegistryCapability)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetCapabilities() ([]CapabilityRegistryCapability, error) {
	return _CapabilityRegistry.Contract.GetCapabilities(&_CapabilityRegistry.CallOpts)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetCapabilities() ([]CapabilityRegistryCapability, error) {
	return _CapabilityRegistry.Contract.GetCapabilities(&_CapabilityRegistry.CallOpts)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) GetCapability(opts *bind.CallOpts, capabilityID [32]byte) (CapabilityRegistryCapability, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getCapability", capabilityID)

	if err != nil {
		return *new(CapabilityRegistryCapability), err
	}

	out0 := *abi.ConvertType(out[0], new(CapabilityRegistryCapability)).(*CapabilityRegistryCapability)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetCapability(capabilityID [32]byte) (CapabilityRegistryCapability, error) {
	return _CapabilityRegistry.Contract.GetCapability(&_CapabilityRegistry.CallOpts, capabilityID)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetCapability(capabilityID [32]byte) (CapabilityRegistryCapability, error) {
	return _CapabilityRegistry.Contract.GetCapability(&_CapabilityRegistry.CallOpts, capabilityID)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) GetCapabilityID(opts *bind.CallOpts, capabilityType [32]byte, version [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getCapabilityID", capabilityType, version)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetCapabilityID(capabilityType [32]byte, version [32]byte) ([32]byte, error) {
	return _CapabilityRegistry.Contract.GetCapabilityID(&_CapabilityRegistry.CallOpts, capabilityType, version)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetCapabilityID(capabilityType [32]byte, version [32]byte) ([32]byte, error) {
	return _CapabilityRegistry.Contract.GetCapabilityID(&_CapabilityRegistry.CallOpts, capabilityType, version)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) GetNodeOperator(opts *bind.CallOpts, nodeOperatorId *big.Int) (CapabilityRegistryNodeOperator, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getNodeOperator", nodeOperatorId)

	if err != nil {
		return *new(CapabilityRegistryNodeOperator), err
	}

	out0 := *abi.ConvertType(out[0], new(CapabilityRegistryNodeOperator)).(*CapabilityRegistryNodeOperator)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetNodeOperator(nodeOperatorId *big.Int) (CapabilityRegistryNodeOperator, error) {
	return _CapabilityRegistry.Contract.GetNodeOperator(&_CapabilityRegistry.CallOpts, nodeOperatorId)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetNodeOperator(nodeOperatorId *big.Int) (CapabilityRegistryNodeOperator, error) {
	return _CapabilityRegistry.Contract.GetNodeOperator(&_CapabilityRegistry.CallOpts, nodeOperatorId)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) IsCapabilityDeprecated(opts *bind.CallOpts, capabilityId [32]byte) (bool, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "isCapabilityDeprecated", capabilityId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) IsCapabilityDeprecated(capabilityId [32]byte) (bool, error) {
	return _CapabilityRegistry.Contract.IsCapabilityDeprecated(&_CapabilityRegistry.CallOpts, capabilityId)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) IsCapabilityDeprecated(capabilityId [32]byte) (bool, error) {
	return _CapabilityRegistry.Contract.IsCapabilityDeprecated(&_CapabilityRegistry.CallOpts, capabilityId)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) Owner() (common.Address, error) {
	return _CapabilityRegistry.Contract.Owner(&_CapabilityRegistry.CallOpts)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) Owner() (common.Address, error) {
	return _CapabilityRegistry.Contract.Owner(&_CapabilityRegistry.CallOpts)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) TypeAndVersion() (string, error) {
	return _CapabilityRegistry.Contract.TypeAndVersion(&_CapabilityRegistry.CallOpts)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) TypeAndVersion() (string, error) {
	return _CapabilityRegistry.Contract.TypeAndVersion(&_CapabilityRegistry.CallOpts)
}

func (_CapabilityRegistry *CapabilityRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_CapabilityRegistry *CapabilityRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AcceptOwnership(&_CapabilityRegistry.TransactOpts)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AcceptOwnership(&_CapabilityRegistry.TransactOpts)
}

func (_CapabilityRegistry *CapabilityRegistryTransactor) AddCapability(opts *bind.TransactOpts, capability CapabilityRegistryCapability) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "addCapability", capability)
}

func (_CapabilityRegistry *CapabilityRegistrySession) AddCapability(capability CapabilityRegistryCapability) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddCapability(&_CapabilityRegistry.TransactOpts, capability)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) AddCapability(capability CapabilityRegistryCapability) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddCapability(&_CapabilityRegistry.TransactOpts, capability)
}

func (_CapabilityRegistry *CapabilityRegistryTransactor) AddNodeOperators(opts *bind.TransactOpts, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "addNodeOperators", nodeOperators)
}

func (_CapabilityRegistry *CapabilityRegistrySession) AddNodeOperators(nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddNodeOperators(&_CapabilityRegistry.TransactOpts, nodeOperators)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) AddNodeOperators(nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddNodeOperators(&_CapabilityRegistry.TransactOpts, nodeOperators)
}

func (_CapabilityRegistry *CapabilityRegistryTransactor) DeprecateCapability(opts *bind.TransactOpts, capabilityId [32]byte) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "deprecateCapability", capabilityId)
}

func (_CapabilityRegistry *CapabilityRegistrySession) DeprecateCapability(capabilityId [32]byte) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.DeprecateCapability(&_CapabilityRegistry.TransactOpts, capabilityId)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) DeprecateCapability(capabilityId [32]byte) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.DeprecateCapability(&_CapabilityRegistry.TransactOpts, capabilityId)
}

func (_CapabilityRegistry *CapabilityRegistryTransactor) RemoveNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []*big.Int) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "removeNodeOperators", nodeOperatorIds)
}

func (_CapabilityRegistry *CapabilityRegistrySession) RemoveNodeOperators(nodeOperatorIds []*big.Int) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.RemoveNodeOperators(&_CapabilityRegistry.TransactOpts, nodeOperatorIds)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) RemoveNodeOperators(nodeOperatorIds []*big.Int) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.RemoveNodeOperators(&_CapabilityRegistry.TransactOpts, nodeOperatorIds)
}

func (_CapabilityRegistry *CapabilityRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_CapabilityRegistry *CapabilityRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.TransferOwnership(&_CapabilityRegistry.TransactOpts, to)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.TransferOwnership(&_CapabilityRegistry.TransactOpts, to)
}

func (_CapabilityRegistry *CapabilityRegistryTransactor) UpdateNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []*big.Int, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "updateNodeOperators", nodeOperatorIds, nodeOperators)
}

func (_CapabilityRegistry *CapabilityRegistrySession) UpdateNodeOperators(nodeOperatorIds []*big.Int, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.UpdateNodeOperators(&_CapabilityRegistry.TransactOpts, nodeOperatorIds, nodeOperators)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) UpdateNodeOperators(nodeOperatorIds []*big.Int, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.UpdateNodeOperators(&_CapabilityRegistry.TransactOpts, nodeOperatorIds, nodeOperators)
}

type CapabilityRegistryCapabilityAddedIterator struct {
	Event *CapabilityRegistryCapabilityAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilityRegistryCapabilityAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilityRegistryCapabilityAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(CapabilityRegistryCapabilityAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *CapabilityRegistryCapabilityAddedIterator) Error() error {
	return it.fail
}

func (it *CapabilityRegistryCapabilityAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilityRegistryCapabilityAdded struct {
	CapabilityId [32]byte
	Raw          types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterCapabilityAdded(opts *bind.FilterOpts, capabilityId [][32]byte) (*CapabilityRegistryCapabilityAddedIterator, error) {

	var capabilityIdRule []interface{}
	for _, capabilityIdItem := range capabilityId {
		capabilityIdRule = append(capabilityIdRule, capabilityIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "CapabilityAdded", capabilityIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryCapabilityAddedIterator{contract: _CapabilityRegistry.contract, event: "CapabilityAdded", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchCapabilityAdded(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryCapabilityAdded, capabilityId [][32]byte) (event.Subscription, error) {

	var capabilityIdRule []interface{}
	for _, capabilityIdItem := range capabilityId {
		capabilityIdRule = append(capabilityIdRule, capabilityIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "CapabilityAdded", capabilityIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilityRegistryCapabilityAdded)
				if err := _CapabilityRegistry.contract.UnpackLog(event, "CapabilityAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) ParseCapabilityAdded(log types.Log) (*CapabilityRegistryCapabilityAdded, error) {
	event := new(CapabilityRegistryCapabilityAdded)
	if err := _CapabilityRegistry.contract.UnpackLog(event, "CapabilityAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilityRegistryCapabilityDeprecatedIterator struct {
	Event *CapabilityRegistryCapabilityDeprecated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilityRegistryCapabilityDeprecatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilityRegistryCapabilityDeprecated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(CapabilityRegistryCapabilityDeprecated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *CapabilityRegistryCapabilityDeprecatedIterator) Error() error {
	return it.fail
}

func (it *CapabilityRegistryCapabilityDeprecatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilityRegistryCapabilityDeprecated struct {
	CapabilityId [32]byte
	Raw          types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterCapabilityDeprecated(opts *bind.FilterOpts, capabilityId [][32]byte) (*CapabilityRegistryCapabilityDeprecatedIterator, error) {

	var capabilityIdRule []interface{}
	for _, capabilityIdItem := range capabilityId {
		capabilityIdRule = append(capabilityIdRule, capabilityIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "CapabilityDeprecated", capabilityIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryCapabilityDeprecatedIterator{contract: _CapabilityRegistry.contract, event: "CapabilityDeprecated", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchCapabilityDeprecated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryCapabilityDeprecated, capabilityId [][32]byte) (event.Subscription, error) {

	var capabilityIdRule []interface{}
	for _, capabilityIdItem := range capabilityId {
		capabilityIdRule = append(capabilityIdRule, capabilityIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "CapabilityDeprecated", capabilityIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilityRegistryCapabilityDeprecated)
				if err := _CapabilityRegistry.contract.UnpackLog(event, "CapabilityDeprecated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) ParseCapabilityDeprecated(log types.Log) (*CapabilityRegistryCapabilityDeprecated, error) {
	event := new(CapabilityRegistryCapabilityDeprecated)
	if err := _CapabilityRegistry.contract.UnpackLog(event, "CapabilityDeprecated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilityRegistryNodeOperatorAddedIterator struct {
	Event *CapabilityRegistryNodeOperatorAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilityRegistryNodeOperatorAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilityRegistryNodeOperatorAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(CapabilityRegistryNodeOperatorAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *CapabilityRegistryNodeOperatorAddedIterator) Error() error {
	return it.fail
}

func (it *CapabilityRegistryNodeOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilityRegistryNodeOperatorAdded struct {
	NodeOperatorId *big.Int
	Admin          common.Address
	Name           string
	Raw            types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterNodeOperatorAdded(opts *bind.FilterOpts, admin []common.Address) (*CapabilityRegistryNodeOperatorAddedIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "NodeOperatorAdded", adminRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryNodeOperatorAddedIterator{contract: _CapabilityRegistry.contract, event: "NodeOperatorAdded", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchNodeOperatorAdded(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeOperatorAdded, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "NodeOperatorAdded", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilityRegistryNodeOperatorAdded)
				if err := _CapabilityRegistry.contract.UnpackLog(event, "NodeOperatorAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) ParseNodeOperatorAdded(log types.Log) (*CapabilityRegistryNodeOperatorAdded, error) {
	event := new(CapabilityRegistryNodeOperatorAdded)
	if err := _CapabilityRegistry.contract.UnpackLog(event, "NodeOperatorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilityRegistryNodeOperatorRemovedIterator struct {
	Event *CapabilityRegistryNodeOperatorRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilityRegistryNodeOperatorRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilityRegistryNodeOperatorRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(CapabilityRegistryNodeOperatorRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *CapabilityRegistryNodeOperatorRemovedIterator) Error() error {
	return it.fail
}

func (it *CapabilityRegistryNodeOperatorRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilityRegistryNodeOperatorRemoved struct {
	NodeOperatorId *big.Int
	Raw            types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterNodeOperatorRemoved(opts *bind.FilterOpts) (*CapabilityRegistryNodeOperatorRemovedIterator, error) {

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "NodeOperatorRemoved")
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryNodeOperatorRemovedIterator{contract: _CapabilityRegistry.contract, event: "NodeOperatorRemoved", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchNodeOperatorRemoved(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeOperatorRemoved) (event.Subscription, error) {

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "NodeOperatorRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilityRegistryNodeOperatorRemoved)
				if err := _CapabilityRegistry.contract.UnpackLog(event, "NodeOperatorRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) ParseNodeOperatorRemoved(log types.Log) (*CapabilityRegistryNodeOperatorRemoved, error) {
	event := new(CapabilityRegistryNodeOperatorRemoved)
	if err := _CapabilityRegistry.contract.UnpackLog(event, "NodeOperatorRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilityRegistryNodeOperatorUpdatedIterator struct {
	Event *CapabilityRegistryNodeOperatorUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilityRegistryNodeOperatorUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilityRegistryNodeOperatorUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(CapabilityRegistryNodeOperatorUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *CapabilityRegistryNodeOperatorUpdatedIterator) Error() error {
	return it.fail
}

func (it *CapabilityRegistryNodeOperatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilityRegistryNodeOperatorUpdated struct {
	NodeOperatorId *big.Int
	Admin          common.Address
	Name           string
	Raw            types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterNodeOperatorUpdated(opts *bind.FilterOpts, admin []common.Address) (*CapabilityRegistryNodeOperatorUpdatedIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "NodeOperatorUpdated", adminRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryNodeOperatorUpdatedIterator{contract: _CapabilityRegistry.contract, event: "NodeOperatorUpdated", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchNodeOperatorUpdated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeOperatorUpdated, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "NodeOperatorUpdated", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilityRegistryNodeOperatorUpdated)
				if err := _CapabilityRegistry.contract.UnpackLog(event, "NodeOperatorUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) ParseNodeOperatorUpdated(log types.Log) (*CapabilityRegistryNodeOperatorUpdated, error) {
	event := new(CapabilityRegistryNodeOperatorUpdated)
	if err := _CapabilityRegistry.contract.UnpackLog(event, "NodeOperatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilityRegistryOwnershipTransferRequestedIterator struct {
	Event *CapabilityRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilityRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilityRegistryOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(CapabilityRegistryOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *CapabilityRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *CapabilityRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilityRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CapabilityRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryOwnershipTransferRequestedIterator{contract: _CapabilityRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilityRegistryOwnershipTransferRequested)
				if err := _CapabilityRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*CapabilityRegistryOwnershipTransferRequested, error) {
	event := new(CapabilityRegistryOwnershipTransferRequested)
	if err := _CapabilityRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilityRegistryOwnershipTransferredIterator struct {
	Event *CapabilityRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilityRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilityRegistryOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(CapabilityRegistryOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *CapabilityRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *CapabilityRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilityRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CapabilityRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryOwnershipTransferredIterator{contract: _CapabilityRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilityRegistryOwnershipTransferred)
				if err := _CapabilityRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*CapabilityRegistryOwnershipTransferred, error) {
	event := new(CapabilityRegistryOwnershipTransferred)
	if err := _CapabilityRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_CapabilityRegistry *CapabilityRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _CapabilityRegistry.abi.Events["CapabilityAdded"].ID:
		return _CapabilityRegistry.ParseCapabilityAdded(log)
	case _CapabilityRegistry.abi.Events["CapabilityDeprecated"].ID:
		return _CapabilityRegistry.ParseCapabilityDeprecated(log)
	case _CapabilityRegistry.abi.Events["NodeOperatorAdded"].ID:
		return _CapabilityRegistry.ParseNodeOperatorAdded(log)
	case _CapabilityRegistry.abi.Events["NodeOperatorRemoved"].ID:
		return _CapabilityRegistry.ParseNodeOperatorRemoved(log)
	case _CapabilityRegistry.abi.Events["NodeOperatorUpdated"].ID:
		return _CapabilityRegistry.ParseNodeOperatorUpdated(log)
	case _CapabilityRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _CapabilityRegistry.ParseOwnershipTransferRequested(log)
	case _CapabilityRegistry.abi.Events["OwnershipTransferred"].ID:
		return _CapabilityRegistry.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CapabilityRegistryCapabilityAdded) Topic() common.Hash {
	return common.HexToHash("0x65610e5677eedff94555572640e442f89848a109ef8593fa927ac30b2565ff06")
}

func (CapabilityRegistryCapabilityDeprecated) Topic() common.Hash {
	return common.HexToHash("0xdcea1b78b6ddc31592a94607d537543fcaafda6cc52d6d5cc7bbfca1422baf21")
}

func (CapabilityRegistryNodeOperatorAdded) Topic() common.Hash {
	return common.HexToHash("0xda6697b182650034bd205cdc2dbfabb06bdb3a0a83a2b45bfefa3c4881284e0b")
}

func (CapabilityRegistryNodeOperatorRemoved) Topic() common.Hash {
	return common.HexToHash("0x1e5877d7b3001d1569bf733b76c7eceda58bd6c031e5b8d0b7042308ba2e9d4f")
}

func (CapabilityRegistryNodeOperatorUpdated) Topic() common.Hash {
	return common.HexToHash("0x14c8f513e8a6d86d2d16b0cb64976de4e72386c4f8068eca3b7354373f8fe97a")
}

func (CapabilityRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (CapabilityRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_CapabilityRegistry *CapabilityRegistry) Address() common.Address {
	return _CapabilityRegistry.address
}

type CapabilityRegistryInterface interface {
	GetCapabilities(opts *bind.CallOpts) ([]CapabilityRegistryCapability, error)

	GetCapability(opts *bind.CallOpts, capabilityID [32]byte) (CapabilityRegistryCapability, error)

	GetCapabilityID(opts *bind.CallOpts, capabilityType [32]byte, version [32]byte) ([32]byte, error)

	GetNodeOperator(opts *bind.CallOpts, nodeOperatorId *big.Int) (CapabilityRegistryNodeOperator, error)

	IsCapabilityDeprecated(opts *bind.CallOpts, capabilityId [32]byte) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddCapability(opts *bind.TransactOpts, capability CapabilityRegistryCapability) (*types.Transaction, error)

	AddNodeOperators(opts *bind.TransactOpts, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error)

	DeprecateCapability(opts *bind.TransactOpts, capabilityId [32]byte) (*types.Transaction, error)

	RemoveNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []*big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []*big.Int, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error)

	FilterCapabilityAdded(opts *bind.FilterOpts, capabilityId [][32]byte) (*CapabilityRegistryCapabilityAddedIterator, error)

	WatchCapabilityAdded(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryCapabilityAdded, capabilityId [][32]byte) (event.Subscription, error)

	ParseCapabilityAdded(log types.Log) (*CapabilityRegistryCapabilityAdded, error)

	FilterCapabilityDeprecated(opts *bind.FilterOpts, capabilityId [][32]byte) (*CapabilityRegistryCapabilityDeprecatedIterator, error)

	WatchCapabilityDeprecated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryCapabilityDeprecated, capabilityId [][32]byte) (event.Subscription, error)

	ParseCapabilityDeprecated(log types.Log) (*CapabilityRegistryCapabilityDeprecated, error)

	FilterNodeOperatorAdded(opts *bind.FilterOpts, admin []common.Address) (*CapabilityRegistryNodeOperatorAddedIterator, error)

	WatchNodeOperatorAdded(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeOperatorAdded, admin []common.Address) (event.Subscription, error)

	ParseNodeOperatorAdded(log types.Log) (*CapabilityRegistryNodeOperatorAdded, error)

	FilterNodeOperatorRemoved(opts *bind.FilterOpts) (*CapabilityRegistryNodeOperatorRemovedIterator, error)

	WatchNodeOperatorRemoved(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeOperatorRemoved) (event.Subscription, error)

	ParseNodeOperatorRemoved(log types.Log) (*CapabilityRegistryNodeOperatorRemoved, error)

	FilterNodeOperatorUpdated(opts *bind.FilterOpts, admin []common.Address) (*CapabilityRegistryNodeOperatorUpdatedIterator, error)

	WatchNodeOperatorUpdated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeOperatorUpdated, admin []common.Address) (event.Subscription, error)

	ParseNodeOperatorUpdated(log types.Log) (*CapabilityRegistryNodeOperatorUpdated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CapabilityRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*CapabilityRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CapabilityRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*CapabilityRegistryOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
