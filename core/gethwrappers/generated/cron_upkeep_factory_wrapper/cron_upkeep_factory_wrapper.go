// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package cron_upkeep_factory_wrapper

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

var CronUpkeepFactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"upkeep\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"NewCronUpkeepCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cronDelegateAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"handler\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"cronString\",\"type\":\"string\"}],\"name\":\"encodeCronJob\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"cronString\",\"type\":\"string\"}],\"name\":\"encodeCronString\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newCronUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedJob\",\"type\":\"bytes\"}],\"name\":\"newCronUpkeepWithJob\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_maxJobs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxJobs\",\"type\":\"uint256\"}],\"name\":\"setMaxJobs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var CronUpkeepFactoryABI = CronUpkeepFactoryMetaData.ABI

type CronUpkeepFactory struct {
	address common.Address
	abi     abi.ABI
	CronUpkeepFactoryCaller
	CronUpkeepFactoryTransactor
	CronUpkeepFactoryFilterer
}

type CronUpkeepFactoryCaller struct {
	contract *bind.BoundContract
}

type CronUpkeepFactoryTransactor struct {
	contract *bind.BoundContract
}

type CronUpkeepFactoryFilterer struct {
	contract *bind.BoundContract
}

type CronUpkeepFactorySession struct {
	Contract     *CronUpkeepFactory
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CronUpkeepFactoryCallerSession struct {
	Contract *CronUpkeepFactoryCaller
	CallOpts bind.CallOpts
}

type CronUpkeepFactoryTransactorSession struct {
	Contract     *CronUpkeepFactoryTransactor
	TransactOpts bind.TransactOpts
}

type CronUpkeepFactoryRaw struct {
	Contract *CronUpkeepFactory
}

type CronUpkeepFactoryCallerRaw struct {
	Contract *CronUpkeepFactoryCaller
}

type CronUpkeepFactoryTransactorRaw struct {
	Contract *CronUpkeepFactoryTransactor
}

func NewCronUpkeepFactory(address common.Address, backend bind.ContractBackend) (*CronUpkeepFactory, error) {
	abi, err := abi.JSON(strings.NewReader(CronUpkeepFactoryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCronUpkeepFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepFactory{address: address, abi: abi, CronUpkeepFactoryCaller: CronUpkeepFactoryCaller{contract: contract}, CronUpkeepFactoryTransactor: CronUpkeepFactoryTransactor{contract: contract}, CronUpkeepFactoryFilterer: CronUpkeepFactoryFilterer{contract: contract}}, nil
}

func NewCronUpkeepFactoryCaller(address common.Address, caller bind.ContractCaller) (*CronUpkeepFactoryCaller, error) {
	contract, err := bindCronUpkeepFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepFactoryCaller{contract: contract}, nil
}

func NewCronUpkeepFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*CronUpkeepFactoryTransactor, error) {
	contract, err := bindCronUpkeepFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepFactoryTransactor{contract: contract}, nil
}

func NewCronUpkeepFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*CronUpkeepFactoryFilterer, error) {
	contract, err := bindCronUpkeepFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepFactoryFilterer{contract: contract}, nil
}

func bindCronUpkeepFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CronUpkeepFactoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_CronUpkeepFactory *CronUpkeepFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CronUpkeepFactory.Contract.CronUpkeepFactoryCaller.contract.Call(opts, result, method, params...)
}

func (_CronUpkeepFactory *CronUpkeepFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.CronUpkeepFactoryTransactor.contract.Transfer(opts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.CronUpkeepFactoryTransactor.contract.Transact(opts, method, params...)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CronUpkeepFactory.Contract.contract.Call(opts, result, method, params...)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.contract.Transfer(opts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.contract.Transact(opts, method, params...)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCaller) CronDelegateAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CronUpkeepFactory.contract.Call(opts, &out, "cronDelegateAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CronUpkeepFactory *CronUpkeepFactorySession) CronDelegateAddress() (common.Address, error) {
	return _CronUpkeepFactory.Contract.CronDelegateAddress(&_CronUpkeepFactory.CallOpts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCallerSession) CronDelegateAddress() (common.Address, error) {
	return _CronUpkeepFactory.Contract.CronDelegateAddress(&_CronUpkeepFactory.CallOpts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCaller) EncodeCronJob(opts *bind.CallOpts, target common.Address, handler []byte, cronString string) ([]byte, error) {
	var out []interface{}
	err := _CronUpkeepFactory.contract.Call(opts, &out, "encodeCronJob", target, handler, cronString)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_CronUpkeepFactory *CronUpkeepFactorySession) EncodeCronJob(target common.Address, handler []byte, cronString string) ([]byte, error) {
	return _CronUpkeepFactory.Contract.EncodeCronJob(&_CronUpkeepFactory.CallOpts, target, handler, cronString)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCallerSession) EncodeCronJob(target common.Address, handler []byte, cronString string) ([]byte, error) {
	return _CronUpkeepFactory.Contract.EncodeCronJob(&_CronUpkeepFactory.CallOpts, target, handler, cronString)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCaller) EncodeCronString(opts *bind.CallOpts, cronString string) ([]byte, error) {
	var out []interface{}
	err := _CronUpkeepFactory.contract.Call(opts, &out, "encodeCronString", cronString)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_CronUpkeepFactory *CronUpkeepFactorySession) EncodeCronString(cronString string) ([]byte, error) {
	return _CronUpkeepFactory.Contract.EncodeCronString(&_CronUpkeepFactory.CallOpts, cronString)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCallerSession) EncodeCronString(cronString string) ([]byte, error) {
	return _CronUpkeepFactory.Contract.EncodeCronString(&_CronUpkeepFactory.CallOpts, cronString)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CronUpkeepFactory.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CronUpkeepFactory *CronUpkeepFactorySession) Owner() (common.Address, error) {
	return _CronUpkeepFactory.Contract.Owner(&_CronUpkeepFactory.CallOpts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCallerSession) Owner() (common.Address, error) {
	return _CronUpkeepFactory.Contract.Owner(&_CronUpkeepFactory.CallOpts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCaller) SMaxJobs(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CronUpkeepFactory.contract.Call(opts, &out, "s_maxJobs")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_CronUpkeepFactory *CronUpkeepFactorySession) SMaxJobs() (*big.Int, error) {
	return _CronUpkeepFactory.Contract.SMaxJobs(&_CronUpkeepFactory.CallOpts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCallerSession) SMaxJobs() (*big.Int, error) {
	return _CronUpkeepFactory.Contract.SMaxJobs(&_CronUpkeepFactory.CallOpts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CronUpkeepFactory.contract.Transact(opts, "acceptOwnership")
}

func (_CronUpkeepFactory *CronUpkeepFactorySession) AcceptOwnership() (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.AcceptOwnership(&_CronUpkeepFactory.TransactOpts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.AcceptOwnership(&_CronUpkeepFactory.TransactOpts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactor) NewCronUpkeep(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CronUpkeepFactory.contract.Transact(opts, "newCronUpkeep")
}

func (_CronUpkeepFactory *CronUpkeepFactorySession) NewCronUpkeep() (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.NewCronUpkeep(&_CronUpkeepFactory.TransactOpts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactorSession) NewCronUpkeep() (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.NewCronUpkeep(&_CronUpkeepFactory.TransactOpts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactor) NewCronUpkeepWithJob(opts *bind.TransactOpts, encodedJob []byte) (*types.Transaction, error) {
	return _CronUpkeepFactory.contract.Transact(opts, "newCronUpkeepWithJob", encodedJob)
}

func (_CronUpkeepFactory *CronUpkeepFactorySession) NewCronUpkeepWithJob(encodedJob []byte) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.NewCronUpkeepWithJob(&_CronUpkeepFactory.TransactOpts, encodedJob)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactorSession) NewCronUpkeepWithJob(encodedJob []byte) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.NewCronUpkeepWithJob(&_CronUpkeepFactory.TransactOpts, encodedJob)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactor) SetMaxJobs(opts *bind.TransactOpts, maxJobs *big.Int) (*types.Transaction, error) {
	return _CronUpkeepFactory.contract.Transact(opts, "setMaxJobs", maxJobs)
}

func (_CronUpkeepFactory *CronUpkeepFactorySession) SetMaxJobs(maxJobs *big.Int) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.SetMaxJobs(&_CronUpkeepFactory.TransactOpts, maxJobs)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactorSession) SetMaxJobs(maxJobs *big.Int) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.SetMaxJobs(&_CronUpkeepFactory.TransactOpts, maxJobs)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _CronUpkeepFactory.contract.Transact(opts, "transferOwnership", to)
}

func (_CronUpkeepFactory *CronUpkeepFactorySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.TransferOwnership(&_CronUpkeepFactory.TransactOpts, to)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.TransferOwnership(&_CronUpkeepFactory.TransactOpts, to)
}

type CronUpkeepFactoryNewCronUpkeepCreatedIterator struct {
	Event *CronUpkeepFactoryNewCronUpkeepCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CronUpkeepFactoryNewCronUpkeepCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CronUpkeepFactoryNewCronUpkeepCreated)
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
		it.Event = new(CronUpkeepFactoryNewCronUpkeepCreated)
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

func (it *CronUpkeepFactoryNewCronUpkeepCreatedIterator) Error() error {
	return it.fail
}

func (it *CronUpkeepFactoryNewCronUpkeepCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CronUpkeepFactoryNewCronUpkeepCreated struct {
	Upkeep common.Address
	Owner  common.Address
	Raw    types.Log
}

func (_CronUpkeepFactory *CronUpkeepFactoryFilterer) FilterNewCronUpkeepCreated(opts *bind.FilterOpts) (*CronUpkeepFactoryNewCronUpkeepCreatedIterator, error) {

	logs, sub, err := _CronUpkeepFactory.contract.FilterLogs(opts, "NewCronUpkeepCreated")
	if err != nil {
		return nil, err
	}
	return &CronUpkeepFactoryNewCronUpkeepCreatedIterator{contract: _CronUpkeepFactory.contract, event: "NewCronUpkeepCreated", logs: logs, sub: sub}, nil
}

func (_CronUpkeepFactory *CronUpkeepFactoryFilterer) WatchNewCronUpkeepCreated(opts *bind.WatchOpts, sink chan<- *CronUpkeepFactoryNewCronUpkeepCreated) (event.Subscription, error) {

	logs, sub, err := _CronUpkeepFactory.contract.WatchLogs(opts, "NewCronUpkeepCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CronUpkeepFactoryNewCronUpkeepCreated)
				if err := _CronUpkeepFactory.contract.UnpackLog(event, "NewCronUpkeepCreated", log); err != nil {
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

func (_CronUpkeepFactory *CronUpkeepFactoryFilterer) ParseNewCronUpkeepCreated(log types.Log) (*CronUpkeepFactoryNewCronUpkeepCreated, error) {
	event := new(CronUpkeepFactoryNewCronUpkeepCreated)
	if err := _CronUpkeepFactory.contract.UnpackLog(event, "NewCronUpkeepCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CronUpkeepFactoryOwnershipTransferRequestedIterator struct {
	Event *CronUpkeepFactoryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CronUpkeepFactoryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CronUpkeepFactoryOwnershipTransferRequested)
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
		it.Event = new(CronUpkeepFactoryOwnershipTransferRequested)
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

func (it *CronUpkeepFactoryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *CronUpkeepFactoryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CronUpkeepFactoryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CronUpkeepFactory *CronUpkeepFactoryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CronUpkeepFactoryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CronUpkeepFactory.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepFactoryOwnershipTransferRequestedIterator{contract: _CronUpkeepFactory.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_CronUpkeepFactory *CronUpkeepFactoryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CronUpkeepFactoryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CronUpkeepFactory.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CronUpkeepFactoryOwnershipTransferRequested)
				if err := _CronUpkeepFactory.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_CronUpkeepFactory *CronUpkeepFactoryFilterer) ParseOwnershipTransferRequested(log types.Log) (*CronUpkeepFactoryOwnershipTransferRequested, error) {
	event := new(CronUpkeepFactoryOwnershipTransferRequested)
	if err := _CronUpkeepFactory.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CronUpkeepFactoryOwnershipTransferredIterator struct {
	Event *CronUpkeepFactoryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CronUpkeepFactoryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CronUpkeepFactoryOwnershipTransferred)
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
		it.Event = new(CronUpkeepFactoryOwnershipTransferred)
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

func (it *CronUpkeepFactoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *CronUpkeepFactoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CronUpkeepFactoryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CronUpkeepFactory *CronUpkeepFactoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CronUpkeepFactoryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CronUpkeepFactory.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepFactoryOwnershipTransferredIterator{contract: _CronUpkeepFactory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_CronUpkeepFactory *CronUpkeepFactoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CronUpkeepFactoryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CronUpkeepFactory.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CronUpkeepFactoryOwnershipTransferred)
				if err := _CronUpkeepFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_CronUpkeepFactory *CronUpkeepFactoryFilterer) ParseOwnershipTransferred(log types.Log) (*CronUpkeepFactoryOwnershipTransferred, error) {
	event := new(CronUpkeepFactoryOwnershipTransferred)
	if err := _CronUpkeepFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_CronUpkeepFactory *CronUpkeepFactory) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _CronUpkeepFactory.abi.Events["NewCronUpkeepCreated"].ID:
		return _CronUpkeepFactory.ParseNewCronUpkeepCreated(log)
	case _CronUpkeepFactory.abi.Events["OwnershipTransferRequested"].ID:
		return _CronUpkeepFactory.ParseOwnershipTransferRequested(log)
	case _CronUpkeepFactory.abi.Events["OwnershipTransferred"].ID:
		return _CronUpkeepFactory.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CronUpkeepFactoryNewCronUpkeepCreated) Topic() common.Hash {
	return common.HexToHash("0x959d571686b1c9343b61bdc3c0459760cb9695fcd4c4c64845e3b2cdd6865ced")
}

func (CronUpkeepFactoryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (CronUpkeepFactoryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_CronUpkeepFactory *CronUpkeepFactory) Address() common.Address {
	return _CronUpkeepFactory.address
}

type CronUpkeepFactoryInterface interface {
	CronDelegateAddress(opts *bind.CallOpts) (common.Address, error)

	EncodeCronJob(opts *bind.CallOpts, target common.Address, handler []byte, cronString string) ([]byte, error)

	EncodeCronString(opts *bind.CallOpts, cronString string) ([]byte, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SMaxJobs(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	NewCronUpkeep(opts *bind.TransactOpts) (*types.Transaction, error)

	NewCronUpkeepWithJob(opts *bind.TransactOpts, encodedJob []byte) (*types.Transaction, error)

	SetMaxJobs(opts *bind.TransactOpts, maxJobs *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterNewCronUpkeepCreated(opts *bind.FilterOpts) (*CronUpkeepFactoryNewCronUpkeepCreatedIterator, error)

	WatchNewCronUpkeepCreated(opts *bind.WatchOpts, sink chan<- *CronUpkeepFactoryNewCronUpkeepCreated) (event.Subscription, error)

	ParseNewCronUpkeepCreated(log types.Log) (*CronUpkeepFactoryNewCronUpkeepCreated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CronUpkeepFactoryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CronUpkeepFactoryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*CronUpkeepFactoryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CronUpkeepFactoryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CronUpkeepFactoryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*CronUpkeepFactoryOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
