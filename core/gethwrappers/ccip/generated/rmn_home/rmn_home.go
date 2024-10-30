// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rmn_home

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

type RMNHomeDynamicConfig struct {
	SourceChains   []RMNHomeSourceChain
	OffchainConfig []byte
}

type RMNHomeNode struct {
	PeerId            [32]byte
	OffchainPublicKey [32]byte
}

type RMNHomeSourceChain struct {
	ChainSelector       uint64
	F                   uint64
	ObserverNodesBitmap *big.Int
}

type RMNHomeStaticConfig struct {
	Nodes          []RMNHomeNode
	OffchainConfig []byte
}

type RMNHomeVersionedConfig struct {
	Version       uint32
	ConfigDigest  [32]byte
	StaticConfig  RMNHomeStaticConfig
	DynamicConfig RMNHomeDynamicConfig
}

var RMNHomeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expectedConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"gotConfigDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateOffchainPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicatePeerId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSourceChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoOpStateTransitionNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughObservers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfBoundsNodesLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfBoundsObserverNodeIndex\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RevokingZeroDigestNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ActiveConfigRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"CandidateConfigRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigPromoted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"DynamicConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllConfigs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"internalType\":\"structRMNHome.VersionedConfig\",\"name\":\"activeConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"internalType\":\"structRMNHome.VersionedConfig\",\"name\":\"candidateConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCandidateDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"internalType\":\"structRMNHome.VersionedConfig\",\"name\":\"versionedConfig\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"ok\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigDigests\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"activeConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"candidateConfigDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"digestToPromote\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"digestToRevoke\",\"type\":\"bytes32\"}],\"name\":\"promoteCandidateAndRevokeActive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"revokeCandidate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"digestToOverwrite\",\"type\":\"bytes32\"}],\"name\":\"setCandidate\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"newConfigDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"newDynamicConfig\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"currentDigest\",\"type\":\"bytes32\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600e80546001600160401b03191690553480156200002157600080fd5b50336000816200004457604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b0384811691909117909155811615620000775762000077816200007f565b5050620000f9565b336001600160a01b03821603620000a957604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6128cb80620001096000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80636dd5b69d1161008c5780638c76967f116100665780638c76967f146101d45780638da5cb5b146101e7578063f2fde38b1461020f578063fb4022d41461022257600080fd5b80636dd5b69d14610196578063736be802146101b757806379ba5097146101cc57600080fd5b80633567e6b4116100bd5780633567e6b41461015b57806338354c5c14610178578063635079561461018057600080fd5b8063118dbac5146100e4578063123e65db1461010a578063181f5a7714610112575b600080fd5b6100f76100f236600461186a565b610235565b6040519081526020015b60405180910390f35b6100f7610418565b61014e6040518060400160405280601181526020017f524d4e486f6d6520312e362e302d64657600000000000000000000000000000081525081565b6040516101019190611945565b610163610457565b60408051928352602083019190915201610101565b6100f76104d8565b6101886104f7565b604051610101929190611ab0565b6101a96101a4366004611ad5565b610a79565b604051610101929190611aee565b6101ca6101c5366004611b12565b610d5d565b005b6101ca610e79565b6101ca6101e2366004611b57565b610f47565b60015460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610101565b6101ca61021d366004611b79565b61115a565b6101ca610230366004611ad5565b61116e565b600061023f61128a565b61025961024b85611d2b565b61025485611e2c565b6112dd565b60006102636104d8565b90508281146102ad576040517f93df584c00000000000000000000000000000000000000000000000000000000815260048101829052602481018490526044015b60405180910390fd5b80156102df5760405183907f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b90600090a25b600e80546000919082906102f89063ffffffff16611f41565b91906101000a81548163ffffffff021916908363ffffffff160217905590506103408660405160200161032b91906120ec565b60405160208183030381529060405282611455565b600e54909350600090600290640100000000900463ffffffff1660011863ffffffff1660028110610373576103736120ff565b600602016001810185905580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff8416178155905086600282016103bd828261234c565b50869050600482016103cf828261254b565b905050837ff6c6d1be15ba0acc8ee645c1ec613c360ef786d2d3200eb8e695b6dec757dbf08389896040516104069392919061278f565b60405180910390a25050509392505050565b60006002610434600e5463ffffffff6401000000009091041690565b63ffffffff166002811061044a5761044a6120ff565b6006020160010154905090565b6000806002610474600e5463ffffffff6401000000009091041690565b63ffffffff166002811061048a5761048a6120ff565b600602016001015460026104b2600e54600163ffffffff640100000000909204919091161890565b63ffffffff16600281106104c8576104c86120ff565b6006020160010154915091509091565b600e54600090600290640100000000900463ffffffff16600118610434565b6104ff6117ec565b6105076117ec565b60006002610523600e5463ffffffff6401000000009091041690565b63ffffffff1660028110610539576105396120ff565b6040805160808101825260069290920292909201805463ffffffff16825260018101546020808401919091528351600283018054606093810283018401875282870181815295969495948701949293919284929091849160009085015b828210156105dc57838290600052602060002090600202016040518060400160405290816000820154815260200160018201548152505081526020019060010190610596565b5050505081526020016001820180546105f490612193565b80601f016020809104026020016040519081016040528092919081815260200182805461062090612193565b801561066d5780601f106106425761010080835404028352916020019161066d565b820191906000526020600020905b81548152906001019060200180831161065057829003601f168201915b50505050508152505081526020016004820160405180604001604052908160008201805480602002602001604051908101604052809291908181526020016000905b8282101561070f5760008481526020908190206040805160608101825260028602909201805467ffffffffffffffff80821685526801000000000000000090910416838501526001908101549183019190915290835290920191016106af565b50505050815260200160018201805461072790612193565b80601f016020809104026020016040519081016040528092919081815260200182805461075390612193565b80156107a05780601f10610775576101008083540402835291602001916107a0565b820191906000526020600020905b81548152906001019060200180831161078357829003601f168201915b505050919092525050509052506020810151909150156107be578092505b600e54600090600290640100000000900463ffffffff1660011863ffffffff16600281106107ee576107ee6120ff565b6040805160808101825260069290920292909201805463ffffffff16825260018101546020808401919091528351600283018054606093810283018401875282870181815295969495948701949293919284929091849160009085015b828210156108915783829060005260206000209060020201604051806040016040529081600082015481526020016001820154815250508152602001906001019061084b565b5050505081526020016001820180546108a990612193565b80601f01602080910402602001604051908101604052809291908181526020018280546108d590612193565b80156109225780601f106108f757610100808354040283529160200191610922565b820191906000526020600020905b81548152906001019060200180831161090557829003601f168201915b50505050508152505081526020016004820160405180604001604052908160008201805480602002602001604051908101604052809291908181526020016000905b828210156109c45760008481526020908190206040805160608101825260028602909201805467ffffffffffffffff8082168552680100000000000000009091041683850152600190810154918301919091529083529092019101610964565b5050505081526020016001820180546109dc90612193565b80601f0160208091040260200160405190810160405280929190818152602001828054610a0890612193565b8015610a555780601f10610a2a57610100808354040283529160200191610a55565b820191906000526020600020905b815481529060010190602001808311610a3857829003601f168201915b50505091909252505050905250602081015190915015610a73578092505b50509091565b610a816117ec565b6000805b6002811015610d52578360028260028110610aa257610aa26120ff565b6006020160010154148015610ab657508315155b15610d4a5760028160028110610ace57610ace6120ff565b6040805160808101825260069290920292909201805463ffffffff16825260018082015460208085019190915284516002840180546060938102830184018852828801818152959794969588958701948492849160009085015b82821015610b6e57838290600052602060002090600202016040518060400160405290816000820154815260200160018201548152505081526020019060010190610b28565b505050508152602001600182018054610b8690612193565b80601f0160208091040260200160405190810160405280929190818152602001828054610bb290612193565b8015610bff5780601f10610bd457610100808354040283529160200191610bff565b820191906000526020600020905b815481529060010190602001808311610be257829003601f168201915b50505050508152505081526020016004820160405180604001604052908160008201805480602002602001604051908101604052809291908181526020016000905b82821015610ca15760008481526020908190206040805160608101825260028602909201805467ffffffffffffffff8082168552680100000000000000009091041683850152600190810154918301919091529083529092019101610c41565b505050508152602001600182018054610cb990612193565b80601f0160208091040260200160405190810160405280929190818152602001828054610ce590612193565b8015610d325780601f10610d0757610100808354040283529160200191610d32565b820191906000526020600020905b815481529060010190602001808311610d1557829003601f168201915b50505091909252505050905250969095509350505050565b600101610a85565b509092600092509050565b610d6561128a565b60005b6002811015610e3f578160028260028110610d8557610d856120ff565b6006020160010154148015610d9957508115155b15610e3757610dd0610daa84611e2c565b60028360028110610dbd57610dbd6120ff565b600602016002016000018054905061155d565b8260028260028110610de457610de46120ff565b600602016004018181610df7919061254b565b905050817f1f69d1a2edb327babc986b3deb80091f101b9105d42a6c30db4d99c31d7e629484604051610e2a91906127ca565b60405180910390a2505050565b600101610d68565b506040517fd0b2c031000000000000000000000000000000000000000000000000000000008152600481018290526024016102a4565b5050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610eca576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610f4f61128a565b81158015610f5b575080155b15610f92576040517f7b4d1e4f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e54600163ffffffff6401000000009092048216181682600282818110610fbc57610fbc6120ff565b6006020160010154146110225760028160028110610fdc57610fdc6120ff565b6006020160010154836040517f93df584c0000000000000000000000000000000000000000000000000000000081526004016102a4929190918252602082015260400190565b6000600261103e600e5463ffffffff6401000000009091041690565b63ffffffff1660028110611054576110546120ff565b600602019050828160010154146110a75760018101546040517f93df584c0000000000000000000000000000000000000000000000000000000081526004810191909152602481018490526044016102a4565b6000600180830191909155600e805463ffffffff6401000000008083048216909418169092027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff90921691909117905582156111295760405183907f0b31c0055e2d464bef7781994b98c4ff9ef4ae0d05f59feb6a68c42de5e201b890600090a25b60405184907ffc3e98dbbd47c3fa7c1c05b6ec711caeaf70eca4554192b9ada8fc11a37f298e90600090a250505050565b61116261128a565b61116b81611728565b50565b61117661128a565b806111ad576040517f0849d8cc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e54600163ffffffff64010000000090920482161816816002828181106111d7576111d76120ff565b60060201600101541461123d57600281600281106111f7576111f76120ff565b6006020160010154826040517f93df584c0000000000000000000000000000000000000000000000000000000081526004016102a4929190918252602082015260400190565b60405182907f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b90600090a26002816002811061127b5761127b6120ff565b60060201600101600090555050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146112db576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b815151610100101561131b576040517faf26d5e300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251518110156114455760006113358260016127dd565b90505b83515181101561143c578351805182908110611356576113566120ff565b60200260200101516000015184600001518381518110611378576113786120ff565b602002602001015160000151036113bb576040517f221a8ae800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83518051829081106113cf576113cf6120ff565b602002602001015160200151846000015183815181106113f1576113f16120ff565b60200260200101516020015103611434576040517fae00651d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600101611338565b5060010161131e565b50610e758183600001515161155d565b604080517f45564d00000000000000000000000000000000000000000000000000000000006020820152469181019190915230606082015263ffffffff821660808201526000907dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060a001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290526114fc9186906020016127f0565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190528051602090910120167e0b0000000000000000000000000000000000000000000000000000000000001790505b92915050565b81515160005b8181101561172257600084600001518281518110611583576115836120ff565b60200260200101519050600082600161159c91906127dd565b90505b8381101561161f5785518051829081106115bb576115bb6120ff565b60200260200101516000015167ffffffffffffffff16826000015167ffffffffffffffff1603611617576040517f3857f84d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60010161159f565b506040810151806116328661010061281f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff901c82161461168e576040517f2847b60600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b81156116b6576116a260018361281f565b909116906116af81612832565b9050611691565b60208301516116c690600261286a565b6116d1906001612896565b67ffffffffffffffff16811015611714576040517fa804bcb300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050806001019050611563565b50505050565b3373ffffffffffffffffffffffffffffffffffffffff821603611777576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6040518060800160405280600063ffffffff1681526020016000801916815260200161182b604051806040016040528060608152602001606081525090565b815260200161184d604051806040016040528060608152602001606081525090565b905290565b60006040828403121561186457600080fd5b50919050565b60008060006060848603121561187f57600080fd5b833567ffffffffffffffff8082111561189757600080fd5b6118a387838801611852565b945060208601359150808211156118b957600080fd5b506118c686828701611852565b925050604084013590509250925092565b60005b838110156118f25781810151838201526020016118da565b50506000910152565b600081518084526119138160208601602086016118d7565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061195860208301846118fb565b9392505050565b8051604080845281518482018190526000926060916020918201918388019190865b828110156119bb578451805167ffffffffffffffff9081168652838201511683860152870151878501529381019392850192600101611981565b50808801519550888303818a015250506119d581856118fb565b979650505050505050565b63ffffffff81511682526000602080830151818501526040808401516080604087015260c0860181516040608089015281815180845260e08a0191508683019350600092505b80831015611a4f5783518051835287015187830152928601926001929092019190850190611a26565b50948301518886037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff800160a08a015294611a8981876118fb565b9550505050505060608301518482036060860152611aa7828261195f565b95945050505050565b604081526000611ac360408301856119e0565b8281036020840152611aa781856119e0565b600060208284031215611ae757600080fd5b5035919050565b604081526000611b0160408301856119e0565b905082151560208301529392505050565b60008060408385031215611b2557600080fd5b823567ffffffffffffffff811115611b3c57600080fd5b611b4885828601611852565b95602094909401359450505050565b60008060408385031215611b6a57600080fd5b50508035926020909101359150565b600060208284031215611b8b57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461195857600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715611c0157611c01611baf565b60405290565b6040516060810167ffffffffffffffff81118282101715611c0157611c01611baf565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611c7157611c71611baf565b604052919050565b600067ffffffffffffffff821115611c9357611c93611baf565b5060051b60200190565b600082601f830112611cae57600080fd5b813567ffffffffffffffff811115611cc857611cc8611baf565b611cf960207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611c2a565b818152846020838601011115611d0e57600080fd5b816020850160208301376000918101602001919091529392505050565b60006040808336031215611d3e57600080fd5b611d46611bde565b833567ffffffffffffffff80821115611d5e57600080fd5b9085019036601f830112611d7157600080fd5b81356020611d86611d8183611c79565b611c2a565b82815260069290921b84018101918181019036841115611da557600080fd5b948201945b83861015611de557878636031215611dc25760008081fd5b611dca611bde565b86358152838701358482015282529487019490820190611daa565b86525087810135955082861115611dfb57600080fd5b611e0736878a01611c9d565b90850152509195945050505050565b67ffffffffffffffff8116811461116b57600080fd5b60006040808336031215611e3f57600080fd5b611e47611bde565b833567ffffffffffffffff80821115611e5f57600080fd5b9085019036601f830112611e7257600080fd5b81356020611e82611d8183611c79565b82815260609283028501820192828201919036851115611ea157600080fd5b958301955b84871015611efb57808736031215611ebe5760008081fd5b611ec6611c07565b8735611ed181611e16565b815287850135611ee081611e16565b81860152878a01358a82015283529586019591830191611ea6565b5086525087810135955082861115611dfb57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff808316818103611f5a57611f5a611f12565b6001019392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112611f9957600080fd5b830160208101925035905067ffffffffffffffff811115611fb957600080fd5b803603821315611fc857600080fd5b9250929050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6000604080840183357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe185360301811261205157600080fd5b8401602081810191359067ffffffffffffffff82111561207057600080fd5b8160061b360383131561208257600080fd5b6040885292819052909160009190606088015b828410156120bb5784358152818501358282015293850193600193909301928501612095565b6120c86020890189611f64565b9650945088810360208a01526120df818787611fcf565b9998505050505050505050565b6020815260006119586020830184612018565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261216357600080fd5b83018035915067ffffffffffffffff82111561217e57600080fd5b602001915036819003821315611fc857600080fd5b600181811c908216806121a757607f821691505b602082108103611864577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b601f82111561222c576000816000526020600020601f850160051c810160208610156122095750805b601f850160051c820191505b8181101561222857828155600101612215565b5050505b505050565b67ffffffffffffffff83111561224957612249611baf565b61225d836122578354612193565b836121e0565b6000601f8411600181146122af57600085156122795750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355612345565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156122fe57868501358255602094850194600190920191016122de565b5086821015612339577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b81357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe183360301811261237e57600080fd5b8201803567ffffffffffffffff81111561239757600080fd5b6020820191508060061b36038213156123af57600080fd5b680100000000000000008111156123c8576123c8611baf565b8254818455808210156124555760017f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff808316831461240957612409611f12565b808416841461241a5761241a611f12565b5060008560005260206000208360011b81018560011b820191505b80821015612450578282558284830155600282019150612435565b505050505b5060008381526020902060005b8281101561248e5783358255602084013560018301556040939093019260029190910190600101612462565b5050505061249f602083018361212e565b611722818360018601612231565b81356124b881611e16565b67ffffffffffffffff811690508154817fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000821617835560208401356124fc81611e16565b6fffffffffffffffff00000000000000008160401b16837fffffffffffffffffffffffffffffffff00000000000000000000000000000000841617178455505050604082013560018201555050565b81357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe183360301811261257d57600080fd5b8201803567ffffffffffffffff81111561259657600080fd5b602082019150606080820236038313156125af57600080fd5b680100000000000000008211156125c8576125c8611baf565b8354828555808310156126555760017f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff808316831461260957612609611f12565b808516851461261a5761261a611f12565b5060008660005260206000208360011b81018660011b820191505b80821015612650578282558284830155600282019150612635565b505050505b5060008481526020902060005b838110156126875761267485836124ad565b9382019360029190910190600101612662565b505050505061249f602083018361212e565b6000604080840183357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18536030181126126d257600080fd5b8401602081810191359067ffffffffffffffff808311156126f257600080fd5b6060808402360385131561270557600080fd5b60408a529483905292936000939060608a015b8486101561275c57863561272b81611e16565b831681528684013561273c81611e16565b831681850152868801358882015295810195600195909501948101612718565b61276960208b018b611f64565b985096508a810360208c0152612780818989611fcf565b9b9a5050505050505050505050565b63ffffffff841681526060602082015260006127ae6060830185612018565b82810360408401526127c08185612699565b9695505050505050565b6020815260006119586020830184612699565b8082018082111561155757611557611f12565b600083516128028184602088016118d7565b8351908301906128168183602088016118d7565b01949350505050565b8181038181111561155757611557611f12565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361286357612863611f12565b5060010190565b67ffffffffffffffff81811683821602808216919082811461288e5761288e611f12565b505092915050565b67ffffffffffffffff8181168382160190808211156128b7576128b7611f12565b509291505056fea164736f6c6343000818000a",
}

var RMNHomeABI = RMNHomeMetaData.ABI

var RMNHomeBin = RMNHomeMetaData.Bin

func DeployRMNHome(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RMNHome, error) {
	parsed, err := RMNHomeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RMNHomeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RMNHome{address: address, abi: *parsed, RMNHomeCaller: RMNHomeCaller{contract: contract}, RMNHomeTransactor: RMNHomeTransactor{contract: contract}, RMNHomeFilterer: RMNHomeFilterer{contract: contract}}, nil
}

type RMNHome struct {
	address common.Address
	abi     abi.ABI
	RMNHomeCaller
	RMNHomeTransactor
	RMNHomeFilterer
}

type RMNHomeCaller struct {
	contract *bind.BoundContract
}

type RMNHomeTransactor struct {
	contract *bind.BoundContract
}

type RMNHomeFilterer struct {
	contract *bind.BoundContract
}

type RMNHomeSession struct {
	Contract     *RMNHome
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type RMNHomeCallerSession struct {
	Contract *RMNHomeCaller
	CallOpts bind.CallOpts
}

type RMNHomeTransactorSession struct {
	Contract     *RMNHomeTransactor
	TransactOpts bind.TransactOpts
}

type RMNHomeRaw struct {
	Contract *RMNHome
}

type RMNHomeCallerRaw struct {
	Contract *RMNHomeCaller
}

type RMNHomeTransactorRaw struct {
	Contract *RMNHomeTransactor
}

func NewRMNHome(address common.Address, backend bind.ContractBackend) (*RMNHome, error) {
	abi, err := abi.JSON(strings.NewReader(RMNHomeABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindRMNHome(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RMNHome{address: address, abi: abi, RMNHomeCaller: RMNHomeCaller{contract: contract}, RMNHomeTransactor: RMNHomeTransactor{contract: contract}, RMNHomeFilterer: RMNHomeFilterer{contract: contract}}, nil
}

func NewRMNHomeCaller(address common.Address, caller bind.ContractCaller) (*RMNHomeCaller, error) {
	contract, err := bindRMNHome(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RMNHomeCaller{contract: contract}, nil
}

func NewRMNHomeTransactor(address common.Address, transactor bind.ContractTransactor) (*RMNHomeTransactor, error) {
	contract, err := bindRMNHome(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RMNHomeTransactor{contract: contract}, nil
}

func NewRMNHomeFilterer(address common.Address, filterer bind.ContractFilterer) (*RMNHomeFilterer, error) {
	contract, err := bindRMNHome(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RMNHomeFilterer{contract: contract}, nil
}

func bindRMNHome(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RMNHomeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_RMNHome *RMNHomeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNHome.Contract.RMNHomeCaller.contract.Call(opts, result, method, params...)
}

func (_RMNHome *RMNHomeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNHome.Contract.RMNHomeTransactor.contract.Transfer(opts)
}

func (_RMNHome *RMNHomeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNHome.Contract.RMNHomeTransactor.contract.Transact(opts, method, params...)
}

func (_RMNHome *RMNHomeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNHome.Contract.contract.Call(opts, result, method, params...)
}

func (_RMNHome *RMNHomeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNHome.Contract.contract.Transfer(opts)
}

func (_RMNHome *RMNHomeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNHome.Contract.contract.Transact(opts, method, params...)
}

func (_RMNHome *RMNHomeCaller) GetActiveDigest(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getActiveDigest")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_RMNHome *RMNHomeSession) GetActiveDigest() ([32]byte, error) {
	return _RMNHome.Contract.GetActiveDigest(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) GetActiveDigest() ([32]byte, error) {
	return _RMNHome.Contract.GetActiveDigest(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) GetAllConfigs(opts *bind.CallOpts) (GetAllConfigs,

	error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getAllConfigs")

	outstruct := new(GetAllConfigs)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ActiveConfig = *abi.ConvertType(out[0], new(RMNHomeVersionedConfig)).(*RMNHomeVersionedConfig)
	outstruct.CandidateConfig = *abi.ConvertType(out[1], new(RMNHomeVersionedConfig)).(*RMNHomeVersionedConfig)

	return *outstruct, err

}

func (_RMNHome *RMNHomeSession) GetAllConfigs() (GetAllConfigs,

	error) {
	return _RMNHome.Contract.GetAllConfigs(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) GetAllConfigs() (GetAllConfigs,

	error) {
	return _RMNHome.Contract.GetAllConfigs(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) GetCandidateDigest(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getCandidateDigest")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_RMNHome *RMNHomeSession) GetCandidateDigest() ([32]byte, error) {
	return _RMNHome.Contract.GetCandidateDigest(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) GetCandidateDigest() ([32]byte, error) {
	return _RMNHome.Contract.GetCandidateDigest(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) GetConfig(opts *bind.CallOpts, configDigest [32]byte) (GetConfig,

	error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getConfig", configDigest)

	outstruct := new(GetConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.VersionedConfig = *abi.ConvertType(out[0], new(RMNHomeVersionedConfig)).(*RMNHomeVersionedConfig)
	outstruct.Ok = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

func (_RMNHome *RMNHomeSession) GetConfig(configDigest [32]byte) (GetConfig,

	error) {
	return _RMNHome.Contract.GetConfig(&_RMNHome.CallOpts, configDigest)
}

func (_RMNHome *RMNHomeCallerSession) GetConfig(configDigest [32]byte) (GetConfig,

	error) {
	return _RMNHome.Contract.GetConfig(&_RMNHome.CallOpts, configDigest)
}

func (_RMNHome *RMNHomeCaller) GetConfigDigests(opts *bind.CallOpts) (GetConfigDigests,

	error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getConfigDigests")

	outstruct := new(GetConfigDigests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ActiveConfigDigest = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.CandidateConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_RMNHome *RMNHomeSession) GetConfigDigests() (GetConfigDigests,

	error) {
	return _RMNHome.Contract.GetConfigDigests(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) GetConfigDigests() (GetConfigDigests,

	error) {
	return _RMNHome.Contract.GetConfigDigests(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_RMNHome *RMNHomeSession) Owner() (common.Address, error) {
	return _RMNHome.Contract.Owner(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) Owner() (common.Address, error) {
	return _RMNHome.Contract.Owner(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_RMNHome *RMNHomeSession) TypeAndVersion() (string, error) {
	return _RMNHome.Contract.TypeAndVersion(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) TypeAndVersion() (string, error) {
	return _RMNHome.Contract.TypeAndVersion(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "acceptOwnership")
}

func (_RMNHome *RMNHomeSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNHome.Contract.AcceptOwnership(&_RMNHome.TransactOpts)
}

func (_RMNHome *RMNHomeTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNHome.Contract.AcceptOwnership(&_RMNHome.TransactOpts)
}

func (_RMNHome *RMNHomeTransactor) PromoteCandidateAndRevokeActive(opts *bind.TransactOpts, digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "promoteCandidateAndRevokeActive", digestToPromote, digestToRevoke)
}

func (_RMNHome *RMNHomeSession) PromoteCandidateAndRevokeActive(digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.PromoteCandidateAndRevokeActive(&_RMNHome.TransactOpts, digestToPromote, digestToRevoke)
}

func (_RMNHome *RMNHomeTransactorSession) PromoteCandidateAndRevokeActive(digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.PromoteCandidateAndRevokeActive(&_RMNHome.TransactOpts, digestToPromote, digestToRevoke)
}

func (_RMNHome *RMNHomeTransactor) RevokeCandidate(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "revokeCandidate", configDigest)
}

func (_RMNHome *RMNHomeSession) RevokeCandidate(configDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.RevokeCandidate(&_RMNHome.TransactOpts, configDigest)
}

func (_RMNHome *RMNHomeTransactorSession) RevokeCandidate(configDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.RevokeCandidate(&_RMNHome.TransactOpts, configDigest)
}

func (_RMNHome *RMNHomeTransactor) SetCandidate(opts *bind.TransactOpts, staticConfig RMNHomeStaticConfig, dynamicConfig RMNHomeDynamicConfig, digestToOverwrite [32]byte) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "setCandidate", staticConfig, dynamicConfig, digestToOverwrite)
}

func (_RMNHome *RMNHomeSession) SetCandidate(staticConfig RMNHomeStaticConfig, dynamicConfig RMNHomeDynamicConfig, digestToOverwrite [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.SetCandidate(&_RMNHome.TransactOpts, staticConfig, dynamicConfig, digestToOverwrite)
}

func (_RMNHome *RMNHomeTransactorSession) SetCandidate(staticConfig RMNHomeStaticConfig, dynamicConfig RMNHomeDynamicConfig, digestToOverwrite [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.SetCandidate(&_RMNHome.TransactOpts, staticConfig, dynamicConfig, digestToOverwrite)
}

func (_RMNHome *RMNHomeTransactor) SetDynamicConfig(opts *bind.TransactOpts, newDynamicConfig RMNHomeDynamicConfig, currentDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "setDynamicConfig", newDynamicConfig, currentDigest)
}

func (_RMNHome *RMNHomeSession) SetDynamicConfig(newDynamicConfig RMNHomeDynamicConfig, currentDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.SetDynamicConfig(&_RMNHome.TransactOpts, newDynamicConfig, currentDigest)
}

func (_RMNHome *RMNHomeTransactorSession) SetDynamicConfig(newDynamicConfig RMNHomeDynamicConfig, currentDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.SetDynamicConfig(&_RMNHome.TransactOpts, newDynamicConfig, currentDigest)
}

func (_RMNHome *RMNHomeTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "transferOwnership", to)
}

func (_RMNHome *RMNHomeSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNHome.Contract.TransferOwnership(&_RMNHome.TransactOpts, to)
}

func (_RMNHome *RMNHomeTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNHome.Contract.TransferOwnership(&_RMNHome.TransactOpts, to)
}

type RMNHomeActiveConfigRevokedIterator struct {
	Event *RMNHomeActiveConfigRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeActiveConfigRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeActiveConfigRevoked)
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
		it.Event = new(RMNHomeActiveConfigRevoked)
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

func (it *RMNHomeActiveConfigRevokedIterator) Error() error {
	return it.fail
}

func (it *RMNHomeActiveConfigRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeActiveConfigRevoked struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterActiveConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeActiveConfigRevokedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "ActiveConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeActiveConfigRevokedIterator{contract: _RMNHome.contract, event: "ActiveConfigRevoked", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchActiveConfigRevoked(opts *bind.WatchOpts, sink chan<- *RMNHomeActiveConfigRevoked, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "ActiveConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeActiveConfigRevoked)
				if err := _RMNHome.contract.UnpackLog(event, "ActiveConfigRevoked", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseActiveConfigRevoked(log types.Log) (*RMNHomeActiveConfigRevoked, error) {
	event := new(RMNHomeActiveConfigRevoked)
	if err := _RMNHome.contract.UnpackLog(event, "ActiveConfigRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeCandidateConfigRevokedIterator struct {
	Event *RMNHomeCandidateConfigRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeCandidateConfigRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeCandidateConfigRevoked)
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
		it.Event = new(RMNHomeCandidateConfigRevoked)
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

func (it *RMNHomeCandidateConfigRevokedIterator) Error() error {
	return it.fail
}

func (it *RMNHomeCandidateConfigRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeCandidateConfigRevoked struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterCandidateConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeCandidateConfigRevokedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "CandidateConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeCandidateConfigRevokedIterator{contract: _RMNHome.contract, event: "CandidateConfigRevoked", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchCandidateConfigRevoked(opts *bind.WatchOpts, sink chan<- *RMNHomeCandidateConfigRevoked, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "CandidateConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeCandidateConfigRevoked)
				if err := _RMNHome.contract.UnpackLog(event, "CandidateConfigRevoked", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseCandidateConfigRevoked(log types.Log) (*RMNHomeCandidateConfigRevoked, error) {
	event := new(RMNHomeCandidateConfigRevoked)
	if err := _RMNHome.contract.UnpackLog(event, "CandidateConfigRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeConfigPromotedIterator struct {
	Event *RMNHomeConfigPromoted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeConfigPromotedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeConfigPromoted)
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
		it.Event = new(RMNHomeConfigPromoted)
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

func (it *RMNHomeConfigPromotedIterator) Error() error {
	return it.fail
}

func (it *RMNHomeConfigPromotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeConfigPromoted struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterConfigPromoted(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeConfigPromotedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "ConfigPromoted", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeConfigPromotedIterator{contract: _RMNHome.contract, event: "ConfigPromoted", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchConfigPromoted(opts *bind.WatchOpts, sink chan<- *RMNHomeConfigPromoted, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "ConfigPromoted", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeConfigPromoted)
				if err := _RMNHome.contract.UnpackLog(event, "ConfigPromoted", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseConfigPromoted(log types.Log) (*RMNHomeConfigPromoted, error) {
	event := new(RMNHomeConfigPromoted)
	if err := _RMNHome.contract.UnpackLog(event, "ConfigPromoted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeConfigSetIterator struct {
	Event *RMNHomeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeConfigSet)
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
		it.Event = new(RMNHomeConfigSet)
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

func (it *RMNHomeConfigSetIterator) Error() error {
	return it.fail
}

func (it *RMNHomeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeConfigSet struct {
	ConfigDigest  [32]byte
	Version       uint32
	StaticConfig  RMNHomeStaticConfig
	DynamicConfig RMNHomeDynamicConfig
	Raw           types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeConfigSetIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "ConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeConfigSetIterator{contract: _RMNHome.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *RMNHomeConfigSet, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "ConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeConfigSet)
				if err := _RMNHome.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseConfigSet(log types.Log) (*RMNHomeConfigSet, error) {
	event := new(RMNHomeConfigSet)
	if err := _RMNHome.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeDynamicConfigSetIterator struct {
	Event *RMNHomeDynamicConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeDynamicConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeDynamicConfigSet)
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
		it.Event = new(RMNHomeDynamicConfigSet)
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

func (it *RMNHomeDynamicConfigSetIterator) Error() error {
	return it.fail
}

func (it *RMNHomeDynamicConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeDynamicConfigSet struct {
	ConfigDigest  [32]byte
	DynamicConfig RMNHomeDynamicConfig
	Raw           types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterDynamicConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeDynamicConfigSetIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "DynamicConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeDynamicConfigSetIterator{contract: _RMNHome.contract, event: "DynamicConfigSet", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *RMNHomeDynamicConfigSet, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "DynamicConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeDynamicConfigSet)
				if err := _RMNHome.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseDynamicConfigSet(log types.Log) (*RMNHomeDynamicConfigSet, error) {
	event := new(RMNHomeDynamicConfigSet)
	if err := _RMNHome.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeOwnershipTransferRequestedIterator struct {
	Event *RMNHomeOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeOwnershipTransferRequested)
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
		it.Event = new(RMNHomeOwnershipTransferRequested)
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

func (it *RMNHomeOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *RMNHomeOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNHomeOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeOwnershipTransferRequestedIterator{contract: _RMNHome.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNHomeOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeOwnershipTransferRequested)
				if err := _RMNHome.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseOwnershipTransferRequested(log types.Log) (*RMNHomeOwnershipTransferRequested, error) {
	event := new(RMNHomeOwnershipTransferRequested)
	if err := _RMNHome.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeOwnershipTransferredIterator struct {
	Event *RMNHomeOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeOwnershipTransferred)
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
		it.Event = new(RMNHomeOwnershipTransferred)
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

func (it *RMNHomeOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *RMNHomeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNHomeOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeOwnershipTransferredIterator{contract: _RMNHome.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNHomeOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeOwnershipTransferred)
				if err := _RMNHome.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseOwnershipTransferred(log types.Log) (*RMNHomeOwnershipTransferred, error) {
	event := new(RMNHomeOwnershipTransferred)
	if err := _RMNHome.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetAllConfigs struct {
	ActiveConfig    RMNHomeVersionedConfig
	CandidateConfig RMNHomeVersionedConfig
}
type GetConfig struct {
	VersionedConfig RMNHomeVersionedConfig
	Ok              bool
}
type GetConfigDigests struct {
	ActiveConfigDigest    [32]byte
	CandidateConfigDigest [32]byte
}

func (_RMNHome *RMNHome) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _RMNHome.abi.Events["ActiveConfigRevoked"].ID:
		return _RMNHome.ParseActiveConfigRevoked(log)
	case _RMNHome.abi.Events["CandidateConfigRevoked"].ID:
		return _RMNHome.ParseCandidateConfigRevoked(log)
	case _RMNHome.abi.Events["ConfigPromoted"].ID:
		return _RMNHome.ParseConfigPromoted(log)
	case _RMNHome.abi.Events["ConfigSet"].ID:
		return _RMNHome.ParseConfigSet(log)
	case _RMNHome.abi.Events["DynamicConfigSet"].ID:
		return _RMNHome.ParseDynamicConfigSet(log)
	case _RMNHome.abi.Events["OwnershipTransferRequested"].ID:
		return _RMNHome.ParseOwnershipTransferRequested(log)
	case _RMNHome.abi.Events["OwnershipTransferred"].ID:
		return _RMNHome.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (RMNHomeActiveConfigRevoked) Topic() common.Hash {
	return common.HexToHash("0x0b31c0055e2d464bef7781994b98c4ff9ef4ae0d05f59feb6a68c42de5e201b8")
}

func (RMNHomeCandidateConfigRevoked) Topic() common.Hash {
	return common.HexToHash("0x53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b")
}

func (RMNHomeConfigPromoted) Topic() common.Hash {
	return common.HexToHash("0xfc3e98dbbd47c3fa7c1c05b6ec711caeaf70eca4554192b9ada8fc11a37f298e")
}

func (RMNHomeConfigSet) Topic() common.Hash {
	return common.HexToHash("0xf6c6d1be15ba0acc8ee645c1ec613c360ef786d2d3200eb8e695b6dec757dbf0")
}

func (RMNHomeDynamicConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1f69d1a2edb327babc986b3deb80091f101b9105d42a6c30db4d99c31d7e6294")
}

func (RMNHomeOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (RMNHomeOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_RMNHome *RMNHome) Address() common.Address {
	return _RMNHome.address
}

type RMNHomeInterface interface {
	GetActiveDigest(opts *bind.CallOpts) ([32]byte, error)

	GetAllConfigs(opts *bind.CallOpts) (GetAllConfigs,

		error)

	GetCandidateDigest(opts *bind.CallOpts) ([32]byte, error)

	GetConfig(opts *bind.CallOpts, configDigest [32]byte) (GetConfig,

		error)

	GetConfigDigests(opts *bind.CallOpts) (GetConfigDigests,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	PromoteCandidateAndRevokeActive(opts *bind.TransactOpts, digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error)

	RevokeCandidate(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error)

	SetCandidate(opts *bind.TransactOpts, staticConfig RMNHomeStaticConfig, dynamicConfig RMNHomeDynamicConfig, digestToOverwrite [32]byte) (*types.Transaction, error)

	SetDynamicConfig(opts *bind.TransactOpts, newDynamicConfig RMNHomeDynamicConfig, currentDigest [32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterActiveConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeActiveConfigRevokedIterator, error)

	WatchActiveConfigRevoked(opts *bind.WatchOpts, sink chan<- *RMNHomeActiveConfigRevoked, configDigest [][32]byte) (event.Subscription, error)

	ParseActiveConfigRevoked(log types.Log) (*RMNHomeActiveConfigRevoked, error)

	FilterCandidateConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeCandidateConfigRevokedIterator, error)

	WatchCandidateConfigRevoked(opts *bind.WatchOpts, sink chan<- *RMNHomeCandidateConfigRevoked, configDigest [][32]byte) (event.Subscription, error)

	ParseCandidateConfigRevoked(log types.Log) (*RMNHomeCandidateConfigRevoked, error)

	FilterConfigPromoted(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeConfigPromotedIterator, error)

	WatchConfigPromoted(opts *bind.WatchOpts, sink chan<- *RMNHomeConfigPromoted, configDigest [][32]byte) (event.Subscription, error)

	ParseConfigPromoted(log types.Log) (*RMNHomeConfigPromoted, error)

	FilterConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *RMNHomeConfigSet, configDigest [][32]byte) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*RMNHomeConfigSet, error)

	FilterDynamicConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeDynamicConfigSetIterator, error)

	WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *RMNHomeDynamicConfigSet, configDigest [][32]byte) (event.Subscription, error)

	ParseDynamicConfigSet(log types.Log) (*RMNHomeDynamicConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNHomeOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNHomeOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*RMNHomeOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNHomeOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNHomeOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*RMNHomeOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
