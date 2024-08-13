// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_l1_standard_bridge

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

var OptimismL1StandardBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_minGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"depositETHTo\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

var OptimismL1StandardBridgeABI = OptimismL1StandardBridgeMetaData.ABI

type OptimismL1StandardBridge struct {
	address common.Address
	abi     abi.ABI
	OptimismL1StandardBridgeCaller
	OptimismL1StandardBridgeTransactor
	OptimismL1StandardBridgeFilterer
}

type OptimismL1StandardBridgeCaller struct {
	contract *bind.BoundContract
}

type OptimismL1StandardBridgeTransactor struct {
	contract *bind.BoundContract
}

type OptimismL1StandardBridgeFilterer struct {
	contract *bind.BoundContract
}

type OptimismL1StandardBridgeSession struct {
	Contract     *OptimismL1StandardBridge
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismL1StandardBridgeCallerSession struct {
	Contract *OptimismL1StandardBridgeCaller
	CallOpts bind.CallOpts
}

type OptimismL1StandardBridgeTransactorSession struct {
	Contract     *OptimismL1StandardBridgeTransactor
	TransactOpts bind.TransactOpts
}

type OptimismL1StandardBridgeRaw struct {
	Contract *OptimismL1StandardBridge
}

type OptimismL1StandardBridgeCallerRaw struct {
	Contract *OptimismL1StandardBridgeCaller
}

type OptimismL1StandardBridgeTransactorRaw struct {
	Contract *OptimismL1StandardBridgeTransactor
}

func NewOptimismL1StandardBridge(address common.Address, backend bind.ContractBackend) (*OptimismL1StandardBridge, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismL1StandardBridgeABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismL1StandardBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismL1StandardBridge{address: address, abi: abi, OptimismL1StandardBridgeCaller: OptimismL1StandardBridgeCaller{contract: contract}, OptimismL1StandardBridgeTransactor: OptimismL1StandardBridgeTransactor{contract: contract}, OptimismL1StandardBridgeFilterer: OptimismL1StandardBridgeFilterer{contract: contract}}, nil
}

func NewOptimismL1StandardBridgeCaller(address common.Address, caller bind.ContractCaller) (*OptimismL1StandardBridgeCaller, error) {
	contract, err := bindOptimismL1StandardBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismL1StandardBridgeCaller{contract: contract}, nil
}

func NewOptimismL1StandardBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*OptimismL1StandardBridgeTransactor, error) {
	contract, err := bindOptimismL1StandardBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismL1StandardBridgeTransactor{contract: contract}, nil
}

func NewOptimismL1StandardBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*OptimismL1StandardBridgeFilterer, error) {
	contract, err := bindOptimismL1StandardBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismL1StandardBridgeFilterer{contract: contract}, nil
}

func bindOptimismL1StandardBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismL1StandardBridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismL1StandardBridge *OptimismL1StandardBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismL1StandardBridge.Contract.OptimismL1StandardBridgeCaller.contract.Call(opts, result, method, params...)
}

func (_OptimismL1StandardBridge *OptimismL1StandardBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL1StandardBridge.Contract.OptimismL1StandardBridgeTransactor.contract.Transfer(opts)
}

func (_OptimismL1StandardBridge *OptimismL1StandardBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismL1StandardBridge.Contract.OptimismL1StandardBridgeTransactor.contract.Transact(opts, method, params...)
}

func (_OptimismL1StandardBridge *OptimismL1StandardBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismL1StandardBridge.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismL1StandardBridge *OptimismL1StandardBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL1StandardBridge.Contract.contract.Transfer(opts)
}

func (_OptimismL1StandardBridge *OptimismL1StandardBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismL1StandardBridge.Contract.contract.Transact(opts, method, params...)
}

func (_OptimismL1StandardBridge *OptimismL1StandardBridgeTransactor) DepositETHTo(opts *bind.TransactOpts, _to common.Address, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _OptimismL1StandardBridge.contract.Transact(opts, "depositETHTo", _to, _minGasLimit, _extraData)
}

func (_OptimismL1StandardBridge *OptimismL1StandardBridgeSession) DepositETHTo(_to common.Address, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _OptimismL1StandardBridge.Contract.DepositETHTo(&_OptimismL1StandardBridge.TransactOpts, _to, _minGasLimit, _extraData)
}

func (_OptimismL1StandardBridge *OptimismL1StandardBridgeTransactorSession) DepositETHTo(_to common.Address, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _OptimismL1StandardBridge.Contract.DepositETHTo(&_OptimismL1StandardBridge.TransactOpts, _to, _minGasLimit, _extraData)
}

func (_OptimismL1StandardBridge *OptimismL1StandardBridge) Address() common.Address {
	return _OptimismL1StandardBridge.address
}

type OptimismL1StandardBridgeInterface interface {
	DepositETHTo(opts *bind.TransactOpts, _to common.Address, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error)

	Address() common.Address
}
