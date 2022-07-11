// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2_wrapper_interface

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

var VRFV2WrapperInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"calculateRequestPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"_requestGasPriceWei\",\"type\":\"uint256\"}],\"name\":\"estimateRequestPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var VRFV2WrapperInterfaceABI = VRFV2WrapperInterfaceMetaData.ABI

type VRFV2WrapperInterface struct {
	address common.Address
	abi     abi.ABI
	VRFV2WrapperInterfaceCaller
	VRFV2WrapperInterfaceTransactor
	VRFV2WrapperInterfaceFilterer
}

type VRFV2WrapperInterfaceCaller struct {
	contract *bind.BoundContract
}

type VRFV2WrapperInterfaceTransactor struct {
	contract *bind.BoundContract
}

type VRFV2WrapperInterfaceFilterer struct {
	contract *bind.BoundContract
}

type VRFV2WrapperInterfaceSession struct {
	Contract     *VRFV2WrapperInterface
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2WrapperInterfaceCallerSession struct {
	Contract *VRFV2WrapperInterfaceCaller
	CallOpts bind.CallOpts
}

type VRFV2WrapperInterfaceTransactorSession struct {
	Contract     *VRFV2WrapperInterfaceTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2WrapperInterfaceRaw struct {
	Contract *VRFV2WrapperInterface
}

type VRFV2WrapperInterfaceCallerRaw struct {
	Contract *VRFV2WrapperInterfaceCaller
}

type VRFV2WrapperInterfaceTransactorRaw struct {
	Contract *VRFV2WrapperInterfaceTransactor
}

func NewVRFV2WrapperInterface(address common.Address, backend bind.ContractBackend) (*VRFV2WrapperInterface, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2WrapperInterfaceABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2WrapperInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperInterface{address: address, abi: abi, VRFV2WrapperInterfaceCaller: VRFV2WrapperInterfaceCaller{contract: contract}, VRFV2WrapperInterfaceTransactor: VRFV2WrapperInterfaceTransactor{contract: contract}, VRFV2WrapperInterfaceFilterer: VRFV2WrapperInterfaceFilterer{contract: contract}}, nil
}

func NewVRFV2WrapperInterfaceCaller(address common.Address, caller bind.ContractCaller) (*VRFV2WrapperInterfaceCaller, error) {
	contract, err := bindVRFV2WrapperInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperInterfaceCaller{contract: contract}, nil
}

func NewVRFV2WrapperInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2WrapperInterfaceTransactor, error) {
	contract, err := bindVRFV2WrapperInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperInterfaceTransactor{contract: contract}, nil
}

func NewVRFV2WrapperInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2WrapperInterfaceFilterer, error) {
	contract, err := bindVRFV2WrapperInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperInterfaceFilterer{contract: contract}, nil
}

func bindVRFV2WrapperInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFV2WrapperInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2WrapperInterface.Contract.VRFV2WrapperInterfaceCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2WrapperInterface.Contract.VRFV2WrapperInterfaceTransactor.contract.Transfer(opts)
}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2WrapperInterface.Contract.VRFV2WrapperInterfaceTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2WrapperInterface.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2WrapperInterface.Contract.contract.Transfer(opts)
}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2WrapperInterface.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceCaller) CalculateRequestPrice(opts *bind.CallOpts, _callbackGasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2WrapperInterface.contract.Call(opts, &out, "calculateRequestPrice", _callbackGasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceSession) CalculateRequestPrice(_callbackGasLimit uint32) (*big.Int, error) {
	return _VRFV2WrapperInterface.Contract.CalculateRequestPrice(&_VRFV2WrapperInterface.CallOpts, _callbackGasLimit)
}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceCallerSession) CalculateRequestPrice(_callbackGasLimit uint32) (*big.Int, error) {
	return _VRFV2WrapperInterface.Contract.CalculateRequestPrice(&_VRFV2WrapperInterface.CallOpts, _callbackGasLimit)
}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceCaller) EstimateRequestPrice(opts *bind.CallOpts, _callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2WrapperInterface.contract.Call(opts, &out, "estimateRequestPrice", _callbackGasLimit, _requestGasPriceWei)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceSession) EstimateRequestPrice(_callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	return _VRFV2WrapperInterface.Contract.EstimateRequestPrice(&_VRFV2WrapperInterface.CallOpts, _callbackGasLimit, _requestGasPriceWei)
}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceCallerSession) EstimateRequestPrice(_callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	return _VRFV2WrapperInterface.Contract.EstimateRequestPrice(&_VRFV2WrapperInterface.CallOpts, _callbackGasLimit, _requestGasPriceWei)
}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceCaller) LastRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2WrapperInterface.contract.Call(opts, &out, "lastRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceSession) LastRequestId() (*big.Int, error) {
	return _VRFV2WrapperInterface.Contract.LastRequestId(&_VRFV2WrapperInterface.CallOpts)
}

func (_VRFV2WrapperInterface *VRFV2WrapperInterfaceCallerSession) LastRequestId() (*big.Int, error) {
	return _VRFV2WrapperInterface.Contract.LastRequestId(&_VRFV2WrapperInterface.CallOpts)
}

func (_VRFV2WrapperInterface *VRFV2WrapperInterface) Address() common.Address {
	return _VRFV2WrapperInterface.address
}

type VRFV2WrapperInterfaceInterface interface {
	CalculateRequestPrice(opts *bind.CallOpts, _callbackGasLimit uint32) (*big.Int, error)

	EstimateRequestPrice(opts *bind.CallOpts, _callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error)

	LastRequestId(opts *bind.CallOpts) (*big.Int, error)

	Address() common.Address
}
