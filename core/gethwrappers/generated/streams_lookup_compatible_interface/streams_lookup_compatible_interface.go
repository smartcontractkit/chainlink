// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package streams_lookup_compatible_interface

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

var StreamsLookupCompatibleInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"checkCallback\",\"inputs\":[{\"name\":\"values\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkErrorHandler\",\"inputs\":[{\"name\":\"errCode\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"StreamsLookup\",\"inputs\":[{\"name\":\"feedParamKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"feeds\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"timeParamKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"time\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]",
}

var StreamsLookupCompatibleInterfaceABI = StreamsLookupCompatibleInterfaceMetaData.ABI

type StreamsLookupCompatibleInterface struct {
	address common.Address
	abi     abi.ABI
	StreamsLookupCompatibleInterfaceCaller
	StreamsLookupCompatibleInterfaceTransactor
	StreamsLookupCompatibleInterfaceFilterer
}

type StreamsLookupCompatibleInterfaceCaller struct {
	contract *bind.BoundContract
}

type StreamsLookupCompatibleInterfaceTransactor struct {
	contract *bind.BoundContract
}

type StreamsLookupCompatibleInterfaceFilterer struct {
	contract *bind.BoundContract
}

type StreamsLookupCompatibleInterfaceSession struct {
	Contract     *StreamsLookupCompatibleInterface
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type StreamsLookupCompatibleInterfaceCallerSession struct {
	Contract *StreamsLookupCompatibleInterfaceCaller
	CallOpts bind.CallOpts
}

type StreamsLookupCompatibleInterfaceTransactorSession struct {
	Contract     *StreamsLookupCompatibleInterfaceTransactor
	TransactOpts bind.TransactOpts
}

type StreamsLookupCompatibleInterfaceRaw struct {
	Contract *StreamsLookupCompatibleInterface
}

type StreamsLookupCompatibleInterfaceCallerRaw struct {
	Contract *StreamsLookupCompatibleInterfaceCaller
}

type StreamsLookupCompatibleInterfaceTransactorRaw struct {
	Contract *StreamsLookupCompatibleInterfaceTransactor
}

func NewStreamsLookupCompatibleInterface(address common.Address, backend bind.ContractBackend) (*StreamsLookupCompatibleInterface, error) {
	abi, err := abi.JSON(strings.NewReader(StreamsLookupCompatibleInterfaceABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindStreamsLookupCompatibleInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StreamsLookupCompatibleInterface{address: address, abi: abi, StreamsLookupCompatibleInterfaceCaller: StreamsLookupCompatibleInterfaceCaller{contract: contract}, StreamsLookupCompatibleInterfaceTransactor: StreamsLookupCompatibleInterfaceTransactor{contract: contract}, StreamsLookupCompatibleInterfaceFilterer: StreamsLookupCompatibleInterfaceFilterer{contract: contract}}, nil
}

func NewStreamsLookupCompatibleInterfaceCaller(address common.Address, caller bind.ContractCaller) (*StreamsLookupCompatibleInterfaceCaller, error) {
	contract, err := bindStreamsLookupCompatibleInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StreamsLookupCompatibleInterfaceCaller{contract: contract}, nil
}

func NewStreamsLookupCompatibleInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*StreamsLookupCompatibleInterfaceTransactor, error) {
	contract, err := bindStreamsLookupCompatibleInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StreamsLookupCompatibleInterfaceTransactor{contract: contract}, nil
}

func NewStreamsLookupCompatibleInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*StreamsLookupCompatibleInterfaceFilterer, error) {
	contract, err := bindStreamsLookupCompatibleInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StreamsLookupCompatibleInterfaceFilterer{contract: contract}, nil
}

func bindStreamsLookupCompatibleInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StreamsLookupCompatibleInterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_StreamsLookupCompatibleInterface *StreamsLookupCompatibleInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StreamsLookupCompatibleInterface.Contract.StreamsLookupCompatibleInterfaceCaller.contract.Call(opts, result, method, params...)
}

func (_StreamsLookupCompatibleInterface *StreamsLookupCompatibleInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StreamsLookupCompatibleInterface.Contract.StreamsLookupCompatibleInterfaceTransactor.contract.Transfer(opts)
}

func (_StreamsLookupCompatibleInterface *StreamsLookupCompatibleInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StreamsLookupCompatibleInterface.Contract.StreamsLookupCompatibleInterfaceTransactor.contract.Transact(opts, method, params...)
}

func (_StreamsLookupCompatibleInterface *StreamsLookupCompatibleInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StreamsLookupCompatibleInterface.Contract.contract.Call(opts, result, method, params...)
}

func (_StreamsLookupCompatibleInterface *StreamsLookupCompatibleInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StreamsLookupCompatibleInterface.Contract.contract.Transfer(opts)
}

func (_StreamsLookupCompatibleInterface *StreamsLookupCompatibleInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StreamsLookupCompatibleInterface.Contract.contract.Transact(opts, method, params...)
}

func (_StreamsLookupCompatibleInterface *StreamsLookupCompatibleInterfaceCaller) CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	var out []interface{}
	err := _StreamsLookupCompatibleInterface.contract.Call(opts, &out, "checkCallback", values, extraData)

	outstruct := new(CheckCallback)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_StreamsLookupCompatibleInterface *StreamsLookupCompatibleInterfaceSession) CheckCallback(values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _StreamsLookupCompatibleInterface.Contract.CheckCallback(&_StreamsLookupCompatibleInterface.CallOpts, values, extraData)
}

func (_StreamsLookupCompatibleInterface *StreamsLookupCompatibleInterfaceCallerSession) CheckCallback(values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _StreamsLookupCompatibleInterface.Contract.CheckCallback(&_StreamsLookupCompatibleInterface.CallOpts, values, extraData)
}

func (_StreamsLookupCompatibleInterface *StreamsLookupCompatibleInterfaceCaller) CheckErrorHandler(opts *bind.CallOpts, errCode *big.Int, extraData []byte) (CheckErrorHandler,

	error) {
	var out []interface{}
	err := _StreamsLookupCompatibleInterface.contract.Call(opts, &out, "checkErrorHandler", errCode, extraData)

	outstruct := new(CheckErrorHandler)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_StreamsLookupCompatibleInterface *StreamsLookupCompatibleInterfaceSession) CheckErrorHandler(errCode *big.Int, extraData []byte) (CheckErrorHandler,

	error) {
	return _StreamsLookupCompatibleInterface.Contract.CheckErrorHandler(&_StreamsLookupCompatibleInterface.CallOpts, errCode, extraData)
}

func (_StreamsLookupCompatibleInterface *StreamsLookupCompatibleInterfaceCallerSession) CheckErrorHandler(errCode *big.Int, extraData []byte) (CheckErrorHandler,

	error) {
	return _StreamsLookupCompatibleInterface.Contract.CheckErrorHandler(&_StreamsLookupCompatibleInterface.CallOpts, errCode, extraData)
}

type CheckCallback struct {
	UpkeepNeeded bool
	PerformData  []byte
}
type CheckErrorHandler struct {
	UpkeepNeeded bool
	PerformData  []byte
}

func (_StreamsLookupCompatibleInterface *StreamsLookupCompatibleInterface) Address() common.Address {
	return _StreamsLookupCompatibleInterface.address
}

type StreamsLookupCompatibleInterfaceInterface interface {
	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (CheckCallback,

		error)

	CheckErrorHandler(opts *bind.CallOpts, errCode *big.Int, extraData []byte) (CheckErrorHandler,

		error)

	Address() common.Address
}
