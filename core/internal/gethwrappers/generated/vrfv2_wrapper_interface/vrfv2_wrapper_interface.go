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

var IVRFV2WrapperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"calculateRequestPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"_requestGasPriceWei\",\"type\":\"uint256\"}],\"name\":\"estimateRequestPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var IVRFV2WrapperABI = IVRFV2WrapperMetaData.ABI

type IVRFV2Wrapper struct {
	address common.Address
	abi     abi.ABI
	IVRFV2WrapperCaller
	IVRFV2WrapperTransactor
	IVRFV2WrapperFilterer
}

type IVRFV2WrapperCaller struct {
	contract *bind.BoundContract
}

type IVRFV2WrapperTransactor struct {
	contract *bind.BoundContract
}

type IVRFV2WrapperFilterer struct {
	contract *bind.BoundContract
}

type IVRFV2WrapperSession struct {
	Contract     *IVRFV2Wrapper
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type IVRFV2WrapperCallerSession struct {
	Contract *IVRFV2WrapperCaller
	CallOpts bind.CallOpts
}

type IVRFV2WrapperTransactorSession struct {
	Contract     *IVRFV2WrapperTransactor
	TransactOpts bind.TransactOpts
}

type IVRFV2WrapperRaw struct {
	Contract *IVRFV2Wrapper
}

type IVRFV2WrapperCallerRaw struct {
	Contract *IVRFV2WrapperCaller
}

type IVRFV2WrapperTransactorRaw struct {
	Contract *IVRFV2WrapperTransactor
}

func NewIVRFV2Wrapper(address common.Address, backend bind.ContractBackend) (*IVRFV2Wrapper, error) {
	abi, err := abi.JSON(strings.NewReader(IVRFV2WrapperABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindIVRFV2Wrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IVRFV2Wrapper{address: address, abi: abi, IVRFV2WrapperCaller: IVRFV2WrapperCaller{contract: contract}, IVRFV2WrapperTransactor: IVRFV2WrapperTransactor{contract: contract}, IVRFV2WrapperFilterer: IVRFV2WrapperFilterer{contract: contract}}, nil
}

func NewIVRFV2WrapperCaller(address common.Address, caller bind.ContractCaller) (*IVRFV2WrapperCaller, error) {
	contract, err := bindIVRFV2Wrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IVRFV2WrapperCaller{contract: contract}, nil
}

func NewIVRFV2WrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*IVRFV2WrapperTransactor, error) {
	contract, err := bindIVRFV2Wrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IVRFV2WrapperTransactor{contract: contract}, nil
}

func NewIVRFV2WrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*IVRFV2WrapperFilterer, error) {
	contract, err := bindIVRFV2Wrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IVRFV2WrapperFilterer{contract: contract}, nil
}

func bindIVRFV2Wrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IVRFV2WrapperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_IVRFV2Wrapper *IVRFV2WrapperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IVRFV2Wrapper.Contract.IVRFV2WrapperCaller.contract.Call(opts, result, method, params...)
}

func (_IVRFV2Wrapper *IVRFV2WrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IVRFV2Wrapper.Contract.IVRFV2WrapperTransactor.contract.Transfer(opts)
}

func (_IVRFV2Wrapper *IVRFV2WrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IVRFV2Wrapper.Contract.IVRFV2WrapperTransactor.contract.Transact(opts, method, params...)
}

func (_IVRFV2Wrapper *IVRFV2WrapperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IVRFV2Wrapper.Contract.contract.Call(opts, result, method, params...)
}

func (_IVRFV2Wrapper *IVRFV2WrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IVRFV2Wrapper.Contract.contract.Transfer(opts)
}

func (_IVRFV2Wrapper *IVRFV2WrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IVRFV2Wrapper.Contract.contract.Transact(opts, method, params...)
}

func (_IVRFV2Wrapper *IVRFV2WrapperCaller) CalculateRequestPrice(opts *bind.CallOpts, _callbackGasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _IVRFV2Wrapper.contract.Call(opts, &out, "calculateRequestPrice", _callbackGasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IVRFV2Wrapper *IVRFV2WrapperSession) CalculateRequestPrice(_callbackGasLimit uint32) (*big.Int, error) {
	return _IVRFV2Wrapper.Contract.CalculateRequestPrice(&_IVRFV2Wrapper.CallOpts, _callbackGasLimit)
}

func (_IVRFV2Wrapper *IVRFV2WrapperCallerSession) CalculateRequestPrice(_callbackGasLimit uint32) (*big.Int, error) {
	return _IVRFV2Wrapper.Contract.CalculateRequestPrice(&_IVRFV2Wrapper.CallOpts, _callbackGasLimit)
}

func (_IVRFV2Wrapper *IVRFV2WrapperCaller) EstimateRequestPrice(opts *bind.CallOpts, _callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IVRFV2Wrapper.contract.Call(opts, &out, "estimateRequestPrice", _callbackGasLimit, _requestGasPriceWei)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IVRFV2Wrapper *IVRFV2WrapperSession) EstimateRequestPrice(_callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	return _IVRFV2Wrapper.Contract.EstimateRequestPrice(&_IVRFV2Wrapper.CallOpts, _callbackGasLimit, _requestGasPriceWei)
}

func (_IVRFV2Wrapper *IVRFV2WrapperCallerSession) EstimateRequestPrice(_callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	return _IVRFV2Wrapper.Contract.EstimateRequestPrice(&_IVRFV2Wrapper.CallOpts, _callbackGasLimit, _requestGasPriceWei)
}

func (_IVRFV2Wrapper *IVRFV2WrapperCaller) LastRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IVRFV2Wrapper.contract.Call(opts, &out, "lastRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IVRFV2Wrapper *IVRFV2WrapperSession) LastRequestId() (*big.Int, error) {
	return _IVRFV2Wrapper.Contract.LastRequestId(&_IVRFV2Wrapper.CallOpts)
}

func (_IVRFV2Wrapper *IVRFV2WrapperCallerSession) LastRequestId() (*big.Int, error) {
	return _IVRFV2Wrapper.Contract.LastRequestId(&_IVRFV2Wrapper.CallOpts)
}

func (_IVRFV2Wrapper *IVRFV2Wrapper) Address() common.Address {
	return _IVRFV2Wrapper.address
}

type IVRFV2WrapperInterface interface {
	CalculateRequestPrice(opts *bind.CallOpts, _callbackGasLimit uint32) (*big.Int, error)

	EstimateRequestPrice(opts *bind.CallOpts, _callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error)

	LastRequestId(opts *bind.CallOpts) (*big.Int, error)

	Address() common.Address
}
