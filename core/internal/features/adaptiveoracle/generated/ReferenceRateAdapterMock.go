// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package adaptiveoracle

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
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
	_ = time.Tick
	_ = context.Background
)

// ReferenceRateAdapterMockMetaData contains all meta data concerning the ReferenceRateAdapterMock contract.
var ReferenceRateAdapterMockMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"decimals_\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReferenceRate\",\"inputs\":[],\"outputs\":[{\"name\":\"rate\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"isSafe\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setRate\",\"inputs\":[{\"name\":\"rate\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"isSafe\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x60a06040526001805460ff191681179055348015601b57600080fd5b5060405161017c38038061017c8339810160408190526038916042565b60ff16608052606a565b600060208284031215605357600080fd5b815160ff81168114606357600080fd5b9392505050565b60805160fa6100826000396000606e015260fa6000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c8063078250ce146041578063313ce5671460675780637f70b3a114609d575b600080fd5b6065604c36600460ba565b6000919091556001805460ff1916911515919091179055565b005b60405160ff7f00000000000000000000000000000000000000000000000000000000000000001681526020015b60405180910390f35b60005460015460ff16604080519283529015156020830152016094565b6000806040838503121560cc57600080fd5b823591506020830135801515811460e257600080fd5b80915050925092905056fea164736f6c634300081a000a",
}

// ReferenceRateAdapterMockABI is the input ABI used to generate the binding from.
// Deprecated: Use ReferenceRateAdapterMockMetaData.ABI instead.
var ReferenceRateAdapterMockABI = ReferenceRateAdapterMockMetaData.ABI

// ReferenceRateAdapterMockBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ReferenceRateAdapterMockMetaData.Bin instead.
var ReferenceRateAdapterMockBin = ReferenceRateAdapterMockMetaData.Bin

// DeployReferenceRateAdapterMock deploys a new Ethereum contract, binding an instance of ReferenceRateAdapterMock to it.
func DeployReferenceRateAdapterMock(auth *bind.TransactOpts, backend bind.ContractBackend, decimals_ uint8) (common.Address, *types.Transaction, *ReferenceRateAdapterMock, error) {
	parsed, err := ReferenceRateAdapterMockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ReferenceRateAdapterMockBin), backend, decimals_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ReferenceRateAdapterMock{ReferenceRateAdapterMockCaller: ReferenceRateAdapterMockCaller{contract: contract}, ReferenceRateAdapterMockTransactor: ReferenceRateAdapterMockTransactor{contract: contract}, ReferenceRateAdapterMockFilterer: ReferenceRateAdapterMockFilterer{contract: contract}}, nil
}

// ReferenceRateAdapterMock is an auto generated Go binding around an Ethereum contract.
type ReferenceRateAdapterMock struct {
	ReferenceRateAdapterMockCaller     // Read-only binding to the contract
	ReferenceRateAdapterMockTransactor // Write-only binding to the contract
	ReferenceRateAdapterMockFilterer   // Log filterer for contract events
}

// ReferenceRateAdapterMockCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReferenceRateAdapterMockCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferenceRateAdapterMockTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReferenceRateAdapterMockTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferenceRateAdapterMockFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReferenceRateAdapterMockFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferenceRateAdapterMockSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReferenceRateAdapterMockSession struct {
	Contract     *ReferenceRateAdapterMock // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ReferenceRateAdapterMockCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReferenceRateAdapterMockCallerSession struct {
	Contract *ReferenceRateAdapterMockCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// ReferenceRateAdapterMockTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReferenceRateAdapterMockTransactorSession struct {
	Contract     *ReferenceRateAdapterMockTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// ReferenceRateAdapterMockRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReferenceRateAdapterMockRaw struct {
	Contract *ReferenceRateAdapterMock // Generic contract binding to access the raw methods on
}

// ReferenceRateAdapterMockCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReferenceRateAdapterMockCallerRaw struct {
	Contract *ReferenceRateAdapterMockCaller // Generic read-only contract binding to access the raw methods on
}

// ReferenceRateAdapterMockTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReferenceRateAdapterMockTransactorRaw struct {
	Contract *ReferenceRateAdapterMockTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReferenceRateAdapterMock creates a new instance of ReferenceRateAdapterMock, bound to a specific deployed contract.
func NewReferenceRateAdapterMock(address common.Address, backend bind.ContractBackend) (*ReferenceRateAdapterMock, error) {
	contract, err := bindReferenceRateAdapterMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReferenceRateAdapterMock{ReferenceRateAdapterMockCaller: ReferenceRateAdapterMockCaller{contract: contract}, ReferenceRateAdapterMockTransactor: ReferenceRateAdapterMockTransactor{contract: contract}, ReferenceRateAdapterMockFilterer: ReferenceRateAdapterMockFilterer{contract: contract}}, nil
}

// NewReferenceRateAdapterMockCaller creates a new read-only instance of ReferenceRateAdapterMock, bound to a specific deployed contract.
func NewReferenceRateAdapterMockCaller(address common.Address, caller bind.ContractCaller) (*ReferenceRateAdapterMockCaller, error) {
	contract, err := bindReferenceRateAdapterMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReferenceRateAdapterMockCaller{contract: contract}, nil
}

// NewReferenceRateAdapterMockTransactor creates a new write-only instance of ReferenceRateAdapterMock, bound to a specific deployed contract.
func NewReferenceRateAdapterMockTransactor(address common.Address, transactor bind.ContractTransactor) (*ReferenceRateAdapterMockTransactor, error) {
	contract, err := bindReferenceRateAdapterMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReferenceRateAdapterMockTransactor{contract: contract}, nil
}

// NewReferenceRateAdapterMockFilterer creates a new log filterer instance of ReferenceRateAdapterMock, bound to a specific deployed contract.
func NewReferenceRateAdapterMockFilterer(address common.Address, filterer bind.ContractFilterer) (*ReferenceRateAdapterMockFilterer, error) {
	contract, err := bindReferenceRateAdapterMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReferenceRateAdapterMockFilterer{contract: contract}, nil
}

// bindReferenceRateAdapterMock binds a generic wrapper to an already deployed contract.
func bindReferenceRateAdapterMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReferenceRateAdapterMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReferenceRateAdapterMock.Contract.ReferenceRateAdapterMockCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReferenceRateAdapterMock.Contract.ReferenceRateAdapterMockTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReferenceRateAdapterMock.Contract.ReferenceRateAdapterMockTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReferenceRateAdapterMock.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReferenceRateAdapterMock.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReferenceRateAdapterMock.Contract.contract.Transact(opts, method, params...)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ReferenceRateAdapterMock.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockSession) Decimals() (uint8, error) {
	return _ReferenceRateAdapterMock.Contract.Decimals(&_ReferenceRateAdapterMock.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockCallerSession) Decimals() (uint8, error) {
	return _ReferenceRateAdapterMock.Contract.Decimals(&_ReferenceRateAdapterMock.CallOpts)
}

// GetReferenceRate is a free data retrieval call binding the contract method 0x7f70b3a1.
//
// Solidity: function getReferenceRate() view returns(int256 rate, bool isSafe)
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockCaller) GetReferenceRate(opts *bind.CallOpts) (struct {
	Rate   *big.Int
	IsSafe bool
}, error) {
	var out []interface{}
	err := _ReferenceRateAdapterMock.contract.Call(opts, &out, "getReferenceRate")

	outstruct := new(struct {
		Rate   *big.Int
		IsSafe bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Rate = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.IsSafe = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// GetReferenceRate is a free data retrieval call binding the contract method 0x7f70b3a1.
//
// Solidity: function getReferenceRate() view returns(int256 rate, bool isSafe)
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockSession) GetReferenceRate() (struct {
	Rate   *big.Int
	IsSafe bool
}, error) {
	return _ReferenceRateAdapterMock.Contract.GetReferenceRate(&_ReferenceRateAdapterMock.CallOpts)
}

// GetReferenceRate is a free data retrieval call binding the contract method 0x7f70b3a1.
//
// Solidity: function getReferenceRate() view returns(int256 rate, bool isSafe)
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockCallerSession) GetReferenceRate() (struct {
	Rate   *big.Int
	IsSafe bool
}, error) {
	return _ReferenceRateAdapterMock.Contract.GetReferenceRate(&_ReferenceRateAdapterMock.CallOpts)
}

// SetRate is a paid mutator transaction binding the contract method 0x078250ce.
//
// Solidity: function setRate(int256 rate, bool isSafe) returns()
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockTransactor) SetRate(opts *bind.TransactOpts, rate *big.Int, isSafe bool) (*types.Transaction, error) {
	return _ReferenceRateAdapterMock.contract.Transact(opts, "setRate", rate, isSafe)
}

// SetRate is a paid mutator transaction binding the contract method 0x078250ce.
//
// Solidity: function setRate(int256 rate, bool isSafe) returns()
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockSession) SetRate(rate *big.Int, isSafe bool) (*types.Transaction, error) {
	return _ReferenceRateAdapterMock.Contract.SetRate(&_ReferenceRateAdapterMock.TransactOpts, rate, isSafe)
}

// SetRate is a paid mutator transaction binding the contract method 0x078250ce.
//
// Solidity: function setRate(int256 rate, bool isSafe) returns()
func (_ReferenceRateAdapterMock *ReferenceRateAdapterMockTransactorSession) SetRate(rate *big.Int, isSafe bool) (*types.Transaction, error) {
	return _ReferenceRateAdapterMock.Contract.SetRate(&_ReferenceRateAdapterMock.TransactOpts, rate, isSafe)
}
