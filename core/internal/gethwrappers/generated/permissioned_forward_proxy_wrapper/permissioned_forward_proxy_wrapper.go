// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package permissioned_forward_proxy_wrapper

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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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

var PermissionedForwardProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"PermissionNotSet\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"PermissionRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"PermissionSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"handler\",\"type\":\"bytes\"}],\"name\":\"forward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"getPermission\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"removePermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"setPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610987806101576000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c80638da5cb5b1161005b5780638da5cb5b14610101578063db7ea6491461011f578063e074bb4714610132578063f2fde38b1461014557600080fd5b80636fadcf721461008257806379ba50971461009757806383c1cd8a1461009f575b600080fd5b610095610090366004610810565b610158565b005b610095610216565b6100d86100ad366004610893565b73ffffffffffffffffffffffffffffffffffffffff9081166000908152600260205260409020541690565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b60005473ffffffffffffffffffffffffffffffffffffffff166100d8565b61009561012d3660046108ae565b610318565b610095610140366004610893565b6103ac565b610095610153366004610893565b610428565b3360009081526002602052604090205473ffffffffffffffffffffffffffffffffffffffff8481169116146101b9576040517ff805512f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61021082828080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505073ffffffffffffffffffffffffffffffffffffffff87169291505061043c565b50505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461029c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610320610485565b73ffffffffffffffffffffffffffffffffffffffff82811660008181526002602090815260409182902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055905192835290917f85de22bdfe368b52aea0d56feaf9bff7d08c7fc7648d4d353d9837584b5d8b58910160405180910390a25050565b6103b4610485565b73ffffffffffffffffffffffffffffffffffffffff811660008181526002602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055517fe004b76fd67c7384775b2fdace54d5f1ef5d5c65b88ba9b62af22df1f620aaf19190a250565b610430610485565b61043981610508565b50565b606061047e83836040518060400160405280601e81526020017f416464726573733a206c6f772d6c6576656c2063616c6c206661696c656400008152506105fd565b9392505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610506576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610293565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610587576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610293565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b606061060c8484600085610614565b949350505050565b6060824710156106a6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610293565b843b61070e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610293565b6000808673ffffffffffffffffffffffffffffffffffffffff168587604051610737919061090d565b60006040518083038185875af1925050503d8060008114610774576040519150601f19603f3d011682016040523d82523d6000602084013e610779565b606091505b5091509150610789828286610794565b979650505050505050565b606083156107a357508161047e565b8251156107b35782518084602001fd5b816040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102939190610929565b803573ffffffffffffffffffffffffffffffffffffffff8116811461080b57600080fd5b919050565b60008060006040848603121561082557600080fd5b61082e846107e7565b9250602084013567ffffffffffffffff8082111561084b57600080fd5b818601915086601f83011261085f57600080fd5b81358181111561086e57600080fd5b87602082850101111561088057600080fd5b6020830194508093505050509250925092565b6000602082840312156108a557600080fd5b61047e826107e7565b600080604083850312156108c157600080fd5b6108ca836107e7565b91506108d8602084016107e7565b90509250929050565b60005b838110156108fc5781810151838201526020016108e4565b838111156102105750506000910152565b6000825161091f8184602087016108e1565b9190910192915050565b60208152600082518060208401526109488160408501602087016108e1565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016919091016040019291505056fea164736f6c634300080d000a",
}

var PermissionedForwardProxyABI = PermissionedForwardProxyMetaData.ABI

var PermissionedForwardProxyBin = PermissionedForwardProxyMetaData.Bin

func DeployPermissionedForwardProxy(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PermissionedForwardProxy, error) {
	parsed, err := PermissionedForwardProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PermissionedForwardProxyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PermissionedForwardProxy{PermissionedForwardProxyCaller: PermissionedForwardProxyCaller{contract: contract}, PermissionedForwardProxyTransactor: PermissionedForwardProxyTransactor{contract: contract}, PermissionedForwardProxyFilterer: PermissionedForwardProxyFilterer{contract: contract}}, nil
}

type PermissionedForwardProxy struct {
	address common.Address
	abi     abi.ABI
	PermissionedForwardProxyCaller
	PermissionedForwardProxyTransactor
	PermissionedForwardProxyFilterer
}

type PermissionedForwardProxyCaller struct {
	contract *bind.BoundContract
}

type PermissionedForwardProxyTransactor struct {
	contract *bind.BoundContract
}

type PermissionedForwardProxyFilterer struct {
	contract *bind.BoundContract
}

type PermissionedForwardProxySession struct {
	Contract     *PermissionedForwardProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type PermissionedForwardProxyCallerSession struct {
	Contract *PermissionedForwardProxyCaller
	CallOpts bind.CallOpts
}

type PermissionedForwardProxyTransactorSession struct {
	Contract     *PermissionedForwardProxyTransactor
	TransactOpts bind.TransactOpts
}

type PermissionedForwardProxyRaw struct {
	Contract *PermissionedForwardProxy
}

type PermissionedForwardProxyCallerRaw struct {
	Contract *PermissionedForwardProxyCaller
}

type PermissionedForwardProxyTransactorRaw struct {
	Contract *PermissionedForwardProxyTransactor
}

func NewPermissionedForwardProxy(address common.Address, backend bind.ContractBackend) (*PermissionedForwardProxy, error) {
	abi, err := abi.JSON(strings.NewReader(PermissionedForwardProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindPermissionedForwardProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PermissionedForwardProxy{address: address, abi: abi, PermissionedForwardProxyCaller: PermissionedForwardProxyCaller{contract: contract}, PermissionedForwardProxyTransactor: PermissionedForwardProxyTransactor{contract: contract}, PermissionedForwardProxyFilterer: PermissionedForwardProxyFilterer{contract: contract}}, nil
}

func NewPermissionedForwardProxyCaller(address common.Address, caller bind.ContractCaller) (*PermissionedForwardProxyCaller, error) {
	contract, err := bindPermissionedForwardProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PermissionedForwardProxyCaller{contract: contract}, nil
}

func NewPermissionedForwardProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*PermissionedForwardProxyTransactor, error) {
	contract, err := bindPermissionedForwardProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PermissionedForwardProxyTransactor{contract: contract}, nil
}

func NewPermissionedForwardProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*PermissionedForwardProxyFilterer, error) {
	contract, err := bindPermissionedForwardProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PermissionedForwardProxyFilterer{contract: contract}, nil
}

func bindPermissionedForwardProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PermissionedForwardProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_PermissionedForwardProxy *PermissionedForwardProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PermissionedForwardProxy.Contract.PermissionedForwardProxyCaller.contract.Call(opts, result, method, params...)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.PermissionedForwardProxyTransactor.contract.Transfer(opts)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.PermissionedForwardProxyTransactor.contract.Transact(opts, method, params...)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PermissionedForwardProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.contract.Transfer(opts)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.contract.Transact(opts, method, params...)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyCaller) GetPermission(opts *bind.CallOpts, sender common.Address) (common.Address, error) {
	var out []interface{}
	err := _PermissionedForwardProxy.contract.Call(opts, &out, "getPermission", sender)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_PermissionedForwardProxy *PermissionedForwardProxySession) GetPermission(sender common.Address) (common.Address, error) {
	return _PermissionedForwardProxy.Contract.GetPermission(&_PermissionedForwardProxy.CallOpts, sender)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyCallerSession) GetPermission(sender common.Address) (common.Address, error) {
	return _PermissionedForwardProxy.Contract.GetPermission(&_PermissionedForwardProxy.CallOpts, sender)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PermissionedForwardProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_PermissionedForwardProxy *PermissionedForwardProxySession) Owner() (common.Address, error) {
	return _PermissionedForwardProxy.Contract.Owner(&_PermissionedForwardProxy.CallOpts)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyCallerSession) Owner() (common.Address, error) {
	return _PermissionedForwardProxy.Contract.Owner(&_PermissionedForwardProxy.CallOpts)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PermissionedForwardProxy.contract.Transact(opts, "acceptOwnership")
}

func (_PermissionedForwardProxy *PermissionedForwardProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.AcceptOwnership(&_PermissionedForwardProxy.TransactOpts)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.AcceptOwnership(&_PermissionedForwardProxy.TransactOpts)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyTransactor) Forward(opts *bind.TransactOpts, target common.Address, handler []byte) (*types.Transaction, error) {
	return _PermissionedForwardProxy.contract.Transact(opts, "forward", target, handler)
}

func (_PermissionedForwardProxy *PermissionedForwardProxySession) Forward(target common.Address, handler []byte) (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.Forward(&_PermissionedForwardProxy.TransactOpts, target, handler)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyTransactorSession) Forward(target common.Address, handler []byte) (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.Forward(&_PermissionedForwardProxy.TransactOpts, target, handler)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyTransactor) RemovePermission(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error) {
	return _PermissionedForwardProxy.contract.Transact(opts, "removePermission", sender)
}

func (_PermissionedForwardProxy *PermissionedForwardProxySession) RemovePermission(sender common.Address) (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.RemovePermission(&_PermissionedForwardProxy.TransactOpts, sender)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyTransactorSession) RemovePermission(sender common.Address) (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.RemovePermission(&_PermissionedForwardProxy.TransactOpts, sender)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyTransactor) SetPermission(opts *bind.TransactOpts, sender common.Address, target common.Address) (*types.Transaction, error) {
	return _PermissionedForwardProxy.contract.Transact(opts, "setPermission", sender, target)
}

func (_PermissionedForwardProxy *PermissionedForwardProxySession) SetPermission(sender common.Address, target common.Address) (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.SetPermission(&_PermissionedForwardProxy.TransactOpts, sender, target)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyTransactorSession) SetPermission(sender common.Address, target common.Address) (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.SetPermission(&_PermissionedForwardProxy.TransactOpts, sender, target)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _PermissionedForwardProxy.contract.Transact(opts, "transferOwnership", to)
}

func (_PermissionedForwardProxy *PermissionedForwardProxySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.TransferOwnership(&_PermissionedForwardProxy.TransactOpts, to)
}

func (_PermissionedForwardProxy *PermissionedForwardProxyTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PermissionedForwardProxy.Contract.TransferOwnership(&_PermissionedForwardProxy.TransactOpts, to)
}

type PermissionedForwardProxyOwnershipTransferRequestedIterator struct {
	Event *PermissionedForwardProxyOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PermissionedForwardProxyOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PermissionedForwardProxyOwnershipTransferRequested)
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
		it.Event = new(PermissionedForwardProxyOwnershipTransferRequested)
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

func (it *PermissionedForwardProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *PermissionedForwardProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PermissionedForwardProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_PermissionedForwardProxy *PermissionedForwardProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PermissionedForwardProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PermissionedForwardProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PermissionedForwardProxyOwnershipTransferRequestedIterator{contract: _PermissionedForwardProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_PermissionedForwardProxy *PermissionedForwardProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PermissionedForwardProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PermissionedForwardProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PermissionedForwardProxyOwnershipTransferRequested)
				if err := _PermissionedForwardProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_PermissionedForwardProxy *PermissionedForwardProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*PermissionedForwardProxyOwnershipTransferRequested, error) {
	event := new(PermissionedForwardProxyOwnershipTransferRequested)
	if err := _PermissionedForwardProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PermissionedForwardProxyOwnershipTransferredIterator struct {
	Event *PermissionedForwardProxyOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PermissionedForwardProxyOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PermissionedForwardProxyOwnershipTransferred)
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
		it.Event = new(PermissionedForwardProxyOwnershipTransferred)
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

func (it *PermissionedForwardProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *PermissionedForwardProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PermissionedForwardProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_PermissionedForwardProxy *PermissionedForwardProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PermissionedForwardProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PermissionedForwardProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PermissionedForwardProxyOwnershipTransferredIterator{contract: _PermissionedForwardProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_PermissionedForwardProxy *PermissionedForwardProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PermissionedForwardProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PermissionedForwardProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PermissionedForwardProxyOwnershipTransferred)
				if err := _PermissionedForwardProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_PermissionedForwardProxy *PermissionedForwardProxyFilterer) ParseOwnershipTransferred(log types.Log) (*PermissionedForwardProxyOwnershipTransferred, error) {
	event := new(PermissionedForwardProxyOwnershipTransferred)
	if err := _PermissionedForwardProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PermissionedForwardProxyPermissionRemovedIterator struct {
	Event *PermissionedForwardProxyPermissionRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PermissionedForwardProxyPermissionRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PermissionedForwardProxyPermissionRemoved)
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
		it.Event = new(PermissionedForwardProxyPermissionRemoved)
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

func (it *PermissionedForwardProxyPermissionRemovedIterator) Error() error {
	return it.fail
}

func (it *PermissionedForwardProxyPermissionRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PermissionedForwardProxyPermissionRemoved struct {
	Sender common.Address
	Raw    types.Log
}

func (_PermissionedForwardProxy *PermissionedForwardProxyFilterer) FilterPermissionRemoved(opts *bind.FilterOpts, sender []common.Address) (*PermissionedForwardProxyPermissionRemovedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PermissionedForwardProxy.contract.FilterLogs(opts, "PermissionRemoved", senderRule)
	if err != nil {
		return nil, err
	}
	return &PermissionedForwardProxyPermissionRemovedIterator{contract: _PermissionedForwardProxy.contract, event: "PermissionRemoved", logs: logs, sub: sub}, nil
}

func (_PermissionedForwardProxy *PermissionedForwardProxyFilterer) WatchPermissionRemoved(opts *bind.WatchOpts, sink chan<- *PermissionedForwardProxyPermissionRemoved, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PermissionedForwardProxy.contract.WatchLogs(opts, "PermissionRemoved", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PermissionedForwardProxyPermissionRemoved)
				if err := _PermissionedForwardProxy.contract.UnpackLog(event, "PermissionRemoved", log); err != nil {
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

func (_PermissionedForwardProxy *PermissionedForwardProxyFilterer) ParsePermissionRemoved(log types.Log) (*PermissionedForwardProxyPermissionRemoved, error) {
	event := new(PermissionedForwardProxyPermissionRemoved)
	if err := _PermissionedForwardProxy.contract.UnpackLog(event, "PermissionRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PermissionedForwardProxyPermissionSetIterator struct {
	Event *PermissionedForwardProxyPermissionSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PermissionedForwardProxyPermissionSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PermissionedForwardProxyPermissionSet)
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
		it.Event = new(PermissionedForwardProxyPermissionSet)
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

func (it *PermissionedForwardProxyPermissionSetIterator) Error() error {
	return it.fail
}

func (it *PermissionedForwardProxyPermissionSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PermissionedForwardProxyPermissionSet struct {
	Sender common.Address
	Target common.Address
	Raw    types.Log
}

func (_PermissionedForwardProxy *PermissionedForwardProxyFilterer) FilterPermissionSet(opts *bind.FilterOpts, sender []common.Address) (*PermissionedForwardProxyPermissionSetIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PermissionedForwardProxy.contract.FilterLogs(opts, "PermissionSet", senderRule)
	if err != nil {
		return nil, err
	}
	return &PermissionedForwardProxyPermissionSetIterator{contract: _PermissionedForwardProxy.contract, event: "PermissionSet", logs: logs, sub: sub}, nil
}

func (_PermissionedForwardProxy *PermissionedForwardProxyFilterer) WatchPermissionSet(opts *bind.WatchOpts, sink chan<- *PermissionedForwardProxyPermissionSet, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PermissionedForwardProxy.contract.WatchLogs(opts, "PermissionSet", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PermissionedForwardProxyPermissionSet)
				if err := _PermissionedForwardProxy.contract.UnpackLog(event, "PermissionSet", log); err != nil {
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

func (_PermissionedForwardProxy *PermissionedForwardProxyFilterer) ParsePermissionSet(log types.Log) (*PermissionedForwardProxyPermissionSet, error) {
	event := new(PermissionedForwardProxyPermissionSet)
	if err := _PermissionedForwardProxy.contract.UnpackLog(event, "PermissionSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_PermissionedForwardProxy *PermissionedForwardProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _PermissionedForwardProxy.abi.Events["OwnershipTransferRequested"].ID:
		return _PermissionedForwardProxy.ParseOwnershipTransferRequested(log)
	case _PermissionedForwardProxy.abi.Events["OwnershipTransferred"].ID:
		return _PermissionedForwardProxy.ParseOwnershipTransferred(log)
	case _PermissionedForwardProxy.abi.Events["PermissionRemoved"].ID:
		return _PermissionedForwardProxy.ParsePermissionRemoved(log)
	case _PermissionedForwardProxy.abi.Events["PermissionSet"].ID:
		return _PermissionedForwardProxy.ParsePermissionSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (PermissionedForwardProxyOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (PermissionedForwardProxyOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (PermissionedForwardProxyPermissionRemoved) Topic() common.Hash {
	return common.HexToHash("0xe004b76fd67c7384775b2fdace54d5f1ef5d5c65b88ba9b62af22df1f620aaf1")
}

func (PermissionedForwardProxyPermissionSet) Topic() common.Hash {
	return common.HexToHash("0x85de22bdfe368b52aea0d56feaf9bff7d08c7fc7648d4d353d9837584b5d8b58")
}

func (_PermissionedForwardProxy *PermissionedForwardProxy) Address() common.Address {
	return _PermissionedForwardProxy.address
}

type PermissionedForwardProxyInterface interface {
	GetPermission(opts *bind.CallOpts, sender common.Address) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Forward(opts *bind.TransactOpts, target common.Address, handler []byte) (*types.Transaction, error)

	RemovePermission(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error)

	SetPermission(opts *bind.TransactOpts, sender common.Address, target common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PermissionedForwardProxyOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PermissionedForwardProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*PermissionedForwardProxyOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PermissionedForwardProxyOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PermissionedForwardProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*PermissionedForwardProxyOwnershipTransferred, error)

	FilterPermissionRemoved(opts *bind.FilterOpts, sender []common.Address) (*PermissionedForwardProxyPermissionRemovedIterator, error)

	WatchPermissionRemoved(opts *bind.WatchOpts, sink chan<- *PermissionedForwardProxyPermissionRemoved, sender []common.Address) (event.Subscription, error)

	ParsePermissionRemoved(log types.Log) (*PermissionedForwardProxyPermissionRemoved, error)

	FilterPermissionSet(opts *bind.FilterOpts, sender []common.Address) (*PermissionedForwardProxyPermissionSetIterator, error)

	WatchPermissionSet(opts *bind.WatchOpts, sink chan<- *PermissionedForwardProxyPermissionSet, sender []common.Address) (event.Subscription, error)

	ParsePermissionSet(log types.Log) (*PermissionedForwardProxyPermissionSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
