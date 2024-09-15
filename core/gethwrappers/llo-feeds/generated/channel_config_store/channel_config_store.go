// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package channel_config_store

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

var ChannelConfigStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"donId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"sha\",\"type\":\"bytes32\"}],\"name\":\"NewChannelDefinition\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"sha\",\"type\":\"bytes32\"}],\"name\":\"setChannelDefinitions\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6106d2806101576000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c806379ba50971161005057806379ba5097146101355780638da5cb5b1461013d578063f2fde38b1461016557600080fd5b806301ffc9a714610077578063181f5a77146100e15780635ba5bac214610120575b600080fd5b6100cc610085366004610483565b7fffffffff00000000000000000000000000000000000000000000000000000000167f5ba5bac2000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b604080518082018252601881527f4368616e6e656c436f6e66696753746f726520302e302e310000000000000000602082015290516100d891906104cc565b61013361012e366004610538565b610178565b005b6101336101f5565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100d8565b6101336101733660046105cc565b6102f7565b61018061030b565b63ffffffff84166000908152600260205260408120805482906101a290610602565b91905081905590508463ffffffff167fe5b641a7879fb491e4e5a35a1ce950f0237b2537ee9b1b1e4fb65e29aff1f5e8828686866040516101e69493929190610661565b60405180910390a25050505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461027b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6102ff61030b565b6103088161038e565b50565b60005473ffffffffffffffffffffffffffffffffffffffff16331461038c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610272565b565b3373ffffffffffffffffffffffffffffffffffffffff82160361040d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610272565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561049557600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146104c557600080fd5b9392505050565b600060208083528351808285015260005b818110156104f9578581018301518582016040015282016104dd565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b6000806000806060858703121561054e57600080fd5b843563ffffffff8116811461056257600080fd5b9350602085013567ffffffffffffffff8082111561057f57600080fd5b818701915087601f83011261059357600080fd5b8135818111156105a257600080fd5b8860208285010111156105b457600080fd5b95986020929092019750949560400135945092505050565b6000602082840312156105de57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146104c557600080fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361065a577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b63ffffffff851681526060602082015282606082015282846080830137600060808483010152600060807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f86011683010190508260408301529594505050505056fea164736f6c6343000813000a",
}

var ChannelConfigStoreABI = ChannelConfigStoreMetaData.ABI

var ChannelConfigStoreBin = ChannelConfigStoreMetaData.Bin

func DeployChannelConfigStore(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ChannelConfigStore, error) {
	parsed, err := ChannelConfigStoreMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChannelConfigStoreBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChannelConfigStore{address: address, abi: *parsed, ChannelConfigStoreCaller: ChannelConfigStoreCaller{contract: contract}, ChannelConfigStoreTransactor: ChannelConfigStoreTransactor{contract: contract}, ChannelConfigStoreFilterer: ChannelConfigStoreFilterer{contract: contract}}, nil
}

type ChannelConfigStore struct {
	address common.Address
	abi     abi.ABI
	ChannelConfigStoreCaller
	ChannelConfigStoreTransactor
	ChannelConfigStoreFilterer
}

type ChannelConfigStoreCaller struct {
	contract *bind.BoundContract
}

type ChannelConfigStoreTransactor struct {
	contract *bind.BoundContract
}

type ChannelConfigStoreFilterer struct {
	contract *bind.BoundContract
}

type ChannelConfigStoreSession struct {
	Contract     *ChannelConfigStore
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ChannelConfigStoreCallerSession struct {
	Contract *ChannelConfigStoreCaller
	CallOpts bind.CallOpts
}

type ChannelConfigStoreTransactorSession struct {
	Contract     *ChannelConfigStoreTransactor
	TransactOpts bind.TransactOpts
}

type ChannelConfigStoreRaw struct {
	Contract *ChannelConfigStore
}

type ChannelConfigStoreCallerRaw struct {
	Contract *ChannelConfigStoreCaller
}

type ChannelConfigStoreTransactorRaw struct {
	Contract *ChannelConfigStoreTransactor
}

func NewChannelConfigStore(address common.Address, backend bind.ContractBackend) (*ChannelConfigStore, error) {
	abi, err := abi.JSON(strings.NewReader(ChannelConfigStoreABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindChannelConfigStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStore{address: address, abi: abi, ChannelConfigStoreCaller: ChannelConfigStoreCaller{contract: contract}, ChannelConfigStoreTransactor: ChannelConfigStoreTransactor{contract: contract}, ChannelConfigStoreFilterer: ChannelConfigStoreFilterer{contract: contract}}, nil
}

func NewChannelConfigStoreCaller(address common.Address, caller bind.ContractCaller) (*ChannelConfigStoreCaller, error) {
	contract, err := bindChannelConfigStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreCaller{contract: contract}, nil
}

func NewChannelConfigStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*ChannelConfigStoreTransactor, error) {
	contract, err := bindChannelConfigStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreTransactor{contract: contract}, nil
}

func NewChannelConfigStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*ChannelConfigStoreFilterer, error) {
	contract, err := bindChannelConfigStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreFilterer{contract: contract}, nil
}

func bindChannelConfigStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChannelConfigStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ChannelConfigStore *ChannelConfigStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChannelConfigStore.Contract.ChannelConfigStoreCaller.contract.Call(opts, result, method, params...)
}

func (_ChannelConfigStore *ChannelConfigStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.ChannelConfigStoreTransactor.contract.Transfer(opts)
}

func (_ChannelConfigStore *ChannelConfigStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.ChannelConfigStoreTransactor.contract.Transact(opts, method, params...)
}

func (_ChannelConfigStore *ChannelConfigStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChannelConfigStore.Contract.contract.Call(opts, result, method, params...)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.contract.Transfer(opts)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.contract.Transact(opts, method, params...)
}

func (_ChannelConfigStore *ChannelConfigStoreCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChannelConfigStore.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ChannelConfigStore *ChannelConfigStoreSession) Owner() (common.Address, error) {
	return _ChannelConfigStore.Contract.Owner(&_ChannelConfigStore.CallOpts)
}

func (_ChannelConfigStore *ChannelConfigStoreCallerSession) Owner() (common.Address, error) {
	return _ChannelConfigStore.Contract.Owner(&_ChannelConfigStore.CallOpts)
}

func (_ChannelConfigStore *ChannelConfigStoreCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ChannelConfigStore.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ChannelConfigStore *ChannelConfigStoreSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ChannelConfigStore.Contract.SupportsInterface(&_ChannelConfigStore.CallOpts, interfaceId)
}

func (_ChannelConfigStore *ChannelConfigStoreCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ChannelConfigStore.Contract.SupportsInterface(&_ChannelConfigStore.CallOpts, interfaceId)
}

func (_ChannelConfigStore *ChannelConfigStoreCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ChannelConfigStore.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_ChannelConfigStore *ChannelConfigStoreSession) TypeAndVersion() (string, error) {
	return _ChannelConfigStore.Contract.TypeAndVersion(&_ChannelConfigStore.CallOpts)
}

func (_ChannelConfigStore *ChannelConfigStoreCallerSession) TypeAndVersion() (string, error) {
	return _ChannelConfigStore.Contract.TypeAndVersion(&_ChannelConfigStore.CallOpts)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelConfigStore.contract.Transact(opts, "acceptOwnership")
}

func (_ChannelConfigStore *ChannelConfigStoreSession) AcceptOwnership() (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.AcceptOwnership(&_ChannelConfigStore.TransactOpts)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.AcceptOwnership(&_ChannelConfigStore.TransactOpts)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactor) SetChannelDefinitions(opts *bind.TransactOpts, donId uint32, url string, sha [32]byte) (*types.Transaction, error) {
	return _ChannelConfigStore.contract.Transact(opts, "setChannelDefinitions", donId, url, sha)
}

func (_ChannelConfigStore *ChannelConfigStoreSession) SetChannelDefinitions(donId uint32, url string, sha [32]byte) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.SetChannelDefinitions(&_ChannelConfigStore.TransactOpts, donId, url, sha)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorSession) SetChannelDefinitions(donId uint32, url string, sha [32]byte) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.SetChannelDefinitions(&_ChannelConfigStore.TransactOpts, donId, url, sha)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ChannelConfigStore.contract.Transact(opts, "transferOwnership", to)
}

func (_ChannelConfigStore *ChannelConfigStoreSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.TransferOwnership(&_ChannelConfigStore.TransactOpts, to)
}

func (_ChannelConfigStore *ChannelConfigStoreTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ChannelConfigStore.Contract.TransferOwnership(&_ChannelConfigStore.TransactOpts, to)
}

type ChannelConfigStoreNewChannelDefinitionIterator struct {
	Event *ChannelConfigStoreNewChannelDefinition

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelConfigStoreNewChannelDefinitionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelConfigStoreNewChannelDefinition)
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
		it.Event = new(ChannelConfigStoreNewChannelDefinition)
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

func (it *ChannelConfigStoreNewChannelDefinitionIterator) Error() error {
	return it.fail
}

func (it *ChannelConfigStoreNewChannelDefinitionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelConfigStoreNewChannelDefinition struct {
	DonId   *big.Int
	Version uint32
	Url     string
	Sha     [32]byte
	Raw     types.Log
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) FilterNewChannelDefinition(opts *bind.FilterOpts, donId []*big.Int) (*ChannelConfigStoreNewChannelDefinitionIterator, error) {

	var donIdRule []interface{}
	for _, donIdItem := range donId {
		donIdRule = append(donIdRule, donIdItem)
	}

	logs, sub, err := _ChannelConfigStore.contract.FilterLogs(opts, "NewChannelDefinition", donIdRule)
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreNewChannelDefinitionIterator{contract: _ChannelConfigStore.contract, event: "NewChannelDefinition", logs: logs, sub: sub}, nil
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) WatchNewChannelDefinition(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreNewChannelDefinition, donId []*big.Int) (event.Subscription, error) {

	var donIdRule []interface{}
	for _, donIdItem := range donId {
		donIdRule = append(donIdRule, donIdItem)
	}

	logs, sub, err := _ChannelConfigStore.contract.WatchLogs(opts, "NewChannelDefinition", donIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelConfigStoreNewChannelDefinition)
				if err := _ChannelConfigStore.contract.UnpackLog(event, "NewChannelDefinition", log); err != nil {
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

func (_ChannelConfigStore *ChannelConfigStoreFilterer) ParseNewChannelDefinition(log types.Log) (*ChannelConfigStoreNewChannelDefinition, error) {
	event := new(ChannelConfigStoreNewChannelDefinition)
	if err := _ChannelConfigStore.contract.UnpackLog(event, "NewChannelDefinition", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelConfigStoreOwnershipTransferRequestedIterator struct {
	Event *ChannelConfigStoreOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelConfigStoreOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelConfigStoreOwnershipTransferRequested)
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
		it.Event = new(ChannelConfigStoreOwnershipTransferRequested)
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

func (it *ChannelConfigStoreOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *ChannelConfigStoreOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelConfigStoreOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelConfigStoreOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelConfigStore.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreOwnershipTransferRequestedIterator{contract: _ChannelConfigStore.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelConfigStore.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelConfigStoreOwnershipTransferRequested)
				if err := _ChannelConfigStore.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_ChannelConfigStore *ChannelConfigStoreFilterer) ParseOwnershipTransferRequested(log types.Log) (*ChannelConfigStoreOwnershipTransferRequested, error) {
	event := new(ChannelConfigStoreOwnershipTransferRequested)
	if err := _ChannelConfigStore.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelConfigStoreOwnershipTransferredIterator struct {
	Event *ChannelConfigStoreOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelConfigStoreOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelConfigStoreOwnershipTransferred)
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
		it.Event = new(ChannelConfigStoreOwnershipTransferred)
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

func (it *ChannelConfigStoreOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *ChannelConfigStoreOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelConfigStoreOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelConfigStoreOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelConfigStore.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ChannelConfigStoreOwnershipTransferredIterator{contract: _ChannelConfigStore.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_ChannelConfigStore *ChannelConfigStoreFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelConfigStore.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelConfigStoreOwnershipTransferred)
				if err := _ChannelConfigStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_ChannelConfigStore *ChannelConfigStoreFilterer) ParseOwnershipTransferred(log types.Log) (*ChannelConfigStoreOwnershipTransferred, error) {
	event := new(ChannelConfigStoreOwnershipTransferred)
	if err := _ChannelConfigStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_ChannelConfigStore *ChannelConfigStore) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ChannelConfigStore.abi.Events["NewChannelDefinition"].ID:
		return _ChannelConfigStore.ParseNewChannelDefinition(log)
	case _ChannelConfigStore.abi.Events["OwnershipTransferRequested"].ID:
		return _ChannelConfigStore.ParseOwnershipTransferRequested(log)
	case _ChannelConfigStore.abi.Events["OwnershipTransferred"].ID:
		return _ChannelConfigStore.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ChannelConfigStoreNewChannelDefinition) Topic() common.Hash {
	return common.HexToHash("0xe5b641a7879fb491e4e5a35a1ce950f0237b2537ee9b1b1e4fb65e29aff1f5e8")
}

func (ChannelConfigStoreOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (ChannelConfigStoreOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_ChannelConfigStore *ChannelConfigStore) Address() common.Address {
	return _ChannelConfigStore.address
}

type ChannelConfigStoreInterface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetChannelDefinitions(opts *bind.TransactOpts, donId uint32, url string, sha [32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterNewChannelDefinition(opts *bind.FilterOpts, donId []*big.Int) (*ChannelConfigStoreNewChannelDefinitionIterator, error)

	WatchNewChannelDefinition(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreNewChannelDefinition, donId []*big.Int) (event.Subscription, error)

	ParseNewChannelDefinition(log types.Log) (*ChannelConfigStoreNewChannelDefinition, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelConfigStoreOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*ChannelConfigStoreOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelConfigStoreOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ChannelConfigStoreOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*ChannelConfigStoreOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
