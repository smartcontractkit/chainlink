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

type IChannelConfigStoreChannelDefinition struct {
	ReportFormat  [8]byte
	ChainSelector uint64
	StreamIDs     []uint32
}

var ChannelConfigStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ChannelDefinitionNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyStreamIDs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByEOA\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StagingConfigAlreadyPromoted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelector\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroReportFormat\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"channelId\",\"type\":\"uint32\"}],\"name\":\"ChannelDefinitionRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"channelId\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"bytes8\",\"name\":\"reportFormat\",\"type\":\"bytes8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint32[]\",\"name\":\"streamIDs\",\"type\":\"uint32[]\"}],\"indexed\":false,\"internalType\":\"structIChannelConfigStore.ChannelDefinition\",\"name\":\"channelDefinition\",\"type\":\"tuple\"}],\"name\":\"NewChannelDefinition\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"channelId\",\"type\":\"uint32\"}],\"name\":\"PromoteStagingConfig\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"channelId\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"bytes8\",\"name\":\"reportFormat\",\"type\":\"bytes8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint32[]\",\"name\":\"streamIDs\",\"type\":\"uint32[]\"}],\"internalType\":\"structIChannelConfigStore.ChannelDefinition\",\"name\":\"channelDefinition\",\"type\":\"tuple\"}],\"name\":\"addChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"channelId\",\"type\":\"uint32\"}],\"name\":\"getChannelDefinitions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes8\",\"name\":\"reportFormat\",\"type\":\"bytes8\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint32[]\",\"name\":\"streamIDs\",\"type\":\"uint32[]\"}],\"internalType\":\"structIChannelConfigStore.ChannelDefinition\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"channelId\",\"type\":\"uint32\"}],\"name\":\"removeChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610f1e806101576000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80637e37e7191161005b5780637e37e719146101535780638da5cb5b14610166578063f2fde38b1461018e578063f5810719146101a157600080fd5b806301ffc9a71461008d578063181f5a77146100f757806322d9780c1461013657806379ba50971461014b575b600080fd5b6100e261009b366004610816565b7fffffffff00000000000000000000000000000000000000000000000000000000167fa96f980c000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b604080518082018252601881527f4368616e6e656c436f6e66696753746f726520302e302e300000000000000000602082015290516100ee919061085f565b6101496101443660046108dd565b6101c1565b005b61014961032d565b610149610161366004610934565b61042f565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100ee565b61014961019c366004610951565b61050f565b6101b46101af366004610934565b610523565b6040516100ee9190610987565b6101c9610665565b6101d66040820182610a1f565b9050600003610211576040517f4b620e2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6102216040820160208301610aa4565b67ffffffffffffffff16600003610264576040517ff89d762900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6102716020820182610aef565b7fffffffffffffffff000000000000000000000000000000000000000000000000166000036102cc576040517febd3ef0200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff8216600090815260026020526040902081906102ed8282610ce5565b9050507fbf2cd44714205d633d3f888ac72ea66d53cd12d4c4e8723a80d9c0bc36484a548282604051610321929190610e2e565b60405180910390a15050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146103b3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610437610665565b63ffffffff81166000908152600260205260408120600101549003610488576040517fd1a751e200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff8116600090815260026020526040812080547fffffffffffffffffffffffffffffffff00000000000000000000000000000000168155906104d160018301826107dd565b505060405163ffffffff821681527f334e877e9691ecae0660510061973bebaa8b4fb37332ed6090052e630c9798619060200160405180910390a150565b610517610665565b610520816106e8565b50565b60408051606080820183526000808352602083015291810191909152333214610578576040517f74e2cd5100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff82166000908152600260209081526040918290208251606081018452815460c081901b7fffffffffffffffff00000000000000000000000000000000000000000000000016825268010000000000000000900467ffffffffffffffff16818401526001820180548551818602810186018752818152929593949386019383018282801561065557602002820191906000526020600020906000905b82829054906101000a900463ffffffff1663ffffffff16815260200190600401906020826003010492830192600103820291508084116106185790505b5050505050815250509050919050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146106e6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103aa565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610767576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103aa565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b50805460008255600701600890049060005260206000209081019061052091905b8082111561081257600081556001016107fe565b5090565b60006020828403121561082857600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461085857600080fd5b9392505050565b600060208083528351808285015260005b8181101561088c57858101830151858201604001528201610870565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b63ffffffff8116811461052057600080fd5b600080604083850312156108f057600080fd5b82356108fb816108cb565b9150602083013567ffffffffffffffff81111561091757600080fd5b83016060818603121561092957600080fd5b809150509250929050565b60006020828403121561094657600080fd5b8135610858816108cb565b60006020828403121561096357600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461085857600080fd5b60006020808352608083017fffffffffffffffff0000000000000000000000000000000000000000000000008551168285015267ffffffffffffffff82860151166040850152604085015160608086015281815180845260a0870191508483019350600092505b80831015610a1457835163ffffffff1682529284019260019290920191908401906109ee565b509695505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610a5457600080fd5b83018035915067ffffffffffffffff821115610a6f57600080fd5b6020019150600581901b3603821315610a8757600080fd5b9250929050565b67ffffffffffffffff8116811461052057600080fd5b600060208284031215610ab657600080fd5b813561085881610a8e565b7fffffffffffffffff0000000000000000000000000000000000000000000000008116811461052057600080fd5b600060208284031215610b0157600080fd5b813561085881610ac1565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b68010000000000000000821115610b5457610b54610b0c565b805482825580831015610bd9576000828152602081206007850160031c81016007840160031c82019150601c8660021b168015610bc0577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8083018054828460200360031b1c16815550505b505b81811015610bd557828155600101610bc2565b5050505b505050565b60008135610beb816108cb565b92915050565b67ffffffffffffffff831115610c0957610c09610b0c565b610c138382610b3b565b60008181526020902082908460031c60005b81811015610c7e576000805b6008811015610c7157610c60610c4687610bde565b63ffffffff908116600584901b90811b91901b1984161790565b602096909601959150600101610c31565b5083820155600101610c25565b507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff88616808703818814610cdb576000805b82811015610cd557610cc4610c4688610bde565b602097909701969150600101610cb0565b50848401555b5050505050505050565b8135610cf081610ac1565b8060c01c90508154817fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000082161783556020840135610d2d81610a8e565b6fffffffffffffffff00000000000000008160401b16837fffffffffffffffffffffffffffffffff0000000000000000000000000000000084161717845550505060408201357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1833603018112610da357600080fd5b8201803567ffffffffffffffff811115610dbc57600080fd5b6020820191508060051b3603821315610dd457600080fd5b610de2818360018601610bf1565b50505050565b8183526000602080850194508260005b85811015610e23578135610e0b816108cb565b63ffffffff1687529582019590820190600101610df8565b509495945050505050565b63ffffffff831681526040602082015260008235610e4b81610ac1565b7fffffffffffffffff0000000000000000000000000000000000000000000000001660408301526020830135610e8081610a8e565b67ffffffffffffffff8082166060850152604085013591507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1853603018212610ec857600080fd5b6020918501918201913581811115610edf57600080fd5b8060051b3603831315610ef157600080fd5b60606080860152610f0660a086018285610de8565b97965050505050505056fea164736f6c6343000813000a",
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

func (_ChannelConfigStore *ChannelConfigStoreCaller) GetChannelDefinitions(opts *bind.CallOpts, channelId uint32) (IChannelConfigStoreChannelDefinition, error) {
	var out []interface{}
	err := _ChannelConfigStore.contract.Call(opts, &out, "getChannelDefinitions", channelId)

	if err != nil {
		return *new(IChannelConfigStoreChannelDefinition), err
	}

	out0 := *abi.ConvertType(out[0], new(IChannelConfigStoreChannelDefinition)).(*IChannelConfigStoreChannelDefinition)

	return out0, err

}

func (_ChannelConfigStore *ChannelConfigStoreSession) GetChannelDefinitions(channelId uint32) (IChannelConfigStoreChannelDefinition, error) {
	return _ChannelConfigStore.Contract.GetChannelDefinitions(&_ChannelConfigStore.CallOpts, channelId)
}

func (_ChannelConfigStore *ChannelConfigStoreCallerSession) GetChannelDefinitions(channelId uint32) (IChannelConfigStoreChannelDefinition, error) {
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

func (_ChannelConfigStore *ChannelConfigStoreTransactor) AddChannel(opts *bind.TransactOpts, channelId uint32, channelDefinition IChannelConfigStoreChannelDefinition) (*types.Transaction, error) {
	return _ChannelConfigStore.contract.Transact(opts, "addChannel", channelId, channelDefinition)
}

func (_ChannelConfigStore *ChannelConfigStoreSession) AddChannel(channelId uint32, channelDefinition IChannelConfigStoreChannelDefinition) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.AddChannel(&_ChannelConfigStore.TransactOpts, channelId, channelDefinition)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorSession) AddChannel(channelId uint32, channelDefinition IChannelConfigStoreChannelDefinition) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.AddChannel(&_ChannelConfigStore.TransactOpts, channelId, channelDefinition)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactor) RemoveChannel(opts *bind.TransactOpts, channelId uint32) (*types.Transaction, error) {
	return _ChannelConfigStore.contract.Transact(opts, "removeChannel", channelId)
}

func (_ChannelConfigStore *ChannelConfigStoreSession) RemoveChannel(channelId uint32) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.RemoveChannel(&_ChannelConfigStore.TransactOpts, channelId)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorSession) RemoveChannel(channelId uint32) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.RemoveChannel(&_ChannelConfigStore.TransactOpts, channelId)
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
	ChannelId uint32
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
	ChannelId         uint32
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
	ChannelId uint32
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
	return common.HexToHash("0x334e877e9691ecae0660510061973bebaa8b4fb37332ed6090052e630c979861")
}

func (ChannelConfigStoreNewChannelDefinition) Topic() common.Hash {
	return common.HexToHash("0xbf2cd44714205d633d3f888ac72ea66d53cd12d4c4e8723a80d9c0bc36484a54")
}

func (ChannelConfigStoreOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (ChannelConfigStoreOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (ChannelConfigStorePromoteStagingConfig) Topic() common.Hash {
	return common.HexToHash("0xbdd8ee023f9979bf23e8af6fd7241f484024e83fb0fabd11bb7fd5e9bed7308a")
}

func (_ChannelConfigStore *ChannelConfigStore) Address() common.Address {
	return _ChannelConfigStore.address
}

type ChannelConfigStoreInterface interface {
	GetChannelDefinitions(opts *bind.CallOpts, channelId uint32) (IChannelConfigStoreChannelDefinition, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddChannel(opts *bind.TransactOpts, channelId uint32, channelDefinition IChannelConfigStoreChannelDefinition) (*types.Transaction, error)

	RemoveChannel(opts *bind.TransactOpts, channelId uint32) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterChannelDefinitionRemoved(opts *bind.FilterOpts) (*ChannelConfigStoreChannelDefinitionRemovedIterator, error)

	WatchChannelDefinitionRemoved(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreChannelDefinitionRemoved) (event.Subscription, error)

	ParseChannelDefinitionRemoved(log types.Log) (*ChannelConfigStoreChannelDefinitionRemoved, error)

	FilterNewChannelDefinition(opts *bind.FilterOpts) (*ChannelConfigStoreNewChannelDefinitionIterator, error)

	WatchNewChannelDefinition(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreNewChannelDefinition) (event.Subscription, error)

	ParseNewChannelDefinition(log types.Log) (*ChannelConfigStoreNewChannelDefinition, error)

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
