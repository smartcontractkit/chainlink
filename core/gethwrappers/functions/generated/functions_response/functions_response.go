// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_response

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

var FunctionsResponseMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x602d6037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea164736f6c6343000813000a",
}

var FunctionsResponseABI = FunctionsResponseMetaData.ABI

var FunctionsResponseBin = FunctionsResponseMetaData.Bin

func DeployFunctionsResponse(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FunctionsResponse, error) {
	parsed, err := FunctionsResponseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsResponseBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FunctionsResponse{FunctionsResponseCaller: FunctionsResponseCaller{contract: contract}, FunctionsResponseTransactor: FunctionsResponseTransactor{contract: contract}, FunctionsResponseFilterer: FunctionsResponseFilterer{contract: contract}}, nil
}

type FunctionsResponse struct {
	address common.Address
	abi     abi.ABI
	FunctionsResponseCaller
	FunctionsResponseTransactor
	FunctionsResponseFilterer
}

type FunctionsResponseCaller struct {
	contract *bind.BoundContract
}

type FunctionsResponseTransactor struct {
	contract *bind.BoundContract
}

type FunctionsResponseFilterer struct {
	contract *bind.BoundContract
}

type FunctionsResponseSession struct {
	Contract     *FunctionsResponse
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsResponseCallerSession struct {
	Contract *FunctionsResponseCaller
	CallOpts bind.CallOpts
}

type FunctionsResponseTransactorSession struct {
	Contract     *FunctionsResponseTransactor
	TransactOpts bind.TransactOpts
}

type FunctionsResponseRaw struct {
	Contract *FunctionsResponse
}

type FunctionsResponseCallerRaw struct {
	Contract *FunctionsResponseCaller
}

type FunctionsResponseTransactorRaw struct {
	Contract *FunctionsResponseTransactor
}

func NewFunctionsResponse(address common.Address, backend bind.ContractBackend) (*FunctionsResponse, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsResponseABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctionsResponse(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionsResponse{address: address, abi: abi, FunctionsResponseCaller: FunctionsResponseCaller{contract: contract}, FunctionsResponseTransactor: FunctionsResponseTransactor{contract: contract}, FunctionsResponseFilterer: FunctionsResponseFilterer{contract: contract}}, nil
}

func NewFunctionsResponseCaller(address common.Address, caller bind.ContractCaller) (*FunctionsResponseCaller, error) {
	contract, err := bindFunctionsResponse(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsResponseCaller{contract: contract}, nil
}

func NewFunctionsResponseTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsResponseTransactor, error) {
	contract, err := bindFunctionsResponse(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsResponseTransactor{contract: contract}, nil
}

func NewFunctionsResponseFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsResponseFilterer, error) {
	contract, err := bindFunctionsResponse(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsResponseFilterer{contract: contract}, nil
}

func bindFunctionsResponse(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionsResponseMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FunctionsResponse *FunctionsResponseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsResponse.Contract.FunctionsResponseCaller.contract.Call(opts, result, method, params...)
}

func (_FunctionsResponse *FunctionsResponseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsResponse.Contract.FunctionsResponseTransactor.contract.Transfer(opts)
}

func (_FunctionsResponse *FunctionsResponseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsResponse.Contract.FunctionsResponseTransactor.contract.Transact(opts, method, params...)
}

func (_FunctionsResponse *FunctionsResponseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsResponse.Contract.contract.Call(opts, result, method, params...)
}

func (_FunctionsResponse *FunctionsResponseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsResponse.Contract.contract.Transfer(opts)
}

func (_FunctionsResponse *FunctionsResponseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsResponse.Contract.contract.Transact(opts, method, params...)
}

func (_FunctionsResponse *FunctionsResponse) Address() common.Address {
	return _FunctionsResponse.address
}

type FunctionsResponseInterface interface {
	Address() common.Address
}
