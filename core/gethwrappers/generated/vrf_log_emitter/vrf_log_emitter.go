// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_log_emitter

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

var VRFLogEmitterMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"emitRandomWordsFulfilled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"emitRandomWordsRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061027f806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063ca920adb1461003b578063fe62d3e914610050575b600080fd5b61004e61004936600461015b565b610063565b005b61004e61005e366004610212565b6100eb565b604080518881526020810188905261ffff86168183015263ffffffff858116606083015284166080820152905173ffffffffffffffffffffffffffffffffffffffff83169167ffffffffffffffff8816918b917f63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772919081900360a00190a45050505050505050565b604080518481526bffffffffffffffffffffffff8416602082015282151581830152905185917f7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4919081900360600190a250505050565b803563ffffffff8116811461015657600080fd5b919050565b600080600080600080600080610100898b03121561017857600080fd5b883597506020890135965060408901359550606089013567ffffffffffffffff811681146101a557600080fd5b9450608089013561ffff811681146101bc57600080fd5b93506101ca60a08a01610142565b92506101d860c08a01610142565b915060e089013573ffffffffffffffffffffffffffffffffffffffff8116811461020157600080fd5b809150509295985092959890939650565b6000806000806080858703121561022857600080fd5b843593506020850135925060408501356bffffffffffffffffffffffff8116811461025257600080fd5b91506060850135801515811461026757600080fd5b93969295509093505056fea164736f6c6343000813000a",
}

var VRFLogEmitterABI = VRFLogEmitterMetaData.ABI

var VRFLogEmitterBin = VRFLogEmitterMetaData.Bin

func DeployVRFLogEmitter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFLogEmitter, error) {
	parsed, err := VRFLogEmitterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFLogEmitterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFLogEmitter{address: address, abi: *parsed, VRFLogEmitterCaller: VRFLogEmitterCaller{contract: contract}, VRFLogEmitterTransactor: VRFLogEmitterTransactor{contract: contract}, VRFLogEmitterFilterer: VRFLogEmitterFilterer{contract: contract}}, nil
}

type VRFLogEmitter struct {
	address common.Address
	abi     abi.ABI
	VRFLogEmitterCaller
	VRFLogEmitterTransactor
	VRFLogEmitterFilterer
}

type VRFLogEmitterCaller struct {
	contract *bind.BoundContract
}

type VRFLogEmitterTransactor struct {
	contract *bind.BoundContract
}

type VRFLogEmitterFilterer struct {
	contract *bind.BoundContract
}

type VRFLogEmitterSession struct {
	Contract     *VRFLogEmitter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFLogEmitterCallerSession struct {
	Contract *VRFLogEmitterCaller
	CallOpts bind.CallOpts
}

type VRFLogEmitterTransactorSession struct {
	Contract     *VRFLogEmitterTransactor
	TransactOpts bind.TransactOpts
}

type VRFLogEmitterRaw struct {
	Contract *VRFLogEmitter
}

type VRFLogEmitterCallerRaw struct {
	Contract *VRFLogEmitterCaller
}

type VRFLogEmitterTransactorRaw struct {
	Contract *VRFLogEmitterTransactor
}

func NewVRFLogEmitter(address common.Address, backend bind.ContractBackend) (*VRFLogEmitter, error) {
	abi, err := abi.JSON(strings.NewReader(VRFLogEmitterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFLogEmitter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFLogEmitter{address: address, abi: abi, VRFLogEmitterCaller: VRFLogEmitterCaller{contract: contract}, VRFLogEmitterTransactor: VRFLogEmitterTransactor{contract: contract}, VRFLogEmitterFilterer: VRFLogEmitterFilterer{contract: contract}}, nil
}

func NewVRFLogEmitterCaller(address common.Address, caller bind.ContractCaller) (*VRFLogEmitterCaller, error) {
	contract, err := bindVRFLogEmitter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFLogEmitterCaller{contract: contract}, nil
}

func NewVRFLogEmitterTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFLogEmitterTransactor, error) {
	contract, err := bindVRFLogEmitter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFLogEmitterTransactor{contract: contract}, nil
}

func NewVRFLogEmitterFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFLogEmitterFilterer, error) {
	contract, err := bindVRFLogEmitter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFLogEmitterFilterer{contract: contract}, nil
}

func bindVRFLogEmitter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFLogEmitterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFLogEmitter *VRFLogEmitterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFLogEmitter.Contract.VRFLogEmitterCaller.contract.Call(opts, result, method, params...)
}

func (_VRFLogEmitter *VRFLogEmitterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFLogEmitter.Contract.VRFLogEmitterTransactor.contract.Transfer(opts)
}

func (_VRFLogEmitter *VRFLogEmitterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFLogEmitter.Contract.VRFLogEmitterTransactor.contract.Transact(opts, method, params...)
}

func (_VRFLogEmitter *VRFLogEmitterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFLogEmitter.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFLogEmitter *VRFLogEmitterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFLogEmitter.Contract.contract.Transfer(opts)
}

func (_VRFLogEmitter *VRFLogEmitterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFLogEmitter.Contract.contract.Transact(opts, method, params...)
}

func (_VRFLogEmitter *VRFLogEmitterTransactor) EmitRandomWordsFulfilled(opts *bind.TransactOpts, requestId *big.Int, outputSeed *big.Int, payment *big.Int, success bool) (*types.Transaction, error) {
	return _VRFLogEmitter.contract.Transact(opts, "emitRandomWordsFulfilled", requestId, outputSeed, payment, success)
}

func (_VRFLogEmitter *VRFLogEmitterSession) EmitRandomWordsFulfilled(requestId *big.Int, outputSeed *big.Int, payment *big.Int, success bool) (*types.Transaction, error) {
	return _VRFLogEmitter.Contract.EmitRandomWordsFulfilled(&_VRFLogEmitter.TransactOpts, requestId, outputSeed, payment, success)
}

func (_VRFLogEmitter *VRFLogEmitterTransactorSession) EmitRandomWordsFulfilled(requestId *big.Int, outputSeed *big.Int, payment *big.Int, success bool) (*types.Transaction, error) {
	return _VRFLogEmitter.Contract.EmitRandomWordsFulfilled(&_VRFLogEmitter.TransactOpts, requestId, outputSeed, payment, success)
}

func (_VRFLogEmitter *VRFLogEmitterTransactor) EmitRandomWordsRequested(opts *bind.TransactOpts, keyHash [32]byte, requestId *big.Int, preSeed *big.Int, subId uint64, minimumRequestConfirmations uint16, callbackGasLimit uint32, numWords uint32, sender common.Address) (*types.Transaction, error) {
	return _VRFLogEmitter.contract.Transact(opts, "emitRandomWordsRequested", keyHash, requestId, preSeed, subId, minimumRequestConfirmations, callbackGasLimit, numWords, sender)
}

func (_VRFLogEmitter *VRFLogEmitterSession) EmitRandomWordsRequested(keyHash [32]byte, requestId *big.Int, preSeed *big.Int, subId uint64, minimumRequestConfirmations uint16, callbackGasLimit uint32, numWords uint32, sender common.Address) (*types.Transaction, error) {
	return _VRFLogEmitter.Contract.EmitRandomWordsRequested(&_VRFLogEmitter.TransactOpts, keyHash, requestId, preSeed, subId, minimumRequestConfirmations, callbackGasLimit, numWords, sender)
}

func (_VRFLogEmitter *VRFLogEmitterTransactorSession) EmitRandomWordsRequested(keyHash [32]byte, requestId *big.Int, preSeed *big.Int, subId uint64, minimumRequestConfirmations uint16, callbackGasLimit uint32, numWords uint32, sender common.Address) (*types.Transaction, error) {
	return _VRFLogEmitter.Contract.EmitRandomWordsRequested(&_VRFLogEmitter.TransactOpts, keyHash, requestId, preSeed, subId, minimumRequestConfirmations, callbackGasLimit, numWords, sender)
}

type VRFLogEmitterRandomWordsFulfilledIterator struct {
	Event *VRFLogEmitterRandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFLogEmitterRandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFLogEmitterRandomWordsFulfilled)
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
		it.Event = new(VRFLogEmitterRandomWordsFulfilled)
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

func (it *VRFLogEmitterRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFLogEmitterRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFLogEmitterRandomWordsFulfilled struct {
	RequestId  *big.Int
	OutputSeed *big.Int
	Payment    *big.Int
	Success    bool
	Raw        types.Log
}

func (_VRFLogEmitter *VRFLogEmitterFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*VRFLogEmitterRandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFLogEmitter.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFLogEmitterRandomWordsFulfilledIterator{contract: _VRFLogEmitter.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFLogEmitter *VRFLogEmitterFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFLogEmitterRandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFLogEmitter.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFLogEmitterRandomWordsFulfilled)
				if err := _VRFLogEmitter.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_VRFLogEmitter *VRFLogEmitterFilterer) ParseRandomWordsFulfilled(log types.Log) (*VRFLogEmitterRandomWordsFulfilled, error) {
	event := new(VRFLogEmitterRandomWordsFulfilled)
	if err := _VRFLogEmitter.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFLogEmitterRandomWordsRequestedIterator struct {
	Event *VRFLogEmitterRandomWordsRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFLogEmitterRandomWordsRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFLogEmitterRandomWordsRequested)
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
		it.Event = new(VRFLogEmitterRandomWordsRequested)
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

func (it *VRFLogEmitterRandomWordsRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFLogEmitterRandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFLogEmitterRandomWordsRequested struct {
	KeyHash                     [32]byte
	RequestId                   *big.Int
	PreSeed                     *big.Int
	SubId                       uint64
	MinimumRequestConfirmations uint16
	CallbackGasLimit            uint32
	NumWords                    uint32
	Sender                      common.Address
	Raw                         types.Log
}

func (_VRFLogEmitter *VRFLogEmitterFilterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*VRFLogEmitterRandomWordsRequestedIterator, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFLogEmitter.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFLogEmitterRandomWordsRequestedIterator{contract: _VRFLogEmitter.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

func (_VRFLogEmitter *VRFLogEmitterFilterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFLogEmitterRandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFLogEmitter.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFLogEmitterRandomWordsRequested)
				if err := _VRFLogEmitter.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

func (_VRFLogEmitter *VRFLogEmitterFilterer) ParseRandomWordsRequested(log types.Log) (*VRFLogEmitterRandomWordsRequested, error) {
	event := new(VRFLogEmitterRandomWordsRequested)
	if err := _VRFLogEmitter.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFLogEmitter *VRFLogEmitter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFLogEmitter.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFLogEmitter.ParseRandomWordsFulfilled(log)
	case _VRFLogEmitter.abi.Events["RandomWordsRequested"].ID:
		return _VRFLogEmitter.ParseRandomWordsRequested(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFLogEmitterRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4")
}

func (VRFLogEmitterRandomWordsRequested) Topic() common.Hash {
	return common.HexToHash("0x63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772")
}

func (_VRFLogEmitter *VRFLogEmitter) Address() common.Address {
	return _VRFLogEmitter.address
}

type VRFLogEmitterInterface interface {
	EmitRandomWordsFulfilled(opts *bind.TransactOpts, requestId *big.Int, outputSeed *big.Int, payment *big.Int, success bool) (*types.Transaction, error)

	EmitRandomWordsRequested(opts *bind.TransactOpts, keyHash [32]byte, requestId *big.Int, preSeed *big.Int, subId uint64, minimumRequestConfirmations uint16, callbackGasLimit uint32, numWords uint32, sender common.Address) (*types.Transaction, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*VRFLogEmitterRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFLogEmitterRandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFLogEmitterRandomWordsFulfilled, error)

	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*VRFLogEmitterRandomWordsRequestedIterator, error)

	WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFLogEmitterRandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsRequested(log types.Log) (*VRFLogEmitterRandomWordsRequested, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
