// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

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
)

// MockETHLINKAggregatorMetaData contains all meta data concerning the MockETHLINKAggregator contract.
var MockETHLINKAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"answer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"ans\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"ans\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161024538038061024583398101604081905261002f91610037565b600055610050565b60006020828403121561004957600080fd5b5051919050565b6101e68061005f6000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c8063313ce5671461006757806354fd4d501461007b5780637284e4161461008c57806385bb7d69146100c35780639a6fc8f5146100cc578063feaf968c1461011d575b600080fd5b604051601281526020015b60405180910390f35b60015b604051908152602001610072565b604080518082018252601581527426b7b1b5a2aa242624a725a0b3b3b932b3b0ba37b960591b60208201529051610072919061012b565b61007e60005481565b6100e96100da366004610180565b50600054600191429081908490565b604080516001600160501b03968716815260208101959095528401929092526060830152909116608082015260a001610072565b6000546001904280836100e9565b600060208083528351808285015260005b818110156101585785810183015185820160400152820161013c565b8181111561016a576000604083870101525b50601f01601f1916929092016040019392505050565b60006020828403121561019257600080fd5b81356001600160501b03811681146101a957600080fd5b939250505056fea2646970667358221220221f919ef5ae81ca8cddf75f7e7bea18da77a41e7ae7a5415c1dde1cffdc935364736f6c634300080d0033",
}

// MockETHLINKAggregatorABI is the input ABI used to generate the binding from.
// Deprecated: Use MockETHLINKAggregatorMetaData.ABI instead.
var MockETHLINKAggregatorABI = MockETHLINKAggregatorMetaData.ABI

// MockETHLINKAggregatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockETHLINKAggregatorMetaData.Bin instead.
var MockETHLINKAggregatorBin = MockETHLINKAggregatorMetaData.Bin

// DeployMockETHLINKAggregator deploys a new Ethereum contract, binding an instance of MockETHLINKAggregator to it.
func DeployMockETHLINKAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _answer *big.Int) (common.Address, *types.Transaction, *MockETHLINKAggregator, error) {
	parsed, err := MockETHLINKAggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockETHLINKAggregatorBin), backend, _answer)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockETHLINKAggregator{MockETHLINKAggregatorCaller: MockETHLINKAggregatorCaller{contract: contract}, MockETHLINKAggregatorTransactor: MockETHLINKAggregatorTransactor{contract: contract}, MockETHLINKAggregatorFilterer: MockETHLINKAggregatorFilterer{contract: contract}}, nil
}

// MockETHLINKAggregator is an auto generated Go binding around an Ethereum contract.
type MockETHLINKAggregator struct {
	MockETHLINKAggregatorCaller     // Read-only binding to the contract
	MockETHLINKAggregatorTransactor // Write-only binding to the contract
	MockETHLINKAggregatorFilterer   // Log filterer for contract events
}

// MockETHLINKAggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockETHLINKAggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockETHLINKAggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockETHLINKAggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockETHLINKAggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockETHLINKAggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockETHLINKAggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockETHLINKAggregatorSession struct {
	Contract     *MockETHLINKAggregator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// MockETHLINKAggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockETHLINKAggregatorCallerSession struct {
	Contract *MockETHLINKAggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// MockETHLINKAggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockETHLINKAggregatorTransactorSession struct {
	Contract     *MockETHLINKAggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// MockETHLINKAggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockETHLINKAggregatorRaw struct {
	Contract *MockETHLINKAggregator // Generic contract binding to access the raw methods on
}

// MockETHLINKAggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockETHLINKAggregatorCallerRaw struct {
	Contract *MockETHLINKAggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// MockETHLINKAggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockETHLINKAggregatorTransactorRaw struct {
	Contract *MockETHLINKAggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockETHLINKAggregator creates a new instance of MockETHLINKAggregator, bound to a specific deployed contract.
func NewMockETHLINKAggregator(address common.Address, backend bind.ContractBackend) (*MockETHLINKAggregator, error) {
	contract, err := bindMockETHLINKAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockETHLINKAggregator{MockETHLINKAggregatorCaller: MockETHLINKAggregatorCaller{contract: contract}, MockETHLINKAggregatorTransactor: MockETHLINKAggregatorTransactor{contract: contract}, MockETHLINKAggregatorFilterer: MockETHLINKAggregatorFilterer{contract: contract}}, nil
}

// NewMockETHLINKAggregatorCaller creates a new read-only instance of MockETHLINKAggregator, bound to a specific deployed contract.
func NewMockETHLINKAggregatorCaller(address common.Address, caller bind.ContractCaller) (*MockETHLINKAggregatorCaller, error) {
	contract, err := bindMockETHLINKAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockETHLINKAggregatorCaller{contract: contract}, nil
}

// NewMockETHLINKAggregatorTransactor creates a new write-only instance of MockETHLINKAggregator, bound to a specific deployed contract.
func NewMockETHLINKAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*MockETHLINKAggregatorTransactor, error) {
	contract, err := bindMockETHLINKAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockETHLINKAggregatorTransactor{contract: contract}, nil
}

// NewMockETHLINKAggregatorFilterer creates a new log filterer instance of MockETHLINKAggregator, bound to a specific deployed contract.
func NewMockETHLINKAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*MockETHLINKAggregatorFilterer, error) {
	contract, err := bindMockETHLINKAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockETHLINKAggregatorFilterer{contract: contract}, nil
}

// bindMockETHLINKAggregator binds a generic wrapper to an already deployed contract.
func bindMockETHLINKAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MockETHLINKAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockETHLINKAggregator *MockETHLINKAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockETHLINKAggregator.Contract.MockETHLINKAggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockETHLINKAggregator *MockETHLINKAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockETHLINKAggregator.Contract.MockETHLINKAggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockETHLINKAggregator *MockETHLINKAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockETHLINKAggregator.Contract.MockETHLINKAggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockETHLINKAggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockETHLINKAggregator *MockETHLINKAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockETHLINKAggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockETHLINKAggregator *MockETHLINKAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockETHLINKAggregator.Contract.contract.Transact(opts, method, params...)
}

// Answer is a free data retrieval call binding the contract method 0x85bb7d69.
//
// Solidity: function answer() view returns(int256)
func (_MockETHLINKAggregator *MockETHLINKAggregatorCaller) Answer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockETHLINKAggregator.contract.Call(opts, &out, "answer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Answer is a free data retrieval call binding the contract method 0x85bb7d69.
//
// Solidity: function answer() view returns(int256)
func (_MockETHLINKAggregator *MockETHLINKAggregatorSession) Answer() (*big.Int, error) {
	return _MockETHLINKAggregator.Contract.Answer(&_MockETHLINKAggregator.CallOpts)
}

// Answer is a free data retrieval call binding the contract method 0x85bb7d69.
//
// Solidity: function answer() view returns(int256)
func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerSession) Answer() (*big.Int, error) {
	return _MockETHLINKAggregator.Contract.Answer(&_MockETHLINKAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_MockETHLINKAggregator *MockETHLINKAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockETHLINKAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_MockETHLINKAggregator *MockETHLINKAggregatorSession) Decimals() (uint8, error) {
	return _MockETHLINKAggregator.Contract.Decimals(&_MockETHLINKAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerSession) Decimals() (uint8, error) {
	return _MockETHLINKAggregator.Contract.Decimals(&_MockETHLINKAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() pure returns(string)
func (_MockETHLINKAggregator *MockETHLINKAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockETHLINKAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() pure returns(string)
func (_MockETHLINKAggregator *MockETHLINKAggregatorSession) Description() (string, error) {
	return _MockETHLINKAggregator.Contract.Description(&_MockETHLINKAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() pure returns(string)
func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerSession) Description() (string, error) {
	return _MockETHLINKAggregator.Contract.Description(&_MockETHLINKAggregator.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 ans, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockETHLINKAggregator *MockETHLINKAggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Ans             *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _MockETHLINKAggregator.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Ans             *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Ans = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 ans, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockETHLINKAggregator *MockETHLINKAggregatorSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Ans             *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockETHLINKAggregator.Contract.GetRoundData(&_MockETHLINKAggregator.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 ans, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Ans             *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockETHLINKAggregator.Contract.GetRoundData(&_MockETHLINKAggregator.CallOpts, _roundId)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 ans, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockETHLINKAggregator *MockETHLINKAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Ans             *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _MockETHLINKAggregator.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Ans             *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Ans = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 ans, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockETHLINKAggregator *MockETHLINKAggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Ans             *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockETHLINKAggregator.Contract.LatestRoundData(&_MockETHLINKAggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 ans, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Ans             *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockETHLINKAggregator.Contract.LatestRoundData(&_MockETHLINKAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(uint256)
func (_MockETHLINKAggregator *MockETHLINKAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockETHLINKAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(uint256)
func (_MockETHLINKAggregator *MockETHLINKAggregatorSession) Version() (*big.Int, error) {
	return _MockETHLINKAggregator.Contract.Version(&_MockETHLINKAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(uint256)
func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerSession) Version() (*big.Int, error) {
	return _MockETHLINKAggregator.Contract.Version(&_MockETHLINKAggregator.CallOpts)
}
