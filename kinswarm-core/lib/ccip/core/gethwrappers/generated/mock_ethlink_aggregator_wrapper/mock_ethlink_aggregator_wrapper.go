// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mock_ethlink_aggregator_wrapper

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

var MockETHLINKAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"answer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"ans\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"ans\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161025c38038061025c8339818101604052602081101561003357600080fd5b5051600055610215806100476000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c806385bb7d691161005057806385bb7d691461012c5780639a6fc8f514610134578063feaf968c1461019c57610072565b8063313ce5671461007757806354fd4d50146100955780637284e416146100af575b600080fd5b61007f6101a4565b6040805160ff9092168252519081900360200190f35b61009d6101a9565b60408051918252519081900360200190f35b6100b76101ae565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100f15781810151838201526020016100d9565b50505050905090810190601f16801561011e5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61009d6101e5565b61015d6004803603602081101561014a57600080fd5b503569ffffffffffffffffffff166101eb565b6040805169ffffffffffffffffffff96871681526020810195909552848101939093526060840191909152909216608082015290519081900360a00190f35b61015d6101fa565b601290565b600190565b60408051808201909152601581527f4d6f636b4554484c494e4b41676772656761746f720000000000000000000000602082015290565b60005481565b50600054600191429081908490565b60005460019142908190849056fea164736f6c6343000606000a",
}

var MockETHLINKAggregatorABI = MockETHLINKAggregatorMetaData.ABI

var MockETHLINKAggregatorBin = MockETHLINKAggregatorMetaData.Bin

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
	return address, tx, &MockETHLINKAggregator{address: address, abi: *parsed, MockETHLINKAggregatorCaller: MockETHLINKAggregatorCaller{contract: contract}, MockETHLINKAggregatorTransactor: MockETHLINKAggregatorTransactor{contract: contract}, MockETHLINKAggregatorFilterer: MockETHLINKAggregatorFilterer{contract: contract}}, nil
}

type MockETHLINKAggregator struct {
	address common.Address
	abi     abi.ABI
	MockETHLINKAggregatorCaller
	MockETHLINKAggregatorTransactor
	MockETHLINKAggregatorFilterer
}

type MockETHLINKAggregatorCaller struct {
	contract *bind.BoundContract
}

type MockETHLINKAggregatorTransactor struct {
	contract *bind.BoundContract
}

type MockETHLINKAggregatorFilterer struct {
	contract *bind.BoundContract
}

type MockETHLINKAggregatorSession struct {
	Contract     *MockETHLINKAggregator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MockETHLINKAggregatorCallerSession struct {
	Contract *MockETHLINKAggregatorCaller
	CallOpts bind.CallOpts
}

type MockETHLINKAggregatorTransactorSession struct {
	Contract     *MockETHLINKAggregatorTransactor
	TransactOpts bind.TransactOpts
}

type MockETHLINKAggregatorRaw struct {
	Contract *MockETHLINKAggregator
}

type MockETHLINKAggregatorCallerRaw struct {
	Contract *MockETHLINKAggregatorCaller
}

type MockETHLINKAggregatorTransactorRaw struct {
	Contract *MockETHLINKAggregatorTransactor
}

func NewMockETHLINKAggregator(address common.Address, backend bind.ContractBackend) (*MockETHLINKAggregator, error) {
	abi, err := abi.JSON(strings.NewReader(MockETHLINKAggregatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMockETHLINKAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockETHLINKAggregator{address: address, abi: abi, MockETHLINKAggregatorCaller: MockETHLINKAggregatorCaller{contract: contract}, MockETHLINKAggregatorTransactor: MockETHLINKAggregatorTransactor{contract: contract}, MockETHLINKAggregatorFilterer: MockETHLINKAggregatorFilterer{contract: contract}}, nil
}

func NewMockETHLINKAggregatorCaller(address common.Address, caller bind.ContractCaller) (*MockETHLINKAggregatorCaller, error) {
	contract, err := bindMockETHLINKAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockETHLINKAggregatorCaller{contract: contract}, nil
}

func NewMockETHLINKAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*MockETHLINKAggregatorTransactor, error) {
	contract, err := bindMockETHLINKAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockETHLINKAggregatorTransactor{contract: contract}, nil
}

func NewMockETHLINKAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*MockETHLINKAggregatorFilterer, error) {
	contract, err := bindMockETHLINKAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockETHLINKAggregatorFilterer{contract: contract}, nil
}

func bindMockETHLINKAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockETHLINKAggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockETHLINKAggregator.Contract.MockETHLINKAggregatorCaller.contract.Call(opts, result, method, params...)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockETHLINKAggregator.Contract.MockETHLINKAggregatorTransactor.contract.Transfer(opts)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockETHLINKAggregator.Contract.MockETHLINKAggregatorTransactor.contract.Transact(opts, method, params...)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockETHLINKAggregator.Contract.contract.Call(opts, result, method, params...)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockETHLINKAggregator.Contract.contract.Transfer(opts)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockETHLINKAggregator.Contract.contract.Transact(opts, method, params...)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorCaller) Answer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockETHLINKAggregator.contract.Call(opts, &out, "answer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockETHLINKAggregator *MockETHLINKAggregatorSession) Answer() (*big.Int, error) {
	return _MockETHLINKAggregator.Contract.Answer(&_MockETHLINKAggregator.CallOpts)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerSession) Answer() (*big.Int, error) {
	return _MockETHLINKAggregator.Contract.Answer(&_MockETHLINKAggregator.CallOpts)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockETHLINKAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_MockETHLINKAggregator *MockETHLINKAggregatorSession) Decimals() (uint8, error) {
	return _MockETHLINKAggregator.Contract.Decimals(&_MockETHLINKAggregator.CallOpts)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerSession) Decimals() (uint8, error) {
	return _MockETHLINKAggregator.Contract.Decimals(&_MockETHLINKAggregator.CallOpts)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockETHLINKAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MockETHLINKAggregator *MockETHLINKAggregatorSession) Description() (string, error) {
	return _MockETHLINKAggregator.Contract.Description(&_MockETHLINKAggregator.CallOpts)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerSession) Description() (string, error) {
	return _MockETHLINKAggregator.Contract.Description(&_MockETHLINKAggregator.CallOpts)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _MockETHLINKAggregator.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(GetRoundData)
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

func (_MockETHLINKAggregator *MockETHLINKAggregatorSession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _MockETHLINKAggregator.Contract.GetRoundData(&_MockETHLINKAggregator.CallOpts, _roundId)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerSession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _MockETHLINKAggregator.Contract.GetRoundData(&_MockETHLINKAggregator.CallOpts, _roundId)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _MockETHLINKAggregator.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(LatestRoundData)
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

func (_MockETHLINKAggregator *MockETHLINKAggregatorSession) LatestRoundData() (LatestRoundData,

	error) {
	return _MockETHLINKAggregator.Contract.LatestRoundData(&_MockETHLINKAggregator.CallOpts)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _MockETHLINKAggregator.Contract.LatestRoundData(&_MockETHLINKAggregator.CallOpts)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockETHLINKAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockETHLINKAggregator *MockETHLINKAggregatorSession) Version() (*big.Int, error) {
	return _MockETHLINKAggregator.Contract.Version(&_MockETHLINKAggregator.CallOpts)
}

func (_MockETHLINKAggregator *MockETHLINKAggregatorCallerSession) Version() (*big.Int, error) {
	return _MockETHLINKAggregator.Contract.Version(&_MockETHLINKAggregator.CallOpts)
}

type GetRoundData struct {
	RoundId         *big.Int
	Ans             *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
type LatestRoundData struct {
	RoundId         *big.Int
	Ans             *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

func (_MockETHLINKAggregator *MockETHLINKAggregator) Address() common.Address {
	return _MockETHLINKAggregator.address
}

type MockETHLINKAggregatorInterface interface {
	Answer(opts *bind.CallOpts) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

		error)

	LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

		error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	Address() common.Address
}
