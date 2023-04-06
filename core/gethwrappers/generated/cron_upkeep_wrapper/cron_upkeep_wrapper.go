// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package cron_upkeep_wrapper

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

var CronUpkeepMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"maxJobs\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"firstJob\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"CronJobIDNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExceedsMaxJobs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidHandler\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TickDoesntMatchSpec\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TickInFuture\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TickTooOld\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnknownFieldType\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"handler\",\"type\":\"bytes\"}],\"name\":\"CronJobCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"CronJobDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"CronJobExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"handler\",\"type\":\"bytes\"}],\"name\":\"CronJobUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"handler\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"encodedCronSpec\",\"type\":\"bytes\"}],\"name\":\"createCronJobFromEncodedSpec\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"deleteCronJob\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveCronJobIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getCronJob\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"handler\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"cronString\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"nextTick\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_maxJobs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newTarget\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"newHandler\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"newEncodedCronSpec\",\"type\":\"bytes\"}],\"name\":\"updateCronJob\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

var CronUpkeepABI = CronUpkeepMetaData.ABI

type CronUpkeep struct {
	address common.Address
	abi     abi.ABI
	CronUpkeepCaller
	CronUpkeepTransactor
	CronUpkeepFilterer
}

type CronUpkeepCaller struct {
	contract *bind.BoundContract
}

type CronUpkeepTransactor struct {
	contract *bind.BoundContract
}

type CronUpkeepFilterer struct {
	contract *bind.BoundContract
}

type CronUpkeepSession struct {
	Contract     *CronUpkeep
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CronUpkeepCallerSession struct {
	Contract *CronUpkeepCaller
	CallOpts bind.CallOpts
}

type CronUpkeepTransactorSession struct {
	Contract     *CronUpkeepTransactor
	TransactOpts bind.TransactOpts
}

type CronUpkeepRaw struct {
	Contract *CronUpkeep
}

type CronUpkeepCallerRaw struct {
	Contract *CronUpkeepCaller
}

type CronUpkeepTransactorRaw struct {
	Contract *CronUpkeepTransactor
}

func NewCronUpkeep(address common.Address, backend bind.ContractBackend) (*CronUpkeep, error) {
	abi, err := abi.JSON(strings.NewReader(CronUpkeepABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCronUpkeep(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CronUpkeep{address: address, abi: abi, CronUpkeepCaller: CronUpkeepCaller{contract: contract}, CronUpkeepTransactor: CronUpkeepTransactor{contract: contract}, CronUpkeepFilterer: CronUpkeepFilterer{contract: contract}}, nil
}

func NewCronUpkeepCaller(address common.Address, caller bind.ContractCaller) (*CronUpkeepCaller, error) {
	contract, err := bindCronUpkeep(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepCaller{contract: contract}, nil
}

func NewCronUpkeepTransactor(address common.Address, transactor bind.ContractTransactor) (*CronUpkeepTransactor, error) {
	contract, err := bindCronUpkeep(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepTransactor{contract: contract}, nil
}

func NewCronUpkeepFilterer(address common.Address, filterer bind.ContractFilterer) (*CronUpkeepFilterer, error) {
	contract, err := bindCronUpkeep(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepFilterer{contract: contract}, nil
}

func bindCronUpkeep(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CronUpkeepMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_CronUpkeep *CronUpkeepRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CronUpkeep.Contract.CronUpkeepCaller.contract.Call(opts, result, method, params...)
}

func (_CronUpkeep *CronUpkeepRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CronUpkeep.Contract.CronUpkeepTransactor.contract.Transfer(opts)
}

func (_CronUpkeep *CronUpkeepRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CronUpkeep.Contract.CronUpkeepTransactor.contract.Transact(opts, method, params...)
}

func (_CronUpkeep *CronUpkeepCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CronUpkeep.Contract.contract.Call(opts, result, method, params...)
}

func (_CronUpkeep *CronUpkeepTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CronUpkeep.Contract.contract.Transfer(opts)
}

func (_CronUpkeep *CronUpkeepTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CronUpkeep.Contract.contract.Transact(opts, method, params...)
}

func (_CronUpkeep *CronUpkeepCaller) GetActiveCronJobIDs(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _CronUpkeep.contract.Call(opts, &out, "getActiveCronJobIDs")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_CronUpkeep *CronUpkeepSession) GetActiveCronJobIDs() ([]*big.Int, error) {
	return _CronUpkeep.Contract.GetActiveCronJobIDs(&_CronUpkeep.CallOpts)
}

func (_CronUpkeep *CronUpkeepCallerSession) GetActiveCronJobIDs() ([]*big.Int, error) {
	return _CronUpkeep.Contract.GetActiveCronJobIDs(&_CronUpkeep.CallOpts)
}

func (_CronUpkeep *CronUpkeepCaller) GetCronJob(opts *bind.CallOpts, id *big.Int) (GetCronJob,

	error) {
	var out []interface{}
	err := _CronUpkeep.contract.Call(opts, &out, "getCronJob", id)

	outstruct := new(GetCronJob)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Target = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Handler = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.CronString = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.NextTick = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_CronUpkeep *CronUpkeepSession) GetCronJob(id *big.Int) (GetCronJob,

	error) {
	return _CronUpkeep.Contract.GetCronJob(&_CronUpkeep.CallOpts, id)
}

func (_CronUpkeep *CronUpkeepCallerSession) GetCronJob(id *big.Int) (GetCronJob,

	error) {
	return _CronUpkeep.Contract.GetCronJob(&_CronUpkeep.CallOpts, id)
}

func (_CronUpkeep *CronUpkeepCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CronUpkeep.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CronUpkeep *CronUpkeepSession) Owner() (common.Address, error) {
	return _CronUpkeep.Contract.Owner(&_CronUpkeep.CallOpts)
}

func (_CronUpkeep *CronUpkeepCallerSession) Owner() (common.Address, error) {
	return _CronUpkeep.Contract.Owner(&_CronUpkeep.CallOpts)
}

func (_CronUpkeep *CronUpkeepCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _CronUpkeep.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_CronUpkeep *CronUpkeepSession) Paused() (bool, error) {
	return _CronUpkeep.Contract.Paused(&_CronUpkeep.CallOpts)
}

func (_CronUpkeep *CronUpkeepCallerSession) Paused() (bool, error) {
	return _CronUpkeep.Contract.Paused(&_CronUpkeep.CallOpts)
}

func (_CronUpkeep *CronUpkeepCaller) SMaxJobs(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CronUpkeep.contract.Call(opts, &out, "s_maxJobs")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_CronUpkeep *CronUpkeepSession) SMaxJobs() (*big.Int, error) {
	return _CronUpkeep.Contract.SMaxJobs(&_CronUpkeep.CallOpts)
}

func (_CronUpkeep *CronUpkeepCallerSession) SMaxJobs() (*big.Int, error) {
	return _CronUpkeep.Contract.SMaxJobs(&_CronUpkeep.CallOpts)
}

func (_CronUpkeep *CronUpkeepTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CronUpkeep.contract.Transact(opts, "acceptOwnership")
}

func (_CronUpkeep *CronUpkeepSession) AcceptOwnership() (*types.Transaction, error) {
	return _CronUpkeep.Contract.AcceptOwnership(&_CronUpkeep.TransactOpts)
}

func (_CronUpkeep *CronUpkeepTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _CronUpkeep.Contract.AcceptOwnership(&_CronUpkeep.TransactOpts)
}

func (_CronUpkeep *CronUpkeepTransactor) CheckUpkeep(opts *bind.TransactOpts, arg0 []byte) (*types.Transaction, error) {
	return _CronUpkeep.contract.Transact(opts, "checkUpkeep", arg0)
}

func (_CronUpkeep *CronUpkeepSession) CheckUpkeep(arg0 []byte) (*types.Transaction, error) {
	return _CronUpkeep.Contract.CheckUpkeep(&_CronUpkeep.TransactOpts, arg0)
}

func (_CronUpkeep *CronUpkeepTransactorSession) CheckUpkeep(arg0 []byte) (*types.Transaction, error) {
	return _CronUpkeep.Contract.CheckUpkeep(&_CronUpkeep.TransactOpts, arg0)
}

func (_CronUpkeep *CronUpkeepTransactor) CreateCronJobFromEncodedSpec(opts *bind.TransactOpts, target common.Address, handler []byte, encodedCronSpec []byte) (*types.Transaction, error) {
	return _CronUpkeep.contract.Transact(opts, "createCronJobFromEncodedSpec", target, handler, encodedCronSpec)
}

func (_CronUpkeep *CronUpkeepSession) CreateCronJobFromEncodedSpec(target common.Address, handler []byte, encodedCronSpec []byte) (*types.Transaction, error) {
	return _CronUpkeep.Contract.CreateCronJobFromEncodedSpec(&_CronUpkeep.TransactOpts, target, handler, encodedCronSpec)
}

func (_CronUpkeep *CronUpkeepTransactorSession) CreateCronJobFromEncodedSpec(target common.Address, handler []byte, encodedCronSpec []byte) (*types.Transaction, error) {
	return _CronUpkeep.Contract.CreateCronJobFromEncodedSpec(&_CronUpkeep.TransactOpts, target, handler, encodedCronSpec)
}

func (_CronUpkeep *CronUpkeepTransactor) DeleteCronJob(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _CronUpkeep.contract.Transact(opts, "deleteCronJob", id)
}

func (_CronUpkeep *CronUpkeepSession) DeleteCronJob(id *big.Int) (*types.Transaction, error) {
	return _CronUpkeep.Contract.DeleteCronJob(&_CronUpkeep.TransactOpts, id)
}

func (_CronUpkeep *CronUpkeepTransactorSession) DeleteCronJob(id *big.Int) (*types.Transaction, error) {
	return _CronUpkeep.Contract.DeleteCronJob(&_CronUpkeep.TransactOpts, id)
}

func (_CronUpkeep *CronUpkeepTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CronUpkeep.contract.Transact(opts, "pause")
}

func (_CronUpkeep *CronUpkeepSession) Pause() (*types.Transaction, error) {
	return _CronUpkeep.Contract.Pause(&_CronUpkeep.TransactOpts)
}

func (_CronUpkeep *CronUpkeepTransactorSession) Pause() (*types.Transaction, error) {
	return _CronUpkeep.Contract.Pause(&_CronUpkeep.TransactOpts)
}

func (_CronUpkeep *CronUpkeepTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _CronUpkeep.contract.Transact(opts, "performUpkeep", performData)
}

func (_CronUpkeep *CronUpkeepSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _CronUpkeep.Contract.PerformUpkeep(&_CronUpkeep.TransactOpts, performData)
}

func (_CronUpkeep *CronUpkeepTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _CronUpkeep.Contract.PerformUpkeep(&_CronUpkeep.TransactOpts, performData)
}

func (_CronUpkeep *CronUpkeepTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _CronUpkeep.contract.Transact(opts, "transferOwnership", to)
}

func (_CronUpkeep *CronUpkeepSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CronUpkeep.Contract.TransferOwnership(&_CronUpkeep.TransactOpts, to)
}

func (_CronUpkeep *CronUpkeepTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CronUpkeep.Contract.TransferOwnership(&_CronUpkeep.TransactOpts, to)
}

func (_CronUpkeep *CronUpkeepTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CronUpkeep.contract.Transact(opts, "unpause")
}

func (_CronUpkeep *CronUpkeepSession) Unpause() (*types.Transaction, error) {
	return _CronUpkeep.Contract.Unpause(&_CronUpkeep.TransactOpts)
}

func (_CronUpkeep *CronUpkeepTransactorSession) Unpause() (*types.Transaction, error) {
	return _CronUpkeep.Contract.Unpause(&_CronUpkeep.TransactOpts)
}

func (_CronUpkeep *CronUpkeepTransactor) UpdateCronJob(opts *bind.TransactOpts, id *big.Int, newTarget common.Address, newHandler []byte, newEncodedCronSpec []byte) (*types.Transaction, error) {
	return _CronUpkeep.contract.Transact(opts, "updateCronJob", id, newTarget, newHandler, newEncodedCronSpec)
}

func (_CronUpkeep *CronUpkeepSession) UpdateCronJob(id *big.Int, newTarget common.Address, newHandler []byte, newEncodedCronSpec []byte) (*types.Transaction, error) {
	return _CronUpkeep.Contract.UpdateCronJob(&_CronUpkeep.TransactOpts, id, newTarget, newHandler, newEncodedCronSpec)
}

func (_CronUpkeep *CronUpkeepTransactorSession) UpdateCronJob(id *big.Int, newTarget common.Address, newHandler []byte, newEncodedCronSpec []byte) (*types.Transaction, error) {
	return _CronUpkeep.Contract.UpdateCronJob(&_CronUpkeep.TransactOpts, id, newTarget, newHandler, newEncodedCronSpec)
}

func (_CronUpkeep *CronUpkeepTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _CronUpkeep.contract.RawTransact(opts, calldata)
}

func (_CronUpkeep *CronUpkeepSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _CronUpkeep.Contract.Fallback(&_CronUpkeep.TransactOpts, calldata)
}

func (_CronUpkeep *CronUpkeepTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _CronUpkeep.Contract.Fallback(&_CronUpkeep.TransactOpts, calldata)
}

func (_CronUpkeep *CronUpkeepTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CronUpkeep.contract.RawTransact(opts, nil)
}

func (_CronUpkeep *CronUpkeepSession) Receive() (*types.Transaction, error) {
	return _CronUpkeep.Contract.Receive(&_CronUpkeep.TransactOpts)
}

func (_CronUpkeep *CronUpkeepTransactorSession) Receive() (*types.Transaction, error) {
	return _CronUpkeep.Contract.Receive(&_CronUpkeep.TransactOpts)
}

type CronUpkeepCronJobCreatedIterator struct {
	Event *CronUpkeepCronJobCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CronUpkeepCronJobCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CronUpkeepCronJobCreated)
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
		it.Event = new(CronUpkeepCronJobCreated)
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

func (it *CronUpkeepCronJobCreatedIterator) Error() error {
	return it.fail
}

func (it *CronUpkeepCronJobCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CronUpkeepCronJobCreated struct {
	Id      *big.Int
	Target  common.Address
	Handler []byte
	Raw     types.Log
}

func (_CronUpkeep *CronUpkeepFilterer) FilterCronJobCreated(opts *bind.FilterOpts, id []*big.Int) (*CronUpkeepCronJobCreatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _CronUpkeep.contract.FilterLogs(opts, "CronJobCreated", idRule)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepCronJobCreatedIterator{contract: _CronUpkeep.contract, event: "CronJobCreated", logs: logs, sub: sub}, nil
}

func (_CronUpkeep *CronUpkeepFilterer) WatchCronJobCreated(opts *bind.WatchOpts, sink chan<- *CronUpkeepCronJobCreated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _CronUpkeep.contract.WatchLogs(opts, "CronJobCreated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CronUpkeepCronJobCreated)
				if err := _CronUpkeep.contract.UnpackLog(event, "CronJobCreated", log); err != nil {
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

func (_CronUpkeep *CronUpkeepFilterer) ParseCronJobCreated(log types.Log) (*CronUpkeepCronJobCreated, error) {
	event := new(CronUpkeepCronJobCreated)
	if err := _CronUpkeep.contract.UnpackLog(event, "CronJobCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CronUpkeepCronJobDeletedIterator struct {
	Event *CronUpkeepCronJobDeleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CronUpkeepCronJobDeletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CronUpkeepCronJobDeleted)
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
		it.Event = new(CronUpkeepCronJobDeleted)
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

func (it *CronUpkeepCronJobDeletedIterator) Error() error {
	return it.fail
}

func (it *CronUpkeepCronJobDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CronUpkeepCronJobDeleted struct {
	Id  *big.Int
	Raw types.Log
}

func (_CronUpkeep *CronUpkeepFilterer) FilterCronJobDeleted(opts *bind.FilterOpts, id []*big.Int) (*CronUpkeepCronJobDeletedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _CronUpkeep.contract.FilterLogs(opts, "CronJobDeleted", idRule)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepCronJobDeletedIterator{contract: _CronUpkeep.contract, event: "CronJobDeleted", logs: logs, sub: sub}, nil
}

func (_CronUpkeep *CronUpkeepFilterer) WatchCronJobDeleted(opts *bind.WatchOpts, sink chan<- *CronUpkeepCronJobDeleted, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _CronUpkeep.contract.WatchLogs(opts, "CronJobDeleted", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CronUpkeepCronJobDeleted)
				if err := _CronUpkeep.contract.UnpackLog(event, "CronJobDeleted", log); err != nil {
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

func (_CronUpkeep *CronUpkeepFilterer) ParseCronJobDeleted(log types.Log) (*CronUpkeepCronJobDeleted, error) {
	event := new(CronUpkeepCronJobDeleted)
	if err := _CronUpkeep.contract.UnpackLog(event, "CronJobDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CronUpkeepCronJobExecutedIterator struct {
	Event *CronUpkeepCronJobExecuted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CronUpkeepCronJobExecutedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CronUpkeepCronJobExecuted)
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
		it.Event = new(CronUpkeepCronJobExecuted)
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

func (it *CronUpkeepCronJobExecutedIterator) Error() error {
	return it.fail
}

func (it *CronUpkeepCronJobExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CronUpkeepCronJobExecuted struct {
	Id      *big.Int
	Success bool
	Raw     types.Log
}

func (_CronUpkeep *CronUpkeepFilterer) FilterCronJobExecuted(opts *bind.FilterOpts, id []*big.Int) (*CronUpkeepCronJobExecutedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _CronUpkeep.contract.FilterLogs(opts, "CronJobExecuted", idRule)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepCronJobExecutedIterator{contract: _CronUpkeep.contract, event: "CronJobExecuted", logs: logs, sub: sub}, nil
}

func (_CronUpkeep *CronUpkeepFilterer) WatchCronJobExecuted(opts *bind.WatchOpts, sink chan<- *CronUpkeepCronJobExecuted, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _CronUpkeep.contract.WatchLogs(opts, "CronJobExecuted", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CronUpkeepCronJobExecuted)
				if err := _CronUpkeep.contract.UnpackLog(event, "CronJobExecuted", log); err != nil {
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

func (_CronUpkeep *CronUpkeepFilterer) ParseCronJobExecuted(log types.Log) (*CronUpkeepCronJobExecuted, error) {
	event := new(CronUpkeepCronJobExecuted)
	if err := _CronUpkeep.contract.UnpackLog(event, "CronJobExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CronUpkeepCronJobUpdatedIterator struct {
	Event *CronUpkeepCronJobUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CronUpkeepCronJobUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CronUpkeepCronJobUpdated)
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
		it.Event = new(CronUpkeepCronJobUpdated)
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

func (it *CronUpkeepCronJobUpdatedIterator) Error() error {
	return it.fail
}

func (it *CronUpkeepCronJobUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CronUpkeepCronJobUpdated struct {
	Id      *big.Int
	Target  common.Address
	Handler []byte
	Raw     types.Log
}

func (_CronUpkeep *CronUpkeepFilterer) FilterCronJobUpdated(opts *bind.FilterOpts, id []*big.Int) (*CronUpkeepCronJobUpdatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _CronUpkeep.contract.FilterLogs(opts, "CronJobUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepCronJobUpdatedIterator{contract: _CronUpkeep.contract, event: "CronJobUpdated", logs: logs, sub: sub}, nil
}

func (_CronUpkeep *CronUpkeepFilterer) WatchCronJobUpdated(opts *bind.WatchOpts, sink chan<- *CronUpkeepCronJobUpdated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _CronUpkeep.contract.WatchLogs(opts, "CronJobUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CronUpkeepCronJobUpdated)
				if err := _CronUpkeep.contract.UnpackLog(event, "CronJobUpdated", log); err != nil {
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

func (_CronUpkeep *CronUpkeepFilterer) ParseCronJobUpdated(log types.Log) (*CronUpkeepCronJobUpdated, error) {
	event := new(CronUpkeepCronJobUpdated)
	if err := _CronUpkeep.contract.UnpackLog(event, "CronJobUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CronUpkeepOwnershipTransferRequestedIterator struct {
	Event *CronUpkeepOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CronUpkeepOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CronUpkeepOwnershipTransferRequested)
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
		it.Event = new(CronUpkeepOwnershipTransferRequested)
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

func (it *CronUpkeepOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *CronUpkeepOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CronUpkeepOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CronUpkeep *CronUpkeepFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CronUpkeepOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CronUpkeep.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepOwnershipTransferRequestedIterator{contract: _CronUpkeep.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_CronUpkeep *CronUpkeepFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CronUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CronUpkeep.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CronUpkeepOwnershipTransferRequested)
				if err := _CronUpkeep.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_CronUpkeep *CronUpkeepFilterer) ParseOwnershipTransferRequested(log types.Log) (*CronUpkeepOwnershipTransferRequested, error) {
	event := new(CronUpkeepOwnershipTransferRequested)
	if err := _CronUpkeep.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CronUpkeepOwnershipTransferredIterator struct {
	Event *CronUpkeepOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CronUpkeepOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CronUpkeepOwnershipTransferred)
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
		it.Event = new(CronUpkeepOwnershipTransferred)
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

func (it *CronUpkeepOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *CronUpkeepOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CronUpkeepOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CronUpkeep *CronUpkeepFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CronUpkeepOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CronUpkeep.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepOwnershipTransferredIterator{contract: _CronUpkeep.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_CronUpkeep *CronUpkeepFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CronUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CronUpkeep.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CronUpkeepOwnershipTransferred)
				if err := _CronUpkeep.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_CronUpkeep *CronUpkeepFilterer) ParseOwnershipTransferred(log types.Log) (*CronUpkeepOwnershipTransferred, error) {
	event := new(CronUpkeepOwnershipTransferred)
	if err := _CronUpkeep.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CronUpkeepPausedIterator struct {
	Event *CronUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CronUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CronUpkeepPaused)
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
		it.Event = new(CronUpkeepPaused)
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

func (it *CronUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *CronUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CronUpkeepPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_CronUpkeep *CronUpkeepFilterer) FilterPaused(opts *bind.FilterOpts) (*CronUpkeepPausedIterator, error) {

	logs, sub, err := _CronUpkeep.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &CronUpkeepPausedIterator{contract: _CronUpkeep.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_CronUpkeep *CronUpkeepFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *CronUpkeepPaused) (event.Subscription, error) {

	logs, sub, err := _CronUpkeep.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CronUpkeepPaused)
				if err := _CronUpkeep.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_CronUpkeep *CronUpkeepFilterer) ParsePaused(log types.Log) (*CronUpkeepPaused, error) {
	event := new(CronUpkeepPaused)
	if err := _CronUpkeep.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CronUpkeepUnpausedIterator struct {
	Event *CronUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CronUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CronUpkeepUnpaused)
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
		it.Event = new(CronUpkeepUnpaused)
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

func (it *CronUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *CronUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CronUpkeepUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_CronUpkeep *CronUpkeepFilterer) FilterUnpaused(opts *bind.FilterOpts) (*CronUpkeepUnpausedIterator, error) {

	logs, sub, err := _CronUpkeep.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &CronUpkeepUnpausedIterator{contract: _CronUpkeep.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_CronUpkeep *CronUpkeepFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *CronUpkeepUnpaused) (event.Subscription, error) {

	logs, sub, err := _CronUpkeep.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CronUpkeepUnpaused)
				if err := _CronUpkeep.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_CronUpkeep *CronUpkeepFilterer) ParseUnpaused(log types.Log) (*CronUpkeepUnpaused, error) {
	event := new(CronUpkeepUnpaused)
	if err := _CronUpkeep.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetCronJob struct {
	Target     common.Address
	Handler    []byte
	CronString string
	NextTick   *big.Int
}

func (_CronUpkeep *CronUpkeep) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _CronUpkeep.abi.Events["CronJobCreated"].ID:
		return _CronUpkeep.ParseCronJobCreated(log)
	case _CronUpkeep.abi.Events["CronJobDeleted"].ID:
		return _CronUpkeep.ParseCronJobDeleted(log)
	case _CronUpkeep.abi.Events["CronJobExecuted"].ID:
		return _CronUpkeep.ParseCronJobExecuted(log)
	case _CronUpkeep.abi.Events["CronJobUpdated"].ID:
		return _CronUpkeep.ParseCronJobUpdated(log)
	case _CronUpkeep.abi.Events["OwnershipTransferRequested"].ID:
		return _CronUpkeep.ParseOwnershipTransferRequested(log)
	case _CronUpkeep.abi.Events["OwnershipTransferred"].ID:
		return _CronUpkeep.ParseOwnershipTransferred(log)
	case _CronUpkeep.abi.Events["Paused"].ID:
		return _CronUpkeep.ParsePaused(log)
	case _CronUpkeep.abi.Events["Unpaused"].ID:
		return _CronUpkeep.ParseUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CronUpkeepCronJobCreated) Topic() common.Hash {
	return common.HexToHash("0xe66fb0bca0f9d6a395d3eaf5f39c6ac87dd34aff4e3f2f9a9b33a46f15589627")
}

func (CronUpkeepCronJobDeleted) Topic() common.Hash {
	return common.HexToHash("0x7aaa5a7c35e162386d922bd67e91ea476d38d9bb931bc369d8b15ab113250974")
}

func (CronUpkeepCronJobExecuted) Topic() common.Hash {
	return common.HexToHash("0x25d1b235668fd0219da15f5fa6054013a53e59c4f3ea31459dc1d4e0b9f23d26")
}

func (CronUpkeepCronJobUpdated) Topic() common.Hash {
	return common.HexToHash("0xeeaf6ad42034ba5357ffd961b8c80bf6cbf53c224020541e46573a3f19ef09a5")
}

func (CronUpkeepOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (CronUpkeepOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (CronUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (CronUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (_CronUpkeep *CronUpkeep) Address() common.Address {
	return _CronUpkeep.address
}

type CronUpkeepInterface interface {
	GetActiveCronJobIDs(opts *bind.CallOpts) ([]*big.Int, error)

	GetCronJob(opts *bind.CallOpts, id *big.Int) (GetCronJob,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Paused(opts *bind.CallOpts) (bool, error)

	SMaxJobs(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	CheckUpkeep(opts *bind.TransactOpts, arg0 []byte) (*types.Transaction, error)

	CreateCronJobFromEncodedSpec(opts *bind.TransactOpts, target common.Address, handler []byte, encodedCronSpec []byte) (*types.Transaction, error)

	DeleteCronJob(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	UpdateCronJob(opts *bind.TransactOpts, id *big.Int, newTarget common.Address, newHandler []byte, newEncodedCronSpec []byte) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterCronJobCreated(opts *bind.FilterOpts, id []*big.Int) (*CronUpkeepCronJobCreatedIterator, error)

	WatchCronJobCreated(opts *bind.WatchOpts, sink chan<- *CronUpkeepCronJobCreated, id []*big.Int) (event.Subscription, error)

	ParseCronJobCreated(log types.Log) (*CronUpkeepCronJobCreated, error)

	FilterCronJobDeleted(opts *bind.FilterOpts, id []*big.Int) (*CronUpkeepCronJobDeletedIterator, error)

	WatchCronJobDeleted(opts *bind.WatchOpts, sink chan<- *CronUpkeepCronJobDeleted, id []*big.Int) (event.Subscription, error)

	ParseCronJobDeleted(log types.Log) (*CronUpkeepCronJobDeleted, error)

	FilterCronJobExecuted(opts *bind.FilterOpts, id []*big.Int) (*CronUpkeepCronJobExecutedIterator, error)

	WatchCronJobExecuted(opts *bind.WatchOpts, sink chan<- *CronUpkeepCronJobExecuted, id []*big.Int) (event.Subscription, error)

	ParseCronJobExecuted(log types.Log) (*CronUpkeepCronJobExecuted, error)

	FilterCronJobUpdated(opts *bind.FilterOpts, id []*big.Int) (*CronUpkeepCronJobUpdatedIterator, error)

	WatchCronJobUpdated(opts *bind.WatchOpts, sink chan<- *CronUpkeepCronJobUpdated, id []*big.Int) (event.Subscription, error)

	ParseCronJobUpdated(log types.Log) (*CronUpkeepCronJobUpdated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CronUpkeepOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CronUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*CronUpkeepOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CronUpkeepOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CronUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*CronUpkeepOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*CronUpkeepPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *CronUpkeepPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*CronUpkeepPaused, error)

	FilterUnpaused(opts *bind.FilterOpts) (*CronUpkeepUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *CronUpkeepUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*CronUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
