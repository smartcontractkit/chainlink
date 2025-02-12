// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ccip_home

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

type CCIPHomeChainConfig struct {
	Readers [][32]byte
	FChain  uint8
	Config  []byte
}

type CCIPHomeChainConfigArgs struct {
	ChainSelector uint64
	ChainConfig   CCIPHomeChainConfig
}

type CCIPHomeOCR3Config struct {
	PluginType            uint8
	ChainSelector         uint64
	FRoleDON              uint8
	OffchainConfigVersion uint64
	OfframpAddress        []byte
	RmnHomeAddress        []byte
	Nodes                 []CCIPHomeOCR3Node
	OffchainConfig        []byte
}

type CCIPHomeOCR3Node struct {
	P2pId          [32]byte
	SignerKey      []byte
	TransmitterKey []byte
}

type CCIPHomeVersionedConfig struct {
	Version      uint32
	ConfigDigest [32]byte
	Config       CCIPHomeOCR3Config
}

var CCIPHomeMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"capabilitiesRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyChainConfigUpdates\",\"inputs\":[{\"name\":\"chainSelectorRemoves\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"chainConfigAdds\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPHome.ChainConfigArgs[]\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"chainConfig\",\"type\":\"tuple\",\"internalType\":\"structCCIPHome.ChainConfig\",\"components\":[{\"name\":\"readers\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"fChain\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"config\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"beforeCapabilityConfigSet\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"update\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"donId\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getActiveDigest\",\"inputs\":[{\"name\":\"donId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllChainConfigs\",\"inputs\":[{\"name\":\"pageIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pageSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPHome.ChainConfigArgs[]\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"chainConfig\",\"type\":\"tuple\",\"internalType\":\"structCCIPHome.ChainConfig\",\"components\":[{\"name\":\"readers\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"fChain\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"config\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllConfigs\",\"inputs\":[{\"name\":\"donId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"}],\"outputs\":[{\"name\":\"activeConfig\",\"type\":\"tuple\",\"internalType\":\"structCCIPHome.VersionedConfig\",\"components\":[{\"name\":\"version\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structCCIPHome.OCR3Config\",\"components\":[{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"},{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"FRoleDON\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offrampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rmnHomeAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPHome.OCR3Node[]\",\"components\":[{\"name\":\"p2pId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signerKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transmitterKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"candidateConfig\",\"type\":\"tuple\",\"internalType\":\"structCCIPHome.VersionedConfig\",\"components\":[{\"name\":\"version\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structCCIPHome.OCR3Config\",\"components\":[{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"},{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"FRoleDON\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offrampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rmnHomeAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPHome.OCR3Node[]\",\"components\":[{\"name\":\"p2pId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signerKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transmitterKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCandidateDigest\",\"inputs\":[{\"name\":\"donId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCapabilityConfiguration\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"configuration\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getCapabilityRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getChainConfig\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structCCIPHome.ChainConfig\",\"components\":[{\"name\":\"readers\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"fChain\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"config\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfig\",\"inputs\":[{\"name\":\"donId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"versionedConfig\",\"type\":\"tuple\",\"internalType\":\"structCCIPHome.VersionedConfig\",\"components\":[{\"name\":\"version\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structCCIPHome.OCR3Config\",\"components\":[{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"},{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"FRoleDON\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offrampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rmnHomeAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPHome.OCR3Node[]\",\"components\":[{\"name\":\"p2pId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signerKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transmitterKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"ok\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigDigests\",\"inputs\":[{\"name\":\"donId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"}],\"outputs\":[{\"name\":\"activeConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"candidateConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNumChainConfigurations\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"promoteCandidateAndRevokeActive\",\"inputs\":[{\"name\":\"donId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"},{\"name\":\"digestToPromote\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"digestToRevoke\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeCandidate\",\"inputs\":[{\"name\":\"donId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setCandidate\",\"inputs\":[{\"name\":\"donId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"},{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structCCIPHome.OCR3Config\",\"components\":[{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"},{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"FRoleDON\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offrampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rmnHomeAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPHome.OCR3Node[]\",\"components\":[{\"name\":\"p2pId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signerKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transmitterKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"digestToOverwrite\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"newConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ActiveConfigRevoked\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CandidateConfigRevoked\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CapabilityConfigurationSet\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainConfigRemoved\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainConfigSet\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"chainConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structCCIPHome.ChainConfig\",\"components\":[{\"name\":\"readers\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"fChain\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"config\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigPromoted\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"version\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structCCIPHome.OCR3Config\",\"components\":[{\"name\":\"pluginType\",\"type\":\"uint8\",\"internalType\":\"enumInternal.OCRPluginType\"},{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"FRoleDON\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offrampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rmnHomeAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPHome.OCR3Node[]\",\"components\":[{\"name\":\"p2pId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signerKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transmitterKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CanOnlySelfCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainSelectorNotFound\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ChainSelectorNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[{\"name\":\"expectedConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"gotConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DONIdMismatch\",\"inputs\":[{\"name\":\"callDonId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"capabilityRegistryDonId\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"FChainMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FChainTooHigh\",\"inputs\":[{\"name\":\"fChain\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"FRoleDON\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"FTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidNode\",\"inputs\":[{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structCCIPHome.OCR3Node\",\"components\":[{\"name\":\"p2pId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signerKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transmitterKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"type\":\"error\",\"name\":\"InvalidPluginType\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSelector\",\"inputs\":[{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoOpStateTransitionNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotEnoughTransmitters\",\"inputs\":[{\"name\":\"got\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minimum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OfframpAddressCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCapabilitiesRegistryCanCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RMNHomeAddressCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RevokingZeroDigestNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooManySigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x60a03460bf57601f613d6638819003918201601f19168301916001600160401b0383118484101760c45780849260209460405283398101031260bf57516001600160a01b03811680820360bf57331560ae57600180546001600160a01b031916331790556006805463ffffffff1916905515609d57608052604051613c8b90816100db823960805181818161026601528181612f3901526139bd0152f35b6342bcdf7f60e11b60005260046000fd5b639b15e16f60e01b60005260046000fd5b600080fd5b634e487b7160e01b600052604160045260246000fdfe6080604052600436101561001257600080fd5b60003560e01c806301ffc9a714610157578063020330e614610152578063181f5a771461014d57806333d9704a146101485780633df45a72146101435780634851d5491461013e5780635a837f97146101395780635f1edd9c146101345780637524051a1461012f57806379ba50971461012a5780637ac0d41e146101255780638318ed5d146101205780638da5cb5b1461011b578063922ea40614610116578063b149092b14610111578063b74b23561461010c578063bae4e0fa14610107578063f2fde38b14610102578063f442c89a146100fd5763fba64a7c146100f857600080fd5b61172d565b611472565b61134e565b6110e5565b61100a565b610f90565b610ee2565b610e90565b610e2f565b610df3565b610d0a565b610c0f565b610bbb565b610934565b6108ae565b6107ce565b61072c565b610415565b61021b565b346102165760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610216576004357fffffffff00000000000000000000000000000000000000000000000000000000811680910361021657807f78bea72100000000000000000000000000000000000000000000000000000000602092149081156101ec575b506040519015158152f35b7f01ffc9a700000000000000000000000000000000000000000000000000000000915014386101e1565b600080fd5b346102165760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261021657602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6060810190811067ffffffffffffffff8211176102d557604052565b61028a565b610100810190811067ffffffffffffffff8211176102d557604052565b6040810190811067ffffffffffffffff8211176102d557604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176102d557604052565b60405190610363604083610313565b565b6040519061036361010083610313565b67ffffffffffffffff81116102d557601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60005b8381106103c25750506000910152565b81810151838201526020016103b2565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f60209361040e815180928187528780880191016103af565b0116010190565b346102165760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102165761049260408051906104568183610313565b600e82527f43434950486f6d6520312e362e300000000000000000000000000000000000006020830152519182916020835260208301906103d2565b0390f35b63ffffffff81160361021657565b6064359061036382610496565b6002111561021657565b3590610363826104b1565b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc6060910112610216576004356104fc81610496565b90602435610509816104b1565b9060443590565b6002111561051a57565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b90600282101561051a5752565b61058a9181518152604061057960208401516060602085015260608401906103d2565b9201519060408184039101526103d2565b90565b9080602083519182815201916020808360051b8301019401926000915b8383106105b957505050505090565b90919293946020806105f5837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086600196030187528951610556565b970193019301919392906105aa565b90604061058a9263ffffffff81511683526020810151602084015201519060606040820152610637606082018351610549565b602082015167ffffffffffffffff166080820152604082015160ff1660a0820152606082015167ffffffffffffffff1660c082015260e06106f86106c361068e6080860151610100858701526101608601906103d2565b60a08601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0868303016101008701526103d2565b60c08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08583030161012086015261058d565b920151906101407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0828503019101526103d2565b346102165761075a610746610740366104c6565b91611b23565b604051928392604084526040840190610604565b90151560208301520390f35b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc60409101126102165760043561079c81610496565b9060243561058a816104b1565b90916107c061058a93604084526040840190610604565b916020818403910152610604565b34610216576107dc36610766565b906107e56117dc565b906107ee6117dc565b9261083d61083763ffffffff841680600052600560205261081384604060002061183b565b90600052600760205263ffffffff61082f85604060002061183b565b541690611882565b50611a52565b60208101516108a4575b508161087c82610876610837946108716108829763ffffffff166000526005602052604060002090565b61183b565b9261312a565b90611882565b602081015161089c575b50610492604051928392836107a9565b91503861088c565b9250610882610847565b346102165761091f60016108c136610766565b929061087c63ffffffff821694856000526005602052846109056108e983604060002061183b565b88600052600760205263ffffffff61082f85604060002061183b565b50015495600052600560205261087681604060002061183b565b50015460408051928352602083019190915290f35b346102165760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102165760043561096f81610496565b6024359061097c826104b1565b604435916064359161098c613150565b831580610bb3575b610b89576109ae6109a5838361312a565b63ffffffff1690565b8460016109d8836109d3876108718863ffffffff166000526005602052604060002090565b611882565b50015403610b2f57506001610a2f610a04846108718563ffffffff166000526005602052604060002090565b61087c610a25866108718763ffffffff166000526007602052604060002090565b5463ffffffff1690565b50018054848103610af9575091610871610a60926000610aa0955563ffffffff166000526007602052604060002090565b6001610a70825463ffffffff1690565b1863ffffffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000825416179055565b80610ace575b507ffc3e98dbbd47c3fa7c1c05b6ec711caeaf70eca4554192b9ada8fc11a37f298e600080a2005b7f0b31c0055e2d464bef7781994b98c4ff9ef4ae0d05f59feb6a68c42de5e201b8600080a238610aa6565b7f93df584c00000000000000000000000000000000000000000000000000000000600052600452602484905260446000fd5b6000fd5b610b576001916109d3610b2b95610871899663ffffffff166000526005602052604060002090565b5001547f93df584c00000000000000000000000000000000000000000000000000000000600052600452602452604490565b7f7b4d1e4f0000000000000000000000000000000000000000000000000000000060005260046000fd5b508215610994565b346102165760206001610c0463ffffffff8061082f610bd936610766565b9316928360005260058752610bf281604060002061183b565b9360005260078752604060002061183b565b500154604051908152f35b3461021657610c1d366104c6565b91610c26613150565b8215610ce05763ffffffff610c3b838361312a565b169263ffffffff82166000526005602052806001610c61866109d387604060002061183b565b50015403610cb957926109d3600193610871610cb4946000977f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b8980a263ffffffff166000526005602052604060002090565b500155005b6001610b57856109d386610871610b2b9763ffffffff166000526005602052604060002090565b7f0849d8cc0000000000000000000000000000000000000000000000000000000060005260046000fd5b346102165760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102165760005473ffffffffffffffffffffffffffffffffffffffff81163303610dc9577fffffffffffffffffffffffff00000000000000000000000000000000000000006001549133828416176001551660005573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b346102165760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610216576020600354604051908152f35b346102165760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261021657610e69600435610496565b6040516020610e788183610313565b600082526104926040519282849384528301906103d2565b346102165760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261021657602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b34610216576020610efb610ef536610766565b90611bed565b604051908152f35b67ffffffffffffffff81160361021657565b6044359061036382610f03565b359061036382610f03565b9190606081019083519160608252825180915260206080830193019060005b818110610f7a5750505060408460ff602061058a9697015116602084015201519060408184039101526103d2565b8251855260209485019490920191600101610f4c565b346102165760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102165767ffffffffffffffff600435610fd481610f03565b610fdc611c1a565b50166000526002602052610492610ff66040600020611c3a565b604051918291602083526020830190610f2d565b346102165760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261021657611047602435600435611e24565b6040518091602082016020835281518091526040830190602060408260051b8601019301916000905b82821061107f57505050500390f35b919360206110d5827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc060019597998495030186526040838a5167ffffffffffffffff815116845201519181858201520190610f2d565b9601920192018594939192611070565b346102165760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102165760043561112081610496565b60243561112c816104b1565b60443567ffffffffffffffff8111610216576101007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc826004019236030112610216576064359261117b613150565b61118d611188368461206f565b613231565b6111978382611bed565b9380850361131c57917f94f085b7c57ec2a270befd0b7b2ec7452580040edee8bb0fb04609c81f0359c69161087c9493610492966112f1575b506112cf8260026112936111f16111ec60065463ffffffff1690565b612158565b946112278663ffffffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000006006541617600655565b6112728660405161126b8161123f8960208301612418565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101835282610313565b8b846135b0565b99896108768c9b6108718563ffffffff166000526005602052604060002090565b506001810188905580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff861617815501612936565b6112de60405192839283612abe565b0390a26040519081529081906020820190565b7f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b600080a2386111d0565b7f93df584c00000000000000000000000000000000000000000000000000000000600052600485905260245260446000fd5b346102165760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102165760043573ffffffffffffffffffffffffffffffffffffffff8116809103610216576113a66136c5565b33811461141757807fffffffffffffffffffffffff0000000000000000000000000000000000000000600054161760005573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b9181601f840112156102165782359167ffffffffffffffff8311610216576020808501948460051b01011161021657565b346102165760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102165760043567ffffffffffffffff8111610216576114c1903690600401611441565b60243567ffffffffffffffff8111610216576114e1903690600401611441565b9190926114ec6136c5565b60005b8281106116015750505060005b81811061150557005b611525611520611516838587612b7a565b6020810190612662565b612bba565b90611539611534828587612b7a565b612429565b6115438351613948565b61155a611554602085015160ff1690565b60ff1690565b156115d75782816115ab6001956115a67f05dd57854af2c291a94ea52e7c43d80bc3be7fa73022f98b735dea86642fa5e09567ffffffffffffffff166000526002602052604060002090565b612d9a565b6115be67ffffffffffffffff8216613bed565b506115ce60405192839283612e2d565b0390a1016114fc565b7fa9b3766e0000000000000000000000000000000000000000000000000000000060005260046000fd5b61163c611638611625611618611534858888612adb565b67ffffffffffffffff1690565b6000526004602052604060002054151590565b1590565b6116e657806116766116716116576115346001958888612adb565b67ffffffffffffffff166000526002602052604060002090565b612b33565b61168f61168a611618611534848888612adb565b613b09565b507f2a680691fef3b2d105196805935232c661ce703e92d464ef0b94a7bc62d714f06116dd6116c2611534848888612adb565b60405167ffffffffffffffff90911681529081906020820190565b0390a1016114ef565b61153490610b2b936116f793612adb565b7f1bd4d2d20000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff16600452602490565b346102165760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102165760043567ffffffffffffffff81116102165761177c903690600401611441565b505060243567ffffffffffffffff811161021657366023820112156102165780600401359067ffffffffffffffff8211610216573660248383010111610216576117da916117c8610f15565b5060246117d36104a4565b9201612f20565b005b604051906117e9826102b9565b8160008152600060208201526040805191611803836102da565b60008352600060208401526000828401526000606084015260606080840152606060a0840152606060c0840152606060e08401520152565b90600281101561051a57600052602052604060002090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b906002811015611896576007020190600090565b611853565b600282101561051a5752565b90600182811c921680156118f0575b60208310146118c157565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f16916118b6565b906040519182600082549261190e846118a7565b808452936001811690811561197a5750600114611933575b5061036392500383610313565b90506000929192526020600020906000915b81831061195e5750509060206103639282010138611926565b6020919350806001915483858901015201910190918492611945565b602093506103639592507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b82010138611926565b67ffffffffffffffff81116102d55760051b60200190565b9081546119de816119ba565b926119ec6040519485610313565b818452602084019060005260206000206000915b838310611a0d5750505050565b60036020600192604051611a20816102b9565b85548152611a2f8587016118fa565b83820152611a3f600287016118fa565b6040820152815201920192019190611a00565b9060405191611a60836102b9565b60408363ffffffff835416815260018301546020820152611b1a6006835194611a88866102da565b611ae1611ad06002830154611aa060ff82168a61189b565b67ffffffffffffffff600882901c1660208a015260ff604882901c16888a015260501c67ffffffffffffffff1690565b67ffffffffffffffff166060880152565b611aed600382016118fa565b6080870152611afe600482016118fa565b60a0870152611b0f600582016119d2565b60c0870152016118fa565b60e08401520152565b90611b2c6117dc565b9260005b60028110611b42575050505090600090565b63ffffffff8416806000526005602052826001611b67846109d388604060002061183b565b5001541480611ba8575b611b7e5750600101611b30565b611ba2955061083794506109d392506000939193526005602052604060002061183b565b90600190565b50821515611b71565b91611be9918354907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b9055565b611c149061087c60019363ffffffff8316600052600560205261087681604060002061183b565b50015490565b60405190611c27826102b9565b6060604083828152600060208201520152565b90604051611c47816102b9565b809260405180602083549182815201908360005260206000209060005b818110611cab5750505060409282611c83611ca6946002940382610313565b8552611ca0611c96600183015460ff1690565b60ff166020870152565b016118fa565b910152565b8254845260209093019260019283019201611c64565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9081600302916003830403611d0157565b611cc1565b81810292918115918404141715611d0157565b60405190611d28602083610313565b600080835282815b828110611d3c57505050565b602090604051611d4b816102f7565b60008152611d57611c1a565b8382015282828501015201611d30565b90611d71826119ba565b611d7e6040519182610313565b8281527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0611dac82946119ba565b019060005b828110611dbd57505050565b602090604051611dcc816102f7565b60008152611dd8611c1a565b8382015282828501015201611db1565b9060018201809211611d0157565b91908201809211611d0157565b91908203918211611d0157565b80518210156118965760209160051b010190565b611e318260035492611d06565b9180158015611f03575b611ef857611e499083611df6565b90808211611ef0575b50611e65611e608383611e03565b611d67565b91805b828110611e755750505090565b80611ee9611e87611618600194613a41565b611ec8611ea88267ffffffffffffffff166000526002602052604060002090565b611ec3611eb3610354565b67ffffffffffffffff9094168452565b611c3a565b6020820152611ed78584611e03565b90611ee28289611e10565b5286611e10565b5001611e68565b905038611e52565b50505061058a611d19565b5081831015611e3b565b60ff81160361021657565b359061036382611f0d565b81601f8201121561021657803590611f3a82610375565b92611f486040519485610313565b8284526020838301011161021657816000926020809301838601378301015290565b9080601f8301121561021657813591611f82836119ba565b92611f906040519485610313565b80845260208085019160051b830101918383116102165760208101915b838310611fbc57505050505090565b823567ffffffffffffffff81116102165782019060607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083880301126102165760405190612009826102b9565b60208301358252604083013567ffffffffffffffff81116102165787602061203392860101611f23565b602083015260608301359167ffffffffffffffff83116102165761205f88602080969581960101611f23565b6040820152815201920191611fad565b9190916101008184031261021657612085610365565b9261208f826104bb565b845261209d60208301610f22565b60208501526120ae60408301611f18565b60408501526120bf60608301610f22565b6060850152608082013567ffffffffffffffff811161021657816120e4918401611f23565b608085015260a082013567ffffffffffffffff81116102165781612109918401611f23565b60a085015260c082013567ffffffffffffffff8111610216578161212e918401611f6a565b60c085015260e082013567ffffffffffffffff8111610216576121519201611f23565b60e0830152565b63ffffffff1663ffffffff8114611d015760010190565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18236030181121561021657016020813591019167ffffffffffffffff821161021657813603831361021657565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18236030181121561021657016020813591019167ffffffffffffffff8211610216578160051b3603831361021657565b90602083828152019060208160051b85010193836000915b8383106122795750505050505090565b9091929394957fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082820301865286357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa1843603018112156102165760206123206001938683940190813581526123126123076122f78685018561216f565b60608886015260608501916121bf565b92604081019061216f565b9160408185039101526121bf565b980196019493019190612269565b61058a916123448161233f846104bb565b610549565b61236461235360208401610f22565b67ffffffffffffffff166020830152565b61237d61237360408401611f18565b60ff166040830152565b61239d61238c60608401610f22565b67ffffffffffffffff166060830152565b61240a6123ff6123e46123c96123b6608087018761216f565b61010060808801526101008701916121bf565b6123d660a087018761216f565b9086830360a08801526121bf565b6123f160c08601866121fe565b9085830360c0870152612251565b9260e081019061216f565b9160e08185039101526121bf565b90602061058a92818152019061232e565b3561058a81610f03565b3561058a81611f0d565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215610216570180359067ffffffffffffffff82116102165760200191813603831361021657565b818110612499575050565b6000815560010161248e565b9190601f81116124b457505050565b610363926000526020600020906020601f840160051c830193106124e0575b601f0160051c019061248e565b90915081906124d3565b90929167ffffffffffffffff81116102d5576125108161250a84546118a7565b846124a5565b6000601f821160011461256a578190611be993949560009261255f575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b01359050388061252d565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082169461259d84600052602060002090565b91805b8781106125f65750836001959697106125be575b505050811b019055565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88560031b161c199101351690553880806125b4565b909260206001819286860135815501940191016125a0565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215610216570180359067ffffffffffffffff821161021657602001918160051b3603831361021657565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa181360301821215610216570190565b61269f81546118a7565b90816126a9575050565b81601f600093116001146126bb575055565b818352602083206126d791601f0160051c81019060010161248e565b808252602082209081548360011b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8560031b1c191617905555565b90803582556001820161272a602083018361243d565b9067ffffffffffffffff82116102d55761274e8261274885546118a7565b856124a5565b600090601f83116001146127bd57926127a7836127b49460029794610363999760009261255f5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b604081019061243d565b929091016124ea565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08316916127f085600052602060002090565b92815b81811061285957509360029693610363989693600193836127b49810612821575b505050811b0190556127aa565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88560031b161c19910135169055388080612814565b919360206001819287870135815501950192016127f3565b6801000000000000000083116102d55780548382558084106128d9575b50906128a08192600052602060002090565b906000925b8484106128b3575050505050565b60036020826128cd6128c760019587612662565b87612714565b019301930192916128a5565b80600302906003820403611d015783600302600381048503611d015782600052602060002091820191015b818110612911575061288e565b6003906000815561292460018201612695565b61293060028201612695565b01612904565b90803591612943836104b1565b600283101561051a576127b46004926103639460ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0085541691161783556129ca61299060208301612429565b84547fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000ff1660089190911b68ffffffffffffffff0016178455565b612a146129d960408301612433565b84547fffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffff1660489190911b69ff00000000000000000016178455565b612a66612a2360608301612429565b84547fffffffffffffffffffffffffffff0000000000000000ffffffffffffffffffff1660509190911b71ffffffffffffffff0000000000000000000016178455565b612a80612a76608083018361243d565b90600186016124ea565b612a9a612a9060a083018361243d565b90600286016124ea565b612ab4612aaa60c083018361260e565b9060038601612871565b60e081019061243d565b60409063ffffffff61058a9493168152816020820152019061232e565b91908110156118965760051b0190565b906801000000000000000081116102d557815491818155828210612b0e57505050565b600052602060002091820191015b818110612b27575050565b60008155600101612b1c565b80546000825580612b53575b506002816000600161036394015501612695565b816000526020600020908101905b818110612b6e5750612b3f565b60008155600101612b61565b91908110156118965760051b810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc181360301821215610216570190565b6060813603126102165760405190612bd1826102b9565b803567ffffffffffffffff811161021657810136601f8201121561021657803590612bfb826119ba565b91612c096040519384610313565b80835260208084019160051b8301019136831161021657602001905b828210612c6b575050508252612c3d60208201611f18565b602083015260408101359067ffffffffffffffff821161021657612c6391369101611f23565b604082015290565b8135815260209182019101612c25565b919091825167ffffffffffffffff81116102d557612c9d8161250a84546118a7565b6020601f8211600114612cf6578190611be9939495600092612ceb5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b01519050388061252d565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0821690612d2984600052602060002090565b9160005b818110612d8257509583600195969710612d4b57505050811b019055565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690553880806125b4565b9192602060018192868b015181550194019201612d2d565b90805180519067ffffffffffffffff82116102d557602090612dbc8386612aeb565b0183600052602060002060005b838110612e1957505050509060026040610363936001840160ff6020830151167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082541617905501519101612c7b565b600190602084519401938184015501612dc9565b60409067ffffffffffffffff61058a94931681528160208201520190610f2d565b906004116102165790600490565b906024116102165760040190602090565b919091357fffffffff0000000000000000000000000000000000000000000000000000000081169260048110612ea1575050565b7fffffffff00000000000000000000000000000000000000000000000000000000929350829060040360031b1b161690565b90816020910312610216573590565b908092918237016000815290565b3d15612f1b573d90612f0182610375565b91612f0f6040519384610313565b82523d6000602084013e565b606090565b909173ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016330361310057612f72612f6c8484612e4e565b90612e6d565b7fffffffff0000000000000000000000000000000000000000000000000000000081167fbae4e0fa0000000000000000000000000000000000000000000000000000000081141590816130d5575b816130aa575b5061305b5750612fe1612fd98484612e5c565b810190612ed3565b63ffffffff82168103613022575050600091829161300460405180938193612ee2565b039082305af1613012612ef0565b901561301b5750565b60203d9101fd5b7f8a6e4ce80000000000000000000000000000000000000000000000000000000060005263ffffffff9081166004521660245260446000fd5b7f12ba286f000000000000000000000000000000000000000000000000000000006000527fffffffff000000000000000000000000000000000000000000000000000000001660045260246000fd5b7f5a837f97000000000000000000000000000000000000000000000000000000009150141538612fc6565b7f7524051a000000000000000000000000000000000000000000000000000000008114159150612fc0565b7fac7a7efd0000000000000000000000000000000000000000000000000000000060005260046000fd5b61314a60019263ffffffff8093166000526007602052604060002061183b565b54161890565b30330361315957565b7f371a73280000000000000000000000000000000000000000000000000000000060005260046000fd5b6040516020810190600082526020815261319e604082610313565b51902090565b906131ae826119ba565b6131bb6040519182610313565b8281527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06131e982946119ba565b0190602036910137565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114611d015760010190565b90602061058a928181520190610556565b6020810167ffffffffffffffff613250825167ffffffffffffffff1690565b161561358657815161326181610510565b61326a81610510565b151580613568575b61353e57608082015180518015918215613528575b50506134fe5760a0820151805180159182156134e8575b50506134be576132bf611638611625611618845167ffffffffffffffff1690565b6134aa576132f961155460016132f16116576132e2611554604089015160ff1690565b955167ffffffffffffffff1690565b015460ff1690565b918183116134765760c0019182515191610100831161344c5761331b90611cf0565b8211156134225760009161332e816131a4565b9360005b8281106133905750505061334861334d91611cf0565b611de8565b9081811061336057505061036390613948565b7f548dd21f0000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b61339b818351611e10565b51604081015151613412575b602081015151158015613409575b6133cf5790600191516133c88289611e10565b5201613332565b613405906040519182917f9fa4031400000000000000000000000000000000000000000000000000000000835260048301613220565b0390fd5b508051156133b5565b9461341c906131f3565b946133a7565b7f4856694e0000000000000000000000000000000000000000000000000000000060005260046000fd5b7f1b925da60000000000000000000000000000000000000000000000000000000060005260046000fd5b507f2db220400000000000000000000000000000000000000000000000000000000060005260049190915260245260446000fd5b51610b2b9067ffffffffffffffff166116f7565b7fdee985740000000000000000000000000000000000000000000000000000000060005260046000fd5b6020012090506134f6613183565b14388061329e565b7f358c19270000000000000000000000000000000000000000000000000000000060005260046000fd5b602001209050613536613183565b143880613287565b7f3302dbd70000000000000000000000000000000000000000000000000000000060005260046000fd5b506001825161357681610510565b61357f81610510565b1415613272565b7f698cf8e00000000000000000000000000000000000000000000000000000000060005260046000fd5b90613677929361360663ffffffff9283604051957f45564d0000000000000000000000000000000000000000000000000000000000602088015246604088015230606088015216608086015260a0850190610549565b1660c082015260c0815261361b60e082610313565b6020604051938261363586945180928580880191016103af565b8301613649825180938580850191016103af565b0101037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101835282610313565b602081519101207fffff00000000000000000000000000000000000000000000000000000000000019167e0a0000000000000000000000000000000000000000000000000000000000001790565b73ffffffffffffffffffffffffffffffffffffffff6001541633036136e657565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b519061036382610496565b9080601f83011215610216578151613732816119ba565b926137406040519485610313565b81845260208085019260051b82010192831161021657602001905b8282106137685750505090565b815181526020918201910161375b565b9080601f8301121561021657815161378f816119ba565b9261379d6040519485610313565b81845260208085019260051b82010192831161021657602001905b8282106137c55750505090565b81518152602091820191016137b8565b6020818303126102165780519067ffffffffffffffff821161021657019080601f830112156102165781519161380a836119ba565b926138186040519485610313565b80845260208085019160051b830101918383116102165760208101915b83831061384457505050505090565b825167ffffffffffffffff8111610216578201906101007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083880301126102165761388d610365565b9061389a60208401613710565b82526138a860408401613710565b60208301526138b960608401613710565b60408301526080830151606083015260a0830151608083015260c083015160a083015260e083015167ffffffffffffffff8111610216578760206138ff9286010161371b565b60c08301526101008301519167ffffffffffffffff83116102165761392c88602080969581960101613778565b60e0820152815201920191613835565b6040513d6000823e3d90fd5b80516139515750565b60405180917f05a519660000000000000000000000000000000000000000000000000000000082526024820160206004840152815180915260206044840192019060005b818110613a10575050509080600092038173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa8015613a0b576139ec5750565b613a08903d806000833e613a008183610313565b8101906137d5565b50565b61393c565b8251845285945060209384019390920191600101613995565b80548210156118965760005260206000200190600090565b6003548110156118965760036000527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b015490565b80548015613ada577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190613aab8282613a29565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b1916905555565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b600081815260046020526040902054908115613be6577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820190828211611d0157600354927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8401938411611d01578383600095613ba59503613bab575b505050613b946003613a76565b600490600052602052604060002090565b55600190565b613b94613bd791613bcd613bc3613bdd956003613a29565b90549060031b1c90565b9283916003613a29565b90611bb1565b55388080613b87565b5050600090565b600081815260046020526040902054613c7857600354680100000000000000008110156102d557613c5f613c2a8260018594016003556003613a29565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b9055600354906000526004602052604060002055600190565b5060009056fea164736f6c634300081a000a",
}

var CCIPHomeABI = CCIPHomeMetaData.ABI

var CCIPHomeBin = CCIPHomeMetaData.Bin

func DeployCCIPHome(auth *bind.TransactOpts, backend bind.ContractBackend, capabilitiesRegistry common.Address) (common.Address, *types.Transaction, *CCIPHome, error) {
	parsed, err := CCIPHomeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CCIPHomeBin), backend, capabilitiesRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CCIPHome{address: address, abi: *parsed, CCIPHomeCaller: CCIPHomeCaller{contract: contract}, CCIPHomeTransactor: CCIPHomeTransactor{contract: contract}, CCIPHomeFilterer: CCIPHomeFilterer{contract: contract}}, nil
}

type CCIPHome struct {
	address common.Address
	abi     abi.ABI
	CCIPHomeCaller
	CCIPHomeTransactor
	CCIPHomeFilterer
}

type CCIPHomeCaller struct {
	contract *bind.BoundContract
}

type CCIPHomeTransactor struct {
	contract *bind.BoundContract
}

type CCIPHomeFilterer struct {
	contract *bind.BoundContract
}

type CCIPHomeSession struct {
	Contract     *CCIPHome
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CCIPHomeCallerSession struct {
	Contract *CCIPHomeCaller
	CallOpts bind.CallOpts
}

type CCIPHomeTransactorSession struct {
	Contract     *CCIPHomeTransactor
	TransactOpts bind.TransactOpts
}

type CCIPHomeRaw struct {
	Contract *CCIPHome
}

type CCIPHomeCallerRaw struct {
	Contract *CCIPHomeCaller
}

type CCIPHomeTransactorRaw struct {
	Contract *CCIPHomeTransactor
}

func NewCCIPHome(address common.Address, backend bind.ContractBackend) (*CCIPHome, error) {
	abi, err := abi.JSON(strings.NewReader(CCIPHomeABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCCIPHome(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CCIPHome{address: address, abi: abi, CCIPHomeCaller: CCIPHomeCaller{contract: contract}, CCIPHomeTransactor: CCIPHomeTransactor{contract: contract}, CCIPHomeFilterer: CCIPHomeFilterer{contract: contract}}, nil
}

func NewCCIPHomeCaller(address common.Address, caller bind.ContractCaller) (*CCIPHomeCaller, error) {
	contract, err := bindCCIPHome(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CCIPHomeCaller{contract: contract}, nil
}

func NewCCIPHomeTransactor(address common.Address, transactor bind.ContractTransactor) (*CCIPHomeTransactor, error) {
	contract, err := bindCCIPHome(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CCIPHomeTransactor{contract: contract}, nil
}

func NewCCIPHomeFilterer(address common.Address, filterer bind.ContractFilterer) (*CCIPHomeFilterer, error) {
	contract, err := bindCCIPHome(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CCIPHomeFilterer{contract: contract}, nil
}

func bindCCIPHome(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CCIPHomeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_CCIPHome *CCIPHomeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CCIPHome.Contract.CCIPHomeCaller.contract.Call(opts, result, method, params...)
}

func (_CCIPHome *CCIPHomeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPHome.Contract.CCIPHomeTransactor.contract.Transfer(opts)
}

func (_CCIPHome *CCIPHomeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CCIPHome.Contract.CCIPHomeTransactor.contract.Transact(opts, method, params...)
}

func (_CCIPHome *CCIPHomeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CCIPHome.Contract.contract.Call(opts, result, method, params...)
}

func (_CCIPHome *CCIPHomeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPHome.Contract.contract.Transfer(opts)
}

func (_CCIPHome *CCIPHomeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CCIPHome.Contract.contract.Transact(opts, method, params...)
}

func (_CCIPHome *CCIPHomeCaller) GetActiveDigest(opts *bind.CallOpts, donId uint32, pluginType uint8) ([32]byte, error) {
	var out []interface{}
	err := _CCIPHome.contract.Call(opts, &out, "getActiveDigest", donId, pluginType)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_CCIPHome *CCIPHomeSession) GetActiveDigest(donId uint32, pluginType uint8) ([32]byte, error) {
	return _CCIPHome.Contract.GetActiveDigest(&_CCIPHome.CallOpts, donId, pluginType)
}

func (_CCIPHome *CCIPHomeCallerSession) GetActiveDigest(donId uint32, pluginType uint8) ([32]byte, error) {
	return _CCIPHome.Contract.GetActiveDigest(&_CCIPHome.CallOpts, donId, pluginType)
}

func (_CCIPHome *CCIPHomeCaller) GetAllChainConfigs(opts *bind.CallOpts, pageIndex *big.Int, pageSize *big.Int) ([]CCIPHomeChainConfigArgs, error) {
	var out []interface{}
	err := _CCIPHome.contract.Call(opts, &out, "getAllChainConfigs", pageIndex, pageSize)

	if err != nil {
		return *new([]CCIPHomeChainConfigArgs), err
	}

	out0 := *abi.ConvertType(out[0], new([]CCIPHomeChainConfigArgs)).(*[]CCIPHomeChainConfigArgs)

	return out0, err

}

func (_CCIPHome *CCIPHomeSession) GetAllChainConfigs(pageIndex *big.Int, pageSize *big.Int) ([]CCIPHomeChainConfigArgs, error) {
	return _CCIPHome.Contract.GetAllChainConfigs(&_CCIPHome.CallOpts, pageIndex, pageSize)
}

func (_CCIPHome *CCIPHomeCallerSession) GetAllChainConfigs(pageIndex *big.Int, pageSize *big.Int) ([]CCIPHomeChainConfigArgs, error) {
	return _CCIPHome.Contract.GetAllChainConfigs(&_CCIPHome.CallOpts, pageIndex, pageSize)
}

func (_CCIPHome *CCIPHomeCaller) GetAllConfigs(opts *bind.CallOpts, donId uint32, pluginType uint8) (GetAllConfigs,

	error) {
	var out []interface{}
	err := _CCIPHome.contract.Call(opts, &out, "getAllConfigs", donId, pluginType)

	outstruct := new(GetAllConfigs)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ActiveConfig = *abi.ConvertType(out[0], new(CCIPHomeVersionedConfig)).(*CCIPHomeVersionedConfig)
	outstruct.CandidateConfig = *abi.ConvertType(out[1], new(CCIPHomeVersionedConfig)).(*CCIPHomeVersionedConfig)

	return *outstruct, err

}

func (_CCIPHome *CCIPHomeSession) GetAllConfigs(donId uint32, pluginType uint8) (GetAllConfigs,

	error) {
	return _CCIPHome.Contract.GetAllConfigs(&_CCIPHome.CallOpts, donId, pluginType)
}

func (_CCIPHome *CCIPHomeCallerSession) GetAllConfigs(donId uint32, pluginType uint8) (GetAllConfigs,

	error) {
	return _CCIPHome.Contract.GetAllConfigs(&_CCIPHome.CallOpts, donId, pluginType)
}

func (_CCIPHome *CCIPHomeCaller) GetCandidateDigest(opts *bind.CallOpts, donId uint32, pluginType uint8) ([32]byte, error) {
	var out []interface{}
	err := _CCIPHome.contract.Call(opts, &out, "getCandidateDigest", donId, pluginType)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_CCIPHome *CCIPHomeSession) GetCandidateDigest(donId uint32, pluginType uint8) ([32]byte, error) {
	return _CCIPHome.Contract.GetCandidateDigest(&_CCIPHome.CallOpts, donId, pluginType)
}

func (_CCIPHome *CCIPHomeCallerSession) GetCandidateDigest(donId uint32, pluginType uint8) ([32]byte, error) {
	return _CCIPHome.Contract.GetCandidateDigest(&_CCIPHome.CallOpts, donId, pluginType)
}

func (_CCIPHome *CCIPHomeCaller) GetCapabilityConfiguration(opts *bind.CallOpts, arg0 uint32) ([]byte, error) {
	var out []interface{}
	err := _CCIPHome.contract.Call(opts, &out, "getCapabilityConfiguration", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_CCIPHome *CCIPHomeSession) GetCapabilityConfiguration(arg0 uint32) ([]byte, error) {
	return _CCIPHome.Contract.GetCapabilityConfiguration(&_CCIPHome.CallOpts, arg0)
}

func (_CCIPHome *CCIPHomeCallerSession) GetCapabilityConfiguration(arg0 uint32) ([]byte, error) {
	return _CCIPHome.Contract.GetCapabilityConfiguration(&_CCIPHome.CallOpts, arg0)
}

func (_CCIPHome *CCIPHomeCaller) GetCapabilityRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CCIPHome.contract.Call(opts, &out, "getCapabilityRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CCIPHome *CCIPHomeSession) GetCapabilityRegistry() (common.Address, error) {
	return _CCIPHome.Contract.GetCapabilityRegistry(&_CCIPHome.CallOpts)
}

func (_CCIPHome *CCIPHomeCallerSession) GetCapabilityRegistry() (common.Address, error) {
	return _CCIPHome.Contract.GetCapabilityRegistry(&_CCIPHome.CallOpts)
}

func (_CCIPHome *CCIPHomeCaller) GetChainConfig(opts *bind.CallOpts, chainSelector uint64) (CCIPHomeChainConfig, error) {
	var out []interface{}
	err := _CCIPHome.contract.Call(opts, &out, "getChainConfig", chainSelector)

	if err != nil {
		return *new(CCIPHomeChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(CCIPHomeChainConfig)).(*CCIPHomeChainConfig)

	return out0, err

}

func (_CCIPHome *CCIPHomeSession) GetChainConfig(chainSelector uint64) (CCIPHomeChainConfig, error) {
	return _CCIPHome.Contract.GetChainConfig(&_CCIPHome.CallOpts, chainSelector)
}

func (_CCIPHome *CCIPHomeCallerSession) GetChainConfig(chainSelector uint64) (CCIPHomeChainConfig, error) {
	return _CCIPHome.Contract.GetChainConfig(&_CCIPHome.CallOpts, chainSelector)
}

func (_CCIPHome *CCIPHomeCaller) GetConfig(opts *bind.CallOpts, donId uint32, pluginType uint8, configDigest [32]byte) (GetConfig,

	error) {
	var out []interface{}
	err := _CCIPHome.contract.Call(opts, &out, "getConfig", donId, pluginType, configDigest)

	outstruct := new(GetConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.VersionedConfig = *abi.ConvertType(out[0], new(CCIPHomeVersionedConfig)).(*CCIPHomeVersionedConfig)
	outstruct.Ok = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

func (_CCIPHome *CCIPHomeSession) GetConfig(donId uint32, pluginType uint8, configDigest [32]byte) (GetConfig,

	error) {
	return _CCIPHome.Contract.GetConfig(&_CCIPHome.CallOpts, donId, pluginType, configDigest)
}

func (_CCIPHome *CCIPHomeCallerSession) GetConfig(donId uint32, pluginType uint8, configDigest [32]byte) (GetConfig,

	error) {
	return _CCIPHome.Contract.GetConfig(&_CCIPHome.CallOpts, donId, pluginType, configDigest)
}

func (_CCIPHome *CCIPHomeCaller) GetConfigDigests(opts *bind.CallOpts, donId uint32, pluginType uint8) (GetConfigDigests,

	error) {
	var out []interface{}
	err := _CCIPHome.contract.Call(opts, &out, "getConfigDigests", donId, pluginType)

	outstruct := new(GetConfigDigests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ActiveConfigDigest = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.CandidateConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_CCIPHome *CCIPHomeSession) GetConfigDigests(donId uint32, pluginType uint8) (GetConfigDigests,

	error) {
	return _CCIPHome.Contract.GetConfigDigests(&_CCIPHome.CallOpts, donId, pluginType)
}

func (_CCIPHome *CCIPHomeCallerSession) GetConfigDigests(donId uint32, pluginType uint8) (GetConfigDigests,

	error) {
	return _CCIPHome.Contract.GetConfigDigests(&_CCIPHome.CallOpts, donId, pluginType)
}

func (_CCIPHome *CCIPHomeCaller) GetNumChainConfigurations(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CCIPHome.contract.Call(opts, &out, "getNumChainConfigurations")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_CCIPHome *CCIPHomeSession) GetNumChainConfigurations() (*big.Int, error) {
	return _CCIPHome.Contract.GetNumChainConfigurations(&_CCIPHome.CallOpts)
}

func (_CCIPHome *CCIPHomeCallerSession) GetNumChainConfigurations() (*big.Int, error) {
	return _CCIPHome.Contract.GetNumChainConfigurations(&_CCIPHome.CallOpts)
}

func (_CCIPHome *CCIPHomeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CCIPHome.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CCIPHome *CCIPHomeSession) Owner() (common.Address, error) {
	return _CCIPHome.Contract.Owner(&_CCIPHome.CallOpts)
}

func (_CCIPHome *CCIPHomeCallerSession) Owner() (common.Address, error) {
	return _CCIPHome.Contract.Owner(&_CCIPHome.CallOpts)
}

func (_CCIPHome *CCIPHomeCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _CCIPHome.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_CCIPHome *CCIPHomeSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CCIPHome.Contract.SupportsInterface(&_CCIPHome.CallOpts, interfaceId)
}

func (_CCIPHome *CCIPHomeCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CCIPHome.Contract.SupportsInterface(&_CCIPHome.CallOpts, interfaceId)
}

func (_CCIPHome *CCIPHomeCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CCIPHome.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_CCIPHome *CCIPHomeSession) TypeAndVersion() (string, error) {
	return _CCIPHome.Contract.TypeAndVersion(&_CCIPHome.CallOpts)
}

func (_CCIPHome *CCIPHomeCallerSession) TypeAndVersion() (string, error) {
	return _CCIPHome.Contract.TypeAndVersion(&_CCIPHome.CallOpts)
}

func (_CCIPHome *CCIPHomeTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPHome.contract.Transact(opts, "acceptOwnership")
}

func (_CCIPHome *CCIPHomeSession) AcceptOwnership() (*types.Transaction, error) {
	return _CCIPHome.Contract.AcceptOwnership(&_CCIPHome.TransactOpts)
}

func (_CCIPHome *CCIPHomeTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _CCIPHome.Contract.AcceptOwnership(&_CCIPHome.TransactOpts)
}

func (_CCIPHome *CCIPHomeTransactor) ApplyChainConfigUpdates(opts *bind.TransactOpts, chainSelectorRemoves []uint64, chainConfigAdds []CCIPHomeChainConfigArgs) (*types.Transaction, error) {
	return _CCIPHome.contract.Transact(opts, "applyChainConfigUpdates", chainSelectorRemoves, chainConfigAdds)
}

func (_CCIPHome *CCIPHomeSession) ApplyChainConfigUpdates(chainSelectorRemoves []uint64, chainConfigAdds []CCIPHomeChainConfigArgs) (*types.Transaction, error) {
	return _CCIPHome.Contract.ApplyChainConfigUpdates(&_CCIPHome.TransactOpts, chainSelectorRemoves, chainConfigAdds)
}

func (_CCIPHome *CCIPHomeTransactorSession) ApplyChainConfigUpdates(chainSelectorRemoves []uint64, chainConfigAdds []CCIPHomeChainConfigArgs) (*types.Transaction, error) {
	return _CCIPHome.Contract.ApplyChainConfigUpdates(&_CCIPHome.TransactOpts, chainSelectorRemoves, chainConfigAdds)
}

func (_CCIPHome *CCIPHomeTransactor) BeforeCapabilityConfigSet(opts *bind.TransactOpts, arg0 [][32]byte, update []byte, arg2 uint64, donId uint32) (*types.Transaction, error) {
	return _CCIPHome.contract.Transact(opts, "beforeCapabilityConfigSet", arg0, update, arg2, donId)
}

func (_CCIPHome *CCIPHomeSession) BeforeCapabilityConfigSet(arg0 [][32]byte, update []byte, arg2 uint64, donId uint32) (*types.Transaction, error) {
	return _CCIPHome.Contract.BeforeCapabilityConfigSet(&_CCIPHome.TransactOpts, arg0, update, arg2, donId)
}

func (_CCIPHome *CCIPHomeTransactorSession) BeforeCapabilityConfigSet(arg0 [][32]byte, update []byte, arg2 uint64, donId uint32) (*types.Transaction, error) {
	return _CCIPHome.Contract.BeforeCapabilityConfigSet(&_CCIPHome.TransactOpts, arg0, update, arg2, donId)
}

func (_CCIPHome *CCIPHomeTransactor) PromoteCandidateAndRevokeActive(opts *bind.TransactOpts, donId uint32, pluginType uint8, digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error) {
	return _CCIPHome.contract.Transact(opts, "promoteCandidateAndRevokeActive", donId, pluginType, digestToPromote, digestToRevoke)
}

func (_CCIPHome *CCIPHomeSession) PromoteCandidateAndRevokeActive(donId uint32, pluginType uint8, digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error) {
	return _CCIPHome.Contract.PromoteCandidateAndRevokeActive(&_CCIPHome.TransactOpts, donId, pluginType, digestToPromote, digestToRevoke)
}

func (_CCIPHome *CCIPHomeTransactorSession) PromoteCandidateAndRevokeActive(donId uint32, pluginType uint8, digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error) {
	return _CCIPHome.Contract.PromoteCandidateAndRevokeActive(&_CCIPHome.TransactOpts, donId, pluginType, digestToPromote, digestToRevoke)
}

func (_CCIPHome *CCIPHomeTransactor) RevokeCandidate(opts *bind.TransactOpts, donId uint32, pluginType uint8, configDigest [32]byte) (*types.Transaction, error) {
	return _CCIPHome.contract.Transact(opts, "revokeCandidate", donId, pluginType, configDigest)
}

func (_CCIPHome *CCIPHomeSession) RevokeCandidate(donId uint32, pluginType uint8, configDigest [32]byte) (*types.Transaction, error) {
	return _CCIPHome.Contract.RevokeCandidate(&_CCIPHome.TransactOpts, donId, pluginType, configDigest)
}

func (_CCIPHome *CCIPHomeTransactorSession) RevokeCandidate(donId uint32, pluginType uint8, configDigest [32]byte) (*types.Transaction, error) {
	return _CCIPHome.Contract.RevokeCandidate(&_CCIPHome.TransactOpts, donId, pluginType, configDigest)
}

func (_CCIPHome *CCIPHomeTransactor) SetCandidate(opts *bind.TransactOpts, donId uint32, pluginType uint8, config CCIPHomeOCR3Config, digestToOverwrite [32]byte) (*types.Transaction, error) {
	return _CCIPHome.contract.Transact(opts, "setCandidate", donId, pluginType, config, digestToOverwrite)
}

func (_CCIPHome *CCIPHomeSession) SetCandidate(donId uint32, pluginType uint8, config CCIPHomeOCR3Config, digestToOverwrite [32]byte) (*types.Transaction, error) {
	return _CCIPHome.Contract.SetCandidate(&_CCIPHome.TransactOpts, donId, pluginType, config, digestToOverwrite)
}

func (_CCIPHome *CCIPHomeTransactorSession) SetCandidate(donId uint32, pluginType uint8, config CCIPHomeOCR3Config, digestToOverwrite [32]byte) (*types.Transaction, error) {
	return _CCIPHome.Contract.SetCandidate(&_CCIPHome.TransactOpts, donId, pluginType, config, digestToOverwrite)
}

func (_CCIPHome *CCIPHomeTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _CCIPHome.contract.Transact(opts, "transferOwnership", to)
}

func (_CCIPHome *CCIPHomeSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CCIPHome.Contract.TransferOwnership(&_CCIPHome.TransactOpts, to)
}

func (_CCIPHome *CCIPHomeTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CCIPHome.Contract.TransferOwnership(&_CCIPHome.TransactOpts, to)
}

type CCIPHomeActiveConfigRevokedIterator struct {
	Event *CCIPHomeActiveConfigRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPHomeActiveConfigRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPHomeActiveConfigRevoked)
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
		it.Event = new(CCIPHomeActiveConfigRevoked)
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

func (it *CCIPHomeActiveConfigRevokedIterator) Error() error {
	return it.fail
}

func (it *CCIPHomeActiveConfigRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPHomeActiveConfigRevoked struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_CCIPHome *CCIPHomeFilterer) FilterActiveConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*CCIPHomeActiveConfigRevokedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _CCIPHome.contract.FilterLogs(opts, "ActiveConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &CCIPHomeActiveConfigRevokedIterator{contract: _CCIPHome.contract, event: "ActiveConfigRevoked", logs: logs, sub: sub}, nil
}

func (_CCIPHome *CCIPHomeFilterer) WatchActiveConfigRevoked(opts *bind.WatchOpts, sink chan<- *CCIPHomeActiveConfigRevoked, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _CCIPHome.contract.WatchLogs(opts, "ActiveConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPHomeActiveConfigRevoked)
				if err := _CCIPHome.contract.UnpackLog(event, "ActiveConfigRevoked", log); err != nil {
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

func (_CCIPHome *CCIPHomeFilterer) ParseActiveConfigRevoked(log types.Log) (*CCIPHomeActiveConfigRevoked, error) {
	event := new(CCIPHomeActiveConfigRevoked)
	if err := _CCIPHome.contract.UnpackLog(event, "ActiveConfigRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPHomeCandidateConfigRevokedIterator struct {
	Event *CCIPHomeCandidateConfigRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPHomeCandidateConfigRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPHomeCandidateConfigRevoked)
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
		it.Event = new(CCIPHomeCandidateConfigRevoked)
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

func (it *CCIPHomeCandidateConfigRevokedIterator) Error() error {
	return it.fail
}

func (it *CCIPHomeCandidateConfigRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPHomeCandidateConfigRevoked struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_CCIPHome *CCIPHomeFilterer) FilterCandidateConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*CCIPHomeCandidateConfigRevokedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _CCIPHome.contract.FilterLogs(opts, "CandidateConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &CCIPHomeCandidateConfigRevokedIterator{contract: _CCIPHome.contract, event: "CandidateConfigRevoked", logs: logs, sub: sub}, nil
}

func (_CCIPHome *CCIPHomeFilterer) WatchCandidateConfigRevoked(opts *bind.WatchOpts, sink chan<- *CCIPHomeCandidateConfigRevoked, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _CCIPHome.contract.WatchLogs(opts, "CandidateConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPHomeCandidateConfigRevoked)
				if err := _CCIPHome.contract.UnpackLog(event, "CandidateConfigRevoked", log); err != nil {
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

func (_CCIPHome *CCIPHomeFilterer) ParseCandidateConfigRevoked(log types.Log) (*CCIPHomeCandidateConfigRevoked, error) {
	event := new(CCIPHomeCandidateConfigRevoked)
	if err := _CCIPHome.contract.UnpackLog(event, "CandidateConfigRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPHomeCapabilityConfigurationSetIterator struct {
	Event *CCIPHomeCapabilityConfigurationSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPHomeCapabilityConfigurationSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPHomeCapabilityConfigurationSet)
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
		it.Event = new(CCIPHomeCapabilityConfigurationSet)
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

func (it *CCIPHomeCapabilityConfigurationSetIterator) Error() error {
	return it.fail
}

func (it *CCIPHomeCapabilityConfigurationSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPHomeCapabilityConfigurationSet struct {
	Raw types.Log
}

func (_CCIPHome *CCIPHomeFilterer) FilterCapabilityConfigurationSet(opts *bind.FilterOpts) (*CCIPHomeCapabilityConfigurationSetIterator, error) {

	logs, sub, err := _CCIPHome.contract.FilterLogs(opts, "CapabilityConfigurationSet")
	if err != nil {
		return nil, err
	}
	return &CCIPHomeCapabilityConfigurationSetIterator{contract: _CCIPHome.contract, event: "CapabilityConfigurationSet", logs: logs, sub: sub}, nil
}

func (_CCIPHome *CCIPHomeFilterer) WatchCapabilityConfigurationSet(opts *bind.WatchOpts, sink chan<- *CCIPHomeCapabilityConfigurationSet) (event.Subscription, error) {

	logs, sub, err := _CCIPHome.contract.WatchLogs(opts, "CapabilityConfigurationSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPHomeCapabilityConfigurationSet)
				if err := _CCIPHome.contract.UnpackLog(event, "CapabilityConfigurationSet", log); err != nil {
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

func (_CCIPHome *CCIPHomeFilterer) ParseCapabilityConfigurationSet(log types.Log) (*CCIPHomeCapabilityConfigurationSet, error) {
	event := new(CCIPHomeCapabilityConfigurationSet)
	if err := _CCIPHome.contract.UnpackLog(event, "CapabilityConfigurationSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPHomeChainConfigRemovedIterator struct {
	Event *CCIPHomeChainConfigRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPHomeChainConfigRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPHomeChainConfigRemoved)
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
		it.Event = new(CCIPHomeChainConfigRemoved)
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

func (it *CCIPHomeChainConfigRemovedIterator) Error() error {
	return it.fail
}

func (it *CCIPHomeChainConfigRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPHomeChainConfigRemoved struct {
	ChainSelector uint64
	Raw           types.Log
}

func (_CCIPHome *CCIPHomeFilterer) FilterChainConfigRemoved(opts *bind.FilterOpts) (*CCIPHomeChainConfigRemovedIterator, error) {

	logs, sub, err := _CCIPHome.contract.FilterLogs(opts, "ChainConfigRemoved")
	if err != nil {
		return nil, err
	}
	return &CCIPHomeChainConfigRemovedIterator{contract: _CCIPHome.contract, event: "ChainConfigRemoved", logs: logs, sub: sub}, nil
}

func (_CCIPHome *CCIPHomeFilterer) WatchChainConfigRemoved(opts *bind.WatchOpts, sink chan<- *CCIPHomeChainConfigRemoved) (event.Subscription, error) {

	logs, sub, err := _CCIPHome.contract.WatchLogs(opts, "ChainConfigRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPHomeChainConfigRemoved)
				if err := _CCIPHome.contract.UnpackLog(event, "ChainConfigRemoved", log); err != nil {
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

func (_CCIPHome *CCIPHomeFilterer) ParseChainConfigRemoved(log types.Log) (*CCIPHomeChainConfigRemoved, error) {
	event := new(CCIPHomeChainConfigRemoved)
	if err := _CCIPHome.contract.UnpackLog(event, "ChainConfigRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPHomeChainConfigSetIterator struct {
	Event *CCIPHomeChainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPHomeChainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPHomeChainConfigSet)
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
		it.Event = new(CCIPHomeChainConfigSet)
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

func (it *CCIPHomeChainConfigSetIterator) Error() error {
	return it.fail
}

func (it *CCIPHomeChainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPHomeChainConfigSet struct {
	ChainSelector uint64
	ChainConfig   CCIPHomeChainConfig
	Raw           types.Log
}

func (_CCIPHome *CCIPHomeFilterer) FilterChainConfigSet(opts *bind.FilterOpts) (*CCIPHomeChainConfigSetIterator, error) {

	logs, sub, err := _CCIPHome.contract.FilterLogs(opts, "ChainConfigSet")
	if err != nil {
		return nil, err
	}
	return &CCIPHomeChainConfigSetIterator{contract: _CCIPHome.contract, event: "ChainConfigSet", logs: logs, sub: sub}, nil
}

func (_CCIPHome *CCIPHomeFilterer) WatchChainConfigSet(opts *bind.WatchOpts, sink chan<- *CCIPHomeChainConfigSet) (event.Subscription, error) {

	logs, sub, err := _CCIPHome.contract.WatchLogs(opts, "ChainConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPHomeChainConfigSet)
				if err := _CCIPHome.contract.UnpackLog(event, "ChainConfigSet", log); err != nil {
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

func (_CCIPHome *CCIPHomeFilterer) ParseChainConfigSet(log types.Log) (*CCIPHomeChainConfigSet, error) {
	event := new(CCIPHomeChainConfigSet)
	if err := _CCIPHome.contract.UnpackLog(event, "ChainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPHomeConfigPromotedIterator struct {
	Event *CCIPHomeConfigPromoted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPHomeConfigPromotedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPHomeConfigPromoted)
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
		it.Event = new(CCIPHomeConfigPromoted)
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

func (it *CCIPHomeConfigPromotedIterator) Error() error {
	return it.fail
}

func (it *CCIPHomeConfigPromotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPHomeConfigPromoted struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_CCIPHome *CCIPHomeFilterer) FilterConfigPromoted(opts *bind.FilterOpts, configDigest [][32]byte) (*CCIPHomeConfigPromotedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _CCIPHome.contract.FilterLogs(opts, "ConfigPromoted", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &CCIPHomeConfigPromotedIterator{contract: _CCIPHome.contract, event: "ConfigPromoted", logs: logs, sub: sub}, nil
}

func (_CCIPHome *CCIPHomeFilterer) WatchConfigPromoted(opts *bind.WatchOpts, sink chan<- *CCIPHomeConfigPromoted, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _CCIPHome.contract.WatchLogs(opts, "ConfigPromoted", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPHomeConfigPromoted)
				if err := _CCIPHome.contract.UnpackLog(event, "ConfigPromoted", log); err != nil {
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

func (_CCIPHome *CCIPHomeFilterer) ParseConfigPromoted(log types.Log) (*CCIPHomeConfigPromoted, error) {
	event := new(CCIPHomeConfigPromoted)
	if err := _CCIPHome.contract.UnpackLog(event, "ConfigPromoted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPHomeConfigSetIterator struct {
	Event *CCIPHomeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPHomeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPHomeConfigSet)
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
		it.Event = new(CCIPHomeConfigSet)
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

func (it *CCIPHomeConfigSetIterator) Error() error {
	return it.fail
}

func (it *CCIPHomeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPHomeConfigSet struct {
	ConfigDigest [32]byte
	Version      uint32
	Config       CCIPHomeOCR3Config
	Raw          types.Log
}

func (_CCIPHome *CCIPHomeFilterer) FilterConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*CCIPHomeConfigSetIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _CCIPHome.contract.FilterLogs(opts, "ConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &CCIPHomeConfigSetIterator{contract: _CCIPHome.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_CCIPHome *CCIPHomeFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *CCIPHomeConfigSet, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _CCIPHome.contract.WatchLogs(opts, "ConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPHomeConfigSet)
				if err := _CCIPHome.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_CCIPHome *CCIPHomeFilterer) ParseConfigSet(log types.Log) (*CCIPHomeConfigSet, error) {
	event := new(CCIPHomeConfigSet)
	if err := _CCIPHome.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPHomeOwnershipTransferRequestedIterator struct {
	Event *CCIPHomeOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPHomeOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPHomeOwnershipTransferRequested)
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
		it.Event = new(CCIPHomeOwnershipTransferRequested)
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

func (it *CCIPHomeOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *CCIPHomeOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPHomeOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CCIPHome *CCIPHomeFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CCIPHomeOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CCIPHome.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CCIPHomeOwnershipTransferRequestedIterator{contract: _CCIPHome.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_CCIPHome *CCIPHomeFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CCIPHomeOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CCIPHome.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPHomeOwnershipTransferRequested)
				if err := _CCIPHome.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_CCIPHome *CCIPHomeFilterer) ParseOwnershipTransferRequested(log types.Log) (*CCIPHomeOwnershipTransferRequested, error) {
	event := new(CCIPHomeOwnershipTransferRequested)
	if err := _CCIPHome.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPHomeOwnershipTransferredIterator struct {
	Event *CCIPHomeOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPHomeOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPHomeOwnershipTransferred)
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
		it.Event = new(CCIPHomeOwnershipTransferred)
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

func (it *CCIPHomeOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *CCIPHomeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPHomeOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CCIPHome *CCIPHomeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CCIPHomeOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CCIPHome.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CCIPHomeOwnershipTransferredIterator{contract: _CCIPHome.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_CCIPHome *CCIPHomeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CCIPHomeOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CCIPHome.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPHomeOwnershipTransferred)
				if err := _CCIPHome.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_CCIPHome *CCIPHomeFilterer) ParseOwnershipTransferred(log types.Log) (*CCIPHomeOwnershipTransferred, error) {
	event := new(CCIPHomeOwnershipTransferred)
	if err := _CCIPHome.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetAllConfigs struct {
	ActiveConfig    CCIPHomeVersionedConfig
	CandidateConfig CCIPHomeVersionedConfig
}
type GetConfig struct {
	VersionedConfig CCIPHomeVersionedConfig
	Ok              bool
}
type GetConfigDigests struct {
	ActiveConfigDigest    [32]byte
	CandidateConfigDigest [32]byte
}

func (_CCIPHome *CCIPHome) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _CCIPHome.abi.Events["ActiveConfigRevoked"].ID:
		return _CCIPHome.ParseActiveConfigRevoked(log)
	case _CCIPHome.abi.Events["CandidateConfigRevoked"].ID:
		return _CCIPHome.ParseCandidateConfigRevoked(log)
	case _CCIPHome.abi.Events["CapabilityConfigurationSet"].ID:
		return _CCIPHome.ParseCapabilityConfigurationSet(log)
	case _CCIPHome.abi.Events["ChainConfigRemoved"].ID:
		return _CCIPHome.ParseChainConfigRemoved(log)
	case _CCIPHome.abi.Events["ChainConfigSet"].ID:
		return _CCIPHome.ParseChainConfigSet(log)
	case _CCIPHome.abi.Events["ConfigPromoted"].ID:
		return _CCIPHome.ParseConfigPromoted(log)
	case _CCIPHome.abi.Events["ConfigSet"].ID:
		return _CCIPHome.ParseConfigSet(log)
	case _CCIPHome.abi.Events["OwnershipTransferRequested"].ID:
		return _CCIPHome.ParseOwnershipTransferRequested(log)
	case _CCIPHome.abi.Events["OwnershipTransferred"].ID:
		return _CCIPHome.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CCIPHomeActiveConfigRevoked) Topic() common.Hash {
	return common.HexToHash("0x0b31c0055e2d464bef7781994b98c4ff9ef4ae0d05f59feb6a68c42de5e201b8")
}

func (CCIPHomeCandidateConfigRevoked) Topic() common.Hash {
	return common.HexToHash("0x53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b")
}

func (CCIPHomeCapabilityConfigurationSet) Topic() common.Hash {
	return common.HexToHash("0x84ad7751b744c9e2ee77da1d902b428aec7f0a343d67a24bbe2142e6f58a8d0f")
}

func (CCIPHomeChainConfigRemoved) Topic() common.Hash {
	return common.HexToHash("0x2a680691fef3b2d105196805935232c661ce703e92d464ef0b94a7bc62d714f0")
}

func (CCIPHomeChainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x05dd57854af2c291a94ea52e7c43d80bc3be7fa73022f98b735dea86642fa5e0")
}

func (CCIPHomeConfigPromoted) Topic() common.Hash {
	return common.HexToHash("0xfc3e98dbbd47c3fa7c1c05b6ec711caeaf70eca4554192b9ada8fc11a37f298e")
}

func (CCIPHomeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x94f085b7c57ec2a270befd0b7b2ec7452580040edee8bb0fb04609c81f0359c6")
}

func (CCIPHomeOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (CCIPHomeOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_CCIPHome *CCIPHome) Address() common.Address {
	return _CCIPHome.address
}

type CCIPHomeInterface interface {
	GetActiveDigest(opts *bind.CallOpts, donId uint32, pluginType uint8) ([32]byte, error)

	GetAllChainConfigs(opts *bind.CallOpts, pageIndex *big.Int, pageSize *big.Int) ([]CCIPHomeChainConfigArgs, error)

	GetAllConfigs(opts *bind.CallOpts, donId uint32, pluginType uint8) (GetAllConfigs,

		error)

	GetCandidateDigest(opts *bind.CallOpts, donId uint32, pluginType uint8) ([32]byte, error)

	GetCapabilityConfiguration(opts *bind.CallOpts, arg0 uint32) ([]byte, error)

	GetCapabilityRegistry(opts *bind.CallOpts) (common.Address, error)

	GetChainConfig(opts *bind.CallOpts, chainSelector uint64) (CCIPHomeChainConfig, error)

	GetConfig(opts *bind.CallOpts, donId uint32, pluginType uint8, configDigest [32]byte) (GetConfig,

		error)

	GetConfigDigests(opts *bind.CallOpts, donId uint32, pluginType uint8) (GetConfigDigests,

		error)

	GetNumChainConfigurations(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyChainConfigUpdates(opts *bind.TransactOpts, chainSelectorRemoves []uint64, chainConfigAdds []CCIPHomeChainConfigArgs) (*types.Transaction, error)

	BeforeCapabilityConfigSet(opts *bind.TransactOpts, arg0 [][32]byte, update []byte, arg2 uint64, donId uint32) (*types.Transaction, error)

	PromoteCandidateAndRevokeActive(opts *bind.TransactOpts, donId uint32, pluginType uint8, digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error)

	RevokeCandidate(opts *bind.TransactOpts, donId uint32, pluginType uint8, configDigest [32]byte) (*types.Transaction, error)

	SetCandidate(opts *bind.TransactOpts, donId uint32, pluginType uint8, config CCIPHomeOCR3Config, digestToOverwrite [32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterActiveConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*CCIPHomeActiveConfigRevokedIterator, error)

	WatchActiveConfigRevoked(opts *bind.WatchOpts, sink chan<- *CCIPHomeActiveConfigRevoked, configDigest [][32]byte) (event.Subscription, error)

	ParseActiveConfigRevoked(log types.Log) (*CCIPHomeActiveConfigRevoked, error)

	FilterCandidateConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*CCIPHomeCandidateConfigRevokedIterator, error)

	WatchCandidateConfigRevoked(opts *bind.WatchOpts, sink chan<- *CCIPHomeCandidateConfigRevoked, configDigest [][32]byte) (event.Subscription, error)

	ParseCandidateConfigRevoked(log types.Log) (*CCIPHomeCandidateConfigRevoked, error)

	FilterCapabilityConfigurationSet(opts *bind.FilterOpts) (*CCIPHomeCapabilityConfigurationSetIterator, error)

	WatchCapabilityConfigurationSet(opts *bind.WatchOpts, sink chan<- *CCIPHomeCapabilityConfigurationSet) (event.Subscription, error)

	ParseCapabilityConfigurationSet(log types.Log) (*CCIPHomeCapabilityConfigurationSet, error)

	FilterChainConfigRemoved(opts *bind.FilterOpts) (*CCIPHomeChainConfigRemovedIterator, error)

	WatchChainConfigRemoved(opts *bind.WatchOpts, sink chan<- *CCIPHomeChainConfigRemoved) (event.Subscription, error)

	ParseChainConfigRemoved(log types.Log) (*CCIPHomeChainConfigRemoved, error)

	FilterChainConfigSet(opts *bind.FilterOpts) (*CCIPHomeChainConfigSetIterator, error)

	WatchChainConfigSet(opts *bind.WatchOpts, sink chan<- *CCIPHomeChainConfigSet) (event.Subscription, error)

	ParseChainConfigSet(log types.Log) (*CCIPHomeChainConfigSet, error)

	FilterConfigPromoted(opts *bind.FilterOpts, configDigest [][32]byte) (*CCIPHomeConfigPromotedIterator, error)

	WatchConfigPromoted(opts *bind.WatchOpts, sink chan<- *CCIPHomeConfigPromoted, configDigest [][32]byte) (event.Subscription, error)

	ParseConfigPromoted(log types.Log) (*CCIPHomeConfigPromoted, error)

	FilterConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*CCIPHomeConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *CCIPHomeConfigSet, configDigest [][32]byte) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*CCIPHomeConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CCIPHomeOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CCIPHomeOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*CCIPHomeOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CCIPHomeOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CCIPHomeOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*CCIPHomeOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
