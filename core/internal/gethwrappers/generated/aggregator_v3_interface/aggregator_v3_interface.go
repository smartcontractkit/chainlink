// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package aggregator_v3_interface

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

var IAggregatorV3MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var IAggregatorV3ABI = IAggregatorV3MetaData.ABI

type IAggregatorV3 struct {
	address common.Address
	abi     abi.ABI
	IAggregatorV3Caller
	IAggregatorV3Transactor
	IAggregatorV3Filterer
}

type IAggregatorV3Caller struct {
	contract *bind.BoundContract
}

type IAggregatorV3Transactor struct {
	contract *bind.BoundContract
}

type IAggregatorV3Filterer struct {
	contract *bind.BoundContract
}

type IAggregatorV3Session struct {
	Contract     *IAggregatorV3
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type IAggregatorV3CallerSession struct {
	Contract *IAggregatorV3Caller
	CallOpts bind.CallOpts
}

type IAggregatorV3TransactorSession struct {
	Contract     *IAggregatorV3Transactor
	TransactOpts bind.TransactOpts
}

type IAggregatorV3Raw struct {
	Contract *IAggregatorV3
}

type IAggregatorV3CallerRaw struct {
	Contract *IAggregatorV3Caller
}

type IAggregatorV3TransactorRaw struct {
	Contract *IAggregatorV3Transactor
}

func NewIAggregatorV3(address common.Address, backend bind.ContractBackend) (*IAggregatorV3, error) {
	abi, err := abi.JSON(strings.NewReader(IAggregatorV3ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindIAggregatorV3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAggregatorV3{address: address, abi: abi, IAggregatorV3Caller: IAggregatorV3Caller{contract: contract}, IAggregatorV3Transactor: IAggregatorV3Transactor{contract: contract}, IAggregatorV3Filterer: IAggregatorV3Filterer{contract: contract}}, nil
}

func NewIAggregatorV3Caller(address common.Address, caller bind.ContractCaller) (*IAggregatorV3Caller, error) {
	contract, err := bindIAggregatorV3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAggregatorV3Caller{contract: contract}, nil
}

func NewIAggregatorV3Transactor(address common.Address, transactor bind.ContractTransactor) (*IAggregatorV3Transactor, error) {
	contract, err := bindIAggregatorV3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAggregatorV3Transactor{contract: contract}, nil
}

func NewIAggregatorV3Filterer(address common.Address, filterer bind.ContractFilterer) (*IAggregatorV3Filterer, error) {
	contract, err := bindIAggregatorV3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAggregatorV3Filterer{contract: contract}, nil
}

func bindIAggregatorV3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IAggregatorV3ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_IAggregatorV3 *IAggregatorV3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAggregatorV3.Contract.IAggregatorV3Caller.contract.Call(opts, result, method, params...)
}

func (_IAggregatorV3 *IAggregatorV3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAggregatorV3.Contract.IAggregatorV3Transactor.contract.Transfer(opts)
}

func (_IAggregatorV3 *IAggregatorV3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAggregatorV3.Contract.IAggregatorV3Transactor.contract.Transact(opts, method, params...)
}

func (_IAggregatorV3 *IAggregatorV3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAggregatorV3.Contract.contract.Call(opts, result, method, params...)
}

func (_IAggregatorV3 *IAggregatorV3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAggregatorV3.Contract.contract.Transfer(opts)
}

func (_IAggregatorV3 *IAggregatorV3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAggregatorV3.Contract.contract.Transact(opts, method, params...)
}

func (_IAggregatorV3 *IAggregatorV3Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IAggregatorV3.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IAggregatorV3 *IAggregatorV3Session) Decimals() (uint8, error) {
	return _IAggregatorV3.Contract.Decimals(&_IAggregatorV3.CallOpts)
}

func (_IAggregatorV3 *IAggregatorV3CallerSession) Decimals() (uint8, error) {
	return _IAggregatorV3.Contract.Decimals(&_IAggregatorV3.CallOpts)
}

func (_IAggregatorV3 *IAggregatorV3Caller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IAggregatorV3.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_IAggregatorV3 *IAggregatorV3Session) Description() (string, error) {
	return _IAggregatorV3.Contract.Description(&_IAggregatorV3.CallOpts)
}

func (_IAggregatorV3 *IAggregatorV3CallerSession) Description() (string, error) {
	return _IAggregatorV3.Contract.Description(&_IAggregatorV3.CallOpts)
}

func (_IAggregatorV3 *IAggregatorV3Caller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _IAggregatorV3.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(GetRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_IAggregatorV3 *IAggregatorV3Session) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _IAggregatorV3.Contract.GetRoundData(&_IAggregatorV3.CallOpts, _roundId)
}

func (_IAggregatorV3 *IAggregatorV3CallerSession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _IAggregatorV3.Contract.GetRoundData(&_IAggregatorV3.CallOpts, _roundId)
}

func (_IAggregatorV3 *IAggregatorV3Caller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _IAggregatorV3.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(LatestRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_IAggregatorV3 *IAggregatorV3Session) LatestRoundData() (LatestRoundData,

	error) {
	return _IAggregatorV3.Contract.LatestRoundData(&_IAggregatorV3.CallOpts)
}

func (_IAggregatorV3 *IAggregatorV3CallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _IAggregatorV3.Contract.LatestRoundData(&_IAggregatorV3.CallOpts)
}

func (_IAggregatorV3 *IAggregatorV3Caller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAggregatorV3.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAggregatorV3 *IAggregatorV3Session) Version() (*big.Int, error) {
	return _IAggregatorV3.Contract.Version(&_IAggregatorV3.CallOpts)
}

func (_IAggregatorV3 *IAggregatorV3CallerSession) Version() (*big.Int, error) {
	return _IAggregatorV3.Contract.Version(&_IAggregatorV3.CallOpts)
}

type GetRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
type LatestRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

func (_IAggregatorV3 *IAggregatorV3) Address() common.Address {
	return _IAggregatorV3.address
}

type IAggregatorV3Interface interface {
	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

		error)

	LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

		error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	Address() common.Address
}
