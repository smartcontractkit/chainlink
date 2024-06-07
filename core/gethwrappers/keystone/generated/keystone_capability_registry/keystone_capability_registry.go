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
	F                        uint8
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityIsDeprecated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"DONDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONCapability\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONNode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedConfigurationContract\",\"type\":\"address\"}],\"name\":\"InvalidCapabilityConfigurationContractInterface\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"nodeCount\",\"type\":\"uint256\"}],\"name\":\"InvalidFaultTolerance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"InvalidNodeCapabilities\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeOperatorAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"InvalidNodeP2PId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lengthOne\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lengthTwo\",\"type\":\"uint256\"}],\"name\":\"LengthMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"NodeDoesNotSupportCapability\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"NodeOperatorDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodePartOfDON\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDeprecated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"NodeOperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"NodeRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability[]\",\"name\":\"capabilities\",\"type\":\"tuple[]\"}],\"name\":\"addCapabilities\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"name\":\"addDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"addNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"addNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"deprecateCapabilities\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCapabilities\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedId\",\"type\":\"bytes32\"}],\"name\":\"getCapability\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityType\",\"name\":\"capabilityType\",\"type\":\"uint8\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"getCapabilityConfigs\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"getDON\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilityRegistry.DONInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilityRegistry.DONInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"labelledName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"}],\"name\":\"getHashedCapabilityId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"getNode\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"getNodeOperator\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeOperators\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodes\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"internalType\":\"uint32[]\",\"name\":\"\",\"type\":\"uint32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"isCapabilityDeprecated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"donIds\",\"type\":\"uint32[]\"}],\"name\":\"removeDONs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint32[]\"}],\"name\":\"removeNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"removedNodeP2PIds\",\"type\":\"bytes32[]\"}],\"name\":\"removeNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"acceptsWorkflows\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"name\":\"updateDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint32[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"updateNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"updateNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600e80546001600160401b0319166401000000011790553480156200002857600080fd5b503380600081620000805760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000b357620000b381620000bc565b50505062000167565b336001600160a01b03821603620001165760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000077565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b614dc480620001776000396000f3fe608060405234801561001057600080fd5b50600436106101ae5760003560e01c80635e65e309116100ee5780638da5cb5b11610097578063d8bc7b6811610071578063d8bc7b68146103f7578063ddbe4f821461040a578063e29581aa1461041f578063f2fde38b1461043557600080fd5b80638da5cb5b1461039c5780639cb7c5f4146103c4578063d59a79f6146103e457600080fd5b806373ac22b4116100c857806373ac22b41461036e57806379ba50971461038157806386fa42461461038957600080fd5b80635e65e3091461033357806366acaa3314610346578063715f52951461035b57600080fd5b8063235374051161015b578063398f377311610135578063398f3773146102cb5780633f2a13c9146102de57806350c946fe146102ff5780635d83d9671461032057600080fd5b80632353740514610285578063275459f2146102a55780632c01a1e8146102b857600080fd5b80631d05394c1161018c5780631d05394c1461023b578063214502431461025057806322bdbcbc1461026557600080fd5b80630fe5800a146101b357806312570011146101d9578063181f5a77146101fc575b600080fd5b6101c66101c1366004613bc4565b610448565b6040519081526020015b60405180910390f35b6101ec6101e7366004613c28565b61047c565b60405190151581526020016101d0565b604080518082018252601881527f4361706162696c697479526567697374727920312e302e300000000000000000602082015290516101d09190613caf565b61024e610249366004613d07565b610489565b005b610258610645565b6040516101d09190613e67565b610278610273366004613f00565b6107aa565b6040516101d09190613f58565b610298610293366004613f00565b610897565b6040516101d09190613f6b565b61024e6102b3366004613d07565b6108db565b61024e6102c6366004613d07565b6109b2565b61024e6102d9366004613d07565b610bd5565b6102f16102ec366004613f7e565b610d9d565b6040516101d0929190613fa8565b61031261030d366004613c28565b610f89565b6040516101d0929190614006565b61024e61032e366004613d07565b61102e565b61024e610341366004613d07565b61117c565b61034e6115ee565b6040516101d0919061402e565b61024e610369366004613d07565b6117d9565b61024e61037c366004613d07565b611894565b61024e611d05565b61024e6103973660046140a1565b611e02565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101d0565b6103d76103d2366004613c28565b612148565b6040516101d091906141c9565b61024e6103f23660046141fb565b61234d565b61024e6104053660046142b0565b612416565b6104126124e0565b6040516101d09190614355565b6104276127de565b6040516101d09291906143c8565b61024e61044336600461449e565b612956565b6000828260405160200161045d929190613fa8565b6040516020818303038152906040528051906020012090505b92915050565b600061047660058361296a565b610491612985565b60005b818110156106405760008383838181106104b0576104b06144b9565b90506020020160208101906104c59190613f00565b63ffffffff8082166000908152600d60209081526040808320805464010000000090049094168084526001850190925282209394509192905b61050782612a08565b81101561055f5761054e8563ffffffff16600c600061052f8587612a1290919063ffffffff16565b8152602001908152602001600020600401612a1e90919063ffffffff16565b5061055881614517565b90506104fe565b508254640100000000900463ffffffff166000036105b6576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff851660048201526024015b60405180910390fd5b63ffffffff84166000818152600d6020908152604080832080547fffffffffffffffffffffffffffffffffffffffffff00000000000000000000001690558051938452908301919091527ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c1581703651910160405180910390a1505050508061063990614517565b9050610494565b505050565b600e54606090640100000000900463ffffffff16600061066660018361454f565b63ffffffff1667ffffffffffffffff81111561068457610684613a5e565b60405190808252806020026020018201604052801561070b57816020015b6040805160e081018252600080825260208083018290529282018190526060808301829052608083019190915260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816106a25790505b509050600060015b8363ffffffff168163ffffffff1610156107875763ffffffff8082166000908152600d602052604090205416156107775761074d81612a2a565b83838151811061075f5761075f6144b9565b60200260200101819052508161077490614517565b91505b6107808161456c565b9050610713565b5061079360018461454f565b63ffffffff1681146107a3578082525b5092915050565b60408051808201909152600081526060602082015263ffffffff82166000908152600b60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff168352600181018054919284019161080e9061458f565b80601f016020809104026020016040519081016040528092919081815260200182805461083a9061458f565b80156108875780601f1061085c57610100808354040283529160200191610887565b820191906000526020600020905b81548152906001019060200180831161086a57829003601f168201915b5050505050815250509050919050565b6040805160e0810182526000808252602082018190529181018290526060808201839052608082019290925260a0810182905260c081019190915261047682612a2a565b6108e3612985565b60005b63ffffffff811682111561064057600083838363ffffffff1681811061090e5761090e6144b9565b90506020020160208101906109239190613f00565b63ffffffff81166000908152600b6020526040812080547fffffffffffffffffffffffff000000000000000000000000000000000000000016815591925061096e60018301826139f1565b505060405163ffffffff8216907fa59268ca81d40429e65ccea5385b59cf2d3fc6519371dee92f8eb1dae5107a7a90600090a2506109ab8161456c565b90506108e6565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610bcf5760008484838181106109ec576109ec6144b9565b602090810292909201356000818152600c90935260409092206001810154929350919050610a49576040517f64e2ee92000000000000000000000000000000000000000000000000000000008152600481018390526024016105ad565b6000610a5782600401612a08565b1115610a92576040517f34a4a3f6000000000000000000000000000000000000000000000000000000008152600481018390526024016105ad565b83158015610acc5750805463ffffffff166000908152600b602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b15610b05576040517f9473075d0000000000000000000000000000000000000000000000000000000081523360048201526024016105ad565b6001810154610b1690600790612a1e565b506002810154610b2890600990612a1e565b506000828152600c6020526040812080547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016815560018101829055600281018290559060048201818181610b7d8282613a2b565b5050505050507f5254e609a97bab37b7cc79fe128f85c097bd6015c6e1624ae0ba392eb975320582604051610bb491815260200190565b60405180910390a1505080610bc890614517565b90506109d0565b50505050565b610bdd612985565b60005b81811015610640576000838383818110610bfc57610bfc6144b9565b9050602002810190610c0e91906145e2565b610c1790614620565b805190915073ffffffffffffffffffffffffffffffffffffffff16610c68576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e54604080518082018252835173ffffffffffffffffffffffffffffffffffffffff908116825260208086015181840190815263ffffffff9095166000818152600b909252939020825181547fffffffffffffffffffffffff00000000000000000000000000000000000000001692169190911781559251919290916001820190610cf490826146d8565b5050600e8054909150600090610d0f9063ffffffff1661456c565b91906101000a81548163ffffffff021916908363ffffffff160217905550816000015173ffffffffffffffffffffffffffffffffffffffff168163ffffffff167f78e94ca80be2c30abc061b99e7eb8583b1254781734b1e3ce339abb57da2fe8e8460200151604051610d829190613caf565b60405180910390a3505080610d9690614517565b9050610be0565b63ffffffff8083166000908152600d60209081526040808320805464010000000090049094168084526001909401825280832085845260030190915281208054606093849390929091610def9061458f565b80601f0160208091040260200160405190810160405280929190818152602001828054610e1b9061458f565b8015610e685780601f10610e3d57610100808354040283529160200191610e68565b820191906000526020600020905b815481529060010190602001808311610e4b57829003601f168201915b5050506000888152600260208190526040909120015492935060609262010000900473ffffffffffffffffffffffffffffffffffffffff16159150610f7b905057600086815260026020819052604091829020015490517f8318ed5d00000000000000000000000000000000000000000000000000000000815263ffffffff891660048201526201000090910473ffffffffffffffffffffffffffffffffffffffff1690638318ed5d90602401600060405180830381865afa158015610f32573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052610f7891908101906147f2565b90505b9093509150505b9250929050565b6040805160808101825260008082526020820181905291810191909152606080820152604080516080810182526000848152600c6020908152838220805463ffffffff80821686526001830154848701526002830154868801526401000000009091041683526003019052918220606082019061100590612cfe565b90526000938452600c602052604090932054929364010000000090930463ffffffff1692915050565b611036612985565b60005b81811015610640576000838383818110611055576110556144b9565b90506020020135905061107281600361296a90919063ffffffff16565b6110ab576040517fe181733f000000000000000000000000000000000000000000000000000000008152600481018290526024016105ad565b6110b6600582612d0b565b6110ef576040517ff7d7a294000000000000000000000000000000000000000000000000000000008152600481018290526024016105ad565b60008181526002602052604081209061110882826139f1565b6111166001830160006139f1565b5060020180547fffffffffffffffffffff0000000000000000000000000000000000000000000016905560405181907fdcea1b78b6ddc31592a94607d537543fcaafda6cc52d6d5cc7bbfca1422baf2190600090a25061117581614517565b9050611039565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610bcf5760008484838181106111b6576111b66144b9565b90506020028101906111c89190614860565b6111d190614894565b805163ffffffff166000908152600b602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff1682526001810180549596509394919390928401916112279061458f565b80601f01602080910402602001604051908101604052809291908181526020018280546112539061458f565b80156112a05780601f10611275576101008083540402835291602001916112a0565b820191906000526020600020905b81548152906001019060200180831161128357829003601f168201915b5050505050815250509050831580156112d05750805173ffffffffffffffffffffffffffffffffffffffff163314155b15611309576040517f9473075d0000000000000000000000000000000000000000000000000000000081523360048201526024016105ad565b6040808301516000908152600c60205220600181015461135d5782604001516040517f64e2ee920000000000000000000000000000000000000000000000000000000081526004016105ad91815260200190565b6020830151158061138d5750826020015181600101541415801561138d5750602083015161138d9060079061296a565b156113c4576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6060830151805160000361140657806040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016105ad9190614967565b8154640100000000900463ffffffff168260046114228361456c565b82546101009290920a63ffffffff8181021990931691831602179091558354640100000000900416905060005b825181101561151a5761148583828151811061146d5761146d6144b9565b6020026020010151600361296a90919063ffffffff16565b6114bd57826040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016105ad9190614967565b6115098382815181106114d2576114d26144b9565b60200260200101518560030160008563ffffffff1663ffffffff168152602001908152602001600020612d0b90919063ffffffff16565b5061151381614517565b905061144f565b50845183547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9091161783556040850151600284015560018301546020860151811461158c57611571600782612a1e565b5060208601516001850181905561158a90600790612d0b565b505b85516040808801516020808a015183519283529082015263ffffffff909216917f4b5b465e22eea0c3d40c30e936643245b80d19b2dcf75788c0699fe8d8db645b910160405180910390a2505050505050806115e790614517565b905061119a565b600e5460609063ffffffff16600061160760018361454f565b63ffffffff1667ffffffffffffffff81111561162557611625613a5e565b60405190808252806020026020018201604052801561166b57816020015b6040805180820190915260008152606060208201528152602001906001900390816116435790505b509050600060015b8363ffffffff168163ffffffff1610156117c35763ffffffff81166000908152600b602052604090205473ffffffffffffffffffffffffffffffffffffffff16156117b35763ffffffff81166000908152600b60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff16835260018101805491928401916117079061458f565b80601f01602080910402602001604051908101604052809291908181526020018280546117339061458f565b80156117805780601f1061175557610100808354040283529160200191611780565b820191906000526020600020905b81548152906001019060200180831161176357829003601f168201915b50505050508152505083838151811061179b5761179b6144b9565b6020026020010181905250816117b090614517565b91505b6117bc8161456c565b9050611673565b50600e546107939060019063ffffffff1661454f565b6117e1612985565b60005b81811015610640576000838383818110611800576118006144b9565b905060200281019061181291906149ab565b61181b906149ee565b9050600061183182600001518360200151610448565b905061183e600382612d0b565b611877576040517febf52551000000000000000000000000000000000000000000000000000000008152600481018290526024016105ad565b6118818183612d17565b50508061188d90614517565b90506117e4565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610bcf5760008484838181106118ce576118ce6144b9565b90506020028101906118e09190614860565b6118e990614894565b805163ffffffff166000908152600b602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff16825260018101805495965093949193909284019161193f9061458f565b80601f016020809104026020016040519081016040528092919081815260200182805461196b9061458f565b80156119b85780601f1061198d576101008083540402835291602001916119b8565b820191906000526020600020905b81548152906001019060200180831161199b57829003601f168201915b50505091909252505081519192505073ffffffffffffffffffffffffffffffffffffffff16611a1e5781516040517fadd9ae1e00000000000000000000000000000000000000000000000000000000815263ffffffff90911660048201526024016105ad565b83158015611a435750805173ffffffffffffffffffffffffffffffffffffffff163314155b15611a7c576040517f9473075d0000000000000000000000000000000000000000000000000000000081523360048201526024016105ad565b6040808301516000908152600c602052206001810154151580611aa157506040830151155b15611ae05782604001516040517f64e2ee920000000000000000000000000000000000000000000000000000000081526004016105ad91815260200190565b60208301511580611afd57506020830151611afd9060079061296a565b15611b34576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608301518051600003611b7657806040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016105ad9190614967565b81548290600490611b9490640100000000900463ffffffff1661456c565b82546101009290920a63ffffffff818102199093169183160217909155825464010000000090041660005b8251811015611c3b57611bdd83828151811061146d5761146d6144b9565b611c1557826040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016105ad9190614967565b611c2a8382815181106114d2576114d26144b9565b50611c3481614517565b9050611bbf565b50845183547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff91821617845560408601516002850155602086015160018501819055611c919160079190612d0b16565b506040850151611ca390600990612d0b565b50845160408087015160208089015183519283529082015263ffffffff909216917f74becb12a5e8fd0e98077d02dfba8f647c9670c9df177e42c2418cf17a636f05910160405180910390a2505050505080611cfe90614517565b90506118b2565b60015473ffffffffffffffffffffffffffffffffffffffff163314611d86576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016105ad565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b828114611e45576040517fab8b67c600000000000000000000000000000000000000000000000000000000815260048101849052602481018290526044016105ad565b6000805473ffffffffffffffffffffffffffffffffffffffff16905b84811015612140576000868683818110611e7d57611e7d6144b9565b9050602002016020810190611e929190613f00565b63ffffffff81166000908152600b6020526040902080549192509073ffffffffffffffffffffffffffffffffffffffff16611f01576040517fadd9ae1e00000000000000000000000000000000000000000000000000000000815263ffffffff831660048201526024016105ad565b6000868685818110611f1557611f156144b9565b9050602002810190611f2791906145e2565b611f3090614620565b805190915073ffffffffffffffffffffffffffffffffffffffff16611f81576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805173ffffffffffffffffffffffffffffffffffffffff163314801590611fbe57503373ffffffffffffffffffffffffffffffffffffffff861614155b15611ff7576040517f9473075d0000000000000000000000000000000000000000000000000000000081523360048201526024016105ad565b8051825473ffffffffffffffffffffffffffffffffffffffff908116911614158061207357506020808201516040516120309201613caf565b604051602081830303815290604052805190602001208260010160405160200161205a9190614a94565b6040516020818303038152906040528051906020012014155b1561212c57805182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602081015160018301906120cd90826146d8565b50806000015173ffffffffffffffffffffffffffffffffffffffff168363ffffffff167f86f41145bde5dd7f523305452e4aad3685508c181432ec733d5f345009358a2883602001516040516121239190613caf565b60405180910390a35b5050508061213990614517565b9050611e61565b505050505050565b6121786040805160a081018252606080825260208201529081016000815260200160008152600060209091015290565b60008281526002602052604090819020815160a081019092528054829082906121a09061458f565b80601f01602080910402602001604051908101604052809291908181526020018280546121cc9061458f565b80156122195780601f106121ee57610100808354040283529160200191612219565b820191906000526020600020905b8154815290600101906020018083116121fc57829003601f168201915b505050505081526020016001820180546122329061458f565b80601f016020809104026020016040519081016040528092919081815260200182805461225e9061458f565b80156122ab5780601f10612280576101008083540402835291602001916122ab565b820191906000526020600020905b81548152906001019060200180831161228e57829003601f168201915b5050509183525050600282015460209091019060ff1660038111156122d2576122d261410d565b60038111156122e3576122e361410d565b81526020016002820160019054906101000a900460ff16600181111561230b5761230b61410d565b600181111561231c5761231c61410d565b81526002919091015462010000900473ffffffffffffffffffffffffffffffffffffffff1660209091015292915050565b612355612985565b63ffffffff8089166000908152600d60205260408120546401000000009004909116908190036123b9576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff8a1660048201526024016105ad565b61240b888888886040518060a001604052808f63ffffffff168152602001876123e19061456c565b97508763ffffffff1681526020018a1515815260200189151581526020018860ff16815250612fab565b505050505050505050565b61241e612985565b600e805460009164010000000090910463ffffffff169060046124408361456c565b82546101009290920a63ffffffff81810219909316918316021790915581166000818152600d602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001684179055815160a08101835292835260019083015286151590820152841515606082015260ff841660808201529091506124d6908990899089908990612fab565b5050505050505050565b606060006124ee6003612cfe565b905060006124fc6005612a08565b82516125089190614b3d565b67ffffffffffffffff81111561252057612520613a5e565b60405190808252806020026020018201604052801561258157816020015b61256e6040805160a081018252606080825260208201529081016000815260200160008152600060209091015290565b81526020019060019003908161253e5790505b5090506000805b83518110156127d55760008482815181106125a5576125a56144b9565b602002602001015190506125c381600561296a90919063ffffffff16565b6127c45760008181526002602052604090819020815160a081019092528054829082906125ef9061458f565b80601f016020809104026020016040519081016040528092919081815260200182805461261b9061458f565b80156126685780601f1061263d57610100808354040283529160200191612668565b820191906000526020600020905b81548152906001019060200180831161264b57829003601f168201915b505050505081526020016001820180546126819061458f565b80601f01602080910402602001604051908101604052809291908181526020018280546126ad9061458f565b80156126fa5780601f106126cf576101008083540402835291602001916126fa565b820191906000526020600020905b8154815290600101906020018083116126dd57829003601f168201915b5050509183525050600282015460209091019060ff1660038111156127215761272161410d565b60038111156127325761273261410d565b81526020016002820160019054906101000a900460ff16600181111561275a5761275a61410d565b600181111561276b5761276b61410d565b81526002919091015462010000900473ffffffffffffffffffffffffffffffffffffffff1660209091015284518590859081106127aa576127aa6144b9565b602002602001018190525082806127c090614517565b9350505b506127ce81614517565b9050612588565b50909392505050565b60608060006127ed6009612cfe565b90506000815167ffffffffffffffff81111561280b5761280b613a5e565b60405190808252806020026020018201604052801561287a57816020015b60408051608081018252600080825260208083018290529282015260608082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816128295790505b5090506000825167ffffffffffffffff81111561289957612899613a5e565b6040519080825280602002602001820160405280156128c2578160200160208202803683370190505b50905060005b835181101561294b576128f38482815181106128e6576128e66144b9565b6020026020010151610f89565b848381518110612905576129056144b9565b6020026020010184848151811061291e5761291e6144b9565b602002602001018263ffffffff1663ffffffff1681525082905250508061294490614517565b90506128c8565b509094909350915050565b61295e612985565b61296781613653565b50565b600081815260018301602052604081205415155b9392505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612a06576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016105ad565b565b6000610476825490565b600061297e8383613748565b600061297e8383613772565b6040805160e0810182526000808252602080830182905282840182905260608084018390526080840183905260a0840181905260c084015263ffffffff8581168352600d8252848320805464010000000090049091168084526001909101825284832060028101805487518186028101860190985280885295969295919493909190830182828015612adb57602002820191906000526020600020905b815481526020019060010190808311612ac7575b505050505090506000815167ffffffffffffffff811115612afe57612afe613a5e565b604051908082528060200260200182016040528015612b4457816020015b604080518082019091526000815260606020820152815260200190600190039081612b1c5790505b50905060005b8151811015612c65576040518060400160405280848381518110612b7057612b706144b9565b60200260200101518152602001856003016000868581518110612b9557612b956144b9565b602002602001015181526020019081526020016000208054612bb69061458f565b80601f0160208091040260200160405190810160405280929190818152602001828054612be29061458f565b8015612c2f5780601f10612c0457610100808354040283529160200191612c2f565b820191906000526020600020905b815481529060010190602001808311612c1257829003601f168201915b5050505050815250828281518110612c4957612c496144b9565b602002602001018190525080612c5e90614517565b9050612b4a565b506040805160e08101825263ffffffff8089166000818152600d6020818152868320548086168752948b168187015260ff680100000000000000008604811697870197909752690100000000000000000085048716151560608701529290915290526a010000000000000000000090049091161515608082015260a08101612cec85612cfe565b81526020019190915295945050505050565b6060600061297e83613865565b600061297e83836138c1565b608081015173ffffffffffffffffffffffffffffffffffffffff1615612e6557608081015173ffffffffffffffffffffffffffffffffffffffff163b1580612e10575060808101516040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527f78bea72100000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff909116906301ffc9a790602401602060405180830381865afa158015612dea573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e0e9190614b50565b155b15612e655760808101516040517fabb5e3fd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911660048201526024016105ad565b600082815260026020526040902081518291908190612e8490826146d8565b5060208201516001820190612e9990826146d8565b5060408201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001836003811115612edb57612edb61410d565b021790555060608201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16610100836001811115612f2257612f2261410d565b0217905550608091909101516002909101805473ffffffffffffffffffffffffffffffffffffffff90921662010000027fffffffffffffffffffff0000000000000000000000000000000000000000ffff90921691909117905560405182907f04f0a9bcf3f3a3b42a4d7ca081119755f82ebe43e0d30c8f7292c4fe0dc4a2ae90600090a25050565b805163ffffffff9081166000908152600d602090815260408083208286015190941683526001909301905220608082015160ff161580612ffd575060808201518590612ff8906001614b6d565b60ff16115b156130465760808201516040517f25b4d61800000000000000000000000000000000000000000000000000000000815260ff9091166004820152602481018690526044016105ad565b6001826020015163ffffffff1611156130f257815163ffffffff166000908152600d602090815260408220908401516001918201918391613087919061454f565b63ffffffff1663ffffffff168152602001908152602001600020905060005b6130af82612a08565b8110156130ef576130de846000015163ffffffff16600c600061052f8587600001612a1290919063ffffffff16565b506130e881614517565b90506130a6565b50505b60005b858110156131dc57613122878783818110613112576131126144b9565b8592602090910201359050612d0b565b61318357825187878381811061313a5761313a6144b9565b6040517f636e405700000000000000000000000000000000000000000000000000000000815263ffffffff909416600485015260200291909101356024830152506044016105ad565b82516131cb9063ffffffff16600c60008a8a868181106131a5576131a56144b9565b905060200201358152602001908152602001600020600401612d0b90919063ffffffff16565b506131d581614517565b90506130f5565b5060005b838110156134c857368585838181106131fb576131fb6144b9565b905060200281019061320d91906145e2565b905061321b6003823561296a565b613254576040517fe181733f000000000000000000000000000000000000000000000000000000008152813560048201526024016105ad565b6132606005823561296a565b1561329a576040517ff7d7a294000000000000000000000000000000000000000000000000000000008152813560048201526024016105ad565b80356000908152600384016020526040812080546132b79061458f565b905011156133035783516040517f3927d08000000000000000000000000000000000000000000000000000000000815263ffffffff9091166004820152813560248201526044016105ad565b60005b87811015613415576133aa8235600c60008c8c86818110613329576133296144b9565b9050602002013581526020019081526020016000206003016000600c60008e8e88818110613359576133596144b9565b90506020020135815260200190815260200160002060000160049054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002061296a90919063ffffffff16565b613405578888828181106133c0576133c06144b9565b6040517fa7e7925000000000000000000000000000000000000000000000000000000000815260209091029290920135600483015250823560248201526044016105ad565b61340e81614517565b9050613306565b506002830180546001810182556000918252602091829020833591015561343e90820182614b86565b8235600090815260038601602052604090209161345c919083614beb565b5083516020808601516134b792918435908c908c9061347d90880188614b86565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061391092505050565b506134c181614517565b90506131e0565b50604080830151835163ffffffff9081166000908152600d602090815284822080549415156901000000000000000000027fffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffff90951694909417909355606086015186518316825284822080549115156a0100000000000000000000027fffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffff9092169190911790556080860151865183168252848220805460ff9290921668010000000000000000027fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff909216919091179055918501805186518316845292849020805493909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9093169290921790558351905191517ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c158170365192613643929163ffffffff92831681529116602082015260400190565b60405180910390a1505050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036136d2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016105ad565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600082600001828154811061375f5761375f6144b9565b9060005260206000200154905092915050565b6000818152600183016020526040812054801561385b576000613796600183614b3d565b85549091506000906137aa90600190614b3d565b905081811461380f5760008660000182815481106137ca576137ca6144b9565b90600052602060002001549050808760000184815481106137ed576137ed6144b9565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061382057613820614d06565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610476565b6000915050610476565b6060816000018054806020026020016040519081016040528092919081815260200182805480156138b557602002820191906000526020600020905b8154815260200190600101908083116138a1575b50505050509050919050565b600081815260018301602052604081205461390857508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610476565b506000610476565b6000848152600260208190526040909120015462010000900473ffffffffffffffffffffffffffffffffffffffff161561214057600084815260026020819052604091829020015490517ffba64a7c0000000000000000000000000000000000000000000000000000000081526201000090910473ffffffffffffffffffffffffffffffffffffffff169063fba64a7c906139b7908690869086908b908d90600401614d35565b600060405180830381600087803b1580156139d157600080fd5b505af11580156139e5573d6000803e3d6000fd5b50505050505050505050565b5080546139fd9061458f565b6000825580601f10613a0d575050565b601f0160209004906000526020600020908101906129679190613a45565b508054600082559060005260206000209081019061296791905b5b80821115613a5a5760008155600101613a46565b5090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516080810167ffffffffffffffff81118282101715613ab057613ab0613a5e565b60405290565b60405160a0810167ffffffffffffffff81118282101715613ab057613ab0613a5e565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613b2057613b20613a5e565b604052919050565b600067ffffffffffffffff821115613b4257613b42613a5e565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112613b7f57600080fd5b8135613b92613b8d82613b28565b613ad9565b818152846020838601011115613ba757600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215613bd757600080fd5b823567ffffffffffffffff80821115613bef57600080fd5b613bfb86838701613b6e565b93506020850135915080821115613c1157600080fd5b50613c1e85828601613b6e565b9150509250929050565b600060208284031215613c3a57600080fd5b5035919050565b60005b83811015613c5c578181015183820152602001613c44565b50506000910152565b60008151808452613c7d816020860160208601613c41565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061297e6020830184613c65565b60008083601f840112613cd457600080fd5b50813567ffffffffffffffff811115613cec57600080fd5b6020830191508360208260051b8501011115610f8257600080fd5b60008060208385031215613d1a57600080fd5b823567ffffffffffffffff811115613d3157600080fd5b613d3d85828601613cc2565b90969095509350505050565b600081518084526020808501945080840160005b83811015613d7957815187529582019590820190600101613d5d565b509495945050505050565b600081518084526020808501808196508360051b8101915082860160005b85811015613de05782840389528151805185528501516040868601819052613dcc81870183613c65565b9a87019a9550505090840190600101613da2565b5091979650505050505050565b600063ffffffff8083511684528060208401511660208501525060ff604083015116604084015260608201511515606084015260808201511515608084015260a082015160e060a0850152613e4560e0850182613d49565b905060c083015184820360c0860152613e5e8282613d84565b95945050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015613eda577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452613ec8858351613ded565b94509285019290850190600101613e8e565b5092979650505050505050565b803563ffffffff81168114613efb57600080fd5b919050565b600060208284031215613f1257600080fd5b61297e82613ee7565b73ffffffffffffffffffffffffffffffffffffffff81511682526000602082015160406020850152613f506040850182613c65565b949350505050565b60208152600061297e6020830184613f1b565b60208152600061297e6020830184613ded565b60008060408385031215613f9157600080fd5b613f9a83613ee7565b946020939093013593505050565b604081526000613fbb6040830185613c65565b8281036020840152613e5e8185613c65565b63ffffffff815116825260208101516020830152604081015160408301526000606082015160806060850152613f506080850182613d49565b6040815260006140196040830185613fcd565b905063ffffffff831660208301529392505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015613eda577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc088860301845261408f858351613f1b565b94509285019290850190600101614055565b600080600080604085870312156140b757600080fd5b843567ffffffffffffffff808211156140cf57600080fd5b6140db88838901613cc2565b909650945060208701359150808211156140f457600080fd5b5061410187828801613cc2565b95989497509550505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6000815160a0845261415160a0850182613c65565b90506020830151848203602086015261416a8282613c65565b9150506040830151600481106141825761418261410d565b604085015260608301516002811061419c5761419c61410d565b606085015260809283015173ffffffffffffffffffffffffffffffffffffffff1692909301919091525090565b60208152600061297e602083018461413c565b801515811461296757600080fd5b803560ff81168114613efb57600080fd5b60008060008060008060008060c0898b03121561421757600080fd5b61422089613ee7565b9750602089013567ffffffffffffffff8082111561423d57600080fd5b6142498c838d01613cc2565b909950975060408b013591508082111561426257600080fd5b5061426f8b828c01613cc2565b9096509450506060890135614283816141dc565b92506080890135614293816141dc565b91506142a160a08a016141ea565b90509295985092959890939650565b600080600080600080600060a0888a0312156142cb57600080fd5b873567ffffffffffffffff808211156142e357600080fd5b6142ef8b838c01613cc2565b909950975060208a013591508082111561430857600080fd5b506143158a828b01613cc2565b9096509450506040880135614329816141dc565b92506060880135614339816141dc565b9150614347608089016141ea565b905092959891949750929550565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015613eda577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526143b685835161413c565b9450928501929085019060010161437c565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b8381101561443d577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa088870301855261442b868351613fcd565b955093820193908201906001016143f1565b50508584038187015286518085528782019482019350915060005b82811015613de057845163ffffffff1684529381019392810192600101614458565b803573ffffffffffffffffffffffffffffffffffffffff81168114613efb57600080fd5b6000602082840312156144b057600080fd5b61297e8261447a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614548576145486144e8565b5060010190565b63ffffffff8281168282160390808211156107a3576107a36144e8565b600063ffffffff808316818103614585576145856144e8565b6001019392505050565b600181811c908216806145a357607f821691505b6020821081036145dc577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc183360301811261461657600080fd5b9190910192915050565b60006040823603121561463257600080fd5b6040516040810167ffffffffffffffff828210818311171561465657614656613a5e565b816040526146638561447a565b8352602085013591508082111561467957600080fd5b5061468636828601613b6e565b60208301525092915050565b601f82111561064057600081815260208120601f850160051c810160208610156146b95750805b601f850160051c820191505b81811015612140578281556001016146c5565b815167ffffffffffffffff8111156146f2576146f2613a5e565b61470681614700845461458f565b84614692565b602080601f83116001811461475957600084156147235750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555612140565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156147a657888601518255948401946001909101908401614787565b50858210156147e257878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b60006020828403121561480457600080fd5b815167ffffffffffffffff81111561481b57600080fd5b8201601f8101841361482c57600080fd5b805161483a613b8d82613b28565b81815285602083850101111561484f57600080fd5b613e5e826020830160208601613c41565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8183360301811261461657600080fd5b6000608082360312156148a657600080fd5b6148ae613a8d565b6148b783613ee7565b81526020808401358183015260408401356040830152606084013567ffffffffffffffff808211156148e857600080fd5b9085019036601f8301126148fb57600080fd5b81358181111561490d5761490d613a5e565b8060051b915061491e848301613ad9565b818152918301840191848101903684111561493857600080fd5b938501935b838510156149565784358252938501939085019061493d565b606087015250939695505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561499f57835183529284019291840191600101614983565b50909695505050505050565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6183360301811261461657600080fd5b803560028110613efb57600080fd5b600060a08236031215614a0057600080fd5b614a08613ab6565b823567ffffffffffffffff80821115614a2057600080fd5b614a2c36838701613b6e565b83526020850135915080821115614a4257600080fd5b50614a4f36828601613b6e565b602083015250604083013560048110614a6757600080fd5b6040820152614a78606084016149df565b6060820152614a896080840161447a565b608082015292915050565b6000602080835260008454614aa88161458f565b80848701526040600180841660008114614ac95760018114614b0157614b2f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838a01528284151560051b8a01019550614b2f565b896000528660002060005b85811015614b275781548b8201860152908301908801614b0c565b8a0184019650505b509398975050505050505050565b81810381811115610476576104766144e8565b600060208284031215614b6257600080fd5b815161297e816141dc565b60ff8181168382160190811115610476576104766144e8565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112614bbb57600080fd5b83018035915067ffffffffffffffff821115614bd657600080fd5b602001915036819003821315610f8257600080fd5b67ffffffffffffffff831115614c0357614c03613a5e565b614c1783614c11835461458f565b83614692565b6000601f841160018114614c695760008515614c335750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355614cff565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015614cb85786850135825560209485019460019092019101614c98565b5086821015614cf3577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6080815284608082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff861115614d6e57600080fd5b8560051b808860a0850137820182810360a09081016020850152614d9490820187613c65565b91505063ffffffff8085166040840152808416606084015250969550505050505056fea164736f6c6343000813000a",
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) AddCapabilities(opts *bind.TransactOpts, capabilities []CapabilityRegistryCapability) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "addCapabilities", capabilities)
}

func (_CapabilityRegistry *CapabilityRegistrySession) AddCapabilities(capabilities []CapabilityRegistryCapability) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddCapabilities(&_CapabilityRegistry.TransactOpts, capabilities)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) AddCapabilities(capabilities []CapabilityRegistryCapability) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddCapabilities(&_CapabilityRegistry.TransactOpts, capabilities)
}

func (_CapabilityRegistry *CapabilityRegistryTransactor) AddDON(opts *bind.TransactOpts, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "addDON", nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilityRegistry *CapabilityRegistrySession) AddDON(nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddDON(&_CapabilityRegistry.TransactOpts, nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) AddDON(nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error) {
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) UpdateDON(opts *bind.TransactOpts, donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "updateDON", donId, nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilityRegistry *CapabilityRegistrySession) UpdateDON(donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.UpdateDON(&_CapabilityRegistry.TransactOpts, donId, nodes, capabilityConfigurations, isPublic, acceptsWorkflows, f)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) UpdateDON(donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error) {
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
	GetCapabilities(opts *bind.CallOpts) ([]CapabilityRegistryCapability, error)

	GetCapability(opts *bind.CallOpts, hashedId [32]byte) (CapabilityRegistryCapability, error)

	GetCapabilityConfigs(opts *bind.CallOpts, donId uint32, capabilityId [32]byte) ([]byte, []byte, error)

	GetDON(opts *bind.CallOpts, donId uint32) (CapabilityRegistryDONInfo, error)

	GetDONs(opts *bind.CallOpts) ([]CapabilityRegistryDONInfo, error)

	GetHashedCapabilityId(opts *bind.CallOpts, labelledName string, version string) ([32]byte, error)

	GetNode(opts *bind.CallOpts, p2pId [32]byte) (CapabilityRegistryNodeInfo, uint32, error)

	GetNodeOperator(opts *bind.CallOpts, nodeOperatorId uint32) (CapabilityRegistryNodeOperator, error)

	GetNodeOperators(opts *bind.CallOpts) ([]CapabilityRegistryNodeOperator, error)

	GetNodes(opts *bind.CallOpts) ([]CapabilityRegistryNodeInfo, []uint32, error)

	IsCapabilityDeprecated(opts *bind.CallOpts, hashedCapabilityId [32]byte) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddCapabilities(opts *bind.TransactOpts, capabilities []CapabilityRegistryCapability) (*types.Transaction, error)

	AddDON(opts *bind.TransactOpts, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error)

	AddNodeOperators(opts *bind.TransactOpts, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error)

	AddNodes(opts *bind.TransactOpts, nodes []CapabilityRegistryNodeInfo) (*types.Transaction, error)

	DeprecateCapabilities(opts *bind.TransactOpts, hashedCapabilityIds [][32]byte) (*types.Transaction, error)

	RemoveDONs(opts *bind.TransactOpts, donIds []uint32) (*types.Transaction, error)

	RemoveNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []uint32) (*types.Transaction, error)

	RemoveNodes(opts *bind.TransactOpts, removedNodeP2PIds [][32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateDON(opts *bind.TransactOpts, donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool, acceptsWorkflows bool, f uint8) (*types.Transaction, error)

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
