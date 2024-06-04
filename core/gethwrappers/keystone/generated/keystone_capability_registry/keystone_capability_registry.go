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

type CapabilityRegistryCapabilityConfiguration struct {
	CapabilityId [32]byte
	Config       []byte
}

type CapabilityRegistryDONInfo struct {
	Id                       uint32
	ConfigCount              uint32
	IsPublic                 bool
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
	ABI: "[{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CapabilityAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityIsDeprecated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"DONDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONCapability\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateDONNode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedConfigurationContract\",\"type\":\"address\"}],\"name\":\"InvalidCapabilityConfigurationContractInterface\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"InvalidNodeCapabilities\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeOperatorAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"InvalidNodeP2PId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNodeSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lengthOne\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lengthTwo\",\"type\":\"uint256\"}],\"name\":\"LengthMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"NodeDoesNotSupportCapability\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"NodeOperatorDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeP2PId\",\"type\":\"bytes32\"}],\"name\":\"NodePartOfDON\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"CapabilityDeprecated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"NodeOperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NodeOperatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"NodeRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"}],\"name\":\"NodeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"labelledName\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability[]\",\"name\":\"capabilities\",\"type\":\"tuple[]\"}],\"name\":\"addCapabilities\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"}],\"name\":\"addDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"addNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"addNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"name\":\"deprecateCapabilities\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCapabilities\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"labelledName\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedId\",\"type\":\"bytes32\"}],\"name\":\"getCapability\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"labelledName\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"getDON\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilityRegistry.DONInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"}],\"name\":\"getDONCapabilityConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodeP2PIds\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"}],\"internalType\":\"structCapabilityRegistry.DONInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"labelledName\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"}],\"name\":\"getHashedCapabilityId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"}],\"name\":\"getNode\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"}],\"name\":\"getNodeOperator\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeOperators\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodes\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"internalType\":\"uint32[]\",\"name\":\"\",\"type\":\"uint32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedCapabilityId\",\"type\":\"bytes32\"}],\"name\":\"isCapabilityDeprecated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"donIds\",\"type\":\"uint32[]\"}],\"name\":\"removeDONs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint32[]\"}],\"name\":\"removeNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"removedNodeP2PIds\",\"type\":\"bytes32[]\"}],\"name\":\"removeNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"labelledName\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"enumCapabilityRegistry.CapabilityResponseType\",\"name\":\"responseType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"configurationContract\",\"type\":\"address\"}],\"internalType\":\"structCapabilityRegistry.Capability[]\",\"name\":\"capabilities\",\"type\":\"tuple[]\"}],\"name\":\"updateCapabilities\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"nodes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"capabilityId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCapabilityRegistry.CapabilityConfiguration[]\",\"name\":\"capabilityConfigurations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"}],\"name\":\"updateDON\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"nodeOperatorIds\",\"type\":\"uint32[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structCapabilityRegistry.NodeOperator[]\",\"name\":\"nodeOperators\",\"type\":\"tuple[]\"}],\"name\":\"updateNodeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nodeOperatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"signer\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashedCapabilityIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCapabilityRegistry.NodeInfo[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"name\":\"updateNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600e80546001600160401b0319166401000000011790553480156200002857600080fd5b503380600081620000805760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000b357620000b381620000bc565b50505062000167565b336001600160a01b03821603620001165760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000077565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61491980620001776000396000f3fe608060405234801561001057600080fd5b50600436106101b95760003560e01c80635e65e309116100f95780639cb7c5f411610097578063caef831d11610071578063caef831d14610431578063ddbe4f8214610444578063e29581aa14610459578063f2fde38b1461046f57600080fd5b80639cb7c5f4146103eb578063b06e07a71461040b578063c63239c51461041e57600080fd5b806373ac22b4116100d357806373ac22b41461039557806379ba5097146103a857806386fa4246146103b05780638da5cb5b146103c357600080fd5b80635e65e3091461035a57806366acaa331461036d5780636ae5c5911461038257600080fd5b8063275459f21161016657806336b402fb1161014057806336b402fb146102cb578063398f37731461031357806350c946fe146103265780635d83d9671461034757600080fd5b8063275459f2146102925780632c01a1e8146102a55780632c3bc7a0146102b857600080fd5b80632145024311610197578063214502431461023d57806322bdbcbc14610252578063235374051461027257600080fd5b806312570011146101be578063181f5a77146101e65780631d05394c14610228575b600080fd5b6101d16101cc36600461378e565b610482565b60405190151581526020015b60405180910390f35b60408051808201909152601881527f4361706162696c697479526567697374727920312e302e30000000000000000060208201525b6040516101dd919061380b565b61023b61023636600461386a565b610495565b005b610245610651565b6040516101dd91906139bc565b610265610260366004613a55565b6107a5565b6040516101dd9190613aad565b610285610280366004613a55565b610892565b6040516101dd9190613ac0565b61023b6102a036600461386a565b6108c5565b61023b6102b336600461386a565b6109a4565b61023b6102c6366004613ad3565b610c8b565b6103056102d9366004613b48565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b6040519081526020016101dd565b61023b61032136600461386a565b610d5b565b61033961033436600461378e565b610f1e565b6040516101dd929190613ba3565b61023b61035536600461386a565b610f53565b61023b61036836600461386a565b61109c565b61037561159f565b6040516101dd9190613bcb565b61023b610390366004613c4c565b61178a565b61023b6103a336600461386a565b61182d565b61023b611cee565b61023b6103be366004613cd0565b611deb565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101dd565b6103fe6103f936600461378e565b612185565b6040516101dd9190613ddb565b61021b610419366004613de9565b61222f565b61023b61042c366004613e13565b612304565b61023b61043f366004613ad3565b612393565b61044c61246d565b6040516101dd9190613ea6565b61046161263c565b6040516101dd929190613ef4565b61023b61047d366004613fd7565b6127c7565b600061048f6005836127db565b92915050565b61049d6127f6565b60005b8181101561064c5760008383838181106104bc576104bc613ff2565b90506020020160208101906104d19190613a55565b63ffffffff8082166000908152600d60209081526040808320805464010000000090049094168084526001850190925282209394509192905b61051382612879565b81101561056b5761055a8563ffffffff16600c600061053b858761288390919063ffffffff16565b815260200190815260200160002060040161288f90919063ffffffff16565b5061056481614050565b905061050a565b508254640100000000900463ffffffff166000036105c2576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff851660048201526024015b60405180910390fd5b63ffffffff84166000818152600d6020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000001690558051938452908301919091527ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c1581703651910160405180910390a1505050508061064590614050565b90506104a0565b505050565b600e54606090640100000000900463ffffffff166000610672600183614088565b63ffffffff1667ffffffffffffffff811115610690576106906140a5565b60405190808252806020026020018201604052801561070657816020015b6040805160a08101825260008082526020808301829052928201526060808201819052608082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816106ae5790505b509050600060015b8363ffffffff168163ffffffff1610156107825763ffffffff8082166000908152600d60205260409020541615610772576107488161289b565b83838151811061075a5761075a613ff2565b60200260200101819052508161076f90614050565b91505b61077b816140d4565b905061070e565b5061078e600184614088565b63ffffffff16811461079e578082525b5092915050565b60408051808201909152600081526060602082015263ffffffff82166000908152600b60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff1683526001810180549192840191610809906140f7565b80601f0160208091040260200160405190810160405280929190818152602001828054610835906140f7565b80156108825780601f1061085757610100808354040283529160200191610882565b820191906000526020600020905b81548152906001019060200180831161086557829003601f168201915b5050505050815250509050919050565b6040805160a081018252600080825260208201819052918101919091526060808201819052608082015261048f8261289b565b6108cd6127f6565b60005b63ffffffff811682111561064c57600083838363ffffffff168181106108f8576108f8613ff2565b905060200201602081019061090d9190613a55565b63ffffffff81166000908152600b6020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001681559192506109586001830182613721565b505060405163ffffffff821681527fa59268ca81d40429e65ccea5385b59cf2d3fc6519371dee92f8eb1dae5107a7a9060200160405180910390a15061099d816140d4565b90506108d0565b6000805473ffffffffffffffffffffffffffffffffffffffff163314905b82811015610c855760008484838181106109de576109de613ff2565b602090810292909201356000818152600c90935260409092206001810154929350911515905080610a3e576040517f64e2ee92000000000000000000000000000000000000000000000000000000008152600481018490526024016105b9565b6000610a4c83600401612879565b1115610a87576040517f34a4a3f6000000000000000000000000000000000000000000000000000000008152600481018490526024016105b9565b815463ffffffff166000908152600b602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff1682526001810180549293919291840191610ad9906140f7565b80601f0160208091040260200160405190810160405280929190818152602001828054610b05906140f7565b8015610b525780601f10610b2757610100808354040283529160200191610b52565b820191906000526020600020905b815481529060010190602001808311610b3557829003601f168201915b505050505081525050905085158015610b825750805173ffffffffffffffffffffffffffffffffffffffff163314155b15610bb9576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001830154610bca9060079061288f565b506002830154610bdc9060099061288f565b506000848152600c6020526040812080547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016815560018101829055600281018290559060048201818181610c31828261375b565b5050505050507f5254e609a97bab37b7cc79fe128f85c097bd6015c6e1624ae0ba392eb975320584604051610c6891815260200190565b60405180910390a15050505080610c7e90614050565b90506109c2565b50505050565b610c936127f6565b60005b8181101561064c576000838383818110610cb257610cb2613ff2565b905060800201803603810190610cc891906141e5565b8051602080830151604080518084019490945283810191909152805180840382018152606090930190528151910120909150610d056003826127db565b610d3e576040517fe181733f000000000000000000000000000000000000000000000000000000008152600481018290526024016105b9565b610d488183612b34565b505080610d5490614050565b9050610c96565b610d636127f6565b60005b8181101561064c576000838383818110610d8257610d82613ff2565b9050602002810190610d94919061423c565b610d9d9061427a565b805190915073ffffffffffffffffffffffffffffffffffffffff16610dee576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e54604080518082018252835173ffffffffffffffffffffffffffffffffffffffff908116825260208086015181840190815263ffffffff9095166000818152600b909252939020825181547fffffffffffffffffffffffff00000000000000000000000000000000000000001692169190911781559251919290916001820190610e7a908261438a565b5050600e8054909150600090610e959063ffffffff166140d4565b91906101000a81548163ffffffff021916908363ffffffff160217905550816000015173ffffffffffffffffffffffffffffffffffffffff167f78e94ca80be2c30abc061b99e7eb8583b1254781734b1e3ce339abb57da2fe8e828460200151604051610f039291906144a4565b60405180910390a2505080610f1790614050565b9050610d66565b6040805160808101825260008082526020820181905291810182905260608082015290610f4a83612d5c565b91509150915091565b610f5b6127f6565b60005b8181101561064c576000838383818110610f7a57610f7a613ff2565b905060200201359050610f978160036127db90919063ffffffff16565b610fd0576040517fe181733f000000000000000000000000000000000000000000000000000000008152600481018290526024016105b9565b610fdb6005826127db565b15611015576040517ff7d7a294000000000000000000000000000000000000000000000000000000008152600481018290526024016105b9565b611020600582612e01565b5060008181526002602081905260408083208381556001810184905590910180547fffffffffffffffffffffff0000000000000000000000000000000000000000001690555182917fdcea1b78b6ddc31592a94607d537543fcaafda6cc52d6d5cc7bbfca1422baf2191a25061109581614050565b9050610f5e565b60005b8181101561064c5760008383838181106110bb576110bb613ff2565b90506020028101906110cd91906144c3565b6110d6906144f7565b905060006110f960005473ffffffffffffffffffffffffffffffffffffffff1690565b825163ffffffff166000908152600b602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff90811683526001820180549690911633149650939491939092840191611156906140f7565b80601f0160208091040260200160405190810160405280929190818152602001828054611182906140f7565b80156111cf5780601f106111a4576101008083540402835291602001916111cf565b820191906000526020600020905b8154815290600101906020018083116111b257829003601f168201915b5050505050815250509050811580156111ff5750805173ffffffffffffffffffffffffffffffffffffffff163314155b15611236576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040808401516000908152600c602052206001015415158061128c5783604001516040517f64e2ee920000000000000000000000000000000000000000000000000000000081526004016105b991815260200190565b602084015115806112d257508360200151600c60008660400151815260200190815260200160002060010154141580156112d2575060208401516112d2906007906127db565b15611309576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6060840151805160000361134b57806040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016105b991906145ca565b6040858101516000908152600c602052208054640100000000900463ffffffff16906004611378836140d4565b82546101009290920a63ffffffff8181021990931691831602179091556040878101516000908152600c60205290812054640100000000900490911691505b8251811015611483576113ed8382815181106113d5576113d5613ff2565b602002602001015160036127db90919063ffffffff16565b61142557826040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016105b991906145ca565b61147283828151811061143a5761143a613ff2565b6020908102919091018101516040808b01516000908152600c845281812063ffffffff808916835260039091019094522091612e0116565b5061147c81614050565b90506113b7565b508551604080880180516000908152600c602090815283822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9096169590951790945581518082528382206002015581518152828120600190810154948b01519251825292902090910154146115385761150760078261288f565b50602080880180516040808b01516000908152600c909452909220600101919091555161153690600790612e01565b505b60408781015188516020808b0151845193845263ffffffff909216908301528183015290517f4b5b465e22eea0c3d40c30e936643245b80d19b2dcf75788c0699fe8d8db645b9181900360600190a1505050505050508061159890614050565b905061109f565b600e5460609063ffffffff1660006115b8600183614088565b63ffffffff1667ffffffffffffffff8111156115d6576115d66140a5565b60405190808252806020026020018201604052801561161c57816020015b6040805180820190915260008152606060208201528152602001906001900390816115f45790505b509050600060015b8363ffffffff168163ffffffff1610156117745763ffffffff81166000908152600b602052604090205473ffffffffffffffffffffffffffffffffffffffff16156117645763ffffffff81166000908152600b60209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff16835260018101805491928401916116b8906140f7565b80601f01602080910402602001604051908101604052809291908181526020018280546116e4906140f7565b80156117315780601f1061170657610100808354040283529160200191611731565b820191906000526020600020905b81548152906001019060200180831161171457829003601f168201915b50505050508152505083838151811061174c5761174c613ff2565b60200260200101819052508161176190614050565b91505b61176d816140d4565b9050611624565b50600e5461078e9060019063ffffffff16614088565b6117926127f6565b600e54640100000000900463ffffffff166000818152600d6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000016821790556117e98160018888888888612e0d565b600e805460049061180790640100000000900463ffffffff166140d4565b91906101000a81548163ffffffff021916908363ffffffff160217905550505050505050565b60005b8181101561064c57600083838381811061184c5761184c613ff2565b905060200281019061185e91906144c3565b611867906144f7565b9050600061188a60005473ffffffffffffffffffffffffffffffffffffffff1690565b825163ffffffff166000908152600b602090815260408083208151808301909252805473ffffffffffffffffffffffffffffffffffffffff908116835260018201805496909116331496509394919390928401916118e7906140f7565b80601f0160208091040260200160405190810160405280929190818152602001828054611913906140f7565b80156119605780601f1061193557610100808354040283529160200191611960565b820191906000526020600020905b81548152906001019060200180831161194357829003601f168201915b50505091909252505081519192505073ffffffffffffffffffffffffffffffffffffffff166119c65782516040517fadd9ae1e00000000000000000000000000000000000000000000000000000000815263ffffffff90911660048201526024016105b9565b811580156119eb5750805173ffffffffffffffffffffffffffffffffffffffff163314155b15611a22576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040808401516000908152600c602052206001015415158080611a4757506040840151155b15611a865783604001516040517f64e2ee920000000000000000000000000000000000000000000000000000000081526004016105b991815260200190565b60208401511580611aa357506020840151611aa3906007906127db565b15611ada576040517f8377314600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608401518051600003611b1c57806040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016105b991906145ca565b6040858101516000908152600c602052208054600490611b4990640100000000900463ffffffff166140d4565b82546101009290920a63ffffffff8181021990931691831602179091556040868101516000908152600c602052908120546401000000009004909116905b8251811015611c0357611ba58382815181106113d5576113d5613ff2565b611bdd57826040517f3748d4c60000000000000000000000000000000000000000000000000000000081526004016105b991906145ca565b611bf283828151811061143a5761143a613ff2565b50611bfc81614050565b9050611b87565b508551604080880180516000908152600c602090815283822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9687161790558251808352848320600201558a018051925182529290206001015551611c759160079190612e0116565b506040860151611c8790600990612e01565b5060408681015187516020808a0151845193845263ffffffff909216908301528183015290517f74becb12a5e8fd0e98077d02dfba8f647c9670c9df177e42c2418cf17a636f059181900360600190a150505050505080611ce790614050565b9050611830565b60015473ffffffffffffffffffffffffffffffffffffffff163314611d6f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016105b9565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b828114611e2e576040517fab8b67c600000000000000000000000000000000000000000000000000000000815260048101849052602481018290526044016105b9565b6000805473ffffffffffffffffffffffffffffffffffffffff16905b8481101561217d576000868683818110611e6657611e66613ff2565b9050602002016020810190611e7b9190613a55565b63ffffffff81166000908152600b602052604090205490915073ffffffffffffffffffffffffffffffffffffffff16611ee8576040517fadd9ae1e00000000000000000000000000000000000000000000000000000000815263ffffffff821660048201526024016105b9565b6000858584818110611efc57611efc613ff2565b9050602002810190611f0e919061423c565b611f179061427a565b805190915073ffffffffffffffffffffffffffffffffffffffff16611f68576040517feeacd93900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805173ffffffffffffffffffffffffffffffffffffffff163314801590611fa557503373ffffffffffffffffffffffffffffffffffffffff851614155b15611fdc576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805163ffffffff83166000908152600b602052604090205473ffffffffffffffffffffffffffffffffffffffff908116911614158061209c5750602080820151604051612029920161380b565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152828252805160209182012063ffffffff86166000908152600b8352929092209192612083926001019101614602565b6040516020818303038152906040528051906020012014155b1561216a57805163ffffffff83166000908152600b6020908152604090912080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909316929092178255820151600190910190612110908261438a565b50806000015173ffffffffffffffffffffffffffffffffffffffff167f86f41145bde5dd7f523305452e4aad3685508c181432ec733d5f345009358a288383602001516040516121619291906144a4565b60405180910390a25b50508061217690614050565b9050611e4a565b505050505050565b604080516080808201835260008083526020808401829052838501829052606084018290528582526002808252918590208551938401865280548452600180820154928501929092529182015493949293919284019160ff16908111156121ee576121ee613d3c565b60018111156121ff576121ff613d3c565b815260029190910154610100900473ffffffffffffffffffffffffffffffffffffffff1660209091015292915050565b63ffffffff8083166000908152600d60209081526040808320805464010000000090049094168084526001909401825280832085845260030190915290208054606092919061227d906140f7565b80601f01602080910402602001604051908101604052809291908181526020018280546122a9906140f7565b80156122f65780601f106122cb576101008083540402835291602001916122f6565b820191906000526020600020905b8154815290600101906020018083116122d957829003601f168201915b505050505091505092915050565b61230c6127f6565b63ffffffff8087166000908152600d6020526040812054640100000000900490911690819003612370576040517f2b62be9b00000000000000000000000000000000000000000000000000000000815263ffffffff881660048201526024016105b9565b61238a8761237d836140d4565b9250828888888888612e0d565b50505050505050565b61239b6127f6565b60005b8181101561064c5760008383838181106123ba576123ba613ff2565b9050608002018036038101906123d091906141e5565b805160208083015160408051808401949094528381019190915280518084038201815260609093019052815191012090915061240d6003826127db565b15612444576040517fe288638f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61244f600382612e01565b5061245a8183612b34565b50508061246690614050565b905061239e565b6060600061247b6003613378565b905060006124896005612879565b825161249591906146ab565b67ffffffffffffffff8111156124ad576124ad6140a5565b60405190808252806020026020018201604052801561251d57816020015b6040805160808101825260008082526020808301829052928201819052606082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816124cb5790505b5090506000805b835181101561263357600084828151811061254157612541613ff2565b6020026020010151905061255f8160056127db90919063ffffffff16565b612622576002600082815260200190815260200160002060405180608001604052908160008201548152602001600182015481526020016002820160009054906101000a900460ff1660018111156125b9576125b9613d3c565b60018111156125ca576125ca613d3c565b815260029190910154610100900473ffffffffffffffffffffffffffffffffffffffff16602090910152845185908590811061260857612608613ff2565b6020026020010181905250828061261e90614050565b9350505b5061262c81614050565b9050612524565b50909392505050565b606080600061264b6009613378565b90506000815167ffffffffffffffff811115612669576126696140a5565b6040519080825280602002602001820160405280156126d857816020015b60408051608081018252600080825260208083018290529282015260608082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816126875790505b5090506000825167ffffffffffffffff8111156126f7576126f76140a5565b604051908082528060200260200182016040528015612720578160200160208202803683370190505b50905060005b83518110156127bc57600084828151811061274357612743613ff2565b6020026020010151905060008061275983612d5c565b915091508186858151811061277057612770613ff2565b60200260200101819052508085858151811061278e5761278e613ff2565b602002602001019063ffffffff16908163ffffffff1681525050505050806127b590614050565b9050612726565b509094909350915050565b6127cf6127f6565b6127d881613385565b50565b600081815260018301602052604081205415155b9392505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612877576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016105b9565b565b600061048f825490565b60006127ef838361347a565b60006127ef83836134a4565b6040805160a081018252600080825260208083018290528284018290526060808401819052608084015263ffffffff8581168352600d825284832080546401000000009004909116808452600190910182528483206002810180548751818602810186019098528088529596929591949390919083018282801561293e57602002820191906000526020600020905b81548152602001906001019080831161292a575b505050505090506000815167ffffffffffffffff811115612961576129616140a5565b6040519080825280602002602001820160405280156129a757816020015b60408051808201909152600081526060602082015281526020019060019003908161297f5790505b50905060005b8151811015612ac85760405180604001604052808483815181106129d3576129d3613ff2565b602002602001015181526020018560030160008685815181106129f8576129f8613ff2565b602002602001015181526020019081526020016000208054612a19906140f7565b80601f0160208091040260200160405190810160405280929190818152602001828054612a45906140f7565b8015612a925780601f10612a6757610100808354040283529160200191612a92565b820191906000526020600020905b815481529060010190602001808311612a7557829003601f168201915b5050505050815250828281518110612aac57612aac613ff2565b602002602001018190525080612ac190614050565b90506129ad565b506040805160a08101825263ffffffff8881166000818152600d6020818152868320548086168752948b168187015292909152905268010000000000000000900460ff1615159181019190915260608101612b2285613378565b81526020019190915295945050505050565b606081015173ffffffffffffffffffffffffffffffffffffffff1615612c8257606081015173ffffffffffffffffffffffffffffffffffffffff163b1580612c2d575060608101516040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527f78bea72100000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff909116906301ffc9a790602401602060405180830381865afa158015612c07573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c2b91906146be565b155b15612c825760608101516040517fabb5e3fd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911660048201526024016105b9565b8060026000848152602001908152602001600020600082015181600001556020820151816001015560408201518160020160006101000a81548160ff02191690836001811115612cd457612cd4613d3c565b0217905550606091909101516002909101805473ffffffffffffffffffffffffffffffffffffffff909216610100027fffffffffffffffffffffff0000000000000000000000000000000000000000ff90921691909117905560405182907f04f0a9bcf3f3a3b42a4d7ca081119755f82ebe43e0d30c8f7292c4fe0dc4a2ae90600090a25050565b6040805160808101825260008082526020820181905291810191909152606080820152604080516080810182526000848152600c6020908152838220805463ffffffff808216865260018301548487015260028301548688015264010000000090910416835260030190529182206060820190612dd890613378565b90526000938452600c602052604090932054929364010000000090930463ffffffff1692915050565b60006127ef8383613597565b63ffffffff8088166000908152600d60209081526040808320938a168084526001948501909252909120911015612ecc5763ffffffff88166000908152600d602052604081206001908101908290612e65908b614088565b63ffffffff1663ffffffff168152602001908152602001600020905060005b612e8d82612879565b811015612ec957612eb88a63ffffffff16600c600061053b858760000161288390919063ffffffff16565b50612ec281614050565b9050612e84565b50505b60005b85811015612fda57612efc878783818110612eec57612eec613ff2565b85926020909102013590506127db565b15612f5d5788878783818110612f1457612f14613ff2565b6040517f636e405700000000000000000000000000000000000000000000000000000000815263ffffffff909416600485015260200291909101356024830152506044016105b9565b612f82878783818110612f7257612f72613ff2565b8592602090910201359050612e01565b50612fc98963ffffffff16600c60008a8a86818110612fa357612fa3613ff2565b905060200201358152602001908152602001600020600401612e0190919063ffffffff16565b50612fd381614050565b9050612ecf565b5060005b838110156132b95736858583818110612ff957612ff9613ff2565b905060200281019061300b919061423c565b9050613019600382356127db565b613052576040517fe181733f000000000000000000000000000000000000000000000000000000008152813560048201526024016105b9565b61305e600582356127db565b15613098576040517ff7d7a294000000000000000000000000000000000000000000000000000000008152813560048201526024016105b9565b80356000908152600384016020526040812080546130b5906140f7565b905011156130fe576040517f3927d08000000000000000000000000000000000000000000000000000000000815263ffffffff8b166004820152813560248201526044016105b9565b60005b87811015613210576131a58235600c60008c8c8681811061312457613124613ff2565b9050602002013581526020019081526020016000206003016000600c60008e8e8881811061315457613154613ff2565b90506020020135815260200190815260200160002060000160049054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020019081526020016000206127db90919063ffffffff16565b613200578888828181106131bb576131bb613ff2565b6040517fa7e7925000000000000000000000000000000000000000000000000000000000815260209091029290920135600483015250823560248201526044016105b9565b61320981614050565b9050613101565b5060028301805460018101825560009182526020918290208335910155613239908201826146db565b82356000908152600386016020526040902091613257919083614740565b506132a88a8a83358b8b61326e60208801886146db565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506135e692505050565b506132b281614050565b9050612fde565b5063ffffffff8881166000818152600d602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffff0000000000ffffffff1668010000000000000000881515027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff1617640100000000958d1695860217905581519283528201929092527ff264aae70bf6a9d90e68e0f9b393f4e7fbea67b063b0f336e0b36c1581703651910160405180910390a15050505050505050565b606060006127ef836136c5565b3373ffffffffffffffffffffffffffffffffffffffff821603613404576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016105b9565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600082600001828154811061349157613491613ff2565b9060005260206000200154905092915050565b6000818152600183016020526040812054801561358d5760006134c86001836146ab565b85549091506000906134dc906001906146ab565b90508181146135415760008660000182815481106134fc576134fc613ff2565b906000526020600020015490508087600001848154811061351f5761351f613ff2565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806135525761355261485b565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061048f565b600091505061048f565b60008181526001830160205260408120546135de5750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561048f565b50600061048f565b60008481526002602081905260409091200154610100900473ffffffffffffffffffffffffffffffffffffffff161561217d57600084815260026020819052604091829020015490517ffba64a7c00000000000000000000000000000000000000000000000000000000815261010090910473ffffffffffffffffffffffffffffffffffffffff169063fba64a7c9061368b908690869086908b908d9060040161488a565b600060405180830381600087803b1580156136a557600080fd5b505af11580156136b9573d6000803e3d6000fd5b50505050505050505050565b60608160000180548060200260200160405190810160405280929190818152602001828054801561371557602002820191906000526020600020905b815481526020019060010190808311613701575b50505050509050919050565b50805461372d906140f7565b6000825580601f1061373d575050565b601f0160209004906000526020600020908101906127d89190613775565b50805460008255906000526020600020908101906127d891905b5b8082111561378a5760008155600101613776565b5090565b6000602082840312156137a057600080fd5b5035919050565b6000815180845260005b818110156137cd576020818501810151868301820152016137b1565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006127ef60208301846137a7565b60008083601f84011261383057600080fd5b50813567ffffffffffffffff81111561384857600080fd5b6020830191508360208260051b850101111561386357600080fd5b9250929050565b6000806020838503121561387d57600080fd5b823567ffffffffffffffff81111561389457600080fd5b6138a08582860161381e565b90969095509350505050565b600081518084526020808501945080840160005b838110156138dc578151875295820195908201906001016138c0565b509495945050505050565b600063ffffffff8083511684526020818185015116818601526040915081840151151582860152606084015160a0606087015261392760a08701826138ac565b9050608085015186820360808801528181518084528484019150848160051b850101858401935060005b828110156139ae578582037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001845284518051835287015187830189905261399b898401826137a7565b9588019594880194925050600101613951565b509998505050505050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015613a2f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452613a1d8583516138e7565b945092850192908501906001016139e3565b5092979650505050505050565b803563ffffffff81168114613a5057600080fd5b919050565b600060208284031215613a6757600080fd5b6127ef82613a3c565b73ffffffffffffffffffffffffffffffffffffffff81511682526000602082015160406020850152613aa560408501826137a7565b949350505050565b6020815260006127ef6020830184613a70565b6020815260006127ef60208301846138e7565b60008060208385031215613ae657600080fd5b823567ffffffffffffffff80821115613afe57600080fd5b818501915085601f830112613b1257600080fd5b813581811115613b2157600080fd5b8660208260071b8501011115613b3657600080fd5b60209290920196919550909350505050565b60008060408385031215613b5b57600080fd5b50508035926020909101359150565b63ffffffff815116825260208101516020830152604081015160408301526000606082015160806060850152613aa560808501826138ac565b604081526000613bb66040830185613b6a565b905063ffffffff831660208301529392505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015613a2f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452613c2c858351613a70565b94509285019290850190600101613bf2565b80151581146127d857600080fd5b600080600080600060608688031215613c6457600080fd5b853567ffffffffffffffff80821115613c7c57600080fd5b613c8889838a0161381e565b90975095506020880135915080821115613ca157600080fd5b50613cae8882890161381e565b9094509250506040860135613cc281613c3e565b809150509295509295909350565b60008060008060408587031215613ce657600080fd5b843567ffffffffffffffff80821115613cfe57600080fd5b613d0a8883890161381e565b90965094506020870135915080821115613d2357600080fd5b50613d308782880161381e565b95989497509550505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8051825260208101516020830152604081015160028110613db5577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b604083015260609081015173ffffffffffffffffffffffffffffffffffffffff16910152565b6080810161048f8284613d6b565b60008060408385031215613dfc57600080fd5b613e0583613a3c565b946020939093013593505050565b60008060008060008060808789031215613e2c57600080fd5b613e3587613a3c565b9550602087013567ffffffffffffffff80821115613e5257600080fd5b613e5e8a838b0161381e565b90975095506040890135915080821115613e7757600080fd5b50613e8489828a0161381e565b9094509250506060870135613e9881613c3e565b809150509295509295509295565b6020808252825182820181905260009190848201906040850190845b81811015613ee857613ed5838551613d6b565b9284019260809290920191600101613ec2565b50909695505050505050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015613f69577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552613f57868351613b6a565b95509382019390820190600101613f1d565b50508584038187015286518085528782019482019350915060005b82811015613fa657845163ffffffff1684529381019392810192600101613f84565b5091979650505050505050565b803573ffffffffffffffffffffffffffffffffffffffff81168114613a5057600080fd5b600060208284031215613fe957600080fd5b6127ef82613fb3565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361408157614081614021565b5060010190565b63ffffffff82811682821603908082111561079e5761079e614021565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600063ffffffff8083168181036140ed576140ed614021565b6001019392505050565b600181811c9082168061410b57607f821691505b602082108103614144577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b6040516080810167ffffffffffffffff8111828210171561416d5761416d6140a5565b60405290565b6040805190810167ffffffffffffffff8111828210171561416d5761416d6140a5565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156141dd576141dd6140a5565b604052919050565b6000608082840312156141f757600080fd5b6141ff61414a565b823581526020830135602082015260408301356002811061421f57600080fd5b604082015261423060608401613fb3565b60608201529392505050565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc183360301811261427057600080fd5b9190910192915050565b60006040823603121561428c57600080fd5b614294614173565b61429d83613fb3565b815260208084013567ffffffffffffffff808211156142bb57600080fd5b9085019036601f8301126142ce57600080fd5b8135818111156142e0576142e06140a5565b614310847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601614196565b9150808252368482850101111561432657600080fd5b80848401858401376000908201840152918301919091525092915050565b601f82111561064c57600081815260208120601f850160051c8101602086101561436b5750805b601f850160051c820191505b8181101561217d57828155600101614377565b815167ffffffffffffffff8111156143a4576143a46140a5565b6143b8816143b284546140f7565b84614344565b602080601f83116001811461440b57600084156143d55750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855561217d565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561445857888601518255948401946001909101908401614439565b508582101561449457878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b63ffffffff83168152604060208201526000613aa560408301846137a7565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8183360301811261427057600080fd5b60006080823603121561450957600080fd5b61451161414a565b61451a83613a3c565b81526020808401358183015260408401356040830152606084013567ffffffffffffffff8082111561454b57600080fd5b9085019036601f83011261455e57600080fd5b813581811115614570576145706140a5565b8060051b9150614581848301614196565b818152918301840191848101903684111561459b57600080fd5b938501935b838510156145b9578435825293850193908501906145a0565b606087015250939695505050505050565b6020808252825182820181905260009190848201906040850190845b81811015613ee8578351835292840192918401916001016145e6565b6000602080835260008454614616816140f7565b80848701526040600180841660008114614637576001811461466f5761469d565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838a01528284151560051b8a0101955061469d565b896000528660002060005b858110156146955781548b820186015290830190880161467a565b8a0184019650505b509398975050505050505050565b8181038181111561048f5761048f614021565b6000602082840312156146d057600080fd5b81516127ef81613c3e565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261471057600080fd5b83018035915067ffffffffffffffff82111561472b57600080fd5b60200191503681900382131561386357600080fd5b67ffffffffffffffff831115614758576147586140a5565b61476c8361476683546140f7565b83614344565b6000601f8411600181146147be57600085156147885750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355614854565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b8281101561480d57868501358255602094850194600190920191016147ed565b5086821015614848577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6080815284608082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8611156148c357600080fd5b8560051b808860a0850137820182810360a090810160208501526148e9908201876137a7565b91505063ffffffff8085166040840152808416606084015250969550505050505056fea164736f6c6343000813000a",
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) AddDON(opts *bind.TransactOpts, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "addDON", nodes, capabilityConfigurations, isPublic)
}

func (_CapabilityRegistry *CapabilityRegistrySession) AddDON(nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddDON(&_CapabilityRegistry.TransactOpts, nodes, capabilityConfigurations, isPublic)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) AddDON(nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.AddDON(&_CapabilityRegistry.TransactOpts, nodes, capabilityConfigurations, isPublic)
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

func (_CapabilityRegistry *CapabilityRegistryTransactor) UpdateCapabilities(opts *bind.TransactOpts, capabilities []CapabilityRegistryCapability) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "updateCapabilities", capabilities)
}

func (_CapabilityRegistry *CapabilityRegistrySession) UpdateCapabilities(capabilities []CapabilityRegistryCapability) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.UpdateCapabilities(&_CapabilityRegistry.TransactOpts, capabilities)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) UpdateCapabilities(capabilities []CapabilityRegistryCapability) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.UpdateCapabilities(&_CapabilityRegistry.TransactOpts, capabilities)
}

func (_CapabilityRegistry *CapabilityRegistryTransactor) UpdateDON(opts *bind.TransactOpts, donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool) (*types.Transaction, error) {
	return _CapabilityRegistry.contract.Transact(opts, "updateDON", donId, nodes, capabilityConfigurations, isPublic)
}

func (_CapabilityRegistry *CapabilityRegistrySession) UpdateDON(donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.UpdateDON(&_CapabilityRegistry.TransactOpts, donId, nodes, capabilityConfigurations, isPublic)
}

func (_CapabilityRegistry *CapabilityRegistryTransactorSession) UpdateDON(donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool) (*types.Transaction, error) {
	return _CapabilityRegistry.Contract.UpdateDON(&_CapabilityRegistry.TransactOpts, donId, nodes, capabilityConfigurations, isPublic)
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
	NodeOperatorId uint32
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
	NodeOperatorId uint32
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
	NodeOperatorId uint32
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
	NodeOperatorId uint32
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

	GetDON(opts *bind.CallOpts, donId uint32) (CapabilityRegistryDONInfo, error)

	GetDONCapabilityConfig(opts *bind.CallOpts, donId uint32, capabilityId [32]byte) ([]byte, error)

	GetDONs(opts *bind.CallOpts) ([]CapabilityRegistryDONInfo, error)

	GetHashedCapabilityId(opts *bind.CallOpts, labelledName [32]byte, version [32]byte) ([32]byte, error)

	GetNode(opts *bind.CallOpts, p2pId [32]byte) (CapabilityRegistryNodeInfo, uint32, error)

	GetNodeOperator(opts *bind.CallOpts, nodeOperatorId uint32) (CapabilityRegistryNodeOperator, error)

	GetNodeOperators(opts *bind.CallOpts) ([]CapabilityRegistryNodeOperator, error)

	GetNodes(opts *bind.CallOpts) ([]CapabilityRegistryNodeInfo, []uint32, error)

	IsCapabilityDeprecated(opts *bind.CallOpts, hashedCapabilityId [32]byte) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddCapabilities(opts *bind.TransactOpts, capabilities []CapabilityRegistryCapability) (*types.Transaction, error)

	AddDON(opts *bind.TransactOpts, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool) (*types.Transaction, error)

	AddNodeOperators(opts *bind.TransactOpts, nodeOperators []CapabilityRegistryNodeOperator) (*types.Transaction, error)

	AddNodes(opts *bind.TransactOpts, nodes []CapabilityRegistryNodeInfo) (*types.Transaction, error)

	DeprecateCapabilities(opts *bind.TransactOpts, hashedCapabilityIds [][32]byte) (*types.Transaction, error)

	RemoveDONs(opts *bind.TransactOpts, donIds []uint32) (*types.Transaction, error)

	RemoveNodeOperators(opts *bind.TransactOpts, nodeOperatorIds []uint32) (*types.Transaction, error)

	RemoveNodes(opts *bind.TransactOpts, removedNodeP2PIds [][32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateCapabilities(opts *bind.TransactOpts, capabilities []CapabilityRegistryCapability) (*types.Transaction, error)

	UpdateDON(opts *bind.TransactOpts, donId uint32, nodes [][32]byte, capabilityConfigurations []CapabilityRegistryCapabilityConfiguration, isPublic bool) (*types.Transaction, error)

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
