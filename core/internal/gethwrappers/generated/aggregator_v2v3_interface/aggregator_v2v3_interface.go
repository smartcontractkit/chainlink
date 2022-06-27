// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package aggregator_v2v3_interface

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

var IAggregatorV2V3MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var IAggregatorV2V3ABI = IAggregatorV2V3MetaData.ABI

type IAggregatorV2V3 struct {
	address common.Address
	abi     abi.ABI
	IAggregatorV2V3Caller
	IAggregatorV2V3Transactor
	IAggregatorV2V3Filterer
}

type IAggregatorV2V3Caller struct {
	contract *bind.BoundContract
}

type IAggregatorV2V3Transactor struct {
	contract *bind.BoundContract
}

type IAggregatorV2V3Filterer struct {
	contract *bind.BoundContract
}

type IAggregatorV2V3Session struct {
	Contract     *IAggregatorV2V3
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type IAggregatorV2V3CallerSession struct {
	Contract *IAggregatorV2V3Caller
	CallOpts bind.CallOpts
}

type IAggregatorV2V3TransactorSession struct {
	Contract     *IAggregatorV2V3Transactor
	TransactOpts bind.TransactOpts
}

type IAggregatorV2V3Raw struct {
	Contract *IAggregatorV2V3
}

type IAggregatorV2V3CallerRaw struct {
	Contract *IAggregatorV2V3Caller
}

type IAggregatorV2V3TransactorRaw struct {
	Contract *IAggregatorV2V3Transactor
}

func NewIAggregatorV2V3(address common.Address, backend bind.ContractBackend) (*IAggregatorV2V3, error) {
	abi, err := abi.JSON(strings.NewReader(IAggregatorV2V3ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindIAggregatorV2V3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAggregatorV2V3{address: address, abi: abi, IAggregatorV2V3Caller: IAggregatorV2V3Caller{contract: contract}, IAggregatorV2V3Transactor: IAggregatorV2V3Transactor{contract: contract}, IAggregatorV2V3Filterer: IAggregatorV2V3Filterer{contract: contract}}, nil
}

func NewIAggregatorV2V3Caller(address common.Address, caller bind.ContractCaller) (*IAggregatorV2V3Caller, error) {
	contract, err := bindIAggregatorV2V3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAggregatorV2V3Caller{contract: contract}, nil
}

func NewIAggregatorV2V3Transactor(address common.Address, transactor bind.ContractTransactor) (*IAggregatorV2V3Transactor, error) {
	contract, err := bindIAggregatorV2V3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAggregatorV2V3Transactor{contract: contract}, nil
}

func NewIAggregatorV2V3Filterer(address common.Address, filterer bind.ContractFilterer) (*IAggregatorV2V3Filterer, error) {
	contract, err := bindIAggregatorV2V3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAggregatorV2V3Filterer{contract: contract}, nil
}

func bindIAggregatorV2V3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IAggregatorV2V3ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_IAggregatorV2V3 *IAggregatorV2V3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAggregatorV2V3.Contract.IAggregatorV2V3Caller.contract.Call(opts, result, method, params...)
}

func (_IAggregatorV2V3 *IAggregatorV2V3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAggregatorV2V3.Contract.IAggregatorV2V3Transactor.contract.Transfer(opts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAggregatorV2V3.Contract.IAggregatorV2V3Transactor.contract.Transact(opts, method, params...)
}

func (_IAggregatorV2V3 *IAggregatorV2V3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAggregatorV2V3.Contract.contract.Call(opts, result, method, params...)
}

func (_IAggregatorV2V3 *IAggregatorV2V3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAggregatorV2V3.Contract.contract.Transfer(opts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAggregatorV2V3.Contract.contract.Transact(opts, method, params...)
}

func (_IAggregatorV2V3 *IAggregatorV2V3Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IAggregatorV2V3.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IAggregatorV2V3 *IAggregatorV2V3Session) Decimals() (uint8, error) {
	return _IAggregatorV2V3.Contract.Decimals(&_IAggregatorV2V3.CallOpts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3CallerSession) Decimals() (uint8, error) {
	return _IAggregatorV2V3.Contract.Decimals(&_IAggregatorV2V3.CallOpts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3Caller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IAggregatorV2V3.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_IAggregatorV2V3 *IAggregatorV2V3Session) Description() (string, error) {
	return _IAggregatorV2V3.Contract.Description(&_IAggregatorV2V3.CallOpts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3CallerSession) Description() (string, error) {
	return _IAggregatorV2V3.Contract.Description(&_IAggregatorV2V3.CallOpts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3Caller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IAggregatorV2V3.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAggregatorV2V3 *IAggregatorV2V3Session) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _IAggregatorV2V3.Contract.GetAnswer(&_IAggregatorV2V3.CallOpts, roundId)
}

func (_IAggregatorV2V3 *IAggregatorV2V3CallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _IAggregatorV2V3.Contract.GetAnswer(&_IAggregatorV2V3.CallOpts, roundId)
}

func (_IAggregatorV2V3 *IAggregatorV2V3Caller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _IAggregatorV2V3.contract.Call(opts, &out, "getRoundData", _roundId)

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

func (_IAggregatorV2V3 *IAggregatorV2V3Session) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _IAggregatorV2V3.Contract.GetRoundData(&_IAggregatorV2V3.CallOpts, _roundId)
}

func (_IAggregatorV2V3 *IAggregatorV2V3CallerSession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _IAggregatorV2V3.Contract.GetRoundData(&_IAggregatorV2V3.CallOpts, _roundId)
}

func (_IAggregatorV2V3 *IAggregatorV2V3Caller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IAggregatorV2V3.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAggregatorV2V3 *IAggregatorV2V3Session) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _IAggregatorV2V3.Contract.GetTimestamp(&_IAggregatorV2V3.CallOpts, roundId)
}

func (_IAggregatorV2V3 *IAggregatorV2V3CallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _IAggregatorV2V3.Contract.GetTimestamp(&_IAggregatorV2V3.CallOpts, roundId)
}

func (_IAggregatorV2V3 *IAggregatorV2V3Caller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAggregatorV2V3.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAggregatorV2V3 *IAggregatorV2V3Session) LatestAnswer() (*big.Int, error) {
	return _IAggregatorV2V3.Contract.LatestAnswer(&_IAggregatorV2V3.CallOpts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3CallerSession) LatestAnswer() (*big.Int, error) {
	return _IAggregatorV2V3.Contract.LatestAnswer(&_IAggregatorV2V3.CallOpts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3Caller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAggregatorV2V3.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAggregatorV2V3 *IAggregatorV2V3Session) LatestRound() (*big.Int, error) {
	return _IAggregatorV2V3.Contract.LatestRound(&_IAggregatorV2V3.CallOpts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3CallerSession) LatestRound() (*big.Int, error) {
	return _IAggregatorV2V3.Contract.LatestRound(&_IAggregatorV2V3.CallOpts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3Caller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _IAggregatorV2V3.contract.Call(opts, &out, "latestRoundData")

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

func (_IAggregatorV2V3 *IAggregatorV2V3Session) LatestRoundData() (LatestRoundData,

	error) {
	return _IAggregatorV2V3.Contract.LatestRoundData(&_IAggregatorV2V3.CallOpts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3CallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _IAggregatorV2V3.Contract.LatestRoundData(&_IAggregatorV2V3.CallOpts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3Caller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAggregatorV2V3.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAggregatorV2V3 *IAggregatorV2V3Session) LatestTimestamp() (*big.Int, error) {
	return _IAggregatorV2V3.Contract.LatestTimestamp(&_IAggregatorV2V3.CallOpts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3CallerSession) LatestTimestamp() (*big.Int, error) {
	return _IAggregatorV2V3.Contract.LatestTimestamp(&_IAggregatorV2V3.CallOpts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3Caller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAggregatorV2V3.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAggregatorV2V3 *IAggregatorV2V3Session) Version() (*big.Int, error) {
	return _IAggregatorV2V3.Contract.Version(&_IAggregatorV2V3.CallOpts)
}

func (_IAggregatorV2V3 *IAggregatorV2V3CallerSession) Version() (*big.Int, error) {
	return _IAggregatorV2V3.Contract.Version(&_IAggregatorV2V3.CallOpts)
}

type IAggregatorV2V3AnswerUpdatedIterator struct {
	Event *IAggregatorV2V3AnswerUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAggregatorV2V3AnswerUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAggregatorV2V3AnswerUpdated)
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
		it.Event = new(IAggregatorV2V3AnswerUpdated)
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

func (it *IAggregatorV2V3AnswerUpdatedIterator) Error() error {
	return it.fail
}

func (it *IAggregatorV2V3AnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAggregatorV2V3AnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log
}

func (_IAggregatorV2V3 *IAggregatorV2V3Filterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*IAggregatorV2V3AnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _IAggregatorV2V3.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &IAggregatorV2V3AnswerUpdatedIterator{contract: _IAggregatorV2V3.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

func (_IAggregatorV2V3 *IAggregatorV2V3Filterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *IAggregatorV2V3AnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _IAggregatorV2V3.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAggregatorV2V3AnswerUpdated)
				if err := _IAggregatorV2V3.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

func (_IAggregatorV2V3 *IAggregatorV2V3Filterer) ParseAnswerUpdated(log types.Log) (*IAggregatorV2V3AnswerUpdated, error) {
	event := new(IAggregatorV2V3AnswerUpdated)
	if err := _IAggregatorV2V3.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAggregatorV2V3NewRoundIterator struct {
	Event *IAggregatorV2V3NewRound

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAggregatorV2V3NewRoundIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAggregatorV2V3NewRound)
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
		it.Event = new(IAggregatorV2V3NewRound)
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

func (it *IAggregatorV2V3NewRoundIterator) Error() error {
	return it.fail
}

func (it *IAggregatorV2V3NewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAggregatorV2V3NewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log
}

func (_IAggregatorV2V3 *IAggregatorV2V3Filterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*IAggregatorV2V3NewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _IAggregatorV2V3.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &IAggregatorV2V3NewRoundIterator{contract: _IAggregatorV2V3.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

func (_IAggregatorV2V3 *IAggregatorV2V3Filterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *IAggregatorV2V3NewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _IAggregatorV2V3.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAggregatorV2V3NewRound)
				if err := _IAggregatorV2V3.contract.UnpackLog(event, "NewRound", log); err != nil {
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

func (_IAggregatorV2V3 *IAggregatorV2V3Filterer) ParseNewRound(log types.Log) (*IAggregatorV2V3NewRound, error) {
	event := new(IAggregatorV2V3NewRound)
	if err := _IAggregatorV2V3.contract.UnpackLog(event, "NewRound", log); err != nil {
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

func (_IAggregatorV2V3 *IAggregatorV2V3) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _IAggregatorV2V3.abi.Events["AnswerUpdated"].ID:
		return _IAggregatorV2V3.ParseAnswerUpdated(log)
	case _IAggregatorV2V3.abi.Events["NewRound"].ID:
		return _IAggregatorV2V3.ParseNewRound(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (IAggregatorV2V3AnswerUpdated) Topic() common.Hash {
	return common.HexToHash("0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f")
}

func (IAggregatorV2V3NewRound) Topic() common.Hash {
	return common.HexToHash("0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271")
}

func (_IAggregatorV2V3 *IAggregatorV2V3) Address() common.Address {
	return _IAggregatorV2V3.address
}

type IAggregatorV2V3Interface interface {
	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error)

	GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

		error)

	GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error)

	LatestAnswer(opts *bind.CallOpts) (*big.Int, error)

	LatestRound(opts *bind.CallOpts) (*big.Int, error)

	LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

		error)

	LatestTimestamp(opts *bind.CallOpts) (*big.Int, error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*IAggregatorV2V3AnswerUpdatedIterator, error)

	WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *IAggregatorV2V3AnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error)

	ParseAnswerUpdated(log types.Log) (*IAggregatorV2V3AnswerUpdated, error)

	FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*IAggregatorV2V3NewRoundIterator, error)

	WatchNewRound(opts *bind.WatchOpts, sink chan<- *IAggregatorV2V3NewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error)

	ParseNewRound(log types.Log) (*IAggregatorV2V3NewRound, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
