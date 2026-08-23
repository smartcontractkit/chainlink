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

// CappedAdaptiveRateLogicMetaData contains all meta data concerning the CappedAdaptiveRateLogic contract.
var CappedAdaptiveRateLogicMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getAdaptiveRate\",\"inputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"referenceRate\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"marketRate\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"adaptiveRate\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"pure\"}]",
	Bin: "0x6080604052348015600f57600080fd5b5060ae80601d6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063f96978f414602d575b600080fd5b603c60383660046067565b604e565b60405190815260200160405180910390f35b6000848412605b5784605d565b835b9695505050505050565b600080600080600060a08688031215607e57600080fd5b50508335956020850135955060408501359460608101359450608001359250905056fea164736f6c634300081a000a",
}

// CappedAdaptiveRateLogicABI is the input ABI used to generate the binding from.
// Deprecated: Use CappedAdaptiveRateLogicMetaData.ABI instead.
var CappedAdaptiveRateLogicABI = CappedAdaptiveRateLogicMetaData.ABI

// CappedAdaptiveRateLogicBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CappedAdaptiveRateLogicMetaData.Bin instead.
var CappedAdaptiveRateLogicBin = CappedAdaptiveRateLogicMetaData.Bin

// DeployCappedAdaptiveRateLogic deploys a new Ethereum contract, binding an instance of CappedAdaptiveRateLogic to it.
func DeployCappedAdaptiveRateLogic(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CappedAdaptiveRateLogic, error) {
	parsed, err := CappedAdaptiveRateLogicMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CappedAdaptiveRateLogicBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CappedAdaptiveRateLogic{CappedAdaptiveRateLogicCaller: CappedAdaptiveRateLogicCaller{contract: contract}, CappedAdaptiveRateLogicTransactor: CappedAdaptiveRateLogicTransactor{contract: contract}, CappedAdaptiveRateLogicFilterer: CappedAdaptiveRateLogicFilterer{contract: contract}}, nil
}

// CappedAdaptiveRateLogic is an auto generated Go binding around an Ethereum contract.
type CappedAdaptiveRateLogic struct {
	CappedAdaptiveRateLogicCaller     // Read-only binding to the contract
	CappedAdaptiveRateLogicTransactor // Write-only binding to the contract
	CappedAdaptiveRateLogicFilterer   // Log filterer for contract events
}

// CappedAdaptiveRateLogicCaller is an auto generated read-only Go binding around an Ethereum contract.
type CappedAdaptiveRateLogicCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CappedAdaptiveRateLogicTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CappedAdaptiveRateLogicTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CappedAdaptiveRateLogicFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CappedAdaptiveRateLogicFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CappedAdaptiveRateLogicSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CappedAdaptiveRateLogicSession struct {
	Contract     *CappedAdaptiveRateLogic // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// CappedAdaptiveRateLogicCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CappedAdaptiveRateLogicCallerSession struct {
	Contract *CappedAdaptiveRateLogicCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// CappedAdaptiveRateLogicTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CappedAdaptiveRateLogicTransactorSession struct {
	Contract     *CappedAdaptiveRateLogicTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// CappedAdaptiveRateLogicRaw is an auto generated low-level Go binding around an Ethereum contract.
type CappedAdaptiveRateLogicRaw struct {
	Contract *CappedAdaptiveRateLogic // Generic contract binding to access the raw methods on
}

// CappedAdaptiveRateLogicCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CappedAdaptiveRateLogicCallerRaw struct {
	Contract *CappedAdaptiveRateLogicCaller // Generic read-only contract binding to access the raw methods on
}

// CappedAdaptiveRateLogicTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CappedAdaptiveRateLogicTransactorRaw struct {
	Contract *CappedAdaptiveRateLogicTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCappedAdaptiveRateLogic creates a new instance of CappedAdaptiveRateLogic, bound to a specific deployed contract.
func NewCappedAdaptiveRateLogic(address common.Address, backend bind.ContractBackend) (*CappedAdaptiveRateLogic, error) {
	contract, err := bindCappedAdaptiveRateLogic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CappedAdaptiveRateLogic{CappedAdaptiveRateLogicCaller: CappedAdaptiveRateLogicCaller{contract: contract}, CappedAdaptiveRateLogicTransactor: CappedAdaptiveRateLogicTransactor{contract: contract}, CappedAdaptiveRateLogicFilterer: CappedAdaptiveRateLogicFilterer{contract: contract}}, nil
}

// NewCappedAdaptiveRateLogicCaller creates a new read-only instance of CappedAdaptiveRateLogic, bound to a specific deployed contract.
func NewCappedAdaptiveRateLogicCaller(address common.Address, caller bind.ContractCaller) (*CappedAdaptiveRateLogicCaller, error) {
	contract, err := bindCappedAdaptiveRateLogic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CappedAdaptiveRateLogicCaller{contract: contract}, nil
}

// NewCappedAdaptiveRateLogicTransactor creates a new write-only instance of CappedAdaptiveRateLogic, bound to a specific deployed contract.
func NewCappedAdaptiveRateLogicTransactor(address common.Address, transactor bind.ContractTransactor) (*CappedAdaptiveRateLogicTransactor, error) {
	contract, err := bindCappedAdaptiveRateLogic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CappedAdaptiveRateLogicTransactor{contract: contract}, nil
}

// NewCappedAdaptiveRateLogicFilterer creates a new log filterer instance of CappedAdaptiveRateLogic, bound to a specific deployed contract.
func NewCappedAdaptiveRateLogicFilterer(address common.Address, filterer bind.ContractFilterer) (*CappedAdaptiveRateLogicFilterer, error) {
	contract, err := bindCappedAdaptiveRateLogic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CappedAdaptiveRateLogicFilterer{contract: contract}, nil
}

// bindCappedAdaptiveRateLogic binds a generic wrapper to an already deployed contract.
func bindCappedAdaptiveRateLogic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CappedAdaptiveRateLogicMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CappedAdaptiveRateLogic *CappedAdaptiveRateLogicRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CappedAdaptiveRateLogic.Contract.CappedAdaptiveRateLogicCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CappedAdaptiveRateLogic *CappedAdaptiveRateLogicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CappedAdaptiveRateLogic.Contract.CappedAdaptiveRateLogicTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CappedAdaptiveRateLogic *CappedAdaptiveRateLogicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CappedAdaptiveRateLogic.Contract.CappedAdaptiveRateLogicTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CappedAdaptiveRateLogic *CappedAdaptiveRateLogicCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CappedAdaptiveRateLogic.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CappedAdaptiveRateLogic *CappedAdaptiveRateLogicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CappedAdaptiveRateLogic.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CappedAdaptiveRateLogic *CappedAdaptiveRateLogicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CappedAdaptiveRateLogic.Contract.contract.Transact(opts, method, params...)
}

// GetAdaptiveRate is a free data retrieval call binding the contract method 0xf96978f4.
//
// Solidity: function getAdaptiveRate(int256 , int256 referenceRate, int256 marketRate, uint256 , uint256 ) pure returns(int256 adaptiveRate)
func (_CappedAdaptiveRateLogic *CappedAdaptiveRateLogicCaller) GetAdaptiveRate(opts *bind.CallOpts, arg0 *big.Int, referenceRate *big.Int, marketRate *big.Int, arg3 *big.Int, arg4 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _CappedAdaptiveRateLogic.contract.Call(opts, &out, "getAdaptiveRate", arg0, referenceRate, marketRate, arg3, arg4)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAdaptiveRate is a free data retrieval call binding the contract method 0xf96978f4.
//
// Solidity: function getAdaptiveRate(int256 , int256 referenceRate, int256 marketRate, uint256 , uint256 ) pure returns(int256 adaptiveRate)
func (_CappedAdaptiveRateLogic *CappedAdaptiveRateLogicSession) GetAdaptiveRate(arg0 *big.Int, referenceRate *big.Int, marketRate *big.Int, arg3 *big.Int, arg4 *big.Int) (*big.Int, error) {
	return _CappedAdaptiveRateLogic.Contract.GetAdaptiveRate(&_CappedAdaptiveRateLogic.CallOpts, arg0, referenceRate, marketRate, arg3, arg4)
}

// GetAdaptiveRate is a free data retrieval call binding the contract method 0xf96978f4.
//
// Solidity: function getAdaptiveRate(int256 , int256 referenceRate, int256 marketRate, uint256 , uint256 ) pure returns(int256 adaptiveRate)
func (_CappedAdaptiveRateLogic *CappedAdaptiveRateLogicCallerSession) GetAdaptiveRate(arg0 *big.Int, referenceRate *big.Int, marketRate *big.Int, arg3 *big.Int, arg4 *big.Int) (*big.Int, error) {
	return _CappedAdaptiveRateLogic.Contract.GetAdaptiveRate(&_CappedAdaptiveRateLogic.CallOpts, arg0, referenceRate, marketRate, arg3, arg4)
}
