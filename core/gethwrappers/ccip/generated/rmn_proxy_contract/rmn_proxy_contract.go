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

var RMNProxyContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"arm\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"arm\",\"type\":\"address\"}],\"name\":\"ARMSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getARM\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"arm\",\"type\":\"address\"}],\"name\":\"setARM\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60803461012f5760009061079a3881900390601f8201601f191683016001600160401b0381118482101761011b579180849260209460405283398101031261011757516001600160a01b03811691908290036101145733156100cf5780546001600160a01b0319163317905580156100be57600280546001600160a01b031916821790556040519081527fef31f568d741a833c6a9dc85a6e1c65e06fa772740d5dc94d1da21827a4e0cab90602090a160405161066590816101358239f35b6342bcdf7f60e11b60005260046000fd5b60405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f00000000000000006044820152606490fd5b80fd5b5080fd5b634e487b7160e01b85526041600452602485fd5b600080fdfe6080604052600436101561001a575b3415610598575b600080fd5b60003560e01c8063181f5a771461007a5780632e90aa2114610075578063458fec3b1461007057806379ba50971461006b5780638da5cb5b146100665763f2fde38b0361000e5761049b565b610449565b6102ee565b61023b565b61019c565b346100155760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261001557604051604081019080821067ffffffffffffffff8311176101055761010191604052600e81527f41524d50726f787920312e302e30000000000000000000000000000000000000602082015260405191829182610134565b0390f35b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b9190916020815282519283602083015260005b8481106101865750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006040809697860101520116010190565b8060208092840101516040828601015201610147565b346100155760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261001557602073ffffffffffffffffffffffffffffffffffffffff60025416604051908152f35b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc60209101126100155760043573ffffffffffffffffffffffffffffffffffffffff811681036100155790565b346100155773ffffffffffffffffffffffffffffffffffffffff61025e366101ee565b6102666105d9565b1680156102c4576020817fef31f568d741a833c6a9dc85a6e1c65e06fa772740d5dc94d1da21827a4e0cab927fffffffffffffffffffffffff00000000000000000000000000000000000000006002541617600255604051908152a1005b7f8579befe0000000000000000000000000000000000000000000000000000000060005260046000fd5b346100155760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100155773ffffffffffffffffffffffffffffffffffffffff6001541633036103eb5760005473ffffffffffffffffffffffffffffffffffffffff16600080547fffffffffffffffffffffffff000000000000000000000000000000000000000016331790556103ac7fffffffffffffffffffffffff000000000000000000000000000000000000000060015416600155565b73ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152fd5b346100155760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261001557602073ffffffffffffffffffffffffffffffffffffffff60005416604051908152f35b346100155773ffffffffffffffffffffffffffffffffffffffff6104be366101ee565b6104c66105d9565b1633811461053a57807fffffffffffffffffffffffff0000000000000000000000000000000000000000600154161760015573ffffffffffffffffffffffffffffffffffffffff8060005416167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152fd5b73ffffffffffffffffffffffffffffffffffffffff60025416803b15610015576000809136828037818036925af13d6000803e6105d4573d6000fd5b3d6000f35b73ffffffffffffffffffffffffffffffffffffffff6000541633036105fa57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152fdfea164736f6c634300081a000a",
}

var RMNProxyContractABI = RMNProxyContractMetaData.ABI

var RMNProxyContractBin = RMNProxyContractMetaData.Bin

func DeployRMNProxyContract(auth *bind.TransactOpts, backend bind.ContractBackend, arm common.Address) (common.Address, *types.Transaction, *RMNProxyContract, error) {
	parsed, err := RMNProxyContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RMNProxyContractBin), backend, arm)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RMNProxyContract{address: address, abi: *parsed, RMNProxyContractCaller: RMNProxyContractCaller{contract: contract}, RMNProxyContractTransactor: RMNProxyContractTransactor{contract: contract}, RMNProxyContractFilterer: RMNProxyContractFilterer{contract: contract}}, nil
}

type RMNProxyContract struct {
	address common.Address
	abi     abi.ABI
	RMNProxyContractCaller
	RMNProxyContractTransactor
	RMNProxyContractFilterer
}

type RMNProxyContractCaller struct {
	contract *bind.BoundContract
}

type RMNProxyContractTransactor struct {
	contract *bind.BoundContract
}

type RMNProxyContractFilterer struct {
	contract *bind.BoundContract
}

type RMNProxyContractSession struct {
	Contract     *RMNProxyContract
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type RMNProxyContractCallerSession struct {
	Contract *RMNProxyContractCaller
	CallOpts bind.CallOpts
}

type RMNProxyContractTransactorSession struct {
	Contract     *RMNProxyContractTransactor
	TransactOpts bind.TransactOpts
}

type RMNProxyContractRaw struct {
	Contract *RMNProxyContract
}

type RMNProxyContractCallerRaw struct {
	Contract *RMNProxyContractCaller
}

type RMNProxyContractTransactorRaw struct {
	Contract *RMNProxyContractTransactor
}

func NewRMNProxyContract(address common.Address, backend bind.ContractBackend) (*RMNProxyContract, error) {
	abi, err := abi.JSON(strings.NewReader(RMNProxyContractABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindRMNProxyContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RMNProxyContract{address: address, abi: abi, RMNProxyContractCaller: RMNProxyContractCaller{contract: contract}, RMNProxyContractTransactor: RMNProxyContractTransactor{contract: contract}, RMNProxyContractFilterer: RMNProxyContractFilterer{contract: contract}}, nil
}

func NewRMNProxyContractCaller(address common.Address, caller bind.ContractCaller) (*RMNProxyContractCaller, error) {
	contract, err := bindRMNProxyContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RMNProxyContractCaller{contract: contract}, nil
}

func NewRMNProxyContractTransactor(address common.Address, transactor bind.ContractTransactor) (*RMNProxyContractTransactor, error) {
	contract, err := bindRMNProxyContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RMNProxyContractTransactor{contract: contract}, nil
}

func NewRMNProxyContractFilterer(address common.Address, filterer bind.ContractFilterer) (*RMNProxyContractFilterer, error) {
	contract, err := bindRMNProxyContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RMNProxyContractFilterer{contract: contract}, nil
}

func bindRMNProxyContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RMNProxyContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_RMNProxyContract *RMNProxyContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNProxyContract.Contract.RMNProxyContractCaller.contract.Call(opts, result, method, params...)
}

func (_RMNProxyContract *RMNProxyContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNProxyContract.Contract.RMNProxyContractTransactor.contract.Transfer(opts)
}

func (_RMNProxyContract *RMNProxyContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNProxyContract.Contract.RMNProxyContractTransactor.contract.Transact(opts, method, params...)
}

func (_RMNProxyContract *RMNProxyContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNProxyContract.Contract.contract.Call(opts, result, method, params...)
}

func (_RMNProxyContract *RMNProxyContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNProxyContract.Contract.contract.Transfer(opts)
}

func (_RMNProxyContract *RMNProxyContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNProxyContract.Contract.contract.Transact(opts, method, params...)
}

func (_RMNProxyContract *RMNProxyContractCaller) GetARM(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RMNProxyContract.contract.Call(opts, &out, "getARM")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_RMNProxyContract *RMNProxyContractSession) GetARM() (common.Address, error) {
	return _RMNProxyContract.Contract.GetARM(&_RMNProxyContract.CallOpts)
}

func (_RMNProxyContract *RMNProxyContractCallerSession) GetARM() (common.Address, error) {
	return _RMNProxyContract.Contract.GetARM(&_RMNProxyContract.CallOpts)
}

func (_RMNProxyContract *RMNProxyContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RMNProxyContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_RMNProxyContract *RMNProxyContractSession) Owner() (common.Address, error) {
	return _RMNProxyContract.Contract.Owner(&_RMNProxyContract.CallOpts)
}

func (_RMNProxyContract *RMNProxyContractCallerSession) Owner() (common.Address, error) {
	return _RMNProxyContract.Contract.Owner(&_RMNProxyContract.CallOpts)
}

func (_RMNProxyContract *RMNProxyContractCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RMNProxyContract.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_RMNProxyContract *RMNProxyContractSession) TypeAndVersion() (string, error) {
	return _RMNProxyContract.Contract.TypeAndVersion(&_RMNProxyContract.CallOpts)
}

func (_RMNProxyContract *RMNProxyContractCallerSession) TypeAndVersion() (string, error) {
	return _RMNProxyContract.Contract.TypeAndVersion(&_RMNProxyContract.CallOpts)
}

func (_RMNProxyContract *RMNProxyContractTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNProxyContract.contract.Transact(opts, "acceptOwnership")
}

func (_RMNProxyContract *RMNProxyContractSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNProxyContract.Contract.AcceptOwnership(&_RMNProxyContract.TransactOpts)
}

func (_RMNProxyContract *RMNProxyContractTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNProxyContract.Contract.AcceptOwnership(&_RMNProxyContract.TransactOpts)
}

func (_RMNProxyContract *RMNProxyContractTransactor) SetARM(opts *bind.TransactOpts, arm common.Address) (*types.Transaction, error) {
	return _RMNProxyContract.contract.Transact(opts, "setARM", arm)
}

func (_RMNProxyContract *RMNProxyContractSession) SetARM(arm common.Address) (*types.Transaction, error) {
	return _RMNProxyContract.Contract.SetARM(&_RMNProxyContract.TransactOpts, arm)
}

func (_RMNProxyContract *RMNProxyContractTransactorSession) SetARM(arm common.Address) (*types.Transaction, error) {
	return _RMNProxyContract.Contract.SetARM(&_RMNProxyContract.TransactOpts, arm)
}

func (_RMNProxyContract *RMNProxyContractTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _RMNProxyContract.contract.Transact(opts, "transferOwnership", to)
}

func (_RMNProxyContract *RMNProxyContractSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNProxyContract.Contract.TransferOwnership(&_RMNProxyContract.TransactOpts, to)
}

func (_RMNProxyContract *RMNProxyContractTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNProxyContract.Contract.TransferOwnership(&_RMNProxyContract.TransactOpts, to)
}

func (_RMNProxyContract *RMNProxyContractTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _RMNProxyContract.contract.RawTransact(opts, calldata)
}

func (_RMNProxyContract *RMNProxyContractSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _RMNProxyContract.Contract.Fallback(&_RMNProxyContract.TransactOpts, calldata)
}

func (_RMNProxyContract *RMNProxyContractTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _RMNProxyContract.Contract.Fallback(&_RMNProxyContract.TransactOpts, calldata)
}

type RMNProxyContractARMSetIterator struct {
	Event *RMNProxyContractARMSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNProxyContractARMSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNProxyContractARMSet)
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
		it.Event = new(RMNProxyContractARMSet)
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

func (it *RMNProxyContractARMSetIterator) Error() error {
	return it.fail
}

func (it *RMNProxyContractARMSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNProxyContractARMSet struct {
	Arm common.Address
	Raw types.Log
}

func (_RMNProxyContract *RMNProxyContractFilterer) FilterARMSet(opts *bind.FilterOpts) (*RMNProxyContractARMSetIterator, error) {

	logs, sub, err := _RMNProxyContract.contract.FilterLogs(opts, "ARMSet")
	if err != nil {
		return nil, err
	}
	return &RMNProxyContractARMSetIterator{contract: _RMNProxyContract.contract, event: "ARMSet", logs: logs, sub: sub}, nil
}

func (_RMNProxyContract *RMNProxyContractFilterer) WatchARMSet(opts *bind.WatchOpts, sink chan<- *RMNProxyContractARMSet) (event.Subscription, error) {

	logs, sub, err := _RMNProxyContract.contract.WatchLogs(opts, "ARMSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNProxyContractARMSet)
				if err := _RMNProxyContract.contract.UnpackLog(event, "ARMSet", log); err != nil {
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

func (_RMNProxyContract *RMNProxyContractFilterer) ParseARMSet(log types.Log) (*RMNProxyContractARMSet, error) {
	event := new(RMNProxyContractARMSet)
	if err := _RMNProxyContract.contract.UnpackLog(event, "ARMSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNProxyContractOwnershipTransferRequestedIterator struct {
	Event *RMNProxyContractOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNProxyContractOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNProxyContractOwnershipTransferRequested)
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
		it.Event = new(RMNProxyContractOwnershipTransferRequested)
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

func (it *RMNProxyContractOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *RMNProxyContractOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNProxyContractOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNProxyContract *RMNProxyContractFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNProxyContractOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNProxyContract.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNProxyContractOwnershipTransferRequestedIterator{contract: _RMNProxyContract.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_RMNProxyContract *RMNProxyContractFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNProxyContractOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNProxyContract.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNProxyContractOwnershipTransferRequested)
				if err := _RMNProxyContract.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_RMNProxyContract *RMNProxyContractFilterer) ParseOwnershipTransferRequested(log types.Log) (*RMNProxyContractOwnershipTransferRequested, error) {
	event := new(RMNProxyContractOwnershipTransferRequested)
	if err := _RMNProxyContract.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNProxyContractOwnershipTransferredIterator struct {
	Event *RMNProxyContractOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNProxyContractOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNProxyContractOwnershipTransferred)
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
		it.Event = new(RMNProxyContractOwnershipTransferred)
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

func (it *RMNProxyContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *RMNProxyContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNProxyContractOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNProxyContract *RMNProxyContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNProxyContractOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNProxyContract.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNProxyContractOwnershipTransferredIterator{contract: _RMNProxyContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_RMNProxyContract *RMNProxyContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNProxyContractOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNProxyContract.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNProxyContractOwnershipTransferred)
				if err := _RMNProxyContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_RMNProxyContract *RMNProxyContractFilterer) ParseOwnershipTransferred(log types.Log) (*RMNProxyContractOwnershipTransferred, error) {
	event := new(RMNProxyContractOwnershipTransferred)
	if err := _RMNProxyContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_RMNProxyContract *RMNProxyContract) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _RMNProxyContract.abi.Events["ARMSet"].ID:
		return _RMNProxyContract.ParseARMSet(log)
	case _RMNProxyContract.abi.Events["OwnershipTransferRequested"].ID:
		return _RMNProxyContract.ParseOwnershipTransferRequested(log)
	case _RMNProxyContract.abi.Events["OwnershipTransferred"].ID:
		return _RMNProxyContract.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (RMNProxyContractARMSet) Topic() common.Hash {
	return common.HexToHash("0xef31f568d741a833c6a9dc85a6e1c65e06fa772740d5dc94d1da21827a4e0cab")
}

func (RMNProxyContractOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (RMNProxyContractOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_RMNProxyContract *RMNProxyContract) Address() common.Address {
	return _RMNProxyContract.address
}

type RMNProxyContractInterface interface {
	GetARM(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetARM(opts *bind.TransactOpts, arm common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterARMSet(opts *bind.FilterOpts) (*RMNProxyContractARMSetIterator, error)

	WatchARMSet(opts *bind.WatchOpts, sink chan<- *RMNProxyContractARMSet) (event.Subscription, error)

	ParseARMSet(log types.Log) (*RMNProxyContractARMSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNProxyContractOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNProxyContractOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*RMNProxyContractOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNProxyContractOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNProxyContractOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*RMNProxyContractOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
