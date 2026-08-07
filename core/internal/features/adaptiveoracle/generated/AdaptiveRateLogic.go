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

// AdaptiveRateLogicMetaData contains all meta data concerning the AdaptiveRateLogic contract.
var AdaptiveRateLogicMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getAdaptiveRate\",\"inputs\":[{\"name\":\"lastAdaptiveRate\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"marketRate\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"adaptiveRate\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"pure\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506101398061001f6000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063f96978f414610030575b600080fd5b61004361003e366004610077565b610055565b60405190815260200160405180910390f35b6000600261006385886100c8565b61006d91906100f0565b9695505050505050565b600080600080600060a0868803121561008f57600080fd5b505083359560208501359550604085013594606081013594506080013592509050565b634e487b7160e01b600052601160045260246000fd5b80820182811260008312801582168215821617156100e8576100e86100b2565b505092915050565b60008261010d57634e487b7160e01b600052601260045260246000fd5b600160ff1b821460001984141615610127576101276100b2565b50059056fea164736f6c634300081a000a",
}

// AdaptiveRateLogicABI is the input ABI used to generate the binding from.
// Deprecated: Use AdaptiveRateLogicMetaData.ABI instead.
var AdaptiveRateLogicABI = AdaptiveRateLogicMetaData.ABI

// AdaptiveRateLogicBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AdaptiveRateLogicMetaData.Bin instead.
var AdaptiveRateLogicBin = AdaptiveRateLogicMetaData.Bin

// DeployAdaptiveRateLogic deploys a new Ethereum contract, binding an instance of AdaptiveRateLogic to it.
func DeployAdaptiveRateLogic(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AdaptiveRateLogic, error) {
	parsed, err := AdaptiveRateLogicMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AdaptiveRateLogicBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AdaptiveRateLogic{AdaptiveRateLogicCaller: AdaptiveRateLogicCaller{contract: contract}, AdaptiveRateLogicTransactor: AdaptiveRateLogicTransactor{contract: contract}, AdaptiveRateLogicFilterer: AdaptiveRateLogicFilterer{contract: contract}}, nil
}

// AdaptiveRateLogic is an auto generated Go binding around an Ethereum contract.
type AdaptiveRateLogic struct {
	AdaptiveRateLogicCaller     // Read-only binding to the contract
	AdaptiveRateLogicTransactor // Write-only binding to the contract
	AdaptiveRateLogicFilterer   // Log filterer for contract events
}

// AdaptiveRateLogicCaller is an auto generated read-only Go binding around an Ethereum contract.
type AdaptiveRateLogicCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdaptiveRateLogicTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AdaptiveRateLogicTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdaptiveRateLogicFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AdaptiveRateLogicFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdaptiveRateLogicSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AdaptiveRateLogicSession struct {
	Contract     *AdaptiveRateLogic // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// AdaptiveRateLogicCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AdaptiveRateLogicCallerSession struct {
	Contract *AdaptiveRateLogicCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// AdaptiveRateLogicTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AdaptiveRateLogicTransactorSession struct {
	Contract     *AdaptiveRateLogicTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// AdaptiveRateLogicRaw is an auto generated low-level Go binding around an Ethereum contract.
type AdaptiveRateLogicRaw struct {
	Contract *AdaptiveRateLogic // Generic contract binding to access the raw methods on
}

// AdaptiveRateLogicCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AdaptiveRateLogicCallerRaw struct {
	Contract *AdaptiveRateLogicCaller // Generic read-only contract binding to access the raw methods on
}

// AdaptiveRateLogicTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AdaptiveRateLogicTransactorRaw struct {
	Contract *AdaptiveRateLogicTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAdaptiveRateLogic creates a new instance of AdaptiveRateLogic, bound to a specific deployed contract.
func NewAdaptiveRateLogic(address common.Address, backend bind.ContractBackend) (*AdaptiveRateLogic, error) {
	contract, err := bindAdaptiveRateLogic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AdaptiveRateLogic{AdaptiveRateLogicCaller: AdaptiveRateLogicCaller{contract: contract}, AdaptiveRateLogicTransactor: AdaptiveRateLogicTransactor{contract: contract}, AdaptiveRateLogicFilterer: AdaptiveRateLogicFilterer{contract: contract}}, nil
}

// NewAdaptiveRateLogicCaller creates a new read-only instance of AdaptiveRateLogic, bound to a specific deployed contract.
func NewAdaptiveRateLogicCaller(address common.Address, caller bind.ContractCaller) (*AdaptiveRateLogicCaller, error) {
	contract, err := bindAdaptiveRateLogic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AdaptiveRateLogicCaller{contract: contract}, nil
}

// NewAdaptiveRateLogicTransactor creates a new write-only instance of AdaptiveRateLogic, bound to a specific deployed contract.
func NewAdaptiveRateLogicTransactor(address common.Address, transactor bind.ContractTransactor) (*AdaptiveRateLogicTransactor, error) {
	contract, err := bindAdaptiveRateLogic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AdaptiveRateLogicTransactor{contract: contract}, nil
}

// NewAdaptiveRateLogicFilterer creates a new log filterer instance of AdaptiveRateLogic, bound to a specific deployed contract.
func NewAdaptiveRateLogicFilterer(address common.Address, filterer bind.ContractFilterer) (*AdaptiveRateLogicFilterer, error) {
	contract, err := bindAdaptiveRateLogic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AdaptiveRateLogicFilterer{contract: contract}, nil
}

// bindAdaptiveRateLogic binds a generic wrapper to an already deployed contract.
func bindAdaptiveRateLogic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AdaptiveRateLogicMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AdaptiveRateLogic *AdaptiveRateLogicRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AdaptiveRateLogic.Contract.AdaptiveRateLogicCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AdaptiveRateLogic *AdaptiveRateLogicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AdaptiveRateLogic.Contract.AdaptiveRateLogicTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AdaptiveRateLogic *AdaptiveRateLogicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AdaptiveRateLogic.Contract.AdaptiveRateLogicTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AdaptiveRateLogic *AdaptiveRateLogicCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AdaptiveRateLogic.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AdaptiveRateLogic *AdaptiveRateLogicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AdaptiveRateLogic.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AdaptiveRateLogic *AdaptiveRateLogicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AdaptiveRateLogic.Contract.contract.Transact(opts, method, params...)
}

// GetAdaptiveRate is a free data retrieval call binding the contract method 0xf96978f4.
//
// Solidity: function getAdaptiveRate(int256 lastAdaptiveRate, int256 , int256 marketRate, uint256 , uint256 ) pure returns(int256 adaptiveRate)
func (_AdaptiveRateLogic *AdaptiveRateLogicCaller) GetAdaptiveRate(opts *bind.CallOpts, lastAdaptiveRate *big.Int, arg1 *big.Int, marketRate *big.Int, arg3 *big.Int, arg4 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AdaptiveRateLogic.contract.Call(opts, &out, "getAdaptiveRate", lastAdaptiveRate, arg1, marketRate, arg3, arg4)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAdaptiveRate is a free data retrieval call binding the contract method 0xf96978f4.
//
// Solidity: function getAdaptiveRate(int256 lastAdaptiveRate, int256 , int256 marketRate, uint256 , uint256 ) pure returns(int256 adaptiveRate)
func (_AdaptiveRateLogic *AdaptiveRateLogicSession) GetAdaptiveRate(lastAdaptiveRate *big.Int, arg1 *big.Int, marketRate *big.Int, arg3 *big.Int, arg4 *big.Int) (*big.Int, error) {
	return _AdaptiveRateLogic.Contract.GetAdaptiveRate(&_AdaptiveRateLogic.CallOpts, lastAdaptiveRate, arg1, marketRate, arg3, arg4)
}

// GetAdaptiveRate is a free data retrieval call binding the contract method 0xf96978f4.
//
// Solidity: function getAdaptiveRate(int256 lastAdaptiveRate, int256 , int256 marketRate, uint256 , uint256 ) pure returns(int256 adaptiveRate)
func (_AdaptiveRateLogic *AdaptiveRateLogicCallerSession) GetAdaptiveRate(lastAdaptiveRate *big.Int, arg1 *big.Int, marketRate *big.Int, arg3 *big.Int, arg4 *big.Int) (*big.Int, error) {
	return _AdaptiveRateLogic.Contract.GetAdaptiveRate(&_AdaptiveRateLogic.CallOpts, lastAdaptiveRate, arg1, marketRate, arg3, arg4)
}
