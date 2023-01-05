// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions

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
)

var FunctionsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyUrl\",\"type\":\"error\"}]",
	Bin: "0x602d6037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea164736f6c6343000806000a",
}

var FunctionsABI = FunctionsMetaData.ABI

var FunctionsBin = FunctionsMetaData.Bin

func DeployFunctions(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Functions, error) {
	parsed, err := FunctionsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Functions{FunctionsCaller: FunctionsCaller{contract: contract}, FunctionsTransactor: FunctionsTransactor{contract: contract}, FunctionsFilterer: FunctionsFilterer{contract: contract}}, nil
}

type Functions struct {
	address common.Address
	abi     abi.ABI
	FunctionsCaller
	FunctionsTransactor
	FunctionsFilterer
}

type FunctionsCaller struct {
	contract *bind.BoundContract
}

type FunctionsTransactor struct {
	contract *bind.BoundContract
}

type FunctionsFilterer struct {
	contract *bind.BoundContract
}

type FunctionsSession struct {
	Contract     *Functions
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsCallerSession struct {
	Contract *FunctionsCaller
	CallOpts bind.CallOpts
}

type FunctionsTransactorSession struct {
	Contract     *FunctionsTransactor
	TransactOpts bind.TransactOpts
}

type FunctionsRaw struct {
	Contract *Functions
}

type FunctionsCallerRaw struct {
	Contract *FunctionsCaller
}

type FunctionsTransactorRaw struct {
	Contract *FunctionsTransactor
}

func NewFunctions(address common.Address, backend bind.ContractBackend) (*Functions, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctions(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Functions{address: address, abi: abi, FunctionsCaller: FunctionsCaller{contract: contract}, FunctionsTransactor: FunctionsTransactor{contract: contract}, FunctionsFilterer: FunctionsFilterer{contract: contract}}, nil
}

func NewFunctionsCaller(address common.Address, caller bind.ContractCaller) (*FunctionsCaller, error) {
	contract, err := bindFunctions(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsCaller{contract: contract}, nil
}

func NewFunctionsTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsTransactor, error) {
	contract, err := bindFunctions(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsTransactor{contract: contract}, nil
}

func NewFunctionsFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsFilterer, error) {
	contract, err := bindFunctions(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsFilterer{contract: contract}, nil
}

func bindFunctions(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FunctionsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_Functions *FunctionsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Functions.Contract.FunctionsCaller.contract.Call(opts, result, method, params...)
}

func (_Functions *FunctionsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Functions.Contract.FunctionsTransactor.contract.Transfer(opts)
}

func (_Functions *FunctionsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Functions.Contract.FunctionsTransactor.contract.Transact(opts, method, params...)
}

func (_Functions *FunctionsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Functions.Contract.contract.Call(opts, result, method, params...)
}

func (_Functions *FunctionsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Functions.Contract.contract.Transfer(opts)
}

func (_Functions *FunctionsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Functions.Contract.contract.Transact(opts, method, params...)
}

func (_Functions *Functions) Address() common.Address {
	return _Functions.address
}

type FunctionsInterface interface {
	Address() common.Address
}
