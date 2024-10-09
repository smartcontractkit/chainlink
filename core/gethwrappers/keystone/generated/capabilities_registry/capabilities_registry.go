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
	EncryptionPublicKey [32]byte
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
	EncryptionPublicKey [32]byte
	HashedCapabilityIds [][32]byte
}

var CapabilitiesRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityIsDeprecated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"CapabilityRequiredByDON\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"DONDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONCapability\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONNode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedConfigurationContract\",\"type\":\"address\"}],\"name\":\"InvalidCapabilityConfigurationContractInterface\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"nodeCount\",\"type\":\"uint256\"}],\"name\":\"InvalidFaultTolerance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"InvalidNodeCapabilities\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"encryptionPublicKey\",\"type\":\"bytes32\"}],\"name\":\"InvalidNodeEncryptionPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeOperatorAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"InvalidNodeP2PId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lengthOne\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lengthTwo\",\"type\":\"uint256\"}],\"name\":\"LengthMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodeAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodeDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"NodeDoesNotSupportCapability\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"NodeOperatorDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodePartOfCapabilitiesDON\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodePartOfWorkflowDON\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDeprecated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"NodeOperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"NodeRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilitiesRegistry.Capability[]\",\"name\":\"capabilities\",\"type\":\"tuple[]\"}],\"name\":\"addCapabilities\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"name\":\"addDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilitiesRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"addNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"encryptionPublicKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilitiesRegistry.NodeParams[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"addNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"deprecateCapabilities\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCapabilities\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"hashedId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isDeprecated\",\"type\":\"bool\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedId\",\"type\":\"bytes32\"}],\"name\":\"getCapability\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"hashedId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isDeprecated\",\"type\":\"bool\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"getCapabilityConfigs\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"getDON\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilitiesRegistry.DONInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilitiesRegistry.DONInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"}],\"name\":\"getHashedCapabilityId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"getNode\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"workflowDONId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"encryptionPublicKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"capabilitiesDONIds\",\"type\":\"uint256[]\"}],\"internalType\":\"structCapabilitiesRegistry.NodeInfo\",\"name\":\"nodeInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"getNodeOperator\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilitiesRegistry.NodeOperator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeOperators\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilitiesRegistry.NodeOperator[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodes\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"workflowDONId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"encryptionPublicKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"capabilitiesDONIds\",\"type\":\"uint256[]\"}],\"internalType\":\"structCapabilitiesRegistry.NodeInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"isCapabilityDeprecated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"donIds\",\"type\":\"uint32[]\"}],\"name\":\"removeDONs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint32[]\"}],\"name\":\"removeNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"removedNodeP2PIds\",\"type\":\"bytes32[]\"}],\"name\":\"removeNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"name\":\"updateDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint32[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilitiesRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"updateNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"encryptionPublicKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilitiesRegistry.NodeParams[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"updateNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600e80546001600160401b0319166401000000011790553480156200002857600080fd5b503380600081620000805760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000b357620000b381620000bc565b50505062000167565b336001600160a01b03821603620001165760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000077565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61510580620001776000396000f3fe608060405234801561001057600080fd5b50600436106101ae5760003560e01c806350c946fe116100ee57806386fa424611610097578063d8bc7b6811610071578063d8bc7b68146103f6578063ddbe4f8214610409578063e29581aa1461041e578063f2fde38b1461043357600080fd5b806386fa42461461039b5780638da5cb5b146103ae5780639cb7c5f4146103d657600080fd5b8063715f5295116100c8578063715f52951461036d57806379ba50971461038057806384f5ed8a1461038857600080fd5b806350c946fe146103255780635d83d9671461034557806366acaa331461035857600080fd5b8063235374051161015b5780632c01a1e8116101355780632c01a1e8146102cb578063358039f4146102de578063398f3773146102f15780633f2a13c91461030457600080fd5b80632353740514610285578063275459f2146102a55780632a852933146102b857600080fd5b80631d05394c1161018c5780631d05394c1461023b578063214502431461025057806322bdbcbc1461026557600080fd5b80630fe5800a146101b357806312570011146101d9578063181f5a77146101fc575b600080fd5b6101c66101c1366004613f39565b610446565b6040519081526020015b60405180910390f35b6101ec6101e7366004613f9d565b61047a565b60405190151581526020016101d0565b604080518082018252601a81527f4361706162696c6974696573526567697374727920312e302e30000000000000602082015290516101d09190614024565b61024e61024936600461407c565b610487565b005b610258610694565b6040516101d091906141fe565b610278610273366004614299565b6107f1565b6040516101d091906142f1565b610298610293366004614299565b6108de565b6040516101d09190614304565b61024e6102b336600461407c565b610922565b61024e6102c6366004614338565b6109f9565b61024e6102d936600461407c565b610ad9565b61024e6102ec36600461407c565b610d7c565b61024e6102ff36600461407c565b6114a4565b6103176103123660046143da565b611663565b6040516101d0929190614404565b610338610333366004613f9d565b61184f565b6040516101d091906144a1565b61024e61035336600461407c565b61193c565b610360611a31565b6040516101d091906144b4565b61024e61037b36600461407c565b611c14565b61024e611cc6565b61024e61039636600461407c565b611dc3565b61024e6103a9366004614529565b61229b565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101d0565b6103e96103e4366004613f9d565b6125db565b6040516101d09190614678565b61024e61040436600461468b565b612816565b6104116128e0565b6040516101d09190614711565b6104266129d4565b6040516101d09190614786565b61024e61044136600461481f565b612ae5565b6000828260405160200161045b929190614404565b6040516020818303038152906040528051906020012090505b92915050565b6000610474600583612af9565b61048f612b14565b60005b8181101561068f5760008383838181106104ae576104ae61483a565b90506020020160208101906104c39190614299565b63ffffffff8181166000908152600d60209081526040808320805464010000000081049095168085526001820190935290832094955093909290916a010000000000000000000090910460ff16905b61051b83612b97565b8110156105bb57811561057157600c60006105368584612ba1565b8152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff1690556105b3565b6105b18663ffffffff16600c60006105928588612ba190919063ffffffff16565b8152602001908152602001600020600501612bad90919063ffffffff16565b505b600101610512565b508354640100000000900463ffffffff16600003610612576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff861660048201526024015b60405180910390fd5b63ffffffff85166000818152600d6020908152604080832080547fffffffffffffffffffffffffffffffffffffffffff0000000000000000000000169055519182527ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c1581703651910160405180910390a25050505050806001019050610492565b505050565b600e54606090640100000000900463ffffffff1660006106b5600183614898565b63ffffffff1667ffffffffffffffff8111156106d3576106d3613df6565b60405190808252806020026020018201604052801561075a57816020015b6040805160e081018252600080825260208083018290529282018190526060808301829052608083019190915260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816106f15790505b509050600060015b8363ffffffff168163ffffffff1610156107ce5763ffffffff8082166000908152600d602052604090205416156107c65761079c81612bb9565b8383815181106107ae576107ae61483a565b6020026020010181905250816107c3906148b5565b91505b600101610762565b506107da600184614898565b63ffffffff1681146107ea578082525b5092915050565b60408051808201909152600081526060602082015263ffffffff82166000908152600b60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff1683526001810180549192840191610855906148ed565b80601f0160208091040260200160405190810160405280929190818152602001828054610881906148ed565b80156108ce5780601f106108a3576101008083540402835291602001916108ce565b820191906000526020600020905b8154815290600101906020018083116108b157829003601f168201915b5050505050815250509050919050565b6040805160e0810182526000808252602082018190529181018290526060808201839052608082019290925260a0810182905260c081019190915261047482612bb9565b61092a612b14565b60005b63ffffffff811682111561068f57600083838363ffffffff168181106109555761095561483a565b905060200201602081019061096a9190614299565b63ffffffff81166000908152600b6020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001681559192506109b56001830182613d89565b505060405163ffffffff8216907fa59268ca81d40429e65ccea5385b59cf2d3fc6519371dee92f8eb1dae5107a7a90600090a2506109f281614940565b905061092d565b610a01612b14565b63ffffffff8088166000908152600d60205260408120805490926401000000009091041690819003610a67576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff8a166004820152602401610609565b610ace888888886040518060a001604052808f63ffffffff16815260200187610a8f90614940565b63ffffffff811682528b15156020830152895460ff6a01000000000000000000009091048116151560408401528b166060909201919091529650612e84565b505050505050505050565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610d76576000848483818110610b1357610b1361483a565b602090810292909201356000818152600c90935260409092206001810154929350919050610b70576040517fd82f6adb00000000000000000000000000000000000000000000000000000000815260048101839052602401610609565b6000610b7e82600501612b97565b1115610bd357610b916005820184612ba1565b6040517f60a6d89800000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015260248101839052604401610609565b805468010000000000000000900463ffffffff1615610c3b5780546040517f60b9df730000000000000000000000000000000000000000000000000000000081526801000000000000000090910463ffffffff16600482015260248101839052604401610609565b83158015610c755750805463ffffffff166000908152600b602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b15610cae576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610609565b6001810154610cbf90600790612bad565b506002810154610cd190600990612bad565b506000828152600c6020526040812080547fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001681556001810182905560028101829055600381018290559060058201818181610d2d8282613dc3565b5050505050507f5254e609a97bab37b7cc79fe128f85c097bd6015c6e1624ae0ba392eb975320582604051610d6491815260200190565b60405180910390a15050600101610af7565b50505050565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610d76576000848483818110610db657610db661483a565b9050602002810190610dc89190614963565b610dd1906149a1565b6040808201516000908152600c6020908152828220805463ffffffff168352600b82528383208451808601909552805473ffffffffffffffffffffffffffffffffffffffff1685526001810180549697509195939493909284019190610e36906148ed565b80601f0160208091040260200160405190810160405280929190818152602001828054610e62906148ed565b8015610eaf5780601f10610e8457610100808354040283529160200191610eaf565b820191906000526020600020905b815481529060010190602001808311610e9257829003601f168201915b505050919092525050506001830154909150610eff5782604001516040517fd82f6adb00000000000000000000000000000000000000000000000000000000815260040161060991815260200190565b84158015610f245750805173ffffffffffffffffffffffffffffffffffffffff163314155b15610f5d576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610609565b6020830151610f98576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600182015460208401518114611019576020840151610fb990600790612af9565b15610ff0576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208401516001840155611005600782612bad565b5060208401516110179060079061369b565b505b6080840151805160000361105b57806040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106099190614a7e565b8354600090859060049061107c90640100000000900463ffffffff16614940565b91906101000a81548163ffffffff021916908363ffffffff1602179055905060005b8251811015611161576110d48382815181106110bc576110bc61483a565b60200260200101516003612af990919063ffffffff16565b61110c57826040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106099190614a7e565b6111588382815181106111215761112161483a565b60200260200101518760040160008563ffffffff1663ffffffff16815260200190815260200160002061369b90919063ffffffff16565b5060010161109e565b50845468010000000000000000900463ffffffff1680156112c25763ffffffff8082166000908152600d6020908152604080832080546401000000009004909416835260019093018152828220600201805484518184028101840190955280855292939290918301828280156111f657602002820191906000526020600020905b8154815260200190600101908083116111e2575b5050505050905060005b81518110156112bf5761125582828151811061121e5761121e61483a565b60200260200101518960040160008763ffffffff1663ffffffff168152602001908152602001600020612af990919063ffffffff16565b6112b75781818151811061126b5761126b61483a565b6020026020010151836040517f03dcd86200000000000000000000000000000000000000000000000000000000815260040161060992919091825263ffffffff16602082015260400190565b600101611200565b50505b60006112d0876005016136a7565b905060005b81518163ffffffff161015611416576000828263ffffffff16815181106112fe576112fe61483a565b60209081029190910181015163ffffffff8082166000908152600d8452604080822080546401000000009004909316825260019092018452818120600201805483518187028101870190945280845293955090939192909183018282801561138557602002820191906000526020600020905b815481526020019060010190808311611371575b5050505050905060005b8151811015611402576113e48282815181106113ad576113ad61483a565b60200260200101518c60040160008a63ffffffff1663ffffffff168152602001908152602001600020612af990919063ffffffff16565b6113fa5781818151811061126b5761126b61483a565b60010161138f565b5050508061140f90614940565b90506112d5565b50875187547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff90911690811788556040808a015160028a018190556020808c01518351928352908201527f4b5b465e22eea0c3d40c30e936643245b80d19b2dcf75788c0699fe8d8db645b910160405180910390a25050505050505050806001019050610d9a565b6114ac612b14565b60005b8181101561068f5760008383838181106114cb576114cb61483a565b90506020028101906114dd9190614ac2565b6114e690614af6565b805190915073ffffffffffffffffffffffffffffffffffffffff16611537576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e54604080518082018252835173ffffffffffffffffffffffffffffffffffffffff908116825260208086015181840190815263ffffffff9095166000818152600b909252939020825181547fffffffffffffffffffffffff000000000000000000000000000000000000000016921691909117815592519192909160018201906115c39082614bb0565b5050600e80549091506000906115de9063ffffffff16614940565b91906101000a81548163ffffffff021916908363ffffffff160217905550816000015173ffffffffffffffffffffffffffffffffffffffff168163ffffffff167f78e94ca80be2c30abc061b99e7eb8583b1254781734b1e3ce339abb57da2fe8e84602001516040516116519190614024565b60405180910390a350506001016114af565b63ffffffff8083166000908152600d602090815260408083208054640100000000900490941680845260019094018252808320858452600301909152812080546060938493909290916116b5906148ed565b80601f01602080910402602001604051908101604052809291908181526020018280546116e1906148ed565b801561172e5780601f106117035761010080835404028352916020019161172e565b820191906000526020600020905b81548152906001019060200180831161171157829003601f168201915b5050506000888152600260208190526040909120015492935060609262010000900473ffffffffffffffffffffffffffffffffffffffff16159150611841905057600086815260026020819052604091829020015490517f8318ed5d00000000000000000000000000000000000000000000000000000000815263ffffffff891660048201526201000090910473ffffffffffffffffffffffffffffffffffffffff1690638318ed5d90602401600060405180830381865afa1580156117f8573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261183e9190810190614cca565b90505b9093509150505b9250929050565b604080516101008101825260008082526020820181905291810182905260608082018390526080820183905260a082019290925260c0810182905260e081019190915260408051610100810182526000848152600c6020908152838220805463ffffffff8082168652640100000000820481168487018190526801000000000000000090920416858701526001820154606086015260028201546080860152600382015460a0860152835260040190529190912060c0820190611911906136a7565b8152602001611934600c60008681526020019081526020016000206005016136a7565b905292915050565b611944612b14565b60005b8181101561068f5760008383838181106119635761196361483a565b905060200201359050611980816003612af990919063ffffffff16565b6119b9576040517fe181733f00000000000000000000000000000000000000000000000000000000815260048101829052602401610609565b6119c460058261369b565b6119fd576040517ff7d7a29400000000000000000000000000000000000000000000000000000000815260048101829052602401610609565b60405181907fdcea1b78b6ddc31592a94607d537543fcaafda6cc52d6d5cc7bbfca1422baf2190600090a250600101611947565b600e5460609063ffffffff166000611a4a600183614898565b63ffffffff1667ffffffffffffffff811115611a6857611a68613df6565b604051908082528060200260200182016040528015611aae57816020015b604080518082019091526000815260606020820152815260200190600190039081611a865790505b509050600060015b8363ffffffff168163ffffffff161015611bfe5763ffffffff81166000908152600b602052604090205473ffffffffffffffffffffffffffffffffffffffff1615611bf65763ffffffff81166000908152600b60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff1683526001810180549192840191611b4a906148ed565b80601f0160208091040260200160405190810160405280929190818152602001828054611b76906148ed565b8015611bc35780601f10611b9857610100808354040283529160200191611bc3565b820191906000526020600020905b815481529060010190602001808311611ba657829003601f168201915b505050505081525050838381518110611bde57611bde61483a565b602002602001018190525081611bf3906148b5565b91505b600101611ab6565b50600e546107da9060019063ffffffff16614898565b611c1c612b14565b60005b8181101561068f576000838383818110611c3b57611c3b61483a565b9050602002810190611c4d9190614963565b611c5690614d47565b90506000611c6c82600001518360200151610446565b9050611c7960038261369b565b611cb2576040517febf5255100000000000000000000000000000000000000000000000000000000815260048101829052602401610609565b611cbc81836136b4565b5050600101611c1f565b60015473ffffffffffffffffffffffffffffffffffffffff163314611d47576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610609565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610d76576000848483818110611dfd57611dfd61483a565b9050602002810190611e0f9190614963565b611e18906149a1565b805163ffffffff166000908152600b602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff168252600181018054959650939491939092840191611e6e906148ed565b80601f0160208091040260200160405190810160405280929190818152602001828054611e9a906148ed565b8015611ee75780601f10611ebc57610100808354040283529160200191611ee7565b820191906000526020600020905b815481529060010190602001808311611eca57829003601f168201915b50505091909252505081519192505073ffffffffffffffffffffffffffffffffffffffff16611f4d5781516040517fadd9ae1e00000000000000000000000000000000000000000000000000000000815263ffffffff9091166004820152602401610609565b83158015611f725750805173ffffffffffffffffffffffffffffffffffffffff163314155b15611fab576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610609565b6040808301516000908152600c602052206001810154156120005782604001516040517f5461848300000000000000000000000000000000000000000000000000000000815260040161060991815260200190565b60408301516120435782604001516040517f64e2ee9200000000000000000000000000000000000000000000000000000000815260040161060991815260200190565b602083015115806120605750602083015161206090600790612af9565b15612097576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b608083015180516000036120d957806040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106099190614a7e565b815482906004906120f790640100000000900463ffffffff16614940565b82546101009290920a63ffffffff818102199093169183160217909155825464010000000090041660005b82518110156121cd576121408382815181106110bc576110bc61483a565b61217857826040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106099190614a7e565b6121c483828151811061218d5761218d61483a565b60200260200101518560040160008563ffffffff1663ffffffff16815260200190815260200160002061369b90919063ffffffff16565b50600101612122565b5060608501516003840155845183547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9182161784556040860151600285015560208601516001850181905561222d916007919061369b16565b50604085015161223f9060099061369b565b50845160408087015160208089015183519283529082015263ffffffff909216917f74becb12a5e8fd0e98077d02dfba8f647c9670c9df177e42c2418cf17a636f05910160405180910390a25050505050806001019050611de1565b8281146122de576040517fab8b67c60000000000000000000000000000000000000000000000000000000081526004810184905260248101829052604401610609565b6000805473ffffffffffffffffffffffffffffffffffffffff16905b848110156125d35760008686838181106123165761231661483a565b905060200201602081019061232b9190614299565b63ffffffff81166000908152600b6020526040902080549192509073ffffffffffffffffffffffffffffffffffffffff1661239a576040517fadd9ae1e00000000000000000000000000000000000000000000000000000000815263ffffffff83166004820152602401610609565b60008686858181106123ae576123ae61483a565b90506020028101906123c09190614ac2565b6123c990614af6565b805190915073ffffffffffffffffffffffffffffffffffffffff1661241a576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b815473ffffffffffffffffffffffffffffffffffffffff16331480159061245757503373ffffffffffffffffffffffffffffffffffffffff861614155b15612490576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610609565b8051825473ffffffffffffffffffffffffffffffffffffffff908116911614158061250c57506020808201516040516124c99201614024565b60405160208183030381529060405280519060200120826001016040516020016124f39190614ded565b6040516020818303038152906040528051906020012014155b156125c557805182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602081015160018301906125669082614bb0565b50806000015173ffffffffffffffffffffffffffffffffffffffff168363ffffffff167f86f41145bde5dd7f523305452e4aad3685508c181432ec733d5f345009358a2883602001516040516125bc9190614024565b60405180910390a35b5050508060010190506122fa565b505050505050565b61261c6040805160e0810182526000808252606060208301819052928201839052909182019081526020016000815260006020820181905260409091015290565b6040805160e08101825283815260008481526002602090815292902080549192830191612648906148ed565b80601f0160208091040260200160405190810160405280929190818152602001828054612674906148ed565b80156126c15780601f10612696576101008083540402835291602001916126c1565b820191906000526020600020905b8154815290600101906020018083116126a457829003601f168201915b505050505081526020016002600085815260200190815260200160002060010180546126ec906148ed565b80601f0160208091040260200160405190810160405280929190818152602001828054612718906148ed565b80156127655780601f1061273a57610100808354040283529160200191612765565b820191906000526020600020905b81548152906001019060200180831161274857829003601f168201915b50505091835250506000848152600260208181526040909220015491019060ff16600381111561279757612797614595565b815260008481526002602081815260409092200154910190610100900460ff1660018111156127c8576127c8614595565b81526000848152600260208181526040928390209091015462010000900473ffffffffffffffffffffffffffffffffffffffff16908301520161280c600585612af9565b1515905292915050565b61281e612b14565b600e805460009164010000000090910463ffffffff1690600461284083614940565b82546101009290920a63ffffffff81810219909316918316021790915581166000818152600d602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001684179055815160a08101835292835260019083015286151590820152841515606082015260ff841660808201529091506128d6908990899089908990612e84565b5050505050505050565b606060006128ee60036136a7565b90506000815167ffffffffffffffff81111561290c5761290c613df6565b60405190808252806020026020018201604052801561297e57816020015b61296b6040805160e0810182526000808252606060208301819052928201839052909182019081526020016000815260006020820181905260409091015290565b81526020019060019003908161292a5790505b50905060005b82518110156107ea576129af8382815181106129a2576129a261483a565b60200260200101516125db565b8282815181106129c1576129c161483a565b6020908102919091010152600101612984565b606060006129e260096136a7565b90506000815167ffffffffffffffff811115612a0057612a00613df6565b604051908082528060200260200182016040528015612a8f57816020015b60408051610100810182526000808252602080830182905292820181905260608083018290526080830182905260a083019190915260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181612a1e5790505b50905060005b82518110156107ea57612ac0838281518110612ab357612ab361483a565b602002602001015161184f565b828281518110612ad257612ad261483a565b6020908102919091010152600101612a95565b612aed612b14565b612af68161389c565b50565b600081815260018301602052604081205415155b9392505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612b95576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610609565b565b6000610474825490565b6000612b0d8383613991565b6000612b0d83836139bb565b6040805160e0810182526000808252602080830182905282840182905260608084018390526080840183905260a0840181905260c084015263ffffffff8581168352600d8252848320805464010000000090049091168084526001909101825284832060028101805487518186028101860190985280885295969295919493909190830182828015612c6a57602002820191906000526020600020905b815481526020019060010190808311612c56575b505050505090506000815167ffffffffffffffff811115612c8d57612c8d613df6565b604051908082528060200260200182016040528015612cd357816020015b604080518082019091526000815260606020820152815260200190600190039081612cab5790505b50905060005b8151811015612deb576040518060400160405280848381518110612cff57612cff61483a565b60200260200101518152602001856003016000868581518110612d2457612d2461483a565b602002602001015181526020019081526020016000208054612d45906148ed565b80601f0160208091040260200160405190810160405280929190818152602001828054612d71906148ed565b8015612dbe5780601f10612d9357610100808354040283529160200191612dbe565b820191906000526020600020905b815481529060010190602001808311612da157829003601f168201915b5050505050815250828281518110612dd857612dd861483a565b6020908102919091010152600101612cd9565b506040805160e08101825263ffffffff8089166000818152600d6020818152868320548086168752948b168187015260ff680100000000000000008604811697870197909752690100000000000000000085048716151560608701529290915290526a010000000000000000000090049091161515608082015260a08101612e72856136a7565b81526020019190915295945050505050565b805163ffffffff9081166000908152600d602090815260408083208286015190941683526001909301905220608082015160ff161580612ed6575060808201518590612ed1906001614e9b565b60ff16115b15612f1f5760808201516040517f25b4d61800000000000000000000000000000000000000000000000000000000815260ff909116600482015260248101869052604401610609565b6001826020015163ffffffff16111561300757815163ffffffff166000908152600d602090815260408220908401516001918201918391612f609190614898565b63ffffffff1663ffffffff168152602001908152602001600020905060005b612f8882612b97565b81101561300457612fb7846000015163ffffffff16600c60006105928587600001612ba190919063ffffffff16565b50600c6000612fc68484612ba1565b8152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff169055600101612f7f565b50505b60005b85811015613241576130378787838181106130275761302761483a565b859260209091020135905061369b565b61309857825187878381811061304f5761304f61483a565b6040517f636e405700000000000000000000000000000000000000000000000000000000815263ffffffff90941660048501526020029190910135602483015250604401610609565b8260600151156131ef57825163ffffffff16600c60008989858181106130c0576130c061483a565b602090810292909201358352508101919091526040016000205468010000000000000000900463ffffffff161480159061313a5750600c600088888481811061310b5761310b61483a565b602090810292909201358352508101919091526040016000205468010000000000000000900463ffffffff1615155b1561319c5782518787838181106131535761315361483a565b6040517f60b9df7300000000000000000000000000000000000000000000000000000000815263ffffffff90941660048501526020029190910135602483015250604401610609565b8251600c60008989858181106131b4576131b461483a565b90506020020135815260200190815260200160002060000160086101000a81548163ffffffff021916908363ffffffff160217905550613239565b82516132379063ffffffff16600c60008a8a868181106132115761321161483a565b90506020020135815260200190815260200160002060050161369b90919063ffffffff16565b505b60010161300a565b5060005b8381101561364f57368585838181106132605761326061483a565b90506020028101906132729190614ac2565b905061328060038235612af9565b6132b9576040517fe181733f00000000000000000000000000000000000000000000000000000000815281356004820152602401610609565b6132c560058235612af9565b156132ff576040517ff7d7a29400000000000000000000000000000000000000000000000000000000815281356004820152602401610609565b803560009081526003840160205260408120805461331c906148ed565b905011156133685783516040517f3927d08000000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015281356024820152604401610609565b60005b878110156134725761340f8235600c60008c8c8681811061338e5761338e61483a565b9050602002013581526020019081526020016000206004016000600c60008e8e888181106133be576133be61483a565b90506020020135815260200190815260200160002060000160049054906101000a900463ffffffff1663ffffffff1663ffffffff168152602001908152602001600020612af990919063ffffffff16565b61346a578888828181106134255761342561483a565b6040517fa7e792500000000000000000000000000000000000000000000000000000000081526020909102929092013560048301525082356024820152604401610609565b60010161336b565b506002830180546001810182556000918252602091829020833591015561349b90820182614eb4565b823560009081526003860160205260409020916134b9919083614f19565b50604080850151855163ffffffff9081166000908152600d602090815284822080549415156901000000000000000000027fffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffff90951694909417909355606088015188518316825284822080549115156a0100000000000000000000027fffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffff9092169190911790556080880151885183168252848220805460ff9290921668010000000000000000027fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff909216919091179055828801805189518416835294909120805494909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9094169390931790558551915161364692918435908c908c9061360c90880188614eb4565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250613aae92505050565b50600101613245565b50815160208084015160405163ffffffff91821681529216917ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c1581703651910160405180910390a2505050505050565b6000612b0d8383613b8f565b60606000612b0d83613bde565b608081015173ffffffffffffffffffffffffffffffffffffffff16156137565761370281608001517f78bea72100000000000000000000000000000000000000000000000000000000613c3a565b6137565760808101516040517fabb5e3fd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610609565b6000828152600260205260409020815182919081906137759082614bb0565b506020820151600182019061378a9082614bb0565b5060408201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660018360038111156137cc576137cc614595565b021790555060608201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600181111561381357613813614595565b0217905550608091909101516002909101805473ffffffffffffffffffffffffffffffffffffffff90921662010000027fffffffffffffffffffff0000000000000000000000000000000000000000ffff90921691909117905560405182907f04f0a9bcf3f3a3b42a4d7ca081119755f82ebe43e0d30c8f7292c4fe0dc4a2ae90600090a25050565b3373ffffffffffffffffffffffffffffffffffffffff82160361391b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610609565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008260000182815481106139a8576139a861483a565b9060005260206000200154905092915050565b60008181526001830160205260408120548015613aa45760006139df600183615034565b85549091506000906139f390600190615034565b9050818114613a58576000866000018281548110613a1357613a1361483a565b9060005260206000200154905080876000018481548110613a3657613a3661483a565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613a6957613a69615047565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610474565b6000915050610474565b6000848152600260208190526040909120015462010000900473ffffffffffffffffffffffffffffffffffffffff16156125d357600084815260026020819052604091829020015490517ffba64a7c0000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063fba64a7c90613b55908690869086908b908d90600401615076565b600060405180830381600087803b158015613b6f57600080fd5b505af1158015613b83573d6000803e3d6000fd5b50505050505050505050565b6000818152600183016020526040812054613bd657508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610474565b506000610474565b606081600001805480602002602001604051908101604052809291908181526020018280548015613c2e57602002820191906000526020600020905b815481526020019060010190808311613c1a575b50505050509050919050565b6000613c4583613c56565b8015612b0d5750612b0d8383613cba565b6000613c82827f01ffc9a700000000000000000000000000000000000000000000000000000000613cba565b80156104745750613cb3827fffffffff00000000000000000000000000000000000000000000000000000000613cba565b1592915050565b604080517fffffffff000000000000000000000000000000000000000000000000000000008316602480830191909152825180830390910181526044909101909152602080820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01ffc9a700000000000000000000000000000000000000000000000000000000178152825160009392849283928392918391908a617530fa92503d91506000519050828015613d72575060208210155b8015613d7e5750600081115b979650505050505050565b508054613d95906148ed565b6000825580601f10613da5575050565b601f016020900490600052602060002090810190612af69190613ddd565b5080546000825590600052602060002090810190612af691905b5b80821115613df25760008155600101613dde565b5090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715613e4857613e48613df6565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613e9557613e95613df6565b604052919050565b600067ffffffffffffffff821115613eb757613eb7613df6565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112613ef457600080fd5b8135613f07613f0282613e9d565b613e4e565b818152846020838601011115613f1c57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215613f4c57600080fd5b823567ffffffffffffffff80821115613f6457600080fd5b613f7086838701613ee3565b93506020850135915080821115613f8657600080fd5b50613f9385828601613ee3565b9150509250929050565b600060208284031215613faf57600080fd5b5035919050565b60005b83811015613fd1578181015183820152602001613fb9565b50506000910152565b60008151808452613ff2816020860160208601613fb6565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000612b0d6020830184613fda565b60008083601f84011261404957600080fd5b50813567ffffffffffffffff81111561406157600080fd5b6020830191508360208260051b850101111561184857600080fd5b6000806020838503121561408f57600080fd5b823567ffffffffffffffff8111156140a657600080fd5b6140b285828601614037565b90969095509350505050565b60008151808452602080850194506020840160005b838110156140ef578151875295820195908201906001016140d3565b509495945050505050565b600082825180855260208086019550808260051b84010181860160005b84811015614177578583037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018952815180518452840151604085850181905261416381860183613fda565b9a86019a9450505090830190600101614117565b5090979650505050505050565b600063ffffffff8083511684528060208401511660208501525060ff604083015116604084015260608201511515606084015260808201511515608084015260a082015160e060a08501526141dc60e08501826140be565b905060c083015184820360c08601526141f582826140fa565b95945050505050565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015614273577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452614261858351614184565b94509285019290850190600101614227565b5092979650505050505050565b803563ffffffff8116811461429457600080fd5b919050565b6000602082840312156142ab57600080fd5b612b0d82614280565b73ffffffffffffffffffffffffffffffffffffffff815116825260006020820151604060208501526142e96040850182613fda565b949350505050565b602081526000612b0d60208301846142b4565b602081526000612b0d6020830184614184565b8035801515811461429457600080fd5b803560ff8116811461429457600080fd5b600080600080600080600060a0888a03121561435357600080fd5b61435c88614280565b9650602088013567ffffffffffffffff8082111561437957600080fd5b6143858b838c01614037565b909850965060408a013591508082111561439e57600080fd5b506143ab8a828b01614037565b90955093506143be905060608901614317565b91506143cc60808901614327565b905092959891949750929550565b600080604083850312156143ed57600080fd5b6143f683614280565b946020939093013593505050565b6040815260006144176040830185613fda565b82810360208401526141f58185613fda565b600061010063ffffffff80845116855280602085015116602086015280604085015116604086015250606083015160608501526080830151608085015260a083015160a085015260c08301518160c0860152614487828601826140be565b91505060e083015184820360e08601526141f582826140be565b602081526000612b0d6020830184614429565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015614273577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526145178583516142b4565b945092850192908501906001016144dd565b6000806000806040858703121561453f57600080fd5b843567ffffffffffffffff8082111561455757600080fd5b61456388838901614037565b9096509450602087013591508082111561457c57600080fd5b5061458987828801614037565b95989497509550505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b805182526000602082015160e060208501526145e360e0850182613fda565b9050604083015184820360408601526145fc8282613fda565b91505060608301516004811061461457614614614595565b606085015260808301516002811061462e5761462e614595565b8060808601525060a083015161465c60a086018273ffffffffffffffffffffffffffffffffffffffff169052565b5060c083015161467060c086018215159052565b509392505050565b602081526000612b0d60208301846145c4565b600080600080600080600060a0888a0312156146a657600080fd5b873567ffffffffffffffff808211156146be57600080fd5b6146ca8b838c01614037565b909950975060208a01359150808211156146e357600080fd5b506146f08a828b01614037565b9096509450614703905060408901614317565b92506143be60608901614317565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015614273577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526147748583516145c4565b9450928501929085019060010161473a565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015614273577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526147e9858351614429565b945092850192908501906001016147af565b803573ffffffffffffffffffffffffffffffffffffffff8116811461429457600080fd5b60006020828403121561483157600080fd5b612b0d826147fb565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b63ffffffff8281168282160390808211156107ea576107ea614869565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036148e6576148e6614869565b5060010190565b600181811c9082168061490157607f821691505b60208210810361493a577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600063ffffffff80831681810361495957614959614869565b6001019392505050565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6183360301811261499757600080fd5b9190910192915050565b600060a082360312156149b357600080fd5b6149bb613e25565b6149c483614280565b8152602080840135818301526040840135604083015260608401356060830152608084013567ffffffffffffffff808211156149ff57600080fd5b9085019036601f830112614a1257600080fd5b813581811115614a2457614a24613df6565b8060051b9150614a35848301613e4e565b8181529183018401918481019036841115614a4f57600080fd5b938501935b83851015614a6d57843582529385019390850190614a54565b608087015250939695505050505050565b6020808252825182820181905260009190848201906040850190845b81811015614ab657835183529284019291840191600101614a9a565b50909695505050505050565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc183360301811261499757600080fd5b600060408236031215614b0857600080fd5b6040516040810167ffffffffffffffff8282108183111715614b2c57614b2c613df6565b81604052614b39856147fb565b83526020850135915080821115614b4f57600080fd5b50614b5c36828601613ee3565b60208301525092915050565b601f82111561068f576000816000526020600020601f850160051c81016020861015614b915750805b601f850160051c820191505b818110156125d357828155600101614b9d565b815167ffffffffffffffff811115614bca57614bca613df6565b614bde81614bd884546148ed565b84614b68565b602080601f831160018114614c315760008415614bfb5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556125d3565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015614c7e57888601518255948401946001909101908401614c5f565b5085821015614cba57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600060208284031215614cdc57600080fd5b815167ffffffffffffffff811115614cf357600080fd5b8201601f81018413614d0457600080fd5b8051614d12613f0282613e9d565b818152856020838501011115614d2757600080fd5b6141f5826020830160208601613fb6565b80356002811061429457600080fd5b600060a08236031215614d5957600080fd5b614d61613e25565b823567ffffffffffffffff80821115614d7957600080fd5b614d8536838701613ee3565b83526020850135915080821115614d9b57600080fd5b50614da836828601613ee3565b602083015250604083013560048110614dc057600080fd5b6040820152614dd160608401614d38565b6060820152614de2608084016147fb565b608082015292915050565b6000602080835260008454614e01816148ed565b8060208701526040600180841660008114614e235760018114614e5d57614e8d565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00851660408a0152604084151560051b8a01019550614e8d565b89600052602060002060005b85811015614e845781548b8201860152908301908801614e69565b8a016040019650505b509398975050505050505050565b60ff818116838216019081111561047457610474614869565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112614ee957600080fd5b83018035915067ffffffffffffffff821115614f0457600080fd5b60200191503681900382131561184857600080fd5b67ffffffffffffffff831115614f3157614f31613df6565b614f4583614f3f83546148ed565b83614b68565b6000601f841160018114614f975760008515614f615750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b17835561502d565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015614fe65786850135825560209485019460019092019101614fc6565b5086821015615021577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b8181038181111561047457610474614869565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6080815284608082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8611156150af57600080fd5b8560051b808860a0850137820182810360a090810160208501526150d590820187613fda565b91505063ffffffff8085166040840152808416606084015250969550505050505056fea164736f6c6343000818000a",
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

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactor) UpdateDON(opts *bind.TransactOpts, donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration, isPublic bool, f uint8) (*types.Transaction, error) {
	return _CapabilitiesRegistry.contract.Transact(opts, "updateDON", donId, nodes, capabilityConfigurations, isPublic, f)
}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) UpdateDON(donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration, isPublic bool, f uint8) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.UpdateDON(&_CapabilitiesRegistry.TransactOpts, donId, nodes, capabilityConfigurations, isPublic, f)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryTransactorSession) UpdateDON(donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration, isPublic bool, f uint8) (*types.Transaction, error) {
	return _CapabilitiesRegistry.Contract.UpdateDON(&_CapabilitiesRegistry.TransactOpts, donId, nodes, capabilityConfigurations, isPublic, f)
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

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) FilterConfigSet(opts *bind.FilterOpts, donId []uint32) (*CapabilitiesRegistryConfigSetIterator, error) {

	var donIdRule []interface{}
	for _, donIdItem := range donId {
		donIdRule = append(donIdRule, donIdItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.FilterLogs(opts, "ConfigSet", donIdRule)
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryConfigSetIterator{contract: _CapabilitiesRegistry.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_CapabilitiesRegistry *CapabilitiesRegistryFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryConfigSet, donId []uint32) (event.Subscription, error) {

	var donIdRule []interface{}
	for _, donIdItem := range donId {
		donIdRule = append(donIdRule, donIdItem)
	}

	logs, sub, err := _CapabilitiesRegistry.contract.WatchLogs(opts, "ConfigSet", donIdRule)
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

	UpdateDON(opts *bind.TransactOpts, donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration, isPublic bool, f uint8) (*types.Transaction, error)

	UpdateNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []uint32, nodeOperators []CapabilitiesRegistryNodeOperator) (*types.Transaction, error)

	UpdateNodes(opts *bind.TransactOpts, nodes []CapabilitiesRegistryNodeParams) (*types.Transaction, error)

	FilterCapabilityConfigured(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilitiesRegistryCapabilityConfiguredIterator, error)

	WatchCapabilityConfigured(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryCapabilityConfigured, hashedCapabilityId [][32]byte) (event.Subscription, error)

	ParseCapabilityConfigured(log types.Log) (*CapabilitiesRegistryCapabilityConfigured, error)

	FilterCapabilityDeprecated(opts *bind.FilterOpts, hashedCapabilityId [][32]byte) (*CapabilitiesRegistryCapabilityDeprecatedIterator, error)

	WatchCapabilityDeprecated(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryCapabilityDeprecated, hashedCapabilityId [][32]byte) (event.Subscription, error)

	ParseCapabilityDeprecated(log types.Log) (*CapabilitiesRegistryCapabilityDeprecated, error)

	FilterConfigSet(opts *bind.FilterOpts, donId []uint32) (*CapabilitiesRegistryConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *CapabilitiesRegistryConfigSet, donId []uint32) (event.Subscription, error)

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
