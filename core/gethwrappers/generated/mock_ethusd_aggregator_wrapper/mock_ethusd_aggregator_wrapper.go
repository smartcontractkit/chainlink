// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mock_ethusd_aggregator_wrapper

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

var MockETHUSDAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_answer\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"answer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"description\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getRoundData\",\"inputs\":[{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"outputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"ans\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"answeredInRound\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRoundData\",\"inputs\":[],\"outputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"ans\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"answeredInRound\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setBlockTimestampDeduction\",\"inputs\":[{\"name\":\"_blockTimestampDeduction\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"}]",
	Bin: "0x6080604052600060015534801561001557600080fd5b506040516103333803806103338339810160408190526100349161003c565b600055610055565b60006020828403121561004e57600080fd5b5051919050565b6102cf806100646000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c806385bb7d691161005b57806385bb7d69146100e65780639a6fc8f5146100ef578063f0ad37df14610139578063feaf968c1461014e57600080fd5b8063313ce5671461008257806354fd4d50146100965780637284e416146100a7575b600080fd5b604051600881526020015b60405180910390f35b60015b60405190815260200161008d565b604080518082018252601481527f4d6f636b45544855534441676772656761746f720000000000000000000000006020820152905161008d91906101ca565b61009960005481565b6101026100fd366004610236565b610156565b6040805169ffffffffffffffffffff968716815260208101959095528401929092526060830152909116608082015260a00161008d565b61014c610147366004610269565b600155565b005b610102610186565b6000806000806000600160005461016b6101b5565b6101736101b5565b9299919850965090945060019350915050565b6000806000806000600160005461019b6101b5565b6101a36101b5565b92989197509550909350600192509050565b6000600154426101c59190610282565b905090565b600060208083528351808285015260005b818110156101f7578581018301518582016040015282016101db565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60006020828403121561024857600080fd5b813569ffffffffffffffffffff8116811461026257600080fd5b9392505050565b60006020828403121561027b57600080fd5b5035919050565b818103818111156102bc577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000813000a",
}

var MockETHUSDAggregatorABI = MockETHUSDAggregatorMetaData.ABI

var MockETHUSDAggregatorBin = MockETHUSDAggregatorMetaData.Bin

func DeployMockETHUSDAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _answer *big.Int) (common.Address, *types.Transaction, *MockETHUSDAggregator, error) {
	parsed, err := MockETHUSDAggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockETHUSDAggregatorBin), backend, _answer)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockETHUSDAggregator{address: address, abi: *parsed, MockETHUSDAggregatorCaller: MockETHUSDAggregatorCaller{contract: contract}, MockETHUSDAggregatorTransactor: MockETHUSDAggregatorTransactor{contract: contract}, MockETHUSDAggregatorFilterer: MockETHUSDAggregatorFilterer{contract: contract}}, nil
}

type MockETHUSDAggregator struct {
	address common.Address
	abi     abi.ABI
	MockETHUSDAggregatorCaller
	MockETHUSDAggregatorTransactor
	MockETHUSDAggregatorFilterer
}

type MockETHUSDAggregatorCaller struct {
	contract *bind.BoundContract
}

type MockETHUSDAggregatorTransactor struct {
	contract *bind.BoundContract
}

type MockETHUSDAggregatorFilterer struct {
	contract *bind.BoundContract
}

type MockETHUSDAggregatorSession struct {
	Contract     *MockETHUSDAggregator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MockETHUSDAggregatorCallerSession struct {
	Contract *MockETHUSDAggregatorCaller
	CallOpts bind.CallOpts
}

type MockETHUSDAggregatorTransactorSession struct {
	Contract     *MockETHUSDAggregatorTransactor
	TransactOpts bind.TransactOpts
}

type MockETHUSDAggregatorRaw struct {
	Contract *MockETHUSDAggregator
}

type MockETHUSDAggregatorCallerRaw struct {
	Contract *MockETHUSDAggregatorCaller
}

type MockETHUSDAggregatorTransactorRaw struct {
	Contract *MockETHUSDAggregatorTransactor
}

func NewMockETHUSDAggregator(address common.Address, backend bind.ContractBackend) (*MockETHUSDAggregator, error) {
	abi, err := abi.JSON(strings.NewReader(MockETHUSDAggregatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMockETHUSDAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockETHUSDAggregator{address: address, abi: abi, MockETHUSDAggregatorCaller: MockETHUSDAggregatorCaller{contract: contract}, MockETHUSDAggregatorTransactor: MockETHUSDAggregatorTransactor{contract: contract}, MockETHUSDAggregatorFilterer: MockETHUSDAggregatorFilterer{contract: contract}}, nil
}

func NewMockETHUSDAggregatorCaller(address common.Address, caller bind.ContractCaller) (*MockETHUSDAggregatorCaller, error) {
	contract, err := bindMockETHUSDAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockETHUSDAggregatorCaller{contract: contract}, nil
}

func NewMockETHUSDAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*MockETHUSDAggregatorTransactor, error) {
	contract, err := bindMockETHUSDAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockETHUSDAggregatorTransactor{contract: contract}, nil
}

func NewMockETHUSDAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*MockETHUSDAggregatorFilterer, error) {
	contract, err := bindMockETHUSDAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockETHUSDAggregatorFilterer{contract: contract}, nil
}

func bindMockETHUSDAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockETHUSDAggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockETHUSDAggregator.Contract.MockETHUSDAggregatorCaller.contract.Call(opts, result, method, params...)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockETHUSDAggregator.Contract.MockETHUSDAggregatorTransactor.contract.Transfer(opts)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockETHUSDAggregator.Contract.MockETHUSDAggregatorTransactor.contract.Transact(opts, method, params...)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockETHUSDAggregator.Contract.contract.Call(opts, result, method, params...)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockETHUSDAggregator.Contract.contract.Transfer(opts)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockETHUSDAggregator.Contract.contract.Transact(opts, method, params...)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorCaller) Answer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockETHUSDAggregator.contract.Call(opts, &out, "answer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockETHUSDAggregator *MockETHUSDAggregatorSession) Answer() (*big.Int, error) {
	return _MockETHUSDAggregator.Contract.Answer(&_MockETHUSDAggregator.CallOpts)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorCallerSession) Answer() (*big.Int, error) {
	return _MockETHUSDAggregator.Contract.Answer(&_MockETHUSDAggregator.CallOpts)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockETHUSDAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_MockETHUSDAggregator *MockETHUSDAggregatorSession) Decimals() (uint8, error) {
	return _MockETHUSDAggregator.Contract.Decimals(&_MockETHUSDAggregator.CallOpts)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorCallerSession) Decimals() (uint8, error) {
	return _MockETHUSDAggregator.Contract.Decimals(&_MockETHUSDAggregator.CallOpts)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockETHUSDAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MockETHUSDAggregator *MockETHUSDAggregatorSession) Description() (string, error) {
	return _MockETHUSDAggregator.Contract.Description(&_MockETHUSDAggregator.CallOpts)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorCallerSession) Description() (string, error) {
	return _MockETHUSDAggregator.Contract.Description(&_MockETHUSDAggregator.CallOpts)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorCaller) GetRoundData(opts *bind.CallOpts, arg0 *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _MockETHUSDAggregator.contract.Call(opts, &out, "getRoundData", arg0)

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

func (_MockETHUSDAggregator *MockETHUSDAggregatorSession) GetRoundData(arg0 *big.Int) (GetRoundData,

	error) {
	return _MockETHUSDAggregator.Contract.GetRoundData(&_MockETHUSDAggregator.CallOpts, arg0)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorCallerSession) GetRoundData(arg0 *big.Int) (GetRoundData,

	error) {
	return _MockETHUSDAggregator.Contract.GetRoundData(&_MockETHUSDAggregator.CallOpts, arg0)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _MockETHUSDAggregator.contract.Call(opts, &out, "latestRoundData")

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

func (_MockETHUSDAggregator *MockETHUSDAggregatorSession) LatestRoundData() (LatestRoundData,

	error) {
	return _MockETHUSDAggregator.Contract.LatestRoundData(&_MockETHUSDAggregator.CallOpts)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _MockETHUSDAggregator.Contract.LatestRoundData(&_MockETHUSDAggregator.CallOpts)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockETHUSDAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockETHUSDAggregator *MockETHUSDAggregatorSession) Version() (*big.Int, error) {
	return _MockETHUSDAggregator.Contract.Version(&_MockETHUSDAggregator.CallOpts)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorCallerSession) Version() (*big.Int, error) {
	return _MockETHUSDAggregator.Contract.Version(&_MockETHUSDAggregator.CallOpts)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorTransactor) SetBlockTimestampDeduction(opts *bind.TransactOpts, _blockTimestampDeduction *big.Int) (*types.Transaction, error) {
	return _MockETHUSDAggregator.contract.Transact(opts, "setBlockTimestampDeduction", _blockTimestampDeduction)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorSession) SetBlockTimestampDeduction(_blockTimestampDeduction *big.Int) (*types.Transaction, error) {
	return _MockETHUSDAggregator.Contract.SetBlockTimestampDeduction(&_MockETHUSDAggregator.TransactOpts, _blockTimestampDeduction)
}

func (_MockETHUSDAggregator *MockETHUSDAggregatorTransactorSession) SetBlockTimestampDeduction(_blockTimestampDeduction *big.Int) (*types.Transaction, error) {
	return _MockETHUSDAggregator.Contract.SetBlockTimestampDeduction(&_MockETHUSDAggregator.TransactOpts, _blockTimestampDeduction)
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

func (_MockETHUSDAggregator *MockETHUSDAggregator) Address() common.Address {
	return _MockETHUSDAggregator.address
}

type MockETHUSDAggregatorInterface interface {
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
