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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"capabilitiesRegistry\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CanOnlySelfCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainSelectorNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ChainSelectorNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expectedConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"gotConfigDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"callDonId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"capabilityRegistryDonId\",\"type\":\"uint32\"}],\"name\":\"DONIdMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FChainMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fChain\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"FRoleDON\",\"type\":\"uint256\"}],\"name\":\"FChainTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FTooHigh\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signerKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"transmitterKey\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPHome.OCR3Node\",\"name\":\"node\",\"type\":\"tuple\"}],\"name\":\"InvalidNode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPluginType\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"}],\"name\":\"InvalidSelector\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoOpStateTransitionNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minimum\",\"type\":\"uint256\"}],\"name\":\"NotEnoughTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OfframpAddressCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCapabilitiesRegistryCanCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RMNHomeAddressCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RevokingZeroDigestNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManySigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ActiveConfigRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"CandidateConfigRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CapabilityConfigurationSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainConfigRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes32[]\",\"name\":\"readers\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"fChain\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structCCIPHome.ChainConfig\",\"name\":\"chainConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigPromoted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"FRoleDON\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offrampAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rmnHomeAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signerKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"transmitterKey\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPHome.OCR3Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structCCIPHome.OCR3Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"chainSelectorRemoves\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes32[]\",\"name\":\"readers\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"fChain\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPHome.ChainConfig\",\"name\":\"chainConfig\",\"type\":\"tuple\"}],\"internalType\":\"structCCIPHome.ChainConfigArgs[]\",\"name\":\"chainConfigAdds\",\"type\":\"tuple[]\"}],\"name\":\"applyChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"update\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"beforeCapabilityConfigSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"}],\"name\":\"getActiveDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"pageIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pageSize\",\"type\":\"uint256\"}],\"name\":\"getAllChainConfigs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes32[]\",\"name\":\"readers\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"fChain\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPHome.ChainConfig\",\"name\":\"chainConfig\",\"type\":\"tuple\"}],\"internalType\":\"structCCIPHome.ChainConfigArgs[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"}],\"name\":\"getAllConfigs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"FRoleDON\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offrampAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rmnHomeAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signerKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"transmitterKey\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPHome.OCR3Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPHome.OCR3Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"internalType\":\"structCCIPHome.VersionedConfig\",\"name\":\"activeConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"FRoleDON\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offrampAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rmnHomeAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signerKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"transmitterKey\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPHome.OCR3Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPHome.OCR3Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"internalType\":\"structCCIPHome.VersionedConfig\",\"name\":\"candidateConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"}],\"name\":\"getCandidateDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"name\":\"getCapabilityConfiguration\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"configuration\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCapabilityRegistry\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"getChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32[]\",\"name\":\"readers\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"fChain\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPHome.ChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"FRoleDON\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offrampAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rmnHomeAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signerKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"transmitterKey\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPHome.OCR3Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPHome.OCR3Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"internalType\":\"structCCIPHome.VersionedConfig\",\"name\":\"versionedConfig\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"ok\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"}],\"name\":\"getConfigDigests\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"activeConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"candidateConfigDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNumChainConfigurations\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"digestToPromote\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"digestToRevoke\",\"type\":\"bytes32\"}],\"name\":\"promoteCandidateAndRevokeActive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"revokeCandidate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"enumInternal.OCRPluginType\",\"name\":\"pluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"FRoleDON\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offrampAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rmnHomeAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"p2pId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signerKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"transmitterKey\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPHome.OCR3Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structCCIPHome.OCR3Config\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"digestToOverwrite\",\"type\":\"bytes32\"}],\"name\":\"setCandidate\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"newConfigDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a03460bf57601f61380238819003918201601f19168301916001600160401b0383118484101760c45780849260209460405283398101031260bf57516001600160a01b03811680820360bf57331560ae57600180546001600160a01b031916331790556006805463ffffffff1916905515609d5760805260405161372790816100db823960805181818161022a01528181612ae301526134ef0152f35b6342bcdf7f60e11b60005260046000fd5b639b15e16f60e01b60005260046000fd5b600080fd5b634e487b7160e01b600052604160045260246000fdfe6080604052600436101561001257600080fd5b60003560e01c806301ffc9a714610157578063020330e614610152578063181f5a771461014d57806333d9704a146101485780633df45a72146101435780634851d5491461013e5780635a837f97146101395780635f1edd9c146101345780637524051a1461012f57806379ba50971461012a5780637ac0d41e146101255780638318ed5d146101205780638da5cb5b1461011b578063922ea40614610116578063b149092b14610111578063b74b23561461010c578063bae4e0fa14610107578063f2fde38b14610102578063f442c89a146100fd5763fba64a7c146100f857600080fd5b6114b7565b61121a565b611114565b610f05565b610e48565b610dec565b610d3e565b610d0a565b610cc7565b610ca9565b610bde565b610ae3565b610a8f565b610826565b6107a0565b6106c0565b61063c565b61037f565b6101fd565b346101f85760206003193601126101f8576004357fffffffff0000000000000000000000000000000000000000000000000000000081168091036101f857807f78bea72100000000000000000000000000000000000000000000000000000000602092149081156101ce575b506040519015158152f35b7f01ffc9a700000000000000000000000000000000000000000000000000000000915014386101c3565b600080fd5b346101f85760006003193601126101f857602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6060810190811067ffffffffffffffff82111761029957604052565b61024e565b610100810190811067ffffffffffffffff82111761029957604052565b6040810190811067ffffffffffffffff82111761029957604052565b90601f601f19910116810190811067ffffffffffffffff82111761029957604052565b604051906103096040836102d7565b565b60405190610309610100836102d7565b67ffffffffffffffff811161029957601f01601f191660200190565b60005b83811061034a5750506000910152565b818101518382015260200161033a565b90601f19601f60209361037881518092818752878088019101610337565b0116010190565b346101f85760006003193601126101f8576103de60408051906103a281836102d7565b601282527f43434950486f6d6520312e362e302d646576000000000000000000000000000060208301525191829160208352602083019061035a565b0390f35b63ffffffff8116036101f857565b60643590610309826103e2565b600211156101f857565b3590610309826103fd565b60031960609101126101f85760043561042a816103e2565b90602435610437816103fd565b9060443590565b6002111561044857565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9060028210156104485752565b6104b8918151815260406104a7602084015160606020850152606084019061035a565b92015190604081840391015261035a565b90565b9080602083519182815201916020808360051b8301019401926000915b8383106104e757505050505090565b909192939460208061050583601f1986600196030187528951610484565b970193019301919392906104d8565b9060406104b89263ffffffff81511683526020810151602084015201519060606040820152610547606082018351610477565b602082015167ffffffffffffffff166080820152604082015160ff1660a0820152606082015167ffffffffffffffff1660c082015260e06106086105d361059e60808601516101008587015261016086019061035a565b60a08601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08683030161010087015261035a565b60c08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0858303016101208601526104bb565b920151906101407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08285030191015261035a565b346101f85761066a61065661065036610412565b9161188f565b604051928392604084526040840190610514565b90151560208301520390f35b60031960409101126101f85760043561068e816103e2565b906024356104b8816103fd565b90916106b26104b893604084526040840190610514565b916020818403910152610514565b346101f8576106ce36610676565b906106d7611548565b906106e0611548565b9261072f61072963ffffffff84168060005260056020526107058460406000206115a7565b90600052600760205263ffffffff6107218560406000206115a7565b5416906115ee565b506117be565b6020810151610796575b508161076e82610768610729946107636107749763ffffffff166000526005602052604060002090565b6115a7565b92612cd4565b906115ee565b602081015161078e575b506103de6040519283928361069b565b91503861077e565b9250610774610739565b346101f85761081160016107b336610676565b929061076e63ffffffff821694856000526005602052846107f76107db8360406000206115a7565b88600052600760205263ffffffff6107218560406000206115a7565b5001549560005260056020526107688160406000206115a7565b50015460408051928352602083019190915290f35b346101f85760806003193601126101f857600435610843816103e2565b60243590610850826103fd565b6044359160643591610860612cfa565b831580610a87575b610a5d576108826108798383612cd4565b63ffffffff1690565b8460016108ac836108a7876107638863ffffffff166000526005602052604060002090565b6115ee565b50015403610a03575060016109036108d8846107638563ffffffff166000526005602052604060002090565b61076e6108f9866107638763ffffffff166000526007602052604060002090565b5463ffffffff1690565b500180548481036109cd575091610763610934926000610974955563ffffffff166000526007602052604060002090565b6001610944825463ffffffff1690565b1863ffffffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000825416179055565b806109a2575b507ffc3e98dbbd47c3fa7c1c05b6ec711caeaf70eca4554192b9ada8fc11a37f298e600080a2005b7f0b31c0055e2d464bef7781994b98c4ff9ef4ae0d05f59feb6a68c42de5e201b8600080a23861097a565b7f93df584c00000000000000000000000000000000000000000000000000000000600052600452602484905260446000fd5b6000fd5b610a2b6001916108a76109ff95610763899663ffffffff166000526005602052604060002090565b5001547f93df584c00000000000000000000000000000000000000000000000000000000600052600452602452604490565b7f7b4d1e4f0000000000000000000000000000000000000000000000000000000060005260046000fd5b508215610868565b346101f85760206001610ad863ffffffff80610721610aad36610676565b9316928360005260058752610ac68160406000206115a7565b936000526007875260406000206115a7565b500154604051908152f35b346101f857610af136610412565b91610afa612cfa565b8215610bb45763ffffffff610b0f8383612cd4565b169263ffffffff82166000526005602052806001610b35866108a78760406000206115a7565b50015403610b8d57926108a7600193610763610b88946000977f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b8980a263ffffffff166000526005602052604060002090565b500155005b6001610a2b856108a7866107636109ff9763ffffffff166000526005602052604060002090565b7f0849d8cc0000000000000000000000000000000000000000000000000000000060005260046000fd5b346101f85760006003193601126101f85760005473ffffffffffffffffffffffffffffffffffffffff81163303610c7f577fffffffffffffffffffffffff00000000000000000000000000000000000000006001549133828416176001551660005573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b346101f85760006003193601126101f8576020600354604051908152f35b346101f85760206003193601126101f857610ce36004356103e2565b6040516020610cf281836102d7565b600082526103de60405192828493845283019061035a565b346101f85760006003193601126101f857602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b346101f8576020610d57610d5136610676565b9061193b565b604051908152f35b67ffffffffffffffff8116036101f857565b6044359061030982610d5f565b359061030982610d5f565b9190606081019083519160608252825180915260206080830193019060005b818110610dd65750505060408460ff60206104b896970151166020840152015190604081840391015261035a565b8251855260209485019490920191600101610da8565b346101f85760206003193601126101f85767ffffffffffffffff600435610e1281610d5f565b610e1a611968565b501660005260026020526103de610e346040600020611988565b604051918291602083526020830190610d89565b346101f85760406003193601126101f857610e67602435600435611b54565b6040518091602082016020835281518091526040830190602060408260051b8601019301916000905b828210610e9f57505050500390f35b91936020610ef5827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc060019597998495030186526040838a5167ffffffffffffffff815116845201519181858201520190610d89565b9601920192018594939192610e90565b346101f85760806003193601126101f857600435610f22816103e2565b602435610f2e816103fd565b60443567ffffffffffffffff81116101f8576101006003198260040192360301126101f85760643592610f5f612cfa565b610f71610f6c3684611d81565b612d9f565b610f7b838261193b565b938085036110e257917f94f085b7c57ec2a270befd0b7b2ec7452580040edee8bb0fb04609c81f0359c69161076e94936103de966110b7575b50611095826002611059610fd5610fd060065463ffffffff1690565b611e6a565b9461100b8663ffffffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000006006541617600655565b611038866040516110318161102389602083016120ee565b03601f1981018352826102d7565b8b8461311e565b99896107688c9b6107638563ffffffff166000526005602052604060002090565b506001810188905580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff86161781550161253a565b6110a4604051928392836126c2565b0390a26040519081529081906020820190565b7f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b600080a238610fb4565b7f93df584c00000000000000000000000000000000000000000000000000000000600052600485905260245260446000fd5b346101f85760206003193601126101f85760043573ffffffffffffffffffffffffffffffffffffffff81168091036101f85761114e613215565b3381146111bf57807fffffffffffffffffffffffff0000000000000000000000000000000000000000600054161760005573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b9181601f840112156101f85782359167ffffffffffffffff83116101f8576020808501948460051b0101116101f857565b346101f85760406003193601126101f85760043567ffffffffffffffff81116101f85761124b9036906004016111e9565b60243567ffffffffffffffff81116101f85761126b9036906004016111e9565b919092611276613215565b60005b82811061138b5750505060005b81811061128f57005b6112af6112aa6112a083858761277e565b60208101906122de565b6127be565b906112c36112be82858761277e565b6120ff565b6112cd835161347a565b6112e46112de602085015160ff1690565b60ff1690565b156113615782816113356001956113307f05dd57854af2c291a94ea52e7c43d80bc3be7fa73022f98b735dea86642fa5e09567ffffffffffffffff166000526002602052604060002090565b612944565b61134867ffffffffffffffff82166136a7565b50611358604051928392836129d7565b0390a101611286565b7fa9b3766e0000000000000000000000000000000000000000000000000000000060005260046000fd5b6113c66113c26113af6113a26112be8588886126df565b67ffffffffffffffff1690565b6000526004602052604060002054151590565b1590565b61147057806114006113fb6113e16112be60019588886126df565b67ffffffffffffffff166000526002602052604060002090565b612737565b6114196114146113a26112be8488886126df565b6135ff565b507f2a680691fef3b2d105196805935232c661ce703e92d464ef0b94a7bc62d714f061146761144c6112be8488886126df565b60405167ffffffffffffffff90911681529081906020820190565b0390a101611279565b6112be906109ff93611481936126df565b7f1bd4d2d20000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff16600452602490565b346101f85760806003193601126101f85760043567ffffffffffffffff81116101f8576114e89036906004016111e9565b505060243567ffffffffffffffff81116101f857366023820112156101f85780600401359067ffffffffffffffff82116101f85736602483830101116101f85761154691611534610d71565b50602461153f6103f0565b9201612aca565b005b604051906115558261027d565b816000815260006020820152604080519161156f8361029e565b60008352600060208401526000828401526000606084015260606080840152606060a0840152606060c0840152606060e08401520152565b90600281101561044857600052602052604060002090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b906002811015611602576007020190600090565b6115bf565b60028210156104485752565b90600182811c9216801561165c575b602083101461162d57565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f1691611622565b906040519182600082549261167a84611613565b80845293600181169081156116e6575060011461169f575b50610309925003836102d7565b90506000929192526020600020906000915b8183106116ca5750509060206103099282010138611692565b60209193508060019154838589010152019101909184926116b1565b602093506103099592507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b82010138611692565b67ffffffffffffffff81116102995760051b60200190565b90815461174a81611726565b9261175860405194856102d7565b818452602084019060005260206000206000915b8383106117795750505050565b6003602060019260405161178c8161027d565b8554815261179b858701611666565b838201526117ab60028701611666565b604082015281520192019201919061176c565b90604051916117cc8361027d565b60408363ffffffff83541681526001830154602082015261188660068351946117f48661029e565b61184d61183c600283015461180c60ff82168a611607565b67ffffffffffffffff600882901c1660208a015260ff604882901c16888a015260501c67ffffffffffffffff1690565b67ffffffffffffffff166060880152565b61185960038201611666565b608087015261186a60048201611666565b60a087015261187b6005820161173e565b60c087015201611666565b60e08401520152565b90611898611548565b9260005b600281106118ae575050505090600090565b63ffffffff84168060005260056020528260016118d3846108a78860406000206115a7565b5001541480611914575b6118ea575060010161189c565b61190e955061072994506108a79250600093919352600560205260406000206115a7565b90600190565b508215156118dd565b91611937918354906000199060031b92831b921b19161790565b9055565b6119629061076e60019363ffffffff831660005260056020526107688160406000206115a7565b50015490565b604051906119758261027d565b6060604083828152600060208201520152565b906040516119958161027d565b809260405180602083549182815201908360005260206000209060005b8181106119f957505050604092826119d16119f49460029403826102d7565b85526119ee6119e4600183015460ff1690565b60ff166020870152565b01611666565b910152565b82548452602090930192600192830192016119b2565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9081600302916003830403611a4f57565b611a0f565b81810292918115918404141715611a4f57565b60405190611a766020836102d7565b600080835282815b828110611a8a57505050565b602090604051611a99816102bb565b60008152611aa5611968565b8382015282828501015201611a7e565b90611abf82611726565b611acc60405191826102d7565b828152601f19611adc8294611726565b019060005b828110611aed57505050565b602090604051611afc816102bb565b60008152611b08611968565b8382015282828501015201611ae1565b9060018201809211611a4f57565b91908201809211611a4f57565b91908203918211611a4f57565b80518210156116025760209160051b010190565b611b618260035492611a54565b9180158015611c33575b611c2857611b799083611b26565b90808211611c20575b50611b95611b908383611b33565b611ab5565b91805b828110611ba55750505090565b80611c19611bb76113a2600194613573565b611bf8611bd88267ffffffffffffffff166000526002602052604060002090565b611bf3611be36102fa565b67ffffffffffffffff9094168452565b611988565b6020820152611c078584611b33565b90611c128289611b40565b5286611b40565b5001611b98565b905038611b82565b5050506104b8611a67565b5081831015611b6b565b60ff8116036101f857565b359061030982611c3d565b81601f820112156101f857803590611c6a8261031b565b92611c7860405194856102d7565b828452602083830101116101f857816000926020809301838601378301015290565b9080601f830112156101f857813591611cb283611726565b92611cc060405194856102d7565b80845260208085019160051b830101918383116101f85760208101915b838310611cec57505050505090565b823567ffffffffffffffff81116101f8578201906060601f1983880301126101f85760405190611d1b8261027d565b60208301358252604083013567ffffffffffffffff81116101f857876020611d4592860101611c53565b602083015260608301359167ffffffffffffffff83116101f857611d7188602080969581960101611c53565b6040820152815201920191611cdd565b919091610100818403126101f857611d9761030b565b92611da182610407565b8452611daf60208301610d7e565b6020850152611dc060408301611c48565b6040850152611dd160608301610d7e565b6060850152608082013567ffffffffffffffff81116101f85781611df6918401611c53565b608085015260a082013567ffffffffffffffff81116101f85781611e1b918401611c53565b60a085015260c082013567ffffffffffffffff81116101f85781611e40918401611c9a565b60c085015260e082013567ffffffffffffffff81116101f857611e639201611c53565b60e0830152565b63ffffffff1663ffffffff8114611a4f5760010190565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1823603018112156101f857016020813591019167ffffffffffffffff82116101f85781360383136101f857565b601f8260209493601f19938186528686013760008582860101520116010190565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1823603018112156101f857016020813591019167ffffffffffffffff82116101f8578160051b360383136101f857565b90602083828152019060208160051b85010193836000915b838310611f6d5750505050505090565b909192939495601f1982820301865286357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa1843603018112156101f8576020611ff6600193868394019081358152611fe8611fdd611fcd86850185611e81565b6060888601526060850191611ed1565b926040810190611e81565b916040818503910152611ed1565b980196019493019190611f5d565b6104b89161201a8161201584610407565b610477565b61203a61202960208401610d7e565b67ffffffffffffffff166020830152565b61205361204960408401611c48565b60ff166040830152565b61207361206260608401610d7e565b67ffffffffffffffff166060830152565b6120e06120d56120ba61209f61208c6080870187611e81565b6101006080880152610100870191611ed1565b6120ac60a0870187611e81565b9086830360a0880152611ed1565b6120c760c0860186611ef2565b9085830360c0870152611f45565b9260e0810190611e81565b9160e0818503910152611ed1565b9060206104b8928181520190612004565b356104b881610d5f565b356104b881611c3d565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156101f8570180359067ffffffffffffffff82116101f8576020019181360383136101f857565b81811061216f575050565b60008155600101612164565b9190601f811161218a57505050565b610309926000526020600020906020601f840160051c830193106121b6575b601f0160051c0190612164565b90915081906121a9565b90929167ffffffffffffffff8111610299576121e6816121e08454611613565b8461217b565b6000601f8211600114612222578190611937939495600092612217575b50506000198260011b9260031b1c19161790565b013590503880612203565b601f1982169461223784600052602060002090565b91805b878110612272575083600195969710612258575b505050811b019055565b60001960f88560031b161c1991013516905538808061224e565b9092602060018192868601358155019401910161223a565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156101f8570180359067ffffffffffffffff82116101f857602001918160051b360383136101f857565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa1813603018212156101f8570190565b61231b8154611613565b9081612325575050565b81601f60009311600114612337575055565b8183526020832061235391601f0160051c810190600101612164565b808252602082209081548360011b906000198560031b1c191617905555565b9080358255600182016123886020830183612113565b9067ffffffffffffffff8211610299576123ac826123a68554611613565b8561217b565b600090601f83116001146123fd57926123e7836123f4946002979461030999976000926122175750506000198260011b9260031b1c19161790565b90555b6040810190612113565b929091016121c0565b601f1983169161241285600052602060002090565b92815b81811061245d57509360029693610309989693600193836123f49810612443575b505050811b0190556123ea565b60001960f88560031b161c19910135169055388080612436565b91936020600181928787013581550195019201612415565b6801000000000000000083116102995780548382558084106124dd575b50906124a48192600052602060002090565b906000925b8484106124b7575050505050565b60036020826124d16124cb600195876122de565b87612372565b019301930192916124a9565b80600302906003820403611a4f5783600302600381048503611a4f5782600052602060002091820191015b8181106125155750612492565b6003906000815561252860018201612311565b61253460028201612311565b01612508565b90803591612547836103fd565b6002831015610448576123f46004926103099460ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0085541691161783556125ce612594602083016120ff565b84547fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000ff1660089190911b68ffffffffffffffff0016178455565b6126186125dd60408301612109565b84547fffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffff1660489190911b69ff00000000000000000016178455565b61266a612627606083016120ff565b84547fffffffffffffffffffffffffffff0000000000000000ffffffffffffffffffff1660509190911b71ffffffffffffffff0000000000000000000016178455565b61268461267a6080830183612113565b90600186016121c0565b61269e61269460a0830183612113565b90600286016121c0565b6126b86126ae60c083018361228a565b9060038601612475565b60e0810190612113565b60409063ffffffff6104b894931681528160208201520190612004565b91908110156116025760051b0190565b906801000000000000000081116102995781549181815582821061271257505050565b600052602060002091820191015b81811061272b575050565b60008155600101612720565b80546000825580612757575b506002816000600161030994015501612311565b816000526020600020908101905b8181106127725750612743565b60008155600101612765565b91908110156116025760051b810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc1813603018212156101f8570190565b6060813603126101f857604051906127d58261027d565b803567ffffffffffffffff81116101f857810136601f820112156101f8578035906127ff82611726565b9161280d60405193846102d7565b80835260208084019160051b830101913683116101f857602001905b82821061286f57505050825261284160208201611c48565b602083015260408101359067ffffffffffffffff82116101f85761286791369101611c53565b604082015290565b8135815260209182019101612829565b919091825167ffffffffffffffff8111610299576128a1816121e08454611613565b6020601f82116001146128dc5781906119379394956000926128d15750506000198260011b9260031b1c19161790565b015190503880612203565b601f198216906128f184600052602060002090565b9160005b81811061292c5750958360019596971061291357505050811b019055565b015160001960f88460031b161c1916905538808061224e565b9192602060018192868b0151815501940192016128f5565b90805180519067ffffffffffffffff82116102995760209061296683866126ef565b0183600052602060002060005b8381106129c357505050509060026040610309936001840160ff6020830151167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008254161790550151910161287f565b600190602084519401938184015501612973565b60409067ffffffffffffffff6104b894931681528160208201520190610d89565b906004116101f85790600490565b906024116101f85760040190602090565b919091357fffffffff0000000000000000000000000000000000000000000000000000000081169260048110612a4b575050565b7fffffffff00000000000000000000000000000000000000000000000000000000929350829060040360031b1b161690565b908160209103126101f8573590565b908092918237016000815290565b3d15612ac5573d90612aab8261031b565b91612ab960405193846102d7565b82523d6000602084013e565b606090565b909173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000163303612caa57612b1c612b1684846129f8565b90612a17565b7fffffffff0000000000000000000000000000000000000000000000000000000081167fbae4e0fa000000000000000000000000000000000000000000000000000000008114159081612c7f575b81612c54575b50612c055750612b8b612b838484612a06565b810190612a7d565b63ffffffff82168103612bcc5750506000918291612bae60405180938193612a8c565b039082305af1612bbc612a9a565b9015612bc55750565b60203d9101fd5b7f8a6e4ce80000000000000000000000000000000000000000000000000000000060005263ffffffff9081166004521660245260446000fd5b7f12ba286f000000000000000000000000000000000000000000000000000000006000527fffffffff000000000000000000000000000000000000000000000000000000001660045260246000fd5b7f5a837f97000000000000000000000000000000000000000000000000000000009150141538612b70565b7f7524051a000000000000000000000000000000000000000000000000000000008114159150612b6a565b7fac7a7efd0000000000000000000000000000000000000000000000000000000060005260046000fd5b612cf460019263ffffffff809316600052600760205260406000206115a7565b54161890565b303303612d0357565b7f371a73280000000000000000000000000000000000000000000000000000000060005260046000fd5b60405160208101906000825260208152612d486040826102d7565b51902090565b90612d5882611726565b612d6560405191826102d7565b828152601f19612d758294611726565b0190602036910137565b6000198114611a4f5760010190565b9060206104b8928181520190610484565b6020810167ffffffffffffffff612dbe825167ffffffffffffffff1690565b16156130f4578151612dcf8161043e565b612dd88161043e565b1515806130d6575b6130ac57608082015180518015918215613096575b505061306c5760a082015180518015918215613056575b505061302c57612e2d6113c26113af6113a2845167ffffffffffffffff1690565b61301857612e676112de6001612e5f6113e1612e506112de604089015160ff1690565b955167ffffffffffffffff1690565b015460ff1690565b91818311612fe45760c00191825151916101008311612fba57612e8990611a3e565b821115612f9057600091612e9c81612d4e565b9360005b828110612efe57505050612eb6612ebb91611a3e565b611b18565b90818110612ece5750506103099061347a565b7f548dd21f0000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b612f09818351611b40565b51604081015151612f80575b602081015151158015612f77575b612f3d579060019151612f368289611b40565b5201612ea0565b612f73906040519182917f9fa4031400000000000000000000000000000000000000000000000000000000835260048301612d8e565b0390fd5b50805115612f23565b94612f8a90612d7f565b94612f15565b7f4856694e0000000000000000000000000000000000000000000000000000000060005260046000fd5b7f1b925da60000000000000000000000000000000000000000000000000000000060005260046000fd5b507f2db220400000000000000000000000000000000000000000000000000000000060005260049190915260245260446000fd5b516109ff9067ffffffffffffffff16611481565b7fdee985740000000000000000000000000000000000000000000000000000000060005260046000fd5b602001209050613064612d2d565b143880612e0c565b7f358c19270000000000000000000000000000000000000000000000000000000060005260046000fd5b6020012090506130a4612d2d565b143880612df5565b7f3302dbd70000000000000000000000000000000000000000000000000000000060005260046000fd5b50600182516130e48161043e565b6130ed8161043e565b1415612de0565b7f698cf8e00000000000000000000000000000000000000000000000000000000060005260046000fd5b906131c7929361317463ffffffff9283604051957f45564d0000000000000000000000000000000000000000000000000000000000602088015246604088015230606088015216608086015260a0850190610477565b1660c082015260c0815261318960e0826102d7565b602060405193826131a38694518092858088019101610337565b83016131b782518093858085019101610337565b010103601f1981018352826102d7565b602081519101207fffff00000000000000000000000000000000000000000000000000000000000019167e0a0000000000000000000000000000000000000000000000000000000000001790565b73ffffffffffffffffffffffffffffffffffffffff60015416330361323657565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b5190610309826103e2565b9080601f830112156101f857815161328281611726565b9261329060405194856102d7565b81845260208085019260051b8201019283116101f857602001905b8282106132b85750505090565b81518152602091820191016132ab565b9080601f830112156101f85781516132df81611726565b926132ed60405194856102d7565b81845260208085019260051b8201019283116101f857602001905b8282106133155750505090565b8151815260209182019101613308565b6020818303126101f85780519067ffffffffffffffff82116101f857019080601f830112156101f85781519161335a83611726565b9261336860405194856102d7565b80845260208085019160051b830101918383116101f85760208101915b83831061339457505050505090565b825167ffffffffffffffff81116101f857820190610100601f1983880301126101f8576133bf61030b565b906133cc60208401613260565b82526133da60408401613260565b60208301526133eb60608401613260565b60408301526080830151606083015260a0830151608083015260c083015160a083015260e083015167ffffffffffffffff81116101f8578760206134319286010161326b565b60c08301526101008301519167ffffffffffffffff83116101f85761345e886020809695819601016132c8565b60e0820152815201920191613385565b6040513d6000823e3d90fd5b80516134835750565b60405180917f05a519660000000000000000000000000000000000000000000000000000000082526024820160206004840152815180915260206044840192019060005b818110613542575050509080600092038173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa801561353d5761351e5750565b61353a903d806000833e61353281836102d7565b810190613325565b50565b61346e565b82518452859450602093840193909201916001016134c7565b80548210156116025760005260206000200190600090565b6003548110156116025760036000527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b015490565b805480156135d05760001901906135bf828261355b565b60001982549160031b1b1916905555565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6000818152600460205260409020549081156136a057600019820190828211611a4f57600354926000198401938411611a4f57838360009561365f9503613665575b50505061364e60036135a8565b600490600052602052604060002090565b55600190565b61364e6136919161368761367d61369795600361355b565b90549060031b1c90565b928391600361355b565b9061191d565b55388080613641565b5050600090565b6000818152600460205260409020546137145760035468010000000000000000811015610299576136fb6136e4826001859401600355600361355b565b81939154906000199060031b92831b921b19161790565b9055600354906000526004602052604060002055600190565b5060009056fea164736f6c634300081a000a",
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
