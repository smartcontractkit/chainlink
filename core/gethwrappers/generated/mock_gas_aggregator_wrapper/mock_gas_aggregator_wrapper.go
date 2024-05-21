// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mock_gas_aggregator_wrapper

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

var MockGASAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"answer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161025c38038061025c8339818101604052602081101561003357600080fd5b5051600055610215806100476000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c806385bb7d691161005057806385bb7d691461012c5780639a6fc8f514610134578063feaf968c1461019c57610072565b8063313ce5671461007757806354fd4d50146100955780637284e416146100af575b600080fd5b61007f6101a4565b6040805160ff9092168252519081900360200190f35b61009d6101a9565b60408051918252519081900360200190f35b6100b76101ae565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100f15781810151838201526020016100d9565b50505050905090810190601f16801561011e5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61009d6101e5565b61015d6004803603602081101561014a57600080fd5b503569ffffffffffffffffffff166101eb565b6040805169ffffffffffffffffffff96871681526020810195909552848101939093526060840191909152909216608082015290519081900360a00190f35b61015d6101fa565b601290565b600190565b60408051808201909152601181527f4d6f636b47415341676772656761746f72000000000000000000000000000000602082015290565b60005481565b50600190600090429081908490565b60016000428083909192939456fea164736f6c6343000606000a",
}

var MockGASAggregatorABI = MockGASAggregatorMetaData.ABI

var MockGASAggregatorBin = MockGASAggregatorMetaData.Bin

func DeployMockGASAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _answer *big.Int) (common.Address, *types.Transaction, *MockGASAggregator, error) {
	parsed, err := MockGASAggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockGASAggregatorBin), backend, _answer)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockGASAggregator{address: address, abi: *parsed, MockGASAggregatorCaller: MockGASAggregatorCaller{contract: contract}, MockGASAggregatorTransactor: MockGASAggregatorTransactor{contract: contract}, MockGASAggregatorFilterer: MockGASAggregatorFilterer{contract: contract}}, nil
}

type MockGASAggregator struct {
	address common.Address
	abi     abi.ABI
	MockGASAggregatorCaller
	MockGASAggregatorTransactor
	MockGASAggregatorFilterer
}

type MockGASAggregatorCaller struct {
	contract *bind.BoundContract
}

type MockGASAggregatorTransactor struct {
	contract *bind.BoundContract
}

type MockGASAggregatorFilterer struct {
	contract *bind.BoundContract
}

type MockGASAggregatorSession struct {
	Contract     *MockGASAggregator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MockGASAggregatorCallerSession struct {
	Contract *MockGASAggregatorCaller
	CallOpts bind.CallOpts
}

type MockGASAggregatorTransactorSession struct {
	Contract     *MockGASAggregatorTransactor
	TransactOpts bind.TransactOpts
}

type MockGASAggregatorRaw struct {
	Contract *MockGASAggregator
}

type MockGASAggregatorCallerRaw struct {
	Contract *MockGASAggregatorCaller
}

type MockGASAggregatorTransactorRaw struct {
	Contract *MockGASAggregatorTransactor
}

func NewMockGASAggregator(address common.Address, backend bind.ContractBackend) (*MockGASAggregator, error) {
	abi, err := abi.JSON(strings.NewReader(MockGASAggregatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMockGASAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockGASAggregator{address: address, abi: abi, MockGASAggregatorCaller: MockGASAggregatorCaller{contract: contract}, MockGASAggregatorTransactor: MockGASAggregatorTransactor{contract: contract}, MockGASAggregatorFilterer: MockGASAggregatorFilterer{contract: contract}}, nil
}

func NewMockGASAggregatorCaller(address common.Address, caller bind.ContractCaller) (*MockGASAggregatorCaller, error) {
	contract, err := bindMockGASAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockGASAggregatorCaller{contract: contract}, nil
}

func NewMockGASAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*MockGASAggregatorTransactor, error) {
	contract, err := bindMockGASAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockGASAggregatorTransactor{contract: contract}, nil
}

func NewMockGASAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*MockGASAggregatorFilterer, error) {
	contract, err := bindMockGASAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockGASAggregatorFilterer{contract: contract}, nil
}

func bindMockGASAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockGASAggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MockGASAggregator *MockGASAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockGASAggregator.Contract.MockGASAggregatorCaller.contract.Call(opts, result, method, params...)
}

func (_MockGASAggregator *MockGASAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockGASAggregator.Contract.MockGASAggregatorTransactor.contract.Transfer(opts)
}

func (_MockGASAggregator *MockGASAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockGASAggregator.Contract.MockGASAggregatorTransactor.contract.Transact(opts, method, params...)
}

func (_MockGASAggregator *MockGASAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockGASAggregator.Contract.contract.Call(opts, result, method, params...)
}

func (_MockGASAggregator *MockGASAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockGASAggregator.Contract.contract.Transfer(opts)
}

func (_MockGASAggregator *MockGASAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockGASAggregator.Contract.contract.Transact(opts, method, params...)
}

func (_MockGASAggregator *MockGASAggregatorCaller) Answer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockGASAggregator.contract.Call(opts, &out, "answer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockGASAggregator *MockGASAggregatorSession) Answer() (*big.Int, error) {
	return _MockGASAggregator.Contract.Answer(&_MockGASAggregator.CallOpts)
}

func (_MockGASAggregator *MockGASAggregatorCallerSession) Answer() (*big.Int, error) {
	return _MockGASAggregator.Contract.Answer(&_MockGASAggregator.CallOpts)
}

func (_MockGASAggregator *MockGASAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockGASAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_MockGASAggregator *MockGASAggregatorSession) Decimals() (uint8, error) {
	return _MockGASAggregator.Contract.Decimals(&_MockGASAggregator.CallOpts)
}

func (_MockGASAggregator *MockGASAggregatorCallerSession) Decimals() (uint8, error) {
	return _MockGASAggregator.Contract.Decimals(&_MockGASAggregator.CallOpts)
}

func (_MockGASAggregator *MockGASAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockGASAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MockGASAggregator *MockGASAggregatorSession) Description() (string, error) {
	return _MockGASAggregator.Contract.Description(&_MockGASAggregator.CallOpts)
}

func (_MockGASAggregator *MockGASAggregatorCallerSession) Description() (string, error) {
	return _MockGASAggregator.Contract.Description(&_MockGASAggregator.CallOpts)
}

func (_MockGASAggregator *MockGASAggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _MockGASAggregator.contract.Call(opts, &out, "getRoundData", _roundId)

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

func (_MockGASAggregator *MockGASAggregatorSession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _MockGASAggregator.Contract.GetRoundData(&_MockGASAggregator.CallOpts, _roundId)
}

func (_MockGASAggregator *MockGASAggregatorCallerSession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _MockGASAggregator.Contract.GetRoundData(&_MockGASAggregator.CallOpts, _roundId)
}

func (_MockGASAggregator *MockGASAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _MockGASAggregator.contract.Call(opts, &out, "latestRoundData")

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

func (_MockGASAggregator *MockGASAggregatorSession) LatestRoundData() (LatestRoundData,

	error) {
	return _MockGASAggregator.Contract.LatestRoundData(&_MockGASAggregator.CallOpts)
}

func (_MockGASAggregator *MockGASAggregatorCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _MockGASAggregator.Contract.LatestRoundData(&_MockGASAggregator.CallOpts)
}

func (_MockGASAggregator *MockGASAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockGASAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockGASAggregator *MockGASAggregatorSession) Version() (*big.Int, error) {
	return _MockGASAggregator.Contract.Version(&_MockGASAggregator.CallOpts)
}

func (_MockGASAggregator *MockGASAggregatorCallerSession) Version() (*big.Int, error) {
	return _MockGASAggregator.Contract.Version(&_MockGASAggregator.CallOpts)
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

func (_MockGASAggregator *MockGASAggregator) Address() common.Address {
	return _MockGASAggregator.address
}

type MockGASAggregatorInterface interface {
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
