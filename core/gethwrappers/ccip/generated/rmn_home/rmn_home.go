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
	MinObservers        uint64
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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expectedConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"gotConfigDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateOffchainPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicatePeerId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSourceChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MinObserversTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoOpStateTransitionNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfBoundsNodesLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfBoundsObserverNodeIndex\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RevokingZeroDigestNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ActiveConfigRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"CandidateConfigRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigPromoted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minObservers\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minObservers\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"DynamicConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllConfigs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minObservers\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"internalType\":\"structRMNHome.VersionedConfig\",\"name\":\"activeConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minObservers\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"internalType\":\"structRMNHome.VersionedConfig\",\"name\":\"candidateConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCandidateDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minObservers\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"internalType\":\"structRMNHome.VersionedConfig\",\"name\":\"versionedConfig\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"ok\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigDigests\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"activeConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"candidateConfigDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"digestToPromote\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"digestToRevoke\",\"type\":\"bytes32\"}],\"name\":\"promoteCandidateAndRevokeActive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"revokeCandidate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minObservers\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"digestToOverwrite\",\"type\":\"bytes32\"}],\"name\":\"setCandidate\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"newConfigDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minObservers\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"newDynamicConfig\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"currentDigest\",\"type\":\"bytes32\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600e80546001600160401b03191690553480156200002157600080fd5b503380600081620000795760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000ac57620000ac81620000b5565b50505062000160565b336001600160a01b038216036200010f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000070565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6128f180620001706000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80636dd5b69d1161008c5780638c76967f116100665780638c76967f146101d45780638da5cb5b146101e7578063f2fde38b1461020f578063fb4022d41461022257600080fd5b80636dd5b69d14610196578063736be802146101b757806379ba5097146101cc57600080fd5b80633567e6b4116100bd5780633567e6b41461015b57806338354c5c14610178578063635079561461018057600080fd5b8063118dbac5146100e4578063123e65db1461010a578063181f5a7714610112575b600080fd5b6100f76100f23660046118e4565b610235565b6040519081526020015b60405180910390f35b6100f7610418565b61014e6040518060400160405280601181526020017f524d4e486f6d6520312e362e302d64657600000000000000000000000000000081525081565b60405161010191906119bf565b610163610457565b60408051928352602083019190915201610101565b6100f76104d8565b6101886104f7565b604051610101929190611b2a565b6101a96101a4366004611b4f565b610a79565b604051610101929190611b68565b6101ca6101c5366004611b8c565b610d5d565b005b6101ca610e79565b6101ca6101e2366004611bd1565b610f76565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610101565b6101ca61021d366004611bf3565b611189565b6101ca610230366004611b4f565b61119d565b600061023f6112b9565b61025961024b85611da5565b61025485611ea6565b61133c565b60006102636104d8565b90508281146102ad576040517f93df584c00000000000000000000000000000000000000000000000000000000815260048101829052602481018490526044015b60405180910390fd5b80156102df5760405183907f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b90600090a25b600e80546000919082906102f89063ffffffff16611fbb565b91906101000a81548163ffffffff021916908363ffffffff160217905590506103408660405160200161032b9190612166565b604051602081830303815290604052826114b4565b600e54909350600090600290640100000000900463ffffffff1660011863ffffffff166002811061037357610373612179565b600602016001810185905580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff8416178155905086600282016103bd82826123c6565b50869050600482016103cf82826125c5565b905050837ff6c6d1be15ba0acc8ee645c1ec613c360ef786d2d3200eb8e695b6dec757dbf083898960405161040693929190612809565b60405180910390a25050509392505050565b60006002610434600e5463ffffffff6401000000009091041690565b63ffffffff166002811061044a5761044a612179565b6006020160010154905090565b6000806002610474600e5463ffffffff6401000000009091041690565b63ffffffff166002811061048a5761048a612179565b600602016001015460026104b2600e54600163ffffffff640100000000909204919091161890565b63ffffffff16600281106104c8576104c8612179565b6006020160010154915091509091565b600e54600090600290640100000000900463ffffffff16600118610434565b6104ff611866565b610507611866565b60006002610523600e5463ffffffff6401000000009091041690565b63ffffffff166002811061053957610539612179565b6040805160808101825260069290920292909201805463ffffffff16825260018101546020808401919091528351600283018054606093810283018401875282870181815295969495948701949293919284929091849160009085015b828210156105dc57838290600052602060002090600202016040518060400160405290816000820154815260200160018201548152505081526020019060010190610596565b5050505081526020016001820180546105f49061220d565b80601f01602080910402602001604051908101604052809291908181526020018280546106209061220d565b801561066d5780601f106106425761010080835404028352916020019161066d565b820191906000526020600020905b81548152906001019060200180831161065057829003601f168201915b50505050508152505081526020016004820160405180604001604052908160008201805480602002602001604051908101604052809291908181526020016000905b8282101561070f5760008481526020908190206040805160608101825260028602909201805467ffffffffffffffff80821685526801000000000000000090910416838501526001908101549183019190915290835290920191016106af565b5050505081526020016001820180546107279061220d565b80601f01602080910402602001604051908101604052809291908181526020018280546107539061220d565b80156107a05780601f10610775576101008083540402835291602001916107a0565b820191906000526020600020905b81548152906001019060200180831161078357829003601f168201915b505050919092525050509052506020810151909150156107be578092505b600e54600090600290640100000000900463ffffffff1660011863ffffffff16600281106107ee576107ee612179565b6040805160808101825260069290920292909201805463ffffffff16825260018101546020808401919091528351600283018054606093810283018401875282870181815295969495948701949293919284929091849160009085015b828210156108915783829060005260206000209060020201604051806040016040529081600082015481526020016001820154815250508152602001906001019061084b565b5050505081526020016001820180546108a99061220d565b80601f01602080910402602001604051908101604052809291908181526020018280546108d59061220d565b80156109225780601f106108f757610100808354040283529160200191610922565b820191906000526020600020905b81548152906001019060200180831161090557829003601f168201915b50505050508152505081526020016004820160405180604001604052908160008201805480602002602001604051908101604052809291908181526020016000905b828210156109c45760008481526020908190206040805160608101825260028602909201805467ffffffffffffffff8082168552680100000000000000009091041683850152600190810154918301919091529083529092019101610964565b5050505081526020016001820180546109dc9061220d565b80601f0160208091040260200160405190810160405280929190818152602001828054610a089061220d565b8015610a555780601f10610a2a57610100808354040283529160200191610a55565b820191906000526020600020905b815481529060010190602001808311610a3857829003601f168201915b50505091909252505050905250602081015190915015610a73578092505b50509091565b610a81611866565b6000805b6002811015610d52578360028260028110610aa257610aa2612179565b6006020160010154148015610ab657508315155b15610d4a5760028160028110610ace57610ace612179565b6040805160808101825260069290920292909201805463ffffffff16825260018082015460208085019190915284516002840180546060938102830184018852828801818152959794969588958701948492849160009085015b82821015610b6e57838290600052602060002090600202016040518060400160405290816000820154815260200160018201548152505081526020019060010190610b28565b505050508152602001600182018054610b869061220d565b80601f0160208091040260200160405190810160405280929190818152602001828054610bb29061220d565b8015610bff5780601f10610bd457610100808354040283529160200191610bff565b820191906000526020600020905b815481529060010190602001808311610be257829003601f168201915b50505050508152505081526020016004820160405180604001604052908160008201805480602002602001604051908101604052809291908181526020016000905b82821015610ca15760008481526020908190206040805160608101825260028602909201805467ffffffffffffffff8082168552680100000000000000009091041683850152600190810154918301919091529083529092019101610c41565b505050508152602001600182018054610cb99061220d565b80601f0160208091040260200160405190810160405280929190818152602001828054610ce59061220d565b8015610d325780601f10610d0757610100808354040283529160200191610d32565b820191906000526020600020905b815481529060010190602001808311610d1557829003601f168201915b50505091909252505050905250969095509350505050565b600101610a85565b509092600092509050565b610d656112b9565b60005b6002811015610e3f578160028260028110610d8557610d85612179565b6006020160010154148015610d9957508115155b15610e3757610dd0610daa84611ea6565b60028360028110610dbd57610dbd612179565b60060201600201600001805490506115bc565b8260028260028110610de457610de4612179565b600602016004018181610df791906125c5565b905050817f1f69d1a2edb327babc986b3deb80091f101b9105d42a6c30db4d99c31d7e629484604051610e2a9190612844565b60405180910390a2505050565b600101610d68565b506040517fd0b2c031000000000000000000000000000000000000000000000000000000008152600481018290526024016102a4565b5050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610efa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016102a4565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610f7e6112b9565b81158015610f8a575080155b15610fc1576040517f7b4d1e4f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e54600163ffffffff6401000000009092048216181682600282818110610feb57610feb612179565b600602016001015414611051576002816002811061100b5761100b612179565b6006020160010154836040517f93df584c0000000000000000000000000000000000000000000000000000000081526004016102a4929190918252602082015260400190565b6000600261106d600e5463ffffffff6401000000009091041690565b63ffffffff166002811061108357611083612179565b600602019050828160010154146110d65760018101546040517f93df584c0000000000000000000000000000000000000000000000000000000081526004810191909152602481018490526044016102a4565b6000600180830191909155600e805463ffffffff6401000000008083048216909418169092027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff90921691909117905582156111585760405183907f0b31c0055e2d464bef7781994b98c4ff9ef4ae0d05f59feb6a68c42de5e201b890600090a25b60405184907ffc3e98dbbd47c3fa7c1c05b6ec711caeaf70eca4554192b9ada8fc11a37f298e90600090a250505050565b6111916112b9565b61119a81611771565b50565b6111a56112b9565b806111dc576040517f0849d8cc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e54600163ffffffff640100000000909204821618168160028281811061120657611206612179565b60060201600101541461126c576002816002811061122657611226612179565b6006020160010154826040517f93df584c0000000000000000000000000000000000000000000000000000000081526004016102a4929190918252602082015260400190565b60405182907f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b90600090a2600281600281106112aa576112aa612179565b60060201600101600090555050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461133a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102a4565b565b815151610100101561137a576040517faf26d5e300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251518110156114a4576000611394826001612857565b90505b83515181101561149b5783518051829081106113b5576113b5612179565b602002602001015160000151846000015183815181106113d7576113d7612179565b6020026020010151600001510361141a576040517f221a8ae800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b835180518290811061142e5761142e612179565b6020026020010151602001518460000151838151811061145057611450612179565b60200260200101516020015103611493576040517fae00651d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600101611397565b5060010161137d565b50610e75818360000151516115bc565b604080517f45564d00000000000000000000000000000000000000000000000000000000006020820152469181019190915230606082015263ffffffff821660808201526000907dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060a001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261155b91869060200161286a565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190528051602090910120167e0b0000000000000000000000000000000000000000000000000000000000001790505b92915050565b81515160005b8181101561176b576000846000015182815181106115e2576115e2612179565b6020026020010151905060008260016115fb9190612857565b90505b8381101561167e57855180518290811061161a5761161a612179565b60200260200101516000015167ffffffffffffffff16826000015167ffffffffffffffff1603611676576040517f3857f84d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001016115fe565b5060408101518061169186610100612899565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff901c8216146116ed576040517f2847b60600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b811561171557611701600183612899565b9091169061170e816128ac565b90506116f0565b80836020015167ffffffffffffffff16111561175d576040517f4ff924ea00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050508060010190506115c2565b50505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036117f0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102a4565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6040518060800160405280600063ffffffff168152602001600080191681526020016118a5604051806040016040528060608152602001606081525090565b81526020016118c7604051806040016040528060608152602001606081525090565b905290565b6000604082840312156118de57600080fd5b50919050565b6000806000606084860312156118f957600080fd5b833567ffffffffffffffff8082111561191157600080fd5b61191d878388016118cc565b9450602086013591508082111561193357600080fd5b50611940868287016118cc565b925050604084013590509250925092565b60005b8381101561196c578181015183820152602001611954565b50506000910152565b6000815180845261198d816020860160208601611951565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006119d26020830184611975565b9392505050565b8051604080845281518482018190526000926060916020918201918388019190865b82811015611a35578451805167ffffffffffffffff90811686528382015116838601528701518785015293810193928501926001016119fb565b50808801519550888303818a01525050611a4f8185611975565b979650505050505050565b63ffffffff81511682526000602080830151818501526040808401516080604087015260c0860181516040608089015281815180845260e08a0191508683019350600092505b80831015611ac95783518051835287015187830152928601926001929092019190850190611aa0565b50948301518886037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff800160a08a015294611b038187611975565b9550505050505060608301518482036060860152611b2182826119d9565b95945050505050565b604081526000611b3d6040830185611a5a565b8281036020840152611b218185611a5a565b600060208284031215611b6157600080fd5b5035919050565b604081526000611b7b6040830185611a5a565b905082151560208301529392505050565b60008060408385031215611b9f57600080fd5b823567ffffffffffffffff811115611bb657600080fd5b611bc2858286016118cc565b95602094909401359450505050565b60008060408385031215611be457600080fd5b50508035926020909101359150565b600060208284031215611c0557600080fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146119d257600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715611c7b57611c7b611c29565b60405290565b6040516060810167ffffffffffffffff81118282101715611c7b57611c7b611c29565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611ceb57611ceb611c29565b604052919050565b600067ffffffffffffffff821115611d0d57611d0d611c29565b5060051b60200190565b600082601f830112611d2857600080fd5b813567ffffffffffffffff811115611d4257611d42611c29565b611d7360207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611ca4565b818152846020838601011115611d8857600080fd5b816020850160208301376000918101602001919091529392505050565b60006040808336031215611db857600080fd5b611dc0611c58565b833567ffffffffffffffff80821115611dd857600080fd5b9085019036601f830112611deb57600080fd5b81356020611e00611dfb83611cf3565b611ca4565b82815260069290921b84018101918181019036841115611e1f57600080fd5b948201945b83861015611e5f57878636031215611e3c5760008081fd5b611e44611c58565b86358152838701358482015282529487019490820190611e24565b86525087810135955082861115611e7557600080fd5b611e8136878a01611d17565b90850152509195945050505050565b67ffffffffffffffff8116811461119a57600080fd5b60006040808336031215611eb957600080fd5b611ec1611c58565b833567ffffffffffffffff80821115611ed957600080fd5b9085019036601f830112611eec57600080fd5b81356020611efc611dfb83611cf3565b82815260609283028501820192828201919036851115611f1b57600080fd5b958301955b84871015611f7557808736031215611f385760008081fd5b611f40611c81565b8735611f4b81611e90565b815287850135611f5a81611e90565b81860152878a01358a82015283529586019591830191611f20565b5086525087810135955082861115611e7557600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff808316818103611fd457611fd4611f8c565b6001019392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261201357600080fd5b830160208101925035905067ffffffffffffffff81111561203357600080fd5b80360382131561204257600080fd5b9250929050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6000604080840183357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18536030181126120cb57600080fd5b8401602081810191359067ffffffffffffffff8211156120ea57600080fd5b8160061b36038313156120fc57600080fd5b6040885292819052909160009190606088015b82841015612135578435815281850135828201529385019360019390930192850161210f565b6121426020890189611fde565b9650945088810360208a0152612159818787612049565b9998505050505050505050565b6020815260006119d26020830184612092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126121dd57600080fd5b83018035915067ffffffffffffffff8211156121f857600080fd5b60200191503681900382131561204257600080fd5b600181811c9082168061222157607f821691505b6020821081036118de577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b601f8211156122a6576000816000526020600020601f850160051c810160208610156122835750805b601f850160051c820191505b818110156122a25782815560010161228f565b5050505b505050565b67ffffffffffffffff8311156122c3576122c3611c29565b6122d7836122d1835461220d565b8361225a565b6000601f84116001811461232957600085156122f35750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b1783556123bf565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156123785786850135825560209485019460019092019101612358565b50868210156123b3577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b81357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18336030181126123f857600080fd5b8201803567ffffffffffffffff81111561241157600080fd5b6020820191508060061b360382131561242957600080fd5b6801000000000000000081111561244257612442611c29565b8254818455808210156124cf5760017f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff808316831461248357612483611f8c565b808416841461249457612494611f8c565b5060008560005260206000208360011b81018560011b820191505b808210156124ca5782825582848301556002820191506124af565b505050505b5060008381526020902060005b8281101561250857833582556020840135600183015560409390930192600291909101906001016124dc565b5050505061251960208301836121a8565b61176b8183600186016122ab565b813561253281611e90565b67ffffffffffffffff811690508154817fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008216178355602084013561257681611e90565b6fffffffffffffffff00000000000000008160401b16837fffffffffffffffffffffffffffffffff00000000000000000000000000000000841617178455505050604082013560018201555050565b81357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18336030181126125f757600080fd5b8201803567ffffffffffffffff81111561261057600080fd5b6020820191506060808202360383131561262957600080fd5b6801000000000000000082111561264257612642611c29565b8354828555808310156126cf5760017f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff808316831461268357612683611f8c565b808516851461269457612694611f8c565b5060008660005260206000208360011b81018660011b820191505b808210156126ca5782825582848301556002820191506126af565b505050505b5060008481526020902060005b83811015612701576126ee8583612527565b93820193600291909101906001016126dc565b505050505061251960208301836121a8565b6000604080840183357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe185360301811261274c57600080fd5b8401602081810191359067ffffffffffffffff8083111561276c57600080fd5b6060808402360385131561277f57600080fd5b60408a529483905292936000939060608a015b848610156127d65786356127a581611e90565b83168152868401356127b681611e90565b831681850152868801358882015295810195600195909501948101612792565b6127e360208b018b611fde565b985096508a810360208c01526127fa818989612049565b9b9a5050505050505050505050565b63ffffffff841681526060602082015260006128286060830185612092565b828103604084015261283a8185612713565b9695505050505050565b6020815260006119d26020830184612713565b808201808211156115b6576115b6611f8c565b6000835161287c818460208801611951565b835190830190612890818360208801611951565b01949350505050565b818103818111156115b6576115b6611f8c565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036128dd576128dd611f8c565b506001019056fea164736f6c6343000818000a",
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
