// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package derived_price_feed_wrapper

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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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
)

var DerivedPriceFeedMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_base\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_quote\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"base\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"quote\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516107c43803806107c48339818101604052606081101561003357600080fd5b508051602082015160409092015190919060ff81161580159061005a5750601260ff821611155b61009f576040805162461bcd60e51b8152602060048201526011602482015270496e76616c6964205f646563696d616c7360781b604482015290519081900360640190fd5b6000805460ff90921660ff19909216919091179055600780546001600160a01b039384166001600160a01b031991821617909155600880549290931691161790556106d5806100ef6000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c80638205bf6a11610081578063b5ab58dc1161005b578063b5ab58dc14610252578063b633620c1461026f578063feaf968c1461028c576100d4565b80638205bf6a146101cf578063999b93af146101d75780639a6fc8f5146101df576100d4565b806354fd4d50116100b257806354fd4d5014610142578063668a0f021461014a5780637284e41614610152576100d4565b8063313ce567146100d95780635001f3b5146100f757806350d25bcd14610128575b600080fd5b6100e1610294565b6040805160ff9092168252519081900360200190f35b6100ff61029d565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6101306102b9565b60408051918252519081900360200190f35b6101306102bf565b6101306102c4565b61015a6102ca565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561019457818101518382015260200161017c565b50505050905090810190601f1680156101c15780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610130610301565b6100ff610307565b610208600480360360208110156101f557600080fd5b503569ffffffffffffffffffff16610323565b604051808669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1681526020019550505050505060405180910390f35b6101306004803603602081101561026857600080fd5b503561037c565b6101306004803603602081101561028557600080fd5b503561038e565b6102086103a0565b60005460ff1681565b60075473ffffffffffffffffffffffffffffffffffffffff1681565b60015481565b600081565b60035481565b60408051808201909152601481527f446572697665645072696365466565642e736f6c000000000000000000000000602082015290565b60025481565b60085473ffffffffffffffffffffffffffffffffffffffff1681565b60008060008060006040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260278152602001806106a26027913960400191505060405180910390fd5b60046020526000908152604090205481565b60056020526000908152604090205481565b600754600854600080549092839283928392839283926103dd9273ffffffffffffffffffffffffffffffffffffffff90811692169060ff166103f0565b9096909550429450849350600092509050565b6000808260ff16600a0a905060008573ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b15801561044457600080fd5b505afa158015610458573d6000803e3d6000fd5b505050506040513d60a081101561046e57600080fd5b50602090810151604080517f313ce567000000000000000000000000000000000000000000000000000000008152905191935060009273ffffffffffffffffffffffffffffffffffffffff8a169263313ce567926004808201939291829003018186803b1580156104de57600080fd5b505afa1580156104f2573d6000803e3d6000fd5b505050506040513d602081101561050857600080fd5b50519050610517828287610651565b915060008673ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b15801561056157600080fd5b505afa158015610575573d6000803e3d6000fd5b505050506040513d60a081101561058b57600080fd5b50602090810151604080517f313ce567000000000000000000000000000000000000000000000000000000008152905191935060009273ffffffffffffffffffffffffffffffffffffffff8b169263313ce567926004808201939291829003018186803b1580156105fb57600080fd5b505afa15801561060f573d6000803e3d6000fd5b505050506040513d602081101561062557600080fd5b50519050610634828289610651565b9150818585028161064157fe5b05955050505050505b9392505050565b60008160ff168360ff161015610672575060ff82820316600a0a830261064a565b8160ff168360ff1611156106995781830360ff16600a0a848161069157fe5b05905061064a565b50919291505056fe6e6f7420696d706c656d656e746564202d20757365206c6174657374526f756e64446174612829a164736f6c6343000706000a",
}

var DerivedPriceFeedABI = DerivedPriceFeedMetaData.ABI

var DerivedPriceFeedBin = DerivedPriceFeedMetaData.Bin

func DeployDerivedPriceFeed(auth *bind.TransactOpts, backend bind.ContractBackend, _base common.Address, _quote common.Address, _decimals uint8) (common.Address, *types.Transaction, *DerivedPriceFeed, error) {
	parsed, err := DerivedPriceFeedMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DerivedPriceFeedBin), backend, _base, _quote, _decimals)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DerivedPriceFeed{DerivedPriceFeedCaller: DerivedPriceFeedCaller{contract: contract}, DerivedPriceFeedTransactor: DerivedPriceFeedTransactor{contract: contract}, DerivedPriceFeedFilterer: DerivedPriceFeedFilterer{contract: contract}}, nil
}

type DerivedPriceFeed struct {
	address common.Address
	abi     abi.ABI
	DerivedPriceFeedCaller
	DerivedPriceFeedTransactor
	DerivedPriceFeedFilterer
}

type DerivedPriceFeedCaller struct {
	contract *bind.BoundContract
}

type DerivedPriceFeedTransactor struct {
	contract *bind.BoundContract
}

type DerivedPriceFeedFilterer struct {
	contract *bind.BoundContract
}

type DerivedPriceFeedSession struct {
	Contract     *DerivedPriceFeed
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type DerivedPriceFeedCallerSession struct {
	Contract *DerivedPriceFeedCaller
	CallOpts bind.CallOpts
}

type DerivedPriceFeedTransactorSession struct {
	Contract     *DerivedPriceFeedTransactor
	TransactOpts bind.TransactOpts
}

type DerivedPriceFeedRaw struct {
	Contract *DerivedPriceFeed
}

type DerivedPriceFeedCallerRaw struct {
	Contract *DerivedPriceFeedCaller
}

type DerivedPriceFeedTransactorRaw struct {
	Contract *DerivedPriceFeedTransactor
}

func NewDerivedPriceFeed(address common.Address, backend bind.ContractBackend) (*DerivedPriceFeed, error) {
	abi, err := abi.JSON(strings.NewReader(DerivedPriceFeedABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindDerivedPriceFeed(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DerivedPriceFeed{address: address, abi: abi, DerivedPriceFeedCaller: DerivedPriceFeedCaller{contract: contract}, DerivedPriceFeedTransactor: DerivedPriceFeedTransactor{contract: contract}, DerivedPriceFeedFilterer: DerivedPriceFeedFilterer{contract: contract}}, nil
}

func NewDerivedPriceFeedCaller(address common.Address, caller bind.ContractCaller) (*DerivedPriceFeedCaller, error) {
	contract, err := bindDerivedPriceFeed(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DerivedPriceFeedCaller{contract: contract}, nil
}

func NewDerivedPriceFeedTransactor(address common.Address, transactor bind.ContractTransactor) (*DerivedPriceFeedTransactor, error) {
	contract, err := bindDerivedPriceFeed(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DerivedPriceFeedTransactor{contract: contract}, nil
}

func NewDerivedPriceFeedFilterer(address common.Address, filterer bind.ContractFilterer) (*DerivedPriceFeedFilterer, error) {
	contract, err := bindDerivedPriceFeed(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DerivedPriceFeedFilterer{contract: contract}, nil
}

func bindDerivedPriceFeed(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DerivedPriceFeedABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_DerivedPriceFeed *DerivedPriceFeedRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DerivedPriceFeed.Contract.DerivedPriceFeedCaller.contract.Call(opts, result, method, params...)
}

func (_DerivedPriceFeed *DerivedPriceFeedRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DerivedPriceFeed.Contract.DerivedPriceFeedTransactor.contract.Transfer(opts)
}

func (_DerivedPriceFeed *DerivedPriceFeedRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DerivedPriceFeed.Contract.DerivedPriceFeedTransactor.contract.Transact(opts, method, params...)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DerivedPriceFeed.Contract.contract.Call(opts, result, method, params...)
}

func (_DerivedPriceFeed *DerivedPriceFeedTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DerivedPriceFeed.Contract.contract.Transfer(opts)
}

func (_DerivedPriceFeed *DerivedPriceFeedTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DerivedPriceFeed.Contract.contract.Transact(opts, method, params...)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) Base(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "base")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) Base() (common.Address, error) {
	return _DerivedPriceFeed.Contract.Base(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) Base() (common.Address, error) {
	return _DerivedPriceFeed.Contract.Base(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) Decimals() (uint8, error) {
	return _DerivedPriceFeed.Contract.Decimals(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) Decimals() (uint8, error) {
	return _DerivedPriceFeed.Contract.Decimals(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) Description() (string, error) {
	return _DerivedPriceFeed.Contract.Description(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) Description() (string, error) {
	return _DerivedPriceFeed.Contract.Description(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) GetAnswer(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "getAnswer", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) GetAnswer(arg0 *big.Int) (*big.Int, error) {
	return _DerivedPriceFeed.Contract.GetAnswer(&_DerivedPriceFeed.CallOpts, arg0)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) GetAnswer(arg0 *big.Int) (*big.Int, error) {
	return _DerivedPriceFeed.Contract.GetAnswer(&_DerivedPriceFeed.CallOpts, arg0)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) GetRoundData(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "getRoundData", arg0)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	out4 := *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return out0, out1, out2, out3, out4, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) GetRoundData(arg0 *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _DerivedPriceFeed.Contract.GetRoundData(&_DerivedPriceFeed.CallOpts, arg0)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) GetRoundData(arg0 *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _DerivedPriceFeed.Contract.GetRoundData(&_DerivedPriceFeed.CallOpts, arg0)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) GetTimestamp(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "getTimestamp", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) GetTimestamp(arg0 *big.Int) (*big.Int, error) {
	return _DerivedPriceFeed.Contract.GetTimestamp(&_DerivedPriceFeed.CallOpts, arg0)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) GetTimestamp(arg0 *big.Int) (*big.Int, error) {
	return _DerivedPriceFeed.Contract.GetTimestamp(&_DerivedPriceFeed.CallOpts, arg0)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) LatestAnswer() (*big.Int, error) {
	return _DerivedPriceFeed.Contract.LatestAnswer(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) LatestAnswer() (*big.Int, error) {
	return _DerivedPriceFeed.Contract.LatestAnswer(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) LatestRound() (*big.Int, error) {
	return _DerivedPriceFeed.Contract.LatestRound(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) LatestRound() (*big.Int, error) {
	return _DerivedPriceFeed.Contract.LatestRound(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "latestRoundData")

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

func (_DerivedPriceFeed *DerivedPriceFeedSession) LatestRoundData() (LatestRoundData,

	error) {
	return _DerivedPriceFeed.Contract.LatestRoundData(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _DerivedPriceFeed.Contract.LatestRoundData(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) LatestTimestamp() (*big.Int, error) {
	return _DerivedPriceFeed.Contract.LatestTimestamp(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) LatestTimestamp() (*big.Int, error) {
	return _DerivedPriceFeed.Contract.LatestTimestamp(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) Quote(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "quote")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) Quote() (common.Address, error) {
	return _DerivedPriceFeed.Contract.Quote(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) Quote() (common.Address, error) {
	return _DerivedPriceFeed.Contract.Quote(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) Version() (*big.Int, error) {
	return _DerivedPriceFeed.Contract.Version(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) Version() (*big.Int, error) {
	return _DerivedPriceFeed.Contract.Version(&_DerivedPriceFeed.CallOpts)
}

type DerivedPriceFeedAnswerUpdatedIterator struct {
	Event *DerivedPriceFeedAnswerUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DerivedPriceFeedAnswerUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DerivedPriceFeedAnswerUpdated)
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
		it.Event = new(DerivedPriceFeedAnswerUpdated)
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

func (it *DerivedPriceFeedAnswerUpdatedIterator) Error() error {
	return it.fail
}

func (it *DerivedPriceFeedAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DerivedPriceFeedAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log
}

func (_DerivedPriceFeed *DerivedPriceFeedFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*DerivedPriceFeedAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _DerivedPriceFeed.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &DerivedPriceFeedAnswerUpdatedIterator{contract: _DerivedPriceFeed.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

func (_DerivedPriceFeed *DerivedPriceFeedFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *DerivedPriceFeedAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _DerivedPriceFeed.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DerivedPriceFeedAnswerUpdated)
				if err := _DerivedPriceFeed.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

func (_DerivedPriceFeed *DerivedPriceFeedFilterer) ParseAnswerUpdated(log types.Log) (*DerivedPriceFeedAnswerUpdated, error) {
	event := new(DerivedPriceFeedAnswerUpdated)
	if err := _DerivedPriceFeed.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DerivedPriceFeedNewRoundIterator struct {
	Event *DerivedPriceFeedNewRound

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DerivedPriceFeedNewRoundIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DerivedPriceFeedNewRound)
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
		it.Event = new(DerivedPriceFeedNewRound)
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

func (it *DerivedPriceFeedNewRoundIterator) Error() error {
	return it.fail
}

func (it *DerivedPriceFeedNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DerivedPriceFeedNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log
}

func (_DerivedPriceFeed *DerivedPriceFeedFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*DerivedPriceFeedNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _DerivedPriceFeed.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &DerivedPriceFeedNewRoundIterator{contract: _DerivedPriceFeed.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

func (_DerivedPriceFeed *DerivedPriceFeedFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *DerivedPriceFeedNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _DerivedPriceFeed.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DerivedPriceFeedNewRound)
				if err := _DerivedPriceFeed.contract.UnpackLog(event, "NewRound", log); err != nil {
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

func (_DerivedPriceFeed *DerivedPriceFeedFilterer) ParseNewRound(log types.Log) (*DerivedPriceFeedNewRound, error) {
	event := new(DerivedPriceFeedNewRound)
	if err := _DerivedPriceFeed.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LatestRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

func (_DerivedPriceFeed *DerivedPriceFeed) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _DerivedPriceFeed.abi.Events["AnswerUpdated"].ID:
		return _DerivedPriceFeed.ParseAnswerUpdated(log)
	case _DerivedPriceFeed.abi.Events["NewRound"].ID:
		return _DerivedPriceFeed.ParseNewRound(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (DerivedPriceFeedAnswerUpdated) Topic() common.Hash {
	return common.HexToHash("0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f")
}

func (DerivedPriceFeedNewRound) Topic() common.Hash {
	return common.HexToHash("0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271")
}

func (_DerivedPriceFeed *DerivedPriceFeed) Address() common.Address {
	return _DerivedPriceFeed.address
}

type DerivedPriceFeedInterface interface {
	Base(opts *bind.CallOpts) (common.Address, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetAnswer(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	GetRoundData(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error)

	GetTimestamp(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	LatestAnswer(opts *bind.CallOpts) (*big.Int, error)

	LatestRound(opts *bind.CallOpts) (*big.Int, error)

	LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

		error)

	LatestTimestamp(opts *bind.CallOpts) (*big.Int, error)

	Quote(opts *bind.CallOpts) (common.Address, error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*DerivedPriceFeedAnswerUpdatedIterator, error)

	WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *DerivedPriceFeedAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error)

	ParseAnswerUpdated(log types.Log) (*DerivedPriceFeedAnswerUpdated, error)

	FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*DerivedPriceFeedNewRoundIterator, error)

	WatchNewRound(opts *bind.WatchOpts, sink chan<- *DerivedPriceFeedNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error)

	ParseNewRound(log types.Log) (*DerivedPriceFeedNewRound, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
