// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package registry_module_owner_custom

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

var RegistryModuleOwnerCustomMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"CanOnlySelfRegister\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"administrator\",\"type\":\"address\"}],\"name\":\"AdministratorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"registerAdminViaGetCCIPAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"registerAdminViaOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161087b38038061087b83398101604081905261002f91610187565b33806000816100855760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100b5576100b5816100de565b5050600280546001600160a01b0319166001600160a01b039390931692909217909155506101b7565b336001600160a01b038216036101365760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161007c565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561019957600080fd5b81516001600160a01b03811681146101b057600080fd5b9392505050565b6106b5806101c66000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c806396ea2f7a1161005057806396ea2f7a146100c7578063f2fde38b146100da578063ff12c354146100ed57600080fd5b8063181f5a771461007757806379ba5097146100955780638da5cb5b1461009f575b600080fd5b61007f610100565b60405161008c91906105b6565b60405180910390f35b61009d61011c565b005b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161008c565b61009d6100d5366004610644565b61021e565b61009d6100e8366004610644565b610299565b61009d6100fb366004610644565b6102aa565b6040518060600160405280602381526020016106866023913981565b60015473ffffffffffffffffffffffffffffffffffffffff1633146101a2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610296818273ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561026d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102919190610668565b6102f9565b50565b6102a161043e565b610296816104c1565b610296818273ffffffffffffffffffffffffffffffffffffffff16638fd6a6ac6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561026d573d6000803e3d6000fd5b73ffffffffffffffffffffffffffffffffffffffff81163314610368576040517fc454d18200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff808316600483015283166024820152604401610199565b6002546040517fd45ef15700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff848116600483015283811660248301529091169063d45ef15790604401600060405180830381600087803b1580156103dd57600080fd5b505af11580156103f1573d6000803e3d6000fd5b505060405173ffffffffffffffffffffffffffffffffffffffff8085169350851691507f09590fb70af4b833346363965e043a9339e8c7d378b8a2b903c75c277faec4f990600090a35050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146104bf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610199565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610540576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610199565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208083528351808285015260005b818110156105e3578581018301518582016040015282016105c7565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b73ffffffffffffffffffffffffffffffffffffffff8116811461029657600080fd5b60006020828403121561065657600080fd5b813561066181610622565b9392505050565b60006020828403121561067a57600080fd5b81516106618161062256fe52656769737472794d6f64756c654f776e6572437573746f6d20312e352e302d646576a164736f6c6343000813000a",
}

var RegistryModuleOwnerCustomABI = RegistryModuleOwnerCustomMetaData.ABI

var RegistryModuleOwnerCustomBin = RegistryModuleOwnerCustomMetaData.Bin

func DeployRegistryModuleOwnerCustom(auth *bind.TransactOpts, backend bind.ContractBackend, tokenAdminRegistry common.Address) (common.Address, *types.Transaction, *RegistryModuleOwnerCustom, error) {
	parsed, err := RegistryModuleOwnerCustomMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RegistryModuleOwnerCustomBin), backend, tokenAdminRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RegistryModuleOwnerCustom{address: address, abi: *parsed, RegistryModuleOwnerCustomCaller: RegistryModuleOwnerCustomCaller{contract: contract}, RegistryModuleOwnerCustomTransactor: RegistryModuleOwnerCustomTransactor{contract: contract}, RegistryModuleOwnerCustomFilterer: RegistryModuleOwnerCustomFilterer{contract: contract}}, nil
}

type RegistryModuleOwnerCustom struct {
	address common.Address
	abi     abi.ABI
	RegistryModuleOwnerCustomCaller
	RegistryModuleOwnerCustomTransactor
	RegistryModuleOwnerCustomFilterer
}

type RegistryModuleOwnerCustomCaller struct {
	contract *bind.BoundContract
}

type RegistryModuleOwnerCustomTransactor struct {
	contract *bind.BoundContract
}

type RegistryModuleOwnerCustomFilterer struct {
	contract *bind.BoundContract
}

type RegistryModuleOwnerCustomSession struct {
	Contract     *RegistryModuleOwnerCustom
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type RegistryModuleOwnerCustomCallerSession struct {
	Contract *RegistryModuleOwnerCustomCaller
	CallOpts bind.CallOpts
}

type RegistryModuleOwnerCustomTransactorSession struct {
	Contract     *RegistryModuleOwnerCustomTransactor
	TransactOpts bind.TransactOpts
}

type RegistryModuleOwnerCustomRaw struct {
	Contract *RegistryModuleOwnerCustom
}

type RegistryModuleOwnerCustomCallerRaw struct {
	Contract *RegistryModuleOwnerCustomCaller
}

type RegistryModuleOwnerCustomTransactorRaw struct {
	Contract *RegistryModuleOwnerCustomTransactor
}

func NewRegistryModuleOwnerCustom(address common.Address, backend bind.ContractBackend) (*RegistryModuleOwnerCustom, error) {
	abi, err := abi.JSON(strings.NewReader(RegistryModuleOwnerCustomABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindRegistryModuleOwnerCustom(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RegistryModuleOwnerCustom{address: address, abi: abi, RegistryModuleOwnerCustomCaller: RegistryModuleOwnerCustomCaller{contract: contract}, RegistryModuleOwnerCustomTransactor: RegistryModuleOwnerCustomTransactor{contract: contract}, RegistryModuleOwnerCustomFilterer: RegistryModuleOwnerCustomFilterer{contract: contract}}, nil
}

func NewRegistryModuleOwnerCustomCaller(address common.Address, caller bind.ContractCaller) (*RegistryModuleOwnerCustomCaller, error) {
	contract, err := bindRegistryModuleOwnerCustom(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryModuleOwnerCustomCaller{contract: contract}, nil
}

func NewRegistryModuleOwnerCustomTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryModuleOwnerCustomTransactor, error) {
	contract, err := bindRegistryModuleOwnerCustom(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryModuleOwnerCustomTransactor{contract: contract}, nil
}

func NewRegistryModuleOwnerCustomFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistryModuleOwnerCustomFilterer, error) {
	contract, err := bindRegistryModuleOwnerCustom(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistryModuleOwnerCustomFilterer{contract: contract}, nil
}

func bindRegistryModuleOwnerCustom(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RegistryModuleOwnerCustomMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RegistryModuleOwnerCustom.Contract.RegistryModuleOwnerCustomCaller.contract.Call(opts, result, method, params...)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegistryModuleOwnerCustomTransactor.contract.Transfer(opts)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegistryModuleOwnerCustomTransactor.contract.Transact(opts, method, params...)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RegistryModuleOwnerCustom.Contract.contract.Call(opts, result, method, params...)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.contract.Transfer(opts)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.contract.Transact(opts, method, params...)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RegistryModuleOwnerCustom.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomSession) Owner() (common.Address, error) {
	return _RegistryModuleOwnerCustom.Contract.Owner(&_RegistryModuleOwnerCustom.CallOpts)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomCallerSession) Owner() (common.Address, error) {
	return _RegistryModuleOwnerCustom.Contract.Owner(&_RegistryModuleOwnerCustom.CallOpts)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RegistryModuleOwnerCustom.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomSession) TypeAndVersion() (string, error) {
	return _RegistryModuleOwnerCustom.Contract.TypeAndVersion(&_RegistryModuleOwnerCustom.CallOpts)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomCallerSession) TypeAndVersion() (string, error) {
	return _RegistryModuleOwnerCustom.Contract.TypeAndVersion(&_RegistryModuleOwnerCustom.CallOpts)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.contract.Transact(opts, "acceptOwnership")
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomSession) AcceptOwnership() (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.AcceptOwnership(&_RegistryModuleOwnerCustom.TransactOpts)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.AcceptOwnership(&_RegistryModuleOwnerCustom.TransactOpts)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactor) RegisterAdminViaGetCCIPAdmin(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.contract.Transact(opts, "registerAdminViaGetCCIPAdmin", token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomSession) RegisterAdminViaGetCCIPAdmin(token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegisterAdminViaGetCCIPAdmin(&_RegistryModuleOwnerCustom.TransactOpts, token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactorSession) RegisterAdminViaGetCCIPAdmin(token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegisterAdminViaGetCCIPAdmin(&_RegistryModuleOwnerCustom.TransactOpts, token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactor) RegisterAdminViaOwner(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.contract.Transact(opts, "registerAdminViaOwner", token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomSession) RegisterAdminViaOwner(token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegisterAdminViaOwner(&_RegistryModuleOwnerCustom.TransactOpts, token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactorSession) RegisterAdminViaOwner(token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegisterAdminViaOwner(&_RegistryModuleOwnerCustom.TransactOpts, token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.contract.Transact(opts, "transferOwnership", to)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.TransferOwnership(&_RegistryModuleOwnerCustom.TransactOpts, to)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.TransferOwnership(&_RegistryModuleOwnerCustom.TransactOpts, to)
}

type RegistryModuleOwnerCustomAdministratorRegisteredIterator struct {
	Event *RegistryModuleOwnerCustomAdministratorRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RegistryModuleOwnerCustomAdministratorRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryModuleOwnerCustomAdministratorRegistered)
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
		it.Event = new(RegistryModuleOwnerCustomAdministratorRegistered)
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

func (it *RegistryModuleOwnerCustomAdministratorRegisteredIterator) Error() error {
	return it.fail
}

func (it *RegistryModuleOwnerCustomAdministratorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RegistryModuleOwnerCustomAdministratorRegistered struct {
	Token         common.Address
	Administrator common.Address
	Raw           types.Log
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomFilterer) FilterAdministratorRegistered(opts *bind.FilterOpts, token []common.Address, administrator []common.Address) (*RegistryModuleOwnerCustomAdministratorRegisteredIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var administratorRule []interface{}
	for _, administratorItem := range administrator {
		administratorRule = append(administratorRule, administratorItem)
	}

	logs, sub, err := _RegistryModuleOwnerCustom.contract.FilterLogs(opts, "AdministratorRegistered", tokenRule, administratorRule)
	if err != nil {
		return nil, err
	}
	return &RegistryModuleOwnerCustomAdministratorRegisteredIterator{contract: _RegistryModuleOwnerCustom.contract, event: "AdministratorRegistered", logs: logs, sub: sub}, nil
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomFilterer) WatchAdministratorRegistered(opts *bind.WatchOpts, sink chan<- *RegistryModuleOwnerCustomAdministratorRegistered, token []common.Address, administrator []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var administratorRule []interface{}
	for _, administratorItem := range administrator {
		administratorRule = append(administratorRule, administratorItem)
	}

	logs, sub, err := _RegistryModuleOwnerCustom.contract.WatchLogs(opts, "AdministratorRegistered", tokenRule, administratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RegistryModuleOwnerCustomAdministratorRegistered)
				if err := _RegistryModuleOwnerCustom.contract.UnpackLog(event, "AdministratorRegistered", log); err != nil {
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

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomFilterer) ParseAdministratorRegistered(log types.Log) (*RegistryModuleOwnerCustomAdministratorRegistered, error) {
	event := new(RegistryModuleOwnerCustomAdministratorRegistered)
	if err := _RegistryModuleOwnerCustom.contract.UnpackLog(event, "AdministratorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RegistryModuleOwnerCustomOwnershipTransferRequestedIterator struct {
	Event *RegistryModuleOwnerCustomOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RegistryModuleOwnerCustomOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryModuleOwnerCustomOwnershipTransferRequested)
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
		it.Event = new(RegistryModuleOwnerCustomOwnershipTransferRequested)
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

func (it *RegistryModuleOwnerCustomOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *RegistryModuleOwnerCustomOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RegistryModuleOwnerCustomOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RegistryModuleOwnerCustomOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RegistryModuleOwnerCustom.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RegistryModuleOwnerCustomOwnershipTransferRequestedIterator{contract: _RegistryModuleOwnerCustom.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RegistryModuleOwnerCustomOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RegistryModuleOwnerCustom.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RegistryModuleOwnerCustomOwnershipTransferRequested)
				if err := _RegistryModuleOwnerCustom.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomFilterer) ParseOwnershipTransferRequested(log types.Log) (*RegistryModuleOwnerCustomOwnershipTransferRequested, error) {
	event := new(RegistryModuleOwnerCustomOwnershipTransferRequested)
	if err := _RegistryModuleOwnerCustom.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RegistryModuleOwnerCustomOwnershipTransferredIterator struct {
	Event *RegistryModuleOwnerCustomOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RegistryModuleOwnerCustomOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryModuleOwnerCustomOwnershipTransferred)
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
		it.Event = new(RegistryModuleOwnerCustomOwnershipTransferred)
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

func (it *RegistryModuleOwnerCustomOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *RegistryModuleOwnerCustomOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RegistryModuleOwnerCustomOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RegistryModuleOwnerCustomOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RegistryModuleOwnerCustom.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RegistryModuleOwnerCustomOwnershipTransferredIterator{contract: _RegistryModuleOwnerCustom.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RegistryModuleOwnerCustomOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RegistryModuleOwnerCustom.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RegistryModuleOwnerCustomOwnershipTransferred)
				if err := _RegistryModuleOwnerCustom.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomFilterer) ParseOwnershipTransferred(log types.Log) (*RegistryModuleOwnerCustomOwnershipTransferred, error) {
	event := new(RegistryModuleOwnerCustomOwnershipTransferred)
	if err := _RegistryModuleOwnerCustom.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustom) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _RegistryModuleOwnerCustom.abi.Events["AdministratorRegistered"].ID:
		return _RegistryModuleOwnerCustom.ParseAdministratorRegistered(log)
	case _RegistryModuleOwnerCustom.abi.Events["OwnershipTransferRequested"].ID:
		return _RegistryModuleOwnerCustom.ParseOwnershipTransferRequested(log)
	case _RegistryModuleOwnerCustom.abi.Events["OwnershipTransferred"].ID:
		return _RegistryModuleOwnerCustom.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (RegistryModuleOwnerCustomAdministratorRegistered) Topic() common.Hash {
	return common.HexToHash("0x09590fb70af4b833346363965e043a9339e8c7d378b8a2b903c75c277faec4f9")
}

func (RegistryModuleOwnerCustomOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (RegistryModuleOwnerCustomOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustom) Address() common.Address {
	return _RegistryModuleOwnerCustom.address
}

type RegistryModuleOwnerCustomInterface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RegisterAdminViaGetCCIPAdmin(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error)

	RegisterAdminViaOwner(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAdministratorRegistered(opts *bind.FilterOpts, token []common.Address, administrator []common.Address) (*RegistryModuleOwnerCustomAdministratorRegisteredIterator, error)

	WatchAdministratorRegistered(opts *bind.WatchOpts, sink chan<- *RegistryModuleOwnerCustomAdministratorRegistered, token []common.Address, administrator []common.Address) (event.Subscription, error)

	ParseAdministratorRegistered(log types.Log) (*RegistryModuleOwnerCustomAdministratorRegistered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RegistryModuleOwnerCustomOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RegistryModuleOwnerCustomOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*RegistryModuleOwnerCustomOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RegistryModuleOwnerCustomOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RegistryModuleOwnerCustomOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*RegistryModuleOwnerCustomOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
