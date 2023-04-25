// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mercury_lookup_compatible_interface

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

var MercuryLookupCompatibleInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"mercuryCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var MercuryLookupCompatibleInterfaceABI = MercuryLookupCompatibleInterfaceMetaData.ABI

type MercuryLookupCompatibleInterface struct {
	address common.Address
	abi     abi.ABI
	MercuryLookupCompatibleInterfaceCaller
	MercuryLookupCompatibleInterfaceTransactor
	MercuryLookupCompatibleInterfaceFilterer
}

type MercuryLookupCompatibleInterfaceCaller struct {
	contract *bind.BoundContract
}

type MercuryLookupCompatibleInterfaceTransactor struct {
	contract *bind.BoundContract
}

type MercuryLookupCompatibleInterfaceFilterer struct {
	contract *bind.BoundContract
}

type MercuryLookupCompatibleInterfaceSession struct {
	Contract     *MercuryLookupCompatibleInterface
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MercuryLookupCompatibleInterfaceCallerSession struct {
	Contract *MercuryLookupCompatibleInterfaceCaller
	CallOpts bind.CallOpts
}

type MercuryLookupCompatibleInterfaceTransactorSession struct {
	Contract     *MercuryLookupCompatibleInterfaceTransactor
	TransactOpts bind.TransactOpts
}

type MercuryLookupCompatibleInterfaceRaw struct {
	Contract *MercuryLookupCompatibleInterface
}

type MercuryLookupCompatibleInterfaceCallerRaw struct {
	Contract *MercuryLookupCompatibleInterfaceCaller
}

type MercuryLookupCompatibleInterfaceTransactorRaw struct {
	Contract *MercuryLookupCompatibleInterfaceTransactor
}

func NewMercuryLookupCompatibleInterface(address common.Address, backend bind.ContractBackend) (*MercuryLookupCompatibleInterface, error) {
	abi, err := abi.JSON(strings.NewReader(MercuryLookupCompatibleInterfaceABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMercuryLookupCompatibleInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MercuryLookupCompatibleInterface{address: address, abi: abi, MercuryLookupCompatibleInterfaceCaller: MercuryLookupCompatibleInterfaceCaller{contract: contract}, MercuryLookupCompatibleInterfaceTransactor: MercuryLookupCompatibleInterfaceTransactor{contract: contract}, MercuryLookupCompatibleInterfaceFilterer: MercuryLookupCompatibleInterfaceFilterer{contract: contract}}, nil
}

func NewMercuryLookupCompatibleInterfaceCaller(address common.Address, caller bind.ContractCaller) (*MercuryLookupCompatibleInterfaceCaller, error) {
	contract, err := bindMercuryLookupCompatibleInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryLookupCompatibleInterfaceCaller{contract: contract}, nil
}

func NewMercuryLookupCompatibleInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*MercuryLookupCompatibleInterfaceTransactor, error) {
	contract, err := bindMercuryLookupCompatibleInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryLookupCompatibleInterfaceTransactor{contract: contract}, nil
}

func NewMercuryLookupCompatibleInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*MercuryLookupCompatibleInterfaceFilterer, error) {
	contract, err := bindMercuryLookupCompatibleInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MercuryLookupCompatibleInterfaceFilterer{contract: contract}, nil
}

func bindMercuryLookupCompatibleInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MercuryLookupCompatibleInterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MercuryLookupCompatibleInterface *MercuryLookupCompatibleInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryLookupCompatibleInterface.Contract.MercuryLookupCompatibleInterfaceCaller.contract.Call(opts, result, method, params...)
}

func (_MercuryLookupCompatibleInterface *MercuryLookupCompatibleInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryLookupCompatibleInterface.Contract.MercuryLookupCompatibleInterfaceTransactor.contract.Transfer(opts)
}

func (_MercuryLookupCompatibleInterface *MercuryLookupCompatibleInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryLookupCompatibleInterface.Contract.MercuryLookupCompatibleInterfaceTransactor.contract.Transact(opts, method, params...)
}

func (_MercuryLookupCompatibleInterface *MercuryLookupCompatibleInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryLookupCompatibleInterface.Contract.contract.Call(opts, result, method, params...)
}

func (_MercuryLookupCompatibleInterface *MercuryLookupCompatibleInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryLookupCompatibleInterface.Contract.contract.Transfer(opts)
}

func (_MercuryLookupCompatibleInterface *MercuryLookupCompatibleInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryLookupCompatibleInterface.Contract.contract.Transact(opts, method, params...)
}

func (_MercuryLookupCompatibleInterface *MercuryLookupCompatibleInterfaceCaller) MercuryCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (MercuryCallback,

	error) {
	var out []interface{}
	err := _MercuryLookupCompatibleInterface.contract.Call(opts, &out, "mercuryCallback", values, extraData)

	outstruct := new(MercuryCallback)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_MercuryLookupCompatibleInterface *MercuryLookupCompatibleInterfaceSession) MercuryCallback(values [][]byte, extraData []byte) (MercuryCallback,

	error) {
	return _MercuryLookupCompatibleInterface.Contract.MercuryCallback(&_MercuryLookupCompatibleInterface.CallOpts, values, extraData)
}

func (_MercuryLookupCompatibleInterface *MercuryLookupCompatibleInterfaceCallerSession) MercuryCallback(values [][]byte, extraData []byte) (MercuryCallback,

	error) {
	return _MercuryLookupCompatibleInterface.Contract.MercuryCallback(&_MercuryLookupCompatibleInterface.CallOpts, values, extraData)
}

type MercuryCallback struct {
	UpkeepNeeded bool
	PerformData  []byte
}

func (_MercuryLookupCompatibleInterface *MercuryLookupCompatibleInterface) Address() common.Address {
	return _MercuryLookupCompatibleInterface.address
}

type MercuryLookupCompatibleInterfaceInterface interface {
	MercuryCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (MercuryCallback,

		error)

	Address() common.Address
}
