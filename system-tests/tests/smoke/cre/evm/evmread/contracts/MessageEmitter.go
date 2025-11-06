// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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
	_ = abi.ConvertType
)

// MessageEmitterMetaData contains all meta data concerning the MessageEmitter contract.
var MessageEmitterMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"MessageEmitted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"emitMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"getMessage\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"onReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610655806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630cc4e8d8146100465780632ac0df2614610076578063805f213214610092575b600080fd5b610060600480360381019061005b919061033d565b6100ae565b60405161006d9190610486565b60405180910390f35b610090600480360381019061008b91906102f0565b6100d7565b005b6100ac60048036038101906100a7919061026f565b610114565b005b6060816040516020016100c19190610440565b6040516020818303038152906040529050919050565b7f50ede1f15a65bab9edf83cef0d1ffb1f21234653b3e58170594c3d8685d30e7a8282604051610108929190610462565b60405180910390a15050565b7f50ede1f15a65bab9edf83cef0d1ffb1f21234653b3e58170594c3d8685d30e7a8282604051610145929190610462565b60405180910390a150505050565b6000610166610161846104cd565b6104a8565b905082815260208101848484011115610182576101816105d6565b5b61018d848285610525565b509392505050565b60008083601f8401126101ab576101aa6105cc565b5b8235905067ffffffffffffffff8111156101c8576101c76105c7565b5b6020830191508360018202830111156101e4576101e36105d1565b5b9250929050565b60008083601f840112610201576102006105cc565b5b8235905067ffffffffffffffff81111561021e5761021d6105c7565b5b60208301915083600182028301111561023a576102396105d1565b5b9250929050565b600082601f830112610256576102556105cc565b5b8135610266848260208601610153565b91505092915050565b60008060008060408587031215610289576102886105e0565b5b600085013567ffffffffffffffff8111156102a7576102a66105db565b5b6102b387828801610195565b9450945050602085013567ffffffffffffffff8111156102d6576102d56105db565b5b6102e287828801610195565b925092505092959194509250565b60008060208385031215610307576103066105e0565b5b600083013567ffffffffffffffff811115610325576103246105db565b5b610331858286016101eb565b92509250509250929050565b600060208284031215610353576103526105e0565b5b600082013567ffffffffffffffff811115610371576103706105db565b5b61037d84828501610241565b91505092915050565b60006103928385610509565b935061039f838584610525565b6103a8836105e5565b840190509392505050565b60006103be826104fe565b6103c88185610509565b93506103d8818560208601610534565b6103e1816105e5565b840191505092915050565b60006103f7826104fe565b610401818561051a565b9350610411818560208601610534565b80840191505092915050565b600061042a60148361051a565b9150610435826105f6565b601482019050919050565b600061044b8261041d565b915061045782846103ec565b915081905092915050565b6000602082019050818103600083015261047d818486610386565b90509392505050565b600060208201905081810360008301526104a081846103b3565b905092915050565b60006104b26104c3565b90506104be8282610567565b919050565b6000604051905090565b600067ffffffffffffffff8211156104e8576104e7610598565b5b6104f1826105e5565b9050602081019050919050565b600081519050919050565b600082825260208201905092915050565b600081905092915050565b82818337600083830152505050565b60005b83811015610552578082015181840152602081019050610537565b83811115610561576000848401525b50505050565b610570826105e5565b810181811067ffffffffffffffff8211171561058f5761058e610598565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f6765744d6573736167652072657475726e733a2000000000000000000000000060008201525056fea26469706673582212205bff47dd80305b143ff82b31b64767793aaab904551180d84b1f4e5de35d74f564736f6c63430008060033",
}

// MessageEmitterABI is the input ABI used to generate the binding from.
// Deprecated: Use MessageEmitterMetaData.ABI instead.
var MessageEmitterABI = MessageEmitterMetaData.ABI

// MessageEmitterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MessageEmitterMetaData.Bin instead.
var MessageEmitterBin = MessageEmitterMetaData.Bin

// DeployMessageEmitter deploys a new Ethereum contract, binding an instance of MessageEmitter to it.
func DeployMessageEmitter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MessageEmitter, error) {
	parsed, err := MessageEmitterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MessageEmitterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MessageEmitter{MessageEmitterCaller: MessageEmitterCaller{contract: contract}, MessageEmitterTransactor: MessageEmitterTransactor{contract: contract}, MessageEmitterFilterer: MessageEmitterFilterer{contract: contract}}, nil
}

// MessageEmitter is an auto generated Go binding around an Ethereum contract.
type MessageEmitter struct {
	MessageEmitterCaller     // Read-only binding to the contract
	MessageEmitterTransactor // Write-only binding to the contract
	MessageEmitterFilterer   // Log filterer for contract events
}

// MessageEmitterCaller is an auto generated read-only Go binding around an Ethereum contract.
type MessageEmitterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MessageEmitterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MessageEmitterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MessageEmitterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MessageEmitterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MessageEmitterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MessageEmitterSession struct {
	Contract     *MessageEmitter   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MessageEmitterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MessageEmitterCallerSession struct {
	Contract *MessageEmitterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// MessageEmitterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MessageEmitterTransactorSession struct {
	Contract     *MessageEmitterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// MessageEmitterRaw is an auto generated low-level Go binding around an Ethereum contract.
type MessageEmitterRaw struct {
	Contract *MessageEmitter // Generic contract binding to access the raw methods on
}

// MessageEmitterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MessageEmitterCallerRaw struct {
	Contract *MessageEmitterCaller // Generic read-only contract binding to access the raw methods on
}

// MessageEmitterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MessageEmitterTransactorRaw struct {
	Contract *MessageEmitterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMessageEmitter creates a new instance of MessageEmitter, bound to a specific deployed contract.
func NewMessageEmitter(address common.Address, backend bind.ContractBackend) (*MessageEmitter, error) {
	contract, err := bindMessageEmitter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MessageEmitter{MessageEmitterCaller: MessageEmitterCaller{contract: contract}, MessageEmitterTransactor: MessageEmitterTransactor{contract: contract}, MessageEmitterFilterer: MessageEmitterFilterer{contract: contract}}, nil
}

// NewMessageEmitterCaller creates a new read-only instance of MessageEmitter, bound to a specific deployed contract.
func NewMessageEmitterCaller(address common.Address, caller bind.ContractCaller) (*MessageEmitterCaller, error) {
	contract, err := bindMessageEmitter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MessageEmitterCaller{contract: contract}, nil
}

// NewMessageEmitterTransactor creates a new write-only instance of MessageEmitter, bound to a specific deployed contract.
func NewMessageEmitterTransactor(address common.Address, transactor bind.ContractTransactor) (*MessageEmitterTransactor, error) {
	contract, err := bindMessageEmitter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MessageEmitterTransactor{contract: contract}, nil
}

// NewMessageEmitterFilterer creates a new log filterer instance of MessageEmitter, bound to a specific deployed contract.
func NewMessageEmitterFilterer(address common.Address, filterer bind.ContractFilterer) (*MessageEmitterFilterer, error) {
	contract, err := bindMessageEmitter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MessageEmitterFilterer{contract: contract}, nil
}

// bindMessageEmitter binds a generic wrapper to an already deployed contract.
func bindMessageEmitter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MessageEmitterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MessageEmitter *MessageEmitterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageEmitter.Contract.MessageEmitterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MessageEmitter *MessageEmitterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageEmitter.Contract.MessageEmitterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MessageEmitter *MessageEmitterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageEmitter.Contract.MessageEmitterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MessageEmitter *MessageEmitterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageEmitter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MessageEmitter *MessageEmitterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageEmitter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MessageEmitter *MessageEmitterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageEmitter.Contract.contract.Transact(opts, method, params...)
}

// GetMessage is a free data retrieval call binding the contract method 0x0cc4e8d8.
//
// Solidity: function getMessage(string str) pure returns(string)
func (_MessageEmitter *MessageEmitterCaller) GetMessage(opts *bind.CallOpts, str string) (string, error) {
	var out []interface{}
	err := _MessageEmitter.contract.Call(opts, &out, "getMessage", str)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetMessage is a free data retrieval call binding the contract method 0x0cc4e8d8.
//
// Solidity: function getMessage(string str) pure returns(string)
func (_MessageEmitter *MessageEmitterSession) GetMessage(str string) (string, error) {
	return _MessageEmitter.Contract.GetMessage(&_MessageEmitter.CallOpts, str)
}

// GetMessage is a free data retrieval call binding the contract method 0x0cc4e8d8.
//
// Solidity: function getMessage(string str) pure returns(string)
func (_MessageEmitter *MessageEmitterCallerSession) GetMessage(str string) (string, error) {
	return _MessageEmitter.Contract.GetMessage(&_MessageEmitter.CallOpts, str)
}

// EmitMessage is a paid mutator transaction binding the contract method 0x2ac0df26.
//
// Solidity: function emitMessage(string message) returns()
func (_MessageEmitter *MessageEmitterTransactor) EmitMessage(opts *bind.TransactOpts, message string) (*types.Transaction, error) {
	return _MessageEmitter.contract.Transact(opts, "emitMessage", message)
}

// EmitMessage is a paid mutator transaction binding the contract method 0x2ac0df26.
//
// Solidity: function emitMessage(string message) returns()
func (_MessageEmitter *MessageEmitterSession) EmitMessage(message string) (*types.Transaction, error) {
	return _MessageEmitter.Contract.EmitMessage(&_MessageEmitter.TransactOpts, message)
}

// EmitMessage is a paid mutator transaction binding the contract method 0x2ac0df26.
//
// Solidity: function emitMessage(string message) returns()
func (_MessageEmitter *MessageEmitterTransactorSession) EmitMessage(message string) (*types.Transaction, error) {
	return _MessageEmitter.Contract.EmitMessage(&_MessageEmitter.TransactOpts, message)
}

// OnReport is a paid mutator transaction binding the contract method 0x805f2132.
//
// Solidity: function onReport(bytes , bytes report) returns()
func (_MessageEmitter *MessageEmitterTransactor) OnReport(opts *bind.TransactOpts, arg0 []byte, report []byte) (*types.Transaction, error) {
	return _MessageEmitter.contract.Transact(opts, "onReport", arg0, report)
}

// OnReport is a paid mutator transaction binding the contract method 0x805f2132.
//
// Solidity: function onReport(bytes , bytes report) returns()
func (_MessageEmitter *MessageEmitterSession) OnReport(arg0 []byte, report []byte) (*types.Transaction, error) {
	return _MessageEmitter.Contract.OnReport(&_MessageEmitter.TransactOpts, arg0, report)
}

// OnReport is a paid mutator transaction binding the contract method 0x805f2132.
//
// Solidity: function onReport(bytes , bytes report) returns()
func (_MessageEmitter *MessageEmitterTransactorSession) OnReport(arg0 []byte, report []byte) (*types.Transaction, error) {
	return _MessageEmitter.Contract.OnReport(&_MessageEmitter.TransactOpts, arg0, report)
}

// MessageEmitterMessageEmittedIterator is returned from FilterMessageEmitted and is used to iterate over the raw logs and unpacked data for MessageEmitted events raised by the MessageEmitter contract.
type MessageEmitterMessageEmittedIterator struct {
	Event *MessageEmitterMessageEmitted // Event containing the contract specifics and raw log

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
func (it *MessageEmitterMessageEmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageEmitterMessageEmitted)
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
		it.Event = new(MessageEmitterMessageEmitted)
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
func (it *MessageEmitterMessageEmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageEmitterMessageEmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageEmitterMessageEmitted represents a MessageEmitted event raised by the MessageEmitter contract.
type MessageEmitterMessageEmitted struct {
	Message string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMessageEmitted is a free log retrieval operation binding the contract event 0x50ede1f15a65bab9edf83cef0d1ffb1f21234653b3e58170594c3d8685d30e7a.
//
// Solidity: event MessageEmitted(string message)
func (_MessageEmitter *MessageEmitterFilterer) FilterMessageEmitted(opts *bind.FilterOpts) (*MessageEmitterMessageEmittedIterator, error) {

	logs, sub, err := _MessageEmitter.contract.FilterLogs(opts, "MessageEmitted")
	if err != nil {
		return nil, err
	}
	return &MessageEmitterMessageEmittedIterator{contract: _MessageEmitter.contract, event: "MessageEmitted", logs: logs, sub: sub}, nil
}

// WatchMessageEmitted is a free log subscription operation binding the contract event 0x50ede1f15a65bab9edf83cef0d1ffb1f21234653b3e58170594c3d8685d30e7a.
//
// Solidity: event MessageEmitted(string message)
func (_MessageEmitter *MessageEmitterFilterer) WatchMessageEmitted(opts *bind.WatchOpts, sink chan<- *MessageEmitterMessageEmitted) (event.Subscription, error) {

	logs, sub, err := _MessageEmitter.contract.WatchLogs(opts, "MessageEmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageEmitterMessageEmitted)
				if err := _MessageEmitter.contract.UnpackLog(event, "MessageEmitted", log); err != nil {
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

// ParseMessageEmitted is a log parse operation binding the contract event 0x50ede1f15a65bab9edf83cef0d1ffb1f21234653b3e58170594c3d8685d30e7a.
//
// Solidity: event MessageEmitted(string message)
func (_MessageEmitter *MessageEmitterFilterer) ParseMessageEmitted(log types.Log) (*MessageEmitterMessageEmitted, error) {
	event := new(MessageEmitterMessageEmitted)
	if err := _MessageEmitter.contract.UnpackLog(event, "MessageEmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
