// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package exposedocr2aggregator

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
)

// AccessControllerInterfaceMetaData contains all meta data concerning the AccessControllerInterface contract.
var AccessControllerInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AccessControllerInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use AccessControllerInterfaceMetaData.ABI instead.
var AccessControllerInterfaceABI = AccessControllerInterfaceMetaData.ABI

// AccessControllerInterface is an auto generated Go binding around an Ethereum contract.
type AccessControllerInterface struct {
	AccessControllerInterfaceCaller     // Read-only binding to the contract
	AccessControllerInterfaceTransactor // Write-only binding to the contract
	AccessControllerInterfaceFilterer   // Log filterer for contract events
}

// AccessControllerInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccessControllerInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControllerInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccessControllerInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControllerInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccessControllerInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControllerInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccessControllerInterfaceSession struct {
	Contract     *AccessControllerInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// AccessControllerInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccessControllerInterfaceCallerSession struct {
	Contract *AccessControllerInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// AccessControllerInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccessControllerInterfaceTransactorSession struct {
	Contract     *AccessControllerInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// AccessControllerInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccessControllerInterfaceRaw struct {
	Contract *AccessControllerInterface // Generic contract binding to access the raw methods on
}

// AccessControllerInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccessControllerInterfaceCallerRaw struct {
	Contract *AccessControllerInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// AccessControllerInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccessControllerInterfaceTransactorRaw struct {
	Contract *AccessControllerInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccessControllerInterface creates a new instance of AccessControllerInterface, bound to a specific deployed contract.
func NewAccessControllerInterface(address common.Address, backend bind.ContractBackend) (*AccessControllerInterface, error) {
	contract, err := bindAccessControllerInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccessControllerInterface{AccessControllerInterfaceCaller: AccessControllerInterfaceCaller{contract: contract}, AccessControllerInterfaceTransactor: AccessControllerInterfaceTransactor{contract: contract}, AccessControllerInterfaceFilterer: AccessControllerInterfaceFilterer{contract: contract}}, nil
}

// NewAccessControllerInterfaceCaller creates a new read-only instance of AccessControllerInterface, bound to a specific deployed contract.
func NewAccessControllerInterfaceCaller(address common.Address, caller bind.ContractCaller) (*AccessControllerInterfaceCaller, error) {
	contract, err := bindAccessControllerInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControllerInterfaceCaller{contract: contract}, nil
}

// NewAccessControllerInterfaceTransactor creates a new write-only instance of AccessControllerInterface, bound to a specific deployed contract.
func NewAccessControllerInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AccessControllerInterfaceTransactor, error) {
	contract, err := bindAccessControllerInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControllerInterfaceTransactor{contract: contract}, nil
}

// NewAccessControllerInterfaceFilterer creates a new log filterer instance of AccessControllerInterface, bound to a specific deployed contract.
func NewAccessControllerInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AccessControllerInterfaceFilterer, error) {
	contract, err := bindAccessControllerInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccessControllerInterfaceFilterer{contract: contract}, nil
}

// bindAccessControllerInterface binds a generic wrapper to an already deployed contract.
func bindAccessControllerInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AccessControllerInterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessControllerInterface *AccessControllerInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControllerInterface.Contract.AccessControllerInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessControllerInterface *AccessControllerInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.AccessControllerInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessControllerInterface *AccessControllerInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.AccessControllerInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessControllerInterface *AccessControllerInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControllerInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessControllerInterface *AccessControllerInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessControllerInterface *AccessControllerInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.contract.Transact(opts, method, params...)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address user, bytes data) view returns(bool)
func (_AccessControllerInterface *AccessControllerInterfaceCaller) HasAccess(opts *bind.CallOpts, user common.Address, data []byte) (bool, error) {
	var out []interface{}
	err := _AccessControllerInterface.contract.Call(opts, &out, "hasAccess", user, data)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address user, bytes data) view returns(bool)
func (_AccessControllerInterface *AccessControllerInterfaceSession) HasAccess(user common.Address, data []byte) (bool, error) {
	return _AccessControllerInterface.Contract.HasAccess(&_AccessControllerInterface.CallOpts, user, data)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address user, bytes data) view returns(bool)
func (_AccessControllerInterface *AccessControllerInterfaceCallerSession) HasAccess(user common.Address, data []byte) (bool, error) {
	return _AccessControllerInterface.Contract.HasAccess(&_AccessControllerInterface.CallOpts, user, data)
}

// AggregatorInterfaceMetaData contains all meta data concerning the AggregatorInterface contract.
var AggregatorInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AggregatorInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use AggregatorInterfaceMetaData.ABI instead.
var AggregatorInterfaceABI = AggregatorInterfaceMetaData.ABI

// AggregatorInterface is an auto generated Go binding around an Ethereum contract.
type AggregatorInterface struct {
	AggregatorInterfaceCaller     // Read-only binding to the contract
	AggregatorInterfaceTransactor // Write-only binding to the contract
	AggregatorInterfaceFilterer   // Log filterer for contract events
}

// AggregatorInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type AggregatorInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AggregatorInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AggregatorInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AggregatorInterfaceSession struct {
	Contract     *AggregatorInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// AggregatorInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AggregatorInterfaceCallerSession struct {
	Contract *AggregatorInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// AggregatorInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AggregatorInterfaceTransactorSession struct {
	Contract     *AggregatorInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// AggregatorInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type AggregatorInterfaceRaw struct {
	Contract *AggregatorInterface // Generic contract binding to access the raw methods on
}

// AggregatorInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AggregatorInterfaceCallerRaw struct {
	Contract *AggregatorInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// AggregatorInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AggregatorInterfaceTransactorRaw struct {
	Contract *AggregatorInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAggregatorInterface creates a new instance of AggregatorInterface, bound to a specific deployed contract.
func NewAggregatorInterface(address common.Address, backend bind.ContractBackend) (*AggregatorInterface, error) {
	contract, err := bindAggregatorInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterface{AggregatorInterfaceCaller: AggregatorInterfaceCaller{contract: contract}, AggregatorInterfaceTransactor: AggregatorInterfaceTransactor{contract: contract}, AggregatorInterfaceFilterer: AggregatorInterfaceFilterer{contract: contract}}, nil
}

// NewAggregatorInterfaceCaller creates a new read-only instance of AggregatorInterface, bound to a specific deployed contract.
func NewAggregatorInterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorInterfaceCaller, error) {
	contract, err := bindAggregatorInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterfaceCaller{contract: contract}, nil
}

// NewAggregatorInterfaceTransactor creates a new write-only instance of AggregatorInterface, bound to a specific deployed contract.
func NewAggregatorInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorInterfaceTransactor, error) {
	contract, err := bindAggregatorInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterfaceTransactor{contract: contract}, nil
}

// NewAggregatorInterfaceFilterer creates a new log filterer instance of AggregatorInterface, bound to a specific deployed contract.
func NewAggregatorInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorInterfaceFilterer, error) {
	contract, err := bindAggregatorInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterfaceFilterer{contract: contract}, nil
}

// bindAggregatorInterface binds a generic wrapper to an already deployed contract.
func bindAggregatorInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AggregatorInterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorInterface *AggregatorInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorInterface.Contract.AggregatorInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorInterface *AggregatorInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.AggregatorInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorInterface *AggregatorInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.AggregatorInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorInterface *AggregatorInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorInterface *AggregatorInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorInterface *AggregatorInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.contract.Transact(opts, method, params...)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AggregatorInterface *AggregatorInterfaceCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AggregatorInterface *AggregatorInterfaceSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetAnswer(&_AggregatorInterface.CallOpts, roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetAnswer(&_AggregatorInterface.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetTimestamp(&_AggregatorInterface.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetTimestamp(&_AggregatorInterface.CallOpts, roundId)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorInterface *AggregatorInterfaceCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorInterface *AggregatorInterfaceSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestAnswer(&_AggregatorInterface.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestAnswer(&_AggregatorInterface.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceSession) LatestRound() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestRound(&_AggregatorInterface.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestRound() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestRound(&_AggregatorInterface.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestTimestamp(&_AggregatorInterface.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestTimestamp(&_AggregatorInterface.CallOpts)
}

// AggregatorInterfaceAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the AggregatorInterface contract.
type AggregatorInterfaceAnswerUpdatedIterator struct {
	Event *AggregatorInterfaceAnswerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AggregatorInterfaceAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorInterfaceAnswerUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AggregatorInterfaceAnswerUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AggregatorInterfaceAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorInterfaceAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorInterfaceAnswerUpdated represents a AnswerUpdated event raised by the AggregatorInterface contract.
type AggregatorInterfaceAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AggregatorInterface *AggregatorInterfaceFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*AggregatorInterfaceAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorInterface.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterfaceAnswerUpdatedIterator{contract: _AggregatorInterface.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AggregatorInterface *AggregatorInterfaceFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *AggregatorInterfaceAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorInterface.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AggregatorInterfaceAnswerUpdated)
				if err := _AggregatorInterface.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AggregatorInterface *AggregatorInterfaceFilterer) ParseAnswerUpdated(log types.Log) (*AggregatorInterfaceAnswerUpdated, error) {
	event := new(AggregatorInterfaceAnswerUpdated)
	if err := _AggregatorInterface.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AggregatorInterfaceNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the AggregatorInterface contract.
type AggregatorInterfaceNewRoundIterator struct {
	Event *AggregatorInterfaceNewRound // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AggregatorInterfaceNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorInterfaceNewRound)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AggregatorInterfaceNewRound)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AggregatorInterfaceNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorInterfaceNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorInterfaceNewRound represents a NewRound event raised by the AggregatorInterface contract.
type AggregatorInterfaceNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AggregatorInterface *AggregatorInterfaceFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*AggregatorInterfaceNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorInterface.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterfaceNewRoundIterator{contract: _AggregatorInterface.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AggregatorInterface *AggregatorInterfaceFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *AggregatorInterfaceNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorInterface.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AggregatorInterfaceNewRound)
				if err := _AggregatorInterface.contract.UnpackLog(event, "NewRound", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AggregatorInterface *AggregatorInterfaceFilterer) ParseNewRound(log types.Log) (*AggregatorInterfaceNewRound, error) {
	event := new(AggregatorInterfaceNewRound)
	if err := _AggregatorInterface.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AggregatorV2V3InterfaceMetaData contains all meta data concerning the AggregatorV2V3Interface contract.
var AggregatorV2V3InterfaceMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AggregatorV2V3InterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use AggregatorV2V3InterfaceMetaData.ABI instead.
var AggregatorV2V3InterfaceABI = AggregatorV2V3InterfaceMetaData.ABI

// AggregatorV2V3Interface is an auto generated Go binding around an Ethereum contract.
type AggregatorV2V3Interface struct {
	AggregatorV2V3InterfaceCaller     // Read-only binding to the contract
	AggregatorV2V3InterfaceTransactor // Write-only binding to the contract
	AggregatorV2V3InterfaceFilterer   // Log filterer for contract events
}

// AggregatorV2V3InterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type AggregatorV2V3InterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorV2V3InterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AggregatorV2V3InterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorV2V3InterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AggregatorV2V3InterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorV2V3InterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AggregatorV2V3InterfaceSession struct {
	Contract     *AggregatorV2V3Interface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// AggregatorV2V3InterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AggregatorV2V3InterfaceCallerSession struct {
	Contract *AggregatorV2V3InterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// AggregatorV2V3InterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AggregatorV2V3InterfaceTransactorSession struct {
	Contract     *AggregatorV2V3InterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// AggregatorV2V3InterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type AggregatorV2V3InterfaceRaw struct {
	Contract *AggregatorV2V3Interface // Generic contract binding to access the raw methods on
}

// AggregatorV2V3InterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AggregatorV2V3InterfaceCallerRaw struct {
	Contract *AggregatorV2V3InterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// AggregatorV2V3InterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AggregatorV2V3InterfaceTransactorRaw struct {
	Contract *AggregatorV2V3InterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAggregatorV2V3Interface creates a new instance of AggregatorV2V3Interface, bound to a specific deployed contract.
func NewAggregatorV2V3Interface(address common.Address, backend bind.ContractBackend) (*AggregatorV2V3Interface, error) {
	contract, err := bindAggregatorV2V3Interface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3Interface{AggregatorV2V3InterfaceCaller: AggregatorV2V3InterfaceCaller{contract: contract}, AggregatorV2V3InterfaceTransactor: AggregatorV2V3InterfaceTransactor{contract: contract}, AggregatorV2V3InterfaceFilterer: AggregatorV2V3InterfaceFilterer{contract: contract}}, nil
}

// NewAggregatorV2V3InterfaceCaller creates a new read-only instance of AggregatorV2V3Interface, bound to a specific deployed contract.
func NewAggregatorV2V3InterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorV2V3InterfaceCaller, error) {
	contract, err := bindAggregatorV2V3Interface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3InterfaceCaller{contract: contract}, nil
}

// NewAggregatorV2V3InterfaceTransactor creates a new write-only instance of AggregatorV2V3Interface, bound to a specific deployed contract.
func NewAggregatorV2V3InterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorV2V3InterfaceTransactor, error) {
	contract, err := bindAggregatorV2V3Interface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3InterfaceTransactor{contract: contract}, nil
}

// NewAggregatorV2V3InterfaceFilterer creates a new log filterer instance of AggregatorV2V3Interface, bound to a specific deployed contract.
func NewAggregatorV2V3InterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorV2V3InterfaceFilterer, error) {
	contract, err := bindAggregatorV2V3Interface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3InterfaceFilterer{contract: contract}, nil
}

// bindAggregatorV2V3Interface binds a generic wrapper to an already deployed contract.
func bindAggregatorV2V3Interface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AggregatorV2V3InterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorV2V3Interface.Contract.AggregatorV2V3InterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.AggregatorV2V3InterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.AggregatorV2V3InterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorV2V3Interface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.contract.Transact(opts, method, params...)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) Decimals() (uint8, error) {
	return _AggregatorV2V3Interface.Contract.Decimals(&_AggregatorV2V3Interface.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) Decimals() (uint8, error) {
	return _AggregatorV2V3Interface.Contract.Decimals(&_AggregatorV2V3Interface.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) Description() (string, error) {
	return _AggregatorV2V3Interface.Contract.Description(&_AggregatorV2V3Interface.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) Description() (string, error) {
	return _AggregatorV2V3Interface.Contract.Description(&_AggregatorV2V3Interface.CallOpts)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetAnswer(&_AggregatorV2V3Interface.CallOpts, roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetAnswer(&_AggregatorV2V3Interface.CallOpts, roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
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

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV2V3Interface.Contract.GetRoundData(&_AggregatorV2V3Interface.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV2V3Interface.Contract.GetRoundData(&_AggregatorV2V3Interface.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetTimestamp(&_AggregatorV2V3Interface.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetTimestamp(&_AggregatorV2V3Interface.CallOpts, roundId)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestAnswer(&_AggregatorV2V3Interface.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestAnswer(&_AggregatorV2V3Interface.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestRound() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestRound(&_AggregatorV2V3Interface.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestRound() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestRound(&_AggregatorV2V3Interface.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
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

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV2V3Interface.Contract.LatestRoundData(&_AggregatorV2V3Interface.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV2V3Interface.Contract.LatestRoundData(&_AggregatorV2V3Interface.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestTimestamp(&_AggregatorV2V3Interface.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestTimestamp(&_AggregatorV2V3Interface.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) Version() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.Version(&_AggregatorV2V3Interface.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) Version() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.Version(&_AggregatorV2V3Interface.CallOpts)
}

// AggregatorV2V3InterfaceAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the AggregatorV2V3Interface contract.
type AggregatorV2V3InterfaceAnswerUpdatedIterator struct {
	Event *AggregatorV2V3InterfaceAnswerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AggregatorV2V3InterfaceAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorV2V3InterfaceAnswerUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AggregatorV2V3InterfaceAnswerUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AggregatorV2V3InterfaceAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorV2V3InterfaceAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorV2V3InterfaceAnswerUpdated represents a AnswerUpdated event raised by the AggregatorV2V3Interface contract.
type AggregatorV2V3InterfaceAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*AggregatorV2V3InterfaceAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorV2V3Interface.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3InterfaceAnswerUpdatedIterator{contract: _AggregatorV2V3Interface.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *AggregatorV2V3InterfaceAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorV2V3Interface.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AggregatorV2V3InterfaceAnswerUpdated)
				if err := _AggregatorV2V3Interface.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) ParseAnswerUpdated(log types.Log) (*AggregatorV2V3InterfaceAnswerUpdated, error) {
	event := new(AggregatorV2V3InterfaceAnswerUpdated)
	if err := _AggregatorV2V3Interface.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AggregatorV2V3InterfaceNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the AggregatorV2V3Interface contract.
type AggregatorV2V3InterfaceNewRoundIterator struct {
	Event *AggregatorV2V3InterfaceNewRound // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AggregatorV2V3InterfaceNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorV2V3InterfaceNewRound)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AggregatorV2V3InterfaceNewRound)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AggregatorV2V3InterfaceNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorV2V3InterfaceNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorV2V3InterfaceNewRound represents a NewRound event raised by the AggregatorV2V3Interface contract.
type AggregatorV2V3InterfaceNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*AggregatorV2V3InterfaceNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorV2V3Interface.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3InterfaceNewRoundIterator{contract: _AggregatorV2V3Interface.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *AggregatorV2V3InterfaceNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorV2V3Interface.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AggregatorV2V3InterfaceNewRound)
				if err := _AggregatorV2V3Interface.contract.UnpackLog(event, "NewRound", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) ParseNewRound(log types.Log) (*AggregatorV2V3InterfaceNewRound, error) {
	event := new(AggregatorV2V3InterfaceNewRound)
	if err := _AggregatorV2V3Interface.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AggregatorV3InterfaceMetaData contains all meta data concerning the AggregatorV3Interface contract.
var AggregatorV3InterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AggregatorV3InterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use AggregatorV3InterfaceMetaData.ABI instead.
var AggregatorV3InterfaceABI = AggregatorV3InterfaceMetaData.ABI

// AggregatorV3Interface is an auto generated Go binding around an Ethereum contract.
type AggregatorV3Interface struct {
	AggregatorV3InterfaceCaller     // Read-only binding to the contract
	AggregatorV3InterfaceTransactor // Write-only binding to the contract
	AggregatorV3InterfaceFilterer   // Log filterer for contract events
}

// AggregatorV3InterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type AggregatorV3InterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorV3InterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AggregatorV3InterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorV3InterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AggregatorV3InterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorV3InterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AggregatorV3InterfaceSession struct {
	Contract     *AggregatorV3Interface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// AggregatorV3InterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AggregatorV3InterfaceCallerSession struct {
	Contract *AggregatorV3InterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// AggregatorV3InterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AggregatorV3InterfaceTransactorSession struct {
	Contract     *AggregatorV3InterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// AggregatorV3InterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type AggregatorV3InterfaceRaw struct {
	Contract *AggregatorV3Interface // Generic contract binding to access the raw methods on
}

// AggregatorV3InterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AggregatorV3InterfaceCallerRaw struct {
	Contract *AggregatorV3InterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// AggregatorV3InterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AggregatorV3InterfaceTransactorRaw struct {
	Contract *AggregatorV3InterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAggregatorV3Interface creates a new instance of AggregatorV3Interface, bound to a specific deployed contract.
func NewAggregatorV3Interface(address common.Address, backend bind.ContractBackend) (*AggregatorV3Interface, error) {
	contract, err := bindAggregatorV3Interface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorV3Interface{AggregatorV3InterfaceCaller: AggregatorV3InterfaceCaller{contract: contract}, AggregatorV3InterfaceTransactor: AggregatorV3InterfaceTransactor{contract: contract}, AggregatorV3InterfaceFilterer: AggregatorV3InterfaceFilterer{contract: contract}}, nil
}

// NewAggregatorV3InterfaceCaller creates a new read-only instance of AggregatorV3Interface, bound to a specific deployed contract.
func NewAggregatorV3InterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorV3InterfaceCaller, error) {
	contract, err := bindAggregatorV3Interface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorV3InterfaceCaller{contract: contract}, nil
}

// NewAggregatorV3InterfaceTransactor creates a new write-only instance of AggregatorV3Interface, bound to a specific deployed contract.
func NewAggregatorV3InterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorV3InterfaceTransactor, error) {
	contract, err := bindAggregatorV3Interface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorV3InterfaceTransactor{contract: contract}, nil
}

// NewAggregatorV3InterfaceFilterer creates a new log filterer instance of AggregatorV3Interface, bound to a specific deployed contract.
func NewAggregatorV3InterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorV3InterfaceFilterer, error) {
	contract, err := bindAggregatorV3Interface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorV3InterfaceFilterer{contract: contract}, nil
}

// bindAggregatorV3Interface binds a generic wrapper to an already deployed contract.
func bindAggregatorV3Interface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AggregatorV3InterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorV3Interface *AggregatorV3InterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorV3Interface.Contract.AggregatorV3InterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorV3Interface *AggregatorV3InterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.AggregatorV3InterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorV3Interface *AggregatorV3InterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.AggregatorV3InterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorV3Interface *AggregatorV3InterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorV3Interface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorV3Interface *AggregatorV3InterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorV3Interface *AggregatorV3InterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.contract.Transact(opts, method, params...)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AggregatorV3Interface.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorV3Interface *AggregatorV3InterfaceSession) Decimals() (uint8, error) {
	return _AggregatorV3Interface.Contract.Decimals(&_AggregatorV3Interface.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) Decimals() (uint8, error) {
	return _AggregatorV3Interface.Contract.Decimals(&_AggregatorV3Interface.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AggregatorV3Interface.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorV3Interface *AggregatorV3InterfaceSession) Description() (string, error) {
	return _AggregatorV3Interface.Contract.Description(&_AggregatorV3Interface.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) Description() (string, error) {
	return _AggregatorV3Interface.Contract.Description(&_AggregatorV3Interface.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorV3Interface.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
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

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV3Interface *AggregatorV3InterfaceSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV3Interface.Contract.GetRoundData(&_AggregatorV3Interface.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV3Interface.Contract.GetRoundData(&_AggregatorV3Interface.CallOpts, _roundId)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorV3Interface.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
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

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV3Interface *AggregatorV3InterfaceSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV3Interface.Contract.LatestRoundData(&_AggregatorV3Interface.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV3Interface.Contract.LatestRoundData(&_AggregatorV3Interface.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV3Interface.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorV3Interface *AggregatorV3InterfaceSession) Version() (*big.Int, error) {
	return _AggregatorV3Interface.Contract.Version(&_AggregatorV3Interface.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) Version() (*big.Int, error) {
	return _AggregatorV3Interface.Contract.Version(&_AggregatorV3Interface.CallOpts)
}

// AggregatorValidatorInterfaceMetaData contains all meta data concerning the AggregatorValidatorInterface contract.
var AggregatorValidatorInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"previousRoundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"previousAnswer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"currentRoundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"currentAnswer\",\"type\":\"int256\"}],\"name\":\"validate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AggregatorValidatorInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use AggregatorValidatorInterfaceMetaData.ABI instead.
var AggregatorValidatorInterfaceABI = AggregatorValidatorInterfaceMetaData.ABI

// AggregatorValidatorInterface is an auto generated Go binding around an Ethereum contract.
type AggregatorValidatorInterface struct {
	AggregatorValidatorInterfaceCaller     // Read-only binding to the contract
	AggregatorValidatorInterfaceTransactor // Write-only binding to the contract
	AggregatorValidatorInterfaceFilterer   // Log filterer for contract events
}

// AggregatorValidatorInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorValidatorInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorValidatorInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AggregatorValidatorInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorValidatorInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AggregatorValidatorInterfaceSession struct {
	Contract     *AggregatorValidatorInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// AggregatorValidatorInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AggregatorValidatorInterfaceCallerSession struct {
	Contract *AggregatorValidatorInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// AggregatorValidatorInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AggregatorValidatorInterfaceTransactorSession struct {
	Contract     *AggregatorValidatorInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// AggregatorValidatorInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceRaw struct {
	Contract *AggregatorValidatorInterface // Generic contract binding to access the raw methods on
}

// AggregatorValidatorInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceCallerRaw struct {
	Contract *AggregatorValidatorInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// AggregatorValidatorInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceTransactorRaw struct {
	Contract *AggregatorValidatorInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAggregatorValidatorInterface creates a new instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterface(address common.Address, backend bind.ContractBackend) (*AggregatorValidatorInterface, error) {
	contract, err := bindAggregatorValidatorInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterface{AggregatorValidatorInterfaceCaller: AggregatorValidatorInterfaceCaller{contract: contract}, AggregatorValidatorInterfaceTransactor: AggregatorValidatorInterfaceTransactor{contract: contract}, AggregatorValidatorInterfaceFilterer: AggregatorValidatorInterfaceFilterer{contract: contract}}, nil
}

// NewAggregatorValidatorInterfaceCaller creates a new read-only instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorValidatorInterfaceCaller, error) {
	contract, err := bindAggregatorValidatorInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceCaller{contract: contract}, nil
}

// NewAggregatorValidatorInterfaceTransactor creates a new write-only instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorValidatorInterfaceTransactor, error) {
	contract, err := bindAggregatorValidatorInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceTransactor{contract: contract}, nil
}

// NewAggregatorValidatorInterfaceFilterer creates a new log filterer instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorValidatorInterfaceFilterer, error) {
	contract, err := bindAggregatorValidatorInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceFilterer{contract: contract}, nil
}

// bindAggregatorValidatorInterface binds a generic wrapper to an already deployed contract.
func bindAggregatorValidatorInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AggregatorValidatorInterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorValidatorInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.contract.Transact(opts, method, params...)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactor) Validate(opts *bind.TransactOpts, previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.contract.Transact(opts, "validate", previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceSession) Validate(previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.Validate(&_AggregatorValidatorInterface.TransactOpts, previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorSession) Validate(previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.Validate(&_AggregatorValidatorInterface.TransactOpts, previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// ConfirmedOwnerMetaData contains all meta data concerning the ConfirmedOwner contract.
var ConfirmedOwnerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161051438038061051483398101604081905261002f9161016e565b8060006001600160a01b03821661008d5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100bd576100bd816100c5565b50505061019e565b336001600160a01b0382160361011d5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610084565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561018057600080fd5b81516001600160a01b038116811461019757600080fd5b9392505050565b610367806101ad6000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806379ba5097146100465780638da5cb5b14610050578063f2fde38b1461007c575b600080fd5b61004e61008f565b005b6000546040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61004e61008a36600461031d565b610191565b60015473ffffffffffffffffffffffffffffffffffffffff163314610115576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6101996101a5565b6101a281610228565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610226576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161010c565b565b3373ffffffffffffffffffffffffffffffffffffffff8216036102a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161010c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561032f57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461035357600080fd5b939250505056fea164736f6c6343000813000a",
}

// ConfirmedOwnerABI is the input ABI used to generate the binding from.
// Deprecated: Use ConfirmedOwnerMetaData.ABI instead.
var ConfirmedOwnerABI = ConfirmedOwnerMetaData.ABI

// ConfirmedOwnerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConfirmedOwnerMetaData.Bin instead.
var ConfirmedOwnerBin = ConfirmedOwnerMetaData.Bin

// DeployConfirmedOwner deploys a new Ethereum contract, binding an instance of ConfirmedOwner to it.
func DeployConfirmedOwner(auth *bind.TransactOpts, backend bind.ContractBackend, newOwner common.Address) (common.Address, *types.Transaction, *ConfirmedOwner, error) {
	parsed, err := ConfirmedOwnerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConfirmedOwnerBin), backend, newOwner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConfirmedOwner{ConfirmedOwnerCaller: ConfirmedOwnerCaller{contract: contract}, ConfirmedOwnerTransactor: ConfirmedOwnerTransactor{contract: contract}, ConfirmedOwnerFilterer: ConfirmedOwnerFilterer{contract: contract}}, nil
}

// ConfirmedOwner is an auto generated Go binding around an Ethereum contract.
type ConfirmedOwner struct {
	ConfirmedOwnerCaller     // Read-only binding to the contract
	ConfirmedOwnerTransactor // Write-only binding to the contract
	ConfirmedOwnerFilterer   // Log filterer for contract events
}

// ConfirmedOwnerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConfirmedOwnerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConfirmedOwnerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConfirmedOwnerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConfirmedOwnerSession struct {
	Contract     *ConfirmedOwner   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConfirmedOwnerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConfirmedOwnerCallerSession struct {
	Contract *ConfirmedOwnerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ConfirmedOwnerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConfirmedOwnerTransactorSession struct {
	Contract     *ConfirmedOwnerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ConfirmedOwnerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConfirmedOwnerRaw struct {
	Contract *ConfirmedOwner // Generic contract binding to access the raw methods on
}

// ConfirmedOwnerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConfirmedOwnerCallerRaw struct {
	Contract *ConfirmedOwnerCaller // Generic read-only contract binding to access the raw methods on
}

// ConfirmedOwnerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConfirmedOwnerTransactorRaw struct {
	Contract *ConfirmedOwnerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConfirmedOwner creates a new instance of ConfirmedOwner, bound to a specific deployed contract.
func NewConfirmedOwner(address common.Address, backend bind.ContractBackend) (*ConfirmedOwner, error) {
	contract, err := bindConfirmedOwner(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwner{ConfirmedOwnerCaller: ConfirmedOwnerCaller{contract: contract}, ConfirmedOwnerTransactor: ConfirmedOwnerTransactor{contract: contract}, ConfirmedOwnerFilterer: ConfirmedOwnerFilterer{contract: contract}}, nil
}

// NewConfirmedOwnerCaller creates a new read-only instance of ConfirmedOwner, bound to a specific deployed contract.
func NewConfirmedOwnerCaller(address common.Address, caller bind.ContractCaller) (*ConfirmedOwnerCaller, error) {
	contract, err := bindConfirmedOwner(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerCaller{contract: contract}, nil
}

// NewConfirmedOwnerTransactor creates a new write-only instance of ConfirmedOwner, bound to a specific deployed contract.
func NewConfirmedOwnerTransactor(address common.Address, transactor bind.ContractTransactor) (*ConfirmedOwnerTransactor, error) {
	contract, err := bindConfirmedOwner(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerTransactor{contract: contract}, nil
}

// NewConfirmedOwnerFilterer creates a new log filterer instance of ConfirmedOwner, bound to a specific deployed contract.
func NewConfirmedOwnerFilterer(address common.Address, filterer bind.ContractFilterer) (*ConfirmedOwnerFilterer, error) {
	contract, err := bindConfirmedOwner(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerFilterer{contract: contract}, nil
}

// bindConfirmedOwner binds a generic wrapper to an already deployed contract.
func bindConfirmedOwner(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConfirmedOwnerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfirmedOwner *ConfirmedOwnerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfirmedOwner.Contract.ConfirmedOwnerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfirmedOwner *ConfirmedOwnerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.ConfirmedOwnerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfirmedOwner *ConfirmedOwnerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.ConfirmedOwnerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfirmedOwner *ConfirmedOwnerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfirmedOwner.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfirmedOwner *ConfirmedOwnerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfirmedOwner *ConfirmedOwnerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwner *ConfirmedOwnerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConfirmedOwner.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwner *ConfirmedOwnerSession) Owner() (common.Address, error) {
	return _ConfirmedOwner.Contract.Owner(&_ConfirmedOwner.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwner *ConfirmedOwnerCallerSession) Owner() (common.Address, error) {
	return _ConfirmedOwner.Contract.Owner(&_ConfirmedOwner.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwner *ConfirmedOwnerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwner.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwner *ConfirmedOwnerSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.AcceptOwnership(&_ConfirmedOwner.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwner *ConfirmedOwnerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.AcceptOwnership(&_ConfirmedOwner.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwner *ConfirmedOwnerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwner.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwner *ConfirmedOwnerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.TransferOwnership(&_ConfirmedOwner.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwner *ConfirmedOwnerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.TransferOwnership(&_ConfirmedOwner.TransactOpts, to)
}

// ConfirmedOwnerOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the ConfirmedOwner contract.
type ConfirmedOwnerOwnershipTransferRequestedIterator struct {
	Event *ConfirmedOwnerOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConfirmedOwnerOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfirmedOwnerOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConfirmedOwnerOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConfirmedOwnerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfirmedOwnerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfirmedOwnerOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the ConfirmedOwner contract.
type ConfirmedOwnerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfirmedOwnerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwner.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerOwnershipTransferRequestedIterator{contract: _ConfirmedOwner.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ConfirmedOwnerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwner.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfirmedOwnerOwnershipTransferRequested)
				if err := _ConfirmedOwner.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) ParseOwnershipTransferRequested(log types.Log) (*ConfirmedOwnerOwnershipTransferRequested, error) {
	event := new(ConfirmedOwnerOwnershipTransferRequested)
	if err := _ConfirmedOwner.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfirmedOwnerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ConfirmedOwner contract.
type ConfirmedOwnerOwnershipTransferredIterator struct {
	Event *ConfirmedOwnerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConfirmedOwnerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfirmedOwnerOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConfirmedOwnerOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConfirmedOwnerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfirmedOwnerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfirmedOwnerOwnershipTransferred represents a OwnershipTransferred event raised by the ConfirmedOwner contract.
type ConfirmedOwnerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfirmedOwnerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwner.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerOwnershipTransferredIterator{contract: _ConfirmedOwner.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConfirmedOwnerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwner.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfirmedOwnerOwnershipTransferred)
				if err := _ConfirmedOwner.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) ParseOwnershipTransferred(log types.Log) (*ConfirmedOwnerOwnershipTransferred, error) {
	event := new(ConfirmedOwnerOwnershipTransferred)
	if err := _ConfirmedOwner.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfirmedOwnerWithProposalMetaData contains all meta data concerning the ConfirmedOwnerWithProposal contract.
var ConfirmedOwnerWithProposalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pendingOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161052f38038061052f83398101604081905261002f91610186565b6001600160a01b03821661008a5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100ba576100ba816100c1565b50506101b9565b336001600160a01b038216036101195760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610081565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b038116811461018157600080fd5b919050565b6000806040838503121561019957600080fd5b6101a28361016a565b91506101b06020840161016a565b90509250929050565b610367806101c86000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806379ba5097146100465780638da5cb5b14610050578063f2fde38b1461007c575b600080fd5b61004e61008f565b005b6000546040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61004e61008a36600461031d565b610191565b60015473ffffffffffffffffffffffffffffffffffffffff163314610115576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6101996101a5565b6101a281610228565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610226576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161010c565b565b3373ffffffffffffffffffffffffffffffffffffffff8216036102a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161010c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561032f57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461035357600080fd5b939250505056fea164736f6c6343000813000a",
}

// ConfirmedOwnerWithProposalABI is the input ABI used to generate the binding from.
// Deprecated: Use ConfirmedOwnerWithProposalMetaData.ABI instead.
var ConfirmedOwnerWithProposalABI = ConfirmedOwnerWithProposalMetaData.ABI

// ConfirmedOwnerWithProposalBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConfirmedOwnerWithProposalMetaData.Bin instead.
var ConfirmedOwnerWithProposalBin = ConfirmedOwnerWithProposalMetaData.Bin

// DeployConfirmedOwnerWithProposal deploys a new Ethereum contract, binding an instance of ConfirmedOwnerWithProposal to it.
func DeployConfirmedOwnerWithProposal(auth *bind.TransactOpts, backend bind.ContractBackend, newOwner common.Address, pendingOwner common.Address) (common.Address, *types.Transaction, *ConfirmedOwnerWithProposal, error) {
	parsed, err := ConfirmedOwnerWithProposalMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConfirmedOwnerWithProposalBin), backend, newOwner, pendingOwner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConfirmedOwnerWithProposal{ConfirmedOwnerWithProposalCaller: ConfirmedOwnerWithProposalCaller{contract: contract}, ConfirmedOwnerWithProposalTransactor: ConfirmedOwnerWithProposalTransactor{contract: contract}, ConfirmedOwnerWithProposalFilterer: ConfirmedOwnerWithProposalFilterer{contract: contract}}, nil
}

// ConfirmedOwnerWithProposal is an auto generated Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposal struct {
	ConfirmedOwnerWithProposalCaller     // Read-only binding to the contract
	ConfirmedOwnerWithProposalTransactor // Write-only binding to the contract
	ConfirmedOwnerWithProposalFilterer   // Log filterer for contract events
}

// ConfirmedOwnerWithProposalCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerWithProposalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerWithProposalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConfirmedOwnerWithProposalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerWithProposalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConfirmedOwnerWithProposalSession struct {
	Contract     *ConfirmedOwnerWithProposal // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ConfirmedOwnerWithProposalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConfirmedOwnerWithProposalCallerSession struct {
	Contract *ConfirmedOwnerWithProposalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// ConfirmedOwnerWithProposalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConfirmedOwnerWithProposalTransactorSession struct {
	Contract     *ConfirmedOwnerWithProposalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// ConfirmedOwnerWithProposalRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalRaw struct {
	Contract *ConfirmedOwnerWithProposal // Generic contract binding to access the raw methods on
}

// ConfirmedOwnerWithProposalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalCallerRaw struct {
	Contract *ConfirmedOwnerWithProposalCaller // Generic read-only contract binding to access the raw methods on
}

// ConfirmedOwnerWithProposalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalTransactorRaw struct {
	Contract *ConfirmedOwnerWithProposalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConfirmedOwnerWithProposal creates a new instance of ConfirmedOwnerWithProposal, bound to a specific deployed contract.
func NewConfirmedOwnerWithProposal(address common.Address, backend bind.ContractBackend) (*ConfirmedOwnerWithProposal, error) {
	contract, err := bindConfirmedOwnerWithProposal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposal{ConfirmedOwnerWithProposalCaller: ConfirmedOwnerWithProposalCaller{contract: contract}, ConfirmedOwnerWithProposalTransactor: ConfirmedOwnerWithProposalTransactor{contract: contract}, ConfirmedOwnerWithProposalFilterer: ConfirmedOwnerWithProposalFilterer{contract: contract}}, nil
}

// NewConfirmedOwnerWithProposalCaller creates a new read-only instance of ConfirmedOwnerWithProposal, bound to a specific deployed contract.
func NewConfirmedOwnerWithProposalCaller(address common.Address, caller bind.ContractCaller) (*ConfirmedOwnerWithProposalCaller, error) {
	contract, err := bindConfirmedOwnerWithProposal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalCaller{contract: contract}, nil
}

// NewConfirmedOwnerWithProposalTransactor creates a new write-only instance of ConfirmedOwnerWithProposal, bound to a specific deployed contract.
func NewConfirmedOwnerWithProposalTransactor(address common.Address, transactor bind.ContractTransactor) (*ConfirmedOwnerWithProposalTransactor, error) {
	contract, err := bindConfirmedOwnerWithProposal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalTransactor{contract: contract}, nil
}

// NewConfirmedOwnerWithProposalFilterer creates a new log filterer instance of ConfirmedOwnerWithProposal, bound to a specific deployed contract.
func NewConfirmedOwnerWithProposalFilterer(address common.Address, filterer bind.ContractFilterer) (*ConfirmedOwnerWithProposalFilterer, error) {
	contract, err := bindConfirmedOwnerWithProposal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalFilterer{contract: contract}, nil
}

// bindConfirmedOwnerWithProposal binds a generic wrapper to an already deployed contract.
func bindConfirmedOwnerWithProposal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConfirmedOwnerWithProposalMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfirmedOwnerWithProposal.Contract.ConfirmedOwnerWithProposalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.ConfirmedOwnerWithProposalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.ConfirmedOwnerWithProposalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfirmedOwnerWithProposal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConfirmedOwnerWithProposal.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalSession) Owner() (common.Address, error) {
	return _ConfirmedOwnerWithProposal.Contract.Owner(&_ConfirmedOwnerWithProposal.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalCallerSession) Owner() (common.Address, error) {
	return _ConfirmedOwnerWithProposal.Contract.Owner(&_ConfirmedOwnerWithProposal.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.AcceptOwnership(&_ConfirmedOwnerWithProposal.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.AcceptOwnership(&_ConfirmedOwnerWithProposal.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.TransferOwnership(&_ConfirmedOwnerWithProposal.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.TransferOwnership(&_ConfirmedOwnerWithProposal.TransactOpts, to)
}

// ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the ConfirmedOwnerWithProposal contract.
type ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator struct {
	Event *ConfirmedOwnerWithProposalOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfirmedOwnerWithProposalOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConfirmedOwnerWithProposalOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfirmedOwnerWithProposalOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the ConfirmedOwnerWithProposal contract.
type ConfirmedOwnerWithProposalOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwnerWithProposal.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator{contract: _ConfirmedOwnerWithProposal.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ConfirmedOwnerWithProposalOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwnerWithProposal.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfirmedOwnerWithProposalOwnershipTransferRequested)
				if err := _ConfirmedOwnerWithProposal.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) ParseOwnershipTransferRequested(log types.Log) (*ConfirmedOwnerWithProposalOwnershipTransferRequested, error) {
	event := new(ConfirmedOwnerWithProposalOwnershipTransferRequested)
	if err := _ConfirmedOwnerWithProposal.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfirmedOwnerWithProposalOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ConfirmedOwnerWithProposal contract.
type ConfirmedOwnerWithProposalOwnershipTransferredIterator struct {
	Event *ConfirmedOwnerWithProposalOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConfirmedOwnerWithProposalOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfirmedOwnerWithProposalOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConfirmedOwnerWithProposalOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConfirmedOwnerWithProposalOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfirmedOwnerWithProposalOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfirmedOwnerWithProposalOwnershipTransferred represents a OwnershipTransferred event raised by the ConfirmedOwnerWithProposal contract.
type ConfirmedOwnerWithProposalOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfirmedOwnerWithProposalOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwnerWithProposal.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalOwnershipTransferredIterator{contract: _ConfirmedOwnerWithProposal.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConfirmedOwnerWithProposalOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwnerWithProposal.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfirmedOwnerWithProposalOwnershipTransferred)
				if err := _ConfirmedOwnerWithProposal.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) ParseOwnershipTransferred(log types.Log) (*ConfirmedOwnerWithProposalOwnershipTransferred, error) {
	event := new(ConfirmedOwnerWithProposalOwnershipTransferred)
	if err := _ConfirmedOwnerWithProposal.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorMetaData contains all meta data concerning the ExposedOCR2Aggregator contract.
var ExposedOCR2AggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"oldLinkToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"newLinkToken\",\"type\":\"address\"}],\"name\":\"LinkTokenSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"observationsTimestamp\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"RequesterAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"}],\"name\":\"RoundRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"previousValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousGasLimit\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"currentValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"currentGasLimit\",\"type\":\"uint32\"}],\"name\":\"ValidatorConfigSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encodedConfig\",\"type\":\"bytes\"}],\"name\":\"exposedConfigDigestFromConfigData\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBillingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequesterAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId_\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidatorConfig\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer_\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp_\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestNewRound\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"setLinkToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"requesterAccessController\",\"type\":\"address\"}],\"name\":\"setRequesterAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"newValidator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"newGasLimit\",\"type\":\"uint32\"}],\"name\":\"setValidatorConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b506000806000806000806040518060200160405280600081525033806000806001600160a01b0316826001600160a01b031603620000965760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c957620000c98162000168565b5050601180546001600160a01b0319166001600160a01b038a169081179091556040519091506000907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a908290a3620001228462000213565b60ff821660c0526010620001378282620004f0565b5062000143836200028c565b6200015060008062000307565b50505050601791820b608052900b60a05250620005bc565b336001600160a01b03821603620001c25760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6012546001600160a01b0390811690821681146200028857601280546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d4891291015b60405180910390a15b5050565b62000296620003ee565b600f546001600160a01b0390811690821681146200028857600f80546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae63491016200027f565b62000311620003ee565b60408051808201909152600e546001600160a01b03808216808452600160a01b90920463ffffffff16602084015284161415806200035f57508163ffffffff16816020015163ffffffff1614155b15620003e9576040805180820182526001600160a01b0385811680835263ffffffff8681166020948501819052600e80546001600160c01b0319168417600160a01b830217905586518786015187519316835294820152909392909116917fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541910160405180910390a35b505050565b6000546001600160a01b031633146200044a5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016200008d565b565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200047757607f821691505b6020821081036200049857634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620003e957600081815260208120601f850160051c81016020861015620004c75750805b601f850160051c820191505b81811015620004e857828155600101620004d3565b505050505050565b81516001600160401b038111156200050c576200050c6200044c565b62000524816200051d845462000462565b846200049e565b602080601f8311600181146200055c5760008415620005435750858301515b600019600386901b1c1916600185901b178555620004e8565b600085815260208120601f198616915b828110156200058d578886015182559484019460019091019084016200056c565b5085821015620005ac5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805160a05160c05161506162000608600039600061045401526000818161051401528181612264015261366e0152600081816103500152818161223c015261364101526150616000f3fe608060405234801561001057600080fd5b50600436106102de5760003560e01c80639bd2c0b111610186578063d09dc339116100e3578063e76d516811610097578063f2fde38b11610071578063f2fde38b14610866578063fbffd2c114610879578063feaf968c1461088c57600080fd5b8063e76d51681461082f578063eb45716314610840578063eb5dcd6c1461085357600080fd5b8063e3d0e712116100c8578063e3d0e712146107aa578063e4902f82146107bd578063e5fe4577146107e557600080fd5b8063d09dc33914610791578063daffc4b51461079957600080fd5b8063b1dc65a41161013a578063b633620c1161011f578063b633620c1461075a578063c10753291461076d578063c4c92b371461078057600080fd5b8063b1dc65a414610734578063b5ab58dc1461074757600080fd5b80639e3ceeab1161016b5780639e3ceeab146106dd578063afcb95d7146106f0578063b121e1471461072157600080fd5b80639bd2c0b1146106785780639c849b30146106ca57600080fd5b8063668a0f021161023f57806381ff7048116101f35780638da5cb5b116101cd5780638da5cb5b146105e657806398e5b12a1461060b5780639a6fc8f51461062e57600080fd5b806381ff7048146105595780638205bf6a146105895780638ac28d5a146105d357600080fd5b80637284e416116102245780637284e41614610536578063739376c11461053e57806379ba50971461055157600080fd5b8063668a0f02146104f757806370da2f671461050f57600080fd5b80634fb174701161029657806354fd4d501161027b57806354fd4d50146104c7578063643dc105146104cf578063666cab8d146104e257600080fd5b80634fb174701461048857806350d25bcd1461049d57600080fd5b806322adbc78116102c757806322adbc781461034b5780632993726814610385578063313ce5671461044f57600080fd5b80630eafb25b146102e3578063181f5a7714610309575b600080fd5b6102f66102f13660046143b7565b61091b565b6040519081526020015b60405180910390f35b60408051808201909152601481527f4f43523241676772656761746f7220312e302e3000000000000000000000000060208201525b604051610300919061441a565b6103727f000000000000000000000000000000000000000000000000000000000000000081565b60405160179190910b8152602001610300565b610413600b546a0100000000000000000000810463ffffffff908116926e010000000000000000000000000000830482169272010000000000000000000000000000000000008104831692760100000000000000000000000000000000000000000000820416917a01000000000000000000000000000000000000000000000000000090910462ffffff1690565b6040805163ffffffff9687168152948616602086015292851692840192909252909216606082015262ffffff909116608082015260a001610300565b6104767f000000000000000000000000000000000000000000000000000000000000000081565b60405160ff9091168152602001610300565b61049b61049636600461442d565b610a3c565b005b600b546601000000000000900463ffffffff166000908152600c602052604090205460170b6102f6565b6102f6600681565b61049b6104dd366004614478565b610cbc565b6104ea610f77565b6040516103009190614535565b600b546601000000000000900463ffffffff166102f6565b6103727f000000000000000000000000000000000000000000000000000000000000000081565b61033e610fd9565b6102f661054c366004614722565b611062565b61049b6110bc565b600d54600a546040805163ffffffff80851682526401000000009094049093166020840152820152606001610300565b600b546601000000000000900463ffffffff9081166000908152600c60205260409020547c01000000000000000000000000000000000000000000000000000000009004166102f6565b61049b6105e13660046143b7565b611185565b6000546001600160a01b03165b6040516001600160a01b039091168152602001610300565b6106136111fa565b60405169ffffffffffffffffffff9091168152602001610300565b61064161063c366004614825565b611371565b6040805169ffffffffffffffffffff968716815260208101959095528401929092526060830152909116608082015260a001610300565b604080518082018252600e546001600160a01b0381168083527401000000000000000000000000000000000000000090910463ffffffff16602092830181905283519182529181019190915201610300565b61049b6106d8366004614896565b61142f565b61049b6106eb3660046143b7565b611625565b600a54600b546040805160008152602081019390935261010090910460081c63ffffffff1690820152606001610300565b61049b61072f3660046143b7565b6116bc565b61049b610742366004614902565b6117b0565b6102f66107553660046149b9565b611cfb565b6102f66107683660046149b9565b611d2d565b61049b61077b3660046149d2565b611d7f565b6012546001600160a01b03166105f3565b6102f6612053565b600f546001600160a01b03166105f3565b61049b6107b83660046149fe565b6120fc565b6107d06107cb3660046143b7565b612978565b60405163ffffffff9091168152602001610300565b6107ed612a36565b6040805195865263ffffffff909416602086015260ff9092169284019290925260179190910b606083015267ffffffffffffffff16608082015260a001610300565b6011546001600160a01b03166105f3565b61049b61084e366004614acb565b612afb565b61049b61086136600461442d565b612c18565b61049b6108743660046143b7565b612d68565b61049b6108873660046143b7565b612d79565b610641600b546601000000000000900463ffffffff9081166000818152600c602090815260409182902082516060810184529054601781900b8083527801000000000000000000000000000000000000000000000000820487169383018490527c0100000000000000000000000000000000000000000000000000000000909104909516920182905291938490565b6001600160a01b03811660009081526002602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046bffffffffffffffffffffffff1691810191909152906109825750600092915050565b600b5460208201516000917201000000000000000000000000000000000000900463ffffffff169060069060ff16601f81106109c0576109c0614af9565b600881049190910154600b546109f6926007166004026101000a90910463ffffffff908116916601000000000000900416614b57565b63ffffffff16610a069190614b7b565b610a1490633b9aca00614b7b565b905081604001516bffffffffffffffffffffffff1681610a349190614b92565b949350505050565b610a44612d8a565b6011546001600160a01b03908116908316819003610a6157505050565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b038416906370a0823190602401602060405180830381865afa158015610abe573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ae29190614ba5565b50610aeb612de6565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000906001600160a01b038316906370a0823190602401602060405180830381865afa158015610b4b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b6f9190614ba5565b6040517fa9059cbb0000000000000000000000000000000000000000000000000000000081526001600160a01b038581166004830152602482018390529192509083169063a9059cbb906044016020604051808303816000875af1158015610bdb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bff9190614bbe565b610c505760405162461bcd60e51b815260206004820152601f60248201527f7472616e736665722072656d61696e696e672066756e6473206661696c65640060448201526064015b60405180910390fd5b601180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0386811691821790925560405190918416907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a90600090a350505b5050565b6012546001600160a01b0316610cda6000546001600160a01b031690565b6001600160a01b0316336001600160a01b03161480610d7f57506040517f6b14daf80000000000000000000000000000000000000000000000000000000081526001600160a01b03821690636b14daf890610d3e9033906000903690600401614be0565b602060405180830381865afa158015610d5b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d7f9190614bbe565b610dcb5760405162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c6044820152606401610c47565b610dd3612de6565b600b80547fffffffffffffffffffffffffffff0000000000000000ffffffffffffffffffff166a010000000000000000000063ffffffff8981169182027fffffffffffffffffffffffffffff00000000ffffffffffffffffffffffffffff16929092176e010000000000000000000000000000898416908102919091177fffffffffffff0000000000000000ffffffffffffffffffffffffffffffffffff1672010000000000000000000000000000000000008985169081027fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff1691909117760100000000000000000000000000000000000000000000948916948502177fffffff000000ffffffffffffffffffffffffffffffffffffffffffffffffffff167a01000000000000000000000000000000000000000000000000000062ffffff89169081029190911790955560408051938452602084019290925290820152606081019190915260808101919091527f0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f9060a00160405180910390a1505050505050565b60606005805480602002602001604051908101604052809291908181526020018280548015610fcf57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610fb1575b5050505050905090565b606060108054610fe890614c1f565b80601f016020809104026020016040519081016040528092919081815260200182805461101490614c1f565b8015610fcf5780601f1061103657610100808354040283529160200191610fcf565b820191906000526020600020905b81548152906001019060200180831161104457509395945050505050565b60006110ad8b8b8b8b8b8b8b8b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508d92508c91506131ca9050565b9b9a5050505050505050505050565b6001546001600160a01b031633146111165760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610c47565b60008054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6001600160a01b038181166000908152601360205260409020541633146111ee5760405162461bcd60e51b815260206004820152601760248201527f4f6e6c792070617965652063616e2077697468647261770000000000000000006044820152606401610c47565b6111f781613258565b50565b600080546001600160a01b031633148061129e5750600f546040517f6b14daf80000000000000000000000000000000000000000000000000000000081526001600160a01b0390911690636b14daf89061125d9033906000903690600401614be0565b602060405180830381865afa15801561127a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061129e9190614bbe565b6112ea5760405162461bcd60e51b815260206004820152601d60248201527f4f6e6c79206f776e6572267265717565737465722063616e2063616c6c0000006044820152606401610c47565b600b54600a546040805191825263ffffffff6101008404600881901c8216602085015260ff811684840152915164ffffffffff9092169366010000000000009004169133917f41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c9181900360600190a2611364816001614c72565b63ffffffff169250505090565b60008080808063ffffffff69ffffffffffffffffffff871611156113a357506000935083925082915081905080611426565b5050505063ffffffff8281166000908152600c602090815260409182902082516060810184529054601781900b8083527801000000000000000000000000000000000000000000000000820486169383018490527c0100000000000000000000000000000000000000000000000000000000909104909416920182905284935090835b91939590929450565b611437612d8a565b8281146114865760405162461bcd60e51b815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a656044820152606401610c47565b60005b8381101561161e5760008585838181106114a5576114a5614af9565b90506020020160208101906114ba91906143b7565b905060008484848181106114d0576114d0614af9565b90506020020160208101906114e591906143b7565b6001600160a01b0380841660009081526013602052604090205491925016801580806115225750826001600160a01b0316826001600160a01b0316145b61156e5760405162461bcd60e51b815260206004820152601160248201527f706179656520616c7265616479207365740000000000000000000000000000006044820152606401610c47565b6001600160a01b03848116600090815260136020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000168583169081179091559083161461160757826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b50505050808061161690614c8f565b915050611489565b5050505050565b61162d612d8a565b600f546001600160a01b039081169082168114610cb857600f80547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0384811691821790925560408051928416835260208301919091527f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae63491015b60405180910390a15050565b6001600160a01b038181166000908152601460205260409020541633146117255760405162461bcd60e51b815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e20616363657074006044820152606401610c47565b6001600160a01b0381811660008181526013602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556014909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60005a604080516101008082018352600b5460ff8116835290810464ffffffffff90811660208085018290526601000000000000840463ffffffff908116968601969096526a01000000000000000000008404861660608601526e01000000000000000000000000000084048616608086015272010000000000000000000000000000000000008404861660a0860152760100000000000000000000000000000000000000000000840490951660c08501527a01000000000000000000000000000000000000000000000000000090920462ffffff1660e08401529394509092918c0135918216116118e45760405162461bcd60e51b815260206004820152600c60248201527f7374616c65207265706f727400000000000000000000000000000000000000006044820152606401610c47565b3360009081526002602052604090205460ff166119435760405162461bcd60e51b815260206004820152601860248201527f756e617574686f72697a6564207472616e736d697474657200000000000000006044820152606401610c47565b600a548b35146119955760405162461bcd60e51b815260206004820152601560248201527f636f6e666967446967657374206d69736d6174636800000000000000000000006044820152606401610c47565b6119a38a8a8a8a8a8a6134a4565b81516119b0906001614cc7565b60ff168714611a015760405162461bcd60e51b815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e6174757265730000000000006044820152606401610c47565b868514611a505760405162461bcd60e51b815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e00006044820152606401610c47565b60008a8a604051611a62929190614ce0565b604051908190038120611a79918e90602001614cf0565b60408051601f19818403018152828252805160209182012083830190925260008084529083018190529092509060005b8a811015611c1f5760006001858a8460208110611ac857611ac8614af9565b611ad591901a601b614cc7565b8f8f86818110611ae757611ae7614af9565b905060200201358e8e87818110611b0057611b00614af9565b9050602002013560405160008152602001604052604051611b3d949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015611b5f573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526003602090815290849020838501909452925460ff8082161515808552610100909204169383019390935290955092509050611bf85760405162461bcd60e51b815260206004820152600f60248201527f7369676e6174757265206572726f7200000000000000000000000000000000006044820152606401610c47565b826020015160080260ff166001901b84019350508080611c1790614c8f565b915050611aa9565b5081827e010101010101010101010101010101010101010101010101010101010101011614611c905760405162461bcd60e51b815260206004820152601060248201527f6475706c6963617465207369676e6572000000000000000000000000000000006044820152606401610c47565b5060009150611cdf9050838d836020020135848e8e8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061354192505050565b9050611ced83828633613a3e565b505050505050505050505050565b600063ffffffff821115611d1157506000919050565b5063ffffffff166000908152600c602052604090205460170b90565b600063ffffffff821115611d4357506000919050565b5063ffffffff9081166000908152600c60205260409020547c010000000000000000000000000000000000000000000000000000000090041690565b6000546001600160a01b0316331480611e2257506012546040517f6b14daf80000000000000000000000000000000000000000000000000000000081526001600160a01b0390911690636b14daf890611de19033906000903690600401614be0565b602060405180830381865afa158015611dfe573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e229190614bbe565b611e6e5760405162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c6044820152606401610c47565b6000611e78613b89565b6011546040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201529192506000916001600160a01b03909116906370a0823190602401602060405180830381865afa158015611edf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f039190614ba5565b905081811015611f555760405162461bcd60e51b815260206004820152601460248201527f696e73756666696369656e742062616c616e63650000000000000000000000006044820152606401610c47565b6011546001600160a01b031663a9059cbb85611f7a611f748686614d06565b87613d6a565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af1158015611fdd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906120019190614bbe565b61204d5760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610c47565b50505050565b6011546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009182916001600160a01b03909116906370a0823190602401602060405180830381865afa1580156120b9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906120dd9190614ba5565b905060006120e9613b89565b90506120f58183614d19565b9250505090565b612104612d8a565b601f865111156121565760405162461bcd60e51b815260206004820152601060248201527f746f6f206d616e79206f7261636c6573000000000000000000000000000000006044820152606401610c47565b84518651146121a75760405162461bcd60e51b815260206004820152601660248201527f6f7261636c65206c656e677468206d69736d61746368000000000000000000006044820152606401610c47565b85516121b4856003614d39565b60ff16106122045760405162461bcd60e51b815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610c47565b6122108460ff16613d84565b604080517f010000000000000000000000000000000000000000000000000000000000000060208201527f0000000000000000000000000000000000000000000000000000000000000000821b60218201527f000000000000000000000000000000000000000000000000000000000000000090911b6039820152605101604051602081830303815290604052805190602001208380519060200120146122f95760405162461bcd60e51b815260206004820152601560248201527f696e76616c6964206f6e636861696e436f6e66696700000000000000000000006044820152606401610c47565b6040805160c0810182528781526020810187905260ff8616918101919091526060810184905267ffffffffffffffff8316608082015260a08101829052600b80547fffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000ff169055612366612de6565b60045460005b818110156124455760006004828154811061238957612389614af9565b6000918252602082200154600580546001600160a01b03909216935090849081106123b6576123b6614af9565b60009182526020808320909101546001600160a01b03948516835260038252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016905594168252600290529190912080547fffffffffffffffffffffffffffffffffffff0000000000000000000000000000169055508061243d81614c8f565b91505061236c565b506124526004600061424f565b61245e6005600061424f565b60005b82515181101561276b57600360008460000151838151811061248557612485614af9565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff16156124f95760405162461bcd60e51b815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610c47565b604080518082019091526001815260ff82166020820152835180516003916000918590811061252a5761252a614af9565b6020908102919091018101516001600160a01b0316825281810192909252604001600090812083518154948401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169015157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff161761010060ff909516949094029390931790925584015180516002929190849081106125cf576125cf614af9565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff16156126435760405162461bcd60e51b815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610c47565b60405180606001604052806001151581526020018260ff16815260200160006bffffffffffffffffffffffff16815250600260008560200151848151811061268d5761268d614af9565b6020908102919091018101516001600160a01b03168252818101929092526040908101600020835181549385015194909201516bffffffffffffffffffffffff1662010000027fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff931515939093167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009094169390931791909117929092161790558061276381614c8f565b915050612461565b50815180516127829160049160209091019061426d565b50602080830151805161279992600592019061426d565b506040820151600b80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600d80547fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff811664010000000063ffffffff43811682029283178555908304811693600193909260009261282b928692908216911617614c72565b92506101000a81548163ffffffff021916908363ffffffff16021790555061288a4630600d60009054906101000a900463ffffffff1663ffffffff1686600001518760200151886040015189606001518a608001518b60a001516131ca565b600a819055600d5484516020860151604080880151606089015160808a015160a08b015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05986128f3988b98919763ffffffff909116969195909490939092909190614d55565b60405180910390a1600b546601000000000000900463ffffffff1660005b84515181101561296b5781600682601f811061292f5761292f614af9565b600891828204019190066004026101000a81548163ffffffff021916908363ffffffff160217905550808061296390614c8f565b915050612911565b5050505050505050505050565b6001600160a01b03811660009081526002602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046bffffffffffffffffffffffff1691810191909152906129df5750600092915050565b6006816020015160ff16601f81106129f9576129f9614af9565b600881049190910154600b54612a2f926007166004026101000a90910463ffffffff908116916601000000000000900416614b57565b9392505050565b600080808080333214612a8b5760405162461bcd60e51b815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f410000000000000000000000006044820152606401610c47565b5050600a54600b5463ffffffff6601000000000000820481166000908152600c60205260409020549296610100909204600881901c8216965064ffffffffff169450601783900b93507c010000000000000000000000000000000000000000000000000000000090920490911690565b612b03612d8a565b60408051808201909152600e546001600160a01b038082168084527401000000000000000000000000000000000000000090920463ffffffff1660208401528416141580612b6157508163ffffffff16816020015163ffffffff1614155b15612c13576040805180820182526001600160a01b0385811680835263ffffffff8681166020948501819052600e80547fffffffffffffffff00000000000000000000000000000000000000000000000016841774010000000000000000000000000000000000000000830217905586518786015187519316835294820152909392909116917fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541910160405180910390a35b505050565b6001600160a01b03828116600090815260136020526040902054163314612c815760405162461bcd60e51b815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e207570646174650000006044820152606401610c47565b6001600160a01b0381163303612cd95760405162461bcd60e51b815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610c47565b6001600160a01b03808316600090815260146020526040902080548383167fffffffffffffffffffffffff000000000000000000000000000000000000000082168117909255909116908114612c13576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a4505050565b612d70612d8a565b6111f781613dd4565b612d81612d8a565b6111f781613e95565b6000546001600160a01b03163314612de45760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610c47565b565b601154600b54604080516103e08101918290526001600160a01b0390931692660100000000000090920463ffffffff1691600091600690601f908285855b82829054906101000a900463ffffffff1663ffffffff1681526020019060040190602082600301049283019260010382029150808411612e245790505050505050905060006005805480602002602001604051908101604052809291908181526020018280548015612ebf57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612ea1575b5050505050905060005b81518110156131bc57600060026000848481518110612eea57612eea614af9565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160029054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff169050600060026000858581518110612f5657612f56614af9565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160026101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008483601f8110612fc357612fc3614af9565b6020020151600b5490870363ffffffff90811692507201000000000000000000000000000000000000909104168102633b9aca0002820180156131b15760006013600087878151811061301857613018614af9565b6020908102919091018101516001600160a01b0390811683529082019290925260409081016000205490517fa9059cbb00000000000000000000000000000000000000000000000000000000815290821660048201819052602482018590529250908a169063a9059cbb906044016020604051808303816000875af11580156130a5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130c99190614bbe565b6131155760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610c47565b878786601f811061312857613128614af9565b602002019063ffffffff16908163ffffffff1681525050886001600160a01b0316816001600160a01b031687878151811061316557613165614af9565b60200260200101516001600160a01b03167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c856040516131a791815260200190565b60405180910390a4505b505050600101612ec9565b5061161e600683601f6142ea565b6000808a8a8a8a8a8a8a8a8a6040516020016131ee99989796959493929190614deb565b60408051601f1981840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150505b9998505050505050505050565b6001600160a01b0381166000908152600260209081526040918290208251606081018452905460ff80821615158084526101008304909116938301939093526201000090046bffffffffffffffffffffffff16928101929092526132ba575050565b60006132c58361091b565b90508015612c13576001600160a01b03838116600090815260136020526040908190205460115491517fa9059cbb000000000000000000000000000000000000000000000000000000008152908316600482018190526024820185905292919091169063a9059cbb906044016020604051808303816000875af1158015613350573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133749190614bbe565b6133c05760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610c47565b600b60000160069054906101000a900463ffffffff166006846020015160ff16601f81106133f0576133f0614af9565b6008810491909101805460079092166004026101000a63ffffffff8181021990931693909216919091029190911790556001600160a01b0384811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff169055601154915186815291841693851692917fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c910160405180910390a450505050565b60006134b1826020614b7b565b6134bc856020614b7b565b6134c888610144614b92565b6134d29190614b92565b6134dc9190614b92565b6134e7906000614b92565b90503681146135385760405162461bcd60e51b815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d6174636800000000000000006044820152606401610c47565b50505050505050565b60008061354d83613f1c565b9050601f81604001515111156135a55760405162461bcd60e51b815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e647300006044820152606401610c47565b604081015151865160ff16106135fd5760405162461bcd60e51b815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e00006044820152606401610c47565b64ffffffffff84166020870152604081015180516000919061362190600290614e73565b8151811061363157613631614af9565b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b1315801561369757507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b6136e35760405162461bcd60e51b815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e676500006044820152606401610c47565b604087018051906136f382614eae565b63ffffffff1663ffffffff168152505060405180606001604052808260170b8152602001836000015163ffffffff1681526020014263ffffffff16815250600c6000896040015163ffffffff1663ffffffff16815260200190815260200160002060008201518160000160006101000a81548177ffffffffffffffffffffffffffffffffffffffffffffffff021916908360170b77ffffffffffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160186101000a81548163ffffffff021916908363ffffffff160217905550604082015181600001601c6101000a81548163ffffffff021916908363ffffffff16021790555090505086600b60008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548164ffffffffff021916908364ffffffffff16021790555060408201518160000160066101000a81548163ffffffff021916908363ffffffff160217905550606082015181600001600a6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600e6101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160000160126101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160000160166101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600001601a6101000a81548162ffffff021916908362ffffff160217905550905050866040015163ffffffff167fc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a823385600001518660400151876020015188606001518d8d604051613987989796959493929190614ed1565b60405180910390a26040808801518351915163ffffffff9283168152600092909116907f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac602719060200160405180910390a3866040015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f42604051613a1791815260200190565b60405180910390a3613a3087604001518260170b613fc1565b506060015195945050505050565b60008360170b1261204d576000613a70633b9aca003a04866080015163ffffffff16876060015163ffffffff16614111565b90506010360260005a90506000613a998663ffffffff1685858b60e0015162ffffff1686614137565b90506000670de0b6b3a764000077ffffffffffffffffffffffffffffffffffffffffffffffff891683026001600160a01b03881660009081526002602052604090205460c08c01519290910492506201000090046bffffffffffffffffffffffff9081169163ffffffff16633b9aca000282840101908116821115613b24575050505050505061204d565b6001600160a01b038816600090815260026020526040902080546bffffffffffffffffffffffff90921662010000027fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff90921691909117905550505050505050505050565b6000806005805480602002602001604051908101604052809291908181526020018280548015613be257602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311613bc4575b50508351600b54604080516103e08101918290529697509195660100000000000090910463ffffffff169450600093509150600690601f908285855b82829054906101000a900463ffffffff1663ffffffff1681526020019060040190602082600301049283019260010382029150808411613c1e5790505050505050905060005b83811015613cb1578181601f8110613c7e57613c7e614af9565b6020020151613c8d9084614b57565b613c9d9063ffffffff1687614b92565b955080613ca981614c8f565b915050613c64565b50600b54613cdf907201000000000000000000000000000000000000900463ffffffff16633b9aca00614b7b565b613ce99086614b7b565b945060005b83811015613d625760026000868381518110613d0c57613d0c614af9565b6020908102919091018101516001600160a01b0316825281019190915260400160002054613d4e906201000090046bffffffffffffffffffffffff1687614b92565b955080613d5a81614c8f565b915050613cee565b505050505090565b600081831015613d7b575081613d7e565b50805b92915050565b806000106111f75760405162461bcd60e51b815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610c47565b336001600160a01b03821603613e2c5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610c47565b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6012546001600160a01b039081169082168114610cb857601280547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0384811691821790925560408051928416835260208301919091527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d4891291016116b0565b613f506040518060800160405280600063ffffffff1681526020016060815260200160608152602001600060170b81525090565b6000806060600085806020019051810190613f6b9190614f8b565b92965090945092509050613f7f868361419b565b81516040805160208082019690965281519082018252918252805160808101825263ffffffff969096168652938501529183015260170b606082015292915050565b60408051808201909152600e546001600160a01b0381168083527401000000000000000000000000000000000000000090910463ffffffff16602083015261400857505050565b6000614015600185614b57565b63ffffffff8181166000818152600c6020908152604091829020549087015187519251602481019490945260179190910b604484018190528985166064850152608484018990529495506140c593169160a40160408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fbeed9b5100000000000000000000000000000000000000000000000000000000179052614213565b61161e5760405162461bcd60e51b815260206004820152601060248201527f696e73756666696369656e7420676173000000000000000000000000000000006044820152606401610c47565b6000838381101561412457600285850304015b61412e8184613d6a565b95945050505050565b6000818610156141895760405162461bcd60e51b815260206004820181905260248201527f6c6566744761732063616e6e6f742065786365656420696e697469616c4761736044820152606401610c47565b50633b9aca0094039190910101020290565b6000815160206141ab9190614b7b565b6141b69060a0614b92565b6141c1906000614b92565b905080835114612c135760405162461bcd60e51b815260206004820152601660248201527f7265706f7274206c656e677468206d69736d61746368000000000000000000006044820152606401610c47565b60005a61138881106142475761138881039050846040820482031115614247576000808451602086016000888af150600191505b509392505050565b50805460008255906000526020600020908101906111f7919061437d565b8280548282559060005260206000209081019282156142da579160200282015b828111156142da57825182547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0390911617825560209092019160019091019061428d565b506142e692915061437d565b5090565b6004830191839082156142da5791602002820160005b8382111561434457835183826101000a81548163ffffffff021916908363ffffffff1602179055509260200192600401602081600301049283019260010302614300565b80156143745782816101000a81549063ffffffff0219169055600401602081600301049283019260010302614344565b50506142e69291505b5b808211156142e6576000815560010161437e565b6001600160a01b03811681146111f757600080fd5b80356143b281614392565b919050565b6000602082840312156143c957600080fd5b8135612a2f81614392565b6000815180845260005b818110156143fa576020818501810151868301820152016143de565b506000602082860101526020601f19601f83011685010191505092915050565b602081526000612a2f60208301846143d4565b6000806040838503121561444057600080fd5b823561444b81614392565b9150602083013561445b81614392565b809150509250929050565b63ffffffff811681146111f757600080fd5b600080600080600060a0868803121561449057600080fd5b853561449b81614466565b945060208601356144ab81614466565b935060408601356144bb81614466565b925060608601356144cb81614466565b9150608086013562ffffff811681146144e357600080fd5b809150509295509295909350565b600081518084526020808501945080840160005b8381101561452a5781516001600160a01b031687529582019590820190600101614505565b509495945050505050565b602081526000612a2f60208301846144f1565b803567ffffffffffffffff811681146143b257600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff811182821017156145b8576145b8614560565b604052919050565b600067ffffffffffffffff8211156145da576145da614560565b5060051b60200190565b600082601f8301126145f557600080fd5b8135602061460a614605836145c0565b61458f565b82815260059290921b8401810191818101908684111561462957600080fd5b8286015b8481101561464d57803561464081614392565b835291830191830161462d565b509695505050505050565b803560ff811681146143b257600080fd5b60008083601f84011261467b57600080fd5b50813567ffffffffffffffff81111561469357600080fd5b6020830191508360208285010111156146ab57600080fd5b9250929050565b600082601f8301126146c357600080fd5b813567ffffffffffffffff8111156146dd576146dd614560565b6146f06020601f19601f8401160161458f565b81815284602083860101111561470557600080fd5b816020850160208301376000918101602001919091529392505050565b6000806000806000806000806000806101208b8d03121561474257600080fd5b8a35995061475260208c016143a7565b985061476060408c01614548565b975060608b013567ffffffffffffffff8082111561477d57600080fd5b6147898e838f016145e4565b985060808d013591508082111561479f57600080fd5b6147ab8e838f016145e4565b97506147b960a08e01614658565b965060c08d01359150808211156147cf57600080fd5b6147db8e838f01614669565b90965094508491506147ef60e08e01614548565b93506101008d013591508082111561480657600080fd5b506148138d828e016146b2565b9150509295989b9194979a5092959850565b60006020828403121561483757600080fd5b813569ffffffffffffffffffff81168114612a2f57600080fd5b60008083601f84011261486357600080fd5b50813567ffffffffffffffff81111561487b57600080fd5b6020830191508360208260051b85010111156146ab57600080fd5b600080600080604085870312156148ac57600080fd5b843567ffffffffffffffff808211156148c457600080fd5b6148d088838901614851565b909650945060208701359150808211156148e957600080fd5b506148f687828801614851565b95989497509550505050565b60008060008060008060008060e0898b03121561491e57600080fd5b606089018a81111561492f57600080fd5b8998503567ffffffffffffffff8082111561494957600080fd5b6149558c838d01614669565b909950975060808b013591508082111561496e57600080fd5b61497a8c838d01614851565b909750955060a08b013591508082111561499357600080fd5b506149a08b828c01614851565b999c989b50969995989497949560c00135949350505050565b6000602082840312156149cb57600080fd5b5035919050565b600080604083850312156149e557600080fd5b82356149f081614392565b946020939093013593505050565b60008060008060008060c08789031215614a1757600080fd5b863567ffffffffffffffff80821115614a2f57600080fd5b614a3b8a838b016145e4565b97506020890135915080821115614a5157600080fd5b614a5d8a838b016145e4565b9650614a6b60408a01614658565b95506060890135915080821115614a8157600080fd5b614a8d8a838b016146b2565b9450614a9b60808a01614548565b935060a0890135915080821115614ab157600080fd5b50614abe89828a016146b2565b9150509295509295509295565b60008060408385031215614ade57600080fd5b8235614ae981614392565b9150602083013561445b81614466565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b63ffffffff828116828216039080821115614b7457614b74614b28565b5092915050565b8082028115828204841417613d7e57613d7e614b28565b80820180821115613d7e57613d7e614b28565b600060208284031215614bb757600080fd5b5051919050565b600060208284031215614bd057600080fd5b81518015158114612a2f57600080fd5b6001600160a01b038416815260406020820152816040820152818360608301376000818301606090810191909152601f909201601f1916010192915050565b600181811c90821680614c3357607f821691505b602082108103614c6c577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b63ffffffff818116838216019080821115614b7457614b74614b28565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614cc057614cc0614b28565b5060010190565b60ff8181168382160190811115613d7e57613d7e614b28565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b81810381811115613d7e57613d7e614b28565b8181036000831280158383131683831282161715614b7457614b74614b28565b60ff8181168382160290811690818114614b7457614b74614b28565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614d858184018a6144f1565b90508281036080840152614d9981896144f1565b905060ff871660a084015282810360c0840152614db681876143d4565b905067ffffffffffffffff851660e0840152828103610100840152614ddb81856143d4565b9c9b505050505050505050505050565b60006101208b83526001600160a01b038b16602084015267ffffffffffffffff808b166040850152816060850152614e258285018b6144f1565b91508382036080850152614e39828a6144f1565b915060ff881660a085015283820360c0850152614e5682886143d4565b90861660e08501528381036101008501529050614ddb81856143d4565b600082614ea9577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b600063ffffffff808316818103614ec757614ec7614b28565b6001019392505050565b600061010080830160178c810b855260206001600160a01b038d168187015263ffffffff8c1660408701528360608701528293508a5180845261012087019450818c01935060005b81811015614f37578451840b86529482019493820193600101614f19565b50505050508281036080840152614f4e81886143d4565b915050614f6060a083018660170b9052565b8360c083015261324b60e083018464ffffffffff169052565b8051601781900b81146143b257600080fd5b60008060008060808587031215614fa157600080fd5b8451614fac81614466565b809450506020808601519350604086015167ffffffffffffffff811115614fd257600080fd5b8601601f81018813614fe357600080fd5b8051614ff1614605826145c0565b81815260059190911b8201830190838101908a83111561501057600080fd5b928401925b828410156150355761502684614f79565b82529284019290840190615015565b809650505050505061504960608601614f79565b90509295919450925056fea164736f6c6343000813000a",
}

// ExposedOCR2AggregatorABI is the input ABI used to generate the binding from.
// Deprecated: Use ExposedOCR2AggregatorMetaData.ABI instead.
var ExposedOCR2AggregatorABI = ExposedOCR2AggregatorMetaData.ABI

// ExposedOCR2AggregatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ExposedOCR2AggregatorMetaData.Bin instead.
var ExposedOCR2AggregatorBin = ExposedOCR2AggregatorMetaData.Bin

// DeployExposedOCR2Aggregator deploys a new Ethereum contract, binding an instance of ExposedOCR2Aggregator to it.
func DeployExposedOCR2Aggregator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ExposedOCR2Aggregator, error) {
	parsed, err := ExposedOCR2AggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ExposedOCR2AggregatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ExposedOCR2Aggregator{ExposedOCR2AggregatorCaller: ExposedOCR2AggregatorCaller{contract: contract}, ExposedOCR2AggregatorTransactor: ExposedOCR2AggregatorTransactor{contract: contract}, ExposedOCR2AggregatorFilterer: ExposedOCR2AggregatorFilterer{contract: contract}}, nil
}

// ExposedOCR2Aggregator is an auto generated Go binding around an Ethereum contract.
type ExposedOCR2Aggregator struct {
	ExposedOCR2AggregatorCaller     // Read-only binding to the contract
	ExposedOCR2AggregatorTransactor // Write-only binding to the contract
	ExposedOCR2AggregatorFilterer   // Log filterer for contract events
}

// ExposedOCR2AggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExposedOCR2AggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExposedOCR2AggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExposedOCR2AggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExposedOCR2AggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExposedOCR2AggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExposedOCR2AggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExposedOCR2AggregatorSession struct {
	Contract     *ExposedOCR2Aggregator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ExposedOCR2AggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExposedOCR2AggregatorCallerSession struct {
	Contract *ExposedOCR2AggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// ExposedOCR2AggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExposedOCR2AggregatorTransactorSession struct {
	Contract     *ExposedOCR2AggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ExposedOCR2AggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExposedOCR2AggregatorRaw struct {
	Contract *ExposedOCR2Aggregator // Generic contract binding to access the raw methods on
}

// ExposedOCR2AggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExposedOCR2AggregatorCallerRaw struct {
	Contract *ExposedOCR2AggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// ExposedOCR2AggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExposedOCR2AggregatorTransactorRaw struct {
	Contract *ExposedOCR2AggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExposedOCR2Aggregator creates a new instance of ExposedOCR2Aggregator, bound to a specific deployed contract.
func NewExposedOCR2Aggregator(address common.Address, backend bind.ContractBackend) (*ExposedOCR2Aggregator, error) {
	contract, err := bindExposedOCR2Aggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2Aggregator{ExposedOCR2AggregatorCaller: ExposedOCR2AggregatorCaller{contract: contract}, ExposedOCR2AggregatorTransactor: ExposedOCR2AggregatorTransactor{contract: contract}, ExposedOCR2AggregatorFilterer: ExposedOCR2AggregatorFilterer{contract: contract}}, nil
}

// NewExposedOCR2AggregatorCaller creates a new read-only instance of ExposedOCR2Aggregator, bound to a specific deployed contract.
func NewExposedOCR2AggregatorCaller(address common.Address, caller bind.ContractCaller) (*ExposedOCR2AggregatorCaller, error) {
	contract, err := bindExposedOCR2Aggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorCaller{contract: contract}, nil
}

// NewExposedOCR2AggregatorTransactor creates a new write-only instance of ExposedOCR2Aggregator, bound to a specific deployed contract.
func NewExposedOCR2AggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*ExposedOCR2AggregatorTransactor, error) {
	contract, err := bindExposedOCR2Aggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorTransactor{contract: contract}, nil
}

// NewExposedOCR2AggregatorFilterer creates a new log filterer instance of ExposedOCR2Aggregator, bound to a specific deployed contract.
func NewExposedOCR2AggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*ExposedOCR2AggregatorFilterer, error) {
	contract, err := bindExposedOCR2Aggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorFilterer{contract: contract}, nil
}

// bindExposedOCR2Aggregator binds a generic wrapper to an already deployed contract.
func bindExposedOCR2Aggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ExposedOCR2AggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExposedOCR2Aggregator.Contract.ExposedOCR2AggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.ExposedOCR2AggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.ExposedOCR2AggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExposedOCR2Aggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.contract.Transact(opts, method, params...)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) Decimals() (uint8, error) {
	return _ExposedOCR2Aggregator.Contract.Decimals(&_ExposedOCR2Aggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) Decimals() (uint8, error) {
	return _ExposedOCR2Aggregator.Contract.Decimals(&_ExposedOCR2Aggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) Description() (string, error) {
	return _ExposedOCR2Aggregator.Contract.Description(&_ExposedOCR2Aggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) Description() (string, error) {
	return _ExposedOCR2Aggregator.Contract.Description(&_ExposedOCR2Aggregator.CallOpts)
}

// ExposedConfigDigestFromConfigData is a free data retrieval call binding the contract method 0x739376c1.
//
// Solidity: function exposedConfigDigestFromConfigData(uint256 _chainId, address _contractAddress, uint64 _configCount, address[] _signers, address[] _transmitters, uint8 _f, bytes _onchainConfig, uint64 _encodedConfigVersion, bytes _encodedConfig) view returns(bytes32)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) ExposedConfigDigestFromConfigData(opts *bind.CallOpts, _chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "exposedConfigDigestFromConfigData", _chainId, _contractAddress, _configCount, _signers, _transmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ExposedConfigDigestFromConfigData is a free data retrieval call binding the contract method 0x739376c1.
//
// Solidity: function exposedConfigDigestFromConfigData(uint256 _chainId, address _contractAddress, uint64 _configCount, address[] _signers, address[] _transmitters, uint8 _f, bytes _onchainConfig, uint64 _encodedConfigVersion, bytes _encodedConfig) view returns(bytes32)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) ExposedConfigDigestFromConfigData(_chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	return _ExposedOCR2Aggregator.Contract.ExposedConfigDigestFromConfigData(&_ExposedOCR2Aggregator.CallOpts, _chainId, _contractAddress, _configCount, _signers, _transmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)
}

// ExposedConfigDigestFromConfigData is a free data retrieval call binding the contract method 0x739376c1.
//
// Solidity: function exposedConfigDigestFromConfigData(uint256 _chainId, address _contractAddress, uint64 _configCount, address[] _signers, address[] _transmitters, uint8 _f, bytes _onchainConfig, uint64 _encodedConfigVersion, bytes _encodedConfig) view returns(bytes32)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) ExposedConfigDigestFromConfigData(_chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	return _ExposedOCR2Aggregator.Contract.ExposedConfigDigestFromConfigData(&_ExposedOCR2Aggregator.CallOpts, _chainId, _contractAddress, _configCount, _signers, _transmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.GetAnswer(&_ExposedOCR2Aggregator.CallOpts, roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.GetAnswer(&_ExposedOCR2Aggregator.CallOpts, roundId)
}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) GetBilling(opts *bind.CallOpts) (struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
}, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "getBilling")

	outstruct := new(struct {
		MaximumGasPriceGwei       uint32
		ReasonableGasPriceGwei    uint32
		ObservationPaymentGjuels  uint32
		TransmissionPaymentGjuels uint32
		AccountingGas             *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaximumGasPriceGwei = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.ReasonableGasPriceGwei = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ObservationPaymentGjuels = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.TransmissionPaymentGjuels = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.AccountingGas = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) GetBilling() (struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
}, error) {
	return _ExposedOCR2Aggregator.Contract.GetBilling(&_ExposedOCR2Aggregator.CallOpts)
}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) GetBilling() (struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
}, error) {
	return _ExposedOCR2Aggregator.Contract.GetBilling(&_ExposedOCR2Aggregator.CallOpts)
}

// GetBillingAccessController is a free data retrieval call binding the contract method 0xc4c92b37.
//
// Solidity: function getBillingAccessController() view returns(address)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) GetBillingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "getBillingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBillingAccessController is a free data retrieval call binding the contract method 0xc4c92b37.
//
// Solidity: function getBillingAccessController() view returns(address)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) GetBillingAccessController() (common.Address, error) {
	return _ExposedOCR2Aggregator.Contract.GetBillingAccessController(&_ExposedOCR2Aggregator.CallOpts)
}

// GetBillingAccessController is a free data retrieval call binding the contract method 0xc4c92b37.
//
// Solidity: function getBillingAccessController() view returns(address)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) GetBillingAccessController() (common.Address, error) {
	return _ExposedOCR2Aggregator.Contract.GetBillingAccessController(&_ExposedOCR2Aggregator.CallOpts)
}

// GetLinkToken is a free data retrieval call binding the contract method 0xe76d5168.
//
// Solidity: function getLinkToken() view returns(address linkToken)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) GetLinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "getLinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetLinkToken is a free data retrieval call binding the contract method 0xe76d5168.
//
// Solidity: function getLinkToken() view returns(address linkToken)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) GetLinkToken() (common.Address, error) {
	return _ExposedOCR2Aggregator.Contract.GetLinkToken(&_ExposedOCR2Aggregator.CallOpts)
}

// GetLinkToken is a free data retrieval call binding the contract method 0xe76d5168.
//
// Solidity: function getLinkToken() view returns(address linkToken)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) GetLinkToken() (common.Address, error) {
	return _ExposedOCR2Aggregator.Contract.GetLinkToken(&_ExposedOCR2Aggregator.CallOpts)
}

// GetRequesterAccessController is a free data retrieval call binding the contract method 0xdaffc4b5.
//
// Solidity: function getRequesterAccessController() view returns(address)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) GetRequesterAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "getRequesterAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRequesterAccessController is a free data retrieval call binding the contract method 0xdaffc4b5.
//
// Solidity: function getRequesterAccessController() view returns(address)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) GetRequesterAccessController() (common.Address, error) {
	return _ExposedOCR2Aggregator.Contract.GetRequesterAccessController(&_ExposedOCR2Aggregator.CallOpts)
}

// GetRequesterAccessController is a free data retrieval call binding the contract method 0xdaffc4b5.
//
// Solidity: function getRequesterAccessController() view returns(address)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) GetRequesterAccessController() (common.Address, error) {
	return _ExposedOCR2Aggregator.Contract.GetRequesterAccessController(&_ExposedOCR2Aggregator.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80 roundId_, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) GetRoundData(opts *bind.CallOpts, roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "getRoundData", roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
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

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80 roundId_, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) GetRoundData(roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _ExposedOCR2Aggregator.Contract.GetRoundData(&_ExposedOCR2Aggregator.CallOpts, roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80 roundId_, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) GetRoundData(roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _ExposedOCR2Aggregator.Contract.GetRoundData(&_ExposedOCR2Aggregator.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.GetTimestamp(&_ExposedOCR2Aggregator.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.GetTimestamp(&_ExposedOCR2Aggregator.CallOpts, roundId)
}

// GetTransmitters is a free data retrieval call binding the contract method 0x666cab8d.
//
// Solidity: function getTransmitters() view returns(address[])
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) GetTransmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "getTransmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetTransmitters is a free data retrieval call binding the contract method 0x666cab8d.
//
// Solidity: function getTransmitters() view returns(address[])
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) GetTransmitters() ([]common.Address, error) {
	return _ExposedOCR2Aggregator.Contract.GetTransmitters(&_ExposedOCR2Aggregator.CallOpts)
}

// GetTransmitters is a free data retrieval call binding the contract method 0x666cab8d.
//
// Solidity: function getTransmitters() view returns(address[])
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) GetTransmitters() ([]common.Address, error) {
	return _ExposedOCR2Aggregator.Contract.GetTransmitters(&_ExposedOCR2Aggregator.CallOpts)
}

// GetValidatorConfig is a free data retrieval call binding the contract method 0x9bd2c0b1.
//
// Solidity: function getValidatorConfig() view returns(address validator, uint32 gasLimit)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) GetValidatorConfig(opts *bind.CallOpts) (struct {
	Validator common.Address
	GasLimit  uint32
}, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "getValidatorConfig")

	outstruct := new(struct {
		Validator common.Address
		GasLimit  uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.GasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return *outstruct, err

}

// GetValidatorConfig is a free data retrieval call binding the contract method 0x9bd2c0b1.
//
// Solidity: function getValidatorConfig() view returns(address validator, uint32 gasLimit)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) GetValidatorConfig() (struct {
	Validator common.Address
	GasLimit  uint32
}, error) {
	return _ExposedOCR2Aggregator.Contract.GetValidatorConfig(&_ExposedOCR2Aggregator.CallOpts)
}

// GetValidatorConfig is a free data retrieval call binding the contract method 0x9bd2c0b1.
//
// Solidity: function getValidatorConfig() view returns(address validator, uint32 gasLimit)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) GetValidatorConfig() (struct {
	Validator common.Address
	GasLimit  uint32
}, error) {
	return _ExposedOCR2Aggregator.Contract.GetValidatorConfig(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) LatestAnswer() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.LatestAnswer(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.LatestAnswer(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _ExposedOCR2Aggregator.Contract.LatestConfigDetails(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _ExposedOCR2Aggregator.Contract.LatestConfigDetails(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(struct {
		ScanLogs     bool
		ConfigDigest [32]byte
		Epoch        uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _ExposedOCR2Aggregator.Contract.LatestConfigDigestAndEpoch(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _ExposedOCR2Aggregator.Contract.LatestConfigDigestAndEpoch(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) LatestRound() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.LatestRound(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.LatestRound(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
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

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _ExposedOCR2Aggregator.Contract.LatestRoundData(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _ExposedOCR2Aggregator.Contract.LatestRoundData(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.LatestTimestamp(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.LatestTimestamp(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes32 configDigest, uint32 epoch, uint8 round, int192 latestAnswer_, uint64 latestTimestamp_)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) LatestTransmissionDetails(opts *bind.CallOpts) (struct {
	ConfigDigest    [32]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "latestTransmissionDetails")

	outstruct := new(struct {
		ConfigDigest    [32]byte
		Epoch           uint32
		Round           uint8
		LatestAnswer    *big.Int
		LatestTimestamp uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigDigest = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.Round = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.LatestAnswer = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.LatestTimestamp = *abi.ConvertType(out[4], new(uint64)).(*uint64)

	return *outstruct, err

}

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes32 configDigest, uint32 epoch, uint8 round, int192 latestAnswer_, uint64 latestTimestamp_)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [32]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _ExposedOCR2Aggregator.Contract.LatestTransmissionDetails(&_ExposedOCR2Aggregator.CallOpts)
}

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes32 configDigest, uint32 epoch, uint8 round, int192 latestAnswer_, uint64 latestTimestamp_)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [32]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _ExposedOCR2Aggregator.Contract.LatestTransmissionDetails(&_ExposedOCR2Aggregator.CallOpts)
}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.LinkAvailableForPayment(&_ExposedOCR2Aggregator.CallOpts)
}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.LinkAvailableForPayment(&_ExposedOCR2Aggregator.CallOpts)
}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int192)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "maxAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int192)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) MaxAnswer() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.MaxAnswer(&_ExposedOCR2Aggregator.CallOpts)
}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int192)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.MaxAnswer(&_ExposedOCR2Aggregator.CallOpts)
}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int192)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "minAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int192)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) MinAnswer() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.MinAnswer(&_ExposedOCR2Aggregator.CallOpts)
}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int192)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.MinAnswer(&_ExposedOCR2Aggregator.CallOpts)
}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address transmitterAddress) view returns(uint32)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) OracleObservationCount(opts *bind.CallOpts, transmitterAddress common.Address) (uint32, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "oracleObservationCount", transmitterAddress)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address transmitterAddress) view returns(uint32)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _ExposedOCR2Aggregator.Contract.OracleObservationCount(&_ExposedOCR2Aggregator.CallOpts, transmitterAddress)
}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address transmitterAddress) view returns(uint32)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _ExposedOCR2Aggregator.Contract.OracleObservationCount(&_ExposedOCR2Aggregator.CallOpts, transmitterAddress)
}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address transmitterAddress) view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) OwedPayment(opts *bind.CallOpts, transmitterAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "owedPayment", transmitterAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address transmitterAddress) view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.OwedPayment(&_ExposedOCR2Aggregator.CallOpts, transmitterAddress)
}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address transmitterAddress) view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.OwedPayment(&_ExposedOCR2Aggregator.CallOpts, transmitterAddress)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) Owner() (common.Address, error) {
	return _ExposedOCR2Aggregator.Contract.Owner(&_ExposedOCR2Aggregator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) Owner() (common.Address, error) {
	return _ExposedOCR2Aggregator.Contract.Owner(&_ExposedOCR2Aggregator.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) TypeAndVersion() (string, error) {
	return _ExposedOCR2Aggregator.Contract.TypeAndVersion(&_ExposedOCR2Aggregator.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) TypeAndVersion() (string, error) {
	return _ExposedOCR2Aggregator.Contract.TypeAndVersion(&_ExposedOCR2Aggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ExposedOCR2Aggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) Version() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.Version(&_ExposedOCR2Aggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorCallerSession) Version() (*big.Int, error) {
	return _ExposedOCR2Aggregator.Contract.Version(&_ExposedOCR2Aggregator.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.AcceptOwnership(&_ExposedOCR2Aggregator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.AcceptOwnership(&_ExposedOCR2Aggregator.TransactOpts)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "acceptPayeeship", transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.AcceptPayeeship(&_ExposedOCR2Aggregator.TransactOpts, transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.AcceptPayeeship(&_ExposedOCR2Aggregator.TransactOpts, transmitter)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "requestNewRound")
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) RequestNewRound() (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.RequestNewRound(&_ExposedOCR2Aggregator.TransactOpts)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) RequestNewRound() (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.RequestNewRound(&_ExposedOCR2Aggregator.TransactOpts)
}

// SetBilling is a paid mutator transaction binding the contract method 0x643dc105.
//
// Solidity: function setBilling(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) SetBilling(opts *bind.TransactOpts, maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "setBilling", maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

// SetBilling is a paid mutator transaction binding the contract method 0x643dc105.
//
// Solidity: function setBilling(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) SetBilling(maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetBilling(&_ExposedOCR2Aggregator.TransactOpts, maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

// SetBilling is a paid mutator transaction binding the contract method 0x643dc105.
//
// Solidity: function setBilling(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) SetBilling(maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetBilling(&_ExposedOCR2Aggregator.TransactOpts, maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "setBillingAccessController", _billingAccessController)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetBillingAccessController(&_ExposedOCR2Aggregator.TransactOpts, _billingAccessController)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetBillingAccessController(&_ExposedOCR2Aggregator.TransactOpts, _billingAccessController)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetConfig(&_ExposedOCR2Aggregator.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetConfig(&_ExposedOCR2Aggregator.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetLinkToken is a paid mutator transaction binding the contract method 0x4fb17470.
//
// Solidity: function setLinkToken(address linkToken, address recipient) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) SetLinkToken(opts *bind.TransactOpts, linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "setLinkToken", linkToken, recipient)
}

// SetLinkToken is a paid mutator transaction binding the contract method 0x4fb17470.
//
// Solidity: function setLinkToken(address linkToken, address recipient) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) SetLinkToken(linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetLinkToken(&_ExposedOCR2Aggregator.TransactOpts, linkToken, recipient)
}

// SetLinkToken is a paid mutator transaction binding the contract method 0x4fb17470.
//
// Solidity: function setLinkToken(address linkToken, address recipient) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) SetLinkToken(linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetLinkToken(&_ExposedOCR2Aggregator.TransactOpts, linkToken, recipient)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] transmitters, address[] payees) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) SetPayees(opts *bind.TransactOpts, transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "setPayees", transmitters, payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] transmitters, address[] payees) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetPayees(&_ExposedOCR2Aggregator.TransactOpts, transmitters, payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] transmitters, address[] payees) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetPayees(&_ExposedOCR2Aggregator.TransactOpts, transmitters, payees)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address requesterAccessController) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) SetRequesterAccessController(opts *bind.TransactOpts, requesterAccessController common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "setRequesterAccessController", requesterAccessController)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address requesterAccessController) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) SetRequesterAccessController(requesterAccessController common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetRequesterAccessController(&_ExposedOCR2Aggregator.TransactOpts, requesterAccessController)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address requesterAccessController) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) SetRequesterAccessController(requesterAccessController common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetRequesterAccessController(&_ExposedOCR2Aggregator.TransactOpts, requesterAccessController)
}

// SetValidatorConfig is a paid mutator transaction binding the contract method 0xeb457163.
//
// Solidity: function setValidatorConfig(address newValidator, uint32 newGasLimit) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) SetValidatorConfig(opts *bind.TransactOpts, newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "setValidatorConfig", newValidator, newGasLimit)
}

// SetValidatorConfig is a paid mutator transaction binding the contract method 0xeb457163.
//
// Solidity: function setValidatorConfig(address newValidator, uint32 newGasLimit) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) SetValidatorConfig(newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetValidatorConfig(&_ExposedOCR2Aggregator.TransactOpts, newValidator, newGasLimit)
}

// SetValidatorConfig is a paid mutator transaction binding the contract method 0xeb457163.
//
// Solidity: function setValidatorConfig(address newValidator, uint32 newGasLimit) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) SetValidatorConfig(newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.SetValidatorConfig(&_ExposedOCR2Aggregator.TransactOpts, newValidator, newGasLimit)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.TransferOwnership(&_ExposedOCR2Aggregator.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.TransferOwnership(&_ExposedOCR2Aggregator.TransactOpts, to)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.TransferPayeeship(&_ExposedOCR2Aggregator.TransactOpts, transmitter, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.TransferPayeeship(&_ExposedOCR2Aggregator.TransactOpts, transmitter, proposed)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.Transmit(&_ExposedOCR2Aggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.Transmit(&_ExposedOCR2Aggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address recipient, uint256 amount) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "withdrawFunds", recipient, amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address recipient, uint256 amount) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.WithdrawFunds(&_ExposedOCR2Aggregator.TransactOpts, recipient, amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address recipient, uint256 amount) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.WithdrawFunds(&_ExposedOCR2Aggregator.TransactOpts, recipient, amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address transmitter) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.contract.Transact(opts, "withdrawPayment", transmitter)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address transmitter) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.WithdrawPayment(&_ExposedOCR2Aggregator.TransactOpts, transmitter)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address transmitter) returns()
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorTransactorSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _ExposedOCR2Aggregator.Contract.WithdrawPayment(&_ExposedOCR2Aggregator.TransactOpts, transmitter)
}

// ExposedOCR2AggregatorAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorAnswerUpdatedIterator struct {
	Event *ExposedOCR2AggregatorAnswerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorAnswerUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorAnswerUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorAnswerUpdated represents a AnswerUpdated event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*ExposedOCR2AggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorAnswerUpdatedIterator{contract: _ExposedOCR2Aggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorAnswerUpdated)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseAnswerUpdated(log types.Log) (*ExposedOCR2AggregatorAnswerUpdated, error) {
	event := new(ExposedOCR2AggregatorAnswerUpdated)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorBillingAccessControllerSetIterator is returned from FilterBillingAccessControllerSet and is used to iterate over the raw logs and unpacked data for BillingAccessControllerSet events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorBillingAccessControllerSetIterator struct {
	Event *ExposedOCR2AggregatorBillingAccessControllerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorBillingAccessControllerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorBillingAccessControllerSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorBillingAccessControllerSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorBillingAccessControllerSet represents a BillingAccessControllerSet event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBillingAccessControllerSet is a free log retrieval operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*ExposedOCR2AggregatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorBillingAccessControllerSetIterator{contract: _ExposedOCR2Aggregator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}

// WatchBillingAccessControllerSet is a free log subscription operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorBillingAccessControllerSet)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBillingAccessControllerSet is a log parse operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*ExposedOCR2AggregatorBillingAccessControllerSet, error) {
	event := new(ExposedOCR2AggregatorBillingAccessControllerSet)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorBillingSetIterator is returned from FilterBillingSet and is used to iterate over the raw logs and unpacked data for BillingSet events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorBillingSetIterator struct {
	Event *ExposedOCR2AggregatorBillingSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorBillingSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorBillingSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorBillingSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorBillingSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorBillingSet represents a BillingSet event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorBillingSet struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterBillingSet is a free log retrieval operation binding the contract event 0x0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f.
//
// Solidity: event BillingSet(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*ExposedOCR2AggregatorBillingSetIterator, error) {

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorBillingSetIterator{contract: _ExposedOCR2Aggregator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}

// WatchBillingSet is a free log subscription operation binding the contract event 0x0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f.
//
// Solidity: event BillingSet(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorBillingSet) (event.Subscription, error) {

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorBillingSet)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBillingSet is a log parse operation binding the contract event 0x0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f.
//
// Solidity: event BillingSet(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseBillingSet(log types.Log) (*ExposedOCR2AggregatorBillingSet, error) {
	event := new(ExposedOCR2AggregatorBillingSet)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorConfigSetIterator struct {
	Event *ExposedOCR2AggregatorConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorConfigSet represents a ConfigSet event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*ExposedOCR2AggregatorConfigSetIterator, error) {

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorConfigSetIterator{contract: _ExposedOCR2Aggregator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorConfigSet)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConfigSet is a log parse operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseConfigSet(log types.Log) (*ExposedOCR2AggregatorConfigSet, error) {
	event := new(ExposedOCR2AggregatorConfigSet)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorLinkTokenSetIterator is returned from FilterLinkTokenSet and is used to iterate over the raw logs and unpacked data for LinkTokenSet events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorLinkTokenSetIterator struct {
	Event *ExposedOCR2AggregatorLinkTokenSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorLinkTokenSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorLinkTokenSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorLinkTokenSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorLinkTokenSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorLinkTokenSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorLinkTokenSet represents a LinkTokenSet event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorLinkTokenSet struct {
	OldLinkToken common.Address
	NewLinkToken common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLinkTokenSet is a free log retrieval operation binding the contract event 0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a.
//
// Solidity: event LinkTokenSet(address indexed oldLinkToken, address indexed newLinkToken)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterLinkTokenSet(opts *bind.FilterOpts, oldLinkToken []common.Address, newLinkToken []common.Address) (*ExposedOCR2AggregatorLinkTokenSetIterator, error) {

	var oldLinkTokenRule []interface{}
	for _, oldLinkTokenItem := range oldLinkToken {
		oldLinkTokenRule = append(oldLinkTokenRule, oldLinkTokenItem)
	}
	var newLinkTokenRule []interface{}
	for _, newLinkTokenItem := range newLinkToken {
		newLinkTokenRule = append(newLinkTokenRule, newLinkTokenItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "LinkTokenSet", oldLinkTokenRule, newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorLinkTokenSetIterator{contract: _ExposedOCR2Aggregator.contract, event: "LinkTokenSet", logs: logs, sub: sub}, nil
}

// WatchLinkTokenSet is a free log subscription operation binding the contract event 0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a.
//
// Solidity: event LinkTokenSet(address indexed oldLinkToken, address indexed newLinkToken)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchLinkTokenSet(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorLinkTokenSet, oldLinkToken []common.Address, newLinkToken []common.Address) (event.Subscription, error) {

	var oldLinkTokenRule []interface{}
	for _, oldLinkTokenItem := range oldLinkToken {
		oldLinkTokenRule = append(oldLinkTokenRule, oldLinkTokenItem)
	}
	var newLinkTokenRule []interface{}
	for _, newLinkTokenItem := range newLinkToken {
		newLinkTokenRule = append(newLinkTokenRule, newLinkTokenItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "LinkTokenSet", oldLinkTokenRule, newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorLinkTokenSet)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLinkTokenSet is a log parse operation binding the contract event 0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a.
//
// Solidity: event LinkTokenSet(address indexed oldLinkToken, address indexed newLinkToken)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseLinkTokenSet(log types.Log) (*ExposedOCR2AggregatorLinkTokenSet, error) {
	event := new(ExposedOCR2AggregatorLinkTokenSet)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorNewRoundIterator struct {
	Event *ExposedOCR2AggregatorNewRound // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorNewRound)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorNewRound)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorNewRound represents a NewRound event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*ExposedOCR2AggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorNewRoundIterator{contract: _ExposedOCR2Aggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorNewRound)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseNewRound(log types.Log) (*ExposedOCR2AggregatorNewRound, error) {
	event := new(ExposedOCR2AggregatorNewRound)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorNewTransmissionIterator is returned from FilterNewTransmission and is used to iterate over the raw logs and unpacked data for NewTransmission events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorNewTransmissionIterator struct {
	Event *ExposedOCR2AggregatorNewTransmission // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorNewTransmissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorNewTransmission)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorNewTransmission)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorNewTransmissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorNewTransmission represents a NewTransmission event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorNewTransmission struct {
	AggregatorRoundId     uint32
	Answer                *big.Int
	Transmitter           common.Address
	ObservationsTimestamp uint32
	Observations          []*big.Int
	Observers             []byte
	JuelsPerFeeCoin       *big.Int
	ConfigDigest          [32]byte
	EpochAndRound         *big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterNewTransmission is a free log retrieval operation binding the contract event 0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, uint32 observationsTimestamp, int192[] observations, bytes observers, int192 juelsPerFeeCoin, bytes32 configDigest, uint40 epochAndRound)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*ExposedOCR2AggregatorNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorNewTransmissionIterator{contract: _ExposedOCR2Aggregator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

// WatchNewTransmission is a free log subscription operation binding the contract event 0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, uint32 observationsTimestamp, int192[] observations, bytes observers, int192 juelsPerFeeCoin, bytes32 configDigest, uint40 epochAndRound)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorNewTransmission)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewTransmission is a log parse operation binding the contract event 0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, uint32 observationsTimestamp, int192[] observations, bytes observers, int192 juelsPerFeeCoin, bytes32 configDigest, uint40 epochAndRound)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseNewTransmission(log types.Log) (*ExposedOCR2AggregatorNewTransmission, error) {
	event := new(ExposedOCR2AggregatorNewTransmission)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorOraclePaidIterator is returned from FilterOraclePaid and is used to iterate over the raw logs and unpacked data for OraclePaid events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorOraclePaidIterator struct {
	Event *ExposedOCR2AggregatorOraclePaid // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorOraclePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorOraclePaid)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorOraclePaid)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorOraclePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorOraclePaid represents a OraclePaid event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	LinkToken   common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOraclePaid is a free log retrieval operation binding the contract event 0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c.
//
// Solidity: event OraclePaid(address indexed transmitter, address indexed payee, uint256 amount, address indexed linkToken)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*ExposedOCR2AggregatorOraclePaidIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	var linkTokenRule []interface{}
	for _, linkTokenItem := range linkToken {
		linkTokenRule = append(linkTokenRule, linkTokenItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorOraclePaidIterator{contract: _ExposedOCR2Aggregator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}

// WatchOraclePaid is a free log subscription operation binding the contract event 0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c.
//
// Solidity: event OraclePaid(address indexed transmitter, address indexed payee, uint256 amount, address indexed linkToken)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	var linkTokenRule []interface{}
	for _, linkTokenItem := range linkToken {
		linkTokenRule = append(linkTokenRule, linkTokenItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorOraclePaid)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOraclePaid is a log parse operation binding the contract event 0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c.
//
// Solidity: event OraclePaid(address indexed transmitter, address indexed payee, uint256 amount, address indexed linkToken)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseOraclePaid(log types.Log) (*ExposedOCR2AggregatorOraclePaid, error) {
	event := new(ExposedOCR2AggregatorOraclePaid)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorOwnershipTransferRequestedIterator struct {
	Event *ExposedOCR2AggregatorOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ExposedOCR2AggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorOwnershipTransferRequestedIterator{contract: _ExposedOCR2Aggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorOwnershipTransferRequested)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*ExposedOCR2AggregatorOwnershipTransferRequested, error) {
	event := new(ExposedOCR2AggregatorOwnershipTransferRequested)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorOwnershipTransferredIterator struct {
	Event *ExposedOCR2AggregatorOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorOwnershipTransferred represents a OwnershipTransferred event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ExposedOCR2AggregatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorOwnershipTransferredIterator{contract: _ExposedOCR2Aggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorOwnershipTransferred)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*ExposedOCR2AggregatorOwnershipTransferred, error) {
	event := new(ExposedOCR2AggregatorOwnershipTransferred)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorPayeeshipTransferRequestedIterator is returned from FilterPayeeshipTransferRequested and is used to iterate over the raw logs and unpacked data for PayeeshipTransferRequested events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorPayeeshipTransferRequestedIterator struct {
	Event *ExposedOCR2AggregatorPayeeshipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorPayeeshipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorPayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorPayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorPayeeshipTransferRequested represents a PayeeshipTransferRequested event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferRequested is a free log retrieval operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*ExposedOCR2AggregatorPayeeshipTransferRequestedIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorPayeeshipTransferRequestedIterator{contract: _ExposedOCR2Aggregator.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferRequested is a free log subscription operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorPayeeshipTransferRequested)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePayeeshipTransferRequested is a log parse operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*ExposedOCR2AggregatorPayeeshipTransferRequested, error) {
	event := new(ExposedOCR2AggregatorPayeeshipTransferRequested)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorPayeeshipTransferredIterator is returned from FilterPayeeshipTransferred and is used to iterate over the raw logs and unpacked data for PayeeshipTransferred events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorPayeeshipTransferredIterator struct {
	Event *ExposedOCR2AggregatorPayeeshipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorPayeeshipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorPayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorPayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorPayeeshipTransferred represents a PayeeshipTransferred event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferred is a free log retrieval operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*ExposedOCR2AggregatorPayeeshipTransferredIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorPayeeshipTransferredIterator{contract: _ExposedOCR2Aggregator.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferred is a free log subscription operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorPayeeshipTransferred)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePayeeshipTransferred is a log parse operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParsePayeeshipTransferred(log types.Log) (*ExposedOCR2AggregatorPayeeshipTransferred, error) {
	event := new(ExposedOCR2AggregatorPayeeshipTransferred)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorRequesterAccessControllerSetIterator is returned from FilterRequesterAccessControllerSet and is used to iterate over the raw logs and unpacked data for RequesterAccessControllerSet events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorRequesterAccessControllerSetIterator struct {
	Event *ExposedOCR2AggregatorRequesterAccessControllerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorRequesterAccessControllerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorRequesterAccessControllerSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorRequesterAccessControllerSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorRequesterAccessControllerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorRequesterAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorRequesterAccessControllerSet represents a RequesterAccessControllerSet event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorRequesterAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRequesterAccessControllerSet is a free log retrieval operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*ExposedOCR2AggregatorRequesterAccessControllerSetIterator, error) {

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorRequesterAccessControllerSetIterator{contract: _ExposedOCR2Aggregator.contract, event: "RequesterAccessControllerSet", logs: logs, sub: sub}, nil
}

// WatchRequesterAccessControllerSet is a free log subscription operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorRequesterAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorRequesterAccessControllerSet)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequesterAccessControllerSet is a log parse operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseRequesterAccessControllerSet(log types.Log) (*ExposedOCR2AggregatorRequesterAccessControllerSet, error) {
	event := new(ExposedOCR2AggregatorRequesterAccessControllerSet)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorRoundRequestedIterator is returned from FilterRoundRequested and is used to iterate over the raw logs and unpacked data for RoundRequested events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorRoundRequestedIterator struct {
	Event *ExposedOCR2AggregatorRoundRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorRoundRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorRoundRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorRoundRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorRoundRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorRoundRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorRoundRequested represents a RoundRequested event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorRoundRequested struct {
	Requester    common.Address
	ConfigDigest [32]byte
	Epoch        uint32
	Round        uint8
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRoundRequested is a free log retrieval operation binding the contract event 0x41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c.
//
// Solidity: event RoundRequested(address indexed requester, bytes32 configDigest, uint32 epoch, uint8 round)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*ExposedOCR2AggregatorRoundRequestedIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorRoundRequestedIterator{contract: _ExposedOCR2Aggregator.contract, event: "RoundRequested", logs: logs, sub: sub}, nil
}

// WatchRoundRequested is a free log subscription operation binding the contract event 0x41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c.
//
// Solidity: event RoundRequested(address indexed requester, bytes32 configDigest, uint32 epoch, uint8 round)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorRoundRequested, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorRoundRequested)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoundRequested is a log parse operation binding the contract event 0x41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c.
//
// Solidity: event RoundRequested(address indexed requester, bytes32 configDigest, uint32 epoch, uint8 round)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseRoundRequested(log types.Log) (*ExposedOCR2AggregatorRoundRequested, error) {
	event := new(ExposedOCR2AggregatorRoundRequested)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorTransmittedIterator is returned from FilterTransmitted and is used to iterate over the raw logs and unpacked data for Transmitted events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorTransmittedIterator struct {
	Event *ExposedOCR2AggregatorTransmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorTransmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorTransmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorTransmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorTransmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorTransmitted represents a Transmitted event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransmitted is a free log retrieval operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterTransmitted(opts *bind.FilterOpts) (*ExposedOCR2AggregatorTransmittedIterator, error) {

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorTransmittedIterator{contract: _ExposedOCR2Aggregator.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

// WatchTransmitted is a free log subscription operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorTransmitted) (event.Subscription, error) {

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorTransmitted)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "Transmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransmitted is a log parse operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseTransmitted(log types.Log) (*ExposedOCR2AggregatorTransmitted, error) {
	event := new(ExposedOCR2AggregatorTransmitted)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExposedOCR2AggregatorValidatorConfigSetIterator is returned from FilterValidatorConfigSet and is used to iterate over the raw logs and unpacked data for ValidatorConfigSet events raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorValidatorConfigSetIterator struct {
	Event *ExposedOCR2AggregatorValidatorConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExposedOCR2AggregatorValidatorConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOCR2AggregatorValidatorConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExposedOCR2AggregatorValidatorConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExposedOCR2AggregatorValidatorConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExposedOCR2AggregatorValidatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExposedOCR2AggregatorValidatorConfigSet represents a ValidatorConfigSet event raised by the ExposedOCR2Aggregator contract.
type ExposedOCR2AggregatorValidatorConfigSet struct {
	PreviousValidator common.Address
	PreviousGasLimit  uint32
	CurrentValidator  common.Address
	CurrentGasLimit   uint32
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterValidatorConfigSet is a free log retrieval operation binding the contract event 0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541.
//
// Solidity: event ValidatorConfigSet(address indexed previousValidator, uint32 previousGasLimit, address indexed currentValidator, uint32 currentGasLimit)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) FilterValidatorConfigSet(opts *bind.FilterOpts, previousValidator []common.Address, currentValidator []common.Address) (*ExposedOCR2AggregatorValidatorConfigSetIterator, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.FilterLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOCR2AggregatorValidatorConfigSetIterator{contract: _ExposedOCR2Aggregator.contract, event: "ValidatorConfigSet", logs: logs, sub: sub}, nil
}

// WatchValidatorConfigSet is a free log subscription operation binding the contract event 0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541.
//
// Solidity: event ValidatorConfigSet(address indexed previousValidator, uint32 previousGasLimit, address indexed currentValidator, uint32 currentGasLimit)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) WatchValidatorConfigSet(opts *bind.WatchOpts, sink chan<- *ExposedOCR2AggregatorValidatorConfigSet, previousValidator []common.Address, currentValidator []common.Address) (event.Subscription, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _ExposedOCR2Aggregator.contract.WatchLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExposedOCR2AggregatorValidatorConfigSet)
				if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorConfigSet is a log parse operation binding the contract event 0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541.
//
// Solidity: event ValidatorConfigSet(address indexed previousValidator, uint32 previousGasLimit, address indexed currentValidator, uint32 currentGasLimit)
func (_ExposedOCR2Aggregator *ExposedOCR2AggregatorFilterer) ParseValidatorConfigSet(log types.Log) (*ExposedOCR2AggregatorValidatorConfigSet, error) {
	event := new(ExposedOCR2AggregatorValidatorConfigSet)
	if err := _ExposedOCR2Aggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LinkTokenInterfaceMetaData contains all meta data concerning the LinkTokenInterface contract.
var LinkTokenInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"remaining\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"decimalPlaces\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"tokenName\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"tokenSymbol\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalTokensIssued\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"transferAndCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// LinkTokenInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use LinkTokenInterfaceMetaData.ABI instead.
var LinkTokenInterfaceABI = LinkTokenInterfaceMetaData.ABI

// LinkTokenInterface is an auto generated Go binding around an Ethereum contract.
type LinkTokenInterface struct {
	LinkTokenInterfaceCaller     // Read-only binding to the contract
	LinkTokenInterfaceTransactor // Write-only binding to the contract
	LinkTokenInterfaceFilterer   // Log filterer for contract events
}

// LinkTokenInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type LinkTokenInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkTokenInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LinkTokenInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkTokenInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LinkTokenInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkTokenInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LinkTokenInterfaceSession struct {
	Contract     *LinkTokenInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// LinkTokenInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LinkTokenInterfaceCallerSession struct {
	Contract *LinkTokenInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// LinkTokenInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LinkTokenInterfaceTransactorSession struct {
	Contract     *LinkTokenInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// LinkTokenInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type LinkTokenInterfaceRaw struct {
	Contract *LinkTokenInterface // Generic contract binding to access the raw methods on
}

// LinkTokenInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LinkTokenInterfaceCallerRaw struct {
	Contract *LinkTokenInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// LinkTokenInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LinkTokenInterfaceTransactorRaw struct {
	Contract *LinkTokenInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLinkTokenInterface creates a new instance of LinkTokenInterface, bound to a specific deployed contract.
func NewLinkTokenInterface(address common.Address, backend bind.ContractBackend) (*LinkTokenInterface, error) {
	contract, err := bindLinkTokenInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterface{LinkTokenInterfaceCaller: LinkTokenInterfaceCaller{contract: contract}, LinkTokenInterfaceTransactor: LinkTokenInterfaceTransactor{contract: contract}, LinkTokenInterfaceFilterer: LinkTokenInterfaceFilterer{contract: contract}}, nil
}

// NewLinkTokenInterfaceCaller creates a new read-only instance of LinkTokenInterface, bound to a specific deployed contract.
func NewLinkTokenInterfaceCaller(address common.Address, caller bind.ContractCaller) (*LinkTokenInterfaceCaller, error) {
	contract, err := bindLinkTokenInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterfaceCaller{contract: contract}, nil
}

// NewLinkTokenInterfaceTransactor creates a new write-only instance of LinkTokenInterface, bound to a specific deployed contract.
func NewLinkTokenInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*LinkTokenInterfaceTransactor, error) {
	contract, err := bindLinkTokenInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterfaceTransactor{contract: contract}, nil
}

// NewLinkTokenInterfaceFilterer creates a new log filterer instance of LinkTokenInterface, bound to a specific deployed contract.
func NewLinkTokenInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*LinkTokenInterfaceFilterer, error) {
	contract, err := bindLinkTokenInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterfaceFilterer{contract: contract}, nil
}

// bindLinkTokenInterface binds a generic wrapper to an already deployed contract.
func bindLinkTokenInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LinkTokenInterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LinkTokenInterface *LinkTokenInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LinkTokenInterface *LinkTokenInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LinkTokenInterface *LinkTokenInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LinkTokenInterface *LinkTokenInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkTokenInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LinkTokenInterface *LinkTokenInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LinkTokenInterface *LinkTokenInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256 remaining)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256 remaining)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.Allowance(&_LinkTokenInterface.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256 remaining)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.Allowance(&_LinkTokenInterface.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256 balance)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256 balance)
func (_LinkTokenInterface *LinkTokenInterfaceSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.BalanceOf(&_LinkTokenInterface.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256 balance)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.BalanceOf(&_LinkTokenInterface.CallOpts, owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 decimalPlaces)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 decimalPlaces)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Decimals() (uint8, error) {
	return _LinkTokenInterface.Contract.Decimals(&_LinkTokenInterface.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 decimalPlaces)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Decimals() (uint8, error) {
	return _LinkTokenInterface.Contract.Decimals(&_LinkTokenInterface.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string tokenName)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string tokenName)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Name() (string, error) {
	return _LinkTokenInterface.Contract.Name(&_LinkTokenInterface.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string tokenName)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Name() (string, error) {
	return _LinkTokenInterface.Contract.Name(&_LinkTokenInterface.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string tokenSymbol)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string tokenSymbol)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Symbol() (string, error) {
	return _LinkTokenInterface.Contract.Symbol(&_LinkTokenInterface.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string tokenSymbol)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Symbol() (string, error) {
	return _LinkTokenInterface.Contract.Symbol(&_LinkTokenInterface.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256 totalTokensIssued)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256 totalTokensIssued)
func (_LinkTokenInterface *LinkTokenInterfaceSession) TotalSupply() (*big.Int, error) {
	return _LinkTokenInterface.Contract.TotalSupply(&_LinkTokenInterface.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256 totalTokensIssued)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) TotalSupply() (*big.Int, error) {
	return _LinkTokenInterface.Contract.TotalSupply(&_LinkTokenInterface.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Approve(&_LinkTokenInterface.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Approve(&_LinkTokenInterface.TransactOpts, spender, value)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(address spender, uint256 addedValue) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) DecreaseApproval(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "decreaseApproval", spender, addedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(address spender, uint256 addedValue) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) DecreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.DecreaseApproval(&_LinkTokenInterface.TransactOpts, spender, addedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(address spender, uint256 addedValue) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) DecreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.DecreaseApproval(&_LinkTokenInterface.TransactOpts, spender, addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(address spender, uint256 subtractedValue) returns()
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) IncreaseApproval(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "increaseApproval", spender, subtractedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(address spender, uint256 subtractedValue) returns()
func (_LinkTokenInterface *LinkTokenInterfaceSession) IncreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.IncreaseApproval(&_LinkTokenInterface.TransactOpts, spender, subtractedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(address spender, uint256 subtractedValue) returns()
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) IncreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.IncreaseApproval(&_LinkTokenInterface.TransactOpts, spender, subtractedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Transfer(&_LinkTokenInterface.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Transfer(&_LinkTokenInterface.TransactOpts, to, value)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) TransferAndCall(opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "transferAndCall", to, value, data)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) TransferAndCall(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferAndCall(&_LinkTokenInterface.TransactOpts, to, value, data)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) TransferAndCall(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferAndCall(&_LinkTokenInterface.TransactOpts, to, value, data)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferFrom(&_LinkTokenInterface.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferFrom(&_LinkTokenInterface.TransactOpts, from, to, value)
}

// OCR2AbstractMetaData contains all meta data concerning the OCR2Abstract contract.
var OCR2AbstractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// OCR2AbstractABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR2AbstractMetaData.ABI instead.
var OCR2AbstractABI = OCR2AbstractMetaData.ABI

// OCR2Abstract is an auto generated Go binding around an Ethereum contract.
type OCR2Abstract struct {
	OCR2AbstractCaller     // Read-only binding to the contract
	OCR2AbstractTransactor // Write-only binding to the contract
	OCR2AbstractFilterer   // Log filterer for contract events
}

// OCR2AbstractCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR2AbstractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AbstractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR2AbstractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AbstractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR2AbstractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AbstractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR2AbstractSession struct {
	Contract     *OCR2Abstract     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OCR2AbstractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR2AbstractCallerSession struct {
	Contract *OCR2AbstractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// OCR2AbstractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR2AbstractTransactorSession struct {
	Contract     *OCR2AbstractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// OCR2AbstractRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR2AbstractRaw struct {
	Contract *OCR2Abstract // Generic contract binding to access the raw methods on
}

// OCR2AbstractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR2AbstractCallerRaw struct {
	Contract *OCR2AbstractCaller // Generic read-only contract binding to access the raw methods on
}

// OCR2AbstractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR2AbstractTransactorRaw struct {
	Contract *OCR2AbstractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR2Abstract creates a new instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2Abstract(address common.Address, backend bind.ContractBackend) (*OCR2Abstract, error) {
	contract, err := bindOCR2Abstract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2Abstract{OCR2AbstractCaller: OCR2AbstractCaller{contract: contract}, OCR2AbstractTransactor: OCR2AbstractTransactor{contract: contract}, OCR2AbstractFilterer: OCR2AbstractFilterer{contract: contract}}, nil
}

// NewOCR2AbstractCaller creates a new read-only instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2AbstractCaller(address common.Address, caller bind.ContractCaller) (*OCR2AbstractCaller, error) {
	contract, err := bindOCR2Abstract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractCaller{contract: contract}, nil
}

// NewOCR2AbstractTransactor creates a new write-only instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2AbstractTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2AbstractTransactor, error) {
	contract, err := bindOCR2Abstract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractTransactor{contract: contract}, nil
}

// NewOCR2AbstractFilterer creates a new log filterer instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2AbstractFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2AbstractFilterer, error) {
	contract, err := bindOCR2Abstract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractFilterer{contract: contract}, nil
}

// bindOCR2Abstract binds a generic wrapper to an already deployed contract.
func bindOCR2Abstract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR2AbstractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2Abstract *OCR2AbstractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Abstract.Contract.OCR2AbstractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2Abstract *OCR2AbstractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.OCR2AbstractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2Abstract *OCR2AbstractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.OCR2AbstractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2Abstract *OCR2AbstractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Abstract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2Abstract *OCR2AbstractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2Abstract *OCR2AbstractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.contract.Transact(opts, method, params...)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Abstract *OCR2AbstractCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	var out []interface{}
	err := _OCR2Abstract.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Abstract *OCR2AbstractSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDetails(&_OCR2Abstract.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Abstract *OCR2AbstractCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDetails(&_OCR2Abstract.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	var out []interface{}
	err := _OCR2Abstract.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(struct {
		ScanLogs     bool
		ConfigDigest [32]byte
		Epoch        uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDigestAndEpoch(&_OCR2Abstract.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractCallerSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDigestAndEpoch(&_OCR2Abstract.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Abstract *OCR2AbstractCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCR2Abstract.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Abstract *OCR2AbstractSession) TypeAndVersion() (string, error) {
	return _OCR2Abstract.Contract.TypeAndVersion(&_OCR2Abstract.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Abstract *OCR2AbstractCallerSession) TypeAndVersion() (string, error) {
	return _OCR2Abstract.Contract.TypeAndVersion(&_OCR2Abstract.CallOpts)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Abstract *OCR2AbstractTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Abstract.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Abstract *OCR2AbstractSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.SetConfig(&_OCR2Abstract.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Abstract *OCR2AbstractTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.SetConfig(&_OCR2Abstract.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Abstract *OCR2AbstractTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Abstract.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Abstract *OCR2AbstractSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.Transmit(&_OCR2Abstract.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Abstract *OCR2AbstractTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.Transmit(&_OCR2Abstract.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// OCR2AbstractConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the OCR2Abstract contract.
type OCR2AbstractConfigSetIterator struct {
	Event *OCR2AbstractConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AbstractConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AbstractConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AbstractConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AbstractConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AbstractConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AbstractConfigSet represents a ConfigSet event raised by the OCR2Abstract contract.
type OCR2AbstractConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Abstract *OCR2AbstractFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OCR2AbstractConfigSetIterator, error) {

	logs, sub, err := _OCR2Abstract.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractConfigSetIterator{contract: _OCR2Abstract.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Abstract *OCR2AbstractFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2AbstractConfigSet) (event.Subscription, error) {

	logs, sub, err := _OCR2Abstract.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AbstractConfigSet)
				if err := _OCR2Abstract.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConfigSet is a log parse operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Abstract *OCR2AbstractFilterer) ParseConfigSet(log types.Log) (*OCR2AbstractConfigSet, error) {
	event := new(OCR2AbstractConfigSet)
	if err := _OCR2Abstract.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AbstractTransmittedIterator is returned from FilterTransmitted and is used to iterate over the raw logs and unpacked data for Transmitted events raised by the OCR2Abstract contract.
type OCR2AbstractTransmittedIterator struct {
	Event *OCR2AbstractTransmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AbstractTransmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AbstractTransmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AbstractTransmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AbstractTransmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AbstractTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AbstractTransmitted represents a Transmitted event raised by the OCR2Abstract contract.
type OCR2AbstractTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransmitted is a free log retrieval operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractFilterer) FilterTransmitted(opts *bind.FilterOpts) (*OCR2AbstractTransmittedIterator, error) {

	logs, sub, err := _OCR2Abstract.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractTransmittedIterator{contract: _OCR2Abstract.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

// WatchTransmitted is a free log subscription operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OCR2AbstractTransmitted) (event.Subscription, error) {

	logs, sub, err := _OCR2Abstract.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AbstractTransmitted)
				if err := _OCR2Abstract.contract.UnpackLog(event, "Transmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransmitted is a log parse operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractFilterer) ParseTransmitted(log types.Log) (*OCR2AbstractTransmitted, error) {
	event := new(OCR2AbstractTransmitted)
	if err := _OCR2Abstract.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorMetaData contains all meta data concerning the OCR2Aggregator contract.
var OCR2AggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"minAnswer_\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"maxAnswer_\",\"type\":\"int192\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"billingAccessController\",\"type\":\"address\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"requesterAccessController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"description_\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"oldLinkToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"newLinkToken\",\"type\":\"address\"}],\"name\":\"LinkTokenSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"observationsTimestamp\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"RequesterAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"}],\"name\":\"RoundRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"previousValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousGasLimit\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"currentValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"currentGasLimit\",\"type\":\"uint32\"}],\"name\":\"ValidatorConfigSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBillingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequesterAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId_\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidatorConfig\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer_\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp_\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestNewRound\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"setLinkToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"requesterAccessController\",\"type\":\"address\"}],\"name\":\"setRequesterAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"newValidator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"newGasLimit\",\"type\":\"uint32\"}],\"name\":\"setValidatorConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b506040516200564938038062005649833981016040819052620000349162000488565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be816200015d565b5050601180546001600160a01b0319166001600160a01b038a169081179091556040519091506000907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a908290a3620001178462000208565b60ff821660c05260106200012c82826200066a565b50620001388362000281565b62000145600080620002fc565b50505050601791820b608052900b60a0525062000736565b336001600160a01b03821603620001b75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6012546001600160a01b0390811690821681146200027d57601280546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d4891291015b60405180910390a15b5050565b6200028b620003e3565b600f546001600160a01b0390811690821681146200027d57600f80546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634910162000274565b62000306620003e3565b60408051808201909152600e546001600160a01b03808216808452600160a01b90920463ffffffff16602084015284161415806200035457508163ffffffff16816020015163ffffffff1614155b15620003de576040805180820182526001600160a01b0385811680835263ffffffff8681166020948501819052600e80546001600160c01b0319168417600160a01b830217905586518786015187519316835294820152909392909116917fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541910160405180910390a35b505050565b6000546001600160a01b031633146200043f5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b6001600160a01b03811681146200045757600080fd5b50565b8051601781900b81146200046d57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b600080600080600080600060e0888a031215620004a457600080fd5b8751620004b18162000441565b96506020620004c28982016200045a565b9650620004d260408a016200045a565b95506060890151620004e48162000441565b60808a0151909550620004f78162000441565b60a08a015190945060ff811681146200050f57600080fd5b60c08a01519093506001600160401b03808211156200052d57600080fd5b818b0191508b601f8301126200054257600080fd5b81518181111562000557576200055762000472565b604051601f8201601f19908116603f0116810190838211818310171562000582576200058262000472565b816040528281528e868487010111156200059b57600080fd5b600093505b82841015620005bf5784840186015181850187015292850192620005a0565b600086848301015280965050505050505092959891949750929550565b600181811c90821680620005f157607f821691505b6020821081036200061257634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620003de57600081815260208120601f850160051c81016020861015620006415750805b601f850160051c820191505b8181101562000662578281556001016200064d565b505050505050565b81516001600160401b0381111562000686576200068662000472565b6200069e81620006978454620005dc565b8462000618565b602080601f831160018114620006d65760008415620006bd5750858301515b600019600386901b1c1916600185901b17855562000662565b600085815260208120601f198616915b828110156200070757888601518255948401946001909101908401620006e6565b5085821015620007265787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805160a05160c051614ec7620007826000396000610449015260008181610509015281816121ec0152613568015260008181610345015281816121c4015261353b0152614ec76000f3fe608060405234801561001057600080fd5b50600436106102d35760003560e01c80639bd2c0b111610186578063d09dc339116100e3578063e76d516811610097578063f2fde38b11610071578063f2fde38b14610848578063fbffd2c11461085b578063feaf968c1461086e57600080fd5b8063e76d516814610811578063eb45716314610822578063eb5dcd6c1461083557600080fd5b8063e3d0e712116100c8578063e3d0e7121461078c578063e4902f821461079f578063e5fe4577146107c757600080fd5b8063d09dc33914610773578063daffc4b51461077b57600080fd5b8063b1dc65a41161013a578063b633620c1161011f578063b633620c1461073c578063c10753291461074f578063c4c92b371461076257600080fd5b8063b1dc65a414610716578063b5ab58dc1461072957600080fd5b80639e3ceeab1161016b5780639e3ceeab146106bf578063afcb95d7146106d2578063b121e1471461070357600080fd5b80639bd2c0b11461065a5780639c849b30146106ac57600080fd5b8063668a0f02116102345780638205bf6a116101e85780638da5cb5b116101cd5780638da5cb5b146105c857806398e5b12a146105ed5780639a6fc8f51461061057600080fd5b80638205bf6a1461056b5780638ac28d5a146105b557600080fd5b80637284e416116102195780637284e4161461052b57806379ba50971461053357806381ff70481461053b57600080fd5b8063668a0f02146104ec57806370da2f671461050457600080fd5b80634fb174701161028b57806354fd4d501161027057806354fd4d50146104bc578063643dc105146104c4578063666cab8d146104d757600080fd5b80634fb174701461047d57806350d25bcd1461049257600080fd5b806322adbc78116102bc57806322adbc7814610340578063299372681461037a578063313ce5671461044457600080fd5b80630eafb25b146102d8578063181f5a77146102fe575b600080fd5b6102eb6102e636600461432f565b6108fd565b6040519081526020015b60405180910390f35b60408051808201909152601481527f4f43523241676772656761746f7220312e302e3000000000000000000000000060208201525b6040516102f59190614392565b6103677f000000000000000000000000000000000000000000000000000000000000000081565b60405160179190910b81526020016102f5565b610408600b546a0100000000000000000000810463ffffffff908116926e010000000000000000000000000000830482169272010000000000000000000000000000000000008104831692760100000000000000000000000000000000000000000000820416917a01000000000000000000000000000000000000000000000000000090910462ffffff1690565b6040805163ffffffff9687168152948616602086015292851692840192909252909216606082015262ffffff909116608082015260a0016102f5565b61046b7f000000000000000000000000000000000000000000000000000000000000000081565b60405160ff90911681526020016102f5565b61049061048b3660046143a5565b610a1e565b005b600b546601000000000000900463ffffffff166000908152600c602052604090205460170b6102eb565b6102eb600681565b6104906104d23660046143f0565b610c9e565b6104df610f59565b6040516102f591906144ad565b600b546601000000000000900463ffffffff166102eb565b6103677f000000000000000000000000000000000000000000000000000000000000000081565b610333610fbb565b610490611044565b600d54600a546040805163ffffffff808516825264010000000090940490931660208401528201526060016102f5565b600b546601000000000000900463ffffffff9081166000908152600c60205260409020547c01000000000000000000000000000000000000000000000000000000009004166102eb565b6104906105c336600461432f565b61110d565b6000546001600160a01b03165b6040516001600160a01b0390911681526020016102f5565b6105f5611182565b60405169ffffffffffffffffffff90911681526020016102f5565b61062361061e3660046144c0565b6112f9565b6040805169ffffffffffffffffffff968716815260208101959095528401929092526060830152909116608082015260a0016102f5565b604080518082018252600e546001600160a01b0381168083527401000000000000000000000000000000000000000090910463ffffffff166020928301819052835191825291810191909152016102f5565b6104906106ba366004614538565b6113b7565b6104906106cd36600461432f565b6115ad565b600a54600b546040805160008152602081019390935261010090910460081c63ffffffff16908201526060016102f5565b61049061071136600461432f565b611644565b6104906107243660046145a4565b611738565b6102eb610737366004614689565b611c83565b6102eb61074a366004614689565b611cb5565b61049061075d3660046146a2565b611d07565b6012546001600160a01b03166105d5565b6102eb611fdb565b600f546001600160a01b03166105d5565b61049061079a366004614864565b612084565b6107b26107ad36600461432f565b612900565b60405163ffffffff90911681526020016102f5565b6107cf6129be565b6040805195865263ffffffff909416602086015260ff9092169284019290925260179190910b606083015267ffffffffffffffff16608082015260a0016102f5565b6011546001600160a01b03166105d5565b610490610830366004614931565b612a83565b6104906108433660046143a5565b612ba0565b61049061085636600461432f565b612cf0565b61049061086936600461432f565b612d01565b610623600b546601000000000000900463ffffffff9081166000818152600c602090815260409182902082516060810184529054601781900b8083527801000000000000000000000000000000000000000000000000820487169383018490527c0100000000000000000000000000000000000000000000000000000000909104909516920182905291938490565b6001600160a01b03811660009081526002602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046bffffffffffffffffffffffff1691810191909152906109645750600092915050565b600b5460208201516000917201000000000000000000000000000000000000900463ffffffff169060069060ff16601f81106109a2576109a261495f565b600881049190910154600b546109d8926007166004026101000a90910463ffffffff9081169166010000000000009004166149bd565b63ffffffff166109e891906149e1565b6109f690633b9aca006149e1565b905081604001516bffffffffffffffffffffffff1681610a1691906149f8565b949350505050565b610a26612d12565b6011546001600160a01b03908116908316819003610a4357505050565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b038416906370a0823190602401602060405180830381865afa158015610aa0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ac49190614a0b565b50610acd612d6e565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000906001600160a01b038316906370a0823190602401602060405180830381865afa158015610b2d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b519190614a0b565b6040517fa9059cbb0000000000000000000000000000000000000000000000000000000081526001600160a01b038581166004830152602482018390529192509083169063a9059cbb906044016020604051808303816000875af1158015610bbd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610be19190614a24565b610c325760405162461bcd60e51b815260206004820152601f60248201527f7472616e736665722072656d61696e696e672066756e6473206661696c65640060448201526064015b60405180910390fd5b601180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0386811691821790925560405190918416907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a90600090a350505b5050565b6012546001600160a01b0316610cbc6000546001600160a01b031690565b6001600160a01b0316336001600160a01b03161480610d6157506040517f6b14daf80000000000000000000000000000000000000000000000000000000081526001600160a01b03821690636b14daf890610d209033906000903690600401614a46565b602060405180830381865afa158015610d3d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d619190614a24565b610dad5760405162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c6044820152606401610c29565b610db5612d6e565b600b80547fffffffffffffffffffffffffffff0000000000000000ffffffffffffffffffff166a010000000000000000000063ffffffff8981169182027fffffffffffffffffffffffffffff00000000ffffffffffffffffffffffffffff16929092176e010000000000000000000000000000898416908102919091177fffffffffffff0000000000000000ffffffffffffffffffffffffffffffffffff1672010000000000000000000000000000000000008985169081027fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff1691909117760100000000000000000000000000000000000000000000948916948502177fffffff000000ffffffffffffffffffffffffffffffffffffffffffffffffffff167a01000000000000000000000000000000000000000000000000000062ffffff89169081029190911790955560408051938452602084019290925290820152606081019190915260808101919091527f0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f9060a00160405180910390a1505050505050565b60606005805480602002602001604051908101604052809291908181526020018280548015610fb157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610f93575b5050505050905090565b606060108054610fca90614a85565b80601f0160208091040260200160405190810160405280929190818152602001828054610ff690614a85565b8015610fb15780601f1061101857610100808354040283529160200191610fb1565b820191906000526020600020905b81548152906001019060200180831161102657509395945050505050565b6001546001600160a01b0316331461109e5760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610c29565b60008054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6001600160a01b038181166000908152601360205260409020541633146111765760405162461bcd60e51b815260206004820152601760248201527f4f6e6c792070617965652063616e2077697468647261770000000000000000006044820152606401610c29565b61117f81613152565b50565b600080546001600160a01b03163314806112265750600f546040517f6b14daf80000000000000000000000000000000000000000000000000000000081526001600160a01b0390911690636b14daf8906111e59033906000903690600401614a46565b602060405180830381865afa158015611202573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112269190614a24565b6112725760405162461bcd60e51b815260206004820152601d60248201527f4f6e6c79206f776e6572267265717565737465722063616e2063616c6c0000006044820152606401610c29565b600b54600a546040805191825263ffffffff6101008404600881901c8216602085015260ff811684840152915164ffffffffff9092169366010000000000009004169133917f41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c9181900360600190a26112ec816001614ad8565b63ffffffff169250505090565b60008080808063ffffffff69ffffffffffffffffffff8716111561132b575060009350839250829150819050806113ae565b5050505063ffffffff8281166000908152600c602090815260409182902082516060810184529054601781900b8083527801000000000000000000000000000000000000000000000000820486169383018490527c0100000000000000000000000000000000000000000000000000000000909104909416920182905284935090835b91939590929450565b6113bf612d12565b82811461140e5760405162461bcd60e51b815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a656044820152606401610c29565b60005b838110156115a657600085858381811061142d5761142d61495f565b9050602002016020810190611442919061432f565b905060008484848181106114585761145861495f565b905060200201602081019061146d919061432f565b6001600160a01b0380841660009081526013602052604090205491925016801580806114aa5750826001600160a01b0316826001600160a01b0316145b6114f65760405162461bcd60e51b815260206004820152601160248201527f706179656520616c7265616479207365740000000000000000000000000000006044820152606401610c29565b6001600160a01b03848116600090815260136020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000168583169081179091559083161461158f57826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b50505050808061159e90614af5565b915050611411565b5050505050565b6115b5612d12565b600f546001600160a01b039081169082168114610c9a57600f80547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0384811691821790925560408051928416835260208301919091527f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae63491015b60405180910390a15050565b6001600160a01b038181166000908152601460205260409020541633146116ad5760405162461bcd60e51b815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e20616363657074006044820152606401610c29565b6001600160a01b0381811660008181526013602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556014909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60005a604080516101008082018352600b5460ff8116835290810464ffffffffff90811660208085018290526601000000000000840463ffffffff908116968601969096526a01000000000000000000008404861660608601526e01000000000000000000000000000084048616608086015272010000000000000000000000000000000000008404861660a0860152760100000000000000000000000000000000000000000000840490951660c08501527a01000000000000000000000000000000000000000000000000000090920462ffffff1660e08401529394509092918c01359182161161186c5760405162461bcd60e51b815260206004820152600c60248201527f7374616c65207265706f727400000000000000000000000000000000000000006044820152606401610c29565b3360009081526002602052604090205460ff166118cb5760405162461bcd60e51b815260206004820152601860248201527f756e617574686f72697a6564207472616e736d697474657200000000000000006044820152606401610c29565b600a548b351461191d5760405162461bcd60e51b815260206004820152601560248201527f636f6e666967446967657374206d69736d6174636800000000000000000000006044820152606401610c29565b61192b8a8a8a8a8a8a61339e565b8151611938906001614b2d565b60ff1687146119895760405162461bcd60e51b815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e6174757265730000000000006044820152606401610c29565b8685146119d85760405162461bcd60e51b815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e00006044820152606401610c29565b60008a8a6040516119ea929190614b46565b604051908190038120611a01918e90602001614b56565b60408051601f19818403018152828252805160209182012083830190925260008084529083018190529092509060005b8a811015611ba75760006001858a8460208110611a5057611a5061495f565b611a5d91901a601b614b2d565b8f8f86818110611a6f57611a6f61495f565b905060200201358e8e87818110611a8857611a8861495f565b9050602002013560405160008152602001604052604051611ac5949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015611ae7573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526003602090815290849020838501909452925460ff8082161515808552610100909204169383019390935290955092509050611b805760405162461bcd60e51b815260206004820152600f60248201527f7369676e6174757265206572726f7200000000000000000000000000000000006044820152606401610c29565b826020015160080260ff166001901b84019350508080611b9f90614af5565b915050611a31565b5081827e010101010101010101010101010101010101010101010101010101010101011614611c185760405162461bcd60e51b815260206004820152601060248201527f6475706c6963617465207369676e6572000000000000000000000000000000006044820152606401610c29565b5060009150611c679050838d836020020135848e8e8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061343b92505050565b9050611c7583828633613938565b505050505050505050505050565b600063ffffffff821115611c9957506000919050565b5063ffffffff166000908152600c602052604090205460170b90565b600063ffffffff821115611ccb57506000919050565b5063ffffffff9081166000908152600c60205260409020547c010000000000000000000000000000000000000000000000000000000090041690565b6000546001600160a01b0316331480611daa57506012546040517f6b14daf80000000000000000000000000000000000000000000000000000000081526001600160a01b0390911690636b14daf890611d699033906000903690600401614a46565b602060405180830381865afa158015611d86573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611daa9190614a24565b611df65760405162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c6044820152606401610c29565b6000611e00613a83565b6011546040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201529192506000916001600160a01b03909116906370a0823190602401602060405180830381865afa158015611e67573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e8b9190614a0b565b905081811015611edd5760405162461bcd60e51b815260206004820152601460248201527f696e73756666696369656e742062616c616e63650000000000000000000000006044820152606401610c29565b6011546001600160a01b031663a9059cbb85611f02611efc8686614b6c565b87613c64565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af1158015611f65573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f899190614a24565b611fd55760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610c29565b50505050565b6011546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009182916001600160a01b03909116906370a0823190602401602060405180830381865afa158015612041573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906120659190614a0b565b90506000612071613a83565b905061207d8183614b7f565b9250505090565b61208c612d12565b601f865111156120de5760405162461bcd60e51b815260206004820152601060248201527f746f6f206d616e79206f7261636c6573000000000000000000000000000000006044820152606401610c29565b845186511461212f5760405162461bcd60e51b815260206004820152601660248201527f6f7261636c65206c656e677468206d69736d61746368000000000000000000006044820152606401610c29565b855161213c856003614b9f565b60ff161061218c5760405162461bcd60e51b815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610c29565b6121988460ff16613c7e565b604080517f010000000000000000000000000000000000000000000000000000000000000060208201527f0000000000000000000000000000000000000000000000000000000000000000821b60218201527f000000000000000000000000000000000000000000000000000000000000000090911b6039820152605101604051602081830303815290604052805190602001208380519060200120146122815760405162461bcd60e51b815260206004820152601560248201527f696e76616c6964206f6e636861696e436f6e66696700000000000000000000006044820152606401610c29565b6040805160c0810182528781526020810187905260ff8616918101919091526060810184905267ffffffffffffffff8316608082015260a08101829052600b80547fffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000ff1690556122ee612d6e565b60045460005b818110156123cd576000600482815481106123115761231161495f565b6000918252602082200154600580546001600160a01b039092169350908490811061233e5761233e61495f565b60009182526020808320909101546001600160a01b03948516835260038252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016905594168252600290529190912080547fffffffffffffffffffffffffffffffffffff000000000000000000000000000016905550806123c581614af5565b9150506122f4565b506123da600460006141d7565b6123e6600560006141d7565b60005b8251518110156126f357600360008460000151838151811061240d5761240d61495f565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff16156124815760405162461bcd60e51b815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610c29565b604080518082019091526001815260ff8216602082015283518051600391600091859081106124b2576124b261495f565b6020908102919091018101516001600160a01b0316825281810192909252604001600090812083518154948401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169015157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff161761010060ff909516949094029390931790925584015180516002929190849081106125575761255761495f565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff16156125cb5760405162461bcd60e51b815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610c29565b60405180606001604052806001151581526020018260ff16815260200160006bffffffffffffffffffffffff1681525060026000856020015184815181106126155761261561495f565b6020908102919091018101516001600160a01b03168252818101929092526040908101600020835181549385015194909201516bffffffffffffffffffffffff1662010000027fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff931515939093167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093179190911792909216179055806126eb81614af5565b9150506123e9565b508151805161270a916004916020909101906141f5565b5060208083015180516127219260059201906141f5565b506040820151600b80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600d80547fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff811664010000000063ffffffff4381168202928317855590830481169360019390926000926127b3928692908216911617614ad8565b92506101000a81548163ffffffff021916908363ffffffff1602179055506128124630600d60009054906101000a900463ffffffff1663ffffffff1686600001518760200151886040015189606001518a608001518b60a00151613cce565b600a819055600d5484516020860151604080880151606089015160808a015160a08b015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e059861287b988b98919763ffffffff909116969195909490939092909190614bbb565b60405180910390a1600b546601000000000000900463ffffffff1660005b8451518110156128f35781600682601f81106128b7576128b761495f565b600891828204019190066004026101000a81548163ffffffff021916908363ffffffff16021790555080806128eb90614af5565b915050612899565b5050505050505050505050565b6001600160a01b03811660009081526002602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046bffffffffffffffffffffffff1691810191909152906129675750600092915050565b6006816020015160ff16601f81106129815761298161495f565b600881049190910154600b546129b7926007166004026101000a90910463ffffffff9081169166010000000000009004166149bd565b9392505050565b600080808080333214612a135760405162461bcd60e51b815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f410000000000000000000000006044820152606401610c29565b5050600a54600b5463ffffffff6601000000000000820481166000908152600c60205260409020549296610100909204600881901c8216965064ffffffffff169450601783900b93507c010000000000000000000000000000000000000000000000000000000090920490911690565b612a8b612d12565b60408051808201909152600e546001600160a01b038082168084527401000000000000000000000000000000000000000090920463ffffffff1660208401528416141580612ae957508163ffffffff16816020015163ffffffff1614155b15612b9b576040805180820182526001600160a01b0385811680835263ffffffff8681166020948501819052600e80547fffffffffffffffff00000000000000000000000000000000000000000000000016841774010000000000000000000000000000000000000000830217905586518786015187519316835294820152909392909116917fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541910160405180910390a35b505050565b6001600160a01b03828116600090815260136020526040902054163314612c095760405162461bcd60e51b815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e207570646174650000006044820152606401610c29565b6001600160a01b0381163303612c615760405162461bcd60e51b815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610c29565b6001600160a01b03808316600090815260146020526040902080548383167fffffffffffffffffffffffff000000000000000000000000000000000000000082168117909255909116908114612b9b576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a4505050565b612cf8612d12565b61117f81613d5c565b612d09612d12565b61117f81613e1d565b6000546001600160a01b03163314612d6c5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610c29565b565b601154600b54604080516103e08101918290526001600160a01b0390931692660100000000000090920463ffffffff1691600091600690601f908285855b82829054906101000a900463ffffffff1663ffffffff1681526020019060040190602082600301049283019260010382029150808411612dac5790505050505050905060006005805480602002602001604051908101604052809291908181526020018280548015612e4757602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612e29575b5050505050905060005b815181101561314457600060026000848481518110612e7257612e7261495f565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160029054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff169050600060026000858581518110612ede57612ede61495f565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160026101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008483601f8110612f4b57612f4b61495f565b6020020151600b5490870363ffffffff90811692507201000000000000000000000000000000000000909104168102633b9aca00028201801561313957600060136000878781518110612fa057612fa061495f565b6020908102919091018101516001600160a01b0390811683529082019290925260409081016000205490517fa9059cbb00000000000000000000000000000000000000000000000000000000815290821660048201819052602482018590529250908a169063a9059cbb906044016020604051808303816000875af115801561302d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130519190614a24565b61309d5760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610c29565b878786601f81106130b0576130b061495f565b602002019063ffffffff16908163ffffffff1681525050886001600160a01b0316816001600160a01b03168787815181106130ed576130ed61495f565b60200260200101516001600160a01b03167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c8560405161312f91815260200190565b60405180910390a4505b505050600101612e51565b506115a6600683601f614272565b6001600160a01b0381166000908152600260209081526040918290208251606081018452905460ff80821615158084526101008304909116938301939093526201000090046bffffffffffffffffffffffff16928101929092526131b4575050565b60006131bf836108fd565b90508015612b9b576001600160a01b03838116600090815260136020526040908190205460115491517fa9059cbb000000000000000000000000000000000000000000000000000000008152908316600482018190526024820185905292919091169063a9059cbb906044016020604051808303816000875af115801561324a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061326e9190614a24565b6132ba5760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610c29565b600b60000160069054906101000a900463ffffffff166006846020015160ff16601f81106132ea576132ea61495f565b6008810491909101805460079092166004026101000a63ffffffff8181021990931693909216919091029190911790556001600160a01b0384811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff169055601154915186815291841693851692917fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c910160405180910390a450505050565b60006133ab8260206149e1565b6133b68560206149e1565b6133c2886101446149f8565b6133cc91906149f8565b6133d691906149f8565b6133e19060006149f8565b90503681146134325760405162461bcd60e51b815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d6174636800000000000000006044820152606401610c29565b50505050505050565b60008061344783613ea4565b9050601f816040015151111561349f5760405162461bcd60e51b815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e647300006044820152606401610c29565b604081015151865160ff16106134f75760405162461bcd60e51b815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e00006044820152606401610c29565b64ffffffffff84166020870152604081015180516000919061351b90600290614c51565b8151811061352b5761352b61495f565b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b1315801561359157507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b6135dd5760405162461bcd60e51b815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e676500006044820152606401610c29565b604087018051906135ed82614c8c565b63ffffffff1663ffffffff168152505060405180606001604052808260170b8152602001836000015163ffffffff1681526020014263ffffffff16815250600c6000896040015163ffffffff1663ffffffff16815260200190815260200160002060008201518160000160006101000a81548177ffffffffffffffffffffffffffffffffffffffffffffffff021916908360170b77ffffffffffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160186101000a81548163ffffffff021916908363ffffffff160217905550604082015181600001601c6101000a81548163ffffffff021916908363ffffffff16021790555090505086600b60008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548164ffffffffff021916908364ffffffffff16021790555060408201518160000160066101000a81548163ffffffff021916908363ffffffff160217905550606082015181600001600a6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600e6101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160000160126101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160000160166101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600001601a6101000a81548162ffffff021916908362ffffff160217905550905050866040015163ffffffff167fc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a823385600001518660400151876020015188606001518d8d604051613881989796959493929190614caf565b60405180910390a26040808801518351915163ffffffff9283168152600092909116907f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac602719060200160405180910390a3866040015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f4260405161391191815260200190565b60405180910390a361392a87604001518260170b613f49565b506060015195945050505050565b60008360170b12611fd557600061396a633b9aca003a04866080015163ffffffff16876060015163ffffffff16614099565b90506010360260005a905060006139938663ffffffff1685858b60e0015162ffffff16866140bf565b90506000670de0b6b3a764000077ffffffffffffffffffffffffffffffffffffffffffffffff891683026001600160a01b03881660009081526002602052604090205460c08c01519290910492506201000090046bffffffffffffffffffffffff9081169163ffffffff16633b9aca000282840101908116821115613a1e5750505050505050611fd5565b6001600160a01b038816600090815260026020526040902080546bffffffffffffffffffffffff90921662010000027fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff90921691909117905550505050505050505050565b6000806005805480602002602001604051908101604052809291908181526020018280548015613adc57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311613abe575b50508351600b54604080516103e08101918290529697509195660100000000000090910463ffffffff169450600093509150600690601f908285855b82829054906101000a900463ffffffff1663ffffffff1681526020019060040190602082600301049283019260010382029150808411613b185790505050505050905060005b83811015613bab578181601f8110613b7857613b7861495f565b6020020151613b8790846149bd565b613b979063ffffffff16876149f8565b955080613ba381614af5565b915050613b5e565b50600b54613bd9907201000000000000000000000000000000000000900463ffffffff16633b9aca006149e1565b613be390866149e1565b945060005b83811015613c5c5760026000868381518110613c0657613c0661495f565b6020908102919091018101516001600160a01b0316825281019190915260400160002054613c48906201000090046bffffffffffffffffffffffff16876149f8565b955080613c5481614af5565b915050613be8565b505050505090565b600081831015613c75575081613c78565b50805b92915050565b8060001061117f5760405162461bcd60e51b815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610c29565b6000808a8a8a8a8a8a8a8a8a604051602001613cf299989796959493929190614d57565b60408051601f1981840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150505b9998505050505050505050565b336001600160a01b03821603613db45760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610c29565b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6012546001600160a01b039081169082168114610c9a57601280547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0384811691821790925560408051928416835260208301919091527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129101611638565b613ed86040518060800160405280600063ffffffff1681526020016060815260200160608152602001600060170b81525090565b6000806060600085806020019051810190613ef39190614df1565b92965090945092509050613f078683614123565b81516040805160208082019690965281519082018252918252805160808101825263ffffffff969096168652938501529183015260170b606082015292915050565b60408051808201909152600e546001600160a01b0381168083527401000000000000000000000000000000000000000090910463ffffffff166020830152613f9057505050565b6000613f9d6001856149bd565b63ffffffff8181166000818152600c6020908152604091829020549087015187519251602481019490945260179190910b6044840181905289851660648501526084840189905294955061404d93169160a40160408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fbeed9b510000000000000000000000000000000000000000000000000000000017905261419b565b6115a65760405162461bcd60e51b815260206004820152601060248201527f696e73756666696369656e7420676173000000000000000000000000000000006044820152606401610c29565b600083838110156140ac57600285850304015b6140b68184613c64565b95945050505050565b6000818610156141115760405162461bcd60e51b815260206004820181905260248201527f6c6566744761732063616e6e6f742065786365656420696e697469616c4761736044820152606401610c29565b50633b9aca0094039190910101020290565b60008151602061413391906149e1565b61413e9060a06149f8565b6141499060006149f8565b905080835114612b9b5760405162461bcd60e51b815260206004820152601660248201527f7265706f7274206c656e677468206d69736d61746368000000000000000000006044820152606401610c29565b60005a61138881106141cf57611388810390508460408204820311156141cf576000808451602086016000888af150600191505b509392505050565b508054600082559060005260206000209081019061117f9190614305565b828054828255906000526020600020908101928215614262579160200282015b8281111561426257825182547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03909116178255602090920191600190910190614215565b5061426e929150614305565b5090565b6004830191839082156142625791602002820160005b838211156142cc57835183826101000a81548163ffffffff021916908363ffffffff1602179055509260200192600401602081600301049283019260010302614288565b80156142fc5782816101000a81549063ffffffff02191690556004016020816003010492830192600103026142cc565b505061426e9291505b5b8082111561426e5760008155600101614306565b6001600160a01b038116811461117f57600080fd5b60006020828403121561434157600080fd5b81356129b78161431a565b6000815180845260005b8181101561437257602081850181015186830182015201614356565b506000602082860101526020601f19601f83011685010191505092915050565b6020815260006129b7602083018461434c565b600080604083850312156143b857600080fd5b82356143c38161431a565b915060208301356143d38161431a565b809150509250929050565b63ffffffff8116811461117f57600080fd5b600080600080600060a0868803121561440857600080fd5b8535614413816143de565b94506020860135614423816143de565b93506040860135614433816143de565b92506060860135614443816143de565b9150608086013562ffffff8116811461445b57600080fd5b809150509295509295909350565b600081518084526020808501945080840160005b838110156144a25781516001600160a01b03168752958201959082019060010161447d565b509495945050505050565b6020815260006129b76020830184614469565b6000602082840312156144d257600080fd5b813569ffffffffffffffffffff811681146129b757600080fd5b60008083601f8401126144fe57600080fd5b50813567ffffffffffffffff81111561451657600080fd5b6020830191508360208260051b850101111561453157600080fd5b9250929050565b6000806000806040858703121561454e57600080fd5b843567ffffffffffffffff8082111561456657600080fd5b614572888389016144ec565b9096509450602087013591508082111561458b57600080fd5b50614598878288016144ec565b95989497509550505050565b60008060008060008060008060e0898b0312156145c057600080fd5b606089018a8111156145d157600080fd5b8998503567ffffffffffffffff808211156145eb57600080fd5b818b0191508b601f8301126145ff57600080fd5b81358181111561460e57600080fd5b8c602082850101111561462057600080fd5b6020830199508098505060808b013591508082111561463e57600080fd5b61464a8c838d016144ec565b909750955060a08b013591508082111561466357600080fd5b506146708b828c016144ec565b999c989b50969995989497949560c00135949350505050565b60006020828403121561469b57600080fd5b5035919050565b600080604083850312156146b557600080fd5b82356146c08161431a565b946020939093013593505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715614726576147266146ce565b604052919050565b600067ffffffffffffffff821115614748576147486146ce565b5060051b60200190565b600082601f83011261476357600080fd5b813560206147786147738361472e565b6146fd565b82815260059290921b8401810191818101908684111561479757600080fd5b8286015b848110156147bb5780356147ae8161431a565b835291830191830161479b565b509695505050505050565b803560ff811681146147d757600080fd5b919050565b600082601f8301126147ed57600080fd5b813567ffffffffffffffff811115614807576148076146ce565b61481a6020601f19601f840116016146fd565b81815284602083860101111561482f57600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff811681146147d757600080fd5b60008060008060008060c0878903121561487d57600080fd5b863567ffffffffffffffff8082111561489557600080fd5b6148a18a838b01614752565b975060208901359150808211156148b757600080fd5b6148c38a838b01614752565b96506148d160408a016147c6565b955060608901359150808211156148e757600080fd5b6148f38a838b016147dc565b945061490160808a0161484c565b935060a089013591508082111561491757600080fd5b5061492489828a016147dc565b9150509295509295509295565b6000806040838503121561494457600080fd5b823561494f8161431a565b915060208301356143d3816143de565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b63ffffffff8281168282160390808211156149da576149da61498e565b5092915050565b8082028115828204841417613c7857613c7861498e565b80820180821115613c7857613c7861498e565b600060208284031215614a1d57600080fd5b5051919050565b600060208284031215614a3657600080fd5b815180151581146129b757600080fd5b6001600160a01b038416815260406020820152816040820152818360608301376000818301606090810191909152601f909201601f1916010192915050565b600181811c90821680614a9957607f821691505b602082108103614ad2577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b63ffffffff8181168382160190808211156149da576149da61498e565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614b2657614b2661498e565b5060010190565b60ff8181168382160190811115613c7857613c7861498e565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b81810381811115613c7857613c7861498e565b81810360008312801583831316838312821617156149da576149da61498e565b60ff81811683821602908116908181146149da576149da61498e565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614beb8184018a614469565b90508281036080840152614bff8189614469565b905060ff871660a084015282810360c0840152614c1c818761434c565b905067ffffffffffffffff851660e0840152828103610100840152614c41818561434c565b9c9b505050505050505050505050565b600082614c87577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b600063ffffffff808316818103614ca557614ca561498e565b6001019392505050565b600061010080830160178c810b855260206001600160a01b038d168187015263ffffffff8c1660408701528360608701528293508a5180845261012087019450818c01935060005b81811015614d15578451840b86529482019493820193600101614cf7565b50505050508281036080840152614d2c818861434c565b915050614d3e60a083018660170b9052565b8360c0830152613d4f60e083018464ffffffffff169052565b60006101208b83526001600160a01b038b16602084015267ffffffffffffffff808b166040850152816060850152614d918285018b614469565b91508382036080850152614da5828a614469565b915060ff881660a085015283820360c0850152614dc2828861434c565b90861660e08501528381036101008501529050614c41818561434c565b8051601781900b81146147d757600080fd5b60008060008060808587031215614e0757600080fd5b8451614e12816143de565b809450506020808601519350604086015167ffffffffffffffff811115614e3857600080fd5b8601601f81018813614e4957600080fd5b8051614e576147738261472e565b81815260059190911b8201830190838101908a831115614e7657600080fd5b928401925b82841015614e9b57614e8c84614ddf565b82529284019290840190614e7b565b8096505050505050614eaf60608601614ddf565b90509295919450925056fea164736f6c6343000813000a",
}

// OCR2AggregatorABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR2AggregatorMetaData.ABI instead.
var OCR2AggregatorABI = OCR2AggregatorMetaData.ABI

// OCR2AggregatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OCR2AggregatorMetaData.Bin instead.
var OCR2AggregatorBin = OCR2AggregatorMetaData.Bin

// DeployOCR2Aggregator deploys a new Ethereum contract, binding an instance of OCR2Aggregator to it.
func DeployOCR2Aggregator(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, minAnswer_ *big.Int, maxAnswer_ *big.Int, billingAccessController common.Address, requesterAccessController common.Address, decimals_ uint8, description_ string) (common.Address, *types.Transaction, *OCR2Aggregator, error) {
	parsed, err := OCR2AggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR2AggregatorBin), backend, link, minAnswer_, maxAnswer_, billingAccessController, requesterAccessController, decimals_, description_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR2Aggregator{OCR2AggregatorCaller: OCR2AggregatorCaller{contract: contract}, OCR2AggregatorTransactor: OCR2AggregatorTransactor{contract: contract}, OCR2AggregatorFilterer: OCR2AggregatorFilterer{contract: contract}}, nil
}

// OCR2Aggregator is an auto generated Go binding around an Ethereum contract.
type OCR2Aggregator struct {
	OCR2AggregatorCaller     // Read-only binding to the contract
	OCR2AggregatorTransactor // Write-only binding to the contract
	OCR2AggregatorFilterer   // Log filterer for contract events
}

// OCR2AggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR2AggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR2AggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR2AggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR2AggregatorSession struct {
	Contract     *OCR2Aggregator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OCR2AggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR2AggregatorCallerSession struct {
	Contract *OCR2AggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// OCR2AggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR2AggregatorTransactorSession struct {
	Contract     *OCR2AggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// OCR2AggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR2AggregatorRaw struct {
	Contract *OCR2Aggregator // Generic contract binding to access the raw methods on
}

// OCR2AggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR2AggregatorCallerRaw struct {
	Contract *OCR2AggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// OCR2AggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR2AggregatorTransactorRaw struct {
	Contract *OCR2AggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR2Aggregator creates a new instance of OCR2Aggregator, bound to a specific deployed contract.
func NewOCR2Aggregator(address common.Address, backend bind.ContractBackend) (*OCR2Aggregator, error) {
	contract, err := bindOCR2Aggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2Aggregator{OCR2AggregatorCaller: OCR2AggregatorCaller{contract: contract}, OCR2AggregatorTransactor: OCR2AggregatorTransactor{contract: contract}, OCR2AggregatorFilterer: OCR2AggregatorFilterer{contract: contract}}, nil
}

// NewOCR2AggregatorCaller creates a new read-only instance of OCR2Aggregator, bound to a specific deployed contract.
func NewOCR2AggregatorCaller(address common.Address, caller bind.ContractCaller) (*OCR2AggregatorCaller, error) {
	contract, err := bindOCR2Aggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorCaller{contract: contract}, nil
}

// NewOCR2AggregatorTransactor creates a new write-only instance of OCR2Aggregator, bound to a specific deployed contract.
func NewOCR2AggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2AggregatorTransactor, error) {
	contract, err := bindOCR2Aggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorTransactor{contract: contract}, nil
}

// NewOCR2AggregatorFilterer creates a new log filterer instance of OCR2Aggregator, bound to a specific deployed contract.
func NewOCR2AggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2AggregatorFilterer, error) {
	contract, err := bindOCR2Aggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorFilterer{contract: contract}, nil
}

// bindOCR2Aggregator binds a generic wrapper to an already deployed contract.
func bindOCR2Aggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR2AggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2Aggregator *OCR2AggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Aggregator.Contract.OCR2AggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2Aggregator *OCR2AggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.OCR2AggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2Aggregator *OCR2AggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.OCR2AggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2Aggregator *OCR2AggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Aggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2Aggregator *OCR2AggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2Aggregator *OCR2AggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.contract.Transact(opts, method, params...)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_OCR2Aggregator *OCR2AggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_OCR2Aggregator *OCR2AggregatorSession) Decimals() (uint8, error) {
	return _OCR2Aggregator.Contract.Decimals(&_OCR2Aggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) Decimals() (uint8, error) {
	return _OCR2Aggregator.Contract.Decimals(&_OCR2Aggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_OCR2Aggregator *OCR2AggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_OCR2Aggregator *OCR2AggregatorSession) Description() (string, error) {
	return _OCR2Aggregator.Contract.Description(&_OCR2Aggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) Description() (string, error) {
	return _OCR2Aggregator.Contract.Description(&_OCR2Aggregator.CallOpts)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_OCR2Aggregator *OCR2AggregatorCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_OCR2Aggregator *OCR2AggregatorSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _OCR2Aggregator.Contract.GetAnswer(&_OCR2Aggregator.CallOpts, roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _OCR2Aggregator.Contract.GetAnswer(&_OCR2Aggregator.CallOpts, roundId)
}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_OCR2Aggregator *OCR2AggregatorCaller) GetBilling(opts *bind.CallOpts) (struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
}, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getBilling")

	outstruct := new(struct {
		MaximumGasPriceGwei       uint32
		ReasonableGasPriceGwei    uint32
		ObservationPaymentGjuels  uint32
		TransmissionPaymentGjuels uint32
		AccountingGas             *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaximumGasPriceGwei = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.ReasonableGasPriceGwei = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ObservationPaymentGjuels = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.TransmissionPaymentGjuels = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.AccountingGas = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_OCR2Aggregator *OCR2AggregatorSession) GetBilling() (struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
}, error) {
	return _OCR2Aggregator.Contract.GetBilling(&_OCR2Aggregator.CallOpts)
}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetBilling() (struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
}, error) {
	return _OCR2Aggregator.Contract.GetBilling(&_OCR2Aggregator.CallOpts)
}

// GetBillingAccessController is a free data retrieval call binding the contract method 0xc4c92b37.
//
// Solidity: function getBillingAccessController() view returns(address)
func (_OCR2Aggregator *OCR2AggregatorCaller) GetBillingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getBillingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBillingAccessController is a free data retrieval call binding the contract method 0xc4c92b37.
//
// Solidity: function getBillingAccessController() view returns(address)
func (_OCR2Aggregator *OCR2AggregatorSession) GetBillingAccessController() (common.Address, error) {
	return _OCR2Aggregator.Contract.GetBillingAccessController(&_OCR2Aggregator.CallOpts)
}

// GetBillingAccessController is a free data retrieval call binding the contract method 0xc4c92b37.
//
// Solidity: function getBillingAccessController() view returns(address)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetBillingAccessController() (common.Address, error) {
	return _OCR2Aggregator.Contract.GetBillingAccessController(&_OCR2Aggregator.CallOpts)
}

// GetLinkToken is a free data retrieval call binding the contract method 0xe76d5168.
//
// Solidity: function getLinkToken() view returns(address linkToken)
func (_OCR2Aggregator *OCR2AggregatorCaller) GetLinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getLinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetLinkToken is a free data retrieval call binding the contract method 0xe76d5168.
//
// Solidity: function getLinkToken() view returns(address linkToken)
func (_OCR2Aggregator *OCR2AggregatorSession) GetLinkToken() (common.Address, error) {
	return _OCR2Aggregator.Contract.GetLinkToken(&_OCR2Aggregator.CallOpts)
}

// GetLinkToken is a free data retrieval call binding the contract method 0xe76d5168.
//
// Solidity: function getLinkToken() view returns(address linkToken)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetLinkToken() (common.Address, error) {
	return _OCR2Aggregator.Contract.GetLinkToken(&_OCR2Aggregator.CallOpts)
}

// GetRequesterAccessController is a free data retrieval call binding the contract method 0xdaffc4b5.
//
// Solidity: function getRequesterAccessController() view returns(address)
func (_OCR2Aggregator *OCR2AggregatorCaller) GetRequesterAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getRequesterAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRequesterAccessController is a free data retrieval call binding the contract method 0xdaffc4b5.
//
// Solidity: function getRequesterAccessController() view returns(address)
func (_OCR2Aggregator *OCR2AggregatorSession) GetRequesterAccessController() (common.Address, error) {
	return _OCR2Aggregator.Contract.GetRequesterAccessController(&_OCR2Aggregator.CallOpts)
}

// GetRequesterAccessController is a free data retrieval call binding the contract method 0xdaffc4b5.
//
// Solidity: function getRequesterAccessController() view returns(address)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetRequesterAccessController() (common.Address, error) {
	return _OCR2Aggregator.Contract.GetRequesterAccessController(&_OCR2Aggregator.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80 roundId_, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_OCR2Aggregator *OCR2AggregatorCaller) GetRoundData(opts *bind.CallOpts, roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getRoundData", roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
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

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80 roundId_, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_OCR2Aggregator *OCR2AggregatorSession) GetRoundData(roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _OCR2Aggregator.Contract.GetRoundData(&_OCR2Aggregator.CallOpts, roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80 roundId_, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetRoundData(roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _OCR2Aggregator.Contract.GetRoundData(&_OCR2Aggregator.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _OCR2Aggregator.Contract.GetTimestamp(&_OCR2Aggregator.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _OCR2Aggregator.Contract.GetTimestamp(&_OCR2Aggregator.CallOpts, roundId)
}

// GetTransmitters is a free data retrieval call binding the contract method 0x666cab8d.
//
// Solidity: function getTransmitters() view returns(address[])
func (_OCR2Aggregator *OCR2AggregatorCaller) GetTransmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getTransmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetTransmitters is a free data retrieval call binding the contract method 0x666cab8d.
//
// Solidity: function getTransmitters() view returns(address[])
func (_OCR2Aggregator *OCR2AggregatorSession) GetTransmitters() ([]common.Address, error) {
	return _OCR2Aggregator.Contract.GetTransmitters(&_OCR2Aggregator.CallOpts)
}

// GetTransmitters is a free data retrieval call binding the contract method 0x666cab8d.
//
// Solidity: function getTransmitters() view returns(address[])
func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetTransmitters() ([]common.Address, error) {
	return _OCR2Aggregator.Contract.GetTransmitters(&_OCR2Aggregator.CallOpts)
}

// GetValidatorConfig is a free data retrieval call binding the contract method 0x9bd2c0b1.
//
// Solidity: function getValidatorConfig() view returns(address validator, uint32 gasLimit)
func (_OCR2Aggregator *OCR2AggregatorCaller) GetValidatorConfig(opts *bind.CallOpts) (struct {
	Validator common.Address
	GasLimit  uint32
}, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getValidatorConfig")

	outstruct := new(struct {
		Validator common.Address
		GasLimit  uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.GasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return *outstruct, err

}

// GetValidatorConfig is a free data retrieval call binding the contract method 0x9bd2c0b1.
//
// Solidity: function getValidatorConfig() view returns(address validator, uint32 gasLimit)
func (_OCR2Aggregator *OCR2AggregatorSession) GetValidatorConfig() (struct {
	Validator common.Address
	GasLimit  uint32
}, error) {
	return _OCR2Aggregator.Contract.GetValidatorConfig(&_OCR2Aggregator.CallOpts)
}

// GetValidatorConfig is a free data retrieval call binding the contract method 0x9bd2c0b1.
//
// Solidity: function getValidatorConfig() view returns(address validator, uint32 gasLimit)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetValidatorConfig() (struct {
	Validator common.Address
	GasLimit  uint32
}, error) {
	return _OCR2Aggregator.Contract.GetValidatorConfig(&_OCR2Aggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_OCR2Aggregator *OCR2AggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_OCR2Aggregator *OCR2AggregatorSession) LatestAnswer() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LatestAnswer(&_OCR2Aggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LatestAnswer(&_OCR2Aggregator.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Aggregator *OCR2AggregatorCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Aggregator *OCR2AggregatorSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2Aggregator.Contract.LatestConfigDetails(&_OCR2Aggregator.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2Aggregator.Contract.LatestConfigDetails(&_OCR2Aggregator.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Aggregator *OCR2AggregatorCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(struct {
		ScanLogs     bool
		ConfigDigest [32]byte
		Epoch        uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Aggregator *OCR2AggregatorSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2Aggregator.Contract.LatestConfigDigestAndEpoch(&_OCR2Aggregator.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2Aggregator.Contract.LatestConfigDigestAndEpoch(&_OCR2Aggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorSession) LatestRound() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LatestRound(&_OCR2Aggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LatestRound(&_OCR2Aggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_OCR2Aggregator *OCR2AggregatorCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
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

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_OCR2Aggregator *OCR2AggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _OCR2Aggregator.Contract.LatestRoundData(&_OCR2Aggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _OCR2Aggregator.Contract.LatestRoundData(&_OCR2Aggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LatestTimestamp(&_OCR2Aggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LatestTimestamp(&_OCR2Aggregator.CallOpts)
}

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes32 configDigest, uint32 epoch, uint8 round, int192 latestAnswer_, uint64 latestTimestamp_)
func (_OCR2Aggregator *OCR2AggregatorCaller) LatestTransmissionDetails(opts *bind.CallOpts) (struct {
	ConfigDigest    [32]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestTransmissionDetails")

	outstruct := new(struct {
		ConfigDigest    [32]byte
		Epoch           uint32
		Round           uint8
		LatestAnswer    *big.Int
		LatestTimestamp uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigDigest = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.Round = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.LatestAnswer = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.LatestTimestamp = *abi.ConvertType(out[4], new(uint64)).(*uint64)

	return *outstruct, err

}

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes32 configDigest, uint32 epoch, uint8 round, int192 latestAnswer_, uint64 latestTimestamp_)
func (_OCR2Aggregator *OCR2AggregatorSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [32]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _OCR2Aggregator.Contract.LatestTransmissionDetails(&_OCR2Aggregator.CallOpts)
}

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes32 configDigest, uint32 epoch, uint8 round, int192 latestAnswer_, uint64 latestTimestamp_)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [32]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _OCR2Aggregator.Contract.LatestTransmissionDetails(&_OCR2Aggregator.CallOpts)
}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_OCR2Aggregator *OCR2AggregatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_OCR2Aggregator *OCR2AggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LinkAvailableForPayment(&_OCR2Aggregator.CallOpts)
}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LinkAvailableForPayment(&_OCR2Aggregator.CallOpts)
}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int192)
func (_OCR2Aggregator *OCR2AggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "maxAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int192)
func (_OCR2Aggregator *OCR2AggregatorSession) MaxAnswer() (*big.Int, error) {
	return _OCR2Aggregator.Contract.MaxAnswer(&_OCR2Aggregator.CallOpts)
}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int192)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _OCR2Aggregator.Contract.MaxAnswer(&_OCR2Aggregator.CallOpts)
}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int192)
func (_OCR2Aggregator *OCR2AggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "minAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int192)
func (_OCR2Aggregator *OCR2AggregatorSession) MinAnswer() (*big.Int, error) {
	return _OCR2Aggregator.Contract.MinAnswer(&_OCR2Aggregator.CallOpts)
}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int192)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _OCR2Aggregator.Contract.MinAnswer(&_OCR2Aggregator.CallOpts)
}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address transmitterAddress) view returns(uint32)
func (_OCR2Aggregator *OCR2AggregatorCaller) OracleObservationCount(opts *bind.CallOpts, transmitterAddress common.Address) (uint32, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "oracleObservationCount", transmitterAddress)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address transmitterAddress) view returns(uint32)
func (_OCR2Aggregator *OCR2AggregatorSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _OCR2Aggregator.Contract.OracleObservationCount(&_OCR2Aggregator.CallOpts, transmitterAddress)
}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address transmitterAddress) view returns(uint32)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _OCR2Aggregator.Contract.OracleObservationCount(&_OCR2Aggregator.CallOpts, transmitterAddress)
}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address transmitterAddress) view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorCaller) OwedPayment(opts *bind.CallOpts, transmitterAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "owedPayment", transmitterAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address transmitterAddress) view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _OCR2Aggregator.Contract.OwedPayment(&_OCR2Aggregator.CallOpts, transmitterAddress)
}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address transmitterAddress) view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _OCR2Aggregator.Contract.OwedPayment(&_OCR2Aggregator.CallOpts, transmitterAddress)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OCR2Aggregator *OCR2AggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OCR2Aggregator *OCR2AggregatorSession) Owner() (common.Address, error) {
	return _OCR2Aggregator.Contract.Owner(&_OCR2Aggregator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) Owner() (common.Address, error) {
	return _OCR2Aggregator.Contract.Owner(&_OCR2Aggregator.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Aggregator *OCR2AggregatorCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Aggregator *OCR2AggregatorSession) TypeAndVersion() (string, error) {
	return _OCR2Aggregator.Contract.TypeAndVersion(&_OCR2Aggregator.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) TypeAndVersion() (string, error) {
	return _OCR2Aggregator.Contract.TypeAndVersion(&_OCR2Aggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorSession) Version() (*big.Int, error) {
	return _OCR2Aggregator.Contract.Version(&_OCR2Aggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_OCR2Aggregator *OCR2AggregatorCallerSession) Version() (*big.Int, error) {
	return _OCR2Aggregator.Contract.Version(&_OCR2Aggregator.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OCR2Aggregator *OCR2AggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.AcceptOwnership(&_OCR2Aggregator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.AcceptOwnership(&_OCR2Aggregator.TransactOpts)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "acceptPayeeship", transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_OCR2Aggregator *OCR2AggregatorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.AcceptPayeeship(&_OCR2Aggregator.TransactOpts, transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.AcceptPayeeship(&_OCR2Aggregator.TransactOpts, transmitter)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_OCR2Aggregator *OCR2AggregatorTransactor) RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "requestNewRound")
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_OCR2Aggregator *OCR2AggregatorSession) RequestNewRound() (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.RequestNewRound(&_OCR2Aggregator.TransactOpts)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) RequestNewRound() (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.RequestNewRound(&_OCR2Aggregator.TransactOpts)
}

// SetBilling is a paid mutator transaction binding the contract method 0x643dc105.
//
// Solidity: function setBilling(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) SetBilling(opts *bind.TransactOpts, maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setBilling", maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

// SetBilling is a paid mutator transaction binding the contract method 0x643dc105.
//
// Solidity: function setBilling(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas) returns()
func (_OCR2Aggregator *OCR2AggregatorSession) SetBilling(maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetBilling(&_OCR2Aggregator.TransactOpts, maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

// SetBilling is a paid mutator transaction binding the contract method 0x643dc105.
//
// Solidity: function setBilling(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetBilling(maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetBilling(&_OCR2Aggregator.TransactOpts, maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setBillingAccessController", _billingAccessController)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_OCR2Aggregator *OCR2AggregatorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetBillingAccessController(&_OCR2Aggregator.TransactOpts, _billingAccessController)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetBillingAccessController(&_OCR2Aggregator.TransactOpts, _billingAccessController)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Aggregator *OCR2AggregatorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetConfig(&_OCR2Aggregator.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetConfig(&_OCR2Aggregator.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetLinkToken is a paid mutator transaction binding the contract method 0x4fb17470.
//
// Solidity: function setLinkToken(address linkToken, address recipient) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) SetLinkToken(opts *bind.TransactOpts, linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setLinkToken", linkToken, recipient)
}

// SetLinkToken is a paid mutator transaction binding the contract method 0x4fb17470.
//
// Solidity: function setLinkToken(address linkToken, address recipient) returns()
func (_OCR2Aggregator *OCR2AggregatorSession) SetLinkToken(linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetLinkToken(&_OCR2Aggregator.TransactOpts, linkToken, recipient)
}

// SetLinkToken is a paid mutator transaction binding the contract method 0x4fb17470.
//
// Solidity: function setLinkToken(address linkToken, address recipient) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetLinkToken(linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetLinkToken(&_OCR2Aggregator.TransactOpts, linkToken, recipient)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] transmitters, address[] payees) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) SetPayees(opts *bind.TransactOpts, transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setPayees", transmitters, payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] transmitters, address[] payees) returns()
func (_OCR2Aggregator *OCR2AggregatorSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetPayees(&_OCR2Aggregator.TransactOpts, transmitters, payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] transmitters, address[] payees) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetPayees(&_OCR2Aggregator.TransactOpts, transmitters, payees)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address requesterAccessController) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) SetRequesterAccessController(opts *bind.TransactOpts, requesterAccessController common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setRequesterAccessController", requesterAccessController)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address requesterAccessController) returns()
func (_OCR2Aggregator *OCR2AggregatorSession) SetRequesterAccessController(requesterAccessController common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetRequesterAccessController(&_OCR2Aggregator.TransactOpts, requesterAccessController)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address requesterAccessController) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetRequesterAccessController(requesterAccessController common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetRequesterAccessController(&_OCR2Aggregator.TransactOpts, requesterAccessController)
}

// SetValidatorConfig is a paid mutator transaction binding the contract method 0xeb457163.
//
// Solidity: function setValidatorConfig(address newValidator, uint32 newGasLimit) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) SetValidatorConfig(opts *bind.TransactOpts, newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setValidatorConfig", newValidator, newGasLimit)
}

// SetValidatorConfig is a paid mutator transaction binding the contract method 0xeb457163.
//
// Solidity: function setValidatorConfig(address newValidator, uint32 newGasLimit) returns()
func (_OCR2Aggregator *OCR2AggregatorSession) SetValidatorConfig(newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetValidatorConfig(&_OCR2Aggregator.TransactOpts, newValidator, newGasLimit)
}

// SetValidatorConfig is a paid mutator transaction binding the contract method 0xeb457163.
//
// Solidity: function setValidatorConfig(address newValidator, uint32 newGasLimit) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetValidatorConfig(newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetValidatorConfig(&_OCR2Aggregator.TransactOpts, newValidator, newGasLimit)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OCR2Aggregator *OCR2AggregatorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.TransferOwnership(&_OCR2Aggregator.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.TransferOwnership(&_OCR2Aggregator.TransactOpts, to)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_OCR2Aggregator *OCR2AggregatorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.TransferPayeeship(&_OCR2Aggregator.TransactOpts, transmitter, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.TransferPayeeship(&_OCR2Aggregator.TransactOpts, transmitter, proposed)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Aggregator *OCR2AggregatorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.Transmit(&_OCR2Aggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.Transmit(&_OCR2Aggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address recipient, uint256 amount) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "withdrawFunds", recipient, amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address recipient, uint256 amount) returns()
func (_OCR2Aggregator *OCR2AggregatorSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.WithdrawFunds(&_OCR2Aggregator.TransactOpts, recipient, amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address recipient, uint256 amount) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.WithdrawFunds(&_OCR2Aggregator.TransactOpts, recipient, amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address transmitter) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "withdrawPayment", transmitter)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address transmitter) returns()
func (_OCR2Aggregator *OCR2AggregatorSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.WithdrawPayment(&_OCR2Aggregator.TransactOpts, transmitter)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address transmitter) returns()
func (_OCR2Aggregator *OCR2AggregatorTransactorSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.WithdrawPayment(&_OCR2Aggregator.TransactOpts, transmitter)
}

// OCR2AggregatorAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the OCR2Aggregator contract.
type OCR2AggregatorAnswerUpdatedIterator struct {
	Event *OCR2AggregatorAnswerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorAnswerUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorAnswerUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorAnswerUpdated represents a AnswerUpdated event raised by the OCR2Aggregator contract.
type OCR2AggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*OCR2AggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorAnswerUpdatedIterator{contract: _OCR2Aggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorAnswerUpdated)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseAnswerUpdated(log types.Log) (*OCR2AggregatorAnswerUpdated, error) {
	event := new(OCR2AggregatorAnswerUpdated)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorBillingAccessControllerSetIterator is returned from FilterBillingAccessControllerSet and is used to iterate over the raw logs and unpacked data for BillingAccessControllerSet events raised by the OCR2Aggregator contract.
type OCR2AggregatorBillingAccessControllerSetIterator struct {
	Event *OCR2AggregatorBillingAccessControllerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorBillingAccessControllerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorBillingAccessControllerSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorBillingAccessControllerSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorBillingAccessControllerSet represents a BillingAccessControllerSet event raised by the OCR2Aggregator contract.
type OCR2AggregatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBillingAccessControllerSet is a free log retrieval operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*OCR2AggregatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorBillingAccessControllerSetIterator{contract: _OCR2Aggregator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}

// WatchBillingAccessControllerSet is a free log subscription operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorBillingAccessControllerSet)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBillingAccessControllerSet is a log parse operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*OCR2AggregatorBillingAccessControllerSet, error) {
	event := new(OCR2AggregatorBillingAccessControllerSet)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorBillingSetIterator is returned from FilterBillingSet and is used to iterate over the raw logs and unpacked data for BillingSet events raised by the OCR2Aggregator contract.
type OCR2AggregatorBillingSetIterator struct {
	Event *OCR2AggregatorBillingSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorBillingSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorBillingSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorBillingSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorBillingSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorBillingSet represents a BillingSet event raised by the OCR2Aggregator contract.
type OCR2AggregatorBillingSet struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterBillingSet is a free log retrieval operation binding the contract event 0x0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f.
//
// Solidity: event BillingSet(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*OCR2AggregatorBillingSetIterator, error) {

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorBillingSetIterator{contract: _OCR2Aggregator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}

// WatchBillingSet is a free log subscription operation binding the contract event 0x0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f.
//
// Solidity: event BillingSet(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorBillingSet) (event.Subscription, error) {

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorBillingSet)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBillingSet is a log parse operation binding the contract event 0x0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f.
//
// Solidity: event BillingSet(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseBillingSet(log types.Log) (*OCR2AggregatorBillingSet, error) {
	event := new(OCR2AggregatorBillingSet)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the OCR2Aggregator contract.
type OCR2AggregatorConfigSetIterator struct {
	Event *OCR2AggregatorConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorConfigSet represents a ConfigSet event raised by the OCR2Aggregator contract.
type OCR2AggregatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OCR2AggregatorConfigSetIterator, error) {

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorConfigSetIterator{contract: _OCR2Aggregator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorConfigSet)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConfigSet is a log parse operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseConfigSet(log types.Log) (*OCR2AggregatorConfigSet, error) {
	event := new(OCR2AggregatorConfigSet)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorLinkTokenSetIterator is returned from FilterLinkTokenSet and is used to iterate over the raw logs and unpacked data for LinkTokenSet events raised by the OCR2Aggregator contract.
type OCR2AggregatorLinkTokenSetIterator struct {
	Event *OCR2AggregatorLinkTokenSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorLinkTokenSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorLinkTokenSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorLinkTokenSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorLinkTokenSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorLinkTokenSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorLinkTokenSet represents a LinkTokenSet event raised by the OCR2Aggregator contract.
type OCR2AggregatorLinkTokenSet struct {
	OldLinkToken common.Address
	NewLinkToken common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLinkTokenSet is a free log retrieval operation binding the contract event 0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a.
//
// Solidity: event LinkTokenSet(address indexed oldLinkToken, address indexed newLinkToken)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterLinkTokenSet(opts *bind.FilterOpts, oldLinkToken []common.Address, newLinkToken []common.Address) (*OCR2AggregatorLinkTokenSetIterator, error) {

	var oldLinkTokenRule []interface{}
	for _, oldLinkTokenItem := range oldLinkToken {
		oldLinkTokenRule = append(oldLinkTokenRule, oldLinkTokenItem)
	}
	var newLinkTokenRule []interface{}
	for _, newLinkTokenItem := range newLinkToken {
		newLinkTokenRule = append(newLinkTokenRule, newLinkTokenItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "LinkTokenSet", oldLinkTokenRule, newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorLinkTokenSetIterator{contract: _OCR2Aggregator.contract, event: "LinkTokenSet", logs: logs, sub: sub}, nil
}

// WatchLinkTokenSet is a free log subscription operation binding the contract event 0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a.
//
// Solidity: event LinkTokenSet(address indexed oldLinkToken, address indexed newLinkToken)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchLinkTokenSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorLinkTokenSet, oldLinkToken []common.Address, newLinkToken []common.Address) (event.Subscription, error) {

	var oldLinkTokenRule []interface{}
	for _, oldLinkTokenItem := range oldLinkToken {
		oldLinkTokenRule = append(oldLinkTokenRule, oldLinkTokenItem)
	}
	var newLinkTokenRule []interface{}
	for _, newLinkTokenItem := range newLinkToken {
		newLinkTokenRule = append(newLinkTokenRule, newLinkTokenItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "LinkTokenSet", oldLinkTokenRule, newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorLinkTokenSet)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLinkTokenSet is a log parse operation binding the contract event 0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a.
//
// Solidity: event LinkTokenSet(address indexed oldLinkToken, address indexed newLinkToken)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseLinkTokenSet(log types.Log) (*OCR2AggregatorLinkTokenSet, error) {
	event := new(OCR2AggregatorLinkTokenSet)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the OCR2Aggregator contract.
type OCR2AggregatorNewRoundIterator struct {
	Event *OCR2AggregatorNewRound // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorNewRound)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorNewRound)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorNewRound represents a NewRound event raised by the OCR2Aggregator contract.
type OCR2AggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*OCR2AggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorNewRoundIterator{contract: _OCR2Aggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorNewRound)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseNewRound(log types.Log) (*OCR2AggregatorNewRound, error) {
	event := new(OCR2AggregatorNewRound)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorNewTransmissionIterator is returned from FilterNewTransmission and is used to iterate over the raw logs and unpacked data for NewTransmission events raised by the OCR2Aggregator contract.
type OCR2AggregatorNewTransmissionIterator struct {
	Event *OCR2AggregatorNewTransmission // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorNewTransmissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorNewTransmission)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorNewTransmission)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorNewTransmissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorNewTransmission represents a NewTransmission event raised by the OCR2Aggregator contract.
type OCR2AggregatorNewTransmission struct {
	AggregatorRoundId     uint32
	Answer                *big.Int
	Transmitter           common.Address
	ObservationsTimestamp uint32
	Observations          []*big.Int
	Observers             []byte
	JuelsPerFeeCoin       *big.Int
	ConfigDigest          [32]byte
	EpochAndRound         *big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterNewTransmission is a free log retrieval operation binding the contract event 0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, uint32 observationsTimestamp, int192[] observations, bytes observers, int192 juelsPerFeeCoin, bytes32 configDigest, uint40 epochAndRound)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*OCR2AggregatorNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorNewTransmissionIterator{contract: _OCR2Aggregator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

// WatchNewTransmission is a free log subscription operation binding the contract event 0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, uint32 observationsTimestamp, int192[] observations, bytes observers, int192 juelsPerFeeCoin, bytes32 configDigest, uint40 epochAndRound)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorNewTransmission)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewTransmission is a log parse operation binding the contract event 0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, uint32 observationsTimestamp, int192[] observations, bytes observers, int192 juelsPerFeeCoin, bytes32 configDigest, uint40 epochAndRound)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseNewTransmission(log types.Log) (*OCR2AggregatorNewTransmission, error) {
	event := new(OCR2AggregatorNewTransmission)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorOraclePaidIterator is returned from FilterOraclePaid and is used to iterate over the raw logs and unpacked data for OraclePaid events raised by the OCR2Aggregator contract.
type OCR2AggregatorOraclePaidIterator struct {
	Event *OCR2AggregatorOraclePaid // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorOraclePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorOraclePaid)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorOraclePaid)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorOraclePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorOraclePaid represents a OraclePaid event raised by the OCR2Aggregator contract.
type OCR2AggregatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	LinkToken   common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOraclePaid is a free log retrieval operation binding the contract event 0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c.
//
// Solidity: event OraclePaid(address indexed transmitter, address indexed payee, uint256 amount, address indexed linkToken)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*OCR2AggregatorOraclePaidIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	var linkTokenRule []interface{}
	for _, linkTokenItem := range linkToken {
		linkTokenRule = append(linkTokenRule, linkTokenItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorOraclePaidIterator{contract: _OCR2Aggregator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}

// WatchOraclePaid is a free log subscription operation binding the contract event 0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c.
//
// Solidity: event OraclePaid(address indexed transmitter, address indexed payee, uint256 amount, address indexed linkToken)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	var linkTokenRule []interface{}
	for _, linkTokenItem := range linkToken {
		linkTokenRule = append(linkTokenRule, linkTokenItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorOraclePaid)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOraclePaid is a log parse operation binding the contract event 0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c.
//
// Solidity: event OraclePaid(address indexed transmitter, address indexed payee, uint256 amount, address indexed linkToken)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseOraclePaid(log types.Log) (*OCR2AggregatorOraclePaid, error) {
	event := new(OCR2AggregatorOraclePaid)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the OCR2Aggregator contract.
type OCR2AggregatorOwnershipTransferRequestedIterator struct {
	Event *OCR2AggregatorOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the OCR2Aggregator contract.
type OCR2AggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2AggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorOwnershipTransferRequestedIterator{contract: _OCR2Aggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorOwnershipTransferRequested)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*OCR2AggregatorOwnershipTransferRequested, error) {
	event := new(OCR2AggregatorOwnershipTransferRequested)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OCR2Aggregator contract.
type OCR2AggregatorOwnershipTransferredIterator struct {
	Event *OCR2AggregatorOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorOwnershipTransferred represents a OwnershipTransferred event raised by the OCR2Aggregator contract.
type OCR2AggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2AggregatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorOwnershipTransferredIterator{contract: _OCR2Aggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorOwnershipTransferred)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*OCR2AggregatorOwnershipTransferred, error) {
	event := new(OCR2AggregatorOwnershipTransferred)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorPayeeshipTransferRequestedIterator is returned from FilterPayeeshipTransferRequested and is used to iterate over the raw logs and unpacked data for PayeeshipTransferRequested events raised by the OCR2Aggregator contract.
type OCR2AggregatorPayeeshipTransferRequestedIterator struct {
	Event *OCR2AggregatorPayeeshipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorPayeeshipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorPayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorPayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorPayeeshipTransferRequested represents a PayeeshipTransferRequested event raised by the OCR2Aggregator contract.
type OCR2AggregatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferRequested is a free log retrieval operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*OCR2AggregatorPayeeshipTransferRequestedIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorPayeeshipTransferRequestedIterator{contract: _OCR2Aggregator.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferRequested is a free log subscription operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorPayeeshipTransferRequested)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePayeeshipTransferRequested is a log parse operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*OCR2AggregatorPayeeshipTransferRequested, error) {
	event := new(OCR2AggregatorPayeeshipTransferRequested)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorPayeeshipTransferredIterator is returned from FilterPayeeshipTransferred and is used to iterate over the raw logs and unpacked data for PayeeshipTransferred events raised by the OCR2Aggregator contract.
type OCR2AggregatorPayeeshipTransferredIterator struct {
	Event *OCR2AggregatorPayeeshipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorPayeeshipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorPayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorPayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorPayeeshipTransferred represents a PayeeshipTransferred event raised by the OCR2Aggregator contract.
type OCR2AggregatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferred is a free log retrieval operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*OCR2AggregatorPayeeshipTransferredIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorPayeeshipTransferredIterator{contract: _OCR2Aggregator.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferred is a free log subscription operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorPayeeshipTransferred)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePayeeshipTransferred is a log parse operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParsePayeeshipTransferred(log types.Log) (*OCR2AggregatorPayeeshipTransferred, error) {
	event := new(OCR2AggregatorPayeeshipTransferred)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorRequesterAccessControllerSetIterator is returned from FilterRequesterAccessControllerSet and is used to iterate over the raw logs and unpacked data for RequesterAccessControllerSet events raised by the OCR2Aggregator contract.
type OCR2AggregatorRequesterAccessControllerSetIterator struct {
	Event *OCR2AggregatorRequesterAccessControllerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorRequesterAccessControllerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorRequesterAccessControllerSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorRequesterAccessControllerSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorRequesterAccessControllerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorRequesterAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorRequesterAccessControllerSet represents a RequesterAccessControllerSet event raised by the OCR2Aggregator contract.
type OCR2AggregatorRequesterAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRequesterAccessControllerSet is a free log retrieval operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*OCR2AggregatorRequesterAccessControllerSetIterator, error) {

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorRequesterAccessControllerSetIterator{contract: _OCR2Aggregator.contract, event: "RequesterAccessControllerSet", logs: logs, sub: sub}, nil
}

// WatchRequesterAccessControllerSet is a free log subscription operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorRequesterAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorRequesterAccessControllerSet)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequesterAccessControllerSet is a log parse operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseRequesterAccessControllerSet(log types.Log) (*OCR2AggregatorRequesterAccessControllerSet, error) {
	event := new(OCR2AggregatorRequesterAccessControllerSet)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorRoundRequestedIterator is returned from FilterRoundRequested and is used to iterate over the raw logs and unpacked data for RoundRequested events raised by the OCR2Aggregator contract.
type OCR2AggregatorRoundRequestedIterator struct {
	Event *OCR2AggregatorRoundRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorRoundRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorRoundRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorRoundRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorRoundRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorRoundRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorRoundRequested represents a RoundRequested event raised by the OCR2Aggregator contract.
type OCR2AggregatorRoundRequested struct {
	Requester    common.Address
	ConfigDigest [32]byte
	Epoch        uint32
	Round        uint8
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRoundRequested is a free log retrieval operation binding the contract event 0x41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c.
//
// Solidity: event RoundRequested(address indexed requester, bytes32 configDigest, uint32 epoch, uint8 round)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*OCR2AggregatorRoundRequestedIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorRoundRequestedIterator{contract: _OCR2Aggregator.contract, event: "RoundRequested", logs: logs, sub: sub}, nil
}

// WatchRoundRequested is a free log subscription operation binding the contract event 0x41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c.
//
// Solidity: event RoundRequested(address indexed requester, bytes32 configDigest, uint32 epoch, uint8 round)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorRoundRequested, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorRoundRequested)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoundRequested is a log parse operation binding the contract event 0x41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c.
//
// Solidity: event RoundRequested(address indexed requester, bytes32 configDigest, uint32 epoch, uint8 round)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseRoundRequested(log types.Log) (*OCR2AggregatorRoundRequested, error) {
	event := new(OCR2AggregatorRoundRequested)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorTransmittedIterator is returned from FilterTransmitted and is used to iterate over the raw logs and unpacked data for Transmitted events raised by the OCR2Aggregator contract.
type OCR2AggregatorTransmittedIterator struct {
	Event *OCR2AggregatorTransmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorTransmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorTransmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorTransmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorTransmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorTransmitted represents a Transmitted event raised by the OCR2Aggregator contract.
type OCR2AggregatorTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransmitted is a free log retrieval operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterTransmitted(opts *bind.FilterOpts) (*OCR2AggregatorTransmittedIterator, error) {

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorTransmittedIterator{contract: _OCR2Aggregator.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

// WatchTransmitted is a free log subscription operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorTransmitted) (event.Subscription, error) {

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorTransmitted)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "Transmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransmitted is a log parse operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseTransmitted(log types.Log) (*OCR2AggregatorTransmitted, error) {
	event := new(OCR2AggregatorTransmitted)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AggregatorValidatorConfigSetIterator is returned from FilterValidatorConfigSet and is used to iterate over the raw logs and unpacked data for ValidatorConfigSet events raised by the OCR2Aggregator contract.
type OCR2AggregatorValidatorConfigSetIterator struct {
	Event *OCR2AggregatorValidatorConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AggregatorValidatorConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorValidatorConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AggregatorValidatorConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AggregatorValidatorConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AggregatorValidatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AggregatorValidatorConfigSet represents a ValidatorConfigSet event raised by the OCR2Aggregator contract.
type OCR2AggregatorValidatorConfigSet struct {
	PreviousValidator common.Address
	PreviousGasLimit  uint32
	CurrentValidator  common.Address
	CurrentGasLimit   uint32
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterValidatorConfigSet is a free log retrieval operation binding the contract event 0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541.
//
// Solidity: event ValidatorConfigSet(address indexed previousValidator, uint32 previousGasLimit, address indexed currentValidator, uint32 currentGasLimit)
func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterValidatorConfigSet(opts *bind.FilterOpts, previousValidator []common.Address, currentValidator []common.Address) (*OCR2AggregatorValidatorConfigSetIterator, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorValidatorConfigSetIterator{contract: _OCR2Aggregator.contract, event: "ValidatorConfigSet", logs: logs, sub: sub}, nil
}

// WatchValidatorConfigSet is a free log subscription operation binding the contract event 0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541.
//
// Solidity: event ValidatorConfigSet(address indexed previousValidator, uint32 previousGasLimit, address indexed currentValidator, uint32 currentGasLimit)
func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchValidatorConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorValidatorConfigSet, previousValidator []common.Address, currentValidator []common.Address) (event.Subscription, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AggregatorValidatorConfigSet)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorConfigSet is a log parse operation binding the contract event 0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541.
//
// Solidity: event ValidatorConfigSet(address indexed previousValidator, uint32 previousGasLimit, address indexed currentValidator, uint32 currentGasLimit)
func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseValidatorConfigSet(log types.Log) (*OCR2AggregatorValidatorConfigSet, error) {
	event := new(OCR2AggregatorValidatorConfigSet)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnableInterfaceMetaData contains all meta data concerning the OwnableInterface contract.
var OwnableInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// OwnableInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnableInterfaceMetaData.ABI instead.
var OwnableInterfaceABI = OwnableInterfaceMetaData.ABI

// OwnableInterface is an auto generated Go binding around an Ethereum contract.
type OwnableInterface struct {
	OwnableInterfaceCaller     // Read-only binding to the contract
	OwnableInterfaceTransactor // Write-only binding to the contract
	OwnableInterfaceFilterer   // Log filterer for contract events
}

// OwnableInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnableInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnableInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnableInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnableInterfaceSession struct {
	Contract     *OwnableInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnableInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnableInterfaceCallerSession struct {
	Contract *OwnableInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// OwnableInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnableInterfaceTransactorSession struct {
	Contract     *OwnableInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// OwnableInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnableInterfaceRaw struct {
	Contract *OwnableInterface // Generic contract binding to access the raw methods on
}

// OwnableInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnableInterfaceCallerRaw struct {
	Contract *OwnableInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnableInterfaceTransactorRaw struct {
	Contract *OwnableInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnableInterface creates a new instance of OwnableInterface, bound to a specific deployed contract.
func NewOwnableInterface(address common.Address, backend bind.ContractBackend) (*OwnableInterface, error) {
	contract, err := bindOwnableInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnableInterface{OwnableInterfaceCaller: OwnableInterfaceCaller{contract: contract}, OwnableInterfaceTransactor: OwnableInterfaceTransactor{contract: contract}, OwnableInterfaceFilterer: OwnableInterfaceFilterer{contract: contract}}, nil
}

// NewOwnableInterfaceCaller creates a new read-only instance of OwnableInterface, bound to a specific deployed contract.
func NewOwnableInterfaceCaller(address common.Address, caller bind.ContractCaller) (*OwnableInterfaceCaller, error) {
	contract, err := bindOwnableInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableInterfaceCaller{contract: contract}, nil
}

// NewOwnableInterfaceTransactor creates a new write-only instance of OwnableInterface, bound to a specific deployed contract.
func NewOwnableInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableInterfaceTransactor, error) {
	contract, err := bindOwnableInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableInterfaceTransactor{contract: contract}, nil
}

// NewOwnableInterfaceFilterer creates a new log filterer instance of OwnableInterface, bound to a specific deployed contract.
func NewOwnableInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableInterfaceFilterer, error) {
	contract, err := bindOwnableInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableInterfaceFilterer{contract: contract}, nil
}

// bindOwnableInterface binds a generic wrapper to an already deployed contract.
func bindOwnableInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OwnableInterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableInterface *OwnableInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableInterface.Contract.OwnableInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableInterface *OwnableInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableInterface.Contract.OwnableInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableInterface *OwnableInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableInterface.Contract.OwnableInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableInterface *OwnableInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableInterface *OwnableInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableInterface *OwnableInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableInterface.Contract.contract.Transact(opts, method, params...)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnableInterface *OwnableInterfaceTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableInterface.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnableInterface *OwnableInterfaceSession) AcceptOwnership() (*types.Transaction, error) {
	return _OwnableInterface.Contract.AcceptOwnership(&_OwnableInterface.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnableInterface *OwnableInterfaceTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OwnableInterface.Contract.AcceptOwnership(&_OwnableInterface.TransactOpts)
}

// Owner is a paid mutator transaction binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() returns(address)
func (_OwnableInterface *OwnableInterfaceTransactor) Owner(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableInterface.contract.Transact(opts, "owner")
}

// Owner is a paid mutator transaction binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() returns(address)
func (_OwnableInterface *OwnableInterfaceSession) Owner() (*types.Transaction, error) {
	return _OwnableInterface.Contract.Owner(&_OwnableInterface.TransactOpts)
}

// Owner is a paid mutator transaction binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() returns(address)
func (_OwnableInterface *OwnableInterfaceTransactorSession) Owner() (*types.Transaction, error) {
	return _OwnableInterface.Contract.Owner(&_OwnableInterface.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address recipient) returns()
func (_OwnableInterface *OwnableInterfaceTransactor) TransferOwnership(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _OwnableInterface.contract.Transact(opts, "transferOwnership", recipient)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address recipient) returns()
func (_OwnableInterface *OwnableInterfaceSession) TransferOwnership(recipient common.Address) (*types.Transaction, error) {
	return _OwnableInterface.Contract.TransferOwnership(&_OwnableInterface.TransactOpts, recipient)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address recipient) returns()
func (_OwnableInterface *OwnableInterfaceTransactorSession) TransferOwnership(recipient common.Address) (*types.Transaction, error) {
	return _OwnableInterface.Contract.TransferOwnership(&_OwnableInterface.TransactOpts, recipient)
}

// OwnerIsCreatorMetaData contains all meta data concerning the OwnerIsCreator contract.
var OwnerIsCreatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610367806101576000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806379ba5097146100465780638da5cb5b14610050578063f2fde38b1461007c575b600080fd5b61004e61008f565b005b6000546040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61004e61008a36600461031d565b610191565b60015473ffffffffffffffffffffffffffffffffffffffff163314610115576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6101996101a5565b6101a281610228565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610226576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161010c565b565b3373ffffffffffffffffffffffffffffffffffffffff8216036102a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161010c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561032f57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461035357600080fd5b939250505056fea164736f6c6343000813000a",
}

// OwnerIsCreatorABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnerIsCreatorMetaData.ABI instead.
var OwnerIsCreatorABI = OwnerIsCreatorMetaData.ABI

// OwnerIsCreatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OwnerIsCreatorMetaData.Bin instead.
var OwnerIsCreatorBin = OwnerIsCreatorMetaData.Bin

// DeployOwnerIsCreator deploys a new Ethereum contract, binding an instance of OwnerIsCreator to it.
func DeployOwnerIsCreator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OwnerIsCreator, error) {
	parsed, err := OwnerIsCreatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OwnerIsCreatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OwnerIsCreator{OwnerIsCreatorCaller: OwnerIsCreatorCaller{contract: contract}, OwnerIsCreatorTransactor: OwnerIsCreatorTransactor{contract: contract}, OwnerIsCreatorFilterer: OwnerIsCreatorFilterer{contract: contract}}, nil
}

// OwnerIsCreator is an auto generated Go binding around an Ethereum contract.
type OwnerIsCreator struct {
	OwnerIsCreatorCaller     // Read-only binding to the contract
	OwnerIsCreatorTransactor // Write-only binding to the contract
	OwnerIsCreatorFilterer   // Log filterer for contract events
}

// OwnerIsCreatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnerIsCreatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerIsCreatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnerIsCreatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerIsCreatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnerIsCreatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerIsCreatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnerIsCreatorSession struct {
	Contract     *OwnerIsCreator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnerIsCreatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnerIsCreatorCallerSession struct {
	Contract *OwnerIsCreatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// OwnerIsCreatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnerIsCreatorTransactorSession struct {
	Contract     *OwnerIsCreatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// OwnerIsCreatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnerIsCreatorRaw struct {
	Contract *OwnerIsCreator // Generic contract binding to access the raw methods on
}

// OwnerIsCreatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnerIsCreatorCallerRaw struct {
	Contract *OwnerIsCreatorCaller // Generic read-only contract binding to access the raw methods on
}

// OwnerIsCreatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnerIsCreatorTransactorRaw struct {
	Contract *OwnerIsCreatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnerIsCreator creates a new instance of OwnerIsCreator, bound to a specific deployed contract.
func NewOwnerIsCreator(address common.Address, backend bind.ContractBackend) (*OwnerIsCreator, error) {
	contract, err := bindOwnerIsCreator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreator{OwnerIsCreatorCaller: OwnerIsCreatorCaller{contract: contract}, OwnerIsCreatorTransactor: OwnerIsCreatorTransactor{contract: contract}, OwnerIsCreatorFilterer: OwnerIsCreatorFilterer{contract: contract}}, nil
}

// NewOwnerIsCreatorCaller creates a new read-only instance of OwnerIsCreator, bound to a specific deployed contract.
func NewOwnerIsCreatorCaller(address common.Address, caller bind.ContractCaller) (*OwnerIsCreatorCaller, error) {
	contract, err := bindOwnerIsCreator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorCaller{contract: contract}, nil
}

// NewOwnerIsCreatorTransactor creates a new write-only instance of OwnerIsCreator, bound to a specific deployed contract.
func NewOwnerIsCreatorTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnerIsCreatorTransactor, error) {
	contract, err := bindOwnerIsCreator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorTransactor{contract: contract}, nil
}

// NewOwnerIsCreatorFilterer creates a new log filterer instance of OwnerIsCreator, bound to a specific deployed contract.
func NewOwnerIsCreatorFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnerIsCreatorFilterer, error) {
	contract, err := bindOwnerIsCreator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorFilterer{contract: contract}, nil
}

// bindOwnerIsCreator binds a generic wrapper to an already deployed contract.
func bindOwnerIsCreator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OwnerIsCreatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnerIsCreator *OwnerIsCreatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnerIsCreator.Contract.OwnerIsCreatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnerIsCreator *OwnerIsCreatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.OwnerIsCreatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnerIsCreator *OwnerIsCreatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.OwnerIsCreatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnerIsCreator *OwnerIsCreatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnerIsCreator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnerIsCreator *OwnerIsCreatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnerIsCreator *OwnerIsCreatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnerIsCreator *OwnerIsCreatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OwnerIsCreator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnerIsCreator *OwnerIsCreatorSession) Owner() (common.Address, error) {
	return _OwnerIsCreator.Contract.Owner(&_OwnerIsCreator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnerIsCreator *OwnerIsCreatorCallerSession) Owner() (common.Address, error) {
	return _OwnerIsCreator.Contract.Owner(&_OwnerIsCreator.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnerIsCreator *OwnerIsCreatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnerIsCreator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnerIsCreator *OwnerIsCreatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.AcceptOwnership(&_OwnerIsCreator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnerIsCreator *OwnerIsCreatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.AcceptOwnership(&_OwnerIsCreator.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OwnerIsCreator *OwnerIsCreatorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OwnerIsCreator.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OwnerIsCreator *OwnerIsCreatorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.TransferOwnership(&_OwnerIsCreator.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OwnerIsCreator *OwnerIsCreatorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.TransferOwnership(&_OwnerIsCreator.TransactOpts, to)
}

// OwnerIsCreatorOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the OwnerIsCreator contract.
type OwnerIsCreatorOwnershipTransferRequestedIterator struct {
	Event *OwnerIsCreatorOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnerIsCreatorOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnerIsCreatorOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OwnerIsCreatorOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnerIsCreatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnerIsCreatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnerIsCreatorOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the OwnerIsCreator contract.
type OwnerIsCreatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OwnerIsCreatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OwnerIsCreator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorOwnershipTransferRequestedIterator{contract: _OwnerIsCreator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OwnerIsCreatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OwnerIsCreator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnerIsCreatorOwnershipTransferRequested)
				if err := _OwnerIsCreator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*OwnerIsCreatorOwnershipTransferRequested, error) {
	event := new(OwnerIsCreatorOwnershipTransferRequested)
	if err := _OwnerIsCreator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnerIsCreatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OwnerIsCreator contract.
type OwnerIsCreatorOwnershipTransferredIterator struct {
	Event *OwnerIsCreatorOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnerIsCreatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnerIsCreatorOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OwnerIsCreatorOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnerIsCreatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnerIsCreatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnerIsCreatorOwnershipTransferred represents a OwnershipTransferred event raised by the OwnerIsCreator contract.
type OwnerIsCreatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OwnerIsCreatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OwnerIsCreator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorOwnershipTransferredIterator{contract: _OwnerIsCreator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnerIsCreatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OwnerIsCreator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnerIsCreatorOwnershipTransferred)
				if err := _OwnerIsCreator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) ParseOwnershipTransferred(log types.Log) (*OwnerIsCreatorOwnershipTransferred, error) {
	event := new(OwnerIsCreatorOwnershipTransferred)
	if err := _OwnerIsCreator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TypeAndVersionInterfaceMetaData contains all meta data concerning the TypeAndVersionInterface contract.
var TypeAndVersionInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// TypeAndVersionInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use TypeAndVersionInterfaceMetaData.ABI instead.
var TypeAndVersionInterfaceABI = TypeAndVersionInterfaceMetaData.ABI

// TypeAndVersionInterface is an auto generated Go binding around an Ethereum contract.
type TypeAndVersionInterface struct {
	TypeAndVersionInterfaceCaller     // Read-only binding to the contract
	TypeAndVersionInterfaceTransactor // Write-only binding to the contract
	TypeAndVersionInterfaceFilterer   // Log filterer for contract events
}

// TypeAndVersionInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypeAndVersionInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypeAndVersionInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TypeAndVersionInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypeAndVersionInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TypeAndVersionInterfaceSession struct {
	Contract     *TypeAndVersionInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TypeAndVersionInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TypeAndVersionInterfaceCallerSession struct {
	Contract *TypeAndVersionInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// TypeAndVersionInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TypeAndVersionInterfaceTransactorSession struct {
	Contract     *TypeAndVersionInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// TypeAndVersionInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type TypeAndVersionInterfaceRaw struct {
	Contract *TypeAndVersionInterface // Generic contract binding to access the raw methods on
}

// TypeAndVersionInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceCallerRaw struct {
	Contract *TypeAndVersionInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// TypeAndVersionInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceTransactorRaw struct {
	Contract *TypeAndVersionInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTypeAndVersionInterface creates a new instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterface(address common.Address, backend bind.ContractBackend) (*TypeAndVersionInterface, error) {
	contract, err := bindTypeAndVersionInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterface{TypeAndVersionInterfaceCaller: TypeAndVersionInterfaceCaller{contract: contract}, TypeAndVersionInterfaceTransactor: TypeAndVersionInterfaceTransactor{contract: contract}, TypeAndVersionInterfaceFilterer: TypeAndVersionInterfaceFilterer{contract: contract}}, nil
}

// NewTypeAndVersionInterfaceCaller creates a new read-only instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterfaceCaller(address common.Address, caller bind.ContractCaller) (*TypeAndVersionInterfaceCaller, error) {
	contract, err := bindTypeAndVersionInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterfaceCaller{contract: contract}, nil
}

// NewTypeAndVersionInterfaceTransactor creates a new write-only instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*TypeAndVersionInterfaceTransactor, error) {
	contract, err := bindTypeAndVersionInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterfaceTransactor{contract: contract}, nil
}

// NewTypeAndVersionInterfaceFilterer creates a new log filterer instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*TypeAndVersionInterfaceFilterer, error) {
	contract, err := bindTypeAndVersionInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterfaceFilterer{contract: contract}, nil
}

// bindTypeAndVersionInterface binds a generic wrapper to an already deployed contract.
func bindTypeAndVersionInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TypeAndVersionInterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TypeAndVersionInterface.Contract.TypeAndVersionInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersionInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersionInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TypeAndVersionInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.contract.Transact(opts, method, params...)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_TypeAndVersionInterface *TypeAndVersionInterfaceCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TypeAndVersionInterface.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_TypeAndVersionInterface *TypeAndVersionInterfaceSession) TypeAndVersion() (string, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersion(&_TypeAndVersionInterface.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_TypeAndVersionInterface *TypeAndVersionInterfaceCallerSession) TypeAndVersion() (string, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersion(&_TypeAndVersionInterface.CallOpts)
}
