// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr2titlerequest

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
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

// ConfirmedOwnerMetaData contains all meta data concerning the ConfirmedOwner contract.
var ConfirmedOwnerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161051638038061051683398101604081905261002f9161016f565b8060006001600160a01b03821661008d5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100bd576100bd816100c5565b50505061019f565b6001600160a01b03811633141561011e5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610084565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561018157600080fd5b81516001600160a01b038116811461019857600080fd5b9392505050565b610368806101ae6000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806379ba5097146100465780638da5cb5b14610050578063f2fde38b1461007c575b600080fd5b61004e61008f565b005b6000546040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61004e61008a36600461031e565b610191565b60015473ffffffffffffffffffffffffffffffffffffffff163314610115576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6101996101a5565b6101a281610228565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610226576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161010c565b565b73ffffffffffffffffffffffffffffffffffffffff81163314156102a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161010c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561033057600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461035457600080fd5b939250505056fea164736f6c6343000806000a",
}

// ConfirmedOwnerABI is the input ABI used to generate the binding from.
// Deprecated: Use ConfirmedOwnerMetaData.ABI instead.
var ConfirmedOwnerABI = ConfirmedOwnerMetaData.ABI

// ConfirmedOwnerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConfirmedOwnerMetaData.Bin instead.
var ConfirmedOwnerBin = ConfirmedOwnerMetaData.Bin

// DeployConfirmedOwner deploys a new Ethereum contract, binding an instance of ConfirmedOwner to it.
func DeployConfirmedOwner(auth *bind.TransactOpts, backend bind.ContractBackend, newOwner common.Address) (common.Address, *types.Transaction, *ConfirmedOwner, error) {
	parsed, err := ConfirmedOwnerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConfirmedOwnerBin), backend, newOwner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConfirmedOwner{ConfirmedOwnerCaller: ConfirmedOwnerCaller{contract: contract}, ConfirmedOwnerTransactor: ConfirmedOwnerTransactor{contract: contract}, ConfirmedOwnerFilterer: ConfirmedOwnerFilterer{contract: contract}}, nil
}

// ConfirmedOwner is an auto generated Go binding around an Ethereum contract.
type ConfirmedOwner struct {
	ConfirmedOwnerCaller     // Read-only binding to the contract
	ConfirmedOwnerTransactor // Write-only binding to the contract
	ConfirmedOwnerFilterer   // Log filterer for contract events
}

// ConfirmedOwnerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConfirmedOwnerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConfirmedOwnerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConfirmedOwnerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConfirmedOwnerSession struct {
	Contract     *ConfirmedOwner   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConfirmedOwnerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConfirmedOwnerCallerSession struct {
	Contract *ConfirmedOwnerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ConfirmedOwnerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConfirmedOwnerTransactorSession struct {
	Contract     *ConfirmedOwnerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ConfirmedOwnerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConfirmedOwnerRaw struct {
	Contract *ConfirmedOwner // Generic contract binding to access the raw methods on
}

// ConfirmedOwnerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConfirmedOwnerCallerRaw struct {
	Contract *ConfirmedOwnerCaller // Generic read-only contract binding to access the raw methods on
}

// ConfirmedOwnerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConfirmedOwnerTransactorRaw struct {
	Contract *ConfirmedOwnerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConfirmedOwner creates a new instance of ConfirmedOwner, bound to a specific deployed contract.
func NewConfirmedOwner(address common.Address, backend bind.ContractBackend) (*ConfirmedOwner, error) {
	contract, err := bindConfirmedOwner(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwner{ConfirmedOwnerCaller: ConfirmedOwnerCaller{contract: contract}, ConfirmedOwnerTransactor: ConfirmedOwnerTransactor{contract: contract}, ConfirmedOwnerFilterer: ConfirmedOwnerFilterer{contract: contract}}, nil
}

// NewConfirmedOwnerCaller creates a new read-only instance of ConfirmedOwner, bound to a specific deployed contract.
func NewConfirmedOwnerCaller(address common.Address, caller bind.ContractCaller) (*ConfirmedOwnerCaller, error) {
	contract, err := bindConfirmedOwner(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerCaller{contract: contract}, nil
}

// NewConfirmedOwnerTransactor creates a new write-only instance of ConfirmedOwner, bound to a specific deployed contract.
func NewConfirmedOwnerTransactor(address common.Address, transactor bind.ContractTransactor) (*ConfirmedOwnerTransactor, error) {
	contract, err := bindConfirmedOwner(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerTransactor{contract: contract}, nil
}

// NewConfirmedOwnerFilterer creates a new log filterer instance of ConfirmedOwner, bound to a specific deployed contract.
func NewConfirmedOwnerFilterer(address common.Address, filterer bind.ContractFilterer) (*ConfirmedOwnerFilterer, error) {
	contract, err := bindConfirmedOwner(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerFilterer{contract: contract}, nil
}

// bindConfirmedOwner binds a generic wrapper to an already deployed contract.
func bindConfirmedOwner(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConfirmedOwnerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfirmedOwner *ConfirmedOwnerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfirmedOwner.Contract.ConfirmedOwnerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfirmedOwner *ConfirmedOwnerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.ConfirmedOwnerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfirmedOwner *ConfirmedOwnerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.ConfirmedOwnerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfirmedOwner *ConfirmedOwnerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfirmedOwner.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfirmedOwner *ConfirmedOwnerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfirmedOwner *ConfirmedOwnerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwner *ConfirmedOwnerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConfirmedOwner.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwner *ConfirmedOwnerSession) Owner() (common.Address, error) {
	return _ConfirmedOwner.Contract.Owner(&_ConfirmedOwner.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwner *ConfirmedOwnerCallerSession) Owner() (common.Address, error) {
	return _ConfirmedOwner.Contract.Owner(&_ConfirmedOwner.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwner *ConfirmedOwnerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwner.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwner *ConfirmedOwnerSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.AcceptOwnership(&_ConfirmedOwner.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwner *ConfirmedOwnerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.AcceptOwnership(&_ConfirmedOwner.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwner *ConfirmedOwnerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwner.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwner *ConfirmedOwnerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.TransferOwnership(&_ConfirmedOwner.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwner *ConfirmedOwnerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.TransferOwnership(&_ConfirmedOwner.TransactOpts, to)
}

// ConfirmedOwnerOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the ConfirmedOwner contract.
type ConfirmedOwnerOwnershipTransferRequestedIterator struct {
	Event *ConfirmedOwnerOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConfirmedOwnerOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfirmedOwnerOwnershipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConfirmedOwnerOwnershipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConfirmedOwnerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfirmedOwnerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfirmedOwnerOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the ConfirmedOwner contract.
type ConfirmedOwnerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfirmedOwnerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwner.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerOwnershipTransferRequestedIterator{contract: _ConfirmedOwner.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ConfirmedOwnerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwner.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfirmedOwnerOwnershipTransferRequested)
				if err := _ConfirmedOwner.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) ParseOwnershipTransferRequested(log types.Log) (*ConfirmedOwnerOwnershipTransferRequested, error) {
	event := new(ConfirmedOwnerOwnershipTransferRequested)
	if err := _ConfirmedOwner.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfirmedOwnerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ConfirmedOwner contract.
type ConfirmedOwnerOwnershipTransferredIterator struct {
	Event *ConfirmedOwnerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConfirmedOwnerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfirmedOwnerOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConfirmedOwnerOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConfirmedOwnerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfirmedOwnerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfirmedOwnerOwnershipTransferred represents a OwnershipTransferred event raised by the ConfirmedOwner contract.
type ConfirmedOwnerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfirmedOwnerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwner.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerOwnershipTransferredIterator{contract: _ConfirmedOwner.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConfirmedOwnerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwner.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfirmedOwnerOwnershipTransferred)
				if err := _ConfirmedOwner.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) ParseOwnershipTransferred(log types.Log) (*ConfirmedOwnerOwnershipTransferred, error) {
	event := new(ConfirmedOwnerOwnershipTransferred)
	if err := _ConfirmedOwner.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfirmedOwnerWithProposalMetaData contains all meta data concerning the ConfirmedOwnerWithProposal contract.
var ConfirmedOwnerWithProposalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pendingOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161053138038061053183398101604081905261002f91610187565b6001600160a01b03821661008a5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100ba576100ba816100c1565b50506101ba565b6001600160a01b03811633141561011a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610081565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b038116811461018257600080fd5b919050565b6000806040838503121561019a57600080fd5b6101a38361016b565b91506101b16020840161016b565b90509250929050565b610368806101c96000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806379ba5097146100465780638da5cb5b14610050578063f2fde38b1461007c575b600080fd5b61004e61008f565b005b6000546040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61004e61008a36600461031e565b610191565b60015473ffffffffffffffffffffffffffffffffffffffff163314610115576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6101996101a5565b6101a281610228565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610226576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161010c565b565b73ffffffffffffffffffffffffffffffffffffffff81163314156102a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161010c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561033057600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461035457600080fd5b939250505056fea164736f6c6343000806000a",
}

// ConfirmedOwnerWithProposalABI is the input ABI used to generate the binding from.
// Deprecated: Use ConfirmedOwnerWithProposalMetaData.ABI instead.
var ConfirmedOwnerWithProposalABI = ConfirmedOwnerWithProposalMetaData.ABI

// ConfirmedOwnerWithProposalBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConfirmedOwnerWithProposalMetaData.Bin instead.
var ConfirmedOwnerWithProposalBin = ConfirmedOwnerWithProposalMetaData.Bin

// DeployConfirmedOwnerWithProposal deploys a new Ethereum contract, binding an instance of ConfirmedOwnerWithProposal to it.
func DeployConfirmedOwnerWithProposal(auth *bind.TransactOpts, backend bind.ContractBackend, newOwner common.Address, pendingOwner common.Address) (common.Address, *types.Transaction, *ConfirmedOwnerWithProposal, error) {
	parsed, err := ConfirmedOwnerWithProposalMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConfirmedOwnerWithProposalBin), backend, newOwner, pendingOwner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConfirmedOwnerWithProposal{ConfirmedOwnerWithProposalCaller: ConfirmedOwnerWithProposalCaller{contract: contract}, ConfirmedOwnerWithProposalTransactor: ConfirmedOwnerWithProposalTransactor{contract: contract}, ConfirmedOwnerWithProposalFilterer: ConfirmedOwnerWithProposalFilterer{contract: contract}}, nil
}

// ConfirmedOwnerWithProposal is an auto generated Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposal struct {
	ConfirmedOwnerWithProposalCaller     // Read-only binding to the contract
	ConfirmedOwnerWithProposalTransactor // Write-only binding to the contract
	ConfirmedOwnerWithProposalFilterer   // Log filterer for contract events
}

// ConfirmedOwnerWithProposalCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerWithProposalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerWithProposalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConfirmedOwnerWithProposalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerWithProposalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConfirmedOwnerWithProposalSession struct {
	Contract     *ConfirmedOwnerWithProposal // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ConfirmedOwnerWithProposalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConfirmedOwnerWithProposalCallerSession struct {
	Contract *ConfirmedOwnerWithProposalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// ConfirmedOwnerWithProposalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConfirmedOwnerWithProposalTransactorSession struct {
	Contract     *ConfirmedOwnerWithProposalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// ConfirmedOwnerWithProposalRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalRaw struct {
	Contract *ConfirmedOwnerWithProposal // Generic contract binding to access the raw methods on
}

// ConfirmedOwnerWithProposalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalCallerRaw struct {
	Contract *ConfirmedOwnerWithProposalCaller // Generic read-only contract binding to access the raw methods on
}

// ConfirmedOwnerWithProposalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalTransactorRaw struct {
	Contract *ConfirmedOwnerWithProposalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConfirmedOwnerWithProposal creates a new instance of ConfirmedOwnerWithProposal, bound to a specific deployed contract.
func NewConfirmedOwnerWithProposal(address common.Address, backend bind.ContractBackend) (*ConfirmedOwnerWithProposal, error) {
	contract, err := bindConfirmedOwnerWithProposal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposal{ConfirmedOwnerWithProposalCaller: ConfirmedOwnerWithProposalCaller{contract: contract}, ConfirmedOwnerWithProposalTransactor: ConfirmedOwnerWithProposalTransactor{contract: contract}, ConfirmedOwnerWithProposalFilterer: ConfirmedOwnerWithProposalFilterer{contract: contract}}, nil
}

// NewConfirmedOwnerWithProposalCaller creates a new read-only instance of ConfirmedOwnerWithProposal, bound to a specific deployed contract.
func NewConfirmedOwnerWithProposalCaller(address common.Address, caller bind.ContractCaller) (*ConfirmedOwnerWithProposalCaller, error) {
	contract, err := bindConfirmedOwnerWithProposal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalCaller{contract: contract}, nil
}

// NewConfirmedOwnerWithProposalTransactor creates a new write-only instance of ConfirmedOwnerWithProposal, bound to a specific deployed contract.
func NewConfirmedOwnerWithProposalTransactor(address common.Address, transactor bind.ContractTransactor) (*ConfirmedOwnerWithProposalTransactor, error) {
	contract, err := bindConfirmedOwnerWithProposal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalTransactor{contract: contract}, nil
}

// NewConfirmedOwnerWithProposalFilterer creates a new log filterer instance of ConfirmedOwnerWithProposal, bound to a specific deployed contract.
func NewConfirmedOwnerWithProposalFilterer(address common.Address, filterer bind.ContractFilterer) (*ConfirmedOwnerWithProposalFilterer, error) {
	contract, err := bindConfirmedOwnerWithProposal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalFilterer{contract: contract}, nil
}

// bindConfirmedOwnerWithProposal binds a generic wrapper to an already deployed contract.
func bindConfirmedOwnerWithProposal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConfirmedOwnerWithProposalABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfirmedOwnerWithProposal.Contract.ConfirmedOwnerWithProposalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.ConfirmedOwnerWithProposalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.ConfirmedOwnerWithProposalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfirmedOwnerWithProposal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConfirmedOwnerWithProposal.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalSession) Owner() (common.Address, error) {
	return _ConfirmedOwnerWithProposal.Contract.Owner(&_ConfirmedOwnerWithProposal.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalCallerSession) Owner() (common.Address, error) {
	return _ConfirmedOwnerWithProposal.Contract.Owner(&_ConfirmedOwnerWithProposal.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.AcceptOwnership(&_ConfirmedOwnerWithProposal.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.AcceptOwnership(&_ConfirmedOwnerWithProposal.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.TransferOwnership(&_ConfirmedOwnerWithProposal.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.TransferOwnership(&_ConfirmedOwnerWithProposal.TransactOpts, to)
}

// ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the ConfirmedOwnerWithProposal contract.
type ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator struct {
	Event *ConfirmedOwnerWithProposalOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfirmedOwnerWithProposalOwnershipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConfirmedOwnerWithProposalOwnershipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfirmedOwnerWithProposalOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the ConfirmedOwnerWithProposal contract.
type ConfirmedOwnerWithProposalOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwnerWithProposal.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator{contract: _ConfirmedOwnerWithProposal.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ConfirmedOwnerWithProposalOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwnerWithProposal.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfirmedOwnerWithProposalOwnershipTransferRequested)
				if err := _ConfirmedOwnerWithProposal.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) ParseOwnershipTransferRequested(log types.Log) (*ConfirmedOwnerWithProposalOwnershipTransferRequested, error) {
	event := new(ConfirmedOwnerWithProposalOwnershipTransferRequested)
	if err := _ConfirmedOwnerWithProposal.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfirmedOwnerWithProposalOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ConfirmedOwnerWithProposal contract.
type ConfirmedOwnerWithProposalOwnershipTransferredIterator struct {
	Event *ConfirmedOwnerWithProposalOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConfirmedOwnerWithProposalOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfirmedOwnerWithProposalOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConfirmedOwnerWithProposalOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConfirmedOwnerWithProposalOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfirmedOwnerWithProposalOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfirmedOwnerWithProposalOwnershipTransferred represents a OwnershipTransferred event raised by the ConfirmedOwnerWithProposal contract.
type ConfirmedOwnerWithProposalOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfirmedOwnerWithProposalOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwnerWithProposal.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalOwnershipTransferredIterator{contract: _ConfirmedOwnerWithProposal.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConfirmedOwnerWithProposalOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwnerWithProposal.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfirmedOwnerWithProposalOwnershipTransferred)
				if err := _ConfirmedOwnerWithProposal.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) ParseOwnershipTransferred(log types.Log) (*ConfirmedOwnerWithProposalOwnershipTransferred, error) {
	event := new(ConfirmedOwnerWithProposalOwnershipTransferred)
	if err := _ConfirmedOwnerWithProposal.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AbstractMetaData contains all meta data concerning the OCR2Abstract contract.
var OCR2AbstractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// OCR2AbstractABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR2AbstractMetaData.ABI instead.
var OCR2AbstractABI = OCR2AbstractMetaData.ABI

// OCR2Abstract is an auto generated Go binding around an Ethereum contract.
type OCR2Abstract struct {
	OCR2AbstractCaller     // Read-only binding to the contract
	OCR2AbstractTransactor // Write-only binding to the contract
	OCR2AbstractFilterer   // Log filterer for contract events
}

// OCR2AbstractCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR2AbstractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AbstractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR2AbstractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AbstractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR2AbstractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AbstractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR2AbstractSession struct {
	Contract     *OCR2Abstract     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OCR2AbstractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR2AbstractCallerSession struct {
	Contract *OCR2AbstractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// OCR2AbstractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR2AbstractTransactorSession struct {
	Contract     *OCR2AbstractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// OCR2AbstractRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR2AbstractRaw struct {
	Contract *OCR2Abstract // Generic contract binding to access the raw methods on
}

// OCR2AbstractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR2AbstractCallerRaw struct {
	Contract *OCR2AbstractCaller // Generic read-only contract binding to access the raw methods on
}

// OCR2AbstractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR2AbstractTransactorRaw struct {
	Contract *OCR2AbstractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR2Abstract creates a new instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2Abstract(address common.Address, backend bind.ContractBackend) (*OCR2Abstract, error) {
	contract, err := bindOCR2Abstract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2Abstract{OCR2AbstractCaller: OCR2AbstractCaller{contract: contract}, OCR2AbstractTransactor: OCR2AbstractTransactor{contract: contract}, OCR2AbstractFilterer: OCR2AbstractFilterer{contract: contract}}, nil
}

// NewOCR2AbstractCaller creates a new read-only instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2AbstractCaller(address common.Address, caller bind.ContractCaller) (*OCR2AbstractCaller, error) {
	contract, err := bindOCR2Abstract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractCaller{contract: contract}, nil
}

// NewOCR2AbstractTransactor creates a new write-only instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2AbstractTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2AbstractTransactor, error) {
	contract, err := bindOCR2Abstract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractTransactor{contract: contract}, nil
}

// NewOCR2AbstractFilterer creates a new log filterer instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2AbstractFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2AbstractFilterer, error) {
	contract, err := bindOCR2Abstract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractFilterer{contract: contract}, nil
}

// bindOCR2Abstract binds a generic wrapper to an already deployed contract.
func bindOCR2Abstract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OCR2AbstractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2Abstract *OCR2AbstractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Abstract.Contract.OCR2AbstractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2Abstract *OCR2AbstractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.OCR2AbstractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2Abstract *OCR2AbstractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.OCR2AbstractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2Abstract *OCR2AbstractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Abstract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2Abstract *OCR2AbstractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2Abstract *OCR2AbstractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.contract.Transact(opts, method, params...)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Abstract *OCR2AbstractCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	var out []interface{}
	err := _OCR2Abstract.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Abstract *OCR2AbstractSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDetails(&_OCR2Abstract.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Abstract *OCR2AbstractCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDetails(&_OCR2Abstract.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	var out []interface{}
	err := _OCR2Abstract.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(struct {
		ScanLogs     bool
		ConfigDigest [32]byte
		Epoch        uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDigestAndEpoch(&_OCR2Abstract.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractCallerSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDigestAndEpoch(&_OCR2Abstract.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Abstract *OCR2AbstractCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCR2Abstract.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Abstract *OCR2AbstractSession) TypeAndVersion() (string, error) {
	return _OCR2Abstract.Contract.TypeAndVersion(&_OCR2Abstract.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Abstract *OCR2AbstractCallerSession) TypeAndVersion() (string, error) {
	return _OCR2Abstract.Contract.TypeAndVersion(&_OCR2Abstract.CallOpts)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Abstract *OCR2AbstractTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Abstract.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Abstract *OCR2AbstractSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.SetConfig(&_OCR2Abstract.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Abstract *OCR2AbstractTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.SetConfig(&_OCR2Abstract.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Abstract *OCR2AbstractTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Abstract.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Abstract *OCR2AbstractSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.Transmit(&_OCR2Abstract.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Abstract *OCR2AbstractTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.Transmit(&_OCR2Abstract.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// OCR2AbstractConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the OCR2Abstract contract.
type OCR2AbstractConfigSetIterator struct {
	Event *OCR2AbstractConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AbstractConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AbstractConfigSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AbstractConfigSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AbstractConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AbstractConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AbstractConfigSet represents a ConfigSet event raised by the OCR2Abstract contract.
type OCR2AbstractConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Abstract *OCR2AbstractFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OCR2AbstractConfigSetIterator, error) {

	logs, sub, err := _OCR2Abstract.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractConfigSetIterator{contract: _OCR2Abstract.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Abstract *OCR2AbstractFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2AbstractConfigSet) (event.Subscription, error) {

	logs, sub, err := _OCR2Abstract.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AbstractConfigSet)
				if err := _OCR2Abstract.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

// ParseConfigSet is a log parse operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Abstract *OCR2AbstractFilterer) ParseConfigSet(log types.Log) (*OCR2AbstractConfigSet, error) {
	event := new(OCR2AbstractConfigSet)
	if err := _OCR2Abstract.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AbstractTransmittedIterator is returned from FilterTransmitted and is used to iterate over the raw logs and unpacked data for Transmitted events raised by the OCR2Abstract contract.
type OCR2AbstractTransmittedIterator struct {
	Event *OCR2AbstractTransmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2AbstractTransmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AbstractTransmitted)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2AbstractTransmitted)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2AbstractTransmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AbstractTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AbstractTransmitted represents a Transmitted event raised by the OCR2Abstract contract.
type OCR2AbstractTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransmitted is a free log retrieval operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractFilterer) FilterTransmitted(opts *bind.FilterOpts) (*OCR2AbstractTransmittedIterator, error) {

	logs, sub, err := _OCR2Abstract.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractTransmittedIterator{contract: _OCR2Abstract.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

// WatchTransmitted is a free log subscription operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OCR2AbstractTransmitted) (event.Subscription, error) {

	logs, sub, err := _OCR2Abstract.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AbstractTransmitted)
				if err := _OCR2Abstract.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

// ParseTransmitted is a log parse operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractFilterer) ParseTransmitted(log types.Log) (*OCR2AbstractTransmitted, error) {
	event := new(OCR2AbstractTransmitted)
	if err := _OCR2Abstract.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2BaseMetaData contains all meta data concerning the OCR2Base contract.
var OCR2BaseMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// OCR2BaseABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR2BaseMetaData.ABI instead.
var OCR2BaseABI = OCR2BaseMetaData.ABI

// OCR2Base is an auto generated Go binding around an Ethereum contract.
type OCR2Base struct {
	OCR2BaseCaller     // Read-only binding to the contract
	OCR2BaseTransactor // Write-only binding to the contract
	OCR2BaseFilterer   // Log filterer for contract events
}

// OCR2BaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR2BaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2BaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR2BaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2BaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR2BaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2BaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR2BaseSession struct {
	Contract     *OCR2Base         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OCR2BaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR2BaseCallerSession struct {
	Contract *OCR2BaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// OCR2BaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR2BaseTransactorSession struct {
	Contract     *OCR2BaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// OCR2BaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR2BaseRaw struct {
	Contract *OCR2Base // Generic contract binding to access the raw methods on
}

// OCR2BaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR2BaseCallerRaw struct {
	Contract *OCR2BaseCaller // Generic read-only contract binding to access the raw methods on
}

// OCR2BaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR2BaseTransactorRaw struct {
	Contract *OCR2BaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR2Base creates a new instance of OCR2Base, bound to a specific deployed contract.
func NewOCR2Base(address common.Address, backend bind.ContractBackend) (*OCR2Base, error) {
	contract, err := bindOCR2Base(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2Base{OCR2BaseCaller: OCR2BaseCaller{contract: contract}, OCR2BaseTransactor: OCR2BaseTransactor{contract: contract}, OCR2BaseFilterer: OCR2BaseFilterer{contract: contract}}, nil
}

// NewOCR2BaseCaller creates a new read-only instance of OCR2Base, bound to a specific deployed contract.
func NewOCR2BaseCaller(address common.Address, caller bind.ContractCaller) (*OCR2BaseCaller, error) {
	contract, err := bindOCR2Base(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2BaseCaller{contract: contract}, nil
}

// NewOCR2BaseTransactor creates a new write-only instance of OCR2Base, bound to a specific deployed contract.
func NewOCR2BaseTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2BaseTransactor, error) {
	contract, err := bindOCR2Base(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2BaseTransactor{contract: contract}, nil
}

// NewOCR2BaseFilterer creates a new log filterer instance of OCR2Base, bound to a specific deployed contract.
func NewOCR2BaseFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2BaseFilterer, error) {
	contract, err := bindOCR2Base(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2BaseFilterer{contract: contract}, nil
}

// bindOCR2Base binds a generic wrapper to an already deployed contract.
func bindOCR2Base(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OCR2BaseABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2Base *OCR2BaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Base.Contract.OCR2BaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2Base *OCR2BaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Base.Contract.OCR2BaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2Base *OCR2BaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Base.Contract.OCR2BaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2Base *OCR2BaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Base.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2Base *OCR2BaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Base.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2Base *OCR2BaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Base.Contract.contract.Transact(opts, method, params...)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Base *OCR2BaseCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	var out []interface{}
	err := _OCR2Base.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Base *OCR2BaseSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2Base.Contract.LatestConfigDetails(&_OCR2Base.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Base *OCR2BaseCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2Base.Contract.LatestConfigDetails(&_OCR2Base.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Base *OCR2BaseCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	var out []interface{}
	err := _OCR2Base.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(struct {
		ScanLogs     bool
		ConfigDigest [32]byte
		Epoch        uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Base *OCR2BaseSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2Base.Contract.LatestConfigDigestAndEpoch(&_OCR2Base.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Base *OCR2BaseCallerSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2Base.Contract.LatestConfigDigestAndEpoch(&_OCR2Base.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OCR2Base *OCR2BaseCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2Base.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OCR2Base *OCR2BaseSession) Owner() (common.Address, error) {
	return _OCR2Base.Contract.Owner(&_OCR2Base.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OCR2Base *OCR2BaseCallerSession) Owner() (common.Address, error) {
	return _OCR2Base.Contract.Owner(&_OCR2Base.CallOpts)
}

// Transmitters is a free data retrieval call binding the contract method 0x81411834.
//
// Solidity: function transmitters() view returns(address[])
func (_OCR2Base *OCR2BaseCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OCR2Base.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// Transmitters is a free data retrieval call binding the contract method 0x81411834.
//
// Solidity: function transmitters() view returns(address[])
func (_OCR2Base *OCR2BaseSession) Transmitters() ([]common.Address, error) {
	return _OCR2Base.Contract.Transmitters(&_OCR2Base.CallOpts)
}

// Transmitters is a free data retrieval call binding the contract method 0x81411834.
//
// Solidity: function transmitters() view returns(address[])
func (_OCR2Base *OCR2BaseCallerSession) Transmitters() ([]common.Address, error) {
	return _OCR2Base.Contract.Transmitters(&_OCR2Base.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Base *OCR2BaseCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCR2Base.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Base *OCR2BaseSession) TypeAndVersion() (string, error) {
	return _OCR2Base.Contract.TypeAndVersion(&_OCR2Base.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Base *OCR2BaseCallerSession) TypeAndVersion() (string, error) {
	return _OCR2Base.Contract.TypeAndVersion(&_OCR2Base.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OCR2Base *OCR2BaseTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Base.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OCR2Base *OCR2BaseSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2Base.Contract.AcceptOwnership(&_OCR2Base.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OCR2Base *OCR2BaseTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2Base.Contract.AcceptOwnership(&_OCR2Base.TransactOpts)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] _signers, address[] _transmitters, uint8 _f, bytes _onchainConfig, uint64 _offchainConfigVersion, bytes _offchainConfig) returns()
func (_OCR2Base *OCR2BaseTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Base.contract.Transact(opts, "setConfig", _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] _signers, address[] _transmitters, uint8 _f, bytes _onchainConfig, uint64 _offchainConfigVersion, bytes _offchainConfig) returns()
func (_OCR2Base *OCR2BaseSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Base.Contract.SetConfig(&_OCR2Base.TransactOpts, _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] _signers, address[] _transmitters, uint8 _f, bytes _onchainConfig, uint64 _offchainConfigVersion, bytes _offchainConfig) returns()
func (_OCR2Base *OCR2BaseTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Base.Contract.SetConfig(&_OCR2Base.TransactOpts, _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OCR2Base *OCR2BaseTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OCR2Base.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OCR2Base *OCR2BaseSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2Base.Contract.TransferOwnership(&_OCR2Base.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OCR2Base *OCR2BaseTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2Base.Contract.TransferOwnership(&_OCR2Base.TransactOpts, to)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Base *OCR2BaseTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Base.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Base *OCR2BaseSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Base.Contract.Transmit(&_OCR2Base.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Base *OCR2BaseTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Base.Contract.Transmit(&_OCR2Base.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// OCR2BaseConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the OCR2Base contract.
type OCR2BaseConfigSetIterator struct {
	Event *OCR2BaseConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2BaseConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2BaseConfigSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2BaseConfigSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2BaseConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2BaseConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2BaseConfigSet represents a ConfigSet event raised by the OCR2Base contract.
type OCR2BaseConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Base *OCR2BaseFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OCR2BaseConfigSetIterator, error) {

	logs, sub, err := _OCR2Base.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OCR2BaseConfigSetIterator{contract: _OCR2Base.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Base *OCR2BaseFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2BaseConfigSet) (event.Subscription, error) {

	logs, sub, err := _OCR2Base.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2BaseConfigSet)
				if err := _OCR2Base.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

// ParseConfigSet is a log parse operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Base *OCR2BaseFilterer) ParseConfigSet(log types.Log) (*OCR2BaseConfigSet, error) {
	event := new(OCR2BaseConfigSet)
	if err := _OCR2Base.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2BaseOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the OCR2Base contract.
type OCR2BaseOwnershipTransferRequestedIterator struct {
	Event *OCR2BaseOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2BaseOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2BaseOwnershipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2BaseOwnershipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2BaseOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2BaseOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2BaseOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the OCR2Base contract.
type OCR2BaseOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OCR2Base *OCR2BaseFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2BaseOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2Base.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2BaseOwnershipTransferRequestedIterator{contract: _OCR2Base.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OCR2Base *OCR2BaseFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2BaseOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2Base.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2BaseOwnershipTransferRequested)
				if err := _OCR2Base.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OCR2Base *OCR2BaseFilterer) ParseOwnershipTransferRequested(log types.Log) (*OCR2BaseOwnershipTransferRequested, error) {
	event := new(OCR2BaseOwnershipTransferRequested)
	if err := _OCR2Base.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2BaseOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OCR2Base contract.
type OCR2BaseOwnershipTransferredIterator struct {
	Event *OCR2BaseOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2BaseOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2BaseOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2BaseOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2BaseOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2BaseOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2BaseOwnershipTransferred represents a OwnershipTransferred event raised by the OCR2Base contract.
type OCR2BaseOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OCR2Base *OCR2BaseFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2BaseOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2Base.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2BaseOwnershipTransferredIterator{contract: _OCR2Base.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OCR2Base *OCR2BaseFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2BaseOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2Base.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2BaseOwnershipTransferred)
				if err := _OCR2Base.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OCR2Base *OCR2BaseFilterer) ParseOwnershipTransferred(log types.Log) (*OCR2BaseOwnershipTransferred, error) {
	event := new(OCR2BaseOwnershipTransferred)
	if err := _OCR2Base.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2BaseTransmittedIterator is returned from FilterTransmitted and is used to iterate over the raw logs and unpacked data for Transmitted events raised by the OCR2Base contract.
type OCR2BaseTransmittedIterator struct {
	Event *OCR2BaseTransmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2BaseTransmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2BaseTransmitted)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2BaseTransmitted)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2BaseTransmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2BaseTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2BaseTransmitted represents a Transmitted event raised by the OCR2Base contract.
type OCR2BaseTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransmitted is a free log retrieval operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Base *OCR2BaseFilterer) FilterTransmitted(opts *bind.FilterOpts) (*OCR2BaseTransmittedIterator, error) {

	logs, sub, err := _OCR2Base.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &OCR2BaseTransmittedIterator{contract: _OCR2Base.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

// WatchTransmitted is a free log subscription operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Base *OCR2BaseFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OCR2BaseTransmitted) (event.Subscription, error) {

	logs, sub, err := _OCR2Base.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2BaseTransmitted)
				if err := _OCR2Base.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

// ParseTransmitted is a log parse operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Base *OCR2BaseFilterer) ParseTransmitted(log types.Log) (*OCR2BaseTransmitted, error) {
	event := new(OCR2BaseTransmitted)
	if err := _OCR2Base.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2TitleRequestMetaData contains all meta data concerning the OCR2TitleRequest contract.
var OCR2TitleRequestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"}],\"name\":\"TitleFulfillment\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"name\":\"TitleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"fulfilled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"name\":\"request\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506000338082816200006a5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200009d576200009d816200011b565b505050151560f81b6080526040805160608101909152602a808252620000cd9190620025896020830139620001c7565b620000f16040518060800160405280604e8152602001620025b3604e9139620001c7565b620001156040518060a00160405280606e815260200162002601606e9139620001c7565b620002e0565b6001600160a01b038116331415620001765760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000061565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b8051602080830191909120600854604080519384018390528301529060600160408051601f1981840301815291905280516020909101206008805491925060006200021283620002b6565b91905055507f37adadbbe0ac5130611b65b06c5e2cef03817b6563f93855718a80afca1402ef81836040516200024a92919062000256565b60405180910390a15050565b82815260006020604081840152835180604085015260005b818110156200028c578581018301518582016060015282016200026e565b818111156200029f576000606083870101525b50601f01601f191692909201606001949350505050565b6000600019821415620002d957634e487b7160e01b600052601160045260246000fd5b5060010190565b60805160f81c61228a620002ff600039600061057e015261228a6000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c806381ff704811610081578063b1dc65a41161005b578063b1dc65a4146101f3578063e3d0e71214610206578063f2fde38b1461021957600080fd5b806381ff70481461017b5780638da5cb5b146101ab578063afcb95d7146101d357600080fd5b80632c199889116100b25780632c1998891461014957806379ba50971461015e578063814118341461016657600080fd5b8063181f5a77146100ce5780632aa91bfd14610116575b600080fd5b604080518082018252601c81527f4f4352325469746c655265717565737420312e302e302d616c706861000000006020820152905161010d9190611e25565b60405180910390f35b610139610124366004611c47565b60009081526009602052604090205460ff1690565b604051901515815260200161010d565b61015c610157366004611ce3565b61022c565b005b61015c6102d5565b61016e6103d7565b60405161010d9190611ddf565b6004546002546040805163ffffffff8085168252640100000000909404909316602084015282015260600161010d565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161010d565b60408051600181526000602082018190529181019190915260600161010d565b61015c610201366004611b62565b610446565b61015c610214366004611a95565b610aea565b61015c610227366004611a7a565b6114cf565b80516020808301919091206008546040805193840183905283015290606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101206008805491925060006102938361212a565b91905055507f37adadbbe0ac5130611b65b06c5e2cef03817b6563f93855718a80afca1402ef81836040516102c9929190611e0c565b60405180910390a15050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461035b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6060600780548060200260200160405190810160405280929190818152602001828054801561043c57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610411575b5050505050905090565b60005a604080516020601f8b018190048102820181019092528981529192508a3591818c0135916104969184918491908e908e90819084018382808284376000920191909152506114e392505050565b6040805183815263ffffffff600884901c1660208201527fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a16040805160608101825260025480825260035460ff8082166020850152610100909104169282019290925290831461056c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f636f6e666967446967657374206d69736d6174636800000000000000000000006044820152606401610352565b61057a8b8b8b8b8b8b6115f0565b60007f0000000000000000000000000000000000000000000000000000000000000000156105d7576002826020015183604001516105b89190612038565b6105c2919061205d565b6105cd906001612038565b60ff1690506105ed565b60208201516105e7906001612038565b60ff1690505b888114610656576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e6174757265730000000000006044820152606401610352565b8887146106bf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e00006044820152606401610352565b3360009081526005602090815260408083208151808301909252805460ff80821684529293919291840191610100909104166002811115610702576107026121c1565b6002811115610713576107136121c1565b9052509050600281602001516002811115610730576107306121c1565b14801561077757506007816000015160ff16815481106107525761075261221f565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b6107dd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d697474657200000000000000006044820152606401610352565b5050505050600088886040516107f4929190611dcf565b60405190819003812061080b918c90602001611df2565b60405160208183030381529060405280519060200120905061082b6118cb565b604080518082019091526000808252602082015260005b88811015610ac85760006001858884602081106108615761086161221f565b61086e91901a601b612038565b8d8d868181106108805761088061221f565b905060200201358c8c878181106108995761089961221f565b90506020020135604051600081526020016040526040516108d6949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156108f8573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526005602090815290849020838501909452835460ff80821685529296509294508401916101009004166002811115610978576109786121c1565b6002811115610989576109896121c1565b90525092506001836020015160028111156109a6576109a66121c1565b14610a0d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e00006044820152606401610352565b8251849060ff16601f8110610a2457610a2461221f565b602002015115610a90576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e61747572650000000000000000000000006044820152606401610352565b600184846000015160ff16601f8110610aab57610aab61221f565b911515602090920201525080610ac08161212a565b915050610842565b5050505063ffffffff8110610adf57610adf612163565b505050505050505050565b855185518560ff16601f831115610b5d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e657273000000000000000000000000000000006044820152606401610352565b60008111610bc7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610352565b818314610c55576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e000000000000000000000000000000000000000000000000000000006064820152608401610352565b610c608160036120a6565b8311610cc8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610352565b610cd06116a7565b6040805160c0810182528a8152602081018a905260ff8916918101919091526060810187905267ffffffffffffffff8616608082015260a081018590525b60065415610ec357600654600090610d28906001906120e3565b9050600060068281548110610d3f57610d3f61221f565b60009182526020822001546007805473ffffffffffffffffffffffffffffffffffffffff90921693509084908110610d7957610d7961221f565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526005909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090811690915592909116808452922080549091169055600680549192509080610df957610df96121f0565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190556007805480610e6257610e626121f0565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550610d0e915050565b60005b81515181101561132a5760006005600084600001518481518110610eec57610eec61221f565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff166002811115610f3657610f366121c1565b14610f9d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610352565b6040805180820190915260ff82168152600160208201528251805160059160009185908110610fce57610fce61221f565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561106f5761106f6121c1565b02179055506000915061107f9050565b60056000846020015184815181106110995761109961221f565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff1660028111156110e3576110e36121c1565b1461114a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610352565b6040805180820190915260ff82168152602081016002815250600560008460200151848151811061117d5761117d61221f565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561121e5761121e6121c1565b02179055505082518051600692508390811061123c5761123c61221f565b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90931692909217909155820151805160079190839081106112b8576112b861221f565b60209081029190910181015182546001810184556000938452919092200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909216919091179055806113228161212a565b915050610ec6565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600480547fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff811664010000000063ffffffff4381168202928317855590830481169360019390926000926113bc928692908216911617612010565b92506101000a81548163ffffffff021916908363ffffffff16021790555061141b4630600460009054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a0015161172a565b6002819055825180516003805460ff909216610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff90921691909117905560045460208501516040808701516060880151608089015160a08a015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05986114ba988b98919763ffffffff909216969095919491939192611edd565b60405180910390a15050505050505050505050565b6114d76116a7565b6114e0816117d5565b50565b600080828060200190518101906114fa9190611c60565b600082815260096020526040902054919350915060ff1615611578576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f616c72656164792066756c66696c6c65640000000000000000000000000000006044820152606401610352565b6000828152600960205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f7cc5a0960ca99cf39ef66b30fb0dbec840eb2cbbd2ecf40d13c78a10a47bb763906115e19084908490611e0c565b60405180910390a15050505050565b60006115fd8260206120a6565b6116088560206120a6565b61161488610144611ff8565b61161e9190611ff8565b6116289190611ff8565b611633906000611ff8565b905036811461169e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d6174636800000000000000006044820152606401610352565b50505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611728576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610352565b565b6000808a8a8a8a8a8a8a8a8a60405160200161174e99989796959493929190611e38565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116331415611855576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610352565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b604051806103e00160405280601f906020820280368337509192915050565b60006118fd6118f884611fb2565b611f63565b905082815283838301111561191157600080fd5b828260208301376000602084830101529392505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461194c57600080fd5b919050565b600082601f83011261196257600080fd5b8135602067ffffffffffffffff82111561197e5761197e61224e565b8160051b61198d828201611f63565b8381528281019086840183880185018910156119a857600080fd5b600093505b858410156119d2576119be81611928565b8352600193909301929184019184016119ad565b50979650505050505050565b60008083601f8401126119f057600080fd5b50813567ffffffffffffffff811115611a0857600080fd5b6020830191508360208260051b8501011115611a2357600080fd5b9250929050565b600082601f830112611a3b57600080fd5b611a4a838335602085016118ea565b9392505050565b803567ffffffffffffffff8116811461194c57600080fd5b803560ff8116811461194c57600080fd5b600060208284031215611a8c57600080fd5b611a4a82611928565b60008060008060008060c08789031215611aae57600080fd5b863567ffffffffffffffff80821115611ac657600080fd5b611ad28a838b01611951565b97506020890135915080821115611ae857600080fd5b611af48a838b01611951565b9650611b0260408a01611a69565b95506060890135915080821115611b1857600080fd5b611b248a838b01611a2a565b9450611b3260808a01611a51565b935060a0890135915080821115611b4857600080fd5b50611b5589828a01611a2a565b9150509295509295509295565b60008060008060008060008060e0898b031215611b7e57600080fd5b606089018a811115611b8f57600080fd5b8998503567ffffffffffffffff80821115611ba957600080fd5b818b0191508b601f830112611bbd57600080fd5b813581811115611bcc57600080fd5b8c6020828501011115611bde57600080fd5b6020830199508098505060808b0135915080821115611bfc57600080fd5b611c088c838d016119de565b909750955060a08b0135915080821115611c2157600080fd5b50611c2e8b828c016119de565b999c989b50969995989497949560c00135949350505050565b600060208284031215611c5957600080fd5b5035919050565b60008060408385031215611c7357600080fd5b82519150602083015167ffffffffffffffff811115611c9157600080fd5b8301601f81018513611ca257600080fd5b8051611cb06118f882611fb2565b818152866020838501011115611cc557600080fd5b611cd68260208301602086016120fa565b8093505050509250929050565b600060208284031215611cf557600080fd5b813567ffffffffffffffff811115611d0c57600080fd5b8201601f81018413611d1d57600080fd5b611d2c848235602084016118ea565b949350505050565b600081518084526020808501945080840160005b83811015611d7a57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101611d48565b509495945050505050565b60008151808452611d9d8160208601602086016120fa565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8183823760009101908152919050565b602081526000611a4a6020830184611d34565b828152608081016060836020840137600081529392505050565b828152604060208201526000611d2c6040830184611d85565b602081526000611a4a6020830184611d85565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152611e7f8285018b611d34565b91508382036080850152611e93828a611d34565b915060ff881660a085015283820360c0850152611eb08288611d85565b90861660e08501528381036101008501529050611ecd8185611d85565b9c9b505050505050505050505050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152611f0d8184018a611d34565b90508281036080840152611f218189611d34565b905060ff871660a084015282810360c0840152611f3e8187611d85565b905067ffffffffffffffff851660e0840152828103610100840152611ecd8185611d85565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611faa57611faa61224e565b604052919050565b600067ffffffffffffffff821115611fcc57611fcc61224e565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b6000821982111561200b5761200b612192565b500190565b600063ffffffff80831681851680830382111561202f5761202f612192565b01949350505050565b600060ff821660ff84168060ff0382111561205557612055612192565b019392505050565b600060ff831680612097577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b8060ff84160491505092915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156120de576120de612192565b500290565b6000828210156120f5576120f5612192565b500390565b60005b838110156121155781810151838201526020016120fd565b83811115612124576000848401525b50505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561215c5761215c612192565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052600160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a68747470733a2f2f626c6f672e636861696e2e6c696e6b2f776861742d69732d636861696e6c696e6b2f68747470733a2f2f7777772e636f696e6465736b2e636f6d2f6d61726b2d637562616e2d6261636b65642d6e66742d6d61726b6574706c6163652d6d696e7461626c652d7261697365732d31336d68747470733a2f2f7777772e626c6f6f6d626572672e636f6d2f6f70696e696f6e2f61727469636c65732f323032312d30362d32342f666964656c6974792d6d616e616765722d6f776e65642d67616d6573746f702d6275742d6c61636b65642d6469616d6f6e642d68616e6473",
}

// OCR2TitleRequestABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR2TitleRequestMetaData.ABI instead.
var OCR2TitleRequestABI = OCR2TitleRequestMetaData.ABI

// OCR2TitleRequestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OCR2TitleRequestMetaData.Bin instead.
var OCR2TitleRequestBin = OCR2TitleRequestMetaData.Bin

// DeployOCR2TitleRequest deploys a new Ethereum contract, binding an instance of OCR2TitleRequest to it.
func DeployOCR2TitleRequest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OCR2TitleRequest, error) {
	parsed, err := OCR2TitleRequestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR2TitleRequestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR2TitleRequest{OCR2TitleRequestCaller: OCR2TitleRequestCaller{contract: contract}, OCR2TitleRequestTransactor: OCR2TitleRequestTransactor{contract: contract}, OCR2TitleRequestFilterer: OCR2TitleRequestFilterer{contract: contract}}, nil
}

// OCR2TitleRequest is an auto generated Go binding around an Ethereum contract.
type OCR2TitleRequest struct {
	OCR2TitleRequestCaller     // Read-only binding to the contract
	OCR2TitleRequestTransactor // Write-only binding to the contract
	OCR2TitleRequestFilterer   // Log filterer for contract events
}

// OCR2TitleRequestCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR2TitleRequestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2TitleRequestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR2TitleRequestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2TitleRequestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR2TitleRequestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2TitleRequestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR2TitleRequestSession struct {
	Contract     *OCR2TitleRequest // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OCR2TitleRequestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR2TitleRequestCallerSession struct {
	Contract *OCR2TitleRequestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// OCR2TitleRequestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR2TitleRequestTransactorSession struct {
	Contract     *OCR2TitleRequestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// OCR2TitleRequestRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR2TitleRequestRaw struct {
	Contract *OCR2TitleRequest // Generic contract binding to access the raw methods on
}

// OCR2TitleRequestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR2TitleRequestCallerRaw struct {
	Contract *OCR2TitleRequestCaller // Generic read-only contract binding to access the raw methods on
}

// OCR2TitleRequestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR2TitleRequestTransactorRaw struct {
	Contract *OCR2TitleRequestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR2TitleRequest creates a new instance of OCR2TitleRequest, bound to a specific deployed contract.
func NewOCR2TitleRequest(address common.Address, backend bind.ContractBackend) (*OCR2TitleRequest, error) {
	contract, err := bindOCR2TitleRequest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2TitleRequest{OCR2TitleRequestCaller: OCR2TitleRequestCaller{contract: contract}, OCR2TitleRequestTransactor: OCR2TitleRequestTransactor{contract: contract}, OCR2TitleRequestFilterer: OCR2TitleRequestFilterer{contract: contract}}, nil
}

// NewOCR2TitleRequestCaller creates a new read-only instance of OCR2TitleRequest, bound to a specific deployed contract.
func NewOCR2TitleRequestCaller(address common.Address, caller bind.ContractCaller) (*OCR2TitleRequestCaller, error) {
	contract, err := bindOCR2TitleRequest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2TitleRequestCaller{contract: contract}, nil
}

// NewOCR2TitleRequestTransactor creates a new write-only instance of OCR2TitleRequest, bound to a specific deployed contract.
func NewOCR2TitleRequestTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2TitleRequestTransactor, error) {
	contract, err := bindOCR2TitleRequest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2TitleRequestTransactor{contract: contract}, nil
}

// NewOCR2TitleRequestFilterer creates a new log filterer instance of OCR2TitleRequest, bound to a specific deployed contract.
func NewOCR2TitleRequestFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2TitleRequestFilterer, error) {
	contract, err := bindOCR2TitleRequest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2TitleRequestFilterer{contract: contract}, nil
}

// bindOCR2TitleRequest binds a generic wrapper to an already deployed contract.
func bindOCR2TitleRequest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OCR2TitleRequestABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2TitleRequest *OCR2TitleRequestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2TitleRequest.Contract.OCR2TitleRequestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2TitleRequest *OCR2TitleRequestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.OCR2TitleRequestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2TitleRequest *OCR2TitleRequestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.OCR2TitleRequestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2TitleRequest *OCR2TitleRequestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2TitleRequest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2TitleRequest *OCR2TitleRequestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2TitleRequest *OCR2TitleRequestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.contract.Transact(opts, method, params...)
}

// Fulfilled is a free data retrieval call binding the contract method 0x2aa91bfd.
//
// Solidity: function fulfilled(bytes32 requestId) view returns(bool)
func (_OCR2TitleRequest *OCR2TitleRequestCaller) Fulfilled(opts *bind.CallOpts, requestId [32]byte) (bool, error) {
	var out []interface{}
	err := _OCR2TitleRequest.contract.Call(opts, &out, "fulfilled", requestId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Fulfilled is a free data retrieval call binding the contract method 0x2aa91bfd.
//
// Solidity: function fulfilled(bytes32 requestId) view returns(bool)
func (_OCR2TitleRequest *OCR2TitleRequestSession) Fulfilled(requestId [32]byte) (bool, error) {
	return _OCR2TitleRequest.Contract.Fulfilled(&_OCR2TitleRequest.CallOpts, requestId)
}

// Fulfilled is a free data retrieval call binding the contract method 0x2aa91bfd.
//
// Solidity: function fulfilled(bytes32 requestId) view returns(bool)
func (_OCR2TitleRequest *OCR2TitleRequestCallerSession) Fulfilled(requestId [32]byte) (bool, error) {
	return _OCR2TitleRequest.Contract.Fulfilled(&_OCR2TitleRequest.CallOpts, requestId)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2TitleRequest *OCR2TitleRequestCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	var out []interface{}
	err := _OCR2TitleRequest.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2TitleRequest *OCR2TitleRequestSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2TitleRequest.Contract.LatestConfigDetails(&_OCR2TitleRequest.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2TitleRequest *OCR2TitleRequestCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2TitleRequest.Contract.LatestConfigDetails(&_OCR2TitleRequest.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2TitleRequest *OCR2TitleRequestCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	var out []interface{}
	err := _OCR2TitleRequest.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(struct {
		ScanLogs     bool
		ConfigDigest [32]byte
		Epoch        uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2TitleRequest *OCR2TitleRequestSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2TitleRequest.Contract.LatestConfigDigestAndEpoch(&_OCR2TitleRequest.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2TitleRequest *OCR2TitleRequestCallerSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2TitleRequest.Contract.LatestConfigDigestAndEpoch(&_OCR2TitleRequest.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OCR2TitleRequest *OCR2TitleRequestCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2TitleRequest.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OCR2TitleRequest *OCR2TitleRequestSession) Owner() (common.Address, error) {
	return _OCR2TitleRequest.Contract.Owner(&_OCR2TitleRequest.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OCR2TitleRequest *OCR2TitleRequestCallerSession) Owner() (common.Address, error) {
	return _OCR2TitleRequest.Contract.Owner(&_OCR2TitleRequest.CallOpts)
}

// Transmitters is a free data retrieval call binding the contract method 0x81411834.
//
// Solidity: function transmitters() view returns(address[])
func (_OCR2TitleRequest *OCR2TitleRequestCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OCR2TitleRequest.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// Transmitters is a free data retrieval call binding the contract method 0x81411834.
//
// Solidity: function transmitters() view returns(address[])
func (_OCR2TitleRequest *OCR2TitleRequestSession) Transmitters() ([]common.Address, error) {
	return _OCR2TitleRequest.Contract.Transmitters(&_OCR2TitleRequest.CallOpts)
}

// Transmitters is a free data retrieval call binding the contract method 0x81411834.
//
// Solidity: function transmitters() view returns(address[])
func (_OCR2TitleRequest *OCR2TitleRequestCallerSession) Transmitters() ([]common.Address, error) {
	return _OCR2TitleRequest.Contract.Transmitters(&_OCR2TitleRequest.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2TitleRequest *OCR2TitleRequestCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCR2TitleRequest.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2TitleRequest *OCR2TitleRequestSession) TypeAndVersion() (string, error) {
	return _OCR2TitleRequest.Contract.TypeAndVersion(&_OCR2TitleRequest.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2TitleRequest *OCR2TitleRequestCallerSession) TypeAndVersion() (string, error) {
	return _OCR2TitleRequest.Contract.TypeAndVersion(&_OCR2TitleRequest.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OCR2TitleRequest *OCR2TitleRequestTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2TitleRequest.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OCR2TitleRequest *OCR2TitleRequestSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.AcceptOwnership(&_OCR2TitleRequest.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OCR2TitleRequest *OCR2TitleRequestTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.AcceptOwnership(&_OCR2TitleRequest.TransactOpts)
}

// Request is a paid mutator transaction binding the contract method 0x2c199889.
//
// Solidity: function request(string url) returns()
func (_OCR2TitleRequest *OCR2TitleRequestTransactor) Request(opts *bind.TransactOpts, url string) (*types.Transaction, error) {
	return _OCR2TitleRequest.contract.Transact(opts, "request", url)
}

// Request is a paid mutator transaction binding the contract method 0x2c199889.
//
// Solidity: function request(string url) returns()
func (_OCR2TitleRequest *OCR2TitleRequestSession) Request(url string) (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.Request(&_OCR2TitleRequest.TransactOpts, url)
}

// Request is a paid mutator transaction binding the contract method 0x2c199889.
//
// Solidity: function request(string url) returns()
func (_OCR2TitleRequest *OCR2TitleRequestTransactorSession) Request(url string) (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.Request(&_OCR2TitleRequest.TransactOpts, url)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] _signers, address[] _transmitters, uint8 _f, bytes _onchainConfig, uint64 _offchainConfigVersion, bytes _offchainConfig) returns()
func (_OCR2TitleRequest *OCR2TitleRequestTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2TitleRequest.contract.Transact(opts, "setConfig", _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] _signers, address[] _transmitters, uint8 _f, bytes _onchainConfig, uint64 _offchainConfigVersion, bytes _offchainConfig) returns()
func (_OCR2TitleRequest *OCR2TitleRequestSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.SetConfig(&_OCR2TitleRequest.TransactOpts, _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] _signers, address[] _transmitters, uint8 _f, bytes _onchainConfig, uint64 _offchainConfigVersion, bytes _offchainConfig) returns()
func (_OCR2TitleRequest *OCR2TitleRequestTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.SetConfig(&_OCR2TitleRequest.TransactOpts, _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OCR2TitleRequest *OCR2TitleRequestTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OCR2TitleRequest.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OCR2TitleRequest *OCR2TitleRequestSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.TransferOwnership(&_OCR2TitleRequest.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OCR2TitleRequest *OCR2TitleRequestTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.TransferOwnership(&_OCR2TitleRequest.TransactOpts, to)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2TitleRequest *OCR2TitleRequestTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2TitleRequest.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2TitleRequest *OCR2TitleRequestSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.Transmit(&_OCR2TitleRequest.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2TitleRequest *OCR2TitleRequestTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2TitleRequest.Contract.Transmit(&_OCR2TitleRequest.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// OCR2TitleRequestConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the OCR2TitleRequest contract.
type OCR2TitleRequestConfigSetIterator struct {
	Event *OCR2TitleRequestConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2TitleRequestConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2TitleRequestConfigSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2TitleRequestConfigSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2TitleRequestConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2TitleRequestConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2TitleRequestConfigSet represents a ConfigSet event raised by the OCR2TitleRequest contract.
type OCR2TitleRequestConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OCR2TitleRequestConfigSetIterator, error) {

	logs, sub, err := _OCR2TitleRequest.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OCR2TitleRequestConfigSetIterator{contract: _OCR2TitleRequest.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2TitleRequestConfigSet) (event.Subscription, error) {

	logs, sub, err := _OCR2TitleRequest.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2TitleRequestConfigSet)
				if err := _OCR2TitleRequest.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

// ParseConfigSet is a log parse operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) ParseConfigSet(log types.Log) (*OCR2TitleRequestConfigSet, error) {
	event := new(OCR2TitleRequestConfigSet)
	if err := _OCR2TitleRequest.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2TitleRequestOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the OCR2TitleRequest contract.
type OCR2TitleRequestOwnershipTransferRequestedIterator struct {
	Event *OCR2TitleRequestOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2TitleRequestOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2TitleRequestOwnershipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2TitleRequestOwnershipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2TitleRequestOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2TitleRequestOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2TitleRequestOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the OCR2TitleRequest contract.
type OCR2TitleRequestOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2TitleRequestOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2TitleRequest.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2TitleRequestOwnershipTransferRequestedIterator{contract: _OCR2TitleRequest.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2TitleRequestOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2TitleRequest.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2TitleRequestOwnershipTransferRequested)
				if err := _OCR2TitleRequest.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) ParseOwnershipTransferRequested(log types.Log) (*OCR2TitleRequestOwnershipTransferRequested, error) {
	event := new(OCR2TitleRequestOwnershipTransferRequested)
	if err := _OCR2TitleRequest.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2TitleRequestOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OCR2TitleRequest contract.
type OCR2TitleRequestOwnershipTransferredIterator struct {
	Event *OCR2TitleRequestOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2TitleRequestOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2TitleRequestOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2TitleRequestOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2TitleRequestOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2TitleRequestOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2TitleRequestOwnershipTransferred represents a OwnershipTransferred event raised by the OCR2TitleRequest contract.
type OCR2TitleRequestOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2TitleRequestOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2TitleRequest.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2TitleRequestOwnershipTransferredIterator{contract: _OCR2TitleRequest.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2TitleRequestOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2TitleRequest.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2TitleRequestOwnershipTransferred)
				if err := _OCR2TitleRequest.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) ParseOwnershipTransferred(log types.Log) (*OCR2TitleRequestOwnershipTransferred, error) {
	event := new(OCR2TitleRequestOwnershipTransferred)
	if err := _OCR2TitleRequest.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2TitleRequestTitleFulfillmentIterator is returned from FilterTitleFulfillment and is used to iterate over the raw logs and unpacked data for TitleFulfillment events raised by the OCR2TitleRequest contract.
type OCR2TitleRequestTitleFulfillmentIterator struct {
	Event *OCR2TitleRequestTitleFulfillment // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2TitleRequestTitleFulfillmentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2TitleRequestTitleFulfillment)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2TitleRequestTitleFulfillment)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2TitleRequestTitleFulfillmentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2TitleRequestTitleFulfillmentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2TitleRequestTitleFulfillment represents a TitleFulfillment event raised by the OCR2TitleRequest contract.
type OCR2TitleRequestTitleFulfillment struct {
	RequestId [32]byte
	Title     string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTitleFulfillment is a free log retrieval operation binding the contract event 0x7cc5a0960ca99cf39ef66b30fb0dbec840eb2cbbd2ecf40d13c78a10a47bb763.
//
// Solidity: event TitleFulfillment(bytes32 requestId, string title)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) FilterTitleFulfillment(opts *bind.FilterOpts) (*OCR2TitleRequestTitleFulfillmentIterator, error) {

	logs, sub, err := _OCR2TitleRequest.contract.FilterLogs(opts, "TitleFulfillment")
	if err != nil {
		return nil, err
	}
	return &OCR2TitleRequestTitleFulfillmentIterator{contract: _OCR2TitleRequest.contract, event: "TitleFulfillment", logs: logs, sub: sub}, nil
}

// WatchTitleFulfillment is a free log subscription operation binding the contract event 0x7cc5a0960ca99cf39ef66b30fb0dbec840eb2cbbd2ecf40d13c78a10a47bb763.
//
// Solidity: event TitleFulfillment(bytes32 requestId, string title)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) WatchTitleFulfillment(opts *bind.WatchOpts, sink chan<- *OCR2TitleRequestTitleFulfillment) (event.Subscription, error) {

	logs, sub, err := _OCR2TitleRequest.contract.WatchLogs(opts, "TitleFulfillment")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2TitleRequestTitleFulfillment)
				if err := _OCR2TitleRequest.contract.UnpackLog(event, "TitleFulfillment", log); err != nil {
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

// ParseTitleFulfillment is a log parse operation binding the contract event 0x7cc5a0960ca99cf39ef66b30fb0dbec840eb2cbbd2ecf40d13c78a10a47bb763.
//
// Solidity: event TitleFulfillment(bytes32 requestId, string title)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) ParseTitleFulfillment(log types.Log) (*OCR2TitleRequestTitleFulfillment, error) {
	event := new(OCR2TitleRequestTitleFulfillment)
	if err := _OCR2TitleRequest.contract.UnpackLog(event, "TitleFulfillment", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2TitleRequestTitleRequestIterator is returned from FilterTitleRequest and is used to iterate over the raw logs and unpacked data for TitleRequest events raised by the OCR2TitleRequest contract.
type OCR2TitleRequestTitleRequestIterator struct {
	Event *OCR2TitleRequestTitleRequest // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2TitleRequestTitleRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2TitleRequestTitleRequest)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2TitleRequestTitleRequest)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2TitleRequestTitleRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2TitleRequestTitleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2TitleRequestTitleRequest represents a TitleRequest event raised by the OCR2TitleRequest contract.
type OCR2TitleRequestTitleRequest struct {
	RequestId [32]byte
	Url       string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTitleRequest is a free log retrieval operation binding the contract event 0x37adadbbe0ac5130611b65b06c5e2cef03817b6563f93855718a80afca1402ef.
//
// Solidity: event TitleRequest(bytes32 requestId, string url)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) FilterTitleRequest(opts *bind.FilterOpts) (*OCR2TitleRequestTitleRequestIterator, error) {

	logs, sub, err := _OCR2TitleRequest.contract.FilterLogs(opts, "TitleRequest")
	if err != nil {
		return nil, err
	}
	return &OCR2TitleRequestTitleRequestIterator{contract: _OCR2TitleRequest.contract, event: "TitleRequest", logs: logs, sub: sub}, nil
}

// WatchTitleRequest is a free log subscription operation binding the contract event 0x37adadbbe0ac5130611b65b06c5e2cef03817b6563f93855718a80afca1402ef.
//
// Solidity: event TitleRequest(bytes32 requestId, string url)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) WatchTitleRequest(opts *bind.WatchOpts, sink chan<- *OCR2TitleRequestTitleRequest) (event.Subscription, error) {

	logs, sub, err := _OCR2TitleRequest.contract.WatchLogs(opts, "TitleRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2TitleRequestTitleRequest)
				if err := _OCR2TitleRequest.contract.UnpackLog(event, "TitleRequest", log); err != nil {
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

// ParseTitleRequest is a log parse operation binding the contract event 0x37adadbbe0ac5130611b65b06c5e2cef03817b6563f93855718a80afca1402ef.
//
// Solidity: event TitleRequest(bytes32 requestId, string url)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) ParseTitleRequest(log types.Log) (*OCR2TitleRequestTitleRequest, error) {
	event := new(OCR2TitleRequestTitleRequest)
	if err := _OCR2TitleRequest.contract.UnpackLog(event, "TitleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2TitleRequestTransmittedIterator is returned from FilterTransmitted and is used to iterate over the raw logs and unpacked data for Transmitted events raised by the OCR2TitleRequest contract.
type OCR2TitleRequestTransmittedIterator struct {
	Event *OCR2TitleRequestTransmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OCR2TitleRequestTransmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2TitleRequestTransmitted)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OCR2TitleRequestTransmitted)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OCR2TitleRequestTransmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2TitleRequestTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2TitleRequestTransmitted represents a Transmitted event raised by the OCR2TitleRequest contract.
type OCR2TitleRequestTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransmitted is a free log retrieval operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) FilterTransmitted(opts *bind.FilterOpts) (*OCR2TitleRequestTransmittedIterator, error) {

	logs, sub, err := _OCR2TitleRequest.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &OCR2TitleRequestTransmittedIterator{contract: _OCR2TitleRequest.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

// WatchTransmitted is a free log subscription operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OCR2TitleRequestTransmitted) (event.Subscription, error) {

	logs, sub, err := _OCR2TitleRequest.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2TitleRequestTransmitted)
				if err := _OCR2TitleRequest.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

// ParseTransmitted is a log parse operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2TitleRequest *OCR2TitleRequestFilterer) ParseTransmitted(log types.Log) (*OCR2TitleRequestTransmitted, error) {
	event := new(OCR2TitleRequestTransmitted)
	if err := _OCR2TitleRequest.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnableInterfaceMetaData contains all meta data concerning the OwnableInterface contract.
var OwnableInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// OwnableInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnableInterfaceMetaData.ABI instead.
var OwnableInterfaceABI = OwnableInterfaceMetaData.ABI

// OwnableInterface is an auto generated Go binding around an Ethereum contract.
type OwnableInterface struct {
	OwnableInterfaceCaller     // Read-only binding to the contract
	OwnableInterfaceTransactor // Write-only binding to the contract
	OwnableInterfaceFilterer   // Log filterer for contract events
}

// OwnableInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnableInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnableInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnableInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnableInterfaceSession struct {
	Contract     *OwnableInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnableInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnableInterfaceCallerSession struct {
	Contract *OwnableInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// OwnableInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnableInterfaceTransactorSession struct {
	Contract     *OwnableInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// OwnableInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnableInterfaceRaw struct {
	Contract *OwnableInterface // Generic contract binding to access the raw methods on
}

// OwnableInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnableInterfaceCallerRaw struct {
	Contract *OwnableInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnableInterfaceTransactorRaw struct {
	Contract *OwnableInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnableInterface creates a new instance of OwnableInterface, bound to a specific deployed contract.
func NewOwnableInterface(address common.Address, backend bind.ContractBackend) (*OwnableInterface, error) {
	contract, err := bindOwnableInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnableInterface{OwnableInterfaceCaller: OwnableInterfaceCaller{contract: contract}, OwnableInterfaceTransactor: OwnableInterfaceTransactor{contract: contract}, OwnableInterfaceFilterer: OwnableInterfaceFilterer{contract: contract}}, nil
}

// NewOwnableInterfaceCaller creates a new read-only instance of OwnableInterface, bound to a specific deployed contract.
func NewOwnableInterfaceCaller(address common.Address, caller bind.ContractCaller) (*OwnableInterfaceCaller, error) {
	contract, err := bindOwnableInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableInterfaceCaller{contract: contract}, nil
}

// NewOwnableInterfaceTransactor creates a new write-only instance of OwnableInterface, bound to a specific deployed contract.
func NewOwnableInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableInterfaceTransactor, error) {
	contract, err := bindOwnableInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableInterfaceTransactor{contract: contract}, nil
}

// NewOwnableInterfaceFilterer creates a new log filterer instance of OwnableInterface, bound to a specific deployed contract.
func NewOwnableInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableInterfaceFilterer, error) {
	contract, err := bindOwnableInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableInterfaceFilterer{contract: contract}, nil
}

// bindOwnableInterface binds a generic wrapper to an already deployed contract.
func bindOwnableInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnableInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableInterface *OwnableInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableInterface.Contract.OwnableInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableInterface *OwnableInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableInterface.Contract.OwnableInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableInterface *OwnableInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableInterface.Contract.OwnableInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableInterface *OwnableInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableInterface *OwnableInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableInterface *OwnableInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableInterface.Contract.contract.Transact(opts, method, params...)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnableInterface *OwnableInterfaceTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableInterface.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnableInterface *OwnableInterfaceSession) AcceptOwnership() (*types.Transaction, error) {
	return _OwnableInterface.Contract.AcceptOwnership(&_OwnableInterface.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnableInterface *OwnableInterfaceTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OwnableInterface.Contract.AcceptOwnership(&_OwnableInterface.TransactOpts)
}

// Owner is a paid mutator transaction binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() returns(address)
func (_OwnableInterface *OwnableInterfaceTransactor) Owner(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableInterface.contract.Transact(opts, "owner")
}

// Owner is a paid mutator transaction binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() returns(address)
func (_OwnableInterface *OwnableInterfaceSession) Owner() (*types.Transaction, error) {
	return _OwnableInterface.Contract.Owner(&_OwnableInterface.TransactOpts)
}

// Owner is a paid mutator transaction binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() returns(address)
func (_OwnableInterface *OwnableInterfaceTransactorSession) Owner() (*types.Transaction, error) {
	return _OwnableInterface.Contract.Owner(&_OwnableInterface.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address recipient) returns()
func (_OwnableInterface *OwnableInterfaceTransactor) TransferOwnership(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _OwnableInterface.contract.Transact(opts, "transferOwnership", recipient)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address recipient) returns()
func (_OwnableInterface *OwnableInterfaceSession) TransferOwnership(recipient common.Address) (*types.Transaction, error) {
	return _OwnableInterface.Contract.TransferOwnership(&_OwnableInterface.TransactOpts, recipient)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address recipient) returns()
func (_OwnableInterface *OwnableInterfaceTransactorSession) TransferOwnership(recipient common.Address) (*types.Transaction, error) {
	return _OwnableInterface.Contract.TransferOwnership(&_OwnableInterface.TransactOpts, recipient)
}

// OwnerIsCreatorMetaData contains all meta data concerning the OwnerIsCreator contract.
var OwnerIsCreatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610149565b6001600160a01b0381163314156100f85760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610368806101586000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806379ba5097146100465780638da5cb5b14610050578063f2fde38b1461007c575b600080fd5b61004e61008f565b005b6000546040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61004e61008a36600461031e565b610191565b60015473ffffffffffffffffffffffffffffffffffffffff163314610115576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6101996101a5565b6101a281610228565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610226576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161010c565b565b73ffffffffffffffffffffffffffffffffffffffff81163314156102a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161010c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561033057600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461035457600080fd5b939250505056fea164736f6c6343000806000a",
}

// OwnerIsCreatorABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnerIsCreatorMetaData.ABI instead.
var OwnerIsCreatorABI = OwnerIsCreatorMetaData.ABI

// OwnerIsCreatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OwnerIsCreatorMetaData.Bin instead.
var OwnerIsCreatorBin = OwnerIsCreatorMetaData.Bin

// DeployOwnerIsCreator deploys a new Ethereum contract, binding an instance of OwnerIsCreator to it.
func DeployOwnerIsCreator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OwnerIsCreator, error) {
	parsed, err := OwnerIsCreatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OwnerIsCreatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OwnerIsCreator{OwnerIsCreatorCaller: OwnerIsCreatorCaller{contract: contract}, OwnerIsCreatorTransactor: OwnerIsCreatorTransactor{contract: contract}, OwnerIsCreatorFilterer: OwnerIsCreatorFilterer{contract: contract}}, nil
}

// OwnerIsCreator is an auto generated Go binding around an Ethereum contract.
type OwnerIsCreator struct {
	OwnerIsCreatorCaller     // Read-only binding to the contract
	OwnerIsCreatorTransactor // Write-only binding to the contract
	OwnerIsCreatorFilterer   // Log filterer for contract events
}

// OwnerIsCreatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnerIsCreatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerIsCreatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnerIsCreatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerIsCreatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnerIsCreatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerIsCreatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnerIsCreatorSession struct {
	Contract     *OwnerIsCreator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnerIsCreatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnerIsCreatorCallerSession struct {
	Contract *OwnerIsCreatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// OwnerIsCreatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnerIsCreatorTransactorSession struct {
	Contract     *OwnerIsCreatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// OwnerIsCreatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnerIsCreatorRaw struct {
	Contract *OwnerIsCreator // Generic contract binding to access the raw methods on
}

// OwnerIsCreatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnerIsCreatorCallerRaw struct {
	Contract *OwnerIsCreatorCaller // Generic read-only contract binding to access the raw methods on
}

// OwnerIsCreatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnerIsCreatorTransactorRaw struct {
	Contract *OwnerIsCreatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnerIsCreator creates a new instance of OwnerIsCreator, bound to a specific deployed contract.
func NewOwnerIsCreator(address common.Address, backend bind.ContractBackend) (*OwnerIsCreator, error) {
	contract, err := bindOwnerIsCreator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreator{OwnerIsCreatorCaller: OwnerIsCreatorCaller{contract: contract}, OwnerIsCreatorTransactor: OwnerIsCreatorTransactor{contract: contract}, OwnerIsCreatorFilterer: OwnerIsCreatorFilterer{contract: contract}}, nil
}

// NewOwnerIsCreatorCaller creates a new read-only instance of OwnerIsCreator, bound to a specific deployed contract.
func NewOwnerIsCreatorCaller(address common.Address, caller bind.ContractCaller) (*OwnerIsCreatorCaller, error) {
	contract, err := bindOwnerIsCreator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorCaller{contract: contract}, nil
}

// NewOwnerIsCreatorTransactor creates a new write-only instance of OwnerIsCreator, bound to a specific deployed contract.
func NewOwnerIsCreatorTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnerIsCreatorTransactor, error) {
	contract, err := bindOwnerIsCreator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorTransactor{contract: contract}, nil
}

// NewOwnerIsCreatorFilterer creates a new log filterer instance of OwnerIsCreator, bound to a specific deployed contract.
func NewOwnerIsCreatorFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnerIsCreatorFilterer, error) {
	contract, err := bindOwnerIsCreator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorFilterer{contract: contract}, nil
}

// bindOwnerIsCreator binds a generic wrapper to an already deployed contract.
func bindOwnerIsCreator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnerIsCreatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnerIsCreator *OwnerIsCreatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnerIsCreator.Contract.OwnerIsCreatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnerIsCreator *OwnerIsCreatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.OwnerIsCreatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnerIsCreator *OwnerIsCreatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.OwnerIsCreatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnerIsCreator *OwnerIsCreatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnerIsCreator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnerIsCreator *OwnerIsCreatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnerIsCreator *OwnerIsCreatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnerIsCreator *OwnerIsCreatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OwnerIsCreator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnerIsCreator *OwnerIsCreatorSession) Owner() (common.Address, error) {
	return _OwnerIsCreator.Contract.Owner(&_OwnerIsCreator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnerIsCreator *OwnerIsCreatorCallerSession) Owner() (common.Address, error) {
	return _OwnerIsCreator.Contract.Owner(&_OwnerIsCreator.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnerIsCreator *OwnerIsCreatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnerIsCreator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnerIsCreator *OwnerIsCreatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.AcceptOwnership(&_OwnerIsCreator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnerIsCreator *OwnerIsCreatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.AcceptOwnership(&_OwnerIsCreator.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OwnerIsCreator *OwnerIsCreatorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OwnerIsCreator.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OwnerIsCreator *OwnerIsCreatorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.TransferOwnership(&_OwnerIsCreator.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OwnerIsCreator *OwnerIsCreatorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.TransferOwnership(&_OwnerIsCreator.TransactOpts, to)
}

// OwnerIsCreatorOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the OwnerIsCreator contract.
type OwnerIsCreatorOwnershipTransferRequestedIterator struct {
	Event *OwnerIsCreatorOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnerIsCreatorOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnerIsCreatorOwnershipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OwnerIsCreatorOwnershipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnerIsCreatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnerIsCreatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnerIsCreatorOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the OwnerIsCreator contract.
type OwnerIsCreatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OwnerIsCreatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OwnerIsCreator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorOwnershipTransferRequestedIterator{contract: _OwnerIsCreator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OwnerIsCreatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OwnerIsCreator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnerIsCreatorOwnershipTransferRequested)
				if err := _OwnerIsCreator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*OwnerIsCreatorOwnershipTransferRequested, error) {
	event := new(OwnerIsCreatorOwnershipTransferRequested)
	if err := _OwnerIsCreator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnerIsCreatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OwnerIsCreator contract.
type OwnerIsCreatorOwnershipTransferredIterator struct {
	Event *OwnerIsCreatorOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnerIsCreatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnerIsCreatorOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OwnerIsCreatorOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnerIsCreatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnerIsCreatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnerIsCreatorOwnershipTransferred represents a OwnershipTransferred event raised by the OwnerIsCreator contract.
type OwnerIsCreatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OwnerIsCreatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OwnerIsCreator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorOwnershipTransferredIterator{contract: _OwnerIsCreator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnerIsCreatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OwnerIsCreator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnerIsCreatorOwnershipTransferred)
				if err := _OwnerIsCreator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) ParseOwnershipTransferred(log types.Log) (*OwnerIsCreatorOwnershipTransferred, error) {
	event := new(OwnerIsCreatorOwnershipTransferred)
	if err := _OwnerIsCreator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TypeAndVersionInterfaceMetaData contains all meta data concerning the TypeAndVersionInterface contract.
var TypeAndVersionInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// TypeAndVersionInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use TypeAndVersionInterfaceMetaData.ABI instead.
var TypeAndVersionInterfaceABI = TypeAndVersionInterfaceMetaData.ABI

// TypeAndVersionInterface is an auto generated Go binding around an Ethereum contract.
type TypeAndVersionInterface struct {
	TypeAndVersionInterfaceCaller     // Read-only binding to the contract
	TypeAndVersionInterfaceTransactor // Write-only binding to the contract
	TypeAndVersionInterfaceFilterer   // Log filterer for contract events
}

// TypeAndVersionInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypeAndVersionInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypeAndVersionInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TypeAndVersionInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypeAndVersionInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TypeAndVersionInterfaceSession struct {
	Contract     *TypeAndVersionInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TypeAndVersionInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TypeAndVersionInterfaceCallerSession struct {
	Contract *TypeAndVersionInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// TypeAndVersionInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TypeAndVersionInterfaceTransactorSession struct {
	Contract     *TypeAndVersionInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// TypeAndVersionInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type TypeAndVersionInterfaceRaw struct {
	Contract *TypeAndVersionInterface // Generic contract binding to access the raw methods on
}

// TypeAndVersionInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceCallerRaw struct {
	Contract *TypeAndVersionInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// TypeAndVersionInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceTransactorRaw struct {
	Contract *TypeAndVersionInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTypeAndVersionInterface creates a new instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterface(address common.Address, backend bind.ContractBackend) (*TypeAndVersionInterface, error) {
	contract, err := bindTypeAndVersionInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterface{TypeAndVersionInterfaceCaller: TypeAndVersionInterfaceCaller{contract: contract}, TypeAndVersionInterfaceTransactor: TypeAndVersionInterfaceTransactor{contract: contract}, TypeAndVersionInterfaceFilterer: TypeAndVersionInterfaceFilterer{contract: contract}}, nil
}

// NewTypeAndVersionInterfaceCaller creates a new read-only instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterfaceCaller(address common.Address, caller bind.ContractCaller) (*TypeAndVersionInterfaceCaller, error) {
	contract, err := bindTypeAndVersionInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterfaceCaller{contract: contract}, nil
}

// NewTypeAndVersionInterfaceTransactor creates a new write-only instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*TypeAndVersionInterfaceTransactor, error) {
	contract, err := bindTypeAndVersionInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterfaceTransactor{contract: contract}, nil
}

// NewTypeAndVersionInterfaceFilterer creates a new log filterer instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*TypeAndVersionInterfaceFilterer, error) {
	contract, err := bindTypeAndVersionInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterfaceFilterer{contract: contract}, nil
}

// bindTypeAndVersionInterface binds a generic wrapper to an already deployed contract.
func bindTypeAndVersionInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TypeAndVersionInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TypeAndVersionInterface.Contract.TypeAndVersionInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersionInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersionInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TypeAndVersionInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.contract.Transact(opts, method, params...)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_TypeAndVersionInterface *TypeAndVersionInterfaceCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TypeAndVersionInterface.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_TypeAndVersionInterface *TypeAndVersionInterfaceSession) TypeAndVersion() (string, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersion(&_TypeAndVersionInterface.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_TypeAndVersionInterface *TypeAndVersionInterfaceCallerSession) TypeAndVersion() (string, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersion(&_TypeAndVersionInterface.CallOpts)
}
