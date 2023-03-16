// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr2dr_client

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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

type FunctionsRequest struct {
	CodeLocation    uint8
	SecretsLocation uint8
	Language        uint8
	Source          string
	Secrets         []byte
	Args            []string
}

var OCR2DRClientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"RequestIsAlreadyPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestIsNotPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderIsNotRegistry\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enumFunctions.Location\",\"name\":\"codeLocation\",\"type\":\"uint8\"},{\"internalType\":\"enumFunctions.Location\",\"name\":\"secretsLocation\",\"type\":\"uint8\"},{\"internalType\":\"enumFunctions.CodeLanguage\",\"name\":\"language\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"secrets\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"}],\"internalType\":\"structFunctions.Request\",\"name\":\"req\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var OCR2DRClientABI = OCR2DRClientMetaData.ABI

type OCR2DRClient struct {
	address common.Address
	abi     abi.ABI
	OCR2DRClientCaller
	OCR2DRClientTransactor
	OCR2DRClientFilterer
}

type OCR2DRClientCaller struct {
	contract *bind.BoundContract
}

type OCR2DRClientTransactor struct {
	contract *bind.BoundContract
}

type OCR2DRClientFilterer struct {
	contract *bind.BoundContract
}

type OCR2DRClientSession struct {
	Contract     *OCR2DRClient
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OCR2DRClientCallerSession struct {
	Contract *OCR2DRClientCaller
	CallOpts bind.CallOpts
}

type OCR2DRClientTransactorSession struct {
	Contract     *OCR2DRClientTransactor
	TransactOpts bind.TransactOpts
}

type OCR2DRClientRaw struct {
	Contract *OCR2DRClient
}

type OCR2DRClientCallerRaw struct {
	Contract *OCR2DRClientCaller
}

type OCR2DRClientTransactorRaw struct {
	Contract *OCR2DRClientTransactor
}

func NewOCR2DRClient(address common.Address, backend bind.ContractBackend) (*OCR2DRClient, error) {
	abi, err := abi.JSON(strings.NewReader(OCR2DRClientABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOCR2DRClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClient{address: address, abi: abi, OCR2DRClientCaller: OCR2DRClientCaller{contract: contract}, OCR2DRClientTransactor: OCR2DRClientTransactor{contract: contract}, OCR2DRClientFilterer: OCR2DRClientFilterer{contract: contract}}, nil
}

func NewOCR2DRClientCaller(address common.Address, caller bind.ContractCaller) (*OCR2DRClientCaller, error) {
	contract, err := bindOCR2DRClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientCaller{contract: contract}, nil
}

func NewOCR2DRClientTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2DRClientTransactor, error) {
	contract, err := bindOCR2DRClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientTransactor{contract: contract}, nil
}

func NewOCR2DRClientFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2DRClientFilterer, error) {
	contract, err := bindOCR2DRClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientFilterer{contract: contract}, nil
}

func bindOCR2DRClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OCR2DRClientABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_OCR2DRClient *OCR2DRClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DRClient.Contract.OCR2DRClientCaller.contract.Call(opts, result, method, params...)
}

func (_OCR2DRClient *OCR2DRClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRClient.Contract.OCR2DRClientTransactor.contract.Transfer(opts)
}

func (_OCR2DRClient *OCR2DRClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DRClient.Contract.OCR2DRClientTransactor.contract.Transact(opts, method, params...)
}

func (_OCR2DRClient *OCR2DRClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DRClient.Contract.contract.Call(opts, result, method, params...)
}

func (_OCR2DRClient *OCR2DRClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRClient.Contract.contract.Transfer(opts)
}

func (_OCR2DRClient *OCR2DRClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DRClient.Contract.contract.Transact(opts, method, params...)
}

func (_OCR2DRClient *OCR2DRClientCaller) EstimateCost(opts *bind.CallOpts, req FunctionsRequest, subscriptionId uint64, gasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OCR2DRClient.contract.Call(opts, &out, "estimateCost", req, subscriptionId, gasLimit, gasPrice)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2DRClient *OCR2DRClientSession) EstimateCost(req FunctionsRequest, subscriptionId uint64, gasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	return _OCR2DRClient.Contract.EstimateCost(&_OCR2DRClient.CallOpts, req, subscriptionId, gasLimit, gasPrice)
}

func (_OCR2DRClient *OCR2DRClientCallerSession) EstimateCost(req FunctionsRequest, subscriptionId uint64, gasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	return _OCR2DRClient.Contract.EstimateCost(&_OCR2DRClient.CallOpts, req, subscriptionId, gasLimit, gasPrice)
}

func (_OCR2DRClient *OCR2DRClientCaller) GetDONPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _OCR2DRClient.contract.Call(opts, &out, "getDONPublicKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_OCR2DRClient *OCR2DRClientSession) GetDONPublicKey() ([]byte, error) {
	return _OCR2DRClient.Contract.GetDONPublicKey(&_OCR2DRClient.CallOpts)
}

func (_OCR2DRClient *OCR2DRClientCallerSession) GetDONPublicKey() ([]byte, error) {
	return _OCR2DRClient.Contract.GetDONPublicKey(&_OCR2DRClient.CallOpts)
}

func (_OCR2DRClient *OCR2DRClientTransactor) HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _OCR2DRClient.contract.Transact(opts, "handleOracleFulfillment", requestId, response, err)
}

func (_OCR2DRClient *OCR2DRClientSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _OCR2DRClient.Contract.HandleOracleFulfillment(&_OCR2DRClient.TransactOpts, requestId, response, err)
}

func (_OCR2DRClient *OCR2DRClientTransactorSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _OCR2DRClient.Contract.HandleOracleFulfillment(&_OCR2DRClient.TransactOpts, requestId, response, err)
}

type OCR2DRClientRequestFulfilledIterator struct {
	Event *OCR2DRClientRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRClientRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRClientRequestFulfilled)
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
		it.Event = new(OCR2DRClientRequestFulfilled)
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

func (it *OCR2DRClientRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *OCR2DRClientRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRClientRequestFulfilled struct {
	Id  [32]byte
	Raw types.Log
}

func (_OCR2DRClient *OCR2DRClientFilterer) FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*OCR2DRClientRequestFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _OCR2DRClient.contract.FilterLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientRequestFulfilledIterator{contract: _OCR2DRClient.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

func (_OCR2DRClient *OCR2DRClientFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *OCR2DRClientRequestFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _OCR2DRClient.contract.WatchLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRClientRequestFulfilled)
				if err := _OCR2DRClient.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
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

func (_OCR2DRClient *OCR2DRClientFilterer) ParseRequestFulfilled(log types.Log) (*OCR2DRClientRequestFulfilled, error) {
	event := new(OCR2DRClientRequestFulfilled)
	if err := _OCR2DRClient.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRClientRequestSentIterator struct {
	Event *OCR2DRClientRequestSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRClientRequestSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRClientRequestSent)
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
		it.Event = new(OCR2DRClientRequestSent)
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

func (it *OCR2DRClientRequestSentIterator) Error() error {
	return it.fail
}

func (it *OCR2DRClientRequestSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRClientRequestSent struct {
	Id  [32]byte
	Raw types.Log
}

func (_OCR2DRClient *OCR2DRClientFilterer) FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*OCR2DRClientRequestSentIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _OCR2DRClient.contract.FilterLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientRequestSentIterator{contract: _OCR2DRClient.contract, event: "RequestSent", logs: logs, sub: sub}, nil
}

func (_OCR2DRClient *OCR2DRClientFilterer) WatchRequestSent(opts *bind.WatchOpts, sink chan<- *OCR2DRClientRequestSent, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _OCR2DRClient.contract.WatchLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRClientRequestSent)
				if err := _OCR2DRClient.contract.UnpackLog(event, "RequestSent", log); err != nil {
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

func (_OCR2DRClient *OCR2DRClientFilterer) ParseRequestSent(log types.Log) (*OCR2DRClientRequestSent, error) {
	event := new(OCR2DRClientRequestSent)
	if err := _OCR2DRClient.contract.UnpackLog(event, "RequestSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OCR2DRClient *OCR2DRClient) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OCR2DRClient.abi.Events["RequestFulfilled"].ID:
		return _OCR2DRClient.ParseRequestFulfilled(log)
	case _OCR2DRClient.abi.Events["RequestSent"].ID:
		return _OCR2DRClient.ParseRequestSent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OCR2DRClientRequestFulfilled) Topic() common.Hash {
	return common.HexToHash("0x85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e6")
}

func (OCR2DRClientRequestSent) Topic() common.Hash {
	return common.HexToHash("0x1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db8")
}

func (_OCR2DRClient *OCR2DRClient) Address() common.Address {
	return _OCR2DRClient.address
}

type OCR2DRClientInterface interface {
	EstimateCost(opts *bind.CallOpts, req FunctionsRequest, subscriptionId uint64, gasLimit uint32, gasPrice *big.Int) (*big.Int, error)

	GetDONPublicKey(opts *bind.CallOpts) ([]byte, error)

	HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error)

	FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*OCR2DRClientRequestFulfilledIterator, error)

	WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *OCR2DRClientRequestFulfilled, id [][32]byte) (event.Subscription, error)

	ParseRequestFulfilled(log types.Log) (*OCR2DRClientRequestFulfilled, error)

	FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*OCR2DRClientRequestSentIterator, error)

	WatchRequestSent(opts *bind.WatchOpts, sink chan<- *OCR2DRClientRequestSent, id [][32]byte) (event.Subscription, error)

	ParseRequestSent(log types.Log) (*OCR2DRClientRequestSent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
