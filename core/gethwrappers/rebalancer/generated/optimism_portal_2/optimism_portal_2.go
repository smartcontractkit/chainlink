// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_portal_2

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

var OptimismPortal2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"disputeGameFactory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"respectedGameType\",\"outputs\":[{\"internalType\":\"GameType\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var OptimismPortal2ABI = OptimismPortal2MetaData.ABI

type OptimismPortal2 struct {
	address common.Address
	abi     abi.ABI
	OptimismPortal2Caller
	OptimismPortal2Transactor
	OptimismPortal2Filterer
}

type OptimismPortal2Caller struct {
	contract *bind.BoundContract
}

type OptimismPortal2Transactor struct {
	contract *bind.BoundContract
}

type OptimismPortal2Filterer struct {
	contract *bind.BoundContract
}

type OptimismPortal2Session struct {
	Contract     *OptimismPortal2
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismPortal2CallerSession struct {
	Contract *OptimismPortal2Caller
	CallOpts bind.CallOpts
}

type OptimismPortal2TransactorSession struct {
	Contract     *OptimismPortal2Transactor
	TransactOpts bind.TransactOpts
}

type OptimismPortal2Raw struct {
	Contract *OptimismPortal2
}

type OptimismPortal2CallerRaw struct {
	Contract *OptimismPortal2Caller
}

type OptimismPortal2TransactorRaw struct {
	Contract *OptimismPortal2Transactor
}

func NewOptimismPortal2(address common.Address, backend bind.ContractBackend) (*OptimismPortal2, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismPortal2ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismPortal2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismPortal2{address: address, abi: abi, OptimismPortal2Caller: OptimismPortal2Caller{contract: contract}, OptimismPortal2Transactor: OptimismPortal2Transactor{contract: contract}, OptimismPortal2Filterer: OptimismPortal2Filterer{contract: contract}}, nil
}

func NewOptimismPortal2Caller(address common.Address, caller bind.ContractCaller) (*OptimismPortal2Caller, error) {
	contract, err := bindOptimismPortal2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismPortal2Caller{contract: contract}, nil
}

func NewOptimismPortal2Transactor(address common.Address, transactor bind.ContractTransactor) (*OptimismPortal2Transactor, error) {
	contract, err := bindOptimismPortal2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismPortal2Transactor{contract: contract}, nil
}

func NewOptimismPortal2Filterer(address common.Address, filterer bind.ContractFilterer) (*OptimismPortal2Filterer, error) {
	contract, err := bindOptimismPortal2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismPortal2Filterer{contract: contract}, nil
}

func bindOptimismPortal2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismPortal2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismPortal2 *OptimismPortal2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismPortal2.Contract.OptimismPortal2Caller.contract.Call(opts, result, method, params...)
}

func (_OptimismPortal2 *OptimismPortal2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismPortal2.Contract.OptimismPortal2Transactor.contract.Transfer(opts)
}

func (_OptimismPortal2 *OptimismPortal2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismPortal2.Contract.OptimismPortal2Transactor.contract.Transact(opts, method, params...)
}

func (_OptimismPortal2 *OptimismPortal2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismPortal2.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismPortal2 *OptimismPortal2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismPortal2.Contract.contract.Transfer(opts)
}

func (_OptimismPortal2 *OptimismPortal2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismPortal2.Contract.contract.Transact(opts, method, params...)
}

func (_OptimismPortal2 *OptimismPortal2Caller) DisputeGameFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OptimismPortal2.contract.Call(opts, &out, "disputeGameFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OptimismPortal2 *OptimismPortal2Session) DisputeGameFactory() (common.Address, error) {
	return _OptimismPortal2.Contract.DisputeGameFactory(&_OptimismPortal2.CallOpts)
}

func (_OptimismPortal2 *OptimismPortal2CallerSession) DisputeGameFactory() (common.Address, error) {
	return _OptimismPortal2.Contract.DisputeGameFactory(&_OptimismPortal2.CallOpts)
}

func (_OptimismPortal2 *OptimismPortal2Caller) RespectedGameType(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _OptimismPortal2.contract.Call(opts, &out, "respectedGameType")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_OptimismPortal2 *OptimismPortal2Session) RespectedGameType() (uint32, error) {
	return _OptimismPortal2.Contract.RespectedGameType(&_OptimismPortal2.CallOpts)
}

func (_OptimismPortal2 *OptimismPortal2CallerSession) RespectedGameType() (uint32, error) {
	return _OptimismPortal2.Contract.RespectedGameType(&_OptimismPortal2.CallOpts)
}

func (_OptimismPortal2 *OptimismPortal2) Address() common.Address {
	return _OptimismPortal2.address
}

type OptimismPortal2Interface interface {
	DisputeGameFactory(opts *bind.CallOpts) (common.Address, error)

	RespectedGameType(opts *bind.CallOpts) (uint32, error)

	Address() common.Address
}
