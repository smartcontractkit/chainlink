// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stream_config_store

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

type IStreamConfigStoreChannelConfiguration struct {
	ChannelConfigId [32]byte
}

type IStreamConfigStoreChannelDefinition struct {
	ReportFormat  [8]byte
	ChainSelector uint64
	StreamIDs     [][32]byte
}

var StreamConfigStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ChannelDefinitionNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyStreamIDs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByEOA\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StagingConfigAlreadyPromoted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelector\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroReportFormat\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"ChannelDefinitionRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes8\",\"name\":\"reportFormat\",\"type\":\"bytes8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32[]\",\"name\":\"streamIDs\",\"type\":\"bytes32[]\"}],\"indexed\":false,\"internalType\":\"structIStreamConfigStore.ChannelDefinition\",\"name\":\"channelDefinition\",\"type\":\"tuple\"}],\"name\":\"NewChannelDefinition\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelConfigId\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structIStreamConfigStore.ChannelConfiguration\",\"name\":\"channelConfig\",\"type\":\"tuple\"}],\"name\":\"NewProductionConfig\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelConfigId\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structIStreamConfigStore.ChannelConfiguration\",\"name\":\"channelConfig\",\"type\":\"tuple\"}],\"name\":\"NewStagingConfig\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"PromoteStagingConfig\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes8\",\"name\":\"reportFormat\",\"type\":\"bytes8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32[]\",\"name\":\"streamIDs\",\"type\":\"bytes32[]\"}],\"internalType\":\"structIStreamConfigStore.ChannelDefinition\",\"name\":\"channelDefinition\",\"type\":\"tuple\"}],\"name\":\"addChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"getChannelDefinitions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes8\",\"name\":\"reportFormat\",\"type\":\"bytes8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32[]\",\"name\":\"streamIDs\",\"type\":\"bytes32[]\"}],\"internalType\":\"structIStreamConfigStore.ChannelDefinition\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"promoteStagingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"removeChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelConfigId\",\"type\":\"bytes32\"}],\"internalType\":\"structIStreamConfigStore.ChannelConfiguration\",\"name\":\"channelConfig\",\"type\":\"tuple\"}],\"name\":\"setStagingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610eef806101576000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c806379ba5097116100765780639a223b8e1161005b5780639a223b8e146101d7578063dd2ef3ef146101ea578063f2fde38b1461020a57600080fd5b806379ba5097146101a75780638da5cb5b146101af57600080fd5b8063181f5a77116100a7578063181f5a7714610142578063414aa93414610181578063567ed3371461019457600080fd5b806301ffc9a7146100c357806302da0e511461012d575b600080fd5b6101186100d13660046108c8565b7fffffffff00000000000000000000000000000000000000000000000000000000167f52e2bc33000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b61014061013b366004610911565b61021d565b005b604080518082018252601781527f53747265616d436f6e66696753746f726520302e302e3000000000000000000060208201529051610124919061092a565b61014061018f366004610996565b6102e9565b6101406101a2366004610911565b610340565b6101406103a0565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610124565b6101406101e53660046109ec565b6104a2565b6101fd6101f8366004610911565b6105fb565b6040516101249190610a3a565b610140610218366004610acc565b61070a565b61022561071e565b600081815260026020526040812060010154900361026f576040517fd1a751e200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081815260026020526040812080547fffffffffffffffffffffffffffffffff00000000000000000000000000000000168155906102b16001830182610896565b50506040518181527fa65638d745b306456ab0961a502338c1f24d1f962615be3d154c4e86e879fc749060200160405180910390a150565b6102f161071e565b60008281526004602052604090208135815581905050604051813581527f56d7e0e88863044d9b3e139e6e9c18977c6e81c387e44ce9661611f53035c7ed906020015b60405180910390a15050565b61034861071e565b6000818152600460209081526040808320815180840183529054808252858552600384529382902093909355518381527fe644aaaa8169119e133c9b338279b4305419a255ace92b4383df2f45f7daa7a89101610334565b60015473ffffffffffffffffffffffffffffffffffffffff163314610426576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6104aa61071e565b6104b76040820182610b02565b90506000036104f2576040517f4b620e2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105026040820160208301610b87565b67ffffffffffffffff16600003610545576040517ff89d762900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105526020820182610bd2565b7fffffffffffffffff000000000000000000000000000000000000000000000000166000036105ad576040517febd3ef0200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260026020526040902081906105c78282610cb7565b9050507f3f9f883dedd481d6ebe8e3efd5ef57f4b0293a3cf1d85946913dd82723542cc68282604051610334929190610e05565b60408051606080820183526000808352602083015291810191909152333214610650576040517f74e2cd5100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000828152600260209081526040918290208251606081018452815460c081901b7fffffffffffffffff00000000000000000000000000000000000000000000000016825268010000000000000000900467ffffffffffffffff1681840152600182018054855181860281018601875281815292959394938601938301828280156106fa57602002820191906000526020600020905b8154815260200190600101908083116106e6575b5050505050815250509050919050565b61071261071e565b61071b816107a1565b50565b60005473ffffffffffffffffffffffffffffffffffffffff16331461079f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161041d565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610820576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161041d565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b508054600082559060005260206000209081019061071b91905b808211156108c457600081556001016108b0565b5090565b6000602082840312156108da57600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461090a57600080fd5b9392505050565b60006020828403121561092357600080fd5b5035919050565b600060208083528351808285015260005b818110156109575785810183015185820160400152820161093b565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60008082840360408112156109aa57600080fd5b8335925060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0820112156109de57600080fd5b506020830190509250929050565b600080604083850312156109ff57600080fd5b82359150602083013567ffffffffffffffff811115610a1d57600080fd5b830160608186031215610a2f57600080fd5b809150509250929050565b60006020808352608083017fffffffffffffffff0000000000000000000000000000000000000000000000008551168285015267ffffffffffffffff82860151166040850152604085015160608086015281815180845260a0870191508483019350600092505b80831015610ac15783518252928401926001929092019190840190610aa1565b509695505050505050565b600060208284031215610ade57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461090a57600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610b3757600080fd5b83018035915067ffffffffffffffff821115610b5257600080fd5b6020019150600581901b3603821315610b6a57600080fd5b9250929050565b67ffffffffffffffff8116811461071b57600080fd5b600060208284031215610b9957600080fd5b813561090a81610b71565b7fffffffffffffffff0000000000000000000000000000000000000000000000008116811461071b57600080fd5b600060208284031215610be457600080fd5b813561090a81610ba4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff831115610c3657610c36610bef565b68010000000000000000831115610c4f57610c4f610bef565b805483825580841015610c86576000828152602081208581019083015b80821015610c8257828255600182019150610c6c565b5050505b5060008181526020812083915b85811015610caf57823582820155602090920191600101610c93565b505050505050565b8135610cc281610ba4565b8060c01c90508154817fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000082161783556020840135610cff81610b71565b6fffffffffffffffff00000000000000008160401b16837fffffffffffffffffffffffffffffffff0000000000000000000000000000000084161717845550505060408201357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1833603018112610d7557600080fd5b8201803567ffffffffffffffff811115610d8e57600080fd5b6020820191508060051b3603821315610da657600080fd5b610db4818360018601610c1e565b50505050565b81835260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115610dec57600080fd5b8260051b80836020870137939093016020019392505050565b8281526040602082015260008235610e1c81610ba4565b7fffffffffffffffff0000000000000000000000000000000000000000000000001660408301526020830135610e5181610b71565b67ffffffffffffffff8082166060850152604085013591507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1853603018212610e9957600080fd5b6020918501918201913581811115610eb057600080fd5b8060051b3603831315610ec257600080fd5b60606080860152610ed760a086018285610dba565b97965050505050505056fea164736f6c6343000810000a",
}

var StreamConfigStoreABI = StreamConfigStoreMetaData.ABI

var StreamConfigStoreBin = StreamConfigStoreMetaData.Bin

func DeployStreamConfigStore(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StreamConfigStore, error) {
	parsed, err := StreamConfigStoreMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StreamConfigStoreBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StreamConfigStore{address: address, abi: *parsed, StreamConfigStoreCaller: StreamConfigStoreCaller{contract: contract}, StreamConfigStoreTransactor: StreamConfigStoreTransactor{contract: contract}, StreamConfigStoreFilterer: StreamConfigStoreFilterer{contract: contract}}, nil
}

type StreamConfigStore struct {
	address common.Address
	abi     abi.ABI
	StreamConfigStoreCaller
	StreamConfigStoreTransactor
	StreamConfigStoreFilterer
}

type StreamConfigStoreCaller struct {
	contract *bind.BoundContract
}

type StreamConfigStoreTransactor struct {
	contract *bind.BoundContract
}

type StreamConfigStoreFilterer struct {
	contract *bind.BoundContract
}

type StreamConfigStoreSession struct {
	Contract     *StreamConfigStore
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type StreamConfigStoreCallerSession struct {
	Contract *StreamConfigStoreCaller
	CallOpts bind.CallOpts
}

type StreamConfigStoreTransactorSession struct {
	Contract     *StreamConfigStoreTransactor
	TransactOpts bind.TransactOpts
}

type StreamConfigStoreRaw struct {
	Contract *StreamConfigStore
}

type StreamConfigStoreCallerRaw struct {
	Contract *StreamConfigStoreCaller
}

type StreamConfigStoreTransactorRaw struct {
	Contract *StreamConfigStoreTransactor
}

func NewStreamConfigStore(address common.Address, backend bind.ContractBackend) (*StreamConfigStore, error) {
	abi, err := abi.JSON(strings.NewReader(StreamConfigStoreABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindStreamConfigStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StreamConfigStore{address: address, abi: abi, StreamConfigStoreCaller: StreamConfigStoreCaller{contract: contract}, StreamConfigStoreTransactor: StreamConfigStoreTransactor{contract: contract}, StreamConfigStoreFilterer: StreamConfigStoreFilterer{contract: contract}}, nil
}

func NewStreamConfigStoreCaller(address common.Address, caller bind.ContractCaller) (*StreamConfigStoreCaller, error) {
	contract, err := bindStreamConfigStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StreamConfigStoreCaller{contract: contract}, nil
}

func NewStreamConfigStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*StreamConfigStoreTransactor, error) {
	contract, err := bindStreamConfigStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StreamConfigStoreTransactor{contract: contract}, nil
}

func NewStreamConfigStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*StreamConfigStoreFilterer, error) {
	contract, err := bindStreamConfigStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StreamConfigStoreFilterer{contract: contract}, nil
}

func bindStreamConfigStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StreamConfigStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_StreamConfigStore *StreamConfigStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StreamConfigStore.Contract.StreamConfigStoreCaller.contract.Call(opts, result, method, params...)
}

func (_StreamConfigStore *StreamConfigStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.StreamConfigStoreTransactor.contract.Transfer(opts)
}

func (_StreamConfigStore *StreamConfigStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.StreamConfigStoreTransactor.contract.Transact(opts, method, params...)
}

func (_StreamConfigStore *StreamConfigStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StreamConfigStore.Contract.contract.Call(opts, result, method, params...)
}

func (_StreamConfigStore *StreamConfigStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.contract.Transfer(opts)
}

func (_StreamConfigStore *StreamConfigStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.contract.Transact(opts, method, params...)
}

func (_StreamConfigStore *StreamConfigStoreCaller) GetChannelDefinitions(opts *bind.CallOpts, channelId [32]byte) (IStreamConfigStoreChannelDefinition, error) {
	var out []interface{}
	err := _StreamConfigStore.contract.Call(opts, &out, "getChannelDefinitions", channelId)

	if err != nil {
		return *new(IStreamConfigStoreChannelDefinition), err
	}

	out0 := *abi.ConvertType(out[0], new(IStreamConfigStoreChannelDefinition)).(*IStreamConfigStoreChannelDefinition)

	return out0, err

}

func (_StreamConfigStore *StreamConfigStoreSession) GetChannelDefinitions(channelId [32]byte) (IStreamConfigStoreChannelDefinition, error) {
	return _StreamConfigStore.Contract.GetChannelDefinitions(&_StreamConfigStore.CallOpts, channelId)
}

func (_StreamConfigStore *StreamConfigStoreCallerSession) GetChannelDefinitions(channelId [32]byte) (IStreamConfigStoreChannelDefinition, error) {
	return _StreamConfigStore.Contract.GetChannelDefinitions(&_StreamConfigStore.CallOpts, channelId)
}

func (_StreamConfigStore *StreamConfigStoreCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StreamConfigStore.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_StreamConfigStore *StreamConfigStoreSession) Owner() (common.Address, error) {
	return _StreamConfigStore.Contract.Owner(&_StreamConfigStore.CallOpts)
}

func (_StreamConfigStore *StreamConfigStoreCallerSession) Owner() (common.Address, error) {
	return _StreamConfigStore.Contract.Owner(&_StreamConfigStore.CallOpts)
}

func (_StreamConfigStore *StreamConfigStoreCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _StreamConfigStore.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_StreamConfigStore *StreamConfigStoreSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _StreamConfigStore.Contract.SupportsInterface(&_StreamConfigStore.CallOpts, interfaceId)
}

func (_StreamConfigStore *StreamConfigStoreCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _StreamConfigStore.Contract.SupportsInterface(&_StreamConfigStore.CallOpts, interfaceId)
}

func (_StreamConfigStore *StreamConfigStoreCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StreamConfigStore.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_StreamConfigStore *StreamConfigStoreSession) TypeAndVersion() (string, error) {
	return _StreamConfigStore.Contract.TypeAndVersion(&_StreamConfigStore.CallOpts)
}

func (_StreamConfigStore *StreamConfigStoreCallerSession) TypeAndVersion() (string, error) {
	return _StreamConfigStore.Contract.TypeAndVersion(&_StreamConfigStore.CallOpts)
}

func (_StreamConfigStore *StreamConfigStoreTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StreamConfigStore.contract.Transact(opts, "acceptOwnership")
}

func (_StreamConfigStore *StreamConfigStoreSession) AcceptOwnership() (*types.Transaction, error) {
	return _StreamConfigStore.Contract.AcceptOwnership(&_StreamConfigStore.TransactOpts)
}

func (_StreamConfigStore *StreamConfigStoreTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _StreamConfigStore.Contract.AcceptOwnership(&_StreamConfigStore.TransactOpts)
}

func (_StreamConfigStore *StreamConfigStoreTransactor) AddChannel(opts *bind.TransactOpts, channelId [32]byte, channelDefinition IStreamConfigStoreChannelDefinition) (*types.Transaction, error) {
	return _StreamConfigStore.contract.Transact(opts, "addChannel", channelId, channelDefinition)
}

func (_StreamConfigStore *StreamConfigStoreSession) AddChannel(channelId [32]byte, channelDefinition IStreamConfigStoreChannelDefinition) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.AddChannel(&_StreamConfigStore.TransactOpts, channelId, channelDefinition)
}

func (_StreamConfigStore *StreamConfigStoreTransactorSession) AddChannel(channelId [32]byte, channelDefinition IStreamConfigStoreChannelDefinition) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.AddChannel(&_StreamConfigStore.TransactOpts, channelId, channelDefinition)
}

func (_StreamConfigStore *StreamConfigStoreTransactor) PromoteStagingConfig(opts *bind.TransactOpts, channelId [32]byte) (*types.Transaction, error) {
	return _StreamConfigStore.contract.Transact(opts, "promoteStagingConfig", channelId)
}

func (_StreamConfigStore *StreamConfigStoreSession) PromoteStagingConfig(channelId [32]byte) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.PromoteStagingConfig(&_StreamConfigStore.TransactOpts, channelId)
}

func (_StreamConfigStore *StreamConfigStoreTransactorSession) PromoteStagingConfig(channelId [32]byte) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.PromoteStagingConfig(&_StreamConfigStore.TransactOpts, channelId)
}

func (_StreamConfigStore *StreamConfigStoreTransactor) RemoveChannel(opts *bind.TransactOpts, channelId [32]byte) (*types.Transaction, error) {
	return _StreamConfigStore.contract.Transact(opts, "removeChannel", channelId)
}

func (_StreamConfigStore *StreamConfigStoreSession) RemoveChannel(channelId [32]byte) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.RemoveChannel(&_StreamConfigStore.TransactOpts, channelId)
}

func (_StreamConfigStore *StreamConfigStoreTransactorSession) RemoveChannel(channelId [32]byte) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.RemoveChannel(&_StreamConfigStore.TransactOpts, channelId)
}

func (_StreamConfigStore *StreamConfigStoreTransactor) SetStagingConfig(opts *bind.TransactOpts, channelId [32]byte, channelConfig IStreamConfigStoreChannelConfiguration) (*types.Transaction, error) {
	return _StreamConfigStore.contract.Transact(opts, "setStagingConfig", channelId, channelConfig)
}

func (_StreamConfigStore *StreamConfigStoreSession) SetStagingConfig(channelId [32]byte, channelConfig IStreamConfigStoreChannelConfiguration) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.SetStagingConfig(&_StreamConfigStore.TransactOpts, channelId, channelConfig)
}

func (_StreamConfigStore *StreamConfigStoreTransactorSession) SetStagingConfig(channelId [32]byte, channelConfig IStreamConfigStoreChannelConfiguration) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.SetStagingConfig(&_StreamConfigStore.TransactOpts, channelId, channelConfig)
}

func (_StreamConfigStore *StreamConfigStoreTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _StreamConfigStore.contract.Transact(opts, "transferOwnership", to)
}

func (_StreamConfigStore *StreamConfigStoreSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.TransferOwnership(&_StreamConfigStore.TransactOpts, to)
}

func (_StreamConfigStore *StreamConfigStoreTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _StreamConfigStore.Contract.TransferOwnership(&_StreamConfigStore.TransactOpts, to)
}

type StreamConfigStoreChannelDefinitionRemovedIterator struct {
	Event *StreamConfigStoreChannelDefinitionRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StreamConfigStoreChannelDefinitionRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StreamConfigStoreChannelDefinitionRemoved)
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
		it.Event = new(StreamConfigStoreChannelDefinitionRemoved)
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

func (it *StreamConfigStoreChannelDefinitionRemovedIterator) Error() error {
	return it.fail
}

func (it *StreamConfigStoreChannelDefinitionRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StreamConfigStoreChannelDefinitionRemoved struct {
	ChannelId [32]byte
	Raw       types.Log
}

func (_StreamConfigStore *StreamConfigStoreFilterer) FilterChannelDefinitionRemoved(opts *bind.FilterOpts) (*StreamConfigStoreChannelDefinitionRemovedIterator, error) {

	logs, sub, err := _StreamConfigStore.contract.FilterLogs(opts, "ChannelDefinitionRemoved")
	if err != nil {
		return nil, err
	}
	return &StreamConfigStoreChannelDefinitionRemovedIterator{contract: _StreamConfigStore.contract, event: "ChannelDefinitionRemoved", logs: logs, sub: sub}, nil
}

func (_StreamConfigStore *StreamConfigStoreFilterer) WatchChannelDefinitionRemoved(opts *bind.WatchOpts, sink chan<- *StreamConfigStoreChannelDefinitionRemoved) (event.Subscription, error) {

	logs, sub, err := _StreamConfigStore.contract.WatchLogs(opts, "ChannelDefinitionRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StreamConfigStoreChannelDefinitionRemoved)
				if err := _StreamConfigStore.contract.UnpackLog(event, "ChannelDefinitionRemoved", log); err != nil {
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

func (_StreamConfigStore *StreamConfigStoreFilterer) ParseChannelDefinitionRemoved(log types.Log) (*StreamConfigStoreChannelDefinitionRemoved, error) {
	event := new(StreamConfigStoreChannelDefinitionRemoved)
	if err := _StreamConfigStore.contract.UnpackLog(event, "ChannelDefinitionRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StreamConfigStoreNewChannelDefinitionIterator struct {
	Event *StreamConfigStoreNewChannelDefinition

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StreamConfigStoreNewChannelDefinitionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StreamConfigStoreNewChannelDefinition)
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
		it.Event = new(StreamConfigStoreNewChannelDefinition)
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

func (it *StreamConfigStoreNewChannelDefinitionIterator) Error() error {
	return it.fail
}

func (it *StreamConfigStoreNewChannelDefinitionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StreamConfigStoreNewChannelDefinition struct {
	ChannelId         [32]byte
	ChannelDefinition IStreamConfigStoreChannelDefinition
	Raw               types.Log
}

func (_StreamConfigStore *StreamConfigStoreFilterer) FilterNewChannelDefinition(opts *bind.FilterOpts) (*StreamConfigStoreNewChannelDefinitionIterator, error) {

	logs, sub, err := _StreamConfigStore.contract.FilterLogs(opts, "NewChannelDefinition")
	if err != nil {
		return nil, err
	}
	return &StreamConfigStoreNewChannelDefinitionIterator{contract: _StreamConfigStore.contract, event: "NewChannelDefinition", logs: logs, sub: sub}, nil
}

func (_StreamConfigStore *StreamConfigStoreFilterer) WatchNewChannelDefinition(opts *bind.WatchOpts, sink chan<- *StreamConfigStoreNewChannelDefinition) (event.Subscription, error) {

	logs, sub, err := _StreamConfigStore.contract.WatchLogs(opts, "NewChannelDefinition")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StreamConfigStoreNewChannelDefinition)
				if err := _StreamConfigStore.contract.UnpackLog(event, "NewChannelDefinition", log); err != nil {
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

func (_StreamConfigStore *StreamConfigStoreFilterer) ParseNewChannelDefinition(log types.Log) (*StreamConfigStoreNewChannelDefinition, error) {
	event := new(StreamConfigStoreNewChannelDefinition)
	if err := _StreamConfigStore.contract.UnpackLog(event, "NewChannelDefinition", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StreamConfigStoreNewProductionConfigIterator struct {
	Event *StreamConfigStoreNewProductionConfig

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StreamConfigStoreNewProductionConfigIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StreamConfigStoreNewProductionConfig)
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
		it.Event = new(StreamConfigStoreNewProductionConfig)
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

func (it *StreamConfigStoreNewProductionConfigIterator) Error() error {
	return it.fail
}

func (it *StreamConfigStoreNewProductionConfigIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StreamConfigStoreNewProductionConfig struct {
	ChannelConfig IStreamConfigStoreChannelConfiguration
	Raw           types.Log
}

func (_StreamConfigStore *StreamConfigStoreFilterer) FilterNewProductionConfig(opts *bind.FilterOpts) (*StreamConfigStoreNewProductionConfigIterator, error) {

	logs, sub, err := _StreamConfigStore.contract.FilterLogs(opts, "NewProductionConfig")
	if err != nil {
		return nil, err
	}
	return &StreamConfigStoreNewProductionConfigIterator{contract: _StreamConfigStore.contract, event: "NewProductionConfig", logs: logs, sub: sub}, nil
}

func (_StreamConfigStore *StreamConfigStoreFilterer) WatchNewProductionConfig(opts *bind.WatchOpts, sink chan<- *StreamConfigStoreNewProductionConfig) (event.Subscription, error) {

	logs, sub, err := _StreamConfigStore.contract.WatchLogs(opts, "NewProductionConfig")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StreamConfigStoreNewProductionConfig)
				if err := _StreamConfigStore.contract.UnpackLog(event, "NewProductionConfig", log); err != nil {
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

func (_StreamConfigStore *StreamConfigStoreFilterer) ParseNewProductionConfig(log types.Log) (*StreamConfigStoreNewProductionConfig, error) {
	event := new(StreamConfigStoreNewProductionConfig)
	if err := _StreamConfigStore.contract.UnpackLog(event, "NewProductionConfig", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StreamConfigStoreNewStagingConfigIterator struct {
	Event *StreamConfigStoreNewStagingConfig

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StreamConfigStoreNewStagingConfigIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StreamConfigStoreNewStagingConfig)
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
		it.Event = new(StreamConfigStoreNewStagingConfig)
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

func (it *StreamConfigStoreNewStagingConfigIterator) Error() error {
	return it.fail
}

func (it *StreamConfigStoreNewStagingConfigIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StreamConfigStoreNewStagingConfig struct {
	ChannelConfig IStreamConfigStoreChannelConfiguration
	Raw           types.Log
}

func (_StreamConfigStore *StreamConfigStoreFilterer) FilterNewStagingConfig(opts *bind.FilterOpts) (*StreamConfigStoreNewStagingConfigIterator, error) {

	logs, sub, err := _StreamConfigStore.contract.FilterLogs(opts, "NewStagingConfig")
	if err != nil {
		return nil, err
	}
	return &StreamConfigStoreNewStagingConfigIterator{contract: _StreamConfigStore.contract, event: "NewStagingConfig", logs: logs, sub: sub}, nil
}

func (_StreamConfigStore *StreamConfigStoreFilterer) WatchNewStagingConfig(opts *bind.WatchOpts, sink chan<- *StreamConfigStoreNewStagingConfig) (event.Subscription, error) {

	logs, sub, err := _StreamConfigStore.contract.WatchLogs(opts, "NewStagingConfig")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StreamConfigStoreNewStagingConfig)
				if err := _StreamConfigStore.contract.UnpackLog(event, "NewStagingConfig", log); err != nil {
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

func (_StreamConfigStore *StreamConfigStoreFilterer) ParseNewStagingConfig(log types.Log) (*StreamConfigStoreNewStagingConfig, error) {
	event := new(StreamConfigStoreNewStagingConfig)
	if err := _StreamConfigStore.contract.UnpackLog(event, "NewStagingConfig", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StreamConfigStoreOwnershipTransferRequestedIterator struct {
	Event *StreamConfigStoreOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StreamConfigStoreOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StreamConfigStoreOwnershipTransferRequested)
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
		it.Event = new(StreamConfigStoreOwnershipTransferRequested)
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

func (it *StreamConfigStoreOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *StreamConfigStoreOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StreamConfigStoreOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_StreamConfigStore *StreamConfigStoreFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StreamConfigStoreOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StreamConfigStore.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StreamConfigStoreOwnershipTransferRequestedIterator{contract: _StreamConfigStore.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_StreamConfigStore *StreamConfigStoreFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *StreamConfigStoreOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StreamConfigStore.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StreamConfigStoreOwnershipTransferRequested)
				if err := _StreamConfigStore.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_StreamConfigStore *StreamConfigStoreFilterer) ParseOwnershipTransferRequested(log types.Log) (*StreamConfigStoreOwnershipTransferRequested, error) {
	event := new(StreamConfigStoreOwnershipTransferRequested)
	if err := _StreamConfigStore.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StreamConfigStoreOwnershipTransferredIterator struct {
	Event *StreamConfigStoreOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StreamConfigStoreOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StreamConfigStoreOwnershipTransferred)
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
		it.Event = new(StreamConfigStoreOwnershipTransferred)
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

func (it *StreamConfigStoreOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *StreamConfigStoreOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StreamConfigStoreOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_StreamConfigStore *StreamConfigStoreFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StreamConfigStoreOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StreamConfigStore.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StreamConfigStoreOwnershipTransferredIterator{contract: _StreamConfigStore.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_StreamConfigStore *StreamConfigStoreFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StreamConfigStoreOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StreamConfigStore.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StreamConfigStoreOwnershipTransferred)
				if err := _StreamConfigStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_StreamConfigStore *StreamConfigStoreFilterer) ParseOwnershipTransferred(log types.Log) (*StreamConfigStoreOwnershipTransferred, error) {
	event := new(StreamConfigStoreOwnershipTransferred)
	if err := _StreamConfigStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StreamConfigStorePromoteStagingConfigIterator struct {
	Event *StreamConfigStorePromoteStagingConfig

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StreamConfigStorePromoteStagingConfigIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StreamConfigStorePromoteStagingConfig)
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
		it.Event = new(StreamConfigStorePromoteStagingConfig)
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

func (it *StreamConfigStorePromoteStagingConfigIterator) Error() error {
	return it.fail
}

func (it *StreamConfigStorePromoteStagingConfigIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StreamConfigStorePromoteStagingConfig struct {
	ChannelId [32]byte
	Raw       types.Log
}

func (_StreamConfigStore *StreamConfigStoreFilterer) FilterPromoteStagingConfig(opts *bind.FilterOpts) (*StreamConfigStorePromoteStagingConfigIterator, error) {

	logs, sub, err := _StreamConfigStore.contract.FilterLogs(opts, "PromoteStagingConfig")
	if err != nil {
		return nil, err
	}
	return &StreamConfigStorePromoteStagingConfigIterator{contract: _StreamConfigStore.contract, event: "PromoteStagingConfig", logs: logs, sub: sub}, nil
}

func (_StreamConfigStore *StreamConfigStoreFilterer) WatchPromoteStagingConfig(opts *bind.WatchOpts, sink chan<- *StreamConfigStorePromoteStagingConfig) (event.Subscription, error) {

	logs, sub, err := _StreamConfigStore.contract.WatchLogs(opts, "PromoteStagingConfig")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StreamConfigStorePromoteStagingConfig)
				if err := _StreamConfigStore.contract.UnpackLog(event, "PromoteStagingConfig", log); err != nil {
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

func (_StreamConfigStore *StreamConfigStoreFilterer) ParsePromoteStagingConfig(log types.Log) (*StreamConfigStorePromoteStagingConfig, error) {
	event := new(StreamConfigStorePromoteStagingConfig)
	if err := _StreamConfigStore.contract.UnpackLog(event, "PromoteStagingConfig", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_StreamConfigStore *StreamConfigStore) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _StreamConfigStore.abi.Events["ChannelDefinitionRemoved"].ID:
		return _StreamConfigStore.ParseChannelDefinitionRemoved(log)
	case _StreamConfigStore.abi.Events["NewChannelDefinition"].ID:
		return _StreamConfigStore.ParseNewChannelDefinition(log)
	case _StreamConfigStore.abi.Events["NewProductionConfig"].ID:
		return _StreamConfigStore.ParseNewProductionConfig(log)
	case _StreamConfigStore.abi.Events["NewStagingConfig"].ID:
		return _StreamConfigStore.ParseNewStagingConfig(log)
	case _StreamConfigStore.abi.Events["OwnershipTransferRequested"].ID:
		return _StreamConfigStore.ParseOwnershipTransferRequested(log)
	case _StreamConfigStore.abi.Events["OwnershipTransferred"].ID:
		return _StreamConfigStore.ParseOwnershipTransferred(log)
	case _StreamConfigStore.abi.Events["PromoteStagingConfig"].ID:
		return _StreamConfigStore.ParsePromoteStagingConfig(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (StreamConfigStoreChannelDefinitionRemoved) Topic() common.Hash {
	return common.HexToHash("0xa65638d745b306456ab0961a502338c1f24d1f962615be3d154c4e86e879fc74")
}

func (StreamConfigStoreNewChannelDefinition) Topic() common.Hash {
	return common.HexToHash("0x3f9f883dedd481d6ebe8e3efd5ef57f4b0293a3cf1d85946913dd82723542cc6")
}

func (StreamConfigStoreNewProductionConfig) Topic() common.Hash {
	return common.HexToHash("0xf484d8aa0665a5502456cd66a8bf6268922b4da7dc29f3bb1fcf67a7da444a2a")
}

func (StreamConfigStoreNewStagingConfig) Topic() common.Hash {
	return common.HexToHash("0x56d7e0e88863044d9b3e139e6e9c18977c6e81c387e44ce9661611f53035c7ed")
}

func (StreamConfigStoreOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (StreamConfigStoreOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (StreamConfigStorePromoteStagingConfig) Topic() common.Hash {
	return common.HexToHash("0xe644aaaa8169119e133c9b338279b4305419a255ace92b4383df2f45f7daa7a8")
}

func (_StreamConfigStore *StreamConfigStore) Address() common.Address {
	return _StreamConfigStore.address
}

type StreamConfigStoreInterface interface {
	GetChannelDefinitions(opts *bind.CallOpts, channelId [32]byte) (IStreamConfigStoreChannelDefinition, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddChannel(opts *bind.TransactOpts, channelId [32]byte, channelDefinition IStreamConfigStoreChannelDefinition) (*types.Transaction, error)

	PromoteStagingConfig(opts *bind.TransactOpts, channelId [32]byte) (*types.Transaction, error)

	RemoveChannel(opts *bind.TransactOpts, channelId [32]byte) (*types.Transaction, error)

	SetStagingConfig(opts *bind.TransactOpts, channelId [32]byte, channelConfig IStreamConfigStoreChannelConfiguration) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterChannelDefinitionRemoved(opts *bind.FilterOpts) (*StreamConfigStoreChannelDefinitionRemovedIterator, error)

	WatchChannelDefinitionRemoved(opts *bind.WatchOpts, sink chan<- *StreamConfigStoreChannelDefinitionRemoved) (event.Subscription, error)

	ParseChannelDefinitionRemoved(log types.Log) (*StreamConfigStoreChannelDefinitionRemoved, error)

	FilterNewChannelDefinition(opts *bind.FilterOpts) (*StreamConfigStoreNewChannelDefinitionIterator, error)

	WatchNewChannelDefinition(opts *bind.WatchOpts, sink chan<- *StreamConfigStoreNewChannelDefinition) (event.Subscription, error)

	ParseNewChannelDefinition(log types.Log) (*StreamConfigStoreNewChannelDefinition, error)

	FilterNewProductionConfig(opts *bind.FilterOpts) (*StreamConfigStoreNewProductionConfigIterator, error)

	WatchNewProductionConfig(opts *bind.WatchOpts, sink chan<- *StreamConfigStoreNewProductionConfig) (event.Subscription, error)

	ParseNewProductionConfig(log types.Log) (*StreamConfigStoreNewProductionConfig, error)

	FilterNewStagingConfig(opts *bind.FilterOpts) (*StreamConfigStoreNewStagingConfigIterator, error)

	WatchNewStagingConfig(opts *bind.WatchOpts, sink chan<- *StreamConfigStoreNewStagingConfig) (event.Subscription, error)

	ParseNewStagingConfig(log types.Log) (*StreamConfigStoreNewStagingConfig, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StreamConfigStoreOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *StreamConfigStoreOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*StreamConfigStoreOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StreamConfigStoreOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StreamConfigStoreOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*StreamConfigStoreOwnershipTransferred, error)

	FilterPromoteStagingConfig(opts *bind.FilterOpts) (*StreamConfigStorePromoteStagingConfigIterator, error)

	WatchPromoteStagingConfig(opts *bind.WatchOpts, sink chan<- *StreamConfigStorePromoteStagingConfig) (event.Subscription, error)

	ParsePromoteStagingConfig(log types.Log) (*StreamConfigStorePromoteStagingConfig, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
