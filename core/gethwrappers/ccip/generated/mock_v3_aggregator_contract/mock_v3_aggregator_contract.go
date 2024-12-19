// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mock_v3_aggregator_contract

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
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

var MockV3AggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"int256\",\"name\":\"_initialAnswer\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"name\":\"updateAnswer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_startedAt\",\"type\":\"uint256\"}],\"name\":\"updateRoundData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60803460c857601f6106f738819003918201601f19168301916001600160401b0383118484101760cd57808492604094855283398101031260c85780519060ff821680920360c857602001519060ff1960005416176000558060015542600255600354600019811460b25760010180600355600052600460205260406000205560035460005260056020524260406000205560035460005260066020524260406000205560405161061390816100e48239f35b634e487b7160e01b600052601160045260246000fd5b600080fd5b634e487b7160e01b600052604160045260246000fdfe608080604052600436101561001357600080fd5b60003560e01c908163313ce567146105b1575080634aa2011f1461052357806350d25bcd146104e757806354fd4d50146104ad578063668a0f02146104715780637284e4161461035f5780638205bf6a146103235780639a6fc8f514610295578063a87a20ce146101c4578063b5ab58dc1461017a578063b633620c146101305763feaf968c146100a357600080fd5b3461012b5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012b576003546000818152600460209081526040808320546006835281842054600584529382902054825169ffffffffffffffffffff90961680875293860191909152908401929092526060830191909152608082015260a090f35b600080fd5b3461012b5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012b5760043560005260056020526020604060002054604051908152f35b3461012b5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012b5760043560005260046020526020604060002054604051908152f35b3461012b5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012b5760043580600155426002556003547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811461026657600101806003556000526004602052604060002055600354600052600560205242604060002055600354600052600660205242604060002055600080f35b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b3461012b5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012b576102cc6105ed565b69ffffffffffffffffffff166000818152600460209081526040808320546006835281842054600584529382902054825186815293840191909152908201929092526060810191909152608081019190915260a090f35b3461012b5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012b576020600254604051908152f35b3461012b5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012b576040516040810181811067ffffffffffffffff82111761044257604052601f81527f76302e382f74657374732f4d6f636b563341676772656761746f722e736f6c00602082015260405190602082528181519182602083015260005b83811061042a5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f836000604080968601015201168101030190f35b602082820181015160408784010152859350016103ea565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b3461012b5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012b576020600354604051908152f35b3461012b5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012b57602060405160008152f35b3461012b5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012b576020600154604051908152f35b3461012b5760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012b5761055a6105ed565b60243569ffffffffffffffffffff6044359216806003558160015582600255600052600460205260406000205560035460005260056020526040600020556003546000526006602052606435604060002055600080f35b3461012b5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012b5760209060ff600054168152f35b6004359069ffffffffffffffffffff8216820361012b5756fea164736f6c634300081a000a",
}

var MockV3AggregatorABI = MockV3AggregatorMetaData.ABI

var MockV3AggregatorBin = MockV3AggregatorMetaData.Bin

func DeployMockV3Aggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _decimals uint8, _initialAnswer *big.Int) (common.Address, *types.Transaction, *MockV3Aggregator, error) {
	parsed, err := MockV3AggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockV3AggregatorBin), backend, _decimals, _initialAnswer)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockV3Aggregator{address: address, abi: *parsed, MockV3AggregatorCaller: MockV3AggregatorCaller{contract: contract}, MockV3AggregatorTransactor: MockV3AggregatorTransactor{contract: contract}, MockV3AggregatorFilterer: MockV3AggregatorFilterer{contract: contract}}, nil
}

type MockV3Aggregator struct {
	address common.Address
	abi     abi.ABI
	MockV3AggregatorCaller
	MockV3AggregatorTransactor
	MockV3AggregatorFilterer
}

type MockV3AggregatorCaller struct {
	contract *bind.BoundContract
}

type MockV3AggregatorTransactor struct {
	contract *bind.BoundContract
}

type MockV3AggregatorFilterer struct {
	contract *bind.BoundContract
}

type MockV3AggregatorSession struct {
	Contract     *MockV3Aggregator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MockV3AggregatorCallerSession struct {
	Contract *MockV3AggregatorCaller
	CallOpts bind.CallOpts
}

type MockV3AggregatorTransactorSession struct {
	Contract     *MockV3AggregatorTransactor
	TransactOpts bind.TransactOpts
}

type MockV3AggregatorRaw struct {
	Contract *MockV3Aggregator
}

type MockV3AggregatorCallerRaw struct {
	Contract *MockV3AggregatorCaller
}

type MockV3AggregatorTransactorRaw struct {
	Contract *MockV3AggregatorTransactor
}

func NewMockV3Aggregator(address common.Address, backend bind.ContractBackend) (*MockV3Aggregator, error) {
	abi, err := abi.JSON(strings.NewReader(MockV3AggregatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMockV3Aggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockV3Aggregator{address: address, abi: abi, MockV3AggregatorCaller: MockV3AggregatorCaller{contract: contract}, MockV3AggregatorTransactor: MockV3AggregatorTransactor{contract: contract}, MockV3AggregatorFilterer: MockV3AggregatorFilterer{contract: contract}}, nil
}

func NewMockV3AggregatorCaller(address common.Address, caller bind.ContractCaller) (*MockV3AggregatorCaller, error) {
	contract, err := bindMockV3Aggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorCaller{contract: contract}, nil
}

func NewMockV3AggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*MockV3AggregatorTransactor, error) {
	contract, err := bindMockV3Aggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorTransactor{contract: contract}, nil
}

func NewMockV3AggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*MockV3AggregatorFilterer, error) {
	contract, err := bindMockV3Aggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorFilterer{contract: contract}, nil
}

func bindMockV3Aggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockV3AggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MockV3Aggregator *MockV3AggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockV3Aggregator.Contract.MockV3AggregatorCaller.contract.Call(opts, result, method, params...)
}

func (_MockV3Aggregator *MockV3AggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.MockV3AggregatorTransactor.contract.Transfer(opts)
}

func (_MockV3Aggregator *MockV3AggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.MockV3AggregatorTransactor.contract.Transact(opts, method, params...)
}

func (_MockV3Aggregator *MockV3AggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockV3Aggregator.Contract.contract.Call(opts, result, method, params...)
}

func (_MockV3Aggregator *MockV3AggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.contract.Transfer(opts)
}

func (_MockV3Aggregator *MockV3AggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.contract.Transact(opts, method, params...)
}

func (_MockV3Aggregator *MockV3AggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_MockV3Aggregator *MockV3AggregatorSession) Decimals() (uint8, error) {
	return _MockV3Aggregator.Contract.Decimals(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorCallerSession) Decimals() (uint8, error) {
	return _MockV3Aggregator.Contract.Decimals(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MockV3Aggregator *MockV3AggregatorSession) Description() (string, error) {
	return _MockV3Aggregator.Contract.Description(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorCallerSession) Description() (string, error) {
	return _MockV3Aggregator.Contract.Description(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorCaller) GetAnswer(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "getAnswer", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockV3Aggregator *MockV3AggregatorSession) GetAnswer(arg0 *big.Int) (*big.Int, error) {
	return _MockV3Aggregator.Contract.GetAnswer(&_MockV3Aggregator.CallOpts, arg0)
}

func (_MockV3Aggregator *MockV3AggregatorCallerSession) GetAnswer(arg0 *big.Int) (*big.Int, error) {
	return _MockV3Aggregator.Contract.GetAnswer(&_MockV3Aggregator.CallOpts, arg0)
}

func (_MockV3Aggregator *MockV3AggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "getRoundData", _roundId)

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

func (_MockV3Aggregator *MockV3AggregatorSession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _MockV3Aggregator.Contract.GetRoundData(&_MockV3Aggregator.CallOpts, _roundId)
}

func (_MockV3Aggregator *MockV3AggregatorCallerSession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _MockV3Aggregator.Contract.GetRoundData(&_MockV3Aggregator.CallOpts, _roundId)
}

func (_MockV3Aggregator *MockV3AggregatorCaller) GetTimestamp(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "getTimestamp", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockV3Aggregator *MockV3AggregatorSession) GetTimestamp(arg0 *big.Int) (*big.Int, error) {
	return _MockV3Aggregator.Contract.GetTimestamp(&_MockV3Aggregator.CallOpts, arg0)
}

func (_MockV3Aggregator *MockV3AggregatorCallerSession) GetTimestamp(arg0 *big.Int) (*big.Int, error) {
	return _MockV3Aggregator.Contract.GetTimestamp(&_MockV3Aggregator.CallOpts, arg0)
}

func (_MockV3Aggregator *MockV3AggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockV3Aggregator *MockV3AggregatorSession) LatestAnswer() (*big.Int, error) {
	return _MockV3Aggregator.Contract.LatestAnswer(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _MockV3Aggregator.Contract.LatestAnswer(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockV3Aggregator *MockV3AggregatorSession) LatestRound() (*big.Int, error) {
	return _MockV3Aggregator.Contract.LatestRound(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _MockV3Aggregator.Contract.LatestRound(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "latestRoundData")

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

func (_MockV3Aggregator *MockV3AggregatorSession) LatestRoundData() (LatestRoundData,

	error) {
	return _MockV3Aggregator.Contract.LatestRoundData(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _MockV3Aggregator.Contract.LatestRoundData(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockV3Aggregator *MockV3AggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _MockV3Aggregator.Contract.LatestTimestamp(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _MockV3Aggregator.Contract.LatestTimestamp(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MockV3Aggregator *MockV3AggregatorSession) Version() (*big.Int, error) {
	return _MockV3Aggregator.Contract.Version(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorCallerSession) Version() (*big.Int, error) {
	return _MockV3Aggregator.Contract.Version(&_MockV3Aggregator.CallOpts)
}

func (_MockV3Aggregator *MockV3AggregatorTransactor) UpdateAnswer(opts *bind.TransactOpts, _answer *big.Int) (*types.Transaction, error) {
	return _MockV3Aggregator.contract.Transact(opts, "updateAnswer", _answer)
}

func (_MockV3Aggregator *MockV3AggregatorSession) UpdateAnswer(_answer *big.Int) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.UpdateAnswer(&_MockV3Aggregator.TransactOpts, _answer)
}

func (_MockV3Aggregator *MockV3AggregatorTransactorSession) UpdateAnswer(_answer *big.Int) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.UpdateAnswer(&_MockV3Aggregator.TransactOpts, _answer)
}

func (_MockV3Aggregator *MockV3AggregatorTransactor) UpdateRoundData(opts *bind.TransactOpts, _roundId *big.Int, _answer *big.Int, _timestamp *big.Int, _startedAt *big.Int) (*types.Transaction, error) {
	return _MockV3Aggregator.contract.Transact(opts, "updateRoundData", _roundId, _answer, _timestamp, _startedAt)
}

func (_MockV3Aggregator *MockV3AggregatorSession) UpdateRoundData(_roundId *big.Int, _answer *big.Int, _timestamp *big.Int, _startedAt *big.Int) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.UpdateRoundData(&_MockV3Aggregator.TransactOpts, _roundId, _answer, _timestamp, _startedAt)
}

func (_MockV3Aggregator *MockV3AggregatorTransactorSession) UpdateRoundData(_roundId *big.Int, _answer *big.Int, _timestamp *big.Int, _startedAt *big.Int) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.UpdateRoundData(&_MockV3Aggregator.TransactOpts, _roundId, _answer, _timestamp, _startedAt)
}

type MockV3AggregatorAnswerUpdatedIterator struct {
	Event *MockV3AggregatorAnswerUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockV3AggregatorAnswerUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockV3AggregatorAnswerUpdated)
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

	select {
	case log := <-it.logs:
		it.Event = new(MockV3AggregatorAnswerUpdated)
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

func (it *MockV3AggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

func (it *MockV3AggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockV3AggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log
}

func (_MockV3Aggregator *MockV3AggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*MockV3AggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _MockV3Aggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorAnswerUpdatedIterator{contract: _MockV3Aggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

func (_MockV3Aggregator *MockV3AggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *MockV3AggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _MockV3Aggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockV3AggregatorAnswerUpdated)
				if err := _MockV3Aggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

func (_MockV3Aggregator *MockV3AggregatorFilterer) ParseAnswerUpdated(log types.Log) (*MockV3AggregatorAnswerUpdated, error) {
	event := new(MockV3AggregatorAnswerUpdated)
	if err := _MockV3Aggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MockV3AggregatorNewRoundIterator struct {
	Event *MockV3AggregatorNewRound

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockV3AggregatorNewRoundIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockV3AggregatorNewRound)
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

	select {
	case log := <-it.logs:
		it.Event = new(MockV3AggregatorNewRound)
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

func (it *MockV3AggregatorNewRoundIterator) Error() error {
	return it.fail
}

func (it *MockV3AggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockV3AggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log
}

func (_MockV3Aggregator *MockV3AggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*MockV3AggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _MockV3Aggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorNewRoundIterator{contract: _MockV3Aggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

func (_MockV3Aggregator *MockV3AggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *MockV3AggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _MockV3Aggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockV3AggregatorNewRound)
				if err := _MockV3Aggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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

func (_MockV3Aggregator *MockV3AggregatorFilterer) ParseNewRound(log types.Log) (*MockV3AggregatorNewRound, error) {
	event := new(MockV3AggregatorNewRound)
	if err := _MockV3Aggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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

func (_MockV3Aggregator *MockV3Aggregator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MockV3Aggregator.abi.Events["AnswerUpdated"].ID:
		return _MockV3Aggregator.ParseAnswerUpdated(log)
	case _MockV3Aggregator.abi.Events["NewRound"].ID:
		return _MockV3Aggregator.ParseNewRound(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MockV3AggregatorAnswerUpdated) Topic() common.Hash {
	return common.HexToHash("0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f")
}

func (MockV3AggregatorNewRound) Topic() common.Hash {
	return common.HexToHash("0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271")
}

func (_MockV3Aggregator *MockV3Aggregator) Address() common.Address {
	return _MockV3Aggregator.address
}

type MockV3AggregatorInterface interface {
	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetAnswer(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

		error)

	GetTimestamp(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	LatestAnswer(opts *bind.CallOpts) (*big.Int, error)

	LatestRound(opts *bind.CallOpts) (*big.Int, error)

	LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

		error)

	LatestTimestamp(opts *bind.CallOpts) (*big.Int, error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	UpdateAnswer(opts *bind.TransactOpts, _answer *big.Int) (*types.Transaction, error)

	UpdateRoundData(opts *bind.TransactOpts, _roundId *big.Int, _answer *big.Int, _timestamp *big.Int, _startedAt *big.Int) (*types.Transaction, error)

	FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*MockV3AggregatorAnswerUpdatedIterator, error)

	WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *MockV3AggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error)

	ParseAnswerUpdated(log types.Log) (*MockV3AggregatorAnswerUpdated, error)

	FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*MockV3AggregatorNewRoundIterator, error)

	WatchNewRound(opts *bind.WatchOpts, sink chan<- *MockV3AggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error)

	ParseNewRound(log types.Log) (*MockV3AggregatorNewRound, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
