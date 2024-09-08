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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AddressZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"CanOnlySelfRegister\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"administrator\",\"type\":\"address\"}],\"name\":\"AdministratorRegistered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"registerAdminViaGetCCIPAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"registerAdminViaOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161047e38038061047e83398101604081905261002f91610067565b6001600160a01b03811661005657604051639fabe1c160e01b815260040160405180910390fd5b6001600160a01b0316608052610097565b60006020828403121561007957600080fd5b81516001600160a01b038116811461009057600080fd5b9392505050565b6080516103cc6100b2600039600061024a01526103cc6000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c8063181f5a771461004657806396ea2f7a14610098578063ff12c354146100ad575b600080fd5b6100826040518060400160405280601f81526020017f52656769737472794d6f64756c654f776e6572437573746f6d20312e352e300081525081565b60405161008f91906102ef565b60405180910390f35b6100ab6100a636600461037e565b6100c0565b005b6100ab6100bb36600461037e565b61013b565b610138818273ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561010f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061013391906103a2565b61018a565b50565b610138818273ffffffffffffffffffffffffffffffffffffffff16638fd6a6ac6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561010f573d6000803e3d6000fd5b73ffffffffffffffffffffffffffffffffffffffff811633146101fd576040517fc454d18200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff80831660048301528316602482015260440160405180910390fd5b6040517fe677ae3700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff838116600483015282811660248301527f0000000000000000000000000000000000000000000000000000000000000000169063e677ae3790604401600060405180830381600087803b15801561028e57600080fd5b505af11580156102a2573d6000803e3d6000fd5b505060405173ffffffffffffffffffffffffffffffffffffffff8085169350851691507f09590fb70af4b833346363965e043a9339e8c7d378b8a2b903c75c277faec4f990600090a35050565b60006020808352835180602085015260005b8181101561031d57858101830151858201604001528201610301565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b73ffffffffffffffffffffffffffffffffffffffff8116811461013857600080fd5b60006020828403121561039057600080fd5b813561039b8161035c565b9392505050565b6000602082840312156103b457600080fd5b815161039b8161035c56fea164736f6c6343000818000a",
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

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustom) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _RegistryModuleOwnerCustom.abi.Events["AdministratorRegistered"].ID:
		return _RegistryModuleOwnerCustom.ParseAdministratorRegistered(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (RegistryModuleOwnerCustomAdministratorRegistered) Topic() common.Hash {
	return common.HexToHash("0x09590fb70af4b833346363965e043a9339e8c7d378b8a2b903c75c277faec4f9")
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustom) Address() common.Address {
	return _RegistryModuleOwnerCustom.address
}

type RegistryModuleOwnerCustomInterface interface {
	TypeAndVersion(opts *bind.CallOpts) (string, error)

	RegisterAdminViaGetCCIPAdmin(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error)

	RegisterAdminViaOwner(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error)

	FilterAdministratorRegistered(opts *bind.FilterOpts, token []common.Address, administrator []common.Address) (*RegistryModuleOwnerCustomAdministratorRegisteredIterator, error)

	WatchAdministratorRegistered(opts *bind.WatchOpts, sink chan<- *RegistryModuleOwnerCustomAdministratorRegistered, token []common.Address, administrator []common.Address) (event.Subscription, error)

	ParseAdministratorRegistered(log types.Log) (*RegistryModuleOwnerCustomAdministratorRegistered, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
