// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_client

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

var FunctionsClientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"OnlyRouterCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var FunctionsClientABI = FunctionsClientMetaData.ABI

type FunctionsClient struct {
	address common.Address
	abi     abi.ABI
	FunctionsClientCaller
	FunctionsClientTransactor
	FunctionsClientFilterer
}

type FunctionsClientCaller struct {
	contract *bind.BoundContract
}

type FunctionsClientTransactor struct {
	contract *bind.BoundContract
}

type FunctionsClientFilterer struct {
	contract *bind.BoundContract
}

type FunctionsClientSession struct {
	Contract     *FunctionsClient
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsClientCallerSession struct {
	Contract *FunctionsClientCaller
	CallOpts bind.CallOpts
}

type FunctionsClientTransactorSession struct {
	Contract     *FunctionsClientTransactor
	TransactOpts bind.TransactOpts
}

type FunctionsClientRaw struct {
	Contract *FunctionsClient
}

type FunctionsClientCallerRaw struct {
	Contract *FunctionsClientCaller
}

type FunctionsClientTransactorRaw struct {
	Contract *FunctionsClientTransactor
}

func NewFunctionsClient(address common.Address, backend bind.ContractBackend) (*FunctionsClient, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsClientABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctionsClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionsClient{address: address, abi: abi, FunctionsClientCaller: FunctionsClientCaller{contract: contract}, FunctionsClientTransactor: FunctionsClientTransactor{contract: contract}, FunctionsClientFilterer: FunctionsClientFilterer{contract: contract}}, nil
}

func NewFunctionsClientCaller(address common.Address, caller bind.ContractCaller) (*FunctionsClientCaller, error) {
	contract, err := bindFunctionsClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientCaller{contract: contract}, nil
}

func NewFunctionsClientTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsClientTransactor, error) {
	contract, err := bindFunctionsClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientTransactor{contract: contract}, nil
}

func NewFunctionsClientFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsClientFilterer, error) {
	contract, err := bindFunctionsClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientFilterer{contract: contract}, nil
}

func bindFunctionsClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionsClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FunctionsClient *FunctionsClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsClient.Contract.FunctionsClientCaller.contract.Call(opts, result, method, params...)
}

func (_FunctionsClient *FunctionsClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsClient.Contract.FunctionsClientTransactor.contract.Transfer(opts)
}

func (_FunctionsClient *FunctionsClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsClient.Contract.FunctionsClientTransactor.contract.Transact(opts, method, params...)
}

func (_FunctionsClient *FunctionsClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsClient.Contract.contract.Call(opts, result, method, params...)
}

func (_FunctionsClient *FunctionsClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsClient.Contract.contract.Transfer(opts)
}

func (_FunctionsClient *FunctionsClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsClient.Contract.contract.Transact(opts, method, params...)
}

func (_FunctionsClient *FunctionsClientTransactor) HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsClient.contract.Transact(opts, "handleOracleFulfillment", requestId, response, err)
}

func (_FunctionsClient *FunctionsClientSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsClient.Contract.HandleOracleFulfillment(&_FunctionsClient.TransactOpts, requestId, response, err)
}

func (_FunctionsClient *FunctionsClientTransactorSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsClient.Contract.HandleOracleFulfillment(&_FunctionsClient.TransactOpts, requestId, response, err)
}

type FunctionsClientRequestFulfilledIterator struct {
	Event *FunctionsClientRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientRequestFulfilled)
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
		it.Event = new(FunctionsClientRequestFulfilled)
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

func (it *FunctionsClientRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientRequestFulfilled struct {
	Id  [32]byte
	Raw types.Log
}

func (_FunctionsClient *FunctionsClientFilterer) FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientRequestFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClient.contract.FilterLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientRequestFulfilledIterator{contract: _FunctionsClient.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

func (_FunctionsClient *FunctionsClientFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *FunctionsClientRequestFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClient.contract.WatchLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientRequestFulfilled)
				if err := _FunctionsClient.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
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

func (_FunctionsClient *FunctionsClientFilterer) ParseRequestFulfilled(log types.Log) (*FunctionsClientRequestFulfilled, error) {
	event := new(FunctionsClientRequestFulfilled)
	if err := _FunctionsClient.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsClientRequestSentIterator struct {
	Event *FunctionsClientRequestSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientRequestSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientRequestSent)
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
		it.Event = new(FunctionsClientRequestSent)
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

func (it *FunctionsClientRequestSentIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientRequestSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientRequestSent struct {
	Id  [32]byte
	Raw types.Log
}

func (_FunctionsClient *FunctionsClientFilterer) FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientRequestSentIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClient.contract.FilterLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientRequestSentIterator{contract: _FunctionsClient.contract, event: "RequestSent", logs: logs, sub: sub}, nil
}

func (_FunctionsClient *FunctionsClientFilterer) WatchRequestSent(opts *bind.WatchOpts, sink chan<- *FunctionsClientRequestSent, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClient.contract.WatchLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientRequestSent)
				if err := _FunctionsClient.contract.UnpackLog(event, "RequestSent", log); err != nil {
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

func (_FunctionsClient *FunctionsClientFilterer) ParseRequestSent(log types.Log) (*FunctionsClientRequestSent, error) {
	event := new(FunctionsClientRequestSent)
	if err := _FunctionsClient.contract.UnpackLog(event, "RequestSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_FunctionsClient *FunctionsClient) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FunctionsClient.abi.Events["RequestFulfilled"].ID:
		return _FunctionsClient.ParseRequestFulfilled(log)
	case _FunctionsClient.abi.Events["RequestSent"].ID:
		return _FunctionsClient.ParseRequestSent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsClientRequestFulfilled) Topic() common.Hash {
	return common.HexToHash("0x85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e6")
}

func (FunctionsClientRequestSent) Topic() common.Hash {
	return common.HexToHash("0x1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db8")
}

func (_FunctionsClient *FunctionsClient) Address() common.Address {
	return _FunctionsClient.address
}

type FunctionsClientInterface interface {
	HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error)

	FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientRequestFulfilledIterator, error)

	WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *FunctionsClientRequestFulfilled, id [][32]byte) (event.Subscription, error)

	ParseRequestFulfilled(log types.Log) (*FunctionsClientRequestFulfilled, error)

	FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientRequestSentIterator, error)

	WatchRequestSent(opts *bind.WatchOpts, sink chan<- *FunctionsClientRequestSent, id [][32]byte) (event.Subscription, error)

	ParseRequestSent(log types.Log) (*FunctionsClientRequestSent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
