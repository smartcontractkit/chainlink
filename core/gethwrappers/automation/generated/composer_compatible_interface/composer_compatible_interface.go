// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package composer_compatible_interface

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

var ComposerCompatibleInterfaceV1MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"scriptHash\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"functionsArguments\",\"type\":\"string[]\"},{\"internalType\":\"bool\",\"name\":\"useMercury\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"ComposerRequestV1\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var ComposerCompatibleInterfaceV1ABI = ComposerCompatibleInterfaceV1MetaData.ABI

type ComposerCompatibleInterfaceV1 struct {
	address common.Address
	abi     abi.ABI
	ComposerCompatibleInterfaceV1Caller
	ComposerCompatibleInterfaceV1Transactor
	ComposerCompatibleInterfaceV1Filterer
}

type ComposerCompatibleInterfaceV1Caller struct {
	contract *bind.BoundContract
}

type ComposerCompatibleInterfaceV1Transactor struct {
	contract *bind.BoundContract
}

type ComposerCompatibleInterfaceV1Filterer struct {
	contract *bind.BoundContract
}

type ComposerCompatibleInterfaceV1Session struct {
	Contract     *ComposerCompatibleInterfaceV1
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ComposerCompatibleInterfaceV1CallerSession struct {
	Contract *ComposerCompatibleInterfaceV1Caller
	CallOpts bind.CallOpts
}

type ComposerCompatibleInterfaceV1TransactorSession struct {
	Contract     *ComposerCompatibleInterfaceV1Transactor
	TransactOpts bind.TransactOpts
}

type ComposerCompatibleInterfaceV1Raw struct {
	Contract *ComposerCompatibleInterfaceV1
}

type ComposerCompatibleInterfaceV1CallerRaw struct {
	Contract *ComposerCompatibleInterfaceV1Caller
}

type ComposerCompatibleInterfaceV1TransactorRaw struct {
	Contract *ComposerCompatibleInterfaceV1Transactor
}

func NewComposerCompatibleInterfaceV1(address common.Address, backend bind.ContractBackend) (*ComposerCompatibleInterfaceV1, error) {
	abi, err := abi.JSON(strings.NewReader(ComposerCompatibleInterfaceV1ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindComposerCompatibleInterfaceV1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ComposerCompatibleInterfaceV1{address: address, abi: abi, ComposerCompatibleInterfaceV1Caller: ComposerCompatibleInterfaceV1Caller{contract: contract}, ComposerCompatibleInterfaceV1Transactor: ComposerCompatibleInterfaceV1Transactor{contract: contract}, ComposerCompatibleInterfaceV1Filterer: ComposerCompatibleInterfaceV1Filterer{contract: contract}}, nil
}

func NewComposerCompatibleInterfaceV1Caller(address common.Address, caller bind.ContractCaller) (*ComposerCompatibleInterfaceV1Caller, error) {
	contract, err := bindComposerCompatibleInterfaceV1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ComposerCompatibleInterfaceV1Caller{contract: contract}, nil
}

func NewComposerCompatibleInterfaceV1Transactor(address common.Address, transactor bind.ContractTransactor) (*ComposerCompatibleInterfaceV1Transactor, error) {
	contract, err := bindComposerCompatibleInterfaceV1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ComposerCompatibleInterfaceV1Transactor{contract: contract}, nil
}

func NewComposerCompatibleInterfaceV1Filterer(address common.Address, filterer bind.ContractFilterer) (*ComposerCompatibleInterfaceV1Filterer, error) {
	contract, err := bindComposerCompatibleInterfaceV1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ComposerCompatibleInterfaceV1Filterer{contract: contract}, nil
}

func bindComposerCompatibleInterfaceV1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ComposerCompatibleInterfaceV1MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ComposerCompatibleInterfaceV1 *ComposerCompatibleInterfaceV1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ComposerCompatibleInterfaceV1.Contract.ComposerCompatibleInterfaceV1Caller.contract.Call(opts, result, method, params...)
}

func (_ComposerCompatibleInterfaceV1 *ComposerCompatibleInterfaceV1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ComposerCompatibleInterfaceV1.Contract.ComposerCompatibleInterfaceV1Transactor.contract.Transfer(opts)
}

func (_ComposerCompatibleInterfaceV1 *ComposerCompatibleInterfaceV1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ComposerCompatibleInterfaceV1.Contract.ComposerCompatibleInterfaceV1Transactor.contract.Transact(opts, method, params...)
}

func (_ComposerCompatibleInterfaceV1 *ComposerCompatibleInterfaceV1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ComposerCompatibleInterfaceV1.Contract.contract.Call(opts, result, method, params...)
}

func (_ComposerCompatibleInterfaceV1 *ComposerCompatibleInterfaceV1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ComposerCompatibleInterfaceV1.Contract.contract.Transfer(opts)
}

func (_ComposerCompatibleInterfaceV1 *ComposerCompatibleInterfaceV1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ComposerCompatibleInterfaceV1.Contract.contract.Transact(opts, method, params...)
}

func (_ComposerCompatibleInterfaceV1 *ComposerCompatibleInterfaceV1Caller) CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	var out []interface{}
	err := _ComposerCompatibleInterfaceV1.contract.Call(opts, &out, "checkCallback", values, extraData)

	outstruct := new(CheckCallback)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_ComposerCompatibleInterfaceV1 *ComposerCompatibleInterfaceV1Session) CheckCallback(values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _ComposerCompatibleInterfaceV1.Contract.CheckCallback(&_ComposerCompatibleInterfaceV1.CallOpts, values, extraData)
}

func (_ComposerCompatibleInterfaceV1 *ComposerCompatibleInterfaceV1CallerSession) CheckCallback(values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _ComposerCompatibleInterfaceV1.Contract.CheckCallback(&_ComposerCompatibleInterfaceV1.CallOpts, values, extraData)
}

type CheckCallback struct {
	UpkeepNeeded bool
	PerformData  []byte
}

func (_ComposerCompatibleInterfaceV1 *ComposerCompatibleInterfaceV1) Address() common.Address {
	return _ComposerCompatibleInterfaceV1.address
}

type ComposerCompatibleInterfaceV1Interface interface {
	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (CheckCallback,

		error)

	Address() common.Address
}
