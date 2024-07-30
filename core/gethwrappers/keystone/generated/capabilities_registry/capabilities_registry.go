// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package capabilities_registry

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

type CapabilitiesRegistryCapability struct {
	LabelledName          string
	Version               string
	CapabilityType        uint8
	ResponseType          uint8
	ConfigurationContract common.Address
}

type CapabilitiesRegistryCapabilityConfiguration struct {
	CapabilityId [32]byte
	Config       []byte
}

type CapabilitiesRegistryCapabilityInfo struct {
	HashedId              [32]byte
	LabelledName          string
	Version               string
	CapabilityType        uint8
	ResponseType          uint8
	ConfigurationContract common.Address
	IsDeprecated          bool
}

type CapabilitiesRegistryDONInfo struct {
	Id                       uint32
	ConfigCount              uint32
	F                        uint8
	IsPublic                 bool
	AcceptsWorkflows         bool
	NodeP2PIds               [][32]byte
	CapabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration
}

type CapabilitiesRegistryNodeInfo struct {
	NodeOperatorId      uint32
	ConfigCount         uint32
	WorkflowDONId       uint32
	Signer              [32]byte
	P2pId               [32]byte
	HashedCapabilityIds [][32]byte
	CapabilitiesDONIds  []*big.Int
}

type CapabilitiesRegistryNodeOperator struct {
	Admin common.Address
	Name  string
}

type CapabilitiesRegistryNodeParams struct {
	NodeOperatorId      uint32
	Signer              [32]byte
	P2pId               [32]byte
	HashedCapabilityIds [][32]byte
}

var CapabilitiesRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityIsDeprecated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"CapabilityRequiredByDON\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"DONDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONCapability\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONNode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedConfigurationContract\",\"type\":\"address\"}],\"name\":\"InvalidCapabilityConfigurationContractInterface\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"nodeCount\",\"type\":\"uint256\"}],\"name\":\"InvalidFaultTolerance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"InvalidNodeCapabilities\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeOperatorAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"InvalidNodeP2PId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lengthOne\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lengthTwo\",\"type\":\"uint256\"}],\"name\":\"LengthMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodeAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodeDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"NodeDoesNotSupportCapability\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"NodeOperatorDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodePartOfCapabilitiesDON\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodePartOfWorkflowDON\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDeprecated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"NodeOperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"NodeRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilitiesRegistry.Capability[]\",\"name\":\"capabilities\",\"type\":\"tuple[]\"}],\"name\":\"addCapabilities\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"name\":\"addDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilitiesRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"addNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilitiesRegistry.NodeParams[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"addNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"deprecateCapabilities\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCapabilities\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"hashedId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isDeprecated\",\"type\":\"bool\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedId\",\"type\":\"bytes32\"}],\"name\":\"getCapability\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"hashedId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isDeprecated\",\"type\":\"bool\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"getCapabilityConfigs\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"getDON\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilitiesRegistry.DONInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilitiesRegistry.DONInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"}],\"name\":\"getHashedCapabilityId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"getNode\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"workflowDONId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"capabilitiesDONIds\",\"type\":\"uint256[]\"}],\"internalType\":\"structCapabilitiesRegistry.NodeInfo\",\"name\":\"nodeInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"getNodeOperator\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilitiesRegistry.NodeOperator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeOperators\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilitiesRegistry.NodeOperator[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodes\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"workflowDONId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"capabilitiesDONIds\",\"type\":\"uint256[]\"}],\"internalType\":\"structCapabilitiesRegistry.NodeInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"isCapabilityDeprecated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"donIds\",\"type\":\"uint32[]\"}],\"name\":\"removeDONs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint32[]\"}],\"name\":\"removeNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"removedNodeP2PIds\",\"type\":\"bytes32[]\"}],\"name\":\"removeNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"name\":\"updateDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint32[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilitiesRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"updateNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilitiesRegistry.NodeParams[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"updateNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600e80546001600160401b0319166401000000011790553480156200002857600080fd5b503380600081620000805760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000b357620000b381620000bc565b50505062000167565b336001600160a01b03821603620001165760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000077565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6150f680620001776000396000f3fe608060405234801561001057600080fd5b50600436106101ae5760003560e01c80635e65e309116100ee5780638da5cb5b11610097578063d8bc7b6811610071578063d8bc7b68146103f6578063ddbe4f8214610409578063e29581aa1461041e578063f2fde38b1461043357600080fd5b80638da5cb5b1461039b5780639cb7c5f4146103c3578063d59a79f6146103e357600080fd5b806373ac22b4116100c857806373ac22b41461036d57806379ba50971461038057806386fa42461461038857600080fd5b80635e65e3091461033257806366acaa3314610345578063715f52951461035a57600080fd5b8063235374051161015b578063398f377311610135578063398f3773146102cb5780633f2a13c9146102de57806350c946fe146102ff5780635d83d9671461031f57600080fd5b80632353740514610285578063275459f2146102a55780632c01a1e8146102b857600080fd5b80631d05394c1161018c5780631d05394c1461023b578063214502431461025057806322bdbcbc1461026557600080fd5b80630fe5800a146101b357806312570011146101d9578063181f5a77146101fc575b600080fd5b6101c66101c1366004613e8b565b610446565b6040519081526020015b60405180910390f35b6101ec6101e7366004613eef565b61047a565b60405190151581526020016101d0565b604080518082018252601a81527f4361706162696c6974696573526567697374727920312e302e30000000000000602082015290516101d09190613f76565b61024e610249366004613fce565b610487565b005b61025861069c565b6040516101d09190614150565b6102786102733660046141eb565b6107f9565b6040516101d09190614243565b6102986102933660046141eb565b6108e6565b6040516101d09190614256565b61024e6102b3366004613fce565b61092a565b61024e6102c6366004613fce565b610a01565b61024e6102d9366004613fce565b610c9d565b6102f16102ec366004614269565b610e5c565b6040516101d0929190614293565b61031261030d366004613eef565b611048565b6040516101d09190614358565b61024e61032d366004613fce565b611122565b61024e610340366004613fce565b611217565b61034d61193f565b6040516101d0919061436b565b61024e610368366004613fce565b611b22565b61024e61037b366004613fce565b611bd4565b61024e6120a2565b61024e6103963660046143e0565b61219f565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101d0565b6103d66103d1366004613eef565b6124df565b6040516101d0919061452f565b61024e6103f1366004614561565b61271a565b61024e610404366004614616565b6127e3565b6104116128ad565b6040516101d091906146bb565b6104266129a1565b6040516101d09190614730565b61024e6104413660046147c9565b612aaa565b6000828260405160200161045b929190614293565b6040516020818303038152906040528051906020012090505b92915050565b6000610474600583612abe565b61048f612ad9565b60005b818110156106975760008383838181106104ae576104ae6147e4565b90506020020160208101906104c391906141eb565b63ffffffff8181166000908152600d60209081526040808320805464010000000081049095168085526001820190935290832094955093909290916a010000000000000000000090910460ff16905b61051b83612b5c565b8110156105bb57811561057157600c60006105368584612b66565b8152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff1690556105b3565b6105b18663ffffffff16600c60006105928588612b6690919063ffffffff16565b8152602001908152602001600020600401612b7290919063ffffffff16565b505b600101610512565b508354640100000000900463ffffffff16600003610612576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff861660048201526024015b60405180910390fd5b63ffffffff85166000818152600d6020908152604080832080547fffffffffffffffffffffffffffffffffffffffffff00000000000000000000001690558051938452908301919091527ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c1581703651910160405180910390a15050505050806001019050610492565b505050565b600e54606090640100000000900463ffffffff1660006106bd600183614842565b63ffffffff1667ffffffffffffffff8111156106db576106db613d25565b60405190808252806020026020018201604052801561076257816020015b6040805160e081018252600080825260208083018290529282018190526060808301829052608083019190915260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816106f95790505b509050600060015b8363ffffffff168163ffffffff1610156107d65763ffffffff8082166000908152600d602052604090205416156107ce576107a481612b7e565b8383815181106107b6576107b66147e4565b6020026020010181905250816107cb9061485f565b91505b60010161076a565b506107e2600184614842565b63ffffffff1681146107f2578082525b5092915050565b60408051808201909152600081526060602082015263ffffffff82166000908152600b60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff168352600181018054919284019161085d90614897565b80601f016020809104026020016040519081016040528092919081815260200182805461088990614897565b80156108d65780601f106108ab576101008083540402835291602001916108d6565b820191906000526020600020905b8154815290600101906020018083116108b957829003601f168201915b5050505050815250509050919050565b6040805160e0810182526000808252602082018190529181018290526060808201839052608082019290925260a0810182905260c081019190915261047482612b7e565b610932612ad9565b60005b63ffffffff811682111561069757600083838363ffffffff1681811061095d5761095d6147e4565b905060200201602081019061097291906141eb565b63ffffffff81166000908152600b6020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001681559192506109bd6001830182613cb8565b505060405163ffffffff8216907fa59268ca81d40429e65ccea5385b59cf2d3fc6519371dee92f8eb1dae5107a7a90600090a2506109fa816148ea565b9050610935565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610c97576000848483818110610a3b57610a3b6147e4565b602090810292909201356000818152600c90935260409092206001810154929350919050610a98576040517fd82f6adb00000000000000000000000000000000000000000000000000000000815260048101839052602401610609565b6000610aa682600401612b5c565b1115610afb57610ab96004820184612b66565b6040517f60a6d89800000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015260248101839052604401610609565b805468010000000000000000900463ffffffff1615610b635780546040517f60b9df730000000000000000000000000000000000000000000000000000000081526801000000000000000090910463ffffffff16600482015260248101839052604401610609565b83158015610b9d5750805463ffffffff166000908152600b602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b15610bd6576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610609565b6001810154610be790600790612b72565b506002810154610bf990600990612b72565b506000828152600c6020526040812080547fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016815560018101829055600281018290559060048201818181610c4e8282613cf2565b5050505050507f5254e609a97bab37b7cc79fe128f85c097bd6015c6e1624ae0ba392eb975320582604051610c8591815260200190565b60405180910390a15050600101610a1f565b50505050565b610ca5612ad9565b60005b81811015610697576000838383818110610cc457610cc46147e4565b9050602002810190610cd6919061490d565b610cdf9061494b565b805190915073ffffffffffffffffffffffffffffffffffffffff16610d30576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e54604080518082018252835173ffffffffffffffffffffffffffffffffffffffff908116825260208086015181840190815263ffffffff9095166000818152600b909252939020825181547fffffffffffffffffffffffff00000000000000000000000000000000000000001692169190911781559251919290916001820190610dbc9082614a05565b5050600e8054909150600090610dd79063ffffffff166148ea565b91906101000a81548163ffffffff021916908363ffffffff160217905550816000015173ffffffffffffffffffffffffffffffffffffffff168163ffffffff167f78e94ca80be2c30abc061b99e7eb8583b1254781734b1e3ce339abb57da2fe8e8460200151604051610e4a9190613f76565b60405180910390a35050600101610ca8565b63ffffffff8083166000908152600d60209081526040808320805464010000000090049094168084526001909401825280832085845260030190915281208054606093849390929091610eae90614897565b80601f0160208091040260200160405190810160405280929190818152602001828054610eda90614897565b8015610f275780601f10610efc57610100808354040283529160200191610f27565b820191906000526020600020905b815481529060010190602001808311610f0a57829003601f168201915b5050506000888152600260208190526040909120015492935060609262010000900473ffffffffffffffffffffffffffffffffffffffff1615915061103a905057600086815260026020819052604091829020015490517f8318ed5d00000000000000000000000000000000000000000000000000000000815263ffffffff891660048201526201000090910473ffffffffffffffffffffffffffffffffffffffff1690638318ed5d90602401600060405180830381865afa158015610ff1573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526110379190810190614b1f565b90505b9093509150505b9250929050565b6040805160e0810182526000808252602082018190529181018290526060808201839052608082019290925260a0810182905260c08101919091526040805160e0810182526000848152600c6020908152838220805463ffffffff8082168652640100000000820481168487018190526801000000000000000090920416858701526001820154606086015260028201546080860152835260030190529190912060a08201906110f790612e49565b815260200161111a600c6000868152602001908152602001600020600401612e49565b905292915050565b61112a612ad9565b60005b81811015610697576000838383818110611149576111496147e4565b905060200201359050611166816003612abe90919063ffffffff16565b61119f576040517fe181733f00000000000000000000000000000000000000000000000000000000815260048101829052602401610609565b6111aa600582612e56565b6111e3576040517ff7d7a29400000000000000000000000000000000000000000000000000000000815260048101829052602401610609565b60405181907fdcea1b78b6ddc31592a94607d537543fcaafda6cc52d6d5cc7bbfca1422baf2190600090a25060010161112d565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610c97576000848483818110611251576112516147e4565b90506020028101906112639190614b8d565b61126c90614bc1565b6040808201516000908152600c6020908152828220805463ffffffff168352600b82528383208451808601909552805473ffffffffffffffffffffffffffffffffffffffff16855260018101805496975091959394939092840191906112d190614897565b80601f01602080910402602001604051908101604052809291908181526020018280546112fd90614897565b801561134a5780601f1061131f5761010080835404028352916020019161134a565b820191906000526020600020905b81548152906001019060200180831161132d57829003601f168201915b50505091909252505050600183015490915061139a5782604001516040517fd82f6adb00000000000000000000000000000000000000000000000000000000815260040161060991815260200190565b841580156113bf5750805173ffffffffffffffffffffffffffffffffffffffff163314155b156113f8576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610609565b6020830151611433576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001820154602084015181146114b457602084015161145490600790612abe565b1561148b576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602084015160018401556114a0600782612b72565b5060208401516114b290600790612e56565b505b606084015180516000036114f657806040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106099190614c94565b8354600090859060049061151790640100000000900463ffffffff166148ea565b91906101000a81548163ffffffff021916908363ffffffff1602179055905060005b82518110156115fc5761156f838281518110611557576115576147e4565b60200260200101516003612abe90919063ffffffff16565b6115a757826040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106099190614c94565b6115f38382815181106115bc576115bc6147e4565b60200260200101518760030160008563ffffffff1663ffffffff168152602001908152602001600020612e5690919063ffffffff16565b50600101611539565b50845468010000000000000000900463ffffffff16801561175d5763ffffffff8082166000908152600d60209081526040808320805464010000000090049094168352600190930181528282206002018054845181840281018401909552808552929392909183018282801561169157602002820191906000526020600020905b81548152602001906001019080831161167d575b5050505050905060005b815181101561175a576116f08282815181106116b9576116b96147e4565b60200260200101518960030160008763ffffffff1663ffffffff168152602001908152602001600020612abe90919063ffffffff16565b61175257818181518110611706576117066147e4565b6020026020010151836040517f03dcd86200000000000000000000000000000000000000000000000000000000815260040161060992919091825263ffffffff16602082015260400190565b60010161169b565b50505b600061176b87600401612e49565b905060005b81518163ffffffff1610156118b1576000828263ffffffff1681518110611799576117996147e4565b60209081029190910181015163ffffffff8082166000908152600d8452604080822080546401000000009004909316825260019092018452818120600201805483518187028101870190945280845293955090939192909183018282801561182057602002820191906000526020600020905b81548152602001906001019080831161180c575b5050505050905060005b815181101561189d5761187f828281518110611848576118486147e4565b60200260200101518c60030160008a63ffffffff1663ffffffff168152602001908152602001600020612abe90919063ffffffff16565b61189557818181518110611706576117066147e4565b60010161182a565b505050806118aa906148ea565b9050611770565b50875187547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff90911690811788556040808a015160028a018190556020808c01518351928352908201527f4b5b465e22eea0c3d40c30e936643245b80d19b2dcf75788c0699fe8d8db645b910160405180910390a25050505050505050806001019050611235565b600e5460609063ffffffff166000611958600183614842565b63ffffffff1667ffffffffffffffff81111561197657611976613d25565b6040519080825280602002602001820160405280156119bc57816020015b6040805180820190915260008152606060208201528152602001906001900390816119945790505b509050600060015b8363ffffffff168163ffffffff161015611b0c5763ffffffff81166000908152600b602052604090205473ffffffffffffffffffffffffffffffffffffffff1615611b045763ffffffff81166000908152600b60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff1683526001810180549192840191611a5890614897565b80601f0160208091040260200160405190810160405280929190818152602001828054611a8490614897565b8015611ad15780601f10611aa657610100808354040283529160200191611ad1565b820191906000526020600020905b815481529060010190602001808311611ab457829003601f168201915b505050505081525050838381518110611aec57611aec6147e4565b602002602001018190525081611b019061485f565b91505b6001016119c4565b50600e546107e29060019063ffffffff16614842565b611b2a612ad9565b60005b81811015610697576000838383818110611b4957611b496147e4565b9050602002810190611b5b9190614cd8565b611b6490614d1b565b90506000611b7a82600001518360200151610446565b9050611b87600382612e56565b611bc0576040517febf5255100000000000000000000000000000000000000000000000000000000815260048101829052602401610609565b611bca8183612e62565b5050600101611b2d565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610c97576000848483818110611c0e57611c0e6147e4565b9050602002810190611c209190614b8d565b611c2990614bc1565b805163ffffffff166000908152600b602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff168252600181018054959650939491939092840191611c7f90614897565b80601f0160208091040260200160405190810160405280929190818152602001828054611cab90614897565b8015611cf85780601f10611ccd57610100808354040283529160200191611cf8565b820191906000526020600020905b815481529060010190602001808311611cdb57829003601f168201915b50505091909252505081519192505073ffffffffffffffffffffffffffffffffffffffff16611d5e5781516040517fadd9ae1e00000000000000000000000000000000000000000000000000000000815263ffffffff9091166004820152602401610609565b83158015611d835750805173ffffffffffffffffffffffffffffffffffffffff163314155b15611dbc576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610609565b6040808301516000908152600c60205220600181015415611e115782604001516040517f5461848300000000000000000000000000000000000000000000000000000000815260040161060991815260200190565b6040830151611e545782604001516040517f64e2ee9200000000000000000000000000000000000000000000000000000000815260040161060991815260200190565b60208301511580611e7157506020830151611e7190600790612abe565b15611ea8576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608301518051600003611eea57806040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106099190614c94565b81548290600490611f0890640100000000900463ffffffff166148ea565b82546101009290920a63ffffffff818102199093169183160217909155825464010000000090041660005b8251811015611fde57611f51838281518110611557576115576147e4565b611f8957826040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106099190614c94565b611fd5838281518110611f9e57611f9e6147e4565b60200260200101518560030160008563ffffffff1663ffffffff168152602001908152602001600020612e5690919063ffffffff16565b50600101611f33565b50845183547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff918216178455604086015160028501556020860151600185018190556120349160079190612e5616565b50604085015161204690600990612e56565b50845160408087015160208089015183519283529082015263ffffffff909216917f74becb12a5e8fd0e98077d02dfba8f647c9670c9df177e42c2418cf17a636f05910160405180910390a25050505050806001019050611bf2565b60015473ffffffffffffffffffffffffffffffffffffffff163314612123576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610609565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b8281146121e2576040517fab8b67c60000000000000000000000000000000000000000000000000000000081526004810184905260248101829052604401610609565b6000805473ffffffffffffffffffffffffffffffffffffffff16905b848110156124d757600086868381811061221a5761221a6147e4565b905060200201602081019061222f91906141eb565b63ffffffff81166000908152600b6020526040902080549192509073ffffffffffffffffffffffffffffffffffffffff1661229e576040517fadd9ae1e00000000000000000000000000000000000000000000000000000000815263ffffffff83166004820152602401610609565b60008686858181106122b2576122b26147e4565b90506020028101906122c4919061490d565b6122cd9061494b565b805190915073ffffffffffffffffffffffffffffffffffffffff1661231e576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b815473ffffffffffffffffffffffffffffffffffffffff16331480159061235b57503373ffffffffffffffffffffffffffffffffffffffff861614155b15612394576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610609565b8051825473ffffffffffffffffffffffffffffffffffffffff908116911614158061241057506020808201516040516123cd9201613f76565b60405160208183030381529060405280519060200120826001016040516020016123f79190614dc1565b6040516020818303038152906040528051906020012014155b156124c957805182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9091161782556020810151600183019061246a9082614a05565b50806000015173ffffffffffffffffffffffffffffffffffffffff168363ffffffff167f86f41145bde5dd7f523305452e4aad3685508c181432ec733d5f345009358a2883602001516040516124c09190613f76565b60405180910390a35b5050508060010190506121fe565b505050505050565b6125206040805160e0810182526000808252606060208301819052928201839052909182019081526020016000815260006020820181905260409091015290565b6040805160e0810182528381526000848152600260209081529290208054919283019161254c90614897565b80601f016020809104026020016040519081016040528092919081815260200182805461257890614897565b80156125c55780601f1061259a576101008083540402835291602001916125c5565b820191906000526020600020905b8154815290600101906020018083116125a857829003601f168201915b505050505081526020016002600085815260200190815260200160002060010180546125f090614897565b80601f016020809104026020016040519081016040528092919081815260200182805461261c90614897565b80156126695780601f1061263e57610100808354040283529160200191612669565b820191906000526020600020905b81548152906001019060200180831161264c57829003601f168201915b50505091835250506000848152600260208181526040909220015491019060ff16600381111561269b5761269b61444c565b815260008481526002602081815260409092200154910190610100900460ff1660018111156126cc576126cc61444c565b81526000848152600260208181526040928390209091015462010000900473ffffffffffffffffffffffffffffffffffffffff169083015201612710600585612abe565b1515905292915050565b612722612ad9565b63ffffffff8089166000908152600d6020526040812054640100000000900490911690819003612786576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff8a166004820152602401610609565b6127d8888888886040518060a001604052808f63ffffffff168152602001876127ae906148ea565b97508763ffffffff1681526020018a1515815260200189151581526020018860ff168152506130f6565b505050505050505050565b6127eb612ad9565b600e805460009164010000000090910463ffffffff1690600461280d836148ea565b82546101009290920a63ffffffff81810219909316918316021790915581166000818152600d602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001684179055815160a08101835292835260019083015286151590820152841515606082015260ff841660808201529091506128a39089908990899089906130f6565b5050505050505050565b606060006128bb6003612e49565b90506000815167ffffffffffffffff8111156128d9576128d9613d25565b60405190808252806020026020018201604052801561294b57816020015b6129386040805160e0810182526000808252606060208301819052928201839052909182019081526020016000815260006020820181905260409091015290565b8152602001906001900390816128f75790505b50905060005b82518110156107f25761297c83828151811061296f5761296f6147e4565b60200260200101516124df565b82828151811061298e5761298e6147e4565b6020908102919091010152600101612951565b606060006129af6009612e49565b90506000815167ffffffffffffffff8111156129cd576129cd613d25565b604051908082528060200260200182016040528015612a5457816020015b6040805160e081018252600080825260208083018290529282018190526060808301829052608083019190915260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816129eb5790505b50905060005b82518110156107f257612a85838281518110612a7857612a786147e4565b6020026020010151611048565b828281518110612a9757612a976147e4565b6020908102919091010152600101612a5a565b612ab2612ad9565b612abb8161391a565b50565b600081815260018301602052604081205415155b9392505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612b5a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610609565b565b6000610474825490565b6000612ad28383613a0f565b6000612ad28383613a39565b6040805160e0810182526000808252602080830182905282840182905260608084018390526080840183905260a0840181905260c084015263ffffffff8581168352600d8252848320805464010000000090049091168084526001909101825284832060028101805487518186028101860190985280885295969295919493909190830182828015612c2f57602002820191906000526020600020905b815481526020019060010190808311612c1b575b505050505090506000815167ffffffffffffffff811115612c5257612c52613d25565b604051908082528060200260200182016040528015612c9857816020015b604080518082019091526000815260606020820152815260200190600190039081612c705790505b50905060005b8151811015612db0576040518060400160405280848381518110612cc457612cc46147e4565b60200260200101518152602001856003016000868581518110612ce957612ce96147e4565b602002602001015181526020019081526020016000208054612d0a90614897565b80601f0160208091040260200160405190810160405280929190818152602001828054612d3690614897565b8015612d835780601f10612d5857610100808354040283529160200191612d83565b820191906000526020600020905b815481529060010190602001808311612d6657829003601f168201915b5050505050815250828281518110612d9d57612d9d6147e4565b6020908102919091010152600101612c9e565b506040805160e08101825263ffffffff8089166000818152600d6020818152868320548086168752948b168187015260ff680100000000000000008604811697870197909752690100000000000000000085048716151560608701529290915290526a010000000000000000000090049091161515608082015260a08101612e3785612e49565b81526020019190915295945050505050565b60606000612ad283613b2c565b6000612ad28383613b88565b608081015173ffffffffffffffffffffffffffffffffffffffff1615612fb057608081015173ffffffffffffffffffffffffffffffffffffffff163b1580612f5b575060808101516040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527f78bea72100000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff909116906301ffc9a790602401602060405180830381865afa158015612f35573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f599190614e6f565b155b15612fb05760808101516040517fabb5e3fd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610609565b600082815260026020526040902081518291908190612fcf9082614a05565b5060208201516001820190612fe49082614a05565b5060408201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660018360038111156130265761302661444c565b021790555060608201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600181111561306d5761306d61444c565b0217905550608091909101516002909101805473ffffffffffffffffffffffffffffffffffffffff90921662010000027fffffffffffffffffffff0000000000000000000000000000000000000000ffff90921691909117905560405182907f04f0a9bcf3f3a3b42a4d7ca081119755f82ebe43e0d30c8f7292c4fe0dc4a2ae90600090a25050565b805163ffffffff9081166000908152600d602090815260408083208286015190941683526001909301905220608082015160ff161580613148575060808201518590613143906001614e8c565b60ff16115b156131915760808201516040517f25b4d61800000000000000000000000000000000000000000000000000000000815260ff909116600482015260248101869052604401610609565b6001826020015163ffffffff16111561327957815163ffffffff166000908152600d6020908152604082209084015160019182019183916131d29190614842565b63ffffffff1663ffffffff168152602001908152602001600020905060005b6131fa82612b5c565b81101561327657613229846000015163ffffffff16600c60006105928587600001612b6690919063ffffffff16565b50600c60006132388484612b66565b8152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff1690556001016131f1565b50505b60005b858110156134b3576132a9878783818110613299576132996147e4565b8592602090910201359050612e56565b61330a5782518787838181106132c1576132c16147e4565b6040517f636e405700000000000000000000000000000000000000000000000000000000815263ffffffff90941660048501526020029190910135602483015250604401610609565b82606001511561346157825163ffffffff16600c6000898985818110613332576133326147e4565b602090810292909201358352508101919091526040016000205468010000000000000000900463ffffffff16148015906133ac5750600c600088888481811061337d5761337d6147e4565b602090810292909201358352508101919091526040016000205468010000000000000000900463ffffffff1615155b1561340e5782518787838181106133c5576133c56147e4565b6040517f60b9df7300000000000000000000000000000000000000000000000000000000815263ffffffff90941660048501526020029190910135602483015250604401610609565b8251600c6000898985818110613426576134266147e4565b90506020020135815260200190815260200160002060000160086101000a81548163ffffffff021916908363ffffffff1602179055506134ab565b82516134a99063ffffffff16600c60008a8a86818110613483576134836147e4565b905060200201358152602001908152602001600020600401612e5690919063ffffffff16565b505b60010161327c565b5060005b8381101561378f57368585838181106134d2576134d26147e4565b90506020028101906134e4919061490d565b90506134f260038235612abe565b61352b576040517fe181733f00000000000000000000000000000000000000000000000000000000815281356004820152602401610609565b61353760058235612abe565b15613571576040517ff7d7a29400000000000000000000000000000000000000000000000000000000815281356004820152602401610609565b803560009081526003840160205260408120805461358e90614897565b905011156135da5783516040517f3927d08000000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015281356024820152604401610609565b60005b878110156136e4576136818235600c60008c8c86818110613600576136006147e4565b9050602002013581526020019081526020016000206003016000600c60008e8e88818110613630576136306147e4565b90506020020135815260200190815260200160002060000160049054906101000a900463ffffffff1663ffffffff1663ffffffff168152602001908152602001600020612abe90919063ffffffff16565b6136dc57888882818110613697576136976147e4565b6040517fa7e792500000000000000000000000000000000000000000000000000000000081526020909102929092013560048301525082356024820152604401610609565b6001016135dd565b506002830180546001810182556000918252602091829020833591015561370d90820182614ea5565b8235600090815260038601602052604090209161372b919083614f0a565b50835160208086015161378692918435908c908c9061374c90880188614ea5565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250613bd792505050565b506001016134b7565b50604080830151835163ffffffff9081166000908152600d602090815284822080549415156901000000000000000000027fffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffff90951694909417909355606086015186518316825284822080549115156a0100000000000000000000027fffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffff9092169190911790556080860151865183168252848220805460ff9290921668010000000000000000027fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff909216919091179055918501805186518316845292849020805493909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9093169290921790558351905191517ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c15817036519261390a929163ffffffff92831681529116602082015260400190565b60405180910390a1505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603613999576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610609565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110613a2657613a266147e4565b9060005260206000200154905092915050565b60008181526001830160205260408120548015613b22576000613a5d600183615025565b8554909150600090613a7190600190615025565b9050818114613ad6576000866000018281548110613a9157613a916147e4565b9060005260206000200154905080876000018481548110613ab457613ab46147e4565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613ae757613ae7615038565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610474565b6000915050610474565b606081600001805480602002602001604051908101604052809291908181526020018280548015613b7c57602002820191906000526020600020905b815481526020019060010190808311613b68575b50505050509050919050565b6000818152600183016020526040812054613bcf57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610474565b506000610474565b6000848152600260208190526040909120015462010000900473ffffffffffffffffffffffffffffffffffffffff16156124d757600084815260026020819052604091829020015490517ffba64a7c0000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063fba64a7c90613c7e908690869086908b908d90600401615067565b600060405180830381600087803b158015613c9857600080fd5b505af1158015613cac573d6000803e3d6000fd5b50505050505050505050565b508054613cc490614897565b6000825580601f10613cd4575050565b601f016020900490600052602060002090810190612abb9190613d0c565b5080546000825590600052602060002090810190612abb91905b5b80821115613d215760008155600101613d0d565b5090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516080810167ffffffffffffffff81118282101715613d7757613d77613d25565b60405290565b60405160a0810167ffffffffffffffff81118282101715613d7757613d77613d25565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613de757613de7613d25565b604052919050565b600067ffffffffffffffff821115613e0957613e09613d25565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112613e4657600080fd5b8135613e59613e5482613def565b613da0565b818152846020838601011115613e6e57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215613e9e57600080fd5b823567ffffffffffffffff80821115613eb657600080fd5b613ec286838701613e35565b93506020850135915080821115613ed857600080fd5b50613ee585828601613e35565b9150509250929050565b600060208284031215613f0157600080fd5b5035919050565b60005b83811015613f23578181015183820152602001613f0b565b50506000910152565b60008151808452613f44816020860160208601613f08565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000612ad26020830184613f2c565b60008083601f840112613f9b57600080fd5b50813567ffffffffffffffff811115613fb357600080fd5b6020830191508360208260051b850101111561104157600080fd5b60008060208385031215613fe157600080fd5b823567ffffffffffffffff811115613ff857600080fd5b61400485828601613f89565b90969095509350505050565b60008151808452602080850194506020840160005b8381101561404157815187529582019590820190600101614025565b509495945050505050565b600082825180855260208086019550808260051b84010181860160005b848110156140c9578583037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001895281518051845284015160408585018190526140b581860183613f2c565b9a86019a9450505090830190600101614069565b5090979650505050505050565b600063ffffffff8083511684528060208401511660208501525060ff604083015116604084015260608201511515606084015260808201511515608084015260a082015160e060a085015261412e60e0850182614010565b905060c083015184820360c0860152614147828261404c565b95945050505050565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b828110156141c5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526141b38583516140d6565b94509285019290850190600101614179565b5092979650505050505050565b803563ffffffff811681146141e657600080fd5b919050565b6000602082840312156141fd57600080fd5b612ad2826141d2565b73ffffffffffffffffffffffffffffffffffffffff8151168252600060208201516040602085015261423b6040850182613f2c565b949350505050565b602081526000612ad26020830184614206565b602081526000612ad260208301846140d6565b6000806040838503121561427c57600080fd5b614285836141d2565b946020939093013593505050565b6040815260006142a66040830185613f2c565b82810360208401526141478185613f2c565b600063ffffffff808351168452602081818501511681860152816040850151166040860152606084015160608601526080840151608086015260a0840151915060e060a086015261430c60e0860183614010565b60c08581015187830391880191909152805180835290830193506000918301905b8083101561434d578451825293830193600192909201919083019061432d565b509695505050505050565b602081526000612ad260208301846142b8565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b828110156141c5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526143ce858351614206565b94509285019290850190600101614394565b600080600080604085870312156143f657600080fd5b843567ffffffffffffffff8082111561440e57600080fd5b61441a88838901613f89565b9096509450602087013591508082111561443357600080fd5b5061444087828801613f89565b95989497509550505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b805182526000602082015160e0602085015261449a60e0850182613f2c565b9050604083015184820360408601526144b38282613f2c565b9150506060830151600481106144cb576144cb61444c565b60608501526080830151600281106144e5576144e561444c565b8060808601525060a083015161451360a086018273ffffffffffffffffffffffffffffffffffffffff169052565b5060c083015161452760c086018215159052565b509392505050565b602081526000612ad2602083018461447b565b8015158114612abb57600080fd5b803560ff811681146141e657600080fd5b60008060008060008060008060c0898b03121561457d57600080fd5b614586896141d2565b9750602089013567ffffffffffffffff808211156145a357600080fd5b6145af8c838d01613f89565b909950975060408b01359150808211156145c857600080fd5b506145d58b828c01613f89565b90965094505060608901356145e981614542565b925060808901356145f981614542565b915061460760a08a01614550565b90509295985092959890939650565b600080600080600080600060a0888a03121561463157600080fd5b873567ffffffffffffffff8082111561464957600080fd5b6146558b838c01613f89565b909950975060208a013591508082111561466e57600080fd5b5061467b8a828b01613f89565b909650945050604088013561468f81614542565b9250606088013561469f81614542565b91506146ad60808901614550565b905092959891949750929550565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b828110156141c5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc088860301845261471e85835161447b565b945092850192908501906001016146e4565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b828110156141c5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526147938583516142b8565b94509285019290850190600101614759565b803573ffffffffffffffffffffffffffffffffffffffff811681146141e657600080fd5b6000602082840312156147db57600080fd5b612ad2826147a5565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b63ffffffff8281168282160390808211156107f2576107f2614813565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361489057614890614813565b5060010190565b600181811c908216806148ab57607f821691505b6020821081036148e4577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600063ffffffff80831681810361490357614903614813565b6001019392505050565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc183360301811261494157600080fd5b9190910192915050565b60006040823603121561495d57600080fd5b6040516040810167ffffffffffffffff828210818311171561498157614981613d25565b8160405261498e856147a5565b835260208501359150808211156149a457600080fd5b506149b136828601613e35565b60208301525092915050565b601f821115610697576000816000526020600020601f850160051c810160208610156149e65750805b601f850160051c820191505b818110156124d7578281556001016149f2565b815167ffffffffffffffff811115614a1f57614a1f613d25565b614a3381614a2d8454614897565b846149bd565b602080601f831160018114614a865760008415614a505750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556124d7565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015614ad357888601518255948401946001909101908401614ab4565b5085821015614b0f57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600060208284031215614b3157600080fd5b815167ffffffffffffffff811115614b4857600080fd5b8201601f81018413614b5957600080fd5b8051614b67613e5482613def565b818152856020838501011115614b7c57600080fd5b614147826020830160208601613f08565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8183360301811261494157600080fd5b600060808236031215614bd357600080fd5b614bdb613d54565b614be4836141d2565b81526020808401358183015260408401356040830152606084013567ffffffffffffffff80821115614c1557600080fd5b9085019036601f830112614c2857600080fd5b813581811115614c3a57614c3a613d25565b8060051b9150614c4b848301613da0565b8181529183018401918481019036841115614c6557600080fd5b938501935b83851015614c8357843582529385019390850190614c6a565b606087015250939695505050505050565b6020808252825182820181905260009190848201906040850190845b81811015614ccc57835183529284019291840191600101614cb0565b50909695505050505050565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6183360301811261494157600080fd5b8035600281106141e657600080fd5b600060a08236031215614d2d57600080fd5b614d35613d7d565b823567ffffffffffffffff80821115614d4d57600080fd5b614d5936838701613e35565b83526020850135915080821115614d6f57600080fd5b50614d7c36828601613e35565b602083015250604083013560048110614d9457600080fd5b6040820152614da560608401614d0c565b6060820152614db6608084016147a5565b608082015292915050565b6000602080835260008454614dd581614897565b8060208701526040600180841660008114614df75760018114614e3157614e61565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00851660408a0152604084151560051b8a01019550614e61565b89600052602060002060005b85811015614e585781548b8201860152908301908801614e3d565b8a016040019650505b509398975050505050505050565b600060208284031215614e8157600080fd5b8151612ad281614542565b60ff818116838216019081111561047457610474614813565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112614eda57600080fd5b83018035915067ffffffffffffffff821115614ef557600080fd5b60200191503681900382131561104157600080fd5b67ffffffffffffffff831115614f2257614f22613d25565b614f3683614f308354614897565b836149bd565b6000601f841160018114614f885760008515614f525750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b17835561501e565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015614fd75786850135825560209485019460019092019101614fb7565b5086821015615012577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b8181038181111561047457610474614813565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6080815284608082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8611156150a057600080fd5b8560051b808860a0850137820182810360a090810160208501526150c690820187613f2c565b91505063ffffffff8085166040840152808416606084015250969550505050505056fea164736f6c6343000818000a",
}

var CapabilitiesRegistryABI = CapabilitiesRegistryMetaData.ABI

var CapabilitiesRegistryBin = CapabilitiesRegistryMetaData.Bin

func DeployCapabilitiesRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CapabilitiesRegistry, error) {
	parsed, err := CapabilitiesRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CapabilitiesRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CapabilitiesRegistry{address: address, abi: *parsed, CapabilitiesRegistryCaller: CapabilitiesRegistryCaller{contract: contract}, CapabilitiesRegistryTransactor: CapabilitiesRegistryTransactor{contract: contract}, CapabilitiesRegistryFilterer: CapabilitiesRegistryFilterer{contract: contract}}, nil
}

type CapabilitiesRegistry struct {
	address common.Address
	abi     abi.ABI
	CapabilitiesRegistryCaller
	CapabilitiesRegistryTransactor
	CapabilitiesRegistryFilterer
}

type CapabilitiesRegistryCaller struct {
	contract *bind.BoundContract
}

type CapabilitiesRegistryTransactor struct {
	contract *bind.BoundContract
}

type CapabilitiesRegistryFilterer struct {
	contract *bind.BoundContract
}

type CapabilitiesRegistrySession struct {
	Contract     *CapabilitiesRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CapabilitiesRegistryCallerSession struct {
	Contract *CapabilitiesRegistryCaller
	CallOpts bind.CallOpts
}

type CapabilitiesRegistryTransactorSession struct {
	Contract     *CapabilitiesRegistryTransactor
	TransactOpts bind.TransactOpts
}

type CapabilitiesRegistryRaw struct {
	Contract *CapabilitiesRegistry
}

type CapabilitiesRegistryCallerRaw struct {
	Contract *CapabilitiesRegistryCaller
}

type CapabilitiesRegistryTransactorRaw struct {
	Contract *CapabilitiesRegistryTransactor
}

func NewCapabilitiesRegistry(address common.Address, backend bind.ContractBackend) (*CapabilitiesRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(CapabilitiesRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCapabilitiesRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistry{address: address, abi: abi, CapabilitiesRegistryCaller: CapabilitiesRegistryCaller{contract: contract}, CapabilitiesRegistryTransactor: CapabilitiesRegistryTransactor{contract: contract}, CapabilitiesRegistryFilterer: CapabilitiesRegistryFilterer{contract: contract}}, nil
}

func NewCapabilitiesRegistryCaller(address common.Address, caller bind.ContractCaller) (*CapabilitiesRegistryCaller, error) {
	contract, err := bindCapabilitiesRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryCaller{contract: contract}, nil
}

func NewCapabilitiesRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*CapabilitiesRegistryTransactor, error) {
	contract, err := bindCapabilitiesRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryTransactor{contract: contract}, nil
}

func NewCapabilitiesRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*CapabilitiesRegistryFilterer, error) {
	contract, err := bindCapabilitiesRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryFilterer{contract: contract}, nil
}

func bindCapabilitiesRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CapabilitiesRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CapabilitiesRegistry.Contract.CapabilitiesRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.CapabilitiesRegistryTransactor.contract.Transfer(opts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.CapabilitiesRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CapabilitiesRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.contract.Transfer(opts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) GetCapabilities(opts *bind.CallOpts) ([]CapabilitiesRegistryCapabilityInfo, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "getCapabilities")

	if err != nil {
		return *new([]CapabilitiesRegistryCapabilityInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]CapabilitiesRegistryCapabilityInfo)).(*[]CapabilitiesRegistryCapabilityInfo)

	return out0, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) GetCapabilities() ([]CapabilitiesRegistryCapabilityInfo, error) {
	return _CapabilitiesRegistry.Contract.GetCapabilities(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) GetCapabilities() ([]CapabilitiesRegistryCapabilityInfo, error) {
	return _CapabilitiesRegistry.Contract.GetCapabilities(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) GetCapability(opts *bind.CallOpts, hashedId [32]byte) (CapabilitiesRegistryCapabilityInfo, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "getCapability", hashedId)

	if err != nil {
		return *new(CapabilitiesRegistryCapabilityInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(CapabilitiesRegistryCapabilityInfo)).(*CapabilitiesRegistryCapabilityInfo)

	return out0, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) GetCapability(hashedId [32]byte) (CapabilitiesRegistryCapabilityInfo, error) {
	return _CapabilitiesRegistry.Contract.GetCapability(&_CapabilitiesRegistry.CallOpts, hashedId)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) GetCapability(hashedId [32]byte) (CapabilitiesRegistryCapabilityInfo, error) {
	return _CapabilitiesRegistry.Contract.GetCapability(&_CapabilitiesRegistry.CallOpts, hashedId)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) GetCapabilityConfigs(opts *bind.CallOpts, donId uint32, capabilityId [32]byte) ([]byte, []byte, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "getCapabilityConfigs", donId, capabilityId)

	if err != nil {
		return *new([]byte), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) GetCapabilityConfigs(donId uint32, capabilityId [32]byte) ([]byte, []byte, error) {
	return _CapabilitiesRegistry.Contract.GetCapabilityConfigs(&_CapabilitiesRegistry.CallOpts, donId, capabilityId)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) GetCapabilityConfigs(donId uint32, capabilityId [32]byte) ([]byte, []byte, error) {
	return _CapabilitiesRegistry.Contract.GetCapabilityConfigs(&_CapabilitiesRegistry.CallOpts, donId, capabilityId)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) GetDON(opts *bind.CallOpts, donId uint32) (CapabilitiesRegistryDONInfo, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "getDON", donId)

	if err != nil {
		return *new(CapabilitiesRegistryDONInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(CapabilitiesRegistryDONInfo)).(*CapabilitiesRegistryDONInfo)

	return out0, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) GetDON(donId uint32) (CapabilitiesRegistryDONInfo, error) {
	return _CapabilitiesRegistry.Contract.GetDON(&_CapabilitiesRegistry.CallOpts, donId)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) GetDON(donId uint32) (CapabilitiesRegistryDONInfo, error) {
	return _CapabilitiesRegistry.Contract.GetDON(&_CapabilitiesRegistry.CallOpts, donId)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) GetDONs(opts *bind.CallOpts) ([]CapabilitiesRegistryDONInfo, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "getDONs")

	if err != nil {
		return *new([]CapabilitiesRegistryDONInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]CapabilitiesRegistryDONInfo)).(*[]CapabilitiesRegistryDONInfo)

	return out0, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) GetDONs() ([]CapabilitiesRegistryDONInfo, error) {
	return _CapabilitiesRegistry.Contract.GetDONs(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) GetDONs() ([]CapabilitiesRegistryDONInfo, error) {
	return _CapabilitiesRegistry.Contract.GetDONs(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) GetHashedCapabilityId(opts *bind.CallOpts, labelledName string, version string) ([32]byte, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "getHashedCapabilityId", labelledName, version)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) GetHashedCapabilityId(labelledName string, version string) ([32]byte, error) {
	return _CapabilitiesRegistry.Contract.GetHashedCapabilityId(&_CapabilitiesRegistry.CallOpts, labelledName, version)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) GetHashedCapabilityId(labelledName string, version string) ([32]byte, error) {
	return _CapabilitiesRegistry.Contract.GetHashedCapabilityId(&_CapabilitiesRegistry.CallOpts, labelledName, version)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) GetNode(opts *bind.CallOpts, p2pId [32]byte) (CapabilitiesRegistryNodeInfo, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "getNode", p2pId)

	if err != nil {
		return *new(CapabilitiesRegistryNodeInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(CapabilitiesRegistryNodeInfo)).(*CapabilitiesRegistryNodeInfo)

	return out0, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) GetNode(p2pId [32]byte) (CapabilitiesRegistryNodeInfo, error) {
	return _CapabilitiesRegistry.Contract.GetNode(&_CapabilitiesRegistry.CallOpts, p2pId)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) GetNode(p2pId [32]byte) (CapabilitiesRegistryNodeInfo, error) {
	return _CapabilitiesRegistry.Contract.GetNode(&_CapabilitiesRegistry.CallOpts, p2pId)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) GetNodeOperator(opts *bind.CallOpts, nodeOperatorId uint32) (CapabilitiesRegistryNodeOperator, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "getNodeOperator", nodeOperatorId)

	if err != nil {
		return *new(CapabilitiesRegistryNodeOperator), err
	}

	out0 := *abi.ConvertType(out[0], new(CapabilitiesRegistryNodeOperator)).(*CapabilitiesRegistryNodeOperator)

	return out0, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) GetNodeOperator(nodeOperatorId uint32) (CapabilitiesRegistryNodeOperator, error) {
	return _CapabilitiesRegistry.Contract.GetNodeOperator(&_CapabilitiesRegistry.CallOpts, nodeOperatorId)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) GetNodeOperator(nodeOperatorId uint32) (CapabilitiesRegistryNodeOperator, error) {
	return _CapabilitiesRegistry.Contract.GetNodeOperator(&_CapabilitiesRegistry.CallOpts, nodeOperatorId)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) GetNodeOperators(opts *bind.CallOpts) ([]CapabilitiesRegistryNodeOperator, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "getNodeOperators")

	if err != nil {
		return *new([]CapabilitiesRegistryNodeOperator), err
	}

	out0 := *abi.ConvertType(out[0], new([]CapabilitiesRegistryNodeOperator)).(*[]CapabilitiesRegistryNodeOperator)

	return out0, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) GetNodeOperators() ([]CapabilitiesRegistryNodeOperator, error) {
	return _CapabilitiesRegistry.Contract.GetNodeOperators(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) GetNodeOperators() ([]CapabilitiesRegistryNodeOperator, error) {
	return _CapabilitiesRegistry.Contract.GetNodeOperators(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) GetNodes(opts *bind.CallOpts) ([]CapabilitiesRegistryNodeInfo, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "getNodes")

	if err != nil {
		return *new([]CapabilitiesRegistryNodeInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]CapabilitiesRegistryNodeInfo)).(*[]CapabilitiesRegistryNodeInfo)

	return out0, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) GetNodes() ([]CapabilitiesRegistryNodeInfo, error) {
	return _CapabilitiesRegistry.Contract.GetNodes(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) GetNodes() ([]CapabilitiesRegistryNodeInfo, error) {
	return _CapabilitiesRegistry.Contract.GetNodes(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) IsCapabilityDeprecated(opts *bind.CallOpts, hashedCapabilityId [32]byte) (bool, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "isCapabilityDeprecated", hashedCapabilityId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) IsCapabilityDeprecated(hashedCapabilityId [32]byte) (bool, error) {
	return _CapabilitiesRegistry.Contract.IsCapabilityDeprecated(&_CapabilitiesRegistry.CallOpts, hashedCapabilityId)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) IsCapabilityDeprecated(hashedCapabilityId [32]byte) (bool, error) {
	return _CapabilitiesRegistry.Contract.IsCapabilityDeprecated(&_CapabilitiesRegistry.CallOpts, hashedCapabilityId)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) Owner() (common.Address, error) {
	return _CapabilitiesRegistry.Contract.Owner(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) Owner() (common.Address, error) {
	return _CapabilitiesRegistry.Contract.Owner(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) TypeAndVersion() (string, error) {
	return _CapabilitiesRegistry.Contract.TypeAndVersion(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) TypeAndVersion() (string, error) {
	return _CapabilitiesRegistry.Contract.TypeAndVersion(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.AcceptOwnership(&_CapabilitiesRegistry.TransactOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.AcceptOwnership(&_CapabilitiesRegistry.TransactOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) AddCapabilities(opts *bind.TransactOpts, capabilities []CapabilitiesRegistryCapability) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "addCapabilities", capabilities)
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) AddCapabilities(capabilities []CapabilitiesRegistryCapability) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.AddCapabilities(&_CapabilitiesRegistry.TransactOpts, capabilities)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) AddCapabilities(capabilities []CapabilitiesRegistryCapability) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.AddCapabilities(&_CapabilitiesRegistry.TransactOpts, capabilities)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) AddDON(opts *bind.TransactOpts, nodes [][32]byte, capabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "addDON", nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) AddDON(nodes [][32]byte, capabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.AddDON(&_CapabilitiesRegistry.TransactOpts, nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) AddDON(nodes [][32]byte, capabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.AddDON(&_CapabilitiesRegistry.TransactOpts, nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) AddNodeOperators(opts *bind.TransactOpts, nodeOperators []CapabilitiesRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "addNodeOperators", nodeOperators)
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) AddNodeOperators(nodeOperators []CapabilitiesRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.AddNodeOperators(&_CapabilitiesRegistry.TransactOpts, nodeOperators)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) AddNodeOperators(nodeOperators []CapabilitiesRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.AddNodeOperators(&_CapabilitiesRegistry.TransactOpts, nodeOperators)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) AddNodes(opts *bind.TransactOpts, nodes []CapabilitiesRegistryNodeParams) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "addNodes", nodes)
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) AddNodes(nodes []CapabilitiesRegistryNodeParams) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.AddNodes(&_CapabilitiesRegistry.TransactOpts, nodes)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) AddNodes(nodes []CapabilitiesRegistryNodeParams) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.AddNodes(&_CapabilitiesRegistry.TransactOpts, nodes)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) DeprecateCapabilities(opts *bind.TransactOpts, hashedCapabilityIds [][32]byte) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "deprecateCapabilities", hashedCapabilityIds)
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) DeprecateCapabilities(hashedCapabilityIds [][32]byte) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.DeprecateCapabilities(&_CapabilitiesRegistry.TransactOpts, hashedCapabilityIds)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) DeprecateCapabilities(hashedCapabilityIds [][32]byte) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.DeprecateCapabilities(&_CapabilitiesRegistry.TransactOpts, hashedCapabilityIds)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) RemoveDONs(opts *bind.TransactOpts, donIds []uint32) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "removeDONs", donIds)
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) RemoveDONs(donIds []uint32) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.RemoveDONs(&_CapabilitiesRegistry.TransactOpts, donIds)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) RemoveDONs(donIds []uint32) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.RemoveDONs(&_CapabilitiesRegistry.TransactOpts, donIds)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) RemoveNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []uint32) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "removeNodeOperators", nodeOperatorIds)
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) RemoveNodeOperators(nodeOperatorIds []uint32) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.RemoveNodeOperators(&_CapabilitiesRegistry.TransactOpts, nodeOperatorIds)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) RemoveNodeOperators(nodeOperatorIds []uint32) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.RemoveNodeOperators(&_CapabilitiesRegistry.TransactOpts, nodeOperatorIds)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) RemoveNodes(opts *bind.TransactOpts, removedNodeP2PIds [][32]byte) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "removeNodes", removedNodeP2PIds)
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) RemoveNodes(removedNodeP2PIds [][32]byte) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.RemoveNodes(&_CapabilitiesRegistry.TransactOpts, removedNodeP2PIds)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) RemoveNodes(removedNodeP2PIds [][32]byte) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.RemoveNodes(&_CapabilitiesRegistry.TransactOpts, removedNodeP2PIds)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.TransferOwnership(&_CapabilitiesRegistry.TransactOpts, to)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.TransferOwnership(&_CapabilitiesRegistry.TransactOpts, to)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) UpdateDON(opts *bind.TransactOpts, donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "updateDON", donId, nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) UpdateDON(donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.UpdateDON(&_CapabilitiesRegistry.TransactOpts, donId, nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) UpdateDON(donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.UpdateDON(&_CapabilitiesRegistry.TransactOpts, donId, nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) UpdateNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []uint32, nodeOperators []CapabilitiesRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "updateNodeOperators", nodeOperatorIds, nodeOperators)
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) UpdateNodeOperators(nodeOperatorIds []uint32, nodeOperators []CapabilitiesRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.UpdateNodeOperators(&_CapabilitiesRegistry.TransactOpts, nodeOperatorIds, nodeOperators)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) UpdateNodeOperators(nodeOperatorIds []uint32, nodeOperators []CapabilitiesRegistryNodeOperator) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.UpdateNodeOperators(&_CapabilitiesRegistry.TransactOpts, nodeOperatorIds, nodeOperators)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) UpdateNodes(opts *bind.TransactOpts, nodes []CapabilitiesRegistryNodeParams) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "updateNodes", nodes)
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) UpdateNodes(nodes []CapabilitiesRegistryNodeParams) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.UpdateNodes(&_CapabilitiesRegistry.TransactOpts, nodes)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) UpdateNodes(nodes []CapabilitiesRegistryNodeParams) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.UpdateNodes(&_CapabilitiesRegistry.TransactOpts, nodes)
}

type CapabilitiesRegistryCapabilityConfiguredIterator struct {
	Event *CapabilitiesRegistryCapabilityConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilitiesRegistryCapabilityConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilitiesRegistryCapabilityConfigured)
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
		it.Event = new(CapabilitiesRegistryCapabilityConfigured)
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

func (it *CapabilitiesRegistryCapabilityConfiguredIterator) Error() error {
	return it.fail
}

func (it *CapabilitiesRegistryCapabilityConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilitiesRegistryCapabilityConfigured struct {
	HashedCapabilityId [32]byte
	Raw                types.Log
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) FilterCapabilityConfigured(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilitiesRegistryCapabilityConfiguredIterator, error) {

	var hashedCapabilityIdRule []interface{}
	for _, hashedCapabilityIdItem := range hashedCapabilityId {
		hashedCapabilityIdRule = append(hashedCapabilityIdRule, hashedCapabilityIdItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.FilterLogs(opts, "CapabilityConfigured", hashedCapabilityIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryCapabilityConfiguredIterator{contract: _CapabilitiesRegistry.contract, event: "CapabilityConfigured", logs: logs, sub: sub}, nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) WatchCapabilityConfigured(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryCapabilityConfigured, hashedCapabilityId [][32]byte) (event.Subscription, error) {

	var hashedCapabilityIdRule []interface{}
	for _, hashedCapabilityIdItem := range hashedCapabilityId {
		hashedCapabilityIdRule = append(hashedCapabilityIdRule, hashedCapabilityIdItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.WatchLogs(opts, "CapabilityConfigured", hashedCapabilityIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilitiesRegistryCapabilityConfigured)
				if err := _CapabilitiesRegistry.contract.UnpackLog(event, "CapabilityConfigured", log); err != nil {
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

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) ParseCapabilityConfigured(log types.Log) (*CapabilitiesRegistryCapabilityConfigured, error) {
	event := new(CapabilitiesRegistryCapabilityConfigured)
	if err := _CapabilitiesRegistry.contract.UnpackLog(event, "CapabilityConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilitiesRegistryCapabilityDeprecatedIterator struct {
	Event *CapabilitiesRegistryCapabilityDeprecated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilitiesRegistryCapabilityDeprecatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilitiesRegistryCapabilityDeprecated)
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
		it.Event = new(CapabilitiesRegistryCapabilityDeprecated)
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

func (it *CapabilitiesRegistryCapabilityDeprecatedIterator) Error() error {
	return it.fail
}

func (it *CapabilitiesRegistryCapabilityDeprecatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilitiesRegistryCapabilityDeprecated struct {
	HashedCapabilityId [32]byte
	Raw                types.Log
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) FilterCapabilityDeprecated(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilitiesRegistryCapabilityDeprecatedIterator, error) {

	var hashedCapabilityIdRule []interface{}
	for _, hashedCapabilityIdItem := range hashedCapabilityId {
		hashedCapabilityIdRule = append(hashedCapabilityIdRule, hashedCapabilityIdItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.FilterLogs(opts, "CapabilityDeprecated", hashedCapabilityIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryCapabilityDeprecatedIterator{contract: _CapabilitiesRegistry.contract, event: "CapabilityDeprecated", logs: logs, sub: sub}, nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) WatchCapabilityDeprecated(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryCapabilityDeprecated, hashedCapabilityId [][32]byte) (event.Subscription, error) {

	var hashedCapabilityIdRule []interface{}
	for _, hashedCapabilityIdItem := range hashedCapabilityId {
		hashedCapabilityIdRule = append(hashedCapabilityIdRule, hashedCapabilityIdItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.WatchLogs(opts, "CapabilityDeprecated", hashedCapabilityIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilitiesRegistryCapabilityDeprecated)
				if err := _CapabilitiesRegistry.contract.UnpackLog(event, "CapabilityDeprecated", log); err != nil {
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

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) ParseCapabilityDeprecated(log types.Log) (*CapabilitiesRegistryCapabilityDeprecated, error) {
	event := new(CapabilitiesRegistryCapabilityDeprecated)
	if err := _CapabilitiesRegistry.contract.UnpackLog(event, "CapabilityDeprecated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilitiesRegistryConfigSetIterator struct {
	Event *CapabilitiesRegistryConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilitiesRegistryConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilitiesRegistryConfigSet)
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
		it.Event = new(CapabilitiesRegistryConfigSet)
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

func (it *CapabilitiesRegistryConfigSetIterator) Error() error {
	return it.fail
}

func (it *CapabilitiesRegistryConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilitiesRegistryConfigSet struct {
	DonId       uint32
	ConfigCount uint32
	Raw         types.Log
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) FilterConfigSet(opts *bind.FilterOpts) (*CapabilitiesRegistryConfigSetIterator, error) {

	logs, sub, err := _CapabilitiesRegistry.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryConfigSetIterator{contract: _CapabilitiesRegistry.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryConfigSet) (event.Subscription, error) {

	logs, sub, err := _CapabilitiesRegistry.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilitiesRegistryConfigSet)
				if err := _CapabilitiesRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) ParseConfigSet(log types.Log) (*CapabilitiesRegistryConfigSet, error) {
	event := new(CapabilitiesRegistryConfigSet)
	if err := _CapabilitiesRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilitiesRegistryNodeAddedIterator struct {
	Event *CapabilitiesRegistryNodeAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilitiesRegistryNodeAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilitiesRegistryNodeAdded)
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
		it.Event = new(CapabilitiesRegistryNodeAdded)
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

func (it *CapabilitiesRegistryNodeAddedIterator) Error() error {
	return it.fail
}

func (it *CapabilitiesRegistryNodeAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilitiesRegistryNodeAdded struct {
	P2pId          [32]byte
	NodeOperatorId uint32
	Signer         [32]byte
	Raw            types.Log
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) FilterNodeAdded(opts *bind.FilterOpts, nodeOperatorId []uint32) (*CapabilitiesRegistryNodeAddedIterator, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.FilterLogs(opts, "NodeAdded", nodeOperatorIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryNodeAddedIterator{contract: _CapabilitiesRegistry.contract, event: "NodeAdded", logs: logs, sub: sub}, nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) WatchNodeAdded(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryNodeAdded, nodeOperatorId []uint32) (event.Subscription, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.WatchLogs(opts, "NodeAdded", nodeOperatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilitiesRegistryNodeAdded)
				if err := _CapabilitiesRegistry.contract.UnpackLog(event, "NodeAdded", log); err != nil {
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

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) ParseNodeAdded(log types.Log) (*CapabilitiesRegistryNodeAdded, error) {
	event := new(CapabilitiesRegistryNodeAdded)
	if err := _CapabilitiesRegistry.contract.UnpackLog(event, "NodeAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilitiesRegistryNodeOperatorAddedIterator struct {
	Event *CapabilitiesRegistryNodeOperatorAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilitiesRegistryNodeOperatorAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilitiesRegistryNodeOperatorAdded)
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
		it.Event = new(CapabilitiesRegistryNodeOperatorAdded)
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

func (it *CapabilitiesRegistryNodeOperatorAddedIterator) Error() error {
	return it.fail
}

func (it *CapabilitiesRegistryNodeOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilitiesRegistryNodeOperatorAdded struct {
	NodeOperatorId uint32
	Admin          common.Address
	Name           string
	Raw            types.Log
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) FilterNodeOperatorAdded(opts *bind.FilterOpts, nodeOperatorId []uint32, admin []common.Address) (*CapabilitiesRegistryNodeOperatorAddedIterator, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.FilterLogs(opts, "NodeOperatorAdded", nodeOperatorIdRule, adminRule)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryNodeOperatorAddedIterator{contract: _CapabilitiesRegistry.contract, event: "NodeOperatorAdded", logs: logs, sub: sub}, nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) WatchNodeOperatorAdded(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryNodeOperatorAdded, nodeOperatorId []uint32, admin []common.Address) (event.Subscription, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.WatchLogs(opts, "NodeOperatorAdded", nodeOperatorIdRule, adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilitiesRegistryNodeOperatorAdded)
				if err := _CapabilitiesRegistry.contract.UnpackLog(event, "NodeOperatorAdded", log); err != nil {
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

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) ParseNodeOperatorAdded(log types.Log) (*CapabilitiesRegistryNodeOperatorAdded, error) {
	event := new(CapabilitiesRegistryNodeOperatorAdded)
	if err := _CapabilitiesRegistry.contract.UnpackLog(event, "NodeOperatorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilitiesRegistryNodeOperatorRemovedIterator struct {
	Event *CapabilitiesRegistryNodeOperatorRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilitiesRegistryNodeOperatorRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilitiesRegistryNodeOperatorRemoved)
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
		it.Event = new(CapabilitiesRegistryNodeOperatorRemoved)
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

func (it *CapabilitiesRegistryNodeOperatorRemovedIterator) Error() error {
	return it.fail
}

func (it *CapabilitiesRegistryNodeOperatorRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilitiesRegistryNodeOperatorRemoved struct {
	NodeOperatorId uint32
	Raw            types.Log
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) FilterNodeOperatorRemoved(opts *bind.FilterOpts, nodeOperatorId []uint32) (*CapabilitiesRegistryNodeOperatorRemovedIterator, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.FilterLogs(opts, "NodeOperatorRemoved", nodeOperatorIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryNodeOperatorRemovedIterator{contract: _CapabilitiesRegistry.contract, event: "NodeOperatorRemoved", logs: logs, sub: sub}, nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) WatchNodeOperatorRemoved(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryNodeOperatorRemoved, nodeOperatorId []uint32) (event.Subscription, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.WatchLogs(opts, "NodeOperatorRemoved", nodeOperatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilitiesRegistryNodeOperatorRemoved)
				if err := _CapabilitiesRegistry.contract.UnpackLog(event, "NodeOperatorRemoved", log); err != nil {
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

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) ParseNodeOperatorRemoved(log types.Log) (*CapabilitiesRegistryNodeOperatorRemoved, error) {
	event := new(CapabilitiesRegistryNodeOperatorRemoved)
	if err := _CapabilitiesRegistry.contract.UnpackLog(event, "NodeOperatorRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilitiesRegistryNodeOperatorUpdatedIterator struct {
	Event *CapabilitiesRegistryNodeOperatorUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilitiesRegistryNodeOperatorUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilitiesRegistryNodeOperatorUpdated)
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
		it.Event = new(CapabilitiesRegistryNodeOperatorUpdated)
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

func (it *CapabilitiesRegistryNodeOperatorUpdatedIterator) Error() error {
	return it.fail
}

func (it *CapabilitiesRegistryNodeOperatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilitiesRegistryNodeOperatorUpdated struct {
	NodeOperatorId uint32
	Admin          common.Address
	Name           string
	Raw            types.Log
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) FilterNodeOperatorUpdated(opts *bind.FilterOpts, nodeOperatorId []uint32, admin []common.Address) (*CapabilitiesRegistryNodeOperatorUpdatedIterator, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.FilterLogs(opts, "NodeOperatorUpdated", nodeOperatorIdRule, adminRule)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryNodeOperatorUpdatedIterator{contract: _CapabilitiesRegistry.contract, event: "NodeOperatorUpdated", logs: logs, sub: sub}, nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) WatchNodeOperatorUpdated(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryNodeOperatorUpdated, nodeOperatorId []uint32, admin []common.Address) (event.Subscription, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.WatchLogs(opts, "NodeOperatorUpdated", nodeOperatorIdRule, adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilitiesRegistryNodeOperatorUpdated)
				if err := _CapabilitiesRegistry.contract.UnpackLog(event, "NodeOperatorUpdated", log); err != nil {
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

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) ParseNodeOperatorUpdated(log types.Log) (*CapabilitiesRegistryNodeOperatorUpdated, error) {
	event := new(CapabilitiesRegistryNodeOperatorUpdated)
	if err := _CapabilitiesRegistry.contract.UnpackLog(event, "NodeOperatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilitiesRegistryNodeRemovedIterator struct {
	Event *CapabilitiesRegistryNodeRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilitiesRegistryNodeRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilitiesRegistryNodeRemoved)
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
		it.Event = new(CapabilitiesRegistryNodeRemoved)
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

func (it *CapabilitiesRegistryNodeRemovedIterator) Error() error {
	return it.fail
}

func (it *CapabilitiesRegistryNodeRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilitiesRegistryNodeRemoved struct {
	P2pId [32]byte
	Raw   types.Log
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) FilterNodeRemoved(opts *bind.FilterOpts) (*CapabilitiesRegistryNodeRemovedIterator, error) {

	logs, sub, err := _CapabilitiesRegistry.contract.FilterLogs(opts, "NodeRemoved")
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryNodeRemovedIterator{contract: _CapabilitiesRegistry.contract, event: "NodeRemoved", logs: logs, sub: sub}, nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) WatchNodeRemoved(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryNodeRemoved) (event.Subscription, error) {

	logs, sub, err := _CapabilitiesRegistry.contract.WatchLogs(opts, "NodeRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilitiesRegistryNodeRemoved)
				if err := _CapabilitiesRegistry.contract.UnpackLog(event, "NodeRemoved", log); err != nil {
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

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) ParseNodeRemoved(log types.Log) (*CapabilitiesRegistryNodeRemoved, error) {
	event := new(CapabilitiesRegistryNodeRemoved)
	if err := _CapabilitiesRegistry.contract.UnpackLog(event, "NodeRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilitiesRegistryNodeUpdatedIterator struct {
	Event *CapabilitiesRegistryNodeUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilitiesRegistryNodeUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilitiesRegistryNodeUpdated)
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
		it.Event = new(CapabilitiesRegistryNodeUpdated)
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

func (it *CapabilitiesRegistryNodeUpdatedIterator) Error() error {
	return it.fail
}

func (it *CapabilitiesRegistryNodeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilitiesRegistryNodeUpdated struct {
	P2pId          [32]byte
	NodeOperatorId uint32
	Signer         [32]byte
	Raw            types.Log
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) FilterNodeUpdated(opts *bind.FilterOpts, nodeOperatorId []uint32) (*CapabilitiesRegistryNodeUpdatedIterator, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.FilterLogs(opts, "NodeUpdated", nodeOperatorIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryNodeUpdatedIterator{contract: _CapabilitiesRegistry.contract, event: "NodeUpdated", logs: logs, sub: sub}, nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) WatchNodeUpdated(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryNodeUpdated, nodeOperatorId []uint32) (event.Subscription, error) {

	var nodeOperatorIdRule []interface{}
	for _, nodeOperatorIdItem := range nodeOperatorId {
		nodeOperatorIdRule = append(nodeOperatorIdRule, nodeOperatorIdItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.WatchLogs(opts, "NodeUpdated", nodeOperatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilitiesRegistryNodeUpdated)
				if err := _CapabilitiesRegistry.contract.UnpackLog(event, "NodeUpdated", log); err != nil {
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

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) ParseNodeUpdated(log types.Log) (*CapabilitiesRegistryNodeUpdated, error) {
	event := new(CapabilitiesRegistryNodeUpdated)
	if err := _CapabilitiesRegistry.contract.UnpackLog(event, "NodeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilitiesRegistryOwnershipTransferRequestedIterator struct {
	Event *CapabilitiesRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilitiesRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilitiesRegistryOwnershipTransferRequested)
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
		it.Event = new(CapabilitiesRegistryOwnershipTransferRequested)
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

func (it *CapabilitiesRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *CapabilitiesRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilitiesRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CapabilitiesRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryOwnershipTransferRequestedIterator{contract: _CapabilitiesRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilitiesRegistryOwnershipTransferRequested)
				if err := _CapabilitiesRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*CapabilitiesRegistryOwnershipTransferRequested, error) {
	event := new(CapabilitiesRegistryOwnershipTransferRequested)
	if err := _CapabilitiesRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CapabilitiesRegistryOwnershipTransferredIterator struct {
	Event *CapabilitiesRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CapabilitiesRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapabilitiesRegistryOwnershipTransferred)
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
		it.Event = new(CapabilitiesRegistryOwnershipTransferred)
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

func (it *CapabilitiesRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *CapabilitiesRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CapabilitiesRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CapabilitiesRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryOwnershipTransferredIterator{contract: _CapabilitiesRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CapabilitiesRegistryOwnershipTransferred)
				if err := _CapabilitiesRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*CapabilitiesRegistryOwnershipTransferred, error) {
	event := new(CapabilitiesRegistryOwnershipTransferred)
	if err := _CapabilitiesRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _CapabilitiesRegistry.abi.Events["CapabilityConfigured"].ID:
		return _CapabilitiesRegistry.ParseCapabilityConfigured(log)
	case _CapabilitiesRegistry.abi.Events["CapabilityDeprecated"].ID:
		return _CapabilitiesRegistry.ParseCapabilityDeprecated(log)
	case _CapabilitiesRegistry.abi.Events["ConfigSet"].ID:
		return _CapabilitiesRegistry.ParseConfigSet(log)
	case _CapabilitiesRegistry.abi.Events["NodeAdded"].ID:
		return _CapabilitiesRegistry.ParseNodeAdded(log)
	case _CapabilitiesRegistry.abi.Events["NodeOperatorAdded"].ID:
		return _CapabilitiesRegistry.ParseNodeOperatorAdded(log)
	case _CapabilitiesRegistry.abi.Events["NodeOperatorRemoved"].ID:
		return _CapabilitiesRegistry.ParseNodeOperatorRemoved(log)
	case _CapabilitiesRegistry.abi.Events["NodeOperatorUpdated"].ID:
		return _CapabilitiesRegistry.ParseNodeOperatorUpdated(log)
	case _CapabilitiesRegistry.abi.Events["NodeRemoved"].ID:
		return _CapabilitiesRegistry.ParseNodeRemoved(log)
	case _CapabilitiesRegistry.abi.Events["NodeUpdated"].ID:
		return _CapabilitiesRegistry.ParseNodeUpdated(log)
	case _CapabilitiesRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _CapabilitiesRegistry.ParseOwnershipTransferRequested(log)
	case _CapabilitiesRegistry.abi.Events["OwnershipTransferred"].ID:
		return _CapabilitiesRegistry.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CapabilitiesRegistryCapabilityConfigured) Topic() common.Hash {
	return common.HexToHash("0x04f0a9bcf3f3a3b42a4d7ca081119755f82ebe43e0d30c8f7292c4fe0dc4a2ae")
}

func (CapabilitiesRegistryCapabilityDeprecated) Topic() common.Hash {
	return common.HexToHash("0xdcea1b78b6ddc31592a94607d537543fcaafda6cc52d6d5cc7bbfca1422baf21")
}

func (CapabilitiesRegistryConfigSet) Topic() common.Hash {
	return common.HexToHash("0xf264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c1581703651")
}

func (CapabilitiesRegistryNodeAdded) Topic() common.Hash {
	return common.HexToHash("0x74becb12a5e8fd0e98077d02dfba8f647c9670c9df177e42c2418cf17a636f05")
}

func (CapabilitiesRegistryNodeOperatorAdded) Topic() common.Hash {
	return common.HexToHash("0x78e94ca80be2c30abc061b99e7eb8583b1254781734b1e3ce339abb57da2fe8e")
}

func (CapabilitiesRegistryNodeOperatorRemoved) Topic() common.Hash {
	return common.HexToHash("0xa59268ca81d40429e65ccea5385b59cf2d3fc6519371dee92f8eb1dae5107a7a")
}

func (CapabilitiesRegistryNodeOperatorUpdated) Topic() common.Hash {
	return common.HexToHash("0x86f41145bde5dd7f523305452e4aad3685508c181432ec733d5f345009358a28")
}

func (CapabilitiesRegistryNodeRemoved) Topic() common.Hash {
	return common.HexToHash("0x5254e609a97bab37b7cc79fe128f85c097bd6015c6e1624ae0ba392eb9753205")
}

func (CapabilitiesRegistryNodeUpdated) Topic() common.Hash {
	return common.HexToHash("0x4b5b465e22eea0c3d40c30e936643245b80d19b2dcf75788c0699fe8d8db645b")
}

func (CapabilitiesRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (CapabilitiesRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_CapabilitiesRegistry *CapabilitiesRegistry) Address() common.Address {
	return _CapabilitiesRegistry.address
}

type CapabilitiesRegistryInterface interface {
	GetCapabilities(opts *bind.CallOpts) ([]CapabilitiesRegistryCapabilityInfo, error)

	GetCapability(opts *bind.CallOpts, hashedId [32]byte) (CapabilitiesRegistryCapabilityInfo, error)

	GetCapabilityConfigs(opts *bind.CallOpts, donId uint32, capabilityId [32]byte) ([]byte, []byte, error)

	GetDON(opts *bind.CallOpts, donId uint32) (CapabilitiesRegistryDONInfo, error)

	GetDONs(opts *bind.CallOpts) ([]CapabilitiesRegistryDONInfo, error)

	GetHashedCapabilityId(opts *bind.CallOpts, labelledName string, version string) ([32]byte, error)

	GetNode(opts *bind.CallOpts, p2pId [32]byte) (CapabilitiesRegistryNodeInfo, error)

	GetNodeOperator(opts *bind.CallOpts, nodeOperatorId uint32) (CapabilitiesRegistryNodeOperator, error)

	GetNodeOperators(opts *bind.CallOpts) ([]CapabilitiesRegistryNodeOperator, error)

	GetNodes(opts *bind.CallOpts) ([]CapabilitiesRegistryNodeInfo, error)

	IsCapabilityDeprecated(opts *bind.CallOpts, hashedCapabilityId [32]byte) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddCapabilities(opts *bind.TransactOpts, capabilities []CapabilitiesRegistryCapability) (*types.Transaction, error)

	AddDON(opts *bind.TransactOpts, nodes [][32]byte, capabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error)

	AddNodeOperators(opts *bind.TransactOpts, nodeOperators []CapabilitiesRegistryNodeOperator) (*types.Transaction, error)

	AddNodes(opts *bind.TransactOpts, nodes []CapabilitiesRegistryNodeParams) (*types.Transaction, error)

	DeprecateCapabilities(opts *bind.TransactOpts, hashedCapabilityIds [][32]byte) (*types.Transaction, error)

	RemoveDONs(opts *bind.TransactOpts, donIds []uint32) (*types.Transaction, error)

	RemoveNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []uint32) (*types.Transaction, error)

	RemoveNodes(opts *bind.TransactOpts, removedNodeP2PIds [][32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateDON(opts *bind.TransactOpts, donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error)

	UpdateNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []uint32, nodeOperators []CapabilitiesRegistryNodeOperator) (*types.Transaction, error)

	UpdateNodes(opts *bind.TransactOpts, nodes []CapabilitiesRegistryNodeParams) (*types.Transaction, error)

	FilterCapabilityConfigured(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilitiesRegistryCapabilityConfiguredIterator, error)

	WatchCapabilityConfigured(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryCapabilityConfigured, hashedCapabilityId [][32]byte) (event.Subscription, error)

	ParseCapabilityConfigured(log types.Log) (*CapabilitiesRegistryCapabilityConfigured, error)

	FilterCapabilityDeprecated(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilitiesRegistryCapabilityDeprecatedIterator, error)

	WatchCapabilityDeprecated(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryCapabilityDeprecated, hashedCapabilityId [][32]byte) (event.Subscription, error)

	ParseCapabilityDeprecated(log types.Log) (*CapabilitiesRegistryCapabilityDeprecated, error)

	FilterConfigSet(opts *bind.FilterOpts) (*CapabilitiesRegistryConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*CapabilitiesRegistryConfigSet, error)

	FilterNodeAdded(opts *bind.FilterOpts, nodeOperatorId []uint32) (*CapabilitiesRegistryNodeAddedIterator, error)

	WatchNodeAdded(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryNodeAdded, nodeOperatorId []uint32) (event.Subscription, error)

	ParseNodeAdded(log types.Log) (*CapabilitiesRegistryNodeAdded, error)

	FilterNodeOperatorAdded(opts *bind.FilterOpts, nodeOperatorId []uint32, admin []common.Address) (*CapabilitiesRegistryNodeOperatorAddedIterator, error)

	WatchNodeOperatorAdded(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryNodeOperatorAdded, nodeOperatorId []uint32, admin []common.Address) (event.Subscription, error)

	ParseNodeOperatorAdded(log types.Log) (*CapabilitiesRegistryNodeOperatorAdded, error)

	FilterNodeOperatorRemoved(opts *bind.FilterOpts, nodeOperatorId []uint32) (*CapabilitiesRegistryNodeOperatorRemovedIterator, error)

	WatchNodeOperatorRemoved(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryNodeOperatorRemoved, nodeOperatorId []uint32) (event.Subscription, error)

	ParseNodeOperatorRemoved(log types.Log) (*CapabilitiesRegistryNodeOperatorRemoved, error)

	FilterNodeOperatorUpdated(opts *bind.FilterOpts, nodeOperatorId []uint32, admin []common.Address) (*CapabilitiesRegistryNodeOperatorUpdatedIterator, error)

	WatchNodeOperatorUpdated(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryNodeOperatorUpdated, nodeOperatorId []uint32, admin []common.Address) (event.Subscription, error)

	ParseNodeOperatorUpdated(log types.Log) (*CapabilitiesRegistryNodeOperatorUpdated, error)

	FilterNodeRemoved(opts *bind.FilterOpts) (*CapabilitiesRegistryNodeRemovedIterator, error)

	WatchNodeRemoved(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryNodeRemoved) (event.Subscription, error)

	ParseNodeRemoved(log types.Log) (*CapabilitiesRegistryNodeRemoved, error)

	FilterNodeUpdated(opts *bind.FilterOpts, nodeOperatorId []uint32) (*CapabilitiesRegistryNodeUpdatedIterator, error)

	WatchNodeUpdated(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryNodeUpdated, nodeOperatorId []uint32) (event.Subscription, error)

	ParseNodeUpdated(log types.Log) (*CapabilitiesRegistryNodeUpdated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CapabilitiesRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*CapabilitiesRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CapabilitiesRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*CapabilitiesRegistryOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
