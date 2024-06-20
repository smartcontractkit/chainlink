// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mock_usd_based_aggregator_wrapper

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

var MockUSDBasedAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"answer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"ans\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"ans\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_blockTimestampDeduction\",\"type\":\"uint256\"}],\"name\":\"setBlockTimestampDeduction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x6080604052600060015534801561001557600080fd5b506040516103333803806103338339810160408190526100349161003c565b600055610055565b60006020828403121561004e57600080fd5b5051919050565b6102cf806100646000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c806385bb7d691161005b57806385bb7d69146100e65780639a6fc8f5146100ef578063f0ad37df14610139578063feaf968c1461014e57600080fd5b8063313ce5671461008257806354fd4d50146100965780637284e416146100a7575b600080fd5b604051600881526020015b60405180910390f35b60015b60405190815260200161008d565b604080518082018252601681527f4d6f636b555344426173656441676772656761746f72000000000000000000006020820152905161008d91906101ca565b61009960005481565b6101026100fd366004610236565b610156565b6040805169ffffffffffffffffffff968716815260208101959095528401929092526060830152909116608082015260a00161008d565b61014c610147366004610269565b600155565b005b610102610186565b6000806000806000600160005461016b6101b5565b6101736101b5565b9299919850965090945060019350915050565b6000806000806000600160005461019b6101b5565b6101a36101b5565b92989197509550909350600192509050565b6000600154426101c59190610282565b905090565b600060208083528351808285015260005b818110156101f7578581018301518582016040015282016101db565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60006020828403121561024857600080fd5b813569ffffffffffffffffffff8116811461026257600080fd5b9392505050565b60006020828403121561027b57600080fd5b5035919050565b818103818111156102bc577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000813000a",
}

var MockUSDBasedAggregatorABI = MockUSDBasedAggregatorMetaData.ABI

var MockUSDBasedAggregatorBin = MockUSDBasedAggregatorMetaData.Bin

func DeployMockUSDBasedAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _answer *big.Int) (common.Address, *types.Transaction, *MockUSDBasedAggregator, error) {
	parsed, err := MockUSDBasedAggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockUSDBasedAggregatorBin), backend, _answer)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockUSDBasedAggregator{address: address, abi: *parsed, MockUSDBasedAggregatorCaller: MockUSDBasedAggregatorCaller{contract: contract}, MockUSDBasedAggregatorTransactor: MockUSDBasedAggregatorTransactor{contract: contract}, MockUSDBasedAggregatorFilterer: MockUSDBasedAggregatorFilterer{contract: contract}}, nil
}

type MockUSDBasedAggregator struct {
	address common.Address
	abi     abi.ABI
	MockUSDBasedAggregatorCaller
	MockUSDBasedAggregatorTransactor
	MockUSDBasedAggregatorFilterer
}

type MockUSDBasedAggregatorCaller struct {
	contract *bind.BoundContract
}

type MockUSDBasedAggregatorTransactor struct {
	contract *bind.BoundContract
}

type MockUSDBasedAggregatorFilterer struct {
	contract *bind.BoundContract
}

type MockUSDBasedAggregatorSession struct {
	Contract     *MockUSDBasedAggregator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MockUSDBasedAggregatorCallerSession struct {
	Contract *MockUSDBasedAggregatorCaller
	CallOpts bind.CallOpts
}

type MockUSDBasedAggregatorTransactorSession struct {
	Contract     *MockUSDBasedAggregatorTransactor
	TransactOpts bind.TransactOpts
}

type MockUSDBasedAggregatorRaw struct {
	Contract *MockUSDBasedAggregator
}

type MockUSDBasedAggregatorCallerRaw struct {
	Contract *MockUSDBasedAggregatorCaller
}

type MockUSDBasedAggregatorTransactorRaw struct {
	Contract *MockUSDBasedAggregatorTransactor
}

func NewMockUSDBasedAggregator(address common.Address, backend bind.ContractBackend) (*MockUSDBasedAggregator, error) {
	abi, err := abi.JSON(strings.NewReader(MockUSDBasedAggregatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMockUSDBasedAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockUSDBasedAggregator{address: address, abi: abi, MockUSDBasedAggregatorCaller: MockUSDBasedAggregatorCaller{contract: contract}, MockUSDBasedAggregatorTransactor: MockUSDBasedAggregatorTransactor{contract: contract}, MockUSDBasedAggregatorFilterer: MockUSDBasedAggregatorFilterer{contract: contract}}, nil
}

func NewMockUSDBasedAggregatorCaller(address common.Address, caller bind.ContractCaller) (*MockUSDBasedAggregatorCaller, error) {
	contract, err := bindMockUSDBasedAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockUSDBasedAggregatorCaller{contract: contract}, nil
}

func NewMockUSDBasedAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*MockUSDBasedAggregatorTransactor, error) {
	contract, err := bindMockUSDBasedAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockUSDBasedAggregatorTransactor{contract: contract}, nil
}

func NewMockUSDBasedAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*MockUSDBasedAggregatorFilterer, error) {
	contract, err := bindMockUSDBasedAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockUSDBasedAggregatorFilterer{contract: contract}, nil
}

func bindMockUSDBasedAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockUSDBasedAggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockUSDBasedAggregator.Contract.MockUSDBasedAggregatorCaller.contract.Call(opts, result, method, params...)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockUSDBasedAggregator.Contract.MockUSDBasedAggregatorTransactor.contract.Transfer(opts)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockUSDBasedAggregator.Contract.MockUSDBasedAggregatorTransactor.contract.Transact(opts, method, params...)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockUSDBasedAggregator.Contract.contract.Call(opts, result, method, params...)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockUSDBasedAggregator.Contract.contract.Transfer(opts)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockUSDBasedAggregator.Contract.contract.Transact(opts, method, params...)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorCaller) Answer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockUSDBasedAggregator.contract.Call(opts, &out, "answer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorSession) Answer() (*big.Int, error) {
	return _MockUSDBasedAggregator.Contract.Answer(&_MockUSDBasedAggregator.CallOpts)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorCallerSession) Answer() (*big.Int, error) {
	return _MockUSDBasedAggregator.Contract.Answer(&_MockUSDBasedAggregator.CallOpts)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockUSDBasedAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorSession) Decimals() (uint8, error) {
	return _MockUSDBasedAggregator.Contract.Decimals(&_MockUSDBasedAggregator.CallOpts)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorCallerSession) Decimals() (uint8, error) {
	return _MockUSDBasedAggregator.Contract.Decimals(&_MockUSDBasedAggregator.CallOpts)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockUSDBasedAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorSession) Description() (string, error) {
	return _MockUSDBasedAggregator.Contract.Description(&_MockUSDBasedAggregator.CallOpts)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorCallerSession) Description() (string, error) {
	return _MockUSDBasedAggregator.Contract.Description(&_MockUSDBasedAggregator.CallOpts)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorCaller) GetRoundData(opts *bind.CallOpts, arg0 *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _MockUSDBasedAggregator.contract.Call(opts, &out, "getRoundData", arg0)

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

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorSession) GetRoundData(arg0 *big.Int) (GetRoundData,

	error) {
	return _MockUSDBasedAggregator.Contract.GetRoundData(&_MockUSDBasedAggregator.CallOpts, arg0)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorCallerSession) GetRoundData(arg0 *big.Int) (GetRoundData,

	error) {
	return _MockUSDBasedAggregator.Contract.GetRoundData(&_MockUSDBasedAggregator.CallOpts, arg0)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _MockUSDBasedAggregator.contract.Call(opts, &out, "latestRoundData")

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

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorSession) LatestRoundData() (LatestRoundData,

	error) {
	return _MockUSDBasedAggregator.Contract.LatestRoundData(&_MockUSDBasedAggregator.CallOpts)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _MockUSDBasedAggregator.Contract.LatestRoundData(&_MockUSDBasedAggregator.CallOpts)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockUSDBasedAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorSession) Version() (*big.Int, error) {
	return _MockUSDBasedAggregator.Contract.Version(&_MockUSDBasedAggregator.CallOpts)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorCallerSession) Version() (*big.Int, error) {
	return _MockUSDBasedAggregator.Contract.Version(&_MockUSDBasedAggregator.CallOpts)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorTransactor) SetBlockTimestampDeduction(opts *bind.TransactOpts, _blockTimestampDeduction *big.Int) (*types.Transaction, error) {
	return _MockUSDBasedAggregator.contract.Transact(opts, "setBlockTimestampDeduction", _blockTimestampDeduction)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorSession) SetBlockTimestampDeduction(_blockTimestampDeduction *big.Int) (*types.Transaction, error) {
	return _MockUSDBasedAggregator.Contract.SetBlockTimestampDeduction(&_MockUSDBasedAggregator.TransactOpts, _blockTimestampDeduction)
}

func (_MockUSDBasedAggregator *MockUSDBasedAggregatorTransactorSession) SetBlockTimestampDeduction(_blockTimestampDeduction *big.Int) (*types.Transaction, error) {
	return _MockUSDBasedAggregator.Contract.SetBlockTimestampDeduction(&_MockUSDBasedAggregator.TransactOpts, _blockTimestampDeduction)
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

func (_MockUSDBasedAggregator *MockUSDBasedAggregator) Address() common.Address {
	return _MockUSDBasedAggregator.address
}

type MockUSDBasedAggregatorInterface interface {
	Answer(opts *bind.CallOpts) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetRoundData(opts *bind.CallOpts, arg0 *big.Int) (GetRoundData,

		error)

	LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

		error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	SetBlockTimestampDeduction(opts *bind.TransactOpts, _blockTimestampDeduction *big.Int) (*types.Transaction, error)

	Address() common.Address
}
