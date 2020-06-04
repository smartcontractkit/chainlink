// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package flux_aggregator_wrapper

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AggregatorInterfaceABI is the input ABI used to generate the binding from.
const AggregatorInterfaceABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"answeredInRound\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"answeredInRound\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AggregatorInterfaceFuncSigs maps the 4-byte function signature to its string representation.
var AggregatorInterfaceFuncSigs = map[string]string{
	"313ce567": "decimals()",
	"7284e416": "description()",
	"b5ab58dc": "getAnswer(uint256)",
	"0720da52": "getRoundData(uint256)",
	"b633620c": "getTimestamp(uint256)",
	"50d25bcd": "latestAnswer()",
	"668a0f02": "latestRound()",
	"feaf968c": "latestRoundData()",
	"8205bf6a": "latestTimestamp()",
	"54fd4d50": "version()",
}

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
	parsed, err := abi.JSON(strings.NewReader(AggregatorInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorInterface *AggregatorInterfaceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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
func (_AggregatorInterface *AggregatorInterfaceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_AggregatorInterface *AggregatorInterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _AggregatorInterface.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_AggregatorInterface *AggregatorInterfaceSession) Decimals() (uint8, error) {
	return _AggregatorInterface.Contract.Decimals(&_AggregatorInterface.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) Decimals() (uint8, error) {
	return _AggregatorInterface.Contract.Decimals(&_AggregatorInterface.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() constant returns(string)
func (_AggregatorInterface *AggregatorInterfaceCaller) Description(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _AggregatorInterface.contract.Call(opts, out, "description")
	return *ret0, err
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() constant returns(string)
func (_AggregatorInterface *AggregatorInterfaceSession) Description() (string, error) {
	return _AggregatorInterface.Contract.Description(&_AggregatorInterface.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() constant returns(string)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) Description() (string, error) {
	return _AggregatorInterface.Contract.Description(&_AggregatorInterface.CallOpts)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) constant returns(int256)
func (_AggregatorInterface *AggregatorInterfaceCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorInterface.contract.Call(opts, out, "getAnswer", roundId)
	return *ret0, err
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) constant returns(int256)
func (_AggregatorInterface *AggregatorInterfaceSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetAnswer(&_AggregatorInterface.CallOpts, roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) constant returns(int256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetAnswer(&_AggregatorInterface.CallOpts, roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x0720da52.
//
// Solidity: function getRoundData(uint256 _roundId) constant returns(uint256 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint256 answeredInRound)
func (_AggregatorInterface *AggregatorInterfaceCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _AggregatorInterface.contract.Call(opts, out, "getRoundData", _roundId)
	return *ret, err
}

// GetRoundData is a free data retrieval call binding the contract method 0x0720da52.
//
// Solidity: function getRoundData(uint256 _roundId) constant returns(uint256 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint256 answeredInRound)
func (_AggregatorInterface *AggregatorInterfaceSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorInterface.Contract.GetRoundData(&_AggregatorInterface.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x0720da52.
//
// Solidity: function getRoundData(uint256 _roundId) constant returns(uint256 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint256 answeredInRound)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorInterface.Contract.GetRoundData(&_AggregatorInterface.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) constant returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorInterface.contract.Call(opts, out, "getTimestamp", roundId)
	return *ret0, err
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) constant returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetTimestamp(&_AggregatorInterface.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) constant returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetTimestamp(&_AggregatorInterface.CallOpts, roundId)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() constant returns(int256)
func (_AggregatorInterface *AggregatorInterfaceCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorInterface.contract.Call(opts, out, "latestAnswer")
	return *ret0, err
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() constant returns(int256)
func (_AggregatorInterface *AggregatorInterfaceSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestAnswer(&_AggregatorInterface.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() constant returns(int256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestAnswer(&_AggregatorInterface.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() constant returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorInterface.contract.Call(opts, out, "latestRound")
	return *ret0, err
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() constant returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceSession) LatestRound() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestRound(&_AggregatorInterface.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() constant returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestRound() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestRound(&_AggregatorInterface.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() constant returns(uint256 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint256 answeredInRound)
func (_AggregatorInterface *AggregatorInterfaceCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _AggregatorInterface.contract.Call(opts, out, "latestRoundData")
	return *ret, err
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() constant returns(uint256 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint256 answeredInRound)
func (_AggregatorInterface *AggregatorInterfaceSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorInterface.Contract.LatestRoundData(&_AggregatorInterface.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() constant returns(uint256 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint256 answeredInRound)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorInterface.Contract.LatestRoundData(&_AggregatorInterface.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() constant returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorInterface.contract.Call(opts, out, "latestTimestamp")
	return *ret0, err
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() constant returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestTimestamp(&_AggregatorInterface.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() constant returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestTimestamp(&_AggregatorInterface.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorInterface.contract.Call(opts, out, "version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceSession) Version() (*big.Int, error) {
	return _AggregatorInterface.Contract.Version(&_AggregatorInterface.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) Version() (*big.Int, error) {
	return _AggregatorInterface.Contract.Version(&_AggregatorInterface.CallOpts)
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
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp)
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
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp)
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
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp)
func (_AggregatorInterface *AggregatorInterfaceFilterer) ParseAnswerUpdated(log types.Log) (*AggregatorInterfaceAnswerUpdated, error) {
	event := new(AggregatorInterfaceAnswerUpdated)
	if err := _AggregatorInterface.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
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
	return event, nil
}

// FluxAggregatorABI is the input ABI used to generate the binding from.
const FluxAggregatorABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"_paymentAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"_timeout\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"_description\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AvailableFundsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"OracleAdminUpdateRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"OracleAdminUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"whitelisted\",\"type\":\"bool\"}],\"name\":\"OraclePermissionsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"authorized\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"delay\",\"type\":\"uint32\"}],\"name\":\"RequesterPermissionsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint128\",\"name\":\"paymentAmount\",\"type\":\"uint128\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"minSubmissionCount\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"maxSubmissionCount\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"restartDelay\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"timeout\",\"type\":\"uint32\"}],\"name\":\"RoundDetailsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"submission\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"round\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"SubmissionReceived\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"acceptAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_oracles\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_admins\",\"type\":\"address[]\"},{\"internalType\":\"uint32\",\"name\":\"_minSubmissions\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maxSubmissions\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_restartDelay\",\"type\":\"uint32\"}],\"name\":\"addOracles\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"allocatedFunds\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"availableFunds\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracles\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"answeredInRound\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"answeredInRound\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"latestSubmission\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxSubmissionCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minSubmissionCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_queriedRoundId\",\"type\":\"uint32\"}],\"name\":\"oracleRoundState\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_eligibleToSubmit\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"_roundId\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"_latestSubmission\",\"type\":\"int256\"},{\"internalType\":\"uint64\",\"name\":\"_startedAt\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_timeout\",\"type\":\"uint64\"},{\"internalType\":\"uint128\",\"name\":\"_availableFunds\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"_oracleCount\",\"type\":\"uint32\"},{\"internalType\":\"uint128\",\"name\":\"_paymentAmount\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paymentAmount\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_oracles\",\"type\":\"address[]\"},{\"internalType\":\"uint32\",\"name\":\"_minSubmissions\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maxSubmissions\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_restartDelay\",\"type\":\"uint32\"}],\"name\":\"removeOracles\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reportingRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestNewRound\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"restartDelay\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_authorized\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"_delay\",\"type\":\"uint32\"}],\"name\":\"setRequesterPermissions\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"_submission\",\"type\":\"int256\"}],\"name\":\"submit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeout\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateAvailableFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"_paymentAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"_minSubmissions\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maxSubmissions\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_restartDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_timeout\",\"type\":\"uint32\"}],\"name\":\"updateFutureRounds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"withdrawablePayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// FluxAggregatorFuncSigs maps the 4-byte function signature to its string representation.
var FluxAggregatorFuncSigs = map[string]string{
	"628806ef": "acceptAdmin(address)",
	"79ba5097": "acceptOwnership()",
	"bbf0b7e9": "addOracles(address[],address[],uint32,uint32,uint32)",
	"d4cc54e4": "allocatedFunds()",
	"46fcff4c": "availableFunds()",
	"313ce567": "decimals()",
	"7284e416": "description()",
	"64efb22b": "getAdmin(address)",
	"b5ab58dc": "getAnswer(uint256)",
	"40884c52": "getOracles()",
	"0720da52": "getRoundData(uint256)",
	"b633620c": "getTimestamp(uint256)",
	"50d25bcd": "latestAnswer()",
	"668a0f02": "latestRound()",
	"feaf968c": "latestRoundData()",
	"bb07bacd": "latestSubmission(address)",
	"8205bf6a": "latestTimestamp()",
	"57970e93": "linkToken()",
	"58609e44": "maxSubmissionCount()",
	"c9374500": "minSubmissionCount()",
	"a4c0ed36": "onTokenTransfer(address,uint256,bytes)",
	"613d8fcc": "oracleCount()",
	"88aa80e7": "oracleRoundState(address,uint32)",
	"8da5cb5b": "owner()",
	"c35905c6": "paymentAmount()",
	"ebf8571c": "removeOracles(address[],uint32,uint32,uint32)",
	"6fb4bb4e": "reportingRound()",
	"98e5b12a": "requestNewRound()",
	"357ebb02": "restartDelay()",
	"20ed0275": "setRequesterPermissions(address,bool,uint32)",
	"202ee0ed": "submit(uint256,int256)",
	"70dea79a": "timeout()",
	"e9ee6eeb": "transferAdmin(address,address)",
	"f2fde38b": "transferOwnership(address)",
	"4f8fc3b5": "updateAvailableFunds()",
	"38aa4c72": "updateFutureRounds(uint128,uint32,uint32,uint32,uint32)",
	"54fd4d50": "version()",
	"c1075329": "withdrawFunds(address,uint256)",
	"3d3d7714": "withdrawPayment(address,address,uint256)",
	"e2e40317": "withdrawablePayment(address)",
}

// FluxAggregatorBin is the compiled bytecode used for deploying new contracts.
var FluxAggregatorBin = "0x60806040523480156200001157600080fd5b50604051620041ad380380620041ad833981810160405260a08110156200003757600080fd5b81516020830151604080850151606086015160808701805193519597949692959194919392820192846401000000008211156200007357600080fd5b9083019060208201858111156200008957600080fd5b8251640100000000811182820188101715620000a457600080fd5b82525081516020918201929091019080838360005b83811015620000d3578181015183820152602001620000b9565b50505050905090810190601f168015620001015780820380516001836020036101000a031916815260200191505b50604052505060008054336001600160a01b031991821617909155600280549091166001600160a01b03881617905550600480546001600160801b0319166001600160801b038616176001600160e01b0316600160e01b63ffffffff8616021790556005805460ff191660ff84161790558051620001879060069060208401906200026a565b50620001a88363ffffffff16426200020c60201b620028901790919060201c565b6000805260096020527fec8156718a8372b1db44bb411437d0870f3e3790d4a08526d024ce1b0b668f6c80546001600160401b03929092166801000000000000000002600160401b600160801b0319909216919091179055506200030f9350505050565b60008282111562000264576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620002ad57805160ff1916838001178555620002dd565b82800160010185558215620002dd579182015b82811115620002dd578251825591602001919060010190620002c0565b50620002eb929150620002ef565b5090565b6200030c91905b80821115620002eb5760008155600101620002f6565b90565b613e8e806200031f6000396000f3fe608060405234801561001057600080fd5b50600436106102485760003560e01c80637284e4161161013b578063bbf0b7e9116100b8578063e2e403171161007c578063e2e403171461083a578063e9ee6eeb14610860578063ebf8571c1461088e578063f2fde38b14610914578063feaf968c1461093a57610248565b8063bbf0b7e914610720578063c1075329146107f6578063c35905c614610822578063c93745001461082a578063d4cc54e41461083257610248565b806398e5b12a116100ff57806398e5b12a1461061c578063a4c0ed3614610624578063b5ab58dc146106a7578063b633620c146106c4578063bb07bacd146106e157610248565b80637284e416146104f157806379ba50971461056e5780638205bf6a1461057657806388aa80e71461057e5780638da5cb5b1461061457610248565b806350d25bcd116101c9578063628806ef1161018d578063628806ef1461048d57806364efb22b146104b3578063668a0f02146104d95780636fb4bb4e146104e157806370dea79a146104e957610248565b806350d25bcd1461043757806354fd4d501461045157806357970e931461045957806358609e441461047d578063613d8fcc1461048557610248565b806338aa4c721161021057806338aa4c72146103335780633d3d77141461037d57806340884c52146103b357806346fcff4c1461040b5780634f8fc3b51461042f57610248565b80630720da521461024d578063202ee0ed1461029557806320ed0275146102ba578063313ce567146102f4578063357ebb0214610312575b600080fd5b61026a6004803603602081101561026357600080fd5b5035610942565b6040805195865260208601949094528484019290925260608401526080830152519081900360a00190f35b6102b8600480360360408110156102ab57600080fd5b5080359060200135610965565b005b6102b8600480360360608110156102d057600080fd5b5080356001600160a01b03169060208101351515906040013563ffffffff16610a36565b6102fc610b6a565b6040805160ff9092168252519081900360200190f35b61031a610b73565b6040805163ffffffff9092168252519081900360200190f35b6102b8600480360360a081101561034957600080fd5b506001600160801b038135169063ffffffff6020820135811691604081013582169160608201358116916080013516610b86565b6102b86004803603606081101561039357600080fd5b506001600160a01b03813581169160208101359091169060400135610f30565b6103bb611123565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156103f75781810151838201526020016103df565b505050509050019250505060405180910390f35b610413611186565b604080516001600160801b039092168252519081900360200190f35b6102b861119c565b61043f611296565b60408051918252519081900360200190f35b61043f6112a5565b6104616112aa565b604080516001600160a01b039092168252519081900360200190f35b61031a6112b9565b61031a6112cc565b6102b8600480360360208110156104a357600080fd5b50356001600160a01b03166112d2565b610461600480360360208110156104c957600080fd5b50356001600160a01b03166113b8565b61043f6113e2565b61043f6113f5565b61031a611401565b6104f9611414565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561053357818101518382015260200161051b565b50505050905090810190601f1680156105605780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102b86114a2565b61043f611551565b6105b06004803603604081101561059457600080fd5b5080356001600160a01b0316906020013563ffffffff1661155b565b60408051981515895263ffffffff97881660208a0152888101969096526001600160401b0394851660608901529290931660808701526001600160801b0390811660a08701529190931660c08501529190911660e083015251908190036101000190f35b610461611824565b6102b8611833565b6102b86004803603606081101561063a57600080fd5b6001600160a01b0382351691602081013591810190606081016040820135600160201b81111561066957600080fd5b82018360208201111561067b57600080fd5b803590602001918460018302840111600160201b8311171561069c57600080fd5b509092509050611944565b61043f600480360360208110156106bd57600080fd5b50356119a5565b61043f600480360360208110156106da57600080fd5b50356119b6565b610707600480360360208110156106f757600080fd5b50356001600160a01b03166119c1565b6040805192835260208301919091528051918290030190f35b6102b8600480360360a081101561073657600080fd5b810190602081018135600160201b81111561075057600080fd5b82018360208201111561076257600080fd5b803590602001918460208302840111600160201b8311171561078357600080fd5b919390929091602081019035600160201b8111156107a057600080fd5b8201836020820111156107b257600080fd5b803590602001918460208302840111600160201b831117156107d357600080fd5b919350915063ffffffff81358116916020810135821691604090910135166119f2565b6102b86004803603604081101561080c57600080fd5b506001600160a01b038135169060200135611b85565b610413611d2c565b61031a611d3b565b610413611d4e565b61043f6004803603602081101561085057600080fd5b50356001600160a01b0316611d5d565b6102b86004803603604081101561087657600080fd5b506001600160a01b0381358116916020013516611d81565b6102b8600480360360808110156108a457600080fd5b810190602081018135600160201b8111156108be57600080fd5b8201836020820111156108d057600080fd5b803590602001918460208302840111600160201b831117156108f157600080fd5b919350915063ffffffff8135811691602081013582169160409091013516611e64565b6102b86004803603602081101561092a57600080fd5b50356001600160a01b0316611f13565b61026a611fb1565b600080600080600061095386611fd2565b939a9299509097509550909350915050565b60606109713384612123565b905080516000148190610a025760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b838110156109c75781810151838201526020016109af565b50505050905090810190601f1680156109f45780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b50610a0c836123c5565b610a16828461245a565b610a1f83612536565b610a2883612685565b610a3183612788565b505050565b6000546001600160a01b03163314610a83576040805162461bcd60e51b81526020600482015260166024820152600080516020613e39833981519152604482015290519081900360640190fd5b6001600160a01b0383166000908152600a602052604090205460ff1615158215151415610aaf57610a31565b8115610af2576001600160a01b0383166000908152600a60205260409020805460ff19168315151764ffffffff00191661010063ffffffff841602179055610b1b565b6001600160a01b0383166000908152600a60205260409020805468ffffffffffffffffff191690555b60408051831515815263ffffffff8316602082015281516001600160a01b038616927fc3df5a754e002718f2e10804b99e6605e7c701d95cec9552c7680ca2b6f2820a928290030190a2505050565b60055460ff1681565b600454600160c01b900463ffffffff1681565b6000546001600160a01b03163314610bd3576040805162461bcd60e51b81526020600482015260166024820152600080516020613e39833981519152604482015290519081900360640190fd5b6000610bdd6112cc565b90508463ffffffff168463ffffffff161015610c40576040805162461bcd60e51b815260206004820152601960248201527f6d6178206d75737420657175616c2f657863656564206d696e00000000000000604482015290519081900360640190fd5b8363ffffffff168163ffffffff161015610ca1576040805162461bcd60e51b815260206004820152601760248201527f6d61782063616e6e6f742065786365656420746f74616c000000000000000000604482015290519081900360640190fd5b63ffffffff81161580610cbf57508263ffffffff168163ffffffff16115b610d10576040805162461bcd60e51b815260206004820152601960248201527f64656c61792063616e6e6f742065786365656420746f74616c00000000000000604482015290519081900360640190fd5b610d22866001600160801b03166127ee565b600354600160801b90046001600160801b03161015610d88576040805162461bcd60e51b815260206004820152601e60248201527f696e73756666696369656e742066756e647320666f72207061796d656e740000604482015290519081900360640190fd5b6000610d926112cc565b63ffffffff161115610df95760008563ffffffff1611610df9576040805162461bcd60e51b815260206004820152601a60248201527f6d696e206d7573742062652067726561746572207468616e2030000000000000604482015290519081900360640190fd5b85600460006101000a8154816001600160801b0302191690836001600160801b0316021790555084600460146101000a81548163ffffffff021916908363ffffffff16021790555083600460106101000a81548163ffffffff021916908363ffffffff16021790555082600460186101000a81548163ffffffff021916908363ffffffff160217905550816004601c6101000a81548163ffffffff021916908363ffffffff1602179055508363ffffffff168563ffffffff16600460009054906101000a90046001600160801b03166001600160801b03167f56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f8686604051808363ffffffff1663ffffffff1681526020018263ffffffff1663ffffffff1681526020019250505060405180910390a4505050505050565b6001600160a01b03838116600090815260086020526040902060020154620100009004163314610fa0576040805162461bcd60e51b815260206004820152601660248201527537b7363c9031b0b63630b1363290313c9030b236b4b760511b604482015290519081900360640190fd5b6001600160a01b03831660009081526008602052604090205481906001600160801b0390811690821681101561101d576040805162461bcd60e51b815260206004820152601f60248201527f696e73756666696369656e7420776974686472617761626c652066756e647300604482015290519081900360640190fd5b6110366001600160801b0382168363ffffffff61282116565b6001600160a01b038616600090815260086020526040902080546001600160801b0319166001600160801b03928316179055600354611076911683612821565b600380546001600160801b0319166001600160801b039283161790556002546040805163a9059cbb60e01b81526001600160a01b03888116600483015293861660248201529051929091169163a9059cbb916044808201926020929091908290030181600087803b1580156110ea57600080fd5b505af11580156110fe573d6000803e3d6000fd5b505050506040513d602081101561111457600080fd5b505161111c57fe5b5050505050565b6060600b80548060200260200160405190810160405280929190818152602001828054801561117b57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161115d575b505050505090505b90565b600354600160801b90046001600160801b031681565b600354600254604080516370a0823160e01b815230600482015290516001600160801b03600160801b850481169460009461123f9492909116926001600160a01b03909116916370a08231916024808301926020929190829003018186803b15801561120757600080fd5b505afa15801561121b573d6000803e3d6000fd5b505050506040513d602081101561123157600080fd5b50519063ffffffff61289016565b600380546001600160801b03908116600160801b8483160217909155909150821681146112925760405181907ffe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f90600090a25b5050565b60006112a06128e7565b905090565b600381565b6002546001600160a01b031681565b600454600160801b900463ffffffff1681565b600b5490565b6001600160a01b03818116600090815260086020526040902060030154163314611343576040805162461bcd60e51b815260206004820152601e60248201527f6f6e6c792063616c6c61626c652062792070656e64696e672061646d696e0000604482015290519081900360640190fd5b6001600160a01b0381166000818152600860205260408082206003810180546001600160a01b0319169055600201805462010000600160b01b031916336201000081029190911790915590519092917f0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe90491a350565b6001600160a01b03808216600090815260086020526040902060020154620100009004165b919050565b600754600160201b900463ffffffff1690565b60075463ffffffff1690565b600454600160e01b900463ffffffff1681565b6006805460408051602060026001851615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561149a5780601f1061146f5761010080835404028352916020019161149a565b820191906000526020600020905b81548152906001019060200180831161147d57829003601f168201915b505050505081565b6001546001600160a01b031633146114fa576040805162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b604482015290519081900360640190fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60006112a0612909565b6000808080808080803332146115b1576040805162461bcd60e51b81526020600482015260166024820152756f66662d636861696e2072656164696e67206f6e6c7960501b604482015290519081900360640190fd5b63ffffffff891660008181526009602090815260408083206001600160a01b038f168452600890925290912090911561169d576115ee8c8c61293e565b6001600160a01b038d166000908152600860205260409020600190810154908401546003808601549054939d508e9c50919a506001600160401b03169850600160401b900463ffffffff169650600160801b90046001600160801b031694506116556112cc565b60018301549094506001600160401b031661167b576004546001600160801b0316611691565b6003820154600160601b90046001600160801b03165b90920191506118179050565b60075481546000918291600160c01b900463ffffffff908116911614806116d357506007546116d19063ffffffff16612993565b155b6007549091506116e89063ffffffff166129b3565b80156116f15750805b1561174e5760075461170f9063ffffffff908116906001906129ed16565b63ffffffff81166000908152600960205260409020600454919c506001600160801b03909116955093506117438e8c612a45565b9b5060019150611792565b60075463ffffffff1660008181526009602052604090206003810154919c50600160601b9091046001600160801b03169550935061178b8b612993565b9b50600291505b61179c8e8c612123565b51156117ab5760009b50600391505b6001838101549085015460038087015490548f938f9390926001600160401b0390911691600160401b90910463ffffffff1690600160801b90046001600160801b031687016117f86112cc565b8b8363ffffffff1693509b509b509b509b509b509b509b509b50505050505b9295985092959890939650565b6000546001600160a01b031681565b336000908152600a602052604090205460ff16611897576040805162461bcd60e51b815260206004820152601860248201527f6e6f7420617574686f72697a6564207265717565737465720000000000000000604482015290519081900360640190fd5b60075463ffffffff16600081815260096020526040902060010154600160401b90046001600160401b03161515806118d357506118d381612a92565b611924576040805162461bcd60e51b815260206004820152601f60248201527f7072657620726f756e64206d75737420626520737570657273656461626c6500604482015290519081900360640190fd5b61194161193c63ffffffff808416906001906129ed16565b612b14565b50565b8015611997576040805162461bcd60e51b815260206004820181905260248201527f7472616e7366657220646f65736e2774206163636570742063616c6c64617461604482015290519081900360640190fd5b61199f61119c565b50505050565b60006119b082612be2565b92915050565b60006119b082612bfa565b6001600160a01b03166000908152600860205260409020600181015490549091600160c01b90910463ffffffff1690565b6000546001600160a01b03163314611a3f576040805162461bcd60e51b81526020600482015260166024820152600080516020613e39833981519152604482015290519081900360640190fd5b858414611a93576040805162461bcd60e51b815260206004820181905260248201527f6e6565642073616d65206f7261636c6520616e642061646d696e20636f756e74604482015290519081900360640190fd5b604d611ab587611aa16112cc565b63ffffffff16612c2590919063ffffffff16565b1115611afe576040805162461bcd60e51b81526020600482015260136024820152721b585e081bdc9858db195cc8185b1b1bddd959606a1b604482015290519081900360640190fd5b60005b86811015611b5157611b49888883818110611b1857fe5b905060200201356001600160a01b0316878784818110611b3457fe5b905060200201356001600160a01b0316612c6d565b600101611b01565b50600454611b7c906001600160801b03811690859085908590600160e01b900463ffffffff16610b86565b50505050505050565b6000546001600160a01b03163314611bd2576040805162461bcd60e51b81526020600482015260166024820152600080516020613e39833981519152604482015290519081900360640190fd5b6004548190611c0690611bed906001600160801b03166127ee565b600354600160801b90046001600160801b031690612890565b1015611c59576040805162461bcd60e51b815260206004820152601a60248201527f696e73756666696369656e7420726573657276652066756e6473000000000000604482015290519081900360640190fd5b6002546040805163a9059cbb60e01b81526001600160a01b038581166004830152602482018590529151919092169163a9059cbb9160448083019260209291908290030181600087803b158015611caf57600080fd5b505af1158015611cc3573d6000803e3d6000fd5b505050506040513d6020811015611cd957600080fd5b5051611d24576040805162461bcd60e51b81526020600482015260156024820152741d1bdad95b881d1c985b9cd9995c8819985a5b1959605a1b604482015290519081900360640190fd5b61129261119c565b6004546001600160801b031681565b600454600160a01b900463ffffffff1681565b6003546001600160801b031681565b6001600160a01b03166000908152600860205260409020546001600160801b031690565b6001600160a01b03828116600090815260086020526040902060020154620100009004163314611df1576040805162461bcd60e51b815260206004820152601660248201527537b7363c9031b0b63630b1363290313c9030b236b4b760511b604482015290519081900360640190fd5b6001600160a01b0382811660008181526008602090815260409182902060030180546001600160a01b0319169486169485179055815133815290810193909352805191927fb79bf2e89c2d70dde91d2991fb1ea69b7e478061ad7c04ed5b02b96bc52b8104929081900390910190a25050565b6000546001600160a01b03163314611eb1576040805162461bcd60e51b81526020600482015260166024820152600080516020613e39833981519152604482015290519081900360640190fd5b60005b84811015611ee857611ee0868683818110611ecb57fe5b905060200201356001600160a01b0316612efd565b600101611eb4565b5060045461111c906001600160801b03811690859085908590600160e01b900463ffffffff16610b86565b6000546001600160a01b03163314611f60576040805162461bcd60e51b81526020600482015260166024820152600080516020613e39833981519152604482015290519081900360640190fd5b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000806000806000611fc16130d6565b945094509450945094509091929394565b6000806000806000611fe2613d20565b63ffffffff808816600090815260096020908152604091829020825160a080820185528254825260018301546001600160401b0380821684870152600160401b82041683870152600160801b90049095166060820152835160028301805460c09581028301860190965295810185815291959294608087019491939284929091849184018282801561209357602002820191906000526020600020905b81548152602001906001019080831161207f575b50505091835250506001919091015463ffffffff808216602080850191909152600160201b83048216604080860191909152600160401b84048316606080870191909152600160601b9094046001600160801b031660809095019490945293909452855192860151918601519501519b9c919b6001600160401b039182169b509416985092169550909350505050565b6001600160a01b03821660009081526008602052604090205460075460609163ffffffff600160801b909104811691168161218c57604051806040016040528060128152602001716e6f7420656e61626c6564206f7261636c6560701b815250925050506119b0565b8363ffffffff168263ffffffff1611156121d857604051806040016040528060168152602001756e6f742079657420656e61626c6564206f7261636c6560501b815250925050506119b0565b6001600160a01b03851660009081526008602052604090205463ffffffff808616600160a01b909204161015612247576040518060400160405280601881526020017f6e6f206c6f6e67657220616c6c6f776564206f7261636c650000000000000000815250925050506119b0565b6001600160a01b03851660009081526008602052604090205463ffffffff808616600160c01b90920416106122b5576040518060400160405280602081526020017f63616e6e6f74207265706f7274206f6e2070726576696f757320726f756e6473815250925050506119b0565b8063ffffffff168463ffffffff16141580156122f157506122e163ffffffff808316906001906129ed16565b63ffffffff168463ffffffff1614155b8015612304575061230284826130ff565b155b15612348576040518060400160405280601781526020017f696e76616c696420726f756e6420746f207265706f7274000000000000000000815250925050506119b0565b8363ffffffff16600114158015612379575061237761237263ffffffff8087169060019061315f16565b6129b3565b155b156123bd576040518060400160405280601f81526020017f70726576696f757320726f756e64206e6f7420737570657273656461626c6500815250925050506119b0565b505092915050565b6123ce816131c2565b6123d757611941565b3360009081526008602052604090205460045463ffffffff600160e01b909204821691600160c01b909104811682019083161180159061241657508015155b156124215750611941565b61242a826131f3565b50336000908152600860205260409020805463ffffffff8316600160e01b026001600160e01b0390911617905550565b61246381612993565b6124b4576040805162461bcd60e51b815260206004820152601f60248201527f726f756e64206e6f7420616363657074696e67207375626d697373696f6e7300604482015290519081900360640190fd5b63ffffffff8116600081815260096020908152604080832060020180546001808201835591855283852001879055338085526008909352818420805463ffffffff60c01b1916600160c01b8702178155018690555190929185917f92e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68c9190a45050565b63ffffffff80821660009081526009602052604090206003810154600290910154600160201b909104909116111561256d57611941565b63ffffffff8116600090815260096020908152604080832060020180548251818502810185019093528083526125d6938301828280156125cc57602002820191906000526020600020905b8154815260200190600101908083116125b8575b5050505050613338565b63ffffffff8316600081815260096020908152604091829020848155600101805467ffffffffffffffff60401b1916600160401b426001600160401b038116919091029190911763ffffffff60801b1916600160801b8602179091556007805467ffffffff000000001916600160201b860217905582519081529151939450919284927f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f928290030190a35050565b63ffffffff808216600090815260096020526040812060039081015490546001600160801b03600160601b9092048216936126ce92600160801b90920490911690849061282116565b600380546001600160801b03908116600160801b8483160217918290559192506126f99116836133e1565b600380546001600160801b0319166001600160801b03928316179055336000908152600860205260409020546127309116836133e1565b3360009081526008602052604080822080546001600160801b0319166001600160801b0394851617905551918316917ffe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f9190a2505050565b63ffffffff80821660009081526009602052604090206003810154600290910154911611156127b657611941565b63ffffffff81166000908152600960205260408120600201906127d98282613d6c565b5060010180546001600160e01b031916905550565b60006119b060026128156128006112cc565b63ffffffff168561343590919063ffffffff16565b9063ffffffff61343516565b6000826001600160801b0316826001600160801b0316111561288a576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b60008282111561288a576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b600754600160201b900463ffffffff1660009081526009602052604090205490565b600754600160201b900463ffffffff16600090815260096020526040902060010154600160401b90046001600160401b031690565b63ffffffff81166000908152600960205260408120600101546001600160401b0316156129895761296e82612993565b8015612982575061297f8383612123565b51155b90506119b0565b61296e8383612a45565b63ffffffff90811660009081526009602052604090206003015416151590565b63ffffffff8116600090815260096020526040812060010154600160401b90046001600160401b03161515806119b057506119b082612a92565b600082820163ffffffff8085169082161015612a3e576040805162461bcd60e51b815260206004820152601b6024820152600080516020613df8833981519152604482015290519081900360640190fd5b9392505050565b6001600160a01b03821660009081526008602052604081205460045463ffffffff600160e01b909204821691600160c01b909104811682019084161180612a8a575080155b949350505050565b63ffffffff8082166000908152600960205260408120600181015460039091015491926001600160401b0390911691600160401b9004168115801590612ade575060008163ffffffff16115b8015612a8a575042612b026001600160401b03841663ffffffff8085169061348e16565b6001600160401b031610949350505050565b612b1d816131c2565b612b2657611941565b336000908152600a602052604090205463ffffffff6501000000000082048116916101009004811682019083161180612b5d575080155b612ba4576040805162461bcd60e51b81526020600482015260136024820152726d7573742064656c617920726571756573747360681b604482015290519081900360640190fd5b612bad826131f3565b50336000908152600a60205260409020805463ffffffff8316650100000000000268ffffffff00000000001990911617905550565b63ffffffff1660009081526009602052604090205490565b63ffffffff16600090815260096020526040902060010154600160401b90046001600160401b031690565b600082820183811015612a3e576040805162461bcd60e51b815260206004820152601b6024820152600080516020613df8833981519152604482015290519081900360640190fd5b612c76826134e2565b15612cc1576040805162461bcd60e51b81526020600482015260166024820152751bdc9858db1948185b1c9958591e48195b98589b195960521b604482015290519081900360640190fd5b6001600160a01b038116612d14576040805162461bcd60e51b8152602060048201526015602482015274063616e6e6f74207365742061646d696e20746f203605c1b604482015290519081900360640190fd5b6001600160a01b03828116600090815260086020526040902060020154620100009004161580612d6957506001600160a01b038281166000908152600860205260409020600201546201000090048116908216145b612dba576040805162461bcd60e51b815260206004820152601c60248201527f6f776e65722063616e6e6f74206f76657277726974652061646d696e00000000604482015290519081900360640190fd5b612dc38261350d565b6001600160a01b03808416600081815260086020526040808220805463ffffffff60a01b1963ffffffff97909716600160801b0263ffffffff60801b19909116179590951663ffffffff60a01b178555600b80546002909601805461ffff90971661ffff19909716969096178655805460018181019092557f0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db90180546001600160a01b031916851790558383528554948716620100000262010000600160b01b0319909516949094179094559251919290917f18dd09695e4fbdae8d1a5edb11221eb04564269c29a089b9753a6535c54ba92e9190a3806001600160a01b0316826001600160a01b03167f0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe90460405160405180910390a35050565b612f06816134e2565b612f4c576040805162461bcd60e51b81526020600482015260126024820152711bdc9858db19481b9bdd08195b98589b195960721b604482015290519081900360640190fd5b600754612f659063ffffffff908116906001906129ed16565b6001600160a01b0382166000908152600860205260408120805463ffffffff93909316600160a01b0263ffffffff60a01b1990931692909217909155600b612fc46001612fb06112cc565b63ffffffff1661315f90919063ffffffff16565b63ffffffff1681548110612fd457fe5b6000918252602080832091909101546001600160a01b0385811680855260089093526040808520600290810180549390941680875291862001805461ffff90931661ffff199384168117909155939094528154169055600b805492935090918391908390811061304057fe5b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b03160217905550600b80548061307957fe5b600082815260208120820160001990810180546001600160a01b03191690559091019091556040516001600160a01b038516907f18dd09695e4fbdae8d1a5edb11221eb04564269c29a089b9753a6535c54ba92e908390a3505050565b6000806000806000611fc1600760049054906101000a900463ffffffff1663ffffffff16611fd2565b60008163ffffffff1661312260018563ffffffff166129ed90919063ffffffff16565b63ffffffff16148015612a3e57505063ffffffff16600090815260096020526040902060010154600160401b90046001600160401b031615919050565b60008263ffffffff168263ffffffff16111561288a576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b6007546000906131de9063ffffffff908116906001906129ed16565b63ffffffff168263ffffffff16149050919050565b61321061320b63ffffffff8084169060019061315f16565b613570565b6007805463ffffffff1990811663ffffffff84811691821790935560048054600083815260096020908152604091829020600381018054600160801b90950489169490971693909317808755845467ffffffff0000000019909116600160a01b9091048816600160201b021780875584546fffffffffffffffffffffffffffffffff60601b199091166001600160801b03909116600160601b021780875593546bffffffff000000000000000019909416600160e01b909404909616600160401b0292909217909355600192909201805467ffffffffffffffff1916426001600160401b0390811691909117918290558351911681529151339391927f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac6027192908290030190a350565b60008151600010613389576040805162461bcd60e51b81526020600482015260166024820152756c697374206d757374206e6f7420626520656d70747960501b604482015290519081900360640190fd5b815160028104600182166133c8576000806133ae86600060018703600187038761362e565b90925090506133bd828261370c565b9450505050506113dd565b6133d8846000600185038461377a565b925050506113dd565b60008282016001600160801b038085169082161015612a3e576040805162461bcd60e51b815260206004820152601b6024820152600080516020613df8833981519152604482015290519081900360640190fd5b600082613444575060006119b0565b8282028284828161345157fe5b0414612a3e5760405162461bcd60e51b8152600401808060200182810382526021815260200180613e186021913960400191505060405180910390fd5b60008282016001600160401b038085169082161015612a3e576040805162461bcd60e51b815260206004820152601b6024820152600080516020613df8833981519152604482015290519081900360640190fd5b6001600160a01b031660009081526008602052604090205463ffffffff600160a01b90910481161490565b60075460009063ffffffff16801580159061354f57506001600160a01b03831660009081526008602052604090205463ffffffff828116600160a01b90920416145b1561355b5790506113dd565b612a3e63ffffffff808316906001906129ed16565b61357981612a92565b61358257611941565b600061359963ffffffff8084169060019061315f16565b63ffffffff81811660009081526009602052604080822080548785168452918320918255600190810154908201805463ffffffff60801b1916600160801b928390049095169091029390931767ffffffffffffffff60401b1916600160401b426001600160401b03160217909255919250600201906136188282613d6c565b5060010180546001600160e01b03191690555050565b60008082841061363d57600080fd5b83861115801561364d5750848411155b61365657600080fd5b8286111580156136665750848311155b61366f57600080fd5b6007868603101561369057613687878787878761380b565b91509150613702565b600061369d888888613bde565b90508084116136ae578095506136fc565b848110156136c1578060010196506136fc565b8085111580156136d057508381105b6136d657fe5b6136e28888838861377a565b92506136f38882600101888761377a565b91506137029050565b5061366f565b9550959350505050565b6000808312801561371d5750600082135b8061373357506000831380156137335750600082125b156137535760026137448484613cbb565b8161374b57fe5b0590506119b0565b6000600280850781850701059050612a8a6137746002860560028605613cbb565b82613cbb565b60008184111561378957600080fd5b8282111561379657600080fd5b828410156137ed57600784840310156137c25760006137b8868686868761380b565b509150612a8a9050565b60006137cf868686613bde565b90508083116137e0578093506137e7565b8060010194505b50613796565b8484815181106137f957fe5b60200260200101519050949350505050565b60008060008686600101039050600088886000018151811061382957fe5b6020026020010151905060008260011061384a576001600160ff1b03613862565b89896001018151811061385957fe5b60200260200101515b905060008360021061387b576001600160ff1b03613893565b8a8a6002018151811061388a57fe5b60200260200101515b90506000846003106138ac576001600160ff1b036138c4565b8b8b600301815181106138bb57fe5b60200260200101515b90506000856004106138dd576001600160ff1b036138f5565b8c8c600401815181106138ec57fe5b60200260200101515b905060008660051061390e576001600160ff1b03613926565b8d8d6005018151811061391d57fe5b60200260200101515b905060008760061061393f576001600160ff1b03613957565b8e8e6006018151811061394e57fe5b60200260200101515b905085871315613965579495945b83851315613971579293925b8183131561397d579091905b84871315613989579395935b83861315613995579294925b8083131561399f57915b848613156139ab579394935b808213156139b557905b828713156139c1579195915b818613156139cd579094905b808513156139d757935b828613156139e3579194915b808413156139ed57925b828513156139f9579193915b81841315613a05579092905b82841315613a11579192915b8d8c0380613a2157879a50613ac7565b8060011415613a3257869a50613ac7565b8060021415613a4357859a50613ac7565b8060031415613a5457849a50613ac7565b8060041415613a6557839a50613ac7565b8060051415613a7657829a50613ac7565b8060061415613a8757819a50613ac7565b6040805162461bcd60e51b815260206004820152601060248201526f6b31206f7574206f6620626f756e647360801b604482015290519081900360640190fd5b8e8c038d8d1415613ae557508a995061370298505050505050505050565b80613afc5750969850613702975050505050505050565b8060011415613b175750959850613702975050505050505050565b8060021415613b325750949850613702975050505050505050565b8060031415613b4d5750939850613702975050505050505050565b8060041415613b685750929850613702975050505050505050565b8060051415613b835750919850613702975050505050505050565b8060061415613b9e5750909850613702975050505050505050565b6040805162461bcd60e51b815260206004820152601060248201526f6b32206f7574206f6620626f756e647360801b604482015290519081900360640190fd5b6000808460028585010481518110613bf257fe5b602002602001015190506001840393506001830192505b60018401935080858581518110613c1c57fe5b602002602001015112613c09575b60018303925080858481518110613c3d57fe5b602002602001015113613c2a5782841015613cad57848381518110613c5e57fe5b6020026020010151858581518110613c7257fe5b6020026020010151868681518110613c8657fe5b60200260200101878681518110613c9957fe5b602090810291909101019190915252613cb6565b82915050612a3e565b613c09565b6000828201818312801590613cd05750838112155b80613ce55750600083128015613ce557508381125b612a3e5760405162461bcd60e51b8152600401808060200182810382526021815260200180613dd76021913960400191505060405180910390fd5b6040518060a001604052806000815260200160006001600160401b0316815260200160006001600160401b03168152602001600063ffffffff168152602001613d67613d8a565b905290565b50805460008255906000526020600020908101906119419190613db8565b6040805160a08101825260608082526000602083018190529282018390528101829052608081019190915290565b61118391905b80821115613dd25760008155600101613dbe565b509056fe5369676e6564536166654d6174683a206164646974696f6e206f766572666c6f77536166654d6174683a206164646974696f6e206f766572666c6f770000000000536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f774f6e6c792063616c6c61626c65206279206f776e657200000000000000000000a26469706673582212205c8bca4eb815ba42d4dd127a7900162b285bcbb17b9283c956304759a35d7b6a64736f6c63430006060033"

// DeployFluxAggregator deploys a new Ethereum contract, binding an instance of FluxAggregator to it.
func DeployFluxAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address, _paymentAmount *big.Int, _timeout uint32, _decimals uint8, _description string) (common.Address, *types.Transaction, *FluxAggregator, error) {
	parsed, err := abi.JSON(strings.NewReader(FluxAggregatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(FluxAggregatorBin), backend, _link, _paymentAmount, _timeout, _decimals, _description)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FluxAggregator{FluxAggregatorCaller: FluxAggregatorCaller{contract: contract}, FluxAggregatorTransactor: FluxAggregatorTransactor{contract: contract}, FluxAggregatorFilterer: FluxAggregatorFilterer{contract: contract}}, nil
}

// FluxAggregator is an auto generated Go binding around an Ethereum contract.
type FluxAggregator struct {
	FluxAggregatorCaller     // Read-only binding to the contract
	FluxAggregatorTransactor // Write-only binding to the contract
	FluxAggregatorFilterer   // Log filterer for contract events
}

// FluxAggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type FluxAggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FluxAggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FluxAggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FluxAggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FluxAggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FluxAggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FluxAggregatorSession struct {
	Contract     *FluxAggregator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FluxAggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FluxAggregatorCallerSession struct {
	Contract *FluxAggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// FluxAggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FluxAggregatorTransactorSession struct {
	Contract     *FluxAggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// FluxAggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type FluxAggregatorRaw struct {
	Contract *FluxAggregator // Generic contract binding to access the raw methods on
}

// FluxAggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FluxAggregatorCallerRaw struct {
	Contract *FluxAggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// FluxAggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FluxAggregatorTransactorRaw struct {
	Contract *FluxAggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFluxAggregator creates a new instance of FluxAggregator, bound to a specific deployed contract.
func NewFluxAggregator(address common.Address, backend bind.ContractBackend) (*FluxAggregator, error) {
	contract, err := bindFluxAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FluxAggregator{FluxAggregatorCaller: FluxAggregatorCaller{contract: contract}, FluxAggregatorTransactor: FluxAggregatorTransactor{contract: contract}, FluxAggregatorFilterer: FluxAggregatorFilterer{contract: contract}}, nil
}

// NewFluxAggregatorCaller creates a new read-only instance of FluxAggregator, bound to a specific deployed contract.
func NewFluxAggregatorCaller(address common.Address, caller bind.ContractCaller) (*FluxAggregatorCaller, error) {
	contract, err := bindFluxAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorCaller{contract: contract}, nil
}

// NewFluxAggregatorTransactor creates a new write-only instance of FluxAggregator, bound to a specific deployed contract.
func NewFluxAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*FluxAggregatorTransactor, error) {
	contract, err := bindFluxAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorTransactor{contract: contract}, nil
}

// NewFluxAggregatorFilterer creates a new log filterer instance of FluxAggregator, bound to a specific deployed contract.
func NewFluxAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*FluxAggregatorFilterer, error) {
	contract, err := bindFluxAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorFilterer{contract: contract}, nil
}

// bindFluxAggregator binds a generic wrapper to an already deployed contract.
func bindFluxAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FluxAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FluxAggregator *FluxAggregatorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FluxAggregator.Contract.FluxAggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FluxAggregator *FluxAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FluxAggregator.Contract.FluxAggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FluxAggregator *FluxAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FluxAggregator.Contract.FluxAggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FluxAggregator *FluxAggregatorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FluxAggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FluxAggregator *FluxAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FluxAggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FluxAggregator *FluxAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FluxAggregator.Contract.contract.Transact(opts, method, params...)
}

// AllocatedFunds is a free data retrieval call binding the contract method 0xd4cc54e4.
//
// Solidity: function allocatedFunds() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorCaller) AllocatedFunds(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "allocatedFunds")
	return *ret0, err
}

// AllocatedFunds is a free data retrieval call binding the contract method 0xd4cc54e4.
//
// Solidity: function allocatedFunds() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorSession) AllocatedFunds() (*big.Int, error) {
	return _FluxAggregator.Contract.AllocatedFunds(&_FluxAggregator.CallOpts)
}

// AllocatedFunds is a free data retrieval call binding the contract method 0xd4cc54e4.
//
// Solidity: function allocatedFunds() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorCallerSession) AllocatedFunds() (*big.Int, error) {
	return _FluxAggregator.Contract.AllocatedFunds(&_FluxAggregator.CallOpts)
}

// AvailableFunds is a free data retrieval call binding the contract method 0x46fcff4c.
//
// Solidity: function availableFunds() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorCaller) AvailableFunds(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "availableFunds")
	return *ret0, err
}

// AvailableFunds is a free data retrieval call binding the contract method 0x46fcff4c.
//
// Solidity: function availableFunds() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorSession) AvailableFunds() (*big.Int, error) {
	return _FluxAggregator.Contract.AvailableFunds(&_FluxAggregator.CallOpts)
}

// AvailableFunds is a free data retrieval call binding the contract method 0x46fcff4c.
//
// Solidity: function availableFunds() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorCallerSession) AvailableFunds() (*big.Int, error) {
	return _FluxAggregator.Contract.AvailableFunds(&_FluxAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_FluxAggregator *FluxAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_FluxAggregator *FluxAggregatorSession) Decimals() (uint8, error) {
	return _FluxAggregator.Contract.Decimals(&_FluxAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_FluxAggregator *FluxAggregatorCallerSession) Decimals() (uint8, error) {
	return _FluxAggregator.Contract.Decimals(&_FluxAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() constant returns(string)
func (_FluxAggregator *FluxAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "description")
	return *ret0, err
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() constant returns(string)
func (_FluxAggregator *FluxAggregatorSession) Description() (string, error) {
	return _FluxAggregator.Contract.Description(&_FluxAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() constant returns(string)
func (_FluxAggregator *FluxAggregatorCallerSession) Description() (string, error) {
	return _FluxAggregator.Contract.Description(&_FluxAggregator.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x64efb22b.
//
// Solidity: function getAdmin(address _oracle) constant returns(address)
func (_FluxAggregator *FluxAggregatorCaller) GetAdmin(opts *bind.CallOpts, _oracle common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "getAdmin", _oracle)
	return *ret0, err
}

// GetAdmin is a free data retrieval call binding the contract method 0x64efb22b.
//
// Solidity: function getAdmin(address _oracle) constant returns(address)
func (_FluxAggregator *FluxAggregatorSession) GetAdmin(_oracle common.Address) (common.Address, error) {
	return _FluxAggregator.Contract.GetAdmin(&_FluxAggregator.CallOpts, _oracle)
}

// GetAdmin is a free data retrieval call binding the contract method 0x64efb22b.
//
// Solidity: function getAdmin(address _oracle) constant returns(address)
func (_FluxAggregator *FluxAggregatorCallerSession) GetAdmin(_oracle common.Address) (common.Address, error) {
	return _FluxAggregator.Contract.GetAdmin(&_FluxAggregator.CallOpts, _oracle)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) constant returns(int256)
func (_FluxAggregator *FluxAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "getAnswer", _roundId)
	return *ret0, err
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) constant returns(int256)
func (_FluxAggregator *FluxAggregatorSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _FluxAggregator.Contract.GetAnswer(&_FluxAggregator.CallOpts, _roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) constant returns(int256)
func (_FluxAggregator *FluxAggregatorCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _FluxAggregator.Contract.GetAnswer(&_FluxAggregator.CallOpts, _roundId)
}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() constant returns(address[])
func (_FluxAggregator *FluxAggregatorCaller) GetOracles(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "getOracles")
	return *ret0, err
}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() constant returns(address[])
func (_FluxAggregator *FluxAggregatorSession) GetOracles() ([]common.Address, error) {
	return _FluxAggregator.Contract.GetOracles(&_FluxAggregator.CallOpts)
}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() constant returns(address[])
func (_FluxAggregator *FluxAggregatorCallerSession) GetOracles() ([]common.Address, error) {
	return _FluxAggregator.Contract.GetOracles(&_FluxAggregator.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x0720da52.
//
// Solidity: function getRoundData(uint256 _roundId) constant returns(uint256 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint256 answeredInRound)
func (_FluxAggregator *FluxAggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _FluxAggregator.contract.Call(opts, out, "getRoundData", _roundId)
	return *ret, err
}

// GetRoundData is a free data retrieval call binding the contract method 0x0720da52.
//
// Solidity: function getRoundData(uint256 _roundId) constant returns(uint256 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint256 answeredInRound)
func (_FluxAggregator *FluxAggregatorSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _FluxAggregator.Contract.GetRoundData(&_FluxAggregator.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x0720da52.
//
// Solidity: function getRoundData(uint256 _roundId) constant returns(uint256 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint256 answeredInRound)
func (_FluxAggregator *FluxAggregatorCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _FluxAggregator.Contract.GetRoundData(&_FluxAggregator.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "getTimestamp", _roundId)
	return *ret0, err
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _FluxAggregator.Contract.GetTimestamp(&_FluxAggregator.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _FluxAggregator.Contract.GetTimestamp(&_FluxAggregator.CallOpts, _roundId)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() constant returns(int256)
func (_FluxAggregator *FluxAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "latestAnswer")
	return *ret0, err
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() constant returns(int256)
func (_FluxAggregator *FluxAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _FluxAggregator.Contract.LatestAnswer(&_FluxAggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() constant returns(int256)
func (_FluxAggregator *FluxAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _FluxAggregator.Contract.LatestAnswer(&_FluxAggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "latestRound")
	return *ret0, err
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) LatestRound() (*big.Int, error) {
	return _FluxAggregator.Contract.LatestRound(&_FluxAggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _FluxAggregator.Contract.LatestRound(&_FluxAggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() constant returns(uint256 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint256 answeredInRound)
func (_FluxAggregator *FluxAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _FluxAggregator.contract.Call(opts, out, "latestRoundData")
	return *ret, err
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() constant returns(uint256 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint256 answeredInRound)
func (_FluxAggregator *FluxAggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _FluxAggregator.Contract.LatestRoundData(&_FluxAggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() constant returns(uint256 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint256 answeredInRound)
func (_FluxAggregator *FluxAggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _FluxAggregator.Contract.LatestRoundData(&_FluxAggregator.CallOpts)
}

// LatestSubmission is a free data retrieval call binding the contract method 0xbb07bacd.
//
// Solidity: function latestSubmission(address _oracle) constant returns(int256, uint256)
func (_FluxAggregator *FluxAggregatorCaller) LatestSubmission(opts *bind.CallOpts, _oracle common.Address) (*big.Int, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _FluxAggregator.contract.Call(opts, out, "latestSubmission", _oracle)
	return *ret0, *ret1, err
}

// LatestSubmission is a free data retrieval call binding the contract method 0xbb07bacd.
//
// Solidity: function latestSubmission(address _oracle) constant returns(int256, uint256)
func (_FluxAggregator *FluxAggregatorSession) LatestSubmission(_oracle common.Address) (*big.Int, *big.Int, error) {
	return _FluxAggregator.Contract.LatestSubmission(&_FluxAggregator.CallOpts, _oracle)
}

// LatestSubmission is a free data retrieval call binding the contract method 0xbb07bacd.
//
// Solidity: function latestSubmission(address _oracle) constant returns(int256, uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) LatestSubmission(_oracle common.Address) (*big.Int, *big.Int, error) {
	return _FluxAggregator.Contract.LatestSubmission(&_FluxAggregator.CallOpts, _oracle)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "latestTimestamp")
	return *ret0, err
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _FluxAggregator.Contract.LatestTimestamp(&_FluxAggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _FluxAggregator.Contract.LatestTimestamp(&_FluxAggregator.CallOpts)
}

// LinkToken is a free data retrieval call binding the contract method 0x57970e93.
//
// Solidity: function linkToken() constant returns(address)
func (_FluxAggregator *FluxAggregatorCaller) LinkToken(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "linkToken")
	return *ret0, err
}

// LinkToken is a free data retrieval call binding the contract method 0x57970e93.
//
// Solidity: function linkToken() constant returns(address)
func (_FluxAggregator *FluxAggregatorSession) LinkToken() (common.Address, error) {
	return _FluxAggregator.Contract.LinkToken(&_FluxAggregator.CallOpts)
}

// LinkToken is a free data retrieval call binding the contract method 0x57970e93.
//
// Solidity: function linkToken() constant returns(address)
func (_FluxAggregator *FluxAggregatorCallerSession) LinkToken() (common.Address, error) {
	return _FluxAggregator.Contract.LinkToken(&_FluxAggregator.CallOpts)
}

// MaxSubmissionCount is a free data retrieval call binding the contract method 0x58609e44.
//
// Solidity: function maxSubmissionCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCaller) MaxSubmissionCount(opts *bind.CallOpts) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "maxSubmissionCount")
	return *ret0, err
}

// MaxSubmissionCount is a free data retrieval call binding the contract method 0x58609e44.
//
// Solidity: function maxSubmissionCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorSession) MaxSubmissionCount() (uint32, error) {
	return _FluxAggregator.Contract.MaxSubmissionCount(&_FluxAggregator.CallOpts)
}

// MaxSubmissionCount is a free data retrieval call binding the contract method 0x58609e44.
//
// Solidity: function maxSubmissionCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCallerSession) MaxSubmissionCount() (uint32, error) {
	return _FluxAggregator.Contract.MaxSubmissionCount(&_FluxAggregator.CallOpts)
}

// MinSubmissionCount is a free data retrieval call binding the contract method 0xc9374500.
//
// Solidity: function minSubmissionCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCaller) MinSubmissionCount(opts *bind.CallOpts) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "minSubmissionCount")
	return *ret0, err
}

// MinSubmissionCount is a free data retrieval call binding the contract method 0xc9374500.
//
// Solidity: function minSubmissionCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorSession) MinSubmissionCount() (uint32, error) {
	return _FluxAggregator.Contract.MinSubmissionCount(&_FluxAggregator.CallOpts)
}

// MinSubmissionCount is a free data retrieval call binding the contract method 0xc9374500.
//
// Solidity: function minSubmissionCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCallerSession) MinSubmissionCount() (uint32, error) {
	return _FluxAggregator.Contract.MinSubmissionCount(&_FluxAggregator.CallOpts)
}

// OracleCount is a free data retrieval call binding the contract method 0x613d8fcc.
//
// Solidity: function oracleCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCaller) OracleCount(opts *bind.CallOpts) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "oracleCount")
	return *ret0, err
}

// OracleCount is a free data retrieval call binding the contract method 0x613d8fcc.
//
// Solidity: function oracleCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorSession) OracleCount() (uint32, error) {
	return _FluxAggregator.Contract.OracleCount(&_FluxAggregator.CallOpts)
}

// OracleCount is a free data retrieval call binding the contract method 0x613d8fcc.
//
// Solidity: function oracleCount() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCallerSession) OracleCount() (uint32, error) {
	return _FluxAggregator.Contract.OracleCount(&_FluxAggregator.CallOpts)
}

// OracleRoundState is a free data retrieval call binding the contract method 0x88aa80e7.
//
// Solidity: function oracleRoundState(address _oracle, uint32 _queriedRoundId) constant returns(bool _eligibleToSubmit, uint32 _roundId, int256 _latestSubmission, uint64 _startedAt, uint64 _timeout, uint128 _availableFunds, uint32 _oracleCount, uint128 _paymentAmount)
func (_FluxAggregator *FluxAggregatorCaller) OracleRoundState(opts *bind.CallOpts, _oracle common.Address, _queriedRoundId uint32) (struct {
	EligibleToSubmit bool
	RoundId          uint32
	LatestSubmission *big.Int
	StartedAt        uint64
	Timeout          uint64
	AvailableFunds   *big.Int
	OracleCount      uint32
	PaymentAmount    *big.Int
}, error) {
	ret := new(struct {
		EligibleToSubmit bool
		RoundId          uint32
		LatestSubmission *big.Int
		StartedAt        uint64
		Timeout          uint64
		AvailableFunds   *big.Int
		OracleCount      uint32
		PaymentAmount    *big.Int
	})
	out := ret
	err := _FluxAggregator.contract.Call(opts, out, "oracleRoundState", _oracle, _queriedRoundId)
	return *ret, err
}

// OracleRoundState is a free data retrieval call binding the contract method 0x88aa80e7.
//
// Solidity: function oracleRoundState(address _oracle, uint32 _queriedRoundId) constant returns(bool _eligibleToSubmit, uint32 _roundId, int256 _latestSubmission, uint64 _startedAt, uint64 _timeout, uint128 _availableFunds, uint32 _oracleCount, uint128 _paymentAmount)
func (_FluxAggregator *FluxAggregatorSession) OracleRoundState(_oracle common.Address, _queriedRoundId uint32) (struct {
	EligibleToSubmit bool
	RoundId          uint32
	LatestSubmission *big.Int
	StartedAt        uint64
	Timeout          uint64
	AvailableFunds   *big.Int
	OracleCount      uint32
	PaymentAmount    *big.Int
}, error) {
	return _FluxAggregator.Contract.OracleRoundState(&_FluxAggregator.CallOpts, _oracle, _queriedRoundId)
}

// OracleRoundState is a free data retrieval call binding the contract method 0x88aa80e7.
//
// Solidity: function oracleRoundState(address _oracle, uint32 _queriedRoundId) constant returns(bool _eligibleToSubmit, uint32 _roundId, int256 _latestSubmission, uint64 _startedAt, uint64 _timeout, uint128 _availableFunds, uint32 _oracleCount, uint128 _paymentAmount)
func (_FluxAggregator *FluxAggregatorCallerSession) OracleRoundState(_oracle common.Address, _queriedRoundId uint32) (struct {
	EligibleToSubmit bool
	RoundId          uint32
	LatestSubmission *big.Int
	StartedAt        uint64
	Timeout          uint64
	AvailableFunds   *big.Int
	OracleCount      uint32
	PaymentAmount    *big.Int
}, error) {
	return _FluxAggregator.Contract.OracleRoundState(&_FluxAggregator.CallOpts, _oracle, _queriedRoundId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_FluxAggregator *FluxAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_FluxAggregator *FluxAggregatorSession) Owner() (common.Address, error) {
	return _FluxAggregator.Contract.Owner(&_FluxAggregator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_FluxAggregator *FluxAggregatorCallerSession) Owner() (common.Address, error) {
	return _FluxAggregator.Contract.Owner(&_FluxAggregator.CallOpts)
}

// PaymentAmount is a free data retrieval call binding the contract method 0xc35905c6.
//
// Solidity: function paymentAmount() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorCaller) PaymentAmount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "paymentAmount")
	return *ret0, err
}

// PaymentAmount is a free data retrieval call binding the contract method 0xc35905c6.
//
// Solidity: function paymentAmount() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorSession) PaymentAmount() (*big.Int, error) {
	return _FluxAggregator.Contract.PaymentAmount(&_FluxAggregator.CallOpts)
}

// PaymentAmount is a free data retrieval call binding the contract method 0xc35905c6.
//
// Solidity: function paymentAmount() constant returns(uint128)
func (_FluxAggregator *FluxAggregatorCallerSession) PaymentAmount() (*big.Int, error) {
	return _FluxAggregator.Contract.PaymentAmount(&_FluxAggregator.CallOpts)
}

// ReportingRound is a free data retrieval call binding the contract method 0x6fb4bb4e.
//
// Solidity: function reportingRound() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) ReportingRound(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "reportingRound")
	return *ret0, err
}

// ReportingRound is a free data retrieval call binding the contract method 0x6fb4bb4e.
//
// Solidity: function reportingRound() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) ReportingRound() (*big.Int, error) {
	return _FluxAggregator.Contract.ReportingRound(&_FluxAggregator.CallOpts)
}

// ReportingRound is a free data retrieval call binding the contract method 0x6fb4bb4e.
//
// Solidity: function reportingRound() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) ReportingRound() (*big.Int, error) {
	return _FluxAggregator.Contract.ReportingRound(&_FluxAggregator.CallOpts)
}

// RestartDelay is a free data retrieval call binding the contract method 0x357ebb02.
//
// Solidity: function restartDelay() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCaller) RestartDelay(opts *bind.CallOpts) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "restartDelay")
	return *ret0, err
}

// RestartDelay is a free data retrieval call binding the contract method 0x357ebb02.
//
// Solidity: function restartDelay() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorSession) RestartDelay() (uint32, error) {
	return _FluxAggregator.Contract.RestartDelay(&_FluxAggregator.CallOpts)
}

// RestartDelay is a free data retrieval call binding the contract method 0x357ebb02.
//
// Solidity: function restartDelay() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCallerSession) RestartDelay() (uint32, error) {
	return _FluxAggregator.Contract.RestartDelay(&_FluxAggregator.CallOpts)
}

// Timeout is a free data retrieval call binding the contract method 0x70dea79a.
//
// Solidity: function timeout() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCaller) Timeout(opts *bind.CallOpts) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "timeout")
	return *ret0, err
}

// Timeout is a free data retrieval call binding the contract method 0x70dea79a.
//
// Solidity: function timeout() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorSession) Timeout() (uint32, error) {
	return _FluxAggregator.Contract.Timeout(&_FluxAggregator.CallOpts)
}

// Timeout is a free data retrieval call binding the contract method 0x70dea79a.
//
// Solidity: function timeout() constant returns(uint32)
func (_FluxAggregator *FluxAggregatorCallerSession) Timeout() (uint32, error) {
	return _FluxAggregator.Contract.Timeout(&_FluxAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) Version() (*big.Int, error) {
	return _FluxAggregator.Contract.Version(&_FluxAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) Version() (*big.Int, error) {
	return _FluxAggregator.Contract.Version(&_FluxAggregator.CallOpts)
}

// WithdrawablePayment is a free data retrieval call binding the contract method 0xe2e40317.
//
// Solidity: function withdrawablePayment(address _oracle) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCaller) WithdrawablePayment(opts *bind.CallOpts, _oracle common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _FluxAggregator.contract.Call(opts, out, "withdrawablePayment", _oracle)
	return *ret0, err
}

// WithdrawablePayment is a free data retrieval call binding the contract method 0xe2e40317.
//
// Solidity: function withdrawablePayment(address _oracle) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorSession) WithdrawablePayment(_oracle common.Address) (*big.Int, error) {
	return _FluxAggregator.Contract.WithdrawablePayment(&_FluxAggregator.CallOpts, _oracle)
}

// WithdrawablePayment is a free data retrieval call binding the contract method 0xe2e40317.
//
// Solidity: function withdrawablePayment(address _oracle) constant returns(uint256)
func (_FluxAggregator *FluxAggregatorCallerSession) WithdrawablePayment(_oracle common.Address) (*big.Int, error) {
	return _FluxAggregator.Contract.WithdrawablePayment(&_FluxAggregator.CallOpts, _oracle)
}

// AcceptAdmin is a paid mutator transaction binding the contract method 0x628806ef.
//
// Solidity: function acceptAdmin(address _oracle) returns()
func (_FluxAggregator *FluxAggregatorTransactor) AcceptAdmin(opts *bind.TransactOpts, _oracle common.Address) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "acceptAdmin", _oracle)
}

// AcceptAdmin is a paid mutator transaction binding the contract method 0x628806ef.
//
// Solidity: function acceptAdmin(address _oracle) returns()
func (_FluxAggregator *FluxAggregatorSession) AcceptAdmin(_oracle common.Address) (*types.Transaction, error) {
	return _FluxAggregator.Contract.AcceptAdmin(&_FluxAggregator.TransactOpts, _oracle)
}

// AcceptAdmin is a paid mutator transaction binding the contract method 0x628806ef.
//
// Solidity: function acceptAdmin(address _oracle) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) AcceptAdmin(_oracle common.Address) (*types.Transaction, error) {
	return _FluxAggregator.Contract.AcceptAdmin(&_FluxAggregator.TransactOpts, _oracle)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_FluxAggregator *FluxAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_FluxAggregator *FluxAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FluxAggregator.Contract.AcceptOwnership(&_FluxAggregator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FluxAggregator.Contract.AcceptOwnership(&_FluxAggregator.TransactOpts)
}

// AddOracles is a paid mutator transaction binding the contract method 0xbbf0b7e9.
//
// Solidity: function addOracles(address[] _oracles, address[] _admins, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorTransactor) AddOracles(opts *bind.TransactOpts, _oracles []common.Address, _admins []common.Address, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "addOracles", _oracles, _admins, _minSubmissions, _maxSubmissions, _restartDelay)
}

// AddOracles is a paid mutator transaction binding the contract method 0xbbf0b7e9.
//
// Solidity: function addOracles(address[] _oracles, address[] _admins, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorSession) AddOracles(_oracles []common.Address, _admins []common.Address, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.AddOracles(&_FluxAggregator.TransactOpts, _oracles, _admins, _minSubmissions, _maxSubmissions, _restartDelay)
}

// AddOracles is a paid mutator transaction binding the contract method 0xbbf0b7e9.
//
// Solidity: function addOracles(address[] _oracles, address[] _admins, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) AddOracles(_oracles []common.Address, _admins []common.Address, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.AddOracles(&_FluxAggregator.TransactOpts, _oracles, _admins, _minSubmissions, _maxSubmissions, _restartDelay)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 , bytes _data) returns()
func (_FluxAggregator *FluxAggregatorTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int, _data []byte) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "onTokenTransfer", arg0, arg1, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 , bytes _data) returns()
func (_FluxAggregator *FluxAggregatorSession) OnTokenTransfer(arg0 common.Address, arg1 *big.Int, _data []byte) (*types.Transaction, error) {
	return _FluxAggregator.Contract.OnTokenTransfer(&_FluxAggregator.TransactOpts, arg0, arg1, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 , bytes _data) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) OnTokenTransfer(arg0 common.Address, arg1 *big.Int, _data []byte) (*types.Transaction, error) {
	return _FluxAggregator.Contract.OnTokenTransfer(&_FluxAggregator.TransactOpts, arg0, arg1, _data)
}

// RemoveOracles is a paid mutator transaction binding the contract method 0xebf8571c.
//
// Solidity: function removeOracles(address[] _oracles, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorTransactor) RemoveOracles(opts *bind.TransactOpts, _oracles []common.Address, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "removeOracles", _oracles, _minSubmissions, _maxSubmissions, _restartDelay)
}

// RemoveOracles is a paid mutator transaction binding the contract method 0xebf8571c.
//
// Solidity: function removeOracles(address[] _oracles, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorSession) RemoveOracles(_oracles []common.Address, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.RemoveOracles(&_FluxAggregator.TransactOpts, _oracles, _minSubmissions, _maxSubmissions, _restartDelay)
}

// RemoveOracles is a paid mutator transaction binding the contract method 0xebf8571c.
//
// Solidity: function removeOracles(address[] _oracles, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) RemoveOracles(_oracles []common.Address, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.RemoveOracles(&_FluxAggregator.TransactOpts, _oracles, _minSubmissions, _maxSubmissions, _restartDelay)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns()
func (_FluxAggregator *FluxAggregatorTransactor) RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "requestNewRound")
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns()
func (_FluxAggregator *FluxAggregatorSession) RequestNewRound() (*types.Transaction, error) {
	return _FluxAggregator.Contract.RequestNewRound(&_FluxAggregator.TransactOpts)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) RequestNewRound() (*types.Transaction, error) {
	return _FluxAggregator.Contract.RequestNewRound(&_FluxAggregator.TransactOpts)
}

// SetRequesterPermissions is a paid mutator transaction binding the contract method 0x20ed0275.
//
// Solidity: function setRequesterPermissions(address _requester, bool _authorized, uint32 _delay) returns()
func (_FluxAggregator *FluxAggregatorTransactor) SetRequesterPermissions(opts *bind.TransactOpts, _requester common.Address, _authorized bool, _delay uint32) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "setRequesterPermissions", _requester, _authorized, _delay)
}

// SetRequesterPermissions is a paid mutator transaction binding the contract method 0x20ed0275.
//
// Solidity: function setRequesterPermissions(address _requester, bool _authorized, uint32 _delay) returns()
func (_FluxAggregator *FluxAggregatorSession) SetRequesterPermissions(_requester common.Address, _authorized bool, _delay uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.SetRequesterPermissions(&_FluxAggregator.TransactOpts, _requester, _authorized, _delay)
}

// SetRequesterPermissions is a paid mutator transaction binding the contract method 0x20ed0275.
//
// Solidity: function setRequesterPermissions(address _requester, bool _authorized, uint32 _delay) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) SetRequesterPermissions(_requester common.Address, _authorized bool, _delay uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.SetRequesterPermissions(&_FluxAggregator.TransactOpts, _requester, _authorized, _delay)
}

// Submit is a paid mutator transaction binding the contract method 0x202ee0ed.
//
// Solidity: function submit(uint256 _roundId, int256 _submission) returns()
func (_FluxAggregator *FluxAggregatorTransactor) Submit(opts *bind.TransactOpts, _roundId *big.Int, _submission *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "submit", _roundId, _submission)
}

// Submit is a paid mutator transaction binding the contract method 0x202ee0ed.
//
// Solidity: function submit(uint256 _roundId, int256 _submission) returns()
func (_FluxAggregator *FluxAggregatorSession) Submit(_roundId *big.Int, _submission *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.Contract.Submit(&_FluxAggregator.TransactOpts, _roundId, _submission)
}

// Submit is a paid mutator transaction binding the contract method 0x202ee0ed.
//
// Solidity: function submit(uint256 _roundId, int256 _submission) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) Submit(_roundId *big.Int, _submission *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.Contract.Submit(&_FluxAggregator.TransactOpts, _roundId, _submission)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0xe9ee6eeb.
//
// Solidity: function transferAdmin(address _oracle, address _newAdmin) returns()
func (_FluxAggregator *FluxAggregatorTransactor) TransferAdmin(opts *bind.TransactOpts, _oracle common.Address, _newAdmin common.Address) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "transferAdmin", _oracle, _newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0xe9ee6eeb.
//
// Solidity: function transferAdmin(address _oracle, address _newAdmin) returns()
func (_FluxAggregator *FluxAggregatorSession) TransferAdmin(_oracle common.Address, _newAdmin common.Address) (*types.Transaction, error) {
	return _FluxAggregator.Contract.TransferAdmin(&_FluxAggregator.TransactOpts, _oracle, _newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0xe9ee6eeb.
//
// Solidity: function transferAdmin(address _oracle, address _newAdmin) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) TransferAdmin(_oracle common.Address, _newAdmin common.Address) (*types.Transaction, error) {
	return _FluxAggregator.Contract.TransferAdmin(&_FluxAggregator.TransactOpts, _oracle, _newAdmin)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_FluxAggregator *FluxAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_FluxAggregator *FluxAggregatorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _FluxAggregator.Contract.TransferOwnership(&_FluxAggregator.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _FluxAggregator.Contract.TransferOwnership(&_FluxAggregator.TransactOpts, _to)
}

// UpdateAvailableFunds is a paid mutator transaction binding the contract method 0x4f8fc3b5.
//
// Solidity: function updateAvailableFunds() returns()
func (_FluxAggregator *FluxAggregatorTransactor) UpdateAvailableFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "updateAvailableFunds")
}

// UpdateAvailableFunds is a paid mutator transaction binding the contract method 0x4f8fc3b5.
//
// Solidity: function updateAvailableFunds() returns()
func (_FluxAggregator *FluxAggregatorSession) UpdateAvailableFunds() (*types.Transaction, error) {
	return _FluxAggregator.Contract.UpdateAvailableFunds(&_FluxAggregator.TransactOpts)
}

// UpdateAvailableFunds is a paid mutator transaction binding the contract method 0x4f8fc3b5.
//
// Solidity: function updateAvailableFunds() returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) UpdateAvailableFunds() (*types.Transaction, error) {
	return _FluxAggregator.Contract.UpdateAvailableFunds(&_FluxAggregator.TransactOpts)
}

// UpdateFutureRounds is a paid mutator transaction binding the contract method 0x38aa4c72.
//
// Solidity: function updateFutureRounds(uint128 _paymentAmount, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay, uint32 _timeout) returns()
func (_FluxAggregator *FluxAggregatorTransactor) UpdateFutureRounds(opts *bind.TransactOpts, _paymentAmount *big.Int, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32, _timeout uint32) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "updateFutureRounds", _paymentAmount, _minSubmissions, _maxSubmissions, _restartDelay, _timeout)
}

// UpdateFutureRounds is a paid mutator transaction binding the contract method 0x38aa4c72.
//
// Solidity: function updateFutureRounds(uint128 _paymentAmount, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay, uint32 _timeout) returns()
func (_FluxAggregator *FluxAggregatorSession) UpdateFutureRounds(_paymentAmount *big.Int, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32, _timeout uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.UpdateFutureRounds(&_FluxAggregator.TransactOpts, _paymentAmount, _minSubmissions, _maxSubmissions, _restartDelay, _timeout)
}

// UpdateFutureRounds is a paid mutator transaction binding the contract method 0x38aa4c72.
//
// Solidity: function updateFutureRounds(uint128 _paymentAmount, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay, uint32 _timeout) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) UpdateFutureRounds(_paymentAmount *big.Int, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32, _timeout uint32) (*types.Transaction, error) {
	return _FluxAggregator.Contract.UpdateFutureRounds(&_FluxAggregator.TransactOpts, _paymentAmount, _minSubmissions, _maxSubmissions, _restartDelay, _timeout)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_FluxAggregator *FluxAggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "withdrawFunds", _recipient, _amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_FluxAggregator *FluxAggregatorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.Contract.WithdrawFunds(&_FluxAggregator.TransactOpts, _recipient, _amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.Contract.WithdrawFunds(&_FluxAggregator.TransactOpts, _recipient, _amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x3d3d7714.
//
// Solidity: function withdrawPayment(address _oracle, address _recipient, uint256 _amount) returns()
func (_FluxAggregator *FluxAggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, _oracle common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.contract.Transact(opts, "withdrawPayment", _oracle, _recipient, _amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x3d3d7714.
//
// Solidity: function withdrawPayment(address _oracle, address _recipient, uint256 _amount) returns()
func (_FluxAggregator *FluxAggregatorSession) WithdrawPayment(_oracle common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.Contract.WithdrawPayment(&_FluxAggregator.TransactOpts, _oracle, _recipient, _amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x3d3d7714.
//
// Solidity: function withdrawPayment(address _oracle, address _recipient, uint256 _amount) returns()
func (_FluxAggregator *FluxAggregatorTransactorSession) WithdrawPayment(_oracle common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FluxAggregator.Contract.WithdrawPayment(&_FluxAggregator.TransactOpts, _oracle, _recipient, _amount)
}

// FluxAggregatorAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the FluxAggregator contract.
type FluxAggregatorAnswerUpdatedIterator struct {
	Event *FluxAggregatorAnswerUpdated // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorAnswerUpdated)
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
		it.Event = new(FluxAggregatorAnswerUpdated)
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
func (it *FluxAggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorAnswerUpdated represents a AnswerUpdated event raised by the FluxAggregator contract.
type FluxAggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp)
func (_FluxAggregator *FluxAggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*FluxAggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorAnswerUpdatedIterator{contract: _FluxAggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp)
func (_FluxAggregator *FluxAggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *FluxAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorAnswerUpdated)
				if err := _FluxAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp)
func (_FluxAggregator *FluxAggregatorFilterer) ParseAnswerUpdated(log types.Log) (*FluxAggregatorAnswerUpdated, error) {
	event := new(FluxAggregatorAnswerUpdated)
	if err := _FluxAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorAvailableFundsUpdatedIterator is returned from FilterAvailableFundsUpdated and is used to iterate over the raw logs and unpacked data for AvailableFundsUpdated events raised by the FluxAggregator contract.
type FluxAggregatorAvailableFundsUpdatedIterator struct {
	Event *FluxAggregatorAvailableFundsUpdated // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorAvailableFundsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorAvailableFundsUpdated)
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
		it.Event = new(FluxAggregatorAvailableFundsUpdated)
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
func (it *FluxAggregatorAvailableFundsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorAvailableFundsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorAvailableFundsUpdated represents a AvailableFundsUpdated event raised by the FluxAggregator contract.
type FluxAggregatorAvailableFundsUpdated struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAvailableFundsUpdated is a free log retrieval operation binding the contract event 0xfe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f.
//
// Solidity: event AvailableFundsUpdated(uint256 indexed amount)
func (_FluxAggregator *FluxAggregatorFilterer) FilterAvailableFundsUpdated(opts *bind.FilterOpts, amount []*big.Int) (*FluxAggregatorAvailableFundsUpdatedIterator, error) {

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "AvailableFundsUpdated", amountRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorAvailableFundsUpdatedIterator{contract: _FluxAggregator.contract, event: "AvailableFundsUpdated", logs: logs, sub: sub}, nil
}

// WatchAvailableFundsUpdated is a free log subscription operation binding the contract event 0xfe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f.
//
// Solidity: event AvailableFundsUpdated(uint256 indexed amount)
func (_FluxAggregator *FluxAggregatorFilterer) WatchAvailableFundsUpdated(opts *bind.WatchOpts, sink chan<- *FluxAggregatorAvailableFundsUpdated, amount []*big.Int) (event.Subscription, error) {

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "AvailableFundsUpdated", amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorAvailableFundsUpdated)
				if err := _FluxAggregator.contract.UnpackLog(event, "AvailableFundsUpdated", log); err != nil {
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

// ParseAvailableFundsUpdated is a log parse operation binding the contract event 0xfe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f.
//
// Solidity: event AvailableFundsUpdated(uint256 indexed amount)
func (_FluxAggregator *FluxAggregatorFilterer) ParseAvailableFundsUpdated(log types.Log) (*FluxAggregatorAvailableFundsUpdated, error) {
	event := new(FluxAggregatorAvailableFundsUpdated)
	if err := _FluxAggregator.contract.UnpackLog(event, "AvailableFundsUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the FluxAggregator contract.
type FluxAggregatorNewRoundIterator struct {
	Event *FluxAggregatorNewRound // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorNewRound)
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
		it.Event = new(FluxAggregatorNewRound)
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
func (it *FluxAggregatorNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorNewRound represents a NewRound event raised by the FluxAggregator contract.
type FluxAggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_FluxAggregator *FluxAggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*FluxAggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorNewRoundIterator{contract: _FluxAggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_FluxAggregator *FluxAggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *FluxAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorNewRound)
				if err := _FluxAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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
func (_FluxAggregator *FluxAggregatorFilterer) ParseNewRound(log types.Log) (*FluxAggregatorNewRound, error) {
	event := new(FluxAggregatorNewRound)
	if err := _FluxAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorOracleAdminUpdateRequestedIterator is returned from FilterOracleAdminUpdateRequested and is used to iterate over the raw logs and unpacked data for OracleAdminUpdateRequested events raised by the FluxAggregator contract.
type FluxAggregatorOracleAdminUpdateRequestedIterator struct {
	Event *FluxAggregatorOracleAdminUpdateRequested // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorOracleAdminUpdateRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorOracleAdminUpdateRequested)
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
		it.Event = new(FluxAggregatorOracleAdminUpdateRequested)
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
func (it *FluxAggregatorOracleAdminUpdateRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorOracleAdminUpdateRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorOracleAdminUpdateRequested represents a OracleAdminUpdateRequested event raised by the FluxAggregator contract.
type FluxAggregatorOracleAdminUpdateRequested struct {
	Oracle   common.Address
	Admin    common.Address
	NewAdmin common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOracleAdminUpdateRequested is a free log retrieval operation binding the contract event 0xb79bf2e89c2d70dde91d2991fb1ea69b7e478061ad7c04ed5b02b96bc52b8104.
//
// Solidity: event OracleAdminUpdateRequested(address indexed oracle, address admin, address newAdmin)
func (_FluxAggregator *FluxAggregatorFilterer) FilterOracleAdminUpdateRequested(opts *bind.FilterOpts, oracle []common.Address) (*FluxAggregatorOracleAdminUpdateRequestedIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "OracleAdminUpdateRequested", oracleRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorOracleAdminUpdateRequestedIterator{contract: _FluxAggregator.contract, event: "OracleAdminUpdateRequested", logs: logs, sub: sub}, nil
}

// WatchOracleAdminUpdateRequested is a free log subscription operation binding the contract event 0xb79bf2e89c2d70dde91d2991fb1ea69b7e478061ad7c04ed5b02b96bc52b8104.
//
// Solidity: event OracleAdminUpdateRequested(address indexed oracle, address admin, address newAdmin)
func (_FluxAggregator *FluxAggregatorFilterer) WatchOracleAdminUpdateRequested(opts *bind.WatchOpts, sink chan<- *FluxAggregatorOracleAdminUpdateRequested, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "OracleAdminUpdateRequested", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorOracleAdminUpdateRequested)
				if err := _FluxAggregator.contract.UnpackLog(event, "OracleAdminUpdateRequested", log); err != nil {
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

// ParseOracleAdminUpdateRequested is a log parse operation binding the contract event 0xb79bf2e89c2d70dde91d2991fb1ea69b7e478061ad7c04ed5b02b96bc52b8104.
//
// Solidity: event OracleAdminUpdateRequested(address indexed oracle, address admin, address newAdmin)
func (_FluxAggregator *FluxAggregatorFilterer) ParseOracleAdminUpdateRequested(log types.Log) (*FluxAggregatorOracleAdminUpdateRequested, error) {
	event := new(FluxAggregatorOracleAdminUpdateRequested)
	if err := _FluxAggregator.contract.UnpackLog(event, "OracleAdminUpdateRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorOracleAdminUpdatedIterator is returned from FilterOracleAdminUpdated and is used to iterate over the raw logs and unpacked data for OracleAdminUpdated events raised by the FluxAggregator contract.
type FluxAggregatorOracleAdminUpdatedIterator struct {
	Event *FluxAggregatorOracleAdminUpdated // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorOracleAdminUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorOracleAdminUpdated)
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
		it.Event = new(FluxAggregatorOracleAdminUpdated)
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
func (it *FluxAggregatorOracleAdminUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorOracleAdminUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorOracleAdminUpdated represents a OracleAdminUpdated event raised by the FluxAggregator contract.
type FluxAggregatorOracleAdminUpdated struct {
	Oracle   common.Address
	NewAdmin common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOracleAdminUpdated is a free log retrieval operation binding the contract event 0x0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe904.
//
// Solidity: event OracleAdminUpdated(address indexed oracle, address indexed newAdmin)
func (_FluxAggregator *FluxAggregatorFilterer) FilterOracleAdminUpdated(opts *bind.FilterOpts, oracle []common.Address, newAdmin []common.Address) (*FluxAggregatorOracleAdminUpdatedIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "OracleAdminUpdated", oracleRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorOracleAdminUpdatedIterator{contract: _FluxAggregator.contract, event: "OracleAdminUpdated", logs: logs, sub: sub}, nil
}

// WatchOracleAdminUpdated is a free log subscription operation binding the contract event 0x0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe904.
//
// Solidity: event OracleAdminUpdated(address indexed oracle, address indexed newAdmin)
func (_FluxAggregator *FluxAggregatorFilterer) WatchOracleAdminUpdated(opts *bind.WatchOpts, sink chan<- *FluxAggregatorOracleAdminUpdated, oracle []common.Address, newAdmin []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "OracleAdminUpdated", oracleRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorOracleAdminUpdated)
				if err := _FluxAggregator.contract.UnpackLog(event, "OracleAdminUpdated", log); err != nil {
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

// ParseOracleAdminUpdated is a log parse operation binding the contract event 0x0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe904.
//
// Solidity: event OracleAdminUpdated(address indexed oracle, address indexed newAdmin)
func (_FluxAggregator *FluxAggregatorFilterer) ParseOracleAdminUpdated(log types.Log) (*FluxAggregatorOracleAdminUpdated, error) {
	event := new(FluxAggregatorOracleAdminUpdated)
	if err := _FluxAggregator.contract.UnpackLog(event, "OracleAdminUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorOraclePermissionsUpdatedIterator is returned from FilterOraclePermissionsUpdated and is used to iterate over the raw logs and unpacked data for OraclePermissionsUpdated events raised by the FluxAggregator contract.
type FluxAggregatorOraclePermissionsUpdatedIterator struct {
	Event *FluxAggregatorOraclePermissionsUpdated // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorOraclePermissionsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorOraclePermissionsUpdated)
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
		it.Event = new(FluxAggregatorOraclePermissionsUpdated)
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
func (it *FluxAggregatorOraclePermissionsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorOraclePermissionsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorOraclePermissionsUpdated represents a OraclePermissionsUpdated event raised by the FluxAggregator contract.
type FluxAggregatorOraclePermissionsUpdated struct {
	Oracle      common.Address
	Whitelisted bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOraclePermissionsUpdated is a free log retrieval operation binding the contract event 0x18dd09695e4fbdae8d1a5edb11221eb04564269c29a089b9753a6535c54ba92e.
//
// Solidity: event OraclePermissionsUpdated(address indexed oracle, bool indexed whitelisted)
func (_FluxAggregator *FluxAggregatorFilterer) FilterOraclePermissionsUpdated(opts *bind.FilterOpts, oracle []common.Address, whitelisted []bool) (*FluxAggregatorOraclePermissionsUpdatedIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var whitelistedRule []interface{}
	for _, whitelistedItem := range whitelisted {
		whitelistedRule = append(whitelistedRule, whitelistedItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "OraclePermissionsUpdated", oracleRule, whitelistedRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorOraclePermissionsUpdatedIterator{contract: _FluxAggregator.contract, event: "OraclePermissionsUpdated", logs: logs, sub: sub}, nil
}

// WatchOraclePermissionsUpdated is a free log subscription operation binding the contract event 0x18dd09695e4fbdae8d1a5edb11221eb04564269c29a089b9753a6535c54ba92e.
//
// Solidity: event OraclePermissionsUpdated(address indexed oracle, bool indexed whitelisted)
func (_FluxAggregator *FluxAggregatorFilterer) WatchOraclePermissionsUpdated(opts *bind.WatchOpts, sink chan<- *FluxAggregatorOraclePermissionsUpdated, oracle []common.Address, whitelisted []bool) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var whitelistedRule []interface{}
	for _, whitelistedItem := range whitelisted {
		whitelistedRule = append(whitelistedRule, whitelistedItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "OraclePermissionsUpdated", oracleRule, whitelistedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorOraclePermissionsUpdated)
				if err := _FluxAggregator.contract.UnpackLog(event, "OraclePermissionsUpdated", log); err != nil {
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

// ParseOraclePermissionsUpdated is a log parse operation binding the contract event 0x18dd09695e4fbdae8d1a5edb11221eb04564269c29a089b9753a6535c54ba92e.
//
// Solidity: event OraclePermissionsUpdated(address indexed oracle, bool indexed whitelisted)
func (_FluxAggregator *FluxAggregatorFilterer) ParseOraclePermissionsUpdated(log types.Log) (*FluxAggregatorOraclePermissionsUpdated, error) {
	event := new(FluxAggregatorOraclePermissionsUpdated)
	if err := _FluxAggregator.contract.UnpackLog(event, "OraclePermissionsUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the FluxAggregator contract.
type FluxAggregatorOwnershipTransferRequestedIterator struct {
	Event *FluxAggregatorOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorOwnershipTransferRequested)
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
		it.Event = new(FluxAggregatorOwnershipTransferRequested)
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
func (it *FluxAggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the FluxAggregator contract.
type FluxAggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_FluxAggregator *FluxAggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FluxAggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorOwnershipTransferRequestedIterator{contract: _FluxAggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_FluxAggregator *FluxAggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FluxAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorOwnershipTransferRequested)
				if err := _FluxAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_FluxAggregator *FluxAggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*FluxAggregatorOwnershipTransferRequested, error) {
	event := new(FluxAggregatorOwnershipTransferRequested)
	if err := _FluxAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FluxAggregator contract.
type FluxAggregatorOwnershipTransferredIterator struct {
	Event *FluxAggregatorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorOwnershipTransferred)
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
		it.Event = new(FluxAggregatorOwnershipTransferred)
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
func (it *FluxAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorOwnershipTransferred represents a OwnershipTransferred event raised by the FluxAggregator contract.
type FluxAggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_FluxAggregator *FluxAggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FluxAggregatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorOwnershipTransferredIterator{contract: _FluxAggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_FluxAggregator *FluxAggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FluxAggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorOwnershipTransferred)
				if err := _FluxAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_FluxAggregator *FluxAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*FluxAggregatorOwnershipTransferred, error) {
	event := new(FluxAggregatorOwnershipTransferred)
	if err := _FluxAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorRequesterPermissionsSetIterator is returned from FilterRequesterPermissionsSet and is used to iterate over the raw logs and unpacked data for RequesterPermissionsSet events raised by the FluxAggregator contract.
type FluxAggregatorRequesterPermissionsSetIterator struct {
	Event *FluxAggregatorRequesterPermissionsSet // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorRequesterPermissionsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorRequesterPermissionsSet)
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
		it.Event = new(FluxAggregatorRequesterPermissionsSet)
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
func (it *FluxAggregatorRequesterPermissionsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorRequesterPermissionsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorRequesterPermissionsSet represents a RequesterPermissionsSet event raised by the FluxAggregator contract.
type FluxAggregatorRequesterPermissionsSet struct {
	Requester  common.Address
	Authorized bool
	Delay      uint32
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRequesterPermissionsSet is a free log retrieval operation binding the contract event 0xc3df5a754e002718f2e10804b99e6605e7c701d95cec9552c7680ca2b6f2820a.
//
// Solidity: event RequesterPermissionsSet(address indexed requester, bool authorized, uint32 delay)
func (_FluxAggregator *FluxAggregatorFilterer) FilterRequesterPermissionsSet(opts *bind.FilterOpts, requester []common.Address) (*FluxAggregatorRequesterPermissionsSetIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "RequesterPermissionsSet", requesterRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorRequesterPermissionsSetIterator{contract: _FluxAggregator.contract, event: "RequesterPermissionsSet", logs: logs, sub: sub}, nil
}

// WatchRequesterPermissionsSet is a free log subscription operation binding the contract event 0xc3df5a754e002718f2e10804b99e6605e7c701d95cec9552c7680ca2b6f2820a.
//
// Solidity: event RequesterPermissionsSet(address indexed requester, bool authorized, uint32 delay)
func (_FluxAggregator *FluxAggregatorFilterer) WatchRequesterPermissionsSet(opts *bind.WatchOpts, sink chan<- *FluxAggregatorRequesterPermissionsSet, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "RequesterPermissionsSet", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorRequesterPermissionsSet)
				if err := _FluxAggregator.contract.UnpackLog(event, "RequesterPermissionsSet", log); err != nil {
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

// ParseRequesterPermissionsSet is a log parse operation binding the contract event 0xc3df5a754e002718f2e10804b99e6605e7c701d95cec9552c7680ca2b6f2820a.
//
// Solidity: event RequesterPermissionsSet(address indexed requester, bool authorized, uint32 delay)
func (_FluxAggregator *FluxAggregatorFilterer) ParseRequesterPermissionsSet(log types.Log) (*FluxAggregatorRequesterPermissionsSet, error) {
	event := new(FluxAggregatorRequesterPermissionsSet)
	if err := _FluxAggregator.contract.UnpackLog(event, "RequesterPermissionsSet", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorRoundDetailsUpdatedIterator is returned from FilterRoundDetailsUpdated and is used to iterate over the raw logs and unpacked data for RoundDetailsUpdated events raised by the FluxAggregator contract.
type FluxAggregatorRoundDetailsUpdatedIterator struct {
	Event *FluxAggregatorRoundDetailsUpdated // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorRoundDetailsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorRoundDetailsUpdated)
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
		it.Event = new(FluxAggregatorRoundDetailsUpdated)
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
func (it *FluxAggregatorRoundDetailsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorRoundDetailsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorRoundDetailsUpdated represents a RoundDetailsUpdated event raised by the FluxAggregator contract.
type FluxAggregatorRoundDetailsUpdated struct {
	PaymentAmount      *big.Int
	MinSubmissionCount uint32
	MaxSubmissionCount uint32
	RestartDelay       uint32
	Timeout            uint32
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterRoundDetailsUpdated is a free log retrieval operation binding the contract event 0x56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f.
//
// Solidity: event RoundDetailsUpdated(uint128 indexed paymentAmount, uint32 indexed minSubmissionCount, uint32 indexed maxSubmissionCount, uint32 restartDelay, uint32 timeout)
func (_FluxAggregator *FluxAggregatorFilterer) FilterRoundDetailsUpdated(opts *bind.FilterOpts, paymentAmount []*big.Int, minSubmissionCount []uint32, maxSubmissionCount []uint32) (*FluxAggregatorRoundDetailsUpdatedIterator, error) {

	var paymentAmountRule []interface{}
	for _, paymentAmountItem := range paymentAmount {
		paymentAmountRule = append(paymentAmountRule, paymentAmountItem)
	}
	var minSubmissionCountRule []interface{}
	for _, minSubmissionCountItem := range minSubmissionCount {
		minSubmissionCountRule = append(minSubmissionCountRule, minSubmissionCountItem)
	}
	var maxSubmissionCountRule []interface{}
	for _, maxSubmissionCountItem := range maxSubmissionCount {
		maxSubmissionCountRule = append(maxSubmissionCountRule, maxSubmissionCountItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "RoundDetailsUpdated", paymentAmountRule, minSubmissionCountRule, maxSubmissionCountRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorRoundDetailsUpdatedIterator{contract: _FluxAggregator.contract, event: "RoundDetailsUpdated", logs: logs, sub: sub}, nil
}

// WatchRoundDetailsUpdated is a free log subscription operation binding the contract event 0x56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f.
//
// Solidity: event RoundDetailsUpdated(uint128 indexed paymentAmount, uint32 indexed minSubmissionCount, uint32 indexed maxSubmissionCount, uint32 restartDelay, uint32 timeout)
func (_FluxAggregator *FluxAggregatorFilterer) WatchRoundDetailsUpdated(opts *bind.WatchOpts, sink chan<- *FluxAggregatorRoundDetailsUpdated, paymentAmount []*big.Int, minSubmissionCount []uint32, maxSubmissionCount []uint32) (event.Subscription, error) {

	var paymentAmountRule []interface{}
	for _, paymentAmountItem := range paymentAmount {
		paymentAmountRule = append(paymentAmountRule, paymentAmountItem)
	}
	var minSubmissionCountRule []interface{}
	for _, minSubmissionCountItem := range minSubmissionCount {
		minSubmissionCountRule = append(minSubmissionCountRule, minSubmissionCountItem)
	}
	var maxSubmissionCountRule []interface{}
	for _, maxSubmissionCountItem := range maxSubmissionCount {
		maxSubmissionCountRule = append(maxSubmissionCountRule, maxSubmissionCountItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "RoundDetailsUpdated", paymentAmountRule, minSubmissionCountRule, maxSubmissionCountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorRoundDetailsUpdated)
				if err := _FluxAggregator.contract.UnpackLog(event, "RoundDetailsUpdated", log); err != nil {
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

// ParseRoundDetailsUpdated is a log parse operation binding the contract event 0x56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f.
//
// Solidity: event RoundDetailsUpdated(uint128 indexed paymentAmount, uint32 indexed minSubmissionCount, uint32 indexed maxSubmissionCount, uint32 restartDelay, uint32 timeout)
func (_FluxAggregator *FluxAggregatorFilterer) ParseRoundDetailsUpdated(log types.Log) (*FluxAggregatorRoundDetailsUpdated, error) {
	event := new(FluxAggregatorRoundDetailsUpdated)
	if err := _FluxAggregator.contract.UnpackLog(event, "RoundDetailsUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FluxAggregatorSubmissionReceivedIterator is returned from FilterSubmissionReceived and is used to iterate over the raw logs and unpacked data for SubmissionReceived events raised by the FluxAggregator contract.
type FluxAggregatorSubmissionReceivedIterator struct {
	Event *FluxAggregatorSubmissionReceived // Event containing the contract specifics and raw log

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
func (it *FluxAggregatorSubmissionReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FluxAggregatorSubmissionReceived)
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
		it.Event = new(FluxAggregatorSubmissionReceived)
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
func (it *FluxAggregatorSubmissionReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FluxAggregatorSubmissionReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FluxAggregatorSubmissionReceived represents a SubmissionReceived event raised by the FluxAggregator contract.
type FluxAggregatorSubmissionReceived struct {
	Submission *big.Int
	Round      uint32
	Oracle     common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSubmissionReceived is a free log retrieval operation binding the contract event 0x92e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68c.
//
// Solidity: event SubmissionReceived(int256 indexed submission, uint32 indexed round, address indexed oracle)
func (_FluxAggregator *FluxAggregatorFilterer) FilterSubmissionReceived(opts *bind.FilterOpts, submission []*big.Int, round []uint32, oracle []common.Address) (*FluxAggregatorSubmissionReceivedIterator, error) {

	var submissionRule []interface{}
	for _, submissionItem := range submission {
		submissionRule = append(submissionRule, submissionItem)
	}
	var roundRule []interface{}
	for _, roundItem := range round {
		roundRule = append(roundRule, roundItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _FluxAggregator.contract.FilterLogs(opts, "SubmissionReceived", submissionRule, roundRule, oracleRule)
	if err != nil {
		return nil, err
	}
	return &FluxAggregatorSubmissionReceivedIterator{contract: _FluxAggregator.contract, event: "SubmissionReceived", logs: logs, sub: sub}, nil
}

// WatchSubmissionReceived is a free log subscription operation binding the contract event 0x92e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68c.
//
// Solidity: event SubmissionReceived(int256 indexed submission, uint32 indexed round, address indexed oracle)
func (_FluxAggregator *FluxAggregatorFilterer) WatchSubmissionReceived(opts *bind.WatchOpts, sink chan<- *FluxAggregatorSubmissionReceived, submission []*big.Int, round []uint32, oracle []common.Address) (event.Subscription, error) {

	var submissionRule []interface{}
	for _, submissionItem := range submission {
		submissionRule = append(submissionRule, submissionItem)
	}
	var roundRule []interface{}
	for _, roundItem := range round {
		roundRule = append(roundRule, roundItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _FluxAggregator.contract.WatchLogs(opts, "SubmissionReceived", submissionRule, roundRule, oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FluxAggregatorSubmissionReceived)
				if err := _FluxAggregator.contract.UnpackLog(event, "SubmissionReceived", log); err != nil {
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

// ParseSubmissionReceived is a log parse operation binding the contract event 0x92e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68c.
//
// Solidity: event SubmissionReceived(int256 indexed submission, uint32 indexed round, address indexed oracle)
func (_FluxAggregator *FluxAggregatorFilterer) ParseSubmissionReceived(log types.Log) (*FluxAggregatorSubmissionReceived, error) {
	event := new(FluxAggregatorSubmissionReceived)
	if err := _FluxAggregator.contract.UnpackLog(event, "SubmissionReceived", log); err != nil {
		return nil, err
	}
	return event, nil
}

// HistoricAggregatorInterfaceABI is the input ABI used to generate the binding from.
const HistoricAggregatorInterfaceABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// HistoricAggregatorInterfaceFuncSigs maps the 4-byte function signature to its string representation.
var HistoricAggregatorInterfaceFuncSigs = map[string]string{
	"b5ab58dc": "getAnswer(uint256)",
	"b633620c": "getTimestamp(uint256)",
	"50d25bcd": "latestAnswer()",
	"668a0f02": "latestRound()",
	"8205bf6a": "latestTimestamp()",
}

// HistoricAggregatorInterface is an auto generated Go binding around an Ethereum contract.
type HistoricAggregatorInterface struct {
	HistoricAggregatorInterfaceCaller     // Read-only binding to the contract
	HistoricAggregatorInterfaceTransactor // Write-only binding to the contract
	HistoricAggregatorInterfaceFilterer   // Log filterer for contract events
}

// HistoricAggregatorInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type HistoricAggregatorInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HistoricAggregatorInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HistoricAggregatorInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HistoricAggregatorInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HistoricAggregatorInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HistoricAggregatorInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HistoricAggregatorInterfaceSession struct {
	Contract     *HistoricAggregatorInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// HistoricAggregatorInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HistoricAggregatorInterfaceCallerSession struct {
	Contract *HistoricAggregatorInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// HistoricAggregatorInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HistoricAggregatorInterfaceTransactorSession struct {
	Contract     *HistoricAggregatorInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// HistoricAggregatorInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type HistoricAggregatorInterfaceRaw struct {
	Contract *HistoricAggregatorInterface // Generic contract binding to access the raw methods on
}

// HistoricAggregatorInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HistoricAggregatorInterfaceCallerRaw struct {
	Contract *HistoricAggregatorInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// HistoricAggregatorInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HistoricAggregatorInterfaceTransactorRaw struct {
	Contract *HistoricAggregatorInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHistoricAggregatorInterface creates a new instance of HistoricAggregatorInterface, bound to a specific deployed contract.
func NewHistoricAggregatorInterface(address common.Address, backend bind.ContractBackend) (*HistoricAggregatorInterface, error) {
	contract, err := bindHistoricAggregatorInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HistoricAggregatorInterface{HistoricAggregatorInterfaceCaller: HistoricAggregatorInterfaceCaller{contract: contract}, HistoricAggregatorInterfaceTransactor: HistoricAggregatorInterfaceTransactor{contract: contract}, HistoricAggregatorInterfaceFilterer: HistoricAggregatorInterfaceFilterer{contract: contract}}, nil
}

// NewHistoricAggregatorInterfaceCaller creates a new read-only instance of HistoricAggregatorInterface, bound to a specific deployed contract.
func NewHistoricAggregatorInterfaceCaller(address common.Address, caller bind.ContractCaller) (*HistoricAggregatorInterfaceCaller, error) {
	contract, err := bindHistoricAggregatorInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HistoricAggregatorInterfaceCaller{contract: contract}, nil
}

// NewHistoricAggregatorInterfaceTransactor creates a new write-only instance of HistoricAggregatorInterface, bound to a specific deployed contract.
func NewHistoricAggregatorInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*HistoricAggregatorInterfaceTransactor, error) {
	contract, err := bindHistoricAggregatorInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HistoricAggregatorInterfaceTransactor{contract: contract}, nil
}

// NewHistoricAggregatorInterfaceFilterer creates a new log filterer instance of HistoricAggregatorInterface, bound to a specific deployed contract.
func NewHistoricAggregatorInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*HistoricAggregatorInterfaceFilterer, error) {
	contract, err := bindHistoricAggregatorInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HistoricAggregatorInterfaceFilterer{contract: contract}, nil
}

// bindHistoricAggregatorInterface binds a generic wrapper to an already deployed contract.
func bindHistoricAggregatorInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HistoricAggregatorInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _HistoricAggregatorInterface.Contract.HistoricAggregatorInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HistoricAggregatorInterface.Contract.HistoricAggregatorInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HistoricAggregatorInterface.Contract.HistoricAggregatorInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _HistoricAggregatorInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HistoricAggregatorInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HistoricAggregatorInterface.Contract.contract.Transact(opts, method, params...)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) constant returns(int256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _HistoricAggregatorInterface.contract.Call(opts, out, "getAnswer", roundId)
	return *ret0, err
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) constant returns(int256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _HistoricAggregatorInterface.Contract.GetAnswer(&_HistoricAggregatorInterface.CallOpts, roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) constant returns(int256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _HistoricAggregatorInterface.Contract.GetAnswer(&_HistoricAggregatorInterface.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) constant returns(uint256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _HistoricAggregatorInterface.contract.Call(opts, out, "getTimestamp", roundId)
	return *ret0, err
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) constant returns(uint256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _HistoricAggregatorInterface.Contract.GetTimestamp(&_HistoricAggregatorInterface.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) constant returns(uint256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _HistoricAggregatorInterface.Contract.GetTimestamp(&_HistoricAggregatorInterface.CallOpts, roundId)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() constant returns(int256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _HistoricAggregatorInterface.contract.Call(opts, out, "latestAnswer")
	return *ret0, err
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() constant returns(int256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceSession) LatestAnswer() (*big.Int, error) {
	return _HistoricAggregatorInterface.Contract.LatestAnswer(&_HistoricAggregatorInterface.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() constant returns(int256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceCallerSession) LatestAnswer() (*big.Int, error) {
	return _HistoricAggregatorInterface.Contract.LatestAnswer(&_HistoricAggregatorInterface.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() constant returns(uint256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _HistoricAggregatorInterface.contract.Call(opts, out, "latestRound")
	return *ret0, err
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() constant returns(uint256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceSession) LatestRound() (*big.Int, error) {
	return _HistoricAggregatorInterface.Contract.LatestRound(&_HistoricAggregatorInterface.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() constant returns(uint256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceCallerSession) LatestRound() (*big.Int, error) {
	return _HistoricAggregatorInterface.Contract.LatestRound(&_HistoricAggregatorInterface.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() constant returns(uint256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _HistoricAggregatorInterface.contract.Call(opts, out, "latestTimestamp")
	return *ret0, err
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() constant returns(uint256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceSession) LatestTimestamp() (*big.Int, error) {
	return _HistoricAggregatorInterface.Contract.LatestTimestamp(&_HistoricAggregatorInterface.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() constant returns(uint256)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceCallerSession) LatestTimestamp() (*big.Int, error) {
	return _HistoricAggregatorInterface.Contract.LatestTimestamp(&_HistoricAggregatorInterface.CallOpts)
}

// HistoricAggregatorInterfaceAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the HistoricAggregatorInterface contract.
type HistoricAggregatorInterfaceAnswerUpdatedIterator struct {
	Event *HistoricAggregatorInterfaceAnswerUpdated // Event containing the contract specifics and raw log

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
func (it *HistoricAggregatorInterfaceAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HistoricAggregatorInterfaceAnswerUpdated)
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
		it.Event = new(HistoricAggregatorInterfaceAnswerUpdated)
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
func (it *HistoricAggregatorInterfaceAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HistoricAggregatorInterfaceAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HistoricAggregatorInterfaceAnswerUpdated represents a AnswerUpdated event raised by the HistoricAggregatorInterface contract.
type HistoricAggregatorInterfaceAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*HistoricAggregatorInterfaceAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _HistoricAggregatorInterface.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &HistoricAggregatorInterfaceAnswerUpdatedIterator{contract: _HistoricAggregatorInterface.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *HistoricAggregatorInterfaceAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _HistoricAggregatorInterface.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HistoricAggregatorInterfaceAnswerUpdated)
				if err := _HistoricAggregatorInterface.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceFilterer) ParseAnswerUpdated(log types.Log) (*HistoricAggregatorInterfaceAnswerUpdated, error) {
	event := new(HistoricAggregatorInterfaceAnswerUpdated)
	if err := _HistoricAggregatorInterface.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// HistoricAggregatorInterfaceNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the HistoricAggregatorInterface contract.
type HistoricAggregatorInterfaceNewRoundIterator struct {
	Event *HistoricAggregatorInterfaceNewRound // Event containing the contract specifics and raw log

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
func (it *HistoricAggregatorInterfaceNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HistoricAggregatorInterfaceNewRound)
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
		it.Event = new(HistoricAggregatorInterfaceNewRound)
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
func (it *HistoricAggregatorInterfaceNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HistoricAggregatorInterfaceNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HistoricAggregatorInterfaceNewRound represents a NewRound event raised by the HistoricAggregatorInterface contract.
type HistoricAggregatorInterfaceNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*HistoricAggregatorInterfaceNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _HistoricAggregatorInterface.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &HistoricAggregatorInterfaceNewRoundIterator{contract: _HistoricAggregatorInterface.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *HistoricAggregatorInterfaceNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _HistoricAggregatorInterface.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HistoricAggregatorInterfaceNewRound)
				if err := _HistoricAggregatorInterface.contract.UnpackLog(event, "NewRound", log); err != nil {
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
func (_HistoricAggregatorInterface *HistoricAggregatorInterfaceFilterer) ParseNewRound(log types.Log) (*HistoricAggregatorInterfaceNewRound, error) {
	event := new(HistoricAggregatorInterfaceNewRound)
	if err := _HistoricAggregatorInterface.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LinkTokenInterfaceABI is the input ABI used to generate the binding from.
const LinkTokenInterfaceABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"remaining\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"decimalPlaces\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"tokenName\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"tokenSymbol\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalTokensIssued\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"transferAndCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// LinkTokenInterfaceFuncSigs maps the 4-byte function signature to its string representation.
var LinkTokenInterfaceFuncSigs = map[string]string{
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"313ce567": "decimals()",
	"66188463": "decreaseApproval(address,uint256)",
	"d73dd623": "increaseApproval(address,uint256)",
	"06fdde03": "name()",
	"95d89b41": "symbol()",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"4000aea0": "transferAndCall(address,uint256,bytes)",
	"23b872dd": "transferFrom(address,address,uint256)",
}

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
	parsed, err := abi.JSON(strings.NewReader(LinkTokenInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LinkTokenInterface *LinkTokenInterfaceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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
func (_LinkTokenInterface *LinkTokenInterfaceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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
// Solidity: function allowance(address owner, address spender) constant returns(uint256 remaining)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LinkTokenInterface.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) constant returns(uint256 remaining)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.Allowance(&_LinkTokenInterface.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) constant returns(uint256 remaining)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.Allowance(&_LinkTokenInterface.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) constant returns(uint256 balance)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LinkTokenInterface.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) constant returns(uint256 balance)
func (_LinkTokenInterface *LinkTokenInterfaceSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.BalanceOf(&_LinkTokenInterface.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) constant returns(uint256 balance)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.BalanceOf(&_LinkTokenInterface.CallOpts, owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8 decimalPlaces)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _LinkTokenInterface.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8 decimalPlaces)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Decimals() (uint8, error) {
	return _LinkTokenInterface.Contract.Decimals(&_LinkTokenInterface.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8 decimalPlaces)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Decimals() (uint8, error) {
	return _LinkTokenInterface.Contract.Decimals(&_LinkTokenInterface.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string tokenName)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LinkTokenInterface.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string tokenName)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Name() (string, error) {
	return _LinkTokenInterface.Contract.Name(&_LinkTokenInterface.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string tokenName)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Name() (string, error) {
	return _LinkTokenInterface.Contract.Name(&_LinkTokenInterface.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string tokenSymbol)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LinkTokenInterface.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string tokenSymbol)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Symbol() (string, error) {
	return _LinkTokenInterface.Contract.Symbol(&_LinkTokenInterface.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string tokenSymbol)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Symbol() (string, error) {
	return _LinkTokenInterface.Contract.Symbol(&_LinkTokenInterface.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256 totalTokensIssued)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LinkTokenInterface.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256 totalTokensIssued)
func (_LinkTokenInterface *LinkTokenInterfaceSession) TotalSupply() (*big.Int, error) {
	return _LinkTokenInterface.Contract.TotalSupply(&_LinkTokenInterface.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256 totalTokensIssued)
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

// MedianABI is the input ABI used to generate the binding from.
const MedianABI = "[]"

// MedianBin is the compiled bytecode used for deploying new contracts.
var MedianBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212207f9f65a3d153989a1e77b8faecaea976b0e0c8fc6977731ab86ab3f63a1288da64736f6c63430006060033"

// DeployMedian deploys a new Ethereum contract, binding an instance of Median to it.
func DeployMedian(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Median, error) {
	parsed, err := abi.JSON(strings.NewReader(MedianABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MedianBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Median{MedianCaller: MedianCaller{contract: contract}, MedianTransactor: MedianTransactor{contract: contract}, MedianFilterer: MedianFilterer{contract: contract}}, nil
}

// Median is an auto generated Go binding around an Ethereum contract.
type Median struct {
	MedianCaller     // Read-only binding to the contract
	MedianTransactor // Write-only binding to the contract
	MedianFilterer   // Log filterer for contract events
}

// MedianCaller is an auto generated read-only Go binding around an Ethereum contract.
type MedianCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MedianTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MedianTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MedianFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MedianFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MedianSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MedianSession struct {
	Contract     *Median           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MedianCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MedianCallerSession struct {
	Contract *MedianCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// MedianTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MedianTransactorSession struct {
	Contract     *MedianTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MedianRaw is an auto generated low-level Go binding around an Ethereum contract.
type MedianRaw struct {
	Contract *Median // Generic contract binding to access the raw methods on
}

// MedianCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MedianCallerRaw struct {
	Contract *MedianCaller // Generic read-only contract binding to access the raw methods on
}

// MedianTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MedianTransactorRaw struct {
	Contract *MedianTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMedian creates a new instance of Median, bound to a specific deployed contract.
func NewMedian(address common.Address, backend bind.ContractBackend) (*Median, error) {
	contract, err := bindMedian(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Median{MedianCaller: MedianCaller{contract: contract}, MedianTransactor: MedianTransactor{contract: contract}, MedianFilterer: MedianFilterer{contract: contract}}, nil
}

// NewMedianCaller creates a new read-only instance of Median, bound to a specific deployed contract.
func NewMedianCaller(address common.Address, caller bind.ContractCaller) (*MedianCaller, error) {
	contract, err := bindMedian(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MedianCaller{contract: contract}, nil
}

// NewMedianTransactor creates a new write-only instance of Median, bound to a specific deployed contract.
func NewMedianTransactor(address common.Address, transactor bind.ContractTransactor) (*MedianTransactor, error) {
	contract, err := bindMedian(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MedianTransactor{contract: contract}, nil
}

// NewMedianFilterer creates a new log filterer instance of Median, bound to a specific deployed contract.
func NewMedianFilterer(address common.Address, filterer bind.ContractFilterer) (*MedianFilterer, error) {
	contract, err := bindMedian(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MedianFilterer{contract: contract}, nil
}

// bindMedian binds a generic wrapper to an already deployed contract.
func bindMedian(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MedianABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Median *MedianRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Median.Contract.MedianCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Median *MedianRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Median.Contract.MedianTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Median *MedianRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Median.Contract.MedianTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Median *MedianCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Median.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Median *MedianTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Median.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Median *MedianTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Median.Contract.contract.Transact(opts, method, params...)
}

// OwnedABI is the input ABI used to generate the binding from.
const OwnedABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// OwnedFuncSigs maps the 4-byte function signature to its string representation.
var OwnedFuncSigs = map[string]string{
	"79ba5097": "acceptOwnership()",
	"8da5cb5b": "owner()",
	"f2fde38b": "transferOwnership(address)",
}

// OwnedBin is the compiled bytecode used for deploying new contracts.
var OwnedBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b03191633179055610237806100326000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806379ba5097146100465780638da5cb5b14610050578063f2fde38b14610074575b600080fd5b61004e61009a565b005b610058610149565b604080516001600160a01b039092168252519081900360200190f35b61004e6004803603602081101561008a57600080fd5b50356001600160a01b0316610158565b6001546001600160a01b031633146100f2576040805162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b604482015290519081900360640190fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000546001600160a01b031681565b6000546001600160a01b031633146101b0576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a35056fea264697066735822122036360901665034eafe982da55068fce2db52b2d93a84e70804cc8710dbf0cd9a64736f6c63430006060033"

// DeployOwned deploys a new Ethereum contract, binding an instance of Owned to it.
func DeployOwned(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Owned, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnedBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}

// Owned is an auto generated Go binding around an Ethereum contract.
type Owned struct {
	OwnedCaller     // Read-only binding to the contract
	OwnedTransactor // Write-only binding to the contract
	OwnedFilterer   // Log filterer for contract events
}

// OwnedCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnedCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnedTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnedFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnedSession struct {
	Contract     *Owned            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnedCallerSession struct {
	Contract *OwnedCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OwnedTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnedTransactorSession struct {
	Contract     *OwnedTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnedRaw struct {
	Contract *Owned // Generic contract binding to access the raw methods on
}

// OwnedCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnedCallerRaw struct {
	Contract *OwnedCaller // Generic read-only contract binding to access the raw methods on
}

// OwnedTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnedTransactorRaw struct {
	Contract *OwnedTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwned creates a new instance of Owned, bound to a specific deployed contract.
func NewOwned(address common.Address, backend bind.ContractBackend) (*Owned, error) {
	contract, err := bindOwned(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}

// NewOwnedCaller creates a new read-only instance of Owned, bound to a specific deployed contract.
func NewOwnedCaller(address common.Address, caller bind.ContractCaller) (*OwnedCaller, error) {
	contract, err := bindOwned(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedCaller{contract: contract}, nil
}

// NewOwnedTransactor creates a new write-only instance of Owned, bound to a specific deployed contract.
func NewOwnedTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnedTransactor, error) {
	contract, err := bindOwned(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedTransactor{contract: contract}, nil
}

// NewOwnedFilterer creates a new log filterer instance of Owned, bound to a specific deployed contract.
func NewOwnedFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnedFilterer, error) {
	contract, err := bindOwned(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnedFilterer{contract: contract}, nil
}

// bindOwned binds a generic wrapper to an already deployed contract.
func bindOwned(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.OwnedCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Owned.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedCallerSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedSession) AcceptOwnership() (*types.Transaction, error) {
	return _Owned.Contract.AcceptOwnership(&_Owned.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Owned.Contract.AcceptOwnership(&_Owned.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Owned *OwnedTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Owned *OwnedSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _Owned.Contract.TransferOwnership(&_Owned.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Owned *OwnedTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _Owned.Contract.TransferOwnership(&_Owned.TransactOpts, _to)
}

// OwnedOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the Owned contract.
type OwnedOwnershipTransferRequestedIterator struct {
	Event *OwnedOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *OwnedOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedOwnershipTransferRequested)
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
		it.Event = new(OwnedOwnershipTransferRequested)
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
func (it *OwnedOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the Owned contract.
type OwnedOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Owned *OwnedFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OwnedOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Owned.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OwnedOwnershipTransferRequestedIterator{contract: _Owned.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Owned *OwnedFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OwnedOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Owned.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedOwnershipTransferRequested)
				if err := _Owned.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_Owned *OwnedFilterer) ParseOwnershipTransferRequested(log types.Log) (*OwnedOwnershipTransferRequested, error) {
	event := new(OwnedOwnershipTransferRequested)
	if err := _Owned.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OwnedOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Owned contract.
type OwnedOwnershipTransferredIterator struct {
	Event *OwnedOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OwnedOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedOwnershipTransferred)
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
		it.Event = new(OwnedOwnershipTransferred)
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
func (it *OwnedOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedOwnershipTransferred represents a OwnershipTransferred event raised by the Owned contract.
type OwnedOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Owned *OwnedFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OwnedOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Owned.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OwnedOwnershipTransferredIterator{contract: _Owned.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Owned *OwnedFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnedOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Owned.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedOwnershipTransferred)
				if err := _Owned.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Owned *OwnedFilterer) ParseOwnershipTransferred(log types.Log) (*OwnedOwnershipTransferred, error) {
	event := new(OwnedOwnershipTransferred)
	if err := _Owned.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
var SafeMathBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea264697066735822122087b2ebbd1816328ade45ce97dc6bd31ab556df02bdb32494c710d1f61c9edc7364736f6c63430006060033"

// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Ethereum contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}

// SafeMath128ABI is the input ABI used to generate the binding from.
const SafeMath128ABI = "[]"

// SafeMath128Bin is the compiled bytecode used for deploying new contracts.
var SafeMath128Bin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220a2751c85745493a37cf742acd7dc008948bfb732a758b3253e0fe5944210eda564736f6c63430006060033"

// DeploySafeMath128 deploys a new Ethereum contract, binding an instance of SafeMath128 to it.
func DeploySafeMath128(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath128, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMath128ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMath128Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath128{SafeMath128Caller: SafeMath128Caller{contract: contract}, SafeMath128Transactor: SafeMath128Transactor{contract: contract}, SafeMath128Filterer: SafeMath128Filterer{contract: contract}}, nil
}

// SafeMath128 is an auto generated Go binding around an Ethereum contract.
type SafeMath128 struct {
	SafeMath128Caller     // Read-only binding to the contract
	SafeMath128Transactor // Write-only binding to the contract
	SafeMath128Filterer   // Log filterer for contract events
}

// SafeMath128Caller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMath128Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath128Transactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMath128Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath128Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMath128Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath128Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMath128Session struct {
	Contract     *SafeMath128      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMath128CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMath128CallerSession struct {
	Contract *SafeMath128Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// SafeMath128TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMath128TransactorSession struct {
	Contract     *SafeMath128Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// SafeMath128Raw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMath128Raw struct {
	Contract *SafeMath128 // Generic contract binding to access the raw methods on
}

// SafeMath128CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMath128CallerRaw struct {
	Contract *SafeMath128Caller // Generic read-only contract binding to access the raw methods on
}

// SafeMath128TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMath128TransactorRaw struct {
	Contract *SafeMath128Transactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath128 creates a new instance of SafeMath128, bound to a specific deployed contract.
func NewSafeMath128(address common.Address, backend bind.ContractBackend) (*SafeMath128, error) {
	contract, err := bindSafeMath128(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath128{SafeMath128Caller: SafeMath128Caller{contract: contract}, SafeMath128Transactor: SafeMath128Transactor{contract: contract}, SafeMath128Filterer: SafeMath128Filterer{contract: contract}}, nil
}

// NewSafeMath128Caller creates a new read-only instance of SafeMath128, bound to a specific deployed contract.
func NewSafeMath128Caller(address common.Address, caller bind.ContractCaller) (*SafeMath128Caller, error) {
	contract, err := bindSafeMath128(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMath128Caller{contract: contract}, nil
}

// NewSafeMath128Transactor creates a new write-only instance of SafeMath128, bound to a specific deployed contract.
func NewSafeMath128Transactor(address common.Address, transactor bind.ContractTransactor) (*SafeMath128Transactor, error) {
	contract, err := bindSafeMath128(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMath128Transactor{contract: contract}, nil
}

// NewSafeMath128Filterer creates a new log filterer instance of SafeMath128, bound to a specific deployed contract.
func NewSafeMath128Filterer(address common.Address, filterer bind.ContractFilterer) (*SafeMath128Filterer, error) {
	contract, err := bindSafeMath128(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMath128Filterer{contract: contract}, nil
}

// bindSafeMath128 binds a generic wrapper to an already deployed contract.
func bindSafeMath128(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMath128ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath128 *SafeMath128Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath128.Contract.SafeMath128Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath128 *SafeMath128Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath128.Contract.SafeMath128Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath128 *SafeMath128Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath128.Contract.SafeMath128Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath128 *SafeMath128CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath128.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath128 *SafeMath128TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath128.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath128 *SafeMath128TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath128.Contract.contract.Transact(opts, method, params...)
}

// SafeMath32ABI is the input ABI used to generate the binding from.
const SafeMath32ABI = "[]"

// SafeMath32Bin is the compiled bytecode used for deploying new contracts.
var SafeMath32Bin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220290efef6f374dc5102574204a2cbc22f9ef159f46540a87a960a22fcfb61c9fd64736f6c63430006060033"

// DeploySafeMath32 deploys a new Ethereum contract, binding an instance of SafeMath32 to it.
func DeploySafeMath32(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath32, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMath32ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMath32Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath32{SafeMath32Caller: SafeMath32Caller{contract: contract}, SafeMath32Transactor: SafeMath32Transactor{contract: contract}, SafeMath32Filterer: SafeMath32Filterer{contract: contract}}, nil
}

// SafeMath32 is an auto generated Go binding around an Ethereum contract.
type SafeMath32 struct {
	SafeMath32Caller     // Read-only binding to the contract
	SafeMath32Transactor // Write-only binding to the contract
	SafeMath32Filterer   // Log filterer for contract events
}

// SafeMath32Caller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMath32Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath32Transactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMath32Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath32Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMath32Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath32Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMath32Session struct {
	Contract     *SafeMath32       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMath32CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMath32CallerSession struct {
	Contract *SafeMath32Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// SafeMath32TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMath32TransactorSession struct {
	Contract     *SafeMath32Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// SafeMath32Raw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMath32Raw struct {
	Contract *SafeMath32 // Generic contract binding to access the raw methods on
}

// SafeMath32CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMath32CallerRaw struct {
	Contract *SafeMath32Caller // Generic read-only contract binding to access the raw methods on
}

// SafeMath32TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMath32TransactorRaw struct {
	Contract *SafeMath32Transactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath32 creates a new instance of SafeMath32, bound to a specific deployed contract.
func NewSafeMath32(address common.Address, backend bind.ContractBackend) (*SafeMath32, error) {
	contract, err := bindSafeMath32(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath32{SafeMath32Caller: SafeMath32Caller{contract: contract}, SafeMath32Transactor: SafeMath32Transactor{contract: contract}, SafeMath32Filterer: SafeMath32Filterer{contract: contract}}, nil
}

// NewSafeMath32Caller creates a new read-only instance of SafeMath32, bound to a specific deployed contract.
func NewSafeMath32Caller(address common.Address, caller bind.ContractCaller) (*SafeMath32Caller, error) {
	contract, err := bindSafeMath32(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMath32Caller{contract: contract}, nil
}

// NewSafeMath32Transactor creates a new write-only instance of SafeMath32, bound to a specific deployed contract.
func NewSafeMath32Transactor(address common.Address, transactor bind.ContractTransactor) (*SafeMath32Transactor, error) {
	contract, err := bindSafeMath32(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMath32Transactor{contract: contract}, nil
}

// NewSafeMath32Filterer creates a new log filterer instance of SafeMath32, bound to a specific deployed contract.
func NewSafeMath32Filterer(address common.Address, filterer bind.ContractFilterer) (*SafeMath32Filterer, error) {
	contract, err := bindSafeMath32(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMath32Filterer{contract: contract}, nil
}

// bindSafeMath32 binds a generic wrapper to an already deployed contract.
func bindSafeMath32(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMath32ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath32 *SafeMath32Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath32.Contract.SafeMath32Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath32 *SafeMath32Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath32.Contract.SafeMath32Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath32 *SafeMath32Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath32.Contract.SafeMath32Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath32 *SafeMath32CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath32.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath32 *SafeMath32TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath32.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath32 *SafeMath32TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath32.Contract.contract.Transact(opts, method, params...)
}

// SafeMath64ABI is the input ABI used to generate the binding from.
const SafeMath64ABI = "[]"

// SafeMath64Bin is the compiled bytecode used for deploying new contracts.
var SafeMath64Bin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212202b800250c5dfc75c4543207fce95860d755520f9b62468b433e1382786dbf75164736f6c63430006060033"

// DeploySafeMath64 deploys a new Ethereum contract, binding an instance of SafeMath64 to it.
func DeploySafeMath64(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath64, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMath64ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMath64Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath64{SafeMath64Caller: SafeMath64Caller{contract: contract}, SafeMath64Transactor: SafeMath64Transactor{contract: contract}, SafeMath64Filterer: SafeMath64Filterer{contract: contract}}, nil
}

// SafeMath64 is an auto generated Go binding around an Ethereum contract.
type SafeMath64 struct {
	SafeMath64Caller     // Read-only binding to the contract
	SafeMath64Transactor // Write-only binding to the contract
	SafeMath64Filterer   // Log filterer for contract events
}

// SafeMath64Caller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMath64Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath64Transactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMath64Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath64Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMath64Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath64Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMath64Session struct {
	Contract     *SafeMath64       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMath64CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMath64CallerSession struct {
	Contract *SafeMath64Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// SafeMath64TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMath64TransactorSession struct {
	Contract     *SafeMath64Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// SafeMath64Raw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMath64Raw struct {
	Contract *SafeMath64 // Generic contract binding to access the raw methods on
}

// SafeMath64CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMath64CallerRaw struct {
	Contract *SafeMath64Caller // Generic read-only contract binding to access the raw methods on
}

// SafeMath64TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMath64TransactorRaw struct {
	Contract *SafeMath64Transactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath64 creates a new instance of SafeMath64, bound to a specific deployed contract.
func NewSafeMath64(address common.Address, backend bind.ContractBackend) (*SafeMath64, error) {
	contract, err := bindSafeMath64(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath64{SafeMath64Caller: SafeMath64Caller{contract: contract}, SafeMath64Transactor: SafeMath64Transactor{contract: contract}, SafeMath64Filterer: SafeMath64Filterer{contract: contract}}, nil
}

// NewSafeMath64Caller creates a new read-only instance of SafeMath64, bound to a specific deployed contract.
func NewSafeMath64Caller(address common.Address, caller bind.ContractCaller) (*SafeMath64Caller, error) {
	contract, err := bindSafeMath64(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMath64Caller{contract: contract}, nil
}

// NewSafeMath64Transactor creates a new write-only instance of SafeMath64, bound to a specific deployed contract.
func NewSafeMath64Transactor(address common.Address, transactor bind.ContractTransactor) (*SafeMath64Transactor, error) {
	contract, err := bindSafeMath64(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMath64Transactor{contract: contract}, nil
}

// NewSafeMath64Filterer creates a new log filterer instance of SafeMath64, bound to a specific deployed contract.
func NewSafeMath64Filterer(address common.Address, filterer bind.ContractFilterer) (*SafeMath64Filterer, error) {
	contract, err := bindSafeMath64(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMath64Filterer{contract: contract}, nil
}

// bindSafeMath64 binds a generic wrapper to an already deployed contract.
func bindSafeMath64(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMath64ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath64 *SafeMath64Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath64.Contract.SafeMath64Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath64 *SafeMath64Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath64.Contract.SafeMath64Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath64 *SafeMath64Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath64.Contract.SafeMath64Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath64 *SafeMath64CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath64.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath64 *SafeMath64TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath64.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath64 *SafeMath64TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath64.Contract.contract.Transact(opts, method, params...)
}

// SignedSafeMathABI is the input ABI used to generate the binding from.
const SignedSafeMathABI = "[]"

// SignedSafeMathBin is the compiled bytecode used for deploying new contracts.
var SignedSafeMathBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220c9b70617b2db308a31ceec052749444516ee3152d2fea05cac45f31a566fa15964736f6c63430006060033"

// DeploySignedSafeMath deploys a new Ethereum contract, binding an instance of SignedSafeMath to it.
func DeploySignedSafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SignedSafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SignedSafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SignedSafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SignedSafeMath{SignedSafeMathCaller: SignedSafeMathCaller{contract: contract}, SignedSafeMathTransactor: SignedSafeMathTransactor{contract: contract}, SignedSafeMathFilterer: SignedSafeMathFilterer{contract: contract}}, nil
}

// SignedSafeMath is an auto generated Go binding around an Ethereum contract.
type SignedSafeMath struct {
	SignedSafeMathCaller     // Read-only binding to the contract
	SignedSafeMathTransactor // Write-only binding to the contract
	SignedSafeMathFilterer   // Log filterer for contract events
}

// SignedSafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SignedSafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignedSafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SignedSafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignedSafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SignedSafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignedSafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SignedSafeMathSession struct {
	Contract     *SignedSafeMath   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SignedSafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SignedSafeMathCallerSession struct {
	Contract *SignedSafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// SignedSafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SignedSafeMathTransactorSession struct {
	Contract     *SignedSafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// SignedSafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SignedSafeMathRaw struct {
	Contract *SignedSafeMath // Generic contract binding to access the raw methods on
}

// SignedSafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SignedSafeMathCallerRaw struct {
	Contract *SignedSafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SignedSafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SignedSafeMathTransactorRaw struct {
	Contract *SignedSafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSignedSafeMath creates a new instance of SignedSafeMath, bound to a specific deployed contract.
func NewSignedSafeMath(address common.Address, backend bind.ContractBackend) (*SignedSafeMath, error) {
	contract, err := bindSignedSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SignedSafeMath{SignedSafeMathCaller: SignedSafeMathCaller{contract: contract}, SignedSafeMathTransactor: SignedSafeMathTransactor{contract: contract}, SignedSafeMathFilterer: SignedSafeMathFilterer{contract: contract}}, nil
}

// NewSignedSafeMathCaller creates a new read-only instance of SignedSafeMath, bound to a specific deployed contract.
func NewSignedSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SignedSafeMathCaller, error) {
	contract, err := bindSignedSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SignedSafeMathCaller{contract: contract}, nil
}

// NewSignedSafeMathTransactor creates a new write-only instance of SignedSafeMath, bound to a specific deployed contract.
func NewSignedSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SignedSafeMathTransactor, error) {
	contract, err := bindSignedSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SignedSafeMathTransactor{contract: contract}, nil
}

// NewSignedSafeMathFilterer creates a new log filterer instance of SignedSafeMath, bound to a specific deployed contract.
func NewSignedSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SignedSafeMathFilterer, error) {
	contract, err := bindSignedSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SignedSafeMathFilterer{contract: contract}, nil
}

// bindSignedSafeMath binds a generic wrapper to an already deployed contract.
func bindSignedSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SignedSafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SignedSafeMath *SignedSafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SignedSafeMath.Contract.SignedSafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SignedSafeMath *SignedSafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignedSafeMath.Contract.SignedSafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SignedSafeMath *SignedSafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignedSafeMath.Contract.SignedSafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SignedSafeMath *SignedSafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SignedSafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SignedSafeMath *SignedSafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignedSafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SignedSafeMath *SignedSafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignedSafeMath.Contract.contract.Transact(opts, method, params...)
}
