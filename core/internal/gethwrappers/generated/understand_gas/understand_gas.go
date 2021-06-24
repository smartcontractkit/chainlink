// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package understand_gas

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

const UnderstandGasABI = "[{\"inputs\":[],\"name\":\"f\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

var UnderstandGasBin = "0x6080604052348015600f57600080fd5b50603c80601d6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806326121ff014602d575b600080fd5b00fea164736f6c6343000804000a"

func DeployUnderstandGas(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UnderstandGas, error) {
	parsed, err := abi.JSON(strings.NewReader(UnderstandGasABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(UnderstandGasBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UnderstandGas{UnderstandGasCaller: UnderstandGasCaller{contract: contract}, UnderstandGasTransactor: UnderstandGasTransactor{contract: contract}, UnderstandGasFilterer: UnderstandGasFilterer{contract: contract}}, nil
}

type UnderstandGas struct {
	address common.Address
	abi     abi.ABI
	UnderstandGasCaller
	UnderstandGasTransactor
	UnderstandGasFilterer
}

type UnderstandGasCaller struct {
	contract *bind.BoundContract
}

type UnderstandGasTransactor struct {
	contract *bind.BoundContract
}

type UnderstandGasFilterer struct {
	contract *bind.BoundContract
}

type UnderstandGasSession struct {
	Contract     *UnderstandGas
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type UnderstandGasCallerSession struct {
	Contract *UnderstandGasCaller
	CallOpts bind.CallOpts
}

type UnderstandGasTransactorSession struct {
	Contract     *UnderstandGasTransactor
	TransactOpts bind.TransactOpts
}

type UnderstandGasRaw struct {
	Contract *UnderstandGas
}

type UnderstandGasCallerRaw struct {
	Contract *UnderstandGasCaller
}

type UnderstandGasTransactorRaw struct {
	Contract *UnderstandGasTransactor
}

func NewUnderstandGas(address common.Address, backend bind.ContractBackend) (*UnderstandGas, error) {
	abi, err := abi.JSON(strings.NewReader(UnderstandGasABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUnderstandGas(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UnderstandGas{address: address, abi: abi, UnderstandGasCaller: UnderstandGasCaller{contract: contract}, UnderstandGasTransactor: UnderstandGasTransactor{contract: contract}, UnderstandGasFilterer: UnderstandGasFilterer{contract: contract}}, nil
}

func NewUnderstandGasCaller(address common.Address, caller bind.ContractCaller) (*UnderstandGasCaller, error) {
	contract, err := bindUnderstandGas(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UnderstandGasCaller{contract: contract}, nil
}

func NewUnderstandGasTransactor(address common.Address, transactor bind.ContractTransactor) (*UnderstandGasTransactor, error) {
	contract, err := bindUnderstandGas(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UnderstandGasTransactor{contract: contract}, nil
}

func NewUnderstandGasFilterer(address common.Address, filterer bind.ContractFilterer) (*UnderstandGasFilterer, error) {
	contract, err := bindUnderstandGas(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UnderstandGasFilterer{contract: contract}, nil
}

func bindUnderstandGas(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UnderstandGasABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_UnderstandGas *UnderstandGasRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UnderstandGas.Contract.UnderstandGasCaller.contract.Call(opts, result, method, params...)
}

func (_UnderstandGas *UnderstandGasRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UnderstandGas.Contract.UnderstandGasTransactor.contract.Transfer(opts)
}

func (_UnderstandGas *UnderstandGasRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UnderstandGas.Contract.UnderstandGasTransactor.contract.Transact(opts, method, params...)
}

func (_UnderstandGas *UnderstandGasCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UnderstandGas.Contract.contract.Call(opts, result, method, params...)
}

func (_UnderstandGas *UnderstandGasTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UnderstandGas.Contract.contract.Transfer(opts)
}

func (_UnderstandGas *UnderstandGasTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UnderstandGas.Contract.contract.Transact(opts, method, params...)
}

func (_UnderstandGas *UnderstandGasTransactor) F(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UnderstandGas.contract.Transact(opts, "f")
}

func (_UnderstandGas *UnderstandGasSession) F() (*types.Transaction, error) {
	return _UnderstandGas.Contract.F(&_UnderstandGas.TransactOpts)
}

func (_UnderstandGas *UnderstandGasTransactorSession) F() (*types.Transaction, error) {
	return _UnderstandGas.Contract.F(&_UnderstandGas.TransactOpts)
}

func (_UnderstandGas *UnderstandGas) Address() common.Address {
	return _UnderstandGas.address
}

type UnderstandGasInterface interface {
	F(opts *bind.TransactOpts) (*types.Transaction, error)

	Address() common.Address
}
