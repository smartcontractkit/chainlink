// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arm_proxy_contract

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

var ARMProxyContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"arm\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"arm\",\"type\":\"address\"}],\"name\":\"ARMSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getARM\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"arm\",\"type\":\"address\"}],\"name\":\"setARM\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161084038038061084083398101604081905261002f91610255565b33806000816100855760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100b5576100b5816100cd565b5050506100c78161017660201b60201c565b50610285565b336001600160a01b038216036101255760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161007c565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61017e6101f9565b6001600160a01b0381166101a5576040516342bcdf7f60e11b815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0383169081179091556040519081527fef31f568d741a833c6a9dc85a6e1c65e06fa772740d5dc94d1da21827a4e0cab9060200160405180910390a150565b6000546001600160a01b031633146102535760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161007c565b565b60006020828403121561026757600080fd5b81516001600160a01b038116811461027e57600080fd5b9392505050565b6105ac806102946000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c806379ba50971161005057806379ba5097146101615780638da5cb5b14610169578063f2fde38b1461018757610072565b8063181f5a77146100bb5780632e90aa211461010d578063458fec3b1461014c575b60025473ffffffffffffffffffffffffffffffffffffffff16803b61009657600080fd5b366000803760008036600080855af13d6000803e80156100b5573d6000f35b503d6000fd5b6100f76040518060400160405280600e81526020017f41524d50726f787920312e302e3000000000000000000000000000000000000081525081565b60405161010491906104f6565b60405180910390f35b60025473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610104565b61015f61015a366004610562565b61019a565b005b61015f610268565b60005473ffffffffffffffffffffffffffffffffffffffff16610127565b61015f610195366004610562565b61036a565b6101a261037e565b73ffffffffffffffffffffffffffffffffffffffff81166101ef576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527fef31f568d741a833c6a9dc85a6e1c65e06fa772740d5dc94d1da21827a4e0cab9060200160405180910390a150565b60015473ffffffffffffffffffffffffffffffffffffffff1633146102ee576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61037261037e565b61037b81610401565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1633146103ff576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102e5565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610480576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102e5565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208083528351808285015260005b8181101561052357858101830151858201604001528201610507565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60006020828403121561057457600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461059857600080fd5b939250505056fea164736f6c6343000813000a",
}

var ARMProxyContractABI = ARMProxyContractMetaData.ABI

var ARMProxyContractBin = ARMProxyContractMetaData.Bin

func DeployARMProxyContract(auth *bind.TransactOpts, backend bind.ContractBackend, arm common.Address) (common.Address, *types.Transaction, *ARMProxyContract, error) {
	parsed, err := ARMProxyContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ARMProxyContractBin), backend, arm)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ARMProxyContract{address: address, abi: *parsed, ARMProxyContractCaller: ARMProxyContractCaller{contract: contract}, ARMProxyContractTransactor: ARMProxyContractTransactor{contract: contract}, ARMProxyContractFilterer: ARMProxyContractFilterer{contract: contract}}, nil
}

type ARMProxyContract struct {
	address common.Address
	abi     abi.ABI
	ARMProxyContractCaller
	ARMProxyContractTransactor
	ARMProxyContractFilterer
}

type ARMProxyContractCaller struct {
	contract *bind.BoundContract
}

type ARMProxyContractTransactor struct {
	contract *bind.BoundContract
}

type ARMProxyContractFilterer struct {
	contract *bind.BoundContract
}

type ARMProxyContractSession struct {
	Contract     *ARMProxyContract
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ARMProxyContractCallerSession struct {
	Contract *ARMProxyContractCaller
	CallOpts bind.CallOpts
}

type ARMProxyContractTransactorSession struct {
	Contract     *ARMProxyContractTransactor
	TransactOpts bind.TransactOpts
}

type ARMProxyContractRaw struct {
	Contract *ARMProxyContract
}

type ARMProxyContractCallerRaw struct {
	Contract *ARMProxyContractCaller
}

type ARMProxyContractTransactorRaw struct {
	Contract *ARMProxyContractTransactor
}

func NewARMProxyContract(address common.Address, backend bind.ContractBackend) (*ARMProxyContract, error) {
	abi, err := abi.JSON(strings.NewReader(ARMProxyContractABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindARMProxyContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ARMProxyContract{address: address, abi: abi, ARMProxyContractCaller: ARMProxyContractCaller{contract: contract}, ARMProxyContractTransactor: ARMProxyContractTransactor{contract: contract}, ARMProxyContractFilterer: ARMProxyContractFilterer{contract: contract}}, nil
}

func NewARMProxyContractCaller(address common.Address, caller bind.ContractCaller) (*ARMProxyContractCaller, error) {
	contract, err := bindARMProxyContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ARMProxyContractCaller{contract: contract}, nil
}

func NewARMProxyContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ARMProxyContractTransactor, error) {
	contract, err := bindARMProxyContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ARMProxyContractTransactor{contract: contract}, nil
}

func NewARMProxyContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ARMProxyContractFilterer, error) {
	contract, err := bindARMProxyContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ARMProxyContractFilterer{contract: contract}, nil
}

func bindARMProxyContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ARMProxyContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ARMProxyContract *ARMProxyContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ARMProxyContract.Contract.ARMProxyContractCaller.contract.Call(opts, result, method, params...)
}

func (_ARMProxyContract *ARMProxyContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ARMProxyContract.Contract.ARMProxyContractTransactor.contract.Transfer(opts)
}

func (_ARMProxyContract *ARMProxyContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ARMProxyContract.Contract.ARMProxyContractTransactor.contract.Transact(opts, method, params...)
}

func (_ARMProxyContract *ARMProxyContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ARMProxyContract.Contract.contract.Call(opts, result, method, params...)
}

func (_ARMProxyContract *ARMProxyContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ARMProxyContract.Contract.contract.Transfer(opts)
}

func (_ARMProxyContract *ARMProxyContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ARMProxyContract.Contract.contract.Transact(opts, method, params...)
}

func (_ARMProxyContract *ARMProxyContractCaller) GetARM(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ARMProxyContract.contract.Call(opts, &out, "getARM")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ARMProxyContract *ARMProxyContractSession) GetARM() (common.Address, error) {
	return _ARMProxyContract.Contract.GetARM(&_ARMProxyContract.CallOpts)
}

func (_ARMProxyContract *ARMProxyContractCallerSession) GetARM() (common.Address, error) {
	return _ARMProxyContract.Contract.GetARM(&_ARMProxyContract.CallOpts)
}

func (_ARMProxyContract *ARMProxyContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ARMProxyContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ARMProxyContract *ARMProxyContractSession) Owner() (common.Address, error) {
	return _ARMProxyContract.Contract.Owner(&_ARMProxyContract.CallOpts)
}

func (_ARMProxyContract *ARMProxyContractCallerSession) Owner() (common.Address, error) {
	return _ARMProxyContract.Contract.Owner(&_ARMProxyContract.CallOpts)
}

func (_ARMProxyContract *ARMProxyContractCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ARMProxyContract.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_ARMProxyContract *ARMProxyContractSession) TypeAndVersion() (string, error) {
	return _ARMProxyContract.Contract.TypeAndVersion(&_ARMProxyContract.CallOpts)
}

func (_ARMProxyContract *ARMProxyContractCallerSession) TypeAndVersion() (string, error) {
	return _ARMProxyContract.Contract.TypeAndVersion(&_ARMProxyContract.CallOpts)
}

func (_ARMProxyContract *ARMProxyContractTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ARMProxyContract.contract.Transact(opts, "acceptOwnership")
}

func (_ARMProxyContract *ARMProxyContractSession) AcceptOwnership() (*types.Transaction, error) {
	return _ARMProxyContract.Contract.AcceptOwnership(&_ARMProxyContract.TransactOpts)
}

func (_ARMProxyContract *ARMProxyContractTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ARMProxyContract.Contract.AcceptOwnership(&_ARMProxyContract.TransactOpts)
}

func (_ARMProxyContract *ARMProxyContractTransactor) SetARM(opts *bind.TransactOpts, arm common.Address) (*types.Transaction, error) {
	return _ARMProxyContract.contract.Transact(opts, "setARM", arm)
}

func (_ARMProxyContract *ARMProxyContractSession) SetARM(arm common.Address) (*types.Transaction, error) {
	return _ARMProxyContract.Contract.SetARM(&_ARMProxyContract.TransactOpts, arm)
}

func (_ARMProxyContract *ARMProxyContractTransactorSession) SetARM(arm common.Address) (*types.Transaction, error) {
	return _ARMProxyContract.Contract.SetARM(&_ARMProxyContract.TransactOpts, arm)
}

func (_ARMProxyContract *ARMProxyContractTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ARMProxyContract.contract.Transact(opts, "transferOwnership", to)
}

func (_ARMProxyContract *ARMProxyContractSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ARMProxyContract.Contract.TransferOwnership(&_ARMProxyContract.TransactOpts, to)
}

func (_ARMProxyContract *ARMProxyContractTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ARMProxyContract.Contract.TransferOwnership(&_ARMProxyContract.TransactOpts, to)
}

func (_ARMProxyContract *ARMProxyContractTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _ARMProxyContract.contract.RawTransact(opts, calldata)
}

func (_ARMProxyContract *ARMProxyContractSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ARMProxyContract.Contract.Fallback(&_ARMProxyContract.TransactOpts, calldata)
}

func (_ARMProxyContract *ARMProxyContractTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ARMProxyContract.Contract.Fallback(&_ARMProxyContract.TransactOpts, calldata)
}

type ARMProxyContractARMSetIterator struct {
	Event *ARMProxyContractARMSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ARMProxyContractARMSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ARMProxyContractARMSet)
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
		it.Event = new(ARMProxyContractARMSet)
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

func (it *ARMProxyContractARMSetIterator) Error() error {
	return it.fail
}

func (it *ARMProxyContractARMSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ARMProxyContractARMSet struct {
	Arm common.Address
	Raw types.Log
}

func (_ARMProxyContract *ARMProxyContractFilterer) FilterARMSet(opts *bind.FilterOpts) (*ARMProxyContractARMSetIterator, error) {

	logs, sub, err := _ARMProxyContract.contract.FilterLogs(opts, "ARMSet")
	if err != nil {
		return nil, err
	}
	return &ARMProxyContractARMSetIterator{contract: _ARMProxyContract.contract, event: "ARMSet", logs: logs, sub: sub}, nil
}

func (_ARMProxyContract *ARMProxyContractFilterer) WatchARMSet(opts *bind.WatchOpts, sink chan<- *ARMProxyContractARMSet) (event.Subscription, error) {

	logs, sub, err := _ARMProxyContract.contract.WatchLogs(opts, "ARMSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ARMProxyContractARMSet)
				if err := _ARMProxyContract.contract.UnpackLog(event, "ARMSet", log); err != nil {
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

func (_ARMProxyContract *ARMProxyContractFilterer) ParseARMSet(log types.Log) (*ARMProxyContractARMSet, error) {
	event := new(ARMProxyContractARMSet)
	if err := _ARMProxyContract.contract.UnpackLog(event, "ARMSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ARMProxyContractOwnershipTransferRequestedIterator struct {
	Event *ARMProxyContractOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ARMProxyContractOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ARMProxyContractOwnershipTransferRequested)
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
		it.Event = new(ARMProxyContractOwnershipTransferRequested)
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

func (it *ARMProxyContractOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *ARMProxyContractOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ARMProxyContractOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ARMProxyContract *ARMProxyContractFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ARMProxyContractOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ARMProxyContract.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ARMProxyContractOwnershipTransferRequestedIterator{contract: _ARMProxyContract.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_ARMProxyContract *ARMProxyContractFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ARMProxyContractOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ARMProxyContract.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ARMProxyContractOwnershipTransferRequested)
				if err := _ARMProxyContract.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_ARMProxyContract *ARMProxyContractFilterer) ParseOwnershipTransferRequested(log types.Log) (*ARMProxyContractOwnershipTransferRequested, error) {
	event := new(ARMProxyContractOwnershipTransferRequested)
	if err := _ARMProxyContract.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ARMProxyContractOwnershipTransferredIterator struct {
	Event *ARMProxyContractOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ARMProxyContractOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ARMProxyContractOwnershipTransferred)
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
		it.Event = new(ARMProxyContractOwnershipTransferred)
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

func (it *ARMProxyContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *ARMProxyContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ARMProxyContractOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ARMProxyContract *ARMProxyContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ARMProxyContractOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ARMProxyContract.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ARMProxyContractOwnershipTransferredIterator{contract: _ARMProxyContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_ARMProxyContract *ARMProxyContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ARMProxyContractOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ARMProxyContract.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ARMProxyContractOwnershipTransferred)
				if err := _ARMProxyContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_ARMProxyContract *ARMProxyContractFilterer) ParseOwnershipTransferred(log types.Log) (*ARMProxyContractOwnershipTransferred, error) {
	event := new(ARMProxyContractOwnershipTransferred)
	if err := _ARMProxyContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_ARMProxyContract *ARMProxyContract) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ARMProxyContract.abi.Events["ARMSet"].ID:
		return _ARMProxyContract.ParseARMSet(log)
	case _ARMProxyContract.abi.Events["OwnershipTransferRequested"].ID:
		return _ARMProxyContract.ParseOwnershipTransferRequested(log)
	case _ARMProxyContract.abi.Events["OwnershipTransferred"].ID:
		return _ARMProxyContract.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ARMProxyContractARMSet) Topic() common.Hash {
	return common.HexToHash("0xef31f568d741a833c6a9dc85a6e1c65e06fa772740d5dc94d1da21827a4e0cab")
}

func (ARMProxyContractOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (ARMProxyContractOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_ARMProxyContract *ARMProxyContract) Address() common.Address {
	return _ARMProxyContract.address
}

type ARMProxyContractInterface interface {
	GetARM(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetARM(opts *bind.TransactOpts, arm common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterARMSet(opts *bind.FilterOpts) (*ARMProxyContractARMSetIterator, error)

	WatchARMSet(opts *bind.WatchOpts, sink chan<- *ARMProxyContractARMSet) (event.Subscription, error)

	ParseARMSet(log types.Log) (*ARMProxyContractARMSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ARMProxyContractOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ARMProxyContractOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*ARMProxyContractOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ARMProxyContractOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ARMProxyContractOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*ARMProxyContractOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
