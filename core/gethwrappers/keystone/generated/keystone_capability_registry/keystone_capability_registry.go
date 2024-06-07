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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityIsDeprecated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"DONDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONCapability\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONNode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedConfigurationContract\",\"type\":\"address\"}],\"name\":\"InvalidCapabilityConfigurationContractInterface\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"f\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"nodeCount\",\"type\":\"uint256\"}],\"name\":\"InvalidFaultTolerance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"InvalidNodeCapabilities\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeOperatorAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"InvalidNodeP2PId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lengthOne\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lengthTwo\",\"type\":\"uint256\"}],\"name\":\"LengthMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodeAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodeDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"NodeDoesNotSupportCapability\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"NodeOperatorDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodePartOfDON\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDeprecated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"NodeOperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"NodeRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability[]\",\"name\":\"capabilities\",\"type\":\"tuple[]\"}],\"name\":\"addCapabilities\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"f\",\"type\":\"uint32\"}],\"name\":\"addDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"addNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"addNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"deprecateCapabilities\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCapabilities\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedId\",\"type\":\"bytes32\"}],\"name\":\"getCapability\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"getCapabilityConfigs\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"getDON\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"f\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilityRegistry.DONInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"f\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilityRegistry.DONInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"}],\"name\":\"getHashedCapabilityId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"getNode\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"getNodeOperator\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeOperators\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodes\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo[]\",\"name\":\"nodeInfo\",\"type\":\"tuple[]\"},{\"internalType\":\"uint32[]\",\"name\":\"configCounts\",\"type\":\"uint32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"isCapabilityDeprecated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"donIds\",\"type\":\"uint32[]\"}],\"name\":\"removeDONs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint32[]\"}],\"name\":\"removeNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"removedNodeP2PIds\",\"type\":\"bytes32[]\"}],\"name\":\"removeNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"f\",\"type\":\"uint32\"}],\"name\":\"updateDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint32[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"updateNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"updateNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600e80546001600160401b0319166401000000011790553480156200002857600080fd5b503380600081620000805760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000b357620000b381620000bc565b50505062000167565b336001600160a01b03821603620001165760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000077565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b614e6f80620001776000396000f3fe608060405234801561001057600080fd5b50600436106101ae5760003560e01c80635d83d967116100ee57806386fa4246116100975780639cb7c5f4116100715780639cb7c5f4146103ea578063ddbe4f821461040a578063e29581aa14610420578063f2fde38b1461043657600080fd5b806386fa42461461039c5780638da5cb5b146103af57806395864d1f146103d757600080fd5b8063715f5295116100c8578063715f52951461036e57806373ac22b41461038157806379ba50971461039457600080fd5b80635d83d967146103335780635e65e3091461034657806366acaa331461035957600080fd5b806322bdbcbc1161015b5780632c01a1e8116101355780632c01a1e8146102cb578063398f3773146102de5780633f2a13c9146102f157806350c946fe1461031257600080fd5b806322bdbcbc146102785780632353740514610298578063275459f2146102b857600080fd5b8063193ec0061161018c578063193ec0061461023b5780631d05394c14610250578063214502431461026357600080fd5b80630fe5800a146101b357806312570011146101d9578063181f5a77146101fc575b600080fd5b6101c66101c1366004613c70565b610449565b6040519081526020015b60405180910390f35b6101ec6101e7366004613cd4565b61047d565b60405190151581526020016101d0565b604080518082018252601881527f4361706162696c697479526567697374727920312e302e300000000000000000602082015290516101d09190613d5b565b61024e610249366004613dda565b61048a565b005b61024e61025e366004613e7f565b610554565b61026b610710565b6040516101d09190613fe9565b61028b610286366004614069565b610875565b6040516101d091906140c1565b6102ab6102a6366004614069565b610962565b6040516101d091906140d4565b61024e6102c6366004613e7f565b6109a6565b61024e6102d9366004613e7f565b610a7d565b61024e6102ec366004613e7f565b610ca0565b6103046102ff3660046140e7565b610e68565b6040516101d0929190614111565b610325610320366004613cd4565b611054565b6040516101d0929190614178565b61024e610341366004613e7f565b6110f9565b61024e610354366004613e7f565b611247565b6103616116b9565b6040516101d091906141a0565b61024e61037c366004613e7f565b6118a4565b61024e61038f366004613e7f565b61195f565b61024e611e04565b61024e6103aa366004614213565b611f01565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101d0565b61024e6103e536600461427f565b612247565b6103fd6103f8366004613cd4565b612313565b6040516101d091906143f0565b610412612518565b6040516101d0929190614403565b610428612883565b6040516101d0929190614493565b61024e610444366004614576565b6129f2565b6000828260405160200161045e929190614111565b6040516020818303038152906040528051906020012090505b92915050565b6000610477600583612a06565b610492612a21565b600e805460009164010000000090910463ffffffff169060046104b4836145c0565b82546101009290920a63ffffffff8181021990931691831602179091558181166000818152600d602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001684179055815160a081018352928352600190830152871515908201528515156060820152908416608082015290915061054a908990899089908990612aa4565b5050505050505050565b61055c612a21565b60005b8181101561070b57600083838381811061057b5761057b6145e3565b90506020020160208101906105909190614069565b63ffffffff8082166000908152600d60209081526040808320805464010000000090049094168084526001850190925282209394509192905b6105d282613154565b81101561062a576106198563ffffffff16600c60006105fa858761315e90919063ffffffff16565b815260200190815260200160002060040161316a90919063ffffffff16565b5061062381614612565b90506105c9565b508254640100000000900463ffffffff16600003610681576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff851660048201526024015b60405180910390fd5b63ffffffff84166000818152600d6020908152604080832080547fffffffffffffffffffffffffffffffffffff00000000000000000000000000001690558051938452908301919091527ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c1581703651910160405180910390a1505050508061070490614612565b905061055f565b505050565b600e54606090640100000000900463ffffffff16600061073160018361464a565b63ffffffff1667ffffffffffffffff81111561074f5761074f613b0a565b6040519080825280602002602001820160405280156107d657816020015b6040805160e081018252600080825260208083018290529282018190526060808301829052608083019190915260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90920191018161076d5790505b509050600060015b8363ffffffff168163ffffffff1610156108525763ffffffff8082166000908152600d602052604090205416156108425761081881613176565b83838151811061082a5761082a6145e3565b60200260200101819052508161083f90614612565b91505b61084b816145c0565b90506107de565b5061085e60018461464a565b63ffffffff16811461086e578082525b5092915050565b60408051808201909152600081526060602082015263ffffffff82166000908152600b60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff16835260018101805491928401916108d990614667565b80601f016020809104026020016040519081016040528092919081815260200182805461090590614667565b80156109525780601f1061092757610100808354040283529160200191610952565b820191906000526020600020905b81548152906001019060200180831161093557829003601f168201915b5050505050815250509050919050565b6040805160e0810182526000808252602082018190529181018290526060808201839052608082019290925260a0810182905260c081019190915261047782613176565b6109ae612a21565b60005b63ffffffff811682111561070b57600083838363ffffffff168181106109d9576109d96145e3565b90506020020160208101906109ee9190614069565b63ffffffff81166000908152600b6020526040812080547fffffffffffffffffffffffff0000000000000000000000000000000000000000168155919250610a396001830182613a9d565b505060405163ffffffff8216907fa59268ca81d40429e65ccea5385b59cf2d3fc6519371dee92f8eb1dae5107a7a90600090a250610a76816145c0565b90506109b1565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610c9a576000848483818110610ab757610ab76145e3565b602090810292909201356000818152600c90935260409092206001810154929350919050610b14576040517fd82f6adb00000000000000000000000000000000000000000000000000000000815260048101839052602401610678565b6000610b2282600401613154565b1115610b5d576040517f34a4a3f600000000000000000000000000000000000000000000000000000000815260048101839052602401610678565b83158015610b975750805463ffffffff166000908152600b602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b15610bd0576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610678565b6001810154610be19060079061316a565b506002810154610bf39060099061316a565b506000828152600c6020526040812080547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016815560018101829055600281018290559060048201818181610c488282613ad7565b5050505050507f5254e609a97bab37b7cc79fe128f85c097bd6015c6e1624ae0ba392eb975320582604051610c7f91815260200190565b60405180910390a1505080610c9390614612565b9050610a9b565b50505050565b610ca8612a21565b60005b8181101561070b576000838383818110610cc757610cc76145e3565b9050602002810190610cd991906146ba565b610ce2906146f8565b805190915073ffffffffffffffffffffffffffffffffffffffff16610d33576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e54604080518082018252835173ffffffffffffffffffffffffffffffffffffffff908116825260208086015181840190815263ffffffff9095166000818152600b909252939020825181547fffffffffffffffffffffffff00000000000000000000000000000000000000001692169190911781559251919290916001820190610dbf90826147b0565b5050600e8054909150600090610dda9063ffffffff166145c0565b91906101000a81548163ffffffff021916908363ffffffff160217905550816000015173ffffffffffffffffffffffffffffffffffffffff168163ffffffff167f78e94ca80be2c30abc061b99e7eb8583b1254781734b1e3ce339abb57da2fe8e8460200151604051610e4d9190613d5b565b60405180910390a3505080610e6190614612565b9050610cab565b63ffffffff8083166000908152600d60209081526040808320805464010000000090049094168084526001909401825280832085845260030190915281208054606093849390929091610eba90614667565b80601f0160208091040260200160405190810160405280929190818152602001828054610ee690614667565b8015610f335780601f10610f0857610100808354040283529160200191610f33565b820191906000526020600020905b815481529060010190602001808311610f1657829003601f168201915b5050506000888152600260208190526040909120015492935060609262010000900473ffffffffffffffffffffffffffffffffffffffff16159150611046905057600086815260026020819052604091829020015490517f8318ed5d00000000000000000000000000000000000000000000000000000000815263ffffffff891660048201526201000090910473ffffffffffffffffffffffffffffffffffffffff1690638318ed5d90602401600060405180830381865afa158015610ffd573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261104391908101906148ca565b90505b9093509150505b9250929050565b6040805160808101825260008082526020820181905291810191909152606080820152604080516080810182526000848152600c6020908152838220805463ffffffff8082168652600183015484870152600283015486880152640100000000909104168352600301905291822060608201906110d090613452565b90526000938452600c602052604090932054929364010000000090930463ffffffff1692915050565b611101612a21565b60005b8181101561070b576000838383818110611120576111206145e3565b90506020020135905061113d816003612a0690919063ffffffff16565b611176576040517fe181733f00000000000000000000000000000000000000000000000000000000815260048101829052602401610678565b61118160058261345f565b6111ba576040517ff7d7a29400000000000000000000000000000000000000000000000000000000815260048101829052602401610678565b6000818152600260205260408120906111d38282613a9d565b6111e1600183016000613a9d565b5060020180547fffffffffffffffffffff0000000000000000000000000000000000000000000016905560405181907fdcea1b78b6ddc31592a94607d537543fcaafda6cc52d6d5cc7bbfca1422baf2190600090a25061124081614612565b9050611104565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610c9a576000848483818110611281576112816145e3565b90506020028101906112939190614938565b61129c9061496c565b805163ffffffff166000908152600b602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff1682526001810180549596509394919390928401916112f290614667565b80601f016020809104026020016040519081016040528092919081815260200182805461131e90614667565b801561136b5780601f106113405761010080835404028352916020019161136b565b820191906000526020600020905b81548152906001019060200180831161134e57829003601f168201915b50505050508152505090508315801561139b5750805173ffffffffffffffffffffffffffffffffffffffff163314155b156113d4576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610678565b6040808301516000908152600c6020522060018101546114285782604001516040517fd82f6adb00000000000000000000000000000000000000000000000000000000815260040161067891815260200190565b60208301511580611458575082602001518160010154141580156114585750602083015161145890600790612a06565b1561148f576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606083015180516000036114d157806040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106789190614a3f565b8154640100000000900463ffffffff168260046114ed836145c0565b82546101009290920a63ffffffff8181021990931691831602179091558354640100000000900416905060005b82518110156115e557611550838281518110611538576115386145e3565b60200260200101516003612a0690919063ffffffff16565b61158857826040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106789190614a3f565b6115d483828151811061159d5761159d6145e3565b60200260200101518560030160008563ffffffff1663ffffffff16815260200190815260200160002061345f90919063ffffffff16565b506115de81614612565b905061151a565b50845183547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff909116178355604085015160028401556001830154602086015181146116575761163c60078261316a565b506020860151600185018190556116559060079061345f565b505b85516040808801516020808a015183519283529082015263ffffffff909216917f4b5b465e22eea0c3d40c30e936643245b80d19b2dcf75788c0699fe8d8db645b910160405180910390a2505050505050806116b290614612565b9050611265565b600e5460609063ffffffff1660006116d260018361464a565b63ffffffff1667ffffffffffffffff8111156116f0576116f0613b0a565b60405190808252806020026020018201604052801561173657816020015b60408051808201909152600081526060602082015281526020019060019003908161170e5790505b509050600060015b8363ffffffff168163ffffffff16101561188e5763ffffffff81166000908152600b602052604090205473ffffffffffffffffffffffffffffffffffffffff161561187e5763ffffffff81166000908152600b60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff16835260018101805491928401916117d290614667565b80601f01602080910402602001604051908101604052809291908181526020018280546117fe90614667565b801561184b5780601f106118205761010080835404028352916020019161184b565b820191906000526020600020905b81548152906001019060200180831161182e57829003601f168201915b505050505081525050838381518110611866576118666145e3565b60200260200101819052508161187b90614612565b91505b611887816145c0565b905061173e565b50600e5461085e9060019063ffffffff1661464a565b6118ac612a21565b60005b8181101561070b5760008383838181106118cb576118cb6145e3565b90506020028101906118dd9190614a52565b6118e690614a95565b905060006118fc82600001518360200151610449565b905061190960038261345f565b611942576040517febf5255100000000000000000000000000000000000000000000000000000000815260048101829052602401610678565b61194c818361346b565b50508061195890614612565b90506118af565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610c9a576000848483818110611999576119996145e3565b90506020028101906119ab9190614938565b6119b49061496c565b805163ffffffff166000908152600b602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff168252600181018054959650939491939092840191611a0a90614667565b80601f0160208091040260200160405190810160405280929190818152602001828054611a3690614667565b8015611a835780601f10611a5857610100808354040283529160200191611a83565b820191906000526020600020905b815481529060010190602001808311611a6657829003601f168201915b50505091909252505081519192505073ffffffffffffffffffffffffffffffffffffffff16611ae95781516040517fadd9ae1e00000000000000000000000000000000000000000000000000000000815263ffffffff9091166004820152602401610678565b83158015611b0e5750805173ffffffffffffffffffffffffffffffffffffffff163314155b15611b47576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610678565b6040808301516000908152600c60205220600181015415611b9c5782604001516040517f5461848300000000000000000000000000000000000000000000000000000000815260040161067891815260200190565b6040830151611bdf5782604001516040517f64e2ee9200000000000000000000000000000000000000000000000000000000815260040161067891815260200190565b60208301511580611bfc57506020830151611bfc90600790612a06565b15611c33576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608301518051600003611c7557806040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106789190614a3f565b81548290600490611c9390640100000000900463ffffffff166145c0565b82546101009290920a63ffffffff818102199093169183160217909155825464010000000090041660005b8251811015611d3a57611cdc838281518110611538576115386145e3565b611d1457826040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106789190614a3f565b611d2983828151811061159d5761159d6145e3565b50611d3381614612565b9050611cbe565b50845183547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff91821617845560408601516002850155602086015160018501819055611d90916007919061345f16565b506040850151611da29060099061345f565b50845160408087015160208089015183519283529082015263ffffffff909216917f74becb12a5e8fd0e98077d02dfba8f647c9670c9df177e42c2418cf17a636f05910160405180910390a2505050505080611dfd90614612565b905061197d565b60015473ffffffffffffffffffffffffffffffffffffffff163314611e85576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610678565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b828114611f44576040517fab8b67c60000000000000000000000000000000000000000000000000000000081526004810184905260248101829052604401610678565b6000805473ffffffffffffffffffffffffffffffffffffffff16905b8481101561223f576000868683818110611f7c57611f7c6145e3565b9050602002016020810190611f919190614069565b63ffffffff81166000908152600b6020526040902080549192509073ffffffffffffffffffffffffffffffffffffffff16612000576040517fadd9ae1e00000000000000000000000000000000000000000000000000000000815263ffffffff83166004820152602401610678565b6000868685818110612014576120146145e3565b905060200281019061202691906146ba565b61202f906146f8565b805190915073ffffffffffffffffffffffffffffffffffffffff16612080576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805173ffffffffffffffffffffffffffffffffffffffff1633148015906120bd57503373ffffffffffffffffffffffffffffffffffffffff861614155b156120f6576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610678565b8051825473ffffffffffffffffffffffffffffffffffffffff9081169116141580612172575060208082015160405161212f9201613d5b565b60405160208183030381529060405280519060200120826001016040516020016121599190614b3b565b6040516020818303038152906040528051906020012014155b1561222b57805182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602081015160018301906121cc90826147b0565b50806000015173ffffffffffffffffffffffffffffffffffffffff168363ffffffff167f86f41145bde5dd7f523305452e4aad3685508c181432ec733d5f345009358a2883602001516040516122229190613d5b565b60405180910390a35b5050508061223890614612565b9050611f60565b505050505050565b61224f612a21565b63ffffffff8089166000908152600d60205260408120546401000000009004909116908190036122b3576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff8a166004820152602401610678565b612308888888886040518060a001604052808f63ffffffff168152602001876122db906145c0565b97508763ffffffff1681526020018a1515815260200189151581526020018863ffffffff16815250612aa4565b505050505050505050565b6123436040805160a081018252606080825260208201529081016000815260200160008152600060209091015290565b60008281526002602052604090819020815160a0810190925280548290829061236b90614667565b80601f016020809104026020016040519081016040528092919081815260200182805461239790614667565b80156123e45780601f106123b9576101008083540402835291602001916123e4565b820191906000526020600020905b8154815290600101906020018083116123c757829003601f168201915b505050505081526020016001820180546123fd90614667565b80601f016020809104026020016040519081016040528092919081815260200182805461242990614667565b80156124765780601f1061244b57610100808354040283529160200191612476565b820191906000526020600020905b81548152906001019060200180831161245957829003601f168201915b5050509183525050600282015460209091019060ff16600381111561249d5761249d614334565b60038111156124ae576124ae614334565b81526020016002820160019054906101000a900460ff1660018111156124d6576124d6614334565b60018111156124e7576124e7614334565b81526002919091015462010000900473ffffffffffffffffffffffffffffffffffffffff1660209091015292915050565b60608060006125276003613452565b905060006125356005613154565b82516125419190614be4565b905060008167ffffffffffffffff81111561255e5761255e613b0a565b6040519080825280602002602001820160405280156125bf57816020015b6125ac6040805160a081018252606080825260208201529081016000815260200160008152600060209091015290565b81526020019060019003908161257c5790505b50905060008267ffffffffffffffff8111156125dd576125dd613b0a565b604051908082528060200260200182016040528015612606578160200160208202803683370190505b5090506000805b855181101561287557600086828151811061262a5761262a6145e3565b60200260200101519050612648816005612a0690919063ffffffff16565b6128645760008181526002602052604090819020815160a0810190925280548290829061267490614667565b80601f01602080910402602001604051908101604052809291908181526020018280546126a090614667565b80156126ed5780601f106126c2576101008083540402835291602001916126ed565b820191906000526020600020905b8154815290600101906020018083116126d057829003601f168201915b5050505050815260200160018201805461270690614667565b80601f016020809104026020016040519081016040528092919081815260200182805461273290614667565b801561277f5780601f106127545761010080835404028352916020019161277f565b820191906000526020600020905b81548152906001019060200180831161276257829003601f168201915b5050509183525050600282015460209091019060ff1660038111156127a6576127a6614334565b60038111156127b7576127b7614334565b81526020016002820160019054906101000a900460ff1660018111156127df576127df614334565b60018111156127f0576127f0614334565b81526002919091015462010000900473ffffffffffffffffffffffffffffffffffffffff16602090910152855186908590811061282f5761282f6145e3565b60200260200101819052508084848151811061284d5761284d6145e3565b602090810291909101015261286183614612565b92505b5061286e81614612565b905061260d565b509096919550909350505050565b60608060006128926009613452565b9050805167ffffffffffffffff8111156128ae576128ae613b0a565b60405190808252806020026020018201604052801561291d57816020015b60408051608081018252600080825260208083018290529282015260608082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816128cc5790505b509250805167ffffffffffffffff81111561293a5761293a613b0a565b604051908082528060200260200182016040528015612963578160200160208202803683370190505b50915060005b81518110156129ec57612994828281518110612987576129876145e3565b6020026020010151611054565b8583815181106129a6576129a66145e3565b602002602001018584815181106129bf576129bf6145e3565b602002602001018263ffffffff1663ffffffff168152508290525050806129e590614612565b9050612969565b50509091565b6129fa612a21565b612a03816136ff565b50565b600081815260018301602052604081205415155b9392505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612aa2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610678565b565b805163ffffffff9081166000908152600d602090815260408083208286015185168452600101909152902060808301519091161580612af8575060808201518590612af0906001614bf7565b63ffffffff16115b15612b445760808201516040517fd5f5269100000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015260248101869052604401610678565b6001826020015163ffffffff161115612bf057815163ffffffff166000908152600d602090815260408220908401516001918201918391612b85919061464a565b63ffffffff1663ffffffff168152602001908152602001600020905060005b612bad82613154565b811015612bed57612bdc846000015163ffffffff16600c60006105fa858760000161315e90919063ffffffff16565b50612be681614612565b9050612ba4565b50505b60005b85811015612cda57612c20878783818110612c1057612c106145e3565b859260209091020135905061345f565b612c81578251878783818110612c3857612c386145e3565b6040517f636e405700000000000000000000000000000000000000000000000000000000815263ffffffff90941660048501526020029190910135602483015250604401610678565b8251612cc99063ffffffff16600c60008a8a86818110612ca357612ca36145e3565b90506020020135815260200190815260200160002060040161345f90919063ffffffff16565b50612cd381614612565b9050612bf3565b5060005b83811015612fc65736858583818110612cf957612cf96145e3565b9050602002810190612d0b91906146ba565b9050612d1960038235612a06565b612d52576040517fe181733f00000000000000000000000000000000000000000000000000000000815281356004820152602401610678565b612d5e60058235612a06565b15612d98576040517ff7d7a29400000000000000000000000000000000000000000000000000000000815281356004820152602401610678565b8035600090815260038401602052604081208054612db590614667565b90501115612e015783516040517f3927d08000000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015281356024820152604401610678565b60005b87811015612f1357612ea88235600c60008c8c86818110612e2757612e276145e3565b9050602002013581526020019081526020016000206003016000600c60008e8e88818110612e5757612e576145e3565b90506020020135815260200190815260200160002060000160049054906101000a900463ffffffff1663ffffffff1663ffffffff168152602001908152602001600020612a0690919063ffffffff16565b612f0357888882818110612ebe57612ebe6145e3565b6040517fa7e792500000000000000000000000000000000000000000000000000000000081526020909102929092013560048301525082356024820152604401610678565b612f0c81614612565b9050612e04565b5060028301805460018101825560009182526020918290208335910155612f3c90820182614c14565b82356000908152600386016020526040902091612f5a919083614c79565b508351602080860151612fb592918435908c908c90612f7b90880188614c14565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506137f492505050565b50612fbf81614612565b9050612cde565b50604080830151835163ffffffff9081166000908152600d602090815284822080549415156c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff90951694909417909355606086015186518316825284822080549115156d0100000000000000000000000000027fffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffff9092169190911790556080860151865183168252848220805491841668010000000000000000027fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff909216919091179055918501805186518316845292849020805493909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9093169290921790558351905191517ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c158170365192613144929163ffffffff92831681529116602082015260400190565b60405180910390a1505050505050565b6000610477825490565b6000612a1a83836138d5565b6000612a1a83836138ff565b6040805160e0810182526000808252602080830182905282840182905260608084018390526080840183905260a0840181905260c084015263ffffffff8581168352600d825284832080546401000000009004909116808452600190910182528483206002810180548751818602810186019098528088529596929591949390919083018282801561322757602002820191906000526020600020905b815481526020019060010190808311613213575b505050505090506000815167ffffffffffffffff81111561324a5761324a613b0a565b60405190808252806020026020018201604052801561329057816020015b6040805180820190915260008152606060208201528152602001906001900390816132685790505b50905060005b81518110156133b15760405180604001604052808483815181106132bc576132bc6145e3565b602002602001015181526020018560030160008685815181106132e1576132e16145e3565b60200260200101518152602001908152602001600020805461330290614667565b80601f016020809104026020016040519081016040528092919081815260200182805461332e90614667565b801561337b5780601f106133505761010080835404028352916020019161337b565b820191906000526020600020905b81548152906001019060200180831161335e57829003601f168201915b5050505050815250828281518110613395576133956145e3565b6020026020010181905250806133aa90614612565b9050613296565b506040805160e08101825263ffffffff8089166000818152600d60208181528683205480861687528b8616828801526801000000000000000081049095169686019690965260ff6c010000000000000000000000008504811615156060870152929091529093526d010000000000000000000000000090049091161515608082015260a0810161344085613452565b81526020019190915295945050505050565b60606000612a1a836139f2565b6000612a1a8383613a4e565b608081015173ffffffffffffffffffffffffffffffffffffffff16156135b957608081015173ffffffffffffffffffffffffffffffffffffffff163b1580613564575060808101516040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527f78bea72100000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff909116906301ffc9a790602401602060405180830381865afa15801561353e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135629190614d94565b155b156135b95760808101516040517fabb5e3fd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610678565b6000828152600260205260409020815182919081906135d890826147b0565b50602082015160018201906135ed90826147b0565b5060408201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600183600381111561362f5761362f614334565b021790555060608201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600181111561367657613676614334565b0217905550608091909101516002909101805473ffffffffffffffffffffffffffffffffffffffff90921662010000027fffffffffffffffffffff0000000000000000000000000000000000000000ffff90921691909117905560405182907f04f0a9bcf3f3a3b42a4d7ca081119755f82ebe43e0d30c8f7292c4fe0dc4a2ae90600090a25050565b3373ffffffffffffffffffffffffffffffffffffffff82160361377e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610678565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000848152600260208190526040909120015462010000900473ffffffffffffffffffffffffffffffffffffffff161561223f57600084815260026020819052604091829020015490517ffba64a7c0000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063fba64a7c9061389b908690869086908b908d90600401614db1565b600060405180830381600087803b1580156138b557600080fd5b505af11580156138c9573d6000803e3d6000fd5b50505050505050505050565b60008260000182815481106138ec576138ec6145e3565b9060005260206000200154905092915050565b600081815260018301602052604081205480156139e8576000613923600183614be4565b855490915060009061393790600190614be4565b905081811461399c576000866000018281548110613957576139576145e3565b906000526020600020015490508087600001848154811061397a5761397a6145e3565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806139ad576139ad614e33565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610477565b6000915050610477565b606081600001805480602002602001604051908101604052809291908181526020018280548015613a4257602002820191906000526020600020905b815481526020019060010190808311613a2e575b50505050509050919050565b6000818152600183016020526040812054613a9557508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610477565b506000610477565b508054613aa990614667565b6000825580601f10613ab9575050565b601f016020900490600052602060002090810190612a039190613af1565b5080546000825590600052602060002090810190612a0391905b5b80821115613b065760008155600101613af2565b5090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516080810167ffffffffffffffff81118282101715613b5c57613b5c613b0a565b60405290565b60405160a0810167ffffffffffffffff81118282101715613b5c57613b5c613b0a565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613bcc57613bcc613b0a565b604052919050565b600067ffffffffffffffff821115613bee57613bee613b0a565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112613c2b57600080fd5b8135613c3e613c3982613bd4565b613b85565b818152846020838601011115613c5357600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215613c8357600080fd5b823567ffffffffffffffff80821115613c9b57600080fd5b613ca786838701613c1a565b93506020850135915080821115613cbd57600080fd5b50613cca85828601613c1a565b9150509250929050565b600060208284031215613ce657600080fd5b5035919050565b60005b83811015613d08578181015183820152602001613cf0565b50506000910152565b60008151808452613d29816020860160208601613ced565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000612a1a6020830184613d11565b60008083601f840112613d8057600080fd5b50813567ffffffffffffffff811115613d9857600080fd5b6020830191508360208260051b850101111561104d57600080fd5b8015158114612a0357600080fd5b803563ffffffff81168114613dd557600080fd5b919050565b600080600080600080600060a0888a031215613df557600080fd5b873567ffffffffffffffff80821115613e0d57600080fd5b613e198b838c01613d6e565b909950975060208a0135915080821115613e3257600080fd5b50613e3f8a828b01613d6e565b9096509450506040880135613e5381613db3565b92506060880135613e6381613db3565b9150613e7160808901613dc1565b905092959891949750929550565b60008060208385031215613e9257600080fd5b823567ffffffffffffffff811115613ea957600080fd5b613eb585828601613d6e565b90969095509350505050565b600081518084526020808501945080840160005b83811015613ef157815187529582019590820190600101613ed5565b509495945050505050565b600063ffffffff80835116845260208181850151168186015260408281860151168187015260608501511515606087015260808501511515608087015260a0850151925060e060a0870152613f5460e0870184613ec1565b925060c085015186840360c08801528381518086528486019150848160051b870101858401935060005b82811015613fdb578782037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018452845180518352870151878301879052613fc887840182613d11565b9588019594880194925050600101613f7e565b509998505050505050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b8281101561405c577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc088860301845261404a858351613efc565b94509285019290850190600101614010565b5092979650505050505050565b60006020828403121561407b57600080fd5b612a1a82613dc1565b73ffffffffffffffffffffffffffffffffffffffff815116825260006020820151604060208501526140b96040850182613d11565b949350505050565b602081526000612a1a6020830184614084565b602081526000612a1a6020830184613efc565b600080604083850312156140fa57600080fd5b61410383613dc1565b946020939093013593505050565b6040815260006141246040830185613d11565b82810360208401526141368185613d11565b95945050505050565b63ffffffff8151168252602081015160208301526040810151604083015260006060820151608060608501526140b96080850182613ec1565b60408152600061418b604083018561413f565b905063ffffffff831660208301529392505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b8281101561405c577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452614201858351614084565b945092850192908501906001016141c7565b6000806000806040858703121561422957600080fd5b843567ffffffffffffffff8082111561424157600080fd5b61424d88838901613d6e565b9096509450602087013591508082111561426657600080fd5b5061427387828801613d6e565b95989497509550505050565b60008060008060008060008060c0898b03121561429b57600080fd5b6142a489613dc1565b9750602089013567ffffffffffffffff808211156142c157600080fd5b6142cd8c838d01613d6e565b909950975060408b01359150808211156142e657600080fd5b506142f38b828c01613d6e565b909650945050606089013561430781613db3565b9250608089013561431781613db3565b915061432560a08a01613dc1565b90509295985092959890939650565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6000815160a0845261437860a0850182613d11565b9050602083015184820360208601526143918282613d11565b9150506040830151600481106143a9576143a9614334565b60408501526060830151600281106143c3576143c3614334565b606085015260809283015173ffffffffffffffffffffffffffffffffffffffff1692909301919091525090565b602081526000612a1a6020830184614363565b6040815260006144166040830185613ec1565b6020838203818501528185518084528284019150828160051b85010183880160005b83811015614484577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552614472838351614363565b94860194925090850190600101614438565b50909998505050505050505050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015614508577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08887030185526144f686835161413f565b955093820193908201906001016144bc565b50508584038187015286518085528782019482019350915060005b8281101561454557845163ffffffff1684529381019392810192600101614523565b5091979650505050505050565b803573ffffffffffffffffffffffffffffffffffffffff81168114613dd557600080fd5b60006020828403121561458857600080fd5b612a1a82614552565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff8083168181036145d9576145d9614591565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361464357614643614591565b5060010190565b63ffffffff82811682821603908082111561086e5761086e614591565b600181811c9082168061467b57607f821691505b6020821081036146b4577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc18336030181126146ee57600080fd5b9190910192915050565b60006040823603121561470a57600080fd5b6040516040810167ffffffffffffffff828210818311171561472e5761472e613b0a565b8160405261473b85614552565b8352602085013591508082111561475157600080fd5b5061475e36828601613c1a565b60208301525092915050565b601f82111561070b57600081815260208120601f850160051c810160208610156147915750805b601f850160051c820191505b8181101561223f5782815560010161479d565b815167ffffffffffffffff8111156147ca576147ca613b0a565b6147de816147d88454614667565b8461476a565b602080601f83116001811461483157600084156147fb5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855561223f565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561487e5788860151825594840194600190910190840161485f565b50858210156148ba57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b6000602082840312156148dc57600080fd5b815167ffffffffffffffff8111156148f357600080fd5b8201601f8101841361490457600080fd5b8051614912613c3982613bd4565b81815285602083850101111561492757600080fd5b614136826020830160208601613ced565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff818336030181126146ee57600080fd5b60006080823603121561497e57600080fd5b614986613b39565b61498f83613dc1565b81526020808401358183015260408401356040830152606084013567ffffffffffffffff808211156149c057600080fd5b9085019036601f8301126149d357600080fd5b8135818111156149e5576149e5613b0a565b8060051b91506149f6848301613b85565b8181529183018401918481019036841115614a1057600080fd5b938501935b83851015614a2e57843582529385019390850190614a15565b606087015250939695505050505050565b602081526000612a1a6020830184613ec1565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff618336030181126146ee57600080fd5b803560028110613dd557600080fd5b600060a08236031215614aa757600080fd5b614aaf613b62565b823567ffffffffffffffff80821115614ac757600080fd5b614ad336838701613c1a565b83526020850135915080821115614ae957600080fd5b50614af636828601613c1a565b602083015250604083013560048110614b0e57600080fd5b6040820152614b1f60608401614a86565b6060820152614b3060808401614552565b608082015292915050565b6000602080835260008454614b4f81614667565b80848701526040600180841660008114614b705760018114614ba857614bd6565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838a01528284151560051b8a01019550614bd6565b896000528660002060005b85811015614bce5781548b8201860152908301908801614bb3565b8a0184019650505b509398975050505050505050565b8181038181111561047757610477614591565b63ffffffff81811683821601908082111561086e5761086e614591565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112614c4957600080fd5b83018035915067ffffffffffffffff821115614c6457600080fd5b60200191503681900382131561104d57600080fd5b67ffffffffffffffff831115614c9157614c91613b0a565b614ca583614c9f8354614667565b8361476a565b6000601f841160018114614cf75760008515614cc15750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355614d8d565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015614d465786850135825560209485019460019092019101614d26565b5086821015614d81577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b600060208284031215614da657600080fd5b8151612a1a81613db3565b6080815284608082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff861115614dea57600080fd5b8560051b808860a0850137820182810360a09081016020850152614e1090820187613d11565b91505063ffffffff80851660408401528084166060840152509695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
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

func (_CapabilityRegistry *CapabilityRegistryCaller) GetCapabilities(opts *bind.CallOpts) ([][32]byte, []CapabilityRegistryCapability, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getCapabilities")

	if err != nil {
		return *new([][32]byte), *new([]CapabilityRegistryCapability), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	out1 := *abi.ConvertType(out[1], new([]CapabilityRegistryCapability)).(*[]CapabilityRegistryCapability)

	return out0, out1, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetCapabilities() ([][32]byte, []CapabilityRegistryCapability, error) {
	return _CapabilityRegistry.Contract.GetCapabilities(&_CapabilityRegistry.CallOpts)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetCapabilities() ([][32]byte, []CapabilityRegistryCapability, error) {
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

func (_CapabilityRegistry *CapabilityRegistryCaller) GetCapabilityConfigs(opts *bind.CallOpts, donId uint32, capabilityId [32]byte) ([]byte, []byte, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getCapabilityConfigs", donId, capabilityId)

	if err != nil {
		return *new([]byte), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetCapabilityConfigs(donId uint32, capabilityId [32]byte) ([]byte, []byte, error) {
	return _CapabilityRegistry.Contract.GetCapabilityConfigs(&_CapabilityRegistry.CallOpts, donId, capabilityId)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetCapabilityConfigs(donId uint32, capabilityId [32]byte) ([]byte, []byte, error) {
	return _CapabilityRegistry.Contract.GetCapabilityConfigs(&_CapabilityRegistry.CallOpts, donId, capabilityId)
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

func (_CapabilityRegistry *CapabilityRegistryCaller) GetNodeOperator(opts *bind.CallOpts, nodeOperatorId uint32) (CapabilityRegistryNodeOperator, error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getNodeOperator", nodeOperatorId)

	if err != nil {
		return *new(CapabilityRegistryNodeOperator), err
	}

	out0 := *abi.ConvertType(out[0], new(CapabilityRegistryNodeOperator)).(*CapabilityRegistryNodeOperator)

	return out0, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetNodeOperator(nodeOperatorId uint32) (CapabilityRegistryNodeOperator, error) {
	return _CapabilityRegistry.Contract.GetNodeOperator(&_CapabilityRegistry.CallOpts, nodeOperatorId)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetNodeOperator(nodeOperatorId uint32) (CapabilityRegistryNodeOperator, error) {
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

func (_CapabilityRegistry *CapabilityRegistryCaller) GetNodes(opts *bind.CallOpts) (GetNodes,

	error) {
	var out []interface{}
	err := _CapabilityRegistry.contract.Call(opts, &out, "getNodes")

	outstruct := new(GetNodes)
	if err != nil {
		return *outstruct, err
	}

	outstruct.NodeInfo = *abi.ConvertType(out[0], new([]CapabilityRegistryNodeInfo)).(*[]CapabilityRegistryNodeInfo)
	outstruct.ConfigCounts = *abi.ConvertType(out[1], new([]uint32)).(*[]uint32)

	return *outstruct, err

}

func (_CapabilityRegistry *CapabilityRegistrySession) GetNodes() (GetNodes,

	error) {
	return _CapabilityRegistry.Contract.GetNodes(&_CapabilityRegistry.CallOpts)
}

func (_CapabilityRegistry *CapabilityRegistryCallerSession) GetNodes() (GetNodes,

	error) {
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) AddCapabilities(opts *bind.TransactOpts, capabilities []CapabilityRegistryCapability) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "addCapabilities", capabilities)
}

func (_CapabilityRegistry *CapabilityRegistrySession) AddCapabilities(capabilities []CapabilityRegistryCapability) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddCapabilities(&_CapabilityRegistry.TransactOpts, capabilities)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) AddCapabilities(capabilities []CapabilityRegistryCapability) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddCapabilities(&_CapabilityRegistry.TransactOpts, capabilities)
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) DeprecateCapabilities(opts *bind.TransactOpts, hashedCapabilityIds [][32]byte) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "deprecateCapabilities", hashedCapabilityIds)
}

func (_CapabilityRegistry *CapabilityRegistrySession) DeprecateCapabilities(hashedCapabilityIds [][32]byte) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.DeprecateCapabilities(&_CapabilityRegistry.TransactOpts, hashedCapabilityIds)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) DeprecateCapabilities(hashedCapabilityIds [][32]byte) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.DeprecateCapabilities(&_CapabilityRegistry.TransactOpts, hashedCapabilityIds)
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) RemoveNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []uint32) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "removeNodeOperators", nodeOperatorIds)
}

func (_CapabilityRegistry *CapabilityRegistrySession) RemoveNodeOperators(nodeOperatorIds []uint32) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.RemoveNodeOperators(&_CapabilityRegistry.TransactOpts, nodeOperatorIds)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) RemoveNodeOperators(nodeOperatorIds []uint32) (*types.Transaction, error) {
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) UpdateNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []uint32, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "updateNodeOperators", nodeOperatorIds, nodeOperators)
}

func (_CapabilityRegistry *CapabilityRegistrySession) UpdateNodeOperators(nodeOperatorIds []uint32, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.UpdateNodeOperators(&_CapabilityRegistry.TransactOpts, nodeOperatorIds, nodeOperators)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) UpdateNodeOperators(nodeOperatorIds []uint32, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error) {
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

type CapabilityRegistryCapabilityConfiguredIterator struct {
	Event *CapabilityRegistryCapabilityConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilityRegistryCapabilityConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilityRegistryCapabilityConfigured)
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
		it.Event = new(CapabilityRegistryCapabilityConfigured)
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

func (it *CapabilityRegistryCapabilityConfiguredIterator) Error() error {
	return it.fail
}

func (it *CapabilityRegistryCapabilityConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilityRegistryCapabilityConfigured struct {
	HashedCapabilityId [32]byte
	Raw                types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterCapabilityConfigured(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilityRegistryCapabilityConfiguredIterator, error) {

	var hashedCapabilityIdRule []interface{}
	for _, hashedCapabilityIdItem := range hashedCapabilityId {
		hashedCapabilityIdRule = append(hashedCapabilityIdRule, hashedCapabilityIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "CapabilityConfigured", hashedCapabilityIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryCapabilityConfiguredIterator{contract: _CapabilityRegistry.contract, event: "CapabilityConfigured", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchCapabilityConfigured(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryCapabilityConfigured, hashedCapabilityId [][32]byte) (event.Subscription, error) {

	var hashedCapabilityIdRule []interface{}
	for _, hashedCapabilityIdItem := range hashedCapabilityId {
		hashedCapabilityIdRule = append(hashedCapabilityIdRule, hashedCapabilityIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "CapabilityConfigured", hashedCapabilityIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilityRegistryCapabilityConfigured)
				if err := _CapabilityRegistry.contract.UnpackLog(event, "CapabilityConfigured", log); err != nil {
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

func (_CapabilityRegistry *CapabilityRegistryFilterer) ParseCapabilityConfigured(log types.Log) (*CapabilityRegistryCapabilityConfigured, error) {
	event := new(CapabilityRegistryCapabilityConfigured)
	if err := _CapabilityRegistry.contract.UnpackLog(event, "CapabilityConfigured", log); err != nil {
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
	NodeOperatorId uint32
	Signer         [32]byte
	Raw            types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterNodeAdded(opts *bind.FilterOpts, nodeOperatorId []uint32) (*CapabilityRegistryNodeAddedIterator, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "NodeAdded", nodeOperatorIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryNodeAddedIterator{contract: _CapabilityRegistry.contract, event: "NodeAdded", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchNodeAdded(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeAdded, nodeOperatorId []uint32) (event.Subscription, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "NodeAdded", nodeOperatorIdRule)
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
	NodeOperatorId uint32
	Admin          common.Address
	Name           string
	Raw            types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterNodeOperatorAdded(opts *bind.FilterOpts, nodeOperatorId []uint32, admin []common.Address) (*CapabilityRegistryNodeOperatorAddedIterator, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "NodeOperatorAdded", nodeOperatorIdRule, adminRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryNodeOperatorAddedIterator{contract: _CapabilityRegistry.contract, event: "NodeOperatorAdded", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchNodeOperatorAdded(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeOperatorAdded, nodeOperatorId []uint32, admin []common.Address) (event.Subscription, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "NodeOperatorAdded", nodeOperatorIdRule, adminRule)
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
	NodeOperatorId uint32
	Raw            types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterNodeOperatorRemoved(opts *bind.FilterOpts, nodeOperatorId []uint32) (*CapabilityRegistryNodeOperatorRemovedIterator, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "NodeOperatorRemoved", nodeOperatorIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryNodeOperatorRemovedIterator{contract: _CapabilityRegistry.contract, event: "NodeOperatorRemoved", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchNodeOperatorRemoved(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeOperatorRemoved, nodeOperatorId []uint32) (event.Subscription, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "NodeOperatorRemoved", nodeOperatorIdRule)
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
	NodeOperatorId uint32
	Admin          common.Address
	Name           string
	Raw            types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterNodeOperatorUpdated(opts *bind.FilterOpts, nodeOperatorId []uint32, admin []common.Address) (*CapabilityRegistryNodeOperatorUpdatedIterator, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "NodeOperatorUpdated", nodeOperatorIdRule, adminRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryNodeOperatorUpdatedIterator{contract: _CapabilityRegistry.contract, event: "NodeOperatorUpdated", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchNodeOperatorUpdated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeOperatorUpdated, nodeOperatorId []uint32, admin []common.Address) (event.Subscription, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "NodeOperatorUpdated", nodeOperatorIdRule, adminRule)
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
	NodeOperatorId uint32
	Signer         [32]byte
	Raw            types.Log
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) FilterNodeUpdated(opts *bind.FilterOpts, nodeOperatorId []uint32) (*CapabilityRegistryNodeUpdatedIterator, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.FilterLogs(opts, "NodeUpdated", nodeOperatorIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilityRegistryNodeUpdatedIterator{contract: _CapabilityRegistry.contract, event: "NodeUpdated", logs: logs, sub: sub}, nil
}

func (_CapabilityRegistry *CapabilityRegistryFilterer) WatchNodeUpdated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeUpdated, nodeOperatorId []uint32) (event.Subscription, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}

	logs, sub, err := _CapabilityRegistry.contract.WatchLogs(opts, "NodeUpdated", nodeOperatorIdRule)
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

type GetNodes struct {
	NodeInfo     []CapabilityRegistryNodeInfo
	ConfigCounts []uint32
}

func (_CapabilityRegistry *CapabilityRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _CapabilityRegistry.abi.Events["CapabilityConfigured"].ID:
		return _CapabilityRegistry.ParseCapabilityConfigured(log)
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

func (CapabilityRegistryCapabilityConfigured) Topic() common.Hash {
	return common.HexToHash("0x04f0a9bcf3f3a3b42a4d7ca081119755f82ebe43e0d30c8f7292c4fe0dc4a2ae")
}

func (CapabilityRegistryCapabilityDeprecated) Topic() common.Hash {
	return common.HexToHash("0xdcea1b78b6ddc31592a94607d537543fcaafda6cc52d6d5cc7bbfca1422baf21")
}

func (CapabilityRegistryConfigSet) Topic() common.Hash {
	return common.HexToHash("0xf264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c1581703651")
}

func (CapabilityRegistryNodeAdded) Topic() common.Hash {
	return common.HexToHash("0x74becb12a5e8fd0e98077d02dfba8f647c9670c9df177e42c2418cf17a636f05")
}

func (CapabilityRegistryNodeOperatorAdded) Topic() common.Hash {
	return common.HexToHash("0x78e94ca80be2c30abc061b99e7eb8583b1254781734b1e3ce339abb57da2fe8e")
}

func (CapabilityRegistryNodeOperatorRemoved) Topic() common.Hash {
	return common.HexToHash("0xa59268ca81d40429e65ccea5385b59cf2d3fc6519371dee92f8eb1dae5107a7a")
}

func (CapabilityRegistryNodeOperatorUpdated) Topic() common.Hash {
	return common.HexToHash("0x86f41145bde5dd7f523305452e4aad3685508c181432ec733d5f345009358a28")
}

func (CapabilityRegistryNodeRemoved) Topic() common.Hash {
	return common.HexToHash("0x5254e609a97bab37b7cc79fe128f85c097bd6015c6e1624ae0ba392eb9753205")
}

func (CapabilityRegistryNodeUpdated) Topic() common.Hash {
	return common.HexToHash("0x4b5b465e22eea0c3d40c30e936643245b80d19b2dcf75788c0699fe8d8db645b")
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
	GetCapabilities(opts *bind.CallOpts) ([][32]byte, []CapabilityRegistryCapability, error)

	GetCapability(opts *bind.CallOpts, hashedId [32]byte) (CapabilityRegistryCapability, error)

	GetCapabilityConfigs(opts *bind.CallOpts, donId uint32, capabilityId [32]byte) ([]byte, []byte, error)

	GetDON(opts *bind.CallOpts, donId uint32) (CapabilityRegistryDONInfo, error)

	GetDONs(opts *bind.CallOpts) ([]CapabilityRegistryDONInfo, error)

	GetHashedCapabilityId(opts *bind.CallOpts, labelledName string, version string) ([32]byte, error)

	GetNode(opts *bind.CallOpts, p2pId [32]byte) (CapabilityRegistryNodeInfo, uint32, error)

	GetNodeOperator(opts *bind.CallOpts, nodeOperatorId uint32) (CapabilityRegistryNodeOperator, error)

	GetNodeOperators(opts *bind.CallOpts) ([]CapabilityRegistryNodeOperator, error)

	GetNodes(opts *bind.CallOpts) (GetNodes,

		error)

	IsCapabilityDeprecated(opts *bind.CallOpts, hashedCapabilityId [32]byte) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddCapabilities(opts *bind.TransactOpts, capabilities []CapabilityRegistryCapability) (*types.Transaction, error)

	AddDON(opts *bind.TransactOpts, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint32) (*types.Transaction, error)

	AddNodeOperators(opts *bind.TransactOpts, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error)

	AddNodes(opts *bind.TransactOpts, nodes []CapabilityRegistryNodeInfo) (*types.Transaction, error)

	DeprecateCapabilities(opts *bind.TransactOpts, hashedCapabilityIds [][32]byte) (*types.Transaction, error)

	RemoveDONs(opts *bind.TransactOpts, donIds []uint32) (*types.Transaction, error)

	RemoveNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []uint32) (*types.Transaction, error)

	RemoveNodes(opts *bind.TransactOpts, removedNodeP2PIds [][32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateDON(opts *bind.TransactOpts, donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint32) (*types.Transaction, error)

	UpdateNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []uint32, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error)

	UpdateNodes(opts *bind.TransactOpts, nodes []CapabilityRegistryNodeInfo) (*types.Transaction, error)

	FilterCapabilityConfigured(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilityRegistryCapabilityConfiguredIterator, error)

	WatchCapabilityConfigured(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryCapabilityConfigured, hashedCapabilityId [][32]byte) (event.Subscription, error)

	ParseCapabilityConfigured(log types.Log) (*CapabilityRegistryCapabilityConfigured, error)

	FilterCapabilityDeprecated(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilityRegistryCapabilityDeprecatedIterator, error)

	WatchCapabilityDeprecated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryCapabilityDeprecated, hashedCapabilityId [][32]byte) (event.Subscription, error)

	ParseCapabilityDeprecated(log types.Log) (*CapabilityRegistryCapabilityDeprecated, error)

	FilterConfigSet(opts *bind.FilterOpts) (*CapabilityRegistryConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*CapabilityRegistryConfigSet, error)

	FilterNodeAdded(opts *bind.FilterOpts, nodeOperatorId []uint32) (*CapabilityRegistryNodeAddedIterator, error)

	WatchNodeAdded(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeAdded, nodeOperatorId []uint32) (event.Subscription, error)

	ParseNodeAdded(log types.Log) (*CapabilityRegistryNodeAdded, error)

	FilterNodeOperatorAdded(opts *bind.FilterOpts, nodeOperatorId []uint32, admin []common.Address) (*CapabilityRegistryNodeOperatorAddedIterator, error)

	WatchNodeOperatorAdded(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeOperatorAdded, nodeOperatorId []uint32, admin []common.Address) (event.Subscription, error)

	ParseNodeOperatorAdded(log types.Log) (*CapabilityRegistryNodeOperatorAdded, error)

	FilterNodeOperatorRemoved(opts *bind.FilterOpts, nodeOperatorId []uint32) (*CapabilityRegistryNodeOperatorRemovedIterator, error)

	WatchNodeOperatorRemoved(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeOperatorRemoved, nodeOperatorId []uint32) (event.Subscription, error)

	ParseNodeOperatorRemoved(log types.Log) (*CapabilityRegistryNodeOperatorRemoved, error)

	FilterNodeOperatorUpdated(opts *bind.FilterOpts, nodeOperatorId []uint32, admin []common.Address) (*CapabilityRegistryNodeOperatorUpdatedIterator, error)

	WatchNodeOperatorUpdated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeOperatorUpdated, nodeOperatorId []uint32, admin []common.Address) (event.Subscription, error)

	ParseNodeOperatorUpdated(log types.Log) (*CapabilityRegistryNodeOperatorUpdated, error)

	FilterNodeRemoved(opts *bind.FilterOpts) (*CapabilityRegistryNodeRemovedIterator, error)

	WatchNodeRemoved(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeRemoved) (event.Subscription, error)

	ParseNodeRemoved(log types.Log) (*CapabilityRegistryNodeRemoved, error)

	FilterNodeUpdated(opts *bind.FilterOpts, nodeOperatorId []uint32) (*CapabilityRegistryNodeUpdatedIterator, error)

	WatchNodeUpdated(opts *bind.WatchOpts, sink chan<- *CapabilityRegistryNodeUpdated, nodeOperatorId []uint32) (event.Subscription, error)

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
