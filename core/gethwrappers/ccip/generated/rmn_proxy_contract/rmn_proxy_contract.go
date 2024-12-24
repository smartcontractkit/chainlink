// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rmn_proxy_contract

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

var RMNProxyMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"arm\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getARM\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setARM\",\"inputs\":[{\"name\":\"arm\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ARMSet\",\"inputs\":[{\"name\":\"arm\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x60803461012f5760009061079a3881900390601f8201601f191683016001600160401b0381118482101761011b579180849260209460405283398101031261011757516001600160a01b03811691908290036101145733156100cf5780546001600160a01b0319163317905580156100be57600280546001600160a01b031916821790556040519081527fef31f568d741a833c6a9dc85a6e1c65e06fa772740d5dc94d1da21827a4e0cab90602090a160405161066590816101358239f35b6342bcdf7f60e11b60005260046000fd5b60405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f00000000000000006044820152606490fd5b80fd5b5080fd5b634e487b7160e01b85526041600452602485fd5b600080fdfe6080604052600436101561001a575b3415610598575b600080fd5b60003560e01c8063181f5a771461007a5780632e90aa2114610075578063458fec3b1461007057806379ba50971461006b5780638da5cb5b146100665763f2fde38b0361000e5761049b565b610449565b6102ee565b61023b565b61019c565b346100155760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261001557604051604081019080821067ffffffffffffffff8311176101055761010191604052600e81527f41524d50726f787920312e302e30000000000000000000000000000000000000602082015260405191829182610134565b0390f35b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b9190916020815282519283602083015260005b8481106101865750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006040809697860101520116010190565b8060208092840101516040828601015201610147565b346100155760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261001557602073ffffffffffffffffffffffffffffffffffffffff60025416604051908152f35b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc60209101126100155760043573ffffffffffffffffffffffffffffffffffffffff811681036100155790565b346100155773ffffffffffffffffffffffffffffffffffffffff61025e366101ee565b6102666105d9565b1680156102c4576020817fef31f568d741a833c6a9dc85a6e1c65e06fa772740d5dc94d1da21827a4e0cab927fffffffffffffffffffffffff00000000000000000000000000000000000000006002541617600255604051908152a1005b7f8579befe0000000000000000000000000000000000000000000000000000000060005260046000fd5b346100155760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100155773ffffffffffffffffffffffffffffffffffffffff6001541633036103eb5760005473ffffffffffffffffffffffffffffffffffffffff16600080547fffffffffffffffffffffffff000000000000000000000000000000000000000016331790556103ac7fffffffffffffffffffffffff000000000000000000000000000000000000000060015416600155565b73ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152fd5b346100155760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261001557602073ffffffffffffffffffffffffffffffffffffffff60005416604051908152f35b346100155773ffffffffffffffffffffffffffffffffffffffff6104be366101ee565b6104c66105d9565b1633811461053a57807fffffffffffffffffffffffff0000000000000000000000000000000000000000600154161760015573ffffffffffffffffffffffffffffffffffffffff8060005416167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152fd5b73ffffffffffffffffffffffffffffffffffffffff60025416803b15610015576000809136828037818036925af13d6000803e6105d4573d6000fd5b3d6000f35b73ffffffffffffffffffffffffffffffffffffffff6000541633036105fa57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152fdfea164736f6c634300081a000a",
}

var RMNProxyABI = RMNProxyMetaData.ABI

var RMNProxyBin = RMNProxyMetaData.Bin

func DeployRMNProxy(auth *bind.TransactOpts, backend bind.ContractBackend, arm common.Address) (common.Address, *types.Transaction, *RMNProxy, error) {
	parsed, err := RMNProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RMNProxyBin), backend, arm)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RMNProxy{address: address, abi: *parsed, RMNProxyCaller: RMNProxyCaller{contract: contract}, RMNProxyTransactor: RMNProxyTransactor{contract: contract}, RMNProxyFilterer: RMNProxyFilterer{contract: contract}}, nil
}

type RMNProxy struct {
	address common.Address
	abi     abi.ABI
	RMNProxyCaller
	RMNProxyTransactor
	RMNProxyFilterer
}

type RMNProxyCaller struct {
	contract *bind.BoundContract
}

type RMNProxyTransactor struct {
	contract *bind.BoundContract
}

type RMNProxyFilterer struct {
	contract *bind.BoundContract
}

type RMNProxySession struct {
	Contract     *RMNProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type RMNProxyCallerSession struct {
	Contract *RMNProxyCaller
	CallOpts bind.CallOpts
}

type RMNProxyTransactorSession struct {
	Contract     *RMNProxyTransactor
	TransactOpts bind.TransactOpts
}

type RMNProxyRaw struct {
	Contract *RMNProxy
}

type RMNProxyCallerRaw struct {
	Contract *RMNProxyCaller
}

type RMNProxyTransactorRaw struct {
	Contract *RMNProxyTransactor
}

func NewRMNProxy(address common.Address, backend bind.ContractBackend) (*RMNProxy, error) {
	abi, err := abi.JSON(strings.NewReader(RMNProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindRMNProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RMNProxy{address: address, abi: abi, RMNProxyCaller: RMNProxyCaller{contract: contract}, RMNProxyTransactor: RMNProxyTransactor{contract: contract}, RMNProxyFilterer: RMNProxyFilterer{contract: contract}}, nil
}

func NewRMNProxyCaller(address common.Address, caller bind.ContractCaller) (*RMNProxyCaller, error) {
	contract, err := bindRMNProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RMNProxyCaller{contract: contract}, nil
}

func NewRMNProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*RMNProxyTransactor, error) {
	contract, err := bindRMNProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RMNProxyTransactor{contract: contract}, nil
}

func NewRMNProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*RMNProxyFilterer, error) {
	contract, err := bindRMNProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RMNProxyFilterer{contract: contract}, nil
}

func bindRMNProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RMNProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_RMNProxy *RMNProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNProxy.Contract.RMNProxyCaller.contract.Call(opts, result, method, params...)
}

func (_RMNProxy *RMNProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNProxy.Contract.RMNProxyTransactor.contract.Transfer(opts)
}

func (_RMNProxy *RMNProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNProxy.Contract.RMNProxyTransactor.contract.Transact(opts, method, params...)
}

func (_RMNProxy *RMNProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_RMNProxy *RMNProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNProxy.Contract.contract.Transfer(opts)
}

func (_RMNProxy *RMNProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNProxy.Contract.contract.Transact(opts, method, params...)
}

func (_RMNProxy *RMNProxyCaller) GetARM(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RMNProxy.contract.Call(opts, &out, "getARM")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_RMNProxy *RMNProxySession) GetARM() (common.Address, error) {
	return _RMNProxy.Contract.GetARM(&_RMNProxy.CallOpts)
}

func (_RMNProxy *RMNProxyCallerSession) GetARM() (common.Address, error) {
	return _RMNProxy.Contract.GetARM(&_RMNProxy.CallOpts)
}

func (_RMNProxy *RMNProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RMNProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_RMNProxy *RMNProxySession) Owner() (common.Address, error) {
	return _RMNProxy.Contract.Owner(&_RMNProxy.CallOpts)
}

func (_RMNProxy *RMNProxyCallerSession) Owner() (common.Address, error) {
	return _RMNProxy.Contract.Owner(&_RMNProxy.CallOpts)
}

func (_RMNProxy *RMNProxyCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RMNProxy.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_RMNProxy *RMNProxySession) TypeAndVersion() (string, error) {
	return _RMNProxy.Contract.TypeAndVersion(&_RMNProxy.CallOpts)
}

func (_RMNProxy *RMNProxyCallerSession) TypeAndVersion() (string, error) {
	return _RMNProxy.Contract.TypeAndVersion(&_RMNProxy.CallOpts)
}

func (_RMNProxy *RMNProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNProxy.contract.Transact(opts, "acceptOwnership")
}

func (_RMNProxy *RMNProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNProxy.Contract.AcceptOwnership(&_RMNProxy.TransactOpts)
}

func (_RMNProxy *RMNProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNProxy.Contract.AcceptOwnership(&_RMNProxy.TransactOpts)
}

func (_RMNProxy *RMNProxyTransactor) SetARM(opts *bind.TransactOpts, arm common.Address) (*types.Transaction, error) {
	return _RMNProxy.contract.Transact(opts, "setARM", arm)
}

func (_RMNProxy *RMNProxySession) SetARM(arm common.Address) (*types.Transaction, error) {
	return _RMNProxy.Contract.SetARM(&_RMNProxy.TransactOpts, arm)
}

func (_RMNProxy *RMNProxyTransactorSession) SetARM(arm common.Address) (*types.Transaction, error) {
	return _RMNProxy.Contract.SetARM(&_RMNProxy.TransactOpts, arm)
}

func (_RMNProxy *RMNProxyTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _RMNProxy.contract.Transact(opts, "transferOwnership", to)
}

func (_RMNProxy *RMNProxySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNProxy.Contract.TransferOwnership(&_RMNProxy.TransactOpts, to)
}

func (_RMNProxy *RMNProxyTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNProxy.Contract.TransferOwnership(&_RMNProxy.TransactOpts, to)
}

func (_RMNProxy *RMNProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _RMNProxy.contract.RawTransact(opts, calldata)
}

func (_RMNProxy *RMNProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _RMNProxy.Contract.Fallback(&_RMNProxy.TransactOpts, calldata)
}

func (_RMNProxy *RMNProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _RMNProxy.Contract.Fallback(&_RMNProxy.TransactOpts, calldata)
}

type RMNProxyARMSetIterator struct {
	Event *RMNProxyARMSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNProxyARMSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNProxyARMSet)
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
		it.Event = new(RMNProxyARMSet)
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

func (it *RMNProxyARMSetIterator) Error() error {
	return it.fail
}

func (it *RMNProxyARMSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNProxyARMSet struct {
	Arm common.Address
	Raw types.Log
}

func (_RMNProxy *RMNProxyFilterer) FilterARMSet(opts *bind.FilterOpts) (*RMNProxyARMSetIterator, error) {

	logs, sub, err := _RMNProxy.contract.FilterLogs(opts, "ARMSet")
	if err != nil {
		return nil, err
	}
	return &RMNProxyARMSetIterator{contract: _RMNProxy.contract, event: "ARMSet", logs: logs, sub: sub}, nil
}

func (_RMNProxy *RMNProxyFilterer) WatchARMSet(opts *bind.WatchOpts, sink chan<- *RMNProxyARMSet) (event.Subscription, error) {

	logs, sub, err := _RMNProxy.contract.WatchLogs(opts, "ARMSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNProxyARMSet)
				if err := _RMNProxy.contract.UnpackLog(event, "ARMSet", log); err != nil {
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

func (_RMNProxy *RMNProxyFilterer) ParseARMSet(log types.Log) (*RMNProxyARMSet, error) {
	event := new(RMNProxyARMSet)
	if err := _RMNProxy.contract.UnpackLog(event, "ARMSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNProxyOwnershipTransferRequestedIterator struct {
	Event *RMNProxyOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNProxyOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNProxyOwnershipTransferRequested)
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
		it.Event = new(RMNProxyOwnershipTransferRequested)
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

func (it *RMNProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *RMNProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNProxy *RMNProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNProxyOwnershipTransferRequestedIterator{contract: _RMNProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_RMNProxy *RMNProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNProxyOwnershipTransferRequested)
				if err := _RMNProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_RMNProxy *RMNProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*RMNProxyOwnershipTransferRequested, error) {
	event := new(RMNProxyOwnershipTransferRequested)
	if err := _RMNProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNProxyOwnershipTransferredIterator struct {
	Event *RMNProxyOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNProxyOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNProxyOwnershipTransferred)
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
		it.Event = new(RMNProxyOwnershipTransferred)
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

func (it *RMNProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *RMNProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNProxy *RMNProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNProxyOwnershipTransferredIterator{contract: _RMNProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_RMNProxy *RMNProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNProxyOwnershipTransferred)
				if err := _RMNProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_RMNProxy *RMNProxyFilterer) ParseOwnershipTransferred(log types.Log) (*RMNProxyOwnershipTransferred, error) {
	event := new(RMNProxyOwnershipTransferred)
	if err := _RMNProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_RMNProxy *RMNProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _RMNProxy.abi.Events["ARMSet"].ID:
		return _RMNProxy.ParseARMSet(log)
	case _RMNProxy.abi.Events["OwnershipTransferRequested"].ID:
		return _RMNProxy.ParseOwnershipTransferRequested(log)
	case _RMNProxy.abi.Events["OwnershipTransferred"].ID:
		return _RMNProxy.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (RMNProxyARMSet) Topic() common.Hash {
	return common.HexToHash("0xef31f568d741a833c6a9dc85a6e1c65e06fa772740d5dc94d1da21827a4e0cab")
}

func (RMNProxyOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (RMNProxyOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_RMNProxy *RMNProxy) Address() common.Address {
	return _RMNProxy.address
}

type RMNProxyInterface interface {
	GetARM(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetARM(opts *bind.TransactOpts, arm common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterARMSet(opts *bind.FilterOpts) (*RMNProxyARMSetIterator, error)

	WatchARMSet(opts *bind.WatchOpts, sink chan<- *RMNProxyARMSet) (event.Subscription, error)

	ParseARMSet(log types.Log) (*RMNProxyARMSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNProxyOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*RMNProxyOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNProxyOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*RMNProxyOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
