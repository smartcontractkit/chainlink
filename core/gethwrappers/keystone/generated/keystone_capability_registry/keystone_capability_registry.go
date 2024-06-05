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
	LabelledName          string
	Version               string
	CapabilityType        uint8
	ResponseType          uint8
	ConfigurationContract common.Address
}

type CapabilityRegistryCapabilityConfiguration struct {
	CapabilityId [32]byte
	Config       []byte
}

type CapabilityRegistryDONInfo struct {
	Id                       uint32
	ConfigCount              uint32
	F                        uint32
	IsPublic                 bool
	AcceptsWorkflows         bool
	NodeP2PIds               [][32]byte
	CapabilityConfigurations []CapabilityRegistryCapabilityConfiguration
}

type CapabilityRegistryNodeInfo struct {
	NodeOperatorId      uint32
	Signer              [32]byte
	P2pId               [32]byte
	HashedCapabilityIds [][32]byte
}

type CapabilityRegistryNodeOperator struct {
	Admin common.Address
	Name  string
}

var CapabilityRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CapabilityAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityIsDeprecated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"DONDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONCapability\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONNode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedConfigurationContract\",\"type\":\"address\"}],\"name\":\"InvalidCapabilityConfigurationContractInterface\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"f\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"nodeCount\",\"type\":\"uint256\"}],\"name\":\"InvalidFaultTolerance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"InvalidNodeCapabilities\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeOperatorAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"InvalidNodeP2PId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lengthOne\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lengthTwo\",\"type\":\"uint256\"}],\"name\":\"LengthMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"NodeDoesNotSupportCapability\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDeprecated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"}],\"name\":\"NodeOperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"NodeRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability\",\"name\":\"capability\",\"type\":\"tuple\"}],\"name\":\"addCapability\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"f\",\"type\":\"uint32\"}],\"name\":\"addDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"addNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"addNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"deprecateCapability\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCapabilities\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedId\",\"type\":\"bytes32\"}],\"name\":\"getCapability\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"getDON\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"f\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilityRegistry.DONInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"getDONCapabilityConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"f\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilityRegistry.DONInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"}],\"name\":\"getHashedCapabilityId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"getNode\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nodeOperatorId\",\"type\":\"uint256\"}],\"name\":\"getNodeOperator\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeOperators\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodes\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"internalType\":\"uint32[]\",\"name\":\"\",\"type\":\"uint32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"isCapabilityDeprecated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"donIds\",\"type\":\"uint32[]\"}],\"name\":\"removeDONs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint256[]\"}],\"name\":\"removeNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"removedNodeP2PIds\",\"type\":\"bytes32[]\"}],\"name\":\"removeNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"f\",\"type\":\"uint32\"}],\"name\":\"updateDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"updateNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"updateNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052601280546001600160401b0319166401000000011790553480156200002857600080fd5b503380600081620000805760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000b357620000b381620000bc565b50505062000167565b336001600160a01b03821603620001165760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000077565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b614bef80620001776000396000f3fe608060405234801561001057600080fd5b50600436106101ae5760003560e01c80635e65e309116100ee57806395864d1f11610097578063b06e07a711610071578063b06e07a7146103ec578063ddbe4f82146103ff578063e29581aa14610414578063f2fde38b1461042a57600080fd5b806395864d1f146103a65780639cb7c5f4146103b9578063ae3c241c146103d957600080fd5b806373ac22b4116100c857806373ac22b41461036357806379ba5097146103765780638da5cb5b1461037e57600080fd5b80635e65e3091461031b57806365c14dc71461032e57806366acaa331461034e57600080fd5b80631cdf63431161015b578063235374051161013557806323537405146102b45780632c01a1e8146102d4578063398f3773146102e757806350c946fe146102fa57600080fd5b80631cdf6343146102795780631d05394c1461028c578063214502431461029f57600080fd5b8063178962481161018c5780631789624814610211578063181f5a7714610224578063193ec0061461026657600080fd5b80630c5801e3146101b35780630fe5800a146101c857806312570011146101ee575b600080fd5b6101c66101c1366004613755565b61043d565b005b6101db6101d6366004613803565b61074e565b6040519081526020015b60405180910390f35b6102016101fc366004613863565b610787565b60405190151581526020016101e5565b6101c661021f36600461387c565b61079a565b60408051808201909152601881527f4361706162696c697479526567697374727920312e302e30000000000000000060208201525b6040516101e5919061391b565b6101c6610274366004613955565b6109d6565b6101c66102873660046139fa565b610ab1565b6101c661029a3660046139fa565b610b81565b6102a7610cc1565b6040516101e59190613b64565b6102c76102c2366004613be4565b610e21565b6040516101e59190613bff565b6101c66102e23660046139fa565b610e65565b6101c66102f53660046139fa565b611109565b61030d610308366004613863565b6112cc565b6040516101e5929190613c53565b6101c66103293660046139fa565b611301565b61034161033c366004613863565b611804565b6040516101e59190613cb0565b6103566118ea565b6040516101e59190613cc3565b6101c66103713660046139fa565b611aa8565b6101c6611f0e565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101e5565b6101c66103b4366004613d36565b61200b565b6103cc6103c7366004613863565b6120d7565b6040516101e59190613ea7565b6101c66103e7366004613863565b6122dc565b6102596103fa366004613eba565b6123a7565b61040761247c565b6040516101e59190613ee4565b61041c612771565b6040516101e5929190613f57565b6101c6610438366004614038565b6128fc565b828114610485576040517fab8b67c600000000000000000000000000000000000000000000000000000000815260048101849052602481018290526044015b60405180910390fd5b6000805473ffffffffffffffffffffffffffffffffffffffff16905b848110156107465760008686838181106104bd576104bd614055565b90506020020135905060008585848181106104da576104da614055565b90506020028101906104ec9190614084565b6104f59061418c565b805190915073ffffffffffffffffffffffffffffffffffffffff16610546576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805173ffffffffffffffffffffffffffffffffffffffff16331480159061058357503373ffffffffffffffffffffffffffffffffffffffff851614155b156105ba576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80516000838152600f602052604090205473ffffffffffffffffffffffffffffffffffffffff908116911614158061066c5750602080820151604051610600920161391b565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201206000868152600f83529290922091926106539260010191016142ab565b6040516020818303038152906040528051906020012014155b156107335780516000838152600f6020908152604090912080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9093169290921782558201516001909101906106d9908261439a565b50806000015173ffffffffffffffffffffffffffffffffffffffff167f14c8f513e8a6d86d2d16b0cb64976de4e72386c4f8068eca3b7354373f8fe97a83836020015160405161072a9291906144b4565b60405180910390a25b50508061073f906144fc565b90506104a1565b505050505050565b600084848484604051602001610767949392919061457d565b604051602081830303815290604052805190602001209050949350505050565b6000610794600583612910565b92915050565b6107a261292b565b60006107be6107b183806145af565b6101d660208601866145af565b90506107cb600382612910565b15610802576040517fe288638f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061081460a0840160808501614038565b73ffffffffffffffffffffffffffffffffffffffff161461097f5761083f60a0830160808401614038565b73ffffffffffffffffffffffffffffffffffffffff163b158061091f575061086d60a0830160808401614038565b6040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527f73e8b41d00000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff91909116906301ffc9a790602401602060405180830381865afa1580156108f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061091d9190614614565b155b1561097f5761093460a0830160808401614038565b6040517fabb5e3fd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116600482015260240161047c565b61098a6003826129ae565b50600081815260026020526040902082906109a58282614801565b505060405181907f65610e5677eedff94555572640e442f89848a109ef8593fa927ac30b2565ff0690600090a25050565b6109de61292b565b60125463ffffffff640100000000909104811660008181526011602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001684179055815160a081018352838152600191810191909152861515918101919091528415156060820152918316608083015290610a6b9089908990899089906129ba565b60128054600490610a8990640100000000900463ffffffff166149c4565b91906101000a81548163ffffffff021916908363ffffffff1602179055505050505050505050565b610ab961292b565b60005b81811015610b7c576000838383818110610ad857610ad8614055565b602090810292909201356000818152600f9093526040832080547fffffffffffffffffffffffff0000000000000000000000000000000000000000168155909350919050610b2960018301826136bb565b50610b3790506007826129ae565b506040518181527f1e5877d7b3001d1569bf733b76c7eceda58bd6c031e5b8d0b7042308ba2e9d4f9060200160405180910390a150610b75816144fc565b9050610abc565b505050565b610b8961292b565b60005b81811015610b7c576000838383818110610ba857610ba8614055565b9050602002016020810190610bbd9190613be4565b63ffffffff808216600090815260116020526040812080549394509264010000000090049091169003610c24576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff8316600482015260240161047c565b63ffffffff808316600081815260116020526040902080547fffffffffffffffffffffffffffffffffffff0000000000000000000000000000169055610c6e91600991906129ae16565b506040805163ffffffff84168152600060208201527ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c1581703651910160405180910390a1505080610cba906144fc565b9050610b8c565b601254606090640100000000900463ffffffff1660006001610ce36009612f9c565b601254610cfe9190640100000000900463ffffffff166149e7565b610d0891906149e7565b67ffffffffffffffff811115610d2057610d206140c2565b604051908082528060200260200182016040528015610da757816020015b6040805160e081018252600080825260208083018290529282018190526060808301829052608083019190915260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181610d3e5790505b509050600060015b838163ffffffff161015610e1857610dd1600963ffffffff8084169061291016565b610e0857610dde81612fa6565b838381518110610df057610df0614055565b602002602001018190525081610e05906144fc565b91505b610e11816149c4565b9050610daf565b50909392505050565b6040805160e0810182526000808252602082018190529181018290526060808201839052608082019290925260a0810182905260c081019190915261079482612fa6565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015611103576000848483818110610e9f57610e9f614055565b602090810292909201356000818152601090935260409092206001015491925050151580610efc576040517f64e2ee920000000000000000000000000000000000000000000000000000000081526004810183905260240161047c565b60008281526010602090815260408083205463ffffffff168352600f82528083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff1682526001810180549293919291840191610f5690614258565b80601f0160208091040260200160405190810160405280929190818152602001828054610f8290614258565b8015610fcf5780601f10610fa457610100808354040283529160200191610fcf565b820191906000526020600020905b815481529060010190602001808311610fb257829003601f168201915b505050505081525050905084158015610fff5750805173ffffffffffffffffffffffffffffffffffffffff163314155b15611036576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526010602052604090206001015461105490600b90613282565b5060008381526010602052604090206002015461107390600d90613282565b5060008381526010602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001681556001810183905560020191909155517f5254e609a97bab37b7cc79fe128f85c097bd6015c6e1624ae0ba392eb9753205906110e79085815260200190565b60405180910390a1505050806110fc906144fc565b9050610e83565b50505050565b61111161292b565b60005b81811015610b7c57600083838381811061113057611130614055565b90506020028101906111429190614084565b61114b9061418c565b805190915073ffffffffffffffffffffffffffffffffffffffff1661119c576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254604080518082018252835173ffffffffffffffffffffffffffffffffffffffff908116825260208086015181840190815263ffffffff9095166000818152600f909252939020825181547fffffffffffffffffffffffff00000000000000000000000000000000000000001692169190911781559251919290916001820190611228908261439a565b5050601280549091506000906112439063ffffffff166149c4565b91906101000a81548163ffffffff021916908363ffffffff160217905550816000015173ffffffffffffffffffffffffffffffffffffffff167fda6697b182650034bd205cdc2dbfabb06bdb3a0a83a2b45bfefa3c4881284e0b8284602001516040516112b19291906144b4565b60405180910390a25050806112c5906144fc565b9050611114565b60408051608081018252600080825260208201819052918101829052606080820152906112f88361328e565b91509150915091565b60005b81811015610b7c57600083838381811061132057611320614055565b905060200281019061133291906149fa565b61133b90614a2e565b9050600061135e60005473ffffffffffffffffffffffffffffffffffffffff1690565b825163ffffffff166000908152600f602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff908116835260018201805496909116331496509394919390928401916113bb90614258565b80601f01602080910402602001604051908101604052809291908181526020018280546113e790614258565b80156114345780601f1061140957610100808354040283529160200191611434565b820191906000526020600020905b81548152906001019060200180831161141757829003601f168201915b5050505050815250509050811580156114645750805173ffffffffffffffffffffffffffffffffffffffff163314155b1561149b576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040808401516000908152601060205220600101541515806114f15783604001516040517f64e2ee9200000000000000000000000000000000000000000000000000000000815260040161047c91815260200190565b6020840151158061153757508360200151601060008660400151815260200190815260200160002060010154141580156115375750602084015161153790600b90612910565b1561156e576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606084015180516000036115b057806040517f3748d4c600000000000000000000000000000000000000000000000000000000815260040161047c9190614b01565b60408581015160009081526010602052208054640100000000900463ffffffff169060046115dd836149c4565b82546101009290920a63ffffffff8181021990931691831602179091556040878101516000908152601060205290812054640100000000900490911691505b82518110156116e85761165283828151811061163a5761163a614055565b6020026020010151600361291090919063ffffffff16565b61168a57826040517f3748d4c600000000000000000000000000000000000000000000000000000000815260040161047c9190614b01565b6116d783828151811061169f5761169f614055565b6020908102919091018101516040808b015160009081526010845281812063ffffffff8089168352600390910190945220916129ae16565b506116e1816144fc565b905061161c565b5085516040808801805160009081526010602090815283822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9096169590951790945581518082528382206002015581518152828120600190810154948b015192518252929020909101541461179d5761176c600b82613282565b50602080880180516040808b015160009081526010909452909220600101919091555161179b90600b906129ae565b505b60408781015188516020808b0151845193845263ffffffff909216908301528183015290517ff101cfc54994c31624d25789378d71ec4dbdc533e26a4ecc6b7648f4798d09169181900360600190a150505050505050806117fd906144fc565b9050611304565b6040805180820190915260008152606060208201526000828152600f60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff168352600181018054919284019161186190614258565b80601f016020809104026020016040519081016040528092919081815260200182805461188d90614258565b80156118da5780601f106118af576101008083540402835291602001916118da565b820191906000526020600020905b8154815290600101906020018083116118bd57829003601f168201915b5050505050815250509050919050565b60125460609063ffffffff16600060016119046007612f9c565b601254611917919063ffffffff166149e7565b61192191906149e7565b67ffffffffffffffff811115611939576119396140c2565b60405190808252806020026020018201604052801561197f57816020015b6040805180820190915260008152606060208201528152602001906001900390816119575790505b509050600060015b8363ffffffff16811015610e18576119a0600782612910565b611a98576000818152600f60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff16835260018101805491928401916119ec90614258565b80601f0160208091040260200160405190810160405280929190818152602001828054611a1890614258565b8015611a655780601f10611a3a57610100808354040283529160200191611a65565b820191906000526020600020905b815481529060010190602001808311611a4857829003601f168201915b505050505081525050838381518110611a8057611a80614055565b602002602001018190525081611a95906144fc565b91505b611aa1816144fc565b9050611987565b60005b81811015610b7c576000838383818110611ac757611ac7614055565b9050602002810190611ad991906149fa565b611ae290614a2e565b90506000611b0560005473ffffffffffffffffffffffffffffffffffffffff1690565b825163ffffffff166000908152600f602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff90811683526001820180549690911633149650939491939092840191611b6290614258565b80601f0160208091040260200160405190810160405280929190818152602001828054611b8e90614258565b8015611bdb5780601f10611bb057610100808354040283529160200191611bdb565b820191906000526020600020905b815481529060010190602001808311611bbe57829003601f168201915b505050505081525050905081158015611c0b5750805173ffffffffffffffffffffffffffffffffffffffff163314155b15611c42576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408084015160009081526010602052206001015415158080611c6757506040840151155b15611ca65783604001516040517f64e2ee9200000000000000000000000000000000000000000000000000000000815260040161047c91815260200190565b60208401511580611cc357506020840151611cc390600b90612910565b15611cfa576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608401518051600003611d3c57806040517f3748d4c600000000000000000000000000000000000000000000000000000000815260040161047c9190614b01565b60408581015160009081526010602052208054600490611d6990640100000000900463ffffffff166149c4565b82546101009290920a63ffffffff81810219909316918316021790915560408681015160009081526010602052908120546401000000009004909116905b8251811015611e2357611dc583828151811061163a5761163a614055565b611dfd57826040517f3748d4c600000000000000000000000000000000000000000000000000000000815260040161047c9190614b01565b611e1283828151811061169f5761169f614055565b50611e1c816144fc565b9050611da7565b5085516040808801805160009081526010602090815283822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9687161790558251808352848320600201558a018051925182529290206001015551611e9591600b91906129ae16565b506040860151611ea790600d906129ae565b5060408681015187516020808a0151845193845263ffffffff909216908301528183015290517fc9296aa9b0951d8000e8ed7f2b5be30c5106de8df3dbedf9a57c93f5f9e4d7da9181900360600190a150505050505080611f07906144fc565b9050611aab565b60015473ffffffffffffffffffffffffffffffffffffffff163314611f8f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161047c565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61201361292b565b63ffffffff808916600090815260116020526040812054640100000000900490911690819003612077576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff8a16600482015260240161047c565b6120cc888888886040518060a001604052808f63ffffffff1681526020018761209f906149c4565b97508763ffffffff1681526020018a1515815260200189151581526020018863ffffffff168152506129ba565b505050505050505050565b6121076040805160a081018252606080825260208201529081016000815260200160008152600060209091015290565b60008281526002602052604090819020815160a0810190925280548290829061212f90614258565b80601f016020809104026020016040519081016040528092919081815260200182805461215b90614258565b80156121a85780601f1061217d576101008083540402835291602001916121a8565b820191906000526020600020905b81548152906001019060200180831161218b57829003601f168201915b505050505081526020016001820180546121c190614258565b80601f01602080910402602001604051908101604052809291908181526020018280546121ed90614258565b801561223a5780601f1061220f5761010080835404028352916020019161223a565b820191906000526020600020905b81548152906001019060200180831161221d57829003601f168201915b5050509183525050600282015460209091019060ff16600381111561226157612261613deb565b600381111561227257612272613deb565b81526020016002820160019054906101000a900460ff16600181111561229a5761229a613deb565b60018111156122ab576122ab613deb565b81526002919091015462010000900473ffffffffffffffffffffffffffffffffffffffff1660209091015292915050565b6122e461292b565b6122ef600382612910565b612328576040517fe181733f0000000000000000000000000000000000000000000000000000000081526004810182905260240161047c565b612333600582612910565b1561236d576040517ff7d7a2940000000000000000000000000000000000000000000000000000000081526004810182905260240161047c565b6123786005826129ae565b5060405181907fdcea1b78b6ddc31592a94607d537543fcaafda6cc52d6d5cc7bbfca1422baf2190600090a250565b63ffffffff808316600090815260116020908152604080832080546401000000009004909416808452600190940182528083208584526003019091529020805460609291906123f590614258565b80601f016020809104026020016040519081016040528092919081815260200182805461242190614258565b801561246e5780601f106124435761010080835404028352916020019161246e565b820191906000526020600020905b81548152906001019060200180831161245157829003601f168201915b505050505091505092915050565b6060600061248a6003613333565b905060006124986005612f9c565b82516124a491906149e7565b67ffffffffffffffff8111156124bc576124bc6140c2565b60405190808252806020026020018201604052801561251d57816020015b61250a6040805160a081018252606080825260208201529081016000815260200160008152600060209091015290565b8152602001906001900390816124da5790505b5090506000805b8351811015610e1857600084828151811061254157612541614055565b6020026020010151905061255f81600561291090919063ffffffff16565b6127605760008181526002602052604090819020815160a0810190925280548290829061258b90614258565b80601f01602080910402602001604051908101604052809291908181526020018280546125b790614258565b80156126045780601f106125d957610100808354040283529160200191612604565b820191906000526020600020905b8154815290600101906020018083116125e757829003601f168201915b5050505050815260200160018201805461261d90614258565b80601f016020809104026020016040519081016040528092919081815260200182805461264990614258565b80156126965780601f1061266b57610100808354040283529160200191612696565b820191906000526020600020905b81548152906001019060200180831161267957829003601f168201915b5050509183525050600282015460209091019060ff1660038111156126bd576126bd613deb565b60038111156126ce576126ce613deb565b81526020016002820160019054906101000a900460ff1660018111156126f6576126f6613deb565b600181111561270757612707613deb565b81526002919091015462010000900473ffffffffffffffffffffffffffffffffffffffff16602090910152845185908590811061274657612746614055565b6020026020010181905250828061275c906144fc565b9350505b5061276a816144fc565b9050612524565b6060806000612780600d613333565b90506000815167ffffffffffffffff81111561279e5761279e6140c2565b60405190808252806020026020018201604052801561280d57816020015b60408051608081018252600080825260208083018290529282015260608082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816127bc5790505b5090506000825167ffffffffffffffff81111561282c5761282c6140c2565b604051908082528060200260200182016040528015612855578160200160208202803683370190505b50905060005b83518110156128f157600084828151811061287857612878614055565b6020026020010151905060008061288e8361328e565b91509150818685815181106128a5576128a5614055565b6020026020010181905250808585815181106128c3576128c3614055565b602002602001019063ffffffff16908163ffffffff1681525050505050806128ea906144fc565b905061285b565b509094909350915050565b61290461292b565b61290d81613340565b50565b600081815260018301602052604081205415155b9392505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146129ac576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161047c565b565b60006129248383613435565b805163ffffffff90811660009081526011602090815260408083208286015185168452600101909152902060808301519091161580612a0e575060808201518590612a06906001614b14565b63ffffffff16115b15612a5a5760808201516040517fd5f5269100000000000000000000000000000000000000000000000000000000815263ffffffff90911660048201526024810186905260440161047c565b60005b85811015612b2257612a8a878783818110612a7a57612a7a614055565b8592602090910201359050612910565b15612aec578251878783818110612aa357612aa3614055565b6040517f636e405700000000000000000000000000000000000000000000000000000000815263ffffffff9094166004850152602002919091013560248301525060440161047c565b612b11878783818110612b0157612b01614055565b85926020909102013590506129ae565b50612b1b816144fc565b9050612a5d565b5060005b83811015612e0e5736858583818110612b4157612b41614055565b9050602002810190612b539190614084565b9050612b6160038235612910565b612b9a576040517fe181733f0000000000000000000000000000000000000000000000000000000081528135600482015260240161047c565b612ba660058235612910565b15612be0576040517ff7d7a2940000000000000000000000000000000000000000000000000000000081528135600482015260240161047c565b8035600090815260038401602052604081208054612bfd90614258565b90501115612c495783516040517f3927d08000000000000000000000000000000000000000000000000000000000815263ffffffff90911660048201528135602482015260440161047c565b60005b87811015612d5b57612cf08235601060008c8c86818110612c6f57612c6f614055565b9050602002013581526020019081526020016000206003016000601060008e8e88818110612c9f57612c9f614055565b90506020020135815260200190815260200160002060000160049054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002061291090919063ffffffff16565b612d4b57888882818110612d0657612d06614055565b6040517fa7e79250000000000000000000000000000000000000000000000000000000008152602090910292909201356004830152508235602482015260440161047c565b612d54816144fc565b9050612c4c565b5060028301805460018101825560009182526020918290208335910155612d84908201826145af565b82356000908152600386016020526040902091612da2919083614631565b508351602080860151612dfd92918435908c908c90612dc3908801886145af565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061348492505050565b50612e07816144fc565b9050612b26565b50604080830151835163ffffffff90811660009081526011602090815284822080549415156c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff90951694909417909355606086015186518316825284822080549115156d0100000000000000000000000000027fffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffff9092169190911790556080860151865183168252848220805491841668010000000000000000027fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff909216919091179055918501805186518316845292849020805493909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9093169290921790558351905191517ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c158170365192612f8c929163ffffffff92831681529116602082015260400190565b60405180910390a1505050505050565b6000610794825490565b6040805160e0810182526000808252602080830182905282840182905260608084018390526080840183905260a0840181905260c084015263ffffffff85811683526011825284832080546401000000009004909116808452600190910182528483206002810180548751818602810186019098528088529596929591949390919083018282801561305757602002820191906000526020600020905b815481526020019060010190808311613043575b505050505090506000815167ffffffffffffffff81111561307a5761307a6140c2565b6040519080825280602002602001820160405280156130c057816020015b6040805180820190915260008152606060208201528152602001906001900390816130985790505b50905060005b81518110156131e15760405180604001604052808483815181106130ec576130ec614055565b6020026020010151815260200185600301600086858151811061311157613111614055565b60200260200101518152602001908152602001600020805461313290614258565b80601f016020809104026020016040519081016040528092919081815260200182805461315e90614258565b80156131ab5780601f10613180576101008083540402835291602001916131ab565b820191906000526020600020905b81548152906001019060200180831161318e57829003601f168201915b50505050508152508282815181106131c5576131c5614055565b6020026020010181905250806131da906144fc565b90506130c6565b506040805160e08101825263ffffffff8089166000818152601160208181528683205480861687528b8616828801526801000000000000000081049095169686019690965260ff6c010000000000000000000000008504811615156060870152929091529093526d010000000000000000000000000090049091161515608082015260a0810161327085613333565b81526020019190915295945050505050565b60006129248383613565565b604080516080810182526000808252602082018190529181019190915260608082015260408051608081018252600084815260106020908152838220805463ffffffff80821686526001830154848701526002830154868801526401000000009091041683526003019052918220606082019061330a90613333565b905260009384526010602052604090932054929364010000000090930463ffffffff1692915050565b606060006129248361365f565b3373ffffffffffffffffffffffffffffffffffffffff8216036133bf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161047c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600081815260018301602052604081205461347c57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610794565b506000610794565b6000848152600260208190526040909120015462010000900473ffffffffffffffffffffffffffffffffffffffff161561074657600084815260026020819052604091829020015490517ffba64a7c0000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063fba64a7c9061352b908690869086908b908d90600401614b31565b600060405180830381600087803b15801561354557600080fd5b505af1158015613559573d6000803e3d6000fd5b50505050505050505050565b6000818152600183016020526040812054801561364e5760006135896001836149e7565b855490915060009061359d906001906149e7565b90508181146136025760008660000182815481106135bd576135bd614055565b90600052602060002001549050808760000184815481106135e0576135e0614055565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061361357613613614bb3565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610794565b6000915050610794565b5092915050565b6060816000018054806020026020016040519081016040528092919081815260200182805480156136af57602002820191906000526020600020905b81548152602001906001019080831161369b575b50505050509050919050565b5080546136c790614258565b6000825580601f106136d7575050565b601f01602090049060005260206000209081019061290d91905b8082111561370557600081556001016136f1565b5090565b60008083601f84011261371b57600080fd5b50813567ffffffffffffffff81111561373357600080fd5b6020830191508360208260051b850101111561374e57600080fd5b9250929050565b6000806000806040858703121561376b57600080fd5b843567ffffffffffffffff8082111561378357600080fd5b61378f88838901613709565b909650945060208701359150808211156137a857600080fd5b506137b587828801613709565b95989497509550505050565b60008083601f8401126137d357600080fd5b50813567ffffffffffffffff8111156137eb57600080fd5b60208301915083602082850101111561374e57600080fd5b6000806000806040858703121561381957600080fd5b843567ffffffffffffffff8082111561383157600080fd5b61383d888389016137c1565b9096509450602087013591508082111561385657600080fd5b506137b5878288016137c1565b60006020828403121561387557600080fd5b5035919050565b60006020828403121561388e57600080fd5b813567ffffffffffffffff8111156138a557600080fd5b820160a0818503121561292457600080fd5b6000815180845260005b818110156138dd576020818501810151868301820152016138c1565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061292460208301846138b7565b801515811461290d57600080fd5b803563ffffffff8116811461395057600080fd5b919050565b600080600080600080600060a0888a03121561397057600080fd5b873567ffffffffffffffff8082111561398857600080fd5b6139948b838c01613709565b909950975060208a01359150808211156139ad57600080fd5b506139ba8a828b01613709565b90965094505060408801356139ce8161392e565b925060608801356139de8161392e565b91506139ec6080890161393c565b905092959891949750929550565b60008060208385031215613a0d57600080fd5b823567ffffffffffffffff811115613a2457600080fd5b613a3085828601613709565b90969095509350505050565b600081518084526020808501945080840160005b83811015613a6c57815187529582019590820190600101613a50565b509495945050505050565b600063ffffffff80835116845260208181850151168186015260408281860151168187015260608501511515606087015260808501511515608087015260a0850151925060e060a0870152613acf60e0870184613a3c565b925060c085015186840360c08801528381518086528486019150848160051b870101858401935060005b82811015613b56578782037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018452845180518352870151878301879052613b43878401826138b7565b9588019594880194925050600101613af9565b509998505050505050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015613bd7577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452613bc5858351613a77565b94509285019290850190600101613b8b565b5092979650505050505050565b600060208284031215613bf657600080fd5b6129248261393c565b6020815260006129246020830184613a77565b63ffffffff815116825260208101516020830152604081015160408301526000606082015160806060850152613c4b6080850182613a3c565b949350505050565b604081526000613c666040830185613c12565b905063ffffffff831660208301529392505050565b73ffffffffffffffffffffffffffffffffffffffff81511682526000602082015160406020850152613c4b60408501826138b7565b6020815260006129246020830184613c7b565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015613bd7577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452613d24858351613c7b565b94509285019290850190600101613cea565b60008060008060008060008060c0898b031215613d5257600080fd5b613d5b8961393c565b9750602089013567ffffffffffffffff80821115613d7857600080fd5b613d848c838d01613709565b909950975060408b0135915080821115613d9d57600080fd5b50613daa8b828c01613709565b9096509450506060890135613dbe8161392e565b92506080890135613dce8161392e565b9150613ddc60a08a0161393c565b90509295985092959890939650565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6000815160a08452613e2f60a08501826138b7565b905060208301518482036020860152613e4882826138b7565b915050604083015160048110613e6057613e60613deb565b6040850152606083015160028110613e7a57613e7a613deb565b606085015260809283015173ffffffffffffffffffffffffffffffffffffffff1692909301919091525090565b6020815260006129246020830184613e1a565b60008060408385031215613ecd57600080fd5b613ed68361393c565b946020939093013593505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015613bd7577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452613f45858351613e1a565b94509285019290850190600101613f0b565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015613fcc577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552613fba868351613c12565b95509382019390820190600101613f80565b50508584038187015286518085528782019482019350915060005b8281101561400957845163ffffffff1684529381019392810192600101613fe7565b5091979650505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461290d57600080fd5b60006020828403121561404a57600080fd5b813561292481614016565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc18336030181126140b857600080fd5b9190910192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715614114576141146140c2565b60405290565b6040516080810167ffffffffffffffff81118282101715614114576141146140c2565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715614184576141846140c2565b604052919050565b60006040823603121561419e57600080fd5b6141a66140f1565b82356141b181614016565b815260208381013567ffffffffffffffff808211156141cf57600080fd5b9085019036601f8301126141e257600080fd5b8135818111156141f4576141f46140c2565b614224847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161413d565b9150808252368482850101111561423a57600080fd5b80848401858401376000908201840152918301919091525092915050565b600181811c9082168061426c57607f821691505b6020821081036142a5577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60006020808352600084546142bf81614258565b808487015260406001808416600081146142e0576001811461431857614346565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838a01528284151560051b8a01019550614346565b896000528660002060005b8581101561433e5781548b8201860152908301908801614323565b8a0184019650505b509398975050505050505050565b601f821115610b7c57600081815260208120601f850160051c8101602086101561437b5750805b601f850160051c820191505b8181101561074657828155600101614387565b815167ffffffffffffffff8111156143b4576143b46140c2565b6143c8816143c28454614258565b84614354565b602080601f83116001811461441b57600084156143e55750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610746565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561446857888601518255948401946001909101908401614449565b50858210156144a457878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b828152604060208201526000613c4b60408301846138b7565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361452d5761452d6144cd565b5060010190565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b604081526000614591604083018688614534565b82810360208401526145a4818587614534565b979650505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126145e457600080fd5b83018035915067ffffffffffffffff8211156145ff57600080fd5b60200191503681900382131561374e57600080fd5b60006020828403121561462657600080fd5b81516129248161392e565b67ffffffffffffffff831115614649576146496140c2565b61465d836146578354614258565b83614354565b6000601f8411600181146146af57600085156146795750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355614745565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156146fe57868501358255602094850194600190920191016146de565b5086821015614739577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b600081356004811061079457600080fd5b6004821061476d5761476d613deb565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0081541660ff831681178255505050565b600081356002811061079457600080fd5b600282106147bf576147bf613deb565b805461ff008360081b167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff8216178255505050565b6000813561079481614016565b61480b82836145af565b67ffffffffffffffff811115614823576148236140c2565b614837816148318554614258565b85614354565b6000601f82116001811461488957600083156148535750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600385901b1c1916600184901b17855561491f565b6000858152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0841690835b828110156148d857868501358255602094850194600190920191016148b8565b5084821015614913577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88660031b161c19848701351681555b505060018360011b0185555b5050505061493060208301836145af565b61493e818360018601614631565b5050600281016149596149536040850161474c565b8261475d565b61496e6149686060850161479e565b826147af565b610b7c61497d608085016147f4565b82547fffffffffffffffffffff0000000000000000000000000000000000000000ffff1660109190911b75ffffffffffffffffffffffffffffffffffffffff000016178255565b600063ffffffff8083168181036149dd576149dd6144cd565b6001019392505050565b81810381811115610794576107946144cd565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff818336030181126140b857600080fd5b600060808236031215614a4057600080fd5b614a4861411a565b614a518361393c565b81526020808401358183015260408401356040830152606084013567ffffffffffffffff80821115614a8257600080fd5b9085019036601f830112614a9557600080fd5b813581811115614aa757614aa76140c2565b8060051b9150614ab884830161413d565b8181529183018401918481019036841115614ad257600080fd5b938501935b83851015614af057843582529385019390850190614ad7565b606087015250939695505050505050565b6020815260006129246020830184613a3c565b63ffffffff818116838216019080821115613658576136586144cd565b6080815284608082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff861115614b6a57600080fd5b8560051b808860a0850137820182810360a09081016020850152614b90908201876138b7565b91505063ffffffff80851660408401528084166060840152509695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
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

func (_CapabilityRegistry *CapabilityRegistryCaller) GetDON(opts *bind.CallOpts, donId uint32) (CapabilityRegistryDONInfo, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getDON", donId)

	if err != nil {
		return *new(CapabilityRegistryDONInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(CapabilityRegistryDONInfo)).(*CapabilityRegistryDONInfo)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetDON(donId uint32) (CapabilityRegistryDONInfo, error) {
	return _CapabilityRegistry.Contract.GetDON(&_CapabilityRegistry.CallOpts, donId)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetDON(donId uint32) (CapabilityRegistryDONInfo, error) {
	return _CapabilityRegistry.Contract.GetDON(&_CapabilityRegistry.CallOpts, donId)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) GetDONCapabilityConfig(opts *bind.CallOpts, donId uint32, capabilityId [32]byte) ([]byte, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getDONCapabilityConfig", donId, capabilityId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetDONCapabilityConfig(donId uint32, capabilityId [32]byte) ([]byte, error) {
	return _CapabilityRegistry.Contract.GetDONCapabilityConfig(&_CapabilityRegistry.CallOpts, donId, capabilityId)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetDONCapabilityConfig(donId uint32, capabilityId [32]byte) ([]byte, error) {
	return _CapabilityRegistry.Contract.GetDONCapabilityConfig(&_CapabilityRegistry.CallOpts, donId, capabilityId)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) GetDONs(opts *bind.CallOpts) ([]CapabilityRegistryDONInfo, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getDONs")

	if err != nil {
		return *new([]CapabilityRegistryDONInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]CapabilityRegistryDONInfo)).(*[]CapabilityRegistryDONInfo)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetDONs() ([]CapabilityRegistryDONInfo, error) {
	return _CapabilityRegistry.Contract.GetDONs(&_CapabilityRegistry.CallOpts)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetDONs() ([]CapabilityRegistryDONInfo, error) {
	return _CapabilityRegistry.Contract.GetDONs(&_CapabilityRegistry.CallOpts)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) GetHashedCapabilityId(opts *bind.CallOpts, labelledName string, version string) ([32]byte, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getHashedCapabilityId", labelledName, version)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetHashedCapabilityId(labelledName string, version string) ([32]byte, error) {
	return _CapabilityRegistry.Contract.GetHashedCapabilityId(&_CapabilityRegistry.CallOpts, labelledName, version)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetHashedCapabilityId(labelledName string, version string) ([32]byte, error) {
	return _CapabilityRegistry.Contract.GetHashedCapabilityId(&_CapabilityRegistry.CallOpts, labelledName, version)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) GetNode(opts *bind.CallOpts, p2pId [32]byte) (CapabilityRegistryNodeInfo, uint32, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getNode", p2pId)

	if err != nil {
		return *new(CapabilityRegistryNodeInfo), *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(CapabilityRegistryNodeInfo)).(*CapabilityRegistryNodeInfo)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return out0, out1, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetNode(p2pId [32]byte) (CapabilityRegistryNodeInfo, uint32, error) {
	return _CapabilityRegistry.Contract.GetNode(&_CapabilityRegistry.CallOpts, p2pId)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetNode(p2pId [32]byte) (CapabilityRegistryNodeInfo, uint32, error) {
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

func (_CapabilityRegistry *CapabilityRegistryCaller) GetNodeOperators(opts *bind.CallOpts) ([]CapabilityRegistryNodeOperator, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getNodeOperators")

	if err != nil {
		return *new([]CapabilityRegistryNodeOperator), err
	}

	out0 := *abi.ConvertType(out[0], new([]CapabilityRegistryNodeOperator)).(*[]CapabilityRegistryNodeOperator)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetNodeOperators() ([]CapabilityRegistryNodeOperator, error) {
	return _CapabilityRegistry.Contract.GetNodeOperators(&_CapabilityRegistry.CallOpts)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetNodeOperators() ([]CapabilityRegistryNodeOperator, error) {
	return _CapabilityRegistry.Contract.GetNodeOperators(&_CapabilityRegistry.CallOpts)
}

func (_CapabilityRegistry *CapabilityRegistryCaller) GetNodes(opts *bind.CallOpts) ([]CapabilityRegistryNodeInfo, []uint32, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getNodes")

	if err != nil {
		return *new([]CapabilityRegistryNodeInfo), *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([]CapabilityRegistryNodeInfo)).(*[]CapabilityRegistryNodeInfo)
	out1 := *abi.ConvertType(out[1], new([]uint32)).(*[]uint32)

	return out0, out1, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetNodes() ([]CapabilityRegistryNodeInfo, []uint32, error) {
	return _CapabilityRegistry.Contract.GetNodes(&_CapabilityRegistry.CallOpts)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetNodes() ([]CapabilityRegistryNodeInfo, []uint32, error) {
	return _CapabilityRegistry.Contract.GetNodes(&_CapabilityRegistry.CallOpts)
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) AddDON(opts *bind.TransactOpts, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint32) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "addDON", nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilityRegistry *CapabilityRegistrySession) AddDON(nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint32) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddDON(&_CapabilityRegistry.TransactOpts, nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) AddDON(nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint32) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddDON(&_CapabilityRegistry.TransactOpts, nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) AddNodes(opts *bind.TransactOpts, nodes []CapabilityRegistryNodeInfo) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "addNodes", nodes)
}

func (_CapabilityRegistry *CapabilityRegistrySession) AddNodes(nodes []CapabilityRegistryNodeInfo) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddNodes(&_CapabilityRegistry.TransactOpts, nodes)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) AddNodes(nodes []CapabilityRegistryNodeInfo) (*types.Transaction, error) {
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) RemoveDONs(opts *bind.TransactOpts, donIds []uint32) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "removeDONs", donIds)
}

func (_CapabilityRegistry *CapabilityRegistrySession) RemoveDONs(donIds []uint32) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.RemoveDONs(&_CapabilityRegistry.TransactOpts, donIds)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) RemoveDONs(donIds []uint32) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.RemoveDONs(&_CapabilityRegistry.TransactOpts, donIds)
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) RemoveNodes(opts *bind.TransactOpts, removedNodeP2PIds [][32]byte) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "removeNodes", removedNodeP2PIds)
}

func (_CapabilityRegistry *CapabilityRegistrySession) RemoveNodes(removedNodeP2PIds [][32]byte) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.RemoveNodes(&_CapabilityRegistry.TransactOpts, removedNodeP2PIds)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) RemoveNodes(removedNodeP2PIds [][32]byte) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.RemoveNodes(&_CapabilityRegistry.TransactOpts, removedNodeP2PIds)
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) UpdateDON(opts *bind.TransactOpts, donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint32) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "updateDON", donId, nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilityRegistry *CapabilityRegistrySession) UpdateDON(donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint32) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.UpdateDON(&_CapabilityRegistry.TransactOpts, donId, nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) UpdateDON(donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint32) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.UpdateDON(&_CapabilityRegistry.TransactOpts, donId, nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) UpdateNodes(opts *bind.TransactOpts, nodes []CapabilityRegistryNodeInfo) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "updateNodes", nodes)
}

func (_CapabilityRegistry *CapabilityRegistrySession) UpdateNodes(nodes []CapabilityRegistryNodeInfo) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.UpdateNodes(&_CapabilityRegistry.TransactOpts, nodes)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) UpdateNodes(nodes []CapabilityRegistryNodeInfo) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.UpdateNodes(&_CapabilityRegistry.TransactOpts, nodes)
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

type CapabilityRegistryConfigSetIterator struct {
	Event *CapabilityRegistryConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilityRegistryConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilityRegistryConfigSet)
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
		it.Event = new(CapabilityRegistryConfigSet)
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

func (it *CapabilityRegistryConfigSetIterator) Error() error {
	return it.fail
}

func (it *CapabilityRegistryConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilityRegistryConfigSet struct {
	DonId       uint32
	ConfigCount uint32
	Raw         types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterConfigSet(opts *bind.FilterOpts) (*CapabilityRegistryConfigSetIterator, error) {

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryConfigSetIterator{contract: _CapabilityRegistry.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryConfigSet) (event.Subscription, error) {

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilityRegistryConfigSet)
				if err := _CapabilityRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_CapabilityRegistry *CapabilityRegistryFilterer) ParseConfigSet(log types.Log) (*CapabilityRegistryConfigSet, error) {
	event := new(CapabilityRegistryConfigSet)
	if err := _CapabilityRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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
	Signer         [32]byte
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

type CapabilityRegistryNodeRemovedIterator struct {
	Event *CapabilityRegistryNodeRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilityRegistryNodeRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilityRegistryNodeRemoved)
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
		it.Event = new(CapabilityRegistryNodeRemoved)
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

func (it *CapabilityRegistryNodeRemovedIterator) Error() error {
	return it.fail
}

func (it *CapabilityRegistryNodeRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilityRegistryNodeRemoved struct {
	P2pId [32]byte
	Raw   types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterNodeRemoved(opts *bind.FilterOpts) (*CapabilityRegistryNodeRemovedIterator, error) {

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "NodeRemoved")
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryNodeRemovedIterator{contract: _CapabilityRegistry.contract, event: "NodeRemoved", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchNodeRemoved(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeRemoved) (event.Subscription, error) {

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "NodeRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilityRegistryNodeRemoved)
				if err := _CapabilityRegistry.contract.UnpackLog(event, "NodeRemoved", log); err != nil {
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

func (_CapabilityRegistry *CapabilityRegistryFilterer) ParseNodeRemoved(log types.Log) (*CapabilityRegistryNodeRemoved, error) {
	event := new(CapabilityRegistryNodeRemoved)
	if err := _CapabilityRegistry.contract.UnpackLog(event, "NodeRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilityRegistryNodeUpdatedIterator struct {
	Event *CapabilityRegistryNodeUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilityRegistryNodeUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilityRegistryNodeUpdated)
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
		it.Event = new(CapabilityRegistryNodeUpdated)
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

func (it *CapabilityRegistryNodeUpdatedIterator) Error() error {
	return it.fail
}

func (it *CapabilityRegistryNodeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilityRegistryNodeUpdated struct {
	P2pId          [32]byte
	NodeOperatorId *big.Int
	Signer         [32]byte
	Raw            types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterNodeUpdated(opts *bind.FilterOpts) (*CapabilityRegistryNodeUpdatedIterator, error) {

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "NodeUpdated")
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryNodeUpdatedIterator{contract: _CapabilityRegistry.contract, event: "NodeUpdated", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchNodeUpdated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeUpdated) (event.Subscription, error) {

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "NodeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilityRegistryNodeUpdated)
				if err := _CapabilityRegistry.contract.UnpackLog(event, "NodeUpdated", log); err != nil {
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

func (_CapabilityRegistry *CapabilityRegistryFilterer) ParseNodeUpdated(log types.Log) (*CapabilityRegistryNodeUpdated, error) {
	event := new(CapabilityRegistryNodeUpdated)
	if err := _CapabilityRegistry.contract.UnpackLog(event, "NodeUpdated", log); err != nil {
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
	case _CapabilityRegistry.abi.Events["ConfigSet"].ID:
		return _CapabilityRegistry.ParseConfigSet(log)
	case _CapabilityRegistry.abi.Events["NodeAdded"].ID:
		return _CapabilityRegistry.ParseNodeAdded(log)
	case _CapabilityRegistry.abi.Events["NodeOperatorAdded"].ID:
		return _CapabilityRegistry.ParseNodeOperatorAdded(log)
	case _CapabilityRegistry.abi.Events["NodeOperatorRemoved"].ID:
		return _CapabilityRegistry.ParseNodeOperatorRemoved(log)
	case _CapabilityRegistry.abi.Events["NodeOperatorUpdated"].ID:
		return _CapabilityRegistry.ParseNodeOperatorUpdated(log)
	case _CapabilityRegistry.abi.Events["NodeRemoved"].ID:
		return _CapabilityRegistry.ParseNodeRemoved(log)
	case _CapabilityRegistry.abi.Events["NodeUpdated"].ID:
		return _CapabilityRegistry.ParseNodeUpdated(log)
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

func (CapabilityRegistryConfigSet) Topic() common.Hash {
	return common.HexToHash("0xf264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c1581703651")
}

func (CapabilityRegistryNodeAdded) Topic() common.Hash {
	return common.HexToHash("0xc9296aa9b0951d8000e8ed7f2b5be30c5106de8df3dbedf9a57c93f5f9e4d7da")
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

func (CapabilityRegistryNodeRemoved) Topic() common.Hash {
	return common.HexToHash("0x5254e609a97bab37b7cc79fe128f85c097bd6015c6e1624ae0ba392eb9753205")
}

func (CapabilityRegistryNodeUpdated) Topic() common.Hash {
	return common.HexToHash("0xf101cfc54994c31624d25789378d71ec4dbdc533e26a4ecc6b7648f4798d0916")
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

	GetDON(opts *bind.CallOpts, donId uint32) (CapabilityRegistryDONInfo, error)

	GetDONCapabilityConfig(opts *bind.CallOpts, donId uint32, capabilityId [32]byte) ([]byte, error)

	GetDONs(opts *bind.CallOpts) ([]CapabilityRegistryDONInfo, error)

	GetHashedCapabilityId(opts *bind.CallOpts, labelledName string, version string) ([32]byte, error)

	GetNode(opts *bind.CallOpts, p2pId [32]byte) (CapabilityRegistryNodeInfo, uint32, error)

	GetNodeOperator(opts *bind.CallOpts, nodeOperatorId *big.Int) (CapabilityRegistryNodeOperator, error)

	GetNodeOperators(opts *bind.CallOpts) ([]CapabilityRegistryNodeOperator, error)

	GetNodes(opts *bind.CallOpts) ([]CapabilityRegistryNodeInfo, []uint32, error)

	IsCapabilityDeprecated(opts *bind.CallOpts, hashedCapabilityId [32]byte) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddCapability(opts *bind.TransactOpts, capability CapabilityRegistryCapability) (*types.Transaction, error)

	AddDON(opts *bind.TransactOpts, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint32) (*types.Transaction, error)

	AddNodeOperators(opts *bind.TransactOpts, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error)

	AddNodes(opts *bind.TransactOpts, nodes []CapabilityRegistryNodeInfo) (*types.Transaction, error)

	DeprecateCapability(opts *bind.TransactOpts, hashedCapabilityId [32]byte) (*types.Transaction, error)

	RemoveDONs(opts *bind.TransactOpts, donIds []uint32) (*types.Transaction, error)

	RemoveNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []*big.Int) (*types.Transaction, error)

	RemoveNodes(opts *bind.TransactOpts, removedNodeP2PIds [][32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateDON(opts *bind.TransactOpts, donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint32) (*types.Transaction, error)

	UpdateNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []*big.Int, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error)

	UpdateNodes(opts *bind.TransactOpts, nodes []CapabilityRegistryNodeInfo) (*types.Transaction, error)

	FilterCapabilityAdded(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilityRegistryCapabilityAddedIterator, error)

	WatchCapabilityAdded(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryCapabilityAdded, hashedCapabilityId [][32]byte) (event.Subscription, error)

	ParseCapabilityAdded(log types.Log) (*CapabilityRegistryCapabilityAdded, error)

	FilterCapabilityDeprecated(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilityRegistryCapabilityDeprecatedIterator, error)

	WatchCapabilityDeprecated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryCapabilityDeprecated, hashedCapabilityId [][32]byte) (event.Subscription, error)

	ParseCapabilityDeprecated(log types.Log) (*CapabilityRegistryCapabilityDeprecated, error)

	FilterConfigSet(opts *bind.FilterOpts) (*CapabilityRegistryConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*CapabilityRegistryConfigSet, error)

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

	FilterNodeRemoved(opts *bind.FilterOpts) (*CapabilityRegistryNodeRemovedIterator, error)

	WatchNodeRemoved(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeRemoved) (event.Subscription, error)

	ParseNodeRemoved(log types.Log) (*CapabilityRegistryNodeRemoved, error)

	FilterNodeUpdated(opts *bind.FilterOpts) (*CapabilityRegistryNodeUpdatedIterator, error)

	WatchNodeUpdated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeUpdated) (event.Subscription, error)

	ParseNodeUpdated(log types.Log) (*CapabilityRegistryNodeUpdated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CapabilityRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*CapabilityRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CapabilityRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*CapabilityRegistryOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
