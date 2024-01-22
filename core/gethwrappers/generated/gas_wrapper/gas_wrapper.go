// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gas_wrapper

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

var KeeperRegistryCheckUpkeepGasUsageWrapperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistryExecutableInterface\",\"name\":\"keeperRegistry\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getKeeperRegistry\",\"outputs\":[{\"internalType\":\"contractAutomationRegistryExecutableInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"measureCheckGas\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161092f38038061092f83398101604081905261002f91610177565b33806000816100855760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100b5576100b5816100cd565b50505060601b6001600160601b0319166080526101a7565b6001600160a01b0381163314156101265760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161007c565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561018957600080fd5b81516001600160a01b03811681146101a057600080fd5b9392505050565b60805160601c6107646101cb6000396000818160e2015261017001526107646000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c80638da5cb5b116100505780638da5cb5b146100a1578063a33c0660146100e0578063f2fde38b1461010657600080fd5b80636bf490301461006c57806379ba509714610097575b600080fd5b61007f61007a36600461062c565b610119565b60405161008e93929190610658565b60405180910390f35b61009f610262565b005b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161008e565b7f00000000000000000000000000000000000000000000000000000000000000006100bb565b61009f61011436600461051a565b610364565b600060606000805a6040517fc41b813a0000000000000000000000000000000000000000000000000000000081526004810188905273ffffffffffffffffffffffffffffffffffffffff87811660248301529192507f00000000000000000000000000000000000000000000000000000000000000009091169063c41b813a90604401600060405180830381600087803b1580156101b657600080fd5b505af192505050801561020957506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052610206919081019061053c565b60015b6102385760005a61021a90836106ba565b6040805160208101909152600080825296509450925061025b915050565b60005a61024590886106ba565b60019a5095985094965061025b95505050505050565b9250925092565b60015473ffffffffffffffffffffffffffffffffffffffff1633146102e8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61036c610378565b610375816103fb565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1633146103f9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102df565b565b73ffffffffffffffffffffffffffffffffffffffff811633141561047b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102df565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b803573ffffffffffffffffffffffffffffffffffffffff8116811461051557600080fd5b919050565b60006020828403121561052c57600080fd5b610535826104f1565b9392505050565b600080600080600060a0868803121561055457600080fd5b855167ffffffffffffffff8082111561056c57600080fd5b818801915088601f83011261058057600080fd5b81518181111561059257610592610728565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156105d8576105d8610728565b816040528281528b60208487010111156105f157600080fd5b6106028360208301602088016106f8565b60208b015160408c015160608d01516080909d0151929e919d509b9a509098509650505050505050565b6000806040838503121561063f57600080fd5b8235915061064f602084016104f1565b90509250929050565b8315158152606060208201526000835180606084015261067f8160808501602088016106f8565b604083019390935250601f919091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160160800192915050565b6000828210156106f3577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b500390565b60005b838110156107135781810151838201526020016106fb565b83811115610722576000848401525b50505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var KeeperRegistryCheckUpkeepGasUsageWrapperABI = KeeperRegistryCheckUpkeepGasUsageWrapperMetaData.ABI

var KeeperRegistryCheckUpkeepGasUsageWrapperBin = KeeperRegistryCheckUpkeepGasUsageWrapperMetaData.Bin

func DeployKeeperRegistryCheckUpkeepGasUsageWrapper(auth *bind.TransactOpts, backend bind.ContractBackend, keeperRegistry common.Address) (common.Address, *types.Transaction, *KeeperRegistryCheckUpkeepGasUsageWrapper, error) {
	parsed, err := KeeperRegistryCheckUpkeepGasUsageWrapperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryCheckUpkeepGasUsageWrapperBin), backend, keeperRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistryCheckUpkeepGasUsageWrapper{address: address, abi: *parsed, KeeperRegistryCheckUpkeepGasUsageWrapperCaller: KeeperRegistryCheckUpkeepGasUsageWrapperCaller{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapperTransactor: KeeperRegistryCheckUpkeepGasUsageWrapperTransactor{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapperFilterer: KeeperRegistryCheckUpkeepGasUsageWrapperFilterer{contract: contract}}, nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapper struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistryCheckUpkeepGasUsageWrapperCaller
	KeeperRegistryCheckUpkeepGasUsageWrapperTransactor
	KeeperRegistryCheckUpkeepGasUsageWrapperFilterer
}

type KeeperRegistryCheckUpkeepGasUsageWrapperCaller struct {
	contract *bind.BoundContract
}

type KeeperRegistryCheckUpkeepGasUsageWrapperTransactor struct {
	contract *bind.BoundContract
}

type KeeperRegistryCheckUpkeepGasUsageWrapperFilterer struct {
	contract *bind.BoundContract
}

type KeeperRegistryCheckUpkeepGasUsageWrapperSession struct {
	Contract     *KeeperRegistryCheckUpkeepGasUsageWrapper
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistryCheckUpkeepGasUsageWrapperCallerSession struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapperCaller
	CallOpts bind.CallOpts
}

type KeeperRegistryCheckUpkeepGasUsageWrapperTransactorSession struct {
	Contract     *KeeperRegistryCheckUpkeepGasUsageWrapperTransactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistryCheckUpkeepGasUsageWrapperRaw struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapper
}

type KeeperRegistryCheckUpkeepGasUsageWrapperCallerRaw struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapperCaller
}

type KeeperRegistryCheckUpkeepGasUsageWrapperTransactorRaw struct {
	Contract *KeeperRegistryCheckUpkeepGasUsageWrapperTransactor
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapper(address common.Address, backend bind.ContractBackend) (*KeeperRegistryCheckUpkeepGasUsageWrapper, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistryCheckUpkeepGasUsageWrapperABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapper{address: address, abi: abi, KeeperRegistryCheckUpkeepGasUsageWrapperCaller: KeeperRegistryCheckUpkeepGasUsageWrapperCaller{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapperTransactor: KeeperRegistryCheckUpkeepGasUsageWrapperTransactor{contract: contract}, KeeperRegistryCheckUpkeepGasUsageWrapperFilterer: KeeperRegistryCheckUpkeepGasUsageWrapperFilterer{contract: contract}}, nil
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapperCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryCheckUpkeepGasUsageWrapperCaller, error) {
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperCaller{contract: contract}, nil
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryCheckUpkeepGasUsageWrapperTransactor, error) {
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperTransactor{contract: contract}, nil
}

func NewKeeperRegistryCheckUpkeepGasUsageWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryCheckUpkeepGasUsageWrapperFilterer, error) {
	contract, err := bindKeeperRegistryCheckUpkeepGasUsageWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperFilterer{contract: contract}, nil
}

func bindKeeperRegistryCheckUpkeepGasUsageWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistryCheckUpkeepGasUsageWrapperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.KeeperRegistryCheckUpkeepGasUsageWrapperCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.KeeperRegistryCheckUpkeepGasUsageWrapperTransactor.contract.Transfer(opts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.KeeperRegistryCheckUpkeepGasUsageWrapperTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.contract.Transfer(opts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperCaller) GetKeeperRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.Call(opts, &out, "getKeeperRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperSession) GetKeeperRegistry() (common.Address, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.GetKeeperRegistry(&_KeeperRegistryCheckUpkeepGasUsageWrapper.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperCallerSession) GetKeeperRegistry() (common.Address, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.GetKeeperRegistry(&_KeeperRegistryCheckUpkeepGasUsageWrapper.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperSession) Owner() (common.Address, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.Owner(&_KeeperRegistryCheckUpkeepGasUsageWrapper.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperCallerSession) Owner() (common.Address, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.Owner(&_KeeperRegistryCheckUpkeepGasUsageWrapper.CallOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.Transact(opts, "acceptOwnership")
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.AcceptOwnership(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.AcceptOwnership(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactor) MeasureCheckGas(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.Transact(opts, "measureCheckGas", id, from)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperSession) MeasureCheckGas(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.MeasureCheckGas(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts, id, from)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactorSession) MeasureCheckGas(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.MeasureCheckGas(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts, id, from)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.Transact(opts, "transferOwnership", to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.TransferOwnership(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts, to)
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.Contract.TransferOwnership(&_KeeperRegistryCheckUpkeepGasUsageWrapper.TransactOpts, to)
}

type KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested)
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
		it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested)
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

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator{contract: _KeeperRegistryCheckUpkeepGasUsageWrapper.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested)
				if err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested, error) {
	event := new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested)
	if err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator struct {
	Event *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred)
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
		it.Event = new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred)
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

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator{contract: _KeeperRegistryCheckUpkeepGasUsageWrapper.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred)
				if err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapperFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred, error) {
	event := new(KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred)
	if err := _KeeperRegistryCheckUpkeepGasUsageWrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapper) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistryCheckUpkeepGasUsageWrapper.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistryCheckUpkeepGasUsageWrapper.ParseOwnershipTransferRequested(log)
	case _KeeperRegistryCheckUpkeepGasUsageWrapper.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistryCheckUpkeepGasUsageWrapper.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_KeeperRegistryCheckUpkeepGasUsageWrapper *KeeperRegistryCheckUpkeepGasUsageWrapper) Address() common.Address {
	return _KeeperRegistryCheckUpkeepGasUsageWrapper.address
}

type KeeperRegistryCheckUpkeepGasUsageWrapperInterface interface {
	GetKeeperRegistry(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	MeasureCheckGas(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistryCheckUpkeepGasUsageWrapperOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
