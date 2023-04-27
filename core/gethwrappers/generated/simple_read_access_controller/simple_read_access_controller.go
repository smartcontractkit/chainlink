// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package simple_read_access_controller

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

var SimpleReadAccessControllerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600080546001600160a01b031916331790556001805460ff60a01b1916600160a01b1790556109bb806100386000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638823da6c11610076578063a118f2491161005b578063a118f249146101fd578063dc7f012414610230578063f2fde38b14610238576100a3565b80638823da6c146101995780638da5cb5b146101cc576100a3565b80630a756983146100a85780636b14daf8146100b257806379ba5097146101895780638038e4a114610191575b600080fd5b6100b061026b565b005b610175600480360360408110156100c857600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516919081019060408101602082013564010000000081111561010057600080fd5b82018360208201111561011257600080fd5b8035906020019184600183028401116401000000008311171561013457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610368945050505050565b604080519115158252519081900360200190f35b6100b061039b565b6100b061049d565b6100b0600480360360208110156101af57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166105af565b6101d46106e7565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6100b06004803603602081101561021357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610703565b61017561083c565b6100b06004803603602081101561024e57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661085d565b60005473ffffffffffffffffffffffffffffffffffffffff1633146102f157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff161561036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b60006103748383610959565b80610394575073ffffffffffffffffffffffffffffffffffffffff831632145b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461042157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff16331461052357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff1661036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b60005473ffffffffffffffffffffffffffffffffffffffff16331461063557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff16156106e45773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a15b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461078957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff166106e45773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a150565b60015474010000000000000000000000000000000000000000900460ff1681565b60005473ffffffffffffffffffffffffffffffffffffffff1633146108e357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604081205460ff168061039457505060015474010000000000000000000000000000000000000000900460ff16159291505056fea164736f6c6343000606000a",
}

var SimpleReadAccessControllerABI = SimpleReadAccessControllerMetaData.ABI

var SimpleReadAccessControllerBin = SimpleReadAccessControllerMetaData.Bin

func DeploySimpleReadAccessController(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SimpleReadAccessController, error) {
	parsed, err := SimpleReadAccessControllerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimpleReadAccessControllerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimpleReadAccessController{SimpleReadAccessControllerCaller: SimpleReadAccessControllerCaller{contract: contract}, SimpleReadAccessControllerTransactor: SimpleReadAccessControllerTransactor{contract: contract}, SimpleReadAccessControllerFilterer: SimpleReadAccessControllerFilterer{contract: contract}}, nil
}

type SimpleReadAccessController struct {
	address common.Address
	abi     abi.ABI
	SimpleReadAccessControllerCaller
	SimpleReadAccessControllerTransactor
	SimpleReadAccessControllerFilterer
}

type SimpleReadAccessControllerCaller struct {
	contract *bind.BoundContract
}

type SimpleReadAccessControllerTransactor struct {
	contract *bind.BoundContract
}

type SimpleReadAccessControllerFilterer struct {
	contract *bind.BoundContract
}

type SimpleReadAccessControllerSession struct {
	Contract     *SimpleReadAccessController
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type SimpleReadAccessControllerCallerSession struct {
	Contract *SimpleReadAccessControllerCaller
	CallOpts bind.CallOpts
}

type SimpleReadAccessControllerTransactorSession struct {
	Contract     *SimpleReadAccessControllerTransactor
	TransactOpts bind.TransactOpts
}

type SimpleReadAccessControllerRaw struct {
	Contract *SimpleReadAccessController
}

type SimpleReadAccessControllerCallerRaw struct {
	Contract *SimpleReadAccessControllerCaller
}

type SimpleReadAccessControllerTransactorRaw struct {
	Contract *SimpleReadAccessControllerTransactor
}

func NewSimpleReadAccessController(address common.Address, backend bind.ContractBackend) (*SimpleReadAccessController, error) {
	abi, err := abi.JSON(strings.NewReader(SimpleReadAccessControllerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindSimpleReadAccessController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessController{address: address, abi: abi, SimpleReadAccessControllerCaller: SimpleReadAccessControllerCaller{contract: contract}, SimpleReadAccessControllerTransactor: SimpleReadAccessControllerTransactor{contract: contract}, SimpleReadAccessControllerFilterer: SimpleReadAccessControllerFilterer{contract: contract}}, nil
}

func NewSimpleReadAccessControllerCaller(address common.Address, caller bind.ContractCaller) (*SimpleReadAccessControllerCaller, error) {
	contract, err := bindSimpleReadAccessController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerCaller{contract: contract}, nil
}

func NewSimpleReadAccessControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleReadAccessControllerTransactor, error) {
	contract, err := bindSimpleReadAccessController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerTransactor{contract: contract}, nil
}

func NewSimpleReadAccessControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleReadAccessControllerFilterer, error) {
	contract, err := bindSimpleReadAccessController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerFilterer{contract: contract}, nil
}

func bindSimpleReadAccessController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SimpleReadAccessControllerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_SimpleReadAccessController *SimpleReadAccessControllerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleReadAccessController.Contract.SimpleReadAccessControllerCaller.contract.Call(opts, result, method, params...)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.SimpleReadAccessControllerTransactor.contract.Transfer(opts)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.SimpleReadAccessControllerTransactor.contract.Transact(opts, method, params...)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleReadAccessController.Contract.contract.Call(opts, result, method, params...)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.contract.Transfer(opts)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.contract.Transact(opts, method, params...)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SimpleReadAccessController.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_SimpleReadAccessController *SimpleReadAccessControllerSession) CheckEnabled() (bool, error) {
	return _SimpleReadAccessController.Contract.CheckEnabled(&_SimpleReadAccessController.CallOpts)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerCallerSession) CheckEnabled() (bool, error) {
	return _SimpleReadAccessController.Contract.CheckEnabled(&_SimpleReadAccessController.CallOpts)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _SimpleReadAccessController.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_SimpleReadAccessController *SimpleReadAccessControllerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _SimpleReadAccessController.Contract.HasAccess(&_SimpleReadAccessController.CallOpts, _user, _calldata)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _SimpleReadAccessController.Contract.HasAccess(&_SimpleReadAccessController.CallOpts, _user, _calldata)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimpleReadAccessController.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SimpleReadAccessController *SimpleReadAccessControllerSession) Owner() (common.Address, error) {
	return _SimpleReadAccessController.Contract.Owner(&_SimpleReadAccessController.CallOpts)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerCallerSession) Owner() (common.Address, error) {
	return _SimpleReadAccessController.Contract.Owner(&_SimpleReadAccessController.CallOpts)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "acceptOwnership")
}

func (_SimpleReadAccessController *SimpleReadAccessControllerSession) AcceptOwnership() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.AcceptOwnership(&_SimpleReadAccessController.TransactOpts)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.AcceptOwnership(&_SimpleReadAccessController.TransactOpts)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "addAccess", _user)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.AddAccess(&_SimpleReadAccessController.TransactOpts, _user)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.AddAccess(&_SimpleReadAccessController.TransactOpts, _user)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "disableAccessCheck")
}

func (_SimpleReadAccessController *SimpleReadAccessControllerSession) DisableAccessCheck() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.DisableAccessCheck(&_SimpleReadAccessController.TransactOpts)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.DisableAccessCheck(&_SimpleReadAccessController.TransactOpts)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "enableAccessCheck")
}

func (_SimpleReadAccessController *SimpleReadAccessControllerSession) EnableAccessCheck() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.EnableAccessCheck(&_SimpleReadAccessController.TransactOpts)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.EnableAccessCheck(&_SimpleReadAccessController.TransactOpts)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "removeAccess", _user)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.RemoveAccess(&_SimpleReadAccessController.TransactOpts, _user)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.RemoveAccess(&_SimpleReadAccessController.TransactOpts, _user)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "transferOwnership", _to)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.TransferOwnership(&_SimpleReadAccessController.TransactOpts, _to)
}

func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.TransferOwnership(&_SimpleReadAccessController.TransactOpts, _to)
}

type SimpleReadAccessControllerAddedAccessIterator struct {
	Event *SimpleReadAccessControllerAddedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SimpleReadAccessControllerAddedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleReadAccessControllerAddedAccess)
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
		it.Event = new(SimpleReadAccessControllerAddedAccess)
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

func (it *SimpleReadAccessControllerAddedAccessIterator) Error() error {
	return it.fail
}

func (it *SimpleReadAccessControllerAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SimpleReadAccessControllerAddedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*SimpleReadAccessControllerAddedAccessIterator, error) {

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerAddedAccessIterator{contract: _SimpleReadAccessController.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerAddedAccess) (event.Subscription, error) {

	logs, sub, err := _SimpleReadAccessController.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SimpleReadAccessControllerAddedAccess)
				if err := _SimpleReadAccessController.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseAddedAccess(log types.Log) (*SimpleReadAccessControllerAddedAccess, error) {
	event := new(SimpleReadAccessControllerAddedAccess)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SimpleReadAccessControllerCheckAccessDisabledIterator struct {
	Event *SimpleReadAccessControllerCheckAccessDisabled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SimpleReadAccessControllerCheckAccessDisabledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleReadAccessControllerCheckAccessDisabled)
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
		it.Event = new(SimpleReadAccessControllerCheckAccessDisabled)
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

func (it *SimpleReadAccessControllerCheckAccessDisabledIterator) Error() error {
	return it.fail
}

func (it *SimpleReadAccessControllerCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SimpleReadAccessControllerCheckAccessDisabled struct {
	Raw types.Log
}

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*SimpleReadAccessControllerCheckAccessDisabledIterator, error) {

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerCheckAccessDisabledIterator{contract: _SimpleReadAccessController.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerCheckAccessDisabled) (event.Subscription, error) {

	logs, sub, err := _SimpleReadAccessController.contract.WatchLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SimpleReadAccessControllerCheckAccessDisabled)
				if err := _SimpleReadAccessController.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
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

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseCheckAccessDisabled(log types.Log) (*SimpleReadAccessControllerCheckAccessDisabled, error) {
	event := new(SimpleReadAccessControllerCheckAccessDisabled)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SimpleReadAccessControllerCheckAccessEnabledIterator struct {
	Event *SimpleReadAccessControllerCheckAccessEnabled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SimpleReadAccessControllerCheckAccessEnabledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleReadAccessControllerCheckAccessEnabled)
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
		it.Event = new(SimpleReadAccessControllerCheckAccessEnabled)
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

func (it *SimpleReadAccessControllerCheckAccessEnabledIterator) Error() error {
	return it.fail
}

func (it *SimpleReadAccessControllerCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SimpleReadAccessControllerCheckAccessEnabled struct {
	Raw types.Log
}

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*SimpleReadAccessControllerCheckAccessEnabledIterator, error) {

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerCheckAccessEnabledIterator{contract: _SimpleReadAccessController.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerCheckAccessEnabled) (event.Subscription, error) {

	logs, sub, err := _SimpleReadAccessController.contract.WatchLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SimpleReadAccessControllerCheckAccessEnabled)
				if err := _SimpleReadAccessController.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
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

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseCheckAccessEnabled(log types.Log) (*SimpleReadAccessControllerCheckAccessEnabled, error) {
	event := new(SimpleReadAccessControllerCheckAccessEnabled)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SimpleReadAccessControllerOwnershipTransferRequestedIterator struct {
	Event *SimpleReadAccessControllerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SimpleReadAccessControllerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleReadAccessControllerOwnershipTransferRequested)
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
		it.Event = new(SimpleReadAccessControllerOwnershipTransferRequested)
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

func (it *SimpleReadAccessControllerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *SimpleReadAccessControllerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SimpleReadAccessControllerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SimpleReadAccessControllerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerOwnershipTransferRequestedIterator{contract: _SimpleReadAccessController.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleReadAccessController.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SimpleReadAccessControllerOwnershipTransferRequested)
				if err := _SimpleReadAccessController.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseOwnershipTransferRequested(log types.Log) (*SimpleReadAccessControllerOwnershipTransferRequested, error) {
	event := new(SimpleReadAccessControllerOwnershipTransferRequested)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SimpleReadAccessControllerOwnershipTransferredIterator struct {
	Event *SimpleReadAccessControllerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SimpleReadAccessControllerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleReadAccessControllerOwnershipTransferred)
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
		it.Event = new(SimpleReadAccessControllerOwnershipTransferred)
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

func (it *SimpleReadAccessControllerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *SimpleReadAccessControllerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SimpleReadAccessControllerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SimpleReadAccessControllerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerOwnershipTransferredIterator{contract: _SimpleReadAccessController.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleReadAccessController.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SimpleReadAccessControllerOwnershipTransferred)
				if err := _SimpleReadAccessController.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseOwnershipTransferred(log types.Log) (*SimpleReadAccessControllerOwnershipTransferred, error) {
	event := new(SimpleReadAccessControllerOwnershipTransferred)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SimpleReadAccessControllerRemovedAccessIterator struct {
	Event *SimpleReadAccessControllerRemovedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SimpleReadAccessControllerRemovedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleReadAccessControllerRemovedAccess)
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
		it.Event = new(SimpleReadAccessControllerRemovedAccess)
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

func (it *SimpleReadAccessControllerRemovedAccessIterator) Error() error {
	return it.fail
}

func (it *SimpleReadAccessControllerRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SimpleReadAccessControllerRemovedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*SimpleReadAccessControllerRemovedAccessIterator, error) {

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerRemovedAccessIterator{contract: _SimpleReadAccessController.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerRemovedAccess) (event.Subscription, error) {

	logs, sub, err := _SimpleReadAccessController.contract.WatchLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SimpleReadAccessControllerRemovedAccess)
				if err := _SimpleReadAccessController.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
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

func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseRemovedAccess(log types.Log) (*SimpleReadAccessControllerRemovedAccess, error) {
	event := new(SimpleReadAccessControllerRemovedAccess)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_SimpleReadAccessController *SimpleReadAccessController) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _SimpleReadAccessController.abi.Events["AddedAccess"].ID:
		return _SimpleReadAccessController.ParseAddedAccess(log)
	case _SimpleReadAccessController.abi.Events["CheckAccessDisabled"].ID:
		return _SimpleReadAccessController.ParseCheckAccessDisabled(log)
	case _SimpleReadAccessController.abi.Events["CheckAccessEnabled"].ID:
		return _SimpleReadAccessController.ParseCheckAccessEnabled(log)
	case _SimpleReadAccessController.abi.Events["OwnershipTransferRequested"].ID:
		return _SimpleReadAccessController.ParseOwnershipTransferRequested(log)
	case _SimpleReadAccessController.abi.Events["OwnershipTransferred"].ID:
		return _SimpleReadAccessController.ParseOwnershipTransferred(log)
	case _SimpleReadAccessController.abi.Events["RemovedAccess"].ID:
		return _SimpleReadAccessController.ParseRemovedAccess(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (SimpleReadAccessControllerAddedAccess) Topic() common.Hash {
	return common.HexToHash("0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4")
}

func (SimpleReadAccessControllerCheckAccessDisabled) Topic() common.Hash {
	return common.HexToHash("0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638")
}

func (SimpleReadAccessControllerCheckAccessEnabled) Topic() common.Hash {
	return common.HexToHash("0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480")
}

func (SimpleReadAccessControllerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (SimpleReadAccessControllerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (SimpleReadAccessControllerRemovedAccess) Topic() common.Hash {
	return common.HexToHash("0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1")
}

func (_SimpleReadAccessController *SimpleReadAccessController) Address() common.Address {
	return _SimpleReadAccessController.address
}

type SimpleReadAccessControllerInterface interface {
	CheckEnabled(opts *bind.CallOpts) (bool, error)

	HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error)

	DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error)

	EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error)

	RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error)

	FilterAddedAccess(opts *bind.FilterOpts) (*SimpleReadAccessControllerAddedAccessIterator, error)

	WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerAddedAccess) (event.Subscription, error)

	ParseAddedAccess(log types.Log) (*SimpleReadAccessControllerAddedAccess, error)

	FilterCheckAccessDisabled(opts *bind.FilterOpts) (*SimpleReadAccessControllerCheckAccessDisabledIterator, error)

	WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerCheckAccessDisabled) (event.Subscription, error)

	ParseCheckAccessDisabled(log types.Log) (*SimpleReadAccessControllerCheckAccessDisabled, error)

	FilterCheckAccessEnabled(opts *bind.FilterOpts) (*SimpleReadAccessControllerCheckAccessEnabledIterator, error)

	WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerCheckAccessEnabled) (event.Subscription, error)

	ParseCheckAccessEnabled(log types.Log) (*SimpleReadAccessControllerCheckAccessEnabled, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SimpleReadAccessControllerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*SimpleReadAccessControllerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SimpleReadAccessControllerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*SimpleReadAccessControllerOwnershipTransferred, error)

	FilterRemovedAccess(opts *bind.FilterOpts) (*SimpleReadAccessControllerRemovedAccessIterator, error)

	WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerRemovedAccess) (event.Subscription, error)

	ParseRemovedAccess(log types.Log) (*SimpleReadAccessControllerRemovedAccess, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
