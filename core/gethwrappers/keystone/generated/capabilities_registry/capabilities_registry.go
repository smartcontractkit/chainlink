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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityIsDeprecated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"CapabilityRequiredByDON\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"DONDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONCapability\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONNode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedConfigurationContract\",\"type\":\"address\"}],\"name\":\"InvalidCapabilityConfigurationContractInterface\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"nodeCount\",\"type\":\"uint256\"}],\"name\":\"InvalidFaultTolerance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"InvalidNodeCapabilities\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeOperatorAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"InvalidNodeP2PId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lengthOne\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lengthTwo\",\"type\":\"uint256\"}],\"name\":\"LengthMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodeAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodeDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"NodeDoesNotSupportCapability\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"NodeOperatorDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodePartOfCapabilitiesDON\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodePartOfWorkflowDON\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDeprecated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"NodeOperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"NodeRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilitiesRegistry.Capability[]\",\"name\":\"capabilities\",\"type\":\"tuple[]\"}],\"name\":\"addCapabilities\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"name\":\"addDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilitiesRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"addNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilitiesRegistry.NodeParams[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"addNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"deprecateCapabilities\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCapabilities\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilitiesRegistry.Capability[]\",\"name\":\"capabilities\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedId\",\"type\":\"bytes32\"}],\"name\":\"getCapability\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilitiesRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilitiesRegistry.Capability\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"getCapabilityConfigs\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"getDON\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilitiesRegistry.DONInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilitiesRegistry.DONInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"}],\"name\":\"getHashedCapabilityId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"getNode\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"workflowDONId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"capabilitiesDONIds\",\"type\":\"uint256[]\"}],\"internalType\":\"structCapabilitiesRegistry.NodeInfo\",\"name\":\"nodeInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"getNodeOperator\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilitiesRegistry.NodeOperator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeOperators\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilitiesRegistry.NodeOperator[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodes\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"workflowDONId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"capabilitiesDONIds\",\"type\":\"uint256[]\"}],\"internalType\":\"structCapabilitiesRegistry.NodeInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"isCapabilityDeprecated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"donIds\",\"type\":\"uint32[]\"}],\"name\":\"removeDONs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint32[]\"}],\"name\":\"removeNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"removedNodeP2PIds\",\"type\":\"bytes32[]\"}],\"name\":\"removeNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilitiesRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"name\":\"updateDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint32[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilitiesRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"updateNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilitiesRegistry.NodeParams[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"updateNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600e80546001600160401b0319166401000000011790553480156200002857600080fd5b503380600081620000805760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000b357620000b381620000bc565b50505062000167565b336001600160a01b03821603620001165760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000077565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61522880620001776000396000f3fe608060405234801561001057600080fd5b50600436106101ae5760003560e01c80635e65e309116100ee5780638da5cb5b11610097578063d8bc7b6811610071578063d8bc7b68146103f6578063ddbe4f8214610409578063e29581aa1461041f578063f2fde38b1461043457600080fd5b80638da5cb5b1461039b5780639cb7c5f4146103c3578063d59a79f6146103e357600080fd5b806373ac22b4116100c857806373ac22b41461036d57806379ba50971461038057806386fa42461461038857600080fd5b80635e65e3091461033257806366acaa3314610345578063715f52951461035a57600080fd5b8063235374051161015b578063398f377311610135578063398f3773146102cb5780633f2a13c9146102de57806350c946fe146102ff5780635d83d9671461031f57600080fd5b80632353740514610285578063275459f2146102a55780632c01a1e8146102b857600080fd5b80631d05394c1161018c5780631d05394c1461023b578063214502431461025057806322bdbcbc1461026557600080fd5b80630fe5800a146101b357806312570011146101d9578063181f5a77146101fc575b600080fd5b6101c66101c136600461405d565b610447565b6040519081526020015b60405180910390f35b6101ec6101e73660046140c1565b61047b565b60405190151581526020016101d0565b604080518082018252601a81527f4361706162696c6974696573526567697374727920312e302e30000000000000602082015290516101d09190614148565b61024e6102493660046141a0565b610488565b005b6102586106ab565b6040516101d09190614300565b610278610273366004614399565b610810565b6040516101d091906143f1565b610298610293366004614399565b6108fd565b6040516101d09190614404565b61024e6102b33660046141a0565b610941565b61024e6102c63660046141a0565b610a18565b61024e6102d93660046141a0565b610cbd565b6102f16102ec366004614417565b610e85565b6040516101d0929190614441565b61031261030d3660046140c1565b611071565b6040516101d091906144d2565b61024e61032d3660046141a0565b61114b565b61024e6103403660046141a0565b611299565b61034d611881565b6040516101d091906144e5565b61024e6103683660046141a0565b611a6c565b61024e61037b3660046141a0565b611b27565b61024e612003565b61024e610396366004614558565b612100565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101d0565b6103d66103d13660046140c1565b612446565b6040516101d09190614680565b61024e6103f13660046146b2565b61264b565b61024e610404366004614767565b612714565b6104116127de565b6040516101d092919061480c565b610427612b41565b6040516101d0919061489c565b61024e610442366004614933565b612c53565b6000828260405160200161045c929190614441565b6040516020818303038152906040528051906020012090505b92915050565b6000610475600583612c67565b610490612c82565b60005b818110156106a65760008383838181106104af576104af61494e565b90506020020160208101906104c49190614399565b63ffffffff8181166000908152600d60209081526040808320805464010000000081049095168085526001820190935290832094955093909290916a010000000000000000000090910460ff16905b61051c83612d05565b8110156105c457811561057257600c60006105378584612d0f565b8152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff1690556105b4565b6105b28663ffffffff16600c60006105938588612d0f90919063ffffffff16565b8152602001908152602001600020600401612d1b90919063ffffffff16565b505b6105bd816149ac565b9050610513565b508354640100000000900463ffffffff1660000361061b576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff861660048201526024015b60405180910390fd5b63ffffffff85166000818152600d6020908152604080832080547fffffffffffffffffffffffffffffffffffffffffff00000000000000000000001690558051938452908301919091527ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c1581703651910160405180910390a150505050508061069f906149ac565b9050610493565b505050565b600e54606090640100000000900463ffffffff1660006106cc6001836149e4565b63ffffffff1667ffffffffffffffff8111156106ea576106ea613ef7565b60405190808252806020026020018201604052801561077157816020015b6040805160e081018252600080825260208083018290529282018190526060808301829052608083019190915260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816107085790505b509050600060015b8363ffffffff168163ffffffff1610156107ed5763ffffffff8082166000908152600d602052604090205416156107dd576107b381612d27565b8383815181106107c5576107c561494e565b6020026020010181905250816107da906149ac565b91505b6107e681614a01565b9050610779565b506107f96001846149e4565b63ffffffff168114610809578082525b5092915050565b60408051808201909152600081526060602082015263ffffffff82166000908152600b60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff168352600181018054919284019161087490614a24565b80601f01602080910402602001604051908101604052809291908181526020018280546108a090614a24565b80156108ed5780601f106108c2576101008083540402835291602001916108ed565b820191906000526020600020905b8154815290600101906020018083116108d057829003601f168201915b5050505050815250509050919050565b6040805160e0810182526000808252602082018190529181018290526060808201839052608082019290925260a0810182905260c081019190915261047582612d27565b610949612c82565b60005b63ffffffff81168211156106a657600083838363ffffffff168181106109745761097461494e565b90506020020160208101906109899190614399565b63ffffffff81166000908152600b6020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001681559192506109d46001830182613e8a565b505060405163ffffffff8216907fa59268ca81d40429e65ccea5385b59cf2d3fc6519371dee92f8eb1dae5107a7a90600090a250610a1181614a01565b905061094c565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610cb7576000848483818110610a5257610a5261494e565b602090810292909201356000818152600c90935260409092206001810154929350919050610aaf576040517fd82f6adb00000000000000000000000000000000000000000000000000000000815260048101839052602401610612565b6000610abd82600401612d05565b1115610b1257610ad06004820184612d0f565b6040517f60a6d89800000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015260248101839052604401610612565b805468010000000000000000900463ffffffff1615610b7a5780546040517f60b9df730000000000000000000000000000000000000000000000000000000081526801000000000000000090910463ffffffff16600482015260248101839052604401610612565b83158015610bb45750805463ffffffff166000908152600b602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b15610bed576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610612565b6001810154610bfe90600790612d1b565b506002810154610c1090600990612d1b565b506000828152600c6020526040812080547fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016815560018101829055600281018290559060048201818181610c658282613ec4565b5050505050507f5254e609a97bab37b7cc79fe128f85c097bd6015c6e1624ae0ba392eb975320582604051610c9c91815260200190565b60405180910390a1505080610cb0906149ac565b9050610a36565b50505050565b610cc5612c82565b60005b818110156106a6576000838383818110610ce457610ce461494e565b9050602002810190610cf69190614a77565b610cff90614ab5565b805190915073ffffffffffffffffffffffffffffffffffffffff16610d50576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e54604080518082018252835173ffffffffffffffffffffffffffffffffffffffff908116825260208086015181840190815263ffffffff9095166000818152600b909252939020825181547fffffffffffffffffffffffff00000000000000000000000000000000000000001692169190911781559251919290916001820190610ddc9082614b6d565b5050600e8054909150600090610df79063ffffffff16614a01565b91906101000a81548163ffffffff021916908363ffffffff160217905550816000015173ffffffffffffffffffffffffffffffffffffffff168163ffffffff167f78e94ca80be2c30abc061b99e7eb8583b1254781734b1e3ce339abb57da2fe8e8460200151604051610e6a9190614148565b60405180910390a3505080610e7e906149ac565b9050610cc8565b63ffffffff8083166000908152600d60209081526040808320805464010000000090049094168084526001909401825280832085845260030190915281208054606093849390929091610ed790614a24565b80601f0160208091040260200160405190810160405280929190818152602001828054610f0390614a24565b8015610f505780601f10610f2557610100808354040283529160200191610f50565b820191906000526020600020905b815481529060010190602001808311610f3357829003601f168201915b5050506000888152600260208190526040909120015492935060609262010000900473ffffffffffffffffffffffffffffffffffffffff16159150611063905057600086815260026020819052604091829020015490517f8318ed5d00000000000000000000000000000000000000000000000000000000815263ffffffff891660048201526201000090910473ffffffffffffffffffffffffffffffffffffffff1690638318ed5d90602401600060405180830381865afa15801561101a573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526110609190810190614c87565b90505b9093509150505b9250929050565b6040805160e0810182526000808252602082018190529181018290526060808201839052608082019290925260a0810182905260c08101919091526040805160e0810182526000848152600c6020908152838220805463ffffffff8082168652640100000000820481168487018190526801000000000000000090920416858701526001820154606086015260028201546080860152835260030190529190912060a082019061112090612ffb565b8152602001611143600c6000868152602001908152602001600020600401612ffb565b905292915050565b611153612c82565b60005b818110156106a65760008383838181106111725761117261494e565b90506020020135905061118f816003612c6790919063ffffffff16565b6111c8576040517fe181733f00000000000000000000000000000000000000000000000000000000815260048101829052602401610612565b6111d3600582613008565b61120c576040517ff7d7a29400000000000000000000000000000000000000000000000000000000815260048101829052602401610612565b6000818152600260205260408120906112258282613e8a565b611233600183016000613e8a565b5060020180547fffffffffffffffffffff0000000000000000000000000000000000000000000016905560405181907fdcea1b78b6ddc31592a94607d537543fcaafda6cc52d6d5cc7bbfca1422baf2190600090a250611292816149ac565b9050611156565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610cb75760008484838181106112d3576112d361494e565b90506020028101906112e59190614cf5565b6112ee90614d29565b805163ffffffff166000908152600b602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff16825260018101805495965093949193909284019161134490614a24565b80601f016020809104026020016040519081016040528092919081815260200182805461137090614a24565b80156113bd5780601f10611392576101008083540402835291602001916113bd565b820191906000526020600020905b8154815290600101906020018083116113a057829003601f168201915b5050505050815250509050831580156113ed5750805173ffffffffffffffffffffffffffffffffffffffff163314155b15611426576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610612565b6040808301516000908152600c60205220600181015461147a5782604001516040517fd82f6adb00000000000000000000000000000000000000000000000000000000815260040161061291815260200190565b60208301516114b5576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001810154602084015181146115365760208401516114d690600790612c67565b1561150d576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208401516001830155611522600782612d1b565b50602084015161153490600790613008565b505b6060840151805160000361157857806040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106129190614dfc565b8254600090849060049061159990640100000000900463ffffffff16614a01565b91906101000a81548163ffffffff021916908363ffffffff1602179055905060005b8251811015611686576115f18382815181106115d9576115d961494e565b60200260200101516003612c6790919063ffffffff16565b61162957826040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106129190614dfc565b61167583828151811061163e5761163e61494e565b60200260200101518660030160008563ffffffff1663ffffffff16815260200190815260200160002061300890919063ffffffff16565b5061167f816149ac565b90506115bb565b50835468010000000000000000900463ffffffff1680156117ef5763ffffffff8082166000908152600d60209081526040808320805464010000000090049094168352600190930181528282206002018054845181840281018401909552808552929392909183018282801561171b57602002820191906000526020600020905b815481526020019060010190808311611707575b5050505050905060005b81518110156117ec5761177a8282815181106117435761174361494e565b60200260200101518860030160008763ffffffff1663ffffffff168152602001908152602001600020612c6790919063ffffffff16565b6117dc578181815181106117905761179061494e565b6020026020010151836040517f03dcd86200000000000000000000000000000000000000000000000000000000815260040161061292919091825263ffffffff16602082015260400190565b6117e5816149ac565b9050611725565b50505b865185547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9091169081178655604080890151600288018190556020808b01518351928352908201527f4b5b465e22eea0c3d40c30e936643245b80d19b2dcf75788c0699fe8d8db645b910160405180910390a2505050505050508061187a906149ac565b90506112b7565b600e5460609063ffffffff16600061189a6001836149e4565b63ffffffff1667ffffffffffffffff8111156118b8576118b8613ef7565b6040519080825280602002602001820160405280156118fe57816020015b6040805180820190915260008152606060208201528152602001906001900390816118d65790505b509050600060015b8363ffffffff168163ffffffff161015611a565763ffffffff81166000908152600b602052604090205473ffffffffffffffffffffffffffffffffffffffff1615611a465763ffffffff81166000908152600b60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff168352600181018054919284019161199a90614a24565b80601f01602080910402602001604051908101604052809291908181526020018280546119c690614a24565b8015611a135780601f106119e857610100808354040283529160200191611a13565b820191906000526020600020905b8154815290600101906020018083116119f657829003601f168201915b505050505081525050838381518110611a2e57611a2e61494e565b602002602001018190525081611a43906149ac565b91505b611a4f81614a01565b9050611906565b50600e546107f99060019063ffffffff166149e4565b611a74612c82565b60005b818110156106a6576000838383818110611a9357611a9361494e565b9050602002810190611aa59190614e0f565b611aae90614e52565b90506000611ac482600001518360200151610447565b9050611ad1600382613008565b611b0a576040517febf5255100000000000000000000000000000000000000000000000000000000815260048101829052602401610612565b611b148183613014565b505080611b20906149ac565b9050611a77565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610cb7576000848483818110611b6157611b6161494e565b9050602002810190611b739190614cf5565b611b7c90614d29565b805163ffffffff166000908152600b602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff168252600181018054959650939491939092840191611bd290614a24565b80601f0160208091040260200160405190810160405280929190818152602001828054611bfe90614a24565b8015611c4b5780601f10611c2057610100808354040283529160200191611c4b565b820191906000526020600020905b815481529060010190602001808311611c2e57829003601f168201915b50505091909252505081519192505073ffffffffffffffffffffffffffffffffffffffff16611cb15781516040517fadd9ae1e00000000000000000000000000000000000000000000000000000000815263ffffffff9091166004820152602401610612565b83158015611cd65750805173ffffffffffffffffffffffffffffffffffffffff163314155b15611d0f576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610612565b6040808301516000908152600c60205220600181015415611d645782604001516040517f5461848300000000000000000000000000000000000000000000000000000000815260040161061291815260200190565b6040830151611da75782604001516040517f64e2ee9200000000000000000000000000000000000000000000000000000000815260040161061291815260200190565b60208301511580611dc457506020830151611dc490600790612c67565b15611dfb576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608301518051600003611e3d57806040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106129190614dfc565b81548290600490611e5b90640100000000900463ffffffff16614a01565b82546101009290920a63ffffffff818102199093169183160217909155825464010000000090041660005b8251811015611f3957611ea48382815181106115d9576115d961494e565b611edc57826040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016106129190614dfc565b611f28838281518110611ef157611ef161494e565b60200260200101518560030160008563ffffffff1663ffffffff16815260200190815260200160002061300890919063ffffffff16565b50611f32816149ac565b9050611e86565b50845183547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff91821617845560408601516002850155602086015160018501819055611f8f916007919061300816565b506040850151611fa190600990613008565b50845160408087015160208089015183519283529082015263ffffffff909216917f74becb12a5e8fd0e98077d02dfba8f647c9670c9df177e42c2418cf17a636f05910160405180910390a2505050505080611ffc906149ac565b9050611b45565b60015473ffffffffffffffffffffffffffffffffffffffff163314612084576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610612565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b828114612143576040517fab8b67c60000000000000000000000000000000000000000000000000000000081526004810184905260248101829052604401610612565b6000805473ffffffffffffffffffffffffffffffffffffffff16905b8481101561243e57600086868381811061217b5761217b61494e565b90506020020160208101906121909190614399565b63ffffffff81166000908152600b6020526040902080549192509073ffffffffffffffffffffffffffffffffffffffff166121ff576040517fadd9ae1e00000000000000000000000000000000000000000000000000000000815263ffffffff83166004820152602401610612565b60008686858181106122135761221361494e565b90506020028101906122259190614a77565b61222e90614ab5565b805190915073ffffffffffffffffffffffffffffffffffffffff1661227f576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805173ffffffffffffffffffffffffffffffffffffffff1633148015906122bc57503373ffffffffffffffffffffffffffffffffffffffff861614155b156122f5576040517f9473075d000000000000000000000000000000000000000000000000000000008152336004820152602401610612565b8051825473ffffffffffffffffffffffffffffffffffffffff9081169116141580612371575060208082015160405161232e9201614148565b60405160208183030381529060405280519060200120826001016040516020016123589190614ef8565b6040516020818303038152906040528051906020012014155b1561242a57805182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602081015160018301906123cb9082614b6d565b50806000015173ffffffffffffffffffffffffffffffffffffffff168363ffffffff167f86f41145bde5dd7f523305452e4aad3685508c181432ec733d5f345009358a2883602001516040516124219190614148565b60405180910390a35b50505080612437906149ac565b905061215f565b505050505050565b6124766040805160a081018252606080825260208201529081016000815260200160008152600060209091015290565b60008281526002602052604090819020815160a0810190925280548290829061249e90614a24565b80601f01602080910402602001604051908101604052809291908181526020018280546124ca90614a24565b80156125175780601f106124ec57610100808354040283529160200191612517565b820191906000526020600020905b8154815290600101906020018083116124fa57829003601f168201915b5050505050815260200160018201805461253090614a24565b80601f016020809104026020016040519081016040528092919081815260200182805461255c90614a24565b80156125a95780601f1061257e576101008083540402835291602001916125a9565b820191906000526020600020905b81548152906001019060200180831161258c57829003601f168201915b5050509183525050600282015460209091019060ff1660038111156125d0576125d06145c4565b60038111156125e1576125e16145c4565b81526020016002820160019054906101000a900460ff166001811115612609576126096145c4565b600181111561261a5761261a6145c4565b81526002919091015462010000900473ffffffffffffffffffffffffffffffffffffffff1660209091015292915050565b612653612c82565b63ffffffff8089166000908152600d60205260408120546401000000009004909116908190036126b7576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff8a166004820152602401610612565b612709888888886040518060a001604052808f63ffffffff168152602001876126df90614a01565b97508763ffffffff1681526020018a1515815260200189151581526020018860ff168152506132a8565b505050505050505050565b61271c612c82565b600e805460009164010000000090910463ffffffff1690600461273e83614a01565b82546101009290920a63ffffffff81810219909316918316021790915581166000818152600d602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001684179055815160a08101835292835260019083015286151590820152841515606082015260ff841660808201529091506127d49089908990899089906132a8565b5050505050505050565b6060806127eb6003612ffb565b915060006127f96005612d05565b83516128059190614fa1565b90508067ffffffffffffffff81111561282057612820613ef7565b60405190808252806020026020018201604052801561288157816020015b61286e6040805160a081018252606080825260208201529081016000815260200160008152600060209091015290565b81526020019060019003908161283e5790505b50915060008167ffffffffffffffff81111561289f5761289f613ef7565b6040519080825280602002602001820160405280156128c8578160200160208202803683370190505b5090506000805b8551811015612b375760008682815181106128ec576128ec61494e565b6020026020010151905061290a816005612c6790919063ffffffff16565b612b265760008181526002602052604090819020815160a0810190925280548290829061293690614a24565b80601f016020809104026020016040519081016040528092919081815260200182805461296290614a24565b80156129af5780601f10612984576101008083540402835291602001916129af565b820191906000526020600020905b81548152906001019060200180831161299257829003601f168201915b505050505081526020016001820180546129c890614a24565b80601f01602080910402602001604051908101604052809291908181526020018280546129f490614a24565b8015612a415780601f10612a1657610100808354040283529160200191612a41565b820191906000526020600020905b815481529060010190602001808311612a2457829003601f168201915b5050509183525050600282015460209091019060ff166003811115612a6857612a686145c4565b6003811115612a7957612a796145c4565b81526020016002820160019054906101000a900460ff166001811115612aa157612aa16145c4565b6001811115612ab257612ab26145c4565b81526002919091015462010000900473ffffffffffffffffffffffffffffffffffffffff166020909101528651879085908110612af157612af161494e565b602002602001018190525080848481518110612b0f57612b0f61494e565b6020908102919091010152612b23836149ac565b92505b50612b30816149ac565b90506128cf565b5090949293505050565b60606000612b4f6009612ffb565b90506000815167ffffffffffffffff811115612b6d57612b6d613ef7565b604051908082528060200260200182016040528015612bf457816020015b6040805160e081018252600080825260208083018290529282018190526060808301829052608083019190915260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181612b8b5790505b50905060005b825181101561080957612c25838281518110612c1857612c1861494e565b6020026020010151611071565b828281518110612c3757612c3761494e565b602002602001018190525080612c4c906149ac565b9050612bfa565b612c5b612c82565b612c6481613aec565b50565b600081815260018301602052604081205415155b9392505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612d03576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610612565b565b6000610475825490565b6000612c7b8383613be1565b6000612c7b8383613c0b565b6040805160e0810182526000808252602080830182905282840182905260608084018390526080840183905260a0840181905260c084015263ffffffff8581168352600d8252848320805464010000000090049091168084526001909101825284832060028101805487518186028101860190985280885295969295919493909190830182828015612dd857602002820191906000526020600020905b815481526020019060010190808311612dc4575b505050505090506000815167ffffffffffffffff811115612dfb57612dfb613ef7565b604051908082528060200260200182016040528015612e4157816020015b604080518082019091526000815260606020820152815260200190600190039081612e195790505b50905060005b8151811015612f62576040518060400160405280848381518110612e6d57612e6d61494e565b60200260200101518152602001856003016000868581518110612e9257612e9261494e565b602002602001015181526020019081526020016000208054612eb390614a24565b80601f0160208091040260200160405190810160405280929190818152602001828054612edf90614a24565b8015612f2c5780601f10612f0157610100808354040283529160200191612f2c565b820191906000526020600020905b815481529060010190602001808311612f0f57829003601f168201915b5050505050815250828281518110612f4657612f4661494e565b602002602001018190525080612f5b906149ac565b9050612e47565b506040805160e08101825263ffffffff8089166000818152600d6020818152868320548086168752948b168187015260ff680100000000000000008604811697870197909752690100000000000000000085048716151560608701529290915290526a010000000000000000000090049091161515608082015260a08101612fe985612ffb565b81526020019190915295945050505050565b60606000612c7b83613cfe565b6000612c7b8383613d5a565b608081015173ffffffffffffffffffffffffffffffffffffffff161561316257608081015173ffffffffffffffffffffffffffffffffffffffff163b158061310d575060808101516040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527f78bea72100000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff909116906301ffc9a790602401602060405180830381865afa1580156130e7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061310b9190614fb4565b155b156131625760808101516040517fabb5e3fd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610612565b6000828152600260205260409020815182919081906131819082614b6d565b50602082015160018201906131969082614b6d565b5060408201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660018360038111156131d8576131d86145c4565b021790555060608201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600181111561321f5761321f6145c4565b0217905550608091909101516002909101805473ffffffffffffffffffffffffffffffffffffffff90921662010000027fffffffffffffffffffff0000000000000000000000000000000000000000ffff90921691909117905560405182907f04f0a9bcf3f3a3b42a4d7ca081119755f82ebe43e0d30c8f7292c4fe0dc4a2ae90600090a25050565b805163ffffffff9081166000908152600d602090815260408083208286015190941683526001909301905220608082015160ff1615806132fa5750608082015185906132f5906001614fd1565b60ff16115b156133435760808201516040517f25b4d61800000000000000000000000000000000000000000000000000000000815260ff909116600482015260248101869052604401610612565b6001826020015163ffffffff16111561343357815163ffffffff166000908152600d60209081526040822090840151600191820191839161338491906149e4565b63ffffffff1663ffffffff168152602001908152602001600020905060005b6133ac82612d05565b811015613430576133db846000015163ffffffff16600c60006105938587600001612d0f90919063ffffffff16565b50600c60006133ea8484612d0f565b8152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff169055613429816149ac565b90506133a3565b50505b60005b85811015613675576134638787838181106134535761345361494e565b8592602090910201359050613008565b6134c457825187878381811061347b5761347b61494e565b6040517f636e405700000000000000000000000000000000000000000000000000000000815263ffffffff90941660048501526020029190910135602483015250604401610612565b82606001511561361b57825163ffffffff16600c60008989858181106134ec576134ec61494e565b602090810292909201358352508101919091526040016000205468010000000000000000900463ffffffff16148015906135665750600c60008888848181106135375761353761494e565b602090810292909201358352508101919091526040016000205468010000000000000000900463ffffffff1615155b156135c857825187878381811061357f5761357f61494e565b6040517f60b9df7300000000000000000000000000000000000000000000000000000000815263ffffffff90941660048501526020029190910135602483015250604401610612565b8251600c60008989858181106135e0576135e061494e565b90506020020135815260200190815260200160002060000160086101000a81548163ffffffff021916908363ffffffff160217905550613665565b82516136639063ffffffff16600c60008a8a8681811061363d5761363d61494e565b90506020020135815260200190815260200160002060040161300890919063ffffffff16565b505b61366e816149ac565b9050613436565b5060005b8381101561396157368585838181106136945761369461494e565b90506020028101906136a69190614a77565b90506136b460038235612c67565b6136ed576040517fe181733f00000000000000000000000000000000000000000000000000000000815281356004820152602401610612565b6136f960058235612c67565b15613733576040517ff7d7a29400000000000000000000000000000000000000000000000000000000815281356004820152602401610612565b803560009081526003840160205260408120805461375090614a24565b9050111561379c5783516040517f3927d08000000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015281356024820152604401610612565b60005b878110156138ae576138438235600c60008c8c868181106137c2576137c261494e565b9050602002013581526020019081526020016000206003016000600c60008e8e888181106137f2576137f261494e565b90506020020135815260200190815260200160002060000160049054906101000a900463ffffffff1663ffffffff1663ffffffff168152602001908152602001600020612c6790919063ffffffff16565b61389e578888828181106138595761385961494e565b6040517fa7e792500000000000000000000000000000000000000000000000000000000081526020909102929092013560048301525082356024820152604401610612565b6138a7816149ac565b905061379f565b50600283018054600181018255600091825260209182902083359101556138d790820182614fea565b823560009081526003860160205260409020916138f591908361504f565b50835160208086015161395092918435908c908c9061391690880188614fea565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250613da992505050565b5061395a816149ac565b9050613679565b50604080830151835163ffffffff9081166000908152600d602090815284822080549415156901000000000000000000027fffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffff90951694909417909355606086015186518316825284822080549115156a0100000000000000000000027fffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffff9092169190911790556080860151865183168252848220805460ff9290921668010000000000000000027fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff909216919091179055918501805186518316845292849020805493909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9093169290921790558351905191517ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c158170365192613adc929163ffffffff92831681529116602082015260400190565b60405180910390a1505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603613b6b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610612565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110613bf857613bf861494e565b9060005260206000200154905092915050565b60008181526001830160205260408120548015613cf4576000613c2f600183614fa1565b8554909150600090613c4390600190614fa1565b9050818114613ca8576000866000018281548110613c6357613c6361494e565b9060005260206000200154905080876000018481548110613c8657613c8661494e565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613cb957613cb961516a565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610475565b6000915050610475565b606081600001805480602002602001604051908101604052809291908181526020018280548015613d4e57602002820191906000526020600020905b815481526020019060010190808311613d3a575b50505050509050919050565b6000818152600183016020526040812054613da157508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610475565b506000610475565b6000848152600260208190526040909120015462010000900473ffffffffffffffffffffffffffffffffffffffff161561243e57600084815260026020819052604091829020015490517ffba64a7c0000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063fba64a7c90613e50908690869086908b908d90600401615199565b600060405180830381600087803b158015613e6a57600080fd5b505af1158015613e7e573d6000803e3d6000fd5b50505050505050505050565b508054613e9690614a24565b6000825580601f10613ea6575050565b601f016020900490600052602060002090810190612c649190613ede565b5080546000825590600052602060002090810190612c6491905b5b80821115613ef35760008155600101613edf565b5090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516080810167ffffffffffffffff81118282101715613f4957613f49613ef7565b60405290565b60405160a0810167ffffffffffffffff81118282101715613f4957613f49613ef7565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613fb957613fb9613ef7565b604052919050565b600067ffffffffffffffff821115613fdb57613fdb613ef7565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261401857600080fd5b813561402b61402682613fc1565b613f72565b81815284602083860101111561404057600080fd5b816020850160208301376000918101602001919091529392505050565b6000806040838503121561407057600080fd5b823567ffffffffffffffff8082111561408857600080fd5b61409486838701614007565b935060208501359150808211156140aa57600080fd5b506140b785828601614007565b9150509250929050565b6000602082840312156140d357600080fd5b5035919050565b60005b838110156140f55781810151838201526020016140dd565b50506000910152565b600081518084526141168160208601602086016140da565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000612c7b60208301846140fe565b60008083601f84011261416d57600080fd5b50813567ffffffffffffffff81111561418557600080fd5b6020830191508360208260051b850101111561106a57600080fd5b600080602083850312156141b357600080fd5b823567ffffffffffffffff8111156141ca57600080fd5b6141d68582860161415b565b90969095509350505050565b600081518084526020808501945080840160005b83811015614212578151875295820195908201906001016141f6565b509495945050505050565b600081518084526020808501808196508360051b8101915082860160005b858110156142795782840389528151805185528501516040868601819052614265818701836140fe565b9a87019a955050509084019060010161423b565b5091979650505050505050565b600063ffffffff8083511684528060208401511660208501525060ff604083015116604084015260608201511515606084015260808201511515608084015260a082015160e060a08501526142de60e08501826141e2565b905060c083015184820360c08601526142f7828261421d565b95945050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015614373577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452614361858351614286565b94509285019290850190600101614327565b5092979650505050505050565b803563ffffffff8116811461439457600080fd5b919050565b6000602082840312156143ab57600080fd5b612c7b82614380565b73ffffffffffffffffffffffffffffffffffffffff815116825260006020820151604060208501526143e960408501826140fe565b949350505050565b602081526000612c7b60208301846143b4565b602081526000612c7b6020830184614286565b6000806040838503121561442a57600080fd5b61443383614380565b946020939093013593505050565b60408152600061445460408301856140fe565b82810360208401526142f781856140fe565b600063ffffffff80835116845280602084015116602085015280604084015116604085015250606082015160608401526080820151608084015260a082015160e060a08501526144b960e08501826141e2565b905060c083015184820360c08601526142f782826141e2565b602081526000612c7b6020830184614466565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015614373577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526145468583516143b4565b9450928501929085019060010161450c565b6000806000806040858703121561456e57600080fd5b843567ffffffffffffffff8082111561458657600080fd5b6145928883890161415b565b909650945060208701359150808211156145ab57600080fd5b506145b88782880161415b565b95989497509550505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6000815160a0845261460860a08501826140fe565b90506020830151848203602086015261462182826140fe565b915050604083015160048110614639576146396145c4565b6040850152606083015160028110614653576146536145c4565b606085015260809283015173ffffffffffffffffffffffffffffffffffffffff1692909301919091525090565b602081526000612c7b60208301846145f3565b8015158114612c6457600080fd5b803560ff8116811461439457600080fd5b60008060008060008060008060c0898b0312156146ce57600080fd5b6146d789614380565b9750602089013567ffffffffffffffff808211156146f457600080fd5b6147008c838d0161415b565b909950975060408b013591508082111561471957600080fd5b506147268b828c0161415b565b909650945050606089013561473a81614693565b9250608089013561474a81614693565b915061475860a08a016146a1565b90509295985092959890939650565b600080600080600080600060a0888a03121561478257600080fd5b873567ffffffffffffffff8082111561479a57600080fd5b6147a68b838c0161415b565b909950975060208a01359150808211156147bf57600080fd5b506147cc8a828b0161415b565b90965094505060408801356147e081614693565b925060608801356147f081614693565b91506147fe608089016146a1565b905092959891949750929550565b60408152600061481f60408301856141e2565b6020838203818501528185518084528284019150828160051b85010183880160005b8381101561488d577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe087840301855261487b8383516145f3565b94860194925090850190600101614841565b50909998505050505050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015614373577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526148fd858351614466565b945092850192908501906001016148c3565b803573ffffffffffffffffffffffffffffffffffffffff8116811461439457600080fd5b60006020828403121561494557600080fd5b612c7b8261490f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036149dd576149dd61497d565b5060010190565b63ffffffff8281168282160390808211156108095761080961497d565b600063ffffffff808316818103614a1a57614a1a61497d565b6001019392505050565b600181811c90821680614a3857607f821691505b602082108103614a71577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc1833603018112614aab57600080fd5b9190910192915050565b600060408236031215614ac757600080fd5b6040516040810167ffffffffffffffff8282108183111715614aeb57614aeb613ef7565b81604052614af88561490f565b83526020850135915080821115614b0e57600080fd5b50614b1b36828601614007565b60208301525092915050565b601f8211156106a657600081815260208120601f850160051c81016020861015614b4e5750805b601f850160051c820191505b8181101561243e57828155600101614b5a565b815167ffffffffffffffff811115614b8757614b87613ef7565b614b9b81614b958454614a24565b84614b27565b602080601f831160018114614bee5760008415614bb85750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855561243e565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015614c3b57888601518255948401946001909101908401614c1c565b5085821015614c7757878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600060208284031215614c9957600080fd5b815167ffffffffffffffff811115614cb057600080fd5b8201601f81018413614cc157600080fd5b8051614ccf61402682613fc1565b818152856020838501011115614ce457600080fd5b6142f78260208301602086016140da565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81833603018112614aab57600080fd5b600060808236031215614d3b57600080fd5b614d43613f26565b614d4c83614380565b81526020808401358183015260408401356040830152606084013567ffffffffffffffff80821115614d7d57600080fd5b9085019036601f830112614d9057600080fd5b813581811115614da257614da2613ef7565b8060051b9150614db3848301613f72565b8181529183018401918481019036841115614dcd57600080fd5b938501935b83851015614deb57843582529385019390850190614dd2565b606087015250939695505050505050565b602081526000612c7b60208301846141e2565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff61833603018112614aab57600080fd5b80356002811061439457600080fd5b600060a08236031215614e6457600080fd5b614e6c613f4f565b823567ffffffffffffffff80821115614e8457600080fd5b614e9036838701614007565b83526020850135915080821115614ea657600080fd5b50614eb336828601614007565b602083015250604083013560048110614ecb57600080fd5b6040820152614edc60608401614e43565b6060820152614eed6080840161490f565b608082015292915050565b6000602080835260008454614f0c81614a24565b80848701526040600180841660008114614f2d5760018114614f6557614f93565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838a01528284151560051b8a01019550614f93565b896000528660002060005b85811015614f8b5781548b8201860152908301908801614f70565b8a0184019650505b509398975050505050505050565b818103818111156104755761047561497d565b600060208284031215614fc657600080fd5b8151612c7b81614693565b60ff81811683821601908111156104755761047561497d565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261501f57600080fd5b83018035915067ffffffffffffffff82111561503a57600080fd5b60200191503681900382131561106a57600080fd5b67ffffffffffffffff83111561506757615067613ef7565b61507b836150758354614a24565b83614b27565b6000601f8411600181146150cd57600085156150975750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355615163565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b8281101561511c57868501358255602094850194600190920191016150fc565b5086821015615157577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6080815284608082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8611156151d257600080fd5b8560051b808860a0850137820182810360a090810160208501526151f8908201876140fe565b91505063ffffffff8085166040840152808416606084015250969550505050505056fea164736f6c6343000813000a",
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

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) GetCapabilities(opts *bind.CallOpts) (GetCapabilities,

	error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "getCapabilities")

	outstruct := new(GetCapabilities)
	if err != nil {
		return *outstruct, err
	}

	outstruct.HashedCapabilityIds = *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	outstruct.Capabilities = *abi.ConvertType(out[1], new([]CapabilitiesRegistryCapability)).(*[]CapabilitiesRegistryCapability)

	return *outstruct, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) GetCapabilities() (GetCapabilities,

	error) {
	return _CapabilitiesRegistry.Contract.GetCapabilities(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) GetCapabilities() (GetCapabilities,

	error) {
	return _CapabilitiesRegistry.Contract.GetCapabilities(&_CapabilitiesRegistry.CallOpts)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCaller) GetCapability(opts *bind.CallOpts, hashedId [32]byte) (CapabilitiesRegistryCapability, error) {
	var out []interface{}
	err := _CapabilitiesRegistry.contract.Call(opts, &out, "getCapability", hashedId)

	if err != nil {
		return *new(CapabilitiesRegistryCapability), err
	}

	out0 := *abi.ConvertType(out[0], new(CapabilitiesRegistryCapability)).(*CapabilitiesRegistryCapability)

	return out0, err

}

func (_CapabilitiesRegistry *CapabilitiesRegistrySession) GetCapability(hashedId [32]byte) (CapabilitiesRegistryCapability, error) {
	return _CapabilitiesRegistry.Contract.GetCapability(&_CapabilitiesRegistry.CallOpts, hashedId)
}

func (_CapabilitiesRegistry *CapabilitiesRegistryCallerSession) GetCapability(hashedId [32]byte) (CapabilitiesRegistryCapability, error) {
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

type GetCapabilities struct {
	HashedCapabilityIds [][32]byte
	Capabilities        []CapabilitiesRegistryCapability
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
	GetCapabilities(opts *bind.CallOpts) (GetCapabilities,

		error)

	GetCapability(opts *bind.CallOpts, hashedId [32]byte) (CapabilitiesRegistryCapability, error)

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
