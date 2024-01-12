// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package channel_config_store

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

type IChannelConfigStoreChannelConfiguration struct {
	ChannelConfigId [32]byte
}

type IChannelConfigStoreChannelDefinition struct {
	ReportFormat  [8]byte
	ChainSelector uint64
	StreamIDs     [][32]byte
}

var ChannelConfigStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ChannelDefinitionNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyStreamIDs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByEOA\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StagingConfigAlreadyPromoted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelector\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroReportFormat\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"ChannelDefinitionRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes8\",\"name\":\"reportFormat\",\"type\":\"bytes8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32[]\",\"name\":\"streamIDs\",\"type\":\"bytes32[]\"}],\"indexed\":false,\"internalType\":\"structIChannelConfigStore.ChannelDefinition\",\"name\":\"channelDefinition\",\"type\":\"tuple\"}],\"name\":\"NewChannelDefinition\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelConfigId\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structIChannelConfigStore.ChannelConfiguration\",\"name\":\"channelConfig\",\"type\":\"tuple\"}],\"name\":\"NewProductionConfig\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelConfigId\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structIChannelConfigStore.ChannelConfiguration\",\"name\":\"channelConfig\",\"type\":\"tuple\"}],\"name\":\"NewStagingConfig\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"PromoteStagingConfig\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes8\",\"name\":\"reportFormat\",\"type\":\"bytes8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32[]\",\"name\":\"streamIDs\",\"type\":\"bytes32[]\"}],\"internalType\":\"structIChannelConfigStore.ChannelDefinition\",\"name\":\"channelDefinition\",\"type\":\"tuple\"}],\"name\":\"addChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"getChannelDefinitions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes8\",\"name\":\"reportFormat\",\"type\":\"bytes8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32[]\",\"name\":\"streamIDs\",\"type\":\"bytes32[]\"}],\"internalType\":\"structIChannelConfigStore.ChannelDefinition\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"promoteStagingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"removeChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelConfigId\",\"type\":\"bytes32\"}],\"internalType\":\"structIChannelConfigStore.ChannelConfiguration\",\"name\":\"channelConfig\",\"type\":\"tuple\"}],\"name\":\"setStagingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610eef806101576000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c806379ba5097116100765780639a223b8e1161005b5780639a223b8e146101d7578063dd2ef3ef146101ea578063f2fde38b1461020a57600080fd5b806379ba5097146101a75780638da5cb5b146101af57600080fd5b8063181f5a77116100a7578063181f5a7714610142578063414aa93414610181578063567ed3371461019457600080fd5b806301ffc9a7146100c357806302da0e511461012d575b600080fd5b6101186100d13660046108c8565b7fffffffff00000000000000000000000000000000000000000000000000000000167f52e2bc33000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b61014061013b366004610911565b61021d565b005b604080518082018252601881527f4368616e6e656c436f6e66696753746f726520302e302e30000000000000000060208201529051610124919061092a565b61014061018f366004610996565b6102e9565b6101406101a2366004610911565b610340565b6101406103a0565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610124565b6101406101e53660046109ec565b6104a2565b6101fd6101f8366004610911565b6105fb565b6040516101249190610a3a565b610140610218366004610acc565b61070a565b61022561071e565b600081815260026020526040812060010154900361026f576040517fd1a751e200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081815260026020526040812080547fffffffffffffffffffffffffffffffff00000000000000000000000000000000168155906102b16001830182610896565b50506040518181527fa65638d745b306456ab0961a502338c1f24d1f962615be3d154c4e86e879fc749060200160405180910390a150565b6102f161071e565b60008281526004602052604090208135815581905050604051813581527f56d7e0e88863044d9b3e139e6e9c18977c6e81c387e44ce9661611f53035c7ed906020015b60405180910390a15050565b61034861071e565b6000818152600460209081526040808320815180840183529054808252858552600384529382902093909355518381527fe644aaaa8169119e133c9b338279b4305419a255ace92b4383df2f45f7daa7a89101610334565b60015473ffffffffffffffffffffffffffffffffffffffff163314610426576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6104aa61071e565b6104b76040820182610b02565b90506000036104f2576040517f4b620e2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105026040820160208301610b87565b67ffffffffffffffff16600003610545576040517ff89d762900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105526020820182610bd2565b7fffffffffffffffff000000000000000000000000000000000000000000000000166000036105ad576040517febd3ef0200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260026020526040902081906105c78282610cb7565b9050507f3f9f883dedd481d6ebe8e3efd5ef57f4b0293a3cf1d85946913dd82723542cc68282604051610334929190610e05565b60408051606080820183526000808352602083015291810191909152333214610650576040517f74e2cd5100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000828152600260209081526040918290208251606081018452815460c081901b7fffffffffffffffff00000000000000000000000000000000000000000000000016825268010000000000000000900467ffffffffffffffff1681840152600182018054855181860281018601875281815292959394938601938301828280156106fa57602002820191906000526020600020905b8154815260200190600101908083116106e6575b5050505050815250509050919050565b61071261071e565b61071b816107a1565b50565b60005473ffffffffffffffffffffffffffffffffffffffff16331461079f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161041d565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610820576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161041d565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b508054600082559060005260206000209081019061071b91905b808211156108c457600081556001016108b0565b5090565b6000602082840312156108da57600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461090a57600080fd5b9392505050565b60006020828403121561092357600080fd5b5035919050565b600060208083528351808285015260005b818110156109575785810183015185820160400152820161093b565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60008082840360408112156109aa57600080fd5b8335925060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0820112156109de57600080fd5b506020830190509250929050565b600080604083850312156109ff57600080fd5b82359150602083013567ffffffffffffffff811115610a1d57600080fd5b830160608186031215610a2f57600080fd5b809150509250929050565b60006020808352608083017fffffffffffffffff0000000000000000000000000000000000000000000000008551168285015267ffffffffffffffff82860151166040850152604085015160608086015281815180845260a0870191508483019350600092505b80831015610ac15783518252928401926001929092019190840190610aa1565b509695505050505050565b600060208284031215610ade57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461090a57600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610b3757600080fd5b83018035915067ffffffffffffffff821115610b5257600080fd5b6020019150600581901b3603821315610b6a57600080fd5b9250929050565b67ffffffffffffffff8116811461071b57600080fd5b600060208284031215610b9957600080fd5b813561090a81610b71565b7fffffffffffffffff0000000000000000000000000000000000000000000000008116811461071b57600080fd5b600060208284031215610be457600080fd5b813561090a81610ba4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff831115610c3657610c36610bef565b68010000000000000000831115610c4f57610c4f610bef565b805483825580841015610c86576000828152602081208581019083015b80821015610c8257828255600182019150610c6c565b5050505b5060008181526020812083915b85811015610caf57823582820155602090920191600101610c93565b505050505050565b8135610cc281610ba4565b8060c01c90508154817fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000082161783556020840135610cff81610b71565b6fffffffffffffffff00000000000000008160401b16837fffffffffffffffffffffffffffffffff0000000000000000000000000000000084161717845550505060408201357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1833603018112610d7557600080fd5b8201803567ffffffffffffffff811115610d8e57600080fd5b6020820191508060051b3603821315610da657600080fd5b610db4818360018601610c1e565b50505050565b81835260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115610dec57600080fd5b8260051b80836020870137939093016020019392505050565b8281526040602082015260008235610e1c81610ba4565b7fffffffffffffffff0000000000000000000000000000000000000000000000001660408301526020830135610e5181610b71565b67ffffffffffffffff8082166060850152604085013591507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1853603018212610e9957600080fd5b6020918501918201913581811115610eb057600080fd5b8060051b3603831315610ec257600080fd5b60606080860152610ed760a086018285610dba565b97965050505050505056fea164736f6c6343000810000a",
}

var ChannelConfigStoreABI = ChannelConfigStoreMetaData.ABI

var ChannelConfigStoreBin = ChannelConfigStoreMetaData.Bin

func DeployChannelConfigStore(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ChannelConfigStore, error) {
	parsed, err := ChannelConfigStoreMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChannelConfigStoreBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChannelConfigStore{address: address, abi: *parsed, ChannelConfigStoreCaller: ChannelConfigStoreCaller{contract: contract}, ChannelConfigStoreTransactor: ChannelConfigStoreTransactor{contract: contract}, ChannelConfigStoreFilterer: ChannelConfigStoreFilterer{contract: contract}}, nil
}

type ChannelConfigStore struct {
	address common.Address
	abi     abi.ABI
	ChannelConfigStoreCaller
	ChannelConfigStoreTransactor
	ChannelConfigStoreFilterer
}

type ChannelConfigStoreCaller struct {
	contract *bind.BoundContract
}

type ChannelConfigStoreTransactor struct {
	contract *bind.BoundContract
}

type ChannelConfigStoreFilterer struct {
	contract *bind.BoundContract
}

type ChannelConfigStoreSession struct {
	Contract     *ChannelConfigStore
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ChannelConfigStoreCallerSession struct {
	Contract *ChannelConfigStoreCaller
	CallOpts bind.CallOpts
}

type ChannelConfigStoreTransactorSession struct {
	Contract     *ChannelConfigStoreTransactor
	TransactOpts bind.TransactOpts
}

type ChannelConfigStoreRaw struct {
	Contract *ChannelConfigStore
}

type ChannelConfigStoreCallerRaw struct {
	Contract *ChannelConfigStoreCaller
}

type ChannelConfigStoreTransactorRaw struct {
	Contract *ChannelConfigStoreTransactor
}

func NewChannelConfigStore(address common.Address, backend bind.ContractBackend) (*ChannelConfigStore, error) {
	abi, err := abi.JSON(strings.NewReader(ChannelConfigStoreABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindChannelConfigStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStore{address: address, abi: abi, ChannelConfigStoreCaller: ChannelConfigStoreCaller{contract: contract}, ChannelConfigStoreTransactor: ChannelConfigStoreTransactor{contract: contract}, ChannelConfigStoreFilterer: ChannelConfigStoreFilterer{contract: contract}}, nil
}

func NewChannelConfigStoreCaller(address common.Address, caller bind.ContractCaller) (*ChannelConfigStoreCaller, error) {
	contract, err := bindChannelConfigStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreCaller{contract: contract}, nil
}

func NewChannelConfigStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*ChannelConfigStoreTransactor, error) {
	contract, err := bindChannelConfigStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreTransactor{contract: contract}, nil
}

func NewChannelConfigStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*ChannelConfigStoreFilterer, error) {
	contract, err := bindChannelConfigStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreFilterer{contract: contract}, nil
}

func bindChannelConfigStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChannelConfigStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ChannelConfigStore *ChannelConfigStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChannelConfigStore.Contract.ChannelConfigStoreCaller.contract.Call(opts, result, method, params...)
}

func (_ChannelConfigStore *ChannelConfigStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.ChannelConfigStoreTransactor.contract.Transfer(opts)
}

func (_ChannelConfigStore *ChannelConfigStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.ChannelConfigStoreTransactor.contract.Transact(opts, method, params...)
}

func (_ChannelConfigStore *ChannelConfigStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChannelConfigStore.Contract.contract.Call(opts, result, method, params...)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.contract.Transfer(opts)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.contract.Transact(opts, method, params...)
}

func (_ChannelConfigStore *ChannelConfigStoreCaller) GetChannelDefinitions(opts *bind.CallOpts, channelId [32]byte) (IChannelConfigStoreChannelDefinition, error) {
	var out []interface{}
	err := _ChannelConfigStore.contract.Call(opts, &out, "getChannelDefinitions", channelId)

	if err != nil {
		return *new(IChannelConfigStoreChannelDefinition), err
	}

	out0 := *abi.ConvertType(out[0], new(IChannelConfigStoreChannelDefinition)).(*IChannelConfigStoreChannelDefinition)

	return out0, err

}

func (_ChannelConfigStore *ChannelConfigStoreSession) GetChannelDefinitions(channelId [32]byte) (IChannelConfigStoreChannelDefinition, error) {
	return _ChannelConfigStore.Contract.GetChannelDefinitions(&_ChannelConfigStore.CallOpts, channelId)
}

func (_ChannelConfigStore *ChannelConfigStoreCallerSession) GetChannelDefinitions(channelId [32]byte) (IChannelConfigStoreChannelDefinition, error) {
	return _ChannelConfigStore.Contract.GetChannelDefinitions(&_ChannelConfigStore.CallOpts, channelId)
}

func (_ChannelConfigStore *ChannelConfigStoreCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChannelConfigStore.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ChannelConfigStore *ChannelConfigStoreSession) Owner() (common.Address, error) {
	return _ChannelConfigStore.Contract.Owner(&_ChannelConfigStore.CallOpts)
}

func (_ChannelConfigStore *ChannelConfigStoreCallerSession) Owner() (common.Address, error) {
	return _ChannelConfigStore.Contract.Owner(&_ChannelConfigStore.CallOpts)
}

func (_ChannelConfigStore *ChannelConfigStoreCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ChannelConfigStore.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ChannelConfigStore *ChannelConfigStoreSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ChannelConfigStore.Contract.SupportsInterface(&_ChannelConfigStore.CallOpts, interfaceId)
}

func (_ChannelConfigStore *ChannelConfigStoreCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ChannelConfigStore.Contract.SupportsInterface(&_ChannelConfigStore.CallOpts, interfaceId)
}

func (_ChannelConfigStore *ChannelConfigStoreCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ChannelConfigStore.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_ChannelConfigStore *ChannelConfigStoreSession) TypeAndVersion() (string, error) {
	return _ChannelConfigStore.Contract.TypeAndVersion(&_ChannelConfigStore.CallOpts)
}

func (_ChannelConfigStore *ChannelConfigStoreCallerSession) TypeAndVersion() (string, error) {
	return _ChannelConfigStore.Contract.TypeAndVersion(&_ChannelConfigStore.CallOpts)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelConfigStore.contract.Transact(opts, "acceptOwnership")
}

func (_ChannelConfigStore *ChannelConfigStoreSession) AcceptOwnership() (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.AcceptOwnership(&_ChannelConfigStore.TransactOpts)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.AcceptOwnership(&_ChannelConfigStore.TransactOpts)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactor) AddChannel(opts *bind.TransactOpts, channelId [32]byte, channelDefinition IChannelConfigStoreChannelDefinition) (*types.Transaction, error) {
	return _ChannelConfigStore.contract.Transact(opts, "addChannel", channelId, channelDefinition)
}

func (_ChannelConfigStore *ChannelConfigStoreSession) AddChannel(channelId [32]byte, channelDefinition IChannelConfigStoreChannelDefinition) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.AddChannel(&_ChannelConfigStore.TransactOpts, channelId, channelDefinition)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorSession) AddChannel(channelId [32]byte, channelDefinition IChannelConfigStoreChannelDefinition) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.AddChannel(&_ChannelConfigStore.TransactOpts, channelId, channelDefinition)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactor) PromoteStagingConfig(opts *bind.TransactOpts, channelId [32]byte) (*types.Transaction, error) {
	return _ChannelConfigStore.contract.Transact(opts, "promoteStagingConfig", channelId)
}

func (_ChannelConfigStore *ChannelConfigStoreSession) PromoteStagingConfig(channelId [32]byte) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.PromoteStagingConfig(&_ChannelConfigStore.TransactOpts, channelId)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorSession) PromoteStagingConfig(channelId [32]byte) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.PromoteStagingConfig(&_ChannelConfigStore.TransactOpts, channelId)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactor) RemoveChannel(opts *bind.TransactOpts, channelId [32]byte) (*types.Transaction, error) {
	return _ChannelConfigStore.contract.Transact(opts, "removeChannel", channelId)
}

func (_ChannelConfigStore *ChannelConfigStoreSession) RemoveChannel(channelId [32]byte) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.RemoveChannel(&_ChannelConfigStore.TransactOpts, channelId)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorSession) RemoveChannel(channelId [32]byte) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.RemoveChannel(&_ChannelConfigStore.TransactOpts, channelId)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactor) SetStagingConfig(opts *bind.TransactOpts, channelId [32]byte, channelConfig IChannelConfigStoreChannelConfiguration) (*types.Transaction, error) {
	return _ChannelConfigStore.contract.Transact(opts, "setStagingConfig", channelId, channelConfig)
}

func (_ChannelConfigStore *ChannelConfigStoreSession) SetStagingConfig(channelId [32]byte, channelConfig IChannelConfigStoreChannelConfiguration) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.SetStagingConfig(&_ChannelConfigStore.TransactOpts, channelId, channelConfig)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorSession) SetStagingConfig(channelId [32]byte, channelConfig IChannelConfigStoreChannelConfiguration) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.SetStagingConfig(&_ChannelConfigStore.TransactOpts, channelId, channelConfig)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ChannelConfigStore.contract.Transact(opts, "transferOwnership", to)
}

func (_ChannelConfigStore *ChannelConfigStoreSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.TransferOwnership(&_ChannelConfigStore.TransactOpts, to)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.TransferOwnership(&_ChannelConfigStore.TransactOpts, to)
}

type ChannelConfigStoreChannelDefinitionRemovedIterator struct {
	Event *ChannelConfigStoreChannelDefinitionRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelConfigStoreChannelDefinitionRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelConfigStoreChannelDefinitionRemoved)
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
		it.Event = new(ChannelConfigStoreChannelDefinitionRemoved)
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

func (it *ChannelConfigStoreChannelDefinitionRemovedIterator) Error() error {
	return it.fail
}

func (it *ChannelConfigStoreChannelDefinitionRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelConfigStoreChannelDefinitionRemoved struct {
	ChannelId [32]byte
	Raw       types.Log
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) FilterChannelDefinitionRemoved(opts *bind.FilterOpts) (*ChannelConfigStoreChannelDefinitionRemovedIterator, error) {

	logs, sub, err := _ChannelConfigStore.contract.FilterLogs(opts, "ChannelDefinitionRemoved")
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreChannelDefinitionRemovedIterator{contract: _ChannelConfigStore.contract, event: "ChannelDefinitionRemoved", logs: logs, sub: sub}, nil
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) WatchChannelDefinitionRemoved(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreChannelDefinitionRemoved) (event.Subscription, error) {

	logs, sub, err := _ChannelConfigStore.contract.WatchLogs(opts, "ChannelDefinitionRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelConfigStoreChannelDefinitionRemoved)
				if err := _ChannelConfigStore.contract.UnpackLog(event, "ChannelDefinitionRemoved", log); err != nil {
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

func (_ChannelConfigStore *ChannelConfigStoreFilterer) ParseChannelDefinitionRemoved(log types.Log) (*ChannelConfigStoreChannelDefinitionRemoved, error) {
	event := new(ChannelConfigStoreChannelDefinitionRemoved)
	if err := _ChannelConfigStore.contract.UnpackLog(event, "ChannelDefinitionRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelConfigStoreNewChannelDefinitionIterator struct {
	Event *ChannelConfigStoreNewChannelDefinition

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelConfigStoreNewChannelDefinitionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelConfigStoreNewChannelDefinition)
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
		it.Event = new(ChannelConfigStoreNewChannelDefinition)
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

func (it *ChannelConfigStoreNewChannelDefinitionIterator) Error() error {
	return it.fail
}

func (it *ChannelConfigStoreNewChannelDefinitionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelConfigStoreNewChannelDefinition struct {
	ChannelId         [32]byte
	ChannelDefinition IChannelConfigStoreChannelDefinition
	Raw               types.Log
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) FilterNewChannelDefinition(opts *bind.FilterOpts) (*ChannelConfigStoreNewChannelDefinitionIterator, error) {

	logs, sub, err := _ChannelConfigStore.contract.FilterLogs(opts, "NewChannelDefinition")
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreNewChannelDefinitionIterator{contract: _ChannelConfigStore.contract, event: "NewChannelDefinition", logs: logs, sub: sub}, nil
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) WatchNewChannelDefinition(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreNewChannelDefinition) (event.Subscription, error) {

	logs, sub, err := _ChannelConfigStore.contract.WatchLogs(opts, "NewChannelDefinition")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelConfigStoreNewChannelDefinition)
				if err := _ChannelConfigStore.contract.UnpackLog(event, "NewChannelDefinition", log); err != nil {
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

func (_ChannelConfigStore *ChannelConfigStoreFilterer) ParseNewChannelDefinition(log types.Log) (*ChannelConfigStoreNewChannelDefinition, error) {
	event := new(ChannelConfigStoreNewChannelDefinition)
	if err := _ChannelConfigStore.contract.UnpackLog(event, "NewChannelDefinition", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelConfigStoreNewProductionConfigIterator struct {
	Event *ChannelConfigStoreNewProductionConfig

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelConfigStoreNewProductionConfigIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelConfigStoreNewProductionConfig)
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
		it.Event = new(ChannelConfigStoreNewProductionConfig)
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

func (it *ChannelConfigStoreNewProductionConfigIterator) Error() error {
	return it.fail
}

func (it *ChannelConfigStoreNewProductionConfigIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelConfigStoreNewProductionConfig struct {
	ChannelConfig IChannelConfigStoreChannelConfiguration
	Raw           types.Log
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) FilterNewProductionConfig(opts *bind.FilterOpts) (*ChannelConfigStoreNewProductionConfigIterator, error) {

	logs, sub, err := _ChannelConfigStore.contract.FilterLogs(opts, "NewProductionConfig")
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreNewProductionConfigIterator{contract: _ChannelConfigStore.contract, event: "NewProductionConfig", logs: logs, sub: sub}, nil
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) WatchNewProductionConfig(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreNewProductionConfig) (event.Subscription, error) {

	logs, sub, err := _ChannelConfigStore.contract.WatchLogs(opts, "NewProductionConfig")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelConfigStoreNewProductionConfig)
				if err := _ChannelConfigStore.contract.UnpackLog(event, "NewProductionConfig", log); err != nil {
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

func (_ChannelConfigStore *ChannelConfigStoreFilterer) ParseNewProductionConfig(log types.Log) (*ChannelConfigStoreNewProductionConfig, error) {
	event := new(ChannelConfigStoreNewProductionConfig)
	if err := _ChannelConfigStore.contract.UnpackLog(event, "NewProductionConfig", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelConfigStoreNewStagingConfigIterator struct {
	Event *ChannelConfigStoreNewStagingConfig

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelConfigStoreNewStagingConfigIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelConfigStoreNewStagingConfig)
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
		it.Event = new(ChannelConfigStoreNewStagingConfig)
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

func (it *ChannelConfigStoreNewStagingConfigIterator) Error() error {
	return it.fail
}

func (it *ChannelConfigStoreNewStagingConfigIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelConfigStoreNewStagingConfig struct {
	ChannelConfig IChannelConfigStoreChannelConfiguration
	Raw           types.Log
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) FilterNewStagingConfig(opts *bind.FilterOpts) (*ChannelConfigStoreNewStagingConfigIterator, error) {

	logs, sub, err := _ChannelConfigStore.contract.FilterLogs(opts, "NewStagingConfig")
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreNewStagingConfigIterator{contract: _ChannelConfigStore.contract, event: "NewStagingConfig", logs: logs, sub: sub}, nil
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) WatchNewStagingConfig(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreNewStagingConfig) (event.Subscription, error) {

	logs, sub, err := _ChannelConfigStore.contract.WatchLogs(opts, "NewStagingConfig")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelConfigStoreNewStagingConfig)
				if err := _ChannelConfigStore.contract.UnpackLog(event, "NewStagingConfig", log); err != nil {
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

func (_ChannelConfigStore *ChannelConfigStoreFilterer) ParseNewStagingConfig(log types.Log) (*ChannelConfigStoreNewStagingConfig, error) {
	event := new(ChannelConfigStoreNewStagingConfig)
	if err := _ChannelConfigStore.contract.UnpackLog(event, "NewStagingConfig", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelConfigStoreOwnershipTransferRequestedIterator struct {
	Event *ChannelConfigStoreOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelConfigStoreOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelConfigStoreOwnershipTransferRequested)
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
		it.Event = new(ChannelConfigStoreOwnershipTransferRequested)
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

func (it *ChannelConfigStoreOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *ChannelConfigStoreOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelConfigStoreOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelConfigStoreOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelConfigStore.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreOwnershipTransferRequestedIterator{contract: _ChannelConfigStore.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelConfigStore.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelConfigStoreOwnershipTransferRequested)
				if err := _ChannelConfigStore.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_ChannelConfigStore *ChannelConfigStoreFilterer) ParseOwnershipTransferRequested(log types.Log) (*ChannelConfigStoreOwnershipTransferRequested, error) {
	event := new(ChannelConfigStoreOwnershipTransferRequested)
	if err := _ChannelConfigStore.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelConfigStoreOwnershipTransferredIterator struct {
	Event *ChannelConfigStoreOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelConfigStoreOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelConfigStoreOwnershipTransferred)
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
		it.Event = new(ChannelConfigStoreOwnershipTransferred)
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

func (it *ChannelConfigStoreOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *ChannelConfigStoreOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelConfigStoreOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelConfigStoreOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelConfigStore.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreOwnershipTransferredIterator{contract: _ChannelConfigStore.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelConfigStore.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelConfigStoreOwnershipTransferred)
				if err := _ChannelConfigStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_ChannelConfigStore *ChannelConfigStoreFilterer) ParseOwnershipTransferred(log types.Log) (*ChannelConfigStoreOwnershipTransferred, error) {
	event := new(ChannelConfigStoreOwnershipTransferred)
	if err := _ChannelConfigStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelConfigStorePromoteStagingConfigIterator struct {
	Event *ChannelConfigStorePromoteStagingConfig

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelConfigStorePromoteStagingConfigIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelConfigStorePromoteStagingConfig)
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
		it.Event = new(ChannelConfigStorePromoteStagingConfig)
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

func (it *ChannelConfigStorePromoteStagingConfigIterator) Error() error {
	return it.fail
}

func (it *ChannelConfigStorePromoteStagingConfigIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelConfigStorePromoteStagingConfig struct {
	ChannelId [32]byte
	Raw       types.Log
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) FilterPromoteStagingConfig(opts *bind.FilterOpts) (*ChannelConfigStorePromoteStagingConfigIterator, error) {

	logs, sub, err := _ChannelConfigStore.contract.FilterLogs(opts, "PromoteStagingConfig")
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStorePromoteStagingConfigIterator{contract: _ChannelConfigStore.contract, event: "PromoteStagingConfig", logs: logs, sub: sub}, nil
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) WatchPromoteStagingConfig(opts *bind.WatchOpts, sink chan<- *ChannelConfigStorePromoteStagingConfig) (event.Subscription, error) {

	logs, sub, err := _ChannelConfigStore.contract.WatchLogs(opts, "PromoteStagingConfig")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelConfigStorePromoteStagingConfig)
				if err := _ChannelConfigStore.contract.UnpackLog(event, "PromoteStagingConfig", log); err != nil {
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

func (_ChannelConfigStore *ChannelConfigStoreFilterer) ParsePromoteStagingConfig(log types.Log) (*ChannelConfigStorePromoteStagingConfig, error) {
	event := new(ChannelConfigStorePromoteStagingConfig)
	if err := _ChannelConfigStore.contract.UnpackLog(event, "PromoteStagingConfig", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_ChannelConfigStore *ChannelConfigStore) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ChannelConfigStore.abi.Events["ChannelDefinitionRemoved"].ID:
		return _ChannelConfigStore.ParseChannelDefinitionRemoved(log)
	case _ChannelConfigStore.abi.Events["NewChannelDefinition"].ID:
		return _ChannelConfigStore.ParseNewChannelDefinition(log)
	case _ChannelConfigStore.abi.Events["NewProductionConfig"].ID:
		return _ChannelConfigStore.ParseNewProductionConfig(log)
	case _ChannelConfigStore.abi.Events["NewStagingConfig"].ID:
		return _ChannelConfigStore.ParseNewStagingConfig(log)
	case _ChannelConfigStore.abi.Events["OwnershipTransferRequested"].ID:
		return _ChannelConfigStore.ParseOwnershipTransferRequested(log)
	case _ChannelConfigStore.abi.Events["OwnershipTransferred"].ID:
		return _ChannelConfigStore.ParseOwnershipTransferred(log)
	case _ChannelConfigStore.abi.Events["PromoteStagingConfig"].ID:
		return _ChannelConfigStore.ParsePromoteStagingConfig(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ChannelConfigStoreChannelDefinitionRemoved) Topic() common.Hash {
	return common.HexToHash("0xa65638d745b306456ab0961a502338c1f24d1f962615be3d154c4e86e879fc74")
}

func (ChannelConfigStoreNewChannelDefinition) Topic() common.Hash {
	return common.HexToHash("0x3f9f883dedd481d6ebe8e3efd5ef57f4b0293a3cf1d85946913dd82723542cc6")
}

func (ChannelConfigStoreNewProductionConfig) Topic() common.Hash {
	return common.HexToHash("0xf484d8aa0665a5502456cd66a8bf6268922b4da7dc29f3bb1fcf67a7da444a2a")
}

func (ChannelConfigStoreNewStagingConfig) Topic() common.Hash {
	return common.HexToHash("0x56d7e0e88863044d9b3e139e6e9c18977c6e81c387e44ce9661611f53035c7ed")
}

func (ChannelConfigStoreOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (ChannelConfigStoreOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (ChannelConfigStorePromoteStagingConfig) Topic() common.Hash {
	return common.HexToHash("0xe644aaaa8169119e133c9b338279b4305419a255ace92b4383df2f45f7daa7a8")
}

func (_ChannelConfigStore *ChannelConfigStore) Address() common.Address {
	return _ChannelConfigStore.address
}

type ChannelConfigStoreInterface interface {
	GetChannelDefinitions(opts *bind.CallOpts, channelId [32]byte) (IChannelConfigStoreChannelDefinition, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddChannel(opts *bind.TransactOpts, channelId [32]byte, channelDefinition IChannelConfigStoreChannelDefinition) (*types.Transaction, error)

	PromoteStagingConfig(opts *bind.TransactOpts, channelId [32]byte) (*types.Transaction, error)

	RemoveChannel(opts *bind.TransactOpts, channelId [32]byte) (*types.Transaction, error)

	SetStagingConfig(opts *bind.TransactOpts, channelId [32]byte, channelConfig IChannelConfigStoreChannelConfiguration) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterChannelDefinitionRemoved(opts *bind.FilterOpts) (*ChannelConfigStoreChannelDefinitionRemovedIterator, error)

	WatchChannelDefinitionRemoved(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreChannelDefinitionRemoved) (event.Subscription, error)

	ParseChannelDefinitionRemoved(log types.Log) (*ChannelConfigStoreChannelDefinitionRemoved, error)

	FilterNewChannelDefinition(opts *bind.FilterOpts) (*ChannelConfigStoreNewChannelDefinitionIterator, error)

	WatchNewChannelDefinition(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreNewChannelDefinition) (event.Subscription, error)

	ParseNewChannelDefinition(log types.Log) (*ChannelConfigStoreNewChannelDefinition, error)

	FilterNewProductionConfig(opts *bind.FilterOpts) (*ChannelConfigStoreNewProductionConfigIterator, error)

	WatchNewProductionConfig(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreNewProductionConfig) (event.Subscription, error)

	ParseNewProductionConfig(log types.Log) (*ChannelConfigStoreNewProductionConfig, error)

	FilterNewStagingConfig(opts *bind.FilterOpts) (*ChannelConfigStoreNewStagingConfigIterator, error)

	WatchNewStagingConfig(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreNewStagingConfig) (event.Subscription, error)

	ParseNewStagingConfig(log types.Log) (*ChannelConfigStoreNewStagingConfig, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelConfigStoreOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*ChannelConfigStoreOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelConfigStoreOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*ChannelConfigStoreOwnershipTransferred, error)

	FilterPromoteStagingConfig(opts *bind.FilterOpts) (*ChannelConfigStorePromoteStagingConfigIterator, error)

	WatchPromoteStagingConfig(opts *bind.WatchOpts, sink chan<- *ChannelConfigStorePromoteStagingConfig) (event.Subscription, error)

	ParsePromoteStagingConfig(log types.Log) (*ChannelConfigStorePromoteStagingConfig, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
