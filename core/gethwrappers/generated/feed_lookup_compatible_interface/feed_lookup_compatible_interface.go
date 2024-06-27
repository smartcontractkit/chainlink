// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package feed_lookup_compatible_interface

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
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

var FeedLookupCompatibleInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var FeedLookupCompatibleInterfaceABI = FeedLookupCompatibleInterfaceMetaData.ABI

type FeedLookupCompatibleInterface struct {
	address common.Address
	abi     abi.ABI
	FeedLookupCompatibleInterfaceCaller
	FeedLookupCompatibleInterfaceTransactor
	FeedLookupCompatibleInterfaceFilterer
}

type FeedLookupCompatibleInterfaceCaller struct {
	contract *bind.BoundContract
}

type FeedLookupCompatibleInterfaceTransactor struct {
	contract *bind.BoundContract
}

type FeedLookupCompatibleInterfaceFilterer struct {
	contract *bind.BoundContract
}

type FeedLookupCompatibleInterfaceSession struct {
	Contract     *FeedLookupCompatibleInterface
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FeedLookupCompatibleInterfaceCallerSession struct {
	Contract *FeedLookupCompatibleInterfaceCaller
	CallOpts bind.CallOpts
}

type FeedLookupCompatibleInterfaceTransactorSession struct {
	Contract     *FeedLookupCompatibleInterfaceTransactor
	TransactOpts bind.TransactOpts
}

type FeedLookupCompatibleInterfaceRaw struct {
	Contract *FeedLookupCompatibleInterface
}

type FeedLookupCompatibleInterfaceCallerRaw struct {
	Contract *FeedLookupCompatibleInterfaceCaller
}

type FeedLookupCompatibleInterfaceTransactorRaw struct {
	Contract *FeedLookupCompatibleInterfaceTransactor
}

func NewFeedLookupCompatibleInterface(address common.Address, backend bind.ContractBackend) (*FeedLookupCompatibleInterface, error) {
	abi, err := abi.JSON(strings.NewReader(FeedLookupCompatibleInterfaceABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFeedLookupCompatibleInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeedLookupCompatibleInterface{address: address, abi: abi, FeedLookupCompatibleInterfaceCaller: FeedLookupCompatibleInterfaceCaller{contract: contract}, FeedLookupCompatibleInterfaceTransactor: FeedLookupCompatibleInterfaceTransactor{contract: contract}, FeedLookupCompatibleInterfaceFilterer: FeedLookupCompatibleInterfaceFilterer{contract: contract}}, nil
}

func NewFeedLookupCompatibleInterfaceCaller(address common.Address, caller bind.ContractCaller) (*FeedLookupCompatibleInterfaceCaller, error) {
	contract, err := bindFeedLookupCompatibleInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeedLookupCompatibleInterfaceCaller{contract: contract}, nil
}

func NewFeedLookupCompatibleInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*FeedLookupCompatibleInterfaceTransactor, error) {
	contract, err := bindFeedLookupCompatibleInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeedLookupCompatibleInterfaceTransactor{contract: contract}, nil
}

func NewFeedLookupCompatibleInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*FeedLookupCompatibleInterfaceFilterer, error) {
	contract, err := bindFeedLookupCompatibleInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeedLookupCompatibleInterfaceFilterer{contract: contract}, nil
}

func bindFeedLookupCompatibleInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeedLookupCompatibleInterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FeedLookupCompatibleInterface *FeedLookupCompatibleInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeedLookupCompatibleInterface.Contract.FeedLookupCompatibleInterfaceCaller.contract.Call(opts, result, method, params...)
}

func (_FeedLookupCompatibleInterface *FeedLookupCompatibleInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeedLookupCompatibleInterface.Contract.FeedLookupCompatibleInterfaceTransactor.contract.Transfer(opts)
}

func (_FeedLookupCompatibleInterface *FeedLookupCompatibleInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeedLookupCompatibleInterface.Contract.FeedLookupCompatibleInterfaceTransactor.contract.Transact(opts, method, params...)
}

func (_FeedLookupCompatibleInterface *FeedLookupCompatibleInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeedLookupCompatibleInterface.Contract.contract.Call(opts, result, method, params...)
}

func (_FeedLookupCompatibleInterface *FeedLookupCompatibleInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeedLookupCompatibleInterface.Contract.contract.Transfer(opts)
}

func (_FeedLookupCompatibleInterface *FeedLookupCompatibleInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeedLookupCompatibleInterface.Contract.contract.Transact(opts, method, params...)
}

func (_FeedLookupCompatibleInterface *FeedLookupCompatibleInterfaceCaller) CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	var out []interface{}
	err := _FeedLookupCompatibleInterface.contract.Call(opts, &out, "checkCallback", values, extraData)

	outstruct := new(CheckCallback)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_FeedLookupCompatibleInterface *FeedLookupCompatibleInterfaceSession) CheckCallback(values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _FeedLookupCompatibleInterface.Contract.CheckCallback(&_FeedLookupCompatibleInterface.CallOpts, values, extraData)
}

func (_FeedLookupCompatibleInterface *FeedLookupCompatibleInterfaceCallerSession) CheckCallback(values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _FeedLookupCompatibleInterface.Contract.CheckCallback(&_FeedLookupCompatibleInterface.CallOpts, values, extraData)
}

type CheckCallback struct {
	UpkeepNeeded bool
	PerformData  []byte
}

func (_FeedLookupCompatibleInterface *FeedLookupCompatibleInterface) Address() common.Address {
	return _FeedLookupCompatibleInterface.address
}

type FeedLookupCompatibleInterfaceInterface interface {
	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (CheckCallback,

		error)

	Address() common.Address
}
