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
	LabelledName          [32]byte
	Version               [32]byte
	ResponseType          uint8
	ConfigurationContract common.Address
}

type CapabilityRegistryNode struct {
	NodeOperatorId               *big.Int
	P2pId                        [32]byte
	Signer                       common.Address
	SupportedHashedCapabilityIds [][32]byte
}

type CapabilityRegistryNodeOperator struct {
	Admin common.Address
	Name  string
}

var CapabilityRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityAlreadyDeprecated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CapabilityAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedConfigurationContract\",\"type\":\"address\"}],\"name\":\"InvalidCapabilityConfigurationContractInterface\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"InvalidNodeCapabilities\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeOperatorAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"InvalidNodeP2PId\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lengthOne\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lengthTwo\",\"type\":\"uint256\"}],\"name\":\"LengthMismatch\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDeprecated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"}],\"name\":\"NodeAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"}],\"name\":\"NodeOperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"labelledName\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability\",\"name\":\"capability\",\"type\":\"tuple\"}],\"name\":\"addCapability\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"addNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"supportedHashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"addNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"deprecateCapability\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCapabilities\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"labelledName\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedId\",\"type\":\"bytes32\"}],\"name\":\"getCapability\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"labelledName\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"labelledName\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"}],\"name\":\"getHashedCapabilityId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"getNode\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"supportedHashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.Node\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"}],\"name\":\"getNodeOperator\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"isCapabilityDeprecated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint256[]\"}],\"name\":\"removeNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"updateNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b503380600081620000695760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200009c576200009c81620000a5565b50505062000150565b336001600160a01b03821603620000ff5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000060565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61214580620001606000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c806365c14dc711610097578063ae3c241c11610066578063ae3c241c14610292578063c2d483a1146102a5578063ddbe4f82146102b8578063f2fde38b146102cd57600080fd5b806365c14dc71461022257806379ba5097146102425780638da5cb5b1461024a5780639cb7c5f41461027257600080fd5b80631cdf6343116100d35780631cdf63431461019457806336b402fb146101a7578063398f3773146101ef57806350c946fe1461020257600080fd5b80630c5801e314610105578063117392ce1461011a578063125700111461012d578063181f5a7714610155575b600080fd5b6101186101133660046116ec565b6102e0565b005b610118610128366004611758565b6105f1565b61014061013b366004611770565b61083c565b60405190151581526020015b60405180910390f35b604080518082018252601881527f4361706162696c697479526567697374727920312e302e3000000000000000006020820152905161014c91906117ed565b6101186101a2366004611800565b61084f565b6101e16101b5366004611842565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b60405190815260200161014c565b6101186101fd366004611800565b610912565b610215610210366004611770565b610aab565b60405161014c9190611864565b610235610230366004611770565b610b73565b60405161014c91906118ea565b610118610c50565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161014c565b610285610280366004611770565b610d4d565b60405161014c91906119cc565b6101186102a0366004611770565b610df7565b6101186102b3366004611800565b610ec2565b6102c061124c565b60405161014c91906119da565b6101186102db366004611a4a565b611391565b828114610328576040517fab8b67c600000000000000000000000000000000000000000000000000000000815260048101849052602481018290526044015b60405180910390fd5b6000805473ffffffffffffffffffffffffffffffffffffffff16905b848110156105e957600086868381811061036057610360611a67565b905060200201359050600085858481811061037d5761037d611a67565b905060200281019061038f9190611a96565b61039890611b9e565b805190915073ffffffffffffffffffffffffffffffffffffffff166103e9576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805173ffffffffffffffffffffffffffffffffffffffff16331480159061042657503373ffffffffffffffffffffffffffffffffffffffff851614155b1561045d576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160008381526007602052604090205473ffffffffffffffffffffffffffffffffffffffff908116911614158061050f57506020808201516040516104a392016117ed565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201206000868152600783529290922091926104f6926001019101611cb7565b6040516020818303038152906040528051906020012014155b156105d6578051600083815260076020908152604090912080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90931692909217825582015160019091019061057c9082611da6565b50806000015173ffffffffffffffffffffffffffffffffffffffff167f14c8f513e8a6d86d2d16b0cb64976de4e72386c4f8068eca3b7354373f8fe97a8383602001516040516105cd929190611ec0565b60405180910390a25b5050806105e290611f08565b9050610344565b505050505050565b6105f96113a5565b604080518235602082810191909152808401358284015282518083038401815260609092019092528051910120610631600382611428565b15610668576040517fe288638f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061067a6080840160608501611a4a565b73ffffffffffffffffffffffffffffffffffffffff16146107e5576106a56080830160608401611a4a565b73ffffffffffffffffffffffffffffffffffffffff163b158061078557506106d36080830160608401611a4a565b6040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527f884efe6100000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff91909116906301ffc9a790602401602060405180830381865afa15801561075f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107839190611f40565b155b156107e55761079a6080830160608401611a4a565b6040517fabb5e3fd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116600482015260240161031f565b6107f0600382611443565b506000818152600260205260409020829061080b8282611f62565b505060405181907f65610e5677eedff94555572640e442f89848a109ef8593fa927ac30b2565ff0690600090a25050565b6000610849600583611428565b92915050565b6108576113a5565b60005b8181101561090d57600083838381811061087657610876611a67565b60209081029290920135600081815260079093526040832080547fffffffffffffffffffffffff00000000000000000000000000000000000000001681559093509190506108c76001830182611606565b50506040518181527f1e5877d7b3001d1569bf733b76c7eceda58bd6c031e5b8d0b7042308ba2e9d4f9060200160405180910390a15061090681611f08565b905061085a565b505050565b61091a6113a5565b60005b8181101561090d57600083838381811061093957610939611a67565b905060200281019061094b9190611a96565b61095490611b9e565b805190915073ffffffffffffffffffffffffffffffffffffffff166109a5576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600954604080518082018252835173ffffffffffffffffffffffffffffffffffffffff908116825260208086015181840190815260008681526007909252939020825181547fffffffffffffffffffffffff000000000000000000000000000000000000000016921691909117815591519091906001820190610a289082611da6565b50905050600960008154610a3b90611f08565b909155508151602083015160405173ffffffffffffffffffffffffffffffffffffffff909216917fda6697b182650034bd205cdc2dbfabb06bdb3a0a83a2b45bfefa3c4881284e0b91610a9091859190611ec0565b60405180910390a2505080610aa490611f08565b905061091d565b6040805160808101825260008082526020820181905291810191909152606080820152600082815260086020908152604091829020825160808101845281548152600182015481840152600282015473ffffffffffffffffffffffffffffffffffffffff16818501526003820180548551818602810186019096528086529194929360608601939290830182828015610b6357602002820191906000526020600020905b815481526020019060010190808311610b4f575b5050505050815250509050919050565b6040805180820190915260008152606060208201526000828152600760209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff1683526001810180549192840191610bd090611c6a565b80601f0160208091040260200160405190810160405280929190818152602001828054610bfc90611c6a565b8015610b635780601f10610c1e57610100808354040283529160200191610b63565b820191906000526020600020905b815481529060010190602001808311610c2c57505050919092525091949350505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610cd1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161031f565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b604080516080808201835260008083526020808401829052838501829052606084018290528582526002808252918590208551938401865280548452600180820154928501929092529182015493949293919284019160ff1690811115610db657610db661192d565b6001811115610dc757610dc761192d565b815260029190910154610100900473ffffffffffffffffffffffffffffffffffffffff1660209091015292915050565b610dff6113a5565b610e0a600382611428565b610e43576040517fe181733f0000000000000000000000000000000000000000000000000000000081526004810182905260240161031f565b610e4e600582611428565b15610e88576040517f16950d1d0000000000000000000000000000000000000000000000000000000081526004810182905260240161031f565b610e93600582611443565b5060405181907fdcea1b78b6ddc31592a94607d537543fcaafda6cc52d6d5cc7bbfca1422baf2190600090a250565b60005b8181101561090d576000838383818110610ee157610ee1611a67565b9050602002810190610ef39190611fe4565b610efc90612018565b805160009081526007602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff168252600181018054959650939491939092840191610f4c90611c6a565b80601f0160208091040260200160405190810160405280929190818152602001828054610f7890611c6a565b8015610fc55780601f10610f9a57610100808354040283529160200191610fc5565b820191906000526020600020905b815481529060010190602001808311610fa857829003601f168201915b5050505050815250509050806000015173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611039576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6020808301516000908152600890915260409020600301541515808061106157506020830151155b156110a05782602001516040517f64e2ee9200000000000000000000000000000000000000000000000000000000815260040161031f91815260200190565b8260600151516000036110e55782606001516040517f3748d4c600000000000000000000000000000000000000000000000000000000815260040161031f91906120ed565b60005b836060015151811015611172576111268460600151828151811061110e5761110e611a67565b6020026020010151600361142890919063ffffffff16565b6111625783606001516040517f3748d4c600000000000000000000000000000000000000000000000000000000815260040161031f91906120ed565b61116b81611f08565b90506110e8565b506020838101805160009081526008835260409081902086518155915160018301558501516002820180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9092169190911790556060850151805186936111f7926003850192910190611640565b505050602083810151845160408051928352928201527f5bfe8a52ad26ac6ee7b0cd46d2fd92be04735a31c45ef8aa3d4b7ea1b61bbc1f910160405180910390a15050508061124590611f08565b9050610ec5565b6060600061125a600361144f565b90506000611268600561145c565b82516112749190612125565b67ffffffffffffffff81111561128c5761128c611ad4565b6040519080825280602002602001820160405280156112fc57816020015b6040805160808101825260008082526020808301829052928201819052606082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816112aa5790505b5090506000805b835181101561138857600084828151811061132057611320611a67565b6020026020010151905061133e81600561142890919063ffffffff16565b6113775761134b81610d4d565b84848151811061135d5761135d611a67565b6020026020010181905250828061137390611f08565b9350505b5061138181611f08565b9050611303565b50909392505050565b6113996113a5565b6113a281611466565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314611426576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161031f565b565b600081815260018301602052604081205415155b9392505050565b600061143c838361155b565b6060600061143c836115aa565b6000610849825490565b3373ffffffffffffffffffffffffffffffffffffffff8216036114e5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161031f565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008181526001830160205260408120546115a257508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610849565b506000610849565b6060816000018054806020026020016040519081016040528092919081815260200182805480156115fa57602002820191906000526020600020905b8154815260200190600101908083116115e6575b50505050509050919050565b50805461161290611c6a565b6000825580601f10611622575050565b601f0160209004906000526020600020908101906113a2919061168b565b82805482825590600052602060002090810192821561167b579160200282015b8281111561167b578251825591602001919060010190611660565b5061168792915061168b565b5090565b5b80821115611687576000815560010161168c565b60008083601f8401126116b257600080fd5b50813567ffffffffffffffff8111156116ca57600080fd5b6020830191508360208260051b85010111156116e557600080fd5b9250929050565b6000806000806040858703121561170257600080fd5b843567ffffffffffffffff8082111561171a57600080fd5b611726888389016116a0565b9096509450602087013591508082111561173f57600080fd5b5061174c878288016116a0565b95989497509550505050565b60006080828403121561176a57600080fd5b50919050565b60006020828403121561178257600080fd5b5035919050565b6000815180845260005b818110156117af57602081850181015186830182015201611793565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061143c6020830184611789565b6000806020838503121561181357600080fd5b823567ffffffffffffffff81111561182a57600080fd5b611836858286016116a0565b90969095509350505050565b6000806040838503121561185557600080fd5b50508035926020909101359150565b6000602080835260a0830184518285015281850151604085015273ffffffffffffffffffffffffffffffffffffffff6040860151166060850152606085015160808086015281815180845260c0870191508483019350600092505b808310156118df57835182529284019260019290920191908401906118bf565b509695505050505050565b6020815273ffffffffffffffffffffffffffffffffffffffff8251166020820152600060208301516040808401526119256060840182611789565b949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b80518252602081015160208301526040810151600281106119a6577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b604083015260609081015173ffffffffffffffffffffffffffffffffffffffff16910152565b60808101610849828461195c565b6020808252825182820181905260009190848201906040850190845b81811015611a1c57611a0983855161195c565b92840192608092909201916001016119f6565b50909695505050505050565b73ffffffffffffffffffffffffffffffffffffffff811681146113a257600080fd5b600060208284031215611a5c57600080fd5b813561143c81611a28565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc1833603018112611aca57600080fd5b9190910192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715611b2657611b26611ad4565b60405290565b6040516080810167ffffffffffffffff81118282101715611b2657611b26611ad4565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611b9657611b96611ad4565b604052919050565b600060408236031215611bb057600080fd5b611bb8611b03565b8235611bc381611a28565b815260208381013567ffffffffffffffff80821115611be157600080fd5b9085019036601f830112611bf457600080fd5b813581811115611c0657611c06611ad4565b611c36847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611b4f565b91508082523684828501011115611c4c57600080fd5b80848401858401376000908201840152918301919091525092915050565b600181811c90821680611c7e57607f821691505b60208210810361176a577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000602080835260008454611ccb81611c6a565b80848701526040600180841660008114611cec5760018114611d2457611d52565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838a01528284151560051b8a01019550611d52565b896000528660002060005b85811015611d4a5781548b8201860152908301908801611d2f565b8a0184019650505b509398975050505050505050565b601f82111561090d57600081815260208120601f850160051c81016020861015611d875750805b601f850160051c820191505b818110156105e957828155600101611d93565b815167ffffffffffffffff811115611dc057611dc0611ad4565b611dd481611dce8454611c6a565b84611d60565b602080601f831160018114611e275760008415611df15750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556105e9565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015611e7457888601518255948401946001909101908401611e55565b5085821015611eb057878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b8281526040602082015260006119256040830184611789565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611f3957611f39611ed9565b5060010190565b600060208284031215611f5257600080fd5b8151801515811461143c57600080fd5b813581556020820135600182015560028101604083013560028110611f8657600080fd5b81546060850135611f9681611a28565b74ffffffffffffffffffffffffffffffffffffffff008160081b1660ff84167fffffffffffffffffffffff000000000000000000000000000000000000000000841617178455505050505050565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81833603018112611aca57600080fd5b60006080823603121561202a57600080fd5b612032611b2c565b8235815260208084013581830152604084013561204e81611a28565b6040830152606084013567ffffffffffffffff8082111561206e57600080fd5b9085019036601f83011261208157600080fd5b81358181111561209357612093611ad4565b8060051b91506120a4848301611b4f565b81815291830184019184810190368411156120be57600080fd5b938501935b838510156120dc578435825293850193908501906120c3565b606087015250939695505050505050565b6020808252825182820181905260009190848201906040850190845b81811015611a1c57835183529284019291840191600101612109565b8181038181111561084957610849611ed956fea164736f6c6343000813000a",
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

func (_CapabilityRegistry *CapabilityRegistryCaller) GetCapability(opts *bind.CallOpts, hashedId [32]byte) (CapabilityRegistryCapability, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getCapability", hashedId)

	if err != nil {
		return *new(CapabilityRegistryCapability), err
	}

	out0 := *abi.ConvertType(out[0], new(CapabilityRegistryCapability)).(*CapabilityRegistryCapability)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetCapability(hashedId [32]byte) (CapabilityRegistryCapability, error) {
	return _CapabilityRegistry.Contract.GetCapability(&_CapabilityRegistry.CallOpts, hashedId)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetCapability(hashedId [32]byte) (CapabilityRegistryCapability, error) {
	return _CapabilityRegistry.Contract.GetCapability(&_CapabilityRegistry.CallOpts, hashedId)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) GetHashedCapabilityId(opts *bind.CallOpts, labelledName [32]byte, version [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getHashedCapabilityId", labelledName, version)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetHashedCapabilityId(labelledName [32]byte, version [32]byte) ([32]byte, error) {
	return _CapabilityRegistry.Contract.GetHashedCapabilityId(&_CapabilityRegistry.CallOpts, labelledName, version)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetHashedCapabilityId(labelledName [32]byte, version [32]byte) ([32]byte, error) {
	return _CapabilityRegistry.Contract.GetHashedCapabilityId(&_CapabilityRegistry.CallOpts, labelledName, version)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) GetNode(opts *bind.CallOpts, p2pId [32]byte) (CapabilityRegistryNode, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getNode", p2pId)

	if err != nil {
		return *new(CapabilityRegistryNode), err
	}

	out0 := *abi.ConvertType(out[0], new(CapabilityRegistryNode)).(*CapabilityRegistryNode)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetNode(p2pId [32]byte) (CapabilityRegistryNode, error) {
	return _CapabilityRegistry.Contract.GetNode(&_CapabilityRegistry.CallOpts, p2pId)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetNode(p2pId [32]byte) (CapabilityRegistryNode, error) {
	return _CapabilityRegistry.Contract.GetNode(&_CapabilityRegistry.CallOpts, p2pId)
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

func (_CapabilityRegistry *CapabilityRegistryCaller) IsCapabilityDeprecated(opts *bind.CallOpts, hashedCapabilityId [32]byte) (bool, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "isCapabilityDeprecated", hashedCapabilityId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) IsCapabilityDeprecated(hashedCapabilityId [32]byte) (bool, error) {
	return _CapabilityRegistry.Contract.IsCapabilityDeprecated(&_CapabilityRegistry.CallOpts, hashedCapabilityId)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) IsCapabilityDeprecated(hashedCapabilityId [32]byte) (bool, error) {
	return _CapabilityRegistry.Contract.IsCapabilityDeprecated(&_CapabilityRegistry.CallOpts, hashedCapabilityId)
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) AddNodes(opts *bind.TransactOpts, nodes []CapabilityRegistryNode) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "addNodes", nodes)
}

func (_CapabilityRegistry *CapabilityRegistrySession) AddNodes(nodes []CapabilityRegistryNode) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddNodes(&_CapabilityRegistry.TransactOpts, nodes)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) AddNodes(nodes []CapabilityRegistryNode) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddNodes(&_CapabilityRegistry.TransactOpts, nodes)
}

func (_CapabilityRegistry *CapabilityRegistryTransactor) DeprecateCapability(opts *bind.TransactOpts, hashedCapabilityId [32]byte) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "deprecateCapability", hashedCapabilityId)
}

func (_CapabilityRegistry *CapabilityRegistrySession) DeprecateCapability(hashedCapabilityId [32]byte) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.DeprecateCapability(&_CapabilityRegistry.TransactOpts, hashedCapabilityId)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) DeprecateCapability(hashedCapabilityId [32]byte) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.DeprecateCapability(&_CapabilityRegistry.TransactOpts, hashedCapabilityId)
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
	HashedCapabilityId [32]byte
	Raw                types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterCapabilityAdded(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilityRegistryCapabilityAddedIterator, error) {

	var hashedCapabilityIdRule []interface{}
	for _, hashedCapabilityIdItem := range hashedCapabilityId {
		hashedCapabilityIdRule = append(hashedCapabilityIdRule, hashedCapabilityIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "CapabilityAdded", hashedCapabilityIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryCapabilityAddedIterator{contract: _CapabilityRegistry.contract, event: "CapabilityAdded", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchCapabilityAdded(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryCapabilityAdded, hashedCapabilityId [][32]byte) (event.Subscription, error) {

	var hashedCapabilityIdRule []interface{}
	for _, hashedCapabilityIdItem := range hashedCapabilityId {
		hashedCapabilityIdRule = append(hashedCapabilityIdRule, hashedCapabilityIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "CapabilityAdded", hashedCapabilityIdRule)
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
	HashedCapabilityId [32]byte
	Raw                types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterCapabilityDeprecated(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilityRegistryCapabilityDeprecatedIterator, error) {

	var hashedCapabilityIdRule []interface{}
	for _, hashedCapabilityIdItem := range hashedCapabilityId {
		hashedCapabilityIdRule = append(hashedCapabilityIdRule, hashedCapabilityIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "CapabilityDeprecated", hashedCapabilityIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryCapabilityDeprecatedIterator{contract: _CapabilityRegistry.contract, event: "CapabilityDeprecated", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchCapabilityDeprecated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryCapabilityDeprecated, hashedCapabilityId [][32]byte) (event.Subscription, error) {

	var hashedCapabilityIdRule []interface{}
	for _, hashedCapabilityIdItem := range hashedCapabilityId {
		hashedCapabilityIdRule = append(hashedCapabilityIdRule, hashedCapabilityIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "CapabilityDeprecated", hashedCapabilityIdRule)
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

type CapabilityRegistryNodeAddedIterator struct {
	Event *CapabilityRegistryNodeAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilityRegistryNodeAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilityRegistryNodeAdded)
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
		it.Event = new(CapabilityRegistryNodeAdded)
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

func (it *CapabilityRegistryNodeAddedIterator) Error() error {
	return it.fail
}

func (it *CapabilityRegistryNodeAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilityRegistryNodeAdded struct {
	P2pId          [32]byte
	NodeOperatorId *big.Int
	Raw            types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterNodeAdded(opts *bind.FilterOpts) (*CapabilityRegistryNodeAddedIterator, error) {

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "NodeAdded")
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryNodeAddedIterator{contract: _CapabilityRegistry.contract, event: "NodeAdded", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchNodeAdded(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeAdded) (event.Subscription, error) {

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "NodeAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilityRegistryNodeAdded)
				if err := _CapabilityRegistry.contract.UnpackLog(event, "NodeAdded", log); err != nil {
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

func (_CapabilityRegistry *CapabilityRegistryFilterer) ParseNodeAdded(log types.Log) (*CapabilityRegistryNodeAdded, error) {
	event := new(CapabilityRegistryNodeAdded)
	if err := _CapabilityRegistry.contract.UnpackLog(event, "NodeAdded", log); err != nil {
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
	case _CapabilityRegistry.abi.Events["NodeAdded"].ID:
		return _CapabilityRegistry.ParseNodeAdded(log)
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

func (CapabilityRegistryNodeAdded) Topic() common.Hash {
	return common.HexToHash("0x5bfe8a52ad26ac6ee7b0cd46d2fd92be04735a31c45ef8aa3d4b7ea1b61bbc1f")
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

	GetCapability(opts *bind.CallOpts, hashedId [32]byte) (CapabilityRegistryCapability, error)

	GetHashedCapabilityId(opts *bind.CallOpts, labelledName [32]byte, version [32]byte) ([32]byte, error)

	GetNode(opts *bind.CallOpts, p2pId [32]byte) (CapabilityRegistryNode, error)

	GetNodeOperator(opts *bind.CallOpts, nodeOperatorId *big.Int) (CapabilityRegistryNodeOperator, error)

	IsCapabilityDeprecated(opts *bind.CallOpts, hashedCapabilityId [32]byte) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddCapability(opts *bind.TransactOpts, capability CapabilityRegistryCapability) (*types.Transaction, error)

	AddNodeOperators(opts *bind.TransactOpts, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error)

	AddNodes(opts *bind.TransactOpts, nodes []CapabilityRegistryNode) (*types.Transaction, error)

	DeprecateCapability(opts *bind.TransactOpts, hashedCapabilityId [32]byte) (*types.Transaction, error)

	RemoveNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []*big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []*big.Int, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error)

	FilterCapabilityAdded(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilityRegistryCapabilityAddedIterator, error)

	WatchCapabilityAdded(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryCapabilityAdded, hashedCapabilityId [][32]byte) (event.Subscription, error)

	ParseCapabilityAdded(log types.Log) (*CapabilityRegistryCapabilityAdded, error)

	FilterCapabilityDeprecated(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilityRegistryCapabilityDeprecatedIterator, error)

	WatchCapabilityDeprecated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryCapabilityDeprecated, hashedCapabilityId [][32]byte) (event.Subscription, error)

	ParseCapabilityDeprecated(log types.Log) (*CapabilityRegistryCapabilityDeprecated, error)

	FilterNodeAdded(opts *bind.FilterOpts) (*CapabilityRegistryNodeAddedIterator, error)

	WatchNodeAdded(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeAdded) (event.Subscription, error)

	ParseNodeAdded(log types.Log) (*CapabilityRegistryNodeAdded, error)

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
